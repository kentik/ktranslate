package nrm

import (
	"testing"
	"unicode/utf8"

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

	f, err := NewFormat(nil, kt.CompressionNone)
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
