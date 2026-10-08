// The worktree merge verb's `--from` sha route: a sibling's branch tip sha composes the
// same fold its label does, so an operator who reads the handoff can pin the sha the
// handoff names rather than retype the sibling's label.
package worktree

import (
	"path/filepath"
	"strings"
	"testing"
)

// requireSiblingFold composes the fold one `--from` spelling addresses, and it pins what
// the fold publishes: the two parents, the merge kind, and the subject. Both spellings of
// one sibling run it, so the sha route cannot drift from the label route. spell answers
// the `--from` value for the sibling and its branch tip.
func requireSiblingFold(t *testing.T, spell func(sibling Creation, tip string) string) {
	t.Helper()
	f := mergeFixture(t, "integration", "delegate")
	target, sibling := f.created[0], f.created[1]
	commitInWorktree(t, target.Path, "target.txt", "target\n", "target work")
	previous := gitOutput(t, target.Path, "rev-parse", "HEAD")
	commitInWorktree(t, sibling.Path, "sibling.txt", "sibling\n", "sibling work")
	siblingTip := gitOutput(t, sibling.Path, "rev-parse", "HEAD")
	from := spell(sibling, siblingTip)

	r := runVerb(t, verbMerge, f.merge("--from", from, target.Assignment.ID))
	if r.exit != 0 {
		t.Fatalf("merge %q exit = %d, want 0; stdout=%q stderr=%q", from, r.exit, r.stdout, r.stderr)
	}
	tip := gitOutput(t, f.root, "rev-parse", target.Assignment.Branch)
	if first := gitOutput(t, f.root, "rev-parse", tip+"^1"); first != previous {
		t.Fatalf("first parent = %s, want the previous tip %s", first, previous)
	}
	if second := gitOutput(t, f.root, "rev-parse", tip+"^2"); second != siblingTip {
		t.Fatalf("second parent = %s, want the sibling's branch tip %s", second, siblingTip)
	}
	if record := mergedRecord(t, r.stdout); !strings.Contains(record, "kind=merge") {
		t.Fatalf("record = %q, want kind=merge", record)
	}
	want := "merge: compose " + from + " " + siblingTip[:8] + " into " + target.Assignment.Label
	if got := gitOutput(t, f.root, "log", "-1", "--format=%s", tip); got != want {
		t.Fatalf("subject = %q, want %q", got, want)
	}
}

// A sibling's branch tip sha resolves through the same sibling checks its label does, and
// it composes the identical fold. The handoff pins a sha, not a label, so this is the
// route an operator actually types.
func TestMergeFoldsASiblingBranchTipBySha(t *testing.T) {
	t.Parallel()
	requireSiblingFold(t, func(_ Creation, tip string) string { return tip })
}

// A fresh sibling sits at the default-branch tip, and the ordinary fold of that tip stays
// a commit fold. A sha route that runs before the commit lookup reads the fresh sibling as
// a second answer and refuses the ordinary fold as ambiguous.
func TestMergeFoldsTheDefaultBranchTipAFreshSiblingSitsOn(t *testing.T) {
	t.Parallel()
	f := mergeFixture(t, "integration")
	target := f.created[0]
	commitInWorktree(t, target.Path, "target.txt", "target\n", "target work")
	previous := gitOutput(t, target.Path, "rev-parse", "HEAD")
	incoming := commitOnDefault(t, f.root, "incoming.txt", "incoming\n")
	sibling := mustCreate(t, f.root, f.home, "merge-fresh", "fresh")
	if tip := gitOutput(t, f.root, "rev-parse", sibling.Assignment.Branch); tip != incoming {
		t.Fatalf("premise failed: the fresh sibling is at %s, want the default-branch tip %s", tip, incoming)
	}

	r := runVerb(t, verbMerge, f.merge("--from", incoming, target.Assignment.ID))
	if r.exit != 0 {
		t.Fatalf("merge exit = %d, want 0; stdout=%q stderr=%q", r.exit, r.stdout, r.stderr)
	}
	tip := gitOutput(t, f.root, "rev-parse", target.Assignment.Branch)
	if first := gitOutput(t, f.root, "rev-parse", tip+"^1"); first != previous {
		t.Fatalf("first parent = %s, want the previous tip %s", first, previous)
	}
	if second := gitOutput(t, f.root, "rev-parse", tip+"^2"); second != incoming {
		t.Fatalf("second parent = %s, want the default-branch tip %s", second, incoming)
	}
}

