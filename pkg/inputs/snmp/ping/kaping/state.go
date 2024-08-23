package kaping

import (
	"crypto/rand"
	"net/netip"
	"sync"
	"time"

	"golang.org/x/net/icmp"
)

type State struct {
	state map[[16]byte]chan Reply
	mutex sync.Mutex
}

type Reply struct {
	Addr netip.Addr
	Echo icmp.Echo
	Time time.Time
}

func NewState() *State {
	return &State{
		state: map[[16]byte]chan Reply{},
		mutex: sync.Mutex{},
	}
}

func (s *State) Insert(reply chan Reply) [16]byte {
	var token [16]byte
	rand.Read(token[:])

	s.mutex.Lock()
	s.state[token] = reply
	s.mutex.Unlock()

	return token
}

func (s *State) Lookup(token [16]byte) (chan<- Reply, bool) {
	s.mutex.Lock()
	reply, ok := s.state[token]
	s.mutex.Unlock()
	return reply, ok
}

func (s *State) Remove(token [16]byte) {
	s.mutex.Lock()
	delete(s.state, token)
	s.mutex.Unlock()
}
