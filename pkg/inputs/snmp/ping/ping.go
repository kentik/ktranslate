package ping

import (
	"os"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"

	ping "github.com/prometheus-community/pro-bing"
)

const (
	KENTIK_PING_PRIV = "KENTIK_PING_PRIV"
)

type Pinger struct {
	log      logger.ContextL
	target   string
	pinger   *ping.Pinger
	priv     bool
	num      int
	interval time.Duration
}

func NewPinger(log logger.ContextL, target string, pingSec int) (*Pinger, error) {
	p := &Pinger{
		log:      log,
		target:   target,
		interval: time.Second * time.Duration(pingSec), // Send 1 ping every this many seconds.
	}

	if os.Getenv(KENTIK_PING_PRIV) != "false" {
		log.Infof("Running ping service in privileged mode. Ping Interval: %v", p.interval)
		p.priv = true
	} else {
		log.Infof("Running ping service in non privileged mode. Ping Interval: %v", p.interval)
	}

	err := p.Reset(p.interval)
	return p, err
}

func (p *Pinger) Statistics() *ping.Statistics {
	return p.pinger.Statistics()
}

func (p *Pinger) Reset(inter time.Duration) error {
	pinger, err := ping.NewPinger(p.target)
	if err != nil {
		return err
	}

	if inter > 0 {
		pinger.Interval = inter
	} else {
		pinger.Interval = p.interval // Sent 1 packet every X seconds.
	}
	pinger.SetPrivileged(p.priv)
	pinger.OnFinish = func(stats *ping.Statistics) {
		p.log.Debugf("Ping run %d finished.", p.num)
		p.num++
	}

	if p.pinger != nil {
		p.pinger.Stop()
	}
	p.pinger = pinger
	go func() {
		err := p.pinger.Run()
		if err != nil {
			p.log.Errorf("Cannot ping: %v", err)
		}
	}()

	return nil
}

func (p *Pinger) Stop() {
	p.pinger.Stop()
}
