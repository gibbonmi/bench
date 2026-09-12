// What a stopped cleanup-set apply reports. The fixtures here drive an apply that does not
// finish, then read the result: which member completed, which one failed, which ones were
// never started, and the exact command a refusal offers as recovery. The sibling file
// clean_set_apply_test.go owns the other half — when an apply refuses at all — and holds
// the fixture helpers both halves share.
package worktree

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

// faultAtSecondLock stops the second target transaction before it takes its lock, which
// leaves the first member completed and every later member unstarted.
func faultAtSecondLock(j *joins, stop error) {
	locks := 0
	j.cleanupBoundary = func(step LifecycleStep) error {
		if step != StepApplyLocked {
			return nil
		}
		locks++
		if locks == 2 {
			return stop
		}
		return nil
	}
}

// TestCleanSetPartialApply is CL7. A transaction fault after an earlier completed removal
// reports the completed effect as completed and the faulted member as failed. The command
// claims no rollback, so the removed checkout stays removed.
func TestCleanSetPartialApply(t *testing.T) {
	t.Parallel()
	root, _, creations := removableSetFixture(t, 3)
	j := defaultJoins()
	set := planExplicitSet(j, root, explicitIdentities(creations), CleanupOptions{})
	if set.fingerprint == "" || len(set.rows) != 3 {
		t.Fatalf("explicit set = %#v, want three applicable members", set)
	}
	stop := errors.New("stop before the second member")
	faultAtSecondLock(&j, stop)

	plans, err := applyExplicitSet(j, root, set, CleanupOptions{})
	if !errors.Is(err, stop) {
		t.Fatalf("apply error = %v, want the injected transaction fault", err)
	}
	completed := memberByID(t, creations, set.rows[0].assignment.ID)
	faulted := memberByID(t, creations, set.rows[1].assignment.ID)
	if _, statErr := os.Lstat(completed.Path); !os.IsNotExist(statErr) {
		t.Fatalf("the completed member %s came back: %v", completed.Path, statErr)
	}
	if len(plans) != 3 || plans[0].Action != ActionRemoved || plans[1].Action != ActionError {
		t.Fatalf("partial rows = %#v, want a completed row and a failed row", plans)
	}
	if plans[0].Target != completed.Path || plans[1].Target != faulted.Path {
		t.Fatalf("partial row targets = %q and %q, want %q and %q", plans[0].Target, plans[1].Target, completed.Path, faulted.Path)
	}
	if !strings.Contains(plans[1].Reason, stop.Error()) {
		t.Fatalf("failed row detail = %q, want the transaction fault", plans[1].Reason)
	}
}

// TestCleanSetUnstartedOutcomes is CL8. A partial apply accounts for the whole selection:
// every member the set never reached reports its own unstarted outcome rather than dropping
// out of the result.
func TestCleanSetUnstartedOutcomes(t *testing.T) {
	t.Parallel()
	root, home, creations := removableSetFixture(t, 3)
	j := defaultJoins()
	set := planExplicitSet(j, root, explicitIdentities(creations), CleanupOptions{})
	if set.fingerprint == "" || len(set.rows) != 3 {
		t.Fatalf("explicit set = %#v, want three applicable members", set)
	}
	faultAtSecondLock(&j, errors.New("stop before the second member"))
	selectors := append(explicitTargetArguments(creations), "--apply", set.fingerprint)

	stdout, stderr, code := runCleanupWith(t, j, root, home, selectors...)
	if code != 1 || stderr != "" {
		t.Fatalf("partial apply = (%d, %q, %q), want a reported failure", code, stdout, stderr)
	}
	rows := cleanupRows(stdout)
	if len(rows) != len(creations) {
		t.Fatalf("partial apply rows = %#v, want one row per selected target", rows)
	}
	unstarted := memberByID(t, creations, set.rows[2].assignment.ID)
	if !strings.Contains(stdout, unstarted.Path+","+string(ActionNotAttempted)+",") {
		t.Fatalf("partial apply = %q, want %q reported as not attempted", stdout, unstarted.Path)
	}
	if _, statErr := os.Stat(unstarted.Path); statErr != nil {
		t.Fatalf("the unstarted member %s was removed: %v", unstarted.Path, statErr)
	}
}

