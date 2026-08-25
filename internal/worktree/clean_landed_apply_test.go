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

func TestCleanLandedApplyRemovesAndSettles(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	root, home, first, second, dirty := landedSetFixture(t)
	plan, planErr, planCode := runCleanLanded(t, root, home, "--landed")
	if planCode != 0 || planErr != "" {
		t.Fatalf("plan exit=%d stdout=%q stderr=%q", planCode, plan, planErr)
	}

	output, stderr, code := runCleanLanded(t, root, home, "--landed", "--apply", landedRowFingerprint(t, plan))
	if code != 0 || stderr != "" {
		t.Fatalf("apply exit=%d stdout=%q stderr=%q", code, output, stderr)
	}
	for _, removed := range []Creation{first, second} {
		if _, err := os.Lstat(removed.Path); !os.IsNotExist(err) {
			t.Fatalf("removed worktree %q still present: %v", removed.Path, err)
		}
	}
	if _, err := os.Lstat(dirty.Path); err != nil {
		t.Fatalf("retained worktree %q disappeared: %v", dirty.Path, err)
	}
	for _, removed := range []Creation{first, second} {
		if git.OK("-C", root, "show-ref", "--verify", "--quiet", removed.Assignment.Branch) {
			t.Fatalf("removed worktree branch %q remains", removed.Assignment.Branch)
		}
	}
	assignments, err := intent.Assignments(root)
	if err != nil || len(assignments) != 1 || assignments[0].ID != dirty.Assignment.ID {
		t.Fatalf("assignments after apply = %#v, %v; want only dirty row", assignments, err)
	}
	list := descendant(t, binary, "worktree", "list")
	list.Dir = root
	listing, listErr := list.Output()
	if listErr != nil || strings.Contains(string(listing), first.Assignment.ID) || strings.Contains(string(listing), second.Assignment.ID) || !strings.Contains(string(listing), dirty.Assignment.ID) {
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
			root, home, first, second, _ := landedSetFixture(t)
			plan, planErr, planCode := runCleanLanded(t, root, home, "--landed")
			if planCode != 0 || planErr != "" {
				t.Fatalf("plan exit=%d stdout=%q stderr=%q", planCode, plan, planErr)
			}
			tc.mutate(t, root, home, first)
			args := []string{"--landed", "--apply", landedRowFingerprint(t, plan)}
			if tc.args != nil {
				args = tc.args(landedRowFingerprint(t, plan))
			}
			stdout, stderr, code := runCleanLanded(t, root, home, args...)
			if code != 1 || stderr != "" || !strings.HasPrefix(stdout, "worktree_cleanup[") || !strings.Contains(stdout, "unknown,error,unknown,unknown,none,") || strings.Count(stdout, errStaleFingerprint.Error()) != 1 {
				t.Fatalf("apply exit=%d stdout=%q stderr=%q, want stale refusal diagnostic", code, stdout, stderr)
			}
			for _, creation := range []Creation{first, second} {
				if _, err := os.Lstat(creation.Path); err != nil {
					t.Fatalf("stale apply removed %q: %v", creation.Path, err)
				}
			}
			assignments, assignmentErr := intent.Assignments(root)
			if assignmentErr != nil {
				t.Fatal(assignmentErr)
			}
			for _, creation := range []Creation{first, second} {
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
	root, home, first, second, _ := landedSetFixture(t)
	plan, planErr, planCode := runCleanLanded(t, root, home, "--landed")
	if planCode != 0 || planErr != "" {
		t.Fatalf("plan exit=%d stdout=%q stderr=%q", planCode, plan, planErr)
	}
	set, err := planLandedSet(root, CleanupOptions{})
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
	for _, creation := range []Creation{first, second} {
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
	previousBoundary := cleanupTransactionBoundary
	t.Cleanup(func() { cleanupTransactionBoundary = previousBoundary })
	mutated := false
	cleanupTransactionBoundary = func(step LifecycleStep) error {
		if step == StepTerminalReceipt && !mutated {
			mutated = true
			mustWrite(t, filepath.Join(drifted.Path, "drifted.txt"), []byte("drifted\n"), 0o644)
		}
		return nil
	}
	_, stderr, code := runCleanLanded(t, root, home, "--landed", "--apply", landedRowFingerprint(t, plan))
	if code != 1 || stderr != "" || !mutated {
		t.Fatalf("apply exit=%d stderr=%q mutated=%t, want per-row stale refusal", code, stderr, mutated)
	}
	if _, err := os.Lstat(settled.Path); !os.IsNotExist(err) {
		t.Fatalf("first row was not settled before second drift: %v", err)
	}
	if _, err := os.Lstat(drifted.Path); err != nil {
		t.Fatalf("drifted second row disappeared: %v", err)
	}
	assignments, err := intent.Assignments(root)
	if err != nil || len(assignments) != 2 {
		t.Fatalf("assignments after second-row drift = %#v, %v; want the dirty and drifted rows", assignments, err)
	}
	for _, assignment := range assignments {
		if assignment.ID == settled.Assignment.ID {
			t.Fatalf("settled first row %q remains assigned", settled.Assignment.ID)
		}
	}
	_, _, retryCode := runCleanLanded(t, root, home, "--landed", "--apply", landedRowFingerprint(t, plan))
	if retryCode != 1 {
		t.Fatalf("spent fingerprint replay exit=%d, want 1", retryCode)
	}
	fresh, freshErr, freshCode := runCleanLanded(t, root, home, "--landed")
	if freshCode != 0 || freshErr != "" || strings.Contains(fresh, settled.Assignment.ID) || !strings.Contains(fresh, drifted.Assignment.ID) {
		t.Fatalf("fresh plan = (%d, %q, %q), want only second row", freshCode, fresh, freshErr)
	}
}

func TestCleanLandedApplyStopsAfterCompletedRowFault(t *testing.T) {
	root, home, first, second, _ := landedSetFixture(t)
	plan, planErr, planCode := runCleanLanded(t, root, home, "--landed")
	if planCode != 0 || planErr != "" {
		t.Fatalf("plan exit=%d stdout=%q stderr=%q", planCode, plan, planErr)
	}
	previousBoundary := cleanupTransactionBoundary
	t.Cleanup(func() { cleanupTransactionBoundary = previousBoundary })
	locks := 0
	cleanupTransactionBoundary = func(step LifecycleStep) error {
		if step == StepApplyLocked {
			locks++
			if locks == 2 {
				return errors.New("stop before second landed row")
			}
		}
		return nil
	}
	_, stderr, code := runCleanLanded(t, root, home, "--landed", "--apply", landedRowFingerprint(t, plan))
	if code != 1 || stderr != "" || locks != 2 {
		t.Fatalf("apply exit=%d stderr=%q locks=%d, want second-row fault", code, stderr, locks)
	}
	removed, present := 0, 0
	for _, creation := range []Creation{first, second} {
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
	assignments, err := intent.Assignments(root)
	if err != nil || len(assignments) != 2 {
		t.Fatalf("assignments after fault = %#v, %v; want retained dirty row and one unstarted removable row", assignments, err)
	}
	_, _, retryCode := runCleanLanded(t, root, home, "--landed", "--apply", landedRowFingerprint(t, plan))
	if retryCode != 1 {
		t.Fatalf("interrupted fingerprint replay exit=%d, want 1", retryCode)
	}
	fresh, freshErr, freshCode := runCleanLanded(t, root, home, "--landed")
	if freshCode != 0 || freshErr != "" || strings.Count(fresh, ",remove,") != 1 {
		t.Fatalf("fresh plan = (%d, %q, %q), want one remaining removable row", freshCode, fresh, freshErr)
	}
}

func TestCleanLandedPlanRetainsPreservedRowThroughEligibilityOwner(t *testing.T) {
	t.Parallel()
	root, home, _, _, dirty := landedSetFixture(t)
	plan, planErr, planCode := runCleanLanded(t, root, home, "--landed")
	if planCode != 0 || planErr != "" {
		t.Fatalf("plan exit=%d stdout=%q stderr=%q", planCode, plan, planErr)
	}
	if !strings.Contains(plan, dirty.Path+",retain,") ||
		!strings.Contains(plan, ",none,") ||
		strings.Count(plan, "per-path cleanup is required to preserve work") != 1 {
		t.Fatalf("plan = %q, want dirty row retained once via the shared preservation refusal", plan)
	}
}

func TestCleanLandedApplyCarriesModifiersAndDeletesProvenBranches(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	mustWrite(t, filepath.Join(root, ".gitignore"), []byte("ignored.txt\n"), 0o644)
	gitRun(t, root, "add", ".gitignore")
	gitRun(t, root, "commit", "-qm", "ignore landed residue")
	clean := mustCreate(t, root, home, "landed-modifier-clean", "clean")
	ignored := mustCreate(t, root, home, "landed-modifier-ignored", "ignored")
	landAssignment(t, root, clean, "clean.txt")
	landAssignment(t, root, ignored, "landed.txt")
	mustWrite(t, filepath.Join(ignored.Path, "ignored.txt"), []byte("residue\n"), 0o644)

	bare, bareErr, bareCode := runCleanLanded(t, root, home, "--landed")
	if bareCode != 0 || bareErr != "" || !strings.Contains(bare, ignored.Path+",retain,") || !strings.Contains(bare, "ignored residuals require --discard-ignored") {
		t.Fatalf("bare plan = (%d, %q, %q), want ignored retain", bareCode, bare, bareErr)
	}
	widened, widenedErr, widenedCode := runCleanLanded(t, root, home, "--discard-ignored", "--full", "--landed")
	if widenedCode != 0 || widenedErr != "" || !strings.Contains(widened, ignored.Path+",discard-remove,") || !strings.Contains(widened, "ignored_paths[1]") {
		t.Fatalf("widened plan = (%d, %q, %q), want discard removal and preview", widenedCode, widened, widenedErr)
	}
	applied, applyErr, applyCode := runCleanLanded(t, root, home, "--discard-ignored", "--full", "--landed", "--apply", landedRowFingerprint(t, widened))
	if applyCode != 0 || applyErr != "" || strings.Count(applied, ",removed,") != 2 {
		t.Fatalf("modifier apply = (%d, %q, %q), want two removals", applyCode, applied, applyErr)
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
	home := filepath.Join(root, ".bench-home")
	creation := mustCreate(t, root, home, "landed-branch-assertion", "branch assertion")
	landAssignment(t, root, creation, "branch.txt")
	plan, planErr, planCode := runCleanLanded(t, root, home, "--discard-branch", "--landed")
	if planCode != 0 || planErr != "" || !strings.Contains(plan, "discards branch "+creation.Assignment.Branch) {
		t.Fatalf("assertion plan = (%d, %q, %q), want asserted branch detail", planCode, plan, planErr)
	}
	_, applyErr, applyCode := runCleanLanded(t, root, home, "--discard-branch", "--landed", "--apply", landedRowFingerprint(t, plan))
	if applyCode != 0 || applyErr != "" {
		t.Fatalf("assertion apply = (%d, %q), want success", applyCode, applyErr)
	}
	if _, err := os.Lstat(creation.Path); !os.IsNotExist(err) || git.OK("-C", root, "show-ref", "--verify", "--quiet", creation.Assignment.Branch) {
		t.Fatalf("assertion apply tree/branch outcome = %v/%t, want both absent", err, git.OK("-C", root, "show-ref", "--verify", "--quiet", creation.Assignment.Branch))
	}
}

func TestCleanLandedRetainsUnparseableLeaseAndSkipsUnprovableBranch(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	unknown := mustCreate(t, root, home, "landed-unknown-lease", "unknown lease")
	unprovable := mustCreate(t, root, home, "landed-unprovable", "unprovable")
	landAssignment(t, root, unknown, "unknown.txt")
	commitInWorktree(t, unprovable.Path, "unprovable.txt", "unmerged\n", "unprovable")
	lease, err := LeaseFile(unknown.Path)
	mustNoError(t, err)
	mustWrite(t, lease, []byte("not-a-lease\n"), 0o600)
	plan, planErr, planCode := runCleanLanded(t, root, home, "--landed")
	if planCode != 0 || planErr != "" || !strings.Contains(plan, unknown.Path+",retain,") || !strings.Contains(plan, "uncertain") || strings.Contains(plan, unprovable.Path) {
		t.Fatalf("plan = (%d, %q, %q), want unknown retain and no unprovable row", planCode, plan, planErr)
	}
	var summary, summaryErr strings.Builder
	if code := ResumeCleanCommand(root, home, nil, &summary, &summaryErr); code != 0 || !strings.Contains(summary.String(), "landed=1") {
		t.Fatalf("resume = (%d, %q, %q), want unknown lease counted landed", code, summary.String(), summaryErr.String())
	}
	_, applyErr, applyCode := runCleanLanded(t, root, home, "--landed", "--apply", landedRowFingerprint(t, plan))
	if applyCode != 0 || applyErr != "" {
		t.Fatalf("unknown lease apply = (%d, %q), want no-op success", applyCode, applyErr)
	}
	if _, err := os.Lstat(unknown.Path); err != nil {
		t.Fatalf("unknown lease path disappeared: %v", err)
	}
}
