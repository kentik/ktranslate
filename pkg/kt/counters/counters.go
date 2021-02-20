package counters

type (
	// CounterSet represents an interface found during counter push
	// The keys to the Values map are the names of individual counters
	// we're tracking for the interface (for instance, the names of MIB
	// variables) and the values are the absolute values of the counters;
	// consumers are usually interested in the deltas.
	CounterSet struct {
		Values map[string]uint64 `json:"values"`
		IntId  string            `json:"intId"`
	}
)

func NewCounterSet() *CounterSet {
	return &CounterSet{Values: map[string]uint64{}}
}

func NewCounterSetWithId(id string) *CounterSet {
	return &CounterSet{Values: map[string]uint64{}, IntId: id}
}

// Calculate the delta, save the current value, and return the delta.
func (cs *CounterSet) SetValueAndReturnDelta(counterName string, cur uint64) uint64 {
	delta := cs.GetDelta(counterName, cur)
	cs.Values[counterName] = cur
	return delta
}

// Calculate the delta between the previous recorded value (in
// cs.Values[counterName]) and the current (new) value.  Given a nil cs, return 0.
func (cs *CounterSet) GetDelta(counterName string, cur uint64) uint64 {
	if cs == nil {
		return 0
	}

	prev, ok := cs.Values[counterName]
	if !ok { // nil case, just return 0
		return 0
	} else if prev <= cur { // common case: return the diff between the prev & cur values.
		delta := cur - prev
		return delta
	} else {
		// Counter rolled over, or machine rebooted.  But we have no way to
		// determine which.  If we knew for sure it'd rolled over, we could do
		// some math and figure out the real value.  But if the switch just
		// rebooted and reset the counter (to zero or some random / arbitrary
		// value), then really all we can do is restart the delta calculation.
		//
		// You might think we could return "cur", on the theory that regardless
		// of a roll-over or a reboot, we know we've gotten at least "cur" bytes
		// since then, so return that many bytes.  But that's only true if the
		// counters are initialized to zero on reboot, and I don't think we're
		// sure of that.  So just return zero to be on the safe side.
		return 0
	}
}

func (cs *CounterSet) SetValue(validate map[string]string, oid string, val uint64) {
	if vName, ok := validate[oid]; ok {
		cs.Values[vName] = val
	}
}
