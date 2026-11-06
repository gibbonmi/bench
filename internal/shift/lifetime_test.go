package shift

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/otelrecord"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
	"github.com/gibbonmi/bench/internal/worktree"
)

// resistantAdapter installs a sleeping adapter that ignores every catchable shutdown
// signal, and it returns the adapter and the file that marks its start.
func resistantAdapter(t *testing.T) (agent, started string) {
	t.Helper()
	started = sleepingAgent(t, "trap '' INT TERM HUP\n")
	return os.Getenv("BENCH_AGENT"), started
}

// adapterReturn waits for the adapter run on done inside the window that inner derives.
func adapterReturn(t *testing.T, done <-chan error, inner time.Duration, wait string) error {
	t.Helper()
	window := bounds.TestDeadline(inner)
	select {
	case err := <-done:
		return err
	case <-time.After(window):
		t.Fatal(bounds.TestTimeoutVerdict(wait, window))
	}
	return nil
}

// protectWorktree registers one pending user of the checkout at dir. The test process is
// the user's live actor, so the user stays unresolved until the test finishes it.
func protectWorktree(t *testing.T, dir string) *protection.Registration {
	t.Helper()
	p := &protection.Protocol{Validators: protection.Validators{protection.KindWorktree: worktree.ValidateProtection}, ReadLimit: bounds.ControlRecordLimit}
	d, err := worktree.Descriptor(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(d.Store) })
	opened, err := p.OpenResource(d)
	if err != nil {
		t.Fatal(err)
	}
	user, err := p.Prepare([]protection.Descriptor{opened})
	if err != nil {
		t.Fatal(err)
	}
	return user
}

// TestActualStartRace is PL83. The test pauses the production startup transition of an
// adapter run after its lifetime is published and before its Start. A pulled line then
// arrives through the signal path, and the resistant adapter must never run: the run returns
// the no-start result at once. A cancellation that raced process publication would be lost,
// and the adapter would outlive the window.
func TestActualStartRace(t *testing.T) {
	agent, marker := resistantAdapter(t)
	s := &session{agent: agent, root: t.TempDir(), stdout: io.Discard, stderr: io.Discard}
	paused, resume := make(chan struct{}), make(chan struct{})
	armFault(t, func(step shiftStep) error {
		if step == stepAdapterStart {
			close(paused)
			<-resume
		}
		return nil
	})
	done, finished := make(chan error, 1), make(chan struct{})
	go func() { defer close(finished); done <- s.runAdapter(t.Context(), "objective") }()
	// A failed assertion still releases the held fault and joins the run before the
	// fault cleanup restores the production fault point. The test context is cancelled
	// before this cleanup runs, so the run returns.
	resumed := false
	t.Cleanup(func() {
		if !resumed {
			close(resume)
		}
		<-finished
	})
	select {
	case <-paused:
	case <-time.After(helperWindow):
		t.Fatal("the adapter run never reached its startup transition")
	}
	s.interrupt(syscall.SIGINT)
	resumed = true
	close(resume)
	if err := adapterReturn(t, done, 0, "the adapter run to return after its startup cancellation"); !errors.Is(err, subprocess.ErrNoStart) {
		t.Fatalf("adapter run = %v, want the no-start result", err)
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("the resistant adapter ran after its startup cancellation: %v", err)
	}
}

// TestDeadlinePrecedence is PL84. A pulled line reaches a running resistant adapter, and
// the wall deadline fires while that shutdown is in flight. The adapter group dies at the
// end of the grace, and the checkpoint still resolves the deadline's incomplete outcome,
// not interrupted/130.
func TestDeadlinePrecedence(t *testing.T) {
	agent, marker := resistantAdapter(t)
	s := &session{agent: agent, root: t.TempDir(), stdout: io.Discard, stderr: io.Discard}
	done := make(chan error, 1)
	go func() { done <- s.runAdapter(t.Context(), "objective") }()
	waitForFile(t, marker)
	s.interrupt(syscall.SIGINT)
	s.expire()
	err := adapterReturn(t, done, bounds.ProcessShutdownGrace+bounds.ProcessFinalWindow, "the resistant adapter to stop")
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != -1 {
		t.Fatalf("adapter run = %v, want a signaled exit", err)
	}
	outcome, detail, ok := s.checkpointOutcome()
	if !ok || outcome != OutcomeIncomplete || detail != "wall deadline exceeded" {
		t.Fatalf("checkpoint = (%q, %q, %v), want the deadline's incomplete outcome", outcome, detail, ok)
	}
}

