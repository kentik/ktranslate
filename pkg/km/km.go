package km

import (
	"context"
	"strconv"
	"strings"

	nmspb "github.com/kentik/api-schema-public/gen/go/kentik/nms/v202505"
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

func LoadMapper(ctx context.Context, log logger.Underlying, apic *kkapi.KentikApi) (*KMMapper, error) {
	ftm := KMMapper{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "kmMapper"}, log),
		apic:     apic,
		metrics:  map[int64]KentikMetric{},
	}

	num, err := ftm.loadMetrics(ctx)
	if err != nil {
		return nil, err
	}

	ftm.Infof("Loaded %d Kentik Metric Definitions.", num)

	return &ftm, nil
}

type KentikMetric struct {
	Id         int64
	EntityType string
	Def        KMDef
}

type KMDef struct {
	Metrics     map[string]*nmspb.FieldSchema
	Dimensions  map[string]*nmspb.FieldSchema
	Description string
}

func (km *KMMapper) loadMetrics(ctx context.Context) (int, error) {
	mets, err := km.apic.GetKMetricsSchema(ctx)
	if err != nil {
		return 0, err
	}

	for _, met := range mets {
		k := KentikMetric{
			Id:         int64(met.GetId()),
			EntityType: met.GetEntityType(),
			Def: KMDef{
				Metrics:     met.GetMetrics(),
				Dimensions:  met.GetDimensions(),
				Description: met.GetDescription(),
			},
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

func isFieldPresentForMask(fieldMask int64, incomingMask int64, hasIncomingMask bool) bool {
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
		if dim.GetLabel() == "" || !isFieldPresentForMask(dim.Mask, dimensionMask, hasDimensionMask) {
			continue
		}

		if val, ok := msg.CustomStr[dim.Column]; ok {
			msg.CustomStr[dim.GetLabel()] = val
			delete(msg.CustomStr, dim.Column)
		}
	}

	for _, met := range definition.Def.Metrics {
		if met.GetLabel() == "" || !isFieldPresentForMask(met.Mask, metricMask, hasMetricMask) {
			continue
		}

		if mvar, ok := msg.CustomBigInt[met.Column]; ok {
			msg.CustomBigInt[met.GetLabel()] = mvar
			msg.CustomStr["unit"] = met.GetUnit()
			delete(msg.CustomBigInt, met.Column)
		} else if mvar, ok := msg.CustomInt[met.Column]; ok {
			msg.CustomInt[met.GetLabel()] = mvar
			msg.CustomStr["unit"] = met.GetUnit()
			delete(msg.CustomInt, met.Column)
		} else if mvar, ok := msg.CustomFloat[met.Column]; ok {
			msg.CustomFloat[met.GetLabel()] = mvar
			msg.CustomStr["unit"] = met.GetUnit()
			delete(msg.CustomFloat, met.Column)
		}
	}
}
