// What a cleanup-set apply reports. Most fixtures here drive an apply that does not finish,
// then read the result: which member completed, which one failed, which ones were never
// started, and the exact command a refusal offers as recovery. One reads the rows an apply
// reports for a member it will not touch at all. The assertions that read a decoded row live
// here too.
//
// The clean_set_*_test.go files split by the question each answers; this one answers what the
// apply reports.
package worktree

import (
	"errors"
	"maps"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
)

// faultAtSecondLock stops the second target transaction before it takes its lock, which
// leaves the first member completed and every later member unstarted.
func faultAtSecondLock(j *joins, stop error) {
	j.cleanupBoundary = atSecondHit(StepApplyLocked, func() error { return stop })
}

// TestCleanSetPartialApply is CL7. A transaction fault after an earlier completed removal
// reports the completed effect as completed and the faulted member as failed. The command
// claims no rollback, so the removed checkout stays removed.
func TestCleanSetPartialApply(t *testing.T) {
	t.Parallel()
	f := removableSetFixture(t, 3)
	j := defaultJoins()
	set := planExplicitSet(j, f.ambient(), f.root, explicitIdentities(f.creations), CleanupOptions{})
	if set.fingerprint == "" || len(set.rows) != 3 {
		t.Fatalf("explicit set = %#v, want three applicable members", set)
	}
	stop := errors.New("stop before the second member")
	faultAtSecondLock(&j, stop)

	plans, err := applyExplicitSet(j, f.ambient(), f.root, set, CleanupOptions{})
	if !errors.Is(err, stop) {
		t.Fatalf("apply error = %v, want the injected transaction fault", err)
	}
	completed := memberByID(t, f.creations, set.rows[0].assignment.ID)
	faulted := memberByID(t, f.creations, set.rows[1].assignment.ID)
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
	f := removableSetFixture(t, 3)
	j := defaultJoins()
	set := planExplicitSet(j, f.ambient(), f.root, explicitIdentities(f.creations), CleanupOptions{})
	if set.fingerprint == "" || len(set.rows) != 3 {
		t.Fatalf("explicit set = %#v, want three applicable members", set)
	}
	faultAtSecondLock(&j, errors.New("stop before the second member"))
	selectors := append(explicitTargetArguments(f.creations), "--apply", set.fingerprint)

	applied := runVerb(t, verbClean, f.callWith(j, selectors...))
	if applied.exit != 1 || applied.stderr != "" {
		t.Fatalf("partial apply = (%d, %q, %q), want a reported failure", applied.exit, applied.stdout, applied.stderr)
	}
	decoded := applied.mustRows(t, cleanupTable)
	if len(decoded) != len(f.creations) {
		t.Fatalf("partial apply rows = %#v, want one row per selected target", decoded)
	}
	rows := textRowsBy(t, decoded, "target")
	// The literal is the agent-facing field label the spec pins. Reading the token back
	// through its own constant would let a rename pass the gate and break the promised label.
	// The row is read by field, so no fixture-generated path is spliced into an expectation.
	unstarted := memberByID(t, f.creations, set.rows[2].assignment.ID)
	row := requireRowFor(t, rows, unstarted.Path)
	if row["action"] != "not-attempted" || row["detail"] != "not attempted; an earlier target in this set did not complete" {
		t.Fatalf("row for %q = %q/%q, want the unstarted label and detail", unstarted.Path, row["action"], row["detail"])
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
			f := removableSetFixture(t, 2)
			sorted := append([]Creation(nil), f.creations...)
			sort.Slice(sorted, func(a, b int) bool { return sorted[a].Assignment.ID < sorted[b].Assignment.ID })
			selectors := tc.selectors(sorted)
			plan := runVerb(t, verbClean, f.call(selectors...))
			if plan.exit != 0 || plan.stderr != "" {
				t.Fatalf("plan = (%d, %q, %q), want one applicable plan", plan.exit, plan.stdout, plan.stderr)
			}
			driftTracked(t, f.files, sorted[1].Assignment.ID)

			stale := runVerb(t, verbClean, f.call(append(selectors, "--apply", plan.mustFingerprint(t))...))
			if stale.exit != 1 || stale.stderr != "" {
				t.Fatalf("stale apply = (%d, %q, %q), want a refusal", stale.exit, stale.stdout, stale.stderr)
			}
			if !strings.Contains(stale.stdout, tc.want(sorted)) {
				t.Fatalf("stale apply = %q, want the re-plan command %q", stale.stdout, tc.want(sorted))
			}
			if strings.Contains(stale.stdout, "--apply") {
				t.Fatalf("stale apply = %q, want no replay of the refused digest", stale.stdout)
			}
		})
	}
}

