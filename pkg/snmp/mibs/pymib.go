package mibs

import (
	"encoding/json"
	"os"
	"strings"
)

const (
	SysidMap = "SYSID_MAP.json"
)

type PySyntax struct {
	Type  string `json:"type"`
	Class string `json:"class"`
}

type PySMIMib struct {
	Name   string   `json:"name"`
	Oid    string   `json:"oid"`
	Syntax PySyntax `json:"nodetype`
}

type PySysidMap map[string][]string

type PySMIMibFile map[string]*PySMIMib

func (mdb *MibDB) LoadPyMibSet(profileDir string) (int, error) {
	files, err := os.ReadDir(profileDir)
	if err != nil {
		return 0, err
	}

	psysmap := PySysidMap{}
	mibset := map[string][]MIB{}
	for _, file := range files {
		fname := profileDir + string(os.PathSeparator) + file.Name()
		if !strings.HasSuffix(fname, ".json") {
			continue
		}
		data, err := os.ReadFile(fname)
		if err != nil {
			mdb.log.Errorf("Cannot read file %s", fname)
			continue
		}

		if file.Name() == SysidMap { // Special file which maps loads mids to sysoid markers.
			err = json.Unmarshal(data, &psysmap)
			if err != nil {
				mdb.log.Errorf("Cannot unmarshal sysmap file %s", fname)
			}
			continue
		}

		pmib := PySMIMibFile{}
		err = json.Unmarshal(data, &pmib)
		if err != nil {
			mdb.log.Errorf("Cannot unmarshal file %s", fname)
			continue
		}

		// For each sysobjid listed, add this file into our map.
		for _, mib := range pmib {
			name := strings.ReplaceAll(file.Name(), ".json", "")
			if _, ok := mibset[name]; !ok {
				mibset[name] = []MIB{}
			}
			mb := ToMib(name, mib)
			if mb != nil {
				mibset[name] = append(mibset[name], *mb)
			}
		}
	}

	// Now, if we have any, load these into per sysid profiles.
	if len(mibset) > 0 {
		for sysid, mibs := range psysmap {
			t := Profile{ContextL: mdb.log, From: profileDir, Sysobjectid: []string{sysid}, Metrics: []MIB{}}
			for _, mib := range mibs {
				if mb, ok := mibset[mib]; ok {
					t.Metrics = append(t.Metrics, mb...)
				}
			}
			if len(t.Metrics) > 0 {
				mdb.profiles[sysid] = &t
				mdb.log.Infof("Add %s with %d mibs", sysid, len(t.Metrics))
			} else {
				mdb.log.Warnf("Skipping %s, no Mibs found", sysid)
			}
		}
	}

	return len(mibset), nil
}

func ToMib(mibName string, pm *PySMIMib) *MIB {
	switch pm.Syntax.Type {
	case "Counter32", "Counter64", "Gauge32", "Gauge64":
		mib := MIB{
			Mib: mibName,
			Table: OID{
				Oid:  pm.Oid,
				Name: pm.Name,
			},
			Symbol: OID{
				Oid:  pm.Oid,
				Name: pm.Name,
			},
			ForcedType: "monotonic_count",
		}
		return &mib
	case "DisplayString", "INTEGER", "Integer32", "InterfaceIndex", "InterfaceIndexOrZero":
		mib := MIB{
			Mib: mibName,
			MetricTags: []Tag{
				Tag{
					Column: OID{
						Oid:  pm.Oid,
						Name: pm.Name,
					},
				},
			},
			ForcedType: pm.Syntax.Type,
		}
		return &mib
	}

	return nil
}
