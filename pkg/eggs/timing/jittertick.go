// A Ticker with jittering
//

package timing

import (
	"errors"
	"math/rand"
	"time"
)

// Just like a time.Ticker, a JitterTicker holds a channel that delivers
// `ticks' of a clock at intervals. The first tick is jittered and will fire
// after a random period (within a given range).

type JitterTicker struct {
	C      <-chan time.Time // The channel on which the ticks are delivered.
	jitter time.Duration    // for tests
	ticker *time.Ticker
	stop   chan struct{}
}

// Creates a jittered ticker
// windowStart and windowEnd specify the jitter range relative to the d interval
// For example, with d=1000ms start=80 and end=120, the first tick will happen
// after a random delay between 800ms and 1200ms. Subsequent ticks will happen
// every 1000ms.
func NewJitterTicker(d time.Duration, windowStart uint, windowEnd uint) *JitterTicker {
	if d <= 0 {
		panic(errors.New("non-positive interval for NewJitterTicker"))
	}
	if windowStart > windowEnd {
		panic(errors.New("windowStart is bigger than windowEnd for NewJitterTicker"))
	}

	dnanos := d.Nanoseconds()
	start := (dnanos * int64(windowStart)) / 100
	end := (dnanos * int64(windowEnd)) / 100

	var jitter time.Duration
	if end == start {
		jitter = time.Duration(start)
	} else {
		jitter = time.Duration(start+rand.Int63n(end-start)) * time.Nanosecond
	}
	return newJitterTicker(d, jitter)
}

func newJitterTicker(d time.Duration, jitter time.Duration) *JitterTicker {

	c := make(chan time.Time, 1)
	t := &JitterTicker{
		C:      c,
		jitter: jitter,
		stop:   make(chan struct{}),
	}

	go func() {
		// wait for out chosen jittering period and send the 1st tick
		time.Sleep(jitter)
		c <- time.Now()

		// ... then start the actual ticker and relay ticks
		t.ticker = time.NewTicker(d)
		for {
			select {
			case t := <-t.ticker.C:
				c <- t
			case <-t.stop:
				t.ticker.Stop()
				return
			}
		}
	}()

	return t
}

// Stop turns off a ticker. After Stop, no more ticks will be sent.
// Stop does not close the channel, to prevent a read from the channel
// succeeding incorrectly.
func (t *JitterTicker) Stop() {
	close(t.stop)
}
