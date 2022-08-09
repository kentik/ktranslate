package baseserver

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/timing"

	"github.com/google/uuid"
	libhoney "github.com/honeycombio/libhoney-go"
	"github.com/judwhite/go-svc"
	"github.com/kentik/ktranslate/pkg/eggs/olly"

	"github.com/kentik/ktranslate/pkg/eggs/concurrent"
	"github.com/kentik/ktranslate/pkg/eggs/version"

	"github.com/kentik/ktranslate/pkg/eggs/features"
	"github.com/kentik/ktranslate/pkg/eggs/properties"

	"github.com/kentik/ktranslate/pkg/util/cmetrics"
	"github.com/kentik/ktranslate/pkg/util/logger"
)

const (
	FAILURE_CODE                 = -10
	ENV_CH_NUM_CPU               = "CH_NUM_CPU"
	readinessWaitGroupContextKey = "_baseserver_ready_wg"
	subContextNameContextKey     = "_baseserver_subctx"
)

var (
	serviceName  string
	logLevel     string
	logToStdout  bool
	metricsDest  string
	metaListen   string
	ollyDataset  string
	ollyWriteKey string
)

func init() {
	flag.StringVar(&serviceName, "service_name", "", "Service identifier")
	flag.StringVar(&logLevel, "log_level", "info", "Logging Level")
	flag.BoolVar(&logToStdout, "stdout", false, "Log to stdout")
	flag.StringVar(&metricsDest, "metrics", "none", "Metrics Configuration. none|syslog|stderr|graphite:127.0.0.1:2003")
	flag.StringVar(&metaListen, "metalisten", "localhost:0", "HTTP interface and port to bind on")
	flag.StringVar(&ollyDataset, "olly_dataset", "", "Olly dataset name")
	flag.StringVar(&ollyWriteKey, "olly_write_key", "", "Olly dataset name")
}

type BaseServerConfiguration struct {
	// base service properties
	ServiceName string
	VersionInfo version.VersionInfo

	// operational
	ShutdownSettleTime time.Duration

	// logging
	LogToStdout bool
	LogLevel    string
	LogPrefix   string

	// metrics
	MetricsPrefix      string
	MetricsDestination string

	// olly
	OllyWriteKey string
	OllyDataset  string

	// meta server properties
	MetaListen string

	// healthchecks
	HealthCheckStartupDelay time.Duration
	HealthCheckPeriod       time.Duration
	HealthCheckTimeout      time.Duration

	// props
	PropsRefreshPeriod time.Duration

	// Skip env dump
	SkipEnvDump bool
}

var BaseServerConfigurationDefaults = BaseServerConfiguration{
	LogToStdout:             true,
	LogLevel:                "info",
	MetricsDestination:      "none",
	MetaListen:              "",
	ShutdownSettleTime:      1 * time.Second,
	HealthCheckStartupDelay: 5 * time.Second,
	HealthCheckPeriod:       30 * time.Second,
	HealthCheckTimeout:      5 * time.Second,
	PropsRefreshPeriod:      5 * time.Minute,
	OllyDataset:             "", // olly is disabled by default
	OllyWriteKey:            "",
}

type BaseServer struct {
	*BaseServerConfiguration
	hce             *HealthCheckExecutor
	metaServer      *MetaServer
	Logger          *logger.Logger
	initialLogLevel logger.Level
	ctx             context.Context
	cancel          context.CancelFunc
	waitGroup       sync.WaitGroup
	propertyService properties.PropertyService
	featureService  features.FeatureService
	ollyBuilder     *olly.Builder
	config          *ktranslate.ServerConfig
}

// Perform baseserver initialization steps -- hopefully 9 out of 10 services can just call this and Run()
func Boilerplate(serviceName string, versionInfo version.VersionInfo, defaultPropertyBacking properties.PropertyBacking, mextra interface{}, cfg *ktranslate.ServerConfig) *BaseServer {
	bs := NewBaseServer(serviceName, versionInfo, "chf", defaultPropertyBacking, cfg)
	bs.Init(mextra)
	setGlobalBaseServer(bs)
	return bs
}

