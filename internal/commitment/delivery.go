package commitment

import (
	"fmt"
	"slices"
)

// Deliver records the verified delivery of the approved deliverable at path. It returns
// the next policy and the outcome sources that the delivery completely satisfies. An
// outcome with no source is rowless: the outcome itself is its one obligation, so its
// delivery records a fact and satisfies no source. A path that the active milestone does
// not approve, or a binding that names no obligation of an outcome with sources, returns
// policy unchanged and no source.
func Deliver(policy Policy, path, source, evidence string) (Policy, []SourceBinding, error) {
	outcome, binding, found, err := activeDeliverable(policy, path)
	if err != nil || !found || (len(binding.Obligations) == 0 && len(outcome.Sources) > 0) {
		return policy, nil, err
	}
	fact := DeliveryFact{Milestone: policy.ActiveMilestone, Outcome: outcome.ID, Binding: binding.Source.ID, Identity: binding.Source.Identity, Source: source, Evidence: evidence}
	for _, recorded := range policy.Deliveries {
		if recorded.Outcome == fact.Outcome && recorded.Binding == fact.Binding {
			return policy, nil, fmt.Errorf("deliverable %q is already delivered", path)
		}
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

// deliveredOutcomes holds each outcome with a recorded delivery and no unsatisfied
// source. A partial delivery therefore leaves its outcome open.
func deliveredOutcomes(policy Policy) map[string]bool {
	satisfied := Satisfied(policy)
	recorded := map[string]bool{}
	for _, fact := range policy.Deliveries {
		recorded[fact.Outcome] = true
	}
	delivered := map[string]bool{}
	for _, milestone := range policy.Milestones {
		for _, outcome := range milestone.Outcomes {
			delivered[outcome.ID] = recorded[outcome.ID] && !slices.ContainsFunc(outcome.Sources, func(source SourceBinding) bool { return !satisfied[source.ID] })
		}
	}
	return delivered
}

// BindingDelivered reports whether policy records the delivery of the deliverable binding
// that outcome approves in milestone.
func BindingDelivered(policy Policy, milestone, outcome, binding string) bool {
	return slices.ContainsFunc(policy.Deliveries, func(fact DeliveryFact) bool {
		return fact.Milestone == milestone && fact.Outcome == outcome && fact.Binding == binding
	})
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
