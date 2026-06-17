package ringbuffer

import "fmt"

// entry holds a key-value pair in the backing slice.
type entry[V any] struct {
	key   string
	value V
}

// Entry is the public view of a key-value pair, returned by Oldest/Newest/Entries.
type Entry[V any] struct {
	Key   string
	Value V
}

// RingBuffer is a fixed-capacity, slice-backed store of string-keyed values.
// When full, inserting a new key evicts the oldest entry to make room.
//
// Get runs in O(1) via a companion map[string]int that records each key's
// physical slot in the backing slice. The map is kept in sync on every Put
// and every eviction.
type RingBuffer[V any] struct {
	data     []entry[V]
	index    map[string]int // key → physical slot in data
	head     int            // physical slot of the oldest entry
	tail     int            // physical slot where the next write goes
	size     int
	capacity int
}

// New creates a RingBuffer with the given capacity.
// Panics if capacity is less than 1.
func New[V any](capacity int) *RingBuffer[V] {
	if capacity < 1 {
		panic("ringbuffer: capacity must be at least 1")
	}
	return &RingBuffer[V]{
		data:     make([]entry[V], capacity),
		index:    make(map[string]int, capacity),
		capacity: capacity,
	}
}

// Put inserts or updates key with the given value.
// - If key is new and the buffer is full, the oldest entry is evicted first.
// - If key already exists its value is updated in-place; insertion order is unchanged.
func (r *RingBuffer[V]) Put(key string, value V) {
	// Key already present: update value in-place, preserve insertion order.
	if slot, exists := r.index[key]; exists {
		r.data[slot].value = value
		return
	}

	if r.size == r.capacity {
		// Evict oldest: delete its map entry before overwriting the slot.
		evicted := r.data[r.head].key
		delete(r.index, evicted)
		r.head = (r.head + 1) % r.capacity
	} else {
		r.size++
	}

	r.data[r.tail] = entry[V]{key: key, value: value}
	r.index[key] = r.tail
	r.tail = (r.tail + 1) % r.capacity
}

// Get returns the value associated with key and true if present,
// or the zero value of V and false if not. Runs in O(1).
func (r *RingBuffer[V]) Get(key string) (value V, found bool) {
	slot, ok := r.index[key]
	if !ok {
		var zero V
		return zero, false
	}
	return r.data[slot].value, true
}

// Like get but nulls out the value when returning.
func (r *RingBuffer[V]) GetAndDelete(key string) (value V, found bool) {
	slot, ok := r.index[key]
	if !ok {
		var zero V
		return zero, false
	}

	val := r.data[slot].value
	delete(r.index, key) // ensure we don't overlap again on flows.
	return val, true
}

// Contains reports whether key is present. Runs in O(1).
func (r *RingBuffer[V]) Contains(key string) bool {
	_, ok := r.index[key]
	return ok
}

// Oldest returns the Entry that would be evicted next and true,
// or a zero Entry and false if the buffer is empty.
func (r *RingBuffer[V]) Oldest() (Entry[V], bool) {
	if r.size == 0 {
		return Entry[V]{}, false
	}
	e := r.data[r.head]
	return Entry[V]{Key: e.key, Value: e.value}, true
}

// Newest returns the most recently inserted Entry and true,
// or a zero Entry and false if the buffer is empty.
func (r *RingBuffer[V]) Newest() (Entry[V], bool) {
	if r.size == 0 {
		return Entry[V]{}, false
	}
	e := r.data[(r.tail-1+r.capacity)%r.capacity]
	return Entry[V]{Key: e.key, Value: e.value}, true
}

// Entries returns all entries in insertion order (oldest first).
func (r *RingBuffer[V]) Entries() []Entry[V] {
	out := make([]Entry[V], r.size)
	for i := range r.size {
		e := r.data[(r.head+i)%r.capacity]
		out[i] = Entry[V]{Key: e.key, Value: e.value}
	}
	return out
}

// Keys returns all keys in insertion order (oldest first).
func (r *RingBuffer[V]) Keys() []string {
	out := make([]string, r.size)
	for i := range r.size {
		out[i] = r.data[(r.head+i)%r.capacity].key
	}
	return out
}

// Len returns the number of entries currently stored.
func (r *RingBuffer[V]) Len() int { return r.size }

// Cap returns the maximum number of entries the buffer can hold.
func (r *RingBuffer[V]) Cap() int { return r.capacity }

// IsFull reports whether the buffer has reached capacity.
func (r *RingBuffer[V]) IsFull() bool { return r.size == r.capacity }

// String returns a human-readable representation of the buffer.
func (r *RingBuffer[V]) String() string {
	return fmt.Sprintf("RingBuffer{cap: %d, len: %d, entries: %v}",
		r.capacity, r.size, r.Entries())
}
