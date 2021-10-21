package syslog

import (
	"context"
	"flag"
	"fmt"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"gopkg.in/mcuadros/go-syslog.v2"
)

type KentikSyslog struct {
	logger.ContextL
	server  *syslog.Server
	handler *syslog.ChannelHandler
	channel syslog.LogPartsChannel
}

var (
	doUDP  = flag.Bool("syslog.udp", true, "Listen on UDP for syslog messages.")
	doTCP  = flag.Bool("syslog.tcp", false, "Listen on TCP for syslog messages.")
	doUnix = flag.Bool("syslog.unix", false, "Listen on a Unix socket for syslog messages.")
	format = flag.String("syslog.format", "RFC5424", "Format to parse syslog messages with. Options are: RFC3164|RFC5424|RFC6587.")
)

func NewSyslogSource(ctx context.Context, host string, log logger.Underlying, registry go_metrics.Registry) (*KentikSyslog, error) {
	ss := KentikSyslog{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "Syslog"}, log),
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
			ks.Infof("XXX %v", logParts)
		}
	}(ks.channel)

	ks.server.Wait()
}
