package cmetrics

import (
	"log"
	"time"

	metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate/pkg/kt"
)

type OpenJCHFConfig struct {
	OutChan       chan []*kt.JCHF   // send outgoing metrics here.
	Registry      metrics.Registry  // Registry to be exported
	FlushInterval time.Duration     // Flush interval
	DurationUnit  time.Duration     // Time conversion unit for durations
	Prefix        string            // Prefix to be prepended to metric names
	Tags          map[string]string // add these tags to each metric writen
	Extra         map[string]string
}

func OpenJCHFWithConfig(c OpenJCHFConfig) {
	for _ = range time.Tick(c.FlushInterval) {
		if err := sendJCHF(&c); nil != err {
			log.Println(err)
		}
	}
}

func sendJCHF(c *OpenJCHFConfig) error {
	shortHostnameBase := GetShortHostname()
	now := time.Now().UnixNano()
	du := float64(c.DurationUnit)
	base := []*kt.JCHF{}

	c.Registry.Each(func(baseName string, i interface{}) {
		name, tags := ExpandTags(baseName, c.Prefix, shortHostnameBase, c.Tags, c.Extra)
		dst := kt.NewJCHF()
		dst.CustomStr = make(map[string]string)
		dst.CustomInt = make(map[string]int32)
		dst.CustomBigInt = make(map[string]int64)
		dst.EventType = kt.KENTIK_EVENT_KTRANS_METRIC
		dst.Provider = kt.ProviderAgent
		dst.Timestamp = now
		dst.DeviceName = shortHostnameBase
		dst.CustomStr["name"] = name
		dst.CustomInt["du"] = int32(du)

		for k, v := range tags {
			dst.CustomStr[k] = v
		}

		switch metric := i.(type) {
		case metrics.Counter:
			dst.CustomStr["type"] = "counter"
			dst.CustomBigInt["count"] = metric.Count() * 100
		case metrics.Gauge:
			dst.CustomStr["type"] = "gauge"
			dst.CustomBigInt["value"] = metric.Value() * 100
		case metrics.Histogram:
			dst.CustomStr["type"] = "histogram"
			h := metric.Snapshot()
			ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			dst.CustomBigInt["count"] = h.Count() * 100
			dst.CustomBigInt["min"] = h.Min() * 100
			dst.CustomBigInt["max"] = h.Max() * 100
			dst.CustomBigInt["mean"] = int64(h.Mean() * 100)
			dst.CustomBigInt["95-percentile"] = int64(ps[2] * 100)
			dst.CustomBigInt["99-percentile"] = int64(ps[3] * 100)
			metric.Clear()
		case metrics.Meter:
			dst.CustomStr["type"] = "meter"
			m := metric.Snapshot()
			dst.CustomBigInt["count"] = m.Count() * 100
			dst.CustomBigInt["one-minute"] = int64(m.Rate1() * 100)
		case metrics.Timer:
			dst.CustomStr["type"] = "timer"
			t := metric.Snapshot()
			ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			dst.CustomBigInt["count"] = t.Count() * 100
			dst.CustomBigInt["min"] = int64((float64(t.Min()) / du) * 100)
			dst.CustomBigInt["max"] = int64((float64(t.Max()) / du) * 100)
			dst.CustomBigInt["mean"] = int64((t.Mean() / du) * 100)
			dst.CustomBigInt["95-percentile"] = int64((ps[2] / du) * 100)
			dst.CustomBigInt["99-percentile"] = int64((ps[3] / du) * 100)
			dst.CustomBigInt["one-minute"] = int64(t.Rate1() * 100)
			dst.CustomBigInt["five-minute"] = int64(t.Rate5() * 100)
			dst.CustomBigInt["fifteen-minute"] = int64(t.Rate15() * 100)
			metric.Clear()
		}
		base = append(base, dst)
	})

	select {
	case c.OutChan <- base:
	default:
		log.Printf("metrics chan full: %d", len(c.OutChan))
	}

	return nil
}
