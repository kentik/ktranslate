package api

import (
	"context"
	"sync"
	"time"

	kkapi "github.com/kentik/ktranslate/pkg/api"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
)

type ApiTagMapper struct {
	sync.RWMutex
	logger.ContextL
	tags  map[uint32]string
	apic  *kkapi.KentikApi
	check chan uint32
}

const (
	CHAN_SLACK        = 10000
	MAX_LOOKUP_SET    = 1000
	LOOKUP_CHECK_TIME = 30 * time.Second
)

func NewApiTagMapper(log logger.Underlying, apic *kkapi.KentikApi) (*ApiTagMapper, error) {
	atm := ApiTagMapper{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "apiMapper"}, log),
		tags:     map[uint32]string{},
		apic:     apic,
		check:    make(chan uint32, CHAN_SLACK),
	}

	return &atm, nil
}

func (atm *ApiTagMapper) Run(ctx context.Context) {
	go atm.startCheckService(ctx)
}

func (atm *ApiTagMapper) LookupKV(k uint32) string {
	atm.RLock()
	defer atm.RUnlock()
	return atm.tags[k]
}

func (atm *ApiTagMapper) LookupTagValue(cid kt.Cid, tagval uint32, colname string) (string, string, bool) {
	atm.RLock()
	defer atm.RUnlock()
	if v, ok := atm.tags[tagval]; ok {
		return colname, v, ok
	}
	// We don't know about this one. Add to the q to check.
	atm.check <- tagval
	return "", "", false
}

func (atm *ApiTagMapper) LookupTagValueBig(cid kt.Cid, tagval int64, colname string) (string, string, bool) {
	return atm.LookupTagValue(cid, uint32(tagval), colname)
}

func (atm *ApiTagMapper) startCheckService(ctx context.Context) {
	lookupCheck := time.NewTicker(LOOKUP_CHECK_TIME)
	lookups := make([]uint32, 0, MAX_LOOKUP_SET)

	for {
		select {
		case _ = <-lookupCheck.C:
			go atm.doLookup(ctx, lookups)
			lookups = make([]uint32, 0, MAX_LOOKUP_SET)
		case v := <-atm.check:
			lookups = append(lookups, v)
			if len(lookups) >= MAX_LOOKUP_SET {
				go atm.doLookup(ctx, lookups)
				lookups = make([]uint32, 0, MAX_LOOKUP_SET)
			}
		case <-ctx.Done():
			atm.Infof("Lookup loop done")
			lookupCheck.Stop()
			return
		}
	}
}

func (atm *ApiTagMapper) doLookup(ctx context.Context, lookups []uint32) {
	if len(lookups) == 0 {
		return
	}

	vals, err := atm.apic.LookupEnumerationValues(ctx, lookups)
	if err != nil {
		atm.Errorf("Error looking up tag enums: %v", err)
	}

	atm.Lock()
	defer atm.Unlock()
	for k, v := range vals {
		atm.tags[k] = v
	}
}
