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

state = {"count": 0}

def main(n):
    i = 0
    for evt in n:
      print(evt.foo)
      evt.company_id = i
      evt.foo = "aaa"
      evt["foo-one"] = "aaa"
      i += 1
      # state["count"] += 1 # doesn't work for some reason

    print(json.encode({"foo":1,"bar":"123"}))
    print(state["count"])
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

	e, err := NewEnricher("", "", file.Name(), l.GetLogger().GetUnderlyingLogger())
	assert.Nil(err)
	assert.NotNil(e)

	out, err := e.Enrich(context.Background(), kt.InputTestingSnmp)
	assert.Nil(err)
	assert.NotNil(out)
	for i, evt := range kt.InputTestingSnmp {
		assert.Equal(i, int(evt.CompanyId))
		assert.Equal("aaa", evt.CustomStr["foo"])
		assert.Equal("aaa", evt.CustomStr["foo-one"])
	}
}

func TestNone(t *testing.T) {
	assert := assert.New(t)
	l := lt.NewTestContextL(logger.NilContext, t)

	_, err := NewEnricher("", "", "", l.GetLogger().GetUnderlyingLogger())
	assert.NotNil(err)
}

func TestSource(t *testing.T) {
	testDataPy := string(`
load("json.star", "json")

state = {"count": 0}

def main(n):
    i = 0
    for evt in n:
      print(evt.foo)
      evt.company_id = i
      evt.foo = "aaa"
      i += 1
      state["count"] += 1

    print(json.encode({"foo":1,"bar":"123"}))
    print(state["count"])
    return len(n)
`)

	assert := assert.New(t)
	l := lt.NewTestContextL(logger.NilContext, t)
	e, err := NewEnricher("", testDataPy, "", l.GetLogger().GetUnderlyingLogger())
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

func TestCatch(t *testing.T) {
	testDataPy := string(`

def bad():
  return 100/0

def main(n):
  res = catch(bad)
  return res
`)

	assert := assert.New(t)
	l := lt.NewTestContextL(logger.NilContext, t)
	e, err := NewEnricher("", testDataPy, "", l.GetLogger().GetUnderlyingLogger())
	assert.Nil(err)
	assert.NotNil(e)

	out, err := e.Enrich(context.Background(), kt.InputTestingSnmp)
	assert.Nil(err)
	assert.NotNil(out)
}

func TestRE(t *testing.T) {
	testDataPy := string(`
def main(n):
  res = findAllSubmatch("foo(.?)", "seafood fool")
  n[0]["foo"] = res[1][1]

  res = findAllSubmatch("foo(.?)", "seafood fool")
  n[0]["foo"] = res[1][0]

  return True
`)

	assert := assert.New(t)
	l := lt.NewTestContextL(logger.NilContext, t)
	e, err := NewEnricher("", testDataPy, "", l.GetLogger().GetUnderlyingLogger())
	assert.Nil(err)
	assert.NotNil(e)

	out, err := e.Enrich(context.Background(), kt.InputTestingSnmp)
	assert.Nil(err)
	assert.Equal("fool", out[0].CustomStr["foo"])
}
