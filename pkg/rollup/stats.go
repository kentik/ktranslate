package rollup

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/kentik/ktranslate/pkg/eggs/logger"

	"github.com/montanaflynn/stats"
)

type StatsRollup struct {
	logger.ContextL
	rollupBase
	state  map[string][]float64
	statr  func(stats.Float64Data) (float64, error)
	statr2 func(stats.Float64Data, float64) (float64, error)
	arg2   float64
}

func newStatsRollup(log logger.Underlying, rd RollupDef) (*StatsRollup, error) {
	r := &StatsRollup{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "sumRollup"}, log),
		state:    map[string][]float64{},
	}

	switch rd.Method {
	case "sum":
		r.statr = stats.Sum
	case "min":
		r.statr = stats.Min
	case "max":
		r.statr = stats.Max
	case "mean":
		r.statr = stats.Mean
	case "median":
		r.statr = stats.Median
	case "entropy":
		r.statr = stats.Entropy
	default:
		if strings.HasPrefix(string(rd.Method), "percentilerank") {
			r.statr2 = stats.PercentileNearestRank
			val, err := strconv.Atoi(string(rd.Method[len("percentilerank"):]))
			if err != nil {
				return nil, fmt.Errorf("Unknown rollup: %s", rd.String())
			}
			r.arg2 = float64(val)
		} else if strings.HasPrefix(string(rd.Method), "percentile") {
			r.statr2 = stats.Percentile
			val, err := strconv.Atoi(string(rd.Method[len("percentile"):]))
			if err != nil {
				return nil, fmt.Errorf("Unknown rollup: %s", rd.String())
			}
			r.arg2 = float64(val)
		} else {
			return nil, fmt.Errorf("Unknown rollup: %s", rd.String())
		}
	}

	err := r.init(rd)
	if err != nil {
		return nil, err
	}
	r.Infof("New Rollup: %s -> %s, value of %f", r.eventType, rd.String(), r.arg2)
	return r, nil
}

func (r *StatsRollup) Add(in []*kt.JCHF) {
	toAdd := map[string][]float64{}
	for i, val := range in {
		mapr := val.ToMap()
		key := r.getKey(mapr)
		if _, ok := toAdd[key]; !ok {
			toAdd[key] = make([]float64, len(in))
		}

		for _, metric := range r.metrics {
			if mm, ok := mapr[metric]; ok {
				if m, ok := mm.(int64); ok {
					toAdd[key][i] = float64(m)
				}
			}
		}
		for _, m := range r.multiMetrics { // Now handle the 2 level deep metrics
			if m1, ok := mapr[m[0]]; ok {
				switch mm := m1.(type) {
				case map[string]int32:
					toAdd[key][i] = float64(mm[m[1]])
				case map[string]int64:
					toAdd[key][i] = float64(mm[m[1]])
				}
			}
		}
		if r.sample && val.SampleRate > 0 { // If we are adjusting for sample rate for this rollup, do so now.
			toAdd[key][i] *= float64(val.SampleRate)
		}
	}

	// Now need a lock for actually updating the current rollup state.
	r.mux.Lock()
	for k, v := range toAdd {
		if _, ok := r.state[k]; !ok {
			r.state[k] = []float64{}
		}
		r.state[k] = append(r.state[k], v...)
	}
	r.mux.Unlock()
}

func (r *StatsRollup) Export() []Rollup {
	r.mux.Lock()
	os := r.state
	r.state = map[string][]float64{}
	r.mux.Unlock()

	ot := r.dtime
	r.dtime = time.Now()
	keys := make([]Rollup, len(os))
	next := 0
	for k, v := range os {
		var value float64
		var err error
		if r.statr2 != nil {
			value, err = r.statr2(v, r.arg2)
		} else {
			value, err = r.statr(v)
		}
		if err != nil {
			r.Errorf("Error calculating: %v", err)
		} else {
			keys[next] = Rollup{EventType: r.eventType, Dimension: k, Metric: value, KeyJoin: r.keyJoin, dims: combo(r.dims, r.multiDims), Interval: r.dtime.Sub(ot)}
			next++
		}
	}

	sort.Sort(byValue(keys))
	if len(keys) > r.topK {
		return keys[0:r.topK]
	}

	return keys
}
