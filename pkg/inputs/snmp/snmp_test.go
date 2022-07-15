package snmp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kentik/ktranslate/pkg/kt"
)

func TestMatchesPrefix(t *testing.T) {
	assert := assert.New(t)

	tests := map[string][]string{
		"tagName":               []string{"tagName", "true"},
		"provider:tagName":      []string{"provider:tagName", "true"},
		"provider:foo:tagName":  []string{"tagName", "true"},
		"provider:bar:tagName":  []string{"", "false"},
		"provider:foo:tagName1": []string{"tagName1", "true"},
	}

	provider := kt.Provider("foo")

	for in, expt := range tests {
		newTag, res := matchesPrefix(in, provider)
		assert.Equal(expt[0], newTag, "%s <-> %s", in, provider)
		assert.Equal(expt[1], fmt.Sprintf("%v", res), "%s <-> %s", in, provider)
	}
}
