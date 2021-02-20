package timing

// Inspired by honeycomb beeline's timer code (apache licensed)

import "time"

type Chrono struct {
	start time.Time
}

// Start a chrono now
func StartChrono() Chrono {
	return Chrono{
		start: time.Now(),
	}
}

// Finish a chrono and return the elapsed time in milliseconds
func (t Chrono) Finish() float64 {
	return float64(t.FinishDur()) / float64(time.Millisecond)
}

// Finish a chrono and return the elapsed time as a time.Duration
func (t Chrono) FinishDur() time.Duration {
	if t.start.IsZero() {
		return 0
	}
	return time.Since(t.start)
}

func (t Chrono) Duration() time.Duration {
	if t.start.IsZero() {
		return 0
	}
	return time.Since(t.start)
}
