package prom

import (
	"flag"
	"strings"

	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"

	"github.com/kentik/ktranslate/pkg/eggs/logger"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	doCollectorStats = flag.Bool("info_collector", false, "Also send stats about this collector")
)

type PromFormat struct {
	logger.ContextL
	vecs map[string]*prometheus.CounterVec
}

func NewFormat(log logger.Underlying, compression kt.Compression) (*PromFormat, error) {
	jf := &PromFormat{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "influxFormat"}, log),
		vecs:     make(map[string]*prometheus.CounterVec),
	}

	if *doCollectorStats {
		prometheus.MustRegister(prometheus.NewBuildInfoCollector())
	}

	return jf, nil
}

// Not supported.
func (f *PromFormat) To(msgs []*kt.JCHF, serBuf []byte) ([]byte, error) {
	// Noop here because we only support rollups in promethius format.
	return nil, nil
}

// Not supported.
func (f *PromFormat) From(raw []byte) ([]map[string]interface{}, error) {
	values := make([]map[string]interface{}, 0)
	return values, nil
}

func (f *PromFormat) Rollup(rolls []rollup.Rollup) ([]byte, error) {
	for _, r := range rolls {
		pkts := strings.Split(r.EventType, ":")
		if _, ok := f.vecs[r.EventType]; !ok {
			f.vecs[r.EventType] = prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: strings.Join(pkts[0:2], ":"),
				},
				pkts[2:],
			)
			prometheus.MustRegister(f.vecs[r.EventType])
		}
		f.vecs[r.EventType].WithLabelValues(strings.Split(r.Dimension, r.KeyJoin)...).Add(float64(r.Metric))
	}

	return nil, nil
}
