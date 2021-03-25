package netflow

import (
	"testing"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	assert := assert.New(t)
	serBuf := make([]byte, 0)
	l := lt.NewTestContextL(logger.NilContext, t).GetLogger().GetUnderlyingLogger()
	_, err := NewFormat(l, kt.CompressionSnappy)
	assert.Error(err)

	vv := "netflow5"
	Version = &vv
	_, err = NewFormat(l, kt.CompressionNone)
	assert.Error(err)

	vv = "ipfix"
	f, err := NewFormat(l, kt.CompressionNone)
	assert.NoError(err)

	res, err := f.To(kt.InputTesting, serBuf)
	assert.NoError(err)
	assert.NotNil(res)

	out, err := f.From(res)
	assert.NoError(err)
	assert.Equal(len(kt.InputTesting), len(out))
	for i, _ := range out {
		assert.Equal(kt.InputTesting[i].Protocol, out[i]["protocol"].(string))
		assert.Equal(kt.InputTesting[i].SrcAddr, out[i]["src_addr"])
		assert.Equal(kt.InputTesting[i].DstAddr, out[i]["dst_addr"])
		assert.Equal(int(kt.InputTesting[i].L4DstPort), int(out[i]["l4_dst_port"].(int64)))
		assert.Equal(int(kt.InputTesting[i].OutputPort), int(out[i]["output_port"].(int64)))
		assert.Equal(int(kt.InputTesting[i].InBytes), int(out[i]["in_bytes"].(int64)))
		assert.Equal(kt.InputTesting[i].SrcEthMac, out[i]["src_eth_mac"])
		assert.Equal(int(kt.InputTesting[i].OutBytes), int(out[i]["out_bytes"].(int64)))
	}
}
