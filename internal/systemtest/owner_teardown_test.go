//go:build system

package systemtest

import (
	"errors"
	"fmt"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
)

func TestRedTimeoutAndDescendantTeardown(t *testing.T) {
	if red := owner.runSelected(owner.repos[1], "not-a-command"); red.code == 0 {
		t.Fatal("unknown command unexpectedly green")
	}
	owner.markTerminal("red")
	interrupted, err := owner.interruptProcessGroup()
	if err != nil {
		t.Fatal(err)
	}
	requireProcessGone(t, interrupted, "interrupt")
	owner.markTerminal("interrupt")
	descendant, err := owner.timeoutProcessGroup(100 * time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	requireProcessGone(t, descendant, "timeout")
	owner.markTerminal("timeout")
}

// requireProcessGone waits for a signalled descendant to leave the process table. The
// signal is already delivered when the wait starts, so the wait contains no window of its
// own and takes the floor the derivation gives a zero bound.
func requireProcessGone(t *testing.T, pid int, outcome string) {
	t.Helper()
	window := bounds.TestDeadline(0)
	deadline := time.Now().Add(window)
	for {
		err := syscall.Kill(pid, 0)
		if errors.Is(err, syscall.ESRCH) {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal(bounds.TestTimeoutVerdict(fmt.Sprintf("descendant %d to exit after process-group %s", pid, outcome), window))
		}
		time.Sleep(10 * time.Millisecond)
	}
}
