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

// TestCleanupTraversalEmptyModes is R01. A selection with no removable member reaches no
// lifecycle boundary, changes no ref or record, and reports each retained row as planned.
func TestCleanupTraversalEmptyModes(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name  string
		build func(*testing.T) (repoHome, []string)
	}{
		{"explicit/empty", func(t *testing.T) (repoHome, []string) {
			root := newWorktreeRepo(t)
			return repoHome{root, filepath.Join(root, ".bench-home")}, nil
		}},
		{"explicit/retained-only", func(t *testing.T) (repoHome, []string) {
			f := retainedMemberFixture(t)
			return f.repoHome, []string{f.retained.Assignment.ID}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f, targets := tc.build(t)
			j := defaultJoins()
			set := planExplicitSet(j, f.ambient(), f.root, targets, CleanupOptions{})
			steps, trace := stepTrace()
			j.cleanupBoundary = trace
			before := strings.Join(baselineRefs(t, f.root), "\n") + repositoryState(t, f.root)

			plans, err := applyExplicitSet(j, f.ambient(), f.root, set, CleanupOptions{})
			mustNoError(t, err)
			if len(*steps) != 0 {
				t.Fatalf("apply reached %v, want no lifecycle boundary", *steps)
			}
			if want := set.plans(); len(plans) != len(want) || (len(want) != 0 && !reflect.DeepEqual(plans, want)) {
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
	t.Run("explicit/all-success", func(t *testing.T) {
		t.Parallel()
		f := removableSetFixture(t, 2)
		j := defaultJoins()
		set := planExplicitSet(j, f.ambient(), f.root, explicitIdentities(f.creations), CleanupOptions{})
		steps, trace := stepTrace()
		j.cleanupBoundary = trace

		_, err := applyExplicitSet(j, f.ambient(), f.root, set, CleanupOptions{})
		mustNoError(t, err)
		member := []LifecycleStep{StepMemberRequalify, StepApplyLocked, StepReceipt, StepUnlock, StepRemovalAttempt, StepRemoval, StepBranch}
		if want := append(append([]LifecycleStep{}, member...), member...); !reflect.DeepEqual(*steps, want) {
			t.Fatalf("boundaries = %v, want %v", *steps, want)
		}
	})
}

// TestCleanupTraversalAuthorityPreservation is R28. Ignored residue that appears inside a
// member's own lock takes its removal away: the locked recheck refuses, the member and its
// residue stay, and the retained member behind it keeps its planned verdict.
func TestCleanupTraversalAuthorityPreservation(t *testing.T) {
	t.Parallel()
	f := retainedMemberFixture(t)
	j := defaultJoins()
	set := planExplicitSet(j, f.ambient(), f.root, []string{f.removable.Assignment.ID, f.retained.Assignment.ID}, CleanupOptions{})
	if set.fingerprint == "" || len(set.rows) != 2 {
		t.Fatalf("explicit set = %#v, want two applicable members", set.rows)
	}
	// The set orders its members by identity, so the removable member's row index is read
	// from the plan rather than assumed.
	removable, retained := 0, 1
	if set.rows[0].assignment.ID != f.removable.Assignment.ID {
		removable, retained = 1, 0
	}
	if !set.rows[removable].plan.Action.Removes() || set.rows[retained].plan.Action.Removes() {
		t.Fatalf("explicit set = %#v, want one removable and one retained member", set.rows)
	}
	residue := filepath.Join(f.removable.Path, "ignored-late.txt")
	j.cleanupBoundary = atStep(StepApplyLocked, func() error {
		mustWrite(t, residue, []byte("residue\n"), 0o644)
		return nil
	})

	plans, err := applyExplicitSet(j, f.ambient(), f.root, set, CleanupOptions{})
	if !errors.Is(err, errStaleFingerprint) {
		t.Fatalf("apply error = %v, want the locked recheck to refuse", err)
	}
	if _, statErr := os.Stat(residue); statErr != nil {
		t.Fatalf("the preserved residue %s was removed: %v", residue, statErr)
	}
	requireMembersPresent(t, []Creation{f.removable, f.retained})
	if len(plans) != 2 || plans[removable].Action != ActionError || !reflect.DeepEqual(plans[retained], set.rows[retained].plan) {
		t.Fatalf("rows = %#v, want the refused member failed and the retained member as planned", plans)
	}
}
