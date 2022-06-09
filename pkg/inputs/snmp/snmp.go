package snmp

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"time"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate/pkg/api"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/inputs/snmp/metadata"
	snmp_metrics "github.com/kentik/ktranslate/pkg/inputs/snmp/metrics"
	"github.com/kentik/ktranslate/pkg/inputs/snmp/mibs"
	"github.com/kentik/ktranslate/pkg/inputs/snmp/traps"
	snmp_util "github.com/kentik/ktranslate/pkg/inputs/snmp/util"
	"github.com/kentik/ktranslate/pkg/kt"

	"gopkg.in/yaml.v2"
)

var (
	mibdb          *mibs.MibDB // Global singleton instance here.
	dumpMibTable   = flag.Bool("snmp_dump_mibs", false, "If true, dump the list of possible mibs on start.")
	flowOnly       = flag.Bool("flow_only", false, "If true, don't poll snmp devices.")
	jsonToYaml     = flag.String("snmp_json2yaml", "", "If set, convert the passed in json file to a yaml profile.")
	snmpWalk       = flag.String("snmp_do_walk", "", "If set, try to perform a snmp walk against the targeted device.")
	snmpWalkOid    = flag.String("snmp_walk_oid", ".1.3.6.1.2.1", "Walk this oid if -snmp_do_walk is set.")
	snmpWalkFormat = flag.String("snmp_walk_format", "", "use this format for walked values if -snmp_do_walk is set.")
	snmpOutFile    = flag.String("snmp_out_file", "", "If set, write updated snmp file here.")
	snmpPollNow    = flag.String("snmp_poll_now", "", "If set, run one snmp poll for the specified device and then exit.")
	snmpDiscoDur   = flag.Int("snmp_discovery_min", 0, "If set, run snmp discovery on this interval (in minutes).")
	snmpDiscoSt    = flag.Bool("snmp_discovery_on_start", false, "If set, run snmp discovery on application start.")
	validateMib    = flag.Bool("snmp_validate", false, "If true, validate mib profiles and exit.")
	ServiceName    = ""
)

func StartSNMPPolls(ctx context.Context, snmpFile string, jchfChan chan []*kt.JCHF, metrics *kt.SnmpMetricSet, registry go_metrics.Registry, apic *api.KentikApi, log logger.ContextL) error {
	// Do this once here just to see if we need to exit right away.
	conf, connectTimeout, retries, err := initSnmp(ctx, snmpFile, log)
	if err != nil || conf == nil || conf.Global == nil { // If no global, we're turning off all snmp polling.
		return err
	}

	if *jsonToYaml != "" { // If this flag is set, convert a passed in json mib file to a yaml profile and call it a day.
		return mibs.ConvertJson2Yaml(*jsonToYaml, log)
	}

	if *snmpWalk != "" { // If this flag is set, do just a snmp walk on the targeted device and exit.
		return snmp_util.DoWalk(*snmpWalk, *snmpWalkOid, *snmpWalkFormat, conf, connectTimeout, retries, log)
	}

	// Load a mibdb if we have one.
	if conf.Global != nil {
		mdb, err := mibs.NewMibDB(conf.Global.MibDB, conf.Global.MibProfileDir, *validateMib, log)
		if err != nil {
			time.Sleep(2 * time.Second) // Give logs time to get sent back.
			return fmt.Errorf("There was an error when setting up the %s mibDB database and the %s profiles: %v.", conf.Global.MibDB, conf.Global.MibProfileDir, err)
		}
		mibdb = mdb
		if *validateMib {
			// We just want to validate that this was ok so time to exit now.
			os.Exit(0)
		}
	} else {
		log.Infof("Skipping configurable mibs")
	}

	// If we just want to poll one device and exit, do this here.
	if *snmpPollNow != "" {
		return pollOnce(ctx, *snmpPollNow, conf, connectTimeout, retries, jchfChan, metrics, registry, log)
	}

	// Now, launch a metadata and metrics server for each configured or discovered device.
	go wrapSnmpPolling(ctx, snmpFile, jchfChan, metrics, registry, apic, log, 0)

	// Run a trap listener?
	if conf.Trap != nil && !*flowOnly {
		err := launchSnmpTrap(conf, jchfChan, metrics, log)
		if err != nil {
			return err
		}
	}

	return nil
}

func Close() {
	if mibdb != nil {
		mibdb.Close()
	}
}

