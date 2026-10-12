package repository

import (
	"fmt"
	"slices"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// OverrideTarget is what the reviewer confirms before an override. It holds the label of the
// assignment, and the outcome of each open build that refuses the assignment's landing.
type OverrideTarget struct {
	Label  string
	Builds []string
}

// OverrideTarget reads the target of an override for assignment from one ledger read. It
// refuses an id that names no active assignment. The landing it reads is the one that closes
// the assignment's bound deliverable, or a spec-less one when the assignment has no binding.
func (store Store) OverrideTarget(assignment string) (OverrideTarget, error) {
	if !intent.ValidIdentity(assignment) {
		return OverrideTarget{}, notOverridable(assignment)
	}
	ledger, err := intent.Read(store.Root)
	if err != nil {
		return OverrideTarget{}, err
	}
	owner, active := activeByID(ledger, assignment)
	if !active {
		return OverrideTarget{}, notOverridable(assignment)
	}
	source := Publication{Assignment: assignment}
	for _, binding := range storedState(ledger).Bindings {
		if binding.Assignment == assignment {
			source.Deliverable = binding.Deliverable
		}
	}
	refusing, err := store.refusingBuilds(ledger, source)
	if err != nil {
		return OverrideTarget{}, err
	}
	target := OverrideTarget{Label: owner.Label}
	for _, build := range refusing {
		target.Builds = append(target.Builds, build.claim.Outcome)
	}
	return target, nil
}

// Override records the reviewer's override for assignment with reason in one decision under
// the intent lock. A later grant for the same assignment replaces the reason of the first.
func (store Store) Override(assignment, reason string) error {
	if err := commitment.ValidReason(reason); err != nil {
		return err
	}
	return intent.Transact(store.Root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
		if _, active := activeByID(ledger, assignment); !active {
			return ledger, false, notOverridable(assignment)
		}
		state := storedState(ledger)
		state.Overrides = append(withoutOverride(state.Overrides, assignment), intent.LandingOverride{Assignment: assignment, Reason: reason})
		return withRuntime(ledger, state)
	}, nil)
}

// holdsOverride reports whether overrides admit the landing of assignment.
func holdsOverride(overrides []intent.LandingOverride, assignment string) bool {
	return slices.ContainsFunc(overrides, overrideOf(assignment))
}

// withoutOverride drops the override of assignment from overrides.
func withoutOverride(overrides []intent.LandingOverride, assignment string) []intent.LandingOverride {
	return slices.DeleteFunc(slices.Clone(overrides), overrideOf(assignment))
}

// overrideOf matches the override record of assignment.
func overrideOf(assignment string) func(intent.LandingOverride) bool {
	return func(override intent.LandingOverride) bool { return override.Assignment == assignment }
}

// notOverridable is the refusal of an override for an id that names no active assignment.
func notOverridable(assignment string) error {
	return fmt.Errorf("commitment override assignment %s is not active", sanitize.Controls(assignment))
}
