// Which member a refused apply names, and what that member's row says. A refusal has to land
// on the member that caused it and report what that member is now: a fault names the row whose
// requalification raised it, a member that drifted to retained still reports retained, and one
// still marked removable reports that its removal never started.
//
// The clean_set_*_test.go files split by the question each answers; this one answers which
// member.
package worktree

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

// TestCleanSetPreflightFaultNamesItsMember is COV-11. A requalify can fail for a reason that
// is not drift, when a repository read the re-plan depends on breaks under it. That fault is
// not a property of the approved set, so it belongs to the member that raised it: that row
// carries the reason, and the rest report only that the set never qualified.
func TestCleanSetPreflightFaultNamesItsMember(t *testing.T) {
	t.Parallel()
	root, _, creations := removableSetFixture(t, 2)
	j := defaultJoins()
	set, planErr := planLandedSet(j, root, CleanupOptions{}, "")
	mustNoError(t, planErr)
	if len(set.rows) != 2 {
		t.Fatalf("landed set = %#v, want two applicable members", set.rows)
	}
	// The ledger the re-plan reads becomes unreadable after the plan. Every member is still
	// exactly as approved, so no drift explains the refusal.
	ledger, err := intent.Address(root)
	mustNoError(t, err)
	mustWrite(t, ledger, []byte("{not a ledger\n"), 0o644)

	plans, applyErr := applyLandedSet(j, root, set, CleanupOptions{}, "")
	if applyErr == nil || errors.Is(applyErr, errStaleFingerprint) {
		t.Fatalf("apply error = %v, want a requalify fault that is not drift", applyErr)
	}
	requireMembersPresent(t, creations)
	if len(plans) != 2 {
		t.Fatalf("refused rows = %#v, want one row per member", plans)
	}
	if plans[0].Action != ActionError || plans[0].Reason == "" {
		t.Fatalf("offending row = %q/%q, want the fault named on its own member", plans[0].Action, plans[0].Reason)
	}
	if plans[0].Reason == errStaleFingerprint.Error() {
		t.Fatalf("offending row detail = %q, want the fault's own reason", plans[0].Reason)
	}
	// The detail literals are agent-facing text these rows promise. Reading one back through
	// the constant that produces it would let a rewrite pass the gate.
	if plans[1].Action != ActionNotAttempted || plans[1].Reason != "not attempted; the set did not qualify before the first removal" {
		t.Fatalf("unstarted row = %q/%q, want the unqualified detail", plans[1].Action, plans[1].Reason)
	}
}

