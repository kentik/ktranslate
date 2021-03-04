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
	log      logger.ContextL
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

func NewMibDB(mibpath string, profileDir string, log logger.ContextL) (*MibDB, error) {
	log.Infof("Loading db from %s", mibpath)
	db, err := leveldb.OpenFile(mibpath, &opt.Options{})
	if err != nil {
		return nil, err
	}

	mdb := &MibDB{
		db:       db,
		log:      log,
		profiles: map[string]*Profile{},
	}

	num, err := mdb.LoadProfiles(profileDir)
	if err != nil {
		return nil, err
	}
	log.Infof("Loaded %d profiles from %s", num, profileDir)

	return mdb, nil
}

func (db *MibDB) Close() {
	db.db.Close()
}

func (db *MibDB) GetForOid(oid string) (map[string]*kt.Mib, error) {
	mibs := map[string]*kt.Mib{}
	iter := db.db.NewIterator(util.BytesPrefix([]byte(oid)), nil)
	for iter.Next() {
		pts := strings.SplitN(string(iter.Value()), " ", 2)
		if len(pts) >= 2 {
			res := reType.FindAllStringSubmatch(pts[1], -1)
			if len(res) > 0 {
				dt, err := strconv.Atoi(res[0][1])
				if err == nil {
					extra := strings.SplitN(pts[1], " ", 2)
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
		return nil, err
	}

	return mibs, nil
}

// Walk up the oid tree until we get something.
func (db *MibDB) GetForOidRecur(oid string) (map[string]*kt.Mib, error) {
	pts := strings.Split(oid, ".")
	for i := len(pts); i > 1; i-- {
		check := strings.Join(pts[0:i], ".")
		res, err := db.GetForOid(check)
		if err != nil {
			return nil, err
		}
		if len(res) > 0 {
			return res, nil
		}
	}

	return nil, nil
}
