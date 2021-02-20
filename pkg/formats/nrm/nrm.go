package nrm

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

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

type LastMetadata struct {
	deviceInfo    map[string]string
	interfaceInfo map[kt.IfaceID]map[string]string
}

type NRMFormat struct {
	logger.ContextL
	compression  kt.Compression
	doGz         bool
	lastMetadata map[string]*LastMetadata
	invalids     map[string]bool
	mux          sync.RWMutex
}

type NRMetricSet struct {
	Metrics []NRMetric `json:"metrics"`
}

type NRMetric struct {
	Name       string            `json:"name"`
	Type       string            `json:"type"`
	Value      int64             `json:"value"`
	Timestamp  int64             `json:"timestamp"`
	Interval   int64             `json:"interval.ms"`
	Attributes map[string]string `json:"attributes"`
}

func NewFormat(log logger.Underlying, compression kt.Compression) (*NRMFormat, error) {
	jf := &NRMFormat{
		compression:  compression,
		ContextL:     logger.NewContextLFromUnderlying(logger.SContext{S: "nrmFormat"}, log),
		doGz:         false,
		invalids:     map[string]bool{},
		lastMetadata: map[string]*LastMetadata{},
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
	if !f.doGz {
		return json.Marshal(rolls)
	}

	serBuf := make([]byte, 0)
	buf := bytes.NewBuffer(serBuf)
	buf.Reset()
	zw, err := gzip.NewWriterLevel(buf, gzip.DefaultCompression)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(rolls)
	if err != nil {
		return nil, err
	}

	_, err = zw.Write(b)
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
	default:
		f.mux.Lock()
		defer f.mux.Unlock()
		if !f.invalids[in.EventType] {
			f.Warnf("Invalid EventType A: %s", in.EventType)
			f.invalids[in.EventType] = true
		}
	}

	return nil
}

func (f *NRMFormat) fromSnmpMetadata(in *kt.JCHF, ts int64) []NRMetric {

	if in.DeviceName == "" { // Only run if this is set.
		return nil
	}
	lm := LastMetadata{
		deviceInfo:    map[string]string{},
		interfaceInfo: map[kt.IfaceID]map[string]string{},
	}
	for k, v := range in.CustomStr {
		if strings.HasPrefix(k, "if.") {
			pts := strings.SplitN(k, ".", 3)
			if len(pts) == 3 {
				if ifint, err := strconv.Atoi(pts[1]); err == nil {
					if _, ok := lm.interfaceInfo[kt.IfaceID(ifint)]; !ok {
						lm.interfaceInfo[kt.IfaceID(ifint)] = map[string]string{}
					}
					lm.interfaceInfo[kt.IfaceID(ifint)][pts[2]] = v
				}
			}
		} else {
			lm.deviceInfo[k] = v
		}
	}

	f.lastMetadata[in.DeviceName] = &lm
	return nil
}

func (f *NRMFormat) fromKSynth(in *kt.JCHF, ts int64) []NRMetric {

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
			"Fetch Size | Ping Min RTT": "Min RTT", "Ping Max RTT": "Max RTT", "Ping Avg RTT": "Avg RTT", "Ping Std RTT": "Std RTT", "Ping Jit RTT": "Jit RTT"}
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
			"Fetch Size | Ping Min RTT": "Min RTT", "Ping Max RTT": "Max RTT", "Ping Avg RTT": "Avg RTT", "Ping Std RTT": "Std RTT", "Ping Jit RTT": "Jit RTT"}
	case 6: // Query
		metrics = map[string]bool{"Fetch Status | Ping Sent | Trace Time": true, "Fetch TTLB | Ping Lost": true}
		names = map[string]string{"Fetch Status | Ping Sent | Trace Time": "Time", "Fetch TTLB | Ping Lost": "Code"}
	case 7: // Shake
		metrics = map[string]bool{"Fetch Status | Ping Sent | Trace Time": true, "Lat/Long Dest": true}
		names = map[string]string{"Fetch Status | Ping Sent | Trace Time": "Time", "Lat/Long Dest": "Port"}
	}

	attr := map[string]string{}
	f.setAttr(attr, in, metrics)
	ms := make([]NRMetric, len(metrics))
	i := 0

	for m, _ := range metrics {
		switch m {
		case "Error", "Timeout":
			ms[i] = NRMetric{
				Name:       "kentik." + m,
				Type:       NR_GAUGE_TYPE,
				Value:      1,
				Timestamp:  ts,
				Attributes: attr,
			}
		default:
			ms[i] = NRMetric{
				Name:       "kentik." + names[m],
				Type:       NR_GAUGE_TYPE,
				Value:      int64(in.CustomInt[m]),
				Timestamp:  ts,
				Attributes: attr,
			}
		}
		i++
	}

	return ms
}

