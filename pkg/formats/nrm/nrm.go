package nrm

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kentik/ktranslate/pkg/formats/nrm/events"
	"github.com/kentik/ktranslate/pkg/formats/util"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

const (
	NR_COUNT_TYPE   = "count"
	NR_DIST_TYPE    = "distribution"
	NR_GAUGE_TYPE   = "gauge"
	NR_SUMMARY_TYPE = "summary"
	NR_UNIQUE_TYPE  = "uniqueCount"

	DO_DEMO_PERIOD = "NRM_DO_DEMO_PERIOD"
	DO_DEMO_AMP    = "NRM_DO_DEMO_AMPLITIDE"

	InstNameVPC       = "vpc-logs"
	InstNameNetflow   = "netflow"
	InstNameSNMP      = "snmp"
	InstNameSynthetic = "synthetic"
)

type NRMFormat struct {
	logger.ContextL
	compression  kt.Compression
	doGz         bool
	lastMetadata map[string]*kt.LastMetadata
	invalids     map[string]bool
	mux          sync.RWMutex
	demo         *Demozer

	EventChan chan []byte
}

type NRCommon struct {
	Timestamp  int64             `json:"timestamp"`
	Attributes map[string]string `json:"attributes"`
}

type NRMetricSet struct {
	Metrics []NRMetric `json:"metrics"`
	Common  *NRCommon  `json:"common"`
}

type NRMetric struct {
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Value      interface{}            `json:"value,omitempty"`
	Timestamp  int64                  `json:"timestamp,omitempty"`
	Interval   int64                  `json:"interval.ms,omitempty"`
	Attributes map[string]interface{} `json:"attributes"`
}

func NewFormat(log logger.Underlying, compression kt.Compression) (*NRMFormat, error) {
	jf := &NRMFormat{
		compression:  compression,
		ContextL:     logger.NewContextLFromUnderlying(logger.SContext{S: "nrmFormat"}, log),
		doGz:         false,
		invalids:     map[string]bool{},
		lastMetadata: map[string]*kt.LastMetadata{},
		EventChan:    make(chan []byte, 100), // Used for sending events to the event API.
	}

	dp := os.Getenv(DO_DEMO_PERIOD)
	if per, err := strconv.Atoi(dp); err == nil {
		da := os.Getenv(DO_DEMO_AMP)
		if amp, err := strconv.Atoi(da); err != nil {
			return nil, fmt.Errorf("Bad value for demo amplitude: %s. Be sure to set %s", da, DO_DEMO_AMP)
		} else {
			jf.demo = NewDemozer(jf, uint32(per), uint32(amp))
			jf.Infof("Running Demo System with period of %d and amplitude of %d", per, amp)
		}
	}

	switch compression {
	case kt.CompressionNone:
		jf.doGz = false
	case kt.CompressionGzip:
		jf.doGz = true
	default:
		return nil, fmt.Errorf("Invalid compression (%s): format nrm only supports none|gzip", compression)
	}

	return jf, nil
}

func (f *NRMFormat) To(msgs []*kt.JCHF, serBuf []byte) (*kt.Output, error) {
	ms := NRMetricSet{
		Metrics: make([]NRMetric, 0, len(msgs)*4),
		Common:  newNRCommon(),
	}
	for _, m := range msgs {
		ms.Metrics = append(ms.Metrics, f.toNRMetric(m)...)
	}

	if len(ms.Metrics) == 0 {
		return nil, nil
	}

	target, err := json.Marshal([]NRMetricSet{ms}) // Has to be an array here, no idea why.
	if err != nil {
		return nil, err
	}

	if !f.doGz {
		return kt.NewOutputWithProvider(target, msgs[0].Provider, kt.MetricOutput), nil
	}

	buf := bytes.NewBuffer(serBuf)
	buf.Reset()
	zw, err := gzip.NewWriterLevel(buf, gzip.DefaultCompression)
	if err != nil {
		return nil, err
	}

	_, err = zw.Write(target)
	if err != nil {
		return nil, err
	}

	err = zw.Close()
	if err != nil {
		return nil, err
	}

	return kt.NewOutputWithProvider(buf.Bytes(), msgs[0].Provider, kt.MetricOutput), nil
}

func (f *NRMFormat) From(raw *kt.Output) ([]map[string]interface{}, error) {
	values := make([]map[string]interface{}, 0)
	return values, nil
}

