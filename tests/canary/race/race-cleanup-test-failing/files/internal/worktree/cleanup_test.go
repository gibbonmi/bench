package worktree

import (
	"sync"
	"testing"
)

// recorded stands in for the cleanup transaction log. Two goroutines append to it with
// nothing ordering them, which is a genuine data race: the detector reports it from the
// vector clocks alone, so the fixture bites deterministically rather than on a timing
// window.
var recorded []string

// The race phase probes for a registered declaration before it materializes, and its
// -run filter exits 0 on a package that never declares it, so the fixture has to carry
// the real name. The assertion itself passes: the plain test phase runs this same test
// without -race and has to stay green there, or the fixture would bite from a phase that
// does not own it and the race phase could be deleted unnoticed.
func TestConcurrentCleanupRecordsOneTransaction(t *testing.T) {
	var wg sync.WaitGroup
	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			recorded = append(recorded, "cleanup")
		}()
	}
	wg.Wait()
	if len(recorded) == 0 {
		t.Fatal("no cleanup transaction was recorded")
	}
}
