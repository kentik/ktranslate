//go:build windows
// +build windows

package snmp

import (
	"context"
	"os"
	"os/signal"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/api"
	"github.com/kentik/ktranslate/pkg/config"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
)

func wrapSnmpPolling(ctx context.Context, snmpFile string, jchfChan chan []*kt.JCHF, metrics *kt.SnmpMetricSet, registry go_metrics.Registry, apic *api.KentikApi, log logger.ContextL, restartCount int, cfg *ktranslate.SNMPInputConfig, confMgr config.ConfigManager) {
	ctxSnmp, cancel := context.WithCancel(ctx)
	err := runSnmpPolling(ctxSnmp, snmpFile, jchfChan, metrics, registry, apic, log, restartCount, cfg)
	if err != nil {
		log.Errorf("There was an error when polling for SNMP devices: %v.", err)
	}

	// We only want to run a disco on start when restartCount is 0. Otherwise you end up doing 2 discos if a new device is found on start.
	runOnStart := cfg.DiscoveryOnStart
	if restartCount > 0 {
		runOnStart = false
	}

	// Now, wait for sigusr2 to re-do or if there's a discovery with new devices.
	c := make(chan os.Signal, 1)
	signal.Notify(c, kt.SIGUSR2)
	if v := cfg.DiscoveryIntervalMinutes; v > 0 || runOnStart { // If we are re-running snmp discovery every interval AND/OR running on start, start the ticker here.
		go RunDiscoOnTimer(ctxSnmp, c, log, v, runOnStart, cfg, apic, confMgr)
	}

	// Block here
	_ = <-c

	// If we got this signal, redo the snmp system.
	cancel()

	go wrapSnmpPolling(ctx, snmpFile, jchfChan, metrics, registry, apic, log, restartCount+1, cfg, confMgr) // Track how many times through here we've been.
}