func (f *NRMFormat) Rollup(rolls []rollup.Rollup) (*kt.Output, error) {
	ms := NRMetricSet{
		Metrics: f.toNRMetricRollup(rolls),
		Common:  newNRCommon(),
	}

	if len(ms.Metrics) == 0 {
		return nil, nil
	}

	target, err := json.Marshal([]NRMetricSet{ms}) // Has to be an array here, no idea why.
	if err != nil {
		return nil, err
	}

	if !f.doGz {
		return kt.NewOutputWithProvider(target, rolls[0].Provider, kt.RollupOutput), nil
	}

	serBuf := make([]byte, 0)
	buf := bytes.NewBuffer(serBuf)
	buf.Reset()
	zw, err := gzip.NewWriterLevel(buf, gzip.DefaultCompression)
	if err != nil {
		return nil, err
	}

	_, err = zw.Write(target)
	if err != nil {
		return nil, err
	}

	err = zw.Close()
	if err != nil {
		return nil, err
	}

	return kt.NewOutputWithProvider(buf.Bytes(), rolls[0].Provider, kt.RollupOutput), nil
}

func (f *NRMFormat) toNRMetric(in *kt.JCHF) []NRMetric {
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
	case kt.KENTIK_EVENT_SNMP_TRAP:
		// This is actually an event, send out as an event to sink directly.
		err := events.SendEvent(in, f.doGz, f.EventChan)
		if err != nil {
			f.Errorf("Cannot send event on -- %v", err)
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
	rollupMap = map[string]string{
		"source_account": "vpc_account",
		"source_region":  "vpc_region",
		"source_vpc":     "vpc_name",
		"src_addr":       "ip_address",
		"src_as":         "asn_name",
		"src_as_name":    "asn_name",
		"src_geo":        "country",
		"l4_src_port":    "port",
		"dest_account":   "vpc_account",
		"dest_region":    "vpc_region",
		"dest_vpc":       "vpc_name",
		"dst_addr":       "ip_address",
		"dst_as":         "asn_name",
		"dst_as_name":    "asn_name",
		"dst_geo":        "country",
		"l4_dst_port":    "port",
	}
)

func (f *NRMFormat) toNRMetricRollup(in []rollup.Rollup) []NRMetric {
	ms := make([]NRMetric, 0, len(in))
	for _, roll := range in {
		if roll.Metric == 0 {
			continue
		}

		dims := roll.GetDims()
		attr := map[string]interface{}{
			"provider":             roll.Provider,
			"instrumentation.name": toInstName(roll.Provider),
		}

		// Override here for router to map to flowdevice
		if attr["provider"].(kt.Provider) == kt.ProviderRouter {
			attr["provider"] = kt.ProviderFlowDevice
		}

		for i, pt := range strings.Split(roll.Dimension, roll.KeyJoin) {
			aname := dims[i]
			if n, ok := rollupMap[aname]; ok {
				aname = n
			}
			attr[aname] = pt
			if pt == "0" || pt == "" || pt == "-" || pt == "--" {
				delete(attr, aname)
			}
			if aname == "port" { // Remap efemeral ports down here.
				port, _ := strconv.Atoi(pt)
				if port > 32768 {
					attr[aname] = 32768
				}
			}
		}

		// Finally, combine vpc_account:vpc_name pre-shipping and call it vpc_identification
		acct, istra := attr["vpc_account"].(string)
		name, istrn := attr["vpc_name"].(string)
		if istra && istrn && acct != "" && name != "" {
			attr["vpc_identification"] = acct + ":" + name
		}

		ms = append(ms, NRMetric{
			Name: "kentik.rollup." + roll.Name,
			Type: NR_SUMMARY_TYPE,
			Value: map[string]uint64{
				"count": roll.Count,
				"sum":   uint64(roll.Metric),
				"min":   roll.Min,
				"max":   roll.Max,
			},
			Interval:   roll.Interval.Microseconds(),
			Attributes: attr,
		})
	}

	// Tweak any values if we have a demoizer set.
	if f.demo != nil {
		f.demo.demoize(ms)
	}

	return ms
}

func (f *NRMFormat) fromSnmpMetadata(in *kt.JCHF) []NRMetric {
	if in.DeviceName == "" { // Only run if this is set.
		return nil
	}

	lm := util.SetMetadata(in)

	f.mux.Lock()
	defer f.mux.Unlock()
	f.lastMetadata[in.DeviceName] = lm

	return nil
}

var (
	synthWLAttr = map[string]bool{
		"agent_id":    true,
		"agent_name":  true,
		"dst_addr":    true,
		"dst_cdn_int": true,
		"dst_geo":     true,
		"provider":    true,
		"src_addr":    true,
		"src_cdn_int": true,
		"src_geo":     true,
		"test_id":     true,
		"test_name":   true,
		"test_type":   true,
		"test_url":    true,
	}
)

func (f *NRMFormat) fromKSynth(in *kt.JCHF) []NRMetric {
	if in.CustomInt["result_type"] <= 1 {
		return nil // Don't worry about timeouts and errrors for now.
	}

	metrics := util.GetSynMetricNameSet(in.CustomInt["result_type"])
	attr := map[string]interface{}{}
	f.mux.RLock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName])
	f.mux.RUnlock()
	ms := make([]NRMetric, 0, len(metrics))
	lost := 0.0
	sent := 0.0

	// Hard code these.
	attr["instrumentation.name"] = InstNameSynthetic

	for k, _ := range attr { // White list only a few attributes here.
		if !synthWLAttr[k] {
			delete(attr, k)
		} else {
			switch tr := attr[k].(type) {
			case int:
				// Force this to string.
				attr[k] = strconv.Itoa(tr)
			case int32:
				// Force this to string.
				attr[k] = strconv.Itoa(int(tr))
			case int64:
				// Force this to string.
				attr[k] = strconv.Itoa(int(tr))
			default:
				// noop.
			}
		}
	}

	for m, name := range metrics {
		switch name {
		case "avg_rtt", "jit_rtt":
			ms = append(ms, NRMetric{
				Name:       "kentik.synth." + name,
				Type:       NR_GAUGE_TYPE,
				Value:      int64(in.CustomInt[m]),
				Attributes: attr,
			})
		case "lost":
			lost = float64(in.CustomInt[m])
		case "sent":
			sent = float64(in.CustomInt[m])
		}
	}

	if sent > 0 {
		ms = append(ms, NRMetric{
			Name:       "kentik.synth.lost_pct",
			Type:       NR_GAUGE_TYPE,
			Value:      (lost / sent) * 100.,
			Attributes: attr,
		})
	}

	return ms
}

