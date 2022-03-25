package mibs

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type MibDB struct {
	db       *leveldb.DB
	profiles map[string]*Profile
	trapMibs map[string]*kt.Mib
	traps    map[string]Trap
	log      logger.ContextL
	validate bool // Error out if a profile is not correct.
}

var (
	reType = regexp.MustCompile(`type=(\d+)`)

	validOids = map[kt.Oidtype]bool{ // Hack
		kt.ObjID:     true,
		kt.String:    true,
		kt.INTEGER:   true,
		kt.NetAddr:   true,
		kt.IpAddr:    true,
		kt.Counter:   true,
		kt.Gauge:     true,
		kt.TimeTicks: true,
		kt.Counter64: true,
		kt.BitString: true,
		kt.Index:     true,
		kt.Integer32: true,
	}
)

func NewMibDB(mibpath string, profileDir string, validate bool, log logger.ContextL) (*MibDB, error) {
	mdb := &MibDB{
		log:      log,
		profiles: map[string]*Profile{},
		traps:    map[string]Trap{},
		trapMibs: map[string]*kt.Mib{},
		validate: validate,
	}

	if mibpath != "" {
		log.Infof("Loading db from %s", mibpath)
		db, err := leveldb.OpenFile(mibpath, &opt.Options{ReadOnly: true})
		if err != nil {
			return nil, err
		}
		mdb.db = db
	}

	if profileDir != "" {
		num, err := mdb.LoadProfiles(profileDir)
		if err != nil {
			return nil, err
		}
		log.Infof("Loaded %d profiles from %s", num, profileDir)
	}

	return mdb, nil
}

func (db *MibDB) Close() {
	if db.db != nil {
		db.db.Close()
	}
}

func (db *MibDB) GetTrap(oid string) *Trap {
	if t, ok := db.traps[oid]; ok {
		return &t
	}
	return nil
}

func (db *MibDB) GetForKey(oid string) (*kt.Mib, error) {
	if res, ok := db.trapMibs[oid]; ok {
		return res, nil
	}

	if db.db == nil { // We might not have set up a db here.
		return nil, nil
	}
	data, err := db.db.Get([]byte(oid), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	pts := strings.SplitN(string(data), " ", 2)
	if len(pts) >= 2 {
		res := reType.FindAllStringSubmatch(pts[1], -1)
		if len(res) > 0 {
			dt, err := strconv.Atoi(res[0][1])
			if err == nil {
				return &kt.Mib{
					Oid:  oid,
					Name: strings.SplitN(pts[0], "(", 2)[0],
					Type: kt.Oidtype(dt),
				}, nil
			}
		}
	}

	return nil, nil
}

func (db *MibDB) GetForOid(oid string, profile string, description string) (map[string]*kt.Mib, kt.Provider, error) {

	getProvider := func(moid string) (kt.Provider, bool) {
		if prov, ok := db.checkForProvider(moid, profile, description); ok {
			db.log.Infof("Provider: %s %s %s %s -> %s", oid, moid, profile, description, prov)
			return prov, true
		}
		return "", false
	}

	if db.db == nil { // We might not have set up a db here.
		if prov, ok := getProvider(oid); ok {
			return nil, prov, nil
		}

		return nil, "", nil
	}
	mibs := map[string]*kt.Mib{}
	iter := db.db.NewIterator(util.BytesPrefix([]byte(oid)), nil)
	provider := kt.ProviderDefault
	foundProv := false
	for iter.Next() {
		pts := strings.SplitN(string(iter.Value()), " ", 2)
		if len(pts) >= 2 {
			res := reType.FindAllStringSubmatch(pts[1], -1)
			if len(res) > 0 {
				dt, err := strconv.Atoi(res[0][1])
				if err == nil {
					extra := strings.SplitN(pts[1], " ", 2)
					if !foundProv {
						if prov, ok := getProvider(pts[0]); ok {
							provider = prov
							foundProv = true
						}
					}
					if validOids[kt.Oidtype(dt)] {
						mb := kt.Mib{
							Oid:  string(iter.Key()),
							Name: strings.SplitN(pts[0], "(", 2)[0],
							Type: kt.Oidtype(dt),
						}
						if len(extra) > 1 {
							mb.Extra = extra[1]
						}
						mibs[mb.Oid] = &mb
					}
				}
			}
		}
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return nil, "", err
	}

	return mibs, provider, nil
}

// Walk up the oid tree until we get something.
func (db *MibDB) GetForOidRecur(oid string, profile string, description string) (map[string]*kt.Mib, kt.Provider, bool, error) {
	pts := strings.Split(oid, ".")
	for i := len(pts); i > 1; i-- {
		check := strings.Join(pts[0:i], ".")
		res, pro, err := db.GetForOid(check, profile, description)
		if err != nil {
			return nil, "", (i == len(pts)), err
		}
		if len(res) > 0 || (db.db == nil && pro != "") {
			return res, pro, (i == len(pts)), nil
		}
	}

	return nil, kt.ProviderDefault, false, nil
}

func (db *MibDB) checkForProvider(name string, profile string, description string) (kt.Provider, bool) {
	// Check for some common patterns, see if we can guess what provider this oid is for.
	name = strings.ToLower(name)
	description = strings.ToLower(description)
	profile = strings.ToLower(profile)

	combo := name + "^" + description
	if strings.Contains(profile, "base.yml") {
		return kt.ProviderDefault, true
	}
	if !strings.Contains(profile, "cisco-catalyst") && (strings.Contains(combo, "router") || strings.Contains(combo, "ios xr") ||
		strings.Contains(combo, "freebsd") || strings.Contains(profile, "cisco-asr")) {
		return kt.ProviderRouter, true
	}
	if strings.Contains(combo, "switch") || strings.Contains(profile, "cisco-catalyst") {
		return kt.ProviderSwitch, true
	}
	if strings.Contains(combo, "firewall") {
		return kt.ProviderFirewall, true
	}
	if strings.Contains(profile, "ups") {
		return kt.ProviderUPS, true
	}
	if strings.Contains(profile, "pdu") {
		return kt.ProviderPDU, true
	}
	if strings.Contains(profile, "fc-switch") {
		return kt.ProviderFibreChannel, true
	}
	if strings.Contains(combo, "iot") {
		return kt.ProviderIOT, true
	}
	if strings.Contains(combo, "printer") {
		return kt.ProviderIOT, true
	}
	if strings.Contains(combo, "nas") {
		return kt.ProviderNas, true
	}
	if strings.Contains(combo, "san") {
		return kt.ProviderSan, true
	}
	if strings.Contains(combo, "wireless") {
		return kt.ProviderWirelessController, true
	}

	return "", false
}
