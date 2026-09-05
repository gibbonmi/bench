// Package gittest owns the shared git test scaffolds package tests use.
// Repository setup and file-backed command fixtures live here. Package-specific
// fixture bytes stay with the test that grades them. Only tests import this package.
package gittest

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
)

// Repo initializes an empty repository in a fresh temporary directory. It returns that
// directory's root. A test that commits wants RepoOnBranch instead. RepoOnBranch also
// pins the branch and an identity.
func Repo(t testing.TB) string {
	t.Helper()
	return initialize(t)
}

// StubGit installs a pure file-backed git stub on the process PATH for
// worktree-resolution tests. It returns the common-directory path emitted by
// successful resolution modes. A caller that starts a child process wants
// StubGitDir instead, which binds no environment and so stays parallel-eligible.
func StubGit(t testing.TB, root, mode, logPath string) string {
	t.Helper()
	dir, commonDir := StubGitDir(t, root, mode, logPath)
	t.Setenv("PATH", dir)
	return commonDir
}

// StubGitDir writes the same stub into a fresh directory and returns that
// directory with the common-directory path. It mutates no environment, so the
// caller hands the directory to one child on its own PATH.
//
// The directory-query modes are bad-git-dir, empty-git-dir, symlink-git-dir,
// file-git-dir, block-git-dir, and fail-git-dir, deterministically answering
// (or refusing) the checkout administration directory query. A caller that
// needs a symlink or a regular file at the bad-git-dir, symlink-git-dir, or
// file-git-dir answer creates it at root joined with "missing-admin",
// "symlink-admin", or "file-admin" respectively — the same join the stub uses.
// The file-query modes are block-git-path, relative-git-path, and
// fail-git-path. In fail-git-dir, relative-git-path, and fail-git-path mode,
// every invocation other than the mode's own targeted query passes through to
// the real git the stub locates with exec.LookPath before the test replaces
// PATH.
func StubGitDir(t testing.TB, root, mode, logPath string) (string, string) {
	t.Helper()
	dir := t.TempDir()
	script := filepath.Join(dir, "git")
	commonDir := filepath.Join(root, ".git")
	switch mode {
	case "bad-rev-parse":
		commonDir = filepath.Join(root, "missing-common")
	case "empty-rev-parse":
		commonDir = ""
	case "symlink-rev-parse":
		commonDir = filepath.Join(root, "symlink-common")
	}
	badAdminDir := filepath.Join(root, "missing-admin")
	symlinkAdminDir := filepath.Join(root, "symlink-admin")
	fileAdminDir := filepath.Join(root, "file-admin")
	realGit, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("stub git: locate real git: %v", err)
	}
	body := fmt.Sprintf(`#!/bin/sh
printf '%%s\n' "$*" >> %[1]q
if [ %[2]q = fail-git-dir ]; then
  case "$*" in
    *'--git-dir'*) exit 1;;
  esac
  exec %[3]q "$@"
fi
if [ %[2]q = fail-git-path ]; then
  case "$*" in
    *'--git-path'*) exit 1;;
  esac
  exec %[3]q "$@"
fi
if [ %[2]q = relative-git-path ]; then
  case "$*" in
    *'--git-path'*)
      eval last=\${$#}
      printf '.git/%%s\n' "$last"
      exit 0;;
  esac
  exec %[3]q "$@"
fi
case "$*" in
  *'--show-toplevel'*) printf '%%s\n' %[4]q;;
  *'--git-common-dir'*)
    if [ %[2]q = fail-rev-parse ]; then exit 1; fi
    if [ %[2]q = fail-rev-parse-noisy ]; then printf 'fatal: common directory unavailable\n' >&2; exit 7; fi
    if [ %[2]q = block-rev-parse ]; then /bin/sleep 600 & wait; fi
    if [ %[2]q = noisy-list ]; then printf 'rev-parse noise\n' >&2; fi
    if [ %[2]q = vanish-after-rev-parse ]; then /bin/rm -- "$0"; fi
    if [ %[2]q = empty-rev-parse ]; then exit 0; fi
    printf '%%s\n' %[5]q;;
  *'--git-dir'*)
    if [ %[2]q = bad-git-dir ]; then printf '%%s\n' %[6]q; exit 0; fi
    if [ %[2]q = empty-git-dir ]; then exit 0; fi
    if [ %[2]q = symlink-git-dir ]; then printf '%%s\n' %[7]q; exit 0; fi
    if [ %[2]q = file-git-dir ]; then printf '%%s\n' %[8]q; exit 0; fi
    if [ %[2]q = block-git-dir ]; then /bin/sleep 600 & wait; fi
    printf '%%s\n' %[5]q;;
  *'--git-path'*)
    if [ %[2]q = block-git-path ]; then /bin/sleep 600 & wait; fi
    printf '%%s\n' %[5]q;;
  *worktree*)
    if [ %[2]q = fail-worktree ]; then exit 1; fi
    if [ %[2]q = block-worktree ]; then /bin/sleep 600 & wait; fi
    if [ %[2]q = noisy-list ]; then printf 'worktree noise\n' >&2; fi
    if [ %[2]q = noisy-list ]; then printf 'worktree %%s\000\000' %[4]q; fi
    exit 0;;
esac
`, logPath, mode, realGit, root, commonDir, badAdminDir, symlinkAdminDir, fileAdminDir)
	if err := os.WriteFile(script, []byte(body), 0o755); err != nil {
		t.Fatalf("stub git: %v", err)
	}
	return dir, commonDir
}

