package guards

import (
	"context"
	"sync"
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
	result, finished := make(chan ScanResult, 1), make(chan struct{})
	go func() { defer close(finished); result <- Scan(ctx, "/repo") }()
	// A failed assertion still releases the held worker and joins the scan before the
	// cleanup above restores the replaced collaborators.
	released := false
	t.Cleanup(func() {
		cancel()
		if !released {
			close(release)
		}
		<-finished
	})
	<-started
	cancel()
	<-cancelled
	select {
	case got := <-result:
		t.Fatalf("Scan returned before worker cleanup: %#v", got)
	default:
	}
	released = true
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

// TestScanOverrideJoinsBeforeRestore proves the join-before-restore order of a replaced
// scan collaborator. A scan worker is held behind a barrier while the subtest leaves. The
// subtest's join releases and waits for that worker before its earlier cleanup restores
// the prior collaborators, so the worker still observes the replacement. Without the join,
// the restoration runs first and the released worker observes the prior value.
func TestScanOverrideJoinsBeforeRestore(t *testing.T) {
	marker := func(name string) func(context.Context, string) ([]candidate, error) {
		return func(context.Context, string) ([]candidate, error) { return []candidate{{fallback: name}}, nil }
	}
	installed := func() string {
		candidates, _ := enumerateGuards(context.Background(), "")
		return candidates[0].fallback
	}
	oldEnumerate, oldInspect := enumerateGuards, inspectGuard
	t.Cleanup(func() { enumerateGuards, inspectGuard = oldEnumerate, oldInspect })
	enumerateGuards = marker("prior")
	started, release, finished := make(chan struct{}), make(chan struct{}), make(chan struct{})
	var once sync.Once
	releaseWorker := func() { once.Do(func() { close(release) }) }
	var observed string
	t.Run("held worker", func(t *testing.T) {
		priorEnumerate, priorInspect := enumerateGuards, inspectGuard
		t.Cleanup(func() { enumerateGuards, inspectGuard = priorEnumerate, priorInspect })
		enumerateGuards = marker("replacement")
		inspectGuard = func(context.Context, string, candidate) [][]string {
			close(started)
			<-release
			observed = installed()
			return nil
		}
		go func() { defer close(finished); Scan(context.Background(), "/repo") }()
		t.Cleanup(func() { releaseWorker(); <-finished })
		<-started
		if got := installed(); got != "replacement" {
			t.Fatalf("while the worker is live the installed collaborator is %q, want the replacement", got)
		}
	})
	releaseWorker()
	<-finished
	if observed != "replacement" {
		t.Fatalf("the held worker observed %q after its release, want the replacement it ran with", observed)
	}
	if got := installed(); got != "prior" {
		t.Fatalf("after the subtest cleanup the installed collaborator is %q, want the prior value", got)
	}
}
