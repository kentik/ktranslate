package syslog

import (
	"context"
	"flag"
	"fmt"
	"net"
	"strings"
	"time"

	go_metrics "github.com/kentik/go-metrics"
	"gopkg.in/mcuadros/go-syslog.v2"
	sfmt "gopkg.in/mcuadros/go-syslog.v2/format"

	"github.com/kentik/ktranslate/pkg/api"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/util/resolv"
)

type KentikSyslog struct {
	logger.ContextL
	server   *syslog.Server
	handler  *syslog.ChannelHandler
	channel  syslog.LogPartsChannel
	logchan  chan string
	metrics  SyslogMetric
	apic     *api.KentikApi
	devices  map[string]*kt.Device
	resolver *resolv.Resolver
}

type SyslogMetric struct {
	Messages go_metrics.Meter
	Errors   go_metrics.Meter
	Queue    go_metrics.Gauge
}

var (
	doUDP   = flag.Bool("syslog.udp", true, "Listen on UDP for syslog messages.")
	doTCP   = flag.Bool("syslog.tcp", true, "Listen on TCP for syslog messages.")
	doUnix  = flag.Bool("syslog.unix", false, "Listen on a Unix socket for syslog messages.")
	format  = flag.String("syslog.format", "Automatic", "Format to parse syslog messages with. Options are: Automatic|RFC3164|RFC5424|RFC6587.")
	threads = flag.Int("syslog.threads", 1, "Number of threads to use to process messages.")
)

const (
	CHAN_SLACK           = 10000
	DeviceUpdateDuration = 1 * time.Hour
	InstNameSyslog       = "ktranslate-syslog"
)

func NewSyslogSource(ctx context.Context, host string, log logger.Underlying, logchan chan string, registry go_metrics.Registry, apic *api.KentikApi, resolver *resolv.Resolver) (*KentikSyslog, error) {
	ks := KentikSyslog{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "Syslog"}, log),
		logchan:  logchan,
		metrics: SyslogMetric{
			Messages: go_metrics.GetOrRegisterMeter(fmt.Sprintf("syslog_messages^force=true"), registry),
			Errors:   go_metrics.GetOrRegisterMeter(fmt.Sprintf("syslog_errors^force=true"), registry),
			Queue:    go_metrics.GetOrRegisterGauge(fmt.Sprintf("syslog_queue^force=true"), registry),
		},
		apic:     apic,
		devices:  apic.GetDevicesAsMap(0),
		resolver: resolver,
	}

	if logchan == nil {
		return nil, fmt.Errorf("Log Sink is not set.")
	}

	channel := make(syslog.LogPartsChannel, CHAN_SLACK)
	handler := syslog.NewChannelHandler(channel)
	server := syslog.NewServer()
	switch *format {
	case "RFC3164":
		server.SetFormat(syslog.RFC3164)
	case "RFC5424":
		server.SetFormat(syslog.RFC5424)
	case "RFC6587":
		server.SetFormat(syslog.RFC6587)
	case "Automatic":
		server.SetFormat(syslog.Automatic)
	default:
		return nil, fmt.Errorf("Invalid syslog format (%s). Options are RFC3164|RFC5424|RFC6587", *format)
	}

	server.SetHandler(handler)
	if *doUDP {
		err := server.ListenUDP(host)
		if err != nil {
			return nil, fmt.Errorf("Cannot listen for syslog udp: %v", err)
		}
		ks.Infof("Listening for UDP on %s", host)
	}
	if *doTCP {
		err := server.ListenTCP(host)
		if err != nil {
			return nil, fmt.Errorf("Cannot listen for syslog with tcp: %v", err)
		}
		ks.Infof("Listening for TCP on %s", host)
	}
	if *doUnix {
		err := server.ListenUnixgram(host)
		if err != nil {
			return nil, fmt.Errorf("Cannot listen for syslog with unixgram: %v", err)
		}
		ks.Infof("Listening for Unixgram on %s", host)
	}

	err := server.Boot()
	if err != nil {
		return nil, fmt.Errorf("Cannot boot syslog server: %v", err)
	}
	ks.server = server
	ks.channel = channel
	ks.handler = handler

	go ks.run(ctx, host)

	return &ks, nil
}

func (ks *KentikSyslog) Close() {}

func (ks *KentikSyslog) HttpInfo() map[string]float64 {
	msgs := map[string]float64{
		"messages": ks.metrics.Messages.Rate1(),
		"errors":   ks.metrics.Errors.Rate1(),
		"queue":    float64(ks.metrics.Queue.Value()),
	}
	return msgs
}

func (ks *KentikSyslog) process(ctx context.Context, id int, channel syslog.LogPartsChannel) {
	deviceTicker := time.NewTicker(DeviceUpdateDuration)
	defer deviceTicker.Stop()
	checkTicker := time.NewTicker(1 * time.Second)
	defer checkTicker.Stop()

	ks.Infof("thread %d running", id)
	for {
		select {
		case logParts := <-channel:
			ks.metrics.Messages.Mark(1)
			msg, err := ks.formatMessage(ctx, logParts)
			if err != nil {
				ks.Errorf("Cannot format syslog: %v", err)
			}
			select {
			case ks.logchan <- string(msg):
			default:
				ks.metrics.Errors.Mark(1)
			}
		case <-deviceTicker.C: // Run these only on the 1st thread.
			if id == 1 {
				go func() {
					ks.Infof("Updating the device list.")
					ks.devices = ks.apic.GetDevicesAsMap(0)
				}()
			}
		case <-checkTicker.C:
			if id == 1 {
				ks.metrics.Queue.Update(int64(len(channel)))
			}
		case <-ctx.Done():
			ks.Infof("thread %d done", id)
			return
		}
	}
}

func (ks *KentikSyslog) run(ctx context.Context, host string) {
	ks.Infof("Server ready on %s", host)

	for i := 1; i <= *threads; i++ {
		go ks.process(ctx, i, ks.channel)
	}

	// Wait forever here.
	ks.server.Wait()
}

func (ks *KentikSyslog) formatMessage(ctx context.Context, msg sfmt.LogParts) ([]byte, error) {
	if client, ok := msg["client"].(string); ok { // Look up device_name here.
		pts := strings.Split(client, ":")
		if dev, ok := ks.devices[pts[0]]; ok {
			msg["device_name"] = dev.Name // Copy in any of these info we get
			dev.SetMsgUserTags(msg)
		}
		if ks.resolver != nil {
			msg["client_name"] = ks.resolver.Resolve(ctx, pts[0])
		}
	}

	// Fall back to hostname if this is set.
	if _, ok := msg["device_name"]; !ok {
		if hostname, ok := msg["hostname"]; ok {
			if hs, ok := hostname.(string); ok {
				if ipr := net.ParseIP(hs); ipr != nil {
					if ks.resolver != nil {
						msg["device_name"] = ks.resolver.Resolve(ctx, hs)
					} else {
						msg["device_name"] = hs
					}
				} else {
					msg["device_name"] = hs
				}
			}
		}
	}

	// One more time for sure.
	if _, ok := msg["device_name"]; !ok {
		msg["device_name"] = msg["client_name"]
	}

	msg["message"] = msg["content"] // Swap these around for NR.
	delete(msg, "content")
	msg["plugin.type"] = kt.PluginSyslog         // NR Processing.
	msg["instrumentation.name"] = InstNameSyslog // NR Processing.

	// Remove any empty strings.
	for k, v := range msg {
		if vs, ok := v.(string); ok {
			if vs == "" {
				delete(msg, k)
			}
		}
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return b, nil
}
