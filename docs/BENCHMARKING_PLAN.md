# Benchmarking Plan — Measuring Before Optimizing

Branch: `investigation`. Companion to `docs/DISCOVERY_PERFORMANCE_PLAN.md`. That document
catalogs suspected bottlenecks with file:line citations (IDs `A1`–`A56` in its Appendix);
this document is the measurement layer that turns "we think X is slow because of Y" into
"here is the before/after number for Y, with statistical confidence." No remediation from
the other plan should be described as done until it has a benchmark result attached.

---

## 0. Principle

Every phase in `docs/DISCOVERY_PERFORMANCE_PLAN.md` §4 gets a benchmark run before the change and
after the change, compared with `benchstat` (§3 below). "It should be faster" is not
evidence; a `benchstat old.txt new.txt` table is.

Two tiers, because "network-bound" and "code-bound" need different treatment:

- **Tier A** — deterministic Go micro-benchmarks that isolate the actual bottleneck *logic*
  (semaphore sizing, serial loops, regex recompilation, channel fan-in) behind an injectable
  fake in place of the real network call. Fast (seconds), zero flakiness, safe to run on
  every PR.
- **Tier B** — a small synthetic SNMP device farm, run as real NixOS virtual machines on a
  real virtual network, driving the actual `Discover()`/`runSnmpPolling` code paths with
  real SNMP/TCP traffic. Slower, has real-world timing variance, run on demand rather than
  gating every PR — but it's the only tier that can faithfully reproduce the *specific*
  network behavior driving the reported symptom (see §2.2).

---

## 1. Nix usage — scope decision (for now)

Nix is being introduced for:

1. **Tier B's benchmark harness** — the NixOS VM test that stands up the synthetic device
   farm (§2.2).
2. **A `devShell`** — a reproducible local dev environment (Go toolchain version, `benchstat`,
   lint tools, `libpcap-dev` for cgo — see `.github/workflows/test-on-pr.yml`'s
   `sudo apt-get install make libpcap-dev` step, which a `devShell` should make unnecessary
   to remember/re-run manually) so anyone picking up this repo gets the same tool versions
   without hand-installing things.
3. **`packages.*.network-agent`** (`nix/network-agent.nix`) — a real ktranslate binary buildable via
   `nix build`, since the binary is now fully static (upstream #14, "go static") and
   distributable on its own. Its `buildPhase` literally shells out to `make all` rather than
   reimplementing the build, so Make remains the single source of truth for *how* to build;
   Nix's job here is limited to vendoring Go module deps reproducibly (`vendorHash`, fetched
   with network access, same as any `buildGoModule` package) and dispatching the build to a
   matching-architecture builder when the evaluating host doesn't have one (see §2.2 on why
   that matters differently for *building* this binary vs. *running* the VM test that uses
   it). Tier B's `collector` VM reuses this same package for its own binary rather than
   building a private copy — one definition of "how to build ktranslate via Nix," not two
   that could drift apart.

Nix is explicitly **not** being adopted for the *release* pipeline. The existing `Makefile` +
`Dockerfile` + `.github/workflows/{create-release,publish-*,ci-build}.yml` remain the only
supported way to produce official released artifacts (Docker images, `.deb`/`.rpm` packages,
etc.) — `packages.*.network-agent` is an additional distribution path and dev convenience, not
a replacement.

This scope is intentionally narrow for now. Revisit if/when there's a concrete reason to
consider Nix for the release pipeline itself (e.g. reproducible release artifacts,
cross-compilation pain) — that's a separate decision with its own tradeoffs (two of the
publish workflows are self-hosted-runner production pipelines with remote buildx and
packagecloud publishing), not a side effect of adopting Nix for benchmarking or packaging a
standalone binary.

---

## 2. The two tiers, in detail

### 2.1 Tier A — deterministic micro-benchmarks

Standard `func BenchmarkX(b *testing.B)`, run via `go test -bench=. -benchmem -run=^$`.
No Nix, no new infra — pure Go, same toolchain as the rest of the repo.

