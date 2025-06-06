package influx

import (
	"bytes"
	"flag"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/influxdata/line-protocol/v2/lineprotocol"
	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/formats/util"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

var (
	measurementPrefix string
	namespaceTokenSep string
)

func init() {
	flag.StringVar(&measurementPrefix, "influxdb_measurement_prefix", "", "Prefix metric names with this")
	flag.StringVar(&namespaceTokenSep, "influxdb_namespace_token", ":", "Use this token to seperate namespaces")
}

type InfluxFormat struct {
	logger.ContextL
	invalids     map[string]bool
	invalidsMux  sync.RWMutex
	lastMetadata map[string]*kt.LastMetadata
	mux          sync.RWMutex
	metrics      *InfluxMetrics
	config       *ktranslate.InfluxDBFormatConfig
}

type InfluxMetrics struct {
	EncoderErrors go_metrics.Counter
}

type InfluxData struct {
	Name         string
	FieldsFloat  map[string]float64
	Fields       map[string]int64
	FieldsString map[string]string
	Tags         map[string]interface{}
	Timestamp    int64
}

func (d *InfluxData) GetTags() string {
	var tags []string
	for key, v := range d.Tags {
		kval := strings.ReplaceAll(key, " ", "_")
		switch t := v.(type) {
		case string:
			tags = append(tags, fmt.Sprintf("%s=%s", kval, t))
		default:
			tags = append(tags, fmt.Sprintf("%s=%v", kval, v))
		}
	}

	sort.Strings(tags) // Need to be sorted for dedupe.
	return strings.Join(tags, ",")
}

func (d *InfluxData) Prefix() string {
	return fmt.Sprintf("%s,%s",
		d.Name,
		d.GetTags(),
	)
}

func NewMergedInfluxData(s InfluxDataSet) *InfluxData {
	if len(s) == 0 {
		return nil
	}
	d := InfluxData{
		Name:         s[0].Name,
		Tags:         s[0].Tags,
		Timestamp:    s[0].Timestamp,
		FieldsFloat:  map[string]float64{},
		Fields:       map[string]int64{},
		FieldsString: map[string]string{},
	}

	for _, f := range s {
		for k, v := range f.FieldsFloat {
			d.FieldsFloat[k] = v
		}
		for k, v := range f.Fields {
			d.Fields[k] = v
		}
		for k, v := range f.FieldsString {
			d.FieldsString[k] = v
		}
	}

	return &d
}

type InfluxDataSet []InfluxData

func (f InfluxDataSet) Len() int           { return len(f) }
func (f InfluxDataSet) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f InfluxDataSet) Less(i, j int) bool { return f[i].Name < f[j].Name }

func (f *InfluxFormat) report(err error, format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	f.invalidsMux.RLock()
	if !f.invalids[str] {
		f.invalidsMux.RUnlock()
		f.invalidsMux.Lock()
		// recheck since another thread may have written while we were upgrading lock
		if !f.invalids[str] {
			f.invalids[str] = true
			f.Errorf("influx encoding error on %s: %v", str, err)
		}
		f.invalidsMux.Unlock()
	} else {
		f.invalidsMux.RUnlock()
	}
	f.metrics.EncoderErrors.Inc(1)
}