// TestCleanSetStaleReplanAction is CL9. Every stale refusal names the exact command that
// re-plans the same selection: the same mode, the same discard modifiers, and one canonical
// identity per explicit member. The action never carries the digest the call just refused.
func TestCleanSetStaleReplanAction(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name      string
		selectors func([]Creation) []string
		want      func([]Creation) string
	}{
		{
			name: "explicit",
			selectors: func(c []Creation) []string {
				return append([]string{"--discard-ignored"}, explicitTargetArguments(c)...)
			},
			want: func(c []Creation) string {
				return "bench worktree clean --discard-ignored " + strings.Join(explicitTargetArguments(c), " ")
			},
		},
		{
			name:      "landed",
			selectors: func([]Creation) []string { return []string{"--discard-branch", "--landed"} },
			want:      func([]Creation) string { return "bench worktree clean --discard-branch --landed" },
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root, home, creations := removableSetFixture(t, 2)
			files := landedFileByAssignment(creations, "member-0.txt", "member-1.txt")
			sorted := append([]Creation(nil), creations...)
			sort.Slice(sorted, func(a, b int) bool { return sorted[a].Assignment.ID < sorted[b].Assignment.ID })
			selectors := tc.selectors(sorted)
			plan, planErr, planCode := runCleanup(t, root, home, selectors...)
			if planCode != 0 || planErr != "" {
				t.Fatalf("plan = (%d, %q, %q), want one applicable plan", planCode, plan, planErr)
			}
			driftTracked(t, files, sorted[1].Assignment.ID)

			stdout, stderr, code := runCleanup(t, root, home, append(selectors, "--apply", cleanupRowFingerprint(t, plan))...)
			if code != 1 || stderr != "" {
				t.Fatalf("stale apply = (%d, %q, %q), want a refusal", code, stdout, stderr)
			}
			if !strings.Contains(stdout, tc.want(sorted)) {
				t.Fatalf("stale apply = %q, want the re-plan command %q", stdout, tc.want(sorted))
			}
			if strings.Contains(stdout, "--apply") {
				t.Fatalf("stale apply = %q, want no replay of the refused digest", stdout)
			}
		})
	}
}

// TestCleanSetUnclaimedStaleReplanAction is CL9 for the unclaimed branch selector, which
// keeps its own tracked and ignored spellings and therefore its own refusal row.
func TestCleanSetUnclaimedStaleReplanAction(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	for _, owner := range []string{"a", "b"} {
		branch := intent.AssignmentBranchRef(strings.Repeat(owner, 32), strings.Repeat("f", 32))
		gitRun(t, root, "branch", strings.TrimPrefix(branch, "refs/heads/"))
	}
	plan, planErr, planCode := runCleanup(t, root, home, "--discard-branch", "--unclaimed")
	if planCode != 0 || planErr != "" {
		t.Fatalf("plan = (%d, %q, %q), want one applicable plan", planCode, plan, planErr)
	}
	extra := intent.AssignmentBranchRef(strings.Repeat("c", 32), strings.Repeat("f", 32))
	gitRun(t, root, "branch", strings.TrimPrefix(extra, "refs/heads/"))

	stdout, stderr, code := runCleanup(t, root, home, "--discard-branch", "--unclaimed", "--apply", cleanupRowFingerprint(t, plan))
	if code != 1 || stderr != "" {
		t.Fatalf("stale apply = (%d, %q, %q), want a refusal", code, stdout, stderr)
	}
	if !strings.Contains(stdout, "bench worktree clean --discard-branch --unclaimed") || strings.Contains(stdout, "--apply") {
		t.Fatalf("stale apply = %q, want the selector-preserving re-plan command alone", stdout)
	}
	if !git.OK("-C", root, "show-ref", "--verify", "--quiet", extra) {
		t.Fatalf("stale unclaimed apply deleted %q", extra)
	}
}
