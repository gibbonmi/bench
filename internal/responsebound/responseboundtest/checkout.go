package responseboundtest

import (
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
