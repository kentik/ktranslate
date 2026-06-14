package ringbuffer

import "testing"

// ── Basic put/get with ints ───────────────────────────────────────────────────

func TestPutAndGet(t *testing.T) {
	rb := New[int](3)
	rb.Put("alpha", 1)
	rb.Put("beta", 2)
	rb.Put("gamma", 3)

	for key, want := range map[string]int{"alpha": 1, "beta": 2, "gamma": 3} {
		got, ok := rb.Get(key)
		if !ok {
			t.Errorf("Get(%q): expected found=true", key)
		}
		if got != want {
			t.Errorf("Get(%q): got %d, want %d", key, got, want)
		}
	}
}

func TestGetMissingKey(t *testing.T) {
	rb := New[string](3)
	rb.Put("foo", "bar")
	v, ok := rb.Get("missing")
	if ok {
		t.Error("Get on missing key should return false")
	}
	if v != "" {
		t.Errorf("Get on missing key should return zero value, got %q", v)
	}
}

func TestGetZeroValueForMissingStruct(t *testing.T) {
	type Point struct{ X, Y int }
	rb := New[Point](2)
	rb.Put("origin", Point{0, 0})
	got, ok := rb.Get("nowhere")
	if ok || got != (Point{}) {
		t.Errorf("expected zero Point and false, got %v %v", got, ok)
	}
}

// ── Update in-place ───────────────────────────────────────────────────────────

func TestUpdateExistingKeyInPlace(t *testing.T) {
	rb := New[int](3)
	rb.Put("a", 1)
	rb.Put("b", 2)
	rb.Put("a", 99) // update, must not shift insertion order

	if rb.Len() != 2 {
		t.Errorf("Len: got %d, want 2", rb.Len())
	}
	keys := rb.Keys()
	if keys[0] != "a" || keys[1] != "b" {
		t.Errorf("Keys after update: got %v, want [a b]", keys)
	}
	v, _ := rb.Get("a")
	if v != 99 {
		t.Errorf("Get(\"a\") after update: got %d, want 99", v)
	}
}

// ── Eviction ──────────────────────────────────────────────────────────────────

func TestEvictsOldest(t *testing.T) {
	rb := New[string](3)
	rb.Put("alpha", "a")
	rb.Put("beta", "b")
	rb.Put("gamma", "g")
	rb.Put("delta", "d") // evicts "alpha"

	if rb.Contains("alpha") {
		t.Error("\"alpha\" should have been evicted")
	}
	for key, want := range map[string]string{"beta": "b", "gamma": "g", "delta": "d"} {
		got, ok := rb.Get(key)
		if !ok || got != want {
			t.Errorf("Get(%q): got %q %v, want %q true", key, got, ok, want)
		}
	}

	oldest, _ := rb.Oldest()
	if oldest.Key != "beta" || oldest.Value != "b" {
		t.Errorf("Oldest: got %+v, want {beta b}", oldest)
	}
}

func TestOverflowMultipleTimes(t *testing.T) {
	rb := New[int](2)
	pairs := []struct {
		k string
		v int
	}{{"a", 1}, {"b", 2}, {"c", 3}, {"d", 4}, {"e", 5}}
	for _, p := range pairs {
		rb.Put(p.k, p.v)
	}
	// Only "d":4 and "e":5 should remain.
	entries := rb.Entries()
	want := []Entry[int]{{"d", 4}, {"e", 5}}
	if len(entries) != len(want) {
		t.Fatalf("Entries() len: got %d, want %d", len(entries), len(want))
	}
	for i, w := range want {
		if entries[i] != w {
			t.Errorf("Entries()[%d]: got %+v, want %+v", i, entries[i], w)
		}
	}
}

func TestMapStaysInSyncAfterManyEvictions(t *testing.T) {
	rb := New[int](3)
	for i, k := range []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"} {
		rb.Put(k, i)
	}
	// Only g:6, h:7, i:8 should survive.
	for _, e := range []Entry[int]{{"g", 6}, {"h", 7}, {"i", 8}} {
		got, ok := rb.Get(e.Key)
		if !ok {
			t.Errorf("expected %q to be present", e.Key)
		} else if got != e.Value {
			t.Errorf("Get(%q): got %d, want %d", e.Key, got, e.Value)
		}
	}
	for _, k := range []string{"a", "b", "c", "d", "e", "f"} {
		if rb.Contains(k) {
			t.Errorf("%q should have been evicted", k)
		}
	}
}

// ── Edge cases ────────────────────────────────────────────────────────────────

func TestCapacityOne(t *testing.T) {
	rb := New[bool](1)
	rb.Put("only", true)
	if !rb.Contains("only") {
		t.Error("expected \"only\" to be present")
	}
	rb.Put("replaced", false)
	if rb.Contains("only") {
		t.Error("\"only\" should have been evicted")
	}
	v, ok := rb.Get("replaced")
	if !ok || v != false {
		t.Errorf("Get(\"replaced\"): got %v %v, want false true", v, ok)
	}
}

func TestEmptyBuffer(t *testing.T) {
	rb := New[int](3)
	if _, ok := rb.Oldest(); ok {
		t.Error("Oldest on empty buffer should return false")
	}
	if _, ok := rb.Newest(); ok {
		t.Error("Newest on empty buffer should return false")
	}
	if rb.Len() != 0 {
		t.Error("Len on empty buffer should be 0")
	}
}

func TestIsFull(t *testing.T) {
	rb := New[float64](2)
	if rb.IsFull() {
		t.Error("new buffer should not be full")
	}
	rb.Put("pi", 3.14)
	rb.Put("e", 2.71)
	if !rb.IsFull() {
		t.Error("buffer should be full after 2 inserts into cap-2")
	}
}

func TestOrderPreserved(t *testing.T) {
	rb := New[int](4)
	pairs := []Entry[int]{{"w", 10}, {"x", 20}, {"y", 30}, {"z", 40}}
	for _, p := range pairs {
		rb.Put(p.Key, p.Value)
	}
	got := rb.Entries()
	for i, want := range pairs {
		if got[i] != want {
			t.Errorf("Entries()[%d]: got %+v, want %+v", i, got[i], want)
		}
	}
}

// ── Generic flexibility ───────────────────────────────────────────────────────

func TestWithStructValues(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}
	rb := New[User](2)
	rb.Put("u1", User{"Alice", 30})
	rb.Put("u2", User{"Bob", 25})

	u, ok := rb.Get("u1")
	if !ok || u.Name != "Alice" || u.Age != 30 {
		t.Errorf("Get(\"u1\"): got %+v %v, want {Alice 30} true", u, ok)
	}

	rb.Put("u3", User{"Carol", 28}) // evicts "u1"
	if rb.Contains("u1") {
		t.Error("\"u1\" should have been evicted")
	}
}

func TestWithSliceValues(t *testing.T) {
	rb := New[[]int](2)
	rb.Put("evens", []int{2, 4, 6})
	rb.Put("odds", []int{1, 3, 5})

	evens, ok := rb.Get("evens")
	if !ok || len(evens) != 3 || evens[0] != 2 {
		t.Errorf("Get(\"evens\"): unexpected result %v %v", evens, ok)
	}
}

func TestNewest(t *testing.T) {
	rb := New[string](3)
	rb.Put("first", "a")
	rb.Put("second", "b")

	newest, ok := rb.Newest()
	if !ok || newest.Key != "second" || newest.Value != "b" {
		t.Errorf("Newest: got %+v %v, want {second b} true", newest, ok)
	}
}