func (f *NRMFormat) fromKflow(in *kt.JCHF) []NRMetric {
	// Map the basic strings into here.
	attr := map[string]interface{}{}
	metrics := map[string]string{"in_bytes": "", "out_bytes": "", "in_pkts": "", "out_pkts": "", "latency_ms": ""}
	f.mux.RLock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName])
	f.mux.RUnlock()
	ms := make([]NRMetric, 0)

	// Hard code these.
	attr["instrumentation.name"] = InstNameNetflow

	for m, _ := range metrics {
		var value int64
		switch m {
		case "in_bytes":
			value = int64(in.InBytes * uint64(in.SampleRate))
		case "out_bytes":
			value = int64(in.OutBytes * uint64(in.SampleRate))
		case "in_pkts":
			value = int64(in.InPkts * uint64(in.SampleRate))
		case "out_pkts":
			value = int64(in.OutPkts * uint64(in.SampleRate))
		case "latency_ms":
			value = int64(in.CustomInt["appl_latency_ms"])
		}
		if value > 0 {
			ms = append(ms, NRMetric{
				Name:       "kentik.flow." + m,
				Type:       NR_GAUGE_TYPE,
				Value:      value,
				Attributes: attr,
			})
		}
	}
	return ms
}

func (f *NRMFormat) fromSnmpDeviceMetric(in *kt.JCHF) []NRMetric {
	metrics := in.CustomMetrics
	attr := map[string]interface{}{}
	f.mux.RLock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName])
	f.mux.RUnlock()
	ms := make([]NRMetric, 0, len(metrics))
	for m, name := range metrics {
		if _, ok := in.CustomBigInt[m]; ok {
			attrNew := copyAttrForSnmp(attr, name)
			ms = append(ms, NRMetric{
				Name:       "kentik.snmp." + m,
				Type:       NR_GAUGE_TYPE,
				Value:      int64(in.CustomBigInt[m]),
				Attributes: attrNew,
			})
		}
	}

	return ms
}

