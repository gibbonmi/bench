package repository

import (
	"fmt"
	"slices"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/refusalroute"
)

// openBuild is one chunked build that holds the default branch. bound maps the assignment
// of each qualifying binding to its deliverable. deliverable is the first one, which the
// refusal names.
type openBuild struct {
	claim       intent.OutcomeClaim
	bound       map[string]string
	deliverable string
}

// AdmitWindow decides the build window for source from one ledger read. The landing asks it
// before the lease, and PublishAdmitted decides again under the intent lock.
func (store Store) AdmitWindow(source Publication) error {
	ledger, err := intent.Read(store.Root)
	if err != nil {
		return err
	}
	return store.window(ledger, source)
}

// window refuses source unless each open build admits it. A build admits its own landing and
// every landing during its open gap. An override for the assignment of source admits it
// through every build. The first refusing build in claim order names the refusal.
func (store Store) window(ledger intent.Ledger, source Publication) error {
	state := storedState(ledger)
	if slices.ContainsFunc(state.Overrides, func(override intent.LandingOverride) bool { return override.Assignment == source.Assignment }) {
		return nil
	}
	for _, build := range store.openBuilds(ledger) {
		if deliverable, own := build.bound[source.Assignment]; own && deliverable == source.Deliverable {
			continue
		}
		gap, recorded := chunkGap(state, build.claim)
		if recorded && gap.Open {
			continue
		}
		return refusalroute.Raised{Name: refusalroute.LandBuildWindow, Err: fmt.Errorf("landing destination %s is held by open build %s (%s); wait for %s", source.Destination, build.claim.Outcome, build.deliverable, waitFor(gap, recorded))}
	}
	return nil
}

// openBuilds derives the open builds of ledger in claim order. A qualifying binding names an
// assignment that is active in the ledger. The start's own validation reads its deliverable
// as a staged spec at the default-branch tip.
func (store Store) openBuilds(ledger intent.Ledger) []openBuild {
	state := runtimeState(ledger)
	var builds []openBuild
	for _, claim := range state.Claims {
		build := openBuild{claim: claim, bound: map[string]string{}}
		for _, binding := range state.Bindings {
			if binding.Milestone != claim.Milestone || binding.Outcome != claim.Outcome || !activeAssignment(ledger, binding.Assignment) {
				continue
			}
			if staged, err := store.validateDeliverable(commitment.SourceBinding{ID: binding.Deliverable, Path: binding.Deliverable, Identity: binding.Identity}); err != nil || !staged {
				continue
			}
			if len(build.bound) == 0 {
				build.deliverable = binding.Deliverable
			}
			build.bound[binding.Assignment] = binding.Deliverable
		}
		if len(build.bound) > 0 {
			builds = append(builds, build)
		}
	}
	return builds
}

// activeAssignment reports whether the ledger holds id as an active assignment.
func activeAssignment(ledger intent.Ledger, id string) bool {
	return slices.ContainsFunc(ledger.Assignments, func(owner intent.Assignment) bool {
		return owner.ID == id && owner.State == intent.StateActive
	})
}

// chunkGap is the gap record of the claimed outcome, if the ledger holds one.
func chunkGap(state intent.CommitmentState, claim intent.OutcomeClaim) (intent.ChunkGap, bool) {
	for _, gap := range state.Gaps {
		if gap.Milestone == claim.Milestone && gap.Outcome == claim.Outcome {
			return gap, true
		}
	}
	return intent.ChunkGap{}, false
}

// waitFor names the gap that a closed or absent gap record leaves the landing to wait for.
func waitFor(gap intent.ChunkGap, recorded bool) string {
	switch {
	case !recorded:
		return "its first chunk checkpoint"
	case gap.Last:
		return "its completion landing"
	}
	return "the checkpoint of chunk " + gap.Chunk
}
