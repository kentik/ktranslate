package otel

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats/util"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
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
	inputs       map[string][]OtelData
}

var (
	endpoint string
	protocol string
	otelm    metric.Meter
)

func init() {
	flag.StringVar(&endpoint, "otel.endpoint", "", "Send data to this endpoint.")
	flag.StringVar(&protocol, "otel.protocol", "stdout", "Send data using this protocol. (grpc,http,https,stdout)")
}

func NewFormat(log logger.Underlying, cfg *ktranslate.OtelFormatConfig) (*OtelFormat, error) {
	jf := &OtelFormat{
		ContextL:     logger.NewContextLFromUnderlying(logger.SContext{S: "otel"}, log),
		lastMetadata: map[string]*kt.LastMetadata{},
		invalids:     map[string]bool{},
		ctx:          context.Background(),
		vecs:         map[string]metric.Float64ObservableGauge{},
		config:       cfg,
		inputs:       map[string][]OtelData{},
	}

	var exp sdkmetric.Exporter
	switch cfg.Protocol {
	case "stdout":
		metricExporter, err := stdoutmetric.New()
		if err != nil {
			return nil, err
		}
		exp = metricExporter
	case "http", "https": // Use OTEL_EXPORTER_OTLP_COMPRESSION env var to turn on gzip compression.
		if cfg.Endpoint == "" {
			return nil, fmt.Errorf("-otel.endpoint required for http(s) exports.")
		}
		metricExporter, err := otlpmetrichttp.New(jf.ctx, otlpmetrichttp.WithEndpoint(cfg.Endpoint))
		if err != nil {
			return nil, err
		}
		exp = metricExporter
	case "grpc": // Same, use OTEL_EXPORTER_OTLP_COMPRESSION env var to turn on gzip compression.
		if cfg.Endpoint == "" {
			return nil, fmt.Errorf("-otel.endpoint required for grpc exports.")
		}
		metricExporter, err := otlpmetricgrpc.New(jf.ctx, otlpmetricgrpc.WithEndpoint(cfg.Endpoint))
		if err != nil {
			return nil, err
		}
		exp = metricExporter
	default:
		return nil, fmt.Errorf("Invalid otel.protocol value. Currently supported: grpc,http,https,stdout")
	}

	meterProvider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp)))
	otel.SetMeterProvider(meterProvider)
	jf.exp = exp

	otelm = otel.Meter("ktranslate")
	jf.Infof("Running exporting via %s to %s", cfg.Protocol, cfg.Endpoint)

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
			lm := m
			cv, err := otelm.Float64ObservableGauge(
				lm.Name,
				metric.WithFloat64Callback(func(_ context.Context, o metric.Float64Observer) error {
					for _, mm := range f.getLatestInputs(lm.Name) {
						o.Observe(mm.Value, metric.WithAttributeSet(mm.GetTagValues()))
					}
					return nil
				}),
			)
			if err != nil {
				return nil, err
			}
			f.vecs[lm.Name] = cv
			f.inputs[lm.Name] = make([]OtelData, 0)
		}
		// Save this for later, for the next time the async callback is run.
		f.inputs[m.Name] = append(f.inputs[m.Name], m)
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

func (f *OtelFormat) getLatestInputs(name string) []OtelData {
	f.mux.Lock()
	defer f.mux.Unlock()

	nv := f.inputs[name]
	f.inputs[name] = make([]OtelData, 0)
	return nv
}

func (f *OtelFormat) toOtelMetric(in *kt.JCHF) []OtelData {
	switch in.EventType {
	case kt.KENTIK_EVENT_TYPE:
		return f.fromKflow(in)
	case kt.KENTIK_EVENT_SYNTH:
		return f.fromKSynth(in)
	case kt.KENTIK_EVENT_SYNTH_GEST:
		return f.fromKSyngest(in)
	case kt.KENTIK_EVENT_SNMP_DEV_METRIC:
		return f.fromSnmpDeviceMetric(in)
	case kt.KENTIK_EVENT_SNMP_INT_METRIC:
		return f.fromSnmpInterfaceMetric(in)
	case kt.KENTIK_EVENT_SNMP_METADATA:
		return f.fromSnmpMetadata(in)
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

func (f *OtelFormat) fromSnmpDeviceMetric(in *kt.JCHF) []OtelData {
	metrics := in.CustomMetrics
	attr := map[string]interface{}{}
	f.mux.RLock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName], false)
	f.mux.RUnlock()

	ms := make([]OtelData, 0, len(metrics))
	for m, name := range metrics {
		if m == "" {
			f.Errorf("Missing metric name, skipping %v", attr)
			continue
		}
		if _, ok := in.CustomBigInt[m]; ok {
			attrNew := util.CopyAttrForSnmp(attr, m, name, f.lastMetadata[in.DeviceName], true, false)
			if util.DropOnFilter(attrNew, f.lastMetadata[in.DeviceName], false) {
				continue // This Metric isn't in the white list so lets drop it.
			}

			mtype := name.GetType()
			if name.Format == kt.FloatMS {
				ms = append(ms, OtelData{
					Name:  "kentik." + mtype + "." + m,
					Value: float64(float64(in.CustomBigInt[m]) / 1000),
					Tags:  attrNew,
				})
			} else {
				ms = append(ms, OtelData{
					Name:  "kentik." + mtype + "." + m,
					Value: float64(in.CustomBigInt[m]),
					Tags:  attrNew,
				})
			}
		}
	}

	return ms
}

