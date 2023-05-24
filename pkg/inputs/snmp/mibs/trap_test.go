package mibs

import (
	"io/fs"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
)

func TestLoadTraps(t *testing.T) {
	l := lt.NewTestContextL(logger.NilContext, t)
	mdb, err := NewMibDB("", "", false, l)
	assert.NoError(t, err)
	defer mdb.Close()
	content := []byte(`
traps:
  # Traps from Visual Data Center (VDC), formerly owned by Optimum Path
  - trap_oid: 1.3.6.1.4.1.34510.2.1.1.6.1
    trap_name: vdcAlarms
    drop_undefined: true
    events:
      - name: vdcAlarmId
        OID: 1.3.6.1.4.1.34510.2.1.1.3.2.1.1
      - name: vdcAlarmDesc
        OID: 1.3.6.1.4.1.34510.2.1.1.3.2.1.2
      - name: vdcAlarmTime
        OID: 1.3.6.1.4.1.34510.2.1.1.3.2.1.3
      - name: vdcLevel
        OID: 1.3.6.1.4.1.34510.2.1.1.3.2.1.4
      - name: vdcDevice
        OID: 1.3.6.1.4.1.34510.2.1.1.3.2.1.5
      - name: vdcMonAttribWild
        OID: 1.3.6.1.4.1.34510.22.*
      - name: cbgpPeerLastErrorTxt
        OID: 1.3.6.1.4.1.9.9.187.1.2.1.1.7.{bgpPeerRemoteAddr:*}.{ifIndex:1}
      - name: chsrpTrapVarBing
        OID: 1.3.6.1.4.1.9.9.187.1.3.4.5.{ifIndex:1}.{cHsrpGrpTable:1}
      - name: chsrpTwoOne
        OID: 1.3.6.1.4.1.9.9.666.1.3.4.5.{cHsrpGrpTable:2}.{ifIndex:1}
      - name: chsrpTrapVarYandex
        OID: 1.3.6.1.4.1.9.9.1811.1.3.4.5.{ifIndex:1}.{cHsrpGrpTable:2}
    attributes:
      responsible_team: vmware
      cHsrpGrpTable: 777
`)
	// Save test data to local.
	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.FailNow()
	}
	if _, err := file.Write(content); err != nil {
		t.FailNow()
	}
	defer os.Remove(file.Name())
	fe, _ := file.Stat()
	err = mdb.parseTrapsFromYml(file.Name(), fs.FileInfoToDirEntry(fe), nil)
	assert.NoError(t, err)
	assert.Equal(t, 10, len(mdb.trapMibs))

	tr, attr, err := mdb.GetForKey(".1.3.6.1.4.1.34510.2.1.1.3.2.1.1")
	assert.NoError(t, err)
	assert.Equal(t, "vdcAlarmId", tr.Name)
	assert.Equal(t, 2, len(attr))

	tr, attr, err = mdb.GetForKey(".1.3.6.1.4.1.34510.22.1.2.3.4.5.6")
	assert.NoError(t, err)
	assert.Equal(t, "vdcMonAttribWild", tr.Name)
	assert.Equal(t, 2, len(attr))

	tr, attr, err = mdb.GetForKey(".1.3.6.1.4.1.9.9.187.1.2.1.1.7.2.2.2.3")
	assert.NoError(t, err)
	assert.Equal(t, "cbgpPeerLastErrorTxt", tr.Name)
	assert.Equal(t, "2.2.2.3", attr["bgpPeerRemoteAddr"])
	assert.Equal(t, "", attr["ifIndex"]) // Since this is after a wildcard, assume nothing because there are no valid tokens to consume.
	assert.Equal(t, "vmware", attr["responsible_team"])
	assert.Equal(t, "777", attr["cHsrpGrpTable"])

	tr, attr, err = mdb.GetForKey(".1.3.6.1.4.1.9.9.187.1.3.4.5.999.666")
	assert.NoError(t, err)
	assert.Equal(t, "chsrpTrapVarBing", tr.Name)
	assert.Equal(t, "999", attr["ifIndex"])
	assert.Equal(t, "666", attr["cHsrpGrpTable"])
	assert.Equal(t, "vmware", attr["responsible_team"])

	tr, attr, err = mdb.GetForKey(".1.3.6.1.4.1.9.9.666.1.3.4.5.66.99.333")
	assert.NoError(t, err)
	assert.Equal(t, "chsrpTwoOne", tr.Name)
	assert.Equal(t, "333", attr["ifIndex"])
	assert.Equal(t, "66.99", attr["cHsrpGrpTable"])
	assert.Equal(t, "vmware", attr["responsible_team"])

	tr, attr, err = mdb.GetForKey(".1.3.6.1.4.1.9.9.1811.1.3.4.5.66.99.333")
	assert.NoError(t, err)
	assert.Equal(t, "chsrpTrapVarYandex", tr.Name)
	assert.Equal(t, "66", attr["ifIndex"])
	assert.Equal(t, "99.333", attr["cHsrpGrpTable"])
	assert.Equal(t, "vmware", attr["responsible_team"])
}
