package replay

import (
	"iter"
	"sync"
	"time"
)

// SyncReplay is a generic, thread safe replay buffer.
type SyncReplay[T any] struct {
	replay *Replay[T]
	mu     sync.RWMutex
}

// New initializes a sync replay buffer.
func NewSyncReplay[T any](maxSize int, expiry time.Duration) *SyncReplay[T] {
	return &SyncReplay[T]{
		replay: New[T](maxSize, expiry),
	}
}

// Iter returns an iterator that yields all non-expired events.
func (sr *SyncReplay[T]) Iter() iter.Seq[T] {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	return sr.replay.Iter()
}

// Add adds one or more events to the sync replay buffer.
// If the sync replay buffer is full, then the oldest event will be overwritten.
func (sr *SyncReplay[T]) Add(evts ...T) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	sr.replay.Add(evts...)
}

// Clear clears the sync replay buffer.
func (sr *SyncReplay[T]) Clear() {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	sr.replay.Clear()
}
