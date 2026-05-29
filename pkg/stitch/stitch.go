package stitch

import (
	"flag"
	"time"

	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/jellydator/ttlcache/v3"
)

var (
	enable bool
	ttlSec int
)

func init() {
	flag.IntVar(&ttlSec, "stitch.ttl.sec", 30, "TTL for holding flows to try and stitch together.")
	flag.BoolVar(&enable, "stitch.enable", false, "Turn on flow stitching.")
}

type Stitcher struct {
	logger.ContextL
	cache *ttlcache.Cache[string, *kt.JCHF]
}

func NewStitcher(log logger.Underlying, cfg *ktranslate.StitchConfig) (*Stitcher, error) {
	if !cfg.Enable {
		return nil, nil
	}

	cache := ttlcache.New[string, *kt.JCHF](
		ttlcache.WithTTL[string, *kt.JCHF](time.Duration(cfg.TTLSec) * time.Second),
	)

	go cache.Start() // starts automatic expired item deletion

	s := &Stitcher{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "flowStitch"}, log),
		cache:    cache,
	}

	s.Infof("Starting flow unification system with a ttl of %v", time.Duration(cfg.TTLSec)*time.Second)
	return s, nil
}

/*
If theres a matching ingress / egress flow, record it here.
*/
func (s *Stitcher) Stitch(msg *kt.JCHF) bool {
	key := msg.GetKey()

	s.Infof("%s", key)

	item, retrieved := s.cache.GetOrSet(key, msg, ttlcache.WithTTL[string, *kt.JCHF](ttlcache.DefaultTTL))
	if retrieved {
		msg.Pair = item.Value()
		return true
	}

	return false
}
