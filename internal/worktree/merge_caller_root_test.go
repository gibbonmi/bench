// The worktree merge verb's caller-root row: the lane grades the composed tree whatever
// directory the caller runs the verb from.
package worktree

import (
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
)

// FT335: a merge run from the primary checkout grades the prose the sibling brings in
// from the composed tree. The lane resolves against the target worktree, so its prose
// check names the target as its file root, as the built-in lane anchors `gate-prose`.
// An authority rooted at the caller's checkout leaves that anchor unmoved, and the check
// reads the incoming Markdown from the target checkout, where the merge has not put it.
func TestMergeGradesIncomingProseFromTheComposedTreeWhateverTheCallerRoot(t *testing.T) {
	t.Parallel()
	j, root, home, _, created := mergeFixture(t, "integration", "sibling")
	target, sibling := created[0], created[1]
	commitInWorktree(t, target.Path, "target-only.md", "target prose\n", "target prose")
	commitInWorktree(t, sibling.Path, "sibling-only.md", "sibling prose\n", "sibling prose")
	j.mergeLane = func(anchor string) (*gate.Lane, error) {
		return &gate.Lane{Checks: []gate.Phase{{Name: "prose", Argv: []string{"sh", "-c",
			`anchor=$1; shift; for path; do test -f "$anchor/$path" || { echo "unreadable $path"; exit 1; }; done`,
			"prose", anchor, gate.LaneNamedMarkdownToken}}}}, nil
	}
	if root == target.Path {
		t.Fatal("the fixture runs the merge from the target, so the caller root never differs")
	}

	code, stdout, stderr := runMerge(t, j, root, home, "--from", sibling.Assignment.Label, target.Assignment.ID)
	if code != 0 {
		t.Fatalf("merge exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	mergedRecord(t, stdout)
}
