package rollup

import (
	"encoding/binary"
	"sort"
	"strings"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"

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
	kvs       chan *uset
	exportKvs chan chan []Rollup
}

type uset struct {
	uniques map[string]gohll.HLL
	count   map[string]uint64
	prov    map[string]kt.Provider
}

func newUniqueRollup(log logger.Underlying, rd RollupDef) (*UniqueRollup, error) {
	r := &UniqueRollup{
		ContextL:  logger.NewContextLFromUnderlying(logger.SContext{S: "uniqueRollup"}, log),
		kvs:       make(chan *uset, CHAN_SLACK),
		exportKvs: make(chan chan []Rollup),
	}

	err := r.init(rd)
	if err != nil {
		return nil, err
	}
	go r.uniqueKvs()

	r.Infof("New Rollup: %s -> %s", r.eventType, rd.String())
	return r, nil
}

func (r *UniqueRollup) Add(in []map[string]interface{}) {
	uniques := map[string]gohll.HLL{}
	count := map[string]uint64{}
	prov := map[string]kt.Provider{}

	for _, mapr := range in {
		key := r.getKey(mapr)
		count[key]++
		prov[key] = mapr["provider"].(kt.Provider)
		if _, ok := uniques[key]; !ok {
			uniques[key] = make(gohll.HLL, hllSizeForErrRateHigh)
		}

		for _, metric := range r.metrics {
			if mm, ok := mapr[metric]; ok {
				switch m := mm.(type) {
				case string:
					uniques[key].Add(siphash.Hash(sipk0, sipk1, []byte(m)))
				case int64:
					target := make([]byte, binary.MaxVarintLen64)
					len := binary.PutVarint(target, m)
					uniques[key].Add(siphash.Hash(sipk0, sipk1, target[0:len]))
				}
			}
		}
		for _, m := range r.multiMetrics { // Now handle the 2 level deep metrics
			if m1, ok := mapr[m[0]]; ok {
				switch mm := m1.(type) {
				case map[string]string:
					uniques[key].Add(siphash.Hash(sipk0, sipk1, []byte(mm[m[1]])))
				case map[string]int32:
					target := make([]byte, binary.MaxVarintLen64)
					len := binary.PutVarint(target, int64(mm[m[1]]))
					uniques[key].Add(siphash.Hash(sipk0, sipk1, target[0:len]))
				case map[string]int64:
					target := make([]byte, binary.MaxVarintLen64)
					len := binary.PutVarint(target, mm[m[1]])
					uniques[key].Add(siphash.Hash(sipk0, sipk1, target[0:len]))
				}
			}
		}
	}

	// Pass on to be globally agregated.
	select {
	case r.kvs <- &uset{uniques: uniques, count: count, prov: prov}:
	default:
		r.Warnf("kvs chan full")
	}
}

func (r *UniqueRollup) uniqueKvs() {
	uniques := map[string]gohll.HLL{}
	count := map[string]uint64{}
	prov := map[string]kt.Provider{}

	for {
		select {
		case itm := <-r.kvs: // Just add to our map
			for k, v := range itm.uniques {
				if _, ok := uniques[k]; ok {
					uniques[k].Merge(v)
				} else {
					uniques[k] = v
				}
			}
			for k, v := range itm.count {
				count[k] += v
			}
			for k, v := range itm.prov {
				prov[k] = v
			}
		case rc := <-r.exportKvs: // Return the top results
			go r.exportUnique(uniques, count, prov, rc)
			uniques = map[string]gohll.HLL{}
			count = map[string]uint64{}
			prov = map[string]kt.Provider{}
		}
	}
}

func (r *UniqueRollup) Export() []Rollup {
	rc := make(chan []Rollup)
	r.exportKvs <- rc
	return <-rc
}

func (r *UniqueRollup) exportUnique(uniques map[string]gohll.HLL, count map[string]uint64, prov map[string]kt.Provider, rc chan []Rollup) {
	if len(uniques) == 0 {
		rc <- nil
		return
	}

	ot := r.dtime
	r.dtime = time.Now()
	keys := make([]Rollup, 0, len(uniques))
	totalc := uint64(0)
	var provt kt.Provider
	for k, v := range uniques {
		keys = append(keys, Rollup{
			Name: r.name, EventType: r.eventType, Dimension: k,
			Metric: float64(v.EstimateCardinality()), KeyJoin: r.keyJoin, dims: combo(r.dims, r.multiDims), Interval: r.dtime.Sub(ot),
			Count: count[k], Provider: prov[k],
		})
		totalc += count[k]
		provt = prov[k]
	}

	sort.Sort(byValue(keys))
	if len(keys) > r.topK {
		r.getTopkUniques(keys, totalc, ot, provt, rc)
	} else {
		rc <- keys
	}

	return
}

func (r *UniqueRollup) getTopkUniques(keys []Rollup, totalc uint64, ot time.Time, prov kt.Provider, rc chan []Rollup) {
	top := make([]Rollup, 0, len(keys))
	seen := map[string]int{}

	for _, roll := range keys {
		pts := strings.Split(roll.Dimension, r.keyJoin)
		if seen[pts[r.primaryDim]] < r.topK { // If the primary key for this rollup has less than the topk set, add it to the list.
			top = append(top, roll)
		}
		seen[pts[r.primaryDim]]++
	}

	rc <- top
}
