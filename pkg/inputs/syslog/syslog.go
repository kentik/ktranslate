package syslog

import (
	"context"
	"flag"
	"fmt"
	"net"
	"strings"
	"time"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate"
	syslog "github.com/kentik/the-library-formally-known-as-go-syslog"
	sfmt "github.com/kentik/the-library-formally-known-as-go-syslog/format"

	"github.com/kentik/ktranslate/pkg/api"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/inputs/snmp"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/util/resolv"
)

var (
	doUDP   bool
	doTCP   bool
	doUnix  bool
	format  string
	threads int
)

func init() {
	flag.BoolVar(&doUDP, "syslog.udp", true, "Listen on UDP for syslog messages.")
	flag.BoolVar(&doTCP, "syslog.tcp", true, "Listen on TCP for syslog messages.")
	flag.BoolVar(&doUnix, "syslog.unix", false, "Listen on a Unix socket for syslog messages.")
	flag.StringVar(&format, "syslog.format", "Automatic", "Format to parse syslog messages with. Options are: Automatic|RFC3164|RFC5424|RFC6587|NoFormat.")
	flag.IntVar(&threads, "syslog.threads", 1, "Number of threads to use to process messages.")
}

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
	config   *ktranslate.SyslogInputConfig
}

type SyslogMetric struct {
	Messages go_metrics.Meter
	Errors   go_metrics.Meter
	Queue    go_metrics.Gauge
}

const (
	CHAN_SLACK           = 10000
	DeviceUpdateDuration = 1 * time.Hour
	InstNameSyslog       = "ktranslate-syslog"
	ErrorCheckDuration   = 1 * time.Minute
)

func NewSyslogSource(ctx context.Context, log logger.Underlying, logchan chan string, registry go_metrics.Registry, apic *api.KentikApi, resolver *resolv.Resolver, cfg *ktranslate.SyslogInputConfig) (*KentikSyslog, error) {
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
		config:   cfg,
	}

	if logchan == nil {
		return nil, fmt.Errorf("Log Sink is not set.")
	}

	channel := make(syslog.LogPartsChannel, CHAN_SLACK)
	handler := syslog.NewChannelHandler(channel)
	server := syslog.NewServer()
	switch cfg.Format {
	case "RFC3164":
		server.SetFormat(syslog.RFC3164)
	case "RFC5424":
		server.SetFormat(syslog.RFC5424)
	case "RFC6587":
		server.SetFormat(syslog.RFC6587)
	case "Automatic":
		server.SetFormat(syslog.Automatic)
	case "NoFormat":
		server.SetFormat(syslog.NoFormat)
	default:
		return nil, fmt.Errorf("Invalid syslog format (%s). Options are Automatic|RFC3164|RFC5424|RFC6587|NoFormat", cfg.Format)
	}

	server.SetHandler(handler)
	if cfg.EnableUDP {
		err := server.ListenUDP(cfg.ListenAddr)
		if err != nil {
			return nil, fmt.Errorf("Cannot listen for syslog udp: %v", err)
		}
		ks.Infof("Listening for UDP on %s", cfg.ListenAddr)
	}
	if cfg.EnableTCP {
		err := server.ListenTCP(cfg.ListenAddr)
		if err != nil {
			return nil, fmt.Errorf("Cannot listen for syslog with tcp: %v", err)
		}
		ks.Infof("Listening for TCP on %s", cfg.ListenAddr)
	}
	if cfg.EnableUnix {
		err := server.ListenUnixgram(cfg.ListenAddr)
		if err != nil {
			return nil, fmt.Errorf("Cannot listen for syslog with unixgram: %v", err)
		}
		ks.Infof("Listening for Unixgram on %s", cfg.ListenAddr)
	}

	err := server.Boot()
	if err != nil {
		return nil, fmt.Errorf("Cannot boot syslog server: %v", err)
	}
	ks.server = server
	ks.channel = channel
	ks.handler = handler

	go ks.run(ctx, cfg.ListenAddr)

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
	errorTicker := time.NewTicker(ErrorCheckDuration)
	defer errorTicker.Stop()
	var lastErr error

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
		case <-errorTicker.C: // See if there's any new errors?
			lm := ks.server.GetLastError()
			if lm != nil && lm != lastErr {
				lastErr = lm
				ks.Errorf("%v", lastErr)
				ks.metrics.Errors.Mark(1)
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

	for i := 1; i <= ks.config.Threads; i++ {
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
		} else if snmp.ServiceName != "" {
			msg["tags.container_service"] = snmp.ServiceName
		}

		if ks.resolver != nil {
			msg["client_name"] = ks.resolver.Resolve(ctx, pts[0], true)
		}
	}

	// Fall back to hostname if this is set.
	if _, ok := msg["device_name"]; !ok {
		if hostname, ok := msg["hostname"]; ok {
			if hs, ok := hostname.(string); ok {
				if ipr := net.ParseIP(hs); ipr != nil {
					// First check if this ip is in our devices list.
					if dev, ok := ks.devices[hs]; ok {
						msg["device_name"] = dev.Name // Copy in any of these info we get
						dev.SetMsgUserTags(msg)
					} else { // If not, try to resolve via dns.
						if ks.resolver != nil {
							resolved_name := ks.resolver.Resolve(ctx, hs, true)
							if resolved_name != "" {
								msg["device_name"] = resolved_name
							} else {
								msg["device_name"] = hs
							}
						} else {
							msg["device_name"] = hs
						}
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

	// Swap these around for NR.
	if _, ok := msg["message"]; !ok {
		if _, ok := msg["content"]; ok {
			msg["message"] = msg["content"]
			delete(msg, "content")
		}
	}

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
