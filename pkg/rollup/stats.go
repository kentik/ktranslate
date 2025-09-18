package rollup

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"

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
	kvs       chan *sumset
	exportKvs chan chan []Rollup
	config    *ktranslate.RollupConfig
}

type sumset struct {
	sum   map[string]uint64
	count map[string]uint64
	min   map[string]uint64
	max   map[string]uint64
	prov  map[string]kt.Provider
}

func newStatsRollup(log logger.Underlying, rd RollupDef, cfg *ktranslate.RollupConfig) (*StatsRollup, error) {
	r := &StatsRollup{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "sumRollup"}, log),
		state:    map[string][]float64{},
		config:   cfg,
	}

	r.keyJoin = cfg.JoinKey
	r.topK = cfg.TopK

	switch rd.Method {
	case "sum":
		r.statr = stats.Sum
		r.isSum = true
		r.kvs = make(chan *sumset, CHAN_SLACK)
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
	sum := map[string]uint64{}
	count := map[string]uint64{}
	min := map[string]uint64{}
	max := map[string]uint64{}
	prov := map[string]kt.Provider{}

	for _, mapr := range in {
		key := r.getKey(mapr)
		sr := uint64(mapr["sample_rate"].(int64))
		value := uint64(0)

		for _, metric := range r.metrics { // 1 level deap one first.
			if mm, ok := mapr[metric]; ok {
				if m, ok := mm.(int64); ok {
					value += uint64(m)
				}
			}
		}

		for _, m := range r.multiMetrics { // Now handle the 2 level deep metrics
			if m1, ok := mapr[m[0]]; ok {
				switch mm := m1.(type) {
				case map[string]int32:
					value += uint64(mm[m[1]])
				case map[string]int64:
					value += uint64(mm[m[1]])
				}
			}
		}

		if r.sample && sr > 0 { // If we are adjusting for sample rate for this rollup, do so now.
			value *= sr
		}
		sum[key] += value
		count[key]++
		if _, ok := min[key]; !ok {
			min[key] = value
		} else if min[key] > value {
			min[key] = value
		}
		if max[key] < value {
			max[key] = value
		}
		prov[key] = mapr["provider"].(kt.Provider)
	}

	// Dump into our hash map here
	select {
	case r.kvs <- &sumset{sum: sum, count: count, min: min, max: max, prov: prov}:
	default:
		r.Warnf("kvs chan full")
	}
}

func (r *StatsRollup) Add(in []map[string]interface{}) {
	if r.hasFilters {
		in = r.filter(in)
	}

	if r.isSum { // this is a fast path for pure additive rollups.
		r.addSum(in)
		return
	}

	// And this is the slow path for more fine grained rollups.
	toAdd := map[string][]float64{}
	for i, mapr := range in {
		key := r.getKey(mapr)
		sr := mapr["sample_rate"].(int64)
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
		if r.sample && sr > 0 { // If we are adjusting for sample rate for this rollup, do so now.
			toAdd[key][i] *= float64(sr)
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
	fullName := strings.Join(r.nameSet, ";")
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
			keys[next] = Rollup{Name: fullName, EventType: r.eventType, Dimension: k, Metric: value, KeyJoin: r.keyJoin, dims: combo(r.multiDims), Interval: r.dtime.Sub(ot)}
			next++
		}
	}

	sort.Sort(byValue(keys))
	if r.config.TopK > 0 && len(keys) > r.config.TopK {
		return keys[0:r.config.TopK]
	}

	return keys
}

func (r *StatsRollup) sumKvs() {
	sum := map[string]uint64{}
	count := map[string]uint64{}
	min := map[string]uint64{}
	max := map[string]uint64{}
	prov := map[string]kt.Provider{}

	for {
		select {
		case itm := <-r.kvs: // Just add to our map // 	case r.kvs <- []map[string]int64{sum, count, min, max}:
			for k, v := range itm.sum {
				sum[k] += v
			}
			for k, v := range itm.count {
				count[k] += v
			}
			for k, v := range itm.min {
				if _, ok := min[k]; !ok {
					min[k] = v
				} else if min[k] > v {
					min[k] = v
				}
			}
			for k, v := range itm.max {
				if max[k] < v {
					max[k] = v
				}
			}
			for k, v := range itm.prov {
				prov[k] = v
			}
		case rc := <-r.exportKvs: // Return the top results
			go r.exportSum(sum, count, min, max, prov, rc)
			sum = map[string]uint64{}
			count = map[string]uint64{}
			min = map[string]uint64{}
			max = map[string]uint64{}
			prov = map[string]kt.Provider{}
		}
	}
}

func (r *StatsRollup) exportSum(sum map[string]uint64, count map[string]uint64, min map[string]uint64, max map[string]uint64, prov map[string]kt.Provider, rc chan []Rollup) {
	if len(sum) == 0 {
		rc <- nil
		return
	}

	ot := r.dtime
	r.dtime = time.Now()
	for _, name := range r.nameSet {
		r.Infof("XXX %v", name)
		keys := make([]Rollup, 0, len(sum))
		total := uint64(0)
		totalc := uint64(0)
		for k, v := range sum {
			keys = append(keys, Rollup{
				Name: name, EventType: r.eventType, Dimension: k,
				Metric: float64(v), KeyJoin: r.keyJoin, dims: combo(r.multiDims), Interval: r.dtime.Sub(ot),
				Count: count[k], Min: min[k], Max: max[k], Provider: prov[k],
			})
			total += v
			totalc += count[k]
		}

		sort.Sort(byValue(keys))
		if r.config.TopK > 0 && len(keys) > r.config.TopK {
			rc <- keys[0:r.config.TopK] // Return only the expected number, as sorted.
		} else {
			rc <- keys
		}
	}

	return
}
