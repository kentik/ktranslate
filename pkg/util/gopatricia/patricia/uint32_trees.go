package patricia

// string implementation of IPv4/IPv6 network patricia trees.
// Note: this is copy/paste/typed the same as PatriciaString - any special handling should be in another .go file

/*
#include "go_patricia.h"
#include <sys/file.h>
*/
import "C"

import (
	"time"

	"github.com/kentik/golog/logger"
	"github.com/kentik/patricia"
	"github.com/kentik/patricia/uint32_tree"
)

// PatriciaUin32 is a IPv4/IPv6 pair of trees that hold uint32 payloads
type Uint32Trees struct {
	Trees
	tree4     *uint32_tree.TreeV4
	tree6     *uint32_tree.TreeV6
	file4     string
	file6     string
	loadTime4 time.Time
	loadTime6 time.Time
}

// NewUint32Trees returns a new Uint32Trees
func NewUint32Trees(log *logger.Logger, file4 string, file6 string) (*Uint32Trees, error) {
	return &Uint32Trees{
		Trees: Trees{
			Length: 0,
			log:    log,
		},
		tree4:     uint32_tree.NewTreeV4(),
		tree6:     uint32_tree.NewTreeV6(),
		file4:     file4,
		file6:     file6,
		loadTime4: time.Now(),
		loadTime6: time.Now(),
	}, nil
}

func (p *Uint32Trees) AddPrefix(addr string, content uint32) {
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
func (p *Uint32Trees) RemovePrefix(addr string, content uint32) bool {
	matchFunc := func(tagData uint32, val uint32) bool {
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
			p.log.Errorf(LOG_PREFIX, "Error removing v4 node with address %s: %s", addr, err)
		}
	} else if v6Addr != nil {
		addressRemoveCount, err = p.tree6.Delete(*v6Addr, matchFunc, content)
		if err != nil {
			p.log.Errorf(LOG_PREFIX, "Error removing v6 node with address %s: %s", addr, err)
		}
	}
	p.Length -= addressRemoveCount
	return addressRemoveCount > 0
}

// Finds the most-specific tag with the input address (relying on tree traversing behavior)
func (p *Uint32Trees) FindBestMatch(addr uint32, addr6 []byte) (bool, uint32, error) {
	if addr > 0 {
		address := patricia.NewIPv4Address(addr, 32)
		return p.tree4.FindDeepestTag(address)
	} else {
		address := patricia.NewIPv6Address(addr6, 128)
		return p.tree6.FindDeepestTag(address)
	}
}
