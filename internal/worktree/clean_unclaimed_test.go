package worktree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/usage"
)

func TestPlanUnclaimedAssignmentSetExcludesClaimedCheckedOutAndForeignRefs(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	orphan := intent.AssignmentBranchRef(strings.Repeat("a", 32), strings.Repeat("b", 32))
	checkedOut := intent.AssignmentBranchRef(strings.Repeat("c", 32), strings.Repeat("d", 32))
	for index, branch := range []string{strings.TrimPrefix(orphan, "refs/heads/"), strings.TrimPrefix(checkedOut, "refs/heads/")} {
		gitRun(t, root, "checkout", "-qb", branch)
		commitInWorktree(t, root, "assignment-"+string(rune('a'+index))+".txt", "x\n", branch)
		gitRun(t, root, "checkout", "-q", "main")
	}
	checkedPath := filepath.Join(t.TempDir(), "checked-out")
	gitRun(t, root, "worktree", "add", "-q", checkedPath, strings.TrimPrefix(checkedOut, "refs/heads/"))
	gitRun(t, root, "branch", "bench/foreign/kept")
	claimed := mustCreate(t, root, home, "claimed-assignment-ref", "claimed")

	set, err := planUnclaimedAssignmentSet(root, unclaimedOptions())
	if err != nil || len(set.rows) != 1 || set.rows[0].ref != orphan {
		t.Fatalf("plan = %#v, err=%v; want only %q", set, err, orphan)
	}
	if set.fingerprint == "" || claimed.Assignment.Branch == orphan {
		t.Fatalf("plan fingerprint=%q claimed=%q", set.fingerprint, claimed.Assignment.Branch)
	}
}

// TestPlanUnclaimedShiftResidueBranch proves the shift namespace joins the unclaimed set
// under the same protections: a checked-out shift branch and a foreign ref stay out, and
// the plan holds one row per namespace (RI1). The landed shift residue row carries the shift
// reason and its removal, and the unique assignment row survives the apply.
func TestPlanUnclaimedShiftResidueBranch(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	residue := strings.TrimPrefix(intent.ShiftBranchPrefix(), "refs/heads/") + "20260101-000000"
	checkedOut := strings.TrimPrefix(intent.ShiftBranchPrefix(), "refs/heads/") + "20260101-000001"
	gitRun(t, root, "branch", residue)
	gitRun(t, root, "checkout", "-qb", checkedOut)
	commitInWorktree(t, root, "shift.txt", "x\n", checkedOut)
	gitRun(t, root, "checkout", "-q", "main")
	gitRun(t, root, "worktree", "add", "-q", filepath.Join(t.TempDir(), "checked-out"), checkedOut)
	gitRun(t, root, "branch", "archive/x")
	unique := addUnclaimedBranch(t, root, "a")
	commitOnBranch(t, root, unique, "unique.txt", "unique\n")

	set, err := planUnclaimedAssignmentSet(root, unclaimedOptions())
	if err != nil || len(set.rows) != 2 || set.rows[0].ref != unique || set.rows[1].ref != "refs/heads/"+residue {
		t.Fatalf("plan = %#v, err=%v; want %q and %q", set, err, unique, residue)
	}
	if set.rows[1].reason != "shift residue branch" {
		t.Fatalf("plan reason = %q, want shift residue branch", set.rows[1].reason)
	}

	plan, _, code := runCleanup(t, root, home, "--discard-branch", "--unclaimed")
	if code != 0 || !strings.Contains(plan, "class=landed; shift residue branch") {
		t.Fatalf("plan exit=%d stdout=%q, want the landed shift reason", code, plan)
	}
	output, _, code := runCleanup(t, root, home, "--discard-branch", "--unclaimed", "--apply", cleanupRowFingerprint(t, plan))
	if code != 0 || !strings.Contains(output, "refs/heads/"+residue+",removed,") || strings.Contains(output, unique+",removed,") {
		t.Fatalf("apply exit=%d stdout=%q, want the residue removed and the unique ref kept", code, output)
	}
	if git.OK("-C", root, "show-ref", "--verify", "--quiet", "refs/heads/"+residue) {
		t.Fatalf("apply retained %q", residue)
	}
	for _, ref := range []string{unique, "refs/heads/" + checkedOut, "refs/heads/archive/x"} {
		if !git.OK("-C", root, "show-ref", "--verify", "--quiet", ref) {
			t.Fatalf("apply deleted kept ref %q", ref)
		}
	}
}

