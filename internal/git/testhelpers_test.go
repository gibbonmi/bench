package git

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
)

// runGit runs `git -C root <args>` and fails the test on a nonzero exit.
func runGit(t *testing.T, root string, args ...string) string {
	t.Helper()
	out, err := exec.Command("git", append([]string{"-C", root}, args...)...).Output()
	if err != nil {
		t.Fatalf("git %v: %v", args, err)
	}
	return string(out)
}

func worktreesWithin(t *testing.T, root string) ([]Worktree, error) {
	t.Helper()
	type result struct {
		worktrees []Worktree
		err       error
	}
	done := make(chan result, 1)
	go func() {
		worktrees, err := Worktrees(root)
		done <- result{worktrees: worktrees, err: err}
	}()
	select {
	case result := <-done:
		return result.worktrees, result.err
	case <-time.After(bounds.TestDeadline(0)):
		t.Fatal("git.Worktrees blocked on worktree admin entry")
		return nil, nil
	}
}

func requireAdminRefusal(t *testing.T, err error, path, shape string) {
	t.Helper()
	var got *WorktreeAdminError
	if !errors.As(err, &got) || got.Path != path || got.Shape != shape || got.Action != "inspect and remove it" || !strings.Contains(err.Error(), path) {
		t.Fatalf("admin refusal = %v, want path=%q shape=%q", err, path, shape)
	}
}

// newRepo initialises a repo with one commit and returns its root.
func newRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	runGit(t, root, "init")
	if err := os.WriteFile(filepath.Join(root, "a.txt"), []byte("hello\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "add", "-A")
	runGit(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-m", "init")
	return root
}

// newTwoBranchRepo initialises a repo with one commit and exactly two local branches.
// Neither branch is a resolvable default candidate: there is no origin/HEAD, and no
// branch named "main" for the candidate probe to verify.
func newTwoBranchRepo(t *testing.T) string {
	t.Helper()
	root := newRepo(t)
	runGit(t, root, "branch", "-M", "master")
	runGit(t, root, "branch", "feature")
	return root
}
