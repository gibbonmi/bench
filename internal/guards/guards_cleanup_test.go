package guards

import (
	"context"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
)

func TestScanWaitsForCancelledWorkerCleanup(t *testing.T) {
	oldEnumerate, oldInspect := enumerateGuards, inspectGuard
	t.Cleanup(func() { enumerateGuards, inspectGuard = oldEnumerate, oldInspect })
	enumerateGuards = func(context.Context, string) ([]candidate, error) {
		return []candidate{{fallback: "one"}}, nil
	}
	started := make(chan struct{})
	cancelled := make(chan struct{})
	release := make(chan struct{})
	cleaned := make(chan struct{})
	inspectGuard = func(ctx context.Context, _ string, _ candidate) [][]string {
		close(started)
		<-ctx.Done()
		close(cancelled)
		<-release
		close(cleaned)
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result := make(chan ScanResult, 1)
	go func() { result <- Scan(ctx, "/repo") }()
	<-started
	cancel()
	<-cancelled
	select {
	case got := <-result:
		t.Fatalf("Scan returned before worker cleanup: %#v", got)
	default:
	}
	close(release)
	window := bounds.TestDeadline(0)
	var got ScanResult
	select {
	case got = <-result:
	case <-time.After(window):
		t.Fatal(bounds.TestTimeoutVerdict("Scan to return after worker cleanup", window))
	}
	<-cleaned
	if got.Status != "incomplete" || got.Inspected != "0" || got.Omitted != "1" || got.Reason != "timeout" {
		t.Fatalf("Scan = %#v", got)
	}
}
