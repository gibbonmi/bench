package responseboundtest

import (
	"io/fs"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/poolkey"
)

// AssignmentCheckout makes a linked worktree of a new repository at an assignment segment
// of the pool below home. It answers the repository, the checkout, and the assignment id.
// A test that needs the spill scope of an assignment worktree takes its tree from here.
func AssignmentCheckout(t testing.TB, home string) (repo, checkout, id string) {
	t.Helper()
	repo = gittest.RepoOnBranch(t, "main")
	id = strings.Repeat("b", 32)
	checkout = filepath.Join(poolkey.Pool(home, repo), poolkey.AssignmentSegment(strings.Repeat("a", 32), id))
	for _, args := range [][]string{{"commit", "-q", "--allow-empty", "-m", "base"}, {"worktree", "add", "-q", "--detach", checkout}} {
		if out, err := exec.Command("git", append([]string{"-C", repo}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %q: %v: %s", args, err, out)
		}
	}
	return repo, checkout, id
}

// StoreFiles answers every regular file below home outside the worktree pools, in lexical
// order, so a test sees where the response store wrote without naming its layout.
func StoreFiles(t testing.TB, home string) []string {
	t.Helper()
	var files []string
	err := filepath.WalkDir(home, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() && path == poolkey.Pools(home) {
			return filepath.SkipDir
		}
		if entry.Type().IsRegular() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}
