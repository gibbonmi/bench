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
// expected quoting is written out here rather than taken from the renderer's own helper.
// A change to that helper cannot then make these tests agree with it by construction.
// Every path a test hands it is free of quotes, which is what keeps the literal this simple.
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
	bindEnv(t, "BENCH_HOME", filepath.Join(root, ".bench-home"))
	return root
}

func mustSweep(t *testing.T, root string) ResumeResult {
	t.Helper()
	result, err := ConservativeCleanup(root)
	mustNoError(t, err)
	return result
}

func makeUnlandedAssignment(t *testing.T, creation Creation) {
	t.Helper()
	commitInWorktree(t, creation.Path, "unlanded.txt", "unlanded\n", "unlanded")
}

func TestResumeSummaryNamesCleanCommand(t *testing.T) {
	root := newSweepRepo(t)
	orphan := mustCreate(t, root, Home(), "summary-clean-command", "aged, tree present")
	makeUnlandedAssignment(t, orphan)
	backdate(t, root, orphan.Assignment, 8*24*time.Hour)

	summary := renderResumeSummary(mustSweep(t, root))
	requireTest(t, strings.Contains(summary, cleanLineFor(orphan.Path)),
		"summary names no retirement command for the orphan at %s:\n%s", orphan.Path, summary)
	requireTest(t, strings.Contains(summary, "--apply"),
		"summary presents the clean command as a single step, hiding the plan/apply pair:\n%s", summary)
}

// TestResumeSummaryReportsOrphanWithIgnoredResidue holds the sweep to the ledger rather
// than to a cleanup plan. Ignored build output is the normal state of a worktree a shift
// ran in.
// PlanAutomatic returns at that retain reason long before it reaches the assignment's state.
// A sweep that reads the plan's reason code therefore reports nothing for exactly the
// population this listing exists for.
func TestResumeSummaryReportsOrphanWithIgnoredResidue(t *testing.T) {
	root := newSweepRepo(t)
	mustWrite(t, filepath.Join(root, ".gitignore"), []byte("ignored.txt\n"), 0o644)
	gitRun(t, root, "add", ".gitignore")
	gitRun(t, root, "commit", "-qm", "ignore")
	orphan := mustCreate(t, root, Home(), "summary-ignored-residue", "aged, ignored residue")
	makeUnlandedAssignment(t, orphan)
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
// whose request-less form orphans the assignment (FT93b).
// That suggestion would manufacture the next generation of the residue this listing reports.
func TestResumeSummaryNeverSuggestsDiscardIgnored(t *testing.T) {
	root := newSweepRepo(t)
	orphan := mustCreate(t, root, Home(), "summary-no-discard", "aged, tree present")
	makeUnlandedAssignment(t, orphan)
	backdate(t, root, orphan.Assignment, 8*24*time.Hour)

	summary := renderResumeSummary(mustSweep(t, root))
	requireTest(t, strings.Contains(summary, cleanLineFor(orphan.Path)),
		"summary names no retirement command, so its wording is untested:\n%s", summary)
	requireTest(t, !strings.Contains(summary, "--discard-ignored"),
		"summary suggests --discard-ignored:\n%s", summary)
}

// TestSweepCompactsOrphanedActiveResidue covers both tree-gone verdicts for an orphan.
// One git no longer registers is the sweep's to compact, and one git still registers is
// the prune path's.
// A registration is not compacted while git is about to prune it, because that would
// leave the ledger and the registration disagreeing.
func TestSweepCompactsOrphanedActiveResidue(t *testing.T) {
	root := newSweepRepo(t)
	residue := mustCreate(t, root, Home(), "orphan-residue-gone", "aged, tree gone, unregistered")
	prunable := mustCreate(t, root, Home(), "orphan-residue-prunable", "aged, tree gone, still registered")
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

// unstamp rewrites a record with no creation stamp. created_at is omitempty, so this is
// the shape on disk of every assignment cut before the field existed: the key is absent
// rather than null.
func unstamp(t *testing.T, root string, assignment intent.Assignment) {
	t.Helper()
	assignment.CreatedAt = nil
	mustNoError(t, intent.PutAssignment(root, assignment))
}

// TestSweepHandlesPreStampLedgerRecords drives the pre-stamp ledger shape through the
// sweep's destructive path. A record with no created_at key is the standing population.
// Its age is inferred from the absent stamp rather than read.
// Both verdicts for an aged record therefore have to hold for it: a record whose tree is
// gone and unregistered is compacted; one whose tree is present is reported and left alone.
func TestSweepHandlesPreStampLedgerRecords(t *testing.T) {
	root := newSweepRepo(t)
	present := mustCreate(t, root, Home(), "prestamp-present", "unstamped, tree present")
	gone := mustCreate(t, root, Home(), "prestamp-gone", "unstamped, tree gone")
	makeUnlandedAssignment(t, present)
	makeUnlandedAssignment(t, gone)
	unstamp(t, root, present.Assignment)
	unstamp(t, root, gone.Assignment)
	address, err := intent.Address(root)
	mustNoError(t, err)
	body, err := os.ReadFile(address)
	mustNoError(t, err)
	requireTest(t, !strings.Contains(string(body), "created_at"),
		"the fixture stamped a record, so it is not the pre-stamp shape:\n%s", body)
	gitRun(t, root, "worktree", "remove", "-f", "-f", gone.Path)

	result := mustSweep(t, root)
	requireTest(t, result.Reconciled == 1,
		"Reconciled=%d, want 1: the unstamped tree-gone residue was not compacted", result.Reconciled)
	if _, err := assignmentByID(root, gone.Assignment.ID); err == nil {
		t.Fatal("the unstamped tree-gone residue record survived the sweep")
	}
	requireTest(t, len(result.Orphans) == 1 && result.Orphans[0].ID == present.Assignment.ID,
		"Orphans=%+v, want only the unstamped record whose tree is still present", result.Orphans)
}

// unregisterWorktree drops git's registration for path while leaving the tree and its
// uncommitted work on disk. No `git worktree` subcommand reaches that state — remove
// takes the tree with it — so the admin directory is deleted directly.
func unregisterWorktree(t *testing.T, root, path string) {
	t.Helper()
	mustNoError(t, os.RemoveAll(filepath.Join(root, ".git", "worktrees", filepath.Base(path))))
}

// denyStat makes every stat below dir fail for a reason other than absence.
// It restores the mode before the enclosing TempDir cleanup tries to walk it.
func denyStat(t *testing.T, dir string) {
	t.Helper()
	mustNoError(t, os.Chmod(dir, 0o000))
	t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })
}

