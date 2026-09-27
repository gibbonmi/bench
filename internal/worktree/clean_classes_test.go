package worktree

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// commitOnBranch adds one commit to a branch through the primary checkout and returns to
// main, so the fixture keeps main checked out. It returns the branch's new tip.
func commitOnBranch(t *testing.T, root, ref, name, body string) string {
	t.Helper()
	gitRun(t, root, "checkout", "-q", strings.TrimPrefix(ref, "refs/heads/"))
	commitInWorktree(t, root, name, body, "commit "+name)
	gitRun(t, root, "checkout", "-q", "main")
	return gitOutput(t, root, "rev-parse", ref)
}

// squashIntoMain folds a branch's content into main as one single-parent commit, the shape
// that only the content proofs can see.
func squashIntoMain(t *testing.T, root, ref string) {
	t.Helper()
	gitRun(t, root, "merge", "-q", "--squash", strings.TrimPrefix(ref, "refs/heads/"))
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "squash "+ref)
}

// unclaimedVerdicts reads an unclaimed plan into one "<action> <detail head>" cell per
// target. The head is the detail before its first semicolon: the class and the holder.
func unclaimedVerdicts(t *testing.T, output string) map[string]string {
	t.Helper()
	verdicts := map[string]string{}
	for _, row := range cleanupRows(output) {
		fields := cleanupRowFields(row)
		if len(fields) != len(cleanupFields) {
			t.Fatalf("row %q does not carry the %d cleanup columns", row, len(cleanupFields))
		}
		head, _, _ := strings.Cut(cleanupRowValue(fields[6]), ";")
		verdicts[cleanupRowValue(fields[0])] = fields[1] + " " + head
	}
	return verdicts
}

