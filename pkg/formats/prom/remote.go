package prom

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	ptime "github.com/prometheus/prometheus/model/timestamp"
	"github.com/prometheus/prometheus/prompb"

	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats/util"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"
)

type RemotePromFormat struct {
	logger.ContextL
	compression  kt.Compression
	doSnappy     bool
	lastMetadata map[string]*kt.LastMetadata
	invalids     map[string]bool
	config       *ktranslate.PrometheusFormatConfig
	seenInvalid  bool

	sync.RWMutex
}

func NewRemoteFormat(log logger.Underlying, compression kt.Compression, cfg *ktranslate.PrometheusFormatConfig) (*RemotePromFormat, error) {
	jf := &RemotePromFormat{
		compression:  compression,
		ContextL:     logger.NewContextLFromUnderlying(logger.SContext{S: "remotePromFormat"}, log),
		invalids:     map[string]bool{},
		lastMetadata: map[string]*kt.LastMetadata{},
		config:       cfg,
	}

	switch compression {
	case kt.CompressionSnappy:
		jf.doSnappy = true
	default:
		return nil, fmt.Errorf("You used an unsupported compression format: %s. For remote_prom, use snappy only.", compression)
	}

	return jf, nil
}

func (f *RemotePromFormat) To(msgs []*kt.JCHF, serBuf []byte) (*kt.Output, error) {
	// First, we need to turn this into a set of []prompb.TimeSeries
	res := make([]prompb.TimeSeries, 0, len(msgs)*4)
	for _, m := range msgs {
		res = append(res, f.toMetric(m)...)
	}

	if len(res) == 0 {
		return nil, nil
	}

	// Marshal proto and compress.
	pbBytes, err := proto.Marshal(&prompb.WriteRequest{
		Timeseries: res,
		Metadata: []prompb.MetricMetadata{prompb.MetricMetadata{
			Type:             prompb.MetricMetadata_GAUGE,
			MetricFamilyName: "test",
		}},
	})
	if err != nil {
		return nil, fmt.Errorf("promwrite: marshaling remote write request proto: %w", err)
	}

	// No Compression.
	if !f.doSnappy {
		return kt.NewOutputWithProviderAndCompanySender(pbBytes, msgs[0].Provider, msgs[0].CompanyId, kt.MetricOutput, ""), nil
	}

	// Snappy Compression.
	if f.doSnappy {
		compressedBytes := snappy.Encode(serBuf, pbBytes)
		return kt.NewOutputWithProviderAndCompanySender(compressedBytes, msgs[0].Provider, msgs[0].CompanyId, kt.MetricOutput, ""), nil
	}

	buf := bytes.NewBuffer(serBuf)
	buf.Reset()
	zw, err := gzip.NewWriterLevel(buf, gzip.DefaultCompression)
	if err != nil {
		return nil, err
	}

	_, err = zw.Write(pbBytes)
	if err != nil {
		return nil, err
	}

	err = zw.Close()
	if err != nil {
		return nil, err
	}

	return kt.NewOutputWithProviderAndCompanySender(buf.Bytes(), msgs[0].Provider, msgs[0].CompanyId, kt.MetricOutput, ""), nil
}

func (f *RemotePromFormat) From(raw *kt.Output) ([]map[string]interface{}, error) {
	values := make([]map[string]interface{}, 0)
	return values, nil
}

func (f *RemotePromFormat) Rollup(rolls []rollup.Rollup) (*kt.Output, error) {
	return nil, nil
}

func (f *RemotePromFormat) toMetric(in *kt.JCHF) []prompb.TimeSeries {
	switch in.EventType {
	case kt.KENTIK_EVENT_TYPE:
		return f.fromKflow(in)
	case kt.KENTIK_EVENT_SNMP_DEV_METRIC:
		return f.fromSnmpDeviceMetric(in)
	case kt.KENTIK_EVENT_SNMP_INT_METRIC:
		return f.fromSnmpInterfaceMetric(in)
	case kt.KENTIK_EVENT_SNMP_METADATA:
		return f.fromSnmpMetadata(in)
	default:
		f.Lock()
		defer f.Unlock()
		if !f.invalids[in.EventType] {
			f.Warnf("Invalid EventType: %s", in.EventType)
			f.invalids[in.EventType] = true
		}
	}

	return nil
}

