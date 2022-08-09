package prom

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	listen string
)

func init() {
	flag.StringVar(&listen, "prom_listen", ":8082", "Bind to listen for prometheus requests on.")
}

type PromSink struct {
	logger.ContextL
	registry go_metrics.Registry
	metrics  *PromMetric
	config   *ktranslate.PrometheusSinkConfig
}

type PromMetric struct {
}

func NewSink(log logger.Underlying, registry go_metrics.Registry, cfg *ktranslate.PrometheusSinkConfig) (*PromSink, error) {
	if cfg == nil {
		return nil, fmt.Errorf("prometheus sink config cannot be nil")
	}
	nr := PromSink{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "promSink"}, log),
		registry: registry,
		metrics:  &PromMetric{},
		config:   cfg,
	}

	return &nr, nil
}

func (s *PromSink) Init(ctx context.Context, format formats.Format, compression kt.Compression, fmtr formats.Formatter) error {

	if format != formats.FORMAT_PROM {
		return fmt.Errorf("Prometheus only supports prometheus format, not %s", format)
	}

	go s.listen(ctx)

	return nil
}

func (s *PromSink) HttpInfo() map[string]float64 {
	return map[string]float64{}
}

func (s *PromSink) Send(ctx context.Context, payload *kt.Output) {
	// Noop because already registered in the rollup phase.
}

func (s *PromSink) Close() {}

func (s *PromSink) PromInfo() map[string]float64 {
	return map[string]float64{}
}

func (s *PromSink) listen(ctx context.Context) {
	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))
	s.Infof("Prometheus listening on %s", s.config.ListenAddr)
	err := http.ListenAndServe(s.config.ListenAddr, nil)
	if err != nil {
		s.Errorf("Error with Prometheus -- %v", err)
	}
}
