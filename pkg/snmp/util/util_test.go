package util

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
	"github.com/kentik/gosnmp"
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
