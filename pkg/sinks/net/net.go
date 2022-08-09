package net

import (
	"context"
	"flag"
	"fmt"
	"net"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"
)

var (
	server   string
	protocol string
)

func init() {
	flag.StringVar(&server, "net_server", "", "Write flows seen to this address (host and port)")
	flag.StringVar(&protocol, "net_protocol", "udp", "Use this protocol for writing data (udp|tcp|unix)")
}

type NetSink struct {
	logger.ContextL
	conn     net.Conn
	registry go_metrics.Registry
	metrics  *NetMetric
	config   *ktranslate.NetSinkConfig
}

type NetMetric struct {
	DeliveryErr go_metrics.Meter
	DeliveryWin go_metrics.Meter
}

func NewSink(log logger.Underlying, registry go_metrics.Registry, cfg *ktranslate.NetSinkConfig) (*NetSink, error) {
	return &NetSink{
		registry: registry,
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "netSink"}, log),
		metrics: &NetMetric{
			DeliveryErr: go_metrics.GetOrRegisterMeter("delivery_errors_net", registry),
			DeliveryWin: go_metrics.GetOrRegisterMeter("delivery_wins_net", registry),
		},
		config: cfg,
	}, nil
}

func (s *NetSink) Init(ctx context.Context, format formats.Format, compression kt.Compression, fmtr formats.Formatter) error {
	if s.config.Endpoint == "" {
		return fmt.Errorf("Net requires -net_server or NetSink.Endpoint to be set")
	}

	var serverAddr net.Addr
	var err error
	switch s.config.Protocol {
	case "udp":
		serverAddr, err = net.ResolveUDPAddr(s.config.Protocol, s.config.Endpoint)
	case "tcp":
		serverAddr, err = net.ResolveTCPAddr(s.config.Protocol, s.config.Endpoint)
	case "unix":
		serverAddr, err = net.ResolveUnixAddr(s.config.Protocol, s.config.Endpoint)
	default:
		err = fmt.Errorf("Invalid protocol: %s. Supported: udp|tcp|unix", s.config.Protocol)

	}
	if err != nil {
		return err
	}

	conn, err := (&net.Dialer{}).DialContext(ctx, s.config.Protocol, serverAddr.String())
	if err != nil {
		return err
	}

	s.conn = conn
	s.Infof("Network: sending to %s:%s", s.config.Protocol, s.config.Endpoint)

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
