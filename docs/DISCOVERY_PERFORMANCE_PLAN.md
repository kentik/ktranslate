# SNMP Discovery & Polling Performance — Bottleneck Analysis and Remediation Plan

Branch: `investigation`. Written in response to the field report: *"IP Range - 65,000 ips
in 2 hrs, may be only 5000 devices."* Roadmap sub-tasks this maps to: **fix multi-threading**,
**parallel discovery**, **concurrent polling**, **adaptive back pressure**.

The linked roadmap spreadsheet
(`docs.google.com/spreadsheets/d/1pzp7vPc6hbCaRldT1jZ_WDxjnucG95BkjWSrl5W0PVA`) requires
Google sign-in and could not be fetched by the agent — none of its content is reflected
here. If it contains constraints (deadlines, already-decided approaches, prior discussion)
that should shape this plan, paste the relevant rows in and the plan will be revised.

All line numbers below were read directly from the working tree on branch `investigation`
at the time of writing. Re-verify with `grep -n` before editing if any surrounding file has
changed.

---

## 1. Executive summary

Discovery and steady-state polling share three architectural problems:

1. **Discovery concurrency is capped at 4 workers by config**, and the shipped
   `check_all_ips: true` default forces *every* address in the range through that 4-worker
   pool with a full SNMP timeout per address, instead of only the addresses that a pre-scan
   found alive — while still paying for the pre-scan.
2. **The pre-scan itself is unbounded** (one goroutine per IP, no cap) and uses a
   TCP-connect-to-port-1 trick with the *SNMP* timeout as its liveness probe, plus an inline
   blocking reverse-DNS lookup per ARP hit — so under `check_all_ips: true` it does a lot of
   expensive, discarded work.
3. **Steady-state polling is parallel per device already**, but (a) any newly discovered
   device restarts the *entire* fleet's pollers and serially re-initializes all of them, and
   (b) all devices' output funnels into a single consumer goroutine by default, with an
   autoscaler that is inert unless explicitly configured otherwise.

(1)+(2) plausibly account for the reported "65,000 IPs / 2 hours" number on their own
(back-of-envelope in §2.2). (3) is what will hurt specifically when scaling *up* to ~5,000
devices, independent of the initial scan.

---

## 2. Verified bottlenecks

### 2.1 Pre-scan ("is this host up?") is unbounded and reuses the SNMP timeout for a TCP dial

File: `pkg/inputs/snmp/disco.go`
- `runScanCheckDisco`, lines 139–186. Line 157: `scan.NewTargetIterator(ipr)`. Line 158:
  `timeout := time.Millisecond * time.Duration(conf.Global.TimeoutMS)` — the *SNMP* timeout
  is reused as the network liveness-probe timeout. Line 163:
  `results, err := scanner.Scan(ctx, conf.Disco.Ports)` — this call blocks until the entire
  CIDR has been probed.

The `Scan` implementation is in the vendored dependency
`github.com/liamg/furious@v0.0.0-20191231090757-c295c872d6c1`, resolved locally at
`~/.go/pkg/mod/github.com/liamg/furious@v0.0.0-20191231090757-c295c872d6c1/scan/scan-device.go`
(this file ships inside the built binary; it is not New Relic/Kentik-owned code, but its
behavior is fully in scope for a fix — either by wrapping it, replacing it, or vendoring a
patched copy):

- Lines 56–125: for every IP produced by the target iterator, a **new goroutine is spawned
  with no concurrency limit whatsoever** — a `/16` scan launches 65,536 goroutines
  essentially at once.
- Line 78: `arp.Search(ip.String())` — reads/parses the OS ARP table once per IP, per
  goroutine (65k times per scan).
- Line 98: `net.LookupAddr(ip.String())` — a **synchronous, blocking reverse-DNS lookup**,
  performed inline whenever there's an ARP hit, with no timeout of its own (subject to
  system resolver defaults, which can be several seconds per query on retry).
- Line 106: `net.DialTimeout("tcp", fmt.Sprintf("%s:1", ip.String()), s.timeout)` — "is this
  host up" is decided by attempting a TCP connection to **port 1** and waiting up to
  `s.timeout` (== SNMP `timeout_ms`, 3000ms in the shipped configs — see §2.2) for a
  response. Hosts that are silently dropped by a firewall/ACL (the common case on
  enterprise networks, vs. an explicit RST) wait out the **full** timeout.
- Line 127: `wg.Wait()` — the function does not return until literally every one of the
  65,536 dials/lookups has completed or timed out.

At the scale of a real fan-out this large, the "unbounded" goroutines don't actually run
concurrently — the OS network stack (file descriptor limits, ephemeral port range,
conntrack table size) throttles them, and what you get in practice is partial, uncontrolled
serialization plus retransmits, not a clean N-way parallel scan.

### 2.2 `check_all_ips: true` discards the pre-scan and forces every address through a 4-worker SNMP-credential probe

Config field: `pkg/kt/snmp.go:261` — `CheckAll bool \`yaml:"check_all_ips,omitempty"\``.
Shipped **on** in `deployment/docker/snmp-base-nr.yaml:20` (`check_all_ips: true`).

