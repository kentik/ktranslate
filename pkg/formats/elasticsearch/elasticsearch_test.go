package elasticsearch

import (
	"testing"

	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/stretchr/testify/assert"
)

func TestSerializeElasticsearch(t *testing.T) {
	serBuf := make([]byte, 0)
	assert := assert.New(t)

	f, err := NewFormat(nil, kt.CompressionNone)
	assert.NoError(err)

	res, err := f.To(kt.InputTesting, serBuf)
	assert.NoError(err)
	assert.NotNil(res)

	out, err := f.From(res)
	assert.NoError(err)
	assert.Equal(len(kt.InputTesting), len(out))
	for i, _ := range out {
		assert.Equal(kt.InputTesting[i].SrcAddr, out[i]["src_addr"])
		assert.Equal(int(0), int(out[i]["timestamp"].(int64)))
	}
}

func TestSerializeElasticsearchGzip(t *testing.T) {
	serBuf := make([]byte, 0)
	assert := assert.New(t)
	f, err := NewFormat(nil, kt.CompressionGzip)
	assert.NoError(err)
	res, err := f.To(kt.InputTesting, serBuf)
	assert.NoError(err)
	assert.NotNil(res)
	out, err := f.From(res)
	assert.NoError(err)
	assert.Equal(len(kt.InputTesting), len(out))
	for i, _ := range out {
		assert.Equal(kt.InputTesting[i].SrcAddr, out[i]["src_addr"])
		assert.Equal(int(0), int(out[i]["timestamp"].(int64)))
	}
}