// For when you need to set metrics prefix.
func BoilerplateWithPrefix(serviceName string, versionInfo version.VersionInfo, metricsPrefix string, defaultPropertyBacking properties.PropertyBacking, mextra interface{}, cfg *ktranslate.ServerConfig) *BaseServer {
	bs := NewBaseServer(serviceName, versionInfo, metricsPrefix, defaultPropertyBacking, cfg)
	bs.Init(mextra)
	setGlobalBaseServer(bs)
	return bs
}

func NewBaseServer(serviceName string, version version.VersionInfo, metricsPrefix string, defaultPropertyBacking properties.PropertyBacking, cfg *ktranslate.ServerConfig) *BaseServer {
	conf := BaseServerConfigurationDefaults
	conf.ServiceName = cfg.ServiceName
	conf.VersionInfo = version
	conf.MetricsPrefix = metricsPrefix
	conf.LogPrefix = serviceName + " "

	conf.LogLevel = cfg.LogLevel
	conf.LogToStdout = cfg.LogToStdout
	conf.MetricsDestination = cfg.MetricsEndpoint
	conf.MetaListen = cfg.MetaListenAddr
	conf.OllyDataset = cfg.OllyDataset
	conf.OllyWriteKey = cfg.OllyWriteKey

	props := properties.NewPropertyService(
		properties.NewFileSystemPropertyBacking("/props"), // highest prio: dynamic FS props
		properties.NewEnvPropertyBacking(),                // env variables can override static defaults
		defaultPropertyBacking,                            // lowest prio: static default values
	)

	bs := &BaseServer{
		BaseServerConfiguration: &conf,
		propertyService:         props,
		featureService:          features.NewFeatureService(props),
	}
	bs.waitGroup.Add(1)
	return bs
}

func (bs *BaseServer) GetPropertyService() properties.PropertyService {
	return bs.propertyService
}

func (bs *BaseServer) GetFeatureService() features.FeatureService {
	return bs.featureService
}

