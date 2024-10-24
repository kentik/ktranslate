package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	nurl "net/url"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	go_metrics "github.com/kentik/go-metrics"
	"github.com/pkg/errors"

	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"
)

var json = jsoniter.ConfigFastest

var (
	targetURL          string
	targetLogURL       string
	insecureSkipVerify bool
	timeoutSec         int
	headers            HeaderFlag
)

const (
	HttpHostname = "ktranslate"
)

func init() {
	flag.StringVar(&targetURL, "http_url", "http://localhost:8086/write?db=kentik", "URL to post to")
	flag.StringVar(&targetLogURL, "http_log_url", "http://localhost:8088/services/collector/event", "URL to post logs to")
	flag.BoolVar(&insecureSkipVerify, "http_insecure", false, "Allow insecure urls.")
	flag.IntVar(&timeoutSec, "http_timeout_sec", 30, "Timeout each request after this long.")
	flag.Var(&headers, "http_header", "Any custom http headers to set on outbound requests")
}

type HttpSink struct {
	logger.ContextL
	TargetUrl  string
	HttpLogUrl string

	client          *http.Client
	tr              *http.Transport
	registry        go_metrics.Registry
	metrics         *HttpMetric
	headers         map[string]string
	targetUrls      []string
	sendMaxDuration time.Duration
	config          *ktranslate.HTTPSinkConfig
	logTee          chan string
}

type HttpMetric struct {
	DeliveryErr  go_metrics.Meter
	DeliveryWin  go_metrics.Meter
	DeliveryLogs go_metrics.Meter
}

type HeaderFlag []string

func (h *HeaderFlag) String() string {
	return strings.Join(*h, ",")
}

func (h *HeaderFlag) Set(value string) error {
	*h = append(*h, value)
	return nil
}

func NewSink(log logger.Underlying, registry go_metrics.Registry, sink string, cfg *ktranslate.HTTPSinkConfig, logTee chan string) (*HttpSink, error) {
	nr := HttpSink{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "httpSink"}, log),
		registry: registry,
		metrics: &HttpMetric{
			DeliveryErr:  go_metrics.GetOrRegisterMeter("delivery_errors_http", registry),
			DeliveryWin:  go_metrics.GetOrRegisterMeter("delivery_wins_http", registry),
			DeliveryLogs: go_metrics.GetOrRegisterMeter("delivery_logs_http", registry),
		},
		headers:         map[string]string{},
		targetUrls:      []string{},
		sendMaxDuration: time.Duration(cfg.TimeoutInSeconds) * time.Second,
		config:          cfg,
		logTee:          logTee,
		HttpLogUrl:      cfg.TargetLogs,
	}

	for _, u := range strings.Split(cfg.Target, ",") {
		if _, err := nurl.Parse(u); err != nil {
			return nil, errors.Wrapf(err, "invalid url for http sink: %s", u)
		}
		nr.Infof("Exporting HTTP to %s", u)
		nr.targetUrls = append(nr.targetUrls, u)
	}

	for _, header := range cfg.Headers {
		pts := strings.SplitN(header, ":", 2)
		if len(pts) > 1 {
			nr.headers[strings.TrimSpace(pts[0])] = strings.TrimSpace(pts[1])
		} else {
			return nil, fmt.Errorf("The %s header contains unsupported content.", header)
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

func (s *HttpSink) Init(ctx context.Context, format formats.Format, compression kt.Compression, fmtr formats.Formatter) error {
	s.tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: s.config.InsecureSkipVerify},
	}
	s.client = &http.Client{Transport: s.tr}

	if compression == kt.CompressionGzip {
		s.headers["Content-Encoding"] = "GZIP"
	}

	// Send logs on if this is set.
	if s.logTee != nil {
		go s.watchLogs(ctx)
	}

	s.Infof("Exporting via HTTP at %v, logs %s", s.targetUrls, s.HttpLogUrl)

	return nil
}

func (s *HttpSink) Send(ctx context.Context, payload *kt.Output) {
	for _, url := range s.targetUrls {
		go func(u string) {
			ctxC, cancel := context.WithTimeout(ctx, s.sendMaxDuration)
			defer cancel()
			s.sendHttp(ctxC, payload.Body, u)
		}(url)
	}
}

func (s *HttpSink) Close() {}

func (s *HttpSink) HttpInfo() map[string]float64 {
	return map[string]float64{
		"DeliveryErr":  s.metrics.DeliveryErr.Rate1(),
		"DeliveryWin":  s.metrics.DeliveryWin.Rate1(),
		"DeliveryLogs": s.metrics.DeliveryLogs.Rate1(),
	}
}

func (s *HttpSink) sendHttp(ctx context.Context, payload []byte, url string) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		s.Errorf("There was an error when creating an HTTP request: %v.", err)
		return
	}

	for k, v := range s.headers {
		req.Header.Set(k, v)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.Errorf("There was an error when sending via HTTP at %s: %v.", url, err)
		s.client = &http.Client{Transport: s.tr}
	} else {
		defer resp.Body.Close()
		bdy, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			s.Errorf("There was an error when getting the HTTP at %s response body: %v.", url, err)
			s.metrics.DeliveryErr.Mark(1)
		} else {
			if resp.StatusCode >= 400 {
				s.Errorf("There was an error when converting to HTTP at %s: %d. Body: %s.", url, resp.StatusCode, string(bdy))
				s.metrics.DeliveryErr.Mark(1)
			} else {
				s.metrics.DeliveryWin.Mark(1)
			}
		}
	}
}

// Forwards any logs recieved to the NR log API.
func (s *HttpSink) watchLogs(ctx context.Context) {
	s.Infof("Receiving logs...")
	logTicker := time.NewTicker(1 * time.Second)
	defer logTicker.Stop()
	batch := make([]string, 0, 100)
	for {
		select {
		case log := <-s.logTee:
			batch = append(batch, log)
			s.metrics.DeliveryLogs.Mark(1)
		case _ = <-logTicker.C:
			if len(batch) > 0 {
				ob := batch
				batch = make([]string, 0, 100)
				go s.sendLogBatch(ctx, ob)
			}
		case <-ctx.Done():
			s.Infof("Logs received")
			return
		}
	}
}

type log struct {
	Host       string `json:"host"`
	SourceType string `json:"sourcetype"`
	Timestamp  int64  `json:"time"`
	Event      string `json:"event"`
}

func (s *HttpSink) sendLogBatch(ctx context.Context, logs []string) {
	ts := time.Now().Unix()
	var buf bytes.Buffer
	for _, l := range logs {
		ll := log{
			Timestamp:  ts,
			SourceType: string(kt.ProviderLogs),
			Event:      l,
			Host:       HttpHostname,
		}
		target, err := json.Marshal(ll)
		if err != nil {
			s.Errorf("There was an error with logs: %v.", err)
			return
		}
		buf.Write(target)
		buf.WriteString("\n") // Because splunk
	}

	s.sendHttp(ctx, buf.Bytes(), s.HttpLogUrl)
}
