package snmp

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kentik/ktranslate/pkg/api"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
	snmp_util "github.com/kentik/ktranslate/pkg/inputs/snmp/util"
	"github.com/kentik/ktranslate/pkg/kt"
)

// BenchmarkDeviceInitLoop models the synchronous per-device work in
// runSnmpPolling's loop (snmp.go:176-235), which runs serially for every device on
// every restart -- not just newly discovered ones. Any discovery delta triggers a
// SIGUSR2 (disco.go:213-215), which snmp_unix.go:44-46 turns into a full teardown and
// serial re-launch of the entire fleet, not just the delta. This benchmark measures
// how that per-device cost scales with fleet size (100/1,000/5,000).
//
// It calls the real snmp_util.InitSNMP twice per device, exactly as launchSnmp does
// (snmp.go:310,316) -- confirmed safe to call directly with no real target listening:
// gosnmp's Connect() only opens a local UDP socket (no handshake, no network round
// trip for UDP), so this exercises genuine production code at genuine cost, not a
// fake. Each connection is closed immediately after creation (GoSNMP.Close(),
// gosnmp.go:294) purely to avoid exhausting file descriptors across repeated
// benchmark iterations -- production code holds these open for the poller's
// lifetime instead.
//
// It also calls the real, unmodified apic.EnsureDevice (api.go:580-624), which
// short-circuits at api.go:585-588 because KentikPlan is never set anywhere in this
// fork's config (confirmed via grep -rn "KentikPlan" pkg/) -- so this is genuinely a
// no-op today, exercised as-is rather than stubbed.
//
// It deliberately does NOT call the real launchSnmp/metricPoller.Poll: that spawns a
// background goroutine which would attempt a real SNMP round trip against a
// nonexistent device and only resolve after a multi-second timeout -- leaking
// goroutines across iterations rather than measuring anything useful. The
// synchronous per-device loop body is what §2.5 in docs/DISCOVERY_PERFORMANCE_PLAN.md
// identifies as the bottleneck; the async poll itself is unaffected by fleet size.
func BenchmarkDeviceInitLoop(b *testing.B) {
	log := lt.NewBenchContextL(logger.NilContext, b)
	apic := &api.KentikApi{} // zero-value: EnsureDevice short-circuits, api.go:585-588
	ctx := context.Background()
	timeout := 3 * time.Second // matches shipped timeout_ms: 3000
	const retries = 0

	for _, fleetSize := range []int{100, 1000, 5000} {
		b.Run(fmt.Sprintf("devices=%d", fleetSize), func(b *testing.B) {
			devices := make([]*kt.SnmpDeviceConfig, fleetSize)
			for i := range devices {
				devices[i] = &kt.SnmpDeviceConfig{
					DeviceName: fmt.Sprintf("device-%d", i),
					DeviceIP:   fmt.Sprintf("10.%d.%d.%d", (i>>16)&0xff, (i>>8)&0xff, i&0xff),
					Community:  "public",
				}
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for _, device := range devices {
					metadataServer, err := snmp_util.InitSNMP(device, timeout, retries, "", log)
					if err != nil {
						b.Fatalf("InitSNMP (metadata): %v", err)
					}
					metricsServer, err := snmp_util.InitSNMP(device, timeout, retries, "", log)
					if err != nil {
						b.Fatalf("InitSNMP (metrics): %v", err)
					}
					if err := apic.EnsureDevice(ctx, device); err != nil {
						b.Fatalf("EnsureDevice: %v", err)
					}
					metadataServer.Close()
					metricsServer.Close()
				}
			}
		})
	}
}
