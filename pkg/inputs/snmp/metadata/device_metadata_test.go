package metadata

import (
	"context"
	"errors"
	"testing"

	"github.com/gosnmp/gosnmp"
	"github.com/stretchr/testify/assert"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
	"github.com/kentik/ktranslate/pkg/kt"
)

type testWalker struct {
	results []gosnmp.SnmpPDU
	err     error
	dm      string
}

func (w testWalker) WalkAll(oid string) ([]gosnmp.SnmpPDU, error) {
	return w.results, w.err
}

func TestPoll(t *testing.T) {
	l := lt.NewTestContextL(logger.NilContext, t)
	conf := &kt.SnmpDeviceConfig{}
	conf.SetTestWalker(testWalker{err: errors.New("Error"), dm: ""}) // Non nil error
	dm := NewDeviceMetadata(nil, conf, kt.NewSnmpDeviceMetric(nil, "deviceA"), l)

	res, err := dm.poll(context.Background(), nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(res.Tables))

	res, err = dm.poll(context.Background(), nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(res.Tables))

	conf.SetTestWalker(testWalker{results: []gosnmp.SnmpPDU{}, dm: ""}) // No results.
	res, err = dm.poll(context.Background(), nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(res.Tables))

	// Test random string
	conf.SetTestWalker(testWalker{results: []gosnmp.SnmpPDU{{Value: "not a byte-slice", Name: ".1.1.2.3"}}, dm: ""}) // string value.
	res, err = dm.poll(context.Background(), nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(res.Tables))

	// Test sysdesc
	conf.SetTestWalker(testWalker{results: []gosnmp.SnmpPDU{{Value: []byte("sysdesc"), Name: ".1.3.6.1.2.1.1.1.0"}}, dm: ""}) // SysDescr.
	res, err = dm.poll(context.Background(), nil)
	assert.NoError(t, err)
	assert.Equal(t, "sysdesc", res.SysDescr)

	meta := map[string]*kt.Mib{
		"1.1.1": &kt.Mib{Name: "foo"},
	}
	dm = NewDeviceMetadata(meta, conf, kt.NewSnmpDeviceMetric(nil, "deviceA"), l)
	conf.SetTestWalker(testWalker{results: []gosnmp.SnmpPDU{{Value: "not a byte-slice", Name: ".1.1.1"}}, dm: ""}) // string value.
	res, err = dm.poll(context.Background(), nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(res.Tables))

	conf.SetTestWalker(testWalker{results: []gosnmp.SnmpPDU{{Value: "not a byte-slice", Name: ".1.1.1.2"}}, dm: ""}) // table.
	res, err = dm.poll(context.Background(), nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res.Tables))
}
