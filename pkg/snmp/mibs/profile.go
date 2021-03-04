package mibs

import (
	"os"
	"strings"

	"github.com/kentik/ktranslate/pkg/eggs/logger"

	"gopkg.in/yaml.v2"
)

type OID struct {
	Oid  string `yaml:"OID"`
	Name string `yaml:"name"`
}

type Tag struct {
	Column OID    `yaml:"column"`
	Tag    string `yaml:"tag"`
	Symbol string `yaml:"symbol"`
	Index  int    `yaml:"index"`
}

type MIB struct {
	Mib        string `yaml:"MIB"`
	Table      OID    `yaml:"table"`
	Symbols    []OID  `yaml:"symbols"`
	MetricTags []Tag  `yaml:"metric_tags"`
	ForcedType string `yaml:"forced_type"`
	Symbol     OID    `yaml:"symbol"`
}

type Device struct {
	Vendor string `yaml:"vendor"`
}

type StringArray []string

func (a *StringArray) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(&multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		*a = []string{single}
	} else {
		*a = multi
	}
	return nil
}

type Profile struct {
	logger.ContextL
	Metrics     []MIB       `yaml:"metrics"`
	Extends     []string    `yaml:"extends"`
	Device      Device      `yaml:"device"`
	MetricTags  []Tag       `yaml:"metric_tags"`
	Sysobjectid StringArray `yaml:"sysobjectid"`
	extended    bool
	From        string
}

func (p *Profile) extend(extends map[string]*Profile) error {
	if p.extended { // Don't extend multiple times, also halt recursion.
		return nil
	}
	p.extended = true

	for _, name := range p.Extends {
		if ep, ok := extends[name]; !ok {
			p.Errorf("Missing extended profile %s", name)
			continue
		} else {
			// Verify this guy is extended.
			err := ep.extend(extends) // recursive, watch out.
			if err != nil {
				p.Errorf("Cannot extend profile %s %v", name, err)
				continue
			}

			// Merge in
			p.merge(ep)
		}
	}

	return nil
}

func (p *Profile) merge(ep *Profile) {
	p.Metrics = append(p.Metrics, ep.Metrics...)
	p.MetricTags = append(p.MetricTags, ep.MetricTags...)
}

func (mdb *MibDB) LoadProfiles(profileDir string) (int, error) {
	files, err := os.ReadDir(profileDir)
	if err != nil {
		return 0, err
	}

	// Load each profile into a parsed form.
	extends := map[string]*Profile{}
	for _, file := range files {
		fname := profileDir + string(os.PathSeparator) + file.Name()
		if !strings.HasSuffix(fname, ".yaml") {
			continue
		}

		t := Profile{ContextL: mdb.log, From: file.Name()}
		data, err := os.ReadFile(fname)
		if err != nil {
			mdb.log.Errorf("Cannot read file %s", fname)
			continue
		}

		err = yaml.Unmarshal(data, &t)
		if err != nil {
			mdb.log.Errorf("Cannot unmarshal file %s", fname)
			continue
		}

		// Keep this in case someone references this file
		extends[file.Name()] = &t

		// For each sysobjid listed, add this file into our map.
		for _, sysid := range t.Sysobjectid {
			mdb.profiles[sysid] = &t
			mdb.log.Debugf("Adding profile for %s: %s", sysid, t.Device.Vendor)
		}
	}

	// Merge any extended data into the referenced profiles
	for _, pro := range mdb.profiles {
		pro.extend(extends)
	}

	return len(mdb.profiles), nil
}

func (mdb *MibDB) FindProfile(sysid string) *Profile {
	if strings.HasPrefix(sysid, ".") {
		sysid = sysid[1:] // Strip this out if we start with .
	}

	// If we have one directly matching, just return this.
	if p, ok := mdb.profiles[sysid]; ok {
		return p
	}

	// Now walk resursivly up the tree, seeing what profiles are found via a wildcard
	pts := strings.Split(sysid, ".")
	for i := len(pts); i > 1; i-- {
		check := strings.Join(pts[0:i], ".") + ".*"
		if p, ok := mdb.profiles[check]; ok {
			return p
		}
	}

	// Didn't match anything so return nil
	return nil
}
