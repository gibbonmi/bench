package commitment

import (
	"fmt"
	"slices"
)

// Deliver records the verified delivery of the approved deliverable at path. It returns
// the next policy and the outcome sources that the delivery completely satisfies. An
// outcome with no source is rowless: each of its approved deliverables is an obligation of
// the outcome itself, so its delivery records a fact and satisfies no source. A path that
// the active milestone does not approve, or a binding that names no obligation of an
// outcome with sources, returns policy unchanged and no source.
func Deliver(policy Policy, path, source, evidence string) (Policy, []SourceBinding, error) {
	outcome, binding, found, err := activeDeliverable(policy, path)
	if err != nil || !found || obligationFree(outcome, binding) {
		return policy, nil, err
	}
	fact := DeliveryFact{Milestone: policy.ActiveMilestone, Outcome: outcome.ID, Binding: binding.Source.ID, Identity: binding.Source.Identity, Source: source, Evidence: evidence}
	if deliveredKeys(policy)[fact.key()] {
		return policy, nil, fmt.Errorf("deliverable %q is already delivered", path)
	}
	var closed []SourceBinding
	for _, candidate := range outcome.Sources {
		if slices.Contains(binding.Obligations, candidate.ID) {
			closed = append(closed, candidate)
		}
	}
	policy.Deliveries = append(slices.Clone(policy.Deliveries), fact)
	if err := Validate(policy); err != nil {
		return Policy{}, nil, err
	}
	return policy, closed, nil
}

// obligationFree reports whether binding can close no source of outcome.
func obligationFree(outcome Outcome, binding DeliveryBinding) bool {
	return len(binding.Obligations) == 0 && len(outcome.Sources) > 0
}

// Satisfied returns the source identities that recorded deliveries completely satisfy.
func Satisfied(policy Policy) map[string]bool {
	satisfied := map[string]bool{}
	for _, fact := range policy.Deliveries {
		binding, _ := deliveredBinding(policy, fact)
		for _, obligation := range binding.Obligations {
			satisfied[obligation] = true
		}
	}
	return satisfied
}

// Remaining returns the active milestone's ordered outcomes that are not yet delivered.
// The recommended sequence projects exactly these outcomes.
func Remaining(policy Policy) []string {
	delivered := deliveredOutcomes(policy)
	return slices.DeleteFunc(Selection(policy).Outcomes, func(outcome string) bool { return delivered[outcome] })
}

// ScopeDelivered reports whether a legacy scope is finished: the policy records a delivery
// for every scope path that the active milestone approves as a deliverable, and the scope
// names at least one such path. A partly delivered scope stays open.
func ScopeDelivered(policy Policy, scope []string) bool {
	delivered := deliveredPaths(policy)
	approved := 0
	for _, path := range scope {
		if _, _, found, _ := activeDeliverable(policy, path); !found {
			continue
		}
		if !delivered[path] {
			return false
		}
		approved++
	}
	return approved > 0
}

// PathDelivered reports whether policy records a verified delivery of the deliverable at
// path. The publication of that delivery closed every row that it satisfied.
func PathDelivered(policy Policy, path string) bool { return deliveredPaths(policy)[path] }

// deliveredPaths holds the path of each deliverable that a recorded fact delivers.
func deliveredPaths(policy Policy) map[string]bool {
	delivered := map[string]bool{}
	for _, fact := range policy.Deliveries {
		if binding, found := deliveredBinding(policy, fact); found {
			delivered[binding.Source.Path] = true
		}
	}
	return delivered
}

// deliveredOutcomes holds each outcome with a recorded delivery and nothing unsettled. An
// outcome with sources owes each source. A rowless outcome owes each of its approved
// deliverables. A partial delivery therefore leaves its outcome open.
func deliveredOutcomes(policy Policy) map[string]bool {
	settled := settle(policy)
	recorded := map[string]bool{}
	for _, fact := range policy.Deliveries {
		recorded[fact.Outcome] = true
	}
	delivered := map[string]bool{}
	for _, milestone := range policy.Milestones {
		for _, outcome := range milestone.Outcomes {
			sources, bindings := settled.unsettled(outcome)
			open := len(sources) > 0
			if len(outcome.Sources) == 0 {
				open = len(bindings) > 0
			}
			delivered[outcome.ID] = recorded[outcome.ID] && !open
		}
	}
	return delivered
}

