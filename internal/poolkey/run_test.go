package poolkey

import (
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

// TestCanonicalInReadsTheCommonDirOnce proves that the canonical root and a later
// common-directory read of one run share one git start.
func TestCanonicalInReadsTheCommonDirOnce(t *testing.T) {
	root := gittest.Repo(t)
	logPath := filepath.Join(t.TempDir(), "argv")
	gittest.StubGit(t, root, "clean", logPath)
	run := new(git.Run)
	if got := CanonicalIn(run, root); got != root {
		t.Fatalf("CanonicalIn(run, %s) = %q, want the root", root, got)
	}
	if _, err := run.CommonDir(root); err != nil {
		t.Fatalf("run.CommonDir(%s): %v", root, err)
	}
	if got := gittest.StubGitStarts(t, logPath, "--git-common-dir"); got != 1 {
		t.Fatalf("common-directory starts = %d, want 1", got)
	}
}