// TestClassifyUnclaimedRefsOverTheEdgeInventory is the class table over the ref shapes of
// the edge inventory. The class and holder spellings are written out here on purpose: a
// changed spelling in production is what RI13 and RI14 exist to catch.
func TestClassifyUnclaimedRefsOverTheEdgeInventory(t *testing.T) {
	t.Parallel()
	const landed, unique = "discard-remove class=landed", "retain class=unique"
	subsumed := func(holder string) string { return "discard-remove class=subsumed holder=" + holder }
	for _, tc := range []struct {
		name  string
		setup func(t *testing.T, root, home string) map[string]string
	}{
		{"RI4 a ref at the main tip is landed", func(t *testing.T, root, _ string) map[string]string {
			return map[string]string{addUnclaimedBranch(t, root, "a"): landed}
		}},
		{"RI62 two refs at the main tip are both landed", func(t *testing.T, root, _ string) map[string]string {
			return map[string]string{addUnclaimedBranch(t, root, "a"): landed, addUnclaimedBranch(t, root, "b"): landed}
		}},
		{"RI9 a ref under an ancestry-landed ref is landed with no holder", func(t *testing.T, root, _ string) map[string]string {
			commitInWorktree(t, root, "main.txt", "main\n", "advance main")
			holder := addUnclaimedBranch(t, root, "b")
			return map[string]string{addUnclaimedBranchAt(t, root, "a", "main~1"): landed, holder: landed}
		}},
		{"RI5 a squash-folded ref is landed", func(t *testing.T, root, _ string) map[string]string {
			ref := addUnclaimedBranch(t, root, "a")
			commitOnBranch(t, root, ref, "one.txt", "one\n")
			commitOnBranch(t, root, ref, "two.txt", "two\n")
			squashIntoMain(t, root, ref)
			return map[string]string{ref: landed}
		}},
		{"RI6 equal unique tips keep the lexically first as the root", func(t *testing.T, root, _ string) map[string]string {
			first := addUnclaimedBranch(t, root, "a")
			second := addUnclaimedBranchAt(t, root, "b", commitOnBranch(t, root, first, "pair.txt", "pair\n"))
			return map[string]string{first: unique, second: subsumed(first)}
		}},
		{"RI7 a strict ancestor of a unique ref is subsumed by it", func(t *testing.T, root, _ string) map[string]string {
			holder := addUnclaimedBranch(t, root, "b")
			below := commitOnBranch(t, root, holder, "below.txt", "below\n")
			commitOnBranch(t, root, holder, "above.txt", "above\n")
			ref := addUnclaimedBranchAt(t, root, "a", below)
			return map[string]string{ref: subsumed(holder), holder: unique}
		}},
		{"RI61 a chain names its root for every member", func(t *testing.T, root, _ string) map[string]string {
			top := addUnclaimedBranch(t, root, "c")
			first := addUnclaimedBranchAt(t, root, "a", commitOnBranch(t, root, top, "a.txt", "a\n"))
			second := addUnclaimedBranchAt(t, root, "b", commitOnBranch(t, root, top, "b.txt", "b\n"))
			commitOnBranch(t, root, top, "c.txt", "c\n")
			return map[string]string{first: subsumed(top), second: subsumed(top), top: unique}
		}},
		{"RI2 a recorded branch with no checkout prints no row", func(t *testing.T, root, home string) map[string]string {
			recorded := mustCreate(t, root, home, "classes-recorded", "recorded")
			commitInWorktree(t, recorded.Path, "recorded.txt", "recorded\n", "recorded")
			gitRun(t, root, "worktree", "remove", "-f", "-f", recorded.Path)
			return map[string]string{}
		}},
		{"RI8 a ref under an active recorded branch is subsumed by it", func(t *testing.T, root, home string) map[string]string {
			active := mustCreate(t, root, home, "classes-active-below", "active")
			commitInWorktree(t, active.Path, "below.txt", "below\n", "below")
			ref := addUnclaimedBranchAt(t, root, "a", gitOutput(t, active.Path, "rev-parse", "HEAD"))
			commitInWorktree(t, active.Path, "above.txt", "above\n", "above")
			return map[string]string{ref: subsumed(active.Assignment.Branch)}
		}},
		{"RI59 a ref at an active recorded tip is subsumed by it", func(t *testing.T, root, home string) map[string]string {
			active := mustCreate(t, root, home, "classes-active-equal", "active")
			commitInWorktree(t, active.Path, "equal.txt", "equal\n", "equal")
			ref := addUnclaimedBranchAt(t, root, "a", gitOutput(t, active.Path, "rev-parse", "HEAD"))
			return map[string]string{ref: subsumed(active.Assignment.Branch)}
		}},
		{"RI60 a ref under a checked-out foreign branch is unique", func(t *testing.T, root, _ string) map[string]string {
			gitRun(t, root, "branch", "archive/feature")
			below := commitOnBranch(t, root, "refs/heads/archive/feature", "below.txt", "below\n")
			commitOnBranch(t, root, "refs/heads/archive/feature", "above.txt", "above\n")
			gitRun(t, root, "worktree", "add", "-q", filepath.Join(t.TempDir(), "feature"), "archive/feature")
			return map[string]string{addUnclaimedBranchAt(t, root, "a", below): unique}
		}},
		{"RI81 a ref under a content-landed ref is unique", func(t *testing.T, root, _ string) map[string]string {
			first := addUnclaimedBranch(t, root, "a")
			second := addUnclaimedBranchAt(t, root, "b", commitOnBranch(t, root, first, "x.txt", "x=1\n"))
			commitOnBranch(t, root, second, "x.txt", "x=2\n")
			squashIntoMain(t, root, second)
			return map[string]string{first: unique, second: landed}
		}},
		{"RI10 a ref with content main lacks and no holder is unique", func(t *testing.T, root, _ string) map[string]string {
			ref := addUnclaimedBranch(t, root, "a")
			commitOnBranch(t, root, ref, "unique.txt", "unique\n")
			return map[string]string{ref: unique}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root, home := unclaimedBranchFixture(t)
			want := tc.setup(t, root, home)
			plan, stderr, code := runCleanup(t, root, home, "--discard-branch", "--unclaimed")
			if got := unclaimedVerdicts(t, plan); code != 0 || stderr != "" || !reflect.DeepEqual(got, want) {
				t.Fatalf("plan exit=%d stderr=%q verdicts=%q, want %q\n%s", code, stderr, got, want, plan)
			}
		})
	}
}
