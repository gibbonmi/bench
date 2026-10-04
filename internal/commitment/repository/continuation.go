package repository

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

// listedRuns checks each run that a plan lists against the run record it names. The run
// must be known, its request digest must match, and each scope path must exist on the
// run's branch, so a continuation never grants scope that the run does not hold.
func (store Store) listedRuns(ledger intent.Ledger, continuations []intent.LegacyContinuation) error {
	for _, continuation := range continuations {
		index := slices.IndexFunc(ledger.Assignments, func(run intent.Assignment) bool { return run.ID == continuation.Assignment })
		if index < 0 {
			return fmt.Errorf("commitment continuation refused: run %q is unknown; run bench commitment inventory", continuation.Assignment)
		}
		run := ledger.Assignments[index]
		if run.Request != continuation.Request {
			return fmt.Errorf("commitment continuation refused: request does not match run %q", run.ID)
		}
		for _, path := range continuation.Scope {
			listing, err := git.Output("-C", store.Root, "--literal-pathspecs", "ls-tree", "-z", run.Branch, "--", path)
			if err != nil || strings.TrimSuffix(listing, "\x00") == "" {
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
