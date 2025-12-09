package replay

import (
	"iter"
	"time"
)

// Can be overridden for testing
var nowFn = time.Now

type replayedEvent[T any] struct {
	event   T
	expires time.Time
}

// Replay is a generic, non-thread safe replay buffer.
type Replay[T any] struct {
	events []replayedEvent[T]
	size   int
	head   int
	tail   int
	expiry time.Duration
}

// New initializes a replay buffer.
func New[T any](maxSize int, expiry time.Duration) *Replay[T] {
	return &Replay[T]{
		events: make([]replayedEvent[T], maxSize),
		size:   0,
		head:   0,
		tail:   0,
		expiry: expiry,
	}
}

// All returns an iterator that yields all non-expired events.
func (r *Replay[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		now := nowFn()
		maxSize := len(r.events)

		for i := 0; i < r.size; i++ {
			idx := (r.head + i) % maxSize
			if evt := r.events[idx]; evt.expires.After(now) {
				if !yield(evt.event) {
					return
				}
			}
		}
	}
}

// Add adds one or more events to the replay buffer.
// If the replay buffer is full, then the oldest event will be overwritten.
func (r *Replay[T]) Add(evts ...T) {
	maxSize := len(r.events)
	for _, evt := range evts {
		isFull := r.size == maxSize
		r.events[r.tail] = replayedEvent[T]{
			event:   evt,
			expires: nowFn().Add(r.expiry),
		}
		r.tail = (r.tail + 1) % maxSize

		if isFull {
			r.head = (r.head + 1) % maxSize
		} else {
			r.size++
		}
	}
}

// Clear clears the replay buffer.
func (r *Replay[T]) Clear() {
	clear(r.events)
	r.size = 0
	r.head = 0
	r.tail = 0
}
