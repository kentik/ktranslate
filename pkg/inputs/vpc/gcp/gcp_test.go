package gcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSampleRate(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(int(1), int(getSampleRate(2.0))) // Upper bound.
	assert.Equal(int(1), int(getSampleRate(1.0)))
	assert.Equal(int(2), int(getSampleRate(.5)))
	assert.Equal(int(10), int(getSampleRate(.1)))
	assert.Equal(int(100), int(getSampleRate(.01)))
	assert.Equal(int(1000), int(getSampleRate(.001)))
	assert.Equal(int(100000), int(getSampleRate(.0000000001))) // Lower bound.
}
