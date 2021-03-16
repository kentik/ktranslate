package mibs

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
	"github.com/kentik/ktranslate/pkg/kt"
)

func TestCheckForProvider(t *testing.T) {
	l := lt.NewTestContextL(logger.NilContext, t)
	inputs := map[string]kt.Provider{
		"net-printer":   kt.ProviderIOT,
		"hpEtherSwitch": kt.ProviderSwitch,
		"Firewall":      kt.ProviderFirewall,
		"Router":        kt.ProviderRouter,
		"dsfsdfaa":      "",
	}

	mdb := &MibDB{
		log: l,
	}

	for input, prov := range inputs {
		res, ok := mdb.checkForProvider(input, "", "")
		assert.Equal(t, prov, res, "input: %s", input)
		assert.Equal(t, res != "", ok, "input: %s", input)
	}
}

func TestFullProvider(t *testing.T) {
	l := lt.NewTestContextL(logger.NilContext, t)
	mdb, err := NewMibDB("mibs.db", "", "", l)
	assert.NoError(t, err)
	defer mdb.Close()

	info := map[string][]string{
		".1.3.6.1.4.1.11.2.3.7.11.162.8": []string{"", ""},
		".1.3.6.1.4.1.2435.2.3.9.1":      []string{"", ""},
		".1.3.6.1.4.1.9.1.2494":          []string{"generic-router.yaml", "Cisco IOS Software [Everest], Catalyst L3 Switch Software (CAT9K_IOSXE)"},
		".1.3.6.1.4.1.9.1.1639":          []string{"generic-router.yaml", "Cisco IOS XR Software (Cisco ASR9K Series),  Version 6.4.2[Default]\r\nCopyright"},
		".1.3.6.1.4.1.9.1.449":           []string{"cisco-catalyst.yaml", "Cisco IOS Software, s72033_rp Software (s72033_rp-ADVIPSERVICESK9_WAN-M)"},
	}

	inputs := map[string]kt.Provider{
		".1.3.6.1.4.1.11.2.3.7.11.162.8": kt.ProviderSwitch,
		".1.3.6.1.4.1.2435.2.3.9.1":      kt.ProviderIOT,
		".1.3.6.1.4.1.9.1.2494":          kt.ProviderSwitch,
		".1.3.6.1.4.1.9.1.1639":          kt.ProviderRouter,
		".1.3.6.1.4.1.9.1.449":           kt.ProviderSwitch,
	}

	for input, prov := range inputs {
		_, res, _, err := mdb.GetForOidRecur(input, info[input][0], info[input][1])
		assert.NoError(t, err)
		assert.Equal(t, prov, res, "input: %s", input)
	}
}
