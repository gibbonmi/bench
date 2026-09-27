package worktree

import (
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/intent"
)

// discardedFixtureBranch is the assignment branch every discarded-ref fixture names.
var discardedFixtureBranch = intent.AssignmentBranchRef(strings.Repeat("a", 32), strings.Repeat("b", 32))

// discardedFixtureDate is the day D that each boundary row plants its ref under.
var discardedFixtureDate = time.Date(2026, 8, 28, 0, 0, 0, 0, time.UTC)

// plantDiscardedRef writes a discarded ref dated day at HEAD and returns its name.
func plantDiscardedRef(t *testing.T, root string, day time.Time) string {
	t.Helper()
	ref := intent.DiscardedRef(day, discardedFixtureBranch)
	gitRun(t, root, "update-ref", ref, "HEAD")
	return ref
}

// TestDiscardedRefNamesTheDateAndTheBranchPath is RI42.
func TestDiscardedRefNamesTheDateAndTheBranchPath(t *testing.T) {
	t.Parallel()
	instant := time.Date(2026, 9, 27, 15, 4, 5, 0, time.UTC)
	want := "refs/bench/discarded/20260927/bench/assign/" + strings.Repeat("a", 32) + "/" + strings.Repeat("b", 32)
	got := intent.DiscardedRef(instant, discardedFixtureBranch)
	requireTest(t, got == want, "discarded ref = %q, want %q", got, want)
}

// TestSweepDeletesADiscardedRefAtThirtyDays is RI43 and RI44: the ref dated D survives
// the last second of D plus 29 days and goes at D plus 30 days, 00:00:00Z.
func TestSweepDeletesADiscardedRefAtThirtyDays(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name    string
		instant time.Time
		swept   int
	}{
		{"RI43 at D plus 30 days", time.Date(2026, 9, 27, 0, 0, 0, 0, time.UTC), 1},
		{"RI44 at D plus 29 days 23:59:59", time.Date(2026, 9, 26, 23, 59, 59, 0, time.UTC), 0},
	} {
		t.Run(row.name, func(t *testing.T) {
			root := newWorktreeRepo(t)
			ref := plantDiscardedRef(t, root, discardedFixtureDate)
			swept, _, err := reconcileLifecycleDebris(defaultJoins(), root, nil, row.instant)
			mustNoError(t, err)
			requireTest(t, swept == row.swept, "swept count = %d, want %d", swept, row.swept)
			requireTest(t, (refsUnder(t, root, ref) == "") == (row.swept == 1), "ref %s after the sweep: %q", ref, refsUnder(t, root, ref))
		})
	}
}

// TestSweepKeepsATodayDiscardedRefWhileItEmptiesRecovery is RI45.
func TestSweepKeepsATodayDiscardedRefWhileItEmptiesRecovery(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	now := time.Date(2026, 9, 27, 12, 0, 0, 0, time.UTC)
	ref := plantDiscardedRef(t, root, now)
	recovery := intent.RecoveryRefPrefix(strings.Repeat("a", 32), strings.Repeat("b", 32)) + "1"
	gitRun(t, root, "update-ref", recovery, "HEAD")
	swept, _, err := reconcileLifecycleDebris(defaultJoins(), root, nil, now)
	mustNoError(t, err)
	requireTest(t, swept == 1, "swept count = %d, want the recovery ref only", swept)
	requireTest(t, refsUnder(t, root, recovery) == "", "recovery ref survived the emptying pass")
	requireTest(t, refsUnder(t, root, ref) != "", "today's discarded ref was emptied with the lifecycle namespaces")
}

// TestSweepKeepsADiscardedRefWithAnUnparseableDate is RI46.
func TestSweepKeepsADiscardedRefWithAnUnparseableDate(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	ref := intent.DiscardedRefNamespace + "latest/" + strings.TrimPrefix(discardedFixtureBranch, "refs/heads/")
	gitRun(t, root, "update-ref", ref, "HEAD")
	swept, _, err := reconcileLifecycleDebris(defaultJoins(), root, nil, time.Date(2026, 9, 27, 0, 0, 0, 0, time.UTC))
	mustNoError(t, err)
	requireTest(t, swept == 0, "swept count = %d, want 0", swept)
	requireTest(t, refsUnder(t, root, ref) != "", "a discarded ref with the date segment latest was deleted")
}

// TestSweepRefusesADiscardedRefMovedAfterListing is RI47.
func TestSweepRefusesADiscardedRefMovedAfterListing(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	ref := plantDiscardedRef(t, root, discardedFixtureDate)
	moved, err := commitTree(root, gitOutput(t, root, "rev-parse", "HEAD^{tree}"), nil, "concurrent ref\n")
	mustNoError(t, err)
	j := defaultJoins()
	j.cleanupBoundary = func(step LifecycleStep) error {
		if step == StepLifecycleSweep {
			gitRun(t, root, "update-ref", ref, moved)
		}
		return nil
	}
	swept, _, err := reconcileLifecycleDebris(j, root, nil, time.Date(2026, 9, 27, 0, 0, 0, 0, time.UTC))
	requireTest(t, err != nil && strings.Contains(err.Error(), "delete lifecycle ref "+ref), "moved discarded ref error = %v", err)
	requireTest(t, swept == 0, "failed delete entered the swept count: %d", swept)
	requireTest(t, gitOutput(t, root, "rev-parse", ref) == moved, "the sweep deleted the moved discarded ref")
}
