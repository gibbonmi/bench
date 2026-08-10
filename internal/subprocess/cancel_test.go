package subprocess

import (
	"os"
	"syscall"
	"testing"
)

// TestCancelSignalsMembershipIsExact pins the set every owner registers. An
// authored expectation is the point: it is independent of CancelSignals' own
// literal, so adding or dropping a signal there reds here instead of silently
// changing what thirteen call sites trap.
func TestCancelSignalsMembershipIsExact(t *testing.T) {
	want := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGHUP}
	if len(CancelSignals) != len(want) {
		t.Fatalf("CancelSignals = %v, want exactly %v", CancelSignals, want)
	}
	for index, signal := range want {
		if CancelSignals[index] != signal {
			t.Errorf("CancelSignals[%d] = %v, want %v", index, CancelSignals[index], signal)
		}
	}
}
