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

// TestDiscardedRefNamesTheDateAndTheBranchPath is RI42 and RI95: the date is the UTC
// date of the instant, also when the instant's own zone is on a different day.
func TestDiscardedRefNamesTheDateAndTheBranchPath(t *testing.T) {
	t.Parallel()
	tail := "/bench/assign/" + strings.Repeat("a", 32) + "/" + strings.Repeat("b", 32)
	for _, row := range []struct {
		name    string
		instant time.Time
		want    string
	}{
		{"RI42 a UTC instant", time.Date(2026, 9, 27, 15, 4, 5, 0, time.UTC), "refs/bench/discarded/20260927" + tail},
		{"RI95 a local instant a day behind UTC", time.Date(2026, 9, 27, 22, 0, 0, 0, time.FixedZone("UTC-4", -4*60*60)), "refs/bench/discarded/20260928" + tail},
	} {
		got := intent.DiscardedRef(row.instant, discardedFixtureBranch)
		requireTest(t, got == row.want, "%s: discarded ref = %q, want %q", row.name, got, row.want)
	}
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

// TestSweepKeepsADiscardedRefWithAnUnparseableDate is RI46 and RI94. The segment
// 20200101x starts with a date far past the window, so a parser that reads only the
// first eight bytes deletes it.
func TestSweepKeepsADiscardedRefWithAnUnparseableDate(t *testing.T) {
	t.Parallel()
	for _, segment := range []string{"latest", "20200101x"} {
		t.Run(segment, func(t *testing.T) {
			root := newWorktreeRepo(t)
			ref := intent.DiscardedRefNamespace + segment + "/branch"
			gitRun(t, root, "update-ref", ref, "HEAD")
			swept, _, err := reconcileLifecycleDebris(defaultJoins(), root, nil, time.Date(2026, 9, 27, 0, 0, 0, 0, time.UTC))
			mustNoError(t, err)
			requireTest(t, swept == 0, "swept count = %d, want 0", swept)
			requireTest(t, refsUnder(t, root, ref) != "", "a discarded ref with the date segment %s was deleted", segment)
		})
	}
}

// TestSweepDeletesADiscardedSymrefAndNotItsTarget is RI93. The sweep deletes the
// expired symref itself and does not follow it to the branch it names.
func TestSweepDeletesADiscardedSymrefAndNotItsTarget(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	head := gitOutput(t, root, "rev-parse", "HEAD")
	gitRun(t, root, "update-ref", "refs/heads/keep", head)
	symref := intent.DiscardedRefNamespace + "20200101/sym"
	gitRun(t, root, "symbolic-ref", symref, "refs/heads/keep")
	swept, _, err := reconcileLifecycleDebris(defaultJoins(), root, nil, time.Date(2026, 9, 27, 0, 0, 0, 0, time.UTC))
	mustNoError(t, err)
	requireTest(t, refsUnder(t, root, "refs/heads/keep") != "" && gitOutput(t, root, "rev-parse", "refs/heads/keep") == head,
		"the sweep followed the discarded symref and deleted refs/heads/keep")
	requireTest(t, swept == 1 && refsUnder(t, root, symref) == "", "swept count = %d, symref after the sweep: %q", swept, refsUnder(t, root, symref))
}

// TestSweepRefusesADiscardedRefMovedAfterListing is RI47.
func TestSweepRefusesADiscardedRefMovedAfterListing(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	ref := plantDiscardedRef(t, root, discardedFixtureDate)
	j, moved := movedRefJoins(t, root, ref)
	swept, _, err := reconcileLifecycleDebris(j, root, nil, time.Date(2026, 9, 27, 0, 0, 0, 0, time.UTC))
	requireTest(t, err != nil && strings.Contains(err.Error(), "delete lifecycle ref "+ref), "moved discarded ref error = %v", err)
	requireTest(t, swept == 0, "failed delete entered the swept count: %d", swept)
	requireTest(t, gitOutput(t, root, "rev-parse", ref) == moved, "the sweep deleted the moved discarded ref")
}
