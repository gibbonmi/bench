// Mode witnesses of the shared cleanup traversal: what each adapter does with an empty or
// retained-only selection, a completed prefix, its injected fault boundaries, and an authority
// refusal inside the locked transaction. Each table names its mode and case, so a later
// adapter migration adds its own rows.
package worktree

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

// stepTrace is a cleanup boundary that records every lifecycle step it reaches, in order.
func stepTrace() (*[]LifecycleStep, Fault) {
	steps := &[]LifecycleStep{}
	return steps, func(step LifecycleStep) error {
		*steps = append(*steps, step)
		return nil
	}
}

// modeApply plans one selection through the explicit adapter, or through the landed adapter
// over the whole repository, which ignores targets. It returns each recorded row's assignment
// identity, the planned rows, and the apply, which runs under the joins j holds at that time.
func modeApply(t *testing.T, j *joins, f repoHome, landed bool, targets []string) ([]string, []CleanupPlan, func() ([]CleanupPlan, error)) {
	t.Helper()
	var ids []string
	if landed {
		set, err := planLandedSet(*j, nil, f.root, CleanupOptions{}, "")
		mustNoError(t, err)
		for _, row := range set.rows {
			ids = append(ids, row.assignment.ID)
		}
		return ids, landedRowPlans(set.rows), func() ([]CleanupPlan, error) {
			return applyLandedSet(*j, f.ambient(), f.root, set, CleanupOptions{}, "")
		}
	}
	set := planExplicitSet(*j, f.ambient(), f.root, targets, CleanupOptions{})
	for _, row := range set.rows {
		ids = append(ids, row.assignment.ID)
	}
	return ids, set.plans(), func() ([]CleanupPlan, error) { return applyExplicitSet(*j, f.ambient(), f.root, set, CleanupOptions{}) }
}

// TestCleanupTraversalEmptyModes is R01. A selection with no removable member reaches no
// lifecycle boundary, changes no ref or record, and reports each retained row as planned.
func TestCleanupTraversalEmptyModes(t *testing.T) {
	t.Parallel()
	empty := func(t *testing.T) (repoHome, []string) {
		root := newWorktreeRepo(t)
		return repoHome{root, filepath.Join(root, ".bench-home")}, nil
	}
	for _, tc := range []struct {
		name   string
		landed bool
		rows   int
		build  func(*testing.T) (repoHome, []string)
	}{
		{"explicit/empty", false, 0, empty},
		{"explicit/retained-only", false, 1, func(t *testing.T) (repoHome, []string) {
			f := retainedMemberFixture(t)
			return f.repoHome, []string{f.retained.Assignment.ID}
		}},
		{"landed/empty", true, 0, empty},
		{"landed/retained-only", true, 1, func(t *testing.T) (repoHome, []string) {
			root, home := ignoringRepo(t)
			retained := landedMember(t, root, home, "set-retained-residue", "retained.txt")
			mustWrite(t, filepath.Join(retained.Path, "ignored-one.txt"), []byte("residue\n"), 0o644)
			return repoHome{root, home}, nil
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f, targets := tc.build(t)
			j := defaultJoins()
			_, want, apply := modeApply(t, &j, f, tc.landed, targets)
			if len(want) != tc.rows {
				t.Fatalf("planned rows = %#v, want %d", want, tc.rows)
			}
			steps, trace := stepTrace()
			j.cleanupBoundary = trace
			before := strings.Join(baselineRefs(t, f.root), "\n") + repositoryState(t, f.root)

			plans, err := apply()
			mustNoError(t, err)
			if len(*steps) != 0 {
				t.Fatalf("apply reached %v, want no lifecycle boundary", *steps)
			}
			if len(plans) != len(want) || (len(want) != 0 && !reflect.DeepEqual(plans, want)) {
				t.Fatalf("rows = %#v, want the planned rows %#v unchanged", plans, want)
			}
			if after := strings.Join(baselineRefs(t, f.root), "\n") + repositoryState(t, f.root); after != before {
				t.Fatalf("apply changed the repository: %q -> %q", before, after)
			}
		})
	}
}

