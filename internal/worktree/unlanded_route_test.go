package worktree

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// unlandedWork is one shape of work that main does not hold. The release refuses each
// shape as unmerged, so the refusal and the list help row name the discard clean.
type unlandedWork struct {
	name  string
	leave func(t *testing.T, path string)
	plan  CleanupAction
}

// unlandedShapes are the FT345 shape, an own commit that main lacks, and the FT408 shape,
// no own commit and uncommitted output. plan is the action the discard clean plans: the
// uncommitted output goes to a recovery ref before the removal.
var unlandedShapes = []unlandedWork{
	{"unlanded commit", func(t *testing.T, path string) {
		commitInWorktree(t, path, "unique.txt", "throwaway\n", "unlanded work")
	}, ActionRemove},
	{"no own commit with output", func(t *testing.T, path string) {
		mustWrite(t, filepath.Join(path, "research.txt"), []byte("output\n"), 0o644)
	}, ActionRecoverRemove},
}

// refusedUnlandedRelease builds an assignment that holds work, whose release refused and
// left it cleanup-pending. It carries the refusal's stderr.
func refusedUnlandedRelease(t *testing.T, request string, work unlandedWork) refusedRelease {
	t.Helper()
	f := newOwnedAssignment(t, request)
	work.leave(t, f.creation.Path)
	release := runVerb(t, verbRelease, f.call("--request", "landed-"+request, f.creation.Path))
	requireTest(t, release.exit == 1 && strings.Contains(release.stderr, "worktree retained (unmerged)"), "release exit=%d stderr=%q", release.exit, release.stderr)
	assignment, err := assignmentByID(f.root, f.creation.Assignment.ID)
	mustNoError(t, err)
	requireTest(t, assignment.State == intent.StateCleanupPending, "state = %q, want cleanup-pending", assignment.State)
	return refusedRelease{ownedAssignment: f, stderr: release.stderr}
}

// requireDiscardRoute runs the clean the surface named and requires it to plan the action
// want and apply the removal, so the named route is one that succeeds.
func requireDiscardRoute(t *testing.T, f refusedRelease, want CleanupAction, selector ...string) {
	t.Helper()
	args := append([]string{"--discard-branch"}, selector...)
	plan := runVerb(t, verbClean, f.call(args...))
	requireTest(t, plan.exit == 0 && strings.Contains(plan.stdout, ","+string(want)+","), "plan = (%d, %q, %q), want %s", plan.exit, plan.stdout, plan.stderr, want)
	applied := runVerb(t, verbClean, f.call(append(args, "--apply", plan.mustFingerprint(t))...))
	requireTest(t, applied.exit == 0 && strings.Contains(applied.stdout, ",removed,"), "apply = (%d, %q, %q), want removed", applied.exit, applied.stdout, applied.stderr)
	_, statErr := os.Lstat(f.creation.Path)
	requireTest(t, os.IsNotExist(statErr), "apply left %s: %v", f.creation.Path, statErr)
}

func TestUnlandedReleaseRefusalNamesTheDiscardClean(t *testing.T) {
	t.Parallel()
	for _, work := range unlandedShapes {
		t.Run(work.name, func(t *testing.T) {
			t.Parallel()
			f := refusedUnlandedRelease(t, "unlanded-next", work)
			want := "next=bench worktree clean --discard-branch " + axi.ShellQuote(f.creation.Path)
			requireTest(t, strings.Contains(f.stderr, want+"\n"), "release refusal = %q, want %q", f.stderr, want)
			requireDiscardRoute(t, f, work.plan, f.creation.Path)
		})
	}
}

func TestUnlandedCleanupPendingListRowNamesTheDiscardClean(t *testing.T) {
	t.Parallel()
	for _, work := range unlandedShapes {
		t.Run(work.name, func(t *testing.T) {
			t.Parallel()
			f := refusedUnlandedRelease(t, "unlanded-list", work)
			list := runVerb(t, verbList, f.call())
			requireTest(t, list.exit == 0, "list exit=%d output=%q", list.exit, list.stdout)
			argv, err := axitest.RecoverHelpCommandArgv(list.stdout)
			mustNoError(t, err)
			want := []string{"bench", "worktree", "clean", "--discard-branch", f.creation.Path}
			requireTest(t, slices.Equal(argv, want), "list help argv = %q, want %q", argv, want)
			requireDiscardRoute(t, f, work.plan, argv[4:]...)
		})
	}
}

// newLandedDirtyAssignment builds the landed neighbor of the FT408 shape: an assignment
// whose own commit main holds, with an untracked file in its tree.
func newLandedDirtyAssignment(t *testing.T, request string) ownedAssignment {
	t.Helper()
	f := newOwnedAssignment(t, request)
	landAssignment(t, f.root, f.creation, "own-commit.txt")
	mustWrite(t, filepath.Join(f.creation.Path, "extra.txt"), []byte("output\n"), 0o644)
	return f
}

// TestLandedDirtyReleaseKeepsTheReleaseRoute pins the neighbor of the FT408 shape. An
// assignment whose own commit main holds has landed. Its release refuses the uncommitted
// output as dirty and names the release route, which the operator runs after the output is
// resolved.
func TestLandedDirtyReleaseKeepsTheReleaseRoute(t *testing.T) {
	t.Parallel()
	f := newLandedDirtyAssignment(t, "landed-dirty")
	release := runVerb(t, verbRelease, f.call("--request", "landed-landed-dirty", f.creation.Path))
	want := "next=bench worktree release --request <request> " + sanitize.ShellQuote(f.creation.Path) + "\n"
	requireTest(t, release.exit == 1 && strings.Contains(release.stderr, "worktree retained (dirty)") && strings.Contains(release.stderr, want),
		"release exit=%d stderr=%q, want a dirty retain and %q", release.exit, release.stderr, want)
}

func TestCleanTargetAcceptsCleanupPendingRow(t *testing.T) {
	t.Parallel()
	work := unlandedShapes[0]
	f := refusedUnlandedRelease(t, "unlanded-target", work)
	requireDiscardRoute(t, f, work.plan, "--target", f.creation.Assignment.Label)
}
