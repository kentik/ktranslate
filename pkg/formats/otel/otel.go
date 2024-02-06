package otel

import (
	"context"
	"flag"
	"strconv"
	"sync"
	"time"

	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats/util"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type OtelFormat struct {
	logger.ContextL
	lastMetadata map[string]*kt.LastMetadata
	mux          sync.RWMutex
	exp          sdkmetric.Exporter
	invalids     map[string]bool
	config       *ktranslate.OtelFormatConfig
	vecs         map[string]metric.Float64ObservableGauge
	ctx          context.Context
	queues       map[string]chan OtelData
}

var (
	endpoint string
	otelm    metric.Meter
)

func init() {
	flag.StringVar(&endpoint, "otel.endpoint", "", "Send data to this endpoint or stdout")
}

func NewFormat(log logger.Underlying, cfg *ktranslate.OtelFormatConfig) (*OtelFormat, error) {
	jf := &OtelFormat{
		ContextL:     logger.NewContextLFromUnderlying(logger.SContext{S: "otel"}, log),
		lastMetadata: map[string]*kt.LastMetadata{},
		invalids:     map[string]bool{},
		ctx:          context.Background(),
		vecs:         map[string]metric.Float64ObservableGauge{},
		config:       cfg,
		queues:       map[string]chan OtelData{},
	}

	var exp sdkmetric.Exporter
	if cfg.Endpoint == "stdout" {
		metricExporter, err := stdoutmetric.New()
		if err != nil {
			return nil, err
		}
		exp = metricExporter
	} else {
		metricExporter, err := otlpmetrichttp.New(jf.ctx, otlpmetrichttp.WithEndpoint(cfg.Endpoint))
		if err != nil {
			return nil, err
		}
		exp = metricExporter
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp, sdkmetric.WithInterval(30*time.Second))))
	otel.SetMeterProvider(meterProvider)
	jf.exp = exp

	otelm = otel.Meter("ktranslate")
	jf.Infof("Running exporting to %s", cfg.Endpoint)

	return jf, nil
}

func (f *OtelFormat) To(msgs []*kt.JCHF, serBuf []byte) (*kt.Output, error) {
	res := make([]OtelData, 0, len(msgs))
	for _, m := range msgs {
		res = append(res, f.toOtelMetric(m)...)
	}

	if len(res) == 0 {
		return nil, nil
	}

	f.mux.Lock()
	defer f.mux.Unlock()

	for _, m := range res {
		if _, ok := f.vecs[m.Name]; !ok {
			cv, err := otelm.Float64ObservableGauge(m.Name)
			if err != nil {
				return nil, err
			}
			f.vecs[m.Name] = cv
			f.queues[m.Name] = make(chan OtelData, 10000)
			f.Infof("Adding %s", m.Name)
			_, err = otelm.RegisterCallback(func(_ context.Context, o metric.Observer) error {
				num := len(f.queues[m.Name])
				for i := 0; i < num; i++ {
					m := <-f.queues[m.Name]
					o.ObserveFloat64(cv, m.Value, metric.WithAttributeSet(m.GetTagValues()))
				}
				return nil
			}, f.vecs[m.Name])
			if err != nil {
				return nil, err
			}
		}
		f.queues[m.Name] <- m
	}

	return nil, nil
}

func (f *OtelFormat) From(raw *kt.Output) ([]map[string]interface{}, error) {
	values := make([]map[string]interface{}, 0)
	return values, nil
}

func (f *OtelFormat) Rollup(rolls []rollup.Rollup) (*kt.Output, error) {
	return nil, nil
}

func (f *OtelFormat) toOtelMetric(in *kt.JCHF) []OtelData {
	switch in.EventType {
	case kt.KENTIK_EVENT_TYPE:
		return f.fromKflow(in)
	case kt.KENTIK_EVENT_SYNTH:
		return f.fromKSynth(in)
	case kt.KENTIK_EVENT_SYNTH_GEST:
		return f.fromKSyngest(in)
	default:
		f.mux.Lock()
		defer f.mux.Unlock()
		if !f.invalids[in.EventType] {
			f.Warnf("Invalid EventType: %s", in.EventType)
			f.invalids[in.EventType] = true
		}
	}

	return nil
}

var (
	synthWLAttr = map[string]bool{
		"agent_id":               true,
		"agent_name":             true,
		"dst_addr":               true,
		"dst_cdn_int":            true,
		"dst_geo":                true,
		"provider":               true,
		"src_addr":               true,
		"src_cdn_int":            true,
		"src_as_name":            true,
		"src_geo":                true,
		"test_id":                true,
		"test_name":              true,
		"test_type":              true,
		"test_url":               true,
		"src_host":               true,
		"dst_host":               true,
		"src_cloud_region":       true,
		"src_cloud_provider":     true,
		"src_site":               true,
		"dst_cloud_region":       true,
		"dst_cloud_provider":     true,
		"dst_site":               true,
		"statusMessage":          true,
		"statusEncoding":         true,
		"https_validity":         true,
		"https_expiry_timestamp": true,
		"dest_ip":                true,
	}

	synthAttrKeys = []string{
		"statusMessage",
		"statusEncoding",
		"https_validity",
		"https_expiry_timestamp",
	}
)

