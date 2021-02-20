package properties

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidPropertyName(t *testing.T) {
	assert.True(t, IsValidPropertyName("cool"))
	assert.True(t, IsValidPropertyName("cool.prop"))
	assert.True(t, IsValidPropertyName("yep-yep"))

	assert.False(t, IsValidPropertyName("cool.prop/42"))
	assert.False(t, IsValidPropertyName(""))
	assert.False(t, IsValidPropertyName("not cool"))
}
