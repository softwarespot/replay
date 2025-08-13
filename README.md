# Replay

[![Go Reference](https://pkg.go.dev/badge/github.com/softwarespot/replay.svg)](https://pkg.go.dev/github.com/softwarespot/replay) ![Go Tests](https://github.com/softwarespot/replay/actions/workflows/go.yml/badge.svg)

**Replay** is a generic compatible module with a thread (**SyncReplay**) and non-thread (**Replay**) safe implementation. It's designed to store and manage a fixed-size buffer of events with expiration logic in-built.

This is particularly useful for cases where you need to keep track of recent events, while automatically discarding those which are outdated.
For example, this could be used in a chat application, where new users who join, should see the recently sent messages.

Examples of using this module can be found from the [./examples](./examples) directory.

## Prerequisites

- Go 1.25.0 or above

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
	r.Add("Event 2", "Event 3")
	r.Add("Event 4")

	for evt := range r.All() {
		fmt.Printf("event: %s\n", evt)
	}
}
```

A basic example of using **SyncReplay**.

```Go
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
```

## License

The code has been licensed under the [MIT](https://opensource.org/license/mit) license.