func TestPlanUnclaimedAssignmentSetExcludesDefaultBranchInAssignmentNamespace(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	defaultBranch := "bench/assign/default/branch"
	gitRun(t, root, "branch", "-m", defaultBranch)
	gitRun(t, root, "checkout", "--detach", "-q")

	set, err := planUnclaimedAssignmentSet(root, unclaimedOptions())
	if err != nil || len(set.rows) != 0 {
		t.Fatalf("plan = %#v, err=%v; want configured default branch excluded", set, err)
	}
}

func TestCleanUnclaimedApplyCurrent(t *testing.T) {
	t.Parallel()
	t.Run("renders plan before removal", func(t *testing.T) {
		root, home := unclaimedBranchFixture(t)
		branch := addUnclaimedBranch(t, root, "c")
		unique := addUnclaimedBranch(t, root, "e")
		commitOnBranch(t, root, unique, "unique.txt", "unique\n")

		output, stderr, code := runCleanup(t, root, home, "--discard-branch", "--unclaimed", "--apply-current")
		planned := strings.Index(output, branch+",discard-remove,")
		removed := strings.Index(output, branch+",removed,")
		if code != 0 || stderr != "" || planned < 0 || removed < 0 || planned >= removed {
			t.Fatalf("cleanup = exit %d stdout=%q stderr=%q, want plan before removal", code, output, stderr)
		}
		if git.OK("-C", root, "show-ref", "--verify", "--quiet", branch) {
			t.Fatalf("apply-current retained %q", branch)
		}
		if !strings.Contains(output, unique+",retain,") || !git.OK("-C", root, "show-ref", "--verify", "--quiet", unique) {
			t.Fatalf("apply-current did not retain unique %q: %q", unique, output)
		}
	})

	for _, tc := range []struct {
		name string
		args []string
	}{
		{name: "without discard branch", args: []string{"--unclaimed", "--apply-current"}},
		{name: "without unclaimed", args: []string{"--discard-branch", "--apply-current", "."}},
		{name: "with fingerprint", args: []string{"--discard-branch", "--unclaimed", "--apply-current", "--apply", strings.Repeat("a", 64)}},
		{name: "with discard ignored", args: []string{"--discard-branch", "--unclaimed", "--apply-current", "--discard-ignored"}},
		{name: "with full", args: []string{"--discard-branch", "--unclaimed", "--apply-current", "--full"}},
		{name: "with landed", args: []string{"--discard-branch", "--unclaimed", "--apply-current", "--landed"}},
		{name: "with path", args: []string{"--discard-branch", "--unclaimed", "--apply-current", "."}},
		{name: "duplicate", args: []string{"--discard-branch", "--unclaimed", "--apply-current", "--apply-current"}},
	} {
		t.Run("refuses "+tc.name, func(t *testing.T) {
			root := newWorktreeRepo(t)
			home := filepath.Join(root, ".bench-home")
			branch := intent.AssignmentBranchRef(strings.Repeat("5", 32), strings.Repeat("6", 32))
			gitRun(t, root, "branch", strings.TrimPrefix(branch, "refs/heads/"))

			stdout, stderr, code := runCleanup(t, root, home, tc.args...)
			if code != 2 || stderr != "" || !strings.Contains(stdout, usage.WorktreeClean) {
				t.Fatalf("cleanup = exit %d stdout=%q stderr=%q, want usage refusal", code, stdout, stderr)
			}
			if !git.OK("-C", root, "show-ref", "--verify", "--quiet", branch) {
				t.Fatalf("invalid invocation deleted %q", branch)
			}
		})
	}
}

func TestCleanUnclaimedDiscardBranchRefusesMovedPlan(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	branch := intent.AssignmentBranchRef(strings.Repeat("e", 32), strings.Repeat("f", 32))
	short := strings.TrimPrefix(branch, "refs/heads/")
	gitRun(t, root, "checkout", "-qb", short)
	commitInWorktree(t, root, "first.txt", "first\n", "first")
	gitRun(t, root, "checkout", "-q", "main")
	plan, _, code := runCleanup(t, root, home, "--discard-branch", "--unclaimed")
	if code != 0 {
		t.Fatalf("plan exit=%d stdout=%q", code, plan)
	}
	gitRun(t, root, "checkout", "-q", short)
	commitInWorktree(t, root, "second.txt", "second\n", "second")
	gitRun(t, root, "checkout", "-q", "main")
	output, _, code := runCleanup(t, root, home, "--discard-branch", "--unclaimed", "--apply", cleanupRowFingerprint(t, plan))
	if code != 1 || !strings.Contains(output, errStaleFingerprint.Error()) {
		t.Fatalf("apply exit=%d stdout=%q, want stale refusal", code, output)
	}
	if !git.OK("-C", root, "show-ref", "--verify", "--quiet", branch) {
		t.Fatalf("stale apply deleted %q", branch)
	}
}

