package snmp

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate/pkg/snmp/metadata"
	snmp_metrics "github.com/kentik/ktranslate/pkg/snmp/metrics"
	"github.com/kentik/ktranslate/pkg/snmp/traps"
	snmp_util "github.com/kentik/ktranslate/pkg/snmp/util"
)

const (
	CHF_SNMP_TIMEOUT = "CHF_SNMP_TIMEOUT"
	CHF_SNMP_RETRY   = "CHF_SNMP_RETRY"
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

	// Now, launch a metadata and metrics server for each configured or discovered device.
	for _, device := range conf.Devices {
		log.Infof("Client SNMP: Running SNMP for %s on %s", device.DeviceName, device.DeviceIP)
		metrics.Mux.Lock()
		nm := kt.NewSnmpDeviceMetric(registry, device.DeviceName)
		metrics.Devices[device.DeviceName] = nm
		metrics.Mux.Unlock()
		err := launchSnmp(device, jchfChan, connectTimeout, retries, nm, log)
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

func launchSnmp(device *kt.SnmpDeviceConfig, jchfChan chan []*kt.JCHF, connectTimeout time.Duration, retries int, metrics *kt.SnmpDeviceMetric, log logger.ContextL) error {

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

	metadataPoller := metadata.NewPoller(metadataServer, device, jchfChan, metrics, log)
	metricPoller := snmp_metrics.NewPoller(metricsServer, device, jchfChan, metrics, log)

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
