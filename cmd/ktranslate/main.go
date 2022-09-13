package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/cmd/version"
	"github.com/kentik/ktranslate/pkg/cat"
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/imdario/mergo"
	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate/pkg/eggs/baseserver"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/eggs/properties"
	yaml "gopkg.in/yaml.v3"
)

var (
	listenIPPort   string
	mappingFile    string
	udrs           string
	geo            string
	asn            string
	dns            string
	threads        int
	threadsInput   int
	maxThreads     int
	format         string
	formatRollup   string
	formatMetric   string
	compression    string
	sinks          string
	maxFlows       int
	dumpRollups    int
	rollupAndAlpha bool
	sample         int
	sampleMin      int
	apiDevices     string
	snmpFile       string
	snmpDisco      bool
	kentikEmail    string
	kentikPlan     int
	apiRoot        string
	sslCertFile    string
	sslKeyFile     string
	tagMapType     string
	vpcSource      string
	flowSource     string
	teeLog         bool
	appMap         string
	syslog         string
	httpInput      bool
	enricher       string
)

func init() {
	// runtime options
	flag.StringVar(&listenIPPort, "listen", "127.0.0.1:8081", "IP:Port to listen on")
	flag.StringVar(&mappingFile, "mapping", "", "Mapping file to use for enums")
	flag.StringVar(&udrs, "udrs", "", "UDR mapping file")
	flag.StringVar(&geo, "geo", "", "Geo mapping file")
	flag.StringVar(&asn, "asn", "", "Asn mapping file")
	flag.StringVar(&dns, "dns", "", "Resolve IPs at this ip:port")
	flag.IntVar(&threads, "threads", 1, "Number of threads to run for processing")
	flag.IntVar(&threadsInput, "input_threads", 1, "Number of threads to run for input processing")
	flag.IntVar(&maxThreads, "max_threads", 1, "Dynamically grow threads up to this number")
	flag.StringVar(&format, "format", "flat_json", "Format to convert kflow to: (json|flat_json|avro|netflow|influx|carbon|prometheus|new_relic|new_relic_metric|splunk|elasticsearch|kflow)")
	flag.StringVar(&formatRollup, "format_rollup", "", "Format to convert rollups to: (json|avro|netflow|influx|prometheus|new_relic|new_relic_metric|splunk|elasticsearch|kflow)")
	flag.StringVar(&formatMetric, "format_metric", "", "Format to convert metrics to: (json|avro|netflow|influx|prometheus|new_relic|new_relic_metric|splunk|elasticsearch|kflow)")
	flag.StringVar(&compression, "compression", "none", "compression algo to use (none|gzip|snappy|deflate|null)")
	flag.StringVar(&sinks, "sinks", "stdout", "List of sinks to send data to. Options: (kafka|stdout|new_relic|kentik|net|http|splunk|prometheus|file|s3|gcloud)")
	flag.IntVar(&maxFlows, "max_flows_per_message", 10000, "Max number of flows to put in each emitted message")
	flag.IntVar(&dumpRollups, "rollup_interval", 0, "Export timer for rollups in seconds")
	flag.BoolVar(&rollupAndAlpha, "rollup_and_alpha", false, "Send both rollups and alpha inputs to sinks")
	flag.IntVar(&sample, "sample_rate", kt.LookupEnvInt("KENTIK_SAMPLE_RATE", 1), "Sampling rate to use. 1 -> 1:1 sampling, 2 -> 1:2 sampling and so on.")
	flag.IntVar(&sampleMin, "max_before_sample", 1, "Only sample when a set of inputs is at least this many")
	flag.StringVar(&apiDevices, "api_devices", "", "json file containing dumy devices to use for the stub Kentik API")
	flag.StringVar(&snmpFile, "snmp", "", "yaml file containing snmp config to use")
	flag.BoolVar(&snmpDisco, "snmp_discovery", false, "If true, try to discover snmp devices on this network as configured.")
	flag.StringVar(&kentikEmail, "kentik_email", "", "Kentik email to use for API calls")
	flag.IntVar(&kentikPlan, "kentik_plan", 0, "Kentik plan id to use for creating devices")
	flag.StringVar(&apiRoot, "api_root", "https://api.kentik.com", "API url prefix. If not set, defaults to https://api.kentik.com")
	flag.StringVar(&sslCertFile, "ssl_cert_file", "", "SSL Cert file to use for serving HTTPS traffic")
	flag.StringVar(&sslKeyFile, "ssl_key_file", "", "SSL Key file to use for serving HTTPS traffic")
	flag.StringVar(&tagMapType, "tag_map_type", "", "type of mapping to use for tag values. file|null")
	flag.StringVar(&vpcSource, "vpc", kt.LookupEnvString("KENTIK_VPC", ""), "Run VPC Flow Ingest")
	flag.StringVar(&flowSource, "nf.source", "", "Run NetFlow Ingest Directly. Valid values here are netflow5|netflow9|ipfix|sflow|nbar|asa|pan")
	flag.BoolVar(&teeLog, "tee_logs", false, "Tee log messages to sink")
	flag.StringVar(&appMap, "application_map", "", "File containing custom application mappings")
	flag.StringVar(&syslog, "syslog.source", "", "Run Syslog Server at this IP:Port or unix socket.")
	flag.BoolVar(&httpInput, "http.source", false, "Listen for content sent via http.")
	flag.StringVar(&enricher, "enricher", "", "Send data to this http url for enrichment.")
}

