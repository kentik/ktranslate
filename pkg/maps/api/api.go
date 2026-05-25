package api

import (
	"bufio"
	"context"
	"os"
	"strings"
	"sync"
	"time"

	kkapi "github.com/kentik/ktranslate/pkg/api"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
)

type ApiTagMapper struct {
	sync.RWMutex
	logger.ContextL
	tags     map[uint32]string
	apic     *kkapi.KentikApi
	check    chan uint32
	searched map[uint32]bool
	colNames map[string]string
}

const (
	CHAN_SLACK        = 10000
	MAX_LOOKUP_SET    = 2000
	LOOKUP_CHECK_TIME = 30 * time.Second

	KT_COL_NAME_MAP_FILE = "KT_COL_NAME_MAP_FILE"
)

func NewApiTagMapper(log logger.Underlying, apic *kkapi.KentikApi) (*ApiTagMapper, error) {
	atm := ApiTagMapper{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "apiMapper"}, log),
		tags:     map[uint32]string{},
		apic:     apic,
		check:    make(chan uint32, CHAN_SLACK),
		searched: map[uint32]bool{},
		colNames: map[string]string{},
	}

	if colNameFile, ok := os.LookupEnv(KT_COL_NAME_MAP_FILE); ok {
		f, err := os.Open(colNameFile)
		if err != nil {
			return nil, err
		}
		defer func() {
			f.Close()
		}()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			pts := strings.SplitN(scanner.Text(), ",", 2)
			atm.colNames[kt.FixupName(pts[0])] = kt.FixupName(pts[1])
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		atm.Infof("Loaded %d tag column names", len(atm.colNames))
	}

	return &atm, nil
}

func (atm *ApiTagMapper) Run(ctx context.Context) {
	if atm != nil {
		go atm.startCheckService(ctx)
	}
}

func (atm *ApiTagMapper) LookupKV(k uint32) string {
	atm.RLock()
	defer atm.RUnlock()
	return atm.tags[k]
}

func (atm *ApiTagMapper) LookupTagValue(cid kt.Cid, tagval uint32, colname string) (string, string, bool) {

	fullColName := colname
	if cn, ok := atm.colNames[colname]; ok {
		fullColName = cn
	} else { // We don't know about this column so gonna skip for now.
		return colname, "", false
	}

	if tagval == 0 { // 0 is a null value here.
		return fullColName, "", false
	}

	atm.RLock()
	defer atm.RUnlock()
	if v, ok := atm.tags[tagval]; ok {
		return fullColName, v, ok
	}
	if atm.searched[tagval] { // Only hit api one time.
		return fullColName, "", false
	}

	// We don't know about this one. Add to the q to check.
	select {
	case atm.check <- tagval:
	default:
		atm.Debugf("Lookup channel full %d", len(atm.check))
	}

	return fullColName, "", false
}

func (atm *ApiTagMapper) LookupTagValueBig(cid kt.Cid, tagval int64, colname string) (string, string, bool) {
	return atm.LookupTagValue(cid, uint32(tagval), colname)
}

func (atm *ApiTagMapper) startCheckService(ctx context.Context) {
	lookupCheck := time.NewTicker(LOOKUP_CHECK_TIME)
	lookups := map[uint32]bool{}
	atm.Infof("Starting lookup loop")

	for {
		select {
		case _ = <-lookupCheck.C:
			go atm.doLookup(ctx, lookups)
			lookups = map[uint32]bool{}
		case v := <-atm.check:
			lookups[v] = true
			if len(lookups) >= MAX_LOOKUP_SET {
				go atm.doLookup(ctx, lookups)
				lookups = map[uint32]bool{}
			}
		case <-ctx.Done():
			atm.Infof("Lookup loop done")
			lookupCheck.Stop()
			return
		}
	}
}

func (atm *ApiTagMapper) doLookup(ctx context.Context, lookups map[uint32]bool) {

	if len(lookups) == 0 {
		return
	}

	atm.Debugf("Doing lookup check with %d lookups", len(lookups))
	vals, err := atm.apic.LookupEnumerationValues(ctx, lookups)
	if err != nil { // Error case, remove the seen marks for ones here since we don't know it yet.
		atm.Errorf("Error looking up tag enums: %v", err)
		return
	}

	// Record everything we searched for and also the actual values.
	atm.Lock()
	for v, _ := range lookups {
		atm.searched[v] = true
	}
	for k, v := range vals {
		atm.tags[k] = v
	}
	atm.Debugf("Finished lookup check with %d lookups, %d searched and %d tags", len(lookups), len(atm.searched), len(atm.tags))
	atm.Unlock()
}
