package rule

import (
	"io/ioutil"
	"net"
	"os"
	"testing"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
	"github.com/stretchr/testify/assert"
)

func TestIsInternal(t *testing.T) {
	sut, err := NewRuleSet("", nil)
	assert.NoError(t, err)

	tests := map[string]bool{
		"1.2.3.4":            false,
		"10.2.1.2":           true,
		"2001:d00:abcd::/64": false,
		"2001:d01:abcd::/64": true,
		"10.2.1.3":           true,
		"172.16.0.1":         true,
		"10.19.38.78":        true,
		"foo":                false,
	}

	asns := map[string]uint32{
		"1.2.3.4":            3,
		"10.2.1.2":           3,
		"2001:d00:abcd::/64": 3,
		"2001:d01:abcd::/64": 64512,
		"10.2.1.3":           4294967294,
		"172.16.0.1":         3,
		"foo":                3,
		"10.19.38.78":        23,
	}

	for ip, isInternal := range tests {
		assert.Equal(t, isInternal, sut.IsInternal(net.ParseIP(ip), asns[ip]), ip)
	}
}

func TestGetService(t *testing.T) {
	l := lt.NewTestContextL(logger.NilContext, t)
	sut, err := NewRuleSet("", l)
	assert.NoError(t, err)

	tests := map[uint32]string{
		443:    "https",
		80:     "http",
		9092:   "XmlIpcRegSvc",
		100000: "nothing",
	}

	for port, service := range tests {
		app, ok := sut.GetService(nil, port, 6)
		if ok {
			assert.Equal(t, service, app, port)
		} else {
			assert.Equal(t, service, "nothing", port)
		}
	}

	// Now, run again but with some overrides
	content := []byte(`
applications:
  - ports: [9092,9093]
    name: kafka
`)

	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.FailNow()
	}
	if _, err := file.Write(content); err != nil {
		t.FailNow()
	}
	defer os.Remove(file.Name())

	sut, err = NewRuleSet(file.Name(), l)
	assert.NoError(t, err)

	tests[9092] = "kafka"

	for port, service := range tests {
		app, ok := sut.GetService(nil, port, 6)
		if ok {
			assert.Equal(t, service, app, port)
		} else {
			assert.Equal(t, service, "nothing", port)
		}
	}
}