func (f *NRMFormat) fromSnmpInterfaceMetric(in *kt.JCHF) []NRMetric {
	metrics := in.CustomMetrics
	attr := map[string]interface{}{}
	f.mux.RLock()
	defer f.mux.RUnlock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName])
	ms := make([]NRMetric, 0, len(metrics))
	for m, name := range metrics {
		if _, ok := in.CustomBigInt[m]; ok {
			attrNew := copyAttrForSnmp(attr, name)
			ms = append(ms, NRMetric{
				Name:       "kentik.snmp." + m,
				Type:       NR_GAUGE_TYPE,
				Value:      int64(in.CustomBigInt[m]),
				Attributes: attrNew,
			})
		}
	}

	// Grap capacity utilization if possible.
	if f.lastMetadata[in.DeviceName] != nil {
		if ii, ok := f.lastMetadata[in.DeviceName].InterfaceInfo[in.InputPort]; ok {
			if speed, ok := ii["Speed"]; ok {
				if ispeed, ok := speed.(int32); ok {
					uptimeSpeed := in.CustomBigInt["Uptime"] * (int64(ispeed) * 1000000) // Convert into bits here, from megabits.
					if uptimeSpeed > 0 {
						ms = append(ms, NRMetric{
							Name:       "kentik.snmp.IfInUtilization",
							Type:       NR_GAUGE_TYPE,
							Value:      float64(in.CustomBigInt["ifHCInOctets"]*8*100) / float64(uptimeSpeed),
							Attributes: attr,
						})
					}
				}
			}
		}
		if oi, ok := f.lastMetadata[in.DeviceName].InterfaceInfo[in.OutputPort]; ok {
			if speed, ok := oi["Speed"]; ok {
				if ispeed, ok := speed.(int32); ok {
					uptimeSpeed := in.CustomBigInt["Uptime"] * (int64(ispeed) * 1000000) // Convert into bits here, from megabits.
					if uptimeSpeed > 0 {
						ms = append(ms, NRMetric{
							Name:       "kentik.snmp.IfOutUtilization",
							Type:       NR_GAUGE_TYPE,
							Value:      float64(in.CustomBigInt["ifHCOutOctets"]*8*100) / float64(uptimeSpeed),
							Attributes: attr,
						})
					}
				}
			}
		}
	}

	return ms
}

func (f *NRMFormat) fromKtranslate(in *kt.JCHF) []NRMetric {
	// Map the basic strings into here.
	attr := map[string]interface{}{}
	metrics := map[string]string{"name": "", "value": "", "count": "", "one-minute": "", "95-percentile": "", "du": ""}
	f.mux.RLock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName])
	f.mux.RUnlock()
	ms := make([]NRMetric, 0)

	// Hard code these.
	attr["instrumentation.name"] = InstNameNetflow

	switch in.CustomStr["type"] {
	case "counter":
		if in.CustomBigInt["count"] > 0 {
			ms = append(ms, NRMetric{
				Name:       "kentik.ktranslate." + in.CustomStr["name"],
				Type:       NR_GAUGE_TYPE,
				Value:      in.CustomBigInt["count"],
				Attributes: attr,
			})
		}
	case "gauge":
		if in.CustomBigInt["value"] > 0 {
			ms = append(ms, NRMetric{
				Name:       "kentik.ktranslate." + in.CustomStr["name"],
				Type:       NR_GAUGE_TYPE,
				Value:      in.CustomBigInt["value"],
				Attributes: attr,
			})
		}
	case "histogram":
		if in.CustomBigInt["95-percentile"] > 0 {
			ms = append(ms, NRMetric{
				Name:       "kentik.ktranslate." + in.CustomStr["name"],
				Type:       NR_GAUGE_TYPE,
				Value:      in.CustomBigInt["95-percentile"],
				Attributes: attr,
			})
		}
	case "meter":
		if in.CustomBigInt["one-minute"] > 0 {
			ms = append(ms, NRMetric{
				Name:       "kentik.ktranslate." + in.CustomStr["name"],
				Type:       NR_GAUGE_TYPE,
				Value:      in.CustomBigInt["one-minute"],
				Attributes: attr,
			})
		}
	case "timer":
		if in.CustomBigInt["value"] > 0 {
			ms = append(ms, NRMetric{
				Name:       "kentik.ktranslate." + in.CustomStr["95-percentile"],
				Type:       NR_GAUGE_TYPE,
				Value:      in.CustomBigInt["95-percentile"],
				Attributes: attr,
			})
		}
	}
	return ms
}

func copyAttrForSnmp(attr map[string]interface{}, name string) map[string]interface{} {
	attrNew := map[string]interface{}{
		"objectIdentifier":     name,
		"instrumentation.name": InstNameSNMP,
	}
	for k, v := range attr {
		attrNew[k] = v
	}
	return attrNew
}

func toInstName(prov kt.Provider) string {
	if strings.Contains(string(prov), "vpc") {
		return InstNameVPC
	}
	return InstNameNetflow
}

func newNRCommon() *NRCommon {
	return &NRCommon{
		Timestamp: time.Now().UnixNano() / 1e+6, // Convert to milliseconds
		Attributes: map[string]string{
			"instrumentation.provider": kt.InstProvider,
			"collector.name":           kt.CollectorName,
		},
	}
}
