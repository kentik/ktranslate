package rollup

import (
	"strings"
	"testing"

	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/stretchr/testify/assert"
)

func TestRollup(t *testing.T) {
	l := lt.NewTestContextL(logger.NilContext, t).GetLogger().GetUnderlyingLogger()
	assert := assert.New(t)
	// filters are type,dimension,operator,value
	rolls := []ktranslate.RollupConfig{
		ktranslate.RollupConfig{
			JoinKey:       "^",
			TopK:          2,
			Formats:       []string{"sum,sum_bytes_in,in_bytes,foo"},
			KeepUndefined: false,
		},
		ktranslate.RollupConfig{
			JoinKey:       "^",
			TopK:          1,
			Formats:       []string{"sum,sum_bytes_in,in_bytes,foo,bar"},
			KeepUndefined: false,
		},
		ktranslate.RollupConfig{
			JoinKey:       "^",
			TopK:          1,
			Formats:       []string{"sum,sum_bytes_in,in_bytes,custom_str.foo,bar"},
			KeepUndefined: false,
		},
	}

	inputs := [][]map[string]interface{}{
		[]map[string]interface{}{
			map[string]interface{}{
				"in_bytes":    int64(10),
				"foo":         "aaa",
				"sample_rate": int64(1),
				"provider":    kt.Provider("pp"),
			},
			map[string]interface{}{
				"in_bytes":    int64(20),
				"foo":         "aaa",
				"sample_rate": int64(1),
				"provider":    kt.Provider("pp"),
			},
			map[string]interface{}{
				"in_bytes":    int64(20),
				"foo":         "bbb",
				"sample_rate": int64(1),
				"provider":    kt.Provider("pp"),
			},
			map[string]interface{}{
				"in_bytes":    int64(2),
				"foo":         "ccc",
				"sample_rate": int64(1),
				"provider":    kt.Provider("pp"),
			},
		},
		[]map[string]interface{}{
			map[string]interface{}{
				"in_bytes":    int64(10),
				"foo":         "bbb",
				"bar":         "ccc",
				"sample_rate": int64(1),
				"provider":    kt.Provider("pp"),
			},
			map[string]interface{}{
				"in_bytes":    int64(34),
				"foo":         "bbb",
				"bar":         "ccc",
				"sample_rate": int64(1),
				"provider":    kt.Provider("pp"),
			},
		},
		[]map[string]interface{}{
			map[string]interface{}{
				"in_bytes":    int64(10),
				"custom_str":  map[string]string{"foo": "ddd"},
				"bar":         "ccc",
				"sample_rate": int64(1),
				"provider":    kt.Provider("pp"),
			},
			map[string]interface{}{
				"in_bytes":    int64(55),
				"custom_str":  map[string]string{"foo": "ddd"},
				"bar":         "ccc",
				"sample_rate": int64(1),
				"provider":    kt.Provider("pp"),
			},
		},
	}

	outputs := []map[string]interface{}{
		map[string]interface{}{
			"metric":     30,
			"dimensions": []string{"aaa"},
		},
		map[string]interface{}{
			"metric":     44,
			"dimensions": []string{"bbb", "ccc"},
		},
		map[string]interface{}{
			"metric":     65,
			"dimensions": []string{"ddd", "ccc"},
		},
	}

	for i, roll := range rolls {
		rd, err := GetRollups(l, &roll)
		assert.NoError(err)
		assert.Equal(len(roll.Formats), len(rd))

		for _, ri := range rd {
			ri.Add(inputs[i])
			res := ri.Export()

			if len(res) > 0 {
				assert.Equal(outputs[i]["metric"].(int), int(res[0].Metric), res)
			}
			assert.Equal(roll.TopK, len(res), i)
			dims := strings.Split(res[0].Dimension, res[0].KeyJoin)
			for j, dim := range dims {
				assert.Equal(outputs[i]["dimensions"].([]string)[j], dim, res)
			}
		}
	}
}
