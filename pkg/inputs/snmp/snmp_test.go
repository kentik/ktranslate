package snmp

import (
	"fmt"
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kentik/ktranslate/pkg/kt"
)

func TestMatchesPrefix(t *testing.T) {
	assert := assert.New(t)

	tests := map[string][]string{
		"tagName":               []string{"tagName", "true"},
		"provider:tagName":      []string{"provider:tagName", "true"},
		"provider:foo:tagName":  []string{"tagName", "true"},
		"provider:bar:tagName":  []string{"", "false"},
		"provider:foo:tagName1": []string{"tagName1", "true"},
	}

	provider := kt.Provider("foo")

	for in, expt := range tests {
		newTag, res := matchesPrefix(in, provider)
		assert.Equal(expt[0], newTag, "%s <-> %s", in, provider)
		assert.Equal(expt[1], fmt.Sprintf("%v", res), "%s <-> %s", in, provider)
	}
}

func TestSetTagsMatch(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]kt.SnmpConfig{
		"one": kt.SnmpConfig{
			Global: &kt.SnmpGlobalConfig{
				UserTags: map[string]string{
					"tag": "global",
				},
				MatchAttr: map[string]string{
					"match": "global",
				},
				ProviderMap: map[string]kt.ProviderMap{
					"foo": kt.ProviderMap{
						UserTags: map[string]string{
							"tag": "provider",
						},
						MatchAttr: map[string]string{
							"match": "provider",
						},
					},
					"bar": kt.ProviderMap{
						UserTags: map[string]string{
							"tag": "provider",
						},
						MatchAttr: map[string]string{
							"match": "provider",
						},
					},
				},
			},
			Devices: map[string]*kt.SnmpDeviceConfig{
				"device": &kt.SnmpDeviceConfig{
					Provider: "foo",
					UserTags: map[string]string{
						"tag": "device", // This should be device tag because its set at device level.
					},
					MatchAttr: map[string]string{
						"match": "device",
					},
				},
				"provider": &kt.SnmpDeviceConfig{
					Provider: "foo",
					UserTags: map[string]string{
						"tagA": "device", // This should fall back to provider level.
					},
					MatchAttr: map[string]string{
						"matchA": "device",
					},
				},
				"global": &kt.SnmpDeviceConfig{
					Provider: "fooA",
					UserTags: map[string]string{
						"tagA": "device", // This should fall back to global level because niether provider or device set.
					},
					MatchAttr: map[string]string{
						"matchA": "device",
					},
				},
			},
		},
		"two": kt.SnmpConfig{ // No provider, just gobal and device.
			Global: &kt.SnmpGlobalConfig{
				UserTags: map[string]string{
					"tag": "global",
				},
				MatchAttr: map[string]string{
					"match": "global",
				},
			},
			Devices: map[string]*kt.SnmpDeviceConfig{
				"device": &kt.SnmpDeviceConfig{
					Provider: "foo",
					UserTags: map[string]string{
						"tag": "device", // This should be device tag because its set at device level.
					},
					MatchAttr: map[string]string{
						"match": "device",
					},
				},
				"global": &kt.SnmpDeviceConfig{
					Provider: "fooA",
					UserTags: map[string]string{
						"tagA": "device", // This should fall back to global level because niether provider or device set.
					},
					MatchAttr: map[string]string{
						"matchA": "device",
					},
				},
			},
		},
		"three": kt.SnmpConfig{ // No provider, just gobal and device.
			Global: &kt.SnmpGlobalConfig{},
			Devices: map[string]*kt.SnmpDeviceConfig{
				"device": &kt.SnmpDeviceConfig{
					Provider: "foo",
					UserTags: map[string]string{
						"tag": "device", // This should be device tag because its set at device level.
					},
					MatchAttr: map[string]string{
						"match": "device",
					},
				},
				"": &kt.SnmpDeviceConfig{
					Provider: "fooA",
					UserTags: map[string]string{
						"tagA": "device", // This should fall back to global level because niether provider or device set.
					},
					MatchAttr: map[string]string{
						"matchA": "device",
					},
				},
			},
		},
	}

	for test, ms := range tests {
		for p, m := range ms.Global.ProviderMap {
			m.Init(p, &ms) // Set up any provider based user and match tags here.
		}
		for k, v := range ms.Global.UserTags {
			for _, device := range ms.Devices {
				if device.UserTags == nil {
					device.UserTags = map[string]string{}
				}
				if _, ok := device.UserTags[k]; !ok {
					device.UserTags[k] = v
				}
			}
		}
		for k, v := range ms.Global.MatchAttr {
			for _, device := range ms.Devices {
				if device.MatchAttr == nil {
					device.MatchAttr = map[string]string{}
				}
				if _, ok := device.MatchAttr[k]; !ok {
					device.MatchAttr[k] = v
				}
			}
		}

		for expt, device := range ms.Devices {
			setDeviceTagsAndMatch(device)
			assert.Equal(expt, device.UserTags["tag"], "%s -> %s %v", test, device.Provider, device.UserTags)
			assert.Equal(expt, device.MatchAttr["match"], "%s -> %s %v", test, device.Provider, device.MatchAttr)
		}
	}
}

func TestIgnoreList(t *testing.T) {
	assert := assert.New(t)

	ignores := []string{"10.2.3.1", "10.2.10.0/24"}

	tests := map[string]bool{
		"10.2.3.1":   true,
		"10.2.10.10": true,
		"10.2.11.10": false,
	}

	ignoreMap := map[string]bool{}
	ignoreList := []netip.Prefix{}
	for _, ip := range ignores {
		if ipr, err := netip.ParsePrefix(ip); err != nil {
			ignoreMap[ip] = true
		} else {
			ignoreList = append(ignoreList, ipr)
		}
	}

	for tst, res := range tests {
		assert.Equal(res, checkIfIgnored(tst, ignoreMap, ignoreList), tst)
	}
}
