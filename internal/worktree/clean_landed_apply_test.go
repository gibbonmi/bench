package worktree

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

func TestCleanUnclaimedDiscardBranchPlansThenDeletesExactAssignmentRef(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	f := repoHome{root, filepath.Join(root, ".bench-home")}
	branch := intent.AssignmentBranchRef(strings.Repeat("a", 32), strings.Repeat("b", 32))
	short := strings.TrimPrefix(branch, "refs/heads/")
	gitRun(t, root, "checkout", "-qb", short)
	commitInWorktree(t, root, "orphan.txt", "orphan\n", "orphan assignment")
	gitRun(t, root, "checkout", "-q", "main")
	gitRun(t, root, "merge", "-q", "--ff-only", short)
	unique := addUnclaimedBranch(t, root, "c")
	commitOnBranch(t, root, unique, "unique.txt", "unique\n")

	plan := runVerb(t, verbClean, f.call("--discard-branch", "--unclaimed"))
	if plan.exit != 0 || plan.stderr != "" || !strings.Contains(plan.stdout, branch) {
		t.Fatalf("plan exit=%d stdout=%q stderr=%q, want unclaimed assignment branch plan", plan.exit, plan.stdout, plan.stderr)
	}
	if !git.OK("-C", root, "show-ref", "--verify", "--quiet", branch) {
		t.Fatalf("plan deleted %q", branch)
	}

	applied := runVerb(t, verbClean, f.call("--discard-branch", "--unclaimed", "--apply", plan.mustFingerprint(t)))
	if applied.exit != 0 || applied.stderr != "" || !strings.Contains(applied.stdout, branch) {
		t.Fatalf("apply exit=%d stdout=%q stderr=%q, want exact branch deletion", applied.exit, applied.stdout, applied.stderr)
	}
	if git.OK("-C", root, "show-ref", "--verify", "--quiet", branch) || !git.OK("-C", root, "show-ref", "--verify", "--quiet", unique) {
		t.Fatalf("apply kept landed %q or removed unique %q", branch, unique)
	}
}

func TestCleanLandedApplyRemovesAndSettles(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	f := landedSetFixture(t)
	plan := runVerb(t, verbClean, f.call("--landed"))
	if plan.exit != 0 || plan.stderr != "" {
		t.Fatalf("plan exit=%d stdout=%q stderr=%q", plan.exit, plan.stdout, plan.stderr)
	}

	applied := runVerb(t, verbClean, f.call("--landed", "--apply", plan.mustFingerprint(t)))
	if applied.exit != 0 || applied.stderr != "" {
		t.Fatalf("apply exit=%d stdout=%q stderr=%q", applied.exit, applied.stdout, applied.stderr)
	}
	for _, removed := range []Creation{f.first, f.second} {
		if _, err := os.Lstat(removed.Path); !os.IsNotExist(err) {
			t.Fatalf("removed worktree %q still present: %v", removed.Path, err)
		}
	}
	if _, err := os.Lstat(f.dirty.Path); err != nil {
		t.Fatalf("retained worktree %q disappeared: %v", f.dirty.Path, err)
	}
	for _, removed := range []Creation{f.first, f.second} {
		if git.OK("-C", f.root, "show-ref", "--verify", "--quiet", removed.Assignment.Branch) {
			t.Fatalf("removed worktree branch %q remains", removed.Assignment.Branch)
		}
	}
	assignments, err := intent.Assignments(f.root)
	if err != nil || len(assignments) != 1 || assignments[0].ID != f.dirty.Assignment.ID {
		t.Fatalf("assignments after apply = %#v, %v; want only dirty row", assignments, err)
	}
	list := descendant(t, binary, "worktree", "list")
	list.Dir = f.root
	listing, listErr := list.Output()
	if listErr != nil || strings.Contains(string(listing), f.first.Assignment.ID) || strings.Contains(string(listing), f.second.Assignment.ID) || !strings.Contains(string(listing), f.dirty.Assignment.ID) {
		t.Fatalf("fresh list after apply = (%v, %q), want only dirty assignment", listErr, listing)
	}
}