func main() {
	var (
		configFilePath = flag.String("config", "", "path to ktranslate config")
		generateConfig = flag.Bool("generate-config", false, "generate ktranslate config and exit")
	)

	// this is needed in order to catch the config options
	flag.Parse()

	// dump default config to stdout and exit
	if *generateConfig {
		if err := yaml.NewEncoder(os.Stdout).Encode(ktranslate.DefaultConfig()); err != nil {
			panic(err)
		}
		os.Exit(0)
	}

	cfg := ktranslate.DefaultConfig()

	// apply initial flags
	if err := applyFlags(cfg); err != nil {
		panic(err)
	}

	// if config specified, merge config
	if v := *configFilePath; v != "" {
		ktCfg, err := ktranslate.LoadConfig(v)
		if err != nil {
			panic(err)
		}

		// merge passed with default
		if err := mergo.Merge(ktCfg, cfg); err != nil {
			panic(err)
		}

		cfg = ktCfg
	}

	if err := applyMode(cfg, kt.LookupEnvString("KENTIK_MODE", flag.Arg(0))); err != nil {
		panic(err)
	}

	metricsChan := make(chan []*kt.JCHF, cat.CHAN_SLACK)
	bs := baseserver.BoilerplateWithPrefix("ktranslate", version.Version, "chf.kkc", properties.NewEnvPropertyBacking(), metricsChan, cfg.Server)
	bs.BaseServerConfiguration.SkipEnvDump = true // Turn off dumping the envs on panic

	// Set up NR logging if configured.
	logTee := make(chan string, cat.CHAN_SLACK)
	if cfg.EnableTeeLogs {
		bs.SetLogTee(logTee)
	}

	prefix := fmt.Sprintf("KTranslate")
	lc := logger.NewContextLFromUnderlying(logger.SContext{S: prefix}, bs.Logger)

	if cfg.ListenAddr == "" {
		bs.Fail("Invalid --listen value")
	}

	// and set this if overridden
	if dumpRollups > 0 {
		cat.RollupsSendDuration = time.Duration(dumpRollups) * time.Second
	}

	kc, err := cat.NewKTranslate(cfg, lc, go_metrics.DefaultRegistry, version.Version.Version, cfg.Sinks, bs.ServiceName, logTee, metricsChan)
	if err != nil {
		bs.Fail(fmt.Sprintf("Cannot start ktranslate: %v", err))
	}

	lc.Infof("Running -- Version %s; Build %s", version.Version.Version, version.Version.Date)
	lc.Infof("CLI: %v", os.Args)
	bs.Run(kc)
}