// TestCleanSetUnclaimedStaleReplanAction is CL9 for the unclaimed branch selector, which
// keeps its own tracked and ignored spellings and therefore its own refusal row.
func TestCleanSetUnclaimedStaleReplanAction(t *testing.T) {
	t.Parallel()
	root, home := unclaimedBranchFixture(t, "a", "b")
	f := repoHome{root, home}
	plan := runVerb(t, verbClean, f.call("--discard-branch", "--unclaimed"))
	if plan.exit != 0 || plan.stderr != "" {
		t.Fatalf("plan = (%d, %q, %q), want one applicable plan", plan.exit, plan.stdout, plan.stderr)
	}
	extra := addUnclaimedBranch(t, root, "c")

	stale := runVerb(t, verbClean, f.call("--discard-branch", "--unclaimed", "--apply", plan.mustFingerprint(t)))
	if stale.exit != 1 || stale.stderr != "" {
		t.Fatalf("stale apply = (%d, %q, %q), want a refusal", stale.exit, stale.stdout, stale.stderr)
	}
	if !strings.Contains(stale.stdout, "bench worktree clean --discard-branch --unclaimed") || strings.Contains(stale.stdout, "--apply") {
		t.Fatalf("stale apply = %q, want the selector-preserving re-plan command alone", stale.stdout)
	}
	if !git.OK("-C", root, "show-ref", "--verify", "--quiet", extra) {
		t.Fatalf("stale unclaimed apply deleted %q", extra)
	}
}

// requireStaleRefusalRow holds the refusal row to the digest the call rejected, reading each
// field out of the decoded row rather than matching a spliced substring. A substring built
// around a raw digest passes or fails by the run's random identity, because the encoder
// quotes a digest that could read as a number.
//
// The action and the detail are literals. Both are agent-facing text this row promises, and
// reading either back through the value that produces it would let a rename pass the gate.
func requireStaleRefusalRow(t *testing.T, result verbResult, tracked, ignored, fingerprint string) {
	t.Helper()
	row := requireRowFor(t, textRowsBy(t, result.mustRows(t, cleanupTable), "target"), "unknown")
	got := []string{row["action"], row["tracked"], row["ignored"], row["recovery"], row["fingerprint"], row["detail"]}
	want := []string{"error", tracked, ignored, "none", fingerprint, "cleanup fingerprint is stale"}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("refusal row = %#v, want %#v", got, want)
	}
}

// requireRowFor returns the decoded cleanup row that names target, or fails t.
func requireRowFor(t *testing.T, rows map[string]textRow, target string) textRow {
	t.Helper()
	row, ok := rows[target]
	if !ok {
		t.Fatalf("rows = %q, want a row for %q", rows, target)
	}
	return row
}

