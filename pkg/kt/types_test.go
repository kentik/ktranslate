package kt

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnifyFlatten(t *testing.T) {
	assert := assert.New(t)

	for _, dst := range InputTestingUnify {
		dst.CustomStr["src_endpoint"] = dst.SrcAddr + ":" + strconv.Itoa(int(dst.L4SrcPort))
		dst.CustomStr["dst_endpoint"] = dst.DstAddr + ":" + strconv.Itoa(int(dst.L4DstPort))
	}

	InputTestingUnify[0].Pair = InputTestingUnify[1]
	res := InputTestingUnify[0].Flatten()

	assert.Equal(InputTestingUnify[0].CustomStr["src_endpoint"], res["src_endpoint"])
	assert.Equal(InputTestingUnify[0].Pair.CustomStr["src_endpoint"], res["pair_src_endpoint"], "map: %v", res)
	assert.Equal(int64(InputTestingUnify[0].Pair.L4DstPort), res["pair_l4_dst_port"], "map: %v", res)
	assert.Equal(int64(InputTestingUnify[0].L4DstPort), res["l4_dst_port"], "map: %v", res)
}