// TestCleanSetMemberDriftAfterPreflight covers the window StepMemberRequalify names: the
// preflight cleared every member, an earlier member's transaction ran, and only then did this
// member change. No transaction opens on it, so its row says what it is now rather than
// claiming a removal that never started.
func TestCleanSetMemberDriftAfterPreflight(t *testing.T) {
	t.Parallel()
	t.Run("drifts to retained", func(t *testing.T) {
		t.Parallel()
		root, home := ignoringRepo(t)
		first := landedMember(t, root, home, "drift-retained-first", "first.txt")
		second := landedMember(t, root, home, "drift-retained-second", "second.txt")
		j := defaultJoins()
		set := planExplicitSet(j, root, explicitIdentities([]Creation{first, second}), CleanupOptions{})
		if set.fingerprint == "" || len(set.rows) != 2 {
			t.Fatalf("explicit set = %#v, want two applicable members", set)
		}
		drifted := memberByID(t, []Creation{first, second}, set.rows[1].assignment.ID)
		j.cleanupBoundary = driftAtSecondRequalify(t, func() {
			mustWrite(t, filepath.Join(drifted.Path, "ignored-late.txt"), []byte("residue\n"), 0o644)
		})

		plans, err := applyExplicitSet(j, root, set, CleanupOptions{})
		if !errors.Is(err, errStaleFingerprint) {
			t.Fatalf("apply error = %v, want the member's own requalify to refuse", err)
		}
		if len(plans) != 2 || plans[1].Action != ActionRetain {
			t.Fatalf("refused rows = %#v, want the drifted member reported as retained", plans)
		}
		if _, statErr := os.Stat(drifted.Path); statErr != nil {
			t.Fatalf("the retained member %s was removed: %v", drifted.Path, statErr)
		}
	})
	t.Run("drifts while still removable", func(t *testing.T) {
		t.Parallel()
		root, _, creations := removableSetFixture(t, 2)
		j := defaultJoins()
		set, planErr := planLandedSet(j, root, CleanupOptions{}, "")
		mustNoError(t, planErr)
		if len(set.rows) != 2 {
			t.Fatalf("landed set = %#v, want two applicable members", set.rows)
		}
		// The record the selector reads disappears, so the member leaves the selection while
		// the approved row still says remove. That row must not report a removal.
		drifted := memberByID(t, creations, set.rows[1].assignment.ID)
		j.cleanupBoundary = driftAtSecondRequalify(t, func() {
			mustNoError(t, intent.DeleteAssignment(root, drifted.Assignment.ID))
		})

		plans, err := applyLandedSet(j, root, set, CleanupOptions{}, "")
		if !errors.Is(err, errStaleFingerprint) {
			t.Fatalf("apply error = %v, want the member's own requalify to refuse", err)
		}
		if len(plans) != 2 || plans[1].Action != ActionNotAttempted || plans[1].Reason != "not attempted; this target no longer matches the approved plan" {
			t.Fatalf("refused rows = %#v, want the drifted member reported as not attempted", plans)
		}
		requireMembersPresent(t, []Creation{drifted})
	})
	t.Run("faults while requalifying", func(t *testing.T) {
		t.Parallel()
		root, _, creations := removableSetFixture(t, 2)
		j := defaultJoins()
		set, planErr := planLandedSet(j, root, CleanupOptions{}, "")
		mustNoError(t, planErr)
		if len(set.rows) != 2 {
			t.Fatalf("landed set = %#v, want two applicable members", set.rows)
		}
		// The preflight cleared every member, the first removal completed, and only then does
		// the ledger the next requalify reads break. That fault is not drift, so the row it
		// belongs to carries its reason.
		survivor := memberByID(t, creations, set.rows[1].assignment.ID)
		j.cleanupBoundary = driftAtSecondRequalify(t, func() {
			ledger, err := intent.Address(root)
			mustNoError(t, err)
			mustWrite(t, ledger, []byte("{not a ledger\n"), 0o644)
		})

		plans, err := applyLandedSet(j, root, set, CleanupOptions{}, "")
		if err == nil || errors.Is(err, errStaleFingerprint) {
			t.Fatalf("apply error = %v, want an in-loop requalify fault that is not drift", err)
		}
		if len(plans) != 2 || plans[1].Action != ActionError || plans[1].Reason == errStaleFingerprint.Error() {
			t.Fatalf("refused rows = %#v, want the fault named on the member that raised it", plans)
		}
		requireMembersPresent(t, []Creation{survivor})
	})
}

// driftAtSecondRequalify runs change in the window before the second member requalifies, so
// the first member's transaction has already completed when the drift lands.
func driftAtSecondRequalify(t *testing.T, change func()) Fault {
	t.Helper()
	windows := 0
	return func(step LifecycleStep) error {
		if step != StepMemberRequalify {
			return nil
		}
		windows++
		if windows == 2 {
			change()
		}
		return nil
	}
}

