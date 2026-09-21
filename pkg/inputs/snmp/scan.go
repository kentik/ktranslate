package snmp

import (
	"context"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/google/gopacket/macs"
	"github.com/mostlygeek/arp"
)

// netScanResult describes what was learned about one IP during a network scan.
type netScanResult struct {
	Host         net.IP
	MAC          string
	Manufacturer string
	Name         string
	Latency      time.Duration
}

func (r netScanResult) IsHostUp() bool {
	return r.Latency > -1
}

// scanCIDR probes every host in cidr: an ARP-table MAC lookup, a MAC-prefix
// vendor lookup, a reverse DNS lookup, and a TCP dial to port 1 to gauge
// liveness/latency. One goroutine per IP.
func scanCIDR(ctx context.Context, cidr string, timeout time.Duration) ([]netScanResult, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	ip = ip.Mask(ipnet.Mask)

	var wg sync.WaitGroup
	var mux sync.Mutex
	results := []netScanResult{}

	for ; ipnet.Contains(ip); incrementIP(ip) {
		select {
		case <-ctx.Done():
			wg.Wait()
			return results, ctx.Err()
		default:
		}

		target := make(net.IP, len(ip))
		copy(target, ip)

		wg.Add(1)
		go func(target net.IP) {
			defer wg.Done()
			r := probeHost(target, timeout)
			mux.Lock()
			results = append(results, r)
			mux.Unlock()
		}(target)
	}
	wg.Wait()

	return results, nil
}

func probeHost(ip net.IP, timeout time.Duration) netScanResult {
	r := netScanResult{Host: ip, Latency: -1}

	if macStr := arp.Search(ip.String()); macStr != "" && macStr != "00:00:00:00:00:00" {
		if mac, err := net.ParseMAC(macStr); err == nil {
			r.MAC = mac.String()

			prefix := [3]byte{mac[0], mac[1], mac[2]}
			if manufacturer, ok := macs.ValidMACPrefixMap[prefix]; ok {
				r.Manufacturer = manufacturer
			}

			if addrs, err := net.LookupAddr(ip.String()); err == nil && len(addrs) > 0 {
				r.Name = addrs[0]
			}
		}
	}

	start := time.Now()
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip.String(), "1"), timeout)
	if err != nil {
		if !strings.Contains(err.Error(), "timeout") {
			r.Latency = time.Since(start)
		}
	} else {
		r.Latency = time.Since(start)
		conn.Close()
	}

	return r
}

// incrementIP walks ip to the next address in place, matching furious's
// TargetIterator semantics (wraps within the byte slice's own length).
func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
