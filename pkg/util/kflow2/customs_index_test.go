package chf

import (
	"github.com/stretchr/testify/assert"
	"testing"
	capnp "zombiezen.com/go/capnproto2"
)

func TestCustomsIndex(t *testing.T) {
	valuesForColumn := map[uint32]string{
		1001: "1001",
		1002: "1002",
		1003: "1003",
		1004: "1004",
		1005: "1005",
		1006: "1006",
		1007: "1007",
		1008: "1008",
		1009: "1009",
		1010: "1010",
	}

	// instantiate the sut
	sut := NewCustomsIndex([]uint32{9999, 1002, 1003, 1001, 1005, 1004, 1006, 1007, 1008, 1009, 1010})

	// swizzel currentFlowNumber to make sure we roll over properly
	sut.currentFlowNumber = maxUint32 - 2

	prefixes := []string{"a_", "b_", "c_", "d_"}

	// index and test a bunch of flows
	for _, prefix := range prefixes {
		// create and index the flow
		flow := newFlow(prefix, valuesForColumn)
		sut.IndexFlow(flow)

		// test each column with a value
		for columnID, value := range valuesForColumn {
			customValue, found := sut.CustomColumnWithID(columnID)
			assert.True(t, found)
			strVal, err := customValue.StrVal()
			assert.NoError(t, err)
			assert.Equal(t, prefix+value, strVal)
		}

		// test column the company doesn't have
		_, found := sut.CustomColumnWithID(9999)
		assert.False(t, found)

		// test column the company has, but flow has no value for
		_, found = sut.CustomColumnWithID(8888)
		assert.False(t, found)
	}

	// now try a flow with fewer columns, showing that we ignore garbage from the previous flow
	flow := newFlow("z_", map[uint32]string{1003: "foo", 1010: "bar"})
	sut.IndexFlow(flow)

	// column 1: found
	customValue, found := sut.CustomColumnWithID(1003)
	assert.True(t, found)
	strVal, err := customValue.StrVal()
	assert.NoError(t, err)
	assert.Equal(t, "z_foo", strVal)

	// column 2: found
	customValue, found = sut.CustomColumnWithID(1010)
	assert.True(t, found)
	strVal, err = customValue.StrVal()
	assert.NoError(t, err)
	assert.Equal(t, "z_bar", strVal)

	// columns not found:
	for _, columnID := range []uint32{8888, 1002, 1001, 1009, 1008, 1004, 1005, 1007, 1006} {
		customValue, found = sut.CustomColumnWithID(columnID)
		assert.False(t, found)
	}
}

// benchmark the CustomsIndex, where we index a flow, then fetch the values of the 11 columns 5 times each
func BenchmarkCustomsIndex(b *testing.B) {
	// build the sut & flow
	valuesForColumn := map[uint32]string{
		1001: "1001",
		1002: "1002",
		1003: "1003",
		1004: "1004",
		1005: "1005",
		1006: "1006",
		1007: "1007",
		1008: "1008",
		1009: "1009",
		1010: "1010",
	}
	sut := NewCustomsIndex([]uint32{9999, 1002, 1003, 1001, 1005, 1004, 1006, 1007, 1008, 1009, 1010})
	flow := newFlow("", valuesForColumn)

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		// index the flow
		sut.IndexFlow(flow)

		// fetch the values of each custom dimension 5 times
		for i := 0; i < 5; i++ {
			sut.CustomColumnWithID(8888)
			sut.CustomColumnWithID(1001)
			sut.CustomColumnWithID(1002)
			sut.CustomColumnWithID(1003)
			sut.CustomColumnWithID(1004)
			sut.CustomColumnWithID(1005)
			sut.CustomColumnWithID(1006)
			sut.CustomColumnWithID(1007)
			sut.CustomColumnWithID(1008)
			sut.CustomColumnWithID(1009)
			sut.CustomColumnWithID(1010)
		}
	}
}

