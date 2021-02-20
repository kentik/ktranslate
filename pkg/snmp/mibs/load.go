package mibs

import (
	"strings"

	"github.com/kentik/ktranslate/pkg/eggs/logger"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type MibDB struct {
	db  *leveldb.DB
	log logger.ContextL
}

type Mib struct {
	Oid  string
	Name string
	Type string
}

func NewMibDB(path string, log logger.ContextL) (*MibDB, error) {
	db, err := leveldb.OpenFile(path, &opt.Options{})
	if err != nil {
		return nil, err
	}

	log.Infof("Loading db from %s", path)
	return &MibDB{
		db:  db,
		log: log,
	}, nil
}

func (db *MibDB) Close() {
	db.db.Close()
}

func (db *MibDB) GetForOid(oid string) ([]Mib, error) {
	mibs := make([]Mib, 0)
	iter := db.db.NewIterator(util.BytesPrefix([]byte(oid)), nil)
	for iter.Next() {
		pts := strings.SplitN(string(iter.Value()), " ", 2)
		mibs = append(mibs, Mib{
			Oid:  string(iter.Key()),
			Name: pts[0],
			Type: pts[1],
		})
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return nil, err
	}

	return mibs, nil
}

// Walk up the oid tree until we get something.
func (db *MibDB) GetForOidRecur(oid string) ([]Mib, error) {
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
