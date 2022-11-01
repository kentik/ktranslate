package nr

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats"
	"github.com/kentik/ktranslate/pkg/formats/nrm"
	"github.com/kentik/ktranslate/pkg/kt"
)

const (
	EnvNrApiKey = "NEW_RELIC_API_KEY"

	NR_DATA_PROVIDER = "kentik"
	NR_USER_AGENT    = "kentik-ktranslate/0.1.0"

	REGION_EU         = "eu"
	REGION_US         = "us"
	REGION_GOV        = "gov"
	REGION_US_STAGING = "us_stage"
)

var (
	nrAccount    string
	estimateSize bool
	nrRegion     string
	nrCheckJson  bool

	NrUrl        = "https://insights-collector.newrelic.com/v1/accounts/%s/events"
	NrMetricsUrl = "https://metric-api.newrelic.com/metric/v1"
	NrLogUrl     = "https://log-api.newrelic.com/log/v1"
	regions      = map[string]map[string]string{
		REGION_US: map[string]string{
			"events":  "https://insights-collector.newrelic.com/v1/accounts/%s/events",
			"metrics": "https://metric-api.newrelic.com/metric/v1",
			"logs":    "https://log-api.newrelic.com/log/v1",
		},
		REGION_EU: map[string]string{
			"events":  "https://insights-collector.eu01.nr-data.net/v1/accounts/%s/events",
			"metrics": "https://metric-api.eu.newrelic.com/metric/v1",
			"logs":    "https://log-api.eu.newrelic.com/log/v1",
		},
		REGION_GOV: map[string]string{
			"events":  "https://gov-insights-collector.newrelic.com/v1/accounts/%s/events",
			"metrics": "https://gov-metric-api.newrelic.com/metric/v1",
			"logs":    "https://gov-log-api.newrelic.com/log/v1",
		},
		REGION_US_STAGING: map[string]string{
			"events":  "https://staging-insights-collector.newrelic.com/v1/accounts/%s/events",
			"metrics": "https://staging-metric-api.newrelic.com/metric/v1",
			"logs":    "https://staging-log-api.newrelic.com/log/v1",
		},
	}
)

func init() {
	flag.StringVar(&nrAccount, "nr_account_id", kt.LookupEnvString("NR_ACCOUNT_ID", ""), "If set, sends flow to New Relic")
	flag.BoolVar(&estimateSize, "nr_estimate_only", false, "If true, record size of inputs to NR but don't actually send anything")
	flag.StringVar(&nrRegion, "nr_region", kt.LookupEnvString("NR_REGION", ""), "NR Region to use. US|EU")
	flag.BoolVar(&nrCheckJson, "nr_check_json", false, "Verify body is valid json before sending on")
}

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
	config      *ktranslate.NewRelicSinkConfig
}

type NRMetric struct {
	DeliveryErr     go_metrics.Meter
	DeliveryWin     go_metrics.Meter
	DeliveryMetrics go_metrics.Meter
	DeliveryLogs    go_metrics.Meter
}

type NRResponce struct {
	Success   bool   `json:"success"`
	Uuid      string `json:"uuid"`
	RequestId string `json:"requestId"`
}

func NewSink(log logger.Underlying, registry go_metrics.Registry, tooBig chan int, logTee chan string, cfg *ktranslate.NewRelicSinkConfig) (*NRSink, error) {
	nr := NRSink{
		ContextL:  logger.NewContextLFromUnderlying(logger.SContext{S: "nrSink"}, log),
		NRApiKey:  os.Getenv(EnvNrApiKey),
		NRAccount: cfg.Account,
		registry:  registry,
		metrics: &NRMetric{
			DeliveryErr:     go_metrics.GetOrRegisterMeter("delivery_errors_nr", registry),
			DeliveryWin:     go_metrics.GetOrRegisterMeter("delivery_wins_nr", registry),
			DeliveryMetrics: go_metrics.GetOrRegisterMeter("delivery_metrics_nr", registry),
			DeliveryLogs:    go_metrics.GetOrRegisterMeter("delivery_logs_nr", registry),
		},
		estimate:  cfg.EstimateOnly,
		checkJson: cfg.ValidateJSON,
		tooBig:    tooBig,
		logTee:    logTee,
		config:    cfg,
	}

	return &nr, nil
}

