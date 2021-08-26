package mibs

import (
	"fmt"
	"os"
	"strings"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/inputs/snmp/mibs/apc"
	"github.com/kentik/ktranslate/pkg/kt"

	"gopkg.in/yaml.v2"
)

type OID struct {
	Oid        string           `yaml:"OID,omitempty"`
	Name       string           `yaml:"name,omitempty"`
	Enum       map[string]int64 `yaml:"enum,omitempty"`
	Tag        string           `yaml:"tag,omitempty"`
	Desc       string           `yaml:"description,omitempty"`
	Conversion string           `yaml:"conversion,omitempty"`
}

type Tag struct {
	Column OID    `yaml:"column,omitempty"`
	Tag    string `yaml:"tag,omitempty"`
	Symbol string `yaml:"symbol,omitempty"`
	Index  int    `yaml:"index,omitempty"`
}

type MIB struct {
	Mib        string `yaml:"MIB,omitempty"`
	Table      OID    `yaml:"table,omitempty"`
	Symbols    []OID  `yaml:"symbols,omitempty"`
	MetricTags []Tag  `yaml:"metric_tags,omitempty"`
	ForcedType string `yaml:"forced_type,omitempty"`
	Symbol     OID    `yaml:"symbol,omitempty"`
	sortKey    string
}

type Device struct {
	Vendor string `yaml:"vendor,omitempty"`
}