func (f *OtelFormat) fromSnmpInterfaceMetric(in *kt.JCHF) []OtelData {
	metrics := in.CustomMetrics
	attr := map[string]interface{}{}
	f.mux.RLock()
	defer f.mux.RUnlock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName], false)
	if f.lastMetadata[in.DeviceName] == nil {
		f.Debugf("Missing interface metadata for %s", in.DeviceName)
	}

	ms := make([]OtelData, 0, len(metrics))
	profileName := "snmp"
	for m, name := range metrics {
		if m == "" {
			f.Errorf("Missing metric name, skipping %v", attr)
			continue
		}
		if strings.HasSuffix(m, "_counter") {
			// Skip these counters which are not needed.
			continue
		}
		profileName = name.Profile
		if _, ok := in.CustomBigInt[m]; ok {
			attrNew := util.CopyAttrForSnmp(attr, m, name, f.lastMetadata[in.DeviceName], true, false)
			if util.DropOnFilter(attrNew, f.lastMetadata[in.DeviceName], true) {
				continue // This Metric isn't in the white list so lets drop it.
			}
			ms = append(ms, OtelData{
				Name:  "kentik.snmp." + m,
				Value: float64(in.CustomBigInt[m]),
				Tags:  attrNew,
			})
		}
	}

	// Grap capacity utilization if possible.
	if f.lastMetadata[in.DeviceName] != nil {
		if ii, ok := f.lastMetadata[in.DeviceName].InterfaceInfo[in.InputPort]; ok {
			if speed, ok := ii["Speed"]; ok {
				if ispeed, ok := speed.(int32); ok {
					uptimeSpeed := in.CustomBigInt["Uptime"] * (int64(ispeed) * 10000) // Convert into bits here, from megabits. Also divide by 100 to convert uptime into seconds, from centi-seconds.
					if uptimeSpeed > 0 {
						attrNew := util.CopyAttrForSnmp(attr, "IfInUtilization", kt.MetricInfo{Oid: "computed", Mib: "computed", Profile: profileName, Table: "if"}, f.lastMetadata[in.DeviceName], true, false)
						if inBytes, ok := in.CustomBigInt["ifHCInOctets"]; ok {
							if !util.DropOnFilter(attrNew, f.lastMetadata[in.DeviceName], true) {
								ms = append(ms, OtelData{
									Name:  "kentik.snmp.IfInUtilization",
									Value: float64(inBytes*8*100) / float64(uptimeSpeed),
									Tags:  attrNew,
								})
							}
						}
					}
				}
			}
		}
		if oi, ok := f.lastMetadata[in.DeviceName].InterfaceInfo[in.OutputPort]; ok {
			if speed, ok := oi["Speed"]; ok {
				if ispeed, ok := speed.(int32); ok {
					uptimeSpeed := in.CustomBigInt["Uptime"] * (int64(ispeed) * 10000) // Convert into bits here, from megabits. Also divide by 100 to convert uptime into seconds, from centi-seconds.
					if uptimeSpeed > 0 {
						attrNew := util.CopyAttrForSnmp(attr, "IfOutUtilization", kt.MetricInfo{Oid: "computed", Mib: "computed", Profile: profileName, Table: "if"}, f.lastMetadata[in.DeviceName], true, false)
						if outBytes, ok := in.CustomBigInt["ifHCOutOctets"]; ok {
							if !util.DropOnFilter(attrNew, f.lastMetadata[in.DeviceName], true) {
								ms = append(ms, OtelData{
									Name:  "kentik.snmp.IfOutUtilization",
									Value: float64(outBytes*8*100) / float64(uptimeSpeed),
									Tags:  attrNew,
								})
							}
						}
					}
				}
			}
		}
	}

	return ms
}

func (f *OtelFormat) fromSnmpMetadata(in *kt.JCHF) []OtelData {
	if in.DeviceName == "" { // Only run if this is set.
		return nil
	}

	lm := util.SetMetadata(in)

	f.mux.Lock()
	defer f.mux.Unlock()
	if f.lastMetadata[in.DeviceName] == nil || lm.Size() >= f.lastMetadata[in.DeviceName].Size() {
		f.Infof("New Metadata for %s", in.DeviceName)
		f.lastMetadata[in.DeviceName] = lm
	} else {
		f.Infof("The metadata for %s was not updated since the attribute size is smaller. New = %d < Old = %d, Size difference = %v.",
			in.DeviceName, lm.Size(), f.lastMetadata[in.DeviceName].Size(), f.lastMetadata[in.DeviceName].Missing(lm))
	}

	return nil
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
