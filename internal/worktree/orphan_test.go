package worktree

import (
	"fmt"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/intent"
	"io"
	"path/filepath"
	"testing"
	"time"
)

func TestCreateStampsAssignment(t *testing.T) {
	root := newWorktreeRepo(t)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	before := time.Now().UTC().Truncate(time.Second)
	creation := mustCreate(t, root, "stamped", "creation stamp")
	after := time.Now().UTC()
	stored, err := assignmentByID(root, creation.Assignment.ID)
	mustNoError(t, err)
	requireTest(t, creation.Assignment.CreatedAt != nil && stored.CreatedAt != nil && *stored.CreatedAt == *creation.Assignment.CreatedAt,
		"created assignment stamp = %v, ledger stamp = %v", creation.Assignment.CreatedAt, stored.CreatedAt)
	stamp, err := time.Parse(time.RFC3339, *stored.CreatedAt)
	mustNoError(t, err)
	requireTest(t, !stamp.Before(before) && !stamp.After(after), "stamp %s outside [%s, %s]", stamp, before, after)
}

// orphanNow is the fixed instant every predicate case is judged against, so no case
// depends on the wall clock.
var orphanNow = time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)

func stampedAt(when time.Time) *string {
	stamp := when.UTC().Format(time.RFC3339)
	return &stamp
}

func TestOrphanedOnAgedRecord(t *testing.T) {
	aged := intent.Assignment{State: intent.StateActive, CreatedAt: stampedAt(orphanNow.Add(-8 * 24 * time.Hour))}
	requireTest(t, orphaned(aged, orphanNow), "orphaned(active, stamped 8 days ago) = false, want true")
}

func TestOrphanedTreatsAbsentStampAsAged(t *testing.T) {
	unstamped := intent.Assignment{State: intent.StateActive}
	requireTest(t, orphaned(unstamped, orphanNow), "orphaned(active, unstamped) = false, want true")
}

func TestOrphanedRequiresAge(t *testing.T) {
	young := intent.Assignment{State: intent.StateActive, CreatedAt: stampedAt(orphanNow.Add(-6 * 24 * time.Hour))}
	requireTest(t, !orphaned(young, orphanNow), "orphaned(active, stamped 6 days ago) = true, want false")
}

// The window's edge is a decision rather than an accident of the comparison: a record
// aged exactly AssignmentStale is still inside the window, and only a strictly older one
// is abandoned.
func TestOrphanedExcludesTheExactWindowEdge(t *testing.T) {
	edge := intent.Assignment{State: intent.StateActive, CreatedAt: stampedAt(orphanNow.Add(-bounds.AssignmentStale))}
	requireTest(t, !orphaned(edge, orphanNow), "orphaned(active, stamped exactly AssignmentStale ago) = true, want false")
}

func TestOrphanedRejectsFutureStamp(t *testing.T) {
	skewed := intent.Assignment{State: intent.StateActive, CreatedAt: stampedAt(orphanNow.Add(30 * 24 * time.Hour))}
	requireTest(t, !orphaned(skewed, orphanNow), "orphaned(active, stamped 30 days ahead) = true, want false")
}

func TestOrphanedOnlyActiveState(t *testing.T) {
	for _, state := range []intent.AssignmentState{intent.StateCleanupPending, intent.StateRecovered, intent.StateComplete} {
		t.Run(string(state), func(t *testing.T) {
			aged := intent.Assignment{State: state, CreatedAt: stampedAt(orphanNow.Add(-8 * 24 * time.Hour))}
			requireTest(t, !orphaned(aged, orphanNow), "orphaned(%s, stamped 8 days ago) = true, want false", state)
		})
	}
}

// An unparseable stamp is unknown age, and unknown must not read as abandoned:
// ValidateAssignment rejects such a record on every ledger read, so reaching here at
// all means the caller bypassed the ledger.
func TestOrphanedRejectsUnparseableStamp(t *testing.T) {
	for _, stamp := range []string{"", "yesterday", "2026-07-\x0027T12:00:00Z"} {
		t.Run(fmt.Sprintf("%q", stamp), func(t *testing.T) {
			bad := intent.Assignment{State: intent.StateActive, CreatedAt: &stamp}
			requireTest(t, !orphaned(bad, orphanNow), "orphaned(active, created_at=%q) = true, want false", stamp)
		})
	}
}

// backdate ages a live assignment by rewriting its ledger stamp, which is how a test
// reaches an aged record through PlanAutomatic — the planner reads the ledger and the
// clock itself, so neither is an argument a caller can set.
func backdate(t *testing.T, root string, assignment intent.Assignment, age time.Duration) {
	t.Helper()
	assignment.CreatedAt = stampedAt(time.Now().Add(-age))
	mustNoError(t, intent.PutAssignment(root, assignment))
}

func TestPlanAutomaticLabelsOrphaned(t *testing.T) {
	root, creation := newOwnedAssignment(t, "orphan-label")
	backdate(t, root, creation.Assignment, 8*24*time.Hour)
	plan, err := PlanAutomatic(root, creation.Path)
	requireTest(t, err == nil && plan.Action == ActionRetain && plan.ReasonCode == ReasonOrphaned,
		"PlanAutomatic over an aged clean assignment = action %q reason %q, %v", plan.Action, plan.ReasonCode, err)
}

func TestPlanAutomaticKeepsEarlierRetainReason(t *testing.T) {
	root := newWorktreeRepo(t)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	mustWrite(t, filepath.Join(root, ".gitignore"), []byte("ignored.txt\n"), 0o644)
	gitRun(t, root, "add", ".gitignore")
	gitRun(t, root, "commit", "-qm", "ignore")
	creation := mustCreate(t, root, "orphan-residue", "ignored residue")
	mustWrite(t, filepath.Join(creation.Path, "ignored.txt"), []byte("residue\n"), 0o644)
	backdate(t, root, creation.Assignment, 8*24*time.Hour)
	plan, err := PlanAutomatic(root, creation.Path)
	requireTest(t, err == nil && plan.Action == ActionRetain && plan.ReasonCode == ReasonIgnored,
		"PlanAutomatic over an aged assignment holding ignored residue = action %q reason %q, %v", plan.Action, plan.ReasonCode, err)
}

// A record carrying no creation stamp is what every worktree cut before the field
// existed looks like, and its git lock was written from the same fields. Release and
// explicit cleanup both re-derive that lock string, so this pins that neither path
// reads the stamp.
func TestReleaseAndPlanExplicitAcceptUnstampedAssignment(t *testing.T) {
	root, creation := newOwnedAssignment(t, "unstamped-lock")
	unstamped := creation.Assignment
	unstamped.CreatedAt = nil
	mustNoError(t, intent.PutAssignment(root, unstamped))
	plan, err := PlanExplicit(root, creation.Path)
	requireTest(t, err == nil && plan.Action == ActionRemove && plan.ReasonCode == "",
		"PlanExplicit over an unstamped assignment = %#v, %v", plan, err)
	code := ReleaseCommand(root, []string{"--request", "landed-unstamped-lock", creation.Path}, io.Discard, io.Discard)
	requireTest(t, code == 0, "ReleaseCommand over an unstamped assignment exit=%d", code)
}
