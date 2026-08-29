package landing

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate/authorization"
)

// mergeRequest binds one target checkout's observed fingerprint to the pair under test,
// so every row below states only what it varies.
func mergeRequest(t *testing.T, root, branch, previous, incoming string) MergeRequest {
	t.Helper()
	fingerprint, err := CheckoutFingerprint(root)
	if err != nil {
		t.Fatal(err)
	}
	return MergeRequest{
		Root: root, Branch: "refs/heads/" + branch, PreviousTip: previous, Incoming: incoming,
		Worktree: root, Fingerprint: fingerprint, Subject: "merge: compose incoming into target",
	}
}

// linearPair commits one descendant on "source" and leaves "destination" at the base,
// so the pair is linear and the caller chooses which end is the target.
func linearPair(t *testing.T, root string) (base, tip string) {
	t.Helper()
	base = git(t, root, "rev-parse", "HEAD")
	git(t, root, "checkout", "-qb", "source", base)
	write(t, root, "added-by-source", "source\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "source")
	git(t, root, "branch", "destination", base)
	return base, git(t, root, "rev-parse", "HEAD")
}

// divergedPair commits one edit on each side of base and checks the target out.
func divergedPair(t *testing.T, root string, destinationChange, sourceChange func()) (destination, source string) {
	t.Helper()
	destination, source = commitSides(t, root, git(t, root, "rev-parse", "HEAD"), destinationChange, sourceChange)
	git(t, root, "checkout", "-q", "destination")
	return destination, source
}

func TestMergePublishesPreviousTipAsFirstParent(t *testing.T) {
	root := fixture(t)
	destination, source := divergedPair(t, root,
		func() { write(t, root, "named", "destination\n") },
		func() { write(t, root, "foreign", "source\n") })
	got, err := greenOwner().Merge(context.Background(), mergeRequest(t, root, "destination", destination, source))
	if err != nil {
		t.Fatal(err)
	}
	if got.Kind != MergeKindMerge {
		t.Fatalf("kind = %q, want %q", got.Kind, MergeKindMerge)
	}
	if first := git(t, root, "rev-parse", got.Tip+"^1"); first != destination {
		t.Fatalf("first parent = %s, want the previous tip %s", first, destination)
	}
	if second := git(t, root, "rev-parse", got.Tip+"^2"); second != source {
		t.Fatalf("second parent = %s, want the incoming commit %s", second, source)
	}
	if tip := git(t, root, "rev-parse", "refs/heads/destination"); tip != got.Tip {
		t.Fatalf("branch tip = %s, want the published commit %s", tip, got.Tip)
	}
	if tree := git(t, root, "rev-parse", got.Tip+"^{tree}"); tree != got.Tree {
		t.Fatalf("published tree = %s, want the authorized tree %s", tree, got.Tree)
	}
	if got.PreviousTip != destination {
		t.Fatalf("previous tip = %s, want %s", got.PreviousTip, destination)
	}
}

func TestMergeRefusesWhenTheBranchTipMovesUnderTheRefUpdate(t *testing.T) {
	root := fixture(t)
	destination, source := divergedPair(t, root,
		func() { write(t, root, "named", "destination\n") },
		func() { write(t, root, "foreign", "source\n") })
	owner := greenOwner()
	real := owner.updateRef
	// The concurrent commit lands between the fingerprint recheck and the update, so the
	// compare-and-swap alone can refuse it.
	owner.updateRef = func(root, ref, new, old string) error {
		git(t, root, "update-ref", ref, source, destination)
		return real(root, ref, new, old)
	}
	_, err := owner.Merge(context.Background(), mergeRequest(t, root, "destination", destination, source))
	if err == nil {
		t.Fatal("merge published over a moved branch tip")
	}
	if !strings.Contains(err.Error(), "compare-and-swap refused") {
		t.Fatalf("error = %v, want the destination compare-and-swap refusal", err)
	}
	if tip := git(t, root, "rev-parse", "refs/heads/destination"); tip != source {
		t.Fatalf("branch tip = %s, want the moved tip %s", tip, source)
	}
}

func TestMergeFastForwardsWithoutANewObject(t *testing.T) {
	root := fixture(t)
	base, tip := linearPair(t, root)
	git(t, root, "checkout", "-q", "destination")
	got, err := greenOwner().Merge(context.Background(), mergeRequest(t, root, "destination", base, tip))
	if err != nil {
		t.Fatal(err)
	}
	if got.Kind != MergeKindFastForward {
		t.Fatalf("kind = %q, want %q", got.Kind, MergeKindFastForward)
	}
	if got.Tip != tip {
		t.Fatalf("tip = %s, want the incoming commit %s with no new object", got.Tip, tip)
	}
	if got.Tree != git(t, root, "rev-parse", tip+"^{tree}") {
		t.Fatalf("tree = %s, want the incoming commit's tree", got.Tree)
	}
	if branch := git(t, root, "rev-parse", "refs/heads/destination"); branch != tip {
		t.Fatalf("branch tip = %s, want %s", branch, tip)
	}
}

func TestMergeReportsCurrentWithoutAuthorizingOrUpdating(t *testing.T) {
	root := fixture(t)
	base, tip := linearPair(t, root)
	owner := New()
	owner.authorize = func(context.Context, string, string, io.Writer, io.Writer) authorization.Result {
		t.Error("a current target authorized a tree")
		return authorization.Result{Kind: authorization.Green}
	}
	owner.updateRef = func(string, string, string, string) error {
		t.Error("a current target updated the branch ref")
		return nil
	}
	got, err := owner.Merge(context.Background(), mergeRequest(t, root, "source", tip, base))
	if err != nil {
		t.Fatal(err)
	}
	if got.Kind != MergeKindCurrent {
		t.Fatalf("kind = %q, want %q", got.Kind, MergeKindCurrent)
	}
	if got.Tip != tip || got.PreviousTip != tip {
		t.Fatalf("tip = %s, previous = %s, want both at %s", got.Tip, got.PreviousTip, tip)
	}
	if branch := git(t, root, "rev-parse", "refs/heads/source"); branch != tip {
		t.Fatalf("branch tip = %s, want the unchanged tip %s", branch, tip)
	}
}

func TestMergeRefusesAConflictOutsideTheCaptureTable(t *testing.T) {
	root := fixture(t)
	destination, _ := divergedPair(t, root,
		func() { write(t, root, "named", "destination\n") },
		func() { write(t, root, "named", "source\n") })
	source := git(t, root, "rev-parse", "refs/heads/source")
	_, err := greenOwner().Merge(context.Background(), mergeRequest(t, root, "destination", destination, source))
	var conflict ConflictError
	if !errors.As(err, &conflict) {
		t.Fatalf("error = %v, want a composition conflict", err)
	}
	if conflict.Kind != "textual" || strings.Join(conflict.Paths, ",") != "named" {
		t.Fatalf("conflict = %q %v, want textual at named", conflict.Kind, conflict.Paths)
	}
	if tip := git(t, root, "rev-parse", "refs/heads/destination"); tip != destination {
		t.Fatalf("branch tip = %s, want the unchanged tip %s", tip, destination)
	}
}

func TestMergeResolvesAConflictedLearningsJournalAsAUnion(t *testing.T) {
	root := fixture(t)
	write(t, root, "capture/learnings.md", "shared\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "capture base")
	destination, source := divergedPair(t, root,
		func() { write(t, root, "capture/learnings.md", "shared\ndestination\n") },
		func() { write(t, root, "capture/learnings.md", "shared\nsource\n") })
	got, err := greenOwner().Merge(context.Background(), mergeRequest(t, root, "destination", destination, source))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(got.Resolved, ",") != "capture/learnings.md:union" {
		t.Fatalf("resolved = %v, want the union of the learnings journal", got.Resolved)
	}
	content := git(t, root, "show", got.Tip+":capture/learnings.md")
	if !strings.Contains(content, "destination") || !strings.Contains(content, "source") {
		t.Fatalf("published journal = %q, want both sides", content)
	}
}

// commitObjectCount counts every commit object the repository holds, reachable or not.
// Composition writes trees and blobs, so only this count isolates a minted commit.
func commitObjectCount(t *testing.T, root string) int {
	t.Helper()
	types := git(t, root, "cat-file", "--batch-all-objects", "--batch-check=%(objecttype)")
	count := 0
	for _, line := range strings.Fields(types) {
		if line == "commit" {
			count++
		}
	}
	return count
}

func TestMergeRefusesAnUnauthorizedTree(t *testing.T) {
	// Both kinds grade a tree, so both must refuse one: the fast-forward path grades the
	// incoming tree rather than a composed one.
	cases := []struct {
		name  string
		setup func(t *testing.T, root string) (branch, previous, incoming string)
	}{
		{name: "merge", setup: func(t *testing.T, root string) (string, string, string) {
			destination, source := divergedPair(t, root,
				func() { write(t, root, "named", "destination\n") },
				func() { write(t, root, "foreign", "source\n") })
			return "destination", destination, source
		}},
		{name: "fast-forward", setup: func(t *testing.T, root string) (string, string, string) {
			base, tip := linearPair(t, root)
			git(t, root, "checkout", "-q", "destination")
			return "destination", base, tip
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := fixture(t)
			branch, previous, incoming := tc.setup(t, root)
			owner := New()
			owner.authorize = func(context.Context, string, string, io.Writer, io.Writer) authorization.Result {
				return authorization.Result{Kind: authorization.LaneFail}
			}
			owner.updateRef = func(string, string, string, string) error {
				t.Error("a refused authorization updated the branch ref")
				return nil
			}
			commits := commitObjectCount(t, root)
			_, err := owner.Merge(context.Background(), mergeRequest(t, root, branch, previous, incoming))
			if err == nil {
				t.Fatal("merge published a tree the authority refused")
			}
			if !strings.Contains(err.Error(), string(authorization.LaneFail)) {
				t.Fatalf("error = %v, want the refused authorization kind", err)
			}
			if tip := git(t, root, "rev-parse", "refs/heads/"+branch); tip != previous {
				t.Fatalf("branch tip = %s, want the unchanged tip %s", tip, previous)
			}
			if got := commitObjectCount(t, root); got != commits {
				t.Fatalf("commit objects = %d, want the unchanged %d", got, commits)
			}
		})
	}
}
