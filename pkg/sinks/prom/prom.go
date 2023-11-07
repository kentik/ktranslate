package prom

import (
	"bytes"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	listen    string
	remoteUrl string
)

func init() {
	flag.StringVar(&listen, "prom_listen", "127.0.0.1:8083", "Bind to listen for prometheus requests on.")
	flag.StringVar(&remoteUrl, "prom_remote_write", "", "Pass on remote write to this address.")
}

type PromSink struct {
	logger.ContextL
	registry    go_metrics.Registry
	metrics     *PromMetric
	config      *ktranslate.PrometheusSinkConfig
	client      *http.Client
	tr          *http.Transport
	compression kt.Compression
	remoteUrl   string
}

type PromMetric struct {
	DeliveryErr     go_metrics.Meter
	DeliveryWin     go_metrics.Meter
	DeliveryMetrics go_metrics.Meter
	DeliveryLogs    go_metrics.Meter
}

func NewSink(log logger.Underlying, registry go_metrics.Registry, cfg *ktranslate.PrometheusSinkConfig) (*PromSink, error) {
	if cfg == nil {
		return nil, fmt.Errorf("prometheus sink config cannot be nil")
	}
	nr := PromSink{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "promSink"}, log),
		registry: registry,
		metrics: &PromMetric{
			DeliveryErr:     go_metrics.GetOrRegisterMeter("delivery_errors_prom", registry),
			DeliveryWin:     go_metrics.GetOrRegisterMeter("delivery_wins_prom", registry),
			DeliveryMetrics: go_metrics.GetOrRegisterMeter("delivery_metrics_prom", registry),
			DeliveryLogs:    go_metrics.GetOrRegisterMeter("delivery_logs_prom", registry),
		},
		config: cfg,
	}

	return &nr, nil
}

func (s *PromSink) Init(ctx context.Context, format formats.Format, compression kt.Compression, fmtr formats.Formatter) error {
	switch format {
	case formats.FORMAT_PROM:
		go s.listen(ctx)
	case formats.FORMAT_PROM_REMOTE:
		s.tr = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: false}, // TODO, any time that we want this to be false?
		}
		s.client = &http.Client{Transport: s.tr}
		if remoteUrl == "" {
			return fmt.Errorf("You must set the -prom_remote_write flag to make this work.")
		}

		if compression != kt.CompressionSnappy {
			return fmt.Errorf("You used the %s unsupported compression format. Use snappy only.", compression)
		}

		s.Infof("Sending to remote_write endpoint %s", remoteUrl)
	default:
		return fmt.Errorf("Prometheus only supports %s and %s formats, not %s", formats.FORMAT_PROM, formats.FORMAT_PROM_REMOTE, format)
	}

	s.remoteUrl = remoteUrl
	s.compression = compression

	return nil
}

func (s *PromSink) HttpInfo() map[string]float64 {
	return map[string]float64{
		"DeliveryErr":       s.metrics.DeliveryErr.Rate1(),
		"DeliveryWin":       s.metrics.DeliveryWin.Rate1(),
		"DeliveryMetrics1":  s.metrics.DeliveryMetrics.Rate1(),
		"DeliveryMetrics15": s.metrics.DeliveryMetrics.Rate15(),
		"DeliveryLogs":      s.metrics.DeliveryLogs.Rate1(),
	}
}

func (s *PromSink) Send(ctx context.Context, payload *kt.Output) {
	if s.client != nil {
		go s.sendRemote(ctx, payload, s.remoteUrl)
	}
}

func (s *PromSink) sendRemote(ctx context.Context, payload *kt.Output, url string) {
	var cbErr error = nil
	if payload.CB != nil { // Let anyone who asked know that this has been sent
		defer func() {
			payload.CB(cbErr) // This needs to get wrapped in its own function to get cbErr right.
		}()
	}

	s.metrics.DeliveryMetrics.Mark(1) // Compression will effect this, but we can do our best.
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload.Body))
	if err != nil {
		s.Errorf("There was an error when communicating to Remote Prometheus: %v.", err)
		cbErr = err
		return
	}

	req.Header.Add("X-Prometheus-Remote-Write-Version", "0.1.0")
	req.Header.Set("Content-Type", "application/x-protobuf")

	if s.compression == kt.CompressionSnappy {
		req.Header.Add("Content-Encoding", "snappy")
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.Errorf("There was an error when creating a new client in Prometheus: %v.", err)
		cbErr = err
		s.client = &http.Client{Transport: s.tr}
	} else {
		defer resp.Body.Close()
		bdy, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			s.Errorf("There was an error when communicating to Prometheus: %v.", err)
			cbErr = err
			s.metrics.DeliveryErr.Mark(1)
		} else {
			if resp.StatusCode >= 400 {
				s.Errorf("There was an error when communicating to Prometheus: %v.", resp.StatusCode)
				cbErr = fmt.Errorf("There was an error when communicating to Prometheus: %v.", resp.StatusCode)
				s.metrics.DeliveryErr.Mark(1)
			} else {
				s.Debugf("Prom Success: %d %s", resp.StatusCode, string(bdy))
				s.metrics.DeliveryWin.Mark(1)
			}
		}
	}
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
