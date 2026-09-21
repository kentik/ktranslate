package version

import (
	"runtime"

	"github.com/kentik/ktranslate/pkg/eggs/version"
)

// versionStr, dateStr, and buildStr are overridden at link time, e.g.:
//   -ldflags "-X github.com/kentik/ktranslate/pkg/version.versionStr=v2.5.0 \
//             -X github.com/kentik/ktranslate/pkg/version.dateStr=2026-09-09 \
//             -X github.com/kentik/ktranslate/pkg/version.buildStr=ci-1234"
// See the Makefile's NETWORK_AGENT_VERSION/-DATE/-BUILD-derived LDFLAGS. buildStr is
// optional and empty by default -- only CI/nix-ci builds set it.
var (
	versionStr = "dev"
	dateStr    = "unknown"
	buildStr   = ""
)

var Version = version.VersionInfo{
	Version: versionStr,
	Date:    dateStr,
	Build:   buildStr,
	// Derived at runtime rather than baked in at build time -- correct even
	// for cross-compiled builds (e.g. `make arm` from a Darwin host).
	Platform: runtime.GOOS + "/" + runtime.GOARCH,
	Distro:   runtime.Version(),
}
