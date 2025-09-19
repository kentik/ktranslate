package otel

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/go-logr/stdr"
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
	"google.golang.org/grpc/credentials"
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
	inputs       map[string]chan OtelData
	trapLog      *OtelLogger
	logTee       chan string
}

const (
	CHAN_SLACK = 10000
)

var (
	endpoint   string
	protocol   string
	otelm      metric.Meter
	clientCert string
	clientKey  string
	rootCA     string
)

func init() {
	flag.StringVar(&endpoint, "otel.endpoint", "", "Send data to this endpoint.")
	flag.StringVar(&protocol, "otel.protocol", "stdout", "Send data using this protocol. (grpc,http,https,stdout)")
	flag.StringVar(&clientCert, "otel.tls_cert", "", "Load TLS client cert from file.")
	flag.StringVar(&clientKey, "otel.tls_key", "", "Load TLS client key from file.")
	flag.StringVar(&rootCA, "otel.root_ca", "", "Load TLS root CA from file.")
}

/*
*
Some usefule env vars to think about setting:

* OTEL_METRIC_EXPORT_INTERVAL=30000 -- time in ms to export. Default 60,000 (1 min).
* OTEL_EXPORTER_OTLP_COMPRESSION=gzip -- turn on gzip compression.
*/
func NewFormat(ctx context.Context, log logger.Underlying, cfg *ktranslate.OtelFormatConfig, logTee chan string) (*OtelFormat, error) {
	jf := &OtelFormat{
		ContextL:     logger.NewContextLFromUnderlying(logger.SContext{S: "otel"}, log),
		lastMetadata: map[string]*kt.LastMetadata{},
		invalids:     map[string]bool{},
		ctx:          context.Background(),
		vecs:         map[string]metric.Float64ObservableGauge{},
		config:       cfg,
		inputs:       map[string]chan OtelData{},
		logTee:       logTee,
	}

	var tlsC *tls.Config = nil
	if cfg.ClientCert != "" && cfg.ClientKey != "" && cfg.RootCA != "" {
		jf.Infof("Loading TLS certs from (cert=%s, key=%s) CA=%s", cfg.ClientCert, cfg.ClientKey, cfg.RootCA)
		c, err := getTls(cfg.ClientCert, cfg.ClientKey, cfg.RootCA)
		if err != nil {
			return nil, err
		}
		tlsC = c
	}

	// Set up a logger for otel.
	stdr.SetVerbosity(int(stdr.Info))
	logger := stdr.NewWithOptions(jf, stdr.Options{Depth: 2})
	otel.SetLogger(logger)

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

		if tlsC != nil {
			metricExporter, err := otlpmetrichttp.New(jf.ctx,
				otlpmetrichttp.WithEndpointURL(cfg.Endpoint),
				otlpmetrichttp.WithTLSClientConfig(tlsC),
			)
			if err != nil {
				return nil, err
			}
			exp = metricExporter
		} else {
			metricExporter, err := otlpmetrichttp.New(jf.ctx, otlpmetrichttp.WithEndpointURL(cfg.Endpoint))
			if err != nil {
				return nil, err
			}
			exp = metricExporter
		}
	case "grpc": // Same, use OTEL_EXPORTER_OTLP_COMPRESSION env var to turn on gzip compression.
		if cfg.Endpoint == "" {
			return nil, fmt.Errorf("-otel.endpoint required for grpc exports.")
		}

		if tlsC != nil {
			metricExporter, err := otlpmetricgrpc.New(jf.ctx,
				otlpmetricgrpc.WithEndpointURL(cfg.Endpoint),
				otlpmetricgrpc.WithTLSCredentials(
					// mutual tls.
					credentials.NewTLS(tlsC),
				),
			)
			if err != nil {
				return nil, err
			}
			exp = metricExporter
		} else {
			metricExporter, err := otlpmetricgrpc.New(jf.ctx, otlpmetricgrpc.WithEndpointURL(cfg.Endpoint))
			if err != nil {
				return nil, err
			}
			exp = metricExporter
		}
	default:
		return nil, fmt.Errorf("Invalid otel.protocol value. Currently supported: grpc,http,https,stdout")
	}

	meterProvider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp)))
	otel.SetMeterProvider(meterProvider)
	jf.exp = exp

	ol, err := NewLogger(ctx, jf, cfg, tlsC)
	if err != nil {
		return nil, err
	}
	jf.trapLog = ol

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

	f.mux.RLock()
	for _, m := range res {
		if _, ok := f.vecs[m.Name]; !ok {
			lm := m
			cv, err := otelm.Float64ObservableGauge(
				lm.Name,
				metric.WithFloat64Callback(func(_ context.Context, o metric.Float64Observer) error {
					f.Debugf("Exporting data from otel for %s", lm.Name)
					for {
						select {
						case mm := <-f.inputs[lm.Name]:
							o.Observe(mm.Value, metric.WithAttributeSet(mm.GetTagValues()))
						default:
							return nil
						}
					}
				}),
			)
			if err != nil {
				f.mux.RUnlock()
				return nil, err
			}
			f.mux.RUnlock()
			f.mux.Lock()
			f.vecs[lm.Name] = cv
			f.inputs[lm.Name] = make(chan OtelData, CHAN_SLACK)
			f.mux.Unlock()
			f.mux.RLock()
			f.Infof("Creating otel export for %s", m.Name)
		}
		// Save this for later, for the next time the async callback is run.
		f.inputs[m.Name] <- m
	}

	f.mux.RUnlock()
	return nil, nil
}

