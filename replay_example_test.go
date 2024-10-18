package replay_test

import (
	"fmt"
	"time"

	"github.com/softwarespot/replay"
)

func ExampleNew() {
	r := replay.New[string](2048, 256*time.Second)
	r.Add("Event 1")
	r.Add("Event 2")
	r.Add("Event 3")
}

func ExampleReplay_All() {
	r := replay.New[string](64, 128*time.Second)
	r.Add("Event 1")
	r.Add("Event 2")
	r.Add("Event 3")

	for evt := range r.All() {
		fmt.Printf("event: %s\n", evt)
	}

	// output: event: Event 1
	// event: Event 2
	// event: Event 3
}
