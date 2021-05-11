package influx

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kentik/ktranslate/pkg/formats/util"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

var (
	Measurement = flag.String("measurement", "kflow", "Measurement to use for rollups.")
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
		fields[i] = k + "=" + strconv.FormatInt(v, 10)
		i++
	}
	for k, v := range d.FieldsFloat {
		fields[i] = k + "=" + strconv.FormatFloat(v, 'f', 4, 64)
		i++
	}

	tags := make([]string, len(d.Tags))
	i = 0
	for key, v := range d.Tags {
		kval := strings.ReplaceAll(key, " ", "_")
		switch t := v.(type) {
		case string:
			if strings.ContainsAny(t, " ") {
				tags[i] = fmt.Sprintf("%s=\"%s\"", kval, t)
			} else {
				tags[i] = fmt.Sprintf("%s=%s", kval, t)
			}
		default:
			tags[i] = fmt.Sprintf("%s=%v", kval, v)
		}

		i++
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

func (f *InfluxFormat) To(msgs []*kt.JCHF, serBuf []byte) ([]byte, error) {
	res := make([]InfluxData, 0, len(msgs))
	for _, m := range msgs {
		res = append(res, f.toInfluxMetric(m)...)
	}

	if len(res) == 0 {
		return nil, nil
	}

	return InfluxDataSet(res).Bytes(), nil
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
func (f *InfluxFormat) From(raw []byte) ([]map[string]interface{}, error) {
	values := make([]map[string]interface{}, 0)
	return values, nil
}

func (f *InfluxFormat) Rollup(rolls []rollup.Rollup) ([]byte, error) {
	res := make([]string, len(rolls))
	ts := time.Now()
	for i, r := range rolls {
		pkts := strings.Split(r.EventType, ":")
		if len(pkts) > 2 {
			res[i] = fmt.Sprintf("%s,%s=%s %s=%d %d", *Measurement, strings.Join(pkts[2:], ":"), r.Dimension, pkts[1], uint64(r.Metric), ts.UnixNano()) // Time to nano
		}
	}

	return []byte(strings.Join(res, "\n")), nil
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
			ms[name] = 1
		default:
			if in.CustomInt["result_type"] > 1 {
				ms[name] = int64(in.CustomInt[m])
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
	metrics := map[string]string{"in_bytes": "", "out_bytes": "", "in_pkts": "", "out_pkts": "", "latency_ms": ""}
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
	ms := map[string]int64{}
	for m, _ := range metrics {
		if _, ok := in.CustomBigInt[m]; ok {
			ms[m] = in.CustomBigInt[m]
		}
	}

	for k, v := range attr { // Weed out any spaces which might break things.
		if sv, ok := v.(string); ok {
			if strings.Contains(sv, " ") {
				delete(attr, k)
			}
		}
	}

	return []InfluxData{InfluxData{
		Name:      "kentik.snmp.device",
		Fields:    ms,
		Timestamp: in.Timestamp * 1000000000,
		Tags:      attr,
	}}
}

func (f *InfluxFormat) fromSnmpInterfaceMetric(in *kt.JCHF) []InfluxData {
	metrics := in.CustomMetrics
	attr := map[string]interface{}{}
	f.mux.RLock()
	defer f.mux.RUnlock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName])
	ms := map[string]int64{}
	msF := map[string]float64{}
	for m, _ := range metrics {
		if _, ok := in.CustomBigInt[m]; ok {
			ms[m] = in.CustomBigInt[m]
		}
	}

	for k, v := range attr { // Weed out any spaces which might break things.
		if sv, ok := v.(string); ok {
			if strings.Contains(sv, " ") {
				delete(attr, k)
			}
		}
	}

	// Grap capacity utilization if possible.
	if f.lastMetadata[in.DeviceName] != nil {
		if ii, ok := f.lastMetadata[in.DeviceName].InterfaceInfo[in.InputPort]; ok {
			if speed, ok := ii["Speed"]; ok {
				if ispeed, ok := speed.(int32); ok {
					uptimeSpeed := in.CustomBigInt["Uptime"] * (int64(ispeed) * 1000000) // Convert into bits here, from megabits.
					if uptimeSpeed > 0 {
						msF["IfInUtilization"] = float64(in.CustomBigInt["ifHCInOctets"]*8*100) / float64(uptimeSpeed)
					}
				}
			}
		}
		if oi, ok := f.lastMetadata[in.DeviceName].InterfaceInfo[in.OutputPort]; ok {
			if speed, ok := oi["Speed"]; ok {
				if ispeed, ok := speed.(int32); ok {
					uptimeSpeed := in.CustomBigInt["Uptime"] * (int64(ispeed) * 1000000) // Convert into bits here, from megabits.
					if uptimeSpeed > 0 {
						msF["IfOutUtilization"] = float64(in.CustomBigInt["ifHCOutOctets"]*8*100) / float64(uptimeSpeed)
					}
				}
			}
		}
	}

	return []InfluxData{InfluxData{
		Name:        "kentik.snmp.interface",
		Fields:      ms,
		FieldsFloat: msF,
		Timestamp:   in.Timestamp * 1000000000,
		Tags:        attr,
	}}
}