func (s *NRSink) Init(ctx context.Context, format formats.Format, compression kt.Compression, fmtr formats.Formatter) error {
	// set region if this is set.
	rval := strings.ToLower(s.config.Region)
	switch rval {
	case "": // noop
	case REGION_US, REGION_EU, REGION_GOV, REGION_US_STAGING:
		NrUrl = regions[rval]["events"]
		NrMetricsUrl = regions[rval]["metrics"]
		s.NRUrlLog = regions[rval]["logs"]
	default:
		return fmt.Errorf("You used an unsupported New Relic One region: %s. The possible values are EU, US, GOV and US_STAGE.", s.config.Region)
	}

	s.NRUrl = NrUrl
	s.format = format
	s.compression = compression

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

	if s.NRAccount == "" || s.NRApiKey == "" {
		return fmt.Errorf("You need to set up your New Relic One account ID in the -nr_account_id variable and your API key in the NEW_RELIC_API_KEY environment variable.")
	}
	if s.format != formats.FORMAT_NR && s.format != formats.FORMAT_JSON_FLAT && s.format != formats.FORMAT_NRM {
		return fmt.Errorf("You used the %s unsupported format. Use flat_json, new_relic, new_relic_metric.", s.format)
	}
	if s.compression != kt.CompressionGzip && s.compression != kt.CompressionNone {
		return fmt.Errorf("You used the %s unsupported compression format. Use gzip or no compression at all.", s.compression)
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
	s.NRUrlMetric = NrMetricsUrl
	if s.format == formats.FORMAT_NRM {
		s.NRUrl = NrMetricsUrl
	}

	if s.NRUrlLog == "" {
		s.NRUrlLog = NrLogUrl
	}

	// Send logs on to NR if this is set.
	if s.logTee != nil {
		go s.watchLogs(ctx)
	}

	s.Infof("Exporting to New Relic at main: %s, events: %s, metrics: %s, logs %s", s.NRUrl, s.NRUrlEvent, s.NRUrlMetric, s.NRUrlLog)

	err := s.test(ctx)
	if err != nil {
		return err
	}
	s.Infof("New Relic sink connection confirmed good.")

	return nil
}

func (s *NRSink) Send(ctx context.Context, payload *kt.Output) {
	if s.checkJson {
		if err := s.doCheckJson(payload); err != nil {
			s.Errorf("You used an unsupported payload: %d. %v.", len(payload.Body), err)
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
		"DeliveryErr":       s.metrics.DeliveryErr.Rate1(),
		"DeliveryWin":       s.metrics.DeliveryWin.Rate1(),
		"DeliveryMetrics1":  s.metrics.DeliveryMetrics.Rate1(),
		"DeliveryMetrics15": s.metrics.DeliveryMetrics.Rate15(),
		"DeliveryLogs":      s.metrics.DeliveryLogs.Rate1(),
	}
}

func (s *NRSink) test(ctx context.Context) error {
	urls := []string{s.NRUrlEvent, s.NRUrlMetric, s.NRUrlLog}
	payload := kt.NewOutput([]byte("{}"))
	errChan := make(chan error)
	cb := func(err error) {
		go func() { // We're still in the same thread as below so need to get a new one.
			errChan <- err
		}()
	}
	payload.CB = cb // sendNR uses a callback system instead of returning directly.
	for _, url := range urls {
		s.sendNR(ctx, payload, url)
		err := <-errChan
		if err != nil {
			return fmt.Errorf("Error testing %s: %v", url, err)
		}
	}

	return nil
}

func (s *NRSink) sendNR(ctx context.Context, payload *kt.Output, url string) {
	var cbErr error = nil
	if payload.CB != nil { // Let anyone who asked know that this has been sent
		defer func() {
			payload.CB(cbErr) // This needs to get wrapped in its own function to get cbErr right.
		}()
	}

	s.metrics.DeliveryMetrics.Mark(1) // Compression will effect this, but we can do our best.
	if s.estimate {                   // here, just mark how much we would have sent.
		return
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload.Body))
	if err != nil {
		s.Errorf("There was an error when communicating to New Relic One: %v.", err)
		cbErr = err
		return
	}

	req.Header.Set("Api-Key", s.NRApiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("NR-Data-Provider", NR_DATA_PROVIDER)
	req.Header.Set("NR-Data-Type", payload.GetDataType())
	req.Header.Set("User-Agent", NR_USER_AGENT)
	if s.compression == kt.CompressionGzip {
		req.Header.Set("Content-Encoding", "GZIP")
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.Errorf("There was an error when creating a new client in New Relic One: %v.", err)
		cbErr = err
		s.client = &http.Client{Transport: s.tr}
	} else {
		defer resp.Body.Close()
		bdy, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			s.Errorf("There was an error when communicating to New Relic One: %v.", err)
			cbErr = err
			s.metrics.DeliveryErr.Mark(1)
		} else {
			if resp.StatusCode == 413 {
				s.Errorf("There was an error when communicating to New Relic One. The message was too big.")
				s.metrics.DeliveryErr.Mark(1)
				s.tooBig <- len(payload.Body)
			} else if resp.StatusCode >= 400 {
				s.Errorf("There was an error when communicating to New Relic One: %v.", resp.StatusCode)
				cbErr = fmt.Errorf("There was an error when communicating to New Relic One: %v.", resp.StatusCode)
				s.metrics.DeliveryErr.Mark(1)
			} else {
				var nr NRResponce
				err = json.Unmarshal(bdy, &nr)
				if err != nil {
					s.Errorf("There was an error when parsing the response from New Relic One: %v.", err)
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
	hasSyslog := false
	for i, l := range logs {
		ls.Logs[i] = log{
			Timestamp: ts,
			Message:   l,
		}
		if !hasSyslog && strings.Contains(l, kt.PluginSyslog) {
			hasSyslog = true
		}
	}
	if !hasSyslog {
		ls.Common.Attributes["plugin.type"] = kt.PluginHealth
		ls.Common.Attributes["logtype"] = "ktranslate-health"
	}

	target, err := json.Marshal([]logSet{ls}) // Has to be an array here, no idea why.
	if err != nil {
		s.Errorf("There was an error with logs: %v.", err)
		return
	}

	if s.compression != kt.CompressionGzip {
		s.sendNR(ctx, kt.NewOutput(target), s.NRUrlLog)
		return
	}

	serBuf := []byte{}
	buf := bytes.NewBuffer(serBuf)
	buf.Reset()
	zw, err := gzip.NewWriterLevel(buf, gzip.DefaultCompression)
	if err != nil {
		s.Errorf("There was an error when compressing logs: %v.", err)
		return
	}

	_, err = zw.Write(target)
	if err != nil {
		s.Errorf("There was an error when writing logs: %v.", err)
		return
	}

	err = zw.Close()
	if err != nil {
		s.Errorf("There was an error when closing logs: %v.", err)
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
