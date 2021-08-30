package net

import (
	"context"
	"flag"
	"fmt"
	"net"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"
)

type NetSink struct {
	logger.ContextL
	conn     net.Conn
	registry go_metrics.Registry
	metrics  *NetMetric
}

type NetMetric struct {
	DeliveryErr go_metrics.Meter
	DeliveryWin go_metrics.Meter
}

var (
	server   = flag.String("net_server", "", "Write flows seen to this address (host and port)")
	protocol = flag.String("net_protocol", "udp", "Use this protocol for writing data (udp|tcp|unix)")
)

func NewSink(log logger.Underlying, registry go_metrics.Registry) (*NetSink, error) {
	return &NetSink{
		registry: registry,
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "netSink"}, log),
		metrics: &NetMetric{
			DeliveryErr: go_metrics.GetOrRegisterMeter("delivery_errors_net", registry),
			DeliveryWin: go_metrics.GetOrRegisterMeter("delivery_wins_net", registry),
		},
	}, nil
}

func (s *NetSink) Init(ctx context.Context, format formats.Format, compression kt.Compression, fmtr formats.Formatter) error {
	if *server == "" {
		return fmt.Errorf("Net requires -net_server to be set")
	}

	var serverAddr net.Addr
	var err error
	switch *protocol {
	case "udp":
		serverAddr, err = net.ResolveUDPAddr(*protocol, *server)
	case "tcp":
		serverAddr, err = net.ResolveTCPAddr(*protocol, *server)
	case "unix":
		serverAddr, err = net.ResolveUnixAddr(*protocol, *server)
	default:
		err = fmt.Errorf("Invalid protocol: %s. Supported: udp|tcp|unix", *protocol)

	}
	if err != nil {
		return err
	}

	conn, err := (&net.Dialer{}).DialContext(ctx, *protocol, serverAddr.String())
	if err != nil {
		return err
	}

	s.conn = conn
	s.Infof("Network: sending to %s:%s", *protocol, *server)

	return nil
}

func (s *NetSink) Send(ctx context.Context, payload *kt.Output) {
	_, err := s.conn.Write(payload.Body)
	if err != nil {
		s.Errorf("There was an error when writing: %v.", err)
		s.metrics.DeliveryErr.Mark(1)
	} else {
		s.metrics.DeliveryWin.Mark(1)
	}
}

func (s *NetSink) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *NetSink) HttpInfo() map[string]float64 {
	return map[string]float64{
		"DeliveryErr": s.metrics.DeliveryErr.Rate1(),
		"DeliveryWin": s.metrics.DeliveryWin.Rate1(),
	}
}
