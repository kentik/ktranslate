package kaping

import "time"

type Result struct {
	Sent int
	Lost int
	RTT  []time.Duration
}
