package rollup

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
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
	rollups string
	keyJoin string
	topK    int
)

func init() {
	flag.StringVar(&rollups, "rollups", "", "Any rollups to use. Format: type, name, metric, dimension 1, dimension 2, ..., dimension n: sum,bytes,in_bytes,dst_addr")
	flag.StringVar(&keyJoin, "rollup_key_join", "^", "Token to use to join dimension keys together")
	flag.IntVar(&topK, "rollup_top_k", 10, "Export only these top values")
}

type Roller interface {
	Add([]map[string]interface{})
	Export() []Rollup
}

type Rollup struct {
	Dimension string        `json:"dimension"`
	Metric    float64       `json:"metric"`
	EventType string        `json:"eventType"`
	KeyJoin   string        `json:"keyJoin"`
	Interval  time.Duration `json:"interval"`
	dims      []string
	Name      string
	Count     uint64
	Min       uint64
	Max       uint64
	Provider  kt.Provider
}

type Method string

type RollupDef struct {
	Sample     bool
	Method     Method
	Metrics    []string
	Dimensions []string
	Name       string
}

func (r *RollupDef) String() string {
	return fmt.Sprintf("Name: %s, Method: %s, Adjust Sample Rate: %v, Metric: %v, Dimensions: %v", r.Name, r.Method, r.Sample, r.Metrics, r.Dimensions)
}

type RollupDefs []RollupDef

func (rf *RollupDefs) String() string {
	pts := make([]string, len(*rf))
	for i, r := range *rf {
		pts[i] = r.String()
	}
	return strings.Join(pts, "\n")
}

func (i *RollupDefs) Set(value string) error {
	pts := strings.Split(value, ",")
	if len(pts) < 3 {
		return fmt.Errorf("Rollup flag is defined by type, name, metric, dimension 1, dimension 2, ..., dimension n")
	}
	ptn := make([]string, len(pts))
	for i, p := range pts {
		ptn[i] = strings.TrimSpace(p)
	}
	if len(ptn[0]) > 2 && ptn[0][0:2] == "s_" {
		*i = append(*i, RollupDef{
			Method:     Method(ptn[0][2:]),
			Name:       ptn[1],
			Metrics:    strings.Split(ptn[2], "+"),
			Dimensions: ptn[3:],
			Sample:     true,
		})
	} else {
		*i = append(*i, RollupDef{
			Method:     Method(ptn[0]),
			Name:       ptn[1],
			Metrics:    strings.Split(ptn[2], "+"),
			Dimensions: ptn[3:],
		})
	}
	return nil
}

func GetRollups(log logger.Underlying, cfg *ktranslate.RollupConfig) ([]Roller, error) {
	rollups := RollupDefs{}
	for _, r := range cfg.Formats {
		if err := rollups.Set(r); err != nil {
			return nil, err
		}
	}
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
			statr, err := newStatsRollup(log, rf, cfg)
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
	dtime        time.Time
	name         string
	primaryDim   int
}

func (r *rollupBase) init(rd RollupDef) error {
	r.metrics = make([]string, 0)
	r.multiMetrics = make([][]string, 0)
	r.dims = make([]string, 0)
	r.multiDims = make([][]string, 0)
	r.dtime = time.Now()
	r.name = rd.Name
	r.eventType = strings.ReplaceAll(fmt.Sprintf(KENTIK_EVENT_TYPE, strings.Join(rd.Metrics, "_"), strings.Join(rd.Dimensions, ":")), ".", "_")
	r.sample = rd.Sample

	isMultiPrimary := false
	for i, d := range rd.Dimensions {
		pts := strings.Split(d, ".")
		switch len(pts) {
		case 1:
			r.dims = append(r.dims, d)
		case 2:
			r.multiDims = append(r.multiDims, pts)
			if i == 0 {
				isMultiPrimary = true
			}
		default:
			return fmt.Errorf("Invalid dimension: %s", d)
		}
	}

	if isMultiPrimary { // How do we sort by?
		r.primaryDim = len(r.dims)
	} else {
		r.primaryDim = 0
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
				if keyPts[next] == "" {
					if strings.HasPrefix(d[1], "source_") {
						keyPts[next] = dd["dest_"+d[1][7:]]
					} else if strings.HasPrefix(d[1], "dest_") {
						keyPts[next] = dd["source_"+d[1][5:]]
					}
				}
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

func (r *Rollup) GetDims() []string {
	return r.dims
}

func combo(dims []string, multiDims [][]string) []string {
	ret := make([]string, len(dims))
	for i, d := range dims {
		ret[i] = d
	}
	for _, d := range multiDims {
		ret = append(ret, d[1])
	}
	return ret
}
