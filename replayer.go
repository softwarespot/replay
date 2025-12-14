package replay

import "iter"

// Replayer defines the replay interface.
type Replayer[T any] interface {
	// Iter returns an iterator that yields all non-expired events.
	Iter() iter.Seq[T]

	// Add adds one or more events to the replay buffer.
	// If the replay buffer is full, then the oldest event will be overwritten.
	Add(evts ...T)

	// Clear clears the replay buffer.
	Clear()
}
