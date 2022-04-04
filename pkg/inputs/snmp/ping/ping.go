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
	log    logger.ContextL
	target string
	pinger *ping.Pinger
	count  int
	priv   bool
}

func NewPinger(log logger.ContextL, target string, inter time.Duration) (*Pinger, error) {
	p := &Pinger{
		log:    log,
		target: target,
		count:  int(inter.Seconds()), // Run 1 ping per sec, for this many seconds.
	}

	if os.Getenv(KENTIK_PING_PRIV) == "true" {
		log.Infof("Running ping service in privileged mode.")
		p.priv = true
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

	// pinger.Interval = inter // Run at 1 per second.
	pinger.Count = p.count
	pinger.SetPrivileged(p.priv)
	p.pinger = pinger

	go func() {
		err := p.pinger.Run()
		if err != nil {
			p.log.Errorf("Cannot ping: %v", err)
		}
	}()

	return nil
}
