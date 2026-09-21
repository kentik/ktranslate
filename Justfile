# Dev tooling helpers for this fork. Run inside `nix develop` (flake.nix) to guarantee the
# tools each recipe needs (benchstat via goperf, go-licence-detector, ...) are present.
#
# Benchmarking recipes below -- see BENCHMARKING_PLAN.md #3.

CURRENT_SYSTEM := `nix eval --impure --raw --expr builtins.currentSystem`

default:
    @just --list

# Download the MaxMind GeoLite2 databases into maxmind-dbs/, for a local `docker build`
# (see Dockerfile's maxmind stage -- CI/release builds get these from actions/cache
# instead, via ci-build.yml/publish-release.yml). Needs MM_ACCOUNT_ID/MM_DOWNLOAD_KEY in
# the environment -- the same names those workflows read from repo secrets. A free
# GeoLite2 account works: https://www.maxmind.com
maxmind-dbs dest="maxmind-dbs":
    #!/usr/bin/env bash
    set -euo pipefail
    if [ -z "${MM_DOWNLOAD_KEY:-}" ]; then
        echo "MM_DOWNLOAD_KEY (MaxMind license key) not set" >&2
        exit 1
    fi
    mkdir -p {{dest}}
    tmp="$(mktemp -d)"
    trap 'rm -rf "$tmp"' EXIT
    curl -sfL -o "$tmp/country.tar.gz" -u "${MM_ACCOUNT_ID:-}:$MM_DOWNLOAD_KEY" "https://download.maxmind.com/geoip/databases/GeoLite2-Country/download?suffix=tar.gz"
    tar zxf "$tmp/country.tar.gz" --strip-components 1 -C {{dest}}
    curl -sfL -o "$tmp/asn.tar.gz" -u "${MM_ACCOUNT_ID:-}:$MM_DOWNLOAD_KEY" "https://download.maxmind.com/geoip/databases/GeoLite2-ASN/download?suffix=tar.gz"
    tar zxf "$tmp/asn.tar.gz" --strip-components 1 -C {{dest}}

# Run Tier A benchmarks for a package (default: everything).
bench pkg="./...":
    go test {{pkg}} -bench=. -benchmem -run=^$

# Run a package's benchmarks `count` times -- enough samples for benchstat to compare.
bench-count pkg count="10":
    go test {{pkg}} -bench=. -benchmem -run=^$ -count={{count}}

# Refresh the checked-in baseline for a package.
bench-baseline pkg dest="benchmarks/baseline.txt":
    #!/usr/bin/env bash
    set -euo pipefail
    mkdir -p "$(dirname {{dest}})"
    just bench-count {{pkg}} 10 | tee {{dest}}

# Compare the current code's benchmarks against the checked-in baseline.
bench-diff pkg baseline="benchmarks/baseline.txt":
    #!/usr/bin/env bash
    set -euo pipefail
    tmp="$(mktemp)"
    trap 'rm -f "$tmp"' EXIT
    just bench-count {{pkg}} 10 > "$tmp"
    # -ignore cpu: benchmarks/baseline.txt is captured on linux/amd64 CI hardware,
    # which varies between runs (see BENCHMARKING_PLAN.md) -- and on a Mac this is
    # also a real goos/goarch mismatch, so don't expect a meaningful `vs base` delta
    # from this locally, only from CI's own comparison.
    benchstat -ignore cpu {{baseline}} "$tmp"

# Run the Tier B synthetic SNMP device farm at its small (12-node) scale --
# BENCHMARKING_PLAN.md #2.2. Fast enough for routine local iteration (confirmed
# end-to-end in a few minutes). Defaults to the current host's system -- on Apple
# Silicon this runs the VMs natively via apple-virt/HVF, no linux-builder involved
# (see nix/tests/snmp-discovery-bench.nix); in CI (system=x86_64-linux) it runs
# natively via kvm.
bench-tier-b system=CURRENT_SYSTEM:
    nix build .#checks.{{system}}.snmp-discovery-bench-smoke -L --print-out-paths

# Run the full 40-node/70-20-10 target topology. CI-scale, not laptop-scale: expect
# tens of minutes under any real resource contention -- see snmp-discovery-bench in
# flake.nix's checks output for why this isn't the local default.
bench-tier-b-full system=CURRENT_SYSTEM:
    nix build .#checks.{{system}}.snmp-discovery-bench -L --print-out-paths

# Regenerate THIRD_PARTY_NOTICES.md from go.mod (direct + indirect deps).
third-party-notices out="THIRD_PARTY_NOTICES.md":
    # go-licence-detector reads LICENSE files out of the local module cache -- it doesn't
    # fetch them itself, so on a cold cache (e.g. a fresh CI runner) it silently produces
    # a notices file with zero package entries instead of erroring.
    go mod download all
    go list -mod=mod -m -json all | go-licence-detector \
        -includeIndirect \
        -rules assets/licence/rules.json \
        -overrides assets/licence/overrides.json \
        -noticeTemplate assets/licence/THIRD_PARTY_NOTICES.md.tmpl \
        -noticeOut {{out}}

# Verify THIRD_PARTY_NOTICES.md is up to date with go.mod.
third-party-notices-check:
    #!/usr/bin/env bash
    set -euo pipefail
    tmp="$(mktemp)"
    trap 'rm -f "$tmp"' EXIT
    just third-party-notices "$tmp"
    diff "$tmp" THIRD_PARTY_NOTICES.md

