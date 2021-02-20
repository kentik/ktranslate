package rollup

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

const (
	Sum            Method = "sum"
	Unique                = "unique"
	Min                   = "min"
	Max                   = "max"
	Mean                  = "mean"
	Median                = "median"
	Entropy               = "entropy"
	Percentilerank        = "percentilerank"
	Percentile            = "entropy"

	KENTIK_EVENT_TYPE = "KFlow:%s:%s"
)

var (
	rollups RollupFlag
	keyJoin = flag.String("rollup_key_join", "^", "Token to use to join dimension keys together")
	topK    = flag.Int("rollup_top_k", 10, "Export only these top values")
)

type Roller interface {
	Add([]*kt.JCHF)
	Export() []Rollup
}

type Rollup struct {
	Dimension string  `json:"dimension"`
	Metric    float64 `json:"metric"`
	EventType string  `json:"eventType"`
	KeyJoin   string  `json:"keyJoin"`
}

type Method string

type RollupDef struct {
	Sample     bool
	Method     Method
	Metrics    []string
	Dimensions []string
}

func (r *RollupDef) String() string {
	return fmt.Sprintf("Method: %s, Adjust Sample Rate: %v, Metric: %v, Dimensions: %v", r.Method, r.Sample, r.Metrics, r.Dimensions)
}

type RollupFlag []RollupDef

func (rf *RollupFlag) String() string {
	pts := make([]string, len(*rf))
	for i, r := range *rf {
		pts[i] = r.String()
	}
	return strings.Join(pts, "\n")
}

func (i *RollupFlag) Set(value string) error {
	pts := strings.Split(value, ",")
	if len(pts) < 3 {
		return fmt.Errorf("Rollup flag is defined by type, metric, dimension 1, dimension 2, ..., dimension n")
	}
	ptn := make([]string, len(pts))
	for i, p := range pts {
		ptn[i] = strings.TrimSpace(p)
	}
	if len(ptn[0]) > 2 && ptn[0][0:2] == "s_" {
		*i = append(*i, RollupDef{
			Method:     Method(ptn[0][2:]),
			Metrics:    strings.Split(ptn[1], "+"),
			Dimensions: ptn[2:],
			Sample:     true,
		})
	} else {
		*i = append(*i, RollupDef{
			Method:     Method(ptn[0]),
			Metrics:    strings.Split(ptn[1], "+"),
			Dimensions: ptn[2:],
		})
	}
	return nil
}

func GetRollups(log logger.Underlying) ([]Roller, error) {
	rolls := make([]Roller, 0)
	for _, rf := range rollups {
		switch rf.Method {
		case Unique:
			ur, err := newUniqueRollup(log, rf)
			if err != nil {
				return nil, err
			}
			rolls = append(rolls, ur)
		default:
			statr, err := newStatsRollup(log, rf)
			if err != nil {
				return nil, err
			}
			rolls = append(rolls, statr)
		}
	}

	return rolls, nil
}

type byValue []Rollup

func (a byValue) Len() int           { return len(a) }
func (a byValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byValue) Less(i, j int) bool { return a[j].Metric < a[i].Metric }

type rollupBase struct {
	metrics      []string
	multiMetrics [][]string
	dims         []string
	multiDims    [][]string
	keyJoin      string
	topK         int
	eventType    string
	mux          sync.RWMutex
	sample       bool
}

func (r *rollupBase) init(rd RollupDef) error {
	r.metrics = make([]string, 0)
	r.multiMetrics = make([][]string, 0)
	r.dims = make([]string, 0)
	r.multiDims = make([][]string, 0)
	r.keyJoin = *keyJoin
	r.topK = *topK
	r.eventType = strings.ReplaceAll(fmt.Sprintf(KENTIK_EVENT_TYPE, strings.Join(rd.Metrics, "_"), strings.Join(rd.Dimensions, ":")), ".", "_")
	r.sample = rd.Sample

	for _, d := range rd.Dimensions {
		pts := strings.Split(d, ".")
		switch len(pts) {
		case 1:
			r.dims = append(r.dims, d)
		case 2:
			r.multiDims = append(r.multiDims, pts)
		default:
			return fmt.Errorf("Invalid dimension: %s", d)
		}
	}

	for _, m := range rd.Metrics {
		pts := strings.Split(m, ".")
		switch len(pts) {
		case 1:
			r.metrics = append(r.metrics, m)
		case 2:
			r.multiMetrics = append(r.multiMetrics, pts)
		default:
			return fmt.Errorf("Invalid metric: %s", m)
		}
	}

	return nil
}

func (r *rollupBase) getKey(mapr map[string]interface{}) string {
	keyPts := make([]string, len(r.dims)+len(r.multiDims))
	for i, d := range r.dims {
		if dd, ok := mapr[d]; ok {
			switch v := dd.(type) {
			case string:
				keyPts[i] = v
			case int64:
				keyPts[i] = strconv.Itoa(int(v))
			default:
				// Skip?
			}
		}
	}
	next := len(r.dims)
	for _, d := range r.multiDims { // Now handle the 2 level deep maps
		if d1, ok := mapr[d[0]]; ok {
			switch dd := d1.(type) {
			case map[string]string:
				keyPts[next] = dd[d[1]]
			case map[string]int32:
				keyPts[next] = strconv.Itoa(int(dd[d[1]]))
			case map[string]int64:
				keyPts[next] = strconv.Itoa(int(dd[d[1]]))
			}
		}
		next++
	}

	return strings.Join(keyPts, r.keyJoin)
}
