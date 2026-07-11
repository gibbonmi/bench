package worktree

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
)

func TestOrphanedDelegateBranches(t *testing.T) {
	root := t.TempDir()
	gitRunBranchTest(t, root, "init")
	gitRunBranchTest(t, root, "config", "user.email", "t@example.com")
	gitRunBranchTest(t, root, "config", "user.name", "t")
	if err := os.WriteFile(filepath.Join(root, "f.txt"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRunBranchTest(t, root, "add", "-A")
	gitRunBranchTest(t, root, "commit", "-m", "base")
	gitRunBranchTest(t, root, "branch", "worktree-agent-orphan")
	gitRunBranchTest(t, root, "branch", "bench/shift-review")
	gitRunBranchTest(t, root, "worktree", "add", "-q", "-b", "worktree-agent-active", filepath.Join(root, "..", "active-wt"), "HEAD")

	want := []string{"worktree-agent-orphan"}
	got, err := OrphanedDelegateBranches(root)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("OrphanedDelegateBranches = %#v, want %#v", got, want)
	}
}

func gitRunBranchTest(t *testing.T, root string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v (%s)", args, err, out)
	}
}
