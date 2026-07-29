package worktree

import (
	"os"
	"sync"
	"testing"
)

// recorded stands in for the cleanup transaction log. Two goroutines append to it with
// nothing ordering them, which is a genuine data race: the detector reports it from the
// vector clocks alone, so the fixture bites deterministically rather than on a timing
// window.
var recorded []string

func TestMain(m *testing.M) {
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
		panic("no cleanup transaction was recorded")
	}
	os.Exit(m.Run())
}