func initSnmp(ctx context.Context, snmpFile string, log logger.ContextL) (*kt.SnmpConfig, time.Duration, int, error) {
	// First, parse the config file and see what we're doing.
	log.Infof("Client SNMP: Running SNMP interface polling, loading config from %s", snmpFile)
	conf, err := parseConfig(ctx, snmpFile, log)
	if err != nil {
		return nil, 0, 0, err
	}

	// If there's no global section, just turn ourselves off here.
	if conf.Global == nil {
		return nil, 0, 0, nil
	}

	// Get these bits of info.
	connectTimeout := 5 * time.Second
	retries := 0

	// Update the timeout values, if needed.
	if conf.Global.TimeoutMS > 0 {
		connectTimeout = time.Duration(conf.Global.TimeoutMS) * time.Millisecond
	}
	if conf.Global.Retries > 0 {
		retries = conf.Global.Retries
	}
	log.Infof("Setting timeout to %v", connectTimeout)
	log.Infof("Setting retry to %v", retries)

	return conf, connectTimeout, retries, nil
}

func wrapSnmpPolling(ctx context.Context, snmpFile string, jchfChan chan []*kt.JCHF, metrics *kt.SnmpMetricSet, registry go_metrics.Registry, apic *api.KentikApi, log logger.ContextL, restartCount int) {
	ctxSnmp, cancel := context.WithCancel(ctx)
	err := runSnmpPolling(ctxSnmp, snmpFile, jchfChan, metrics, registry, apic, log, restartCount)
	if err != nil {
		log.Errorf("There was an error when polling for SNMP devices: %v.", err)
	}

	// Now, wait for sigusr2 to re-do or if there's a discovery with new devices.
	c := make(chan os.Signal, 1)
	signal.Notify(c, kt.SIGUSR2)
	if *snmpDiscoDur > 0 { // If we are re-running snmp discovery every interval, start the ticker here.
		go RunDiscoOnTimer(ctxSnmp, c, snmpFile, log, *snmpDiscoDur, *snmpDiscoSt)
	}

	// Block here
	_ = <-c

	// If we got this signal, redo the snmp system.
	cancel()

	go wrapSnmpPolling(ctx, snmpFile, jchfChan, metrics, registry, apic, log, restartCount+1) // Track how many times through here we've been.
}

func runSnmpPolling(ctx context.Context, snmpFile string, jchfChan chan []*kt.JCHF, metrics *kt.SnmpMetricSet, registry go_metrics.Registry, apic *api.KentikApi, log logger.ContextL, restartCount int) error {
	// Parse again to make sure nothing's changed.
	conf, connectTimeout, retries, err := initSnmp(ctx, snmpFile, log)
	if err != nil || conf == nil || conf.Global == nil {
		return err
	}

	log.Infof("Client SNMP: Setting up for %d devices", len(conf.Devices))
	for _, device := range conf.Devices {
		if device.Provider == "" {
			// Default provider to something we can work with.
			device.Provider = kt.ProviderDefault
		}
		if *flowOnly || device.FlowOnly {
			continue
		}

		log.Infof("Client SNMP: Running SNMP for %s on %s (type=%s)", device.DeviceName, device.DeviceIP, device.Provider)

		// Check for duplicate device names here.
		metrics.Mux.Lock()
		var nm *kt.SnmpDeviceMetric
		if _, ok := metrics.Devices[device.DeviceName]; ok {
			nm = metrics.Devices[device.DeviceName]
			if restartCount == 0 {
				log.Errorf("Duplicate device name detected (%s). Is this a misconfiguration?", device.DeviceName)
			}
		} else {
			nm = kt.NewSnmpDeviceMetric(registry, device.DeviceName)
			metrics.Devices[device.DeviceName] = nm
		}
		metrics.Mux.Unlock()

		cl := logger.NewSubContextL(logger.SContext{S: device.DeviceName}, log)
		var profile *mibs.Profile
		if mibdb != nil {
			profile = mibdb.FindProfile(device.OID, device.Description, device.MibProfile)
			if profile != nil {
				cl.Infof("Found profile for %s: %v %s", device.OID, profile.From, device.MibProfile)
				if *dumpMibTable {
					profile.DumpOids(cl)
				}

				if profile.NoUseBulkWalkAll {
					device.NoUseBulkWalkAll = true
					cl.Infof("Turning off BulkWalkAll for device via profile.")
				}

				// Use the profile's provider if it is set.
				if profile.Provider != "" {
					device.Provider = profile.Provider
					cl.Infof("Setting profile of %s for device %s.", device.Provider, device.DeviceName)
				}
			}
		}

		// Create this device in Kentik if the option is set.
		err := apic.EnsureDevice(ctx, device)
		if err != nil {
			return err
		}

		err = launchSnmp(ctx, conf.Global, device, jchfChan, connectTimeout, retries, nm, profile, cl)
		if err != nil {
			return err
		}
	}

	return nil
}

