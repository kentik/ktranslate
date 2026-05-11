package km

import (
	"bufio"
	"encoding/json"
	"flag"
	"os"
	"strconv"
	"strings"

	kkapi "github.com/kentik/ktranslate/pkg/api"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
)

/**
Contains code to map the kentik metrics data into usable format.
*/

type KMMapper struct {
	logger.ContextL
	apic    *kkapi.KentikApi
	metrics map[int64]KentikMetric
}

var (
	tags string
)

func init() {
	flag.StringVar(&tags, "km_value_map", "", "CSV file mapping kentik metric values.")
}

func LoadMapper(log logger.Underlying, tagMapFilePath string, apic *kkapi.KentikApi) (*KMMapper, error) {
	ftm := KMMapper{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "kmMapper"}, log),
		apic:     apic,
		metrics:  map[int64]KentikMetric{},
	}

	num, err := ftm.loadFile(tagMapFilePath)
	if err != nil {
		return nil, err
	}

	ftm.Infof("Loaded %d Kentik Metric Definitions.", num)

	return &ftm, nil
}

type KentikMetric struct {
	Id   int64
	Name string
	Def  KMDef
}

type KMDef struct {
	Imports     []string
	Metrics     map[string]*KMMetric
	Dimensions  map[string]*KMDimension
	Description string
}

type KMMetric struct {
	Mask   int
	Type   string
	Unit   string
	Label  string
	Column string
}

type KMDimension struct {
	Mask   int
	Type   string
	Label  string
	Column string
}

func (km *KMMapper) loadFile(file string) (int, error) {
	f, err := os.Open(file)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		pts := strings.SplitN(scanner.Text(), "|", 3)
		if len(pts) != 3 {
			km.Warnf("Bad line, skipping %s", scanner.Text())
			continue
		}
		id, err := strconv.ParseInt(pts[0], 10, 64)
		if err != nil {
			km.Warnf("Bad id value, skipping: %s", pts[0])
			continue
		}
		k := KentikMetric{Id: id, Name: pts[1], Def: KMDef{}}
		err = json.Unmarshal([]byte(pts[2]), &k.Def)
		if err != nil {
			km.Warnf("Bad km definition, skipping: %s %v", pts[2], err)
			continue
		}

		if len(k.Def.Metrics) > 0 && len(k.Def.Dimensions) > 0 {
			for _, dim := range k.Def.Dimensions {
				dim.Column = kt.FixupName(dim.Column)
			}
			for _, met := range k.Def.Metrics {
				met.Column = kt.FixupName(met.Column)
			}
			km.metrics[k.Id] = k
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return len(km.metrics), nil
}

func readMask(msg *kt.JCHF, key string) (int64, bool) {
	if msg == nil {
		return 0, false
	}

	if mask, ok := msg.CustomBigInt[key]; ok {
		return mask, true
	}

	if mask, ok := msg.CustomInt[key]; ok {
		return int64(mask), true
	}

	if raw, ok := msg.CustomStr[key]; ok {
		mask, err := strconv.ParseInt(strings.TrimSpace(raw), 0, 64)
		if err == nil {
			return mask, true
		}
	}

	return 0, false
}

func isFieldPresentForMask(fieldMask int, incomingMask int64, hasIncomingMask bool) bool {
	if !hasIncomingMask {
		return true
	}

	if fieldMask == 0 {
		return true
	}

	return incomingMask&int64(fieldMask) != 0
}

func (km *KMMapper) Enrich(id int64, msg *kt.JCHF) {
	if km == nil || msg == nil {
		return
	}

	definition, ok := km.metrics[id]
	if !ok {
		return // Metric not found, just return.
	}

	dimensionMask, hasDimensionMask := readMask(msg, "km_dimension_mask")
	metricMask, hasMetricMask := readMask(msg, "km_metric_mask")

	for _, dim := range definition.Def.Dimensions {
		if dim.Label == "" || !isFieldPresentForMask(dim.Mask, dimensionMask, hasDimensionMask) {
			continue
		}

		if val, ok := msg.CustomStr[dim.Column]; ok {
			msg.CustomStr[dim.Label] = val
			delete(msg.CustomStr, dim.Column)
		}
	}

	for _, met := range definition.Def.Metrics {
		if met.Label == "" || !isFieldPresentForMask(met.Mask, metricMask, hasMetricMask) {
			continue
		}

		if mvar, ok := msg.CustomBigInt[met.Column]; ok {
			msg.CustomBigInt[met.Label] = mvar
			msg.CustomStr["unit"] = met.Unit
			delete(msg.CustomBigInt, met.Column)
		} else if mvar, ok := msg.CustomInt[met.Column]; ok {
			msg.CustomInt[met.Label] = mvar
			msg.CustomStr["unit"] = met.Unit
			delete(msg.CustomInt, met.Column)
		} else if mvar, ok := msg.CustomFloat[met.Column]; ok {
			msg.CustomFloat[met.Label] = mvar
			msg.CustomStr["unit"] = met.Unit
			delete(msg.CustomFloat, met.Column)
		}
	}
}
