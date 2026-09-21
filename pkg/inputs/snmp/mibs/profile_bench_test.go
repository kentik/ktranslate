package mibs

import (
	"fmt"
	"testing"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
)

// buildBenchMibDB returns a MibDB with numProfiles sibling OID-wildcard entries (so
// FindProfile's recursive wildcard walk, profile.go:289-296, has a realistically
// populated map to search) plus one profile at matchOid+".*" carrying numMatches
// regex entries in MatchesList, so checkMatch (profile.go:302-337) does real work.
// None of the generated patterns match sysdesc -- the common case in production,
// where most configured match rules don't apply to a given device, but every one of
// them still gets recompiled from scratch on every call (profile.go:305,322) rather
// than once at load time. See docs/DISCOVERY_PERFORMANCE_PLAN.md, finding A36.
func buildBenchMibDB(b *testing.B, numProfiles, numMatches int) (mdb *MibDB, sysid, sysdesc string) {
	b.Helper()
	l := lt.NewBenchContextL(logger.NilContext, b)
	profiles := make(map[string]*Profile, numProfiles+1)
	for i := range numProfiles {
		profiles[fmt.Sprintf("1.3.6.1.4.1.%d.*", i+1)] = &Profile{Device: Device{Vendor: fmt.Sprintf("vendor-%d", i)}}
	}

	matches := make([]Match, numMatches)
	for i := range numMatches {
		matches[i] = Match{Regex: fmt.Sprintf("^does-not-match-pattern-%d", i), Target: "unused.yml"}
	}
	const matchOid = "1.3.6.1.4.1.99999.1.2.3"
	profiles[matchOid+".*"] = &Profile{
		From:        "matched.yml",
		MatchesList: matches,
		Device:      Device{Vendor: "matched"},
	}

	mdb = &MibDB{log: l, profiles: profiles}
	sysid = matchOid + ".7" // resolves to the matchOid wildcard entry via the recursive walk
	sysdesc = "Some Real Device Description, No Match Here"
	return
}

// BenchmarkFindProfile_MatchesList measures FindProfile's cost as the number of
// configured MatchesList regexes on the matched profile grows. This is the call path
// exercised on every discovery hit (disco.go:350) and, more importantly, on every
// config parse -- which happens up to three times per discovery cycle across the
// whole device fleet (snmp.go:205,457-468; disco.go:44,411; snmp.go:161). See
// docs/DISCOVERY_PERFORMANCE_PLAN.md findings A15, A24, A30, A36.
func BenchmarkFindProfile_MatchesList(b *testing.B) {
	for _, numMatches := range []int{0, 5, 20, 50} {
		b.Run(fmt.Sprintf("patterns=%d", numMatches), func(b *testing.B) {
			mdb, sysid, sysdesc := buildBenchMibDB(b, 500, numMatches)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				mdb.FindProfile(sysid, sysdesc, "")
			}
		})
	}
}

// BenchmarkFindProfile_FleetParse approximates the per-parse cost identified in
// snmp.go:457-468 -- a loop calling FindProfile once per device in the fleet -- at
// fleet sizes relevant to the "supporting 5k devices" goal. This isolates the same
// checkMatch cost as BenchmarkFindProfile_MatchesList but reports it as the unit that
// actually matters: total time for one full-fleet parse pass, and (via -benchmem)
// allocations per pass.
func BenchmarkFindProfile_FleetParse(b *testing.B) {
	for _, fleetSize := range []int{100, 1000, 5000} {
		b.Run(fmt.Sprintf("devices=%d", fleetSize), func(b *testing.B) {
			mdb, sysid, sysdesc := buildBenchMibDB(b, 500, 10) // 10 patterns: a realistic mid-size profile
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for range fleetSize {
					mdb.FindProfile(sysid, sysdesc, "")
				}
			}
		})
	}
}
