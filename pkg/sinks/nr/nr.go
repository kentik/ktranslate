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
	"strings"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/formats/nrm"
	"github.com/kentik/ktranslate/pkg/kt"
)

const (
	EnvNrApiKey = "NEW_RELIC_API_KEY"

	NR_DATA_PROVIDER = "kentik"
	NR_USER_AGENT    = "kentik-ktranslate/0.1.0"
)

type NRSink struct {
	logger.ContextL
	NRAccount   string
	NRApiKey    string
	NRUrl       string
	NRUrlEvent  string
	NRUrlMetric string

	client      *http.Client
	tr          *http.Transport
	registry    go_metrics.Registry
	metrics     *NRMetric
	format      formats.Format
	compression kt.Compression
	estimate    bool
	fmtr        *nrm.NRMFormat
	tooBig      chan int
}

type NRMetric struct {
	DeliveryErr   go_metrics.Meter
	DeliveryWin   go_metrics.Meter
	DeliveryBytes go_metrics.Meter
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
	EstimateSize = flag.Bool("nr_estimate_only", false, "If true, record size of inputs to NR but don't actually send anything")
	NrRegion     = flag.String("nr_region", "", "NR Region to use. US|EU")

	regions = map[string]map[string]string{
		"us": map[string]string{
			"events":  "https://insights-collector.newrelic.com/v1/accounts/%s/events",
			"metrics": "https://metric-api.newrelic.com/metric/v1",
		},
		"eu": map[string]string{
			"events":  "https://insights-collector.eu01.nr-data.net/v1/accounts/%s/events",
			"metrics": "https://metric-api.eu.newrelic.com/metric/v1",
		},
	}
)

func NewSink(log logger.Underlying, registry go_metrics.Registry, tooBig chan int) (*NRSink, error) {
	nr := NRSink{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "nrSink"}, log),
		NRApiKey: os.Getenv(EnvNrApiKey),
		registry: registry,
		metrics: &NRMetric{
			DeliveryErr:   go_metrics.GetOrRegisterMeter("delivery_errors_nr", registry),
			DeliveryWin:   go_metrics.GetOrRegisterMeter("delivery_wins_nr", registry),
			DeliveryBytes: go_metrics.GetOrRegisterMeter("delivery_bytes_nr", registry),
		},
		estimate: *EstimateSize,
		tooBig:   tooBig,
	}

	return &nr, nil
}

func (s *NRSink) Init(ctx context.Context, format formats.Format, compression kt.Compression, fmtr formats.Formatter) error {
	// set region if this is set.
	rval := strings.ToLower(*NrRegion)
	switch rval {
	case "": // noop
	case "eu", "us":
		*NrUrl = regions[rval]["events"]
		*NrMetricsUrl = regions[rval]["metrics"]
	default:
		return fmt.Errorf("Invalid NR region %s. Current regions: EU|US", *NrRegion)
	}

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

	// Try to upcast the formater here. This lets us send events also.
	nrmf, ok := fmtr.(*nrm.NRMFormat)
	if ok {
		s.fmtr = nrmf
		go s.checkForEvents(ctx)
	}

	s.NRUrl = fmt.Sprintf(s.NRUrl, s.NRAccount)
	s.NRUrlEvent = s.NRUrl
	s.NRUrlMetric = *NrMetricsUrl
	if s.format == formats.FORMAT_NRM {
		s.NRUrl = *NrMetricsUrl
	}
	s.Infof("Exporting to New Relic at main: %s, events: %s, metrics: %s", s.NRUrl, s.NRUrlEvent, s.NRUrlMetric)

	return nil
}

func (s *NRSink) Send(ctx context.Context, payload *kt.Output) {
	if payload.IsEvent() {
		go s.sendNR(ctx, payload, s.NRUrlEvent)
	} else if payload.IsMetric() {
		go s.sendNR(ctx, payload, s.NRUrlMetric)
	} else {
		go s.sendNR(ctx, payload, s.NRUrl)
	}
}

func (s *NRSink) Close() {}

func (s *NRSink) HttpInfo() map[string]float64 {
	return map[string]float64{
		"DeliveryErr":     s.metrics.DeliveryErr.Rate1(),
		"DeliveryWin":     s.metrics.DeliveryWin.Rate1(),
		"DeliveryBytes1":  s.metrics.DeliveryBytes.Rate1(),
		"DeliveryBytes15": s.metrics.DeliveryBytes.Rate15(),
	}
}

func (s *NRSink) sendNR(ctx context.Context, payload *kt.Output, url string) {
	s.metrics.DeliveryBytes.Mark(int64(len(payload.Body))) // Compression will effect this, but we can do our best.
	if s.estimate {                                        // here, just mark how much we would have sent.
		return
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload.Body))
	if err != nil {
		s.Errorf("Cannot create NR request: %v", err)
		return
	}

	req.Header.Set("Api-Key", s.NRApiKey)
	req.Header.Set("X-Insert-Key", s.NRApiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("NR-Data-Provider", NR_DATA_PROVIDER)
	req.Header.Set("NR-Data-Type", payload.GetDataType())
	req.Header.Set("User-Agent", NR_USER_AGENT)
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
			if resp.StatusCode == 413 {
				s.Errorf("Cannot write to NR, body too big. Adjusting max records down")
				s.metrics.DeliveryErr.Mark(1)
				s.tooBig <- len(payload.Body)
			} else if resp.StatusCode >= 400 {
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

func (s *NRSink) checkForEvents(ctx context.Context) {
	for {
		select {
		case evt := <-s.fmtr.EventChan:
			s.Debugf("Sending event")
			go s.sendNR(ctx, kt.NewOutputWithProvider(evt, kt.ProviderRouter, "event"), s.NRUrlEvent)
		case <-ctx.Done():
			s.Infof("checkForEvents Done")
			return
		}
	}
}
