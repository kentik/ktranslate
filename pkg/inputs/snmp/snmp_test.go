package snmp

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kentik/ktranslate/pkg/kt"
)

func TestMatchesPrefix(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]bool{
		"tagName":               true,
		"provider:tagName":      true,
		"provider:foo:tagName":  true,
		"provider:bar:tagName":  false,
		"provider:foo:tagName1": true,
	}

	provider := kt.Provider("foo")

	for in, expt := range tests {
		res := matchesPrefix(in, provider)
		assert.Equal(expt, res, "%s <-> %s", in, provider)
	}
}
