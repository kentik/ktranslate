package nrm

import (
	"testing"
	"unicode/utf8"

	"github.com/kentik/ktranslate"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeMetricsUTF8(t *testing.T) {
	assert := assert.New(t)

	metrics := []NRMetric{
		{
			Value: "bad\xffvalue",
			Attributes: map[string]any{
				"mac_address": "bad\xffattr",
				"count":       int64(3),
				"ok":          "fine",
			},
		},
		{
			Value:      float64(42),
			Attributes: nil,
		},
	}

	sanitizeMetricsUTF8(metrics)

	assert.Equal("626164ff76616c7565", metrics[0].Value)
	assert.Equal("626164ff61747472", metrics[0].Attributes["mac_address"])
	assert.Equal(int64(3), metrics[0].Attributes["count"])
	assert.Equal("fine", metrics[0].Attributes["ok"])
	assert.Equal(float64(42), metrics[1].Value)
	assert.Nil(metrics[1].Attributes)
}

func TestToSanitizesInvalidUTF8(t *testing.T) {
	assert := assert.New(t)

	f, err := NewFormat(nil, kt.CompressionNone, &ktranslate.NRMFormatConfig{CustomAttributes: map[string]string{}})
	assert.NoError(err)

	in := kt.NewJCHF()
	in.SetMap()
	in.CompanyId = 10
	in.EventType = kt.KENTIK_EVENT_KTRANS_METRIC
	in.CustomStr["type"] = "counter"
	in.CustomStr["force"] = "true"
	in.CustomStr["name"] = "test_metric"
	in.CustomStr["mac_address"] = "bad\xffvalue"

	out, err := f.To([]*kt.JCHF{in}, nil)
	assert.NoError(err)
	if !assert.NotNil(out) {
		return
	}

	assert.True(utf8.Valid(out.Body), "output must be valid UTF-8, got %q", out.Body)

	var sets []NRMetricSet
	assert.NoError(json.Unmarshal(out.Body, &sets))
	if !assert.Len(sets, 1) || !assert.Len(sets[0].Metrics, 1) {
		return
	}

	assert.Equal("626164ff76616c7565", sets[0].Metrics[0].Attributes["mac_address"])
}

// Custom attributes must reach the shared Common block -- unlike per-device
// user_tags, this is the only path that also covers non-device-scoped
// metrics such as this heartbeat/self-instrumentation one. See NR-612348.
func TestNewNRCommonMergesCustomAttributes(t *testing.T) {
	assert := assert.New(t)

	f, err := NewFormat(nil, kt.CompressionNone, &ktranslate.NRMFormatConfig{
		CustomAttributes: map[string]string{"install_id": "test-instance-123"},
	})
	assert.NoError(err)

	in := kt.NewJCHF()
	in.SetMap()
	in.CompanyId = 10
	in.EventType = kt.KENTIK_EVENT_KTRANS_METRIC
	in.CustomStr["type"] = "counter"
	in.CustomStr["force"] = "true"
	in.CustomStr["name"] = "test_metric"

	out, err := f.To([]*kt.JCHF{in}, nil)
	assert.NoError(err)
	if !assert.NotNil(out) {
		return
	}

	var sets []NRMetricSet
	assert.NoError(json.Unmarshal(out.Body, &sets))
	if !assert.Len(sets, 1) {
		return
	}

	// The custom attribute is present alongside the existing hardcoded ones -- it must
	// never clobber instrumentation.provider/collector.name.
	assert.Equal("test-instance-123", sets[0].Common.Attributes["install_id"])
	assert.Equal(kt.InstProvider, sets[0].Common.Attributes[AttrInstrumentationProvider])
	assert.Equal(kt.CollectorName, sets[0].Common.Attributes[AttrCollectorName])
}

// A -nr_custom_attributes entry named after a reserved attribute must not be
// able to override it -- every metric batch would otherwise carry the wrong
// instrumentation.provider/collector.name.
func TestNewFormatDropsReservedCustomAttributeKeys(t *testing.T) {
	assert := assert.New(t)

	cfg := &ktranslate.NRMFormatConfig{
		CustomAttributes: map[string]string{
			AttrInstrumentationProvider: "hijacked",
			AttrCollectorName:           "hijacked",
			"install_id":                "test-instance-123",
		},
	}
	f, err := NewFormat(&lt.Test{T: t}, kt.CompressionNone, cfg)
	assert.NoError(err)

	common := f.newNRCommon()
	assert.Equal(kt.InstProvider, common.Attributes[AttrInstrumentationProvider])
	assert.Equal(kt.CollectorName, common.Attributes[AttrCollectorName])
	assert.Equal("test-instance-123", common.Attributes["install_id"])

	// The reserved keys are stripped from the config itself, not just masked per-call.
	_, hasProvider := cfg.CustomAttributes[AttrInstrumentationProvider]
	_, hasCollector := cfg.CustomAttributes[AttrCollectorName]
	assert.False(hasProvider)
	assert.False(hasCollector)
	assert.Equal("test-instance-123", cfg.CustomAttributes["install_id"])
}

func TestNewNRCommonWithNoCustomAttributes(t *testing.T) {
	assert := assert.New(t)

	f, err := NewFormat(nil, kt.CompressionNone, &ktranslate.NRMFormatConfig{CustomAttributes: map[string]string{}})
	assert.NoError(err)

	common := f.newNRCommon()
	assert.Equal(kt.InstProvider, common.Attributes[AttrInstrumentationProvider])
	assert.Equal(kt.CollectorName, common.Attributes[AttrCollectorName])
	assert.Len(common.Attributes, 2)
}
