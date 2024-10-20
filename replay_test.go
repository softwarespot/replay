package replay

import (
	"testing"
	"time"
)

func Test_New(t *testing.T) {
	r := New[string](5, 5*time.Second)

	// Add events to the replay buffer
	nowFn = func() time.Time { return parseAsDateTime("2024-10-01 00:00:00") }
	r.Add("Event 1")
	r.Add("Event 2")
	assertEqualForAll(t, r, []string{
		"Event 1",
		"Event 2",
	})

	// Add new events 10s later
	nowFn = func() time.Time { return parseAsDateTime("2024-10-01 00:00:10") }
	r.Add("Event 3")
	r.Add("Event 4")

	// Should not return the expired events
	assertEqualForAll(t, r, []string{
		"Event 3",
		"Event 4",
	})

	// Add new events 1s later, which should not return expired events
	nowFn = func() time.Time { return parseAsDateTime("2024-10-01 00:00:11") }
	r.Add("Event 5")
	r.Add("Event 6", "Event 7")
	assertEqualForAll(t, r, []string{
		"Event 3",
		"Event 4",
		"Event 5",
		"Event 6",
		"Event 7",
	})

	// Should clear the events
	r.Clear()
	assertEqualForAll(t, r, nil)
}
