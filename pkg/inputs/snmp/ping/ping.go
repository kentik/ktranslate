package ping

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net"
	"net/netip"
	"os"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/inputs/snmp/ping/kaping"

	probing "github.com/prometheus-community/pro-bing"
	"gonum.org/v1/gonum/stat"
)

const (
	KENTIK_PING_PRIV = "KENTIK_PING_PRIV"
)

type Pinger struct {
	log      logger.ContextL
	target   string
	pinger   *kaping.Pinger
	priv     bool
	count    int
	timeout  time.Duration
	interval time.Duration
	citer    time.Duration
}

func NewPinger(log logger.ContextL, target string, pingSec int, timeout time.Duration) (*Pinger, error) {
	p := &Pinger{
		log:      log,
		target:   target,
		interval: time.Second * time.Duration(pingSec), // Send 1 ping every this many seconds.
		timeout:  timeout,
	}

	// Figure out who are we pinging.
	addrs, err := p.resolve()
	if err != nil {
		return nil, err
	}

	cfg := kaping.DefaultConfig()
	if os.Getenv(KENTIK_PING_PRIV) != "false" {
		log.Infof("Running ping service in privileged mode. Ping Interval: %v", p.interval)
		cfg.RawSocket = true
	} else {
		log.Infof("Running ping service in non privileged mode. Ping Interval: %v", p.interval)
		cfg.RawSocket = false
	}

	pinger, err := kaping.NewPinger(cfg, addrs)
	if err != nil {
		return nil, err
	}
	p.pinger = pinger
	return p, nil
}

func (p *Pinger) Start(ctx context.Context) {
	if p.pinger != nil {
		go p.pinger.Start(ctx)
	}
}

func (p *Pinger) Stop() {
	// Noop
}

func (p *Pinger) Reset(inter time.Duration, count int) error {
	if inter > 0 {
		p.citer = inter
	} else {
		p.citer = p.interval // Sent 1 packet every X seconds.
	}
	p.count = count

	p.log.Infof("Pinger reset to interval: %v, count: %v", p.citer, p.count)
	return nil
}

func (p *Pinger) Ping() (*probing.Statistics, error) {
	if p.pinger == nil {
		return nil, fmt.Errorf("pinger unavailable.")
	}

	addr, err := p.resolveOne()
	if err != nil {
		return nil, fmt.Errorf("host resolution failed: %w", err)
	}
	count := p.count
	iter := p.citer

	result, err := p.pinger.Ping(addr, count, iter, p.timeout)
	if err != nil {
		return nil, err
	}

	return statistics(p.target, addr, count, result), nil
}

func (p *Pinger) resolveOne() (netip.Addr, error) {
	addrs, err := p.resolve()
	if err != nil {
		return netip.Addr{}, err
	}

	return addrs[rand.Intn(len(addrs))], nil
}

func (p *Pinger) resolve() ([]netip.Addr, error) {
	addrs, err := net.LookupHost(p.target)
	switch {
	case err != nil:
		return nil, err
	case len(addrs) == 0:
		return nil, fmt.Errorf("address resolution failed")
	}

	res := make([]netip.Addr, len(addrs))
	for i, addr := range addrs {
		ipr, err := netip.ParseAddr(addr)
		if err != nil {
			return nil, err
		}
		res[i] = ipr
	}

	return res, nil
}

func statistics(host string, addr netip.Addr, count int, result *kaping.Result) *probing.Statistics {
	min := time.Duration(math.MaxInt64)
	max := time.Duration(0)
	rtt := make([]float64, len(result.RTT))

	for i, r := range result.RTT {
		if r < min {
			min = r
		}

		if r > max {
			max = r
		}

		rtt[i] = float64(r)
	}

	mean, stdev := stat.MeanStdDev(rtt, nil)

	return &probing.Statistics{
		PacketsRecv: count - result.Lost,
		PacketsSent: result.Sent,
		PacketLoss:  float64(result.Lost) / float64(result.Sent),
		IPAddr: &net.IPAddr{
			IP:   net.IP(addr.AsSlice()),
			Zone: addr.Zone(),
		},
		Addr:      host,
		Rtts:      result.RTT,
		MinRtt:    min,
		MaxRtt:    max,
		AvgRtt:    time.Duration(mean),
		StdDevRtt: time.Duration(stdev),
	}
}
