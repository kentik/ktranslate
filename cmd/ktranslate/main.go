package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/kentik/ktranslate/cmd/version"
	"github.com/kentik/ktranslate/pkg/cat"
	"github.com/kentik/ktranslate/pkg/formats"
	ktb "github.com/kentik/ktranslate/pkg/kt"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate/pkg/eggs/baseserver"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/eggs/properties"
)

const ()

func main() {
	var (
		// runtime options
		listenIPPort   = flag.String("listen", "127.0.0.1:8081", "IP:Port to listen on")
		mappingFile    = flag.String("mapping", "config.json", "Mapping file to use for enums")
		region         = flag.String("region", "", "Region mapping file")
		city           = flag.String("city", "", "City mapping file")
		interfaces     = flag.String("interfaces", "", "Interface mapping file")
		udrs           = flag.String("udrs", "", "UDR mapping file")
		geo            = flag.String("geo", "", "Geo mapping file")
		asn4           = flag.String("asn4", "", "Asn ipv6 mapping file")
		asn6           = flag.String("asn6", "", "Asn ipv6 mapping file")
		dns            = flag.String("dns", "", "Resolve IPs at this ip:port")
		threads        = flag.Int("threads", 0, "Number of threads to run for processing")
		format         = flag.String("format", "json", "Format to convert kflow to: (json|avro|netflow|influx|prometheus|new_relic)")
		compression    = flag.String("compression", "none", "compression algo to use (none|gzip|snappy|deflate|null)")
		sinks          = flag.String("sinks", "stdout", "List of sinks to send data to. Options: (kafka|stdout|new_relic|kentik|net|http|splunk|prometheus|file|s3|gcloud)")
		maxFlows       = flag.Int("max_flows_per_message", 10000, "Max number of flows to put in each emitted message")
		dumpRollups    = flag.Int("rollup_interval", 0, "Export timer for rollups in seconds")
		rollupAndAlpha = flag.Bool("rollup_and_alpha", false, "Send both rollups and alpha inputs to sinks")
		sample         = flag.Int("sample_rate", 1, "Sampling rate to use. 1 -> 1:1 sampling, 2 -> 1:2 sampling and so on.")
		apiDevices     = flag.String("api_devices", "", "json file containing dumy devices to use for the stub Kentik API")
		snmpFile       = flag.String("snmp", "", "yaml file containing snmp config to use")
		snmpDisco      = flag.Bool("snmp_discovery", false, "If true, try to discover snmp devices on this network as configured.")
		subtype        = flag.String("subtype", "", "Load mappings for this device subtype")
	)

	bs := baseserver.BoilerplateWithPrefix("ktranslate", version.Version, "chf.kkc", properties.NewEnvPropertyBacking())
	bs.BaseServerConfiguration.SkipEnvDump = true // Turn off dumping the envs on panic

	if *listenIPPort == "" {
		bs.Fail("Invalid --listen value")
	}

	prefix := fmt.Sprintf("KTranslate")
	lc := logger.NewContextLFromUnderlying(logger.SContext{S: prefix}, bs.Logger)

	conf := cat.Config{
		Listen:            *listenIPPort,
		MappingFile:       *mappingFile,
		Code2Region:       *region,
		Code2City:         *city,
		Format:            formats.Format(*format),
		Threads:           *threads,
		Compression:       ktb.Compression(*compression),
		MaxFlowPerMessage: *maxFlows,
		RollupAndAlpha:    *rollupAndAlpha,
		DeviceFile:        *interfaces,
		UDRFile:           *udrs,
		GeoMapping:        *geo,
		Asn4:              *asn4,
		Asn6:              *asn6,
		DnsResolver:       *dns,
		SampleRate:        uint32(*sample),
		SNMPFile:          *snmpFile,
		SNMPDisco:         *snmpDisco,
		Subtype:           *subtype,
	}

	if *apiDevices != "" {
		conf.Auth = &cat.AuthConfig{
			DevicesFile: *apiDevices,
		}
	}

	// Default this to 1.
	if conf.Threads <= 0 {
		conf.Threads = 1
	}

	// and set this if overridden
	if *dumpRollups > 0 {
		cat.RollupsSendDuration = time.Duration(*dumpRollups) * time.Second
	}

	kc, err := cat.NewKTranslate(&conf, lc, go_metrics.DefaultRegistry, version.Version.Version, bs.Logger, *sinks)
	if err != nil {
		bs.Fail(fmt.Sprintf("Cannot start ktranslate: %v", err))
	} else {
		lc.Infof("Version %s; Build %s", version.Version.Version, version.Version.Date)
	}

	lc.Infof("Running")
	bs.Run(kc)
}
