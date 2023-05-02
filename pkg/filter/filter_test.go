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
	// filters are type,dimension,operator,value
	filters := []string{
		"string,src_addr,==,10.2.2.1",
		"int,custom_bigint.fooII,==,12",
		"addr,src_addr,%,10.2.2.0/24",
		"addr,src_addr,%,10.2.3.0/24",
		"int,custom_bigint.foo,!=,13",
		"string,src_addr,%,10.2",
		"int,fooII,==,12",
		"string,foo,==,bar",
		"string,fooAAAA,==,",
	}
	fs, err := GetFilters(l, filters)
	assert.NoError(err)
	assert.Equal(len(filters), len(fs))

	results := []bool{true, true, true, false, true, true, true, true, false}
	for i, fs := range fs {
		assert.Equal(results[i], fs.Filter(kt.InputTesting[0]), "%d -> %v", i, filters[i])
	}
}
