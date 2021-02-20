package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"
)

type HttpSink struct {
	logger.ContextL
	TargetUrl string

	client   *http.Client
	tr       *http.Transport
	registry go_metrics.Registry
	metrics  *HttpMetric
	headers  map[string]string
}

type HttpMetric struct {
	DeliveryErr go_metrics.Meter
	DeliveryWin go_metrics.Meter
}

type HeaderFlag []string

func (h *HeaderFlag) String() string {
	return strings.Join(*h, ",")
}

func (h *HeaderFlag) Set(value string) error {
	*h = append(*h, value)
	return nil
}

var (
	TargetUrl = flag.String("http_url", "http://localhost:8086/write?db=kentik", "URL to post to")
	headers   HeaderFlag
)

func NewSink(log logger.Underlying, registry go_metrics.Registry, sink string) (*HttpSink, error) {
	nr := HttpSink{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "httpSink"}, log),
		registry: registry,
		metrics: &HttpMetric{
			DeliveryErr: go_metrics.GetOrRegisterMeter("delivery_errors_http", registry),
			DeliveryWin: go_metrics.GetOrRegisterMeter("delivery_wins_http", registry),
		},
		headers: map[string]string{},
	}

	for _, header := range headers {
		pts := strings.SplitN(header, ":", 2)
		if len(pts) > 1 {
			nr.headers[strings.TrimSpace(pts[0])] = strings.TrimSpace(pts[1])
		} else {
			return nil, fmt.Errorf("Invalid header: %s", header)
		}
	}

	for k, v := range nr.headers {
		nr.Infof(`Adding HTTP header "%s: %s"`, k, v)
	}

	if sink == "splunk" {
		if _, ok := nr.headers["Authorization"]; !ok {
			return nil, fmt.Errorf("Authorization header required for splunk")
		}
	}

	return &nr, nil
}

func (s *HttpSink) Init(ctx context.Context, format formats.Format, compression kt.Compression) error {
	s.TargetUrl = *TargetUrl

	s.tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	s.client = &http.Client{Transport: s.tr}

	if compression == kt.CompressionGzip {
		s.headers["Content-Encoding"] = "GZIP"
	}

	s.Infof("Exporting to Http at %s", s.TargetUrl)

	return nil
}

func (s *HttpSink) Send(ctx context.Context, payload []byte) {
	go s.sendHttp(ctx, payload)
}

func (s *HttpSink) Close() {}

func (s *HttpSink) HttpInfo() map[string]float64 {
	return map[string]float64{
		"DeliveryErr": s.metrics.DeliveryErr.Rate1(),
		"DeliveryWin": s.metrics.DeliveryWin.Rate1(),
	}
}

func (s *HttpSink) sendHttp(ctx context.Context, payload []byte) {
	req, err := http.NewRequestWithContext(ctx, "POST", s.TargetUrl, bytes.NewBuffer(payload))
	if err != nil {
		s.Errorf("Cannot create HTTP request: %v", err)
		return
	}

	for k, v := range s.headers {
		req.Header.Set(k, v)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.Errorf("Cannot write to HTTP: %v, creating new client", err)
		s.client = &http.Client{Transport: s.tr}
	} else {
		defer resp.Body.Close()
		bdy, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			s.Errorf("Cannot get resp body from HTTP: %v", err)
			s.metrics.DeliveryErr.Mark(1)
		} else {
			if resp.StatusCode >= 400 {
				s.Errorf("Cannot write to HTTP, status code %d, bdy: %s", resp.StatusCode, string(bdy))
				s.metrics.DeliveryErr.Mark(1)
			} else {
				s.metrics.DeliveryWin.Mark(1)
			}
		}
	}
}
