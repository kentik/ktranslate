package rule

import (
	"net"
	"strings"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

var (
	INTERNAL_IPS = []string{
		"0.0.0.0/8",
		"127.0.0.0/8",
		"100.64.0.0/10",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"169.254.0.0/16",
		"192.0.0.0/24",
		"192.0.2.0/24",
		"192.18.0.0/15",
		"198.51.100.0/24",
		"203.0.113.0/24",
		"224.0.0.0/4",
		"192.88.99.0/24",
		"240.0.0.0/4",
		"fc00::/7",  // (ula)
		"fe80::/10", // (link local)
	}
)

// RuleSet holds a list of network classification rules
type RuleSet struct {
	log            logger.ContextL
	ipAddressRules *IPAddressRules
}

// NewRuleSet returns a new RuleSet
func NewRuleSet(log logger.ContextL) *RuleSet {
	rs := RuleSet{
		log:            log,
		ipAddressRules: NewIPAddressRules(),
	}

	for _, line := range INTERNAL_IPS {
		line = strings.TrimSpace(line)
		if err := rs.ipAddressRules.AddIPAddress(line, MatchPrivateIP); err != nil {
			continue // TODO: okay to keep going?
		}
	}

	return &rs
}

// check whether the AS matches against our static list of private ASNs
func matchesPrivateASN(as uint32) bool {
	if (as >= 64512 && as <= 65534) || (as >= 4200000000 && as <= 4294967294) {
		return true
	}

	return false
}

// IsInternal returns whether an IP/ASN is from static list of private IPs.
func (r *RuleSet) IsInternal(ip net.IP, as uint32) bool {
	var searchResult Match

	// IP address
	searchResult = r.ipAddressRules.Check(ip)
	if searchResult == MatchPrivateIP {
		return true
	}

	// private ASNs
	if matchesPrivateASN(as) {
		return true
	}

	return false
}
