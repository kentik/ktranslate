package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/kentik/ktranslate/cmd/version"
	"github.com/kentik/ktranslate/pkg/cat"
	"github.com/kentik/ktranslate/pkg/cat/auth"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate/pkg/eggs/baseserver"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/eggs/properties"
	"github.com/kentik/ktranslate/pkg/inputs/flow"
	"github.com/kentik/ktranslate/pkg/inputs/vpc"
	"github.com/kentik/ktranslate/pkg/maps"
)

func main() {
	var (
		// runtime options
		listenIPPort   = flag.String("listen", "127.0.0.1:8081", "IP:Port to listen on")
		mappingFile    = flag.String("mapping", "", "Mapping file to use for enums")
		udrs           = flag.String("udrs", "", "UDR mapping file")
		geo            = flag.String("geo", "", "Geo mapping file")
		asn            = flag.String("asn", "", "Asn mapping file")
		dns            = flag.String("dns", "", "Resolve IPs at this ip:port")
		threads        = flag.Int("threads", 0, "Number of threads to run for processing")
		threadsInput   = flag.Int("input_threads", 0, "Number of threads to run for input processing")
		maxThreads     = flag.Int("max_threads", 0, "Dynamically grow threads up to this number")
		format         = flag.String("format", "flat_json", "Format to convert kflow to: (json|flat_json|avro|netflow|influx|prometheus|new_relic|new_relic_metric|splunk|elasticsearch)")
		formatRollup   = flag.String("format_rollup", "", "Format to convert rollups to: (json|avro|netflow|influx|prometheus|new_relic|new_relic_metric|splunk|elasticsearch)")
		compression    = flag.String("compression", "none", "compression algo to use (none|gzip|snappy|deflate|null)")
		sinks          = flag.String("sinks", "stdout", "List of sinks to send data to. Options: (kafka|stdout|new_relic|kentik|net|http|splunk|prometheus|file|s3|gcloud)")
		maxFlows       = flag.Int("max_flows_per_message", 10000, "Max number of flows to put in each emitted message")
		dumpRollups    = flag.Int("rollup_interval", 0, "Export timer for rollups in seconds")
		rollupAndAlpha = flag.Bool("rollup_and_alpha", false, "Send both rollups and alpha inputs to sinks")
		sample         = flag.Int("sample_rate", kt.LookupEnvInt("KENTIK_SAMPLE_RATE", 0), "Sampling rate to use. 1 -> 1:1 sampling, 2 -> 1:2 sampling and so on.")
		sampleMin      = flag.Int("max_before_sample", 1, "Only sample when a set of inputs is at least this many")
		apiDevices     = flag.String("api_devices", "", "json file containing dumy devices to use for the stub Kentik API")
		snmpFile       = flag.String("snmp", "", "yaml file containing snmp config to use")
		snmpDisco      = flag.Bool("snmp_discovery", false, "If true, try to discover snmp devices on this network as configured.")
		kentikEmail    = flag.String("kentik_email", "", "Kentik email to use for API calls")
		apiRoot        = flag.String("api_root", "https://api.kentik.com", "API url prefix. If not set, defaults to https://api.kentik.com")
		kentikPlan     = flag.Int("kentik_plan", 0, "Kentik plan id to use for creating devices")
		sslCertFile    = flag.String("ssl_cert_file", "", "SSL Cert file to use for serving HTTPS traffic")
		sslKeyFile     = flag.String("ssl_key_file", "", "SSL Key file to use for serving HTTPS traffic")
		tagMapType     = flag.String("tag_map_type", "", "type of mapping to use for tag values. file|null")
		vpcSource      = flag.String("vpc", kt.LookupEnvString("KENTIK_VPC", ""), "Run VPC Flow Ingest")
		flowSource     = flag.String("nf.source", "", "Run NetFlow Ingest Directly. Valid values here are netflow5|netflow9|ipfix|sflow")
		teeLog         = flag.Bool("tee_logs", false, "Tee log messages to sink")
		appMap         = flag.String("application_map", "", "File containing custom application mappings")
		syslog         = flag.String("syslog.source", "", "Run Syslog Server at this IP:Port or unix socket.")
	)

	metricsChan := make(chan []*kt.JCHF, cat.CHAN_SLACK)
	bs := baseserver.BoilerplateWithPrefix("ktranslate", version.Version, "chf.kkc", properties.NewEnvPropertyBacking(), metricsChan)
	bs.BaseServerConfiguration.SkipEnvDump = true // Turn off dumping the envs on panic

	// Set up NR logging if configured.
	logTee := make(chan string, cat.CHAN_SLACK)
	if *teeLog {
		bs.SetLogTee(logTee)
	}

	// If we're running in a given mode, set the flags accordingly.
	setMode(bs, kt.LookupEnvString("KENTIK_MODE", flag.Arg(0)), *sample, *syslog)

	if *listenIPPort == "" {
		bs.Fail("Invalid --listen value")
	}

	prefix := fmt.Sprintf("KTranslate")
	lc := logger.NewContextLFromUnderlying(logger.SContext{S: prefix}, bs.Logger)

	conf := cat.Config{
		Listen:            *listenIPPort,
		SslCertFile:       *sslCertFile,
		SslKeyFile:        *sslKeyFile,
		MappingFile:       *mappingFile,
		Format:            formats.Format(*format),
		FormatRollup:      formats.Format(*formatRollup),
		Threads:           *threads,
		ThreadsInput:      *threadsInput,
		MaxThreads:        *maxThreads,
		Compression:       kt.Compression(*compression),
		MaxFlowPerMessage: *maxFlows,
		RollupAndAlpha:    *rollupAndAlpha,
		UDRFile:           *udrs,
		GeoMapping:        *geo,
		AsnMapping:        *asn,
		DnsResolver:       *dns,
		SampleRate:        uint32(*sample),
		MaxBeforeSample:   *sampleMin,
		SNMPFile:          *snmpFile,
		SNMPDisco:         *snmpDisco,
		TagMapType:        maps.Mapper(*tagMapType),
		VpcSource:         vpc.CloudSource(*vpcSource),
		FlowSource:        flow.FlowSource(*flowSource),
		SyslogSource:      *syslog,
		AppMap:            *appMap,
		Kentik: &kt.KentikConfig{
			ApiEmail: *kentikEmail,
			ApiToken: os.Getenv(kt.KentikAPIToken),
			ApiRoot:  *apiRoot,
			ApiPlan:  *kentikPlan,
		},
		LogTee:      logTee,
		MetricsChan: metricsChan,
	}

	if *apiDevices != "" {
		conf.Auth = &auth.AuthConfig{
			DevicesFile: *apiDevices,
		}
	}

	// Default these to 1.
	if conf.Threads <= 0 {
		conf.Threads = 1
	}
	if conf.ThreadsInput <= 0 {
		conf.ThreadsInput = 1
	}
	if conf.MaxThreads <= 0 {
		conf.MaxThreads = 1
	}
	if conf.SampleRate == 0 {
		conf.SampleRate = 1
	}

	// and set this if overridden
	if *dumpRollups > 0 {
		cat.RollupsSendDuration = time.Duration(*dumpRollups) * time.Second
	}

	kc, err := cat.NewKTranslate(&conf, lc, go_metrics.DefaultRegistry, version.Version.Version, *sinks, bs.ServiceName)
	if err != nil {
		bs.Fail(fmt.Sprintf("Cannot start ktranslate: %v", err))
	}

	lc.Infof("Running -- Version %s; Build %s", version.Version.Version, version.Version.Date)
	lc.Infof("CLI: %v", os.Args)
	bs.Run(kc)
}

