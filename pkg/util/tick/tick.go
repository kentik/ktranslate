package tick

import (
	"time"
)

// An AlignedTicker exposes a channel on which a message is written once at
// a particular time in the future, and at a fixed interval thereafter.  Note
// that this is essentially a combination of features of the standard go
// time.Timer and time.Ticker types -- it's intended to be used as a
// (datarace-free) replacement for a pattern of using those two types in
// multiple Kentik services.

// (Quick note: this type basically replicates the behavior of go's internal
// runtime.Timer type -- but that type is exposed only through time.Ticker
// and time.Timer, each of which exposes only a subset of the functionality
// we need here.  So we're basically using one of each of the public types
// to reconstruct the hidden capabilities of the single internal type.
// Frustrating, but I don't see a better option.)

type AlignedTicker struct {
	C     <-chan time.Time
	c     chan time.Time
	first *time.Timer
	after *time.Ticker
	stop  chan struct{}
}

/** We've gotten rid of all of these now, apparently

func NewAlignedTicker(interval time.Duration) *AlignedTicker {
	now := time.Now()
	firstTime := now.Add(interval).Truncate(interval)
	firstInterval := firstTime.Sub(now)

	return newAlignedTicker(firstInterval, interval)
}
*/

//
func NewFixedTimer(start time.Time, interval time.Duration) *AlignedTicker {

	firstInterval := start.Sub(time.Now())
	return newAlignedTicker(firstInterval, interval)
}

func newAlignedTicker(firstInterval, afterInterval time.Duration) *AlignedTicker {
	c := make(chan time.Time, 1)

	at := &AlignedTicker{
		C:     c,
		c:     c,
		first: time.NewTimer(firstInterval),
		stop:  make(chan struct{}),
	}

	go at.run(afterInterval)

	return at
}

func (at *AlignedTicker) run(interval time.Duration) {

	// The following loop depends on a nil channel never blocking.
	// Exactly one of firstC and afterC should ever be non-nil, so
	// only one of the two loop cases can ever trigger at any given
	// time -- the firstC case before the first timer event is received,
	// and the afterC case thereafter.
	firstC := at.first.C
	var afterC <-chan time.Time

	for {
		select {
		case t := <-firstC:
			// the initial timer has fired.  It's time to create
			// the recurring ticker, stop the timer, and then nil
			// firstC to shut it up
			at.after = time.NewTicker(interval)
			afterC = at.after.C

			at.first.Stop()
			firstC = nil

			at.c <- t

		case t := <-afterC:
			at.c <- t

		case <-at.stop:
			if firstC != nil {
				at.first.Stop()
			}
			if afterC != nil {
				at.after.Stop()
			}
			return
		}
	}
}

// Stop the AlignedTicker from generating any more ticks on its channel.  Like
// time.Ticker.Stop() and time.Timer.Stop(), does not close the public channel,
// so that a read on the channel won't unblock.
func (at *AlignedTicker) Stop() {
	close(at.stop)
}
