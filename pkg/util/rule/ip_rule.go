package rule

import (
	"fmt"
	"net"

	"github.com/kentik/patricia"
	tree "github.com/kentik/patricia/uint32_tree"
)

// compare two uint32s
func uint32Matches(payload uint32, val uint32) bool {
	return payload == val
}

// IPAddressRules manages rules associated with IPAddresses
type IPAddressRules struct {
	treeV4 *tree.TreeV4
	treeV6 *tree.TreeV6
}

// NewIPAddressRules returns a new IPAddressRules
func NewIPAddressRules() *IPAddressRules {
	return &IPAddressRules{
		treeV4: tree.NewTreeV4(),
		treeV6: tree.NewTreeV6(),
	}
}

// AddIPAddress adds a rule match by IP Address
func (r *IPAddressRules) AddIPAddress(ipAddress string, matchType Match) error {
	v4Addr, v6Addr, err := patricia.ParseIPFromString(ipAddress)
	if err != nil {
		return fmt.Errorf("Error parsing IP address: %s", err)
	}

	if v4Addr != nil {
		r.treeV4.Add(*v4Addr, matchType.Uint32(), uint32Matches)
	} else if v6Addr != nil {
		r.treeV6.Add(*v6Addr, matchType.Uint32(), uint32Matches)
	}

	return nil
}

// Check checks which Matches match the input IP
func (r *IPAddressRules) Check(ip net.IP) Match {
	ret := uint32(0)
	var found bool
	var err error

	if ipv4 := ip.To4(); ipv4 != nil {
		address := patricia.NewIPv4AddressFromBytes(ipv4, 32)
		found, ret, err = r.treeV4.FindDeepestTag(address)
	} else {
		address := patricia.NewIPv6Address(ip.To16(), 128)
		found, ret, err = r.treeV6.FindDeepestTag(address)
	}
	if err == nil && found {
		return NewMatch(ret)
	}

	return MatchNone
}
