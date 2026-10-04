package commitment_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/commitment/commitcmd"
	"github.com/gibbonmi/bench/internal/intent"
)

func TestCommitmentConcurrentStarts(t *testing.T) {
	root, first, second := admissionRepo(t, false, false)
	mustCommand(t, first, "block", "--outcome", "A", "--reason", "Waiting")
	realGit, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	barrier := t.TempDir()
	shim := `#!/bin/sh
if [ "$2" = "$DC_PAUSED_ROOT" ] && [ "$3" = ls-tree ] && [ "$7" = specs/B/spec.md ]; then
 : > "$DC_BARRIER/entered"
 while [ ! -f "$DC_BARRIER/release" ]; do sleep 0.01; done
fi
exec "$DC_REAL_GIT" "$@"
`
	if err := os.WriteFile(filepath.Join(barrier, "git"), []byte(shim), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DC_REAL_GIT", realGit)
	t.Setenv("DC_BARRIER", barrier)
	t.Setenv("DC_PAUSED_ROOT", second)
	t.Setenv("PATH", barrier+string(os.PathListSeparator)+os.Getenv("PATH"))
	release := func() {
		if err := os.WriteFile(filepath.Join(barrier, "release"), nil, 0600); err != nil {
			t.Error(err)
		}
	}
	defer release()
	type result struct {
		output string
		code   int
	}
	bDone := make(chan result, 1)
	go func() { out, code := startOutcome(second, "B", "second"); bDone <- result{out, code} }()
	waitFile(t, filepath.Join(barrier, "entered"))
	ledgerPath, err := intent.Address(root)
	if err != nil {
		t.Fatal(err)
	}
	_, lockErr := os.Stat(ledgerPath + ".lock")
	aDone := make(chan result, 1)
	go func() {
		out, code := commitcmd.Command(first, []string{"unblock", "--outcome", "A"})
		if code == 0 {
			out, code = startOutcome(first, "A", "first")
		}
		aDone <- result{out, code}
	}()
	var a result
	// A held transaction must finish before the unblock can run. An unprotected
	// snapshot lets the unblock and A complete before B publishes that snapshot.
	if lockErr == nil {
		release()
		a = awaitStart(t, aDone)
	} else if os.IsNotExist(lockErr) {
		a = awaitStart(t, aDone)
		release()
	} else {
		t.Fatal(lockErr)
	}
	b := awaitStart(t, bDone)
	successes := 0
	for _, item := range []result{a, b} {
		if item.code == 0 {
			successes++
		}
	}
	state, err := intent.Read(root)
	if err != nil {
		t.Fatal(err)
	}
	if successes != 1 || state.Commitment == nil || len(state.Commitment.Claims) != 1 || len(state.Commitment.Bindings) != 1 {
		t.Fatalf("competing starts: A=%+v B=%+v state=%+v", a, b, state.Commitment)
	}
}

func awaitStart[T any](t *testing.T, done <-chan T) T {
	t.Helper()
	select {
	case value := <-done:
		return value
	case <-time.After(bounds.TestDeadline(bounds.IntentLockTimeout)):
		t.Fatal("start did not finish")
	}
	var zero T
	return zero
}

func waitFile(t *testing.T, path string) {
	t.Helper()
	deadline := time.Now().Add(bounds.TestDeadline(bounds.IntentLockTimeout))
	for time.Now().Before(deadline) {
		if _, err := os.Stat(path); err == nil {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("transaction did not reach %s", path)
}
