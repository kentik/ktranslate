package aws

import (
	"math/rand"
	"net"
	"time"

	"github.com/json-iterator/go"
)

// fast JSON encoding
var json = jsoniter.ConfigFastest

// jsonSorted is fast, but still sorts keys
var jsonSorted = jsoniter.Config{
	IndentionStep:          4,
	EscapeHTML:             false,
	SortMapKeys:            true,
	ValidateJsonRawMessage: false,
}.Froze()

var (
	internalIPs []*net.IPNet
)

func init() {
	internalIPs = []*net.IPNet{}
	for _, cidr := range []string{
		"10.0.0.0/8",     // RFC1918
		"192.168.0.0/16", // RFC1918
		"172.16.0.0/12",  // RFC1918
		"127.0.0.0/8",    // IPv4 loopback
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
	} {
		_, block, _ := net.ParseCIDR(cidr)
		internalIPs = append(internalIPs, block)
	}

	rand.Seed(time.Now().UnixNano())
}

func isPrivateIP(ip net.IP) bool {
	for _, block := range internalIPs {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}