Consumed in `pkg/inputs/snmp/disco.go:172`:
```go
if strings.HasSuffix(ipr, "/32") || result.IsHostUp() || conf.Disco.CheckAll {
```
When `CheckAll` is true, **every IP in the range** — not just the ones the §2.1 pre-scan
found alive — is queued into `doubleCheckHost` (line 178: `go doubleCheckHost(...)`).
Concurrency there is capped by a semaphore:

- Lines 58–61: `conf.Disco.Threads` defaults to 4 if unset (`log.Warnf("...defaulting
  threads to %d from 0.", ...)`).
- Lines 82–85: `ctl := make(chan bool, conf.Disco.Threads)` filled with `Threads` tokens —
  this is the *entire* concurrency budget for verifying every address in the range.
- Shipped config sets it explicitly: `config/snmp-base.yaml:18` and
  `deployment/docker/snmp-base-nr.yaml:18` both have `threads: 4`.

Meanwhile the pre-scan from §2.1 still runs over all 65k addresses first — and with
`CheckAll` true, its up/down result is never used for filtering, only `result.Manufacturer`
/`result.Name` are used for a log line and an initial (soon-overwritten) device name
(disco.go:257, 273). **The entire expensive pre-scan phase is wasted work whenever
`check_all_ips` is on.**

**Back-of-envelope cost** using the shipped defaults
(`config/snmp-base.yaml:18,27,28` / `deployment/docker/snmp-base-nr.yaml:18,26,27`:
`threads: 4`, `timeout_ms: 3000`, `retries: 0`, one community `public`):

`doubleCheckHost` (disco.go:244–390) for a non-SNMP host tries exactly one
version×community combination (1 community configured, `retries: 0`) and gives up after one
timeout — confirmed fail-fast in `GetBasicDeviceMetadata`
(`pkg/inputs/snmp/metadata/device_metadata.go:298-302`, `return nil, err` on the *first*
`WalkOID` error, it does not walk every metadata OID before giving up). So each dead address
costs ≈1 × `timeout_ms` = 3s.

```
(65,000 addresses − 5,000 real devices) × 3s ÷ 4 concurrent workers
  = 60,000 × 3 ÷ 4
  ≈ 45,000s ≈ 12.5 hours (upper bound; less if some hosts fail fast via ICMP-unreachable
    instead of a silent timeout — a mix that plausibly lands around the reported 2 hours
    depending on network ACL behavior)
```

This single combination — `check_all_ips: true` + `threads: 4` + a wasted pre-scan — is the
highest-confidence explanation for the reported symptom and should be fixed first.

### 2.3 `doubleCheckHost` exhausts every SNMP version × community combination serially, per host

File: `pkg/inputs/snmp/disco.go:244–390`.
- Line 248: `_ = <-ctl` — blocks until a worker slot is free.
- Lines 261–295: if v3 is configured, loops over every `V3SNMPConfig` in
  `conf.Disco.OtherV3s` + `conf.Disco.DefaultV3`, **serially**, each doing a full
  `InitSNMP` + `GetBasicDeviceMetadata` round trip before trying the next.
- Lines 298–329 (labelled loop `outer:`): if v3 didn't resolve anything, loops over
  `versions := []bool{conf.Disco.UseV1}` (or `{false, true}` if `CheckAllVersions` is set,
  line 300) **×** every string in `conf.Disco.DefaultCommunities`, again fully serially per
  combination.

With more than one community/v3 config configured (a very normal real-world setup — e.g.
`public`, a legacy community, plus a v3 profile), the §2.2 per-host cost multiplies directly
by the number of combinations, still gated through the same 4-worker pool.

### 2.4 CIDR ranges are scanned one at a time, not in parallel

File: `pkg/inputs/snmp/disco.go:139–186`, `runScanCheckDisco`.
- Line 141: `for _, ipr := range conf.Disco.Cidrs {` — a plain, serial loop.
- Line 181: `wg.Wait()` inside the loop body — each CIDR's scan-and-verify cycle must fully
  drain before the next CIDR entry starts.

If the target range is expressed as several smaller CIDR blocks (common when discovery is
scoped per site/VLAN) rather than one large block, there is currently **zero** parallelism
across them, compounding §2.1–§2.3 linearly per block.

### 2.5 Any newly discovered device restarts the entire fleet's pollers, then serially re-initializes every device

This is the mechanism most directly relevant to *"supporting 5k devices"* — its cost scales
with total fleet size on every discovery tick, not with the size of the delta.

- `pkg/inputs/snmp/disco.go:206–242`, `RunDiscoOnTimer`: runs `Discover` on a timer
  (`pt := time.Duration(pollTimeMin) * time.Minute`, line 207). Lines 213–215:
  ```go
  if stats.delta != 0 || stats.added > 0 { // Only restart if there's a different configuration.
      c <- kt.SIGUSR2 // Restart the main loop with a new config.
  ```
  So finding **one** new/changed device anywhere in the range signals a restart.
