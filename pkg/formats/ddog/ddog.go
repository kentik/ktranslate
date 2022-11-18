package ddog

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/formats/util"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
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
	ms := datadogV2.NewMetricPayloadWithDefaults()
	for _, m := range msgs {
		err := f.toDDogMetric(m, ms)
		if err != nil {
			return nil, err
		}
	}

	if len(ms.Series) == 0 {
		return nil, nil
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
	return nil, nil
}

func (f *DDogFormat) toDDogMetric(in *kt.JCHF, ms *datadogV2.MetricPayload) error {
	switch in.EventType {
	case kt.KENTIK_EVENT_SNMP_DEV_METRIC:
		return f.fromSnmpDeviceMetric(in, ms)
	case kt.KENTIK_EVENT_SNMP_INT_METRIC:
		return f.fromSnmpInterfaceMetric(in, ms)
	case kt.KENTIK_EVENT_SNMP_METADATA:
		return f.fromSnmpMetadata(in, ms)
	case kt.KENTIK_EVENT_KTRANS_METRIC:
	}

	return nil
}

func (f *DDogFormat) fromKtranslate(in *kt.JCHF, ms *datadogV2.MetricPayload) error {
	// Map the basic strings into here.
	attr := map[string]interface{}{}
	metrics := map[string]kt.MetricInfo{"name": kt.MetricInfo{}, "value": kt.MetricInfo{}, "count": kt.MetricInfo{}, "one-minute": kt.MetricInfo{}, "95-percentile": kt.MetricInfo{}, "du": kt.MetricInfo{}}
	f.mux.RLock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName], false)
	f.mux.RUnlock()
	tags := getDDMetricTags(attr)
	var value *float64 = nil

	switch in.CustomStr["type"] {
	case "counter":
		if in.CustomStr["force"] == "true" || in.CustomBigInt["count"] > 0 {
			value = datadog.PtrFloat64(float64(in.CustomBigInt["count"]) / 100)
		}
	case "gauge":
		if in.CustomStr["force"] == "true" || in.CustomBigInt["value"] > 0 {
			value = datadog.PtrFloat64(float64(in.CustomBigInt["value"]) / 100)
		}
	case "histogram":
		if in.CustomStr["force"] == "true" || in.CustomBigInt["95-percentile"] > 0 {
			value = datadog.PtrFloat64(float64(in.CustomBigInt["95-percentile"]) / 100)
		}
	case "meter":
		if in.CustomStr["force"] == "true" || in.CustomBigInt["one-minute"] > 0 {
			value = datadog.PtrFloat64(float64(in.CustomBigInt["one-minute"]) / 100)
		}
	case "timer":
		if in.CustomStr["force"] == "true" || in.CustomBigInt["95-percentile"] > 0 {
			value = datadog.PtrFloat64(float64(in.CustomBigInt["95-percentile"]) / 100)
		}
	}

	if value != nil {
		ms.Series = append(ms.Series, datadogV2.MetricSeries{
			Metric: in.CustomStr["name"],
			Type:   datadogV2.METRICINTAKETYPE_GAUGE.Ptr(),
			Points: []datadogV2.MetricPoint{
				{
					Timestamp: datadog.PtrInt64(in.Timestamp),
					Value:     value,
				},
			},
			Tags: tags,
		})
	}

	return nil
}

func (f *DDogFormat) fromSnmpMetadata(in *kt.JCHF, ms *datadogV2.MetricPayload) error {
	if in.DeviceName == "" { // Only run if this is set.
		return nil
	}

	lm := util.SetMetadata(in)

	f.mux.Lock()
	defer f.mux.Unlock()
	f.lastMetadata[in.DeviceName] = lm

	return nil
}