type Profile struct {
	logger.ContextL `yaml:"-"`
	Metrics         []MIB          `yaml:"metrics,omitempty"`
	Extends         []string       `yaml:"extends,omitempty"`
	Device          Device         `yaml:"device,omitempty"`
	MetricTags      []Tag          `yaml:"metric_tags,omitempty"`
	Sysobjectid     kt.StringArray `yaml:"sysobjectid,omitempty"`
	From            string         `yaml:"from,omitempty"`
	Provider        kt.Provider    `yaml:"provider,omitempty"`
	extended        bool
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
		pro.validate()
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

		// Ignore hidden files.
		if file.Name()[0:1] == "." {
			continue
		}

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

		pts := strings.Split(fname, ".")
		switch pts[len(pts)-1] {
		case "xml":
			err := mdb.parseMibFromXml(fname)
			if err != nil {
				mdb.log.Errorf("Cannot parse XML mib %s %v", fname, err)
			}
		case "yaml", "yml":
			err := mdb.parseMibFromYml(fname, file, extends)
			if err != nil {
				mdb.log.Errorf("Cannot parse Yaml mib %s %v", fname, err)
			}
		default:
			if len(pts) > 1 {
				mdb.log.Infof("Ignoring file %s", fname)
			}
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

func (p *Profile) validate() error {
	hasErr := false
	for _, metric := range p.Metrics {
		if metric.Symbol.Oid == "" && len(metric.Symbols) == 0 {
			p.Warnf("Possibly corrupted table? %s %s %v", p.From, metric.Mib, metric)
			hasErr = true
		}
		if metric.Symbol.Oid != "" && metric.Symbol.Name == "" {
			p.Warnf("Possibly corrupted symbol oid? %s %s %v", p.From, metric.Mib, metric.Symbol.Oid)
			hasErr = true
		}
		for _, s := range metric.Symbols {
			if s.Oid != "" && s.Name == "" {
				p.Warnf("Possibly corrupted symbols oid? %s %s %v", p.From, metric.Mib, s.Oid)
				hasErr = true
			}
		}
	}

	for _, tag := range p.MetricTags {
		if tag.Column.Oid == "" {
			p.Warnf("Possibly corrupted metadata table? %s %v", p.From, tag)
			hasErr = true
		}
		if tag.Column.Oid != "" && tag.Column.Name == "" {
			p.Warnf("Possibly corrupted metadata table? %s %v", p.From, tag)
			hasErr = true
		}
	}

	for _, metric := range p.Metrics {
		for _, tag := range metric.MetricTags {
			if tag.Column.Oid != "" && tag.Column.Name == "" {
				p.Warnf("Possibly corrupted metadata table? %s %v", p.From, tag)
				hasErr = true
			}
		}
	}

	if hasErr {
		return fmt.Errorf("Error in %s profile", p.From)
	}
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
	enabledTables := map[string]map[string]bool{}
	for _, mib := range enabledMibs {
		pts := strings.Split(mib, ".")
		if len(pts) == 1 {
			enabled[mib] = true
		} else {
			enabled[pts[0]] = true
			if _, ok := enabledTables[pts[0]]; !ok {
				enabledTables[pts[0]] = map[string]bool{}
			}
			enabledTables[pts[0]][pts[1]] = true
		}
	}

	for _, metric := range p.Metrics {
		if !enabled[metric.Mib] { // Only look at mibs we have opted into caring about.
			continue
		}
		if enabledTables[metric.Mib] != nil {
			if !enabledTables[metric.Mib][metric.Table.Name] { // And also worry about names.
				continue
			}
		}

		var otype kt.Oidtype
		if metric.ForcedType != "" {
			switch metric.ForcedType {
			case "monotonic_count", "monotonic_count_and_rate":
				otype = kt.Counter
			default: // We only are looking for metric type values here.
				continue
			}
		} else {
			otype = kt.Counter
		}

		// TODO -- so we want to collase Symbol and Symbols?
		if metric.Symbol.Oid != "" {
			mib := &kt.Mib{
				Oid:        metric.Symbol.Oid,
				Name:       metric.Symbol.Name,
				Type:       otype,
				Enum:       metric.Symbol.Enum,
				Tag:        metric.Symbol.Tag,
				Conversion: metric.Symbol.Conversion,
			}
			if len(mib.Enum) > 0 {
				mib.EnumRev = make(map[int64]string)
			}
			for k, v := range mib.Enum {
				mib.Enum[strings.ToLower(k)] = v
				mib.EnumRev[v] = k
			}
			if strings.HasPrefix(metric.Symbol.Name, "if") {
				interfaceMetrics[metric.Symbol.Oid] = mib
			} else {
				deviceMetrics[metric.Symbol.Oid] = mib
			}
		}

		for _, s := range metric.Symbols {
			mib := &kt.Mib{
				Oid:        s.Oid,
				Name:       s.Name,
				Type:       otype,
				Enum:       s.Enum,
				Tag:        s.Tag,
				Conversion: s.Conversion,
			}
			if len(mib.Enum) > 0 {
				mib.EnumRev = make(map[int64]string)
			}
			for k, v := range mib.Enum {
				mib.Enum[strings.ToLower(k)] = v
				mib.EnumRev[v] = k
			}
			if strings.HasPrefix(s.Name, "if") {
				interfaceMetrics[s.Oid] = mib
			} else {
				deviceMetrics[s.Oid] = mib
			}
		}
	}

	// If we have any duplicate names, keep the longest OID.
	prune(deviceMetrics)
	prune(interfaceMetrics)

	return deviceMetrics, interfaceMetrics
}

func (p *Profile) GetMetadata(enabledMibs []string) (map[string]*kt.Mib, map[string]*kt.Mib) {
	deviceMetadata := map[string]*kt.Mib{}
	interfaceMetadata := map[string]*kt.Mib{}

	enabled := map[string]bool{}
	enabledTables := map[string]map[string]bool{}
	for _, mib := range enabledMibs {
		pts := strings.Split(mib, ".")
		if len(pts) == 1 {
			enabled[mib] = true
		} else {
			enabled[pts[0]] = true
			if _, ok := enabledTables[pts[0]]; !ok {
				enabledTables[pts[0]] = map[string]bool{}
			}
			enabledTables[pts[0]][pts[1]] = true
		}
	}

	// Top level tags here.
	for _, tag := range p.MetricTags {
		if tag.Column.Oid != "" {
			mib := &kt.Mib{
				Oid:        tag.Column.Oid,
				Name:       tag.Column.Name,
				Type:       kt.String,
				Tag:        tag.Tag,
				Conversion: tag.Column.Conversion,
			}
			deviceMetadata[tag.Column.Oid] = mib
		}
	}

	// per metric tags here.
	for _, metric := range p.Metrics {
		if !enabled[metric.Mib] { // Only look at mibs we have opted into caring about.
			continue
		}
		if enabledTables[metric.Mib] != nil {
			if !enabledTables[metric.Mib][metric.Table.Name] { // And also worry about names.
				continue
			}
		}

		for _, t := range metric.MetricTags {
			if t.Column.Oid != "" {
				mib := &kt.Mib{
					Oid:        t.Column.Oid,
					Name:       t.Column.Name,
					Type:       kt.String,
					Tag:        t.Tag,
					Conversion: t.Column.Conversion,
				}
				if strings.HasPrefix(t.Column.Name, "if") {
					interfaceMetadata[t.Column.Oid] = mib
				} else {
					deviceMetadata[t.Column.Oid] = mib
				}
			}
		}
	}

	// If we have any duplicate names, keep the longest OID.
	prune(deviceMetadata)
	prune(interfaceMetadata)

	return deviceMetadata, interfaceMetadata
}

func (p *Profile) GetMibs() map[string]bool {
	mibs := map[string]bool{}
	for _, metric := range p.Metrics {
		mibs[metric.Mib] = true
	}
	return mibs
}

func (mdb *MibDB) parseMibFromYml(fname string, file os.DirEntry, extends map[string]*Profile) error {
	t := Profile{ContextL: mdb.log, From: file.Name()}
	data, err := os.ReadFile(fname)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &t)
	if err != nil {
		return err
	}

	// Keep this in case someone references this file
	extends[file.Name()] = &t

	// For each sysobjid listed, add this file into our map.
	for _, sysid := range t.Sysobjectid {
		if strings.HasPrefix(sysid, ".") {
			sysid = sysid[1:] // Strip this out if we start with .
		}
		if op, ok := mdb.profiles[sysid]; !ok {
			mdb.profiles[sysid] = &t
			mdb.log.Debugf("Adding profile for %s: %s %s", sysid, t.Device.Vendor, t.From)
		} else {
			mdb.log.Errorf("Profile for %s already exists. New file %s, Old file %s", sysid, t.From, op.From)
		}
	}

	return nil
}

func (mdb *MibDB) parseMibFromXml(file string) error {
	ap, err := apc.ParseApcMib(file, mdb.log)
	if err != nil {
		return err
	}

	profiles := newProfileFromApc(ap, file, mdb.log)
	for _, t := range profiles {
		for _, sysid := range t.Sysobjectid {
			if strings.HasPrefix(sysid, ".") {
				sysid = sysid[1:] // Strip this out if we start with .
			}
			if sysid == "" {
				mdb.log.Warnf("Skipping profile with no OID: %s: %s", sysid, t.Device.Vendor)
				continue
			}
			mdb.profiles[sysid] = t
			mdb.log.Debugf("Adding profile for [%s]: %s %d metrics and %d tags", sysid, t.Device.Vendor, len(t.Metrics), len(t.MetricTags))
		}
	}

	return nil
}

func newProfileFromApc(ap *apc.APC, file string, log logger.ContextL) []*Profile {
	profiles := make([]*Profile, 0, len(ap.Devices))
	for _, device := range ap.Devices {
		t := Profile{
			ContextL:    log,
			From:        file,
			Device:      Device{Vendor: device.Id},
			Sysobjectid: []string{},
			MetricTags:  []Tag{},
			Metrics:     []MIB{},
		}
		for _, oid := range device.OidMustExist {
			t.Sysobjectid = append(t.Sysobjectid, oid.Oid)
		}
		for _, data := range device.SetProductData {
			if data.Oid != "" {
				t.MetricTags = append(t.MetricTags, Tag{
					Column: OID{
						Oid:  data.Oid,
						Name: data.Field,
					},
					Tag:    data.RuleId,
					Symbol: data.Id,
				})
			}
		}
		for _, data := range device.SetLocationData {
			if data.Oid != "" {
				t.MetricTags = append(t.MetricTags, Tag{
					Column: OID{
						Oid:  data.Oid,
						Name: data.Field,
					},
					Tag:    data.RuleId,
					Symbol: data.Id,
				})
			}
		}

		mibSet := map[string]*MIB{}
		for _, sensor := range device.NumSensors {
			if _, ok := mibSet[sensor.SensorSet]; !ok {
				mibSet[sensor.SensorSet] = &MIB{
					//ForcedType: sensor.Type,
					Mib:     sensor.SensorSet,
					Symbols: []OID{},
				}
			}
			mib := mibSet[sensor.SensorSet]
			mib.Symbols = append(mib.Symbols, OID{
				Oid:  sensor.Value.Oid,
				Name: sensor.SensorId,
			})
		}

		for _, mib := range mibSet {
			t.Metrics = append(t.Metrics, *mib)
		}

		profiles = append(profiles, &t)
	}

	return profiles
}

// Sometimes the same tag can be in multiple mibs. Run down the list and keep the longest oid if there are any collisions.
func prune(mibs map[string]*kt.Mib) {
	seenNames := map[string]string{}

	for oid, mib := range mibs {
		oidName := mib.Name
		if mib.Tag != "" {
			oidName = mib.Tag
		}
		if _, ok := seenNames[oidName]; !ok { // We haven't seen this name yet so mark that we have it.
			seenNames[oidName] = oid
		} else {
			// There's a conflict. Keep the longest oid.
			if len(oid) > len(seenNames[oidName]) {
				seenNames[oidName] = oid
			}
		}
	}

	// Reverse this map to know who we're keeping.
	keepNames := map[string]bool{}
	for _, oid := range seenNames {
		keepNames[oid] = true
	}

	// Lastly, delete if we're not keeping.
	for oid, _ := range mibs {
		if !keepNames[oid] {
			delete(mibs, oid)
		}
	}
}
