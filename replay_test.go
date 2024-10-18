package replay

import (
	"testing"
	"time"
)

func Test_New(t *testing.T) {
	r := New[string](5, 5*time.Second)

	// Add events to the replay buffer
	nowFn = func() time.Time { return parseAsDateTime("2024-10-01 00:00:00") }
	r.Add("Event1")
	r.Add("Event2")
	assertEqualForAll(t, r, []string{
		"Event1",
		"Event2",
	})

	// Add new events 10s later
	nowFn = func() time.Time { return parseAsDateTime("2024-10-01 00:00:10") }
	r.Add("Event3")
	r.Add("Event4")

	// Should not return the expired events
	assertEqualForAll(t, r, []string{
		"Event3",
		"Event4",
	})

	// Add new events 1s later, which should not return expired events
	nowFn = func() time.Time { return parseAsDateTime("2024-10-01 00:00:11") }
	r.Add("Event5")
	r.Add("Event6")
	r.Add("Event7")
	assertEqualForAll(t, r, []string{
		"Event3",
		"Event4",
		"Event5",
		"Event6",
		"Event7",
	})

	// Should clear the events
	r.Clear()
	assertEqualForAll(t, r, nil)
}
