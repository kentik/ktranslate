package nrm

import (
	"fmt"
	"testing"

	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/stretchr/testify/assert"
)

func TestCopyAttrforSNMP(t *testing.T) {
	assert := assert.New(t)

	input := map[string]interface{}{}
	for i := 0; i < 10; i++ {
		input[fmt.Sprintf("XXX%d", i)] = i
	}
	name := kt.MetricInfo{Oid: "oid", Mib: "mib"}

	res := copyAttrForSnmp(input, "test", name, nil)
	assert.Equal(len(input)+3, len(res)) // adds in three keys
	assert.Equal("oid", res["objectIdentifier"])

	for i := 0; i < MAX_ATTR_FOR_NR+10; i++ {
		input[fmt.Sprintf("XXX%d", i)] = i
	}
	res = copyAttrForSnmp(input, "test", name, nil)
	assert.Equal(MAX_ATTR_FOR_NR, len(res)) // truncated at MAX_ATTR_FOR_NR
	assert.Equal("oid", res["objectIdentifier"])

	input = map[string]interface{}{kt.StringPrefix + "foo": "one"}
	res = copyAttrForSnmp(input, "test", name, nil)
	assert.Equal("one", res["foo"], res)

	input = map[string]interface{}{kt.StringPrefix + "foo": "one"}
	name = kt.MetricInfo{Oid: "oid", Mib: "mib", Table: "noMatch"}
	res = copyAttrForSnmp(input, "test", name, nil)
	assert.Equal(nil, res["foo"], res)

	input = map[string]interface{}{kt.StringPrefix + "foo": "one"}
	name = kt.MetricInfo{Oid: "oid", Mib: "mib", Table: "foo"}
	res = copyAttrForSnmp(input, "test", name, nil)
	assert.Equal("one", res["foo"], res)

	input = map[string]interface{}{"foo": "one"}
	name = kt.MetricInfo{Oid: "oid", Mib: "mib", Table: "foo"}
	res = copyAttrForSnmp(input, "test", name, nil)
	assert.Equal("one", res["foo"], res)

	input = map[string]interface{}{"foo": "one"}
	name = kt.MetricInfo{Oid: "oid", Mib: "mib", Table: "somethingElse"}
	res = copyAttrForSnmp(input, "test", name, nil)
	assert.Equal(nil, res["foo"], res)
}
