package repository

import (
	"fmt"
	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/intent/admissionpolicy"
	"slices"
)

// Ready verifies the owned assignment against current delivery authority without creating a claim.
func (store Store) Ready(deliverable string) error {
	ledger, err := intent.Read(store.Root)
	if err != nil {
		return err
	}
	return store.ready(ledger, deliverable, "")
}

func (store Store) ready(ledger intent.Ledger, deliverable, outcomeID string) error {
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
		if binding.Milestone != policy.ActiveMilestone || (deliverable != "" && binding.Deliverable != deliverable) || (outcomeID != "" && binding.Outcome != outcomeID) {
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
		if _, err := commitment.Admit(policy, state, binding); err != nil {
			return err
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
		return nil
	}
	return fmt.Errorf("assignment has no current delivery binding; run bench commitment start --outcome <id> --request <request> --deliverable <path>")
}

// Inheritance checks the source assignment before a sibling acquires its worktree.
func (store Store) Inheritance(assignment string) (*intent.DeliveryBinding, error) {
	ledger, err := intent.Read(store.Root)
	if err != nil {
		return nil, err
	}
	return store.inheritance(ledger, assignment)
}

func (store Store) inheritance(ledger intent.Ledger, assignment string) (*intent.DeliveryBinding, error) {
	for _, owner := range ledger.Assignments {
		if owner.ID != assignment {
			continue
		}
		if owner.State != intent.StateActive {
			return nil, fmt.Errorf("source assignment is not active")
		}
		for _, binding := range runtimeState(ledger).Bindings {
			if binding.Assignment != assignment {
				continue
			}
			if err := (Store{Root: owner.Worktree}).ready(ledger, binding.Deliverable, ""); err != nil {
				return nil, err
			}
			return &binding, nil
		}
		return nil, nil
	}
	return nil, fmt.Errorf("source assignment is unknown")
}

// RegisterSibling records ownership and its inherited binding in one transaction.
func (store Store) RegisterSibling(assignment intent.Assignment, expected intent.DeliveryBinding) error {
	if err := intent.ValidateAssignment(assignment); err != nil {
		return err
	}
	return intent.Transact(store.Root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
		inherited, err := store.inheritance(ledger, expected.Assignment)
		if err != nil {
			return ledger, false, err
		}
		if inherited == nil || *inherited != expected {
			return ledger, false, fmt.Errorf("source delivery binding changed during creation")
		}
		policy, err := store.admissionPolicy()
		if err != nil {
			return ledger, false, err
		}
		inherited.Assignment, inherited.Request = assignment.ID, assignment.Request
		next, err := commitment.Admit(policy, runtimeState(ledger), *inherited)
		if err != nil {
			return ledger, false, err
		}
		registered, changed, err := admissionpolicy.PutAssignment(ledger, assignment)
		if err != nil {
			return ledger, false, err
		}
		registered, bound, err := withRuntime(registered, next)
		return registered, changed || bound, err
	}, nil)
}

// RollbackSibling removes only the failed creation's records and retains its outcome claim.
func (store Store) RollbackSibling(assignment intent.Assignment) error {
	return intent.Transact(store.Root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
		for _, current := range ledger.Assignments {
			if current.ID == assignment.ID && current.Request != assignment.Request {
				return ledger, false, fmt.Errorf("failed sibling request changed")
			}
		}
		next, changed, err := admissionpolicy.DeleteAssignment(ledger, assignment.ID)
		if err != nil {
			return ledger, false, err
		}
		state := runtimeState(next)
		state.Bindings = slices.DeleteFunc(slices.Clone(state.Bindings), func(binding intent.DeliveryBinding) bool {
			return binding.Assignment == assignment.ID && binding.Request == assignment.Request
		})
		next, unbound, err := withRuntime(next, state)
		return next, changed || unbound, err
	}, nil)
}

// ReadyOutcome checks the current assignment's existing authority for a shift.
func (store Store) ReadyOutcome(outcome string) error {
	if outcome == "" {
		return fmt.Errorf("shift requires --outcome <id>")
	}
	ledger, err := intent.Read(store.Root)
	if err != nil {
		return err
	}
	return store.ready(ledger, "", outcome)
}

// LegacyScope returns only the current assignment's exact approved continuation scope.
func (store Store) LegacyScope(request string) ([]string, error) {
	ledger, err := intent.Read(store.Root)
	if err != nil {
		return nil, err
	}
	owner, valid := requestedAssignment(ledger, store.Root, request)
	if !valid {
		return nil, fmt.Errorf("continuation requires the active owned assignment and its exact request")
	}
	for _, continuation := range runtimeState(ledger).Continuations {
		if continuation.Assignment == owner.ID && continuation.Request == owner.Request {
			return append([]string(nil), continuation.Scope...), nil
		}
	}
	return nil, fmt.Errorf("legacy continuation is not approved for this assignment")
}
