package snmp

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	//"os/signal"
	"strings"
	//"syscall"
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
	mibdb        *mibs.MibDB // Global singleton instance here.
	dumpMibTable = flag.Bool("snmp_dump_mibs", false, "If true, dump the list of possible mibs on start.")
	flowOnly     = flag.Bool("snmp_flow_only", false, "If true, don't poll snmp devices.")
	jsonToYaml   = flag.String("snmp_json2yaml", "", "If set, convert the passed in json file to a yaml profile.")
	snmpWalk     = flag.String("snmp_do_walk", "", "If set, try to perform a snmp walk against the targeted device.")
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
		return snmp_util.DoWalk(*snmpWalk, conf, connectTimeout, retries, log)
	}

	// Load a mibdb if we have one.
	if conf.Global != nil {
		mdb, err := mibs.NewMibDB(conf.Global.MibDB, conf.Global.MibProfileDir, log)
		if err != nil {
			return fmt.Errorf("Cannot set up mibDB -- db: %s, profiles: %s -> %v", conf.Global.MibDB, conf.Global.MibProfileDir, err)
		}
		mibdb = mdb
	} else {
		log.Infof("Skipping configurable mibs")
	}

	// Now, launch a metadata and metrics server for each configured or discovered device.
	go wrapSnmpPolling(ctx, snmpFile, jchfChan, metrics, registry, apic, log)

	// Run a trap listener?
	if conf.Trap != nil {
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
	conf, err := parseConfig(snmpFile, log)
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

func wrapSnmpPolling(ctx context.Context, snmpFile string, jchfChan chan []*kt.JCHF, metrics *kt.SnmpMetricSet, registry go_metrics.Registry, apic *api.KentikApi, log logger.ContextL) {
	ctxSnmp, cancel := context.WithCancel(ctx)
	err := runSnmpPolling(ctxSnmp, snmpFile, jchfChan, metrics, registry, apic, log)
	if err != nil {
		log.Errorf("Error running snmp polling: %v", err)
	}

	// Now, wait for sigusr1 to re-do.
	c := make(chan os.Signal, 1)
	//signal.Notify(c, syscall.SIGUSR2)

	// Block here
	_ = <-c

	// If we got this signal, redo the snmp system.
	cancel()

	go wrapSnmpPolling(ctx, snmpFile, jchfChan, metrics, registry, apic, log)
}

func runSnmpPolling(ctx context.Context, snmpFile string, jchfChan chan []*kt.JCHF, metrics *kt.SnmpMetricSet, registry go_metrics.Registry, apic *api.KentikApi, log logger.ContextL) error {
	// Parse again to make sure nothing's changed.
	conf, connectTimeout, retries, err := initSnmp(ctx, snmpFile, log)
	if err != nil || conf == nil || conf.Global == nil {
		return err
	}

	log.Infof("Client SNMP: Setting up for %d devices", len(conf.Devices))
	for _, device := range conf.Devices {
		if device.Provider == "" {
			// Default provider to something we can work with.
			device.Provider = kt.ProviderRouter
		}
		if *flowOnly || device.FlowOnly {
			continue
		}

		log.Infof("Client SNMP: Running SNMP for %s on %s (type=%s)", device.DeviceName, device.DeviceIP, device.Provider)
		metrics.Mux.Lock()
		nm := kt.NewSnmpDeviceMetric(registry, device.DeviceName)
		metrics.Devices[device.DeviceName] = nm
		metrics.Mux.Unlock()
		cl := logger.NewSubContextL(logger.SContext{S: device.DeviceName}, log)
		var profile *mibs.Profile
		if mibdb != nil {
			profile = mibdb.FindProfile(device.OID)
			if profile != nil {
				log.Infof("Found profile for %s: %v", device.OID, profile.From)
				if *dumpMibTable {
					profile.DumpOids(cl)
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

	// We need two of these, to avoid concurrent access by the two pollers.
	// gosnmp isn't real clear on its approach to concurrency, but it seems
	// like maintaining separate GoSNMP structs for the two goroutines is safe.
	metadataServer, err := snmp_util.InitSNMP(device, connectTimeout, retries, "", log)
	if err != nil {
		log.Warnf("Init Issue starting SNMP interface component -- %v", err)
		return err
	}
	metricsServer, err := snmp_util.InitSNMP(device, connectTimeout, retries, "", log)
	if err != nil {
		log.Warnf("Init Issue starting SNMP interface component -- %v", err)
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
		_, err := metricPoller.Poll()
		if err != nil {
			log.Warnf("Init Issue polling SNMP counters: %v", err)
		}

		// Having done that, we'll launch additional, separate goroutines for
		// metadata and counter polling
		metadataPoller.StartLoop(ctx)
		metricPoller.StartLoop(ctx)
	}()

	return nil
}

func parseConfig(file string, log logger.ContextL) (*kt.SnmpConfig, error) {
	ms := kt.SnmpConfig{}
	by, err := ioutil.ReadFile(file)
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

	return &ms, nil
}

// Public wrapper for calling this other places.
func ParseConfig(file string, log logger.ContextL) (*kt.SnmpConfig, error) {
	return parseConfig(file, log)
}