- `pkg/inputs/snmp/snmp_unix.go:20–47`, `wrapSnmpPolling` (the unix build; `snmp_windows.go`
  has the equivalent):
  - Line 22: `ctxSnmp, cancel := context.WithCancel(ctx)` — one context shared by every
    device poller launched this generation.
  - Line 23: `runSnmpPolling(ctxSnmp, ...)` — launches all pollers under that context.
  - Line 41: `_ = <-c` — blocks waiting for the `SIGUSR2` signal.
  - **Line 44: `cancel()`** — cancels the context for **every currently running device's**
    metadata/metric poller goroutines (they `select` on `ctx.Done()` inside
    `pkg/inputs/snmp/metrics/poll.go:231` and the equivalent in the metadata poller),
    stopping metric collection fleet-wide.
  - Line 46: `go wrapSnmpPolling(ctx, ..., restartCount+1, ...)` — recurses, which calls
    `runSnmpPolling` again.
- `pkg/inputs/snmp/snmp.go:159–238`, `runSnmpPolling`:
  - Line 161: `initSnmp(...)` → re-parses the config file from scratch (see §2.6).
  - **Line 176: `for _, device := range conf.Devices {`** — iterates **all** devices in the
    merged config, not just the newly discovered one(s). For each device, serially:
    - Lines 189–200: duplicate-name bookkeeping under `metrics.Mux.Lock()` (cheap, not the
      bottleneck).
    - Line 205: `mibdb.FindProfile(...)` — see §2.6 for its cost.
    - Line 226: `apic.EnsureDevice(ctx, device)` — a network call if Kentik device-creation
      is active (see §2.9 for why it's currently a no-op in this fork).
    - Line 231: `launchSnmp(...)` — this itself does two synchronous `InitSNMP` calls
      (lines 310, 316 in `launchSnmp`, `pkg/inputs/snmp/snmp.go:292–348`) before handing off
      to a background goroutine (line 328) for the actual polling loop. `InitSNMP` sets up
      a `gosnmp.GoSNMP` struct — it is local setup, not a network round trip, so this part
      is cheap per device, but it is still one more serial step repeated for every device on
      every restart.

At 5,000 devices: discovering **one** new device anywhere in a large range stops metric
collection for the whole fleet and walks all 5,000 devices serially (profile lookup +
EnsureDevice-if-active + SNMP client setup) before pollers resume. With
`DiscoveryIntervalMinutes` set to re-run discovery periodically, this repeats on every tick
that finds any change.

### 2.6 Every config parse re-derives MIB-profile matches for every device — with regex recompiled on every lookup, and this happens multiple times per discovery cycle

`parseConfig` (`pkg/inputs/snmp/snmp.go:350–482`) is invoked from at least three places in a
single discovery-and-reload cycle:
1. `Discover` → `pkg/inputs/snmp/disco.go:44`.
2. `addDevices` re-reads the file "to get any new changes" →
   `pkg/inputs/snmp/disco.go:411`.
3. The subsequent restart's `initSnmp` → `pkg/inputs/snmp/snmp.go:161` (via `runSnmpPolling`,
   itself triggered by the `SIGUSR2` from §2.5).

Inside `parseConfig`, lines 457–468:
```go
if mibdb != nil {
    for _, device := range ms.Devices {
        profile := mibdb.FindProfile(device.OID, device.Description, device.MibProfile)
        ...
```
— an O(devices) pass every time, ×3 calls per discovery cycle above ⇒ ~15,000 `FindProfile`
calls for a 5,000-device fleet per cycle (plus the one in `runSnmpPolling` itself,
`snmp.go:205`, and the one per newly-found device in `doubleCheckHost`,
`disco.go:350`).

`FindProfile` (`pkg/inputs/snmp/mibs/profile.go:270–300`) can fall through to `checkMatch`
(lines 302–337), which **compiles a regex from scratch on every single call**:
- Line 305: `r, err := regexp.Compile(match.Regex)` (ordered `MatchesList`).
- Line 322: `r, err := regexp.Compile(m)` (legacy `Matches` map).

None of this is cached. It's not the dominant cost next to §2.1–§2.5, but it's pure waste
that scales linearly with device count and repeats 3+ times per discovery cycle — worth
fixing in the same pass since §3's restart-storm fix touches the same code paths.

### 2.7 Per-device polling is already parallel — but funnels into a single consumer goroutine by default

Good news first: once a device's poller is launched, `pkg/inputs/snmp/snmp.go:328–345`
starts genuinely independent goroutines per device
(`metadataPoller.StartLoop(ctx)` / `metricPoller.StartLoop(ctx)`,
`pkg/inputs/snmp/metrics/poll.go:163–244` for the metric side). Nothing in the steady-state
poll loop itself serializes across devices. This part does **not** need "concurrent
polling" work in the sense of adding parallelism — it already has it.

The bottleneck is downstream, in `pkg/cat/kkc.go`:
- Line 779: `kc.inputChan = make(chan []*kt.JCHF, CHAN_SLACK)` (`CHAN_SLACK = 8000`, line
  44) — one shared channel that **every** SNMP device's poller writes into
  (`p.jchfChan <- flows`, `pkg/inputs/snmp/metrics/poll.go:226`, and the status-flow send at
  line 229).