func TestCleanLandedApplyRefusesInitialDriftWithoutMutation(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		mutate func(*testing.T, string, string, Creation)
		args   func(string) []string
	}{
		{name: "new member", mutate: func(t *testing.T, root, home string, _ Creation) {
			created := mustCreate(t, root, home, "landed-new-member", "new member")
			landAssignment(t, root, created, "member.txt")
		}},
		{name: "tracked dirt", mutate: func(t *testing.T, _, _ string, first Creation) {
			mustWrite(t, filepath.Join(first.Path, "first.txt"), []byte("changed\n"), 0o644)
		}},
		{name: "live lease", mutate: func(t *testing.T, _, _ string, first Creation) {
			lease, err := LeaseFile(first.Path)
			mustNoError(t, err)
			mustWrite(t, lease, []byte(strconv.Itoa(os.Getpid())+" 2026-08-16T00:00:00Z\n"), 0o600)
		}},
		{name: "landed head advances", mutate: func(t *testing.T, root, _ string, first Creation) {
			commitInWorktree(t, first.Path, "advance.txt", "advance\n", "advance landed head")
			gitRun(t, root, "cherry-pick", strings.TrimPrefix(first.Assignment.Branch, "refs/heads/"))
		}},
		{name: "option set", mutate: func(*testing.T, string, string, Creation) {}, args: func(fingerprint string) []string {
			return []string{"--discard-ignored", "--landed", "--apply", fingerprint}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := landedSetFixture(t)
			plan := runVerb(t, verbClean, f.call("--landed"))
			if plan.exit != 0 || plan.stderr != "" {
				t.Fatalf("plan exit=%d stdout=%q stderr=%q", plan.exit, plan.stdout, plan.stderr)
			}
			tc.mutate(t, f.root, f.home, f.first)
			args := []string{"--landed", "--apply", plan.mustFingerprint(t)}
			if tc.args != nil {
				args = tc.args(plan.mustFingerprint(t))
			}
			applied := runVerb(t, verbClean, f.call(args...))
			if applied.exit != 1 || applied.stderr != "" || !strings.HasPrefix(applied.stdout, "worktree_cleanup[") || !strings.Contains(applied.stdout, "unknown,error,unknown,unknown,none,") || strings.Count(applied.stdout, errStaleFingerprint.Error()) != 1 {
				t.Fatalf("apply exit=%d stdout=%q stderr=%q, want stale refusal diagnostic", applied.exit, applied.stdout, applied.stderr)
			}
			for _, creation := range []Creation{f.first, f.second} {
				if _, err := os.Lstat(creation.Path); err != nil {
					t.Fatalf("stale apply removed %q: %v", creation.Path, err)
				}
			}
			assignments, assignmentErr := intent.Assignments(f.root)
			if assignmentErr != nil {
				t.Fatal(assignmentErr)
			}
			for _, creation := range []Creation{f.first, f.second} {
				active := false
				for _, assignment := range assignments {
					if assignment.ID == creation.Assignment.ID && assignment.State == intent.StateActive {
						active = true
						break
					}
				}
				if !active {
					t.Fatalf("stale apply settled %q: %#v", creation.Assignment.ID, assignments)
				}
			}
		})
	}
}

