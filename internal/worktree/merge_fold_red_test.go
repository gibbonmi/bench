package worktree

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/refusalroute"
)

// foldRedFixture produces the fold red of one red cause: main adds the red file to a target
// whose tip alone grades green. target makes the target set, and it returns the set.
func foldRedFixture(cause string, target func(t *testing.T) mergeSet) mergeRefusalFixture {
	return mergeRefusalFixture{
		face:  faceMergeFoldRed,
		cause: cause,
		exit:  1,
		build: func(t *testing.T) (mergeSet, []string) {
			f := target(t)
			commitOnDefault(t, f.root, redFile, "red\n")
			return f, mergeTargetArgs(f, "main")
		},
		names: func(f mergeSet) []string {
			target := f.created[0].Assignment
			checkoutGit := conflictFacts("", target.ID, target.Worktree)[refusalroute.FactCheckoutGit]
			return []string{target.Label, checkoutGit, "bench commit", "bench record fold"}
		},
		absent: func(mergeSet) []string { return []string{"FT342"} },
		carry: map[int]func(*testing.T, mergeSet){
			// The reviewer reads the cause: the target tip alone holds no red.
			0: func(t *testing.T, f mergeSet) {
				if _, err := os.Stat(filepath.Join(f.created[0].Path, redFile)); !os.IsNotExist(err) {
					t.Fatalf("target holds %s: %v, want the red from main alone", redFile, err)
				}
			},
			// The hand merge of main repairs the red in the merge commit.
			1: func(t *testing.T, f mergeSet) {
				gitRun(t, f.created[0].Path, "merge", "-q", "--no-ff", "--no-commit", "main")
				gitRun(t, f.created[0].Path, "rm", "-q", "-f", redFile)
			},
			2: func(t *testing.T, f mergeSet) { commitHandMerge(t, f.created[0].Path) },
			// The fixture binds no spec, so no review record takes the fold entry.
			3: func(*testing.T, mergeSet) {},
		},
	}
}
