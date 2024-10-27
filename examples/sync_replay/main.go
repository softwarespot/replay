package main

import (
	"fmt"
	"time"

	"github.com/softwarespot/replay"
)

type Event struct {
	ID int
}

func main() {
	// Create a sync replay buffer with the type "Event"
	sr := replay.NewSyncReplay[Event](64, 256*time.Second)
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

	worker(1, sr)

	fmt.Println("wait 3s for the replay buffer to contain more events")
	time.Sleep(3 * time.Second)

	worker(2, sr)

	fmt.Println("done. Notice that the 2nd worker replayed more events than the 1st worker?")
}

func worker[T Event](id int, sr *replay.SyncReplay[Event]) {
	for evt := range sr.All() {
		fmt.Printf("worker ID: %d, event: %d\n", id, evt.ID)
	}
}