func (f *InfluxFormat) Bytes(s InfluxDataSet) []byte {
	// Sort the output before anything.
	sort.Sort(s)

	// Next map common prefixes.
	prefixes := map[string][]InfluxData{}
	for _, l := range s {
		prefix := l.Prefix()
		if _, ok := prefixes[prefix]; !ok {
			prefixes[prefix] = []InfluxData{}
		}
		prefixes[prefix] = append(prefixes[prefix], l)
	}

	// Now merge down any common prefixes.
	merged := []*InfluxData{}
	for _, l := range prefixes {
		merged = append(merged, NewMergedInfluxData(l))
	}

	// Then format for output.
	var enc lineprotocol.Encoder
	enc.SetPrecision(lineprotocol.Nanosecond)
line:
	for _, l := range merged {
		if l != nil {
			enc.StartLine(l.Name)
			if enc.Err() != nil {
				f.report(enc.Err(), "measurement '%s'", l.Name)
				continue line
			}
			keys := make([]string, 0, len(l.Tags))
			for k := range l.Tags {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				v := l.Tags[k]
				s, ok := v.(string)
				if !ok {
					s = fmt.Sprintf("%v", v)
				}
				if s != "" {
					enc.AddTag(k, prepareTagValueMap(s))
				}
				if enc.Err() != nil {
					f.report(enc.Err(), "tag '%s'='%v'", k, s)
					continue line
				}
			}
			for k, v := range l.Fields {
				enc.AddField(k, lineprotocol.IntValue(v))
				if enc.Err() != nil {
					f.report(enc.Err(), "int field '%s'='%d'", k, v)
					continue line
				}
			}
			for k, v := range l.FieldsFloat {
				if fv, ok := lineprotocol.FloatValue(v); ok {
					enc.AddField(k, fv)
					if enc.Err() != nil {
						f.report(enc.Err(), "float field '%s'='%g'", k, v)
						continue line
					}
				}
			}
			for k, v := range l.FieldsString {
				if fv, ok := lineprotocol.StringValue(prepareTagValueMap(v)); ok {
					enc.AddField(k, fv)
					if enc.Err() != nil {
						f.report(enc.Err(), "string field '%s'='%s'", k, v)
						continue line
					}
				}
			}
			enc.EndLine(time.Unix(l.Timestamp, 0))
			if enc.Err() != nil {
				f.report(enc.Err(), "end of line on measurement '%s'", l.Name)
				continue line
			}
		}
	}
	return enc.Bytes()
}

var invalidTagChars string
var invalidTagCharsOnce sync.Once

func invalidTagCharsInit() {
	b := []byte{}
	for i := 0; i < 32; i++ {
		b = append(b, byte(i))
	}
	b = append(b, 0x7f)
	invalidTagChars = string(b)
}

func prepareTagValue(s string) string {
	invalidTagCharsOnce.Do(invalidTagCharsInit)
	if strings.ContainsAny(s, invalidTagChars) {
		var clean bytes.Buffer
		for _, c := range s {
			if !strings.ContainsRune(invalidTagChars, c) {
				clean.WriteRune(c)
			} else {
				clean.WriteRune(' ')
			}
		}
		return clean.String()
	}
	return s
}

func prepareTagValueMap(s string) string {
	invalidTagCharsOnce.Do(invalidTagCharsInit)
	if strings.ContainsAny(s, invalidTagChars) {
		swap := func(r rune) rune {
			switch {
			case r >= 0x00 && r <= 0x1f:
				return ' '
			case r == 0x7f:
				return ' '
			}
			return r
		}
		return strings.Map(swap, s)
	}
	return s
}

func NewFormat(log logger.Underlying, registry go_metrics.Registry, compression kt.Compression, cfg *ktranslate.InfluxDBFormatConfig) (*InfluxFormat, error) {
	if cfg == nil {
		return nil, fmt.Errorf("influxdb format cannot be nil")
	}
	jf := &InfluxFormat{
		ContextL:     logger.NewContextLFromUnderlying(logger.SContext{S: "influxFormat"}, log),
		invalids:     map[string]bool{},
		lastMetadata: map[string]*kt.LastMetadata{},
		metrics: &InfluxMetrics{
			EncoderErrors: go_metrics.GetOrRegisterCounter("influx_encoder_errors", registry),
		},
		config: cfg,
	}

	namespaceTokenSep = cfg.NamespaceToken

	return jf, nil
}

func (f *InfluxFormat) To(msgs []*kt.JCHF, serBuf []byte) (*kt.Output, error) {
	res := make([]InfluxData, 0, len(msgs))
	for _, m := range msgs {
		res = append(res, f.toInfluxMetric(m)...)
	}

	if len(res) == 0 {
		return nil, nil
	}

	return kt.NewOutput(f.Bytes(InfluxDataSet(res))), nil
}

