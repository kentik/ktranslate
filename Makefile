MODULE := github.com/kentik/ktranslate

# NETWORK_AGENT_VERSION: defaults to the checked-in VERSION file (repo root) -- the same
# semver source of truth flake.nix and nix/network-agent.nix use, so all build paths agree
# unless a Docker --build-arg or CI explicitly overrides it. Falls back to git describe
# only if VERSION itself is ever missing. See `check-version-env-var` below for the one
# place this name must also match.
NETWORK_AGENT_VERSION ?= $(shell cat VERSION 2>/dev/null || git describe --tags --always --dirty 2>/dev/null || echo dev)

# NETWORK_AGENT_DATE: this commit's own timestamp, not wall-clock -- building the same
# commit twice stamps the same date both times (matches nix/network-agent.nix's use of
# self.lastModifiedDate). Falls back to wall-clock only when there's no git history to
# ask at all (e.g. an extracted source tarball with no .git).
NETWORK_AGENT_DATE ?= $(shell git log -1 --format=%cI 2>/dev/null || date -u +%Y-%m-%dT%H:%M:%SZ)

# NETWORK_AGENT_BUILD: optional, empty by default. CI sets this to identify which run
# produced a given artifact (see ci-build.yml) -- unlike VERSION/DATE, it's not meant to
# be the same across builds of the same commit.
NETWORK_AGENT_BUILD ?=

LDFLAGS := -X '$(MODULE)/pkg/version.versionStr=$(NETWORK_AGENT_VERSION)' -X '$(MODULE)/pkg/version.dateStr=$(NETWORK_AGENT_DATE)' -X '$(MODULE)/pkg/version.buildStr=$(NETWORK_AGENT_BUILD)'

.PHONY: all
all:
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/ktranslate ./cmd/ktranslate

.PHONY: windows
windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/ktranslate.exe ./cmd/ktranslate

.PHONY: arm
arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-s -w $(LDFLAGS)" -o bin/ktranslate ./cmd/ktranslate

.PHONY: print-version-env-var
print-version-env-var:
	@echo NETWORK_AGENT_VERSION

.PHONY: check-version-env-var
check-version-env-var:
	@grep -qx "ARG $$(make -s print-version-env-var)" Dockerfile || \
	  { echo "Dockerfile's ARG doesn't match Makefile's version env var" >&2; exit 1; }

.PHONY: test
test: check-version-env-var
	go test ./cmd/... ./pkg/...

.PHONY: bench
bench:
	go test -bench=. ./cmd/... ./pkg/...

.PHONY: ktranslate
ktranslate:
	go install ./cmd/ktranslate

.PHONY: clean
clean:
	rm -f bin/ktranslate

.PHONY: generate
generate:
	go generate ./...

.PHONY: install
install:
	mkdir -p $(DESTDIR)/usr/local/bin
	install -m 0755 bin/ktranslate $(DESTDIR)/usr/local/bin

.PHONY: docker
docker: all
	docker pull ubuntu:20.04
	docker build -t ktranslate:v2 -f Dockerfile .
