// How the clean command wires each selection mode to the refusal renderer. Three sibling
// files own the rest: clean_set_outcomes_test.go owns what a stopped apply reports and proves
// the renderer itself, clean_set_apply_test.go owns which window an apply can refuse in and
// the fixture builders, and clean_set_refusal_test.go owns which member a refusal names. What
// is left, and what lives here, is the wiring between them: every route that can return a
// stale refusal from inside an apply reaches the refusal form through the command, so a call
// site that reverted to the plain outcome rows turns red.
package worktree

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

// TestCleanSetApplyTimeStaleWiring is COV-8. The renderer is proved elsewhere; this holds the
// command's wiring to it. Each case reaches the apply-time stale arm through the command and
// reads the refusal form back, so a wiring that reverted to the plain rows turns red.
func TestCleanSetApplyTimeStaleWiring(t *testing.T) {
	t.Parallel()
	t.Run("landed apply-time drift", func(t *testing.T) {
		t.Parallel()
		root, home, first, second, _ := landedSetFixture(t)
		plan, planErr, planCode := runCleanup(t, root, home, "--landed")
		if planCode != 0 || planErr != "" {
			t.Fatalf("plan = (%d, %q, %q), want one applicable plan", planCode, plan, planErr)
		}
		// The set removes in canonical identity order and the fixture's third member retains,
		// so the later of these two identities is the member the apply reaches second.
		drifted := first
		if second.Assignment.ID > first.Assignment.ID {
			drifted = second
		}
		// The drift lands after the first member's terminal receipt, so the command's entry
		// check passes and the refusal comes from inside the apply.
		j := defaultJoins()
		mutated := false
		j.cleanupBoundary = func(step LifecycleStep) error {
			if step == StepTerminalReceipt && !mutated {
				mutated = true
				mustWrite(t, filepath.Join(drifted.Path, "drifted.txt"), []byte("drifted\n"), 0o644)
			}
			return nil
		}

		digest := cleanupRowFingerprint(t, plan)
		stdout, stderr, code := runCleanupWith(t, j, root, home, "--landed", "--apply", digest)
		if code != 1 || stderr != "" || !mutated {
			t.Fatalf("landed apply = (%d, %q, %q) mutated=%t, want an apply-time refusal", code, stdout, stderr, mutated)
		}
		requireStaleRefusalRow(t, stdout, "unknown", "unknown", digest)
		if !strings.Contains(stdout, "bench worktree clean --landed") {
			t.Fatalf("landed apply = %q, want the selector-preserving re-plan command", stdout)
		}
	})
	// Both unclaimed routes read the branch namespace twice, once for the caller's plan and
	// once inside the apply. A branch appearing between those reads is the race, and the seam
	// stands in that window.
	for _, tc := range []struct {
		name  string
		apply func(string) []string
	}{
		{name: "unclaimed apply", apply: func(d string) []string {
			return []string{"--discard-branch", "--unclaimed", "--apply", d}
		}},
		{name: "unclaimed apply-current", apply: func(string) []string {
			return []string{"--discard-branch", "--unclaimed", "--apply-current"}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := newWorktreeRepo(t)
			home := filepath.Join(root, ".bench-home")
			for _, owner := range []string{"a", "b"} {
				ref := intent.AssignmentBranchRef(strings.Repeat(owner, 32), strings.Repeat("f", 32))
				gitRun(t, root, "branch", strings.TrimPrefix(ref, "refs/heads/"))
			}
			set, err := planUnclaimedAssignmentSet(root, unclaimedOptions())
			mustNoError(t, err)
			if len(set.rows) != 2 {
				t.Fatalf("unclaimed set = %#v, want two selected refs", set.rows)
			}
			j, raced := defaultJoins(), false
			j.cleanupBoundary = func(step LifecycleStep) error {
				if step == StepUnlockedReplan && !raced {
					raced = true
					ref := intent.AssignmentBranchRef(strings.Repeat("d", 32), strings.Repeat("f", 32))
					gitRun(t, root, "branch", strings.TrimPrefix(ref, "refs/heads/"))
				}
				return nil
			}

			stdout, stderr, code := runCleanupWith(t, j, root, home, tc.apply(set.fingerprint)...)
			if code != 1 || stderr != "" || !raced {
				t.Fatalf("%s = (%d, %q, %q) raced=%t, want an apply-time refusal", tc.name, code, stdout, stderr, raced)
			}
			requireStaleRefusalRow(t, stdout, "unclaimed", "none", set.fingerprint)
			if !strings.Contains(stdout, "bench worktree clean --discard-branch --unclaimed") {
				t.Fatalf("%s = %q, want the selector-preserving re-plan command", tc.name, stdout)
			}
			for _, row := range set.rows {
				if !git.OK("-C", root, "show-ref", "--verify", "--quiet", row.ref) {
					t.Fatalf("%s deleted %q", tc.name, row.ref)
				}
			}
		})
	}
}
