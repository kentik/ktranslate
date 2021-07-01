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
	name := "name"

	res := copyAttrForSnmp(input, name)
	assert.Equal(12, len(res)) // adds in two keys
	assert.Equal("name", res["objectIdentifier"])

	for i := 0; i < MAX_ATTR_FOR_NR+10; i++ {
		input[fmt.Sprintf("XXX%d", i)] = i
	}
	res = copyAttrForSnmp(input, name)
	assert.Equal(MAX_ATTR_FOR_NR, len(res)) // truncated at MAX_ATTR_FOR_NR
	assert.Equal("name", res["objectIdentifier"])

	input = map[string]interface{}{kt.StringPrefix + "foo": "one"}
	res = copyAttrForSnmp(input, name)
	assert.Equal("one", res["foo"], res)
}
