package repository

import (
	"fmt"
	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/intent"
)

// Ready verifies the owned assignment against current delivery authority without creating a claim.
func (store Store) Ready(deliverable string) error {
	ledger, err := intent.Read(store.Root)
	if err != nil {
		return err
	}
	return store.ready(ledger, deliverable)
}

func (store Store) ready(ledger intent.Ledger, deliverable string) error {
	owners := intent.AssignmentsOwning(ledger.Assignments, store.Root)
	if len(owners) != 1 || owners[0].State != intent.StateActive {
		return fmt.Errorf("delivery requires an active owned assignment; run bench worktree create")
	}
	policy, err := store.admissionPolicy()
	if err != nil {
		return err
	}
	state := runtimeState(ledger)
	for _, binding := range state.Bindings {
		if binding.Assignment != owners[0].ID || binding.Request != owners[0].Request {
			continue
		}
		if binding.Milestone != policy.ActiveMilestone || (deliverable != "" && binding.Deliverable != deliverable) {
			break
		}
		outcome, err := commitment.ActiveOutcome(policy, binding.Outcome)
		if err != nil {
			return err
		}
		approved := false
		for _, source := range outcome.Deliverables {
			if source.Source.Path == binding.Deliverable && source.Source.Identity == binding.Identity {
				if err := store.validateDeliverable(source.Source); err != nil {
					return err
				}
				approved = true
			}
		}
		if !approved {
			break
		}
		claimed := false
		for _, claim := range state.Claims {
			if claim.Milestone == binding.Milestone && claim.Outcome == binding.Outcome {
				claimed = true
			}
		}
		if !claimed {
			break
		}
		if _, err := commitment.Admit(policy, state, binding); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("assignment has no current delivery binding; run bench commitment start --outcome <id> --request <request> --deliverable <path>")
}
