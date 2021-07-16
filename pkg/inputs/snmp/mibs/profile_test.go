package mibs

import (
	"testing"

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
		},
	}

	inputs := map[string]string{
		".1.3.6.1.4.1.2435.2.3.9.1":      "generic",
		".1.3.5.1.4.5":                   "nil",
		".1.3.6.1.4.1.2636.1.1.1.2.65":   "juniper",
		".1.3.6.1.4.1.29671.2.65":        "meraki",
		".1.3.6.1.4.1.2636.1.1.1.2.65.2": "generic",
		".1.3.6.1.4.1.29671":             "meraki",
		".1.3.6.1.4.1.318":               "apc",
		".1.3.6.1.4.1.8072.3.2.8":        "freebsd",
		".1.3.6.1.4.1.318.1.3.27":        "ups",
		".1.3.6.1.4.1.318.1.3.5.4":       "ups",
		".1.3.6.1.4.1.318.1.3.4.5":       "pdu",
		".1.3.6.1.4.1.318.1.3.4.6":       "pdu",
	}

	for sysid, vendor := range inputs {
		res := mdb.FindProfile(sysid)
		if vendor == "nil" {
			assert.Nil(t, res)
		} else {
			assert.Equal(t, vendor, res.Device.Vendor, "sysid: %s", sysid)
		}
	}
}

func TestPrune(t *testing.T) {
	mibs := map[string]*kt.Mib{
		".1.3.6.1.4.1.9.9.48.1.1.1.5": &kt.Mib{Tag: "MemoryUsed"},
	}
	prune(mibs)
	assert.Equal(t, 1, len(mibs))

	mibs = map[string]*kt.Mib{
		".1.3.6.1.4.1.9.9.48.1.1.1.5":     &kt.Mib{Tag: "MemoryUsed"},
		"1.3.6.1.4.1.9.9.305.1.1.1":       &kt.Mib{Tag: "CPU"},
		".1.3.6.1.4.1.9.9.109.1.1.1.1.7":  &kt.Mib{Tag: "CPU"},
		".1.3.6.1.4.1.9.9.109.1.1.1.1.13": &kt.Mib{Name: "cpmCPUMemoryFree"},
	}
	prune(mibs)
	assert.Equal(t, 3, len(mibs))
}
