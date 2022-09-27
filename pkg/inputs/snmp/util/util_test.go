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
		val, str, _ := GetFromConv(pdu, "powerset_status", l)
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
		_, str, _ := GetFromConv(pdu, "hwaddr", l)
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
		_, str, _ := GetFromConv(pdu, "hextoint:LittleEndian:uint64", l)
		assert.Equal(expt, str, "%s", in)
	}
}

func TestHexToIP(t *testing.T) {
	assert := assert.New(t)
	l := lt.NewTestContextL(logger.NilContext, t)

	tests := map[string][]byte{
		"10.0.100.10":          []byte{0x00, 0x0a, 0x00, 0x00, 0x00, 0x64, 0x00, 0x0a},
		"10.0.100.11":          []byte{0x0a, 0x00, 0x64, 0x0b},
		"2001:504:0:2::6169:1": []byte{0x20, 0x01, 0x05, 0x04, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x61, 0x69, 0x00, 0x01},
	}

	for expt, in := range tests {
		pdu := gosnmp.SnmpPDU{Value: in}
		_, str, _ := GetFromConv(pdu, CONV_HEXTOIP, l)
		assert.Equal(expt, str, "%s", in)
	}
}

func TestEngineID(t *testing.T) {
	assert := assert.New(t)
	l := lt.NewTestContextL(logger.NilContext, t)

	tests := map[string][]byte{
		"38:30:3a:30:30:3a:31:66:3a:38:38:3a:38:30:3a:37:00:00":       []byte{0x38, 0x30, 0x3a, 0x30, 0x30, 0x3a, 0x31, 0x66, 0x3a, 0x38, 0x38, 0x3a, 0x38, 0x30, 0x3a, 0x37, 0x00, 0x00},
		"30:30:3a:31:66:3a:38:38:3a:38:30:3a:37:00:00":                []byte{0x30, 0x30, 0x3a, 0x31, 0x66, 0x3a, 0x38, 0x38, 0x3a, 0x38, 0x30, 0x3a, 0x37, 0x00, 0x00},
		"38:30:3a:30:30:3a:31:66:3a:38:38:3a:38:30:3a:37:00:00:00:00": []byte{0x38, 0x30, 0x3a, 0x30, 0x30, 0x3a, 0x31, 0x66, 0x3a, 0x38, 0x38, 0x3a, 0x38, 0x30, 0x3a, 0x37, 0x00, 0x00, 0x00, 0x00},
	}

	for expt, in := range tests {
		pdu := gosnmp.SnmpPDU{Value: in}
		_, str, _ := GetFromConv(pdu, CONV_ENGINE_ID, l)
		assert.Equal(expt, str, "%s", in)
	}
}

type multiRe struct {
	input   []byte
	re      string
	outputs map[string]string
}

func TestRegex(t *testing.T) {
	assert := assert.New(t)
	l := lt.NewTestContextL(logger.NilContext, t)

	tests := map[int64][]byte{
		62: []byte("    5 Secs ( 96.3762%)   60 Secs ( 62.8549%)  300 Secs ( 25.2877%)"),
		64: []byte("    5 Secs ( 96.3762%)   60 Secs ( 64.8549%)  300 Secs ( 25.2877%)"),
	}

	for expt, in := range tests {
		pdu := gosnmp.SnmpPDU{Value: in}
		ival, _, named := GetFromConv(pdu, CONV_REGEXP+`:60 Secs.*?(\d+)`, l)
		assert.Equal(expt, ival, "%s %v", in, named)
	}

	testStr := []multiRe{
		multiRe{
			input: []byte("Juniper Networks, Inc. mx80-48t , version 13.3R9.13 Build date: 2016-03-01 08:36:50 UTC "),
			re:    `:Juniper Networks, Inc. (?P<model>\S+?) , version (?P<version>\S+)`,
			outputs: map[string]string{
				"model":   "mx80-48t",
				"version": "13.3R9.13",
			},
		},
	}

	for _, in := range testStr {
		pdu := gosnmp.SnmpPDU{Value: in.input}
		_, str, mult := GetFromConv(pdu, CONV_REGEXP+in.re, l)
		for k, v := range in.outputs {
			assert.Equal(v, mult[k], "%s %s %v %s", k, string(in.input), mult, str)
		}
	}
}

func TestToOne(t *testing.T) {
	assert := assert.New(t)
	l := lt.NewTestContextL(logger.NilContext, t)

	tests := map[string]int64{
		"lalalalaal": 1,
		"foofoffofo": 1,
	}

	for in, expt := range tests {
		pdu := gosnmp.SnmpPDU{Value: []byte(in)}
		ival, _, _ := GetFromConv(pdu, CONV_ONE, l)
		assert.Equal(expt, ival, "%s", in)
	}
}
