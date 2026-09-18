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

// startPausedAt starts one Bench process that pauses at stage, and waits for the marker that
// pause writes. This is the one pause harness the cases here share: the stage names the
// marker, so no case carries a marker name of its own.
func (j evidenceJourney) startPausedAt(t *testing.T, worktree systemLandingWorktree, stage string, args ...string) (*exec.Cmd, string, func() string) {
	t.Helper()
	marker := filepath.Join(j.home, worktree.request+" "+stage+".marker")
	cmd, stdout, _ := systemStartSelected(t, worktree.path, j.env(chargeevidence.PauseEnvironment+"="+stage+":"+marker), args...)
	t.Cleanup(func() {
		if cmd.ProcessState == nil {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
		}
	})
	waitForFile(t, marker, "the process to reach stage "+stage)
	return cmd, marker, stdout.String
}

// startPausedRead starts one evidence read that pauses while it holds the shared operation
// lock, so another process meets a live reader.
func (j evidenceJourney) startPausedRead(t *testing.T, worktree systemLandingWorktree, identity string) (*exec.Cmd, string) {
	t.Helper()
	cmd, marker, _ := j.startPausedAt(t, worktree, chargeevidence.StageReaderLock, "preflight", "evidence", identity)
	return cmd, marker
}

// cleanArgs is the one cleanup invocation the cases here run.
func cleanArgs(args ...string) []string {
	return append([]string{"preflight", "evidence-clean"}, args...)
}

// clean runs one cleanup form in the given directory.
func (j evidenceJourney) clean(t *testing.T, dir string, args ...string) processResult {
	t.Helper()
	return systemSelected(t, dir, j.env(), cleanArgs(args...)...)
}

// planFingerprint reads the fingerprint one cleanup plan response committed to.
var planFingerprint = regexp.MustCompile(`(sha256:[0-9a-f]{64})`)

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

// TestEvidenceCleanupInterrupt is CE122. A killed apply leaves the targets it already removed
// removed and reports nothing, so the plan it was applying no longer describes the store:
// that fingerprint refuses, and only a fresh plan covers the targets that survived.
func TestEvidenceCleanupInterrupt(t *testing.T) {
	j := newEvidenceJourney(t)
	worktree := j.assignment(t, "clean-interrupt", "interrupt\n")
	identityOf(t, j.prepare(t, worktree))
	// A second target makes the interruption observable. The apply pauses after each
	// removal, so the kill at the first pause leaves the target it had not reached.
	orphan := chargeevidence.TempPrefix + "interrupt" + chargeevidence.TempSuffix
	if err := os.WriteFile(filepath.Join(j.store(), orphan), []byte("partial bytes\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	plan := j.clean(t, worktree.path)
	fingerprint := planFingerprint.FindString(plan.stdout)
	if plan.code != 0 || fingerprint == "" || len(j.packs(t)) != 2 {
		t.Fatalf("cleanup plan = (%d, %q) over %v", plan.code, plan.stdout, j.packs(t))
	}

	apply, _, _ := j.startPausedAt(t, worktree, chargeevidence.StageRemoved, cleanArgs("--apply", fingerprint)...)
	if err := apply.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	_ = apply.Wait()
	if got := j.packs(t); len(got) != 1 || got[0] == orphan {
		t.Fatalf("the killed apply left %v, want only the target it had not reached", got)
	}
	// The interrupted plan is spent: it named a target the store no longer holds.
	repeat := j.clean(t, worktree.path, "--apply", fingerprint)
	if repeat.code != 1 || !strings.Contains(repeat.stdout, chargeevidence.RefuseStalePlan) {
		t.Fatalf("apply of the interrupted plan = (%d, %q)", repeat.code, repeat.stdout)
	}
	fresh := j.clean(t, worktree.path)
	freshFingerprint := planFingerprint.FindString(fresh.stdout)
	if fresh.code != 0 || freshFingerprint == "" || freshFingerprint == fingerprint {
		t.Fatalf("fresh plan = (%d, %q)", fresh.code, fresh.stdout)
	}
	if applied := j.clean(t, worktree.path, "--apply", freshFingerprint); applied.code != 0 || len(j.packs(t)) != 0 {
		t.Fatalf("apply of the fresh plan = (%d, %q), store %v", applied.code, applied.stdout, j.packs(t))
	}
}
