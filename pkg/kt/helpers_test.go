package kt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt2ip(t *testing.T) {
	inputs := []uint64{1057383603, 183819855, 1057412798, 183819855, 183818207, 183819882}
	outputs := []string{"63.6.100.179", "10.244.222.79", "63.6.214.190", "10.244.222.79", "10.244.215.223", "10.244.222.106"}

	for i, input := range inputs {
		assert.Equal(t, outputs[i], Int2ip(uint32(input)).String(), "%d", i)
	}
}
