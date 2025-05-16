package snmp

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
	"github.com/kentik/ktranslate/pkg/kt"
)

func TestGetIP(t *testing.T) {
	assert := assert.New(t)
	l := lt.NewTestContextL(logger.NilContext, t)

	test := []byte(`{
        "id": 1639742,
        "family": {
          "value": 4,
          "label": "IPv4"
        },
        "address": "10.249.157.132/32"
      }`)

	conf := kt.NetboxConfig{NetboxIP: "primary"}
	nbRes := NBIP{}
	err := json.Unmarshal(test, &nbRes)
	assert.Nil(err)
	ipv, err := getIP(NBResult{PrimaryIp: &nbRes}, &conf, l)
	assert.Nil(err)
	assert.Equal("10.249.157.132", ipv.Addr().String())
}