func TestCleanUnclaimedDiscardBranchRefusesChangedSet(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	first := intent.AssignmentBranchRef(strings.Repeat("1", 32), strings.Repeat("2", 32))
	second := intent.AssignmentBranchRef(strings.Repeat("3", 32), strings.Repeat("4", 32))
	gitRun(t, root, "checkout", "-qb", strings.TrimPrefix(first, "refs/heads/"))
	commitInWorktree(t, root, "first-set.txt", "first\n", "first set member")
	gitRun(t, root, "checkout", "-q", "main")
	plan, _, code := runCleanup(t, root, home, "--discard-branch", "--unclaimed")
	if code != 0 {
		t.Fatalf("plan exit=%d stdout=%q", code, plan)
	}
	gitRun(t, root, "checkout", "-qb", strings.TrimPrefix(second, "refs/heads/"))
	commitInWorktree(t, root, "second-set.txt", "second\n", "second set member")
	gitRun(t, root, "checkout", "-q", "main")

	output, _, code := runCleanup(t, root, home, "--discard-branch", "--unclaimed", "--apply", cleanupRowFingerprint(t, plan))
	if code != 1 || !strings.Contains(output, errStaleFingerprint.Error()) {
		t.Fatalf("apply exit=%d stdout=%q, want changed-set refusal", code, output)
	}
	for _, ref := range []string{first, second} {
		if !git.OK("-C", root, "show-ref", "--verify", "--quiet", ref) {
			t.Fatalf("changed-set apply deleted %q", ref)
		}
	}
}

func TestCleanUnclaimedRequiresExactAuthorizationGrammar(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		args []string
	}{
		{name: "missing discard branch", args: []string{"--unclaimed"}},
		{name: "discard ignored", args: []string{"--discard-branch", "--unclaimed", "--discard-ignored"}},
		{name: "full", args: []string{"--discard-branch", "--unclaimed", "--full"}},
		{name: "landed", args: []string{"--discard-branch", "--unclaimed", "--landed"}},
		{name: "path", args: []string{"--discard-branch", "--unclaimed", "."}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := newWorktreeRepo(t)
			home := filepath.Join(root, ".bench-home")
			ref := intent.AssignmentBranchRef(strings.Repeat("5", 32), strings.Repeat("6", 32))
			gitRun(t, root, "branch", strings.TrimPrefix(ref, "refs/heads/"))

			stdout, stderr, code := runCleanup(t, root, home, tc.args...)
			if code != 2 || stderr != "" || !strings.Contains(stdout, usage.WorktreeClean) {
				t.Fatalf("cleanup = exit %d stdout=%q stderr=%q, want usage refusal", code, stdout, stderr)
			}
			if !git.OK("-C", root, "show-ref", "--verify", "--quiet", ref) {
				t.Fatalf("invalid invocation deleted %q", ref)
			}
		})
	}
}

// reachableFromAHead reports whether tip stays reachable from main or from any surviving
// branch, the promise every bulk deletion keeps.
func reachableFromAHead(t *testing.T, root, tip string) bool {
	t.Helper()
	for _, head := range strings.Split(gitOutput(t, root, "for-each-ref", "--format=%(refname)", "refs/heads/"), "\n") {
		if git.OK("-C", root, "merge-base", "--is-ancestor", tip, head) {
			return true
		}
	}
	return false
}