func launchSnmpTrap(conf *kt.SnmpConfig, jchfChan chan []*kt.JCHF, metrics *kt.SnmpMetricSet, log logger.ContextL) error {
	log.Infof("Client SNMP: Running SNMP Trap listener on %s", conf.Trap.Listen)
	tl, err := traps.NewSnmpTrapListener(conf, jchfChan, metrics, mibdb, log)
	if err != nil {
		return err
	}

	go func() {
		tl.Listen()
	}()

	return nil
}

func launchSnmp(ctx context.Context, conf *kt.SnmpGlobalConfig, device *kt.SnmpDeviceConfig, jchfChan chan []*kt.JCHF, connectTimeout time.Duration, retries int, metrics *kt.SnmpDeviceMetric, profile *mibs.Profile, log logger.ContextL) error {
	// Sometimes this device is pinging only. In this case, start the ping loop and return.
	if device.PingOnly {
		return launchPingOnly(ctx, conf, device, jchfChan, connectTimeout, retries, metrics, profile, log)
	} else if conf.RunPing || device.RunPing {
		if err := launchPingOnly(ctx, conf, device, jchfChan, connectTimeout, retries, metrics, profile, log); err != nil {
			return err
		}
	}

	// Sometimes a device is only going to be running its extention.
	if device.Ext != nil && device.Ext.ExtOnly {
		return launchExtOnly(ctx, conf, device, jchfChan, connectTimeout, retries, metrics, profile, log)
	}

	// We need two of these, to avoid concurrent access by the two pollers.
	// gosnmp isn't real clear on its approach to concurrency, but it seems
	// like maintaining separate GoSNMP structs for the two goroutines is safe.
	metadataServer, err := snmp_util.InitSNMP(device, connectTimeout, retries, "", log)
	if err != nil {
		log.Warnf("There was an error when starting SNMP interface component -- %v.", err)
		metrics.Fail.Update(kt.SNMP_BAD_INIT_METADATA)
		return err
	}
	metricsServer, err := snmp_util.InitSNMP(device, connectTimeout, retries, "", log)
	if err != nil {
		log.Warnf("There was an error when starting SNMP interface component -- %v.", err)
		metrics.Fail.Update(kt.SNMP_BAD_INIT_METRICS)
		return err
	}

	metadataPoller := metadata.NewPoller(metadataServer, conf, device, jchfChan, metrics, profile, log)
	metricPoller := snmp_metrics.NewPoller(metricsServer, conf, device, jchfChan, metrics, profile, log)

	// We've now done everything we can do synchronously -- return to the client initialization
	// code, and do everything else in the background
	go func() {
		// Do a first counter poll to seed the interface tracker with a rough sort of the
		// interfaces by volume.  We need this before the first metadata poll runs, so that
		// we can discard low-volume interfaces
		// NOTE: we run this poll even if the customer has opted for Minimum SNMP polling
		// ALSO: we throw away the CHFs returned here, since they can't possibly have non-zero
		// delta values yet.
		_, err := metricPoller.Poll(ctx)
		if err != nil {
			metrics.Fail.Update(kt.SNMP_BAD_FIRST_METRICS_POLL)
			log.Warnf("There was an error when polling the SNMP counters: %v.", err)
		}

		// Having done that, we'll launch additional, separate goroutines for
		// metadata and counter polling
		metadataPoller.StartLoop(ctx)
		metricPoller.StartLoop(ctx)
	}()

	return nil
}