func (f *InfluxFormat) toInfluxMetric(in *kt.JCHF) []InfluxData {
	switch in.EventType {
	case kt.KENTIK_EVENT_TYPE:
		return f.fromKflow(in)
	case kt.KENTIK_EVENT_SNMP_DEV_METRIC:
		return f.fromSnmpDeviceMetric(in)
	case kt.KENTIK_EVENT_SNMP_INT_METRIC:
		return f.fromSnmpInterfaceMetric(in)
	case kt.KENTIK_EVENT_SYNTH:
		return f.fromKSynth(in)
	case kt.KENTIK_EVENT_SNMP_METADATA:
		return f.fromSnmpMetadata(in)
	case kt.KENTIK_EVENT_KTRANS_METRIC:
		return f.fromKtranslate(in)
	default:
		f.invalidsMux.Lock()
		defer f.invalidsMux.Unlock()
		if !f.invalids[in.EventType] {
			f.Warnf("Invalid EventType: %s", in.EventType)
			f.invalids[in.EventType] = true
		}
	}

	return nil
}

// Not supported.
func (f *InfluxFormat) From(raw *kt.Output) ([]map[string]interface{}, error) {
	values := make([]map[string]interface{}, 0)
	return values, nil
}

func (f *InfluxFormat) Rollup(rolls []rollup.Rollup) (*kt.Output, error) {
	var enc lineprotocol.Encoder
	ts := time.Now()
line:
	for _, roll := range rolls {
		if roll.Metric == 0 {
			continue
		}
		mets := strings.Split(roll.EventType, ":")
		if len(mets) > 2 {
			enc.StartLine(roll.Name)
			if enc.Err() != nil {
				f.report(enc.Err(), "rollup measurement '%s'", roll.Name)
				continue line
			}

			dims := roll.GetDims()
			tags := map[string]string{}
			for i, pt := range strings.Split(roll.Dimension, roll.KeyJoin) {
				tags[dims[i]] = pt
			}
			keys := make([]string, 0, len(tags))
			for k := range tags {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				enc.AddTag(k, tags[k])
				if enc.Err() != nil {
					f.report(enc.Err(), "rollup tag '%s'='%s'", k, tags[k])
					continue line
				}
			}

			enc.AddField(mets[1], lineprotocol.IntValue(int64(roll.Metric)))
			if enc.Err() != nil {
				f.report(enc.Err(), "rollup int field '%s'", mets[1])
				continue line

			}
			enc.AddField("count", lineprotocol.IntValue(int64(roll.Count)))
			if enc.Err() != nil {
				f.report(enc.Err(), "rollup count field")
				continue line
			}
			enc.EndLine(ts)
			if enc.Err() != nil {
				f.report(enc.Err(), "rollup end of line on measurement '%s'", roll.Name)
				continue line
			}
		}
	}

	return kt.NewOutput(enc.Bytes()), nil
}

func (f *InfluxFormat) fromSnmpMetadata(in *kt.JCHF) []InfluxData {
	if in.DeviceName == "" { // Only run if this is set.
		return nil
	}

	lm := util.SetMetadata(in)

	f.mux.Lock()
	defer f.mux.Unlock()
	f.lastMetadata[in.DeviceName] = lm

	return nil
}

