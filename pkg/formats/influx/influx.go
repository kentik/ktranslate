package influx

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kentik/ktranslate/pkg/formats/util"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

type InfluxFormat struct {
	logger.ContextL
	invalids     map[string]bool
	lastMetadata map[string]*kt.LastMetadata
	mux          sync.RWMutex
}

type InfluxData struct {
	Name        string
	FieldsFloat map[string]float64
	Fields      map[string]int64
	Tags        map[string]interface{}
	Timestamp   int64
}

func (d *InfluxData) String() string {
	fields := make([]string, len(d.Fields)+len(d.FieldsFloat))
	i := 0
	for k, v := range d.Fields {
		fields[i] = k + "=" + strconv.FormatInt(v, 10) + "i"
		i++
	}
	for k, v := range d.FieldsFloat {
		fields[i] = k + "=" + strconv.FormatFloat(v, 'f', 4, 64)
		i++
	}

	var tags []string
	for key, v := range d.Tags {
		kval := strings.ReplaceAll(key, " ", "_")
		switch t := v.(type) {
		case string:
			if t != "" {
				tags = append(tags, fmt.Sprintf("%s=%s", kval, influxEscape(t)))
			}
		default:
			tags = append(tags, fmt.Sprintf("%s=%v", kval, v))
		}
	}

	return fmt.Sprintf("%s,%s %s %d",
		d.Name,
		strings.Join(tags, ","),
		strings.Join(fields, ","),
		d.Timestamp)
}

type InfluxDataSet []InfluxData

func (s InfluxDataSet) Bytes() []byte {
	res := make([]string, 0)
	for _, l := range s {
		res = append(res, l.String())
	}
	return []byte(strings.Join(res, "\n"))
}

func NewFormat(log logger.Underlying, compression kt.Compression) (*InfluxFormat, error) {
	jf := &InfluxFormat{
		ContextL:     logger.NewContextLFromUnderlying(logger.SContext{S: "influxFormat"}, log),
		invalids:     map[string]bool{},
		lastMetadata: map[string]*kt.LastMetadata{},
	}

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

	return kt.NewOutput(InfluxDataSet(res).Bytes()), nil
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

// Not supported.
func (f *InfluxFormat) From(raw *kt.Output) ([]map[string]interface{}, error) {
	values := make([]map[string]interface{}, 0)
	return values, nil
}

func (f *InfluxFormat) Rollup(rolls []rollup.Rollup) (*kt.Output, error) {
	res := make([]string, 0)
	ts := time.Now()
	for _, roll := range rolls {
		if roll.Metric == 0 {
			continue
		}

		dims := roll.GetDims()
		mets := strings.Split(roll.EventType, ":")
		attr := []string{}
		for i, pt := range strings.Split(roll.Dimension, roll.KeyJoin) {
			attr = append(attr, dims[i]+"="+influxEscape(pt))
		}
		if len(mets) > 2 {
			res = append(res, fmt.Sprintf("%s,%s %s=%d,count=%d %d", roll.Name, strings.Join(attr, ","), mets[1], uint64(roll.Metric), roll.Count, ts.UnixNano())) // Time to nano
		}
	}

	return kt.NewOutput([]byte(strings.Join(res, "\n"))), nil
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

func (f *InfluxFormat) fromKSynth(in *kt.JCHF) []InfluxData {
	metrics := util.GetSynMetricNameSet(in.CustomInt["result_type"])
	attr := map[string]interface{}{}
	f.mux.RLock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName])
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

	return []InfluxData{InfluxData{
		Name:      "kentik.synth",
		Fields:    ms,
		Timestamp: in.Timestamp * 1000000000,
		Tags:      attr,
	}}
}

func (f *InfluxFormat) fromKflow(in *kt.JCHF) []InfluxData {
	// Map the basic strings into here.
	attr := map[string]interface{}{}
	metrics := map[string]kt.MetricInfo{"in_bytes": kt.MetricInfo{}, "out_bytes": kt.MetricInfo{}, "in_pkts": kt.MetricInfo{}, "out_pkts": kt.MetricInfo{}, "latency_ms": kt.MetricInfo{}}
	f.mux.RLock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName])
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

	return []InfluxData{InfluxData{
		Name:      "kentik.flow",
		Fields:    ms,
		Timestamp: in.Timestamp * 1000000000,
		Tags:      attr,
	}}
}

