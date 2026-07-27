package worktree

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/intent"
)

// cleanLineFor is the retirement command a reader is expected to paste for path. The
// expected quoting is written out here rather than taken from the renderer's own helper,
// so a change to that helper cannot make these tests agree with it by construction. Every
// path a test hands it is free of quotes, which is what keeps the literal this simple.
func cleanLineFor(path string) string { return "bench worktree clean '" + path + "'" }

func ledgerIDs(t *testing.T, root string) string {
	t.Helper()
	assignments, err := intent.Assignments(root)
	mustNoError(t, err)
	ids := make([]string, 0, len(assignments))
	for _, a := range assignments {
		ids = append(ids, a.ID)
	}
	sort.Strings(ids)
	return strings.Join(ids, ",")
}

func newSweepRepo(t *testing.T) string {
	t.Helper()
	root := newWorktreeRepo(t)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	return root
}

func mustSweep(t *testing.T, root string) ResumeResult {
	t.Helper()
	result, err := ConservativeCleanup(root)
	mustNoError(t, err)
	return result
}

func TestResumeSummaryNamesCleanCommand(t *testing.T) {
	root := newSweepRepo(t)
	orphan := mustCreate(t, root, "summary-clean-command", "aged, tree present")
	backdate(t, root, orphan.Assignment, 8*24*time.Hour)

	summary := renderResumeSummary(mustSweep(t, root))
	requireTest(t, strings.Contains(summary, cleanLineFor(orphan.Path)),
		"summary names no retirement command for the orphan at %s:\n%s", orphan.Path, summary)
	requireTest(t, strings.Contains(summary, "--apply"),
		"summary presents the clean command as a single step, hiding the plan/apply pair:\n%s", summary)
}

// TestResumeSummaryReportsOrphanWithIgnoredResidue holds the sweep to the ledger rather
// than to a cleanup plan. Ignored build output is the normal state of a worktree a shift
// ran in, and PlanAutomatic returns at that retain reason long before it reaches the
// assignment's state, so a sweep reading the plan's reason code reports nothing for
// exactly the population this listing exists for.
func TestResumeSummaryReportsOrphanWithIgnoredResidue(t *testing.T) {
	root := newSweepRepo(t)
	mustWrite(t, filepath.Join(root, ".gitignore"), []byte("ignored.txt\n"), 0o644)
	gitRun(t, root, "add", ".gitignore")
	gitRun(t, root, "commit", "-qm", "ignore")
	orphan := mustCreate(t, root, "summary-ignored-residue", "aged, ignored residue")
	mustWrite(t, filepath.Join(orphan.Path, "ignored.txt"), []byte("residue\n"), 0o644)
	backdate(t, root, orphan.Assignment, 8*24*time.Hour)

	plan, err := PlanAutomatic(root, orphan.Path)
	mustNoError(t, err)
	requireTest(t, plan.ReasonCode == ReasonIgnored, "fixture does not retain on ignored residue: reason %q", plan.ReasonCode)

	summary := renderResumeSummary(mustSweep(t, root))
	requireTest(t, strings.Contains(summary, cleanLineFor(orphan.Path)),
		"summary drops the orphan whose plan retains for an earlier reason:\n%s", summary)
}

// TestResumeSummaryNeverSuggestsDiscardIgnored keeps the emitted remedy off the flag
// whose request-less form orphans the assignment (FT93b) — suggesting it would
// manufacture the next generation of the residue this listing reports.
func TestResumeSummaryNeverSuggestsDiscardIgnored(t *testing.T) {
	root := newSweepRepo(t)
	orphan := mustCreate(t, root, "summary-no-discard", "aged, tree present")
	backdate(t, root, orphan.Assignment, 8*24*time.Hour)

	summary := renderResumeSummary(mustSweep(t, root))
	requireTest(t, strings.Contains(summary, cleanLineFor(orphan.Path)),
		"summary names no retirement command, so its wording is untested:\n%s", summary)
	requireTest(t, !strings.Contains(summary, "--discard-ignored"),
		"summary suggests --discard-ignored:\n%s", summary)
}

