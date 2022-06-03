package ping

import (
	"os"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"

	"github.com/go-ping/ping"
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

func NewPinger(log logger.ContextL, target string, inter time.Duration, pingSec int) (*Pinger, error) {
	p := &Pinger{
		log:      log,
		target:   target,
		interval: time.Second * time.Duration(pingSec), // Send 1 ping every this many seconds.
	}

	if os.Getenv(KENTIK_PING_PRIV) == "true" {
		log.Infof("Running ping service in privileged mode. Ping Interval: %v", p.interval)
		p.priv = true
	} else {
		log.Infof("Running ping service in non privileged mode. Ping Interval: %v", p.interval)
	}

	err := p.Reset()
	return p, err
}

func (p *Pinger) Statistics() *ping.Statistics {
	return p.pinger.Statistics()
}

func (p *Pinger) Reset() error {
	pinger, err := ping.NewPinger(p.target)
	if err != nil {
		return err
	}

	pinger.Interval = p.interval // Sent 1 packet every X seconds. Default to 1.
	pinger.SetPrivileged(p.priv)
	pinger.OnFinish = func(stats *ping.Statistics) {
		p.log.Infof("Ping run %d finished.", p.num)
		p.num++
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
