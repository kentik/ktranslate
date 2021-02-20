package timing

/*  Keeping this around, but not running routinely -- I don't think there's a much better test here than
 *  visually scanning the results, and make sure they're close-enough (for whatever value of close-enough)
 *  to the expected times.  Obviously, any automated test of this feature is going to either be prone to
 *  spurious scheduling false-positives, or so lenient as to be almost meaningless.  (time.Ticker's automated
 *  tests allow the tick to arrive anytime within 20% of the interval...)

import (
	"fmt"
	"testing"
	"time"
)


func TestAlignedTicker(t *testing.T) {
	const Count = 10

	// create a channel and spawn off a producer for it, just to interrupt our select loop below
	// as often as possible
	noise := make(chan bool)
	go func() {
		for {
			noise <- false
		}
	}()

	at := NewAlignedTicker(time.Duration(1) * time.Second)
	for i := 0; i < Count; {
		select {
		case ti := <-at.C:
			i++
			fmt.Printf("ticked at %v\n", ti)
		case <-noise:
			// no-op; just consume the signal
		}
	}
}
*/
