package rule

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsInternal(t *testing.T) {
	sut := NewRuleSet(nil)

	tests := map[string]bool{
		"1.2.3.4":            false,
		"10.2.1.2":           true,
		"2001:d00:abcd::/64": false,
		"2001:d01:abcd::/64": true,
		"10.2.1.3":           true,
		"172.16.0.1":         true,
		"10.19.38.78":        true,
		"foo":                false,
	}

	asns := map[string]uint32{
		"1.2.3.4":            3,
		"10.2.1.2":           3,
		"2001:d00:abcd::/64": 3,
		"2001:d01:abcd::/64": 64512,
		"10.2.1.3":           4294967294,
		"172.16.0.1":         3,
		"foo":                3,
		"10.19.38.78":        23,
	}

	for ip, isInternal := range tests {
		assert.Equal(t, isInternal, sut.IsInternal(net.ParseIP(ip), asns[ip]), ip)
	}
}