func (f *InfluxFormat) fromSnmpDeviceMetric(in *kt.JCHF) []InfluxData {
	metrics := in.CustomMetrics
	attr := map[string]interface{}{}
	f.mux.RLock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName])
	f.mux.RUnlock()

	results := []InfluxData{}
	for m, name := range metrics {
		if m == "" {
			f.Errorf("Missing metric name, skipping %v", attr)
			continue
		}
		if _, ok := in.CustomBigInt[m]; ok {
			attrNew := util.CopyAttrForSnmp(attr, m, name, f.lastMetadata[in.DeviceName])
			if util.DropOnFilter(attrNew, f.lastMetadata[in.DeviceName], false) {
				continue // This Metric isn't in the white list so lets drop it.
			}

			mtype := name.GetType()
			if name.Format == kt.FloatMS {
				results = append(results, InfluxData{
					Name:        "kentik." + mtype,
					FieldsFloat: map[string]float64{m: float64(float64(in.CustomBigInt[m]) / 1000)},
					Timestamp:   in.Timestamp * 1000000000,
					Tags:        attrNew,
				})
			} else {
				results = append(results, InfluxData{
					Name:      "kentik." + mtype,
					Fields:    map[string]int64{m: int64(in.CustomBigInt[m])},
					Timestamp: in.Timestamp * 1000000000,
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
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName])

	profileName := "snmp"
	results := []InfluxData{}
	for m, name := range metrics {
		if m == "" {
			f.Errorf("Missing metric name, skipping %v", attr)
			continue
		}
		profileName = name.Profile
		if _, ok := in.CustomBigInt[m]; ok {
			attrNew := util.CopyAttrForSnmp(attr, m, name, f.lastMetadata[in.DeviceName])
			if util.DropOnFilter(attrNew, f.lastMetadata[in.DeviceName], true) {
				continue // This Metric isn't in the white list so lets drop it.
			}

			mtype := name.GetType()
			if name.Format == kt.FloatMS {
				results = append(results, InfluxData{
					Name:        "kentik." + mtype,
					FieldsFloat: map[string]float64{m: float64(float64(in.CustomBigInt[m]) / 1000)},
					Timestamp:   in.Timestamp * 1000000000,
					Tags:        attrNew,
				})
			} else {
				results = append(results, InfluxData{
					Name:      "kentik." + mtype,
					Fields:    map[string]int64{m: int64(in.CustomBigInt[m])},
					Timestamp: in.Timestamp * 1000000000,
					Tags:      attrNew,
				})
			}
		}
	}

	// Grap capacity utilization if possible.
	if f.lastMetadata[in.DeviceName] != nil {
		if ii, ok := f.lastMetadata[in.DeviceName].InterfaceInfo[in.InputPort]; ok {
			if speed, ok := ii["Speed"]; ok {
				if ispeed, ok := speed.(int32); ok {
					uptimeSpeed := in.CustomBigInt["Uptime"] * (int64(ispeed) * 10000) // Convert into bits here, from megabits. Also divide by 100 to convert uptime into seconds, from centi-seconds.
					if uptimeSpeed > 0 {
						attrNew := util.CopyAttrForSnmp(attr, "IfInUtilization", kt.MetricInfo{Oid: "computed", Mib: "computed", Profile: profileName, Table: "if"}, f.lastMetadata[in.DeviceName])
						if inBytes, ok := in.CustomBigInt["ifHCInOctets"]; ok {
							if !util.DropOnFilter(attrNew, f.lastMetadata[in.DeviceName], true) {
								results = append(results, InfluxData{
									Name:        "kentik.snmp",
									FieldsFloat: map[string]float64{"IfInUtilization": float64(inBytes*8*100) / float64(uptimeSpeed)},
									Timestamp:   in.Timestamp * 1000000000,
									Tags:        attrNew,
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
						attrNew := util.CopyAttrForSnmp(attr, "IfOutUtilization", kt.MetricInfo{Oid: "computed", Mib: "computed", Profile: profileName, Table: "if"}, f.lastMetadata[in.DeviceName])
						if outBytes, ok := in.CustomBigInt["ifHCOutOctets"]; ok {
							if !util.DropOnFilter(attrNew, f.lastMetadata[in.DeviceName], true) {
								results = append(results, InfluxData{
									Name:        "kentik.snmp",
									FieldsFloat: map[string]float64{"IfOutUtilization": float64(outBytes*8*100) / float64(uptimeSpeed)},
									Timestamp:   in.Timestamp * 1000000000,
									Tags:        attrNew,
								})
							}
						}
					}
				}
			}
		}
	}

	return results
}

var escaper = regexp.MustCompile("([,= \\s])")

// Escape special characters according to https://docs.influxdata.com/influxdb/v1.8/write_protocols/line_protocol_tutorial/#special-characters-and-keywords
func influxEscape(s string) string {
	if strings.ContainsAny(s, ",= \t\r\n") {
		return string(escaper.ReplaceAll([]byte(s), []byte("\\$1")))
	} else {
		return s
	}
}