func (f *NRMFormat) fromKflow(in *kt.JCHF, ts int64) []NRMetric {
	// Map the basic strings into here.
	attr := map[string]string{}
	metrics := map[string]bool{"in_bytes": true, "out_bytes": true, "in_pkts": true, "out_pkts": true}
	f.setAttr(attr, in, metrics)
	ms := []NRMetric{
		NRMetric{
			Name:       "kentik.in_bytes",
			Type:       NR_COUNT_TYPE,
			Value:      int64(in.InBytes),
			Timestamp:  ts,
			Interval:   60000, // Can we be better here?
			Attributes: attr,
		},
		NRMetric{
			Name:       "kentik.out_bytes",
			Type:       NR_COUNT_TYPE,
			Value:      int64(in.OutBytes),
			Timestamp:  ts,
			Interval:   60000, // Can we be better here?
			Attributes: attr,
		},
		NRMetric{
			Name:       "kentik.in_pkts",
			Type:       NR_COUNT_TYPE,
			Value:      int64(in.InPkts),
			Timestamp:  ts,
			Interval:   60000, // Can we be better here?
			Attributes: attr,
		},
		NRMetric{
			Name:       "kentik.out_pkts",
			Type:       NR_COUNT_TYPE,
			Value:      int64(in.OutPkts),
			Timestamp:  ts,
			Interval:   60000, // Can we be better here?
			Attributes: attr,
		},
	}
	return ms
}

func (f *NRMFormat) fromSnmpDeviceMetric(in *kt.JCHF, ts int64) []NRMetric {
	metrics := map[string]bool{"CPU": true, "MemoryTotal": true, "MemoryUsed": true, "MemoryFree": true, "MemoryUtilization": true, "Uptime": true}
	attr := map[string]string{}
	f.setAttr(attr, in, metrics)
	ms := make([]NRMetric, len(metrics))
	i := 0
	for m, _ := range metrics {
		ms[i] = NRMetric{
			Name:       "kentik." + m,
			Type:       NR_GAUGE_TYPE,
			Value:      in.CustomBigInt[m],
			Timestamp:  ts,
			Attributes: attr,
		}
		i++
	}

	return ms
}

func (f *NRMFormat) fromSnmpInterfaceMetric(in *kt.JCHF, ts int64) []NRMetric {
	metrics := map[string]bool{"ifHCInOctets": true, "ifHCInUcastPkts": true, "ifHCOutOctets": true, "ifHCOutUcastPkts": true, "ifInErrors": true, "ifOutErrors": true,
		"ifInDiscards": true, "ifOutDiscards": true, "ifHCOutMulticastPkts": true, "ifHCOutBroadcastPkts": true, "ifHCInMulticastPkts": true, "ifHCInBroadcastPkts": true}
	attr := map[string]string{}
	f.setAttr(attr, in, metrics)
	ms := make([]NRMetric, len(metrics))
	i := 0
	for m, _ := range metrics {
		ms[i] = NRMetric{
			Name:       "kentik." + m,
			Type:       NR_GAUGE_TYPE,
			Value:      in.CustomBigInt[m],
			Timestamp:  ts,
			Attributes: attr,
		}
		i++
	}

	return ms
}

func (f *NRMFormat) setAttr(attr map[string]string, in *kt.JCHF, metrics map[string]bool) {
	mapr := in.Flatten()
	for k, v := range mapr {
		switch vt := v.(type) {
		case string:
			if !metrics[k] && vt != "" {
				attr[k] = vt
			}
		case int64:
			if !metrics[k] && vt > 0 {
				attr[k] = strconv.Itoa(int(vt))
			}
		case int32:
			if !metrics[k] && vt > 0 {
				attr[k] = strconv.Itoa(int(vt))
			}
		}
	}

	if f.lastMetadata[in.DeviceName] != nil {
		for k, v := range f.lastMetadata[in.DeviceName].deviceInfo {
			attr[k] = v
		}

		if in.OutputPort != in.InputPort {
			if ii, ok := f.lastMetadata[in.DeviceName].interfaceInfo[in.InputPort]; ok {
				for k, v := range ii {
					attr["input_if_"+k] = v
				}
			}
			if ii, ok := f.lastMetadata[in.DeviceName].interfaceInfo[in.OutputPort]; ok {
				for k, v := range ii {
					attr["output_if_"+k] = v
				}
			}
		} else {
			if ii, ok := f.lastMetadata[in.DeviceName].interfaceInfo[in.OutputPort]; ok {
				for k, v := range ii {
					attr["if_"+k] = v
				}
			}
		}
	}
}
