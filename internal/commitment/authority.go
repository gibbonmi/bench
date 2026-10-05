package commitment

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/intent"
)

// Effects is the exact policy displacement produced by a proposal.
type Effects struct {
	Added              []string `json:"added"`
	Reordered          []string `json:"reordered"`
	Delayed            []string `json:"delayed"`
	Removed            []string `json:"removed"`
	Activated          []string `json:"activated"`
	Switched           []string `json:"switched"`
	ParallelAuthorized []string `json:"parallel_authorized"`
	Completed          []string `json:"completed"`
}

// Plan binds one predecessor to canonical proposed policy bytes and effects. Continuations
// lists each already-authorized run that approval lets finish.
type Plan struct {
	ID               string                      `json:"id"`
	Predecessor      string                      `json:"predecessor"`
	ProposalIdentity string                      `json:"proposal_identity"`
	Proposed         []byte                      `json:"proposed"`
	Sources          []SourceBinding             `json:"sources"`
	Effects          Effects                     `json:"effects"`
	Continuations    []intent.LegacyContinuation `json:"continuations,omitempty"`
}

// Proposal is one plan input: the proposed policy and each already-authorized run that
// approval lets finish. The plan identity binds the list, and approval records it in the
// local intent record. The tracked policy never holds it.
type Proposal struct {
	Policy        Policy
	Continuations []intent.LegacyContinuation
}

// BuildPlan computes the exact transition from current to the proposal.
func BuildPlan(current *Policy, proposal Proposal) (Plan, error) {
	proposed := proposal.Policy
	proposedBytes, err := Bytes(proposed)
	if err != nil {
		return Plan{}, err
	}
	if err := refuseObligationFree(current, proposed); err != nil {
		return Plan{}, err
	}
	predecessor := "absent"
	if current != nil {
		currentBytes, err := Bytes(*current)
		if err != nil {
			return Plan{}, err
		}
		predecessor = identity(currentBytes)
	}
	effects := transitionEffects(current, proposed)
	sources := policySources(proposed)
	// Each open source of the current policy that the proposal does not hold stays bound, so
	// the approval of a removal refuses when the removed row changes. A deliverable that the
	// proposal drops binds nothing, so its deletion from the default branch leaves the plan
	// valid.
	if current != nil {
		open, _ := Unsettled(*current)
		for _, source := range open {
			if !slices.Contains(sources, source) {
				sources = append(sources, source)
			}
		}
	}
	effectBytes, err := json.Marshal(effects)
	if err != nil {
		return Plan{}, fmt.Errorf("encode commitment effects: %w", err)
	}
	var binding strings.Builder
	for _, source := range sources {
		fmt.Fprintf(&binding, "%s\x00%s\x00%s\x00", source.ID, source.Path, source.Identity)
	}
	proposalIdentity := identity(proposedBytes)
	bound := predecessor + "\x00" + proposalIdentity + "\x00" + binding.String() + "\x00" + string(effectBytes)
	// A plan that lists no run binds no list, so the policy transition alone decides its identity.
	if len(proposal.Continuations) != 0 {
		listed, err := json.Marshal(proposal.Continuations)
		if err != nil {
			return Plan{}, fmt.Errorf("encode commitment continuations: %w", err)
		}
		bound += "\x00" + string(listed)
	}
	return Plan{
		ID:               identity([]byte(bound)),
		Predecessor:      predecessor,
		ProposalIdentity: proposalIdentity,
		Proposed:         proposedBytes,
		Sources:          sources,
		Effects:          effects,
		Continuations:    proposal.Continuations,
	}, nil
}

// refuseObligationFree refuses each obligation-free binding of proposed, in any milestone,
// unless current holds the same outcome with an equal binding that is already
// obligation-free. Validate does not take this rule, so a legacy policy stays readable and
// a plan can repair it.
func refuseObligationFree(current *Policy, proposed Policy) error {
	for _, milestone := range proposed.Milestones {
		for _, outcome := range milestone.Outcomes {
			for _, binding := range outcome.Deliverables {
				if obligationFree(outcome, binding) && !retainedObligationFree(current, outcome.ID, binding) {
					return fmt.Errorf("outcome %q deliverable %q names no obligation; list each outcome source that it completely satisfies", outcome.ID, binding.Source.ID)
				}
			}
		}
	}
	return nil
}

// retainedObligationFree reports whether current approves binding unchanged for the outcome
// named id, and that binding is already obligation-free there.
func retainedObligationFree(current *Policy, id string, binding DeliveryBinding) bool {
	if current == nil {
		return false
	}
	for _, milestone := range current.Milestones {
		for _, outcome := range milestone.Outcomes {
			if outcome.ID == id {
				return slices.ContainsFunc(outcome.Deliverables, func(kept DeliveryBinding) bool {
					return kept.Source == binding.Source && slices.Equal(kept.Obligations, binding.Obligations) && obligationFree(outcome, kept)
				})
			}
		}
	}
	return false
}

