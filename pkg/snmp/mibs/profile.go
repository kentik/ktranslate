package mibs

import (
	"os"
	"strings"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"

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
	extends := map[string]*Profile{}

	// Recursively get all the profiles found.
	mdb.log.Infof("Looking at mib profiles in %s", profileDir)
	err := mdb.loadProfileDir(profileDir, extends)
	if err != nil {
		return 0, err
	}

	// Merge any extended data into the referenced profiles
	mdb.log.Infof("Now trying to extend profiles")
	for _, pro := range mdb.profiles {
		pro.extend(extends)
	}

	return len(mdb.profiles), nil
}

func (mdb *MibDB) loadProfileDir(profileDir string, extends map[string]*Profile) error {
	files, err := os.ReadDir(profileDir)
	if err != nil {
		return err
	}

	// Load each profile into a parsed form.
	for _, file := range files {
		fname := profileDir + string(os.PathSeparator) + file.Name()

		// Now, recurse down if this file actually a directory
		info, err := os.Stat(fname)
		if err != nil {
			mdb.log.Errorf("Cannot stat dir %s", fname)
			continue
		}
		if info.IsDir() {
			mdb.log.Infof("Recursing into %s", fname)
			err := mdb.loadProfileDir(fname, extends)
			if err != nil {
				mdb.log.Errorf("Cannot recurse into directory %s", fname)
				continue
			}
		}

		if !strings.HasSuffix(fname, ".yaml") && !strings.HasSuffix(fname, ".yml") {
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
			mdb.log.Errorf("Cannot unmarshal file %s -> %v", fname, err)
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

	return nil
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
	for i := len(pts); i > 0; i-- {
		check := strings.Join(pts[0:i], ".") + ".*"
		if p, ok := mdb.profiles[check]; ok {
			return p
		}
	}

	// Didn't match anything so return nil
	return nil
}

func (p *Profile) DumpOids(log logger.ContextL) {
	log.Infof("Device Tags:")
	for _, tag := range p.MetricTags {
		if tag.Column.Oid != "" {
			log.Infof("   -> %s -> %s", tag.Column.Oid, tag.Column.Name)
		}
	}

	log.Infof("Device Metrics:")
	for _, metric := range p.Metrics {
		log.Infof("MIB -> %s | %s %s %s", metric.Mib, metric.Table.Name, metric.ForcedType, metric.Symbol)
		for _, s := range metric.Symbols {
			if s.Oid != "" {
				log.Infof("   -> %s -> %s", s.Oid, s.Name)
			}
		}
		if metric.Symbol.Oid != "" {
			log.Infof("Symbol   -> %s -> %s", metric.Symbol.Oid, metric.Symbol.Name)
		}

		for _, tag := range metric.MetricTags {
			if tag.Column.Oid != "" {
				log.Infof("Tag   -> %s -> %s %s %s", tag.Column.Oid, tag.Column.Name, tag.Tag, tag.Symbol)
			}
		}
	}
}

// Return oids for metrics (type counter) for the enabled mibs
// IF-MIB | ifXTable monotonic_count { } This is an interface metric because it starts with a if
// UDP-MIB |  monotonic_count {1.3.6.1.2.1.7.8.0 udpHCInDatagrams} This is a device metrics because it does't, but still is a counter.

func (p *Profile) GetMetrics(enabledMibs []string) (map[string]*kt.Mib, map[string]*kt.Mib) {
	deviceMetrics := map[string]*kt.Mib{}
	interfaceMetrics := map[string]*kt.Mib{}

	enabled := map[string]bool{}
	for _, mib := range enabledMibs {
		enabled[mib] = true
	}

	for _, metric := range p.Metrics {
		if !enabled[metric.Mib] { // Only look at mibs we have opted into caring about.
			continue
		}

		var otype kt.Oidtype
		switch metric.ForcedType {
		case "monotonic_count", "monotonic_count_and_rate":
			otype = kt.Counter
		default: // We only are looking for metric type values here.
			continue
		}

		// TODO -- so we want to collase Symobol and Symbols?
		if metric.Symbol.Oid != "" {
			mib := &kt.Mib{
				Oid:  metric.Symbol.Oid,
				Name: metric.Symbol.Name,
				Type: otype,
			}
			if strings.HasPrefix(metric.Symbol.Name, "if") {
				interfaceMetrics[metric.Symbol.Oid] = mib
			} else {
				deviceMetrics[metric.Symbol.Oid] = mib
			}
		}

		for _, s := range metric.Symbols {
			mib := &kt.Mib{
				Oid:  s.Oid,
				Name: s.Name,
				Type: otype,
			}
			if strings.HasPrefix(s.Name, "if") {
				interfaceMetrics[s.Oid] = mib
			} else {
				deviceMetrics[s.Oid] = mib
			}
		}
	}

	return deviceMetrics, interfaceMetrics
}

func (p *Profile) GetMetadata(enabledMibs []string) (map[string]*kt.Mib, map[string]*kt.Mib) {
	deviceMetadata := map[string]*kt.Mib{}
	interfaceMetadata := map[string]*kt.Mib{}

	enabled := map[string]bool{}
	for _, mib := range enabledMibs {
		enabled[mib] = true
	}

	for _, tag := range p.MetricTags {
		if tag.Column.Oid != "" {
			mib := &kt.Mib{
				Oid:  tag.Column.Oid,
				Name: tag.Column.Name,
				Type: kt.String,
			}
			deviceMetadata[tag.Column.Oid] = mib
		}
	}

	for _, metric := range p.Metrics {
		if !enabled[metric.Mib] { // Only look at mibs we have opted into caring about.
			continue
		}

		for _, t := range metric.MetricTags {
			mib := &kt.Mib{
				Oid:  t.Column.Oid,
				Name: t.Column.Name,
				Type: kt.String,
			}
			if strings.HasPrefix(t.Column.Name, "if") {
				interfaceMetadata[t.Column.Oid] = mib
			} else {
				deviceMetadata[t.Column.Oid] = mib
			}
		}
	}
	return deviceMetadata, interfaceMetadata
}