// TestCheckpointTeardownFaultKeepsResponse covers the checkpoint's teardown failure. The wall
// deadline stops a sleeping adapter, and the checkpoint's teardown fails. Loop must return
// failed/1 to its caller with the shift_result block written, so a bounded caller that
// prints the response after Loop returns keeps it.
func TestCheckpointTeardownFaultKeepsResponse(t *testing.T) {
	faultFixtureCore(t, greenGate, nil)
	sleepingAgent(t, "")
	t.Setenv("BENCH_MAX_WALL", "1s")
	armFault(t, func(step shiftStep) error {
		if step == stepTeardown {
			return fmt.Errorf("injected teardown failure")
		}
		return nil
	})

	var stdout, stderr bytes.Buffer
	code := Loop("delivery", "checkpoint teardown fault", &stdout, &stderr)

	out := stdout.String()
	if code != 1 || !strings.Contains(out, "shift_result[1]") || !strings.Contains(out, "\n  failed,") || !strings.Contains(out, "teardown failed") {
		t.Fatalf("Loop past its deadline with a teardown fault = (%d, %q, %q), want failed/1 with the shift_result block", code, out, stderr.String())
	}
}

// TestProtectedRollback is PL117. A rollback of a worktree whose user is unresolved refuses
// before any reset or clean, so the tracked change and the untracked sentinel that the user
// can still change stay. Once the user resolves, the rollback discards both.
func TestProtectedRollback(t *testing.T) {
	root := t.TempDir()
	gittest.MustRun(t, root, "init", "-q", "-b", "main")
	tracked, sentinel := filepath.Join(root, "tracked.txt"), filepath.Join(root, "sentinel")
	if err := os.WriteFile(tracked, []byte("base\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gittest.MustRun(t, root, "add", "tracked.txt")
	gittest.MustRun(t, root, append(gittest.IdentityArgs(), "commit", "-qm", "base")...)
	for path, body := range map[string]string{tracked: "changed\n", sentinel: "child output\n"} {
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	user := protectWorktree(t, root)

	if err := rollback(root); !errors.Is(err, protection.ErrRetained) {
		t.Fatalf("rollback beside an unresolved user = %v, want a retention", err)
	}
	if body, _ := os.ReadFile(tracked); string(body) != "changed\n" {
		t.Fatalf("the refused rollback reset the tracked change to %q", body)
	}
	if _, err := os.Stat(sentinel); err != nil {
		t.Fatalf("the refused rollback cleaned the sentinel: %v", err)
	}

	if err := user.Finish(protection.Terminal{NoStart: true}); err != nil {
		t.Fatal(err)
	}
	if err := rollback(root); err != nil {
		t.Fatalf("rollback after the user resolved = %v", err)
	}
	if body, _ := os.ReadFile(tracked); string(body) != "base\n" {
		t.Fatalf("the rollback left the tracked change %q", body)
	}
	if _, err := os.Stat(sentinel); !os.IsNotExist(err) {
		t.Fatalf("the rollback kept the sentinel: %v", err)
	}
}

// TestRecoverRetainsProtectedCleanScratch is the recovery owner's protected-scratch witness.
// A killed shift left only scratch, so its dirty set is empty, but a user of its worktree is
// unresolved. The recovery pass retains and locks the worktree with its notes, and it records
// the worktree pointer.
func TestRecoverRetainsProtectedCleanScratch(t *testing.T) {
	root, entry := crashedShift(t, false, true)
	user := protectWorktree(t, entry.Worktree)
	t.Cleanup(func() { _ = user.Finish(protection.Terminal{NoStart: true}) })
	Recover(root, io.Discard)
	if _, err := os.Stat(filepath.Join(entry.Worktree, notesFile)); err != nil {
		t.Fatalf("the recovery removed the protected scratch: %v", err)
	}
	if !locked(t, root, entry.Worktree) {
		t.Fatalf("worktree %s, want it locked", worktreeState(t, root, entry))
	}
	if got := ledgerEntry(t, root, entry.Key); got.Recovery != recoveryWorktree(entry.Worktree) {
		t.Fatalf("entry = %+v, want the retained worktree pointer", got)
	}
	if spans := recoverySpans(t, root); len(spans) != 1 {
		t.Fatalf("the record holds %d recovery spans, want 1", len(spans))
	} else {
		requireAttr(t, spans[0], otelrecord.AttrCleanup, otelrecord.CleanupRetained)
	}
}
