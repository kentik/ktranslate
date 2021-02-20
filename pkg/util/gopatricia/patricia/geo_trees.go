package patricia

// string implementation of IPv4/IPv6 network patricia trees.
// Note: this is copy/paste/typed the same as PatriciaString - any special handling should be in another .go file

/*
#include "go_patricia.h"
#include <sys/file.h>
*/
import "C"

import (
	"github.com/kentik/ktranslate/pkg/util/fetch"

	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/kentik/golog/logger"
	"github.com/kentik/ktranslate/pkg/util/flock"
)

const (
	FILE_CHECK_TIME = time.Hour
)

// Geo tree only holds geo data
type GeoTrees struct {
	Trees
	packed unsafe.Pointer // *C.packed_geo_t
}

func NewGeoTrees(log *logger.Logger) (*GeoTrees, error) {
	return &GeoTrees{
		Trees: Trees{
			Length: 0,
			log:    log,
		},
	}, nil
}

// Run in a loop, update when the base file changes
func (p *GeoTrees) Update(file string) {

	// read current file's time
	finfo, err := os.Stat(file)
	if err != nil {
		p.log.Errorf(LOG_PREFIX, "Error stating file %s %v", file, err)
		return
	}
	loadTime := finfo.ModTime()

	checker := time.NewTicker(FILE_CHECK_TIME)
	defer checker.Stop()

	p.log.Infof(LOG_PREFIX, "Update starting on %s, checking every %v", file, FILE_CHECK_TIME)
	for {
		select {
		case _ = <-checker.C:
			finfo, err := os.Stat(file)
			if err != nil {
				p.log.Errorf(LOG_PREFIX, "Error stating file %s %v", file, err)
				continue
			}

			if finfo.ModTime() == loadTime {
				continue // same time, don't update
			}

			p.log.Infof(LOG_PREFIX, "Reloading %s", file)
			if err := p.reload(file); err != nil {
				p.log.Errorf(LOG_PREFIX, "Error reloading file %s %v", file, err)
			} else {
				loadTime = finfo.ModTime() // success
			}
		}
	}
}

func (p *GeoTrees) Load(file string) error {
	return p.reload(file)
}

func (p *GeoTrees) setPackedGeo(packedGeo *C.packed_geo_t) {
	newPtr := unsafe.Pointer(packedGeo)
	oldPtr := atomic.SwapPointer(&p.packed, newPtr)
	if oldPtr != nil {
		// To avoid freeing our old copy of the data while we're still using it, do it async after a short while
		go func() {
			time.Sleep(10 * time.Second)
			C.close_packed_geo((*C.struct__packed_geo_t)(oldPtr))
		}()
	}
}

func (p *GeoTrees) loadPacked(file string) bool {
	geoFileMmap := C.CString(file + fetch.MMAP_SUFFIX)
	defer C.free(unsafe.Pointer(geoFileMmap))
	packedGeo := C.new_packed_geo_from_file(geoFileMmap)
	if packedGeo == nil {
		return false
	}
	p.Length = int(packedGeo.length)
	p.setPackedGeo(packedGeo)
	return true
}

func (p *GeoTrees) csvToPacked(file string) error {
	fd, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fd.Close()

	ips := make([]*C.geo_info_t, 0)
	defer func() {
		for _, entry := range ips {
			C.free(unsafe.Pointer(entry))
		}
	}()

	// parse csv to memory
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		pts := strings.Split(scanner.Text(), ",")
		if len(pts) > 3 {
			ipr := pts[0]
			cntd, _ := strconv.Atoi(pts[1])
			citd, _ := strconv.Atoi(pts[2])
			regd, _ := strconv.Atoi(pts[3])

			if ip, _, err := net.ParseCIDR(ipr); err == nil {
				b := PackIPAddr(ip)
				ips = append(ips, C.new_geo((*C.uint8_t)(&b[0]), C.uint16_t(cntd), C.uint16_t(regd), C.uint32_t(citd)))
			}
		} else {
			return fmt.Errorf("invalid Geo Line: %s", scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if len(ips) == 0 {
		return PatriciaError(fmt.Sprintf("LoadGeo() failed: %s, %d", file, p.Length))
	}

	// save to packed file
	cGeo := (**C.geo_info_t)(unsafe.Pointer(&ips[0]))
	geoFileMmap := C.CString(file + fetch.MMAP_SUFFIX)
	defer C.free(unsafe.Pointer(geoFileMmap))
	C.save_packed_geo(cGeo, C.size_t(len(ips)), geoFileMmap)
	return nil
}

func (p *GeoTrees) reload(file string) error {
	// If the geo file is off, go ahead and return null at this point.
	if file == "off" {
		return nil
	}

	st := time.Now()

	return flock.WithLock(file+".lock", func() error {
		// Try loading packed version via mmap
		if p.loadPacked(file) {
			p.log.Infof(LOG_PREFIX, "(Re)-loaded %d GEO entries via mmap in %v", p.Length, time.Now().Sub(st))
			return nil
		}

		// No mmap, load things up line by line and build the mmapable file
		if err := p.csvToPacked(file); err != nil {
			p.log.Infof(LOG_PREFIX, "Could not parse csv %v: %v", file, err)
			return err
		}

		// Now load it again
		if !p.loadPacked(file) {
			return fmt.Errorf("could not load newly created packed file")
		}

		p.log.Infof(LOG_PREFIX, "(Re)-loaded %d GEO entries from CSV in %v", p.Length, time.Now().Sub(st))
		return nil
	})
}

func (p *GeoTrees) SearchBestFromHostGeo(ip net.IP) (*NodeGeo, error) {
	b := PackIPAddr(ip)
	packed := (*C.packed_geo_t)(p.packed)
	n := &NodeGeo{node: C.find_in_packed(packed, (*C.uint8_t)(&b[0]))}
	if n.node == nil {
		return nil, PatriciaError("No result")
	} else {
		return n, nil
	}
}
