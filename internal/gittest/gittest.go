// Package gittest owns the one git test-repository scaffold the package tests share.
// It creates repositories and nothing else: package-specific fixture bytes stay with the
// test that grades them. Only tests import this package; no production package may.
package gittest

import (
	"os/exec"
	"strings"
	"testing"
)

// Repo initializes an empty repository in a fresh temporary directory and returns its root.
// A test that commits wants RepoOnBranch instead, which also pins the branch and an identity.
func Repo(t testing.TB) string {
	t.Helper()
	return initialize(t)
}

// RepoOnBranch initializes an empty repository on branch with a commit identity configured,
// so a committing test neither inherits nor depends on the developer's own git config.
func RepoOnBranch(t testing.TB, branch string) string {
	t.Helper()
	root := initialize(t, "-b", branch)
	run(t, root, "config", "user.email", "bench@example.invalid")
	run(t, root, "config", "user.name", "bench test")
	return root
}

func initialize(t testing.TB, options ...string) string {
	t.Helper()
	root := t.TempDir()
	run(t, root, append([]string{"init", "-q"}, options...)...)
	return root
}

func run(t testing.TB, root string, args ...string) {
	t.Helper()
	if out, err := exec.Command("git", append([]string{"-C", root}, args...)...).CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v: %s", strings.Join(args, " "), err, out)
	}
}
