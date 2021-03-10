package prom

import (
	"flag"
	"fmt"
	"strings"
	"sync"

	"github.com/kentik/ktranslate/pkg/formats/util"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"

	"github.com/kentik/ktranslate/pkg/eggs/logger"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	doCollectorStats = flag.Bool("info_collector", false, "Also send stats about this collector")
)

type PromData struct {
	Name      string
	Value     float64
	Tags      map[string]interface{}
	Timestamp int64
}

func (d *PromData) GetTagLabels(vecTags map[string]map[string]int) []string {
	if _, ok := vecTags[d.Name]; !ok {
		vecTags[d.Name] = map[string]int{}
	}
	i := 0
	tags := make([]string, len(d.Tags))
	for k, _ := range d.Tags {
		vecTags[d.Name][k] = i
		tags[i] = k
		i++
	}
	return tags
}

func (d *PromData) GetTagValues(vecTags map[string]map[string]int) []string {
	tags := make([]string, len(vecTags[d.Name]))
	for k, v := range d.Tags {
		posit, ok := vecTags[d.Name][k]
		if !ok {
			continue
		}
		switch t := v.(type) {
		case string:
			tags[posit] = t
		default:
			tags[posit] = fmt.Sprintf("%v", v)
		}
	}
	return tags
}

type PromFormat struct {
	logger.ContextL
	vecs         map[string]*prometheus.CounterVec
	invalids     map[string]bool
	lastMetadata map[string]*kt.LastMetadata
	vecTags      map[string]map[string]int
	mux          sync.RWMutex
}

func NewFormat(log logger.Underlying, compression kt.Compression) (*PromFormat, error) {
	jf := &PromFormat{
		ContextL:     logger.NewContextLFromUnderlying(logger.SContext{S: "influxFormat"}, log),
		vecs:         make(map[string]*prometheus.CounterVec),
		invalids:     map[string]bool{},
		lastMetadata: map[string]*kt.LastMetadata{},
		vecTags:      map[string]map[string]int{},
	}

	if *doCollectorStats {
		prometheus.MustRegister(prometheus.NewBuildInfoCollector())
	}

	return jf, nil
}

// Not supported.
func (f *PromFormat) To(msgs []*kt.JCHF, serBuf []byte) ([]byte, error) {
	res := make([]PromData, 0, len(msgs))
	for _, m := range msgs {
		res = append(res, f.toPromMetric(m)...)
	}

	if len(res) == 0 {
		return nil, nil
	}

	for _, m := range res {
		if _, ok := f.vecs[m.Name]; !ok {
			f.mux.Lock()
			f.vecs[m.Name] = prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: m.Name,
				},
				m.GetTagLabels(f.vecTags),
			)
			f.mux.Unlock()
			prometheus.MustRegister(f.vecs[m.Name])
		}
		f.mux.RLock()
		f.vecs[m.Name].WithLabelValues(m.GetTagValues(f.vecTags)...).Add(m.Value)
		f.mux.RUnlock()
	}

	return nil, nil
}

// Not supported.
func (f *PromFormat) From(raw []byte) ([]map[string]interface{}, error) {
	values := make([]map[string]interface{}, 0)
	return values, nil
}

func (f *PromFormat) Rollup(rolls []rollup.Rollup) ([]byte, error) {
	for _, r := range rolls {
		pkts := strings.Split(r.EventType, ":")
		if _, ok := f.vecs[r.EventType]; !ok {
			f.vecs[r.EventType] = prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: strings.Join(pkts[0:2], ":"),
				},
				pkts[2:],
			)
			prometheus.MustRegister(f.vecs[r.EventType])
		}
		f.vecs[r.EventType].WithLabelValues(strings.Split(r.Dimension, r.KeyJoin)...).Add(float64(r.Metric))
	}

	return nil, nil
}

func (f *PromFormat) toPromMetric(in *kt.JCHF) []PromData {
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

func (f *PromFormat) fromKSynth(in *kt.JCHF) []PromData {
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
	f.mux.RLock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName])
	f.mux.RUnlock()
	ms := map[string]int64{}

	for m, _ := range metrics {
		switch m {
		case "Error", "Timeout":
			ms[m] = 1
		default:
			ms[names[m]] = int64(in.CustomInt[m])
		}
	}

	res := []PromData{}
	for k, v := range ms {
		res = append(res, PromData{
			Name:      "kentik:synth:" + k,
			Value:     float64(v),
			Timestamp: in.Timestamp * 1000000000,
			Tags:      attr,
		})
	}

	return res
}

func (f *PromFormat) fromKflow(in *kt.JCHF) []PromData {
	// Map the basic strings into here.
	attr := map[string]interface{}{}
	metrics := map[string]bool{"in_bytes": true, "out_bytes": true, "in_pkts": true, "out_pkts": true, "latency_ms": true}
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
			ms[m] = int64(in.CustomInt["APPL_LATENCY_MS"])
		}
	}

	res := []PromData{}
	for k, v := range ms {
		res = append(res, PromData{
			Name:      "kentik:flow:" + k,
			Value:     float64(v),
			Timestamp: in.Timestamp * 1000000000,
			Tags:      attr,
		})
	}

	return res
}

func (f *PromFormat) fromSnmpDeviceMetric(in *kt.JCHF) []PromData {
	var metrics map[string]bool
	if len(in.CustomMetrics) > 0 {
		metrics = in.CustomMetrics
	} else {
		metrics = map[string]bool{"CPU": true, "MemoryTotal": true, "MemoryUsed": true, "MemoryFree": true, "MemoryUtilization": true, "Uptime": true}
	}
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

	res := []PromData{}
	for k, v := range ms {
		res = append(res, PromData{
			Name:      "kentik:snmp:" + k,
			Value:     float64(v),
			Timestamp: in.Timestamp * 1000000000,
			Tags:      attr,
		})
	}

	return res
}

func (f *PromFormat) fromSnmpInterfaceMetric(in *kt.JCHF) []PromData {
	var metrics map[string]bool
	if len(in.CustomMetrics) > 0 {
		metrics = in.CustomMetrics
	} else {
		metrics = map[string]bool{"ifHCInOctets": true, "ifHCInUcastPkts": true, "ifHCOutOctets": true, "ifHCOutUcastPkts": true, "ifInErrors": true, "ifOutErrors": true,
			"ifInDiscards": true, "ifOutDiscards": true, "ifHCOutMulticastPkts": true, "ifHCOutBroadcastPkts": true, "ifHCInMulticastPkts": true, "ifHCInBroadcastPkts": true}
	}
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

	res := []PromData{}
	for k, v := range ms {
		res = append(res, PromData{
			Name:      "kentik:snmp:" + k,
			Value:     float64(v),
			Timestamp: in.Timestamp * 1000000000,
			Tags:      attr,
		})
	}
	for k, v := range msF {
		res = append(res, PromData{
			Name:      "kentik:snmp:" + k,
			Value:     v,
			Timestamp: in.Timestamp * 1000000000,
			Tags:      attr,
		})
	}

	return res
}

func (f *PromFormat) fromSnmpMetadata(in *kt.JCHF) []PromData {
	if in.DeviceName == "" { // Only run if this is set.
		return nil
	}

	lm := util.SetMetadata(in)

	f.mux.Lock()
	defer f.mux.Unlock()
	f.lastMetadata[in.DeviceName] = lm

	return nil
}