// apply config based on mode group
func applyMode(cfg *ktranslate.Config, mode string) error {
	setNr := func() { // Specific settings for NR
		cfg.Format = "new_relic"
		cfg.SampleMin = 100
		cfg.Compression = "gzip"
		cfg.Sinks = "new_relic"

		if cfg.SampleRate == 0 {
			cfg.SampleRate = 1000
		}
	}

	switch mode {
	case "":
		return nil // noop
	case "nr1.vpc.lambda":
		setNr() // Here, we only send the flow in as events to NR.
		cfg.AWSVPCInput.Enable = true
		cfg.AWSVPCInput.IsLambda = true
		cfg.DNS = "local"
	case "vpc":
		cfg.Rollup.Formats = append(cfg.Rollup.Formats, "s_sum,vpc.xmt.bytes,out_bytes,custom_str.source_vpc,custom_str.application_type,custom_str.source_account,custom_str.source_region,src_addr,custom_str.src_as_name,src_geo,l4_src_port,protocol")
		cfg.Rollup.Formats = append(cfg.Rollup.Formats, "s_sum,vpc.rcv.bytes,in_bytes,custom_str.dest_vpc,custom_str.application_type,custom_str.dest_account,custom_str.dest_region,dst_addr,custom_str.dst_as_name,dst_geo,l4_dst_port,protocol")
	case "nr1.vpc":
		cfg.DNS = "local"
		setNr()
	case "flow":
		cfg.Rollup.Formats = append(cfg.Rollup.Formats, "s_sum,bytes.xmt,in_bytes+out_bytes,device_name,src_addr,custom_str.src_as_name,src_geo,l4_src_port,protocol")
		cfg.Rollup.Formats = append(cfg.Rollup.Formats, "s_sum,bytes.rcv,in_bytes+out_bytes,device_name,dst_addr,custom_str.dst_as_name,dst_geo,l4_dst_port,protocol")
		cfg.Rollup.Formats = append(cfg.Rollup.Formats, "s_sum,pkts.xmt,in_pkts+out_pkts,device_name,src_addr,custom_str.src_as_name,src_geo,l4_src_port,protocol")
		cfg.Rollup.Formats = append(cfg.Rollup.Formats, "s_sum,pkts.rcv,in_pkts+out_pkts,device_name,dst_addr,custom_str.dst_as_name,dst_geo,l4_dst_port,protocol")
	case "nr1.flow":
		cfg.SNMPInput.FlowOnly = true
		setNr()
	case "nr1.discovery":
		cfg.EnableSNMPDiscovery = true
		setNr()
	case "nr1.syslog": // Tune for syslog. Don't want any sampling so can't use setNR directly.
		cfg.Compression = "gzip"
		cfg.Sinks = "new_relic"
		cfg.Format = "new_relic_metric"
		cfg.SNMPInput.FlowOnly = true // Don't do snmp polling.
		if cfg.SyslogInput.ListenAddr == "" {
			cfg.SyslogInput.ListenAddr = "0.0.0.0:5143"
		}
		cfg.SyslogInput.Enable = true
	case "nr1.snmp": // Tune for snmp sending.
		cfg.Compression = "gzip"
		cfg.Sinks = "new_relic"
		cfg.Format = "new_relic_metric"
		cfg.MaxFlowsPerMessage = 100
	default:
		return fmt.Errorf("Invalid mode " + mode + ". Options = nr1.vpc|nr1.flow|nr1.snmp|vpc|flow")
	}

	return nil
}

