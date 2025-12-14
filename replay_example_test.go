package replay_test

import (
	"fmt"
	"time"

	"github.com/softwarespot/replay"
)

func ExampleNew() {
	r := replay.New[string](2048, 256*time.Second)
	r.Add("Event 1")
	r.Add("Event 2", "Event 3")
	r.Add("Event 4")
}

func ExampleReplay_Iter() {
	r := replay.New[string](64, 128*time.Second)
	r.Add("Event 1")
	r.Add("Event 2", "Event 3")
	r.Add("Event 4")

	for evt := range r.Iter() {
		fmt.Printf("event: %s\n", evt)
	}
	// output: event: Event 1
	// event: Event 2
	// event: Event 3
	// event: Event 4
}

func ExampleNewSyncReplay() {
	sr := replay.NewSyncReplay[string](2048, 256*time.Second)
	sr.Add("Event 1")
	sr.Add("Event 2", "Event 3")
	sr.Add("Event 4")
}

func ExampleSyncReplay_Iter() {
	sr := replay.NewSyncReplay[string](64, 128*time.Second)
	sr.Add("Event 1")
	sr.Add("Event 2", "Event 3")
	sr.Add("Event 4")

	for evt := range sr.Iter() {
		fmt.Printf("event: %s\n", evt)
	}
	// output: event: Event 1
	// event: Event 2
	// event: Event 3
	// event: Event 4
}
