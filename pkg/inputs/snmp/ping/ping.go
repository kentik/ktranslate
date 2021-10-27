package ping

import (
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"

	"github.com/go-ping/ping"
)

type Pinger struct {
	log    logger.ContextL
	target string
	pinger *ping.Pinger
}

func NewPinger(log logger.ContextL, target string, inter time.Duration) (*Pinger, error) {
	p := &Pinger{
		log:    log,
		target: target,
	}
	pinger, err := ping.NewPinger(target)
	if err != nil {
		return nil, err
	}

	pinger.Interval = inter
	go pinger.Run()
	p.pinger = pinger

	return p, nil
}

func (p *Pinger) Statistics() *ping.Statistics {
	return p.pinger.Statistics()
}
