package influx

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

var (
	Measurement  = flag.String("measurement", "kflow", "Measurement to use for rollups.")
	DroppedAttrs = map[string]bool{
		"timestamp":               true,
		"sampled_packet_size":     true,
		"Lat/Long Dest":           true,
		"MEMBER_ID":               true,
		"dst_eth_mac":             true,
		"src_eth_mac":             true,
		"Manufacturer":            true,
		"Error Cause/Trace Route": true,
		"Hop Data":                true,
		"STR01":                   true,
		"ULT_EXIT_PORT":           true,
		"Task ID":                 true,
		"APP_PROTOCOL":            true,
		"Agent ID":                true,
		"ULT_EXIT_DEVICE_ID":      true,
		"device_id":               true,
		"kt_functional_testing":   true,
		"CLIENT_NW_LATENCY_MS":    true,
		"APPL_LATENCY_MS":         true,
		"SERVER_NW_LATENCY_MS":    true,
		"CONNECTION_ID":           true,
	}
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
	for k, v := range d.Tags {
		switch t := v.(type) {
		case string:
			if strings.ContainsAny(t, " ") {
				tags[i] = fmt.Sprintf("%s=\"%s\"", k, t)
			} else {
				tags[i] = fmt.Sprintf("%s=%s", k, t)
			}
		default:
			tags[i] = fmt.Sprintf("%s=%v", k, v)
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
	lm := kt.LastMetadata{
		DeviceInfo:    map[string]interface{}{},
		InterfaceInfo: map[kt.IfaceID]map[string]interface{}{},
	}
	for k, v := range in.CustomStr {
		if DroppedAttrs[k] {
			continue // Skip because we don't want this messing up cardinality.
		}
		if strings.HasPrefix(k, "if.") {
			pts := strings.SplitN(k, ".", 3)
			if len(pts) == 3 {
				if ifint, err := strconv.Atoi(pts[1]); err == nil {
					if _, ok := lm.InterfaceInfo[kt.IfaceID(ifint)]; !ok {
						lm.InterfaceInfo[kt.IfaceID(ifint)] = map[string]interface{}{}
					}
					if v != "" {
						lm.InterfaceInfo[kt.IfaceID(ifint)][pts[2]] = v
					}
				}
			}
		} else {
			if v != "" {
				lm.DeviceInfo[k] = v
			}
		}
	}
	for k, v := range in.CustomInt {
		if DroppedAttrs[k] {
			continue // Skip because we don't want this messing up cardinality.
		}
		if strings.HasPrefix(k, "if.") {
			pts := strings.SplitN(k, ".", 3)
			if len(pts) == 3 {
				if ifint, err := strconv.Atoi(pts[1]); err == nil {
					if _, ok := lm.InterfaceInfo[kt.IfaceID(ifint)]; !ok {
						lm.InterfaceInfo[kt.IfaceID(ifint)] = map[string]interface{}{}
					}
					lm.InterfaceInfo[kt.IfaceID(ifint)][pts[2]] = v
				}
			}
		} else {
			lm.DeviceInfo[k] = v
		}
	}

	f.lastMetadata[in.DeviceName] = &lm
	return nil
}

func (f *InfluxFormat) fromKSynth(in *kt.JCHF) []InfluxData {
	var metrics map[string]bool
	var names map[string]string
	switch in.CustomInt["Result Type"] {
	case 0: // Error
		metrics = map[string]bool{"Error": true}
	case 1: // Timeout
		metrics = map[string]bool{"Timeout": true}
	case 2: // Ping
		metrics = map[string]bool{"Fetch Status | Ping Sent | Trace Time": true, "Fetch TTLB | Ping Lost": true,
			"Fetch Size | Ping Min RTT": true, "Ping Max RTT": true, "Ping Avg RTT": true, "Ping Std RTT": true, "Ping Jit RTT": true}
		names = map[string]string{"Fetch Status | Ping Sent | Trace Time": "Sent", "Fetch TTLB | Ping Lost": "Lost",
			"Fetch Size | Ping Min RTT": "MinRTT", "Ping Max RTT": "MaxRTT", "Ping Avg RTT": "AvgRTT", "Ping Std RTT": "StdRTT", "Ping Jit RTT": "JitRTT"}
	case 3: // Fetch
		metrics = map[string]bool{"Fetch Status | Ping Sent | Trace Time": true, "Fetch TTLB | Ping Lost": true, "Fetch Size | Ping Min RTT": true}
		names = map[string]string{"Fetch Status | Ping Sent | Trace Time": "Status", "Fetch TTLB | Ping Lost": "TTLB", "Fetch Size | Ping Min RTT": "Size"}
	case 4: // Trace
		metrics = map[string]bool{"Fetch Status | Ping Sent | Trace Time": true}
		names = map[string]string{"Fetch Status | Ping Sent | Trace Time": "Time"}
	case 5: // Knock
		metrics = map[string]bool{"Fetch Status | Ping Sent | Trace Time": true, "Fetch TTLB | Ping Lost": true,
			"Fetch Size | Ping Min RTT": true, "Ping Max RTT": true, "Ping Avg RTT": true, "Ping Std RTT": true, "Ping Jit RTT": true}
		names = map[string]string{"Fetch Status | Ping Sent | Trace Time": "Sent", "Fetch TTLB | Ping Lost": "Lost",
			"Fetch Size | Ping Min RTT": "MinRTT", "Ping Max RTT": "MaxRTT", "Ping Avg RTT": "AvgRTT", "Ping Std RTT": "StdRTT", "Ping Jit RTT": "JitRTT"}
	case 6: // Query
		metrics = map[string]bool{"Fetch Status | Ping Sent | Trace Time": true, "Fetch TTLB | Ping Lost": true}
		names = map[string]string{"Fetch Status | Ping Sent | Trace Time": "Time", "Fetch TTLB | Ping Lost": "Code"}
	case 7: // Shake
		metrics = map[string]bool{"Fetch Status | Ping Sent | Trace Time": true, "Lat/Long Dest": true}
		names = map[string]string{"Fetch Status | Ping Sent | Trace Time": "Time", "Lat/Long Dest": "Port"}
	}

	attr := map[string]interface{}{}
	f.setAttr(attr, in, metrics)
	ms := map[string]int64{}

	for m, _ := range metrics {
		switch m {
		case "Error", "Timeout":
			ms[m] = 1
		default:
			ms[names[m]] = int64(in.CustomInt[m])
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
	metrics := map[string]bool{"in_bytes": true, "out_bytes": true, "in_pkts": true, "out_pkts": true, "latency_ms": true}
	f.setAttr(attr, in, metrics)
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
			ms[m] = int64(in.CustomInt["APPL_LATENCY_MS"])
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
	var metrics map[string]bool
	if len(in.CustomMetrics) > 0 {
		metrics = in.CustomMetrics
	} else {
		metrics = map[string]bool{"CPU": true, "MemoryTotal": true, "MemoryUsed": true, "MemoryFree": true, "MemoryUtilization": true, "Uptime": true}
	}
	attr := map[string]interface{}{}
	f.setAttr(attr, in, metrics)
	ms := map[string]int64{}
	for m, _ := range metrics {
		if _, ok := in.CustomBigInt[m]; ok {
			ms[m] = in.CustomBigInt[m]
		}
	}

	return []InfluxData{InfluxData{
		Name:      "kentik.snmp",
		Fields:    ms,
		Timestamp: in.Timestamp * 1000000000,
		Tags:      attr,
	}}
}

func (f *InfluxFormat) fromSnmpInterfaceMetric(in *kt.JCHF) []InfluxData {
	var metrics map[string]bool
	if len(in.CustomMetrics) > 0 {
		metrics = in.CustomMetrics
	} else {
		metrics = map[string]bool{"ifHCInOctets": true, "ifHCInUcastPkts": true, "ifHCOutOctets": true, "ifHCOutUcastPkts": true, "ifInErrors": true, "ifOutErrors": true,
			"ifInDiscards": true, "ifOutDiscards": true, "ifHCOutMulticastPkts": true, "ifHCOutBroadcastPkts": true, "ifHCInMulticastPkts": true, "ifHCInBroadcastPkts": true}
	}
	attr := map[string]interface{}{}
	f.setAttr(attr, in, metrics)
	ms := map[string]int64{}
	msF := map[string]float64{}
	for m, _ := range metrics {
		if _, ok := in.CustomBigInt[m]; ok {
			ms[m] = in.CustomBigInt[m]
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
		Name:        "kentik.snmp",
		Fields:      ms,
		FieldsFloat: msF,
		Timestamp:   in.Timestamp * 1000000000,
		Tags:        attr,
	}}
}

func (f *InfluxFormat) setAttr(attr map[string]interface{}, in *kt.JCHF, metrics map[string]bool) {
	mapr := in.Flatten()
	for k, v := range mapr {
		if DroppedAttrs[k] {
			continue // Skip because we don't want this messing up cardinality.
		}

		switch vt := v.(type) {
		case string:
			if !metrics[k] && vt != "" {
				attr[k] = vt
			}
		case int64:
			if !metrics[k] && vt > 0 {
				attr[k] = int(vt)
			}
		case int32:
			if !metrics[k] && vt > 0 {
				attr[k] = int(vt)
			}
		}
	}

	if f.lastMetadata[in.DeviceName] != nil {
		for k, v := range f.lastMetadata[in.DeviceName].DeviceInfo {
			attr[k] = v
		}

		if in.OutputPort != in.InputPort {
			if ii, ok := f.lastMetadata[in.DeviceName].InterfaceInfo[in.InputPort]; ok {
				for k, v := range ii {
					if v != "" {
						attr["input_if_"+k] = v
					}
				}
			}
			if ii, ok := f.lastMetadata[in.DeviceName].InterfaceInfo[in.OutputPort]; ok {
				for k, v := range ii {
					if v != "" {
						attr["output_if_"+k] = v
					}
				}
			}
		} else {
			if ii, ok := f.lastMetadata[in.DeviceName].InterfaceInfo[in.OutputPort]; ok {
				for k, v := range ii {
					if v != "" {
						attr["if_"+k] = v
					}
				}
			}
		}
	}
}
