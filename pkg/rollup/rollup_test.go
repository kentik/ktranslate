package rollup

import (
	"strings"
	"testing"
	"time"

	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
	"github.com/kentik/ktranslate/pkg/filter"
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/stretchr/testify/assert"
)

func TestRollup(t *testing.T) {
	l := lt.NewTestContextL(logger.NilContext, t).GetLogger().GetUnderlyingLogger()
	assert := assert.New(t)
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
		ktranslate.RollupConfig{
			JoinKey:       "^",
			TopK:          1,
			Formats:       []string{"sum,sum_bytes_in,in_bytes,ccc,custom_str.foo,bar"},
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
		[]map[string]interface{}{
			map[string]interface{}{
				"in_bytes":    int64(10),
				"custom_str":  map[string]string{"foo": "ddd"},
				"bar":         "ccc",
				"ccc":         "fff",
				"sample_rate": int64(1),
				"provider":    kt.Provider("pp"),
			},
			map[string]interface{}{
				"in_bytes":    int64(65),
				"custom_str":  map[string]string{"foo": "ddd"},
				"bar":         "ccc",
				"ccc":         "fff",
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
		map[string]interface{}{
			"metric":     75,
			"dimensions": []string{"fff", "ddd", "ccc"},
		},
	}

	for i, roll := range rolls {
		rd, err := GetRollups(l, &roll)
		assert.NoError(err)
		assert.Equal(len(roll.Formats), len(rd))

		for _, ri := range rd {
			ri.Add(inputs[i])
			time.Sleep(50 * time.Microsecond)
			res := ri.Export()

			assert.Equal(outputs[i]["metric"].(int), int(res[0].Metric), res)
			assert.Equal(roll.TopK, len(res), i)
			dims := strings.Split(res[0].Dimension, res[0].KeyJoin)
			for j, dim := range dims {
				assert.Equal(outputs[i]["dimensions"].([]string)[j], dim, res)
			}
		}
	}
}

func TestRollupFilter(t *testing.T) {
	l := lt.NewTestContextL(logger.NilContext, t).GetLogger().GetUnderlyingLogger()
	assert := assert.New(t)
	rolls := []ktranslate.RollupConfig{
		ktranslate.RollupConfig{
			JoinKey:       "^",
			TopK:          1,
			Formats:       []string{"s_sum,name_one,in_bytes,foo"},
			KeepUndefined: true,
		},
	}

	inputs := [][]map[string]interface{}{
		[]map[string]interface{}{
			map[string]interface{}{
				"in_bytes":    int64(5),
				"foo":         "aaa",
				"filter":      int64(1),
				"sample_rate": int64(1),
				"provider":    kt.Provider("pp"),
			},
			map[string]interface{}{
				"in_bytes":    int64(15),
				"foo":         "aaa",
				"filter":      int64(1),
				"sample_rate": int64(1),
				"provider":    kt.Provider("pp"),
			},
			map[string]interface{}{
				"in_bytes":    int64(20),
				"foo":         "aaa",
				"filter":      int64(2),
				"sample_rate": int64(1),
				"provider":    kt.Provider("pp"),
			},
			map[string]interface{}{
				"in_bytes":    int64(2),
				"foo":         "aaa",
				"filter":      int64(2),
				"sample_rate": int64(1),
				"provider":    kt.Provider("pp"),
			},
		},
	}

	filters := [][]string{
		[]string{
			"int,filter,==,1,name_one",
		},
	}

	outputs := []map[string]interface{}{
		map[string]interface{}{
			"metric":     20,
			"dimensions": []string{"aaa"},
		},
	}

	for i, roll := range rolls {
		rd, err := GetRollups(l, &roll)
		assert.NoError(err)
		assert.Equal(len(roll.Formats), len(rd))

		fs, err := filter.GetFilters(l, filters[i])
		assert.NoError(err)

		fullSet := []filter.FilterWrapper{}
		for _, filter := range fs {
			if filter.GetName() == "" { // No name means a global application.
				fullSet = append(fullSet, filter)
				continue
			}

			found := false
			for _, ri := range rd {
				if filter.GetName() == ri.GetName() {
					ri.SetFilter(filter)
					found = true
				}
			}
			if !found {
				t.Errorf("No name match for filter %v", filter.GetName())
			}
		}
		assert.Equal(0, len(fullSet))

		for _, ri := range rd {
			ri.Add(inputs[i])
			time.Sleep(50 * time.Microsecond)
			res := ri.Export()

			assert.Equal(roll.TopK, len(res), i)
			assert.Equal(outputs[i]["metric"].(int), int(res[0].Metric), res)
			dims := strings.Split(res[0].Dimension, res[0].KeyJoin)
			for j, dim := range dims {
				assert.Equal(outputs[i]["dimensions"].([]string)[j], dim, res)
			}
		}
	}
}

func BenchmarkRollups(b *testing.B) {
	l := lt.NewBenchContextL(logger.NilContext, b).GetLogger().GetUnderlyingLogger()
	assert := assert.New(b)
	// filters are type,dimension,operator,value
	roll := ktranslate.RollupConfig{
		JoinKey:       "^",
		TopK:          2,
		Formats:       []string{"sum,sum_bytes_in,in_bytes,custom_str.foo,bar"},
		KeepUndefined: false,
	}

	inputs := []map[string]interface{}{
		map[string]interface{}{
			"in_bytes":    int64(10),
			"custom_str":  map[string]string{"foo": "ddd"},
			"bar":         "ccc",
			"sample_rate": int64(1),
			"provider":    kt.Provider("pp"),
		},
		map[string]interface{}{
			"in_bytes":    int64(20),
			"custom_str":  map[string]string{"foo": "ddd"},
			"bar":         "ccc",
			"sample_rate": int64(1),
			"provider":    kt.Provider("pp"),
		},
		map[string]interface{}{
			"in_bytes":    int64(20),
			"custom_str":  map[string]string{"foo": "ddd"},
			"bar":         "eee",
			"sample_rate": int64(1),
			"provider":    kt.Provider("pp"),
		},
		map[string]interface{}{
			"in_bytes":    int64(2),
			"custom_str":  map[string]string{"foo": "ddd"},
			"bar":         "fff",
			"sample_rate": int64(1),
			"provider":    kt.Provider("pp"),
		},
	}

	rd, err := GetRollups(l, &roll)
	assert.NoError(err)
	for n := 0; n < b.N; n++ {
		for _, ri := range rd {
			ri.Add(inputs)
			time.Sleep(50 * time.Microsecond)
			res := ri.Export()
			assert.Equal(roll.TopK, len(res))
		}
	}
}
