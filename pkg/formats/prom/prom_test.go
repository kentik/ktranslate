package prom

import (
	"testing"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/stretchr/testify/assert"
)

func TestSeriToInflux(t *testing.T) {
	serBuf := make([]byte, 0)
	assert := assert.New(t)
	l := lt.NewTestContextL(logger.NilContext, t).GetLogger().GetUnderlyingLogger()

	f, err := NewFormat(l, kt.CompressionNone)
	assert.NoError(err)

	res, err := f.To(kt.InputTesting, serBuf)
	assert.NoError(err)
	assert.Nil(res)

	res, err = f.To(kt.InputTestingSynth, serBuf)
	assert.NoError(err)
	assert.Nil(res)

	res, err = f.To(kt.InputTestingSnmp, serBuf)
	assert.NoError(err)
	assert.Nil(res)
}
