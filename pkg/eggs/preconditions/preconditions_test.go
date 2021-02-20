package preconditions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssertNonNil(t *testing.T) {
	var nilPtr *int = nil
	assert.Panics(t, func() { AssertNonNil(nilPtr, "nilptr") })

	c := 42
	var ptr *int = &c
	assert.NotPanics(t, func() { AssertNonNil(ptr, "nilptr") })
}

func TestAssertNil(t *testing.T) {
	var nilPtr *int = nil
	assert.NotPanics(t, func() { AssertNil(nilPtr, "nilptr") })

	c := 42
	var ptr *int = &c
	assert.Panics(t, func() { AssertNil(ptr, "nilptr") })
}
