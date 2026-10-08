package repository

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/refusalroute"
)

// Start binds the owned assignment and its deliverable under the intent lock.
func (store Store) Start(outcome, request, deliverable string) error {
	return intent.Transact(store.Root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
		owner, valid := requestedAssignment(ledger, store.Root, request)
		if !valid {
			return ledger, false, errors.New("start requires the active owned assignment and its exact request")
		}
		policy, err := store.admissionPolicy()
		if err != nil {
			return ledger, false, err
		}
		selected, err := commitment.ActiveOutcome(policy, outcome)
		if err != nil {
			return ledger, false, err
		}
		var source *commitment.SourceBinding
		for _, approved := range selected.Deliverables {
			if approved.Source.Path == deliverable {
				value := approved.Source
				source = &value
				break
			}
		}
		if source == nil {
			return ledger, false, fmt.Errorf("deliverable %q is not approved for outcome %q", deliverable, outcome)
		}
		if err := store.validateDeliverable(*source); err != nil {
			return ledger, false, err
		}
		binding := intent.DeliveryBinding{Assignment: owner.ID, Request: owner.Request, Milestone: policy.ActiveMilestone, Outcome: outcome, Deliverable: deliverable, Identity: source.Identity}
		next, err := commitment.Admit(policy, runtimeState(ledger), binding)
		if err != nil {
			return ledger, false, err
		}
		return withRuntime(ledger, next)
	}, nil)
}

// Block preserves an obligation and releases its active claim.
func (store Store) Block(outcome, reason string) error {
	return store.setBlocker(outcome, reason, true)
}

// Unblock restores eligibility without preempting an active outcome.
func (store Store) Unblock(outcome string) error { return store.setBlocker(outcome, "", false) }

func (store Store) setBlocker(outcome, reason string, blocked bool) error {
	return intent.Transact(store.Root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
		policy, err := store.admissionPolicy()
		if err != nil {
			return ledger, false, err
		}
		next, err := commitment.SetBlocker(policy, runtimeState(ledger), outcome, reason, blocked)
		if err != nil {
			return ledger, false, err
		}
		return withRuntime(ledger, next)
	}, nil)
}

func (store Store) admissionPolicy() (commitment.Policy, error) {
	policy, exists, err := store.Policy()
	if err != nil {
		return policy, err
	}
	if !exists {
		return policy, refusalroute.Raised{Name: refusalroute.CommitmentDecision, Err: errors.New("commitment adoption required")}
	}
	if strings.TrimSpace(policy.ActiveMilestone) == "" {
		return policy, errors.New("no milestone is active")
	}
	return policy, nil
}

// runtimeState is the commitment state that every decision reads. A continuation lasts
// only while its run is active, so it drops each continuation whose run is complete,
// recovered, cleanup-pending, or absent. commitment.OpenContinuations then drops each
// delivered scope. A write of the result persists the drop.
func runtimeState(ledger intent.Ledger) intent.CommitmentState {
	state := storedState(ledger)
	if len(state.Continuations) == 0 {
		return state
	}
	state.Continuations = slices.DeleteFunc(slices.Clone(state.Continuations), func(continuation intent.LegacyContinuation) bool {
		return !slices.ContainsFunc(ledger.Assignments, func(run intent.Assignment) bool {
			return run.ID == continuation.Assignment && run.State == intent.StateActive
		})
	})
	if len(state.Continuations) == 0 {
		state.Continuations = nil
	}
	return state
}

// storedState is the commitment state as the ledger holds it.
func storedState(ledger intent.Ledger) intent.CommitmentState {
	if ledger.Commitment != nil {
		return *ledger.Commitment
	}
	return intent.CommitmentState{}
}

func withRuntime(ledger intent.Ledger, next intent.CommitmentState) (intent.Ledger, bool, error) {
	changed := !reflect.DeepEqual(storedState(ledger), next)
	if changed {
		ledger.Commitment = &next
	}
	return ledger, changed, nil
}

func requestedAssignment(ledger intent.Ledger, root, request string) (intent.Assignment, bool) {
	owners := intent.AssignmentsOwning(ledger.Assignments, root)
	if len(owners) != 1 || owners[0].State != intent.StateActive || request == "" || owners[0].Request != intent.RequestDigest(request) {
		return intent.Assignment{}, false
	}
	return owners[0], true
}