// TestCleanSetPreflightFaultOnLaterMember is COV-13. The preflight skips a member the plan
// retained, so the member whose requalification faults need not be the first row. The fault
// belongs to the row that raised it, and the retained row ahead of it keeps its own verdict.
func TestCleanSetPreflightFaultOnLaterMember(t *testing.T) {
	t.Parallel()
	root, _, creations := removableSetFixture(t, 2)
	j := defaultJoins()
	ordered, planErr := planLandedSet(j, root, CleanupOptions{}, "")
	mustNoError(t, planErr)
	if len(ordered.rows) != 2 {
		t.Fatalf("landed set = %#v, want two members", ordered.rows)
	}
	// The member the set reaches first becomes dirty, so the plan retains it and the preflight
	// skips it. The fault then lands on the member at index one.
	leading := memberByID(t, creations, ordered.rows[0].assignment.ID)
	mustWrite(t, filepath.Join(leading.Path, "member-0.txt"), []byte("dirty\n"), 0o644)
	mustWrite(t, filepath.Join(leading.Path, "member-1.txt"), []byte("dirty\n"), 0o644)
	set, replanErr := planLandedSet(j, root, CleanupOptions{}, "")
	mustNoError(t, replanErr)
	if len(set.rows) != 2 || set.rows[0].plan.Action.Removes() || !set.rows[1].plan.Action.Removes() {
		t.Fatalf("landed set = %#v, want a retained row ahead of a removable one", set.rows)
	}
	ledger, addressErr := intent.Address(root)
	mustNoError(t, addressErr)
	mustWrite(t, ledger, []byte("{not a ledger\n"), 0o644)

	plans, applyErr := applyLandedSet(j, root, set, CleanupOptions{}, "")
	if applyErr == nil || errors.Is(applyErr, errStaleFingerprint) {
		t.Fatalf("apply error = %v, want a requalify fault that is not drift", applyErr)
	}
	requireMembersPresent(t, creations)
	if len(plans) != 2 || plans[0].Action != ActionRetain {
		t.Fatalf("refused rows = %#v, want the retained row ahead to keep its verdict", plans)
	}
	if plans[1].Action != ActionError || plans[1].Reason == errStaleFingerprint.Error() {
		t.Fatalf("offending row = %q/%q, want the fault named on the later member", plans[1].Action, plans[1].Reason)
	}
}

// TestCleanSetExplicitPreflightFaultOnLaterMember is COV-17. The explicit mode's preflight
// carries its own offender index, and it skips a member the plan retained, so the member that
// faults need not be the first row. The landed mode's index is covered above; this holds the
// explicit one.
func TestCleanSetExplicitPreflightFaultOnLaterMember(t *testing.T) {
	t.Parallel()
	root, home := ignoringRepo(t)
	first := landedMember(t, root, home, "explicit-fault-first", "first.txt")
	second := landedMember(t, root, home, "explicit-fault-second", "second.txt")
	members := []Creation{first, second}
	j := defaultJoins()
	ordered := planExplicitSet(j, root, explicitIdentities(members), CleanupOptions{})
	if len(ordered.rows) != 2 {
		t.Fatalf("explicit set = %#v, want two members", ordered.rows)
	}
	// The member the set reaches first gains ignored residue, so the plan retains it and the
	// preflight skips it. The fault then lands on the member at index one.
	leading := memberByID(t, members, ordered.rows[0].assignment.ID)
	mustWrite(t, filepath.Join(leading.Path, "ignored-one.txt"), []byte("residue\n"), 0o644)
	set := planExplicitSet(j, root, explicitIdentities(members), CleanupOptions{})
	if set.fingerprint == "" || set.rows[0].plan.Action.Removes() || !set.rows[1].plan.Action.Removes() {
		t.Fatalf("explicit set = %#v, want a retained row ahead of a removable one", set.rows)
	}
	// The repository the re-plan reads goes away, so the requalify faults for a reason that is
	// not drift. Every member is still exactly as approved.
	mustNoError(t, os.RemoveAll(filepath.Join(root, ".git")))

	plans, applyErr := applyExplicitSet(j, root, set, CleanupOptions{})
	if applyErr == nil || errors.Is(applyErr, errStaleFingerprint) {
		t.Fatalf("apply error = %v, want a requalify fault that is not drift", applyErr)
	}
	if len(plans) != 2 || plans[0].Action != ActionRetain {
		t.Fatalf("refused rows = %#v, want the retained row ahead to keep its verdict", plans)
	}
	if plans[1].Action != ActionError || plans[1].Reason == errStaleFingerprint.Error() {
		t.Fatalf("offending row = %q/%q, want the fault named on the later member", plans[1].Action, plans[1].Reason)
	}
}
