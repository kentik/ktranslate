package carbon

import (
	"fmt"
	"strings"
	"sync"

	"github.com/kentik/ktranslate/pkg/formats/util"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

type CarbonFormat struct {
	logger.ContextL
	invalids map[string]bool
	mux      sync.RWMutex
}

type CarbonData struct {
	Type          string
	Name          string
	Value         int64
	Unit          string
	IntrinsicTags map[string]interface{}
	MetaTags      map[string]interface{}
	Timestamp     int64
}

// formatted using http://metrics20.org/spec/
func (d *CarbonData) String() string {
	intrinsicTags := parseTags(d.IntrinsicTags)
	metaTags := parseTags(d.MetaTags)

	return fmt.Sprintf("metric=%s mtype=%s unit=%s %s  %s %d %d",
		d.Name,
		d.Type,
		d.Unit,
		strings.Join(intrinsicTags, " "),
		strings.Join(metaTags, " "),
		d.Value,
		d.Timestamp)
}

type CarbonDataSet []*CarbonData

func (s CarbonDataSet) Bytes() []byte {
	res := make([]string, 0)
	for _, l := range s {
		res = append(res, l.String())
	}
	return []byte(strings.Join(res, "\n"))
}

func NewFormat(log logger.Underlying, compression kt.Compression) (*CarbonFormat, error) {
	jf := &CarbonFormat{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "carbonFormat"}, log),
		invalids: map[string]bool{},
	}

	return jf, nil
}

func (f *CarbonFormat) To(msgs []*kt.JCHF, serBuf []byte) (*kt.Output, error) {
	res := make([]*CarbonData, 0, len(msgs))
	for _, m := range msgs {
		res = append(res, f.toCarbonMetric(m)...)
	}

	if len(res) == 0 {
		return nil, nil
	}

	data := CarbonDataSet(res).Bytes()

	return kt.NewOutput(data), nil
}

