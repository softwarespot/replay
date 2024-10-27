package replay

import "iter"

// Replayer defines the replay interface.
type Replayer[T any] interface {
	// All returns an iterator that yields all non-expired events.
	All() iter.Seq[T]

	// Add adds one or more events to the replay buffer.
	// If the replay buffer is full, then the oldest event will be overwritten.
	Add(evts ...T)

	// Clear clears the replay buffer.
	Clear()
}
