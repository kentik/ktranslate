package mibs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
	"github.com/kentik/ktranslate/pkg/kt"
)

func TestFindProfile(t *testing.T) {
	l := lt.NewTestContextL(logger.NilContext, t)
	mdb := &MibDB{
		db:  nil,
		log: l,
		profiles: map[string]*Profile{
			"1.3.6.1.4.1.2636.1.1.1.2.65": &Profile{Device: Device{Vendor: "juniper"}},
			"1.3.6.1.4.1.29671.*":         &Profile{Device: Device{Vendor: "meraki"}},
			"1.3.6.1.4.*":                 &Profile{Device: Device{Vendor: "generic"}},
			"1.3.6.1.4.1.318.*":           &Profile{Device: Device{Vendor: "apc"}},
			"1.3.6.1.4.1.8072.3.2.8":      &Profile{Device: Device{Vendor: "freebsd"}},
			"1.3.6.1.4.1.318.1.3.4.*":     &Profile{Device: Device{Vendor: "pdu"}},
			"1.3.6.1.4.1.318.1.3.*":       &Profile{Device: Device{Vendor: "ups"}},
			"1.3.6.1.4.1.8072.3.2.*":      &Profile{Matches: map[string]string{"^barracuda": "barracuda.yml"}, Device: Device{Vendor: "net-bsd"}},
			"1.3.6.1.4.1.1588.2.1.1.*":    &Profile{From: "barracuda.yml", Device: Device{Vendor: "barracuda"}},
			"1.3.6.1.4.1.1588.2.1.22.*":   &Profile{From: "foo.yml", Device: Device{Vendor: "foo"}},
			"1.3.6.1.4.1.2.3.4":           &Profile{From: "direct", Device: Device{Vendor: "direct"}},
			"1.3.6.1.4.1.8072.3.66.*":     &Profile{MatchesList: []Match{Match{Regex: "^barracuda", Target: "barracuda.yml"}}, Device: Device{Vendor: "net-bsd"}},
			"1.3.6.1.4.1.8072.3.77.*":     &Profile{MatchesList: []Match{Match{Regex: "^foocuda", Target: "foo.yml"}}, Matches: map[string]string{"^barracuda": "barracuda.yml"}, Device: Device{Vendor: "net-bsd"}},
		},
	}

	inputs := []struct {
		Sysoid   string
		Sysdesc  string
		Profile  string
		Expected string
	}{
		{
			Sysoid:   ".1.3.6.1.4.1.2435.2.3.9.1",
			Expected: "generic",
		},
		{
			Sysoid:   ".1.3.5.1.4.5",
			Expected: "nil",
		},
		{
			Sysoid:   ".1.3.6.1.4.1.2636.1.1.1.2.65",
			Expected: "juniper",
		},
		{
			Sysoid:   ".1.3.6.1.4.1.29671.2.65",
			Expected: "meraki",
		},
		{
			Sysoid:   ".1.3.6.1.4.1.2636.1.1.1.2.65.2",
			Expected: "generic",
		},
		{
			Sysoid:   ".1.3.6.1.4.1.29671",
			Expected: "meraki",
		},
		{
			Sysoid:   ".1.3.6.1.4.1.318",
			Expected: "apc",
		},
		{
			Sysoid:   ".1.3.6.1.4.1.8072.3.2.8",
			Expected: "freebsd",
		},
		{
			Sysoid:   ".1.3.6.1.4.1.318.1.3.27",
			Expected: "ups",
		},
		{
			Sysoid:   ".1.3.6.1.4.1.318.1.3.5.4",
			Expected: "ups",
		},
		{
			Sysoid:   ".1.3.6.1.4.1.318.1.3.4.5",
			Expected: "pdu",
		},
		{
			Sysoid:   ".1.3.6.1.4.1.318.1.3.4.6",
			Expected: "pdu",
		},
		{
			Sysoid:   ".1.3.6.1.4.1.318.1.3.4.6",
			Expected: "pdu",
		},
		{
			Sysoid:   ".1.3.6.1.4.1.8072.3.2.10",
			Expected: "barracuda",
			Sysdesc:  "Barracuda Email Security Gateway",
		},
		{
			Sysoid:   ".1.3.6.1.4.1.8072.3.2.10",
			Expected: "barracuda",
			Sysdesc:  "Barracuda Email Security Gateway",
		},
		{
			Sysoid:   "",
			Expected: "direct",
			Sysdesc:  "Barracuda Email Security Gateway",
			Profile:  "!direct",
		},
		{
			Sysoid:   ".1.3.6.1.4.1.8072.3.66.10",
			Expected: "barracuda",
			Sysdesc:  "Barracuda Email Security Gateway",
		},
		{
			Sysoid:   ".1.3.6.1.4.1.8072.3.77.10",
			Expected: "foo",
			Sysdesc:  "Foocuda Email Security Gateway",
		},
	}

	for _, in := range inputs {
		res := mdb.FindProfile(in.Sysoid, in.Sysdesc, in.Profile)
		if in.Expected == "nil" {
			assert.Nil(t, res)
		} else {
			assert.Equal(t, in.Expected, res.Device.Vendor, "sysid: %s", in.Sysoid)
		}
	}
}

