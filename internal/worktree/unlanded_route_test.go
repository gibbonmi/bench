package worktree

import (
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/intent"
)

// refusedUnlandedRelease builds the FT345 row: an assignment with one unlanded commit,
// whose release refused and left it cleanup-pending. It carries the refusal's stderr.
func refusedUnlandedRelease(t *testing.T, request string) refusedRelease {
	t.Helper()
	f := newOwnedAssignment(t, request)
	commitInWorktree(t, f.creation.Path, "unique.txt", "throwaway\n", "unlanded work")
	release := runVerb(t, verbRelease, f.call("--request", "landed-"+request, f.creation.Path))
	requireTest(t, release.exit == 1 && strings.Contains(release.stderr, "worktree retained (unmerged)"), "release exit=%d stderr=%q", release.exit, release.stderr)
	assignment, err := assignmentByID(f.root, f.creation.Assignment.ID)
	mustNoError(t, err)
	requireTest(t, assignment.State == intent.StateCleanupPending, "state = %q, want cleanup-pending", assignment.State)
	return refusedRelease{ownedAssignment: f, stderr: release.stderr}
}

// requireDiscardRoute runs the clean the surface named and requires it to plan and apply
// the removal, so the named route is one that succeeds.
func requireDiscardRoute(t *testing.T, f refusedRelease, selector ...string) {
	t.Helper()
	args := append([]string{"--discard-branch"}, selector...)
	plan := runVerb(t, verbClean, f.call(args...))
	requireTest(t, plan.exit == 0 && strings.Contains(plan.stdout, ",remove,"), "plan = (%d, %q, %q), want remove", plan.exit, plan.stdout, plan.stderr)
	applied := runVerb(t, verbClean, f.call(append(args, "--apply", plan.mustFingerprint(t))...))
	requireTest(t, applied.exit == 0 && strings.Contains(applied.stdout, ",removed,"), "apply = (%d, %q, %q), want removed", applied.exit, applied.stdout, applied.stderr)
	_, statErr := os.Lstat(f.creation.Path)
	requireTest(t, os.IsNotExist(statErr), "apply left %s: %v", f.creation.Path, statErr)
}

func TestUnlandedReleaseRefusalNamesTheDiscardClean(t *testing.T) {
	t.Parallel()
	f := refusedUnlandedRelease(t, "unlanded-next")
	want := "next=bench worktree clean --discard-branch " + axi.ShellQuote(f.creation.Path)
	requireTest(t, strings.Contains(f.stderr, want+"\n"), "release refusal = %q, want %q", f.stderr, want)
	requireDiscardRoute(t, f, f.creation.Path)
}

func TestUnlandedCleanupPendingListRowNamesTheDiscardClean(t *testing.T) {
	t.Parallel()
	f := refusedUnlandedRelease(t, "unlanded-list")
	list := runVerb(t, verbList, f.call())
	requireTest(t, list.exit == 0, "list exit=%d output=%q", list.exit, list.stdout)
	argv, err := axitest.RecoverHelpCommandArgv(list.stdout)
	mustNoError(t, err)
	want := []string{"bench", "worktree", "clean", "--discard-branch", f.creation.Path}
	requireTest(t, slices.Equal(argv, want), "list help argv = %q, want %q", argv, want)
	requireDiscardRoute(t, f, argv[4:]...)
}

func TestCleanTargetAcceptsCleanupPendingRow(t *testing.T) {
	t.Parallel()
	f := refusedUnlandedRelease(t, "unlanded-target")
	requireDiscardRoute(t, f, "--target", f.creation.Assignment.Label)
}
