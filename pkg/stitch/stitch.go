package stitch

import (
	"time"

	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/jellydator/ttlcache/v3"
)

type Stitcher struct {
	cache *ttlcache.Cache[string, *kt.JCHF]
}

func NewStitcher() (*Stitcher, error) {
	cache := ttlcache.New[string, *kt.JCHF](
		ttlcache.WithTTL[string, *kt.JCHF](1 * time.Minute),
	)

	go cache.Start() // starts automatic expired item deletion

	return &Stitcher{cache: cache}, nil
}

/*
If theres a matching ingress / egress flow, record it here.
*/
func (s *Stitcher) Stitch(msg *kt.JCHF) {
	key := msg.GetKey()
	item, retrieved := s.cache.GetOrSet(key, msg, ttlcache.WithTTL[string, *kt.JCHF](ttlcache.DefaultTTL))
	if retrieved {
		msg.Pair = item.Value()
	}
}
