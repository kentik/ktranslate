package mibs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"gopkg.in/yaml.v3"
)

const (
	SysidMap    = "SYSID_MAP.json"
	permissions = 0644
)

type Constraint struct {
	Enumeration map[string]int64 `json:"enumeration"`
}

type Type struct {
	Constraints Constraint `json:"constraints"`
}

type PySyntax struct {
	Type  string `json:"type"`
	Class string `json:"class"`
}

type PySMIMib struct {
	Name     string   `json:"name"`
	Oid      string   `json:"oid"`
	Syntax   PySyntax `json:"syntax"`
	Desc     string   `json:"description"`
	Class    string   `json:"class"`
	Nodetype string   `json:"nodetype"`
	Type     Type     `json:"type"`
}

type PySysidMap map[string][]string

type PySMIMibFile map[string]*PySMIMib

/**
Code to generate a yaml from file a .mib file. To run:

* Install https://github.com/etingof/pysmi. (pip install pysmi)
* Install librenms. (git clone https://github.com/librenms/librenms.git)
* Download your mib target. Place into /tmp/snmp_in directory. Create /tmp/snmp_out directory.
Use https://help.zscaler.com/downloads/zia/documentation-knowledgebase/analytics/snmp-mibs/about-the-zscaler-snmp-mibs/zscaler-nss-mib.mib as a test if you need one.

* Run mibdump:
mibdump.py --mib-source=file:///path/to/src/librenms/mibs --mib-source=file:///tmp/snmp_in --generate-mib-texts --ignore-errors  --destination-format json --destination-directory=/tmp/snmp_out MY_MIB_NAME

For example, using mib ZSCALER-NSS-MIB:
mibdump.py --mib-source=file:///home/pye/src/librenms/mibs --mib-source=file:///tmp/snmp_in --generate-mib-texts --ignore-errors  --destination-format json --destination-directory=/tmp/snmp_out ZSCALER-NSS-MIB

* Verify that /tmp/snmp_out/ZSCALER-NSS-MIB.json exists

* Convert to yaml
docker run --rm -v /tmp/snmp_out:/snmp_out kentik/ktranslate:v2 -snmp /etc/ktranslate/snmp-base.yaml -snmp_json2yaml /snmp_out/MY_MIB_NAME.json

For example, using ZSCALER-NSS-MIB:
docker run --rm -v /tmp/snmp_out:/snmp_out kentik/ktranslate:v2 -snmp /etc/ktranslate/snmp-base.yaml -snmp_json2yaml /snmp_out/ZSCALER-NSS-MIB.json

* The final yaml file is at /tmp/snmp_out/ZSCALER-NSS-MIB.yaml

NOTE: human checking is still needed. Notably, you will need to:

1) Add an extends section.
2) Ensure that the sysobjectid section is present and sane.

*/
func ConvertJson2Yaml(file string, log logger.ContextL) error {
	// Load up from json.
	pmib, err := loadJson(file, log)
	if err != nil {
		return err
	}

	mibset := []*MIB{}
	tables := map[string]*MIB{}
	sysoids := []string{}
	finals := []MIB{}
	name := strings.ReplaceAll(file, ".json", "")
	pts := strings.Split(name, "/")
	mibName := pts[len(pts)-1]
	if pmib != nil {
		for _, mib := range *pmib {
			if mib.Name == "" { // If there's no name, go ahead and skip.
				continue
			}
			log.Infof("Evaluating %s", mib.Name)
			if mib.Class == "objectidentity" {
				log.Infof("Adding as a sysoid: %s", mib.Name)
				sysoids = append(sysoids, mib.Oid)
				continue
			}

			// Else, try to convert this into a mib set.
			mb := ToMib(*pmib, mibName, mib, log)
			if mb != nil {
				if mb.Table.Oid != "" {
					tables[mb.Table.Oid] = mb
				} else {
					mibset = append(mibset, mb)
				}
			}
		}
	}

	// See if these oids are contained in any tables.
	for _, mib := range mibset {
		added := false
		for oid, table := range tables {
			if mib.Symbol.Oid != "" && strings.HasPrefix(mib.Symbol.Oid, oid) {
				table.Symbols = append(table.Symbols, mib.Symbol)
				added = true
				break
			}
			if len(mib.MetricTags) > 0 && strings.HasPrefix(mib.MetricTags[0].Column.Oid, oid) {
				table.MetricTags = append(table.MetricTags, mib.MetricTags...)
				added = true
				break
			}
		}
		// Else, this is a top level mib, just add in here.
		if !added {
			finals = append(finals, *mib)
		}
	}

	// Copy the tables to finals.
	for _, table := range tables {
		if len(table.Symbols) > 0 || table.Symbol.Oid != "" {
			finals = append(finals, *table)
		}
	}

	// Remove any metrics with no metrics (just tags). Needed because otherwise validation will fail.
	gm := []MIB{}
	for _, metric := range finals {
		if metric.Symbol.Oid != "" || len(metric.Symbols) > 0 {
			gm = append(gm, metric)
		}
	}

	// Sort the mibs.
	sort.Sort(MibList(gm))

	// Make a profile.
	pro := Profile{
		Metrics:     gm,
		Sysobjectid: sysoids,
		ContextL:    log,
		From:        name,
	}

	// Does this profile make sense?
	if err := pro.validate(); err != nil {
		time.Sleep(200 * time.Millisecond)
		return err
	}

	// Now, write out the profile.
	pro.From = ""
	t, err := yaml.Marshal(pro)
	if err != nil {
		return err
	}

	// Add a header on
	header := []byte(fmt.Sprintf("# Autogenerated by ktranslate on %v from %s.json\n\n", time.Now().Format(time.UnixDate), mibName))
	full := append(header, t...)

	// Write out.
	err = ioutil.WriteFile(name+".yaml", full, permissions)
	if err != nil {
		return err
	}

	time.Sleep(200 * time.Millisecond)
	return fmt.Errorf("ok")
}

