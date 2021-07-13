// Copyright (C) 2014 CloudHelix.
// Use of this source code is governed by a GPLv3
// license that can be found in the LICENSE file.

package patricia

/*
#include "go_patricia.h"
#include <sys/file.h>
*/
import "C"

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/kentik/golog/logger"
	"github.com/kentik/patricia"
)

var (
	ErrorNotFound     = errors.New("Not found")
	ErrorNoEntries    = errors.New("No entries")
	ErrorNotIPv6      = errors.New("Not an IPv6 address")
	ErrorDuplicateKey = errors.New("Key already exists")
	LOG_PREFIX        = "patricia: "
)

// PatriciaError is used for errors using the sflow library. It implements the
// builtin error interface.
type PatriciaError string

type NodeGeo struct {
	node *C.geo_info_t
}

func (err PatriciaError) Error() string {
	return string(err)
}

func GetCountry(n *NodeGeo) uint32 {
	return uint32(C.get_country(n.node))
}

func GetRegion(n *NodeGeo) uint32 {
	return uint32(C.get_region(n.node))
}

func GetCity(n *NodeGeo) uint32 {
	return uint32(C.get_city(n.node))
}

// Utility to turn Geo -> packed int
func PackGeo(cnty []byte) uint32 {
	if len(cnty) == 2 {
		return uint32(rune((cnty)[0])*256 + rune((cnty)[1]))
	}
	return 0
}

// Utility to turn packed int -> Geo
func UnpackGeo(pked uint32) *string {
	vl := fmt.Sprintf("%c%c", pked>>8, pked&0xFF)
	return &vl
}

func OpenGeo(file string, force bool, log *logger.Logger) (*GeoTrees, error) {

	// If the geo file is off, go ahead and return null at this point.
	if file == "off" {
		return nil, nil
	}

	p, err := NewGeoTrees(log)
	if err != nil {
		return nil, err
	}

	err = p.Load(file)
	if err != nil {
		return nil, err
	}

	go p.Update(file)

	return p, nil
}

func OpenASN(file4 string, file6 string, fileName string, log *logger.Logger) (*Uint32Trees, error) {
	p, err := NewUint32Trees(log, file4, file6)
	if err != nil {
		return nil, err
	}

	loadNames := func(fileName string) (int, error) {
		found := 0
		file, err := os.Open(fileName)
		if err != nil {
			return found, err
		}
		defer file.Close()

		mapr := map[uint32]string{}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			pts := strings.Fields(line)
			if len(pts) > 0 {
				num, err := strconv.Atoi(pts[0][2:])
				if err != nil {
					return found, fmt.Errorf("Error adding name with line %s -> %v", line, err)
				}
				more := strings.Split(pts[1], ",")
				mapr[uint32(num)] = strings.TrimSpace(more[0])
				found++
			}
		}

		if err := scanner.Err(); err != nil {
			return found, err
		}

		p.mapr = mapr
		return found, nil
	}

	loadUp := func(fileName string) (int, error) {
		found := 0
		file, err := os.Open(fileName)
		if err != nil {
			return found, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			pts := strings.Split(line, ",")
			if len(pts) > 0 {
				v4Addr, v6Addr, err := patricia.ParseIPFromString(pts[1])
				if err != nil {
					return found, err
				} else {
					asnNum, err := strconv.Atoi(pts[0])
					if err != nil {
						return found, err
					}

					if v4Addr != nil {
						if countIncreased, _, err := p.tree4.Add(*v4Addr, uint32(asnNum), nil); err != nil {
							return found, fmt.Errorf("Error adding v4 node with address %s: %s", line, err)
						} else {
							if countIncreased {
								found++
							}
						}
					} else if v6Addr != nil {
						if countIncreased, _, err := p.tree6.Add(*v6Addr, uint32(asnNum), nil); err != nil {
							return found, fmt.Errorf("Error adding v6 node with address %s: %s", line, err)
						} else {
							if countIncreased {
								found++
							}
						}
					}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			return found, err
		}

		return found, nil
	}

	found4, err := loadUp(file4)
	if err != nil {
		return nil, err
	}

	found6, err := loadUp(file6)
	if err != nil {
		return nil, err
	}

	_, err = loadNames(fileName)
	if err != nil {
		return nil, err
	}

	p.Length = found4 + found6

	return p, nil
}

func PackIPAddr(ip net.IP) *[16]byte {
	var addr [16]byte
	if ipv4 := ip.To4(); ipv4 != nil {
		addr[10] = 0xff
		addr[11] = 0xff
		copy(addr[12:], ipv4)
	} else {
		copy(addr[:], ip)
	}
	return &addr
}
