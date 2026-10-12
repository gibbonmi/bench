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
	Leases        []LandingLease       `json:"leases,omitempty"`
	Gaps          []ChunkGap           `json:"gaps,omitempty"`
	Overrides     []LandingOverride    `json:"overrides,omitempty"`
}

// LandingLease names the assignment whose landing holds the destination ref, from its
// composition through its publish.
type LandingLease struct {
	Destination string `json:"destination"`
	Assignment  string `json:"assignment"`
}

// ChunkGap is the chunk state of one outcome's build. Open states that the build's last
// checkpointed chunk lets an unrelated landing through, and Last states that Chunk is the
// plan's last chunk.
type ChunkGap struct {
	Milestone string `json:"milestone"`
	Outcome   string `json:"outcome"`
	Chunk     string `json:"chunk"`
	Open      bool   `json:"open"`
	Last      bool   `json:"last"`
}

// LandingOverride is the reviewer's unconsumed grant that lets one assignment's landing
// through every open build, with the reviewer's one-line reason.
type LandingOverride struct {
	Assignment string `json:"assignment"`
	Reason     string `json:"reason"`
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
	leases := map[string]bool{}
	for _, lease := range state.Leases {
		if !clean(lease.Destination) || !ValidIdentity(lease.Assignment) || leases[lease.Destination] {
			return fmt.Errorf("invalid landing lease on %q", lease.Destination)
		}
		leases[lease.Destination] = true
	}
	gaps := map[string]bool{}
	for _, gap := range state.Gaps {
		if !clean(gap.Milestone) || !clean(gap.Outcome) || !clean(gap.Chunk) || gaps[gap.Outcome] {
			return fmt.Errorf("invalid chunk gap for outcome %q", gap.Outcome)
		}
		gaps[gap.Outcome] = true
	}
	overrides := map[string]bool{}
	for _, override := range state.Overrides {
		if !ValidIdentity(override.Assignment) || !clean(override.Reason) || overrides[override.Assignment] {
			return fmt.Errorf("invalid landing override for assignment %q", override.Assignment)
		}
		overrides[override.Assignment] = true
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
