package snmp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate/pkg/snmp/metadata"
	snmp_metrics "github.com/kentik/ktranslate/pkg/snmp/metrics"
	"github.com/kentik/ktranslate/pkg/snmp/mibs"
	"github.com/kentik/ktranslate/pkg/snmp/traps"
	snmp_util "github.com/kentik/ktranslate/pkg/snmp/util"
)

const (
	CHF_SNMP_TIMEOUT = "CHF_SNMP_TIMEOUT"
	CHF_SNMP_RETRY   = "CHF_SNMP_RETRY"
)

var (
	mibdb *mibs.MibDB // Global singleton instance here.
)

func StartSNMPPolls(snmpFile string, jchfChan chan []*kt.JCHF, metrics *kt.SnmpMetricSet, registry go_metrics.Registry, log logger.ContextL) error {

	// First, parse the config file and see what we're doing.
	log.Infof("Client SNMP: Running SNMP interface polling, loading config from %s", snmpFile)
	conf, err := parseConfig(snmpFile)
	if err != nil {
		return err
	}

	// Get these bits of info.
	connectTimeout := 5 * time.Second
	retries := 0

	// Update the timeout values, if needed.
	localTO := os.Getenv(CHF_SNMP_TIMEOUT)
	if lto, err := strconv.Atoi(localTO); err == nil {
		connectTimeout = time.Duration(lto) * time.Second
		log.Infof("Setting timeout to %v", connectTimeout)
	}

	localRt := os.Getenv(CHF_SNMP_RETRY)
	if lto, err := strconv.Atoi(localRt); err == nil {
		retries = lto
		log.Infof("Setting retry to %v", retries)
	}

	// Load a mibdb if we have one.
	if conf.Global != nil && (conf.Global.MibDB != "" && conf.Global.MibProfileDir != "") {
		mdb, err := mibs.NewMibDB(conf.Global.MibDB, conf.Global.MibProfileDir, log)
		if err != nil {
			return fmt.Errorf("Cannot set up mibDB -- db: %s, profiles: %s -> %v", conf.Global.MibDB, conf.Global.MibProfileDir, err)
		}
		mibdb = mdb
	} else {
		log.Infof("Skipping configurable mibs")
	}

	// Now, launch a metadata and metrics server for each configured or discovered device.
	for _, device := range conf.Devices {
		log.Infof("Client SNMP: Running SNMP for %s on %s", device.DeviceName, device.DeviceIP)
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
			}
		}
		err := launchSnmp(conf.Global, device, jchfChan, connectTimeout, retries, nm, profile, cl)
		if err != nil {
			return err
		}
	}

	// Run a trap listener?
	if conf.Trap != nil {
		err := launchSnmpTrap(conf.Trap, jchfChan, metrics, log)
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

func launchSnmpTrap(conf *kt.SnmpTrapConfig, jchfChan chan []*kt.JCHF, metrics *kt.SnmpMetricSet, log logger.ContextL) error {
	log.Infof("Client SNMP: Running SNMP Trap listener on %s", conf.Listen)
	tl, err := traps.NewSnmpTrapListener(conf, jchfChan, metrics, log)
	if err != nil {
		return err
	}

	go func() {
		tl.Listen()
	}()

	return nil
}

func launchSnmp(conf *kt.SnmpGlobalConfig, device *kt.SnmpDeviceConfig, jchfChan chan []*kt.JCHF, connectTimeout time.Duration, retries int, metrics *kt.SnmpDeviceMetric, profile *mibs.Profile, log logger.ContextL) error {

	// We need two of these, to avoid concurrent access by the two pollers.
	// gosnmp isn't real clear on its approach to concurrency, but it seems
	// like maintaining separate GoSNMP structs for the two goroutines is safe.
	metadataServer, err := snmp_util.InitSNMP(device, connectTimeout, retries, log)
	if err != nil {
		log.Warnf("Init Issue starting SNMP interface component -- %v", err)
		return err
	}
	metricsServer, err := snmp_util.InitSNMP(device, connectTimeout, retries, log)
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
		metadataPoller.StartLoop()
		metricPoller.StartLoop()
	}()

	return nil
}

func parseConfig(file string) (*kt.SnmpConfig, error) {
	ms := kt.SnmpConfig{}
	by, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(by, &ms)
	if err != nil {
		return nil, err
	}

	return &ms, nil
}