// The sha route runs the same sibling checks the label route does. A dirty or detached
// sibling refuses with the existing sibling refusal details, because the sha resolves to
// the same assignment the label would.
func TestMergeRefusesADirtyOrDetachedSiblingBySha(t *testing.T) {
	t.Parallel()
	f := mergeFixture(t, "integration", "delegate")
	target, sibling := f.created[0], f.created[1]
	commitInWorktree(t, sibling.Path, "sibling.txt", "sibling\n", "sibling work")
	previous := gitOutput(t, target.Path, "rev-parse", "HEAD")
	siblingTip := gitOutput(t, sibling.Path, "rev-parse", "HEAD")
	mustWrite(t, filepath.Join(sibling.Path, "sibling.txt"), []byte("uncommitted\n"), 0o644)

	r := runVerb(t, verbMerge, f.merge("--from", siblingTip, target.Assignment.ID))
	requireMergeRefusal(t, r, "sibling checkout is not clean",
		"next=bench commit --in '"+sibling.Assignment.Label+"' ", "sibling.txt")
	requireMergeUnchanged(t, f, target, previous)

	gitRun(t, sibling.Path, "checkout", "-q", "--", "sibling.txt")
	gitRun(t, sibling.Path, "checkout", "-q", "--detach", "HEAD")
	r = runVerb(t, verbMerge, f.merge("--from", siblingTip, target.Assignment.ID))
	requireMergeRefusal(t, r, sibling.Assignment.Branch)
	requireMergeUnchanged(t, f, target, previous)
}

// The sha route is the merge verb's alone, so the create verb's `--from` keeps its
// label-only contract: a sibling's tip sha names no active assignment there. A create
// verb that shared the sha route would start a worktree from a bare commit.
func TestCreateFromRefusesASiblingTipSha(t *testing.T) {
	t.Parallel()
	f := mergeFixture(t, "delegate")
	sibling := f.created[0]
	commitInWorktree(t, sibling.Path, "sibling.txt", "sibling\n", "sibling work")
	siblingTip := gitOutput(t, sibling.Path, "rev-parse", "HEAD")
	before := assignmentCount(t, f.root)

	requireCreateFromRefusal(t, runVerb(t, verbCreate, f.call(
		"--request", "create-from-sha", "--label", "dependent", "--from", siblingTip)),
		"bench worktree create: --from names no active assignment\n", "next=bench worktree list\n")
	if after := assignmentCount(t, f.root); after != before {
		t.Fatalf("ledger holds %d records, want the %d it held before the refusal", after, before)
	}
}

// A sha two siblings share names no one sibling, so it refuses and names both ids. A
// first-match resolver folds whichever sibling the ledger happens to list first.
func TestMergeRefusesAShaAtTwoSiblingTips(t *testing.T) {
	t.Parallel()
	f := mergeFixture(t, "integration", "delegate", "second-delegate")
	target, sibling, other := f.created[0], f.created[1], f.created[2]
	previous := gitOutput(t, target.Path, "rev-parse", "HEAD")
	commitInWorktree(t, sibling.Path, "sibling.txt", "sibling\n", "sibling work")
	shared := gitOutput(t, sibling.Path, "rev-parse", "HEAD")
	gitRun(t, other.Path, "merge", "-q", "--ff-only", shared)

	r := runVerb(t, verbMerge, f.merge("--from", shared, target.Assignment.ID))
	requireMergeRefusal(t, r, "--from is the branch tip of more than one sibling",
		"observed="+shared, sibling.Assignment.ID, other.Assignment.ID)
	requireMergeUnchanged(t, f, target, previous)
}

// A sha that is the target's own branch tip refuses the way the target's own label does.
// A self-fold prints `kind=current` and hides the operand error.
func TestMergeRefusesAShaThatIsTheTargetsOwnTip(t *testing.T) {
	t.Parallel()
	f := mergeFixture(t, "integration")
	target := f.created[0]
	commitInWorktree(t, target.Path, "target.txt", "target\n", "target work")
	previous := gitOutput(t, target.Path, "rev-parse", "HEAD")

	r := runVerb(t, verbMerge, f.merge("--from", previous, target.Assignment.ID))
	requireMergeRefusal(t, r, "--from resolves to the target itself",
		"observed="+target.Assignment.ID)
	requireMergeUnchanged(t, f, target, previous)
}
