package kentik

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"
)

const (
	CHF_TYPE = "application/chf"
)

var (
	relayUrl string
)

func init() {
	flag.StringVar(&relayUrl, "kentik_relay_url", "", "If set, override the kentik api url to send flow over here.")
}

type KentikSink struct {
	logger.ContextL
	registry  go_metrics.Registry
	metrics   *KentikMetric
	KentikUrl string
	client    *http.Client
	tr        *http.Transport
	conf      *kt.KentikConfig
	isKentik  bool
	config    *ktranslate.KentikSinkConfig
}

type KentikMetric struct {
	DeliveryErr go_metrics.Meter
	DeliveryWin go_metrics.Meter
}

func NewSink(log logger.Underlying, registry go_metrics.Registry, conf *kt.KentikConfig, cfg *ktranslate.KentikSinkConfig) (*KentikSink, error) {
	return &KentikSink{
		registry: registry,
		conf:     conf,
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "kentikSink"}, log),
		metrics: &KentikMetric{
			DeliveryErr: go_metrics.GetOrRegisterMeter("delivery_errors_kentik", registry),
			DeliveryWin: go_metrics.GetOrRegisterMeter("delivery_wins_kentik", registry),
		},
		config: cfg,
	}, nil
}

func (s *KentikSink) Init(ctx context.Context, format formats.Format, compression kt.Compression, fmtr formats.Formatter) error {
	if s.conf == nil || s.conf.ApiEmail == "" || s.conf.ApiToken == "" {
		return fmt.Errorf("Kentik requires -kentik_email and KENTIK_API_TOKEN env var to be set")
	}
	s.KentikUrl = strings.ReplaceAll(s.conf.ApiRoot, "api.", "flow.") + "/chf"
	if v := s.config.RelayURL; v != "" { // If this is set, override and go directly here instead.
		s.KentikUrl = v
	}

	s.isKentik = strings.Contains(strings.ToLower(s.KentikUrl), "kentik.com") // Make sure we can't feed data back into kentik in a loop.

	s.tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	s.client = &http.Client{Transport: s.tr}

	s.Infof("Exporting to Kentik at %s (isKentik=%v)", s.KentikUrl, s.isKentik)

	return nil
}

func (s *KentikSink) Send(ctx context.Context, payload *kt.Output) {
	// Noop, can't send this way.
}

func (s *KentikSink) Close() {}

func (s *KentikSink) HttpInfo() map[string]float64 {
	return map[string]float64{
		"DeliveryErr": s.metrics.DeliveryErr.Rate1(),
		"DeliveryWin": s.metrics.DeliveryWin.Rate1(),
	}
}

func (s *KentikSink) SendKentik(payload []byte, cid int, senderId string, offset int) {
	if s.isKentik && offset == 0 { // Cut short any flow which is coming from kentik going back to kentik.
		return
	}

	vals := url.Values{}
	vals.Set("sid", strconv.Itoa(cid))
	vals.Set("sender_id", senderId)
	valString := vals.Encode()
	fullUrl := s.KentikUrl + "?" + valString

	gziped, err := s.gzBuf(nil, payload)
	if err != nil {
		s.Errorf("Cannot compress Kentik forward: %v", err)
		return
	}
	req, err := http.NewRequestWithContext(context.Background(), "POST", fullUrl, bytes.NewBuffer(gziped))
	if err != nil {
		s.Errorf("Cannot create Kentik request: %v", err)
		return
	}

	req.Header.Set("X-CH-Auth-Email", s.conf.ApiEmail)
	req.Header.Set("X-CH-Auth-API-Token", s.conf.ApiEmail)
	req.Header.Set("Content-Type", CHF_TYPE)
	req.Header.Set("Content-Encoding", "gzip")

	resp, err := s.client.Do(req)
	if err != nil {
		s.Errorf("Cannot write to Kentik: %v, creating new client, URL=%s", err, fullUrl)
		s.client = &http.Client{Transport: s.tr}
	} else {
		defer resp.Body.Close()
		_, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			s.Errorf("Cannot get resp body from Kentik: %v", err)
			s.metrics.DeliveryErr.Mark(1)
		} else {
			if resp.StatusCode != 200 {
				s.Errorf("Cannot write to Kentik, status code %d", resp.StatusCode)
				s.metrics.DeliveryErr.Mark(1)
			} else {
				s.metrics.DeliveryWin.Mark(1)
			}
		}
	}
}

func (s *KentikSink) gzBuf(serBuf []byte, raw []byte) ([]byte, error) {
	if serBuf == nil {
		serBuf = make([]byte, len(raw))
	}
	buf := bytes.NewBuffer(serBuf)
	buf.Reset()
	zw, err := gzip.NewWriterLevel(buf, gzip.DefaultCompression)
	if err != nil {
		return nil, err
	}

	_, err = zw.Write(raw)
	if err != nil {
		return nil, err
	}

	err = zw.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
