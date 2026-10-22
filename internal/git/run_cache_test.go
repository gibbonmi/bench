package git

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
)

// loggedStub binds the git stub for root on PATH. It returns the stub's common-directory
// answer and the path of its argument log, which gains one line for each git start.
func loggedStub(t *testing.T, root, mode string) (string, string) {
	t.Helper()
	logPath := filepath.Join(t.TempDir(), "argv")
	return gittest.StubGit(t, root, mode, logPath), logPath
}

// gitStarts counts the git starts in the stub argument log whose arguments hold flag.
func gitStarts(t *testing.T, logPath, flag string) int {
	t.Helper()
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	return strings.Count(string(data), flag)
}

// TestRunCachesTheCommonDir proves that one run starts the common-directory query once
// for a root, and that the run's Worktrees form reads the common directory through the
// run. (Coverage row LGP9.)
func TestRunCachesTheCommonDir(t *testing.T) {
	root := newRepo(t)
	want, logPath := loggedStub(t, root, "clean")
	run := new(Run)
	for range 2 {
		if got, err := run.CommonDir(root); err != nil || got != want {
			t.Fatalf("run.CommonDir(%s) = %q, %v, want %q, nil", root, got, err, want)
		}
	}
	if _, err := run.Worktrees(root); err != nil {
		t.Fatalf("run.Worktrees(%s): %v", root, err)
	}
	if got := gitStarts(t, logPath, "--git-common-dir"); got != 1 {
		t.Fatalf("common-directory starts = %d, want 1", got)
	}
}

// TestRunCachesTopLevelAndAdminDir proves that one run starts the top-level query and
// the admin-directory query once each for a root. (Coverage row LGP10.)
func TestRunCachesTopLevelAndAdminDir(t *testing.T) {
	root := newRepo(t)
	admin, logPath := loggedStub(t, root, "clean")
	run := new(Run)
	for range 2 {
		if got, err := run.RootAt(root); err != nil || got != root {
			t.Fatalf("run.RootAt(%s) = %q, %v, want the root", root, got, err)
		}
		if got, err := run.AdminDir(root); err != nil || got != admin {
			t.Fatalf("run.AdminDir(%s) = %q, %v, want %q, nil", root, got, err, admin)
		}
	}
	for _, flag := range []string{"--show-toplevel", "--git-dir"} {
		if got := gitStarts(t, logPath, flag); got != 1 {
			t.Fatalf("%s starts = %d, want 1", flag, got)
		}
	}
}

// TestRunRevalidatesACachedCommonDir proves that a cached common directory, which a
// symlink replaced after the first read, refuses on the second read without a second
// git start. (Coverage row LGP11.)
func TestRunRevalidatesACachedCommonDir(t *testing.T) {
	root := newRepo(t)
	common, logPath := loggedStub(t, root, "clean")
	run := new(Run)
	if _, err := run.CommonDir(root); err != nil {
		t.Fatalf("first read: %v", err)
	}
	moved := common + "-moved"
	if err := os.Rename(common, moved); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(moved, common); err != nil {
		t.Fatal(err)
	}
	_, err := run.CommonDir(root)
	var resolution *ResolutionError
	if !errors.As(err, &resolution) || resolution.Path != common || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("second read = %v, want the symlink refusal of %q", err, common)
	}
	if got := gitStarts(t, logPath, "--git-common-dir"); got != 1 {
		t.Fatalf("common-directory starts = %d, want 1", got)
	}
}

// TestRunReadsRefsLive proves that a branch that update-ref moves inside one run reads
// its new commit, while the run still answers its location reads from the cache.
// (Coverage row LGP12.)
func TestRunReadsRefsLive(t *testing.T) {
	root := newRepo(t)
	runGit(t, root, "branch", "topic")
	old := strings.TrimSpace(runGit(t, root, "rev-parse", "topic"))
	next := strings.TrimSpace(runGit(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit-tree", "-p", old, "-m", "next", old+"^{tree}"))
	_, logPath := loggedStub(t, root, "fail-git-path")
	run := new(Run)
	readInRun := func(want string) {
		t.Helper()
		if _, err := run.CommonDir(root); err != nil {
			t.Fatalf("run.CommonDir: %v", err)
		}
		if _, err := run.RootAt(root); err != nil {
			t.Fatalf("run.RootAt: %v", err)
		}
		if got, err := ResolveCommit(root, "refs/heads/topic"); err != nil || got != want {
			t.Fatalf("topic = %q, %v, want %q", got, err, want)
		}
	}
	readInRun(old)
	runGit(t, root, "update-ref", "refs/heads/topic", next, old)
	readInRun(next)
	if got := gitStarts(t, logPath, "--git-common-dir"); got != 1 {
		t.Fatalf("common-directory starts = %d, want 1 inside the run", got)
	}
}

// TestRunValuesShareNoCache proves that two run values over one root each start their
// own common-directory query. (Coverage row LGP14.)
func TestRunValuesShareNoCache(t *testing.T) {
	root := newRepo(t)
	_, logPath := loggedStub(t, root, "clean")
	for _, run := range []*Run{new(Run), new(Run)} {
		if _, err := run.CommonDir(root); err != nil {
			t.Fatalf("run.CommonDir(%s): %v", root, err)
		}
	}
	if got := gitStarts(t, logPath, "--git-common-dir"); got != 2 {
		t.Fatalf("common-directory starts = %d, want 2", got)
	}
}
