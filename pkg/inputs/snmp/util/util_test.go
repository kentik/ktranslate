package util

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gosnmp/gosnmp"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
)

type testWalker struct {
	results []gosnmp.SnmpPDU
	err     error
	dm      string
}

func (w testWalker) WalkAll(oid string) ([]gosnmp.SnmpPDU, error) {
	return w.results, w.err
}

func TestGetDeviceManufacturer(t *testing.T) {
	l := lt.NewTestContextL(logger.NilContext, t)
	for i, r := range []testWalker{
		// non-nil error
		{err: errors.New("Error"), dm: ""},
		// No results
		{results: []gosnmp.SnmpPDU{}, dm: ""},
		// type assertion failes
		{
			results: []gosnmp.SnmpPDU{{Value: "not a byte-slice"}},
			dm:      ""},
		// Empty string
		{
			results: []gosnmp.SnmpPDU{{Value: []byte{}}},
			dm:      ""},
		// convert base64
		{
			results: []gosnmp.SnmpPDU{{Value: []byte("SnVuaXBlciBOZXR3b3Jrcw==")}},
			dm:      "Juniper Networks"},
		// non-base64 passes through unscathed
		{
			results: []gosnmp.SnmpPDU{{Value: []byte("not base64")}},
			dm:      "not base64"},
		// 128 characters is fine
		{
			// 128 x's, encoded
			results: []gosnmp.SnmpPDU{{Value: []byte(strings.Repeat("eHh4", 43) + "eHg=")}},
			// 128 x's
			dm: strings.Repeat("x", 128)},
		// > 128 characters is truncated
		{
			// 129 x's, encoded
			results: []gosnmp.SnmpPDU{{Value: []byte(strings.Repeat("eHh4", 44))}},
			// Still 128 x's
			dm: strings.Repeat("x", 128)},
		// UTF8 in results.  Note literal `\u` in output.
		{
			results: []gosnmp.SnmpPDU{{Value: []byte("\u1234 utf8 \u2345")}},
			dm:      "\\u1234 utf8 \\u2345"},
		{
			results: []gosnmp.SnmpPDU{{Value: []byte(
				strings.Repeat("x", 127) + string('\u1234'))}},
			// Note that this result ends in "...xxx\u1234", and would be 133
			// bytes.
			dm: strings.Repeat("x", 127) + "\\u1234"},
		{
			results: []gosnmp.SnmpPDU{{Value: []byte(
				strings.Repeat("x", 127) + "\u1234x")}},
			dm: strings.Repeat("x", 127) + "\\u1234"},
	} {
		assert.Equal(t, r.dm, GetDeviceManufacturer(r, l), "failed %d", i)
	}
}

func TestGetIndex(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("4.5.6", GetIndex("1.2.3.4.5.6", ".2.3."))
	assert.Equal("4.5.6", GetIndex("1.2.3.4.5.6", "1.2.3."))
	assert.Equal(".4.5.6", GetIndex("1.2.3.4.5.6", "1.2.3"))
	assert.Equal(".4.5.6", GetIndex(".1.2.3.4.5.6", "1.2.3"))
}

func TestHandlePowerset(t *testing.T) {
	assert := assert.New(t)
	l := lt.NewTestContextL(logger.NilContext, t)

	tests := map[string][]byte{
		"Informational: On Line,Serial Communication Established,Powered On": []byte("0001010000000000001000000000000000000000000000000000000000000000"),
		"Informational: On Line,Powered On":                                  []byte("0001000000000000001000000000000000000000000000000000000000000000"),
		"Disaster: No Batteries Attached":                                    []byte("0001010000000000001000000000000010000000000000000000000000000000"),
		"Warning: Self Test In Progress":                                     []byte("0101010000000000001000000000110000000000000000000000000000000000"),
		"High: On Battery,Low Battery / On Battery":                          []byte("0101010000000000001000000000010000000000000000000000000000000000"), // No self test.
	}

	levels := map[string]int64{
		"Informational": 1,
		"Disaster":      4,
		"Warning":       2,
		"High":          3,
	}

	for expt, in := range tests {
		pdu := gosnmp.SnmpPDU{Value: in}
		val, str := GetFromConv(pdu, "powerset_status", l)
		assert.Equal(expt, str, "%s", in)
		pts := strings.SplitN(str, ":", 2)
		assert.Equal(levels[pts[0]], val, "%s", in)
	}
}

func TestHWAddr(t *testing.T) {
	assert := assert.New(t)
	l := lt.NewTestContextL(logger.NilContext, t)

	tests := map[string][]byte{
		"39:30:3a:36:31:3a:61:65:3a:66:62:3a:63:32:3a:31:39": []byte("90:61:ae:fb:c2:19"),
	}

	for expt, in := range tests {
		pdu := gosnmp.SnmpPDU{Value: in}
		_, str := GetFromConv(pdu, "hwaddr", l)
		assert.Equal(expt, str, "%s", in)
	}
}

func TestHexToInt(t *testing.T) {
	assert := assert.New(t)
	l := lt.NewTestContextL(logger.NilContext, t)

	tests := map[string][]byte{
		"7090466131453292601": []byte("9061aefbc219"),
	}

	for expt, in := range tests {
		pdu := gosnmp.SnmpPDU{Value: in}
		_, str := GetFromConv(pdu, "hextoint:LittleEndian:uint64", l)
		assert.Equal(expt, str, "%s", in)
	}
}