// TestSweepRetainsRecordWhenStatIsUnknown holds the sweep to the difference between a
// tree that is gone and one it cannot see. An unreadable pool answers every stat with a
// permission error.
// If the sweep reads that as absence, it compacts the ledger row out from under a
// worktree that is still on disk holding uncommitted work.
// That is the one loss the sweep's tree-gone verdicts exist to avoid.
func TestSweepRetainsRecordWhenStatIsUnknown(t *testing.T) {
	root := newSweepRepo(t)
	unknown := mustCreate(t, root, Home(), "orphan-stat-unknown", "aged, tree present but unreadable")
	mustWrite(t, filepath.Join(unknown.Path, "work.txt"), []byte("uncommitted\n"), 0o644)
	backdate(t, root, unknown.Assignment, 8*24*time.Hour)
	unregisterWorktree(t, root, unknown.Path)
	denyStat(t, Pool(root))

	_, statErr := os.Stat(unknown.Path)
	requireTest(t, statErr != nil && !os.IsNotExist(statErr),
		"the fixture does not induce an unknown stat (this user reads a 0o000 directory): %v", statErr)

	result := mustSweep(t, root)
	requireTest(t, result.Reconciled == 0, "Reconciled=%d, want 0: an unstattable tree is unknown, not gone", result.Reconciled)
	if _, err := assignmentByID(root, unknown.Assignment.ID); err != nil {
		t.Fatalf("the sweep compacted a record whose worktree it could not stat: %v", err)
	}
	requireTest(t, len(result.Orphans) == 0,
		"Orphans=%+v, want none: the emitted line is a retirement command for a tree nothing here can reach", result.Orphans)
}

// TestSweepPurgesRecoveredRecordsWhoseTreeIsGone pins the boundary of the reconcile's
// destructive branch. No record is ever both orphaned and holding preserved work.
// Orphanhood requires an active state, and ValidateAssignment refuses an active record with
// recovery metadata on every ledger read.
// The closest reachable record is therefore a recovered one, and the standing cleaner drops it.
// The state it names is one only the removed lifecycle produced.
// The ref it points at is one the same run sweeps.
func TestSweepPurgesRecoveredRecordsWhoseTreeIsGone(t *testing.T) {
	root := newSweepRepo(t)
	preserved := mustCreate(t, root, Home(), "orphan-preserved", "aged, holds preserved work")
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
	requireTest(t, result.Reconciled == 1, "Reconciled=%d, want 1", result.Reconciled)
	if _, err := assignmentByID(root, a.ID); err == nil {
		t.Fatal("a recovered record whose tree is gone survived the reconcile")
	}
}

// TestSweepIsIdempotent pins what a reader actually sees: the summary reprints at every
// session start.
// Once the ledger has settled, repeated sweeps must agree exactly and delete nothing further.
// The first sweep is the settling run — it compacts the residue and so legitimately
// differs — and the two runs after it are the compared pair.
func TestSweepIsIdempotent(t *testing.T) {
	root := newSweepRepo(t)
	present := mustCreate(t, root, Home(), "idempotent-orphan", "aged, tree present")
	gone := mustCreate(t, root, Home(), "idempotent-residue", "aged, tree gone")
	makeUnlandedAssignment(t, present)
	makeUnlandedAssignment(t, gone)
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
