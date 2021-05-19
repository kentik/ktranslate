package ddog

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"

	datadog "github.com/DataDog/datadog-api-client-go/api/v1/datadog"
)

var (
	DDOG_RATE_TYPE  = "rate"
	DDOG_GAUGE_TYPE = "gauge"
)

type DDogFormat struct {
	logger.ContextL
	compression  kt.Compression
	doGz         bool
	lastMetadata map[string]*kt.LastMetadata
	invalids     map[string]bool
	mux          sync.RWMutex

	EventChan chan []byte
}

func NewFormat(log logger.Underlying, compression kt.Compression) (*DDogFormat, error) {
	jf := &DDogFormat{
		compression:  compression,
		ContextL:     logger.NewContextLFromUnderlying(logger.SContext{S: "ddogFormat"}, log),
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
		return nil, fmt.Errorf("Invalid compression (%s): format ddog only supports none|gzip", compression)
	}

	return jf, nil
}

func (f *DDogFormat) To(msgs []*kt.JCHF, serBuf []byte) (*kt.Output, error) {
	ms := datadog.NewMetricsPayloadWithDefaults()
	series := map[string]*datadog.Series{}
	for _, m := range msgs {
		f.toDDogMetric(m, series)
	}

	if len(series) == 0 {
		return nil, nil
	}

	for _, s := range series {
		ms.Series = append(ms.Series, *s)
	}

	target, err := json.Marshal(ms)
	if err != nil {
		return nil, err
	}

	if !f.doGz {
		return kt.NewOutput(target), nil
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

	return kt.NewOutput(buf.Bytes()), nil
}

func (f *DDogFormat) From(raw *kt.Output) ([]map[string]interface{}, error) {
	values := make([]map[string]interface{}, 0)
	return values, nil
}

func (f *DDogFormat) Rollup(rolls []rollup.Rollup) (*kt.Output, error) {
	ms := datadog.NewMetricsPayloadWithDefaults()
	series := map[string]*datadog.Series{}
	f.toDDogMetricRollup(rolls, series)

	if len(series) == 0 {
		return nil, nil
	}

	for _, s := range series {
		ms.Series = append(ms.Series, *s)
	}

	target, err := json.Marshal(ms)
	if err != nil {
		return nil, err
	}

	if !f.doGz {
		return kt.NewOutput(target), nil
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

	return kt.NewOutput(buf.Bytes()), nil
}

func (f *DDogFormat) toDDogMetricRollup(in []rollup.Rollup, series map[string]*datadog.Series) {
	for _, roll := range in {
		dims := roll.GetDims()
		attr := map[string]string{
			"provider": string(kt.ProviderRouter),
		}
		host := ""
		bad := false
		for i, pt := range strings.Split(roll.Dimension, roll.KeyJoin) {
			attr[dims[i]] = pt
			if pt == "0" || pt == "" {
				bad = true
			}
			if dims[i] == "device_name" {
				host = pt
			}
		}

		// Turn into k=v tags
		tags := []string{}
		for k, v := range attr {
			tags = append(tags, k+"="+v)
		}

		if !bad {
			seriesName := "kentik.rollup." + roll.Name
			key := strings.Join([]string{seriesName, host}, ":")
			intv := roll.Interval.Microseconds()
			interval := datadog.NewNullableInt64(&intv)
			if _, ok := series[key]; !ok {
				series[key] = datadog.NewSeries(seriesName, [][]float64{})
				series[key].Host = &host
				series[key].Interval = *interval
				series[key].Type = &DDOG_RATE_TYPE
				series[key].Tags = &tags
			}
			series[key].Points = append(series[key].Points, []float64{float64(time.Now().Unix()), roll.Metric})
		}
	}
}

func (f *DDogFormat) toDDogMetric(in *kt.JCHF, series map[string]*datadog.Series) {
	host := in.DeviceName
	if host == "" {
		host = strconv.Itoa(int(in.DeviceId))
	}

	// Map the basic strings into here.
	ts := float64(in.Timestamp)
	if ts == 0 {
		ts = float64(time.Now().Unix())
	}
	metrics := map[string]string{"in_bytes": "", "out_bytes": "", "in_pkts": "", "out_pkts": "", "latency_ms": ""}
	for m, _ := range metrics {
		seriesName := "kentik.flow." + m
		key := strings.Join([]string{seriesName, host}, ":")
		if _, ok := series[key]; !ok {
			series[key] = datadog.NewSeries(seriesName, [][]float64{})
			series[key].Host = &host
			series[key].Type = &DDOG_GAUGE_TYPE
		}
		switch m {
		case "in_bytes":
			if in.InBytes > 0 {
				series[key].Points = append(series[key].Points, []float64{ts, float64(in.InBytes * uint64(in.SampleRate))})
			}
		case "out_bytes":
			if in.OutBytes > 0 {
				series[key].Points = append(series[key].Points, []float64{ts, float64(in.OutBytes * uint64(in.SampleRate))})
			}
		case "in_pkts":
			if in.InPkts > 0 {
				series[key].Points = append(series[key].Points, []float64{ts, float64(in.InPkts * uint64(in.SampleRate))})
			}
		case "out_pkts":
			if in.OutPkts > 0 {
				series[key].Points = append(series[key].Points, []float64{ts, float64(in.OutPkts * uint64(in.SampleRate))})
			}
		case "latency_ms":
			if in.CustomInt["appl_latency_ms"] > 0 {
				series[key].Points = append(series[key].Points, []float64{ts, float64(in.CustomInt["appl_latency_ms"])})
			}
		}
	}
}