// TestCleanUnclaimedBulkSweepKeepsUniqueRefs is RI18, RI19, and RI21: both apply forms remove
// the landed and subsumed rows at their tips, keep every unique root, and leave each deleted
// tip reachable from main or from a surviving ref.
func TestCleanUnclaimedBulkSweepKeepsUniqueRefs(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"--apply", "--apply-current"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			root, home := unclaimedBranchFixture(t)
			commitInWorktree(t, root, "main.txt", "main\n", "advance main")
			pair := addUnclaimedBranch(t, root, "3")
			top := addUnclaimedBranch(t, root, "7")
			removed := []string{
				addUnclaimedBranch(t, root, "1"),
				addUnclaimedBranchAt(t, root, "2", "main~1"),
				addUnclaimedBranchAt(t, root, "4", commitOnBranch(t, root, pair, "pair.txt", "pair\n")),
				addUnclaimedBranchAt(t, root, "5", commitOnBranch(t, root, top, "a.txt", "a\n")),
				addUnclaimedBranchAt(t, root, "6", commitOnBranch(t, root, top, "b.txt", "b\n")),
			}
			commitOnBranch(t, root, top, "c.txt", "c\n")
			tips := make([]string, len(removed))
			for i, ref := range removed {
				tips[i] = gitOutput(t, root, "rev-parse", ref)
			}
			args := []string{"--discard-branch", "--unclaimed", mode}
			if mode == "--apply" {
				plan, _, _ := runCleanup(t, root, home, "--discard-branch", "--unclaimed")
				args = append(args, cleanupRowFingerprint(t, plan))
			}
			output, stderr, code := runCleanup(t, root, home, args...)
			if code != 0 || stderr != "" || strings.Count(output, ",removed,") != len(removed) {
				t.Fatalf("apply exit=%d stderr=%q stdout=%q, want %d removals", code, stderr, output, len(removed))
			}
			for i, ref := range removed {
				if git.OK("-C", root, "show-ref", "--verify", "--quiet", ref) || !reachableFromAHead(t, root, tips[i]) {
					t.Fatalf("apply left %q or lost its tip %s", ref, tips[i])
				}
			}
			for _, ref := range []string{pair, top} {
				if !git.OK("-C", root, "show-ref", "--verify", "--quiet", ref) {
					t.Fatalf("apply deleted unique ref %q", ref)
				}
			}
		})
	}
}

// TestCleanUnclaimedPlanNamesApplyOnlyWhenARowRemoves is RI22 and RI66, and it pins the
// retained detail of a unique row (RI10) in this one test file.
func TestCleanUnclaimedPlanNamesApplyOnlyWhenARowRemoves(t *testing.T) {
	t.Parallel()
	root, home := unclaimedBranchFixture(t)
	unique := addUnclaimedBranch(t, root, "a")
	commitOnBranch(t, root, unique, "unique.txt", "unique\n")
	plan, stderr, code := runCleanup(t, root, home, "--discard-branch", "--unclaimed")
	rows := cleanupRows(plan)
	if code != 0 || stderr != "" || strings.Contains(plan, "--apply") || len(rows) != 1 ||
		!strings.HasSuffix(cleanupRowValue(cleanupRowFields(rows[0])[6]), "retained: content main lacks") {
		t.Fatalf("unique-only plan exit=%d stderr=%q stdout=%q, want a retained row and no apply action", code, stderr, plan)
	}
	addUnclaimedBranch(t, root, "b")
	plan, stderr, code = runCleanup(t, root, home, "--discard-branch", "--unclaimed")
	action := "bench worktree clean --discard-branch --unclaimed --apply " + cleanupRowFingerprint(t, plan)
	if code != 0 || stderr != "" || !strings.Contains(plan, action) {
		t.Fatalf("landed plan exit=%d stderr=%q stdout=%q, want %q", code, stderr, plan, action)
	}
}

// TestCleanUnclaimedErrorRowRefusesTheSet is RI11, RI72, and RI73. The loose ref file names a
// blob directly, because git update-ref can refuse a ref that names no commit.
func TestCleanUnclaimedErrorRowRefusesTheSet(t *testing.T) {
	t.Parallel()
	root, home := unclaimedBranchFixture(t)
	landed := addUnclaimedBranch(t, root, "a")
	broken := intent.AssignmentBranchRef(strings.Repeat("b", 32), strings.Repeat("f", 32))
	loose := filepath.Join(root, ".git", filepath.FromSlash(broken))
	mustNoError(t, os.MkdirAll(filepath.Dir(loose), 0o755))
	mustWrite(t, loose, []byte(gitOutput(t, root, "hash-object", "-w", "tracked.txt")+"\n"), 0o644)

	plan, _, code := runCleanup(t, root, home, "--discard-branch", "--unclaimed")
	if code != 1 || strings.Contains(plan, "--apply") || setFingerprint.MatchString(plan) {
		t.Errorf("plan exit=%d stdout=%q, want exit 1 with no fingerprint and no apply action", code, plan)
	}
	fields := map[string][]string{}
	for _, row := range cleanupRows(plan) {
		fields[cleanupRowValue(cleanupRowFields(row)[0])] = cleanupRowFields(row)
	}
	if row := fields[broken]; len(row) != len(cleanupFields) || row[1] != "error" || !strings.Contains(row[6], broken) || !strings.Contains(row[6], "blob") {
		t.Errorf("broken ref row = %q, want an error naming the ref and the blob", row)
	}
	output, _, code := runCleanup(t, root, home, "--discard-branch", "--unclaimed", "--apply-current")
	if code != 1 || strings.Contains(output, ",removed,") || !git.OK("-C", root, "show-ref", "--verify", "--quiet", landed) {
		t.Fatalf("apply-current exit=%d stdout=%q, want a refusal that keeps %q", code, output, landed)
	}
}

