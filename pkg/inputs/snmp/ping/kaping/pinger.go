package kaping

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/netip"
	"time"

	"golang.org/x/net/icmp"
)

type Pinger struct {
	sock4 Socket
	sock6 Socket
	state *State
	error chan error
}

func NewPinger(cfg Config) (*Pinger, error) {
	mode := ICMP

	if cfg.RawSocket {
		mode = RAW
	}

	sock4, err := NewSocket(cfg.BindAddr4, mode)
	if err != nil {
		return nil, fmt.Errorf("IPv4 socket: %w", err)
	}

	sock6, err := NewSocket(cfg.BindAddr6, mode)
	if err != nil {
		return nil, fmt.Errorf("IPv6 socket: %w", err)
	}

	state := NewState()
	error := make(chan error)

	return &Pinger{sock4, sock6, state, error}, nil
}

func (p *Pinger) Start(ctx context.Context) {
	go p.receive(p.sock4)
	go p.receive(p.sock6)

	<-ctx.Done()

	p.sock4.Close()
	p.sock6.Close()
}

func (p *Pinger) Ping(addr netip.Addr, count int, delay, timeout time.Duration) (*Result, error) {
	result := &Result{}

	id := rand.Intn(65536)

	for i := 0; i < count; i++ {
		echo := icmp.Echo{
			ID:  id,
			Seq: i,
		}

		sent := time.Now()
		reply, err := p.Probe(addr, echo, timeout)
		result.Sent++

		switch {
		case err == nil:
			rtt := reply.Time.Sub(sent)
			result.RTT = append(result.RTT, rtt)
		case err == ErrTimeout:
			result.Lost++
		default:
			return nil, err
		}

		time.Sleep(delay)
	}

	return result, nil
}

func (p *Pinger) Probe(addr netip.Addr, echo icmp.Echo, timeout time.Duration) (Reply, error) {
	reply := make(chan Reply, 1)

	token := p.state.Insert(reply)
	defer p.state.Remove(token)

	echo.Data = token[:]

	err := p.send(addr, echo)
	if err != nil {
		return Reply{}, err
	}

	select {
	case reply := <-reply:
		return reply, nil
	case <-time.After(timeout):
		return Reply{}, ErrTimeout
	}
}

func (p *Pinger) Error() <-chan error {
	return p.error
}

func (p *Pinger) send(addr netip.Addr, echo icmp.Echo) error {
	if addr.Is4() {
		return p.sock4.Send(addr, echo)
	} else {
		return p.sock6.Send(addr, echo)
	}
}

func (p *Pinger) receive(sock Socket) {
	var (
		buffer [1500]byte
		token  [16]byte
	)

	for {
		msg, peer, err := sock.Recv(buffer[:])

		switch {
		case errors.Is(err, net.ErrClosed):
			return
		case err != nil:
			select {
			case p.error <- err:
			default:
			}
			continue
		}

		echo, ok := msg.Body.(*icmp.Echo)
		if !ok {
			continue
		}

		reply := Reply{
			Addr: peer,
			Echo: *echo,
			Time: time.Now(),
		}

		copy(token[:], echo.Data)

		if channel, ok := p.state.Lookup(token); ok {
			select {
			case channel <- reply:
			default:
			}
		}
	}
}

var ErrTimeout = errors.New("probe timeout")
