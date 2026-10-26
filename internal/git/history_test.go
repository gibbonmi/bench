package git

import (
	"slices"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
)

// TestFirstParentsAfterListsTheChainOldestFirst proves that the list holds the first-parent
// commits after the marker up to the destination, oldest first, without the commits that a
// merge brings in through its second parent, and that a marker at the destination lists none.
func TestFirstParentsAfterListsTheChainOldestFirst(t *testing.T) {
	root := gittest.RepoOnBranch(t, "main")
	commit := func(message string) string {
		t.Helper()
		gittest.Output(t, root, "commit", "-q", "--allow-empty", "-m", message)
		return gittest.Output(t, root, "rev-parse", "HEAD")
	}
	marker := commit("marker")
	first := commit("first")
	gittest.Output(t, root, "checkout", "-q", "-b", "side")
	commit("side")
	gittest.Output(t, root, "checkout", "-q", "main")
	gittest.Output(t, root, "merge", "-q", "--no-ff", "-m", "land side", "side")
	landing := gittest.Output(t, root, "rev-parse", "HEAD")
	got, err := FirstParentsAfter(root, marker, landing)
	if err != nil || !slices.Equal(got, []string{first, landing}) {
		t.Fatalf("FirstParentsAfter = %q, %v; want %q", got, err, []string{first, landing})
	}
	if got, err := FirstParentsAfter(root, landing, landing); err != nil || len(got) != 0 {
		t.Fatalf("FirstParentsAfter of one commit = %q, %v; want none", got, err)
	}
}
