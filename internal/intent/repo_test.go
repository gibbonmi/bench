package intent

import (
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

// TestRepoReadsTheCommonDirOnce proves that two ledger reads through one Repo with a run
// start the common-directory query once, and that two root-taking reads start it twice.
// (Coverage rows GQ27 and GQ28.)
func TestRepoReadsTheCommonDirOnce(t *testing.T) {
	root := newRepo(t)
	for _, c := range []struct {
		name   string
		read   func() (Ledger, error)
		starts int
	}{
		{"run", In(new(git.Run), root).Read, 1},
		{"root", func() (Ledger, error) { return Read(root) }, 2},
	} {
		t.Run(c.name, func(t *testing.T) {
			logPath := filepath.Join(t.TempDir(), "argv")
			gittest.StubGit(t, root, "log", logPath)
			for call := range 2 {
				if _, err := c.read(); err != nil {
					t.Fatalf("read %d: %v", call+1, err)
				}
			}
			if got := gittest.StubGitStarts(t, logPath, "--git-common-dir"); got != c.starts {
				t.Fatalf("common-directory starts = %d, want %d", got, c.starts)
			}
		})
	}
}
