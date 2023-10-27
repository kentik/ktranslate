package kt

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt2ip(t *testing.T) {
	inputs := []uint64{1057383603, 183819855, 1057412798, 183819855, 183818207, 183819882}
	outputs := []string{"63.6.100.179", "10.244.222.79", "63.6.214.190", "10.244.222.79", "10.244.215.223", "10.244.222.106"}

	for i, input := range inputs {
		assert.Equal(t, outputs[i], Int2ip(uint32(input)).String(), "%d", i)
	}

	// Test again as a byte[]
	for i, input := range inputs {
		var bytes [4]byte
		bytes[0] = byte(input & 0xFF)
		bytes[1] = byte((input >> 8) & 0xFF)
		bytes[2] = byte((input >> 16) & 0xFF)
		bytes[3] = byte((input >> 24) & 0xFF)

		assert.Equal(t, outputs[i], net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0]).String(), "%d", i)
	}
}
