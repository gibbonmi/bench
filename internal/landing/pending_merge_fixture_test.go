package landing

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// The paths that every pending merge fixture changes: both sides change mergeConflictPath,
// and the incoming commit deletes mergeDeletedPath. mergeResolution is the worktree content
// of the conflicted path after the fixture resolves it.
const (
	mergeConflictPath = "conflict.txt"
	mergeDeletedPath  = "gone.txt"
	mergeResolution   = "resolution"
)

// pendingMerge builds a linked worktree in a conflicted merge of an incoming commit, and it
// returns the worktree and its HEAD before the merge. The incoming commit also adds each
// path of added. The fixture writes the resolution to the worktree file only, so the index
// keeps the unmerged entries of the conflicted path. The landing helper git fails on the
// exit 1 of a conflicted merge, so the merge runs through its own call, as stageHandMerge
// in internal/worktree does; no shared helper crosses the package boundary.
func pendingMerge(t *testing.T, added ...string) (root, base string) {
	t.Helper()
	primary := fixture(t)
	write(t, primary, mergeConflictPath, "base\n")
	write(t, primary, mergeDeletedPath, "base\n")
	git(t, primary, "add", ".")
	git(t, primary, "commit", "-qm", "merge base")
	git(t, primary, "checkout", "-qb", "incoming")
	write(t, primary, mergeConflictPath, "incoming\n")
	for _, path := range added {
		write(t, primary, path, "incoming\n")
	}
	git(t, primary, "rm", "-q", mergeDeletedPath)
	git(t, primary, "add", "-A")
	git(t, primary, "commit", "-qm", "incoming")
	root = filepath.Join(t.TempDir(), "linked")
	git(t, primary, "worktree", "add", "-q", "-b", "topic", root, "main")
	write(t, root, mergeConflictPath, "ours\n")
	git(t, root, "commit", "-qam", "ours")
	base = git(t, root, "rev-parse", "HEAD")
	merge := exec.Command("git", "-C", root, "merge", "--no-commit", "incoming")
	if out, err := merge.CombinedOutput(); err == nil || !strings.Contains(string(out), "CONFLICT") {
		t.Fatalf("merge = %v, %s; want the conflict", err, out)
	}
	write(t, root, mergeConflictPath, mergeResolution)
	return root, base
}
