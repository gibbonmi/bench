// The worktree merge verb's `--from` sha route: a sibling's branch tip sha composes the
// same fold its label does, so an operator who reads the handoff can pin the sha the
// handoff names rather than retype the sibling's label.
package worktree

import (
	"path/filepath"
	"testing"
)

// A sibling's branch tip sha resolves through the same sibling lookup its label
// does, and it composes the identical fold. The handoff pins a sha, not a label, so this
// is the route an operator actually types.
func TestMergeFoldsASiblingBranchTipBySha(t *testing.T) {
	t.Parallel()
	j, root, home, _, created := mergeFixture(t, "integration", "delegate")
	target, sibling := created[0], created[1]
	commitInWorktree(t, target.Path, "target.txt", "target\n", "target work")
	previous := gitOutput(t, target.Path, "rev-parse", "HEAD")
	commitInWorktree(t, sibling.Path, "sibling.txt", "sibling\n", "sibling work")
	siblingTip := gitOutput(t, sibling.Path, "rev-parse", "HEAD")

	code, stdout, stderr := runMerge(t, j, root, home, "--from", siblingTip, target.Assignment.ID)
	if code != 0 {
		t.Fatalf("merge exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	tip := gitOutput(t, root, "rev-parse", target.Assignment.Branch)
	if first := gitOutput(t, root, "rev-parse", tip+"^1"); first != previous {
		t.Fatalf("first parent = %s, want the previous tip %s", first, previous)
	}
	if second := gitOutput(t, root, "rev-parse", tip+"^2"); second != siblingTip {
		t.Fatalf("second parent = %s, want the sibling's branch tip %s", second, siblingTip)
	}
}

// The sha route runs the same sibling checks the label route does. A dirty or
// detached sibling refuses with the existing sibling refusal details, because the sha
// resolves to the same assignment the label would.
func TestMergeRefusesADirtyOrDetachedSiblingBySha(t *testing.T) {
	t.Parallel()
	j, root, home, tally, created := mergeFixture(t, "integration", "delegate")
	target, sibling := created[0], created[1]
	commitInWorktree(t, sibling.Path, "sibling.txt", "sibling\n", "sibling work")
	previous := gitOutput(t, target.Path, "rev-parse", "HEAD")
	siblingTip := gitOutput(t, sibling.Path, "rev-parse", "HEAD")
	mustWrite(t, filepath.Join(sibling.Path, "sibling.txt"), []byte("uncommitted\n"), 0o644)

	code, stdout, stderr := runMerge(t, j, root, home, "--from", siblingTip, target.Assignment.ID)
	requireMergeRefusal(t, code, stdout, stderr, "sibling checkout is not clean",
		"next=bench worktree exec "+sibling.Assignment.ID+" -- bench commit", "sibling.txt")
	requireMergeUnchanged(t, root, target.Path, target.Assignment.Branch, previous, tally)

	gitRun(t, sibling.Path, "checkout", "-q", "--", "sibling.txt")
	gitRun(t, sibling.Path, "checkout", "-q", "--detach", "HEAD")
	code, stdout, stderr = runMerge(t, j, root, home, "--from", siblingTip, target.Assignment.ID)
	requireMergeRefusal(t, code, stdout, stderr, sibling.Assignment.Branch)
	requireMergeUnchanged(t, root, target.Path, target.Assignment.Branch, previous, tally)
}