- Line 798: `assureInput()` called before starting SNMP polling (line 800:
  `snmp.StartSNMPPolls(ctx, kc.inputChan, ...)`), and the closure at line 777:
  ```go
  assureInput := func() { // Start up input processing if any is asked of us.
      if kc.inputChan == nil {
          kc.inputChan = make(chan []*kt.JCHF, CHAN_SLACK)
          for i := 0; i < kc.config.InputThreads; i++ {
              go kc.monitorInput(ctx, i, kc.format.To)
          }
  ```
  `kc.config.InputThreads` (`config.go:308`) **defaults to `1`** (`config.go:420`). Unless a
  deployment explicitly passes `-input_threads N` (wired at `cmd/ktranslate/main.go:289`),
  **exactly one goroutine** (`monitorInput`, `pkg/cat/kkc.go:776–787`) drains the channel
  that every one of the 5,000 devices' pollers writes into, calling `kc.handleInput`
  (line 511) synchronously per batch — which does enrichment, filtering, rollups, and
  serialization (`pkg/cat/kkc.go:511–553`) before the result can be handed to a sink.

### 2.8 Backpressure is a fixed-size shared channel plus a coarse, largely-inert autoscaler

- `pkg/cat/kkc.go:565–576`, `watchInput`: every 60 seconds
  (`checkTicker := time.NewTicker(60 * time.Second)`, line 557) it checks
  `len(kc.inputChan) > CHAN_SLACK-10` (line 566, i.e. the channel is >99.9% full) and, only
  if `kc.config.InputThreads < kc.config.MaxThreads` (line 565), launches one more
  `monitorInput` consumer and increments `InputThreads` (lines 568–569).
- `MaxThreads` (`config.go:309`) **also defaults to `1`** (`config.go:421`), set via
  `-max_threads` (`cmd/ktranslate/main.go:300`). With both defaults at 1, the condition at
  line 565 (`1 < 1`) is always false — **this autoscaler is inert unless an operator
  explicitly overrides both flags.**
- Even when active: it only ever *adds* consumers (never removes them as load drops), it
  checks on a fixed 60s cadence regardless of how fast the queue is actually filling, and it
  has no visibility into *which* devices are producing the backlog — so there is no way to
  prioritize or shed load selectively. This is a shared, blunt instrument, not the
  "adaptive back pressure" called for in the roadmap.
- Producer side: `p.jchfChan <- flows` (`metrics/poll.go:226,229`) is an unbuffered-relative-
  to-load send into a channel of fixed capacity 8000 — once full, **every** device's poller
  goroutine blocks equally on send, with no way to signal a specific noisy/slow device to
  back off rather than stalling the whole fleet's collection cadence.

### 2.9 Minor findings

- **Only the first configured port is used for verification.** `pkg/inputs/snmp/disco.go:279`
  (v3 path) and line 312 (v2c path): `Port: uint16(conf.Disco.Ports[0])`. If an operator
  configures multiple SNMP ports in `discovery.ports`, only `Ports[0]` is ever tried for the
  actual credential check — not a performance bottleneck per se, but worth fixing alongside
  the rest of this file.
- **`EnsureDevice` is O(n) per call and currently a dead code path in this fork.**
  `pkg/api/api.go:580–624`: lines 592–599 do a nested linear scan over `api.devices` (a
  `map[string][]Device`) for every call, making device-registration startup O(n²) across a
  restart if it were active. However, lines 585–588 short-circuit
  (`if api.config.KentikPlan == 0 { return nil }`), and `grep -rn "KentikPlan"` across `pkg/`
  shows it is **only ever read**, never set, anywhere in this fork's config
  (`pkg/config/nr/nr.go`, `pkg/config/config.go`) — so `EnsureDevice` is a no-op today. Flag
  this so it isn't silently reintroduced as a bottleneck if Kentik-cloud device creation is
  ever wired back up for this fork; not worth spending remediation effort on now.

---

## 3. Root cause → roadmap sub-task mapping

| Roadmap sub-task | Relevant findings |
|---|---|
| Fix multi-threading | §2.1 (unbounded pre-scan goroutines), §2.2/§2.3 (`threads: 4` hard cap on verification) |
| Parallel discovery | §2.4 (CIDRs scanned serially), §2.2 (pre-scan wasted when `check_all_ips` is on) |
| Concurrent polling | §2.5 (full-fleet restart on any new device), §2.7 (single default consumer thread) |
| Adaptive back pressure | §2.8 (fixed channel + inert/blunt autoscaler) |
| (bonus, same files) | §2.6 (regex recompilation), §2.9 (port[0]-only, dead EnsureDevice path) |

---

## 4. Remediation plan

Ordered by expected impact-to-effort ratio. Each phase is independently shippable and
testable; none strictly depends on a later phase.

### Phase 1 — Fix the discovery verification worker pool and stop wasting the pre-scan

**Goal:** eliminate the §2.2/§2.3 12-hour-scale cost; make thread count self-tuning instead
of a config knob defaulting to 4.