func TestGetMibAndOid(t *testing.T) {
	fooPro := Profile{Metrics: []MIB{MIB{Table: OID{Oid: ".1.2.3.3", Name: "foo"}}}}

	tests := map[string][]*Profile{
		"foo": []*Profile{
			&Profile{Metrics: []MIB{MIB{Table: OID{Oid: ".1.2.3.3", Name: "foo"}}}},
			&fooPro,
			&fooPro, // Caching.
			&fooPro,
		},
		"computed": []*Profile{
			nil, // Nil case.
			&Profile{Metrics: []MIB{MIB{Table: OID{Oid: ".1.2.3.3.1", Name: "foo"}}}}}, // Not found case.
	}

	for expected, profiles := range tests {
		for _, profile := range profiles {
			mib, _ := profile.GetMibAndOid()
			assert.Equal(t, expected, mib)
		}
	}
}

func TestPrune(t *testing.T) {
	mibs := map[string]*kt.Mib{
		".1.3.6.1.4.1.9.9.48.1.1.1.5": &kt.Mib{Tag: "MemoryUsed"},
	}
	prune(mibs, 0, 0)
	assert.Equal(t, 1, len(mibs))

	mibs = map[string]*kt.Mib{
		".1.3.6.1.4.1.9.9.48.1.1.1.5":     &kt.Mib{Tag: "MemoryUsed"},
		".1.3.6.1.4.1.9.9.305.1.1.1":      &kt.Mib{Tag: "CPU"},
		".1.3.6.1.4.1.9.9.109.1.1.1.1.7":  &kt.Mib{Tag: "CPU"},
		".1.3.6.1.4.1.9.9.109.1.1.1.1.13": &kt.Mib{Name: "cpmCPUMemoryFree"},
	}
	prune(mibs, 0, 0)
	assert.Equal(t, 3, len(mibs))
	assert.Nil(t, mibs[".1.3.6.1.4.1.9.9.305.1.1.1"], len(mibs))
	assert.NotNil(t, mibs[".1.3.6.1.4.1.9.9.109.1.1.1.1.7"], len(mibs))

	mibs = map[string]*kt.Mib{
		".1.3.6.1.4.1.9.9.48.1.1.1.5":     &kt.Mib{Tag: "MemoryUsed"},
		".1.3.6.1.4.1.9.9.305.1.1.1":      &kt.Mib{Tag: "CPU"},
		".1.3.6.1.4.1.9.9.109.1.1.1.1.7":  &kt.Mib{Tag: "CPU", FromExtended: true},
		".1.3.6.1.4.1.9.9.109.1.1.1.1.13": &kt.Mib{Name: "cpmCPUMemoryFree"},
	}
	prune(mibs, 0, 0)
	assert.Equal(t, 3, len(mibs))
	assert.NotNil(t, mibs[".1.3.6.1.4.1.9.9.305.1.1.1"], len(mibs))
	assert.Nil(t, mibs[".1.3.6.1.4.1.9.9.109.1.1.1.1.7"], len(mibs))
	assert.Equal(t, mibs[".1.3.6.1.4.1.9.9.305.1.1.1"].PollDur, time.Duration(0)*time.Second)

	prune(mibs, 60, 30)
	assert.Equal(t, 3, len(mibs))
	assert.NotNil(t, mibs[".1.3.6.1.4.1.9.9.305.1.1.1"], len(mibs))
	assert.Equal(t, mibs[".1.3.6.1.4.1.9.9.305.1.1.1"].PollDur, time.Duration(60-kt.PollAdjustTime)*time.Second)
}

func TestProfileName(t *testing.T) {
	tests := map[string][]*Profile{
		"foo-test": []*Profile{&Profile{From: "foo-test.yml"}},
		"snmp":     []*Profile{nil, &Profile{}},
		"noone":    []*Profile{&Profile{From: "noone"}},
		"override": []*Profile{&Profile{From: "lalalalala"}},
	}
	overrides := map[string]string{
		"override": "override",
	}

	for expected, profiles := range tests {
		for _, profile := range profiles {
			res := profile.GetProfileName(overrides[expected])
			assert.Equal(t, expected, res)
		}
	}
}

func TestGetTableName(t *testing.T) {
	tests := map[string]OID{
		"physicalDisk":      OID{Name: "physicalDiskTable"},
		"systemSlot":        OID{Name: "systemSlotTable"},
		"la":                OID{Name: "laTable"},
		"diskIO":            OID{Name: "diskIOTable"},
		"":                  OID{Name: ""},
		"diskio":            OID{Name: "diskiotable"},
		"if":                OID{Name: "ifXTable"},
		"jnxOperatingEntry": OID{Name: "jnxOperatingEntry"},
	}

	for expected, oid := range tests {
		res := oid.GetTableName()
		assert.Equal(t, expected, res)
	}
}
