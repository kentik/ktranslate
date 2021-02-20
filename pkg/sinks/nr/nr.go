package nr

import (
	"bytes"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/kt"
)

const (
	EnvNrApiKey = "NEW_RELIC_API_KEY"
)

type NRSink struct {
	logger.ContextL
	NRAccount string
	NRApiKey  string
	NRUrl     string

	client      *http.Client
	tr          *http.Transport
	registry    go_metrics.Registry
	metrics     *NRMetric
	format      formats.Format
	compression kt.Compression
}

type NRMetric struct {
	DeliveryErr go_metrics.Meter
	DeliveryWin go_metrics.Meter
}

type NRResponce struct {
	Success   bool   `json:"success"`
	Uuid      string `json:"uuid"`
	RequestId string `json:"requestId"`
}

var (
	NrAccount    = flag.String("nr_account_id", "", "If set, sends flow to New Relic")
	NrUrl        = flag.String("nr_url", "https://insights-collector.newrelic.com/v1/accounts/%s/events", "URL to use to send into NR")
	NrMetricsUrl = flag.String("nr_metrics_url", "https://metric-api.newrelic.com/metric/v1", "URL to use to send into NR Metrics API")
)

func NewSink(log logger.Underlying, registry go_metrics.Registry) (*NRSink, error) {
	nr := NRSink{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "nrSink"}, log),
		NRApiKey: os.Getenv(EnvNrApiKey),
		registry: registry,
		metrics: &NRMetric{
			DeliveryErr: go_metrics.GetOrRegisterMeter("delivery_errors_nr", registry),
			DeliveryWin: go_metrics.GetOrRegisterMeter("delivery_wins_nr", registry),
		},
	}

	return &nr, nil
}

func (s *NRSink) Init(ctx context.Context, format formats.Format, compression kt.Compression) error {
	s.NRAccount = *NrAccount
	s.NRUrl = *NrUrl
	s.format = format
	s.compression = compression

	s.tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	s.client = &http.Client{Transport: s.tr}

	if s.NRAccount == "" || s.NRApiKey == "" {
		return fmt.Errorf("New Relic requires -nr_account_id flag and NEW_RELIC_API_KEY env var to be set")
	}
	if s.format != formats.FORMAT_NR && s.format != formats.FORMAT_JSON && s.format != formats.FORMAT_NRM {
		return fmt.Errorf("New Relic only supports new_relic and json formats, not %s", s.format)
	}
	if s.compression != kt.CompressionGzip && s.compression != kt.CompressionNone {
		return fmt.Errorf("New Relic only supports gzip and none compression, not %s", s.compression)
	}

	s.NRUrl = fmt.Sprintf(s.NRUrl, s.NRAccount)
	if s.format == formats.FORMAT_NRM {
		s.NRUrl = *NrMetricsUrl
	}
	s.Infof("Exporting to New Relic at %s", s.NRUrl)

	return nil
}

func (s *NRSink) Send(ctx context.Context, payload []byte) {
	go s.sendNR(ctx, payload)
}

func (s *NRSink) Close() {}

func (s *NRSink) HttpInfo() map[string]float64 {
	return map[string]float64{
		"DeliveryErr": s.metrics.DeliveryErr.Rate1(),
		"DeliveryWin": s.metrics.DeliveryWin.Rate1(),
	}
}

func (s *NRSink) sendNR(ctx context.Context, payload []byte) {
	req, err := http.NewRequestWithContext(ctx, "POST", s.NRUrl, bytes.NewBuffer(payload))
	if err != nil {
		s.Errorf("Cannot create NR request: %v", err)
		return
	}

	req.Header.Set("Api-Key", s.NRApiKey)
	req.Header.Set("X-Insert-Key", s.NRApiKey)
	req.Header.Set("Content-Type", "application/json")
	if s.compression == kt.CompressionGzip {
		req.Header.Set("Content-Encoding", "GZIP")
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.Errorf("Cannot write to NR: %v, creating new client", err)
		s.client = &http.Client{Transport: s.tr}
	} else {
		defer resp.Body.Close()
		bdy, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			s.Errorf("Cannot get resp body from NR: %v", err)
			s.metrics.DeliveryErr.Mark(1)
		} else {
			if resp.StatusCode >= 400 {
				s.Errorf("Cannot write to NR, status code %d", resp.StatusCode)
				s.metrics.DeliveryErr.Mark(1)
			} else {
				var nr NRResponce
				err = json.Unmarshal(bdy, &nr)
				if err != nil {
					s.Errorf("Cannot parse resp from NR: %v", err)
					s.metrics.DeliveryErr.Mark(1)
				} else {
					s.Debugf("NR Success: %v UUID: %s, RID: %s", nr.Success, nr.Uuid, nr.RequestId)
					s.metrics.DeliveryWin.Mark(1)
				}
			}
		}
	}
}
