package counters

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounters(t *testing.T) {
	cs := NewCounterSet()
	assert.NotNil(t, cs)

	intf := "1"

	d := cs.SetValueAndReturnDelta(intf, 10)
	assert.Equal(t, int(0), int(d))

	d = cs.SetValueAndReturnDelta(intf, 10)
	assert.Equal(t, int(0), int(d))

	d = cs.SetValueAndReturnDelta(intf, 20)
	assert.Equal(t, int(10), int(d))

	// Rollovers and resets all return 0
	d = cs.SetValueAndReturnDelta(intf, 19) // rollover 32
	assert.Equal(t, int(0), int(d))

	cs.SetValueAndReturnDelta(intf, math.MaxUint32-10)
	d = cs.SetValueAndReturnDelta(intf, 20)
	assert.Equal(t, int(0), int(d))

	// Rollover 64b counter
	cs.SetValueAndReturnDelta(intf, math.MaxUint64-10)
	d = cs.SetValueAndReturnDelta(intf, 20)
	assert.Equal(t, int(0), int(d))
}

func TestCounterEqual(t *testing.T) {
	cs := NewCounterSet()
	assert.NotNil(t, cs)

	intf := "1"

	d := cs.SetValueAndReturnDelta(intf, 10)
	assert.Equal(t, int(0), int(d))

	d = cs.SetValueAndReturnDelta(intf, 819504)
	assert.Equal(t, int(819504-10), int(d))

	d = cs.SetValueAndReturnDelta(intf, 819504)
	assert.Equal(t, int(0), int(d))
}

func TestCounterRolloverZero(t *testing.T) {
	cs := NewCounterSet()
	assert.NotNil(t, cs)

	intf := "1"

	d := cs.SetValueAndReturnDelta(intf, 0)
	assert.Equal(t, int(0), int(d))

	d = cs.SetValueAndReturnDelta(intf, math.MaxUint32)
	assert.Equal(t, int(math.MaxUint32), int(d))

	d = cs.SetValueAndReturnDelta(intf, 10)
	assert.Equal(t, int(0), int(d))
}

func TestCounterGetDelta(t *testing.T) {
	intf := "1"

	var cs *CounterSet // nil
	d := cs.GetDelta(intf, 10)
	assert.Equal(t, int(0), int(d))

	cs = NewCounterSet()
	assert.NotNil(t, cs)

	d = cs.GetDelta(intf, 0)
	assert.Equal(t, int(0), int(d))

	cs.Values[intf] = 0
	d = cs.GetDelta(intf, math.MaxUint32)
	assert.Equal(t, int(math.MaxUint32), int(d))

	cs.Values[intf] = 90
	d = cs.GetDelta(intf, 100)
	assert.Equal(t, int(10), int(d))
}
