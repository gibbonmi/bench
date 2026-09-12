// What a cleanup-set apply reports. Most fixtures here drive an apply that does not finish,
// then read the result: which member completed, which one failed, which ones were never
// started, and the exact command a refusal offers as recovery. One reads the rows an apply
// reports for a member it will not touch at all. The assertions that read a rendered row live
// here too.
//
// The clean_set_*_test.go files split by the question each answers; this one answers what the
// apply reports.
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
	// The literal is the agent-facing field label the spec pins. Reading the token back
	// through its own constant would let a rename pass the gate and break the promised label.
	// The row is read by field, so no fixture-generated path is spliced into an expectation.
	unstarted := memberByID(t, creations, set.rows[2].assignment.ID)
	if action := cleanupRowFields(rowForTarget(t, stdout, unstarted.Path))[1]; action != "not-attempted" {
		t.Fatalf("row for %q = %q, want not-attempted", unstarted.Path, action)
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

// requireStaleRefusalRow holds the refusal row to the digest the call rejected, reading each
// field out of the rendered row rather than matching a spliced substring. A substring built
// around a raw digest passes or fails by the run's random identity, because the encoder
// quotes a digest that could read as a number.
func requireStaleRefusalRow(t *testing.T, output, tracked, ignored, fingerprint string) {
	t.Helper()
	fields := cleanupRowFields(rowForTarget(t, output, "unknown"))
	got := []string{fields[1], fields[2], fields[3], fields[4], cleanupRowValue(fields[5]), fields[6]}
	want := []string{string(ActionError), tracked, ignored, "none", fingerprint, errStaleFingerprint.Error()}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("refusal row = %#v, want %#v", got, want)
	}
}

// rowForTarget returns the one rendered cleanup row that names target.
func rowForTarget(t *testing.T, output, target string) string {
	t.Helper()
	for _, row := range cleanupRows(output) {
		if cleanupRowFields(row)[0] == target {
			return row
		}
	}
	t.Fatalf("output = %q, want a row for %q", output, target)
	return ""
}

// TestCleanSetRetainedMember is CL10 and CL11 inside a set apply. A member the plan retained
// is a member the apply will not touch, so the apply neither qualifies it beforehand nor
// rewrites the verdict the plan gave it.
func TestCleanSetRetainedMember(t *testing.T) {
	t.Parallel()
	t.Run("keeps the planned verdict", func(t *testing.T) {
		t.Parallel()
		root, home, removable, retained := retainedMemberFixture(t)
		targets := []string{"--target", removable.Assignment.ID, "--target", retained.Assignment.ID}
		plan, planErr, planCode := runCleanup(t, root, home, targets...)
		if planCode != 0 || planErr != "" {
			t.Fatalf("plan = (%d, %q, %q), want one applicable plan", planCode, plan, planErr)
		}
		planned := rowForTarget(t, plan, retained.Path)
		if !strings.Contains(planned, ",retain,") {
			t.Fatalf("planned retained row = %q, want the ignored-residue refusal", planned)
		}

		applied, applyErr, applyCode := runCleanup(t, root, home, append(targets, "--apply", cleanupRowFingerprint(t, plan))...)
		if applyCode != 0 || applyErr != "" || strings.Count(applied, ",removed,") != 1 {
			t.Fatalf("apply = (%d, %q, %q), want exactly one removal", applyCode, applied, applyErr)
		}
		// The row passes through, so it still carries the set digest the plan answered under
		// rather than a digest a transaction the apply never opened would have produced.
		if got := rowForTarget(t, applied, retained.Path); got != planned {
			t.Fatalf("applied retained row = %q, want the planned row %q", got, planned)
		}
		if _, statErr := os.Stat(retained.Path); statErr != nil {
			t.Fatalf("apply removed the retained member %s: %v", retained.Path, statErr)
		}
		if _, statErr := os.Lstat(removable.Path); !os.IsNotExist(statErr) {
			t.Fatalf("apply left the removable member %s: %v", removable.Path, statErr)
		}
	})
	t.Run("preflight skips retentions", func(t *testing.T) {
		t.Parallel()
		root, _, removable, retained := retainedMemberFixture(t)
		j := defaultJoins()
		set := planExplicitSet(j, root, []string{removable.Assignment.ID, retained.Assignment.ID}, CleanupOptions{})
		if set.fingerprint == "" || len(set.rows) != 2 {
			t.Fatalf("explicit set = %#v, want two members", set)
		}
		// The retained member's own residue changes. The apply was never going to touch that
		// member, so a preflight that qualified it would refuse a removal on evidence that
		// decides nothing about the member being removed.
		mustWrite(t, filepath.Join(retained.Path, "ignored-two.txt"), []byte("more residue\n"), 0o644)

		plans, err := applyExplicitSet(j, root, set, CleanupOptions{})
		if err != nil {
			t.Fatalf("apply error = %v, want the retained member's drift to decide nothing", err)
		}
		if _, statErr := os.Lstat(removable.Path); !os.IsNotExist(statErr) {
			t.Fatalf("the removable member %s survived: %v", removable.Path, statErr)
		}
		if _, statErr := os.Stat(retained.Path); statErr != nil {
			t.Fatalf("apply removed the retained member %s: %v", retained.Path, statErr)
		}
		if len(plans) != 2 {
			t.Fatalf("apply rows = %#v, want one row per member", plans)
		}
	})
	t.Run("a refused apply keeps the retained verdict", func(t *testing.T) {
		t.Parallel()
		root, _, removable, retained := retainedMemberFixture(t)
		j := defaultJoins()
		set := planExplicitSet(j, root, []string{removable.Assignment.ID, retained.Assignment.ID}, CleanupOptions{})
		if set.fingerprint == "" || len(set.rows) != 2 {
			t.Fatalf("explicit set = %#v, want two members", set)
		}
		mustWrite(t, filepath.Join(removable.Path, "removable.txt"), []byte("drifted\n"), 0o644)

		plans, err := applyExplicitSet(j, root, set, CleanupOptions{})
		if !errors.Is(err, errStaleFingerprint) {
			t.Fatalf("apply error = %v, want a preflight refusal", err)
		}
		requireMembersPresent(t, []Creation{removable, retained})
		// The refused member was going to be removed and was not, which is what unstarted
		// means. The retained member was never a removal, so calling it unstarted would erase
		// the CL10 authority the plan spent on it.
		for _, plan := range plans {
			want := ActionNotAttempted
			if plan.Target == retained.Path {
				want = ActionRetain
			}
			if plan.Action != want {
				t.Fatalf("refused row %q = %q, want %q", plan.Target, plan.Action, want)
			}
		}
	})
}

