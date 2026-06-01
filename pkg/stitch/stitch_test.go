package stitch

import (
	"strconv"
	"testing"

	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/stretchr/testify/assert"
)

func TestUnify(t *testing.T) {
	assert := assert.New(t)
	l := lt.NewTestContextL(logger.NilContext, t).GetLogger().GetUnderlyingLogger()

	s, err := NewStitcher(l, &ktranslate.StitchConfig{Enable: true, TTLSec: 30})

	for _, dst := range kt.InputTestingUnify {
		dst.CustomStr["src_endpoint"] = dst.SrcAddr + ":" + strconv.Itoa(int(dst.L4SrcPort))
		dst.CustomStr["dst_endpoint"] = dst.DstAddr + ":" + strconv.Itoa(int(dst.L4DstPort))
	}

	assert.False(s.Stitch(kt.InputTestingUnify[0]))
	assert.True(s.Stitch(kt.InputTestingUnify[1]))
	assert.False(s.Stitch(kt.InputTestingUnify[1]))
}
