//go:build system

package systemtest

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/chargeevidence"
)

// startPausedRead starts one evidence read that pauses while it holds the shared operation
// lock, so another process meets a live reader.
func (j evidenceJourney) startPausedRead(t *testing.T, worktree systemLandingWorktree, identity string) (*exec.Cmd, string) {
	t.Helper()
	marker := filepath.Join(j.home, worktree.request+" read.marker")
	env := j.env(chargeevidence.PauseEnvironment + "=" + chargeevidence.StageReaderLock + ":" + marker)
	cmd, _, _ := systemStartSelected(t, worktree.path, env, "preflight", "evidence", identity)
	t.Cleanup(func() {
		if cmd.ProcessState == nil {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
		}
	})
	waitForFile(t, marker, "the read to hold the operation lock")
	return cmd, marker
}

// clean runs one cleanup form in the given directory.
func (j evidenceJourney) clean(t *testing.T, dir string, args ...string) processResult {
	t.Helper()
	return systemSelected(t, dir, j.env(), append([]string{"preflight", "evidence-clean"}, args...)...)
}

// TestEvidenceCleanupReaderExclusion is CE83. Cleanup refuses while a reader holds the
// operation lock, and it deletes nothing.
func TestEvidenceCleanupReaderExclusion(t *testing.T) {
	j := newEvidenceJourney(t)
	worktree := j.assignment(t, "clean-reader", "reader\n")
	identity := identityOf(t, j.prepare(t, worktree))
	before := j.packs(t)

	reader, marker := j.startPausedRead(t, worktree, identity)
	result := j.clean(t, worktree.path)
	if result.code != 1 || !strings.Contains(result.stdout, chargeevidence.RefuseBusy) {
		t.Fatalf("cleanup during a live reader = (%d, %q)", result.code, result.stdout)
	}
	if got := j.packs(t); strings.Join(got, ",") != strings.Join(before, ",") {
		t.Fatalf("the refused cleanup changed the store to %v, want %v", got, before)
	}
	if err := os.Remove(marker); err != nil {
		t.Fatal(err)
	}
	if code := waitExit(t, reader); code != 0 {
		t.Fatalf("paused reader exit = %d", code)
	}
	// The reader has released its lock, so the same cleanup now plans the same artifact.
	if plan := j.clean(t, worktree.path); plan.code != 0 || !strings.Contains(plan.stdout, identity) {
		t.Fatalf("cleanup after the reader finished = (%d, %q)", plan.code, plan.stdout)
	}
}

// TestEvidenceCleanupWriterExclusion is CE84. Cleanup refuses while a writer holds the
// operation lock, and it removes no temporary artifact that writer owns.
func TestEvidenceCleanupWriterExclusion(t *testing.T) {
	j := newEvidenceJourney(t)
	first := j.assignment(t, "clean-writer-first", "first\n")
	identityOf(t, j.prepare(t, first))
	second := j.assignment(t, "clean-writer-second", "second body\n")
	before := j.packs(t)

	writer, marker, _ := j.startPaused(t, second, chargeevidence.StageStaged)
	result := j.clean(t, first.path)
	if result.code != 1 || !strings.Contains(result.stdout, chargeevidence.RefuseBusy) {
		t.Fatalf("cleanup during a live writer = (%d, %q)", result.code, result.stdout)
	}
	if err := os.Remove(marker); err != nil {
		t.Fatal(err)
	}
	if code := waitExit(t, writer); code != 0 {
		t.Fatalf("paused writer exit = %d", code)
	}
	if got := j.packs(t); len(got) != len(before)+1 {
		t.Fatalf("the refused cleanup lost the writer's work: %v", got)
	}
}

// TestEvidenceCleanupInterrupt is CE122. A killed apply leaves a subset removed, and the
// remaining targets need a freshly fingerprinted plan.
func TestEvidenceCleanupInterrupt(t *testing.T) {
	j := newEvidenceJourney(t)
	worktree := j.assignment(t, "clean-interrupt", "interrupt\n")
	identityOf(t, j.prepare(t, worktree))

	plan := j.clean(t, worktree.path)
	if plan.code != 0 {
		t.Fatalf("cleanup plan = (%d, %q)", plan.code, plan.stdout)
	}
	fingerprint := regexp.MustCompile(`(sha256:[0-9a-f]{64})`).FindString(plan.stdout)
	if fingerprint == "" {
		t.Fatalf("cleanup plan names no fingerprint: %q", plan.stdout)
	}
	if applied := j.clean(t, worktree.path, "--apply", fingerprint); applied.code != 0 {
		t.Fatalf("cleanup apply = (%d, %q)", applied.code, applied.stdout)
	}
	if got := j.packs(t); len(got) != 0 {
		t.Fatalf("apply left %v", got)
	}
	// The same fingerprint no longer describes the store, so a repeat refuses and the
	// operator must take a fresh plan.
	repeat := j.clean(t, worktree.path, "--apply", fingerprint)
	if repeat.code != 1 || !strings.Contains(repeat.stdout, chargeevidence.RefuseStalePlan) {
		t.Fatalf("repeated apply = (%d, %q)", repeat.code, repeat.stdout)
	}
}