// TestCleanSetApplyTimeStaleRefusal is COV-6. An apply can return a stale refusal from
// inside itself, after the command's entry check already passed. That arm renders the
// refusal form rather than the outcome rows, and it is the only place the agent learns the
// plan died mid-apply. Each mode reaches it by its own route.
func TestCleanSetApplyTimeStaleRefusal(t *testing.T) {
	t.Parallel()
	t.Run("explicit late drift through the command", func(t *testing.T) {
		t.Parallel()
		root, home, creations := removableSetFixture(t, 2)
		files := landedFileByAssignment(creations, "member-0.txt", "member-1.txt")
		j := defaultJoins()
		set := planExplicitSet(j, root, explicitIdentities(creations), CleanupOptions{})
		if set.fingerprint == "" || len(set.rows) != 2 {
			t.Fatalf("explicit set = %#v, want two applicable members", set)
		}
		targets := explicitTargetArguments(creations)
		plan, planErr, planCode := runCleanup(t, root, home, targets...)
		if planCode != 0 || planErr != "" {
			t.Fatalf("plan = (%d, %q, %q), want one applicable plan", planCode, plan, planErr)
		}
		// The drift lands inside the second member's locked transaction, so the command's
		// entry check and the preflight both pass and the refusal comes from the apply.
		drifted := memberByID(t, creations, set.rows[1].assignment.ID)
		locks := 0
		j.cleanupBoundary = func(step LifecycleStep) error {
			if step == StepApplyLocked {
				locks++
			}
			if locks == 2 {
				driftTracked(t, files, drifted.Assignment.ID)
			}
			return nil
		}

		digest := cleanupRowFingerprint(t, plan)
		stdout, stderr, code := runCleanupWith(t, j, root, home, append(targets, "--apply", digest)...)
		if code != 1 || stderr != "" {
			t.Fatalf("late-drift apply = (%d, %q, %q), want a refusal", code, stdout, stderr)
		}
		requireStaleRefusalRow(t, stdout, "unknown", "unknown", digest)
		// The renderer names members in canonical identity order, which the fixture draws at
		// random, so the expectation reads that order from the same source the renderer uses.
		want := "bench worktree clean " + strings.Join(set.targetSelectors(), " ")
		if !strings.Contains(stdout, want) {
			t.Fatalf("late-drift apply = %q, want the re-plan command %q", stdout, want)
		}
		if strings.Contains(stdout, "--apply") {
			t.Fatalf("late-drift apply = %q, want no replay of the refused digest", stdout)
		}
	})
	t.Run("unclaimed applier re-plan mismatch", func(t *testing.T) {
		t.Parallel()
		root := newWorktreeRepo(t)
		for _, owner := range []string{"a", "b"} {
			branch := intent.AssignmentBranchRef(strings.Repeat(owner, 32), strings.Repeat("f", 32))
			gitRun(t, root, "branch", strings.TrimPrefix(branch, "refs/heads/"))
		}
		set, err := planUnclaimedAssignmentSet(root, unclaimedOptions())
		mustNoError(t, err)
		if len(set.rows) != 2 {
			t.Fatalf("unclaimed set = %#v, want two selected refs", set.rows)
		}
		// A ref appears between the plan and the apply, which only the applier's own re-plan
		// sees. The rendering of this refusal is covered end to end by the wiring test; what
		// only a direct call can read is that the applier returns no rows of its own.
		extra := intent.AssignmentBranchRef(strings.Repeat("c", 32), strings.Repeat("f", 32))
		gitRun(t, root, "branch", strings.TrimPrefix(extra, "refs/heads/"))

		plans, applyErr := applyUnclaimedAssignmentSet(defaultJoins(), root, set, unclaimedOptions())
		if !errors.Is(applyErr, errStaleFingerprint) {
			t.Fatalf("apply error = %v, want the applier's own stale refusal", applyErr)
		}
		if plans != nil {
			t.Fatalf("stale apply rows = %#v, want none; the caller owns the refusal row", plans)
		}
		for _, ref := range []string{set.rows[0].ref, set.rows[1].ref, extra} {
			if !git.OK("-C", root, "show-ref", "--verify", "--quiet", ref) {
				t.Fatalf("stale apply deleted %q", ref)
			}
		}
	})
}
