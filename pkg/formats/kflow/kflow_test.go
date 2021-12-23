package kflow

import (
	"testing"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/stretchr/testify/assert"
)

func TestSeriToJflow(t *testing.T) {
	serBuf := make([]byte, 0)
	assert := assert.New(t)
	l := lt.NewTestContextL(logger.NilContext, t).GetLogger().GetUnderlyingLogger()

	f, err := NewFormat(l, kt.CompressionNone)
	assert.Error(err)

	f, err = NewFormat(l, kt.CompressionGzip)
	assert.NoError(err)

	res, err := f.To(kt.InputTesting, serBuf)
	assert.NoError(err)
	assert.NotNil(res)
	assert.True(len(res.Body) > 0)

	out, err := f.From(res)
	assert.NoError(err)
	assert.Equal(len(kt.InputTesting), len(out))
	for i, _ := range out {
		assert.Equal(kt.InputTesting[i].Timestamp, out[i]["timestamp"])
		for k, v := range kt.InputTesting[i].CustomStr {
			assert.Equal(v, out[i][k])
		}
	}
}