// RepoOnBranch initializes an empty repository on branch, with a commit identity
// configured. So a committing test neither inherits nor depends on the developer's own
// git config.
func RepoOnBranch(t testing.TB, branch string) string {
	t.Helper()
	root := initialize(t, "-b", branch)
	run(t, root, "config", "user.email", "bench@example.invalid")
	run(t, root, "config", "user.name", "bench test")
	return root
}

// TopicTrackingRepo initializes a repository whose checked-out branch is "topic" and
// whose "topic" branch tracks "origin/main". So the default branch resolves to "main",
// and a bare push has both a checked-out answer and an upstream answer. It returns the
// repository root. The remote-tracking ref is written by hand, so the fixture needs no
// second repository and no network.
func TopicTrackingRepo(t testing.TB) string {
	t.Helper()
	root := RepoOnBranch(t, "main")
	run(t, root, "commit", "-qm", "initial", "--allow-empty")
	run(t, root, "checkout", "-q", "-b", "topic")
	run(t, root, "remote", "add", "origin", "https://example.invalid/r.git")
	run(t, root, "update-ref", "refs/remotes/origin/main", "HEAD")
	run(t, root, "config", "branch.topic.remote", "origin")
	run(t, root, "config", "branch.topic.merge", "refs/heads/main")
	return root
}

// FIFOWorktreeAdmin plants a writerless FIFO gitdir under the real repository's common directory.
func FIFOWorktreeAdmin(t testing.TB, root, id string) string {
	t.Helper()
	commonDir := output(t, root, "rev-parse", "--path-format=absolute", "--git-common-dir")
	worktreeDir := filepath.Join(commonDir, "worktrees", id)
	if err := os.MkdirAll(worktreeDir, 0o755); err != nil {
		t.Fatalf("create worktree admin directory: %v", err)
	}
	gitDir := filepath.Join(worktreeDir, "gitdir")
	if err := syscall.Mkfifo(gitDir, 0o600); err != nil {
		capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
	}
	return gitDir
}

func initialize(t testing.TB, options ...string) string {
	t.Helper()
	root := t.TempDir()
	run(t, root, append([]string{"init", "-q"}, options...)...)
	return root
}

func run(t testing.TB, root string, args ...string) {
	t.Helper()
	_ = output(t, root, args...)
}

func output(t testing.TB, root string, args ...string) string {
	t.Helper()
	out, err := exec.Command("git", append([]string{"-C", root}, args...)...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v: %s", strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(string(out))
}