func (f *InfluxFormat) fromKtranslate(in *kt.JCHF) []InfluxData {
	// Map the basic strings into here.
	attr := map[string]interface{}{}
	metrics := map[string]kt.MetricInfo{"name": kt.MetricInfo{}, "value": kt.MetricInfo{}, "count": kt.MetricInfo{}, "one-minute": kt.MetricInfo{}, "95-percentile": kt.MetricInfo{}, "du": kt.MetricInfo{}}
	f.mux.RLock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName], false)
	f.mux.RUnlock()
	ms := make([]InfluxData, 0)
	name := strings.ReplaceAll(in.CustomStr["name"], "chf.kkc.", "")

	switch in.CustomStr["type"] {
	case "counter":
		if in.CustomStr["force"] == "true" || in.CustomBigInt["count"] > 0 {
			ms = append(ms, InfluxData{
				Name:        f.config.MeasurementPrefix + "ktranslate",
				FieldsFloat: map[string]float64{name: float64(in.CustomBigInt["count"]) / 100},
				Timestamp:   in.Timestamp,
				Tags:        attr,
			})
		}
	case "gauge":
		if in.CustomStr["force"] == "true" || in.CustomBigInt["value"] > 0 {
			ms = append(ms, InfluxData{
				Name:        f.config.MeasurementPrefix + "ktranslate",
				FieldsFloat: map[string]float64{name: float64(in.CustomBigInt["value"]) / 100},
				Timestamp:   in.Timestamp,
				Tags:        attr,
			})
		}
	case "histogram":
		if in.CustomStr["force"] == "true" || in.CustomBigInt["95-percentile"] > 0 {
			ms = append(ms, InfluxData{
				Name:        f.config.MeasurementPrefix + "ktranslate",
				FieldsFloat: map[string]float64{name: float64(in.CustomBigInt["95-percentile"]) / 100},
				Timestamp:   in.Timestamp,
				Tags:        attr,
			})
		}
	case "meter":
		if in.CustomStr["force"] == "true" || in.CustomBigInt["one-minute"] > 0 {
			ms = append(ms, InfluxData{
				Name:        f.config.MeasurementPrefix + "ktranslate",
				FieldsFloat: map[string]float64{name: float64(in.CustomBigInt["one-minute"]) / 100},
				Timestamp:   in.Timestamp,
				Tags:        attr,
			})
		}
	case "timer":
		if in.CustomStr["force"] == "true" || in.CustomBigInt["95-percentile"] > 0 {
			ms = append(ms, InfluxData{
				Name:        f.config.MeasurementPrefix + "ktranslate",
				FieldsFloat: map[string]float64{name: float64(in.CustomBigInt["95-percentile"]) / 100},
				Timestamp:   in.Timestamp,
				Tags:        attr,
			})
		}
	}
	return ms
}

func (f *InfluxFormat) fromKSynth(in *kt.JCHF) []InfluxData {
	metrics := util.GetSynMetricNameSet(in.CustomInt["result_type"])
	attr := map[string]interface{}{}
	f.mux.RLock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName], true)
	f.mux.RUnlock()
	ms := map[string]int64{}

	for m, name := range metrics {
		switch m {
		case "error", "timeout":
			ms[name.Name] = 1
		default:
			if in.CustomInt["result_type"] > 1 {
				ms[name.Name] = int64(in.CustomInt[m])
			}
		}
	}

	return []InfluxData{{
		Name:      f.config.MeasurementPrefix + "ksynth",
		Fields:    ms,
		Timestamp: in.Timestamp,
		Tags:      attr,
	}}
}

func (f *InfluxFormat) fromKflow(in *kt.JCHF) []InfluxData {
	// Map the basic strings into here.
	attr := map[string]interface{}{}
	metrics := map[string]kt.MetricInfo{"in_bytes": {}, "out_bytes": {}, "in_pkts": {}, "out_pkts": {}, "latency_ms": {}}
	f.mux.RLock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName], true)
	f.mux.RUnlock()
	ms := map[string]int64{}
	for m := range metrics {
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

	return []InfluxData{{
		Name:      f.config.MeasurementPrefix + "flow",
		Fields:    ms,
		Timestamp: in.Timestamp,
		Tags:      attr,
	}}
}

