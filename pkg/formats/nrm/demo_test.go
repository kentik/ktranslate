package nrm

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWave(t *testing.T) {
	assert := assert.New(t)
	w := NewSineWave(10, 10)
	assert.NotNil(w)

	for i := 0; i < 6; i++ {
		v := <-w.Output
		assert.Equal(float32(i), float32(v), "%d -> %f", i, v)
	}
	for i := 4; i > 0; i-- {
		v := <-w.Output
		assert.Equal(float32(i), float32(v), "%d -> %f", i, v)
	}
	for i := 0; i < 6; i++ { // Now we go negative.
		v := <-w.Output
		assert.Equal(float32(i*-1), float32(v), "%d -> %f", i, v)
	}
}

func TestWaveSmall(t *testing.T) {
	assert := assert.New(t)
	w := NewSineWave(100, 1)
	assert.NotNil(w)

	for i := 0; i < 50; i++ {
		v := <-w.Output
		assert.True(math.Abs(float64(i)*.01-float64(v)) < .001, "%d -> %f", i, v)
	}
}
