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

	assert.Equal(byte('\n'), res.Body[len(res.Body)-1])

	pts := strings.Split(string(res.Body[:len(res.Body)-1]), "\n")
	assert.Equal(len(pts), len(kt.InputTesting))
}

func TestInfluxEscapeTag(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"asdf", "asdf"},
		{"as df", "as\\ df"},
		{"as,df", "as\\,df"},
		{"as=df", "as\\=df"},
		{"as, df", "as\\,\\ df"},
		{"as\ndf", "as\\\ndf"},
		{"as\r\ndf", "as\\\r\\\ndf"},
	}
	for _, test := range tests {
		if got := influxEscapeTag(test.input); got != test.want {
			t.Errorf("influxEscapeTag(%q) = %q; want %q", test.input, got, test.want)
		}
	}
}

func TestInfluxEscapeField(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"asdf", "asdf"},
		{"as df", "as df"},
		{"as\"df", "as\\\"df"},
		{"as\\df", "as\\\\df"},
	}
	for _, test := range tests {
		if got := influxEscapeField(test.input); got != test.want {
			t.Errorf("influxEscapeField(%q) = %q; want %q", test.input, got, test.want)
		}
	}
}