func (f *OtelFormat) From(raw *kt.Output) ([]map[string]interface{}, error) {
	values := make([]map[string]interface{}, 0)
	return values, nil
}

func (f *OtelFormat) Rollup(rolls []rollup.Rollup) (*kt.Output, error) {
	res := f.toOtelDataRollup(rolls)
	if len(res) == 0 {
		return nil, nil
	}

	f.mux.RLock()
	for _, m := range res {
		if _, ok := f.vecs[m.Name]; !ok {
			lm := m
			cv, err := otelm.Float64ObservableGauge(
				lm.Name,
				metric.WithFloat64Callback(func(_ context.Context, o metric.Float64Observer) error {
					f.Debugf("Exporting data from otel for %s", lm.Name)
					for {
						select {
						case mm := <-f.inputs[lm.Name]:
							o.Observe(mm.Value, metric.WithAttributeSet(mm.GetTagValues()))
						default:
							return nil
						}
					}
				}),
			)
			if err != nil {
				f.mux.RUnlock()
				return nil, err
			}
			f.mux.RUnlock()
			f.mux.Lock()
			f.vecs[lm.Name] = cv
			f.inputs[lm.Name] = make(chan OtelData, CHAN_SLACK)
			f.mux.Unlock()
			f.mux.RLock()
		}
		// Save this for later, for the next time the async callback is run.
		f.inputs[m.Name] <- m
	}

	f.mux.RUnlock()
	return nil, nil
}

var (
	rollupMap = map[string]string{
		"src_geo":     "src_country",
		"l4_src_port": "src_port",
		"dst_geo":     "dst_country",
		"l4_dst_port": "dst_port",
	}
)

