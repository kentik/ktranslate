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
		deviceMetadataMibs: map[string]*kt.Mib{
			"1.1.2": &kt.Mib{Name: "deviceFoo", Tag: "my_device_tag", Table: "myDevice"},
		},
		interfaceMetadataMibs: map[string]*kt.Mib{
			"1.1.1": &kt.Mib{Name: "ifName", Tag: "interface_name", Table: "if"},
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
		res, ok := p.lookupMib(key, mib.Table == "if")
		if key == "" {
			assert.Equal(t, false, ok)
		} else {
			assert.NotNil(t, res)
			assert.Equal(t, mib.Table, res.Table)
		}
	}
}

func TestToFlows(t *testing.T) {
	l := lt.NewTestContextL(logger.NilContext, t)
	conf := &kt.SnmpDeviceConfig{
		Provider: "foo",
		UserTags: map[string]string{
			"foo": "$foo",
			"aaa": "$SysContact",
		},
	}
	conf.InitUserTags("service")

	p := &Poller{
		log:   l,
		conf:  conf,
		gconf: &kt.SnmpGlobalConfig{},
	}

	input := kt.DeviceData{
		Manufacturer: "man",
		DeviceMetricsMetadata: &kt.DeviceMetricsMetadata{
			SysContact: "ddd",
			Customs: map[string]string{
				"foo": "",
				"bar": "",
			},
		},
	}
	res, err := p.toFlows(&input)
	assert.NotNil(t, res)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, "ddd", res[0].CustomStr["tags.aaa"])
	assert.Equal(t, "", res[0].CustomStr["tags.foo"]) // Empty string match for tag here.
}

func TestToFlowsCache(t *testing.T) {
	l := lt.NewTestContextL(logger.NilContext, t)
	conf := &kt.SnmpDeviceConfig{
		Provider: "foo",
		UserTags: map[string]string{},
	}
	conf.InitUserTags("service")

	p := &Poller{
		log:   l,
		conf:  conf,
		gconf: &kt.SnmpGlobalConfig{},
	}

	input := kt.DeviceData{
		DeviceMetricsMetadata: &kt.DeviceMetricsMetadata{
			Customs: map[string]string{
				"foo": "aaa",
				"bar": "bbb",
			},
		},
	}
	res, err := p.toFlows(&input)
	assert.NotNil(t, res)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, "aaa", res[0].CustomStr["foo"])
	assert.Equal(t, "bbb", res[0].CustomStr["bar"])

	// Now without
	input = kt.DeviceData{
		DeviceMetricsMetadata: &kt.DeviceMetricsMetadata{
			Customs: map[string]string{
				"foo": "aaa",
			},
		},
	}
	res, err = p.toFlows(&input)
	assert.NotNil(t, res)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, "aaa", res[0].CustomStr["foo"])
	assert.Equal(t, "", res[0].CustomStr["bar"])

	// And again, but with caching turned on.
	p = &Poller{
		log:   l,
		conf:  conf,
		gconf: &kt.SnmpGlobalConfig{SaveCache: true},
	}

	input = kt.DeviceData{
		DeviceMetricsMetadata: &kt.DeviceMetricsMetadata{
			Customs: map[string]string{
				"foo": "aaa",
				"bar": "bbb",
			},
		},
	}
	res, err = p.toFlows(&input)
	assert.NotNil(t, res)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, "aaa", res[0].CustomStr["foo"])
	assert.Equal(t, "bbb", res[0].CustomStr["bar"])

	// Now without
	input = kt.DeviceData{
		DeviceMetricsMetadata: &kt.DeviceMetricsMetadata{
			Customs: map[string]string{
				"foo": "aaa",
			},
		},
	}
	res, err = p.toFlows(&input)
	assert.NotNil(t, res)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, "aaa", res[0].CustomStr["foo"])
	assert.Equal(t, "bbb", res[0].CustomStr["bar"]) // bbb is picked up from cache.
}
