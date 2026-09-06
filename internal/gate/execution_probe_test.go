package gate

import (
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
