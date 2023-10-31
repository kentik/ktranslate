package prom

import (
	"flag"
	"fmt"
	"strings"
	"sync"

	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/formats/util"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"

	"github.com/kentik/ktranslate/pkg/eggs/logger"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	doCollectorStats bool
	seenNeeded       int
)

func init() {
	flag.BoolVar(&doCollectorStats, "info_collector", false, "Also send stats about this collector")
	flag.IntVar(&seenNeeded, "prom_seen", 10, "Number of flows needed inbound before we start writting to the collector")

}

type PromData struct {
	Name  string
	Value float64
	Tags  map[string]interface{}
}

func (d *PromData) AddTagLabels(vecTags tagVec) {
	if _, ok := vecTags[d.Name]; !ok {
		vecTags[d.Name] = map[string]int{}
	}
	next := len(vecTags[d.Name])
	for k, _ := range d.Tags {
		if _, ok := vecTags[d.Name][k]; !ok {
			vecTags[d.Name][k] = next
			next++
		}
	}
}

func (d *PromData) GetTagValues(vecTags tagVec) []string {
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

type tagVec map[string]map[string]int

type PromFormat struct {
	logger.ContextL
	vecs         map[string]*prometheus.CounterVec
	invalids     map[string]bool
	lastMetadata map[string]*kt.LastMetadata
	vecTags      tagVec
	seen         int
	config       *ktranslate.PrometheusFormatConfig

	mux sync.RWMutex
}

func NewFormat(log logger.Underlying, compression kt.Compression, cfg *ktranslate.PrometheusFormatConfig) (*PromFormat, error) {
	if cfg == nil {
		return nil, fmt.Errorf("prometheus format cannot be nil")
	}
	jf := &PromFormat{
		ContextL:     logger.NewContextLFromUnderlying(logger.SContext{S: "promFormat"}, log),
		vecs:         make(map[string]*prometheus.CounterVec),
		invalids:     map[string]bool{},
		lastMetadata: map[string]*kt.LastMetadata{},
		vecTags:      map[string]map[string]int{},
		config:       cfg,
	}

	if cfg.EnableCollectorStats {
		prometheus.MustRegister(prometheus.NewBuildInfoCollector())
	}

	return jf, nil
}

func (f *PromFormat) addLabels(res []PromData) {
	for _, m := range res {
		m.AddTagLabels(f.vecTags)
	}
}

func (f *PromFormat) toLabels(name string) []string {
	res := make([]string, len(f.vecTags[name]))
	for k, v := range f.vecTags[name] {
		res[v] = strings.ReplaceAll(k, " ", "_")
	}
	return res
}

func (f *PromFormat) To(msgs []*kt.JCHF, serBuf []byte) (*kt.Output, error) {
	res := make([]PromData, 0, len(msgs))
	for _, m := range msgs {
		res = append(res, f.toPromMetric(m)...)
	}

	if len(res) == 0 {
		return nil, nil
	}

	f.mux.Lock()
	defer f.mux.Unlock()

	if f.seen < f.config.FlowsNeeded {
		f.addLabels(res)
		f.seen++
		if f.seen == f.config.FlowsNeeded {
			f.Infof("Seen enough!")
		} else {
			f.Infof("Seen %d", f.seen)
		}
		return nil, nil
	}

	for _, m := range res {
		if _, ok := f.vecs[m.Name]; !ok {
			labels := f.toLabels(m.Name)
			cv := prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: m.Name,
				},
				labels,
			)
			prometheus.MustRegister(cv)
			f.vecs[m.Name] = cv
			f.Infof("Adding %s %v", m.Name, labels)
		}
		f.vecs[m.Name].WithLabelValues(m.GetTagValues(f.vecTags)...).Add(m.Value)
	}

	return nil, nil
}

// Not supported.
func (f *PromFormat) From(raw *kt.Output) ([]map[string]interface{}, error) {
	values := make([]map[string]interface{}, 0)
	return values, nil
}

func (f *PromFormat) Rollup(rolls []rollup.Rollup) (*kt.Output, error) {
	for _, roll := range rolls {
		if roll.Metric == 0 {
			continue
		}
		if _, ok := f.vecs[roll.EventType]; !ok {
			f.vecs[roll.EventType] = prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: strings.ReplaceAll(roll.Name, ".", ":"),
				},
				roll.GetDims(),
			)
			prometheus.MustRegister(f.vecs[roll.EventType])
		}
		f.vecs[roll.EventType].WithLabelValues(strings.Split(roll.Dimension, roll.KeyJoin)...).Add(float64(roll.Metric))
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
	metrics := util.GetSynMetricNameSet(in.CustomInt["result_type"])
	attr := map[string]interface{}{}
	f.mux.RLock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName], false)
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

	res := []PromData{}
	for k, v := range ms {
		res = append(res, PromData{
			Name:  "kentik:synth:" + k,
			Value: float64(v),
			Tags:  attr,
		})
	}

	return res
}

func (f *PromFormat) fromKflow(in *kt.JCHF) []PromData {
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

	res := []PromData{}
	for k, v := range ms {
		if v == 0 { // Drop zero valued metrics here.
			continue
		}
		res = append(res, PromData{
			Name:  "kentik:flow:" + k,
			Value: float64(v),
			Tags:  attr,
		})
	}

	return res
}

func (f *PromFormat) fromSnmpDeviceMetric(in *kt.JCHF) []PromData {
	metrics := in.CustomMetrics
	attr := map[string]interface{}{}
	f.mux.RLock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName], false)
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
			Name:  "kentik:snmp:" + k,
			Value: float64(v),
			Tags:  attr,
		})
	}

	return res
}

func (f *PromFormat) fromSnmpInterfaceMetric(in *kt.JCHF) []PromData {
	metrics := in.CustomMetrics
	attr := map[string]interface{}{}
	f.mux.RLock()
	defer f.mux.RUnlock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName], false)
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
			Name:  "kentik:snmp:" + k,
			Value: float64(v),
			Tags:  attr,
		})
	}
	for k, v := range msF {
		res = append(res, PromData{
			Name:  "kentik:snmp:" + k,
			Value: v,
			Tags:  attr,
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
