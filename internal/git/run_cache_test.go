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
// answer and the path of its argument log.
func loggedStub(t *testing.T, root, mode string) (string, string) {
	t.Helper()
	logPath := filepath.Join(t.TempDir(), "argv")
	return gittest.StubGit(t, root, mode, logPath), logPath
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
	if got := gittest.StubGitStarts(t, logPath, "--git-common-dir"); got != 1 {
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
		if got := gittest.StubGitStarts(t, logPath, flag); got != 1 {
			t.Fatalf("%s starts = %d, want 1", flag, got)
		}
	}
}

// TestRunKeysTheAdminDirApartFromTheCommonDir proves that one run over a linked worktree
// answers the admin directory and the common directory apart. A primary checkout gives
// one answer to both queries, so only a linked worktree shows a read under the wrong key.
func TestRunKeysTheAdminDirApartFromTheCommonDir(t *testing.T) {
	linked := linkedWorktree(t, newRepo(t))
	run := new(Run)
	common, err := run.CommonDir(linked)
	if err != nil {
		t.Fatalf("run.CommonDir(%s): %v", linked, err)
	}
	if admin, err := run.AdminDir(linked); err != nil || admin == common {
		t.Fatalf("run.AdminDir(%s) = %q, %v, want an answer apart from the common directory %q", linked, admin, err, common)
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
	if got := gittest.StubGitStarts(t, logPath, "--git-common-dir"); got != 1 {
		t.Fatalf("common-directory starts = %d, want 1", got)
	}
}

// TestRunReadsRefsLive proves that the default branch that update-ref moves inside one run
// reads its new commit through ResolvedDefault, the package ResolveCommit, and the run's
// ResolveCommit, while the run still answers its location reads from the cache.
// (Coverage rows LGP12 and GQ45.)
func TestRunReadsRefsLive(t *testing.T) {
	root := newRepo(t)
	def, ok := ResolvedDefault(root)
	if !ok {
		t.Fatalf("ResolvedDefault(%s) found no default branch", root)
	}
	ref := "refs/heads/" + def
	old := strings.TrimSpace(runGit(t, root, "rev-parse", ref))
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
		branch, ok := ResolvedDefault(root)
		if !ok || branch != def {
			t.Fatalf("ResolvedDefault = %q, %v, want %q", branch, ok, def)
		}
		if got, err := ResolveCommit(root, "refs/heads/"+branch); err != nil || got != want {
			t.Fatalf("%s = %q, %v, want %q", branch, got, err, want)
		}
		if got, err := run.ResolveCommit(root, "refs/heads/"+branch); err != nil || got != want {
			t.Fatalf("run %s = %q, %v, want %q", branch, got, err, want)
		}
	}
	readInRun(old)
	runGit(t, root, "update-ref", ref, next, old)
	readInRun(next)
	if got := gittest.StubGitStarts(t, logPath, "--git-common-dir"); got != 1 {
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
	if got := gittest.StubGitStarts(t, logPath, "--git-common-dir"); got != 2 {
		t.Fatalf("common-directory starts = %d, want 2", got)
	}
}

// TestRunStoresNoFailedAnswer proves that a failed common-directory query stores nothing,
// so the next read in the run starts git again.
func TestRunStoresNoFailedAnswer(t *testing.T) {
	root := newRepo(t)
	_, logPath := loggedStub(t, root, "fail-rev-parse")
	run := new(Run)
	for read := range 2 {
		if _, err := run.CommonDir(root); err == nil {
			t.Fatalf("read %d: run.CommonDir(%s) succeeded, want the query failure", read+1, root)
		}
	}
	if got := gittest.StubGitStarts(t, logPath, "--git-common-dir"); got != 2 {
		t.Fatalf("common-directory starts = %d, want 2", got)
	}
}