- `pkg/inputs/snmp/disco.go:58–61`: raise the effective default (e.g. scale with
  `runtime.NumCPU()` or a much higher fixed floor like 64–256 for a network-bound workload —
  this is I/O-bound waiting on timeouts, not CPU-bound, so the safe ceiling is much higher
  than 4) when `Threads == 0`, and document that `threads: 4` in the shipped
  `config/snmp-base.yaml:18` / `deployment/docker/snmp-base-nr.yaml:18` should be updated to
  a realistic value (or removed so the new default applies).
- `pkg/inputs/snmp/disco.go:139–186` (`runScanCheckDisco`): when `conf.Disco.CheckAll` is
  true, skip the `scan.NewDeviceScanner`/`scanner.Scan` call (lines 157–166) entirely and
  build the candidate list directly from the target iterator, since its result is discarded
  anyway (§2.2). This alone removes the wasted pre-scan phase for the exact configuration
  the field report was hit with.
- When `CheckAll` is false (pre-scan filtering is actually used): replace the vendored
  `furious` `DeviceScanner.Scan` unbounded fan-out (`scan-device.go:56-125`) with a
  bounded-concurrency version — either vendor a patched copy in-repo (this codebase already
  vendors/wraps other third-party pieces, e.g. `pkg/inputs/snmp/util`) or write a
  replacement pre-scan using a worker pool (`golang.org/x/sync/errgroup` with
  `SetLimit`, or a simple semaphore channel like the one already used for
  `doubleCheckHost`, `disco.go:82-85`). Cap concurrency to something bounded by `ulimit -n`
  headroom (e.g. a few thousand, not 65,536), and give the reverse-DNS lookup
  (`scan-device.go:98`) its own short timeout via `context.WithTimeout` + a
  cancellable resolver call, rather than letting it block unbounded.
- **Bonus side effect of vendoring/replacing this dependency:** `disco.go:159` only ever
  calls `scan.NewDeviceScanner`/`scan.NewTargetIterator` from the vendored `furious/scan`
  package — never `scan.NewSynScanner`. But that package's `scan-syn.go:14` imports
  `github.com/google/gopacket/pcap` (a cgo binding requiring `libpcap` at link time) purely
  to support the SYN-scanner variant we don't use. Go links cgo at package granularity, so
  merely importing `furious/scan` for the scanner we do use drags in `libpcap` for the one we
  don't. Confirmed via `go list -deps -json ./cmd/ktranslate` filtered for `.CgoFiles`: of the
  three cgo-requiring packages in the whole build graph, `gopacket/pcap` is the only one that
  both needs a real external system library and has no pure-Go fallback (`DataDog/zstd`
  statically vendors its own C source by default; `prometheus/client_golang`'s cgo file is
  Darwin-only with an automatic `!cgo` fallback). Trimming `scan-syn.go` out when
  vendoring/replacing `DeviceScanner` removes `libpcap` from the build entirely — on Linux,
  with no cgo left in the graph, this should make the binary fully static with no extra flags
  (today, `CGO_ENABLED=1` by default and `libpcap` is dynamically linked, which is also why
  `.github/workflows/test-on-pr.yml` needs `apt-get install libpcap-dev` at all). Worth
  confirming with `ldd`/`readelf -d` on the resulting Linux binary once this lands.
- Verification: with a lab range configured with `check_all_ips: true`, measure wall-clock
  time for a 65k-address scan before/after. Target: no worse than
  `(dead_addresses × timeout_ms) ÷ new_thread_count`, and the pre-scan phase should not run
  at all in the `check_all_ips: true` case (confirm via a log line count / timing span
  around `disco.go:156-182`).

### Phase 2 — Parallelize across CIDR blocks

- `pkg/inputs/snmp/disco.go:141` (`for _, ipr := range conf.Disco.Cidrs`): convert to a
  bounded-concurrency fan-out (e.g. `errgroup` with a limit, or a semaphore sized
  independently from the per-host `Threads` semaphore at lines 82-85 — the two pools
  compose, so size them so their product stays within the same total-concurrency budget
  established in Phase 1, rather than multiplying unboundedly).
- `foundDevices` (line 107) and its guarding `mux` (line 168) are already safe for
  concurrent writers from multiple `doubleCheckHost` goroutines; extending that same map +
  mutex to be shared across concurrently-scanned CIDRs requires no structural change, just
  moving the mutex/map construction above the per-CIDR loop (it already is, at
  `disco.go:107`, since `runScanCheckDisco` takes `foundDevices` as a parameter) — the `wg`/
  `mux` currently declared *inside* the loop body (lines 167-168) need to move to guard
  the whole per-CIDR fan-out instead, one level up.
- Verification: configure the same address count across N CIDR blocks vs. 1 block; total
  wall-clock time should converge (not multiply by N).

### Phase 3 — Delta-only device (re)launch: stop restarting the whole fleet

This is the change most directly aimed at "supporting 5k devices," since §2.5's cost today
scales with total fleet size on every discovery tick, not with the delta.

- Introduce a variant of `runSnmpPolling` (`pkg/inputs/snmp/snmp.go:159–238`) that accepts
  the specific set of devices to (re)launch rather than always iterating `conf.Devices` in
  full (line 176). `Discover`/`addDevices` already compute exactly this set —
  `SnmpDiscoDeviceStat` (`disco.go:34-38`) and the `foundDevices` map passed into
  `addDevices` (`disco.go:392`, iterated at line 419) — so the delta is available; it's
  discarded today at the `SIGUSR2` handoff (`disco.go:213-215`) instead of being passed
  through to the restart.