// TestCleanupTraversalCompletedPrefix is R19. A fault on a later member leaves the earlier
// removal complete: its tree stays gone, its receipt stays complete, and its row still
// reports the removal its transaction made.
func TestCleanupTraversalCompletedPrefix(t *testing.T) {
	t.Parallel()
	f := removableSetFixture(t, 3)
	j := defaultJoins()
	set := planExplicitSet(j, f.ambient(), f.root, explicitIdentities(f.creations), CleanupOptions{})
	stop := errors.New("stop before the second member")
	faultAtSecondLock(&j, stop)

	plans, err := applyExplicitSet(j, f.ambient(), f.root, set, CleanupOptions{})
	if !errors.Is(err, stop) {
		t.Fatalf("apply error = %v, want the injected fault", err)
	}
	completed := memberByID(t, f.creations, set.rows[0].assignment.ID)
	if len(plans) != 3 || plans[0].Action != ActionRemoved || plans[0].Target != completed.Path || plans[0].Reason != "" {
		t.Fatalf("rows = %#v, want the completed removal of %s first", plans, completed.Path)
	}
	if _, statErr := os.Lstat(completed.Path); !os.IsNotExist(statErr) {
		t.Fatalf("the completed member %s came back: %v", completed.Path, statErr)
	}
	ledger, readErr := intent.Read(f.root)
	mustNoError(t, readErr)
	if len(ledger.CleanupReceipts) != 1 || ledger.CleanupReceipts[0].State != intent.ReceiptComplete || ledger.CleanupReceipts[0].Action != string(ActionRemoved) {
		t.Fatalf("receipts = %#v, want one complete removal receipt", ledger.CleanupReceipts)
	}
	for _, row := range set.rows[1:] {
		requireMembersPresent(t, []Creation{memberByID(t, f.creations, row.assignment.ID)})
	}
}

// TestCleanupTraversalModeFaultBoundaries is R27. Each mode reaches exactly the lifecycle
// boundaries its original adapter reached, in the same order. The expected traces are the
// ones the original adapters produced.
func TestCleanupTraversalModeFaultBoundaries(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		landed bool
		member []LifecycleStep
	}{
		{"explicit/all-success", false, []LifecycleStep{StepMemberRequalify, StepApplyLocked, StepReceipt, StepUnlock, StepRemovalAttempt, StepRemoval, StepBranch}},
		{"landed/all-success", true, []LifecycleStep{StepMemberRequalify, StepApplyLocked, StepReceipt, StepUnlock, StepRemovalAttempt, StepRemoval, StepBranch, StepTerminalReceipt}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := removableSetFixture(t, 2)
			j := defaultJoins()
			_, _, apply := modeApply(t, &j, f.repoHome, tc.landed, explicitIdentities(f.creations))
			steps, trace := stepTrace()
			j.cleanupBoundary = trace

			_, err := apply()
			mustNoError(t, err)
			if want := append(append([]LifecycleStep{}, tc.member...), tc.member...); !reflect.DeepEqual(*steps, want) {
				t.Fatalf("boundaries = %v, want %v", *steps, want)
			}
		})
	}
	// A landed branch writes no discarded ref, so only the unique branch reaches the discard
	// windows, and the unrecorded stage has no requalify boundary.
	t.Run("unrecorded/all-success", func(t *testing.T) {
		t.Parallel()
		root, home := unclaimedBranchFixture(t, "a")
		unique, _ := uniqueBranch(t, root, "e")
		f, options := repoHome{root, home}, CleanupOptions{DiscardBranch: true}
		set := planExplicitSet(defaultJoins(), baselineAmbient(f), root, branchTargets(unclaimedBranchRef("a"), unique), options)
		j := defaultJoins()
		steps, trace := stepTrace()
		j.cleanupBoundary = trace

		_, err := applyExplicitSet(j, baselineAmbient(f), root, set, options)
		mustNoError(t, err)
		if want := []LifecycleStep{StepDiscardedRefWrite, StepDiscardedRefAbsent, StepDiscardedBranchDelete}; !reflect.DeepEqual(*steps, want) {
			t.Fatalf("boundaries = %v, want %v", *steps, want)
		}
	})
}

// TestCleanupTraversalAuthorityPreservation is R28. Ignored residue that appears inside a
// member's own lock takes its removal away: the locked recheck refuses, the member and its
// residue stay, and the retained member behind it keeps its planned verdict.
func TestCleanupTraversalAuthorityPreservation(t *testing.T) {
	t.Parallel()
	for _, landed := range []bool{false, true} {
		t.Run(map[bool]string{false: "explicit", true: "landed"}[landed], func(t *testing.T) {
			t.Parallel()
			f := retainedMemberFixture(t)
			j := defaultJoins()
			ids, planned, apply := modeApply(t, &j, f.repoHome, landed, []string{f.removable.Assignment.ID, f.retained.Assignment.ID})
			if len(ids) != 2 || len(planned) != 2 {
				t.Fatalf("planned rows = %#v, want two applicable members", planned)
			}
			// The set orders its members by identity, so the removable member's row index is read
			// from the plan rather than assumed.
			removable, retained := 0, 1
			if ids[0] != f.removable.Assignment.ID {
				removable, retained = 1, 0
			}
			if !planned[removable].Action.Removes() || planned[retained].Action.Removes() {
				t.Fatalf("planned rows = %#v, want one removable and one retained member", planned)
			}
			residue := filepath.Join(f.removable.Path, "ignored-late.txt")
			j.cleanupBoundary = atStep(StepApplyLocked, func() error {
				mustWrite(t, residue, []byte("residue\n"), 0o644)
				return nil
			})

			plans, err := apply()
			if !errors.Is(err, errStaleFingerprint) {
				t.Fatalf("apply error = %v, want the locked recheck to refuse", err)
			}
			if _, statErr := os.Stat(residue); statErr != nil {
				t.Fatalf("the preserved residue %s was removed: %v", residue, statErr)
			}
			requireMembersPresent(t, []Creation{f.removable, f.retained})
			if len(plans) != 2 || plans[removable].Action != ActionError || !reflect.DeepEqual(plans[retained], planned[retained]) {
				t.Fatalf("rows = %#v, want the refused member failed and the retained member as planned", plans)
			}
		})
	}
}

