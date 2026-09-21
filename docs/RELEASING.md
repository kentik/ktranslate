# Releasing

How to cut a pre-release and promote it to a full release, driving
`.github/workflows/publish-release.yml` from the CLI instead of the GitHub release UI.

## Why not just use the GitHub UI

`publish-release.yml` has two build paths:

- Publishing a **pre-release** (tag has a SemVer `-suffix`, e.g. `v0.1.0-rc1`) builds the
  image from source and pushes `newrelic/network-agent:<version>` and
  `newrelic/network-agent:sha-<commit>` to Docker Hub.
- Publishing (or promoting an existing pre-release to) a **full release** (bare tag, e.g.
  `v0.1.0`) never rebuilds — it retags the `sha-<commit>` image already pushed by that same
  commit's pre-release as `<version>` and `latest`, via `docker buildx imagetools create`.
  This is deliberate: what ships as the real release is byte-identical to what the
  pre-release build tested. If no matching pre-release image exists for that commit, the
  `promote` job fails on purpose.

That last part is where the GitHub UI gets risky: creating a release by typing a tag name
defaults its target to whatever the branch's HEAD happens to be *at that moment*. If you
create the pre-release, then later create the full release the same way, and someone pushed
to the branch in between, the two tags land on different commits — and `promote` fails
because the full release's commit has no matching `sha-<commit>` image.

`just release-rc` / `just release-promote` (in the `Justfile`) exist to make the commit
explicit and to catch the workflow's own failure conditions locally, before a release is
even created.

## Prerequisites

- `gh` authenticated (`gh auth login`) — available in `nix develop`'s devShell, or install
  separately: https://cli.github.com
- `nix` on `PATH`, for the shared SemVer check in `nix/semver.nix`
  (`nix run .#check-semver`) — the exact same rule `publish-release.yml` itself enforces.
- The commit you're cutting a release from must already be pushed to `origin` — `gh release
  create --target` can only tag a commit GitHub already has.

## Cutting a pre-release

```
just release-rc 0.1.0-rc1
```

This publishes a GitHub pre-release tagged `v0.1.0-rc1` at the current `HEAD`, which
triggers `publish-release.yml`'s `publish` job. Pass a second argument to target a
different commit/branch/tag instead of `HEAD`:

```
just release-rc 0.1.0-rc1 some-branch
```

Before creating anything, the recipe checks that:

- `0.1.0-rc1` actually has a prerelease suffix (a bare version is rejected — that's
  `release-promote`'s job, not this one's).
- It's valid SemVer (via `nix run .#check-semver`).
- `v0.1.0-rc1` doesn't already exist as a tag (locally or on `origin`) or as a GitHub
  release.

Once `publish-release.yml` finishes, Docker Hub has `newrelic/network-agent:0.1.0-rc1` and
`newrelic/network-agent:sha-<commit>`. Test that image. If something's wrong, fix it, commit,
and cut `0.1.0-rc2` at the new commit — never move a tag.

## Promoting a tested pre-release to a full release

Once `0.1.0-rc1` (or whichever RC) is verified good:

```
just release-promote 0.1.0 0.1.0-rc1
```

This publishes a full GitHub release tagged `v0.1.0`, targeted at the **exact same commit**
as `v0.1.0-rc1` — resolved from the existing tag, not re-typed — which triggers
`publish-release.yml`'s `promote` job. No image is rebuilt; the existing
`newrelic/network-agent:sha-<commit>` is retagged as `0.1.0` and `latest`.

Before creating anything, the recipe checks that:

- `0.1.0` has no prerelease suffix and is valid SemVer.
- `v0.1.0` doesn't already exist (locally, on `origin`, or as a GitHub release) — you can't
  promote onto an already-released version.
- `0.1.0-rc1` exists as a tag and as a **published GitHub pre-release** specifically (not
  just a bare git tag someone pushed by hand).
- (Best-effort, non-fatal) `publish-release.yml` actually completed successfully for
  `0.1.0-rc1` — if this can't be confirmed, you get a warning, not a hard stop, since
  `promote`'s own Docker Hub check is the real, unskippable gate.

If any of the hard checks fail, fix the underlying issue (cut a new RC, pick the right
`from_tag`, etc.) rather than working around the recipe.

## What if I need to re-run a pre-release build?

Re-running/re-publishing the *same* pre-release (same tag) is expected to work — it rebuilds
and re-pushes, which is fine for retrying something flaky. `release-rc`'s "tag must not
already exist" check only stops you from reusing an RC tag for a *different* commit; if you
genuinely need to retry the exact same tag, do that through the GitHub UI's "re-run" on the
existing release/workflow run rather than through this recipe.

## About the `VERSION` file

`VERSION` (repo root) is a fallback default, not the source of truth for what actually gets
released:

- **Nix's `packages.*.network-agent`** (`nix/network-agent.nix`) reads it as a pinned,
  Nix-evaluation-pure version string, specifically so the derivation's cache isn't
  invalidated on every single commit the way stamping in `self.rev` would be. It will lag
  behind the latest tag between bumps — that's expected, the same way a Nix package's
  version generally lags upstream between packaging updates.
- **Local `make` builds** default to it (`Makefile`'s `NETWORK_AGENT_VERSION ?= ...`).
- **`workflow_dispatch`** runs of `publish-release.yml` fall back to it when no version
  input is given.
- **Real releases ignore it entirely** — `publish-release.yml`'s `release`-event path always
  uses the git tag name, never this file.

It's checked for valid SemVer format by `version-format-check.yml`, but nothing currently
checks that it matches the latest actual release — don't rely on it to answer "what's the
current released version," use the repo's GitHub releases/tags for that.
