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
	Type      string
	Name      string
	Value     int64
	Unit      string
	Tags      map[string]interface{}
	Timestamp int64
}

// formatted using http://metrics20.org/spec/
func (d *CarbonData) String() string {
	tags := make([]string, len(d.Tags))
	i := 0
	for key, v := range d.Tags {
		kval := strings.ReplaceAll(key, " ", "_")
		switch t := v.(type) {
		case string:
			tags[i] = fmt.Sprintf("%s=%s", kval, strings.ReplaceAll(t, " ", "_"))
		default:
			tags[i] = fmt.Sprintf("%s=%v", kval, v)
		}

		i++
	}

	return fmt.Sprintf("metric=%s mtype=%s unit=%s  %s %d %d",
		d.Name,
		d.Type,
		d.Unit,
		strings.Join(tags, " "),
		d.Value,
		d.Timestamp)
}

type CarbonDataSet []CarbonData

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
	res := make([]CarbonData, 0, len(msgs))
	for _, m := range msgs {
		res = append(res, f.toCarbonMetric(m)...)
	}

	if len(res) == 0 {
		return nil, nil
	}

	return kt.NewOutput(CarbonDataSet(res).Bytes()), nil
}

func (f *CarbonFormat) toCarbonMetric(in *kt.JCHF) []CarbonData {
	switch in.EventType {
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

func (f *CarbonFormat) fromKflow(in *kt.JCHF) []CarbonData {
	attr := map[string]interface{}{}
	metrics := map[string]kt.MetricInfo{"in_bytes": kt.MetricInfo{}, "out_bytes": kt.MetricInfo{}, "in_pkts": kt.MetricInfo{}, "out_pkts": kt.MetricInfo{}, "latency_ms": kt.MetricInfo{}}
	util.SetAttr(attr, in, metrics, nil)
	metricData := []CarbonData{}
	for m, _ := range metrics {
		v := int64(0)
		mtype := ""
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
		metricData = append(metricData, CarbonData{
			Type:      mtype,
			Name:      m,
			Value:     v,
			Unit:      unit,
			Timestamp: in.Timestamp,
			Tags:      attr,
		})
	}

	return metricData
}
