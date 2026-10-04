package commitment

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// ActiveOutcome resolves an outcome in the active milestone.
func ActiveOutcome(policy Policy, id string) (Outcome, error) {
	for _, milestone := range policy.Milestones {
		if milestone.ID != policy.ActiveMilestone {
			continue
		}
		for _, outcome := range milestone.Outcomes {
			if outcome.ID == id {
				return outcome, nil
			}
		}
	}
	return Outcome{}, fmt.Errorf("outcome %q is not committed in the active milestone", id)
}

// Admit decides a start from one policy and runtime snapshot.
func Admit(policy Policy, state intent.CommitmentState, binding intent.DeliveryBinding) (intent.CommitmentState, error) {
	if err := eligible(policy, state, binding.Outcome); err != nil {
		return state, err
	}
	active := claimed(state)
	found := false
	for _, previous := range state.Bindings {
		if previous.Assignment != binding.Assignment {
			continue
		}
		if previous != binding {
			return state, fmt.Errorf("assignment is already bound to outcome %q and deliverable %q", previous.Outcome, previous.Deliverable)
		}
		found = true
	}
	if !found {
		state.Bindings = append(slices.Clone(state.Bindings), binding)
	}
	if !active[binding.Outcome] {
		state.Claims = append(slices.Clone(state.Claims), intent.OutcomeClaim{Milestone: policy.ActiveMilestone, Outcome: binding.Outcome})
	}
	return state, nil
}

// eligible is the one admission decision for an outcome start: it refuses an uncommitted,
// delivered, blocked, dependent, displacing, or out-of-order outcome. Admit and the reader
// projection both ask it, so a reader cannot select work that a start would refuse.
func eligible(policy Policy, state intent.CommitmentState, id string) error {
	outcome, err := ActiveOutcome(policy, id)
	if err != nil {
		return err
	}
	delivered := deliveredOutcomes(policy)
	if delivered[outcome.ID] {
		return fmt.Errorf("outcome %q is already delivered", outcome.ID)
	}
	blocked := blockedOutcomes(state)
	if blocked[outcome.ID] {
		return fmt.Errorf("outcome %q is blocked", outcome.ID)
	}
	for _, dependency := range outcome.Dependencies {
		if !delivered[dependency] {
			return fmt.Errorf("outcome %q has unfinished dependency %q", outcome.ID, dependency)
		}
	}
	active := claimed(state)
	if active[outcome.ID] {
		return nil
	}
	if len(active) != 0 && !permitsParallel(policy, active, outcome.ID) {
		return fmt.Errorf("another outcome is active; no parallel grant admits %q", outcome.ID)
	}
	for _, prior := range Selection(policy).Outcomes {
		if prior == outcome.ID {
			break
		}
		if !delivered[prior] && !blocked[prior] && !active[prior] {
			return fmt.Errorf("outcome %q precedes %q", prior, outcome.ID)
		}
	}
	return nil
}

func claimed(state intent.CommitmentState) map[string]bool {
	active := map[string]bool{}
	for _, claim := range state.Claims {
		active[claim.Outcome] = true
	}
	return active
}

func blockedOutcomes(state intent.CommitmentState) map[string]bool {
	blocked := map[string]bool{}
	for _, blocker := range state.Blockers {
		blocked[blocker.Outcome] = true
	}
	return blocked
}

func permitsParallel(policy Policy, active map[string]bool, candidate string) bool {
	for _, grant := range policy.ParallelGrants {
		if !slices.Contains(grant.Outcomes, candidate) {
			continue
		}
		matches := true
		for id := range active {
			if !slices.Contains(grant.Outcomes, id) {
				matches = false
			}
		}
		if matches {
			return true
		}
	}
	return false
}

// SetBlocker records or clears a blocker without displacing another active claim.
func SetBlocker(policy Policy, state intent.CommitmentState, outcome, reason string, blocked bool) (intent.CommitmentState, error) {
	if _, err := ActiveOutcome(policy, outcome); err != nil {
		return state, err
	}
	if deliveredOutcomes(policy)[outcome] {
		return state, fmt.Errorf("outcome %q is already delivered", outcome)
	}
	if blocked && (strings.TrimSpace(reason) == "" || !sanitize.LineSafe(reason)) {
		return state, fmt.Errorf("blocker reason must be one nonempty line")
	}
	state.Blockers = slices.DeleteFunc(slices.Clone(state.Blockers), func(item intent.OutcomeBlocker) bool { return item.Outcome == outcome })
	if blocked {
		state.Blockers = append(state.Blockers, intent.OutcomeBlocker{Outcome: outcome, Reason: reason})
		state.Claims = slices.DeleteFunc(slices.Clone(state.Claims), func(item intent.OutcomeClaim) bool { return item.Outcome == outcome })
	}
	return state, nil
}
