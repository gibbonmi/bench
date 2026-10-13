package gate

import (
	"bytes"
	"context"
	"os"
	"strings"
	"syscall"
	"testing"
)

func TestExecutionInProgressReadsTheLock(t *testing.T) {
	root := failureOutcomeFixture(t)

	if held, err := ExecutionInProgress(root); err != nil || held {
		t.Fatalf("fresh repository = %v, %v, want false, nil", held, err)
	}

	t.Run("while a child process holds the lock", func(t *testing.T) {
		startGateLockHolder(t, root)
		held, err := ExecutionInProgress(root)
		if err != nil || !held {
			t.Fatalf("under holder = %v, %v, want true, nil", held, err)
		}
	})

	if held, err := ExecutionInProgress(root); err != nil || held {
		t.Fatalf("after release = %v, %v, want false, nil", held, err)
	}

	if held, err := ExecutionInProgress(t.TempDir()); err == nil {
		t.Fatalf("outside a repository = %v, nil, want an error", held)
	}
}

func TestGateRunRefusesNonContentionLockErrorWithCause(t *testing.T) {
	root := outcomeFixture(t)
	previous := takeExecutionLock
	takeExecutionLock = func(*os.File) error { return syscall.ENOLCK }
	t.Cleanup(func() { takeExecutionLock = previous })

	var stdout, stderr bytes.Buffer
	result := Execute(context.Background(), root, &stdout, &stderr)
	if result.ActionExit != 1 {
		t.Fatalf("result = %#v, stderr=%q, want ActionExit 1", result, stderr.String())
	}
	want := "gate lock unavailable: " + syscall.ENOLCK.Error() + "\n"
	if got := stderr.String(); got != want || strings.Contains(got, "in progress") {
		t.Fatalf("stderr = %q, want %q and no in-progress message", got, want)
	}
	if _, err := os.Stat(runCountWitness.path(t, root)); !os.IsNotExist(err) {
		t.Fatalf("run counter stat = %v, want no oracle run after a lock error", err)
	}
}
