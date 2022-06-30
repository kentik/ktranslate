package filter

import (
	"testing"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	l := lt.NewTestContextL(logger.NilContext, t).GetLogger().GetUnderlyingLogger()
	assert := assert.New(t)
	filters = []FilterDefWrapper{
		[]FilterDef{
			FilterDef{
				Dimension: "src_addr",
				Operator:  "==",
				Value:     "10.2.2.1",
				FType:     "string",
			},
			FilterDef{
				Dimension: "custom_bigint.fooII",
				Operator:  "==",
				Value:     "12",
				FType:     "int",
			},
			FilterDef{
				Dimension: "src_addr",
				Operator:  "%",
				Value:     "10.2.2.0/24",
				FType:     "addr",
			},
			FilterDef{
				Dimension: "src_addr",
				Operator:  "%",
				Value:     "10.2.3.0/24",
				FType:     "addr",
			},
			FilterDef{
				Dimension: "custom_bigint.foo",
				Operator:  "!=",
				Value:     "13",
				FType:     "int",
			},
			FilterDef{
				Dimension: "src_addr",
				Operator:  "%",
				Value:     "10.2",
				FType:     "string",
			},
		},
	}
	fs, err := GetFilters(l)
	assert.NoError(err)
	assert.Equal(len(filters), len(fs))

	results := []bool{true, true, true, false, true, true}
	for i, fs := range fs {
		assert.Equal(results[i], fs.Filter(kt.InputTesting[0]), "%d -> %v", i, filters[i])
	}
}
