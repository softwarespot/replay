package replay

import (
	"reflect"
	"testing"
	"time"
)

// assertEqual checks if two values are equal. If they are not, it logs using t.Fatalf()
func assertEqual[T any](t testing.TB, got, correct T) {
	t.Helper()
	if !reflect.DeepEqual(got, correct) {
		t.Fatalf("assertEqual: expected values to be equal, got:\n%+v\ncorrect:\n%+v", got, correct)
	}
}

// assertEqualForAll check if the replayed events are equal. If they are not, it logs using t.Fatalf()
func assertEqualForAll[T any](t testing.TB, r *Replay[T], correct []T) {
	t.Helper()
	var got []T
	for evt := range r.All() {
		got = append(got, evt)
	}
	if !reflect.DeepEqual(got, correct) {
		t.Fatalf("assertEqualForAll: expected values to be equal, got:\n%+v\ncorrect:\n%+v", got, correct)
	}
}

// parseAsDateTime parses a string representation of date and time
// in the format "YYYY-MM-DD HH:MM:SS" into a time.Time value
// based on the local time zone
func parseAsDateTime(tt string) time.Time {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", tt, time.Local)
	if err != nil {
		return time.Time{}
	}
	return t
}
