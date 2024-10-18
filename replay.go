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

// A generic, non-thread safe replay buffer.
type Replay[T any] struct {
	idx     int
	events  []replayedEvent[T]
	maxSize int
	expiry  time.Duration
}

// New initializes a replay buffer.
func New[T any](maxSize int, expiry time.Duration) *Replay[T] {
	return &Replay[T]{
		idx:     0,
		events:  make([]replayedEvent[T], maxSize),
		maxSize: maxSize,
		expiry:  expiry,
	}
}

// All returns an iterator that yields all non-expired events.
func (r *Replay[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		now := nowFn()
		for i := 0; i < r.maxSize; i++ {
			idx := (r.idx + i) % r.maxSize
			evt := r.events[idx]
			if evt.expires.After(now) {
				if !yield(evt.event) {
					return
				}
			}
		}
	}
}

// Add adds an event to the replay buffer.
// If the replay buffer is full, then the oldest event will be overwritten with the new event.
func (r *Replay[T]) Add(evt T) {
	r.events[r.idx] = replayedEvent[T]{
		event:   evt,
		expires: nowFn().Add(r.expiry),
	}
	r.idx = (r.idx + 1) % r.maxSize
}

// Clear resets the replay buffer.
func (r *Replay[T]) Clear() {
	clear(r.events)
	r.idx = 0
}