func setMode(bs *baseserver.BaseServer, mode string, sample int, syslog string) {
	setNr := func() { // Specific settings for NR
		flag.Set("format", "new_relic")
		flag.Set("max_before_sample", "100")
		flag.Set("compression", "gzip")
		flag.Set("sinks", "new_relic")

		if sample == 0 {
			flag.Set("sample_rate", "1000")
		}
	}

	switch mode {
	case "":
		return // noop
	case "nr1.vpc.lambda":
		setNr() // Here, we only send the flow in as events to NR.
		flag.Set("vpc", "aws")
		flag.Set("aws_lambda", "true")
	case "vpc":
		flag.Set("rollups", "s_sum,vpc.xmt.bytes,out_bytes,custom_str.source_vpc,custom_str.application_type,custom_str.source_account,custom_str.source_region,src_addr,custom_str.src_as_name,src_geo,l4_src_port,protocol")
		flag.Set("rollups", "s_sum,vpc.rcv.bytes,in_bytes,custom_str.dest_vpc,custom_str.application_type,custom_str.dest_account,custom_str.dest_region,dst_addr,custom_str.dst_as_name,dst_geo,l4_dst_port,protocol")
	case "nr1.vpc":
		setNr()
	case "flow":
		flag.Set("rollups", "s_sum,bytes.xmt,in_bytes+out_bytes,device_name,src_addr,custom_str.src_as_name,src_geo,l4_src_port,protocol")
		flag.Set("rollups", "s_sum,bytes.rcv,in_bytes+out_bytes,device_name,dst_addr,custom_str.dst_as_name,dst_geo,l4_dst_port,protocol")
		flag.Set("rollups", "s_sum,pkts.xmt,in_pkts+out_pkts,device_name,src_addr,custom_str.src_as_name,src_geo,l4_src_port,protocol")
		flag.Set("rollups", "s_sum,pkts.rcv,in_pkts+out_pkts,device_name,dst_addr,custom_str.dst_as_name,dst_geo,l4_dst_port,protocol")
	case "nr1.flow":
		flag.Set("flow_only", "true")
		setNr()
	case "nr1.discovery":
		flag.Set("snmp_discovery", "true")
		setNr()
	case "nr1.syslog": // Tune for syslog.
		flag.Set("compression", "gzip")
		flag.Set("sinks", "new_relic")
		flag.Set("format", "new_relic_metric")
		if syslog == "" {
			flag.Set("syslog.source", "0.0.0.0:5143")
		}
	case "nr1.snmp": // Tune for snmp sending.
		flag.Set("compression", "gzip")
		flag.Set("sinks", "new_relic")
		flag.Set("format", "new_relic_metric")
		flag.Set("max_flows_per_message", "100")
	default:
		bs.Fail("Invalid mode " + mode + ". Options = nr1.vpc|nr1.flow|nr1.snmp|vpc|flow")
	}
}
