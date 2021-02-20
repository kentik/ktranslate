package tick

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBasicJitterTicker(t *testing.T) {
	assert := assert.New(t)

	jticker := NewJitterTicker(100*time.Millisecond, 0, 100)
	counter := uint64(0)
	go func() {
		for range jticker.C {
			atomic.AddUint64(&counter, 1)
		}
	}()

	time.Sleep(300 * time.Millisecond)
	assert.Equal(uint64(3), atomic.LoadUint64(&counter))
}

func TestJitterTickerRange(t *testing.T) {
	assert := assert.New(t)

	for i := 1; i <= 1000; i++ {
		jticker := NewJitterTicker(1000, 50, 75)
		assert.True((jticker.jitter >= 500) && (jticker.jitter <= 750))
	}
}

func TestJitterTickerEmptyRange(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(time.Duration(0), NewJitterTicker(1000, 0, 0).jitter)
	assert.Equal(time.Duration(1500), NewJitterTicker(1000, 150, 150).jitter)
}

func BenchmarkJitterTicker(b *testing.B) {
	jticker := NewJitterTicker(1, 0, 0)

	counter := 0
	for range jticker.C {
		counter++
		if counter >= b.N {
			break
		}
	}
	jticker.Stop()
}