// unrecordedSuffixSet is one recorded removable member ahead of three unique unclaimed
// branches, planned with --discard-branch at the fixed plan day. The recorded member makes the
// branches' indexes in the set's plans differ from their indexes among the branches.
type unrecordedSuffixSet struct {
	removableSet
	a    ambient
	set  explicitCleanupSet
	tips map[string]string
}

func unrecordedSuffixFixture(t *testing.T, options CleanupOptions) unrecordedSuffixSet {
	t.Helper()
	f := unrecordedSuffixSet{removableSet: removableSetFixture(t, 1), tips: map[string]string{}}
	targets := explicitIdentities(f.creations)
	for _, owner := range []string{"a", "b", "c"} {
		ref, tip := uniqueBranch(t, f.root, owner)
		f.tips[ref] = tip
		targets = append(targets, branchTargets(ref)...)
	}
	f.a = baselineAmbient(f.repoHome)
	f.set = planExplicitSet(defaultJoins(), f.a, f.root, targets, options)
	if f.set.fingerprint == "" || len(f.set.rows) != 1 || len(f.set.unrecorded) != 3 {
		t.Fatalf("explicit set = %#v, want one recorded and three unrecorded members", f.set)
	}
	return f
}

// TestCleanupTraversalUnrecordedSuffix is R11, R13, R14, R16, and R19 for the unrecorded stage.
// A stop at the second branch keeps the recorded member and the first branch removed, names the
// stopped branch in its own row, reports the unreached branch as not attempted, and leaves the
// stopped and unreached branches at their tips with no discarded ref.
func TestCleanupTraversalUnrecordedSuffix(t *testing.T) {
	t.Parallel()
	options := CleanupOptions{DiscardBranch: true}
	stop := errors.New("stop at the second branch")
	for _, tc := range []struct {
		name string
		// fault installs the stop. It returns the error the apply must return and the tip the
		// stopped branch must keep.
		fault func(*testing.T, *joins, unrecordedSuffixSet) (error, func() string)
		// stopped checks the stopped branch's own row against its planned row.
		stopped func(planned, got CleanupPlan) bool
		// writes is how many visits reach a discarded-ref write. A fault stops inside the second
		// write, and drift stops the second visit before its write.
		writes int
	}{
		{"R11 a fault in the second discard", func(_ *testing.T, j *joins, f unrecordedSuffixSet) (error, func() string) {
			j.cleanupBoundary = atSecondHit(StepDiscardedRefWrite, func() error { return stop })
			ref := f.set.unrecorded[1].ref
			return stop, func() string { return f.tips[ref] }
		}, func(planned, got CleanupPlan) bool {
			return got.Target == planned.Target && got.Action == ActionError && got.Reason == stop.Error() && got.Recovery == planned.Recovery
		}, 2},
		{"R13 the second branch moves and stays removable", func(t *testing.T, j *joins, f unrecordedSuffixSet) (error, func() string) {
			ref, moved := f.set.unrecorded[1].ref, ""
			j.cleanupBoundary = atStep(StepDiscardedBranchDelete, func() error {
				if moved == "" {
					moved = commitOnBranch(t, f.root, ref, "moved.txt", "moved\n")
				}
				return nil
			})
			return errStaleFingerprint, func() string { return moved }
		}, func(planned, got CleanupPlan) bool {
			return got.Target == planned.Target && got.Action == ActionNotAttempted && got.Reason == driftedDetail
		}, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := unrecordedSuffixFixture(t, options)
			j := defaultJoins()
			want, stoppedTip := tc.fault(t, &j, f)
			steps, trace := stepTrace()
			inner := j.cleanupBoundary
			j.cleanupBoundary = func(step LifecycleStep) error { _ = trace(step); return inner(step) }

			plans, err := applyExplicitSet(j, f.a, f.root, f.set, options)
			if !errors.Is(err, want) {
				t.Fatalf("apply error = %v, want %v", err, want)
			}
			members := f.set.plans()
			if len(plans) != len(members) {
				t.Fatalf("rows = %#v, want one row per selected member", plans)
			}
			first, second, third := f.set.unrecorded[0], f.set.unrecorded[1], f.set.unrecorded[2]
			if plans[0].Action != ActionRemoved || plans[1].Action != ActionRemoved || plans[1].Target != first.plan.Target || plans[1].Recovery != first.discarded {
				t.Fatalf("rows = %#v, want the recorded member and the first branch removed", plans)
			}
			if refTip(f.root, first.ref) != "" || refTip(f.root, first.discarded) != f.tips[first.ref] {
				t.Fatalf("the first branch %s reads %q and its discarded ref %q, want it gone under %s", first.ref, refTip(f.root, first.ref), refTip(f.root, first.discarded), f.tips[first.ref])
			}
			if !tc.stopped(second.plan, plans[2]) {
				t.Fatalf("stopped row = %#v, want the outcome of %s", plans[2], second.ref)
			}
			if got := plans[3]; got.Target != third.plan.Target || got.Action != ActionNotAttempted || got.Reason != notAttemptedDetail {
				t.Fatalf("unreached row = %#v, want %s not attempted", got, third.ref)
			}
			if refTip(f.root, second.ref) != stoppedTip() || refTip(f.root, third.ref) != f.tips[third.ref] || refTip(f.root, second.discarded) != "" || refTip(f.root, third.discarded) != "" {
				t.Fatalf("refs = %q, want the stopped and unreached branches kept with no discarded ref", discardedRefs(t, f.root))
			}
			writes := 0
			for _, step := range *steps {
				if step == StepDiscardedRefWrite {
					writes++
				}
			}
			if writes != tc.writes {
				t.Fatalf("boundaries = %v, want no discard visit after the stopped branch", *steps)
			}
		})
	}
}

