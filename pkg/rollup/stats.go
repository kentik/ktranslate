package rollup

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"

	"github.com/montanaflynn/stats"
)

const (
	CHAN_SLACK = 10000
)

type StatsRollup struct {
	logger.ContextL
	rollupBase
	state     map[string][]float64
	statr     func(stats.Float64Data) (float64, error)
	statr2    func(stats.Float64Data, float64) (float64, error)
	arg2      float64
	isSum     bool
	kvs       chan map[string]float64
	exportKvs chan chan []Rollup
}

func newStatsRollup(log logger.Underlying, rd RollupDef) (*StatsRollup, error) {
	r := &StatsRollup{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "sumRollup"}, log),
		state:    map[string][]float64{},
	}

	switch rd.Method {
	case "sum":
		r.statr = stats.Sum
		r.isSum = true
		r.kvs = make(chan map[string]float64, CHAN_SLACK)
		r.exportKvs = make(chan chan []Rollup)
		go r.sumKvs()
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

func (r *StatsRollup) addSum(in []map[string]interface{}) {
	ks := map[string]float64{}
	for _, mapr := range in {
		key := r.getKey(mapr)
		sr := mapr["sample_rate"].(int64)
		for _, metric := range r.metrics {
			if mm, ok := mapr[metric]; ok {
				if m, ok := mm.(int64); ok {
					if m > 0 {
						if r.sample && sr > 0 { // If we are adjusting for sample rate for this rollup, do so now.
							m *= sr
						}
						ks[key] += float64(m)
					}
				}
			}
		}
	}

	// Dump into our hash map here
	select {
	case r.kvs <- ks:
	default:
		r.Warnf("kvs chan full")
	}
}

func (r *StatsRollup) Add(in []map[string]interface{}) {
	if r.isSum {
		r.addSum(in)
		return
	}
}

func (r *StatsRollup) Export() []Rollup {
	if r.isSum {
		rc := make(chan []Rollup)
		r.exportKvs <- rc
		return <-rc
	}

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
			keys[next] = Rollup{Name: r.name, EventType: r.eventType, Dimension: k, Metric: value, KeyJoin: r.keyJoin, dims: combo(r.dims, r.multiDims), Interval: r.dtime.Sub(ot)}
			next++
		}
	}

	sort.Sort(byValue(keys))
	if len(keys) > r.topK {
		return keys[0:r.topK]
	}

	return keys
}

func (r *StatsRollup) sumKvs() {
	vs := map[string]float64{}
	for {
		select {
		case itm := <-r.kvs: // Just add to our map
			for k, v := range itm {
				vs[k] += v
			}
		case rc := <-r.exportKvs: // Return the top results
			go r.exportSum(vs, rc)
			vs = map[string]float64{}
		}
	}
}

func (r *StatsRollup) exportSum(vs map[string]float64, rc chan []Rollup) {
	if len(vs) == 0 {
		rc <- nil
		return
	}

	ot := r.dtime
	r.dtime = time.Now()
	keys := make([]Rollup, 0, len(vs))
	total := 0.0
	for k, v := range vs {
		keys = append(keys, Rollup{Name: r.name, EventType: r.eventType, Dimension: k, Metric: v, KeyJoin: r.keyJoin, dims: combo(r.dims, r.multiDims), Interval: r.dtime.Sub(ot)})
		total += v
	}

	sort.Sort(byValue(keys))
	if len(keys) > r.topK {
		top := keys[0:r.topK]
		dims := combo(r.dims, r.multiDims)
		totals := make([]string, len(dims))
		for i, _ := range dims {
			totals[i] = "total"
		}
		top = append(top, Rollup{Name: r.name, EventType: r.eventType, Dimension: strings.Join(totals, r.keyJoin), Metric: total, KeyJoin: r.keyJoin, dims: dims, Interval: r.dtime.Sub(ot)})
		rc <- top
		return
	}

	rc <- keys
	return
}