func TestCleanLandedApplyReplansEachRowBeforeMutation(t *testing.T) {
	t.Parallel()
	f := landedSetFixture(t)
	plan := runVerb(t, verbClean, f.call("--landed"))
	if plan.exit != 0 || plan.stderr != "" {
		t.Fatalf("plan exit=%d stdout=%q stderr=%q", plan.exit, plan.stdout, plan.stderr)
	}
	set, err := planLandedSet(defaultJoins(), f.root, CleanupOptions{}, "")
	if err != nil || len(set.rows) != 3 {
		t.Fatalf("landed set = %#v, %v; want three rows", set, err)
	}
	removable := make([]landedCleanupRow, 0, 2)
	for _, row := range set.rows {
		if row.plan.Action.Removes() {
			removable = append(removable, row)
		}
	}
	if len(removable) != 2 {
		t.Fatalf("removable rows = %#v, want two", removable)
	}
	var settled, drifted Creation
	for _, creation := range []Creation{f.first, f.second} {
		if creation.Assignment.ID == removable[0].assignment.ID {
			settled = creation
		}
		if creation.Assignment.ID == removable[1].assignment.ID {
			drifted = creation
		}
	}
	if settled.Path == "" || drifted.Path == "" {
		t.Fatalf("set order did not name two removable rows: %#v", set.rows)
	}
	j := defaultJoins()
	mutated := false
	j.cleanupBoundary = func(step LifecycleStep) error {
		if step == StepTerminalReceipt && !mutated {
			mutated = true
			mustWrite(t, filepath.Join(drifted.Path, "drifted.txt"), []byte("drifted\n"), 0o644)
		}
		return nil
	}
	applied := runVerb(t, verbClean, f.callWith(j, "--landed", "--apply", plan.mustFingerprint(t)))
	if applied.exit != 1 || applied.stderr != "" || !mutated {
		t.Fatalf("apply exit=%d stderr=%q mutated=%t, want per-row stale refusal", applied.exit, applied.stderr, mutated)
	}
	if _, err := os.Lstat(settled.Path); !os.IsNotExist(err) {
		t.Fatalf("first row was not settled before second drift: %v", err)
	}
	if _, err := os.Lstat(drifted.Path); err != nil {
		t.Fatalf("drifted second row disappeared: %v", err)
	}
	assignments, err := intent.Assignments(f.root)
	if err != nil || len(assignments) != 2 {
		t.Fatalf("assignments after second-row drift = %#v, %v; want the dirty and drifted rows", assignments, err)
	}
	for _, assignment := range assignments {
		if assignment.ID == settled.Assignment.ID {
			t.Fatalf("settled first row %q remains assigned", settled.Assignment.ID)
		}
	}
	if retry := runVerb(t, verbClean, f.callWith(j, "--landed", "--apply", plan.mustFingerprint(t))); retry.exit != 1 {
		t.Fatalf("spent fingerprint replay exit=%d, want 1", retry.exit)
	}
	fresh := runVerb(t, verbClean, f.call("--landed"))
	if fresh.exit != 0 || fresh.stderr != "" || strings.Contains(fresh.stdout, settled.Assignment.ID) || !strings.Contains(fresh.stdout, drifted.Assignment.ID) {
		t.Fatalf("fresh plan = (%d, %q, %q), want only second row", fresh.exit, fresh.stdout, fresh.stderr)
	}
}

