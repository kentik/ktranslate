package prom

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
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
	invalidTag       = regexp.MustCompile(`^\d+$`)
)

func init() {
	flag.BoolVar(&doCollectorStats, "info_collector", false, "Also send stats about this collector")
	flag.IntVar(&seenNeeded, "prom_seen", 4, "Number of flows needed inbound before we start writting to the collector")

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
		if invalidTag.MatchString(k) {
			continue
		}
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
	vecs         map[string]*prometheus.GaugeVec
	invalids     map[string]bool
	lastMetadata map[string]*kt.LastMetadata
	vecTags      tagVec
	seen         map[string]int
	config       *ktranslate.PrometheusFormatConfig

	mux sync.RWMutex
}

func NewFormat(log logger.Underlying, compression kt.Compression, cfg *ktranslate.PrometheusFormatConfig) (*PromFormat, error) {
	if cfg == nil {
		return nil, fmt.Errorf("prometheus format cannot be nil")
	}
	jf := &PromFormat{
		ContextL:     logger.NewContextLFromUnderlying(logger.SContext{S: "promFormat"}, log),
		vecs:         make(map[string]*prometheus.GaugeVec),
		invalids:     map[string]bool{},
		lastMetadata: map[string]*kt.LastMetadata{},
		vecTags:      map[string]map[string]int{},
		config:       cfg,
		seen:         map[string]int{},
	}

	if cfg.EnableCollectorStats {
		prometheus.MustRegister(prometheus.NewBuildInfoCollector())
	}

	return jf, nil
}

func (f *PromFormat) toLabels(name string) []string {
	// Map out any chars which prom can't handle as label names.
	makeSafe := func(r rune) rune {
		switch {
		case r == '.':
			return '_'
		case r == ' ':
			return '_'
		}
		return r
	}

	res := make([]string, len(f.vecTags[name]))
	for k, v := range f.vecTags[name] {
		res[v] = strings.Map(makeSafe, k)
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

		if _, ok := f.vecs[m.Name]; !ok {
			labels := f.toLabels(m.Name)
			cv := prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: m.Name,
				},
				labels,
			)
			prometheus.MustRegister(cv)
			f.vecs[m.Name] = cv
			f.Infof("Adding %s %v", m.Name, labels)
		}
		f.vecs[m.Name].WithLabelValues(m.GetTagValues(f.vecTags)...).Set(m.Value)
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
			f.vecs[roll.EventType] = prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: strings.ReplaceAll(roll.Name, ".", ":"),
				},
				roll.GetDims(),
			)
			prometheus.MustRegister(f.vecs[roll.EventType])
		}
		f.vecs[roll.EventType].WithLabelValues(strings.Split(roll.Dimension, roll.KeyJoin)...).Set(float64(roll.Metric))
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
	case kt.KENTIK_EVENT_SYNTH_GEST:
		return f.fromKSyngest(in)
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

func (f *PromFormat) fromKSyngest(in *kt.JCHF) []PromData {
	metrics := util.GetSyngestMetricNameSet()
	attr := map[string]interface{}{}
	f.mux.RLock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName], false)
	f.mux.RUnlock()
	ms := make([]PromData, 0, len(metrics))

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
			ms = append(ms, PromData{
				Name:  "kentik:syngest:" + name.Name,
				Value: float64(in.CustomInt[m]),
				Tags:  attr,
			})
		}
	}

	return ms
}

func (f *PromFormat) fromKSynth(in *kt.JCHF) []PromData {
	if in.CustomInt["result_type"] <= 1 {
		return nil // Don't worry about timeouts and errors for now.
	}

	rawStr := in.CustomStr["error_cause/trace_route"] // Pull this out early.
	metrics := util.GetSynMetricNameSet(in.CustomInt["result_type"])
	attr := map[string]interface{}{}
	f.mux.RLock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName], false)
	f.mux.RUnlock()
	ms := make([]PromData, 0, len(metrics))

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
			ms = append(ms, PromData{
				Name:  "kentik:synth:" + name.Name,
				Value: float64(in.CustomInt[m]),
				Tags:  attr,
			})
		}
	}

	return ms
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