// TestCleanUnclaimedPlanIsReadOnlyAndQuietWhenEmpty is RI23 and RI12.
func TestCleanUnclaimedPlanIsReadOnlyAndQuietWhenEmpty(t *testing.T) {
	t.Parallel()
	root, home := unclaimedBranchFixture(t)
	empty, stderr, code := runCleanup(t, root, home, "--discard-branch", "--unclaimed")
	if code != 0 || stderr != "" || len(cleanupRows(empty)) != 0 || !strings.HasPrefix(empty, "worktree_cleanup[0]") {
		t.Fatalf("empty plan exit=%d stderr=%q stdout=%q, want the empty table", code, stderr, empty)
	}
	addUnclaimedBranch(t, root, "a")
	unique := addUnclaimedBranch(t, root, "b")
	addUnclaimedBranchAt(t, root, "c", commitOnBranch(t, root, unique, "unique.txt", "unique\n"))
	refs := gitOutput(t, root, "for-each-ref")
	first, _, firstCode := runCleanup(t, root, home, "--discard-branch", "--unclaimed")
	second, _, secondCode := runCleanup(t, root, home, "--discard-branch", "--unclaimed")
	if firstCode != 0 || secondCode != 0 || first != second || len(cleanupRows(first)) != 3 || gitOutput(t, root, "for-each-ref") != refs {
		t.Fatalf("plans = (%d, %q) then (%d, %q), want identical rows and unchanged refs", firstCode, first, secondCode, second)
	}
}

// TestCleanUnclaimedStaleClassRefusesTheOldPlan is RI17 and RI63: a moved subsumed ref and a
// class change with the tip unchanged both refuse the old fingerprint before any delete.
func TestCleanUnclaimedStaleClassRefusesTheOldPlan(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		mutate func(t *testing.T, root, home, unique, subsumed string)
	}{
		{"RI17 a new commit on a subsumed ref", func(t *testing.T, root, _, _, subsumed string) {
			commitOnBranch(t, root, subsumed, "moved.txt", "moved\n")
		}},
		{"RI63 a recorded active branch at a descendant of the unique ref", func(t *testing.T, root, home, unique, _ string) {
			active := mustCreate(t, root, home, "stale-class-active", "active")
			gitRun(t, active.Path, "reset", "-q", "--hard", gitOutput(t, root, "rev-parse", unique))
			commitInWorktree(t, active.Path, "descendant.txt", "descendant\n", "descendant")
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root, home := unclaimedBranchFixture(t)
			landed := addUnclaimedBranch(t, root, "a")
			unique := addUnclaimedBranch(t, root, "c")
			subsumed := addUnclaimedBranchAt(t, root, "d", commitOnBranch(t, root, unique, "unique.txt", "unique\n"))
			uniqueTip := gitOutput(t, root, "rev-parse", unique)
			plan, _, _ := runCleanup(t, root, home, "--discard-branch", "--unclaimed")
			tc.mutate(t, root, home, unique, subsumed)
			output, _, code := runCleanup(t, root, home, "--discard-branch", "--unclaimed", "--apply", cleanupRowFingerprint(t, plan))
			if code != 1 || !strings.Contains(output, errStaleFingerprint.Error()) || !strings.Contains(output, "bench worktree clean --discard-branch --unclaimed") {
				t.Fatalf("apply exit=%d stdout=%q, want the stale refusal and the re-plan command", code, output)
			}
			if !git.OK("-C", root, "show-ref", "--verify", "--quiet", landed) || gitOutput(t, root, "rev-parse", unique) != uniqueTip {
				t.Fatalf("stale apply changed %q or %q", landed, unique)
			}
		})
	}
}