- Change the `SIGUSR2` contract (`disco.go:206-242`, `snmp_unix.go:20-47`/
  `snmp_windows.go`) to carry the delta (e.g. via a small struct on a channel instead of a
  bare signal) so `wrapSnmpPolling`/`runSnmpPolling` can:
  - Launch pollers for genuinely new devices only (call `launchSnmp`,
    `snmp.go:231/292-348`, per new device, without touching already-running ones).
  - Only `cancel()` (`snmp_unix.go:44`) and relaunch existing devices whose *config*
    actually changed (`stats.replaced`, `disco.go:35`) — not devices that were merely
    re-`Checked` with no material change.
  - Leave devices untouched when `stats.delta == 0 && stats.added == 0` already short-
    circuits correctly today (`disco.go:213`) — that part is fine, the problem is only what
    happens on the *other* branch.
- This requires the most design care of any phase (context lifetime per device instead of
  one shared `ctxSnmp` for the whole generation, `snmp_unix.go:22`) — recommend spiking this
  as its own design pass before implementation, since it changes a fairly central control
  path. Suggest introducing a per-device `context.CancelFunc` map in `KTranslate`/the SNMP
  package so individual devices can be torn down and relaunched independently of the global
  `ctxSnmp`.
- Verification: with a lab config of ~500+ devices and periodic discovery enabled
  (`DiscoveryIntervalMinutes`), confirm that discovering one new device does not interrupt
  in-flight polling/metric continuity for the existing fleet (watch for gaps in
  `metricPoller`'s counter-delta state, `interfaceMetrics.DiscardDeltaState()` calls in
  `metrics/poll.go` around lines 190/203/217 — these should not fire for unaffected devices
  during a delta-only reload).

### Phase 4 — Consumer-side concurrency defaults and a real backpressure signal

- `config.go:420-421`: reconsider the `InputThreads: 1, MaxThreads: 1` defaults — at
  minimum, raise the *default* `MaxThreads` so the existing autoscaler in
  `pkg/cat/kkc.go:565-576` isn't inert out of the box for SNMP-heavy deployments; consider
  scaling the default by expected device count or `runtime.NumCPU()`.
- Make `watchInput` (`kkc.go:556-576`) responsive rather than only additive: track a rolling
  fill-rate rather than a single 60s-interval snapshot (line 557, 566), and consider scaling
  `InputThreads` back down when the queue stays empty, rather than only ever incrementing
  (line 569 never has a corresponding decrement).
- For genuine "adaptive back pressure" (as opposed to "more consumer threads"): consider
  giving producers (the per-device pollers in `metrics/poll.go:226`) a way to sense
  channel pressure directly — e.g. check `len(p.jchfChan)` before blocking and skip/collapse
  a redundant status-flow send (line 229) under load, or expose a shared rate limiter that
  poll loops consult before starting an expensive counter walk
  (`p.Poll(ctx)`, `poll.go:196`) when the shared channel is already near capacity. This
  turns "block everyone equally" into "the newest/lowest-priority work yields first."
