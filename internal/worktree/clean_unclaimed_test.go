package worktree

import (
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

	set, err := planUnclaimedAssignmentSet(root, CleanupOptions{DiscardBranch: true, Unclaimed: true})
	if err != nil || len(set.rows) != 1 || set.rows[0].ref != orphan {
		t.Fatalf("plan = %#v, err=%v; want only %q", set, err, orphan)
	}
	if set.fingerprint == "" || claimed.Assignment.Branch == orphan {
		t.Fatalf("plan fingerprint=%q claimed=%q", set.fingerprint, claimed.Assignment.Branch)
	}
}

func TestPlanUnclaimedAssignmentSetExcludesDefaultBranchInAssignmentNamespace(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	defaultBranch := "bench/assign/default/branch"
	gitRun(t, root, "branch", "-m", defaultBranch)
	gitRun(t, root, "checkout", "--detach", "-q")

	set, err := planUnclaimedAssignmentSet(root, CleanupOptions{DiscardBranch: true, Unclaimed: true})
	if err != nil || len(set.rows) != 0 {
		t.Fatalf("plan = %#v, err=%v; want configured default branch excluded", set, err)
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