| Benchmark target | Bottleneck ref | File to fake | What's measured | Status |
|---|---|---|---|---|
| Pre-scan fan-out | `A7`, `A52`-`A56` (`disco.go:157-163`; vendored `furious/scan/scan-device.go:37-132`) | Fake probe with configurable latency, gated by a shared `osCeiling` semaphore modeling the finite OS resource a real dial contends for | `BenchmarkPreScanFanOut` (`disco_bench_test.go`): unbounded vs. bounded pool, at 65,536 addresses / ~7.6% live | Done |
| Verification loop `doubleCheckHost` | `A8`, `A9`, `A12`-`A14` (`disco.go:172,176-181,244-390`) | Fake probe, same latency model | `BenchmarkVerificationLoop` (`disco_bench_test.go`): `threads=4` (shipped default, `A49`/`A50`) vs. 64/256 | Done |
| CIDR serialization | `A6` (`disco.go:141`) | Same fake prober | `BenchmarkCIDRSerialization` (`disco_bench_test.go`): serial (current) vs. parallel across 4 CIDRs | Done |
| Restart storm / device (re)init | `A23`-`A27` (`snmp.go:176,205,226,231,310,316`) | **No fake needed** — calls the real `snmp_util.InitSNMP` (confirmed network-free: `gosnmp.Connect()` only opens a local UDP socket) and the real, no-op `apic.EnsureDevice`; does *not* call the real `launchSnmp`, which would leak background goroutines each waiting out a real multi-second SNMP timeout | `BenchmarkDeviceInitLoop` (`snmp_bench_test.go`): devices/sec at fleet sizes 100/1,000/5,000 | Done |
| Regex/profile matching | `A36` (`mibs/profile.go:302-337`, lines 305/322) | None — use real loaded profiles | `BenchmarkFindProfile_MatchesList`, `BenchmarkFindProfile_FleetParse` (`mibs/profile_bench_test.go`) | Done |
| Consumer/backpressure | `A38`, `A40`-`A45` (`kkc.go:226,229,777-786,556-576,511-553`) | Fake per-batch cost in place of real `handleInput` work (constructing a real `*KTranslate` fixture wasn't worth it just for this) | `BenchmarkConsumerThroughput` (`kkc_bench_test.go`): `consumers=1` (shipped default, `A46`) vs. 4/16, at producer counts 100/1,000/5,000 | Done |

All four files exist now:

```
pkg/inputs/snmp/disco_bench_test.go        # pre-scan fan-out, doubleCheckHost, CIDR serialization
pkg/inputs/snmp/snmp_bench_test.go         # restart-storm / device (re)init loop
pkg/inputs/snmp/mibs/profile_bench_test.go # FindProfile / checkMatch
pkg/cat/kkc_bench_test.go                  # consumer throughput / backpressure
```

A first version of the pre-scan fan-out benchmark modeled the probe as a bare
`time.Sleep`, with no shared resource constraint — and found "unbounded" *faster*
than any bounded pool, which would have been a misleading result to leave standing.
`time.Sleep` costs nothing in real OS resources, so Go's scheduler handles tens of
thousands of concurrent sleepers for free; a real dial doesn't get that deal. The
`osCeiling` semaphore in the final version applies the same finite-resource
constraint to every variant, which flips the result to the honest, defensible one:
unbounded fan-out buys no throughput over a correctly-sized pool, only extra
`B/op`/`allocs/op` for goroutines that never needed to exist. Worth remembering when
writing the next Tier A fake — a fake that's *too* free-lunch can quietly measure the
wrong thing.

`benchmarks/baseline.txt` covers all of the above, captured with `-count=5` (not 10,
to keep the combined run under a few minutes — the earlier `-count=10` run across
just the mibs package took ~2 minutes; across all three packages it exceeded 5).
`-count=5` is enough to store as a reference baseline, but `benchstat` will report
`± ∞` (needs ≥6 samples for a confidence interval) — any real before/after comparison
should re-run both sides with `-count=10`+ for statistical confidence, per §3.

**Platform matters more than it looks like it should — and it's not just OS/arch.**
`benchstat` groups results by the full `goos`/`goarch`/`cpu` config line at the top of
each file before comparing anything — two files that differ on *any* of those are
printed as separate, uncompared result sets (no `vs base` column, no error, just
silently no comparison). This bit twice, at two different levels:

1. First version of the CI workflow used a baseline captured on a Mac
   (`darwin/arm64`) compared against CI's `linux/amd64` run — an OS/arch mismatch.
   Fixed by capturing `benchmarks/baseline.txt` on `linux/amd64` itself (via
   `gh run download` against an actual CI run, not a local machine), matching what
   `benchmark.yml` (§3) runs on.
2. That *still* wasn't enough — a second real CI run (same workflow, same
   `ubuntu-latest`) produced the exact same silent-no-comparison symptom again. Cause:
   `ubuntu-latest` is not hardware-uniform. Run #1 landed on `cpu: AMD EPYC 9V45
   96-Core Processor`; run #2 landed on `cpu: INTEL(R) XEON(R) PLATINUM 8573C`. Same
   OS, same arch, different silicon — and `benchstat` treats that as a different
   config just like it would a different OS.

Fix: `benchstat -ignore cpu` (confirmed via `benchstat -h`: `-ignore keys` — "ignore
variations in keys"), applied in both `benchmark.yml` and the `Justfile`'s
`bench-diff` recipe. Verified by reproducing the exact failure locally first
(same baseline file, one copy with its `cpu:` line hand-edited to a different string)
and confirming `-ignore cpu` restores the `vs base` column — not just re-running CI
and hoping. `goos`/`goarch` are deliberately still left as real grouping keys (a
genuine cross-platform difference should stay flagged); only `cpu` — noise on a
shared runner fleet, not signal — is ignored.

Consequence that's still true: running `just bench-diff` on a Mac against this
`linux/amd64` baseline will still hit a real `goos`/`goarch` mismatch (which
`-ignore cpu` doesn't and shouldn't paper over). That's a known, accepted tradeoff
for now (one baseline, not a per-platform set) — locally on a Mac, either capture a
separate local baseline for that session, or just eyeball the absolute numbers
rather than expecting a `vs base` column. CI is the one place this baseline is
actually meant to produce a real comparison.

### 2.2 Tier B — NixOS VM synthetic device farm

**Why not just fake the network in Go for this too:** the dominant real-world cost
identified in the other plan (§2.1/§2.2 there) hinges on a distinction a Go-level fake on
loopback cannot reproduce — a **silently dropped** packet (firewalled/filtered, the common
enterprise case, forces the full `timeout_ms` wait) vs. an **actively rejected** one
(closed port, fast RST/ICMP-unreachable). On loopback, "nothing is listening" always fails
fast — you cannot get the slow case without a real kernel and a real drop rule. NixOS VM
tests give you both, on a real virtual network, reproducibly.

Implemented in `nix/tests/snmp-discovery-bench.nix`, wired into `flake.nix`'s `checks`
output. Confirmed working end-to-end, real numbers below — not a sketch.

**Topology:**
- One `collector` node: runs the real ktranslate binary (§1, `nix/network-agent.nix`) against
  the farm's address range, using a real `snmp.yml`
  discovery config that mirrors the shipped `deployment/docker/snmp-base-nr.yaml` example
  (same `threads`, `timeout_ms`, `retries`, and — deliberately — `check_all_ips: true`, so
  the benchmark measures the actual shipped configuration, including its full-subnet-sweep
  behavior, not a hypothetical narrower one).
- N `device` nodes: run real `net-snmp`'s `snmpd` (packaged in nixpkgs), each on its own
  address in the farm's virtual `/24`, generated programmatically
  (`builtins.listToAttrs (map ... (lib.range 1 N))`) and split by `respondFrac`/`rejectFrac`
  parameters.
- `networking.firewall.rejectPackets` (default `false` = silent DROP) gives the
  respond/reject/drop distinction natively, with no manual iptables rules needed: leave it
  at the default for **drop** nodes (silent, forces the full timeout); set it `true` for
  **respond** and **reject** nodes (fast RST/ICMP-unreachable) — respond nodes need this
  too, not just `services.snmpd.enable`, or the pre-scan's TCP-dial liveness probe gets
  silently dropped there as well, corrupting the "respond nodes are fast" half of the
  measurement. The remainder of non-respond/reject addresses split evenly between an
  explicit **drop** node and simply omitting a node at that address at all
  (**unclaimed**) — see §7 on why these two aren't yet confirmed to behave identically.

**Mechanics:** `pkgs.testers.runNixOSTest` with a standard `nodes = { ... }` (QEMU/KVM)
definition — not the newer systemd-nspawn `containers = { ... }` backend, which was tried
first (needs no virtualization at all, looked appealing) but unconditionally requires the
Nix daemon feature `uid-range`, needing `auto-allocate-uids`/`cgroups` enabled in the
*executing* machine's own `nix.conf` plus a matching feature declaration on the calling
side — an out-of-repo system change that turned out to be unnecessary once the real
execution model below was understood correctly.

**Execution model — the part that wasn't obvious going in.** On Apple Silicon, these VMs
run **natively on the host Mac**, not inside nix-darwin's `linux-builder`. That distinction
matters because the `linux-builder` (itself a VM, via Apple's Virtualization.framework)
cannot do nested virtualization — confirmed both by the official NixOS wiki ("As it
happens, M1 doesn't support nested virtualization. So it can run a Linux builder, or any
other NixOS virtual machine, but it cannot do so inside e.g. the Linux builder") and by
this repo's own attempt to fix it via `uid-range`, which was solving the wrong problem.
The actual fix: since nixpkgs#282401 (2024, see
[nixcademy's writeup](https://nixcademy.com/posts/running-nixos-integration-tests-on-macos/)),
`runNixOSTest`'s qemu process and Python test driver run directly on whichever host
*evaluates* the test — so evaluating this file's `pkgs` as `aarch64-darwin`/`x86_64-darwin`
(not `aarch64-linux`/`x86_64-linux`) makes the test execute right there on the Mac, via
`apple-virt`/HVF acceleration, with zero `linux-builder` involvement for the test-*run*
itself. `flake.nix`'s `checks` output does this by mapping each of the four `forAllSystems`
systems to the matching-arch **Linux** system only for `collectorBin` (`packages.*` stays
Linux-only, since that's a real binary the guest VMs need), while `pkgs.testers.runNixOSTest`
itself gets called with the system's own (possibly Darwin) `pkgs`. `nix.settings.system-features`
needing `nixos-test`/`apple-virt` is auto-detected by Nix 2.19+ with no nix-darwin config
change required at all (confirmed via `nix show-config system-features` locally). See
`nix/tests/minimal-ping.nix` — a fast (~15-20s), permanently-kept two-node ping-pong sanity
check — for the minimal reproduction that confirmed this mechanism before trusting it on
the real, much larger test.

**Scale, confirmed empirically (not assumed):**

| Scale | Where | Result |
|---|---|---|
| 12 nodes (7 respond / 2 reject / 1 drop / 2 unclaimed) | Local (Apple Silicon Mac, native `apple-virt`) | Completes in ~4-4.5 min; discovery itself ~194s, dominated by `check_all_ips: true` sweeping the full `/24` (~240 unclaimed addresses beyond the 12 defined ones), each paying real timeout cost — a small-scale, faithful reproduction of the exact field-reported pattern ("65,000 IPs scanned, ~5,000 real devices") this tier exists to validate against. |
| 40 nodes (70/20/10 split, the original target topology) | Local, same Mac | Ran over an hour under heavy, sustained CPU contention (41 concurrent qemu processes; load average 25-40) without completing — not crashing, genuinely resource-bound. |
| 12 nodes | CI (`ubuntu-latest`, 4 vCPU/16GB) | **Failed** — the collector VM never became interactive: `RuntimeError: Shell did not start in time`. Confirmed by reading nixpkgs source (`nixos/lib/test-driver`'s `connect()`): this is a **hardcoded, non-configurable** 10 retries × 30s = 5-minute wait, not a NixOS option. With `virtualisation.cores` defaulting to 1/VM, 12 concurrent VMs on a 4-vCPU runner is a real 3x oversubscription — genuine CPU starvation during boot-time systemd/dbus activation, not a fluke. |
| 8 nodes (5 respond / 1 reject / 1 drop / 1 unclaimed) | CI, same runner | **Succeeds** — confirmed via two real runs (`device_count: 5, node_count: 8` both times), ~7 min total. This is what `benchmark-tier-b.yml` actually runs. |

So the field is now three scales, not two: 12-node smoke (local iteration), 8-node CI (what actually runs in CI, sized for a standard runner's 4 vCPUs), and the 40-node target (local/aspirational on beefier hardware only — confirmed *not* to fit a standard GitHub-hosted runner at all, not just "not yet confirmed").

Per-VM memory needed hand-tuning to get even the 40-node case to boot without an outright
crash: the test framework's default `virtualisation.memorySize` (~1024 MiB) times dozens of
concurrently-running device VMs first caused severe memory-pressure thrashing (25GB+
compressed), and overcorrecting to 192 MiB caused an actual VM boot failure (a
respond-category node loading `net-snmp` disconnected the test-driver shell entirely).
384 MiB for device nodes / 768 MiB for the collector was the value that held.

Consequently, `flake.nix` exposes **three** checks per system: `snmp-discovery-bench-smoke`
(`deviceCount = 12`, used as the `Justfile`'s `bench-tier-b` default for routine local
iteration), `snmp-discovery-bench-ci` (`deviceCount = 8`, what `benchmark-tier-b.yml`
actually runs), and `snmp-discovery-bench` (the real 40-node/70-20-10 target,
`Justfile`'s `bench-tier-b-full` — local/aspirational use on hardware with cores to spare,
not something CI runs). This is a calibration point, not a literal full-scale
(5,000-65,000) replica — see §5 for how it combines with Tier A to reason about full scale.

**CI feasibility (confirmed, not assumed):** `.github/workflows/benchmark-tier-b.yml` runs
`checks.x86_64-linux.snmp-discovery-bench-ci` on `ubuntu-latest` via
`cachix/install-nix-action@v31` with `enable_kvm: true` — that runner is already
`x86_64-linux`, so the VM test runs its qemu process natively via `kvm` right there, no
`apple-virt`/`linux-builder` cross-system complexity involved at all (that's specific to
running this locally from Apple Silicon). Verified with two real runs (via a temporarily
added `pull_request` trigger on the implementing PR, removed again once confirmed): the
12-node smoke scale fails there (see the table above), the 8-node CI scale succeeds
consistently.

---

## 3. Tooling and reporting

- `go test -bench=. -benchmem -run=^$ ./...` — Tier A. `-run=^$` skips normal tests so only
  benchmarks execute.
- `benchstat` (`golang.org/x/perf/cmd/benchstat`, not currently installed — confirmed via
  `which benchstat`) turns two benchmark runs into a statistical diff (mean, variance,
  percent change with confidence), which is what makes a claim like "Phase 1 made
  verification 40% faster" defensible rather than anecdotal. Workflow for any change:
  1. On the pre-change code, run benchmarks, save `old.txt`.
  2. Make the change.
  3. Run benchmarks again, save `new.txt`.
  4. `benchstat old.txt new.txt` → paste the table into the PR description as evidence.
- Tier B's report is simpler (one real number per run, not a statistical distribution over
  many fast iterations): `elapsed_s`, `device_count`, `node_count`, compared against
  `benchmarks/tier-b-baseline.json` (a checked-in snapshot from a real CI run of
  `snmp-discovery-bench-ci`, not a local one — same "capture on the hardware the
  comparison actually runs on" principle as Tier A's baseline, §2.1). Reported the same
  way as Tier A: a markdown table to `$GITHUB_STEP_SUMMARY` and a sticky PR comment.
  Refresh the baseline the same way Tier A's is refreshed — from an actual CI run's
  result, via `gh run download` or reading the job summary, not a local run.

### CI wiring

This repo's CI on `investigation` is deliberately manual/opt-in
(`.github/workflows/test-on-pr.yml` is `workflow_dispatch`-only; auto-triggers were
disabled; `.github/workflows/ci-build.yml` runs on
`push: [investigation]` + `pull_request`). A new benchmark workflow should match that
convention rather than gate every push:

- `.github/workflows/benchmark.yml` (Tier A) — `workflow_dispatch` + `push: [develop]` +
  `pull_request` (path-filtered). Checkout → `cachix/install-nix-action@v31` →
  `nix develop --command` runs the benchmarks and `benchstat -ignore cpu` against the
  checked-in baseline → writes to `$GITHUB_STEP_SUMMARY` and a sticky PR comment
  (`marocchino/sticky-pull-request-comment@v2`) → uploads raw output as an artifact.
- `.github/workflows/benchmark-tier-b.yml` (Tier B) — **separate workflow**, not a second
  job here: it takes noticeably longer per run than Tier A (VM boot included), so it's
  path-filtered to only the code that plausibly changes what it measures
  (`pkg/inputs/snmp/**`, `nix/tests/**`) rather than running on every PR regardless of
  relevance. Triggers: `workflow_dispatch` + nightly `schedule` (drift tracking without
  needing someone to remember to run it) + `push: [develop]` + `pull_request`
  (path-filtered), mirroring Tier A's own trigger shape. Checkout →
  `cachix/install-nix-action@v31` (`enable_kvm: true`) →
  `nix build .#checks.x86_64-linux.snmp-discovery-bench-ci` (the 8-node CI-sized target,
  §2.2 — **not** the 40-node one, confirmed too large for a standard runner) → compare
  against `benchmarks/tier-b-baseline.json` → write to `$GITHUB_STEP_SUMMARY` and a
  sticky PR comment (`marocchino/sticky-pull-request-comment@v2`) → upload the raw
  result as an artifact.
- Regression policy: start by just reporting the diff for humans to judge (this is an
  investigation branch, not a release pipeline) rather than failing the build on a
  threshold. Tighten later once the numbers are trusted.
- Caching: nixpkgs' own packages (the NixOS base system, `net-snmp`) come from the public
  `cache.nixos.org` substituter by default — no project-specific binary cache (e.g. Cachix)
  is needed to get a fast *base* system. Only worth adding a project cache later if the test
  derivation itself grows expensive to rebuild across runs.

---

## 4. Proposed file layout

```
flake.nix                                   # devShell + network-agent package + Tier B checks
flake.lock
nix/network-agent.nix                       # builds the real ktranslate binary (§1) -- reused by Tier B
nix/tests/minimal-ping.nix                  # fast sanity check for the execution model (§2.2)
nix/tests/snmp-discovery-bench.nix          # the runNixOSTest definition (§2.2)
pkg/inputs/snmp/disco_bench_test.go         # Tier A
pkg/inputs/snmp/snmp_bench_test.go          # Tier A
pkg/inputs/snmp/mibs/profile_bench_test.go  # Tier A
pkg/cat/kkc_bench_test.go                   # Tier A
benchmarks/baseline.txt                     # Tier A baseline for benchstat diffing
benchmarks/tier-b-baseline.json             # Tier B baseline (from a real CI run, not local)
Justfile                                    # bench*, bench-tier-b, bench-tier-b-full recipes
.github/workflows/benchmark.yml             # Tier A
.github/workflows/benchmark-tier-b.yml      # Tier B
```

---

## 5. How Tier A and Tier B combine to reason about full scale

Tier B gives a real, small-scale (tens–hundreds of nodes) throughput number with faithful
network behavior. Tier A gives a cheap, large-scale (thousands of synthetic addresses/
devices) number with faked network behavior but the *same* concurrency-cap/serialization
logic as production. Use Tier B to confirm Tier A's fake is calibrated correctly (e.g. "at
80 nodes with a realistic drop ratio, Tier B says X seconds; does Tier A's fake-latency
model predict something close to X when configured with the same node count and ratio?")
— once that calibration checks out, trust Tier A's extrapolation to 5,000 devices / 65,000
IPs, since literally running that many VMs isn't practical in CI.

---

## 6. Validating the benchmarks themselves

- Tier A: run each benchmark with `-count=10` and check `benchstat`'s own variance output
  before trusting a single run — Go benchmarks can be noisy on shared CI runners; a change
  claimed as "faster" should show up as a statistically significant delta, not just a
  different single sample.
- Tier B: run the same topology twice before changing anything, to establish the run-to-run
  variance baseline for a real network/VM test (expect more variance than Tier A — real
  timeouts, real scheduler jitter) before treating any single before/after pair as
  conclusive.

---

## 7. Open questions / follow-ups

Resolved during implementation:

- Getting the ktranslate binary into the `collector` VM: `environment.systemPackages =
  [ collectorBin ]` with `collectorBin` from `packages.*.network-agent` (`nix/network-agent.nix`)
  — a normal Nix store path closure-referenced into the VM, no manual mount/copy step
  needed.
- Exact discovery-output key shape: confirmed `<name>__<ip>:` (`disco.go:421,429`) — the
  test counts devices via `grep -c '__' <output file>`.
- No CLI flags beyond `-snmp=... -snmp_discovery=true -snmp_out_file=... -log_level=info`
  were needed for `Discover()` to run to completion.
- Where the `benchmarks/baseline.txt` Tier A baseline comes from: a checked-in snapshot
  captured from an actual CI run (see §2.1) — already done, not hypothetical.
- **Whether 40 concurrent VMs (or even 12) fit a standard GitHub-hosted runner.** They
  don't — confirmed with two real `benchmark-tier-b.yml` runs (via a temporarily added
  `pull_request` trigger, removed once confirmed): 12 nodes hit the NixOS test driver's
  hardcoded 5-minute boot-shell timeout (`RuntimeError: Shell did not start in time`,
  `nixos/lib/test-driver`'s `connect()` — not a NixOS option, can't be raised), and the
  root cause (CPU oversubscription: `virtualisation.cores` defaults to 1/VM against the
  runner's 4 physical vCPUs) means the 40-node target almost certainly never fits a
  standard runner at all. Fixed by adding `snmp-discovery-bench-ci` (8 nodes, still hits
  all four respond/reject/drop/unclaimed categories) as what CI actually runs — confirmed
  passing twice in a row.

Still open:

- **Unclaimed address vs. explicit drop node — not yet distinguished empirically.** Both
  currently produce a silent timeout from the collector's point of view (no host to ARP for
  vs. a host that drops at the firewall), and nothing in the confirmed runs above isolated
  their timing from each other. Worth a targeted comparison once someone's looking at
  per-address timing rather than just aggregate elapsed time.
- **Resolved: Tier B now has a correctness assertion, not just a measurement.** Before
  this, the test recorded `device_count`/`elapsed_s` without checking either against an
  expectation — a regression that silently discovered 0 devices, or took 10x longer,
  would not have failed the test at all. `snmp-discovery-bench.nix`'s test script now
  asserts `device_count == <respond-category node count>` right after computing it, so
  `Discover()` actually misbehaving (not just crashing outright) fails the test.
- **Resolved: comparison reporting mirrors Tier A** — `benchmarks/tier-b-baseline.json`
  (a checked-in snapshot from a real CI run) plus a markdown-table step in
  `benchmark-tier-b.yml`, posted as a sticky PR comment on `pull_request` (path-filtered
  to `pkg/inputs/snmp/**`/`nix/tests/**`) in addition to the nightly/dispatch run. A
  single-run number, not a statistical distribution like `benchstat` — expect more
  run-to-run variance than Tier A (real VM boot/scheduler jitter). Whether that variance
  is small enough for the comparison to be trustworthy without multiple runs is not yet
  characterized — see §6's stated intent to run the same topology twice before trusting
  a before/after pair, which hasn't actually been done yet for Tier B.
- Whether the real 40-node/70-20-10 target (`snmp-discovery-bench`) is worth running
  *anywhere* automated, given it doesn't fit standard CI hardware — options include a
  bigger (paid) GitHub-hosted runner tier, or accepting it as a manual/local-only
  calibration check run occasionally on beefier hardware. Not decided yet.

---

## Appendix — cross-reference to `docs/DISCOVERY_PERFORMANCE_PLAN.md`

This document deliberately does not repeat the bottleneck citations (`A1`-`A56`) already
recorded there — see that file's Appendix for the exact `file:line` source of every
bottleneck a Tier A or Tier B benchmark in this document is designed to measure.
