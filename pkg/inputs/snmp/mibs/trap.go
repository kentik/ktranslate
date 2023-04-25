package mibs

import (
	"os"
	"strconv"
	"strings"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"

	"gopkg.in/yaml.v2"
)

type Trap struct {
	Oid           string `yaml:"trap_oid"`
	Name          string `yaml:"trap_name"`
	DropUndefined bool   `yaml:"drop_undefined"`
	Events        []OID  `yaml:"events"`
}

type TrapBase struct {
	logger.ContextL `yaml:"-"`
	Traps           []Trap `yaml:"traps"`
	From            string `yaml:"from,omitempty"`
}

func (t *Trap) DropUndefinedVars() bool {
	if t == nil {
		return false
	}
	return t.DropUndefined
}

func normalizeOid(oid string) string {
	if strings.HasPrefix(oid, ".") {
		return oid
	}
	return "." + oid
}

func (mdb *MibDB) parseTrapsFromYml(fname string, file os.DirEntry, extends map[string]*Profile) error {
	t := TrapBase{ContextL: mdb.log, From: file.Name()}
	data, err := os.ReadFile(fname)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &t)
	if err != nil {
		return err
	}

	added := 0
	for _, trap := range t.Traps {
		for _, event := range trap.Events {
			oid := event.Oid
			kvs := map[string][]int{}
			pts := strings.Split(oid, ".")
			newOid := []string{}
			nextWildcard := 0
			for i := 0; i < len(pts); i++ {
				last := pts[i]
				if len(last) > 2 && last[0:1] == "{" && last[len(last)-1:] == "}" { // Handle this as a variable.
					set := strings.Split(last[1:len(last)-1], ":")
					if len(set) == 2 {
						vlen, err := strconv.Atoi(set[1])
						if err == nil {
							kvs[set[0]] = []int{nextWildcard + 1, vlen}
							nextWildcard += vlen
						} else if set[1] == "*" { // Wildcard means use all the rest.
							kvs[set[0]] = []int{nextWildcard + 1, 0}
							break // don't consume any more variables because the wildcard ate them all.
						} else {
							// Noop?
						}
					}
				} else {
					nextWildcard += 1
					newOid = append(newOid, last) // Put this on as a regular key.
				}
			}
			if len(kvs) > 0 {
				newOid = append(newOid, "*") // End the oid with a wildcard because we're matching on variables.
			}
			oid = strings.Join(newOid, ".")

			mib := &kt.Mib{
				Oid:        oid,
				Name:       event.Name,
				Enum:       event.Enum,
				Tag:        event.Tag,
				Conversion: event.Conversion,
				Extra:      trap.Name,
				Mib:        trap.Oid,
				VarSet:     kvs,
			}
			if len(mib.Enum) > 0 {
				mib.EnumRev = make(map[int64]string)
			}
			for k, v := range mib.Enum {
				mib.Enum[strings.ToLower(k)] = v
				mib.EnumRev[v] = k
			}
			added++
			mdb.trapMibs[normalizeOid(mib.Oid)] = mib
		}
		mdb.traps[normalizeOid(trap.Oid)] = trap
	}

	mdb.log.Infof("Loading %d snmp trap data points and %d traps from %s.", added, len(mdb.traps), fname)

	return nil
}