- Verification: synthetic slow-sink test (throttle a sink's write rate) with ~5,000
  simulated devices; confirm consumer thread count grows (Phase 4a) and that poll cadence
  degrades gracefully/proportionally rather than every device's `StartLoop` ticker
  (`poll.go:172-173`) backing up identically and simultaneously.

### Phase 5 — Small cleanups (bundle with Phase 3, since they touch the same call paths)

- `pkg/inputs/snmp/mibs/profile.go:302-337` (`checkMatch`): pre-compile `MatchesList`/
  `Matches` regexes once when a `Profile` is loaded (`mibs/load.go`/`mibs/profile.go`
  parsing path) instead of calling `regexp.Compile` on every `FindProfile` call
  (lines 305, 322).
- `pkg/inputs/snmp/snmp.go:350-482` (`parseConfig`): avoid the redundant back-to-back parses
  identified in §2.6 (`disco.go:44` → `disco.go:411` → `snmp.go:161`) — e.g. have
  `addDevices` mutate the already-parsed in-memory `conf` instead of re-reading and
  re-parsing the file from disk when nothing external (global/disco/trap sections) could
  plausibly have changed mid-run.
- `pkg/inputs/snmp/disco.go:279,312`: honor all of `conf.Disco.Ports`, not just `Ports[0]`,
  or explicitly document/log that only the first port is used for credential verification.

---

## 5. Suggested order of operations

1. **Phase 1** first — highest confidence, matches the reported symptom almost exactly, and
   is contained entirely within `disco.go` plus (optionally) a vendored-dependency swap.
2. **Phase 2** next — small, low-risk, same file.
3. **Phase 3** — largest design change, biggest payoff for the "5k devices" scaling goal.
   Worth a short design doc/spike given it touches the restart/context-lifetime model.
4. **Phase 4** — can proceed in parallel with Phase 3 since it's a different subsystem
   (`pkg/cat/kkc.go` vs. `pkg/inputs/snmp`).
5. **Phase 5** — bundle into Phase 3's PR since it touches the same call paths and is easy
   to verify together.

## 6. Before changing behavior: measure first

See `docs/BENCHMARKING_PLAN.md` for the full measurement strategy (deterministic micro-benchmarks
for each bottleneck below, plus a NixOS-VM synthetic device farm for a real end-to-end
IPs/sec and devices/sec number). No phase in §4 should be considered done without a
`benchstat`-style before/after comparison attached to its PR.

### 6.1 Instrumentation to confirm the hypothesis in the field

Before landing fixes, consider adding (temporary or permanent) timing logs around:
- `disco.go:156` (`stb := time.Now()`) through `disco.go:182` — already logs total time per
  CIDR (`"Checked %d ips in %v (from start: %v)"`) but doesn't separate pre-scan time
  (`scanner.Scan`, line 163) from verification time (the `doubleCheckHost` fan-out,
  lines 171-181). Splitting that log line into two would directly confirm whether §2.1 or
  §2.2/§2.3 dominates in the field, before/after each phase.
- A counter for how many addresses fail via immediate error (ICMP unreachable / RST) vs. how
  many burn the full timeout, to validate the §2.2 cost model.

## 7. Test plan

- Unit-level: existing tests in `pkg/inputs/snmp/snmp_test.go`,
  `pkg/inputs/snmp/netbox_test.go` should continue to pass unmodified by Phase 1/2 (no
  netbox-path changes). Add unit coverage for the new bounded-concurrency pre-scan
  replacement and for the delta-only relaunch logic in Phase 3 (e.g. a fake `foundDevices`
  set with 1 new device among 500 existing should result in exactly 1 `launchSnmp` call and
  zero `cancel()`-driven restarts for the other 499).
- Integration/lab: a discovery run against a real or simulated `/16`-scale range (a mock
  SNMP responder farm is fine) with `check_all_ips: true`, `threads` at old vs. new default,
  measuring wall-clock before/after Phase 1.
- Load: simulate ~5,000 devices polling concurrently with a deliberately slow sink, to
  validate Phase 4's backpressure behavior and Phase 3's "no fleet-wide restart" property
  under periodic discovery.

## 8. Risks & rollback

- Raising discovery/verification concurrency (Phase 1) increases simultaneous outbound
  SNMP/TCP traffic — on constrained networks or against IDS/IPS this could trigger
  rate-limiting or alerting that wasn't present at `threads: 4`. Make the new default
  configurable and call it out in release notes/docs, not just a silent code change.
- Phase 3's per-device context lifetime change touches a central control path
  (`wrapSnmpPolling`/`runSnmpPolling`); regressions here risk silently dropping devices from
  polling rather than just being slow. Keep the existing full-restart path available behind
  a flag during rollout, and verify with the delta-tracking test above before removing it.
- All changes are isolated to `pkg/inputs/snmp/**` and `pkg/cat/kkc.go` — no changes to
  wire formats, sink protocols, or config schema are required for Phases 1–3 (Phase 4's
  default-value change is the only config-visible behavior change, and only to *defaults*,
  not to the schema).

---

## Appendix — citation index

| # | File | Lines | What |
|---|---|---|---|
| A1 | `pkg/inputs/snmp/disco.go` | 40-137 | `Discover` — top-level discovery entry point |
| A2 | `pkg/inputs/snmp/disco.go` | 44 | `parseConfig` call #1 per discovery cycle |
| A3 | `pkg/inputs/snmp/disco.go` | 58-61 | `Threads` defaults to 4 |
| A4 | `pkg/inputs/snmp/disco.go` | 82-85 | Verification semaphore (`ctl` channel) |
| A5 | `pkg/inputs/snmp/disco.go` | 139-186 | `runScanCheckDisco` |
| A6 | `pkg/inputs/snmp/disco.go` | 141 | Serial `for` over CIDRs |
| A7 | `pkg/inputs/snmp/disco.go` | 157-163 | Pre-scan invocation, SNMP timeout reused |
| A8 | `pkg/inputs/snmp/disco.go` | 172 | `CheckAll` bypass condition |
| A9 | `pkg/inputs/snmp/disco.go` | 176-181 | Verification fan-out + `wg.Wait()` |
| A10 | `pkg/inputs/snmp/disco.go` | 206-242 | `RunDiscoOnTimer` |
| A11 | `pkg/inputs/snmp/disco.go` | 213-215 | `SIGUSR2` trigger on any delta |
| A12 | `pkg/inputs/snmp/disco.go` | 244-390 | `doubleCheckHost` |
| A13 | `pkg/inputs/snmp/disco.go` | 261-295 | Serial v3 config loop |
| A14 | `pkg/inputs/snmp/disco.go` | 298-329 | Serial version×community loop |
| A15 | `pkg/inputs/snmp/disco.go` | 350 | `FindProfile` call (discovery path) |
| A16 | `pkg/inputs/snmp/disco.go` | 392-585 | `addDevices` |
| A17 | `pkg/inputs/snmp/disco.go` | 411 | `parseConfig` call #2 per discovery cycle |
| A18 | `pkg/inputs/snmp/disco.go` | 279, 312 | `Ports[0]`-only usage |
| A19 | `pkg/inputs/snmp/snmp.go` | 62-121 | `StartSNMPPolls` |
| A20 | `pkg/inputs/snmp/snmp.go` | 129-157 | `initSnmp` |
| A21 | `pkg/inputs/snmp/snmp.go` | 159-238 | `runSnmpPolling` |
| A22 | `pkg/inputs/snmp/snmp.go` | 161 | `parseConfig` call #3 per discovery cycle (restart) |
| A23 | `pkg/inputs/snmp/snmp.go` | 176 | Serial `for` over **all** devices on every restart |
| A24 | `pkg/inputs/snmp/snmp.go` | 205 | `FindProfile` call (restart path) |
| A25 | `pkg/inputs/snmp/snmp.go` | 226 | `apic.EnsureDevice` serial call |
| A26 | `pkg/inputs/snmp/snmp.go` | 231, 292-348 | `launchSnmp` |
| A27 | `pkg/inputs/snmp/snmp.go` | 310, 316 | Two serial `InitSNMP` calls per device |
| A28 | `pkg/inputs/snmp/snmp.go` | 328-345 | Per-device polling goroutines (already parallel) |
| A29 | `pkg/inputs/snmp/snmp.go` | 350-482 | `parseConfig` |
| A30 | `pkg/inputs/snmp/snmp.go` | 457-468 | `FindProfile` loop over all devices, per parse |
| A31 | `pkg/inputs/snmp/snmp_unix.go` | 20-47 | `wrapSnmpPolling` (unix) |
| A32 | `pkg/inputs/snmp/snmp_unix.go` | 22 | Shared `ctxSnmp` for the whole generation |
| A33 | `pkg/inputs/snmp/snmp_unix.go` | 44 | `cancel()` — kills every device's pollers |
| A34 | `pkg/inputs/snmp/snmp_unix.go` | 46 | Recursive relaunch → full re-`runSnmpPolling` |
| A35 | `pkg/inputs/snmp/mibs/profile.go` | 270-300 | `FindProfile` |
| A36 | `pkg/inputs/snmp/mibs/profile.go` | 302-337 | `checkMatch` — regex compiled per call (305, 322) |
| A37 | `pkg/inputs/snmp/metrics/poll.go` | 163-244 | `StartLoop` |
| A38 | `pkg/inputs/snmp/metrics/poll.go` | 226, 229 | Blocking sends into shared `jchfChan` |
| A39 | `pkg/cat/kkc.go` | 44 | `CHAN_SLACK = 8000` |
| A40 | `pkg/cat/kkc.go` | 777-786 | `assureInput` closure, `InputThreads` consumers |
| A41 | `pkg/cat/kkc.go` | 779 | `kc.inputChan` creation |
| A42 | `pkg/cat/kkc.go` | 798, 800 | `assureInput()` + `StartSNMPPolls(ctx, kc.inputChan, ...)` |
| A43 | `pkg/cat/kkc.go` | 556-576 | `watchInput` autoscaler (inert by default) |
| A44 | `pkg/cat/kkc.go` | 565-566 | Autoscale trigger condition |
| A45 | `pkg/cat/kkc.go` | 511-553 | `handleInput` — per-batch consumer work |
| A46 | `config.go` | 308-309, 420-421 | `InputThreads`/`MaxThreads` fields + defaults (1, 1) |
| A47 | `cmd/ktranslate/main.go` | 149, 289, 300 | `metricsChan` creation; `-input_threads`/`-max_threads` flags |
| A48 | `pkg/kt/snmp.go` | 258, 261-262, 294-295 | `Threads`, `CheckAll`, `CheckAllVersions`, global `TimeoutMS`/`Retries` fields |
| A49 | `config/snmp-base.yaml` | 18, 27-28 | Shipped `threads: 4`, `timeout_ms: 3000`, `retries: 0` |
| A50 | `deployment/docker/snmp-base-nr.yaml` | 18, 20, 26-27 | Shipped `threads: 4`, `check_all_ips: true`, `timeout_ms: 3000`, `retries: 0` |
| A51 | `pkg/api/api.go` | 580-624 | `EnsureDevice` — O(n) scan, no-op today (`KentikPlan` never set) |
| A52 | `~/.go/pkg/mod/github.com/liamg/furious@.../scan/scan-device.go` | 37-132 | Vendored `DeviceScanner.Scan` — unbounded goroutines |
| A53 | same file | 78 | `arp.Search` per IP |
| A54 | same file | 98 | Synchronous reverse-DNS lookup |
| A55 | same file | 106 | TCP-connect-to-port-1 liveness probe, SNMP timeout reused |
| A56 | same file | 127 | `wg.Wait()` — blocks until all 65k probes finish |
| A57 | `~/.go/pkg/mod/github.com/liamg/furious@.../scan/scan-syn.go` | 14 | Unused `SynScanner` variant imports `gopacket/pcap` (cgo/`libpcap`), pulled in by package-granularity cgo linking even though `disco.go:159` never calls it |
| A58 | `.github/workflows/test-on-pr.yml` | — | `apt-get install libpcap-dev` step — needed only because of A57 |