// TestCleanSetRetainedMember is CL10 and CL11 inside a set apply. A member the plan retained
// is a member the apply will not touch, so the apply neither qualifies it beforehand nor
// rewrites the verdict the plan gave it.
func TestCleanSetRetainedMember(t *testing.T) {
	t.Parallel()
	t.Run("keeps the planned verdict", func(t *testing.T) {
		t.Parallel()
		f := retainedMemberFixture(t)
		targets := []string{"--target", f.removable.Assignment.ID, "--target", f.retained.Assignment.ID}
		plan := runVerb(t, verbClean, f.call(targets...))
		if plan.exit != 0 || plan.stderr != "" {
			t.Fatalf("plan = (%d, %q, %q), want one applicable plan", plan.exit, plan.stdout, plan.stderr)
		}
		planned := requireRowFor(t, textRowsBy(t, plan.mustRows(t, cleanupTable), "target"), f.retained.Path)
		if planned["action"] != string(ActionRetain) {
			t.Fatalf("planned retained row = %q, want the ignored-residue refusal", planned)
		}

		applied := runVerb(t, verbClean, f.call(append(targets, "--apply", plan.mustFingerprint(t))...))
		if applied.exit != 0 || applied.stderr != "" || strings.Count(applied.stdout, ",removed,") != 1 {
			t.Fatalf("apply = (%d, %q, %q), want exactly one removal", applied.exit, applied.stdout, applied.stderr)
		}
		// The row passes through, so it still carries the set digest the plan answered under
		// rather than a digest a transaction the apply never opened would have produced.
		if got := requireRowFor(t, textRowsBy(t, applied.mustRows(t, cleanupTable), "target"), f.retained.Path); !maps.Equal(got, planned) {
			t.Fatalf("applied retained row = %q, want the planned row %q", got, planned)
		}
		if _, statErr := os.Stat(f.retained.Path); statErr != nil {
			t.Fatalf("apply removed the retained member %s: %v", f.retained.Path, statErr)
		}
		if _, statErr := os.Lstat(f.removable.Path); !os.IsNotExist(statErr) {
			t.Fatalf("apply left the removable member %s: %v", f.removable.Path, statErr)
		}
	})
	t.Run("preflight skips retentions", func(t *testing.T) {
		t.Parallel()
		f := retainedMemberFixture(t)
		j := defaultJoins()
		set := planExplicitSet(j, f.ambient(), f.root, []string{f.removable.Assignment.ID, f.retained.Assignment.ID}, CleanupOptions{})
		if set.fingerprint == "" || len(set.rows) != 2 {
			t.Fatalf("explicit set = %#v, want two members", set)
		}
		// The retained member's own residue changes. The apply was never going to touch that
		// member, so a preflight that qualified it would refuse a removal on evidence that
		// decides nothing about the member being removed.
		mustWrite(t, filepath.Join(f.retained.Path, "ignored-two.txt"), []byte("more residue\n"), 0o644)

		plans, err := applyExplicitSet(j, f.ambient(), f.root, set, CleanupOptions{})
		if err != nil {
			t.Fatalf("apply error = %v, want the retained member's drift to decide nothing", err)
		}
		if _, statErr := os.Lstat(f.removable.Path); !os.IsNotExist(statErr) {
			t.Fatalf("the removable member %s survived: %v", f.removable.Path, statErr)
		}
		if _, statErr := os.Stat(f.retained.Path); statErr != nil {
			t.Fatalf("apply removed the retained member %s: %v", f.retained.Path, statErr)
		}
		if len(plans) != 2 {
			t.Fatalf("apply rows = %#v, want one row per member", plans)
		}
	})
	t.Run("a refused apply keeps the retained verdict", func(t *testing.T) {
		t.Parallel()
		f := retainedMemberFixture(t)
		j := defaultJoins()
		set := planExplicitSet(j, f.ambient(), f.root, []string{f.removable.Assignment.ID, f.retained.Assignment.ID}, CleanupOptions{})
		if set.fingerprint == "" || len(set.rows) != 2 {
			t.Fatalf("explicit set = %#v, want two members", set)
		}
		mustWrite(t, filepath.Join(f.removable.Path, "removable.txt"), []byte("drifted\n"), 0o644)

		plans, err := applyExplicitSet(j, f.ambient(), f.root, set, CleanupOptions{})
		if !errors.Is(err, errStaleFingerprint) {
			t.Fatalf("apply error = %v, want a preflight refusal", err)
		}
		requireMembersPresent(t, []Creation{f.removable, f.retained})
		// The refused member was going to be removed and was not, which is what unstarted
		// means. The retained member was never a removal, so calling it unstarted would erase
		// the CL10 authority the plan spent on it.
		for _, plan := range plans {
			want := ActionNotAttempted
			if plan.Target == f.retained.Path {
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
		f := removableSetFixture(t, 2)
		j := defaultJoins()
		set := planExplicitSet(j, f.ambient(), f.root, explicitIdentities(f.creations), CleanupOptions{})
		if set.fingerprint == "" || len(set.rows) != 2 {
			t.Fatalf("explicit set = %#v, want two applicable members", set)
		}
		targets := explicitTargetArguments(f.creations)
		plan := runVerb(t, verbClean, f.call(targets...))
		if plan.exit != 0 || plan.stderr != "" {
			t.Fatalf("plan = (%d, %q, %q), want one applicable plan", plan.exit, plan.stdout, plan.stderr)
		}
		// The drift lands inside the second member's locked transaction, so the command's
		// entry check and the preflight both pass and the refusal comes from the apply.
		drifted := memberByID(t, f.creations, set.rows[1].assignment.ID)
		j.cleanupBoundary = atSecondHit(StepApplyLocked, func() error {
			driftTracked(t, f.files, drifted.Assignment.ID)
			return nil
		})

		digest := plan.mustFingerprint(t)
		stale := runVerb(t, verbClean, f.callWith(j, append(targets, "--apply", digest)...))
		if stale.exit != 1 || stale.stderr != "" {
			t.Fatalf("late-drift apply = (%d, %q, %q), want a refusal", stale.exit, stale.stdout, stale.stderr)
		}
		requireStaleRefusalRow(t, stale, "unknown", "unknown", digest)
		// The renderer names members in canonical identity order, which the fixture draws at
		// random, so the expectation reads that order from the same source the renderer uses.
		want := "bench worktree clean " + strings.Join(set.targetSelectors(), " ")
		if !strings.Contains(stale.stdout, want) {
			t.Fatalf("late-drift apply = %q, want the re-plan command %q", stale.stdout, want)
		}
		if strings.Contains(stale.stdout, "--apply") {
			t.Fatalf("late-drift apply = %q, want no replay of the refused digest", stale.stdout)
		}
	})
	t.Run("unclaimed applier re-plan mismatch", func(t *testing.T) {
		t.Parallel()
		root, _ := unclaimedBranchFixture(t, "a", "b")
		set, err := planUnclaimedAssignmentSet(root, unclaimedOptions())
		mustNoError(t, err)
		if len(set.rows) != 2 {
			t.Fatalf("unclaimed set = %#v, want two selected refs", set.rows)
		}
		// A ref appears between the plan and the apply, which only the applier's own re-plan
		// sees. The rendering of this refusal is covered end to end by the wiring test; what
		// only a direct call can read is that the applier returns no rows of its own.
		extra := addUnclaimedBranch(t, root, "c")

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
