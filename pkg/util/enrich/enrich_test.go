package enrich

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
	"github.com/kentik/ktranslate/pkg/kt"
)

func TestPython(t *testing.T) {

	testDataPy := []byte(`
load("json.star", "json")

def main(n):
    i = 0
    for evt in n:
      print(evt.foo)
      evt.company_id = i
      evt.foo = "aaa"
      i += 1

    print(json.encode({"foo":1,"bar":"123"}))
    return len(n)
`)

	assert := assert.New(t)
	l := lt.NewTestContextL(logger.NilContext, t)
	file, err := os.CreateTemp("", "")
	if err != nil {
		t.FailNow()
	}
	file.Write(testDataPy)
	file.Sync()
	defer os.Remove(file.Name())

	e, err := NewEnricher(file.Name(), l.GetLogger().GetUnderlyingLogger())
	assert.Nil(err)
	assert.NotNil(e)

	out, err := e.Enrich(context.Background(), kt.InputTestingSnmp)
	assert.Nil(err)
	assert.NotNil(out)
	for i, evt := range kt.InputTestingSnmp {
		assert.Equal(i, int(evt.CompanyId))
		assert.Equal("aaa", evt.CustomStr["foo"])
	}
}