func (f *DDogFormat) fromSnmpDeviceMetric(in *kt.JCHF, ms *datadogV2.MetricPayload) error {
	metrics := in.CustomMetrics
	attr := map[string]interface{}{}
	f.mux.RLock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName], true)
	f.mux.RUnlock()

	for m, name := range metrics {
		if m == "" {
			f.Errorf("Missing metric name, skipping %v", attr)
			continue
		}
		if _, ok := in.CustomBigInt[m]; ok {
			attrNew := util.CopyAttrForSnmp(attr, m, name, f.lastMetadata[in.DeviceName], false)
			if util.DropOnFilter(attrNew, f.lastMetadata[in.DeviceName], false) {
				continue // This Metric isn't in the white list so lets drop it.
			}
			seriesName := "kentik.snmp." + m
			tags := getDDMetricTags(attrNew)
			if name.Format == kt.FloatMS {
				ms.Series = append(ms.Series, datadogV2.MetricSeries{
					Metric: seriesName,
					Type:   datadogV2.METRICINTAKETYPE_GAUGE.Ptr(),
					Points: []datadogV2.MetricPoint{
						{
							Timestamp: datadog.PtrInt64(in.Timestamp),
							Value:     datadog.PtrFloat64(float64(float64(in.CustomBigInt[m]) / 1000)),
						},
					},
					Tags: tags,
				})
			} else {
				ms.Series = append(ms.Series, datadogV2.MetricSeries{
					Metric: seriesName,
					Type:   datadogV2.METRICINTAKETYPE_GAUGE.Ptr(),
					Points: []datadogV2.MetricPoint{
						{
							Timestamp: datadog.PtrInt64(in.Timestamp),
							Value:     datadog.PtrFloat64(float64(in.CustomBigInt[m])),
						},
					},
					Tags: tags,
				})
			}
		}
	}

	return nil
}

func (f *DDogFormat) fromSnmpInterfaceMetric(in *kt.JCHF, ms *datadogV2.MetricPayload) error {
	metrics := in.CustomMetrics
	attr := map[string]interface{}{}
	f.mux.RLock()
	defer f.mux.RUnlock()
	util.SetAttr(attr, in, metrics, f.lastMetadata[in.DeviceName], true)

	profileName := "snmp"
	for m, name := range metrics {
		if m == "" {
			f.Errorf("Missing metric name, skipping %v", attr)
			continue
		}
		if _, ok := in.CustomBigInt[m]; ok {
			attrNew := util.CopyAttrForSnmp(attr, m, name, f.lastMetadata[in.DeviceName], false)
			if util.DropOnFilter(attrNew, f.lastMetadata[in.DeviceName], true) {
				continue // This Metric isn't in the white list so lets drop it.
			}

			profileName = name.Profile
			seriesName := "kentik.snmp." + m
			tags := getDDMetricTags(attrNew)
			if name.Format == kt.FloatMS {
				ms.Series = append(ms.Series, datadogV2.MetricSeries{
					Metric: seriesName,
					Type:   datadogV2.METRICINTAKETYPE_GAUGE.Ptr(),
					Points: []datadogV2.MetricPoint{
						{
							Timestamp: datadog.PtrInt64(in.Timestamp),
							Value:     datadog.PtrFloat64(float64(float64(in.CustomBigInt[m]) / 1000)),
						},
					},
					Tags: tags,
				})
			} else {
				ms.Series = append(ms.Series, datadogV2.MetricSeries{
					Metric: seriesName,
					Type:   datadogV2.METRICINTAKETYPE_GAUGE.Ptr(),
					Points: []datadogV2.MetricPoint{
						{
							Timestamp: datadog.PtrInt64(in.Timestamp),
							Value:     datadog.PtrFloat64(float64(in.CustomBigInt[m])),
						},
					},
					Tags: tags,
				})
			}
		}
	}

	// grab rates computed over time here if possible.
	if f.lastMetadata[in.DeviceName] != nil {
		f.setRates(ms, "In", in, attr, profileName)
		f.setRates(ms, "Out", in, attr, profileName)
	}

	return nil
}

