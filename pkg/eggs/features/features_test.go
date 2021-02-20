package features

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidFeatureName(t *testing.T) {
	assert.True(t, IsValidFeatureName("cool"))
	assert.True(t, IsValidFeatureName("cool.prop"))

	assert.False(t, IsValidFeatureName("cool.prop/42"))
	assert.False(t, IsValidFeatureName(""))
	assert.False(t, IsValidFeatureName("not cool"))
	assert.False(t, IsValidFeatureName("nope-nope"))
}