// benchmark a simple index of the custom columns made from a map per flow,
// where we index a flow, then fetch the values of the 11 columns 5 times each
func BenchmarkFlowMapIndex(b *testing.B) {
	valuesForColumn := map[uint32]string{
		1001: "1001",
		1002: "1002",
		1003: "1003",
		1004: "1004",
		1005: "1005",
		1006: "1006",
		1007: "1007",
		1008: "1008",
		1009: "1009",
		1010: "1010",
	}
	flow := newFlow("", valuesForColumn)
	var tmpCustom Custom

	// not 100% sure this belongs here, but if it's just in the loop, benchmark doesn't consider it an alloc
	var customByID map[uint32]Custom

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		// Note: this is setup to approximate streaming repo's GetCustomColumnWithID

		// index the flow into a new map
		customList, _ := flow.Custom()
		customByID = make(map[uint32]Custom, customList.Len())
		for i := 0; i < customList.Len(); i++ {
			c := customList.At(i)
			customByID[c.Id()] = c
		}

		// fetch the values of each custom dimension 5 times
		for i := 0; i < 5; i++ {
			if customByID != nil { // always true, but here to simulate GetCustomColumnWithID
				tmpCustom = customByID[8888]
				tmpCustom = customByID[1001]
				tmpCustom = customByID[1002]
				tmpCustom = customByID[1003]
				tmpCustom = customByID[1004]
				tmpCustom = customByID[1005]
				tmpCustom = customByID[1006]
				tmpCustom = customByID[1007]
				tmpCustom = customByID[1008]
				tmpCustom = customByID[1009]
				tmpCustom = customByID[1010]
			}
		}
	}

	// no-op to allow lookups in customByID without "not used" errors
	tmpCustom.Id()
}

// benchmark a simple index of the custom columns made from a REUSED map,
// where we index a flow, then fetch the values of the 11 columns 5 times each
func BenchmarkReusedFlowMapIndex(b *testing.B) {
	valuesForColumn := map[uint32]string{
		1001: "1001",
		1002: "1002",
		1003: "1003",
		1004: "1004",
		1005: "1005",
		1006: "1006",
		1007: "1007",
		1008: "1008",
		1009: "1009",
		1010: "1010",
	}
	flow := newFlow("", valuesForColumn)
	var tmpCustom Custom

	customByID := make(map[uint32]Custom, 10)

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		// Note: this is setup to approximate streaming repo's GetCustomColumnWithID

		// index the flow into a new map
		customList, _ := flow.Custom()
		for i := 0; i < customList.Len(); i++ {
			c := customList.At(i)
			customByID[c.Id()] = c
		}

		// fetch the values of each custom dimension 5 times
		for i := 0; i < 5; i++ {
			if customByID != nil { // always true, but here to simulate GetCustomColumnWithID
				tmpCustom = customByID[8888]
				tmpCustom = customByID[1001]
				tmpCustom = customByID[1002]
				tmpCustom = customByID[1003]
				tmpCustom = customByID[1004]
				tmpCustom = customByID[1005]
				tmpCustom = customByID[1006]
				tmpCustom = customByID[1007]
				tmpCustom = customByID[1008]
				tmpCustom = customByID[1009]
				tmpCustom = customByID[1010]
			}
		}

		// clear out the map for next iteration
		for k := range customByID {
			delete(customByID, k)
		}
	}

	// no-op to allow lookups in customByID without "not used" errors
	tmpCustom.Id()
}

// build a test flow with a bunch of custom column values
func newFlow(prefix string, valuesForColumn map[uint32]string) CHF {
	_, buf, _ := capnp.NewMessage(capnp.SingleSegment(nil))
	ret, _ := NewCHF(buf)

	list, _ := NewCustom_List(buf, int32(len(valuesForColumn)))
	i := 0
	for k, v := range valuesForColumn {
		cst, _ := NewCustom(buf)
		cst.SetId(k)
		cst.Value().SetStrVal(prefix + v)
		list.Set(i, cst)
		i++
	}
	ret.SetCustom(list)

	return ret
}
