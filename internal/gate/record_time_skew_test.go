package gate

import (
	"strings"
	"testing"
	"time"
)

// A host wall clock can step back by up to a second under load, so a record written
// just before the step reads as slightly ahead of the reader's now. Each record class
// accepts that small step and still refuses a record that is clearly ahead.
func TestRecordTimeToleratesABackwardClockStep(t *testing.T) {
	now := time.Date(2026, 10, 3, 11, 59, 55, 214_000_000, time.UTC)
	tree, digest := strings.Repeat("a", 40), strings.Repeat("b", 64)
	records := map[string]func(stamp string) verdictRecord{
		"full verdict": func(stamp string) verdictRecord {
			return verdictRecord{Schema: verdictSchema, State: Ready, Status: "green", Tree: tree, Oracle: digest, RecordedAt: stamp}
		},
		"pending": func(stamp string) verdictRecord {
			return verdictRecord{Schema: verdictSchema, State: Pending, Tree: tree, Oracle: digest, StartedAt: stamp, OwnerPID: 1}
		},
		"lane record": func(stamp string) verdictRecord {
			return verdictRecord{Schema: verdictSchema, Tree: tree, Lane: "fast", Outcome: lanePass, RunBinary: digest, RecordedAt: stamp}
		},
	}
	for class, build := range records {
		for _, tc := range []struct {
			ahead time.Duration
			valid bool
		}{{time.Second, true}, {3 * time.Second, false}} {
			stamp := now.Add(tc.ahead).Truncate(time.Second).Format(time.RFC3339)
			record := build(stamp)
			got, err := validateRecordBytes(inspectJSON(record), record, now)
			if got.name != class {
				t.Fatalf("%s at +%s selected class %q", class, tc.ahead, got.name)
			}
			if (err == nil) != tc.valid {
				t.Errorf("%s at +%s: err = %v, want valid = %t", class, tc.ahead, err, tc.valid)
			}
		}
	}
}
