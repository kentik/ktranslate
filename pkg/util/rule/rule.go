package rule

import (
	"io/ioutil"
	"net"
	"strings"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/util/service"
	"gopkg.in/yaml.v2"
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
	portRules      map[uint32]string
}

type CustomRule struct {
	Ports []uint32 `yaml:"ports"`
	Name  string   `yaml:"name"`
}

type CustomRuleSet struct {
	Applications []CustomRule `yaml:"applications"`
}

// NewRuleSet returns a new RuleSet
func NewRuleSet(appMap string, log logger.ContextL) (*RuleSet, error) {
	rs := RuleSet{
		log:            log,
		ipAddressRules: NewIPAddressRules(),
		portRules:      map[uint32]string{},
	}

	for _, line := range INTERNAL_IPS {
		line = strings.TrimSpace(line)
		if err := rs.ipAddressRules.AddIPAddress(line, MatchPrivateIP); err != nil {
			continue // TODO: okay to keep going?
		}
	}

	// If there's a custom set, get these here.
	if appMap != "" {
		customs := CustomRuleSet{}
		byc, err := ioutil.ReadFile(appMap)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(byc, &customs)
		if err != nil {
			return nil, err
		}

		log.Infof("Loaded %d custom rules.", len(customs.Applications))
		for _, custom := range customs.Applications {
			for _, port := range custom.Ports {
				rs.portRules[port] = custom.Name
			}
		}
	}

	return &rs, nil
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

// IP is for future proofing if there's demand.
func (r *RuleSet) GetService(ip net.IP, port uint32, protocol uint8) (string, bool) {
	// first see if we have a custom mapping here.
	if app, ok := r.portRules[port]; ok {
		return app, true
	}

	// If we don't, see if there's a global definition.
	if app, ok := service.Services[service.Port{Number: port, Protocol: protocol}]; ok {
		return app, true
	}

	// We couldn't find anything.
	return "", false
}
