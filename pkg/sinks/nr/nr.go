package nr

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

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
	NRUrlLog    string

	client      *http.Client
	tr          *http.Transport
	registry    go_metrics.Registry
	metrics     *NRMetric
	format      formats.Format
	compression kt.Compression
	estimate    bool
	checkJson   bool
	fmtr        *nrm.NRMFormat
	tooBig      chan int
	logTee      chan string
}

type NRMetric struct {
	DeliveryErr   go_metrics.Meter
	DeliveryWin   go_metrics.Meter
	DeliveryBytes go_metrics.Meter
	DeliveryLogs  go_metrics.Meter
}

type NRResponce struct {
	Success   bool   `json:"success"`
	Uuid      string `json:"uuid"`
	RequestId string `json:"requestId"`
}

var (
	NrAccount    = flag.String("nr_account_id", kt.LookupEnvString("NR_ACCOUNT_ID", ""), "If set, sends flow to New Relic")
	NrUrl        = flag.String("nr_url", "https://insights-collector.newrelic.com/v1/accounts/%s/events", "URL to use to send into NR")
	NrMetricsUrl = flag.String("nr_metrics_url", "https://metric-api.newrelic.com/metric/v1", "URL to use to send into NR Metrics API")
	EstimateSize = flag.Bool("nr_estimate_only", false, "If true, record size of inputs to NR but don't actually send anything")
	NrRegion     = flag.String("nr_region", kt.LookupEnvString("NR_REGION", ""), "NR Region to use. US|EU")
	NrCheckJson  = flag.Bool("nr_check_json", false, "Verify body is valid json before sending on")

	regions = map[string]map[string]string{
		"us": map[string]string{
			"events":  "https://insights-collector.newrelic.com/v1/accounts/%s/events",
			"metrics": "https://metric-api.newrelic.com/metric/v1",
			"logs":    "https://log-api.newrelic.com/log/v1",
		},
		"eu": map[string]string{
			"events":  "https://insights-collector.eu01.nr-data.net/v1/accounts/%s/events",
			"metrics": "https://metric-api.eu.newrelic.com/metric/v1",
			"logs":    "https://log-api.eu.newrelic.com/log/v1",
		},
	}
)

func NewSink(log logger.Underlying, registry go_metrics.Registry, tooBig chan int, logTee chan string) (*NRSink, error) {
	nr := NRSink{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "nrSink"}, log),
		NRApiKey: os.Getenv(EnvNrApiKey),
		registry: registry,
		metrics: &NRMetric{
			DeliveryErr:   go_metrics.GetOrRegisterMeter("delivery_errors_nr", registry),
			DeliveryWin:   go_metrics.GetOrRegisterMeter("delivery_wins_nr", registry),
			DeliveryBytes: go_metrics.GetOrRegisterMeter("delivery_bytes_nr", registry),
			DeliveryLogs:  go_metrics.GetOrRegisterMeter("delivery_logs_nr", registry),
		},
		estimate:  *EstimateSize,
		checkJson: *NrCheckJson,
		tooBig:    tooBig,
		logTee:    logTee,
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
		s.NRUrlLog = regions[rval]["logs"]
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

	if strings.Contains(s.NRUrl, "%s") {
		s.NRUrl = fmt.Sprintf(s.NRUrl, s.NRAccount)
	}
	s.NRUrlEvent = s.NRUrl
	s.NRUrlMetric = *NrMetricsUrl
	if s.format == formats.FORMAT_NRM {
		s.NRUrl = *NrMetricsUrl
	}

	if s.NRUrlLog == "" { // TODO -- better default?
		s.NRUrlLog = regions["us"]["logs"]
	}

	// Send logs on to NR if this is set.
	if s.logTee != nil {
		go s.watchLogs(ctx)
	}

	s.Infof("Exporting to New Relic at main: %s, events: %s, metrics: %s, logs %s", s.NRUrl, s.NRUrlEvent, s.NRUrlMetric, s.NRUrlLog)

	return nil
}

func (s *NRSink) Send(ctx context.Context, payload *kt.Output) {
	if s.checkJson {
		if err := s.doCheckJson(payload); err != nil {
			s.Errorf("Invalid payload! Not sending -- len: %d, err: %v", len(payload.Body), err)
			return
		}
	}

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
		"DeliveryLogs":    s.metrics.DeliveryLogs.Rate1(),
	}
}

func (s *NRSink) sendNR(ctx context.Context, payload *kt.Output, url string) {
	var cbErr error = nil
	if payload.CB != nil { // Let anyone who asked know that this has been sent
		defer payload.CB(cbErr)
	}

	s.metrics.DeliveryBytes.Mark(int64(len(payload.Body))) // Compression will effect this, but we can do our best.
	if s.estimate {                                        // here, just mark how much we would have sent.
		return
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload.Body))
	if err != nil {
		s.Errorf("Cannot create NR request: %v", err)
		cbErr = err
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
		cbErr = err
		s.client = &http.Client{Transport: s.tr}
	} else {
		defer resp.Body.Close()
		bdy, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			s.Errorf("Cannot get resp body from NR: %v", err)
			cbErr = err
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
					cbErr = err
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

// Forwards any logs recieved to the NR log API.
func (s *NRSink) watchLogs(ctx context.Context) {
	s.Infof("Watching for logs")
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
			s.Infof("Watching for logs done")
			return
		}
	}
}

// Quick types to pass in log lines.
type logSet struct {
	Common *common `json:"common"`
	Logs   []log   `json:"logs"`
}

type common struct {
	Attributes map[string]string `json:"attributes"`
}

type log struct {
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
}

func (s *NRSink) sendLogBatch(ctx context.Context, logs []string) {
	ts := time.Now().Unix()
	ls := logSet{
		Common: &common{
			Attributes: map[string]string{
				"instrumentation.provider": kt.InstProvider,
				"collector.name":           kt.CollectorName,
			},
		},
		Logs: make([]log, len(logs)),
	}
	for i, l := range logs {
		ls.Logs[i] = log{
			Timestamp: ts,
			Message:   l,
		}
	}

	target, err := json.Marshal([]logSet{ls}) // Has to be an array here, no idea why.
	if err != nil {
		s.Errorf("Cannot marshal log set: %v", err)
		return
	}

	if s.compression != kt.CompressionGzip {
		s.sendNR(ctx, kt.NewOutput(target), s.NRUrlLog)
	}

	serBuf := []byte{}
	buf := bytes.NewBuffer(serBuf)
	buf.Reset()
	zw, err := gzip.NewWriterLevel(buf, gzip.DefaultCompression)
	if err != nil {
		s.Errorf("Cannot gzip log set: %v", err)
		return
	}

	_, err = zw.Write(target)
	if err != nil {
		s.Errorf("Cannot write log set: %v", err)
		return
	}

	err = zw.Close()
	if err != nil {
		s.Errorf("Cannot close log set: %v", err)
		return
	}

	s.sendNR(ctx, kt.NewOutput(buf.Bytes()), s.NRUrlLog)
}

func (s *NRSink) doCheckJson(payload *kt.Output) error {
	var base interface{}
	if s.compression != kt.CompressionGzip { // plain path
		if err := json.Unmarshal(payload.Body, &base); err != nil {
			return err
		}
		return nil
	}

	// Compression check here.
	r, err := gzip.NewReader(bytes.NewBuffer(payload.Body))
	if err != nil {
		return err
	}
	if err := json.NewDecoder(r).Decode(&base); err != nil {
		return err
	}

	return nil
}