func TestCleanLandedApplyStopsAfterCompletedRowFault(t *testing.T) {
	t.Parallel()
	f := landedSetFixture(t)
	plan := runVerb(t, verbClean, f.call("--landed"))
	if plan.exit != 0 || plan.stderr != "" {
		t.Fatalf("plan exit=%d stdout=%q stderr=%q", plan.exit, plan.stdout, plan.stderr)
	}
	j := defaultJoins()
	locks := 0
	j.cleanupBoundary = func(step LifecycleStep) error {
		if step == StepApplyLocked {
			locks++
			if locks == 2 {
				return errors.New("stop before second landed row")
			}
		}
		return nil
	}
	applied := runVerb(t, verbClean, f.callWith(j, "--landed", "--apply", plan.mustFingerprint(t)))
	if applied.exit != 1 || applied.stderr != "" || locks != 2 {
		t.Fatalf("apply exit=%d stderr=%q locks=%d, want second-row fault", applied.exit, applied.stderr, locks)
	}
	removed, present := 0, 0
	for _, creation := range []Creation{f.first, f.second} {
		if _, err := os.Lstat(creation.Path); os.IsNotExist(err) {
			removed++
		} else if err == nil {
			present++
		} else {
			t.Fatalf("inspect %q: %v", creation.Path, err)
		}
	}
	if removed != 1 || present != 1 {
		t.Fatalf("fault outcome removed=%d present=%d, want one completed and one untouched", removed, present)
	}
	assignments, err := intent.Assignments(f.root)
	if err != nil || len(assignments) != 2 {
		t.Fatalf("assignments after fault = %#v, %v; want retained dirty row and one unstarted removable row", assignments, err)
	}
	if retry := runVerb(t, verbClean, f.call("--landed", "--apply", plan.mustFingerprint(t))); retry.exit != 1 {
		t.Fatalf("interrupted fingerprint replay exit=%d, want 1", retry.exit)
	}
	fresh := runVerb(t, verbClean, f.call("--landed"))
	if fresh.exit != 0 || fresh.stderr != "" || strings.Count(fresh.stdout, ",remove,") != 1 {
		t.Fatalf("fresh plan = (%d, %q, %q), want one remaining removable row", fresh.exit, fresh.stdout, fresh.stderr)
	}
}

func TestCleanLandedPlanRetainsPreservedRowThroughEligibilityOwner(t *testing.T) {
	t.Parallel()
	f := landedSetFixture(t)
	plan := runVerb(t, verbClean, f.call("--landed"))
	if plan.exit != 0 || plan.stderr != "" {
		t.Fatalf("plan exit=%d stdout=%q stderr=%q", plan.exit, plan.stdout, plan.stderr)
	}
	if !strings.Contains(plan.stdout, f.dirty.Path+",retain,") ||
		!strings.Contains(plan.stdout, ",none,") ||
		strings.Count(plan.stdout, "per-path cleanup is required to preserve work") != 1 {
		t.Fatalf("plan = %q, want dirty row retained once via the shared preservation refusal", plan.stdout)
	}
}

func TestCleanLandedApplyCarriesModifiersAndDeletesProvenBranches(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	f := repoHome{root, filepath.Join(root, ".bench-home")}
	mustWrite(t, filepath.Join(root, ".gitignore"), []byte("ignored.txt\n"), 0o644)
	gitRun(t, root, "add", ".gitignore")
	gitRun(t, root, "commit", "-qm", "ignore landed residue")
	clean := mustCreate(t, root, f.home, "landed-modifier-clean", "clean")
	ignored := mustCreate(t, root, f.home, "landed-modifier-ignored", "ignored")
	landAssignment(t, root, clean, "clean.txt")
	landAssignment(t, root, ignored, "landed.txt")
	mustWrite(t, filepath.Join(ignored.Path, "ignored.txt"), []byte("residue\n"), 0o644)

	bare := runVerb(t, verbClean, f.call("--landed"))
	if bare.exit != 0 || bare.stderr != "" || !strings.Contains(bare.stdout, ignored.Path+",retain,") || !strings.Contains(bare.stdout, "ignored residuals require --discard-ignored") {
		t.Fatalf("bare plan = (%d, %q, %q), want ignored retain", bare.exit, bare.stdout, bare.stderr)
	}
	widened := runVerb(t, verbClean, f.call("--discard-ignored", "--full", "--landed"))
	if widened.exit != 0 || widened.stderr != "" || !strings.Contains(widened.stdout, ignored.Path+",discard-remove,") || !strings.Contains(widened.stdout, "ignored_paths[1]") {
		t.Fatalf("widened plan = (%d, %q, %q), want discard removal and preview", widened.exit, widened.stdout, widened.stderr)
	}
	applied := runVerb(t, verbClean, f.call("--discard-ignored", "--full", "--landed", "--apply", widened.mustFingerprint(t)))
	if applied.exit != 0 || applied.stderr != "" || strings.Count(applied.stdout, ",removed,") != 2 {
		t.Fatalf("modifier apply = (%d, %q, %q), want two removals", applied.exit, applied.stdout, applied.stderr)
	}
	for _, creation := range []Creation{clean, ignored} {
		if _, err := os.Lstat(creation.Path); !os.IsNotExist(err) {
			t.Fatalf("modifier apply left tree %q: %v", creation.Path, err)
		}
		if git.OK("-C", root, "show-ref", "--verify", "--quiet", creation.Assignment.Branch) {
			t.Fatalf("modifier apply left branch %q", creation.Assignment.Branch)
		}
	}
}