# --- Release helpers ---------------------------------------------------------
#
# Drive publish-release.yml's two entry points from the CLI instead of the GitHub
# release UI, so the commit a pre-release and its eventual full release point at
# is pinned explicitly (not "whatever the target branch's HEAD happens to be when
# I click Publish") and the checks the workflow itself would fail on (tag already
# exists, no matching pre-release to promote, ...) surface here first. See
# docs/RELEASING.md for the full walkthrough. Both need `gh` authenticated
# (`gh auth login`) and `nix` on PATH.

# Cut a pre-release: publishes a GitHub pre-release tagged v<version> at `ref`
# (default: current HEAD), which triggers publish-release.yml's `publish` job to
# build and push newrelic/network-agent:<version> and
# newrelic/network-agent:sha-<commit> to Docker Hub. `version` must carry a
# SemVer prerelease suffix (e.g. "0.1.0-rc1") -- same rule publish-release.yml
# itself enforces, checked here first so a typo fails locally instead of via a
# round trip through a created-then-broken GitHub release.
release-rc version ref="HEAD":
    #!/usr/bin/env bash
    set -euo pipefail
    command -v gh >/dev/null || { echo "error: gh CLI not found -- see https://cli.github.com" >&2; exit 1; }
    tag="v{{version}}"

    case "{{version}}" in
      *-*) ;;
      *)
        echo "error: '{{version}}' has no prerelease suffix (e.g. '{{version}}-rc1') -- a bare version is reserved for 'just release-promote'." >&2
        exit 1
        ;;
    esac
    nix run .#check-semver -- "{{version}}"

    sha="$(git rev-parse "{{ref}}")"

    if git rev-parse -q --verify "refs/tags/$tag" >/dev/null || \
       git ls-remote --exit-code origin "refs/tags/$tag" >/dev/null 2>&1; then
      echo "error: tag $tag already exists (locally or on origin)." >&2
      exit 1
    fi
    if gh release view "$tag" >/dev/null 2>&1; then
      echo "error: a GitHub release already exists for $tag." >&2
      exit 1
    fi

    echo "Cutting pre-release $tag at $sha ..."
    gh release create "$tag" --target "$sha" --prerelease --title "$tag" --generate-notes

    echo
    echo "Pre-release $tag published. Once publish-release.yml finishes, Docker Hub will have:"
    echo "  newrelic/network-agent:{{version}}"
    echo "  newrelic/network-agent:sha-$sha"
    echo
    echo "To promote this exact commit to a full release once it's been tested:"
    echo "  just release-promote <final-version> $tag"

# Promote a tested pre-release to a full release: checks that `from_tag` exists
# and is a published GitHub pre-release, and that `version` (a bare SemVer, no
# prerelease suffix -- e.g. "0.1.0") hasn't already been released, then
# publishes a full GitHub release at the *same commit* as `from_tag`. This
# triggers publish-release.yml's `promote` job, which retags the pre-release's
# already-pushed newrelic/network-agent:sha-<commit> image as <version> and
# `latest` -- no rebuild, so what ships is byte-identical to what the
# pre-release build tested.
release-promote version from_tag:
    #!/usr/bin/env bash
    set -euo pipefail
    command -v gh >/dev/null || { echo "error: gh CLI not found -- see https://cli.github.com" >&2; exit 1; }
    tag="v{{version}}"
    from="{{from_tag}}"
    case "$from" in v*) ;; *) from="v$from" ;; esac

    case "{{version}}" in
      *-*)
        echo "error: '{{version}}' has a prerelease suffix -- release-promote is for the final version (e.g. '0.1.0'), not the pre-release itself." >&2
        exit 1
        ;;
    esac
    nix run .#check-semver -- "{{version}}"

    if git rev-parse -q --verify "refs/tags/$tag" >/dev/null || \
       git ls-remote --exit-code origin "refs/tags/$tag" >/dev/null 2>&1; then
      echo "error: tag $tag already exists -- this version has already been released." >&2
      exit 1
    fi
    if gh release view "$tag" >/dev/null 2>&1; then
      echo "error: a GitHub release already exists for $tag." >&2
      exit 1
    fi

    git fetch origin --tags --quiet
    sha="$(git rev-parse -q --verify "refs/tags/$from^{commit}" 2>/dev/null)" || {
      echo "error: tag $from not found -- cut a pre-release first with 'just release-rc'." >&2
      exit 1
    }

    is_prerelease="$(gh release view "$from" --json isPrerelease --jq '.isPrerelease' 2>/dev/null)" || {
      echo "error: no GitHub release found for $from -- promoting requires a published pre-release, not a bare tag." >&2
      exit 1
    }
    if [ "$is_prerelease" != "true" ]; then
      echo "error: $from exists but is not marked as a pre-release on GitHub." >&2
      exit 1
    fi

    # Best-effort: confirm the pre-release's own build actually succeeded, so a
    # broken/incomplete Docker Hub push doesn't surprise the promote job below.
    # Not a hard gate -- gh's run-listing by tag isn't guaranteed exhaustive, and
    # publish-release.yml's own promote job is the authoritative, unskippable check.
    if ! gh run list --workflow=publish-release.yml --json event,headBranch,conclusion \
         --jq ".[] | select(.event == \"release\" and .headBranch == \"$from\" and .conclusion == \"success\")" \
         | grep -q .; then
      echo "warning: no successful publish-release.yml run found for $from -- its image may not actually be on Docker Hub yet. The promote step below will fail if it isn't." >&2
    fi

    echo "Promoting $from (commit $sha) to full release $tag ..."
    gh release create "$tag" --target "$sha" --title "$tag" --notes "Promoted from $from."

    echo
    echo "Full release $tag published. publish-release.yml's promote job will retag the"
    echo "existing newrelic/network-agent:sha-$sha image as {{version}} and latest -- no rebuild."
