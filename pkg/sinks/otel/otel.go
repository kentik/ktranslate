package otel

import (
	"context"
	"log/slog"

	jsoniter "github.com/json-iterator/go"
	go_metrics "github.com/kentik/go-metrics"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"
)

/**
This is a stub sink because otel format handles most of the outputs.

Just here to allow logs and syslog values from the logTee to go outbound.
*/

var json = jsoniter.ConfigFastest

const (
	HttpHostname = "ktranslate"
)

type OtelSink struct {
	logger.ContextL

	registry go_metrics.Registry
	metrics  *OtelMetric
	logTee   chan string
}

type OtelMetric struct {
	DeliveryLogs go_metrics.Meter
}

func NewSink(log logger.Underlying, registry go_metrics.Registry, logTee chan string) (*OtelSink, error) {
	nr := OtelSink{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "otelSink"}, log),
		registry: registry,
		metrics: &OtelMetric{
			DeliveryLogs: go_metrics.GetOrRegisterMeter("delivery_logs_otel", registry),
		},
		logTee: logTee,
	}

	return &nr, nil
}

func (s *OtelSink) Init(ctx context.Context, format formats.Format, compression kt.Compression, fmtr formats.Formatter) error {

	// Send logs on if this is set.
	if s.logTee != nil {
		go s.watchLogs(ctx)
	}

	s.Infof("Exporting logs via otel")
	return nil
}

func (s *OtelSink) Send(ctx context.Context, payload *kt.Output) {
	// Noop here because we don't expect to send anything from otel format.
}

func (s *OtelSink) Close() {}

func (s *OtelSink) HttpInfo() map[string]float64 {
	return map[string]float64{
		"DeliveryLogs": s.metrics.DeliveryLogs.Rate1(),
	}
}

// Forwards any logs recieved to the NR log API.
func (s *OtelSink) watchLogs(ctx context.Context) {
	s.Infof("Receiving logs...")
	for {
		select {
		case log := <-s.logTee:
			slog.LogAttrs(ctx, slog.LevelInfo, log)
			s.metrics.DeliveryLogs.Mark(1)
		case <-ctx.Done():
			s.Infof("Logs received")
			return
		}
	}
}