func (f *CarbonFormat) toCarbonMetric(in *kt.JCHF) []*CarbonData {
	switch in.EventType {
	case kt.KENTIK_EVENT_SYNTH, kt.KENTIK_EVENT_TRACE:
		return f.fromKsynth(in)
	case kt.KENTIK_EVENT_TYPE:
		return f.fromKflow(in)
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

// Not supported
func (f *CarbonFormat) From(raw *kt.Output) ([]map[string]interface{}, error) {
	values := make([]map[string]interface{}, 0)
	return values, nil
}

// Not supported
func (f *CarbonFormat) Rollup(rolls []rollup.Rollup) (*kt.Output, error) {
	return nil, nil
}

func (f *CarbonFormat) fromKflow(in *kt.JCHF) []*CarbonData {
	attr := map[string]interface{}{}
	metrics := map[string]kt.MetricInfo{
		"in_bytes":   kt.MetricInfo{},
		"out_bytes":  kt.MetricInfo{},
		"in_pkts":    kt.MetricInfo{},
		"out_pkts":   kt.MetricInfo{},
		"latency_ms": kt.MetricInfo{},
	}
	util.SetAttr(attr, in, metrics, nil, false)
	metricData := []*CarbonData{}
	for m, _ := range metrics {
		v := int64(0)
		mtype := "count"
		unit := ""
		switch m {
		case "in_bytes":
			v = int64(in.InBytes * uint64(in.SampleRate))
			mtype = "rate"
			unit = "B/s"
		case "out_bytes":
			v = int64(in.OutBytes * uint64(in.SampleRate))
			mtype = "rate"
			unit = "B/s"
		case "in_pkts":
			v = int64(in.InPkts * uint64(in.SampleRate))
			mtype = "rate"
			unit = "pckt/s"
		case "out_pkts":
			v = int64(in.OutPkts * uint64(in.SampleRate))
			mtype = "rate"
			unit = "pckt/s"
		case "latency_ms":
			v = int64(in.CustomInt["appl_latency_ms"])
			mtype = "gauge"
			unit = "ms"
		}
		// intrinsic tags are the metric dimensions
		intrinsicTags := map[string]interface{}{
			"device": in.DeviceId.Itoa(),
		}
		if v := in.DstAddr; v != "" {
			intrinsicTags["dst_addr"] = v
		}
		if v := in.SrcAddr; v != "" {
			intrinsicTags["src_addr"] = v
		}
		if v := in.Protocol; v != "" {
			intrinsicTags["protocol"] = v
		}

		// meta tags are used to decorate and add context
		metaTags := map[string]interface{}{
			"l4_dst_port":    in.L4DstPort,
			"l4_src_port":    in.L4SrcPort,
			"src_as":         in.SrcAs,
			"dst_as":         in.DstAs,
			"tcp_retransmit": in.TcpRetransmit,
			"vlan_in":        in.VlanIn,
			"vlan_out":       in.VlanOut,
		}
		if v := in.DstGeo; v != "" {
			metaTags["dst_geo"] = v
		}
		if v := in.SrcGeo; v != "" {
			metaTags["src_geo"] = v
		}
		if v := in.DstBgpAsPath; v != "" {
			metaTags["dst_bgp_as_path"] = v
		}
		if v := in.SrcBgpAsPath; v != "" {
			metaTags["src_bgp_as_path"] = v
		}
		for k, v := range in.CustomStr {
			if v == "" {
				continue
			}
			metaTags[k] = v
		}
		metricData = append(metricData, &CarbonData{
			Type:          mtype,
			Name:          m,
			Value:         v,
			Unit:          unit,
			Timestamp:     in.Timestamp,
			IntrinsicTags: intrinsicTags,
			MetaTags:      metaTags,
		})
	}

	return metricData
}

func (f *CarbonFormat) fromKsynth(in *kt.JCHF) []*CarbonData {
	metrics := []string{
		"ping_avg_rtt",
		"ping_jit_rtt",
		"ping_max_rtt",
		"ping_std_rtt",
	}
	metricData := []*CarbonData{}
	for _, m := range metrics {
		val := int64(0)
		mtype := "gauge"
		unit := "μ"
		if v, ok := in.CustomInt[m]; ok {
			val = int64(v)
		}

		// intrinsic tags are the metric dimensions
		intrinsicTags := map[string]interface{}{}
		if v, ok := in.CustomStr["application_type"]; ok && v != "" {
			intrinsicTags["application_type"] = v
		}
		if v := in.DstAddr; v != "" {
			intrinsicTags["dst_addr"] = v
		}
		if v := in.SrcAddr; v != "" {
			intrinsicTags["src_addr"] = v
		}

		// meta tags are used to decorate and add context
		metaTags := map[string]interface{}{
			"src_as": in.SrcAs,
			"dst_as": in.DstAs,
		}
		if v := in.DstGeo; v != "" {
			metaTags["dst_geo"] = v
		}
		if v := in.SrcGeo; v != "" {
			metaTags["src_geo"] = v
		}
		if v := in.SrcGeo; v != "" {
			metaTags["src_geo"] = v
		}
		if v := in.SrcGeoRegion; v != "" {
			metaTags["src_geo_region"] = v
		}
		if v := in.DstGeoRegion; v != "" {
			metaTags["dst_geo_region"] = v
		}
		if v := in.SrcGeoCity; v != "" {
			metaTags["src_geo_city"] = v
		}
		if v := in.DstGeoCity; v != "" {
			metaTags["dst_geo_city"] = v
		}
		for k, v := range in.CustomStr {
			if v == "" {
				continue
			}
			metaTags[k] = v
		}
		metricData = append(metricData, &CarbonData{
			Type:          mtype,
			Name:          m,
			Value:         val,
			Unit:          unit,
			Timestamp:     in.Timestamp,
			IntrinsicTags: intrinsicTags,
			MetaTags:      metaTags,
		})
	}

	return metricData
}

func parseTags(tags map[string]interface{}) []string {
	parsedTags := make([]string, len(tags))
	x := 0
	for key, v := range tags {
		kval := strings.ReplaceAll(key, " ", "_")
		switch t := v.(type) {
		case string:
			parsedTags[x] = fmt.Sprintf("%s=%s", kval, strings.ReplaceAll(t, " ", "_"))
		default:
			parsedTags[x] = fmt.Sprintf("%s=%v", kval, v)
		}

		x++
	}

	return parsedTags
}
