package ledger

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"
)

// CommitmentState survives assignment cleanup and process termination.
type CommitmentState struct {
	Bindings []DeliveryBinding `json:"bindings,omitempty"`
	Claims   []OutcomeClaim    `json:"claims,omitempty"`
	Blockers []OutcomeBlocker  `json:"blockers,omitempty"`
}

type DeliveryBinding struct {
	Assignment  string `json:"assignment"`
	Request     string `json:"request"`
	Milestone   string `json:"milestone"`
	Outcome     string `json:"outcome"`
	Deliverable string `json:"deliverable"`
	Identity    string `json:"identity"`
}

type OutcomeClaim struct {
	Milestone string `json:"milestone"`
	Outcome   string `json:"outcome"`
}

type OutcomeBlocker struct {
	Outcome string `json:"outcome"`
	Reason  string `json:"reason"`
}

// ValidateCommitment checks durable runtime records independently of live assignments.
func ValidateCommitment(state *CommitmentState) error {
	if state == nil {
		return nil
	}
	clean := func(value string) bool {
		return strings.TrimSpace(value) != "" && strings.IndexFunc(value, unicode.IsControl) < 0
	}
	bindings := map[string]bool{}
	for _, binding := range state.Bindings {
		if !ValidIdentity(binding.Assignment) || !digestPattern.MatchString(binding.Request) || !clean(binding.Milestone) || !clean(binding.Outcome) || !clean(binding.Identity) || !clean(binding.Deliverable) || filepath.IsAbs(binding.Deliverable) || filepath.ToSlash(filepath.Clean(binding.Deliverable)) != binding.Deliverable || binding.Deliverable == ".." || strings.HasPrefix(binding.Deliverable, "../") || bindings[binding.Assignment] {
			return fmt.Errorf("invalid commitment delivery binding for assignment %q", binding.Assignment)
		}
		bindings[binding.Assignment] = true
	}
	claims := map[string]bool{}
	for _, claim := range state.Claims {
		if !clean(claim.Milestone) || !clean(claim.Outcome) || claims[claim.Outcome] {
			return fmt.Errorf("invalid commitment claim for outcome %q", claim.Outcome)
		}
		claims[claim.Outcome] = true
	}
	blockers := map[string]bool{}
	for _, blocker := range state.Blockers {
		if !clean(blocker.Outcome) || !clean(blocker.Reason) || blockers[blocker.Outcome] || claims[blocker.Outcome] {
			return fmt.Errorf("invalid commitment blocker for outcome %q", blocker.Outcome)
		}
		blockers[blocker.Outcome] = true
	}
	return nil
}