func (f *OtelFormat) fromKSyngest(in *kt.JCHF) []OtelData {
	metrics := util.GetSyngestMetricNameSet()
	attr := map[string]interface{}{}
	f.mux.RLock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName], false)
	f.mux.RUnlock()
	ms := make([]OtelData, 0, len(metrics))

	for k, v := range attr { // White list only a few attributes here.
		if !synthWLAttr[k] {
			delete(attr, k)
		}
		if k == "test_id" { // Force this to be a string.
			if vi, ok := v.(int); ok {
				attr[k] = strconv.Itoa(vi)
			}
		}
	}

	for m, name := range metrics {
		if in.CustomInt[m] > 0 {
			ms = append(ms, OtelData{
				Name:  "kentik.syngest." + name.Name,
				Value: float64(in.CustomInt[m]),
				Tags:  attr,
			})
		}
	}

	return ms
}

func (f *OtelFormat) fromKSynth(in *kt.JCHF) []OtelData {
	if in.CustomInt["result_type"] <= 1 {
		return nil // Don't worry about timeouts and errors for now.
	}

	rawStr := in.CustomStr["error_cause/trace_route"] // Pull this out early.
	metrics := util.GetSynMetricNameSet(in.CustomInt["result_type"])
	attr := map[string]interface{}{}
	f.mux.RLock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName], false)
	f.mux.RUnlock()
	ms := make([]OtelData, 0, len(metrics))

	// If there's str00 data, try to unserialize and pass in useful bits.
	if rawStr != "" {
		strData := []interface{}{}
		if err := json.Unmarshal([]byte(rawStr), &strData); err == nil {
			if len(strData) > 0 {
				switch sd := strData[0].(type) {
				case map[string]interface{}:
					for _, key := range synthAttrKeys {
						if val, ok := sd[key]; ok {
							attr[key] = val
						}
					}
				}
			}
		}
	}

	for k, v := range attr { // White list only a few attributes here.
		if !synthWLAttr[k] {
			delete(attr, k)
		}
		if k == "test_id" { // Force this to be a string.
			if vi, ok := v.(int); ok {
				attr[k] = strconv.Itoa(vi)
			}
		}
	}

	for m, name := range metrics {
		switch name.Name {
		case "avg_rtt", "jit_rtt", "time", "code", "port", "status", "ttlb", "size", "trx_time", "validation", "lost", "sent":
			ms = append(ms, OtelData{
				Name:  "kentik.synth." + name.Name,
				Value: float64(in.CustomInt[m]),
				Tags:  attr,
			})
		}
	}

	return ms
}

func (f *OtelFormat) fromKflow(in *kt.JCHF) []OtelData {
	// Map the basic strings into here.
	attr := map[string]interface{}{}
	metrics := map[string]kt.MetricInfo{"in_bytes": kt.MetricInfo{}, "out_bytes": kt.MetricInfo{}, "in_pkts": kt.MetricInfo{}, "out_pkts": kt.MetricInfo{}, "latency_ms": kt.MetricInfo{}}
	f.mux.RLock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName], false)
	f.mux.RUnlock()
	ms := map[string]int64{}
	for m, _ := range metrics {
		switch m {
		case "in_bytes":
			ms[m] = int64(in.InBytes * uint64(in.SampleRate))
		case "out_bytes":
			ms[m] = int64(in.OutBytes * uint64(in.SampleRate))
		case "in_pkts":
			ms[m] = int64(in.InPkts * uint64(in.SampleRate))
		case "out_pkts":
			ms[m] = int64(in.OutPkts * uint64(in.SampleRate))
		case "latency_ms":
			ms[m] = int64(in.CustomInt["appl_latency_ms"])
		}
	}

	res := []OtelData{}
	for k, v := range ms {
		if v == 0 { // Drop zero valued metrics here.
			continue
		}
		res = append(res, OtelData{
			Name:  "kentik.flow." + k,
			Value: float64(v),
			Tags:  attr,
		})
	}

	return res
}

type OtelData struct {
	Name  string
	Value float64
	Tags  map[string]interface{}
}

func (d *OtelData) GetTagValues() attribute.Set {
	res := make([]attribute.KeyValue, 0, len(d.Tags))
	for k, v := range d.Tags {
		switch t := v.(type) {
		case string:
			res = append(res, attribute.String(k, t))
		case int64:
			res = append(res, attribute.Int64(k, t))
		case float64:
			res = append(res, attribute.Float64(k, t))
		default:
		}
	}
	s, _ := attribute.NewSetWithFiltered(res, func(kv attribute.KeyValue) bool { return true })
	return s
}