func (f *OtelFormat) toOtelDataRollup(in []rollup.Rollup) []OtelData {
	ms := make([]OtelData, 0, len(in))
	for _, roll := range in {
		dims := roll.GetDims()
		attr := map[string]interface{}{}
		bad := false
		for i, pt := range strings.Split(roll.Dimension, roll.KeyJoin) {
			aname := dims[i]
			if n, ok := rollupMap[aname]; ok {
				aname = n
			}
			attr[aname] = pt
			if pt == "0" || pt == "" || pt == "-" || pt == "--" {
				bad = true
			}
			if aname == "src_port" || aname == "dst_port" { // Remap high ports down here.
				port, _ := strconv.Atoi(pt)
				if port > 32768 {
					attr[aname] = "32768"
				}
			}
		}
		if !bad {
			ms = append(ms, OtelData{
				Name:  "kentik.rollup." + roll.Name,
				Value: roll.Metric,
				Tags:  attr,
			})
		}
	}
	return ms
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
	case kt.KENTIK_EVENT_KTRANS_METRIC:
		return f.fromKtranslate(in)
	case kt.KENTIK_EVENT_SNMP_TRAP, kt.KENTIK_EVENT_EXT:
		// This is actually an event, send out as an event to sink directly.

		//err := f.trapLog.RecordLog(in, "New Trap Event")
		//if err != nil {
		//	f.Errorf("There was an error when sending an event: %v.", err)
		//	}
		// Debug in progress. Again.
		flat := in.Flatten()
		strip(flat)
		b, err := json.Marshal(flat)
		if err != nil {
			f.Errorf("There was an error when sending an event: %v.", err)
		}
		select {
		case f.logTee <- string(b):
		default:
			f.Errorf("There was an error processing a trap, log chan full.")
		}
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

// There's too much data to send as metrics here. Send on as a log instead.
func (f *OtelFormat) fromKflow(in *kt.JCHF) []OtelData {
	err := f.trapLog.RecordLog(in, "KFlow")
	if err != nil {
		f.Errorf("There was an error when sending an event: %v.", err)
	}
	return nil
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

func (f *OtelFormat) fromKtranslate(in *kt.JCHF) []OtelData {
	// Map the basic strings into here.
	attr := map[string]interface{}{}
	metrics := map[string]kt.MetricInfo{"name": kt.MetricInfo{}, "value": kt.MetricInfo{}, "count": kt.MetricInfo{}, "one-minute": kt.MetricInfo{}, "95-percentile": kt.MetricInfo{}, "du": kt.MetricInfo{}}
	f.mux.RLock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName], false)
	f.mux.RUnlock()
	ms := make([]OtelData, 0)

	switch in.CustomStr["type"] {
	case "counter":
		if in.CustomStr["force"] == "true" || in.CustomBigInt["count"] > 0 {
			ms = append(ms, OtelData{
				Name:  "kentik.ktranslate." + in.CustomStr["name"],
				Value: float64(in.CustomBigInt["count"]) / 100,
				Tags:  attr,
			})
		}
	case "gauge":
		if in.CustomStr["force"] == "true" || in.CustomBigInt["value"] > 0 {
			ms = append(ms, OtelData{
				Name:  "kentik.ktranslate." + in.CustomStr["name"],
				Value: float64(in.CustomBigInt["value"]) / 100,
				Tags:  attr,
			})
		}
	case "histogram":
		if in.CustomStr["force"] == "true" || in.CustomBigInt["95-percentile"] > 0 {
			ms = append(ms, OtelData{
				Name:  "kentik.ktranslate." + in.CustomStr["name"],
				Value: float64(in.CustomBigInt["95-percentile"]) / 100,
				Tags:  attr,
			})
		}
	case "meter":
		if in.CustomStr["force"] == "true" || in.CustomBigInt["one-minute"] > 0 {
			ms = append(ms, OtelData{
				Name:  "kentik.ktranslate." + in.CustomStr["name"],
				Value: float64(in.CustomBigInt["one-minute"]) / 100,
				Tags:  attr,
			})
		}
	case "timer":
		if in.CustomStr["force"] == "true" || in.CustomBigInt["95-percentile"] > 0 {
			ms = append(ms, OtelData{
				Name:  "kentik.ktranslate." + in.CustomStr["name"],
				Value: float64(in.CustomBigInt["95-percentile"]) / 100,
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

func (d *OtelData) GetTagValues() attribute.Set {
	res := make([]attribute.KeyValue, 0, len(d.Tags))
	for k, v := range d.Tags {
		switch t := v.(type) {
		case string:
			res = append(res, attribute.String(k, t))
		case int64:
			res = append(res, attribute.Int64(k, t))
		case int32:
			res = append(res, attribute.Int64(k, int64(t)))
		case float64:
			res = append(res, attribute.Float64(k, t))
		case uint32:
			res = append(res, attribute.Int64(k, int64(t)))
		case uint64:
			res = append(res, attribute.Int64(k, int64(t)))
		default:
			// Convert unknown types to string representation
			res = append(res, attribute.String(k, fmt.Sprintf("%v", t)))
		}
	}
	s, _ := attribute.NewSetWithFiltered(res, func(kv attribute.KeyValue) bool { return true })
	return s
}

// getTls returns a configuration that enables the use of mutual TLS.
func getTls(clientCert string, clientKey string, rootCA string) (*tls.Config, error) {
	clientAuth, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		return nil, err
	}

	caCert, err := os.ReadFile(rootCA)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	c := &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{clientAuth},
	}

	return c, nil
}

// From best I can tell, calldepth is always 2 here.
func (f *OtelFormat) Output(calldepth int, logline string) error {
	f.Infof("Otel: %s", logline)
	return nil
}
