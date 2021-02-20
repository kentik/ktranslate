package cat

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

const (
	CacheClearDuration = 60 * 60 * time.Second
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

	go res.clearCache(ctx)

	return res, nil
}

func (r *Resolver) Resolve(ctx context.Context, ip string) string {
	final, ok := r.cache[ip]
	if ok {
		return final
	}

	// Else, look it up on the network.
	if ans, err := r.resolver.LookupAddr(ctx, ip); err == nil {
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