func (bs *BaseServer) GetHealthCheckHandler() func(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 3; i++ {
		if bs.metaServer != nil {
			return bs.metaServer.endpoint_healthcheck
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

// Perform some early initialization steps -- things it makes sense to do before callers start building/initializing
// anything from the actual service. Most things should probably be started from Run().  Called by Boilerplate().
func (bs *BaseServer) Init(mextra interface{}) {
	bs.InitLogger(bs.LogToStdout, bs.LogLevel)
	bs.Logger.Infof(bs.LogPrefix, "version %s starting", bs.VersionInfo.Version)
	bs.InitMaxProcs()
	bs.InitOlly()
	bs.InitMetrics(mextra)
}

func (bs *BaseServer) Fail(msg string) {
	if bs.Logger != nil {
		bs.Logger.Panic(bs.LogPrefix, msg)
	}
	fmt.Printf("%s\n", msg)

	if !bs.SkipEnvDump {
		fmt.Printf("Environment:\n")
		env, redacted := getRedactedEnvironment()
		for k, v := range env {
			fmt.Printf("%s=%s\n", k, v)
		}
		for _, k := range redacted {
			fmt.Printf("%s=%s\n", k, " # redacted")
		}
	}

	os.Exit(FAILURE_CODE)
}

func (bs *BaseServer) WaitUntilReady(timeout time.Duration) {
	concurrent.WgWaitTimeout(&bs.waitGroup, timeout)
}

// Finish initializing and run until signaled otherwise. Spawns sub routines.
func (bs *BaseServer) Run(service Service) {
	bs.ctx, bs.cancel = context.WithCancel(context.WithValue(context.WithValue(context.Background(), readinessWaitGroupContextKey, &bs.waitGroup), subContextNameContextKey, "BaseServer.run"))

	bs.spawnPropsRefresh(bs.ctx)
	bs.spawnHealthCheck(bs.readyAwareSubContext(bs.ctx, "health check"), service)
	bs.spawnLegacyHealthCheck(bs.readyAwareSubContext(bs.ctx, "legacy health check"), service)
	bs.spawnMetaServer(bs.readyAwareSubContext(bs.ctx, "metaserver"), service)

	// If windows, turn over to windows process here
	if runtime.GOOS == "windows" {
		bs.Logger.Infof(bs.LogPrefix, "Running in Windows mode")
		defer bs.cancel()
		if err := svc.Run(service); err != nil {
			bs.Fail(fmt.Sprintf("service Run() error: %v", err))
		}
		return
	}

	// run the actual service
	go func(ctx context.Context) {
		setReady(ctx)
		if err := service.Run(ctx); err != nil {
			if err == context.Canceled {
				bs.Logger.Infof(bs.LogPrefix, "service context cancelled")
			} else {
				bs.Fail(fmt.Sprintf("service Run() error: %v", err))
			}
		}
		bs.cancel()
	}(bs.readyAwareSubContext(bs.ctx, "service run goroutine"))

	s := make(chan os.Signal, 2)
	signal.Notify(s, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

	setReady(bs.ctx) // goes with waitGroup.Add(1) in NewBaseServer
	olly.QuickC(bs, olly.Op("baseserver.start"))

	for {
		select {
		case <-bs.ctx.Done():
			olly.QuickC(bs, olly.Op("baseserver.stop"))
			bs.finish()
			return
		case sig := <-s:
			switch sig {
			case syscall.SIGQUIT:
				bs.Shutdown("SIGQUIT")
			case syscall.SIGINT:
				bs.Shutdown("SIGINT")
			case syscall.SIGTERM:
				bs.Shutdown("SIGTERM")
			}
		}
	}
}

func (bs *BaseServer) finish() {
	bs.Logger.Infof(bs.LogPrefix, "service.Close() called, now waiting for things to settle")

	t := timing.StartChrono()
	bs.FinishLogger()
	bs.FinishOlly()

	time.Sleep(bs.ShutdownSettleTime - t.Duration()) // Give everything enough time to settle.
	bs.Logger.Infof(bs.LogPrefix, "draining logger and exiting main thread")
}

func (bs *BaseServer) Shutdown(reason string) {
	bs.Logger.Info(bs.LogPrefix, "Shutdown('%s')", reason)
	bs.cancel()
}

// Tell base logger to also log messages along this tee
func (bs *BaseServer) SetLogTee(logChan chan string) {
	logger.SetTee(logChan)
}

// Initialize logging.
func (bs *BaseServer) InitLogger(stdout bool, loglevel string) {

	bs.LogPrefix = bs.ServiceName + " "

	if stdout {
		logger.SetStdOut()
	}

	progSvcName := path.Base(os.Args[0])
	if progSvcName != bs.ServiceName {
		progSvcName = fmt.Sprintf("%s/%s", progSvcName, bs.ServiceName)
	}

	pid := os.Getpid()
	if pid > 10 {
		// Note about the above comparison: if our pid is super low, we're probably running inside docker and/or in a
		// context where pid is not likely to be very important, and we omit it.
		progSvcName = fmt.Sprintf("%s(%d)", progSvcName, pid)
	}

	if err := logger.SetLogName(fmt.Sprintf("%s ", progSvcName)); err != nil {
		bs.Fail("Cannot set log name for program")
	}
	ll, ok := logger.CfgLevels[strings.ToLower(loglevel)]
	if !ok {
		bs.Fail("Unsupported log level: " + loglevel)
	}
	bs.initialLogLevel = ll
	if bs.Logger = logger.New(ll); bs.Logger == nil {
		bs.Fail("Cannot start logger")
	}
}

func (bs *BaseServer) FinishLogger() {
	logger.Drain()
}

// Set the number of cpus this process can use.
func (bs *BaseServer) InitMaxProcs() {
	if nc, err := strconv.Atoi(os.Getenv(ENV_CH_NUM_CPU)); err == nil {
		runtime.GOMAXPROCS(nc)
		bs.Logger.Info(bs.LogPrefix, "Setting GOMAXPROCS to %d", nc)
	}
}

// Initialize metrics.
func (bs *BaseServer) InitMetrics(extra interface{}) {
	tags := []string{
		"ver=" + bs.VersionInfo.Version,
		"svc=" + bs.ServiceName,
	}
	cmetrics.SetConf(bs.MetricsDestination, bs.Logger, bs.LogPrefix, bs.MetricsPrefix, nil, tags, nil, nil, extra)
}

// Initialize olly observability.
func (bs *BaseServer) InitOlly() {
	if bs.OllyDataset == "" || bs.OllyWriteKey == "" {
		bs.Logger.Infof(bs.LogPrefix, "olly: disabled")
		bs.ollyBuilder = olly.NewBuilder()
		return
	}
	bs.Logger.Infof(bs.LogPrefix, "olly: enabled")

	hostname, _ := os.Hostname() // nolint:errcheck

	olly.Init(bs.ServiceName, bs.VersionInfo.Version, libhoney.Config{
		WriteKey: bs.OllyWriteKey,
		Dataset:  bs.OllyDataset,
	}, "svc_process_uuid", uuid.New().String(), "node", hostname)
	bs.ollyBuilder = olly.NewBuilder()
}

func (bs *BaseServer) FinishOlly() {
	olly.Close()
}

// Initialize our legacy health check.
func (bs *BaseServer) spawnHealthCheck(ctx context.Context, service Service) {
	bs.hce = NewHealthCheckExecutor(service, bs.HealthCheckStartupDelay, bs.HealthCheckPeriod, bs.HealthCheckTimeout)
	go bs.hce.Run(ctx)
}

// Start legacy healthcheck if needed. Called as part of Init()
func (bs *BaseServer) spawnLegacyHealthCheck(ctx context.Context, service Service) {
	setReady(ctx)
}

func (bs *BaseServer) spawnMetaServer(ctx context.Context, service Service) {
	if bs.MetaListen == "" { // This is turned off for now.
		return
	}

	if bs.hce == nil {
		bs.Fail("initMetaServer: hce cannot be nil")
	}

	go func() {
		bs.metaServer = NewMetaServer(bs.MetaListen, bs.ServiceName, bs.VersionInfo, service, bs.Logger, bs.initialLogLevel, bs.hce)
		if err := bs.metaServer.Run(ctx); err != nil {
			bs.Fail(fmt.Sprintf("Error running meta server: %+v", err))
		}
	}()
}

func (bs *BaseServer) spawnPropsRefresh(ctx context.Context) {
	go func() {
		for {
			select {
			case <-time.After(bs.PropsRefreshPeriod):
				bs.propertyService.Refresh()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (bs *BaseServer) OllyBuilder() *olly.Builder {
	return bs.ollyBuilder
}

func (bs *BaseServer) readyAwareSubContext(ctx context.Context, name string) context.Context {
	val := ctx.Value(readinessWaitGroupContextKey)

	if val == nil {
		bs.Fail("Context is missing a value for readinessWaitGroupContextKey")
	}

	wg := val.(*sync.WaitGroup)
	wg.Add(1)

	// fmt.Printf("+ subcontext(%s) wg(%+v)\n", name, wg)
	return context.WithValue(context.WithValue(ctx, readinessWaitGroupContextKey, wg), subContextNameContextKey, name)
}

func setReady(ctx context.Context) {
	/* subCtxNameStr := "UNKNOWN"
	if val := ctx.Value(subContextNameContextKey); val != nil {
		subCtxNameStr = val.(string)
	}
	*/
	if val := ctx.Value(readinessWaitGroupContextKey); val != nil {
		wg := val.(*sync.WaitGroup)
		// fmt.Printf("- subcontext(%s) wg(%+v)\n", subCtxNameStr, wg)
		wg.Done()
	} /* else {
		 fmt.Printf("- subcontext(%s) wg(%+v)\n", subCtxNameStr, nil)
	} */
}
