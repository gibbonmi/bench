//go:build system

package systemtest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

func TestStatusRouteExecutesUnclaimedCleanup(t *testing.T) {
	repo, err := os.MkdirTemp(owner.root, "status-route [unclaimed]-")
	if err != nil {
		t.Fatal(err)
	}
	if result := owner.runAt(repo, nil, "git", "init", "-q", "-b", "main"); result.code != 0 {
		t.Fatalf("git init = (%d, %q)", result.code, result.stderr)
	}
	for _, args := range [][]string{{"config", "user.email", "route@example.test"}, {"config", "user.name", "Route Test"}} {
		if result := owner.runAt(repo, nil, "git", args...); result.code != 0 {
			t.Fatalf("git %s = (%d, %q)", strings.Join(args, " "), result.code, result.stderr)
		}
	}
	if result := owner.runAt(repo, nil, "git", "commit", "--allow-empty", "-qm", "base"); result.code != 0 {
		t.Fatalf("git commit base = (%d, %q)", result.code, result.stderr)
	}
	branch := intent.AssignmentBranchRef(strings.Repeat("a", 32), strings.Repeat("b", 32))
	shortBranch := strings.TrimPrefix(branch, "refs/heads/")
	if result := owner.runAt(repo, nil, "git", "checkout", "-qb", shortBranch); result.code != 0 {
		t.Fatalf("git checkout = (%d, %q)", result.code, result.stderr)
	}
	if err := os.WriteFile(filepath.Join(repo, "orphan.txt"), []byte("orphan\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"add", "orphan.txt"}, {"commit", "-qm", "orphan assignment"}, {"checkout", "-q", "main"}} {
		if result := owner.runAt(repo, nil, "git", args...); result.code != 0 {
			t.Fatalf("git %s = (%d, %q)", strings.Join(args, " "), result.code, result.stderr)
		}
	}

	route := owner.runSelected(repo, "status", "--route")
	if route.code != 0 {
		t.Fatalf("bench status --route = (%d, %q, %q)", route.code, route.stdout, route.stderr)
	}
	const cleanupPrefix = "bench worktree clean --discard-branch --unclaimed"
	start := strings.Index(route.stdout, cleanupPrefix)
	if start < 0 {
		t.Fatalf("bench status --route returned %q, want %q", route.stdout, cleanupPrefix)
	}
	command := route.stdout[start:]
	if end := strings.IndexByte(command, '\n'); end >= 0 {
		command = command[:end]
	}
	shim := t.TempDir()
	if err := os.Symlink(owner.selected.path, filepath.Join(shim, "bench")); err != nil {
		t.Fatal(err)
	}
	run := owner.runAt(repo, []string{"PATH=" + shim + string(os.PathListSeparator) + os.Getenv("PATH")}, "bash", "-c", command)
	if run.code != 0 {
		t.Fatalf("routed command %q = (%d, %q, %q)", command, run.code, run.stdout, run.stderr)
	}
	if result := owner.runAt(repo, nil, "git", "show-ref", "--verify", "--quiet", branch); result.code == 0 {
		t.Fatalf("routed command %q left unclaimed assignment branch %q", command, branch)
	}
}