func (f *RemotePromFormat) fromSnmpDeviceMetric(in *kt.JCHF) []prompb.TimeSeries {
	metrics := in.CustomMetrics
	attr := map[string]interface{}{}
	f.RLock()
	defer f.RUnlock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName], false)
	if f.lastMetadata[in.DeviceName] == nil {
		f.Debugf("Missing device metadata for %s", in.DeviceName)
	}

	res := []prompb.TimeSeries{}
	for m, name := range metrics {
		if m == "" {
			f.Errorf("Missing metric name, skipping %v", attr)
			continue
		}
		if _, ok := in.CustomBigInt[m]; ok {
			attrNew := util.CopyAttrForSnmp(attr, m, name, f.lastMetadata[in.DeviceName], true, false)
			if util.DropOnFilter(attrNew, f.lastMetadata[in.DeviceName], false) {
				continue // This Metric isn't in the white list so lets drop it.
			}

			mtype := name.GetType()
			labels := []prompb.Label{prompb.Label{
				Name:  "name",
				Value: "kentik." + mtype + "." + m,
			}}
			for k, v := range attrNew {
				switch val := v.(type) {
				case string:
					labels = append(labels, prompb.Label{
						Name:  k,
						Value: val,
					})
				case int64, int32:
					labels = append(labels, prompb.Label{
						Name:  k,
						Value: fmt.Sprintf("%v", val),
					})
				}
			}

			//mtype := name.GetType()
			ms := make([]prompb.Sample, 1)
			if name.Format == kt.FloatMS {
				ms[0] = prompb.Sample{
					Timestamp: ptime.FromFloatSeconds(float64(in.Timestamp)),
					Value:     float64(in.CustomBigInt[m]) / 1000.,
				}
			} else {
				ms[0] = prompb.Sample{
					Timestamp: ptime.FromFloatSeconds(float64(in.Timestamp)),
					Value:     float64(in.CustomBigInt[m]),
				}
			}

			res = append(res, prompb.TimeSeries{
				Labels:  labels,
				Samples: ms,
			})
		}
	}

	return res
}

func (f *RemotePromFormat) fromSnmpInterfaceMetric(in *kt.JCHF) []prompb.TimeSeries {
	metrics := in.CustomMetrics
	attr := map[string]interface{}{}
	f.RLock()
	defer f.RUnlock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName], true)
	res := []prompb.TimeSeries{}
	for m, name := range metrics {
		if m == "" {
			f.Errorf("Missing metric name, skipping %v", attr)
			continue
		}
		if _, ok := in.CustomBigInt[m]; ok {
			attrNew := util.CopyAttrForSnmp(attr, m, name, f.lastMetadata[in.DeviceName], false, true)
			if util.DropOnFilter(attrNew, f.lastMetadata[in.DeviceName], true) {
				continue // This Metric isn't in the white list so lets drop it.
			}

			labels := []prompb.Label{prompb.Label{
				Name:  "name",
				Value: "kentik.snmp." + m,
			}}
			for k, v := range attrNew {
				switch val := v.(type) {
				case string:
					labels = append(labels, prompb.Label{
						Name:  k,
						Value: val,
					})
				case int64, int32:
					labels = append(labels, prompb.Label{
						Name:  k,
						Value: fmt.Sprintf("%v", val),
					})
				}
			}

			//mtype := name.GetType()
			ms := make([]prompb.Sample, 1)
			if name.Format == kt.FloatMS {
				ms[0] = prompb.Sample{
					Timestamp: ptime.FromFloatSeconds(float64(in.Timestamp)),
					Value:     float64(in.CustomBigInt[m]) / 1000.,
				}
			} else {
				ms[0] = prompb.Sample{
					Timestamp: ptime.FromFloatSeconds(float64(in.Timestamp)),
					Value:     float64(in.CustomBigInt[m]),
				}
			}

			res = append(res, prompb.TimeSeries{
				Labels:  labels,
				Samples: ms,
			})
		}
	}

	return res
}

func (f *RemotePromFormat) fromSnmpMetadata(in *kt.JCHF) []prompb.TimeSeries {
	if in.DeviceName == "" { // Only run if this is set.
		return nil
	}

	lm := util.SetMetadata(in)

	f.Lock()
	defer f.Unlock()
	if f.lastMetadata[in.DeviceName] == nil || lm.Size() >= f.lastMetadata[in.DeviceName].Size() {
		f.Infof("New Metadata for %s", in.DeviceName)
		f.lastMetadata[in.DeviceName] = lm
	} else {
		f.Infof("The metadata for %s was not updated since the attribute size is smaller. New = %d < Old = %d, Size difference = %v.",
			in.DeviceName, lm.Size(), f.lastMetadata[in.DeviceName].Size(), f.lastMetadata[in.DeviceName].Missing(lm))
	}

	return nil
}

func (f *RemotePromFormat) fromKflow(in *kt.JCHF) []prompb.TimeSeries {
	// Map the basic strings into here.
	attr := map[string]interface{}{}
	metrics := map[string]kt.MetricInfo{"in_bytes": kt.MetricInfo{}, "out_bytes": kt.MetricInfo{}, "in_pkts": kt.MetricInfo{}, "out_pkts": kt.MetricInfo{}, "latency_ms": kt.MetricInfo{}}
	f.RLock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName], false)
	f.RUnlock()
	ms := make([]prompb.Sample, 0, len(metrics))

	labels := make([]prompb.Label, 0, len(attr))
	for k, v := range attr {
		var value string
		switch val := v.(type) {
		case string:
			value = val
		case int64, int32:
			value = fmt.Sprintf("%v", val)
		}
		labels = append(labels, prompb.Label{
			Name:  k,
			Value: value,
		})
	}

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
			ms = append(ms, prompb.Sample{
				Timestamp: ptime.FromFloatSeconds(float64(in.Timestamp)),
				Value:     float64(value),
			})
		}
	}
	res := prompb.TimeSeries{
		Labels:  labels,
		Samples: ms,
	}

	return []prompb.TimeSeries{res}
}