// TestCleanupTraversalUnrecordedMixed is R04, R05, and R06 for the unrecorded stage. Without
// --discard-branch the branches pass through as planned behind a removed recorded member and
// reach no discard boundary. With it, a branch that moved before the apply refuses the set in
// its preflight, so the recorded member ahead of it stays.
func TestCleanupTraversalUnrecordedMixed(t *testing.T) {
	t.Parallel()
	t.Run("R04 R05 retained branches", func(t *testing.T) {
		t.Parallel()
		f := unrecordedSuffixFixture(t, CleanupOptions{})
		j := defaultJoins()
		steps, trace := stepTrace()
		j.cleanupBoundary = trace

		plans, err := applyExplicitSet(j, f.a, f.root, f.set, CleanupOptions{})
		mustNoError(t, err)
		if len(plans) != 4 || plans[0].Action != ActionRemoved || !reflect.DeepEqual(plans[1:], f.set.plans()[1:]) {
			t.Fatalf("rows = %#v, want the recorded member removed and each branch as planned", plans)
		}
		for _, step := range *steps {
			if step == StepDiscardedRefWrite || step == StepDiscardedRefAbsent || step == StepDiscardedBranchDelete {
				t.Fatalf("boundaries = %v, want no discard visit of a retained branch", *steps)
			}
		}
		for ref, tip := range f.tips {
			if refTip(f.root, ref) != tip {
				t.Fatalf("retained branch %s reads %q, want %s", ref, refTip(f.root, ref), tip)
			}
		}
	})
	t.Run("R06 a branch moved before the apply", func(t *testing.T) {
		t.Parallel()
		options := CleanupOptions{DiscardBranch: true}
		f := unrecordedSuffixFixture(t, options)
		moved := commitOnBranch(t, f.root, f.set.unrecorded[2].ref, "moved.txt", "moved\n")

		plans, err := applyExplicitSet(defaultJoins(), f.a, f.root, f.set, options)
		if !errors.Is(err, errStaleFingerprint) {
			t.Fatalf("apply error = %v, want the preflight to refuse", err)
		}
		requireMembersPresent(t, f.creations)
		for _, plan := range plans {
			if plan.Action == ActionRemoved {
				t.Fatalf("rows = %#v, want no removal", plans)
			}
		}
		if discardedRefs(t, f.root) != "" || refTip(f.root, f.set.unrecorded[2].ref) != moved {
			t.Fatalf("refs = %q, want no discarded ref and the moved branch at %s", discardedRefs(t, f.root), moved)
		}
	})
}
