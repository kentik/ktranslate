package otel

import (
	"context"
	"strconv"
	"sync"

	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats/util"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type OtelFormat struct {
	logger.ContextL
	compression  kt.Compression
	doGz         bool
	lastMetadata map[string]*kt.LastMetadata
	mux          sync.RWMutex
	exp          *otlpmetrichttp.Exporter
	vecTags      tagVec
	seen         map[string]int
	invalids     map[string]bool
	config       *ktranslate.OtelFormatConfig
	vecs         map[string]metric.Float64Counter
	ctx          context.Context
}

var (
	otelm = otel.Meter("base")
)

func NewFormat(log logger.Underlying, compression kt.Compression, cfg *ktranslate.OtelFormatConfig) (*OtelFormat, error) {
	jf := &OtelFormat{
		compression:  compression,
		ContextL:     logger.NewContextLFromUnderlying(logger.SContext{S: "nrmFormat"}, log),
		lastMetadata: map[string]*kt.LastMetadata{},
		ctx:          context.Background(),
	}

	exp, err := otlpmetrichttp.New(jf.ctx, otlpmetrichttp.WithEndpoint(cfg.Endpoint))
	if err != nil {
		return nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp)))
	otel.SetMeterProvider(meterProvider)
	jf.exp = exp

	otelm = otel.Meter("ktranslate")

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
		if f.seen[m.Name] < f.config.FlowsNeeded {
			m.AddTagLabels(f.vecTags)
			f.seen[m.Name]++
			if f.seen[m.Name] == f.config.FlowsNeeded {
				f.Infof("Seen enough %s!", m.Name)
			} else {
				f.Infof("Seen %s -> %d", m.Name, f.seen[m.Name])
			}
			continue
		}

		labels := m.GetTagValues(f.vecTags)
		if _, ok := f.vecs[m.Name]; !ok {
			cv, err := otelm.Float64Counter(m.Name)
			if err != nil {
				return nil, err
			}
			f.vecs[m.Name] = cv
			f.Infof("Adding %s %v", m.Name, labels)
		}
		f.vecs[m.Name].Add(f.ctx, m.Value, metric.WithAttributes(labels...))
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
				Name:  "kentik:syngest:" + name.Name,
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
				Name:  "kentik:synth:" + name.Name,
				Value: float64(in.CustomInt[m]),
				Tags:  attr,
			})
		}
	}

	return ms
}

type OtelData struct {
	Name  string
	Value float64
	Tags  map[string]interface{}
}

func (d *OtelData) AddTagLabels(vecTags tagVec) {
	if _, ok := vecTags[d.Name]; !ok {
		vecTags[d.Name] = make([]attribute.KeyValue, 0)
	}
	for k, v := range d.Tags {
		found := false
		for _, kk := range vecTags[d.Name] {
			if string(kk.Key) == k {
				found = true
				break // Key already exists.
			}
		}
		// If here, new key.
		if !found {
			switch t := v.(type) {
			case string:
				vecTags[d.Name] = append(vecTags[d.Name], attribute.String(k, t))
			case int64:
				vecTags[d.Name] = append(vecTags[d.Name], attribute.Int64(k, t))
			case float64:
				vecTags[d.Name] = append(vecTags[d.Name], attribute.Float64(k, t))
			default:
			}
		}
	}
}

func (d *OtelData) GetTagValues(vecTags tagVec) []attribute.KeyValue {
	return vecTags[d.Name]
}

type tagVec map[string][]attribute.KeyValue
