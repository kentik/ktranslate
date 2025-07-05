package snmp

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
	"github.com/kentik/ktranslate/pkg/kt"
)

type ntest struct {
	nbresult []byte
	target   string
	res      string
}

func TestGetIP(t *testing.T) {
	assert := assert.New(t)
	l := lt.NewTestContextL(logger.NilContext, t)

	tests := []ntest{
		ntest{
			res:    "10.249.157.132",
			target: "primary",
			nbresult: []byte(`{"primary_ip":{
        "id": 1639742,
        "family": {
          "value": 4,
          "label": "IPv4"
        },
        "address": "10.249.157.132/32"
      }}`)},
		ntest{
			res:    "10.249.157.135",
			target: "oob",
			nbresult: []byte(`{"oob_ip":{
        "id": 1639742,
        "family": {
          "value": 4,
          "label": "IPv4"
        },
        "address": "10.249.157.135/32"
      }}`)},
		ntest{
			res:    "10.249.157.135",
			target: "oob,primary_ip4",
			nbresult: []byte(`{"oob_ip":null,"primary_ip4":{
        "id": 1639742,
        "family": {
          "value": 4,
          "label": "IPv4"
        },
        "address": "10.249.157.135/32"
      }}`)},
		ntest{
			res:    "dead::beef",
			target: "oob,primary_ip4,primary_ip6",
			nbresult: []byte(`{"oob_ip":null,"primary_ip4":null,"primary_ip6":{
        "id": 1639742,
        "family": {
          "value": 6,
          "label": "IPv6"
        },
        "address": "dead::beef/64"
      }}`)},
	}

	for _, test := range tests {
		conf := kt.NetboxConfig{NetboxIP: test.target}
		nbRes := NBResult{}
		err := json.Unmarshal(test.nbresult, &nbRes)
		assert.Nil(err)
		ipv, err := getIP(nbRes, &conf, l)
		assert.Nil(err)
		assert.Equal(test.res, ipv.Addr().String())
	}
}
