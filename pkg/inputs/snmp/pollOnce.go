package snmp

import (
	"context"
	"fmt"
	"os"
	"time"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/inputs/snmp/metadata"
	snmp_metrics "github.com/kentik/ktranslate/pkg/inputs/snmp/metrics"
	snmp_util "github.com/kentik/ktranslate/pkg/inputs/snmp/util"
	"github.com/kentik/ktranslate/pkg/kt"
)

func pollOnce(ctx context.Context, tdevice string, conf *kt.SnmpConfig, connectTimeout time.Duration, retries int, jchfChan chan []*kt.JCHF, metrics *kt.SnmpMetricSet, registry go_metrics.Registry, log logger.ContextL, logchan chan string) error {
	device := conf.Devices[tdevice]
	if device == nil {
		for _, dev := range conf.Devices {
			if dev.DeviceName == tdevice {
				device = dev
				break
			}
		}
	}

	if device == nil {
		return fmt.Errorf("The %s device was not found in the SNMP configuration file.", tdevice)
	}

	profile := mibdb.FindProfile(device.OID, device.Description, device.MibProfile)
	if profile == nil {
		return fmt.Errorf("No profile found for %s", tdevice)
	}

	// We need two of these, to avoid concurrent access by the two pollers.
	// gosnmp isn't real clear on its approach to concurrency, but it seems
	// like maintaining separate GoSNMP structs for the two goroutines is safe.
	metadataServer, err := snmp_util.InitSNMP(device, connectTimeout, retries, "", log)
	if err != nil {
		log.Warnf("There was an error when starting SNMP interface component -- %v.", err)
		return err
	}
	metricsServer, err := snmp_util.InitSNMP(device, connectTimeout, retries, "", log)
	if err != nil {
		log.Warnf("There was an error when starting SNMP interface component -- %v.", err)
		return err
	}

	nm := kt.NewSnmpDeviceMetric(registry, device.DeviceName)
	metadataPoller := metadata.NewPoller(metadataServer, conf.Global, device, jchfChan, nm, profile, log)
	metricPoller := snmp_metrics.NewPoller(metricsServer, conf.Global, device, jchfChan, nm, profile, log, logchan)

	metadataPoller.StartLoop(ctx)
	// Give a little time to get this done.
	time.Sleep(1 * time.Second)

	flows, err := metricPoller.Poll(ctx)
	if err != nil {
		return err
	}

	jchfChan <- flows

	// Give some time and then halt the process.
	go func() {
		time.Sleep(3 * time.Second)
		os.Exit(0)
	}()

	return nil
}