// TestSweepCompactsOrphanedActiveResidue covers both tree-gone verdicts for an orphan:
// one git no longer registers is the sweep's to compact, and one git still registers is
// the prune path's, because compacting a registration git is about to prune leaves the
// ledger and the registration disagreeing.
func TestSweepCompactsOrphanedActiveResidue(t *testing.T) {
	root := newSweepRepo(t)
	residue := mustCreate(t, root, "orphan-residue-gone", "aged, tree gone, unregistered")
	prunable := mustCreate(t, root, "orphan-residue-prunable", "aged, tree gone, still registered")
	backdate(t, root, residue.Assignment, 8*24*time.Hour)
	backdate(t, root, prunable.Assignment, 8*24*time.Hour)
	gitRun(t, root, "worktree", "remove", "-f", "-f", residue.Path)
	mustNoError(t, os.RemoveAll(prunable.Path))

	result := mustSweep(t, root)
	requireTest(t, result.Reconciled == 1, "Reconciled=%d, want 1", result.Reconciled)
	if _, err := assignmentByID(root, residue.Assignment.ID); err == nil {
		t.Fatal("orphaned residue record survived the sweep")
	}
	if _, err := assignmentByID(root, prunable.Assignment.ID); err != nil {
		t.Fatalf("orphaned record git still registers was compacted, want the prune path to own it: %v", err)
	}
}

// TestSweepPreservesOrphanedRecoveryRecords pins the boundary of the sweep's one
// destructive branch. No record is ever both orphaned and holding preserved work —
// orphanhood requires state active, and ValidateAssignment refuses an active record with
// recovery metadata on every ledger read — so residualAssignment is the single guard
// standing between the compaction branch and a pointer to preserved work, and the closest
// reachable record is driven through the sweep to hold it there.
func TestSweepPreservesOrphanedRecoveryRecords(t *testing.T) {
	root := newSweepRepo(t)
	preserved := mustCreate(t, root, "orphan-preserved", "aged, holds preserved work")
	a := preserved.Assignment
	ref := intent.RecoveryRefPrefix(a.OwnerID, a.ID) + "1"
	recovery := []intent.Recovery{{Ref: ref, Root: strings.Repeat("a", 40), Payloads: []string{strings.Repeat("b", 40)}}}

	unreachable := a
	unreachable.CreatedAt, unreachable.Recovery = stampedAt(time.Now().Add(-8*24*time.Hour)), recovery
	requireTest(t, intent.PutAssignment(root, unreachable) != nil,
		"the ledger accepted an active record holding recovery metadata, so an orphan can hold preserved work")

	gitRun(t, root, "worktree", "remove", "-f", "-f", preserved.Path)
	a.State, a.Recovery, a.CreatedAt = intent.StateRecovered, recovery, stampedAt(time.Now().Add(-8*24*time.Hour))
	mustNoError(t, intent.PutAssignment(root, a))
	requireTest(t, !orphaned(a, time.Now()), "an aged recovered record reads as orphaned")

	result := mustSweep(t, root)
	got, err := assignmentByID(root, a.ID)
	requireTest(t, err == nil && len(got.Recovery) == 1, "the record holding preserved work was compacted: %v", err)
	requireTest(t, len(result.Preserved) == 1 && result.Preserved[0].Ref == ref,
		"Preserved=%+v, want one entry at %s", result.Preserved, ref)
}

// TestSweepIsIdempotent pins what a reader actually sees: the summary reprints at every
// session start, so once the ledger has settled, repeated sweeps must agree exactly and
// delete nothing further. The first sweep is the settling run — it compacts the residue
// and so legitimately differs — and the two runs after it are the compared pair.
func TestSweepIsIdempotent(t *testing.T) {
	root := newSweepRepo(t)
	present := mustCreate(t, root, "idempotent-orphan", "aged, tree present")
	gone := mustCreate(t, root, "idempotent-residue", "aged, tree gone")
	backdate(t, root, present.Assignment, 8*24*time.Hour)
	backdate(t, root, gone.Assignment, 8*24*time.Hour)
	gitRun(t, root, "worktree", "remove", "-f", "-f", gone.Path)

	requireTest(t, mustSweep(t, root).Reconciled == 1, "settling sweep did not compact the residue")

	second := renderResumeSummary(mustSweep(t, root))
	before := ledgerIDs(t, root)
	third := renderResumeSummary(mustSweep(t, root))
	after := ledgerIDs(t, root)

	requireTest(t, second == third, "settled sweeps disagree:\nsecond:\n%s\nthird:\n%s", second, third)
	requireTest(t, strings.Contains(second, cleanLineFor(present.Path)),
		"a settled sweep stops reporting the standing orphan:\n%s", second)
	requireTest(t, before == after, "a settled sweep changed the ledger: %s -> %s", before, after)
}
