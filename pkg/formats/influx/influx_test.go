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

	res, err := f.To(kt.InputTesting, serBuf)
	assert.NoError(err)
	assert.NotNil(res)

	pts := strings.Split(string(res.Body), "\n")
	assert.Equal(len(pts), len(kt.InputTesting))
}

func TestInfluxEscape(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"asdf", "asdf"},
		{"as df", "as\\ df"},
		{"as,df", "as\\,df"},
		{"as=df", "as\\=df"},
		{"as, df", "as\\,\\ df"},
	}
	for _, test := range tests {
		if got := influxEscape(test.input); got != test.want {
			t.Errorf("influxEscape(%q) = %q; want %q", test.input, got, test.want)
		}
	}
}
