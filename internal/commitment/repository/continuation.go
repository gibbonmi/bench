package repository

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/refusalroute"
)

// listedRuns checks each run that a plan lists against the published policy and the run
// record it names. Only the initial adoption lists runs, so a published policy refuses
// any list. The run must be known and active, its request digest must match, and each
// scope path must exist on the run's branch, so a continuation never grants scope that
// the run does not hold. Plan and approval both ask it.
func (store Store) listedRuns(ledger intent.Ledger, current *commitment.Policy, continuations []intent.LegacyContinuation) error {
	if len(continuations) != 0 && current != nil {
		return refusalroute.Raised{Name: refusalroute.CommitmentDecision, Err: errors.New("commitment continuation refused: a policy is published, so only the initial adoption lists runs")}
	}
	for _, continuation := range continuations {
		index := slices.IndexFunc(ledger.Assignments, func(run intent.Assignment) bool { return run.ID == continuation.Assignment })
		if index < 0 {
			return refusalroute.Raised{Name: refusalroute.CommitmentRunUnknown, Err: fmt.Errorf("commitment continuation refused: run %q is unknown", continuation.Assignment)}
		}
		run := ledger.Assignments[index]
		if run.State != intent.StateActive {
			return fmt.Errorf("commitment continuation refused: run %q is not active (%s)", run.ID, run.State)
		}
		if run.Request != continuation.Request {
			return fmt.Errorf("commitment continuation refused: request does not match run %q", run.ID)
		}
		for _, path := range continuation.Scope {
			listing, err := git.Output("-C", store.Root, "--literal-pathspecs", "ls-tree", "-z", run.Branch, "--", path)
			if err != nil {
				return fmt.Errorf("commitment continuation refused: cannot read the branch of run %q: %w", run.ID, err)
			}
			if strings.TrimSuffix(listing, "\x00") == "" {
				return fmt.Errorf("commitment continuation refused: scope %q is outside run %q", path, run.ID)
			}
		}
	}
	return nil
}

// withContinuations records exactly the listed runs as legacy continuations. A listed run
// replaces its earlier continuation, and every other continuation stays.
func withContinuations(ledger intent.Ledger, continuations []intent.LegacyContinuation) intent.Ledger {
	if len(continuations) == 0 {
		return ledger
	}
	state := runtimeState(ledger)
	state.Continuations = slices.DeleteFunc(slices.Clone(state.Continuations), func(existing intent.LegacyContinuation) bool {
		return slices.ContainsFunc(continuations, func(listed intent.LegacyContinuation) bool { return listed.Assignment == existing.Assignment })
	})
	state.Continuations = append(state.Continuations, continuations...)
	ledger.Commitment = &state
	return ledger
}