func (f *DDogFormat) setRates(ms *datadogV2.MetricPayload, direction string, in *kt.JCHF, attr map[string]interface{}, profileName string) {
	var port kt.IfaceID
	if direction == "In" {
		port = in.InputPort
	} else {
		port = in.OutputPort
	}
	utilName := fmt.Sprintf("If%sUtilization", direction)
	bitRate := fmt.Sprintf("If%sBitRate", direction)
	pktRate := fmt.Sprintf("If%sPktRate", direction)
	totalBytes := in.CustomBigInt[fmt.Sprintf("ifHC%sOctets", direction)]
	totalPkts := in.CustomBigInt[fmt.Sprintf("ifHC%sUcastPkts", direction)] + in.CustomBigInt[fmt.Sprintf("ifHC%sMulticastPkts", direction)] + in.CustomBigInt[fmt.Sprintf("ifHC%sBroadcastPkts", direction)]

	if ii, ok := f.lastMetadata[in.DeviceName].InterfaceInfo[port]; ok {
		if speed, ok := ii["Speed"]; ok {
			if ispeed, ok := speed.(int32); ok {
				uptime := in.CustomBigInt["Uptime"]
				uptimeSpeed := uptime * (int64(ispeed) * 10000) // Convert into bits here, from megabits. Also divide by 100 to convert uptime into seconds, from centi-seconds.
				if uptimeSpeed > 0 {
					attrNew := util.CopyAttrForSnmp(attr, utilName, kt.MetricInfo{Oid: "computed", Mib: "computed", Profile: profileName, Table: "if"}, f.lastMetadata[in.DeviceName], false)
					if !util.DropOnFilter(attrNew, f.lastMetadata[in.DeviceName], true) {
						tags := getDDMetricTags(attrNew)
						if totalBytes > 0 {
							ms.Series = append(ms.Series, datadogV2.MetricSeries{
								Metric: "kentik.snmp." + utilName,
								Type:   datadogV2.METRICINTAKETYPE_GAUGE.Ptr(),
								Points: []datadogV2.MetricPoint{
									{
										Timestamp: datadog.PtrInt64(in.Timestamp),
										Value:     datadog.PtrFloat64(float64(totalBytes*8*100) / float64(uptimeSpeed)),
									},
								},
								Tags: tags,
							},
								datadogV2.MetricSeries{
									Metric: "kentik.snmp." + bitRate,
									Type:   datadogV2.METRICINTAKETYPE_GAUGE.Ptr(),
									Points: []datadogV2.MetricPoint{
										{
											Timestamp: datadog.PtrInt64(in.Timestamp),
											Value:     datadog.PtrFloat64(float64(totalBytes*8*100) / float64(uptime)),
										},
									},
									Tags: tags,
								},
							)
						}
						if totalPkts > 0 {
							ms.Series = append(ms.Series, datadogV2.MetricSeries{
								Metric: "kentik.snmp." + pktRate,
								Type:   datadogV2.METRICINTAKETYPE_GAUGE.Ptr(),
								Points: []datadogV2.MetricPoint{
									{
										Timestamp: datadog.PtrInt64(in.Timestamp),
										Value:     datadog.PtrFloat64(float64(totalPkts*100) / float64(uptime)),
									},
								},
								Tags: tags,
							})
						}
					}
				}
			}
		}
	}
}

func getDDMetricTags(attrNew map[string]interface{}) []string {
	tags := []string{}
	for k, v := range attrNew {
		switch nv := v.(type) {
		case string:
			tags = append(tags, strings.Join([]string{k, nv}, ":"))
		case int64:
			tags = append(tags, strings.Join([]string{k, strconv.Itoa(int(nv))}, ":"))
		case int32:
			tags = append(tags, strings.Join([]string{k, strconv.Itoa(int(nv))}, ":"))
		case kt.Cid:
			tags = append(tags, strings.Join([]string{k, strconv.Itoa(int(nv))}, ":"))
		}
	}
	return tags
}
