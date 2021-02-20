package patricia

// string implementation of IPv4/IPv6 network patricia trees.
// Note: this is copy/paste/typed the same as PatriciaUint32 - any special handling should be in another .go file

import (
	"net"

	"github.com/kentik/golog/logger"
	"github.com/kentik/patricia"
	"github.com/kentik/patricia/string_tree"
)

// StringTrees is a IPv4/IPv6 pair of trees that hold string payloads
type StringTrees struct {
	Trees
	tree4 *string_tree.TreeV4
	tree6 *string_tree.TreeV6
}

// NewStringTrees returns a new StringTrees
func NewStringTrees(log *logger.Logger) (*StringTrees, error) {
	return &StringTrees{
		Trees: Trees{
			Length: 0,
			log:    log,
		},
		tree4: string_tree.NewTreeV4(),
		tree6: string_tree.NewTreeV6(),
	}, nil
}

// Set the value for a prefix - overwrites what's there (technically, the first in the list)
func (p *StringTrees) SetPrefix(addr string, content string) {
	v4Addr, v6Addr, err := patricia.ParseIPFromString(addr)
	if err != nil {
		p.log.Warnf(LOG_PREFIX, "Err0r parsing node with address %s: %s", addr, err)
		return
	} else {
		if v4Addr != nil {
			if countIncreased, _, err := p.tree4.Set(*v4Addr, content); err != nil {
				p.log.Errorf(LOG_PREFIX, "Error adding v4 node with address %s: %s", addr, err)
			} else {
				if countIncreased {
					p.Length++
				}
			}
		} else if v6Addr != nil {
			if countIncreased, _, err := p.tree6.Set(*v6Addr, content); err != nil {
				p.log.Errorf(LOG_PREFIX, "Error adding v6 node with address %s: %s", addr, err)
			} else {
				if countIncreased {
					p.Length++
				}
			}
		}
	}
}

func (p *StringTrees) AddPrefix(addr string, content string) {
	v4Addr, v6Addr, err := patricia.ParseIPFromString(addr)
	if err != nil {
		p.log.Warnf(LOG_PREFIX, "Err0r parsing node with address %s: %s", addr, err)
		return
	} else {
		if v4Addr != nil {
			if countIncreased, _, err := p.tree4.Add(*v4Addr, content, nil); err != nil {
				p.log.Errorf(LOG_PREFIX, "Error adding v4 node with address %s: %s", addr, err)
			} else {
				if countIncreased {
					p.Length++
				}
			}
		} else if v6Addr != nil {
			if countIncreased, _, err := p.tree6.Add(*v6Addr, content, nil); err != nil {
				p.log.Errorf(LOG_PREFIX, "Error adding v6 node with address %s: %s", addr, err)
			} else {
				if countIncreased {
					p.Length++
				}
			}
		}
	}
}

// Always removes content at prefix
// - returns whether the delete removed any prefixes
func (p *StringTrees) RemovePrefix(addr string, content string) bool {
	matchFunc := func(tagData string, val string) bool {
		return tagData == val
	}

	var err error
	addressRemoveCount := 0

	v4Addr, v6Addr, err := patricia.ParseIPFromString(addr)
	if err != nil {
		p.log.Warnf(LOG_PREFIX, "Err0r parsing node with address %s: %s", addr, err)
		return false
	}

	if v4Addr != nil {
		addressRemoveCount, err = p.tree4.Delete(*v4Addr, matchFunc, content)
		if err != nil {
			p.log.Errorf(LOG_PREFIX, "Error removing v4 tag with address %s - err:%s; #removed: %d", addr, err, addressRemoveCount)
		}
	} else if v6Addr != nil {
		addressRemoveCount, err = p.tree6.Delete(*v6Addr, matchFunc, content)
		if err != nil {
			p.log.Errorf(LOG_PREFIX, "Error removing v6 tag with address %s - err:%s; #removed: %d", addr, err, addressRemoveCount)
		}
	}

	p.Length -= addressRemoveCount
	return addressRemoveCount > 0
}

// Finds the most-specific tag with the input address (relying on tree traversing behavior)
func (p *StringTrees) FindBestMatch(addr uint32, addr6 []byte) (bool, string, error) {
	if addr > 0 {
		address := patricia.NewIPv4Address(addr, 32)
		return p.tree4.FindDeepestTag(address)
	} else {
		address := patricia.NewIPv6Address(addr6, 128)
		return p.tree6.FindDeepestTag(address)
	}
}

// Finds all tags with the input address (relying on tree traversing behavior)
func (p *StringTrees) FindAllMatches(addr net.IP, len uint8) ([]string, error) {
	if addr.To4() != nil {
		address := patricia.NewIPv4AddressFromBytes(addr.To4(), uint(len))
		return p.tree4.FindTags(address)
	} else {
		address6 := patricia.NewIPv6Address(addr.To16(), uint(len))
		return p.tree6.FindTags(address6)
	}
}
