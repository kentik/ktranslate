package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gosnmp/gosnmp"
	"github.com/kentik/ktranslate/pkg/inputs/snmp/mibs"
	"github.com/kentik/ktranslate/pkg/kt"
)

func TestCheckCondition(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]bool{
		"":         true,
		"name1=2":  false,
		"name1=66": true,
	}

	resultSet := []wrapper{
		wrapper{
			oid:      "1.2",
			mib:      &kt.Mib{Tag: "name1"},
			variable: gosnmp.SnmpPDU{Value: 66, Name: "1.2.5"},
		},
		wrapper{
			oid:      "1.1",
			mib:      &kt.Mib{Tag: "name2"},
			variable: gosnmp.SnmpPDU{Value: 3, Name: "1.1.5"},
		},
	}

	for in, expt := range tests {
		oid := mibs.OID{Condition: in}
		mib := &kt.Mib{
			Condition: oid.GetCondition(),
			Tag:       "name2",
		}
		w := wrapper{
			variable: gosnmp.SnmpPDU{Value: 3, Name: "1.1.5"},
			mib:      mib,
			oid:      "1.1",
		}

		res := w.checkCondition(".5", resultSet)
		assert.Equal(expt, res, "%s", in)
	}
}
