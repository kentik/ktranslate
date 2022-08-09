package influx

import (
	"strings"
	"testing"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/stretchr/testify/assert"
)

func TestSeriToInflux(t *testing.T) {
	serBuf := make([]byte, 0)
	assert := assert.New(t)
	l := lt.NewTestContextL(logger.NilContext, t).GetLogger().GetUnderlyingLogger()

	f, err := NewFormat(l, go_metrics.DefaultRegistry, kt.CompressionNone)
	assert.NoError(err)

	prefix := "notempty"
	Prefix = &prefix
	res, err := f.To(kt.InputTesting, serBuf)
	assert.NoError(err)
	assert.NotNil(res)

	assert.Equal(byte('\n'), res.Body[len(res.Body)-1])

	pts := strings.Split(string(res.Body[:len(res.Body)-1]), "\n")
	assert.Equal(len(pts), len(kt.InputTesting))
}

func TestNewline(t *testing.T) {
	assert := assert.New(t)
	l := lt.NewTestContextL(logger.NilContext, t).GetLogger().GetUnderlyingLogger()

	f, err := NewFormat(l, go_metrics.DefaultRegistry, kt.CompressionNone)
	assert.NoError(err)

	input := InfluxDataSet{
		InfluxData{
			Name:        "measurement",
			FieldsFloat: map[string]float64{},
			Fields:      map[string]int64{"moltue": 42},
			Tags: map[string]interface{}{
				"tag1": `line1
line2`,
			},
			Timestamp: 0,
		},
	}
	res := f.Bytes(input)
	assert.NotNil(res)
	str := string(res)
	assert.Contains(str, "measurement")
	assert.Contains(str, "line1")
	assert.Contains(str, "line2")
}

const (
	clean = "the quick red fox jumps over the lazy brown dogs"
	dirty = `the quick red fox
			jumps over
			the lazy brown dogs`
)

func BenchmarkPrepareTagValueContainsAnyClean(b *testing.B) {
	benchmarkPrepareTagValueContainsAny(b, clean)
}

func BenchmarkPrepareTagValueContainsAnyDirty(b *testing.B) {
	benchmarkPrepareTagValueContainsAny(b, dirty)
}

func benchmarkPrepareTagValueContainsAny(b *testing.B, s string) {
	for i := 0; i < b.N; i++ {
		prepareTagValue(s)
	}
}

func BenchmarkPrepareTagValueMapClean(b *testing.B) {
	benchmarkPrepareTagValueMap(b, clean)
}

func BenchmarkPrepareTagValueMapDirty(b *testing.B) {
	benchmarkPrepareTagValueMap(b, dirty)
}

func benchmarkPrepareTagValueMap(b *testing.B, s string) {
	for i := 0; i < b.N; i++ {
		prepareTagValueMap(s)
	}
}