func loadJson(fname string, log logger.ContextL) (*PySMIMibFile, error) {
	data, err := os.ReadFile(fname)
	if err != nil {
		return nil, fmt.Errorf("There was an error when accessing the %s file: %v.", fname, err)
	}

	pmib := PySMIMibFile{}
	err = json.Unmarshal(data, &pmib)
	if err != nil {
		return nil, fmt.Errorf("There was an error when unmarshalling the %s file: %v.", fname, err)
	}

	return &pmib, nil
}

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

		var pmib *PySMIMibFile
		if file.Name() == SysidMap { // Special file which maps loads mids to sysoid markers.
			data, err := os.ReadFile(fname)
			if err != nil {
				mdb.log.Errorf("There was an error when accessing the %s file: %v.", fname, err)
				continue
			}

			err = json.Unmarshal(data, &psysmap)
			if err != nil {
				mdb.log.Errorf("There was an error when marshalling the %s sysmap file: %v.", fname, err)
			}
			continue
		} else {
			pm, err := loadJson(fname, mdb.log)
			if err != nil {
				mdb.log.Errorf("%v", err)
				continue
			}
			pmib = pm
		}

		// For each sysobjid listed, add this file into our map.
		if pmib != nil {
			for _, mib := range *pmib {
				name := strings.ReplaceAll(file.Name(), ".json", "")
				if _, ok := mibset[name]; !ok {
					mibset[name] = []MIB{}
				}
				mb := ToMib(*pmib, name, mib, mdb.log)
				if mb != nil {
					mibset[name] = append(mibset[name], *mb)
				}
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

func ToMib(mibSet map[string]*PySMIMib, mibName string, pm *PySMIMib, log logger.ContextL) *MIB {
	if strings.HasSuffix(pm.Name, "Table") {
		mib := MIB{
			Mib: mibName,
			Table: OID{
				Oid:  pm.Oid,
				Name: pm.Name,
				Desc: pm.Desc,
			},
			Symbols:    []OID{},
			MetricTags: []Tag{},
			sortKey:    pm.Name,
		}
		return &mib
	}

	// If there's no oid, return right away.
	if pm.Oid == "" {
		return nil
	}

	switch pm.Syntax.Type {
	case "Counter32", "Counter64", "Gauge32", "Gauge64":
		mib := MIB{
			Mib: mibName,
			Symbol: OID{
				Oid:  pm.Oid,
				Name: pm.Name,
				Desc: pm.Desc,
			},
			sortKey: pm.Name,
		}
		return &mib
	case "DisplayString", "InterfaceIndex", "InterfaceIndexOrZero", "OCTET STRING":
		mib := MIB{
			Mib: mibName,
			MetricTags: []Tag{
				Tag{
					Column: OID{
						Oid:  pm.Oid,
						Name: pm.Name,
						Desc: pm.Desc,
					},
				},
			},
			sortKey: pm.Name,
		}
		return &mib
	case "INTEGER", "Integer32":
		var mib = MIB{}
		switch pm.Nodetype {
		case "column":
			mib = MIB{
				Mib: mibName,
				MetricTags: []Tag{
					Tag{
						Column: OID{
							Oid:  pm.Oid,
							Name: pm.Name,
							Desc: pm.Desc,
						},
					},
				},
				sortKey: pm.Name,
			}
		default:
			mib = MIB{
				Mib: mibName,
				Symbol: OID{
					Oid:  pm.Oid,
					Name: pm.Name,
					Desc: pm.Desc,
				},
				sortKey: pm.Name,
			}
		}
		return &mib
	default:
		if ms, ok := mibSet[pm.Syntax.Type]; ok {
			mib := MIB{
				Mib: mibName,
				Symbol: OID{
					Oid:  pm.Oid,
					Name: pm.Name,
					Desc: pm.Desc,
					Enum: ms.Type.Constraints.Enumeration,
				},
				sortKey: pm.Name,
			}
			return &mib
		} else {
			log.Warnf("Skipping type %s for %s", pm.Syntax.Type, pm.Name)
		}
	}

	return nil
}

// Utils to sort
type MibList []MIB

func (s MibList) Len() int           { return len(s) }
func (s MibList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s MibList) Less(i, j int) bool { return s[i].sortKey < s[j].sortKey }
