package landing

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

// The paths and the branches of every pending merge fixture: both sides change
// mergeConflictPath, and the incoming commit on the branch mergeIncoming deletes
// mergeDeletedPath. mergeResolution is the worktree content of the conflicted path after the
// fixture resolves it. The linked worktree is on the branch mergeTopic.
const (
	mergeConflictPath = "conflict.txt"
	mergeDeletedPath  = "gone.txt"
	mergeResolution   = "resolution"
	mergeIncoming     = "incoming"
	mergeTopic        = "topic"
)

// mergeNameSet is the complete name set of a pendingMerge fixture that adds no path.
var mergeNameSet = []string{mergeConflictPath, mergeDeletedPath}

// mergeStateFiles are the files that Git writes for a conflicted merge in the checkout
// administration directory.
var mergeStateFiles = []string{"MERGE_HEAD", "MERGE_MSG", "MERGE_MODE", "AUTO_MERGE"}

// mergeCheckout builds a linked worktree on the topic branch, whose own commit changes the
// conflicted path, and the incoming branch, which also adds each path of added. It returns
// the worktree and its HEAD, and no merge is pending.
func mergeCheckout(t *testing.T, added ...string) (root, base string) {
	t.Helper()
	primary := fixture(t)
	write(t, primary, mergeConflictPath, "base\n")
	write(t, primary, mergeDeletedPath, "base\n")
	git(t, primary, "add", ".")
	git(t, primary, "commit", "-qm", "merge base")
	git(t, primary, "checkout", "-qb", mergeIncoming)
	write(t, primary, mergeConflictPath, "incoming\n")
	for _, path := range added {
		write(t, primary, path, "incoming\n")
	}
	git(t, primary, "rm", "-q", mergeDeletedPath)
	git(t, primary, "add", "-A")
	git(t, primary, "commit", "-qm", "incoming")
	root = filepath.Join(t.TempDir(), "linked")
	git(t, primary, "worktree", "add", "-q", "-b", mergeTopic, root, "main")
	write(t, root, mergeConflictPath, "ours\n")
	git(t, root, "commit", "-qam", "ours")
	return root, git(t, root, "rev-parse", "HEAD")
}

// pendingMerge starts a conflicted merge of the incoming branch in the worktree of
// mergeCheckout, and it returns the worktree and its HEAD before the merge. The fixture
// writes the resolution to the worktree file only, so the index keeps the unmerged entries
// of the conflicted path. The landing helper git fails on the exit 1 of a conflicted merge,
// so the merge runs through its own call, as stageHandMerge in internal/worktree does; no
// shared helper crosses the package boundary.
func pendingMerge(t *testing.T, added ...string) (root, base string) {
	t.Helper()
	root, base = mergeCheckout(t, added...)
	merge := exec.Command("git", "-C", root, "merge", "--no-commit", mergeIncoming)
	if out, err := merge.CombinedOutput(); err == nil || !strings.Contains(string(out), "CONFLICT") {
		t.Fatalf("merge = %v, %s; want the conflict", err, out)
	}
	write(t, root, mergeConflictPath, mergeResolution)
	return root, base
}

// octopusMerge starts a pending merge of two side commits in the worktree of mergeCheckout.
// Each side commit forks from the worktree's HEAD and adds the path of its own name, so the
// merge has no conflict. It returns the worktree, its HEAD before the merge, the merge heads
// in merge order, and the added paths.
func octopusMerge(t *testing.T) (root, base string, heads, added []string) {
	t.Helper()
	root, base = mergeCheckout(t)
	added = []string{"side-a", "side-b"}
	for _, side := range added {
		git(t, root, "checkout", "-q", "-b", side, base)
		write(t, root, side, "side\n")
		git(t, root, "add", side)
		git(t, root, "commit", "-qm", side)
		heads = append(heads, git(t, root, "rev-parse", "HEAD"))
	}
	git(t, root, "checkout", "-q", mergeTopic)
	git(t, root, append([]string{"merge", "--no-commit", "--no-ff"}, added...)...)
	return root, base, heads, added
}

// topicLanding is the request that lands paths on the topic branch of a pending merge
// fixture from base.
func topicLanding(root, base string, paths ...string) Request {
	return Request{Root: root, Destination: "refs/heads/" + mergeTopic, Expected: base, Message: "merge", Paths: paths}
}

// mergeStatePath is the path of the merge-state file name in the checkout administration
// directory of root.
func mergeStatePath(t *testing.T, root, name string) string {
	t.Helper()
	path, err := benchgit.AdminPath(root, name)
	if err != nil {
		t.Fatal(err)
	}
	return path
}

// mergeStatePresent reports whether the merge-state file name exists for root.
func mergeStatePresent(t *testing.T, root, name string) bool {
	t.Helper()
	_, err := os.Lstat(mergeStatePath(t, root, name))
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	return err == nil
}

// plantMergeHeadDirectory replaces the MERGE_HEAD of root with a directory, which the
// pending-merge reader refuses.
func plantMergeHeadDirectory(t *testing.T, root string) {
	t.Helper()
	mergeHead := mergeStatePath(t, root, "MERGE_HEAD")
	if err := os.Remove(mergeHead); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(mergeHead, 0o755); err != nil {
		t.Fatal(err)
	}
}
