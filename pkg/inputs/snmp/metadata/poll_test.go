package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
	"github.com/kentik/ktranslate/pkg/kt"
)

func TestGetDeviceManufacturer(t *testing.T) {
	l := lt.NewTestContextL(logger.NilContext, t)

	p := &Poller{
		deviceMetadataMibs: map[string]*kt.Mib{},
		interfaceMetadataMibs: map[string]*kt.Mib{
			"1.1.1": &kt.Mib{Name: "ifName", Tag: "interface_name", Table: "if"},
			"1.1.2": &kt.Mib{Name: "deviceFoo", Tag: "my_device_tag", Table: "myDevice"},
		},
		log: l,
	}

	tests := map[string]*kt.Mib{
		"ifName":         &kt.Mib{Table: "if"},
		"interface_name": &kt.Mib{Table: "if"},
		"deviceFoo":      &kt.Mib{Table: "myDevice"},
		"my_device_tag":  &kt.Mib{Table: "myDevice"},
	}

	for key, mib := range tests {
		res, ok := p.lookupMib(key)
		if key == "" {
			assert.Equal(t, false, ok)
		} else {
			assert.NotNil(t, res)
			assert.Equal(t, mib.Table, res.Table)
		}
	}
}