func parseConfig(ctx context.Context, file string, log logger.ContextL) (*kt.SnmpConfig, error) {
	ms := kt.SnmpConfig{}
	by, err := snmp_util.LoadFile(ctx, file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(by, &ms)
	if err != nil {
		return nil, err
	}

	// Expand out any seconds which require it.
	if ms.Disco != nil && len(ms.Disco.Cidrs) > 0 && strings.HasPrefix(ms.Disco.Cidrs[0], "@") {
		cidrList := []string{}
		byc, err := ioutil.ReadFile(ms.Disco.Cidrs[0][1:])
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(byc, &cidrList)
		if err != nil {
			return nil, err
		}
		ms.Disco.CidrOrig = ms.Disco.Cidrs[0][1:]
		ms.Disco.Cidrs = cidrList
	}

	if ms.Disco != nil && len(ms.Disco.IgnoreList) > 0 && strings.HasPrefix(ms.Disco.IgnoreList[0], "@") {
		ignoreList := []string{}
		byc, err := ioutil.ReadFile(ms.Disco.IgnoreList[0][1:])
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(byc, &ignoreList)
		if err != nil {
			return nil, err
		}
		ms.Disco.IgnoreOrig = ms.Disco.IgnoreList[0][1:]
		ms.Disco.IgnoreList = ignoreList
	}

	fullDevices := map[string]*kt.SnmpDeviceConfig{}
	for name, devFile := range ms.Devices {
		if strings.HasPrefix(name, "file_") && strings.HasPrefix(devFile.DeviceName, "@") {
			deviceSet := map[string]*kt.SnmpDeviceConfig{}
			byc, err := ioutil.ReadFile(devFile.DeviceName[1:])
			if err != nil {
				return nil, err
			}
			err = yaml.Unmarshal(byc, &deviceSet)
			if err != nil {
				return nil, err
			}
			ms.DeviceOrig = devFile.DeviceName[1:]
			log.Infof("Loading %d devices from %s", len(deviceSet), devFile.DeviceName[1:])
			for k, v := range deviceSet {
				fullDevices[k] = v
			}
		}
	}
	if len(fullDevices) > 0 {
		ms.Devices = fullDevices
	}

	// If there's a global v3, map it in for any disco, trap and device settings.
	if ms.Global != nil && ms.Global.GlobalV3 != nil {
		if ms.Disco != nil && ms.Disco.DefaultV3.InheritGlobal() {
			ms.Disco.DefaultV3 = ms.Global.GlobalV3
		}
		if ms.Trap != nil && ms.Trap.V3.InheritGlobal() {
			ms.Trap.V3 = ms.Global.GlobalV3
		}
		for _, device := range ms.Devices {
			if device.V3.InheritGlobal() {
				device.V3 = ms.Global.GlobalV3
			}
		}
	}

	// If there's a global user tags and match, add them in here.
	if ms.Global != nil {
		for k, v := range ms.Global.UserTags {
			for _, device := range ms.Devices {
				if device.UserTags == nil {
					device.UserTags = map[string]string{}
				}
				if _, ok := device.UserTags[k]; !ok {
					device.UserTags[k] = v
				}
			}
		}
		for k, v := range ms.Global.MatchAttr {
			for _, device := range ms.Devices {
				if device.MatchAttr == nil {
					device.MatchAttr = map[string]string{}
				}
				if _, ok := device.MatchAttr[k]; !ok {
					device.MatchAttr[k] = v
				}
			}
		}
	}

	// Correctly format all the user tags needed here:
	for _, device := range ms.Devices {
		device.InitUserTags(ServiceName)
	}

	return &ms, nil
}

/**
Handle the case where we're only doing a ping loop of a device.
*/
func launchPingOnly(ctx context.Context, conf *kt.SnmpGlobalConfig, device *kt.SnmpDeviceConfig, jchfChan chan []*kt.JCHF, connectTimeout time.Duration, retries int, metrics *kt.SnmpDeviceMetric, profile *mibs.Profile, log logger.ContextL) error {
	metricPoller := snmp_metrics.NewPollerForPing(conf, device, jchfChan, metrics, profile, log)

	// We've now done everything we can do synchronously -- return to the client initialization
	// code, and do everything else in the background
	go func() {
		metricPoller.StartPingOnlyLoop(ctx)
	}()

	return nil
}

/**
Handle the case where we're only doing a extention loop of a device.
*/
func launchExtOnly(ctx context.Context, conf *kt.SnmpGlobalConfig, device *kt.SnmpDeviceConfig, jchfChan chan []*kt.JCHF, connectTimeout time.Duration, retries int, metrics *kt.SnmpDeviceMetric, profile *mibs.Profile, log logger.ContextL) error {
	metricPoller := snmp_metrics.NewPollerForExtention(conf, device, jchfChan, metrics, profile, log)

	// We've now done everything we can do synchronously -- return to the client initialization
	// code, and do everything else in the background
	go func() {
		metricPoller.StartExtensionOnlyLoop(ctx)
	}()

	return nil
}

// Public wrapper for calling this other places.
func ParseConfig(file string, log logger.ContextL) (*kt.SnmpConfig, error) {
	return parseConfig(context.Background(), file, log)
}
