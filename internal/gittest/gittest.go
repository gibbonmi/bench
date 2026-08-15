// Package gittest owns the shared git test scaffolds package tests use.
// Repository setup and file-backed command fixtures live here; package-specific fixture
// bytes stay with the test that grades them. Only tests import this package.
package gittest

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Repo initializes an empty repository in a fresh temporary directory and returns its root.
// A test that commits wants RepoOnBranch instead, which also pins the branch and an identity.
func Repo(t testing.TB) string {
	t.Helper()
	return initialize(t)
}

// StubGit installs a pure file-backed git stub for worktree-resolution tests.
func StubGit(t testing.TB, root, mode, logPath string) {
	t.Helper()
	dir := t.TempDir()
	script := filepath.Join(dir, "git")
	body := fmt.Sprintf("#!/bin/sh\nprintf '%%s\\n' \"$*\" >> %q\ncase \"$*\" in\n  *'--show-toplevel'*) printf '%%s\\n' %q;;\n  *'--git-common-dir'*)\n    if [ %q = fail-rev-parse ]; then exit 1; fi\n    if [ %q = bad-rev-parse ]; then printf '%%s\\n' %q; else printf '%%s\\n' %q; fi;;\n  *worktree*)\n    if [ %q = fail-worktree ]; then exit 1; fi\n    if [ %q = block-worktree ]; then sleep 600; fi\n    exit 0;;\nesac\n", logPath, root, mode, mode, filepath.Join(root, "missing-common"), filepath.Join(root, ".git"), mode, mode)
	if err := os.WriteFile(script, []byte(body), 0o755); err != nil {
		t.Fatalf("stub git: %v", err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
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
