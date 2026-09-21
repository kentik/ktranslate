package relay

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"
)

const (
	CHF_TYPE = "application/chf"

	DefaultSendTimeout = 30 * time.Second
)

// RelaySink forwards a copy of outbound flow to another ktranslate instance's
// http.source listener, for chaining instances together. It's only ever
// instantiated when Config.TeeFlow is set (see pkg/cat/kkc.go) — it isn't a
// selectable --sinks destination.
type RelaySink struct {
	logger.ContextL
	registry        go_metrics.Registry
	metrics         *RelayMetric
	relayUrl        string
	client          *http.Client
	tr              *http.Transport
	config          *ktranslate.Config
	sendMaxDuration time.Duration
	compression     kt.Compression
}

type RelayMetric struct {
	DeliveryErr go_metrics.Meter
	DeliveryWin go_metrics.Meter
}

func NewSink(log logger.Underlying, registry go_metrics.Registry, cfg *ktranslate.Config) (*RelaySink, error) {
	return &RelaySink{
		registry: registry,
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "relaySink"}, log),
		metrics: &RelayMetric{
			DeliveryErr: go_metrics.GetOrRegisterMeter("delivery_errors_relay", registry),
			DeliveryWin: go_metrics.GetOrRegisterMeter("delivery_wins_relay", registry),
		},
		sendMaxDuration: DefaultSendTimeout,
		config:          cfg,
	}, nil
}

func (s *RelaySink) Init(ctx context.Context, format formats.Format, compression kt.Compression, fmtr formats.Formatter) error {
	if s.config.TeeFlow == "" {
		return fmt.Errorf("relay sink requires tee_flow to be set")
	}
	s.relayUrl = s.config.TeeFlow

	s.tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	s.client = &http.Client{Transport: s.tr}
	s.compression = compression

	s.Infof("Relaying flow to %s", s.relayUrl)

	return nil
}

func (s *RelaySink) Send(ctx context.Context, payload *kt.Output) {
	go func() {
		ctxC, cancel := context.WithTimeout(ctx, s.sendMaxDuration)
		defer cancel()
		s.sendRelay(ctxC, payload.Body, int(payload.Ctx.CompanyId), payload.Ctx.SenderId)
	}()
}

func (s *RelaySink) Close() {}

func (s *RelaySink) HttpInfo() map[string]float64 {
	return map[string]float64{
		"DeliveryErr": s.metrics.DeliveryErr.Rate1(),
		"DeliveryWin": s.metrics.DeliveryWin.Rate1(),
	}
}

func (s *RelaySink) sendRelay(ctx context.Context, payload []byte, cid int, senderId string) {
	vals := url.Values{}
	vals.Set("sid", strconv.Itoa(cid))
	vals.Set("sender_id", senderId)
	valString := vals.Encode()
	fullUrl := s.relayUrl + "?" + valString

	req, err := http.NewRequestWithContext(ctx, "POST", fullUrl, bytes.NewBuffer(payload))
	if err != nil {
		s.Errorf("Cannot create relay request: %v", err)
		return
	}

	req.Header.Set("Content-Type", CHF_TYPE)
	if s.compression == kt.CompressionGzip {
		req.Header.Set("Content-Encoding", "gzip")
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.Errorf("Cannot relay flow: %v, creating new client, URL=%s", err, fullUrl)
		s.client = &http.Client{Transport: s.tr}
	} else {
		defer resp.Body.Close()
		_, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			s.Errorf("Cannot get resp body from relay: %v", err)
			s.metrics.DeliveryErr.Mark(1)
		} else {
			if resp.StatusCode != 200 {
				s.Errorf("Cannot relay flow, status code %d", resp.StatusCode)
				s.metrics.DeliveryErr.Mark(1)
			} else {
				s.metrics.DeliveryWin.Mark(1)
			}
		}
	}
}