func (f *InfluxFormat) fromSnmpDeviceMetric(in *kt.JCHF) []InfluxData {
	metrics := in.CustomMetrics
	attr := map[string]interface{}{}
	f.mux.RLock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName], true)
	f.mux.RUnlock()
	ip := attr["src_addr"]

	results := []InfluxData{}
	for m, name := range metrics {
		if m == "" {
			f.Errorf("Missing metric name, skipping %v", attr)
			continue
		}
		if _, ok := in.CustomBigInt[m]; ok {
			attrNew := util.CopyAttrForSnmp(attr, m, name, f.lastMetadata[in.DeviceName], false, true)
			if util.DropOnFilter(attrNew, f.lastMetadata[in.DeviceName], false) {
				continue // This Metric isn't in the white list so lets drop it.
			}

			mib := getMib(attrNew, ip)
			metricName := m
			if strings.Contains(metricName, namespaceTokenSep) {
				pts := strings.SplitN(metricName, namespaceTokenSep, 2)
				metricName = pts[1]
			}
			if name.Format == kt.FloatMS {
				results = append(results, InfluxData{
					Name:        f.config.MeasurementPrefix + mib,
					FieldsFloat: map[string]float64{metricName: float64(float64(in.CustomBigInt[m]) / 1000)},
					Timestamp:   in.Timestamp,
					Tags:        attrNew,
				})
			} else {
				results = append(results, InfluxData{
					Name:      f.config.MeasurementPrefix + mib,
					Fields:    map[string]int64{metricName: int64(in.CustomBigInt[m])},
					Timestamp: in.Timestamp,
					Tags:      attrNew,
				})
			}
		}
	}

	return results
}

func (f *InfluxFormat) fromSnmpInterfaceMetric(in *kt.JCHF) []InfluxData {
	metrics := in.CustomMetrics
	attr := map[string]interface{}{}
	f.mux.RLock()
	defer f.mux.RUnlock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName], true)
	ip := attr["src_ip"]

	profileName := "snmp"
	results := []InfluxData{}
	for m, name := range metrics {
		if m == "" {
			f.Errorf("Missing metric name, skipping %v", attr)
			continue
		}
		profileName = name.Profile
		if _, ok := in.CustomBigInt[m]; ok {
			attrNew := util.CopyAttrForSnmp(attr, m, name, f.lastMetadata[in.DeviceName], false, true)
			if util.DropOnFilter(attrNew, f.lastMetadata[in.DeviceName], true) {
				continue // This Metric isn't in the white list so lets drop it.
			}

			mib := getMib(attrNew, ip)
			metricName := m
			if strings.Contains(metricName, namespaceTokenSep) {
				pts := strings.SplitN(metricName, namespaceTokenSep, 2)
				metricName = pts[1]
			}
			if name.Format == kt.FloatMS {
				results = append(results, InfluxData{
					Name:        f.config.MeasurementPrefix + mib,
					FieldsFloat: map[string]float64{metricName: float64(float64(in.CustomBigInt[m]) / 1000)},
					Timestamp:   in.Timestamp,
					Tags:        attrNew,
				})
			} else {
				results = append(results, InfluxData{
					Name:      f.config.MeasurementPrefix + mib,
					Fields:    map[string]int64{metricName: int64(in.CustomBigInt[m])},
					Timestamp: in.Timestamp,
					Tags:      attrNew,
				})
				if sv, ok := in.CustomStr[kt.StringPrefix+m]; ok {
					results[len(results)-1].FieldsString = map[string]string{m + "_str": sv}
				}
			}
		}
	}

	// Grap varioius rates computed over time here if possible.
	if f.lastMetadata[in.DeviceName] != nil {
		results = f.setRates("In", in, results, attr, profileName, ip)
		results = f.setRates("Out", in, results, attr, profileName, ip)
	}

	return results
}