// Identity returns the content identity used by plans and source bindings.
func Identity(data []byte) string { return identity(data) }

func identity(data []byte) string {
	digest := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(digest[:])
}

func transitionEffects(current *Policy, proposed Policy) Effects {
	effects := Effects{
		Added:              []string{},
		Reordered:          []string{},
		Delayed:            []string{},
		Removed:            []string{},
		Activated:          []string{},
		Switched:           []string{},
		ParallelAuthorized: []string{},
		Completed:          []string{},
	}
	proposedAll := orderedOutcomes(proposed)
	proposedSet := sliceSet(proposedAll)
	if current == nil {
		effects.Added = append(effects.Added, proposedAll...)
		effects.Activated = append(effects.Activated, Selection(proposed).Outcomes...)
		effects.ParallelAuthorized = parallelKeys(proposed.ParallelGrants)
		return effects
	}
	currentAll := orderedOutcomes(*current)
	currentSet := sliceSet(currentAll)
	for _, outcome := range proposedAll {
		if !currentSet[outcome] {
			effects.Added = append(effects.Added, outcome)
		}
	}
	delivered := deliveredOutcomes(*current)
	for _, outcome := range currentAll {
		if !proposedSet[outcome] && !delivered[outcome] {
			effects.Removed = append(effects.Removed, outcome)
		}
	}
	currentSelection := Selection(*current)
	proposedSelection := Selection(proposed)
	switching := current.ActiveMilestone != "" && current.ActiveMilestone != proposed.ActiveMilestone
	switch {
	case current.ActiveMilestone == "" && proposed.ActiveMilestone != "":
		effects.Activated = append(effects.Activated, proposedSelection.Outcomes...)
	case switching:
		for _, outcome := range currentSelection.Outcomes {
			if !delivered[outcome] {
				effects.Delayed = append(effects.Delayed, outcome)
				effects.Switched = append(effects.Switched, outcome)
			}
		}
	case current.ActiveMilestone == proposed.ActiveMilestone:
		positions := make(map[string]int, len(currentSelection.Outcomes))
		for index, outcome := range currentSelection.Outcomes {
			positions[outcome] = index
		}
		for index, outcome := range proposedSelection.Outcomes {
			previous, exists := positions[outcome]
			if !exists || delivered[outcome] {
				continue
			}
			for _, earlier := range proposedSelection.Outcomes[:index] {
				position, existed := positions[earlier]
				if !existed || position > previous {
					effects.Delayed = append(effects.Delayed, outcome)
					break
				}
			}
		}
		commonCurrent := retainedOrder(currentSelection.Outcomes, proposedSet)
		commonProposed := retainedOrder(proposedSelection.Outcomes, currentSet)
		if !slices.Equal(commonCurrent, commonProposed) {
			effects.Reordered = append(effects.Reordered, commonProposed...)
		}
	}
	for _, completion := range addedCompletions(current.Completions, proposed.Completions) {
		effects.Completed = append(effects.Completed, completion.Milestone)
	}
	oldGrants := sliceSet(parallelKeys(current.ParallelGrants))
	for _, grant := range parallelKeys(proposed.ParallelGrants) {
		if !oldGrants[grant] {
			effects.ParallelAuthorized = append(effects.ParallelAuthorized, grant)
		}
	}
	return effects
}

func orderedOutcomes(policy Policy) []string {
	var outcomes []string
	for _, milestone := range policy.Milestones {
		for _, outcome := range milestone.Outcomes {
			outcomes = append(outcomes, outcome.ID)
		}
	}
	return outcomes
}

// policySources returns the sources and deliverables that a plan binds to current content:
// the content that no recorded delivery settles.
func policySources(policy Policy) []SourceBinding {
	sources, bindings := Unsettled(policy)
	for _, binding := range bindings {
		if !slices.Contains(sources, binding.Source) {
			sources = append(sources, binding.Source)
		}
	}
	return sources
}

// addedCompletions returns each completion of proposed that current does not record.
func addedCompletions(current, proposed []Completion) []Completion {
	var added []Completion
	for _, completion := range proposed {
		if !slices.Contains(current, completion) {
			added = append(added, completion)
		}
	}
	return added
}

func retainedOrder(order []string, retained map[string]bool) []string {
	var result []string
	for _, item := range order {
		if retained[item] {
			result = append(result, item)
		}
	}
	return result
}

func sliceSet(values []string) map[string]bool {
	set := make(map[string]bool, len(values))
	for _, value := range values {
		set[value] = true
	}
	return set
}

func parallelKeys(grants []ParallelGrant) []string {
	keys := make([]string, 0, len(grants))
	for _, grant := range grants {
		outcomes := append([]string(nil), grant.Outcomes...)
		slices.Sort(outcomes)
		continuations := append([]string(nil), grant.Continuations...)
		slices.Sort(continuations)
		keys = append(keys, strings.Join(append(outcomes, continuations...), "+"))
	}
	return keys
}
