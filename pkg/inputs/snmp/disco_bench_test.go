package snmp

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// simulatedProbeLatency stands in for a real network round trip -- both furious's
// TCP-dial-to-port-1 liveness probe (vendored scan-device.go:106) and a single
// doubleCheckHost community/version attempt (disco.go:244-390, which fails fast on
// its first error per metadata/device_metadata.go:298-302, so it costs one simulated
// probe per combination, not one per OID). Real timeout_ms is 3000
// (config/snmp-base.yaml:27, deployment/docker/snmp-base-nr.yaml:26); scaled down by
// ~1000x here so the suite stays fast. Read results as relative comparisons between
// concurrency strategies, not absolute production timing -- docs/BENCHMARKING_PLAN.md's
// Tier B (a real NixOS device farm) is what validates real-world magnitude.
const simulatedProbeLatency = 3 * time.Millisecond

func simulateProbe(up bool) {
	if !up {
		time.Sleep(simulatedProbeLatency) // dead host: pay the full simulated "timeout"
	}
	// live host: real code returns almost immediately on success -- nothing to simulate
}

// BenchmarkPreScanFanOut models furious's DeviceScanner.Scan (A52-A56): one goroutine
// per address with no concurrency cap at all.
//
// A first version of this benchmark modeled the probe as a bare time.Sleep and found
// "unbounded" *faster* than any bounded pool -- which would be a misleading result to
// leave standing. That happens because time.Sleep costs nothing in real OS resources,
// so Go's scheduler handles 65,536 concurrent sleepers for free; a real dial doesn't
// get that deal -- it contends for a finite number of file descriptors/ephemeral
// ports/conntrack entries regardless of how many goroutines the Go program spawns.
// osCeiling below makes that shared, finite resource explicit and applies it to
// every variant equally, so what's actually being compared is fair: given the same
// real ceiling, does spawning more goroutines than could ever usefully run at once
// buy you anything? The expected (and, once osCeiling is introduced, actual) answer
// is no throughput benefit, only extra memory/allocs for the goroutines parked
// waiting on osCeiling that didn't need to exist. This still cannot reproduce what an
// exhausted ceiling actually *does* in production (real dial() errors, retransmits,
// cascading slowness) -- that's what Tier B's NixOS farm is for.
const osCeiling = 1024

func BenchmarkPreScanFanOut(b *testing.B) {
	const addrs = 65536
	liveRatio := 0.076 // ~5,000/65,000 -- the field report's ratio; var, not const, so int() below truncates rather than requiring an exact constant
	live := int(addrs * liveRatio)
	isUp := func(i int) bool { return i < live }

	// runWithPool uses one uniform goroutine shape for every variant (the same
	// closure, with or without a semaphore acquire) so allocation counts are
	// directly comparable across variants -- an earlier version wrapped the pooled
	// case in an extra closure layer the unbounded case didn't have, which made
	// "pool" look like it allocated more even though it holds far fewer goroutines
	// alive at once. poolSize <= 0 means unbounded (no self-imposed limit beyond
	// the shared osCeiling every variant already pays into).
	runWithPool := func(b *testing.B, poolSize int) {
		for i := 0; i < b.N; i++ {
			ceiling := make(chan struct{}, osCeiling)
			var sem chan struct{}
			if poolSize > 0 {
				sem = make(chan struct{}, poolSize)
			}
			var wg sync.WaitGroup
			wg.Add(addrs)
			for a := range addrs {
				up := isUp(a)
				if sem != nil {
					sem <- struct{}{}
				}
				go func() {
					defer wg.Done()
					if sem != nil {
						defer func() { <-sem }()
					}
					ceiling <- struct{}{}
					defer func() { <-ceiling }()
					simulateProbe(up)
				}()
			}
			wg.Wait()
		}
	}

	b.Run("unbounded", func(b *testing.B) { runWithPool(b, 0) })
	for _, pool := range []int{256, osCeiling} {
		b.Run(fmt.Sprintf("pool=%d", pool), func(b *testing.B) {
			runWithPool(b, pool)
		})
	}
}

// BenchmarkVerificationLoop models doubleCheckHost's semaphore-gated worker pool
// (disco.go:82-85, 176-181): conf.Disco.Threads workers processing candidate
// addresses, each costing one simulated probe. threads=4 is the shipped default
// (config/snmp-base.yaml:18, deployment/docker/snmp-base-nr.yaml:18); the other
// values are Phase 1's proposed higher floor.
func BenchmarkVerificationLoop(b *testing.B) {
	const addrs = 2000
	liveRatio := 0.076 // var, not const -- see BenchmarkPreScanFanOut
	live := int(addrs * liveRatio)
	isUp := func(i int) bool { return i < live }

	for _, threads := range []int{4, 64, 256} {
		b.Run(fmt.Sprintf("threads=%d", threads), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ctl := make(chan bool, threads)
				for range threads {
					ctl <- true
				}
				var wg sync.WaitGroup
				wg.Add(addrs)
				for a := range addrs {
					go func(up bool) {
						<-ctl
						defer func() { wg.Done(); ctl <- true }()
						simulateProbe(up)
					}(isUp(a))
				}
				wg.Wait()
			}
		})
	}
}

// BenchmarkCIDRSerialization models runScanCheckDisco's serial `for _, ipr := range
// conf.Disco.Cidrs` (disco.go:141), where each CIDR's scan-and-verify cycle
// (disco.go:167-181) fully completes -- including its own wg.Wait() at line 181 --
// before the next CIDR starts, vs. running independent CIDRs concurrently (Phase 2).
// Each simulated CIDR reuses the same threads=4-style bounded pool internally, since
// Phase 2 is about parallelizing *across* CIDRs, not changing per-CIDR concurrency.
func BenchmarkCIDRSerialization(b *testing.B) {
	const cidrs = 4
	const addrsPerCIDR = 500
	const threads = 4
	liveRatio := 0.076 // var, not const -- see BenchmarkPreScanFanOut
	live := int(addrsPerCIDR * liveRatio)

	scanOneCIDR := func() {
		ctl := make(chan bool, threads)
		for range threads {
			ctl <- true
		}
		var wg sync.WaitGroup
		wg.Add(addrsPerCIDR)
		for a := range addrsPerCIDR {
			up := a < live
			go func(up bool) {
				<-ctl
				defer func() { wg.Done(); ctl <- true }()
				simulateProbe(up)
			}(up)
		}
		wg.Wait()
	}

	b.Run("serial", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for range cidrs {
				scanOneCIDR()
			}
		}
	})

	b.Run("parallel", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var wg sync.WaitGroup
			wg.Add(cidrs)
			for range cidrs {
				go func() {
					defer wg.Done()
					scanOneCIDR()
				}()
			}
			wg.Wait()
		}
	})
}