// TODO: this should be removed when flags are removed in favor of config
func applyFlags(cfg *ktranslate.Config) error {
	errCh := make(chan error, 1)
	doneCh := make(chan bool, 1)
	go func() {
		flag.VisitAll(func(f *flag.Flag) {
			val := f.Value.String()
			if val == "" {
				return
			}

			switch f.Name {
			case "listen":
				cfg.ListenAddr = val
			case "mapping":
				cfg.MappingFile = val
			case "udrs":
				cfg.UDRSFile = val
			case "geo":
				cfg.GeoFile = val
			case "asn":
				cfg.ASNFile = val
			case "dns":
				cfg.DNS = val
			case "threads":
				v, err := strconv.Atoi(val)
				if err != nil {
					errCh <- err
					return
				}
				if v == 0 {
					errCh <- fmt.Errorf("threads must be > 0")
					return
				}
				cfg.ProcessingThreads = v
			case "input_threads":
				v, err := strconv.Atoi(val)
				if err != nil {
					errCh <- err
					return
				}
				if v == 0 {
					errCh <- fmt.Errorf("input_threads must be > 0")
					return
				}
				cfg.InputThreads = v
			case "max_threads":
				v, err := strconv.Atoi(val)
				if err != nil {
					errCh <- err
					return
				}
				if v == 0 {
					errCh <- fmt.Errorf("max_threads must be > 0")
					return
				}
				cfg.MaxThreads = v
			case "format":
				cfg.Format = val
			case "format_rollup":
				cfg.FormatRollup = val
			case "format_metric":
				cfg.FormatMetric = val
			case "compression":
				cfg.Compression = val
			case "sinks":
				cfg.Sinks = val
			case "max_flows_per_message":
				v, err := strconv.Atoi(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.MaxFlowsPerMessage = v
			case "rollup_interval":
				v, err := strconv.Atoi(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.RollupInterval = v
			case "rollup_and_alpha":
				v, err := strconv.ParseBool(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.RollupAndAlpha = v
			case "sample_rate":
				v, err := strconv.Atoi(val)
				if err != nil {
					errCh <- err
					return
				}
				if v == 0 {
					errCh <- fmt.Errorf("sample_rate must be > 0")
					return
				}
				cfg.SampleRate = v
			case "max_before_sample":
				v, err := strconv.Atoi(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.SampleMin = v
			case "api_devices":
				cfg.API.DeviceFile = val
			case "snmp":
				cfg.SNMPInput.Enable = true
				cfg.SNMPInput.SNMPFile = val
			case "snmp_discovery":
				v, err := strconv.ParseBool(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.EnableSNMPDiscovery = v
			case "kentik_email":
				cfg.KentikEmail = val
			case "api_root":
				cfg.APIBaseURL = val
			case "kentik_plan":
				v, err := strconv.Atoi(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.KentikPlan = v
			case "ssl_cert_file":
				cfg.SSLCertFile = val
			case "ssl_key_file":
				cfg.SSLKeyFile = val
			case "tag_map_type":
				cfg.TagMapType = val
			case "vpc":
				switch strings.ToLower(val) {
				case "aws":
					cfg.AWSVPCInput.Enable = true
				case "gcp":
					cfg.GCPVPCInput.Enable = true
				default:
					errCh <- fmt.Errorf("unsupported vpc type: %s", val)
				}
			case "nf.source":
				cfg.FlowInput.Enable = true
				cfg.FlowInput.Protocol = val
			case "tee_logs":
				v, err := strconv.ParseBool(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.EnableTeeLogs = v
			case "application_map":
				cfg.ApplicationFile = val
			case "syslog.source":
				cfg.SyslogInput.Enable = true
				cfg.SyslogInput.ListenAddr = val
			case "http.source":
				v, err := strconv.ParseBool(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.EnableHTTPInput = v
			case "enricher":
				cfg.EnricherURL = val
			// pkg/maps/file
			case "tag_map":
				cfg.TagMapFile = val
			// pkg/formats/netflow
			case "netflow_version":
				cfg.NetflowFormat.Version = val
			// pkg/formats/prom
			case "info_collector":
				v, err := strconv.ParseBool(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.PrometheusFormat.EnableCollectorStats = v
			case "prom_seen":
				v, err := strconv.Atoi(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.PrometheusFormat.FlowsNeeded = v
			// pkg/formats/influxdb
			case "influxdb_measurement_prefix":
				cfg.InfluxDBFormat.MeasurementPrefix = val
			// pkg/sinks/prom
			case "prom_listen":
				cfg.PrometheusSink.ListenAddr = val
			// pkg/sinks/gcloud
			case "gcloud_bucket":
				cfg.GCloudSink.Bucket = val
			case "gcloud_prefix":
				cfg.GCloudSink.Prefix = val
			case "gcloud_content_type":
				cfg.GCloudSink.ContentType = val
			case "gcloud_flush_sec":
				v, err := strconv.Atoi(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.GCloudSink.FlushIntervalSeconds = v
			// pkg/sinks/s3
			case "s3_bucket":
				cfg.S3Sink.Bucket = val
			case "s3_prefix":
				cfg.S3Sink.Prefix = val
			case "s3_flush_sec":
				v, err := strconv.Atoi(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.S3Sink.FlushIntervalSeconds = v
			// pkg/sinks/net
			case "net_server":
				cfg.NetSink.Endpoint = val
			case "net_protocol":
				cfg.NetSink.Protocol = val
			// pkg/sinks/nr
			case "nr_account_id":
				cfg.NewRelicSink.Account = val
			case "nr_estimate_only":
				v, err := strconv.ParseBool(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.NewRelicSink.EstimateOnly = v
			case "nr_region":
				cfg.NewRelicSink.Region = val
			case "nr_check_json":
				v, err := strconv.ParseBool(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.NewRelicSink.ValidateJSON = v
			// pkg/sinks/file
			case "file_out":
				cfg.FileSink.Path = val
			case "file_on":
				v, err := strconv.ParseBool(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.FileSink.EnableImmediateWrite = v
			case "file_flush_sec":
				v, err := strconv.Atoi(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.FileSink.FlushIntervalSeconds = v
			// pkg/sinks/gcppubsub
			case "gcp_pubsub_project_id":
				cfg.GCloudPubSubSink.ProjectID = val
			case "gcp_pubsub_topic":
				cfg.GCloudPubSubSink.Topic = val
			// pkg/sinks/http
			case "http_url":
				cfg.HTTPSink.Target = val
			case "http_insecure":
				v, err := strconv.ParseBool(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.HTTPSink.InsecureSkipVerify = v
			case "http_timeout_sec":
				v, err := strconv.Atoi(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.HTTPSink.TimeoutInSeconds = v
			case "http_header":
				cfg.HTTPSink.Headers = strings.Split(val, ",")
			// pkg/sinks/kafka
			case "kafka_topic":
				cfg.KafkaSink.Topic = val
			case "bootstrap.servers":
				cfg.KafkaSink.BootstrapServers = val
			// pkg/sinks/kentik
			case "kentik_relay_url":
				cfg.KentikSink.RelayURL = val
			// pkg/rollup/rollup
			case "rollup_key_join":
				cfg.Rollup.JoinKey = val
			case "rollup_top_k":
				v, err := strconv.Atoi(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.Rollup.TopK = v
			case "rollups":
				cfg.Rollup.Formats = strings.Split(val, ",")
			// pkg/eggs/kmux
			case "dir":
				cfg.KMux.Dir = val
			// pkg/eggs/baseserver
			case "service_name":
				cfg.Server.ServiceName = val
			case "log_level":
				cfg.Server.LogLevel = val
			case "stdout":
				v, err := strconv.ParseBool(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.Server.LogToStdout = v
			case "metalisten":
				cfg.Server.MetaListenAddr = val
			case "metrics":
				cfg.Server.MetricsEndpoint = val
			case "olly_dataset":
				cfg.Server.OllyDataset = val
			case "olly_write_key":
				cfg.Server.OllyWriteKey = val
			// pkg/api
			case "api_device_file":
				cfg.API.DeviceFile = val
			// pkg/filter
			case "filters":
				cfg.Filters = strings.Split(val, ",")
			// pkg/inputs/syslog
			case "syslog.udp":
				v, err := strconv.ParseBool(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.SyslogInput.EnableUDP = v
			case "syslog.tcp":
				v, err := strconv.ParseBool(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.SyslogInput.EnableTCP = v
			case "syslog.unix":
				v, err := strconv.ParseBool(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.SyslogInput.EnableUnix = v
			case "syslog.format":
				cfg.SyslogInput.Format = val
			case "syslog.threads":
				v, err := strconv.Atoi(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.SyslogInput.Threads = v
			// pkg/inputs/snmp
			case "snmp_dump_mibs":
				v, err := strconv.ParseBool(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.SNMPInput.DumpMIBs = v
			case "flow_only":
				v, err := strconv.ParseBool(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.SNMPInput.FlowOnly = v
			case "snmp_json2yaml":
				cfg.SNMPInput.JSONToYAML = val
			case "snmp_do_walk":
				cfg.SNMPInput.WalkTarget = val
			case "snmp_walk_oid":
				cfg.SNMPInput.WalkOID = val
			case "snmp_walk_format":
				cfg.SNMPInput.WalkFormat = val
			case "snmp_out_file":
				cfg.SNMPInput.OutputFile = val
			case "snmp_poll_now":
				cfg.SNMPInput.PollNowTarget = val
			case "snmp_discovery_min":
				v, err := strconv.Atoi(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.SNMPInput.DiscoveryIntervalMinutes = v
			case "snmp_discovery_on_start":
				v, err := strconv.ParseBool(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.SNMPInput.DiscoveryOnStart = v
			case "snmp_validate":
				v, err := strconv.ParseBool(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.SNMPInput.ValidateMIBs = v
			// pkg/inputs/vpc/gcp
			case "gcp.project":
				cfg.GCPVPCInput.Enable = true
				cfg.GCPVPCInput.ProjectID = val
			case "gcp.sub":
				cfg.GCPVPCInput.Subject = val
			case "gcp.sample":
				v, err := strconv.ParseFloat(val, 64)
				if err != nil {
					errCh <- err
					return
				}
				cfg.GCPVPCInput.SampleRate = v
			// pkg/inputs/vpc/aws
			case "iam_role":
				cfg.AWSVPCInput.IAMRole = val
			case "sqs_name":
				cfg.AWSVPCInput.Enable = true
				cfg.AWSVPCInput.SQSName = val
			case "aws_regions":
				cfg.AWSVPCInput.Regions = strings.Split(val, ",")
			case "aws_lambda":
				v, err := strconv.ParseBool(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.AWSVPCInput.IsLambda = v
			case "aws_local_file":
				cfg.AWSVPCInput.LocalFile = val
			// pkg/inputs/flow
			case "nf.addr":
				cfg.FlowInput.ListenIP = val
			case "nf.port":
				v, err := strconv.Atoi(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.FlowInput.ListenPort = v
			case "nf.reuserport":
				v, err := strconv.ParseBool(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.FlowInput.EnableReusePort = v
			case "nf.workers":
				v, err := strconv.Atoi(val)
				if err != nil {
					errCh <- err
					return
				}
				cfg.FlowInput.Workers = v
			case "nf.message.fields":
				cfg.FlowInput.MessageFields = val
			case "nf.prom.listen":
				cfg.FlowInput.PrometheusListenAddr = val
			case "nf.mapping":
				cfg.FlowInput.MappingFile = val
			// configs
			case "config", "generate-config":
				// ignore
			default:
				// error here to detect flags that are not handled
				// by the config adapter
				errCh <- fmt.Errorf("unhandled flag %s", f.Name)
			}
		})

		doneCh <- true
	}()

	select {
	case <-doneCh:
	case err := <-errCh:
		return err
	}

	return nil
}
