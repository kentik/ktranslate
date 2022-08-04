package influx

import (
	"strings"
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

	prefix := "notempty"
	Prefix = &prefix
	res, err := f.To(kt.InputTesting, serBuf)
	assert.NoError(err)
	assert.NotNil(res)

	assert.Equal(byte('\n'), res.Body[len(res.Body)-1])

	pts := strings.Split(string(res.Body[:len(res.Body)-1]), "\n")
	assert.Equal(len(pts), len(kt.InputTesting))
}
