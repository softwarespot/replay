# Replay [![Go Reference](https://pkg.go.dev/badge/github.com/softwarespot/replay.svg)](https://pkg.go.dev/github.com/softwarespot/replay)

**Replay** is a generic, non-thread-safe implementation, designed to store and manage a fixed-size buffer of events with expiration logic in-built.

This is useful in scenarios where you need to keep track of recent events while automatically discarding those outdated events.
For example in a chat application where new users should see the last N number of messages, then this implementation is perfect for such a use case.

## Prerequisites

-   Go 1.23.0 or above

## Installation

```bash
go get -u github.com/softwarespot/replay
```

## Usage

A basic example of using **Replay**.

```Go
package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/softwarespot/replay"
)

func main() {
	// Create a replay buffer with the type "string"
   	r := replay.New[string](64, 128*time.Second)
	r.Add("Event 1")
	r.Add("Event 2")
	r.Add("Event 3")

	for evt := range r.All() {
		fmt.Printf("event: %s\n", evt)
	}
}
```

**Replay** by default is not thread-safe.
This is an example of creating a thread-safe wrapper, in just a few lines of code.

```Go
package main

import (
	"fmt"
	"iter"
	"sync"
	"time"

	"github.com/softwarespot/replay"
)

// A generic thread-safe wrapper around replay.
type SyncReplay[T any] struct {
	replay *replay.Replay[T]
	mu     sync.RWMutex
}

func NewSyncReplay[T any](maxSize int, expiry time.Duration) *SyncReplay[T] {
	return &SyncReplay[T]{
		replay: replay.New[T](maxSize, expiry),
	}
}

func (sr *SyncReplay[T]) Add(evt T) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	sr.replay.Add(evt)
}

func (sr *SyncReplay[T]) All() iter.Seq[T] {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	return sr.replay.All()
}

type Event struct {
	ID int
}

func main() {
	// Create a sync replay buffer with the type "Event"
	sr := NewSyncReplay[Event](64, 256*time.Second)
	go func() {
		var id int
		for {
			sr.Add(Event{
				ID: id,
			})
			id++
			time.Sleep(256 * time.Millisecond)
		}
	}()

	fmt.Println("wait 5s for the replay buffer to contain events")
	time.Sleep(5 * time.Second)

	var wg sync.WaitGroup

	wg.Add(1)
	go worker(&wg, 1, sr)

	fmt.Println("wait 3s for the replay buffer to contain more events")
	time.Sleep(3 * time.Second)

	wg.Add(1)
	go worker(&wg, 2, sr)

	wg.Wait()
	fmt.Println("done. Notice that the 2nd worker replayed more events than the 1st worker?")
}

func worker[T Event](wg *sync.WaitGroup, id int, sr *SyncReplay[Event]) {
	defer wg.Done()
	for evt := range sr.All() {
		fmt.Printf("worker ID: %d, event: %d\n", id, evt.ID)
	}
}
```

## License

The code has been licensed under the [MIT](https://opensource.org/license/mit) license.
