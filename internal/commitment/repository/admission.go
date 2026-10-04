package repository

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/intent"
)

// Start binds the owned assignment and its deliverable under the intent lock.
func (store Store) Start(outcome, request, deliverable string) error {
	return intent.Transact(store.Root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
		owners := intent.AssignmentsOwning(ledger.Assignments, store.Root)
		if len(owners) != 1 || owners[0].State != intent.StateActive || request == "" || owners[0].Request != intent.RequestDigest(request) {
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
		binding := intent.DeliveryBinding{Assignment: owners[0].ID, Request: owners[0].Request, Milestone: policy.ActiveMilestone, Outcome: outcome, Deliverable: deliverable, Identity: source.Identity}
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
		return policy, errors.New("commitment adoption required: run bench commitment plan --input <file>")
	}
	if strings.TrimSpace(policy.ActiveMilestone) == "" {
		return policy, errors.New("no milestone is active")
	}
	return policy, nil
}

func runtimeState(ledger intent.Ledger) intent.CommitmentState {
	if ledger.Commitment != nil {
		return *ledger.Commitment
	}
	return intent.CommitmentState{}
}

func withRuntime(ledger intent.Ledger, next intent.CommitmentState) (intent.Ledger, bool, error) {
	changed := !reflect.DeepEqual(runtimeState(ledger), next)
	if changed {
		ledger.Commitment = &next
	}
	return ledger, changed, nil
}