func (f *InfluxFormat) setRates(direction string, in *kt.JCHF, results []InfluxData, attr map[string]interface{}, profileName string, ip interface{}) []InfluxData {
	var port kt.IfaceID
	if direction == "In" {
		port = in.InputPort
	} else {
		port = in.OutputPort
	}
	utilName := fmt.Sprintf("If%sUtilization", direction)
	bitRate := fmt.Sprintf("If%sBitRate", direction)
	pktRate := fmt.Sprintf("If%sPktRate", direction)
	totalBytes := in.CustomBigInt[fmt.Sprintf("ifHC%sOctets", direction)]
	totalPkts := in.CustomBigInt[fmt.Sprintf("ifHC%sUcastPkts", direction)] + in.CustomBigInt[fmt.Sprintf("ifHC%sMulticastPkts", direction)] + in.CustomBigInt[fmt.Sprintf("ifHC%sBroadcastPkts", direction)]

	if ii, ok := f.lastMetadata[in.DeviceName].InterfaceInfo[port]; ok {
		if speed, ok := ii["Speed"]; ok {
			if ispeed, ok := speed.(int32); ok {
				uptime := in.CustomBigInt["Uptime"]
				uptimeSpeed := uptime * (int64(ispeed) * 10000) // Convert into bits here, from megabits. Also divide by 100 to convert uptime into seconds, from centi-seconds.
				if uptimeSpeed > 0 {
					attrNew := util.CopyAttrForSnmp(attr, utilName, kt.MetricInfo{Oid: "computed", Mib: "computed", Profile: profileName, Table: "if"}, f.lastMetadata[in.DeviceName], false, true)
					if !util.DropOnFilter(attrNew, f.lastMetadata[in.DeviceName], true) {
						getMib(attrNew, ip)
						if totalBytes > 0 {
							results = append(results, InfluxData{
								Name:        f.config.MeasurementPrefix + "IF-MIB::if",
								FieldsFloat: map[string]float64{utilName: float64(totalBytes*8*100) / float64(uptimeSpeed)},
								Timestamp:   in.Timestamp,
								Tags:        attrNew,
							})
							results = append(results, InfluxData{
								Name:        f.config.MeasurementPrefix + "IF-MIB::if",
								FieldsFloat: map[string]float64{bitRate: float64(totalBytes*8*100) / float64(uptime)},
								Timestamp:   in.Timestamp,
								Tags:        attrNew,
							})
						}
						if totalPkts > 0 {
							results = append(results, InfluxData{
								Name:        f.config.MeasurementPrefix + "IF-MIB::if",
								FieldsFloat: map[string]float64{pktRate: float64(totalPkts*100) / float64(uptime)},
								Timestamp:   in.Timestamp,
								Tags:        attrNew,
							})
						}
					}
				}
			}
		}
	}

	return results
}

func getMib(attr map[string]interface{}, ip interface{}) string {
	// Remove any lingering droppable fields.
	for k, v := range attr {
		if _, ok := dropFields[k]; ok {
			delete(attr, k)
		}
		if k == "Index" {
			delete(attr, k)
			attr["index"] = v
		}

		// If the attribute has a namespace prefix, drop here.
		if strings.Contains(k, namespaceTokenSep) {
			pts := strings.SplitN(k, namespaceTokenSep, 2)
			delete(attr, k)
			attr[pts[1]] = v
		}
	}
	if ip != nil {
		attr["device_ip"] = ip
	}

	// And now figure out what the mib name is.
	mib, ok := attr["mib-name"].(string)
	if !ok {
		return "device"
	}
	delete(attr, "mib-name")
	if mib == "" {
		mib = "device"
	}

	// If there's a table, add this to the info.
	mibTable, ok := attr["mib-table"].(string)
	if ok {
		// If the MIB is normalized use "/" as separator
		if strings.HasPrefix(mib, "/") || strings.HasSuffix(mib, "/") {
			mib = strings.TrimRight(mib, "/") + "/" + strings.TrimLeft(mibTable, "/")
		} else {
			mib = mib + "::" + mibTable
		}
		delete(attr, "mib-table")
	}

	return mib
}

var dropFields = map[string]bool{
	//"device_name": true,
	"objectIdentifier":     true,
	"eventType":            true,
	"provider":             true,
	"SysServices":          true,
	"instrumentation.name": true,
}
