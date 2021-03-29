package nrm

import (
	"bytes"
	"compress/gzip"
	"fmt"
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
)

type NRMFormat struct {
	logger.ContextL
	compression  kt.Compression
	doGz         bool
	lastMetadata map[string]*kt.LastMetadata
	invalids     map[string]bool
	mux          sync.RWMutex

	EventChan chan []byte
}

type NRMetricSet struct {
	Metrics []NRMetric `json:"metrics"`
}

type NRMetric struct {
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Value      interface{}            `json:"value"`
	Timestamp  int64                  `json:"timestamp"`
	Interval   int64                  `json:"interval.ms"`
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

func (f *NRMFormat) To(msgs []*kt.JCHF, serBuf []byte) ([]byte, error) {
	ms := NRMetricSet{
		Metrics: make([]NRMetric, 0, len(msgs)*4),
	}
	ct := time.Now().UnixNano() / 1e+6 // Convert to milliseconds
	for _, m := range msgs {
		ms.Metrics = append(ms.Metrics, f.toNRMetric(m, ct)...)
	}

	if len(ms.Metrics) == 0 {
		return nil, nil
	}

	target, err := json.Marshal([]NRMetricSet{ms}) // Has to be an array here, no idea why.
	if err != nil {
		return nil, err
	}

	if !f.doGz {
		return target, nil
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

	return buf.Bytes(), nil
}

func (f *NRMFormat) From(raw []byte) ([]map[string]interface{}, error) {
	values := make([]map[string]interface{}, 0)
	return values, nil
}

func (f *NRMFormat) Rollup(rolls []rollup.Rollup) ([]byte, error) {
	ct := time.Now().UnixNano() / 1e+6 // Convert to milliseconds
	ms := NRMetricSet{
		Metrics: f.toNRMetricRollup(rolls, ct),
	}

	if len(ms.Metrics) == 0 {
		return nil, nil
	}

	target, err := json.Marshal([]NRMetricSet{ms}) // Has to be an array here, no idea why.
	if err != nil {
		return nil, err
	}

	if !f.doGz {
		return target, nil
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

	return buf.Bytes(), nil
}

func (f *NRMFormat) toNRMetric(in *kt.JCHF, ts int64) []NRMetric {
	switch in.EventType {
	case kt.KENTIK_EVENT_TYPE:
		return f.fromKflow(in, ts)
	case kt.KENTIK_EVENT_SNMP_DEV_METRIC:
		return f.fromSnmpDeviceMetric(in, ts)
	case kt.KENTIK_EVENT_SNMP_INT_METRIC:
		return f.fromSnmpInterfaceMetric(in, ts)
	case kt.KENTIK_EVENT_SYNTH:
		return f.fromKSynth(in, ts)
	case kt.KENTIK_EVENT_SNMP_METADATA:
		return f.fromSnmpMetadata(in, ts)
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

func (f *NRMFormat) toNRMetricRollup(in []rollup.Rollup, ts int64) []NRMetric {
	ms := make([]NRMetric, 0, len(in))
	for _, roll := range in {
		dims := roll.GetDims()
		attr := map[string]interface{}{
			"provider": kt.ProviderRouter,
		}
		for i, pt := range strings.Split(roll.Dimension, roll.KeyJoin) {
			attr[dims[i]] = pt
		}
		ptsm := strings.Split(roll.EventType, ":")
		ms = append(ms, NRMetric{
			Name:       "kentik.rollup." + ptsm[1],
			Type:       NR_GAUGE_TYPE,
			Value:      int64(roll.Metric),
			Interval:   roll.Interval.Microseconds(),
			Timestamp:  ts,
			Attributes: attr,
		})
	}
	return ms
}

func (f *NRMFormat) fromSnmpMetadata(in *kt.JCHF, ts int64) []NRMetric {
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

func (f *NRMFormat) fromKSynth(in *kt.JCHF, ts int64) []NRMetric {
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

	for k, _ := range attr { // White list only a few attributes here.
		if !synthWLAttr[k] {
			delete(attr, k)
		}
	}

	for m, name := range metrics {
		switch name {
		case "avg_rtt", "jit_rtt":
			ms = append(ms, NRMetric{
				Name:       "kentik.synth." + name,
				Type:       NR_GAUGE_TYPE,
				Value:      int64(in.CustomInt[m]),
				Timestamp:  ts,
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
			Timestamp:  ts,
			Attributes: attr,
		})
	}

	return ms
}

func (f *NRMFormat) fromKflow(in *kt.JCHF, ts int64) []NRMetric {
	// Map the basic strings into here.
	attr := map[string]interface{}{}
	metrics := map[string]string{"in_bytes": "", "out_bytes": "", "in_pkts": "", "out_pkts": "", "latency_ms": ""}
	f.mux.RLock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName])
	f.mux.RUnlock()
	ms := make([]NRMetric, 0)
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
				Timestamp:  ts,
				Attributes: attr,
			})
		}
	}
	return ms
}

func (f *NRMFormat) fromSnmpDeviceMetric(in *kt.JCHF, ts int64) []NRMetric {
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
				Timestamp:  ts,
				Attributes: attrNew,
			})
		}
	}

	return ms
}

func (f *NRMFormat) fromSnmpInterfaceMetric(in *kt.JCHF, ts int64) []NRMetric {
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
				Timestamp:  ts,
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
							Timestamp:  ts,
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
							Timestamp:  ts,
							Attributes: attr,
						})
					}
				}
			}
		}
	}

	return ms
}

func copyAttrForSnmp(attr map[string]interface{}, name string) map[string]interface{} {
	attrNew := map[string]interface{}{
		"objectIdentifier": name,
	}
	for k, v := range attr {
		attrNew[k] = v
	}
	return attrNew
}
