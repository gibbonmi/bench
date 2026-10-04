package commitment

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
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
}

// Plan binds one predecessor to canonical proposed policy bytes and effects.
type Plan struct {
	ID               string          `json:"id"`
	Predecessor      string          `json:"predecessor"`
	ProposalIdentity string          `json:"proposal_identity"`
	Proposed         []byte          `json:"proposed"`
	Sources          []SourceBinding `json:"sources"`
	Effects          Effects         `json:"effects"`
}

// BuildPlan computes the exact transition from current to proposed.
func BuildPlan(current *Policy, proposed Policy) (Plan, error) {
	proposedBytes, err := Bytes(proposed)
	if err != nil {
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
	if current != nil {
		for _, source := range policySources(*current) {
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
	planIdentity := identity([]byte(predecessor + "\x00" + proposalIdentity + "\x00" + binding.String() + "\x00" + string(effectBytes)))
	return Plan{
		ID:               planIdentity,
		Predecessor:      predecessor,
		ProposalIdentity: proposalIdentity,
		Proposed:         proposedBytes,
		Sources:          sources,
		Effects:          effects,
	}, nil
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

func policySources(policy Policy) []SourceBinding {
	var sources []SourceBinding
	for _, milestone := range policy.Milestones {
		for _, outcome := range milestone.Outcomes {
			sources = append(sources, outcome.Sources...)
			for _, binding := range outcome.Deliverables {
				if !slices.Contains(sources, binding.Source) {
					sources = append(sources, binding.Source)
				}
			}
		}
	}
	return sources
}

func deliveredOutcomes(policy Policy) map[string]bool {
	delivered := map[string]bool{}
	for _, fact := range policy.Deliveries {
		delivered[fact.Outcome] = true
	}
	return delivered
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
		keys = append(keys, strings.Join(outcomes, "+"))
	}
	return keys
}
