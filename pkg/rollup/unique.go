package rollup

import (
	"encoding/binary"
	"sort"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"

	"github.com/dchest/siphash"
	gohll "github.com/sasha-s/go-hll"
)

var (
	hllErrRateHigh        = 0.02 // Matching whats used in prod
	hllSizeForErrRateHigh = 3080 // See verifyHLLConstantsOrPanic.
	sipk0                 = uint64(0)
	sipk1                 = uint64(0)
)

type UniqueRollup struct {
	logger.ContextL
	rollupBase
	state map[string]gohll.HLL
}

func newUniqueRollup(log logger.Underlying, rd RollupDef) (*UniqueRollup, error) {
	r := &UniqueRollup{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "uniqueRollup"}, log),
		state:    map[string]gohll.HLL{},
	}

	err := r.init(rd)
	if err != nil {
		return nil, err
	}
	r.Infof("New Rollup: %s -> %s", r.eventType, rd.String())
	return r, nil
}

func (r *UniqueRollup) Add(in []map[string]interface{}) {
	toAdd := map[string]gohll.HLL{}
	for _, mapr := range in {
		key := r.getKey(mapr)
		if _, ok := toAdd[key]; !ok {
			toAdd[key] = make(gohll.HLL, hllSizeForErrRateHigh)
		}

		for _, metric := range r.metrics {
			if mm, ok := mapr[metric]; ok {
				switch m := mm.(type) {
				case string:
					toAdd[key].Add(siphash.Hash(sipk0, sipk1, []byte(m)))
				case int64:
					target := make([]byte, binary.MaxVarintLen64)
					len := binary.PutVarint(target, m)
					toAdd[key].Add(siphash.Hash(sipk0, sipk1, target[0:len]))
				}
			}
		}
		for _, m := range r.multiMetrics { // Now handle the 2 level deep metrics
			if m1, ok := mapr[m[0]]; ok {
				switch mm := m1.(type) {
				case map[string]string:
					toAdd[key].Add(siphash.Hash(sipk0, sipk1, []byte(mm[m[1]])))
				case map[string]int32:
					target := make([]byte, binary.MaxVarintLen64)
					len := binary.PutVarint(target, int64(mm[m[1]]))
					toAdd[key].Add(siphash.Hash(sipk0, sipk1, target[0:len]))
				case map[string]int64:
					target := make([]byte, binary.MaxVarintLen64)
					len := binary.PutVarint(target, mm[m[1]])
					toAdd[key].Add(siphash.Hash(sipk0, sipk1, target[0:len]))
				}
			}
		}
	}

	// Now need a lock for actually updating the current rollup state.
	r.mux.Lock()
	for k, v := range toAdd {
		if _, ok := r.state[k]; ok {
			r.state[k].Merge(v)
		} else {
			r.state[k] = v
		}
	}
	r.mux.Unlock()
}

func (r *UniqueRollup) Export() []Rollup {
	r.mux.Lock()
	os := r.state
	r.state = map[string]gohll.HLL{}
	r.mux.Unlock()

	ot := r.dtime
	r.dtime = time.Now()
	keys := make([]Rollup, len(os))
	next := 0
	for k, v := range os {
		keys[next] = Rollup{Name: r.name, EventType: r.eventType, Dimension: k, Metric: float64(v.EstimateCardinality()), KeyJoin: r.keyJoin, dims: combo(r.dims, r.multiDims), Interval: r.dtime.Sub(ot)}
		next++
	}

	sort.Sort(byValue(keys))
	if len(keys) > r.topK {
		return keys[0:r.topK]
	}

	return keys
}