// Unsettled returns, in policy order, each outcome source that no recorded delivery
// satisfies and each approved deliverable that no recorded delivery delivers. A recorded
// delivery settles what it delivered: its closure deleted each satisfied row, and its
// publication flipped a delivered spec or removed a delivered tickets-only folder. The
// fact, not the current tree, then binds that content.
func Unsettled(policy Policy) ([]SourceBinding, []DeliveryBinding) {
	settled := settle(policy)
	var sources []SourceBinding
	var bindings []DeliveryBinding
	for _, milestone := range policy.Milestones {
		for _, outcome := range milestone.Outcomes {
			openSources, openBindings := settled.unsettled(outcome)
			sources, bindings = append(sources, openSources...), append(bindings, openBindings...)
		}
	}
	return sources, bindings
}

// settlement holds what the recorded deliveries of one policy settle.
type settlement struct {
	satisfied map[string]bool
	delivered map[deliveryKey]bool
}

func settle(policy Policy) settlement {
	return settlement{satisfied: Satisfied(policy), delivered: deliveredKeys(policy)}
}

// unsettled returns the sources and the deliverables of outcome that settled leaves open.
func (settled settlement) unsettled(outcome Outcome) ([]SourceBinding, []DeliveryBinding) {
	var sources []SourceBinding
	for _, source := range outcome.Sources {
		if !settled.satisfied[source.ID] {
			sources = append(sources, source)
		}
	}
	var bindings []DeliveryBinding
	for _, binding := range outcome.Deliverables {
		if !settled.delivered[deliveryKey{outcome: outcome.ID, binding: binding.Source.ID}] {
			bindings = append(bindings, binding)
		}
	}
	return sources, bindings
}

// deliveryKey names one delivered binding. Outcome identities are unique across the
// policy, so the outcome and the binding identify it without the milestone.
type deliveryKey struct{ outcome, binding string }

func (fact DeliveryFact) key() deliveryKey {
	return deliveryKey{outcome: fact.Outcome, binding: fact.Binding}
}

// deliveredKeys holds the key of each recorded delivery.
func deliveredKeys(policy Policy) map[deliveryKey]bool {
	keys := map[deliveryKey]bool{}
	for _, fact := range policy.Deliveries {
		keys[fact.key()] = true
	}
	return keys
}

// deliveredBinding resolves the approved deliverable that fact records.
func deliveredBinding(policy Policy, fact DeliveryFact) (DeliveryBinding, bool) {
	for _, milestone := range policy.Milestones {
		if milestone.ID != fact.Milestone {
			continue
		}
		for _, outcome := range milestone.Outcomes {
			if outcome.ID != fact.Outcome {
				continue
			}
			for _, binding := range outcome.Deliverables {
				if binding.Source.ID == fact.Binding {
					return binding, true
				}
			}
		}
	}
	return DeliveryBinding{}, false
}

// activeDeliverable resolves the one active outcome that approves the deliverable at path.
func activeDeliverable(policy Policy, path string) (Outcome, DeliveryBinding, bool, error) {
	var outcome Outcome
	var binding DeliveryBinding
	found := false
	for _, milestone := range policy.Milestones {
		if milestone.ID != policy.ActiveMilestone {
			continue
		}
		for _, candidate := range milestone.Outcomes {
			for _, approved := range candidate.Deliverables {
				if approved.Source.Path != path {
					continue
				}
				if found {
					return Outcome{}, DeliveryBinding{}, false, fmt.Errorf("deliverable %q has more than one approved outcome", path)
				}
				outcome, binding, found = candidate, approved, true
			}
		}
	}
	return outcome, binding, found, nil
}
