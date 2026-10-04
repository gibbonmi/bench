package ledger

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"
)

// CommitmentState survives assignment cleanup and process termination.
type CommitmentState struct {
	Bindings      []DeliveryBinding    `json:"bindings,omitempty"`
	Claims        []OutcomeClaim       `json:"claims,omitempty"`
	Blockers      []OutcomeBlocker     `json:"blockers,omitempty"`
	Continuations []LegacyContinuation `json:"continuations,omitempty"`
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
		if !ValidIdentity(binding.Assignment) || !digestPattern.MatchString(binding.Request) || !clean(binding.Milestone) || !clean(binding.Outcome) || !clean(binding.Identity) || !commitmentPath(binding.Deliverable) || bindings[binding.Assignment] {
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
	continuations := map[string]bool{}
	for _, continuation := range state.Continuations {
		if !ValidIdentity(continuation.Assignment) || !digestPattern.MatchString(continuation.Request) || len(continuation.Scope) == 0 || continuations[continuation.Assignment] {
			return fmt.Errorf("invalid legacy continuation for assignment %q", continuation.Assignment)
		}
		continuations[continuation.Assignment] = true
		paths := map[string]bool{}
		for _, path := range continuation.Scope {
			if !commitmentPath(path) || paths[path] {
				return fmt.Errorf("invalid legacy continuation scope %q", path)
			}
			paths[path] = true
		}
	}
	return nil
}

// LegacyContinuation binds an approved existing run to its original scope.
type LegacyContinuation struct {
	Assignment string   `json:"assignment"`
	Request    string   `json:"request"`
	Scope      []string `json:"scope"`
}

func commitmentPath(value string) bool {
	return strings.TrimSpace(value) != "" && strings.IndexFunc(value, unicode.IsControl) < 0 && !filepath.IsAbs(value) && filepath.ToSlash(filepath.Clean(value)) == value && value != "." && value != ".." && !strings.HasPrefix(value, "../")
}
