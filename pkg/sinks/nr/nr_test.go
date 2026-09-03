package nr

import (
	"bytes"
	"compress/gzip"
	"testing"

	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/stretchr/testify/assert"
)

func TestCheckJson(t *testing.T) {
	assert := assert.New(t)

	// Good no compression case
	input := "{}"
	s := NRSink{}
	err := s.doCheckJson(kt.NewOutput([]byte(input)))
	assert.NoError(err)

	// bad no compression case
	input = "{}aaa"
	err = s.doCheckJson(kt.NewOutput([]byte(input)))
	assert.Error(err)

	// No compresson bad
	s.compression = kt.CompressionGzip
	input = "{}"
	err = s.doCheckJson(kt.NewOutput([]byte(input)))
	assert.Error(err)

	input = "[]"
	serBuf := []byte{}
	buf := bytes.NewBuffer(serBuf)
	buf.Reset()
	zw, _ := gzip.NewWriterLevel(buf, gzip.DefaultCompression)
	zw.Write([]byte(input))
	zw.Close()
	err = s.doCheckJson(kt.NewOutput(buf.Bytes()))
	assert.NoError(err)

	input = "aaaa[]"
	buf.Reset()
	zwa, _ := gzip.NewWriterLevel(buf, gzip.DefaultCompression)
	zwa.Write([]byte(input))
	zwa.Close()
	err = s.doCheckJson(kt.NewOutput(buf.Bytes()))
	assert.Error(err)
}

func TestDetectRegionFromLicenseKey(t *testing.T) {
	assert := assert.New(t)

	// JP-style prefix: "jpx" followed by arbitrary key material.
	assert.Equal(REGION_JP, detectRegionFromLicenseKey("jpxSOMEDUMMYKEYMATERIAL"))

	// EU-style prefix: "eu01x" followed by arbitrary key material.
	assert.Equal(REGION_EU, detectRegionFromLicenseKey("eu01xSOMEDUMMYKEYMATERIAL"))

	// US-style key: no prefix, so no "x" delimiter to match on.
	assert.Equal("", detectRegionFromLicenseKey("SOMEDUMMYKEYMATERIALNOPREFIX"))

	// Empty key.
	assert.Equal("", detectRegionFromLicenseKey(""))

	// Unrecognized region code should not error, just fall through to the default.
	assert.Equal("", detectRegionFromLicenseKey("zz01xSOMEDUMMYKEYMATERIAL"))
}