func TestCleanLandedDiscardBranchOnlyChangesDetail(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	f := repoHome{root, filepath.Join(root, ".bench-home")}
	creation := mustCreate(t, root, f.home, "landed-branch-assertion", "branch assertion")
	landAssignment(t, root, creation, "branch.txt")
	plan := runVerb(t, verbClean, f.call("--discard-branch", "--landed"))
	if plan.exit != 0 || plan.stderr != "" || !strings.Contains(plan.stdout, "discards branch "+creation.Assignment.Branch) {
		t.Fatalf("assertion plan = (%d, %q, %q), want asserted branch detail", plan.exit, plan.stdout, plan.stderr)
	}
	applied := runVerb(t, verbClean, f.call("--discard-branch", "--landed", "--apply", plan.mustFingerprint(t)))
	if applied.exit != 0 || applied.stderr != "" {
		t.Fatalf("assertion apply = (%d, %q), want success", applied.exit, applied.stderr)
	}
	if _, err := os.Lstat(creation.Path); !os.IsNotExist(err) || git.OK("-C", root, "show-ref", "--verify", "--quiet", creation.Assignment.Branch) {
		t.Fatalf("assertion apply tree/branch outcome = %v/%t, want both absent", err, git.OK("-C", root, "show-ref", "--verify", "--quiet", creation.Assignment.Branch))
	}
}

func TestCleanLandedRetainsUnparseableLeaseAndSkipsUnprovableBranch(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	f := repoHome{root, filepath.Join(root, ".bench-home")}
	unknown := mustCreate(t, root, f.home, "landed-unknown-lease", "unknown lease")
	unprovable := mustCreate(t, root, f.home, "landed-unprovable", "unprovable")
	landAssignment(t, root, unknown, "unknown.txt")
	commitInWorktree(t, unprovable.Path, "unprovable.txt", "unmerged\n", "unprovable")
	lease, err := LeaseFile(unknown.Path)
	mustNoError(t, err)
	mustWrite(t, lease, []byte("not-a-lease\n"), 0o600)
	plan := runVerb(t, verbClean, f.call("--landed"))
	if plan.exit != 0 || plan.stderr != "" || !strings.Contains(plan.stdout, unknown.Path+",retain,") || !strings.Contains(plan.stdout, "uncertain") || strings.Contains(plan.stdout, unprovable.Path) {
		t.Fatalf("plan = (%d, %q, %q), want unknown retain and no unprovable row", plan.exit, plan.stdout, plan.stderr)
	}
	if summary := runVerb(t, verbResumeClean, f.call()); summary.exit != 0 || !strings.Contains(summary.stdout, "landed=1") {
		t.Fatalf("resume = (%d, %q, %q), want unknown lease counted landed", summary.exit, summary.stdout, summary.stderr)
	}
	applied := runVerb(t, verbClean, f.call("--landed", "--apply", plan.mustFingerprint(t)))
	if applied.exit != 0 || applied.stderr != "" {
		t.Fatalf("unknown lease apply = (%d, %q), want no-op success", applied.exit, applied.stderr)
	}
	if _, err := os.Lstat(unknown.Path); err != nil {
		t.Fatalf("unknown lease path disappeared: %v", err)
	}
}
