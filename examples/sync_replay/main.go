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
	// Create a sync replay buffer of type "Event"
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
