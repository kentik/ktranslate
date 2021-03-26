package mibs

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
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
		},
	}

	inputs := map[string]string{
		".1.3.6.1.4.1.2435.2.3.9.1":      "generic",
		".1.3.5.1.4.5":                   "nil",
		".1.3.6.1.4.1.2636.1.1.1.2.65":   "juniper",
		".1.3.6.1.4.1.29671.2.65":        "meraki",
		".1.3.6.1.4.1.2636.1.1.1.2.65.2": "generic",
		".1.3.6.1.4.1.29671":             "meraki",
		".1.3.6.1.4.1.318.1.3.4.6":       "apc",
		".1.3.6.1.4.1.318":               "apc",
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
