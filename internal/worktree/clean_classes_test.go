package worktree

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
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

// recordedBranch plants a recorded branch in state that adds x=1 and then changes it to x=2,
// with the unrecorded ref a at its first commit. It returns ref a and the recorded branch.
// A squash-folded recorded branch is landed by content only, and ref a is not landed.
func recordedBranch(t *testing.T, root, home string, state intent.AssignmentState, squashed bool) (string, string) {
	t.Helper()
	recorded := mustCreate(t, root, home, "classes-recorded-"+string(state), "recorded")
	commitInWorktree(t, recorded.Path, "x.txt", "x=1\n", "x=1")
	ref := addUnclaimedBranchAt(t, root, "a", gitOutput(t, recorded.Path, "rev-parse", "HEAD"))
	commitInWorktree(t, recorded.Path, "x.txt", "x=2\n", "x=2")
	recorded.Assignment.State = state
	mustNoError(t, intent.PutAssignment(root, recorded.Assignment))
	if squashed {
		squashLand(t, root, recorded.Assignment.Branch)
		if landed, byContent := verdict(t, root, recorded.Assignment.Branch); !landed || !byContent {
			t.Fatalf("recorded branch landed=%t byContent=%t, want landed by content only", landed, byContent)
		}
		if landed, _ := verdict(t, root, ref); landed {
			t.Fatalf("ref %s is landed, want it not landed", ref)
		}
	}
	return ref, recorded.Assignment.Branch
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
// the edge inventory. Its class and holder spellings are independent expectations pinned
// for RI13 and RI14.
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
			squashLand(t, root, ref)
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
			ref, branch := recordedBranch(t, root, home, intent.StateActive, false)
			return map[string]string{ref: subsumed(branch)}
		}},
		{"RI87 a ref under a cleanup-pending recorded branch is subsumed by it", func(t *testing.T, root, home string) map[string]string {
			ref, branch := recordedBranch(t, root, home, intent.StateCleanupPending, false)
			return map[string]string{ref: subsumed(branch)}
		}},
		{"RI86 a ref under a complete record's branch is unique", func(t *testing.T, root, home string) map[string]string {
			ref, _ := recordedBranch(t, root, home, intent.StateComplete, false)
			return map[string]string{ref: unique}
		}},
		{"RI84 a ref under an active content-landed recorded branch is unique", func(t *testing.T, root, home string) map[string]string {
			ref, _ := recordedBranch(t, root, home, intent.StateActive, true)
			return map[string]string{ref: unique}
		}},
		{"RI84 a ref under a cleanup-pending content-landed recorded branch is unique", func(t *testing.T, root, home string) map[string]string {
			ref, _ := recordedBranch(t, root, home, intent.StateCleanupPending, true)
			return map[string]string{ref: unique}
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
			squashLand(t, root, second)
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

// TestCleanUnclaimedSymrefFailsClosed is RI83. A symref in either namespace, at the
// checked-out main or at a lexically earlier unique root, gives an error row, and both apply
// forms refuse the set while the target and a removable sibling survive.
func TestCleanUnclaimedSymrefFailsClosed(t *testing.T) {
	t.Parallel()
	for _, symref := range []string{unclaimedBranchRef("e"), intent.ShiftBranchPrefix() + "20260101-000000"} {
		for _, at := range []string{"main", "a unique root"} {
			t.Run(strings.TrimPrefix(symref, "refs/heads/")+" at "+at, func(t *testing.T) {
				t.Parallel()
				root, home := unclaimedBranchFixture(t, "b")
				target := "refs/heads/main"
				if at != "main" {
					target = addUnclaimedBranch(t, root, "a")
					commitOnBranch(t, root, target, "unique.txt", "unique\n")
				}
				plan, _, _ := runCleanup(t, root, home, "--discard-branch", "--unclaimed")
				fingerprint := cleanupRowFingerprint(t, plan)
				gitRun(t, root, "symbolic-ref", symref, target)
				identities := func() string {
					return gitOutput(t, root, "for-each-ref", "--format=%(refname) %(objectname) %(symref)")
				}
				refs := identities()
				plan, _, code := runCleanup(t, root, home, "--discard-branch", "--unclaimed")
				row := unclaimedVerdicts(t, plan)[symref]
				if code != 1 || strings.Contains(plan, "--apply") || setFingerprint.MatchString(plan) || !strings.HasPrefix(row, string(ActionError)+" "+symref) || !strings.Contains(row, "symref") {
					t.Fatalf("plan exit=%d row=%q stdout=%q, want exit 1, an error row naming the symref, no fingerprint, and no apply action", code, row, plan)
				}
				for _, apply := range [][]string{{"--apply", fingerprint}, {"--apply-current"}} {
					output, _, code := runCleanup(t, root, home, append([]string{"--discard-branch", "--unclaimed"}, apply...)...)
					if code != 1 || strings.Contains(output, ",removed,") || identities() != refs {
						t.Fatalf("%s exit=%d stdout=%q, want a refusal that changes no ref", apply[0], code, output)
					}
				}
			})
		}
	}
}
