package replay

import (
	"testing"
	"time"
)

func Test_NewSyncReplay(t *testing.T) {
	sr := NewSyncReplay[string](5, 5*time.Second)

	// Add events to the replay buffer
	nowFn = func() time.Time { return parseAsDateTime("2024-10-01 00:00:00") }
	sr.Add("Event 1")
	sr.Add("Event 2")
	assertEqual(t, sr, []string{
		"Event 1",
		"Event 2",
	})

	// Add new events 10s later
	nowFn = func() time.Time { return parseAsDateTime("2024-10-01 00:00:10") }
	sr.Add("Event 3")
	sr.Add("Event 4")

	// Should not return the expired events
	assertEqual(t, sr, []string{
		"Event 3",
		"Event 4",
	})

	// Add new events 1s later, which should not return expired events
	nowFn = func() time.Time { return parseAsDateTime("2024-10-01 00:00:11") }
	sr.Add("Event 5")
	sr.Add("Event 6", "Event 7")
	assertEqual(t, sr, []string{
		"Event 3",
		"Event 4",
		"Event 5",
		"Event 6",
		"Event 7",
	})

	// Should clear the events
	sr.Clear()
	assertEqual(t, sr, nil)
}
