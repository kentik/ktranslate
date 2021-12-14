package cat

import (
	"context"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

const (
	CacheClearDuration  = 60 * 60 * time.Second
	MAX_CACHE_SIZE_NAME = "KentikMaxDnsCacheSize"
	RESOLVE_TIME_MAX    = time.Duration(time.Millisecond * 80)
)

var (
	MAX_CACHE_SIZE = 10000 // Cache up to 10K ips.
)

type Resolver struct {
	logger.ContextL
	resolver *net.Resolver
	cache    map[string]string
}

func NewResolver(ctx context.Context, log logger.Underlying, dsnHost string) (*Resolver, error) {
	res := &Resolver{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "resolver"}, log),
		resolver: &net.Resolver{
			PreferGo:     true,
			StrictErrors: false,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				if dead, ok := ctx.Deadline(); ok {
					return net.DialTimeout(network, dsnHost, dead.Sub(time.Now()))
				} else {
					return net.DialTimeout(network, dsnHost, 1*time.Second)
				}
			},
		},
		cache: map[string]string{},
	}

	if val, ok := os.LookupEnv(MAX_CACHE_SIZE_NAME); ok {
		if ival, err := strconv.Atoi(val); err == nil {
			MAX_CACHE_SIZE = ival
		}
	}
	res.Infof("Running dns with a cache size of %d ips.", MAX_CACHE_SIZE)

	go res.clearCache(ctx)

	return res, nil
}

func (r *Resolver) Resolve(ctx context.Context, ip string) string {
	final, ok := r.cache[ip]
	if ok {
		return final
	}

	// Else, look it up on the network, unless we are full.
	if len(r.cache) > MAX_CACHE_SIZE {
		return ""
	}

	// Cap the time we spend searching for an answer here.
	ctxC, cancel := context.WithTimeout(ctx, RESOLVE_TIME_MAX)
	defer cancel()

	if ans, err := r.resolver.LookupAddr(ctxC, ip); err == nil {
		if len(ans) > 0 {
			final = ans[0]
			if final[len(final)-1:] == "." {
				final = strings.ToLower(final[0 : len(final)-1]) // Remove trailing . and also make lower case (because I always forget this)
			}
		}
	} // ignore errors here
	r.cache[ip] = final // cache.

	return final
}

func (r *Resolver) clearCache(ctx context.Context) {
	clearTicker := time.NewTicker(CacheClearDuration)
	defer func() {
		clearTicker.Stop()
	}()

	for {
		select {
		case <-clearTicker.C:
			r.cache = map[string]string{}

		case <-ctx.Done():
			return
		}
	}
}
