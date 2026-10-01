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
// whose release refused and left it cleanup-pending. It returns the refusal's stderr.
func refusedUnlandedRelease(t *testing.T, request string) (string, Creation, string, string) {
	t.Helper()
	f := newOwnedAssignment(t, request)
	commitInWorktree(t, f.creation.Path, "unique.txt", "throwaway\n", "unlanded work")
	var stdout, stderr strings.Builder
	code := ReleaseCommand(f.root, f.home, []string{"--request", "landed-" + request, f.creation.Path}, &stdout, &stderr)
	requireTest(t, code == 1 && strings.Contains(stderr.String(), "worktree retained (unmerged)"), "release exit=%d stderr=%q", code, stderr.String())
	assignment, err := assignmentByID(f.root, f.creation.Assignment.ID)
	mustNoError(t, err)
	requireTest(t, assignment.State == intent.StateCleanupPending, "state = %q, want cleanup-pending", assignment.State)
	return f.root, f.creation, f.home, stderr.String()
}

// requireDiscardRoute runs the clean the surface named and requires it to plan and apply
// the removal, so the named route is one that succeeds.
func requireDiscardRoute(t *testing.T, root, home, path string, selector ...string) {
	t.Helper()
	args := append([]string{"--discard-branch"}, selector...)
	plan, planErr, planCode := runCleanup(t, root, home, args...)
	requireTest(t, planCode == 0 && strings.Contains(plan, ",remove,"), "plan = (%d, %q, %q), want remove", planCode, plan, planErr)
	applied, applyErr, applyCode := runCleanup(t, root, home, append(args, "--apply", cleanupRowFingerprint(t, plan))...)
	requireTest(t, applyCode == 0 && strings.Contains(applied, ",removed,"), "apply = (%d, %q, %q), want removed", applyCode, applied, applyErr)
	_, statErr := os.Lstat(path)
	requireTest(t, os.IsNotExist(statErr), "apply left %s: %v", path, statErr)
}

func TestUnlandedReleaseRefusalNamesTheDiscardClean(t *testing.T) {
	t.Parallel()
	root, creation, home, stderr := refusedUnlandedRelease(t, "unlanded-next")
	want := "next=bench worktree clean --discard-branch " + axi.ShellQuote(creation.Path)
	requireTest(t, strings.Contains(stderr, want+"\n"), "release refusal = %q, want %q", stderr, want)
	requireDiscardRoute(t, root, home, creation.Path, creation.Path)
}

func TestUnlandedCleanupPendingListRowNamesTheDiscardClean(t *testing.T) {
	t.Parallel()
	root, creation, home, _ := refusedUnlandedRelease(t, "unlanded-list")
	list, code := ListCommand(root, home, nil)
	requireTest(t, code == 0, "list exit=%d output=%q", code, list)
	argv, err := axitest.RecoverHelpCommandArgv(list)
	mustNoError(t, err)
	want := []string{"bench", "worktree", "clean", "--discard-branch", creation.Path}
	requireTest(t, slices.Equal(argv, want), "list help argv = %q, want %q", argv, want)
	requireDiscardRoute(t, root, home, creation.Path, argv[4:]...)
}

func TestCleanTargetAcceptsCleanupPendingRow(t *testing.T) {
	t.Parallel()
	root, creation, home, _ := refusedUnlandedRelease(t, "unlanded-target")
	requireDiscardRoute(t, root, home, creation.Path, "--target", creation.Assignment.Label)
}
