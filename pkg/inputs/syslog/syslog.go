package syslog

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"gopkg.in/mcuadros/go-syslog.v2"
	sfmt "gopkg.in/mcuadros/go-syslog.v2/format"
)

type KentikSyslog struct {
	logger.ContextL
	server  *syslog.Server
	handler *syslog.ChannelHandler
	channel syslog.LogPartsChannel
	logchan chan string
}

var (
	doUDP  = flag.Bool("syslog.udp", true, "Listen on UDP for syslog messages.")
	doTCP  = flag.Bool("syslog.tcp", true, "Listen on TCP for syslog messages.")
	doUnix = flag.Bool("syslog.unix", false, "Listen on a Unix socket for syslog messages.")
	format = flag.String("syslog.format", "Automatic", "Format to parse syslog messages with. Options are: Automatic|RFC3164|RFC5424|RFC6587.")
)

func NewSyslogSource(ctx context.Context, host string, log logger.Underlying, logchan chan string, registry go_metrics.Registry) (*KentikSyslog, error) {
	ss := KentikSyslog{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "Syslog"}, log),
		logchan:  logchan,
	}

	if logchan == nil {
		return nil, fmt.Errorf("Log Sink is not set.")
	}

	channel := make(syslog.LogPartsChannel)
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
		ss.Infof("Listening for UDP on %s", host)
	}
	if *doTCP {
		err := server.ListenTCP(host)
		if err != nil {
			return nil, fmt.Errorf("Cannot listen for syslog with tcp: %v", err)
		}
		ss.Infof("Listening for TCP on %s", host)
	}
	if *doUnix {
		err := server.ListenUnixgram(host)
		if err != nil {
			return nil, fmt.Errorf("Cannot listen for syslog with unixgram: %v", err)
		}
		ss.Infof("Listening for Unixgram on %s", host)
	}

	err := server.Boot()
	if err != nil {
		return nil, fmt.Errorf("Cannot boot syslog server: %v", err)
	}
	ss.server = server
	ss.channel = channel
	ss.handler = handler

	go ss.run(ctx, host)

	return &ss, nil
}

func (ks *KentikSyslog) run(ctx context.Context, host string) {
	ks.Infof("Server ready on %s", host)

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			msg, err := formatMessage(logParts)
			if err != nil {
				ks.Errorf("Cannot format syslog: %v", err)
			}
			ks.logchan <- string(msg)
		}
	}(ks.channel)

	ks.server.Wait()
}

func formatMessage(msg sfmt.LogParts) ([]byte, error) {
	msg["message"] = msg["content"] // Swap these around for NR.
	delete(msg, "content")
	b, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return b, nil
}
