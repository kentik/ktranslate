package cat

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/kentik/ktranslate/pkg/kt"
)

// simulatedBatchCost stands in for handleInput's real per-batch work (enrichment,
// filtering, rollups, serialization -- kkc.go:511-553). Constructing a real
// *KTranslate fixture to call handleInput directly would require wiring up its
// metrics registry, format encoder, and several other fields not relevant to what
// this benchmark measures; instead this models the producer/consumer *shape* of
// kc.inputChan (kkc.go:777-786) using the real CHAN_SLACK constant, with a fake
// per-batch cost standing in for handleInput's actual work. Scaled down from a
// realistic per-batch processing time so the suite stays fast -- read results as
// relative comparisons between consumer counts, not absolute production timing.
const simulatedBatchCost = 100 * time.Microsecond

// BenchmarkConsumerThroughput models the inputChan producer/consumer pattern
// (kkc.go:777-786: kc.inputChan sized CHAN_SLACK, drained by kc.config.InputThreads
// consumer goroutines running monitorInput, kkc.go:776-787) under load from many
// concurrent producers -- one per SNMP device poller sending a batch per poll cycle
// (metrics/poll.go:226). consumers=1 is the shipped default (config.go:308,420:
// InputThreads defaults to 1); watchInput's autoscaler (kkc.go:556-576) only adds
// more if MaxThreads is also raised above its own default of 1 (config.go:309,421),
// so out of the box this is the only consumer count that actually applies.
func BenchmarkConsumerThroughput(b *testing.B) {
	for _, producers := range []int{100, 1000, 5000} {
		for _, consumers := range []int{1, 4, 16} {
			b.Run(fmt.Sprintf("producers=%d/consumers=%d", producers, consumers), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					ch := make(chan []*kt.JCHF, CHAN_SLACK)

					var consumerWG sync.WaitGroup
					consumerWG.Add(consumers)
					for range consumers {
						go func() {
							defer consumerWG.Done()
							for range ch {
								time.Sleep(simulatedBatchCost)
							}
						}()
					}

					var producerWG sync.WaitGroup
					producerWG.Add(producers)
					for range producers {
						go func() {
							defer producerWG.Done()
							ch <- []*kt.JCHF{}
						}()
					}

					producerWG.Wait()
					close(ch)
					consumerWG.Wait()
				}
			})
		}
	}
}
