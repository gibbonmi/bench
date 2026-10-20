package repository

import (
	"errors"
	"fmt"
	"slices"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/spec"
)

// openBuild is one chunked build that holds the default branch. bound maps the assignment
// of each qualifying binding to its deliverable. deliverable and target come from the first
// qualifying binding: the refusal names the deliverable, and target is the bound assignment
// that is the build's integration target.
type openBuild struct {
	claim       intent.OutcomeClaim
	bound       map[string]string
	deliverable string
	target      intent.Assignment
}

// HoldsOpenBuild reports whether the ledger at the store root holds an open build of any
// outcome, whatever its chunk gap.
func (store Store) HoldsOpenBuild() (bool, error) {
	builds, err := store.readOpenBuilds()
	return len(builds) > 0, err
}

// DeclaredBuild reports whether the ledger at the store root holds a declared build, which
// is an open build whose bound spec declares a persisted field at the default-branch tip. It
// answers that build's integration target. A declaration that does not parse counts as
// none, because the install refuses that spec, so no candidate can be installed for it.
func (store Store) DeclaredBuild() (intent.Assignment, bool, error) {
	builds, err := store.readOpenBuilds()
	if err != nil {
		return intent.Assignment{}, false, err
	}
	for _, build := range builds {
		fields, err := store.DefaultDeclaration(build.deliverable)
		if errors.Is(err, spec.ErrPersistedFieldsInvalid) {
			continue
		}
		if err != nil {
			return intent.Assignment{}, false, err
		}
		if len(fields) > 0 {
			return build.target, true, nil
		}
	}
	return intent.Assignment{}, false, nil
}

// readOpenBuilds derives the open builds from one strict read of the ledger at the store
// root. Both exported answers read through it, so they share one derivation.
func (store Store) readOpenBuilds() ([]openBuild, error) {
	ledger, err := intent.Read(store.Root)
	if err != nil {
		return nil, err
	}
	return store.openBuilds(ledger)
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

// window refuses source unless each open build admits it. An override for the assignment of
// source admits it through every build. The first refusing build in claim order names the
// refusal.
func (store Store) window(ledger intent.Ledger, source Publication) error {
	state := storedState(ledger)
	if holdsOverride(state.Overrides, source.Assignment) {
		return nil
	}
	refusing, err := store.refusingBuilds(ledger, source)
	if err != nil {
		return refusalroute.Raised{Name: refusalroute.LandInfrastructure, Err: fmt.Errorf("build window unavailable: %w", err)}
	}
	if len(refusing) == 0 {
		return nil
	}
	build := refusing[0]
	gap, recorded := chunkGap(state, build.claim)
	return refusalroute.Raised{Name: refusalroute.LandBuildWindow, Err: fmt.Errorf("landing destination %s is held by open build %s (%s); wait for %s", source.Destination, build.claim.Outcome, build.deliverable, waitFor(gap, recorded))}
}

// refusingBuilds are the open builds of ledger that refuse source, in claim order, whatever
// override the ledger holds. A build admits its own landing and every landing during its
// open gap.
func (store Store) refusingBuilds(ledger intent.Ledger, source Publication) ([]openBuild, error) {
	builds, err := store.openBuilds(ledger)
	if err != nil {
		return nil, err
	}
	state := storedState(ledger)
	return slices.DeleteFunc(builds, func(build openBuild) bool {
		if deliverable, own := build.bound[source.Assignment]; own && deliverable == source.Deliverable {
			return true
		}
		gap, recorded := chunkGap(state, build.claim)
		return recorded && gap.Open
	}), nil
}

// openBuilds derives the open builds of ledger in claim order. A qualifying binding names an
// assignment that is active in the ledger. The start's own validation reads its deliverable
// as a staged spec at the default-branch tip. A claimed outcome needs that tip, so a fault
// that resolves it is an error. A classification refusal only disqualifies the binding.
func (store Store) openBuilds(ledger intent.Ledger) ([]openBuild, error) {
	state := runtimeState(ledger)
	if len(state.Claims) == 0 {
		return nil, nil
	}
	revision, err := store.sourceRevision()
	if err != nil {
		return nil, err
	}
	var builds []openBuild
	for _, claim := range state.Claims {
		build := openBuild{claim: claim, bound: map[string]string{}}
		for _, binding := range state.Bindings {
			owner, active := activeByID(ledger, binding.Assignment)
			if binding.Milestone != claim.Milestone || binding.Outcome != claim.Outcome || !active {
				continue
			}
			if staged, err := store.validateDeliverableAt(revision, commitment.SourceBinding{ID: binding.Deliverable, Path: binding.Deliverable, Identity: binding.Identity}); err != nil || !staged {
				continue
			}
			if len(build.bound) == 0 {
				build.deliverable, build.target = binding.Deliverable, owner
			}
			build.bound[binding.Assignment] = binding.Deliverable
		}
		if len(build.bound) > 0 {
			builds = append(builds, build)
		}
	}
	return builds, nil
}

// chunkGap is the gap record of the claimed outcome, if the ledger holds one.
func chunkGap(state intent.CommitmentState, claim intent.OutcomeClaim) (intent.ChunkGap, bool) {
	if i := gapIndex(state, claim); i >= 0 {
		return state.Gaps[i], true
	}
	return intent.ChunkGap{}, false
}

// gapIndex is the position of the claimed outcome's gap record in state, or -1.
func gapIndex(state intent.CommitmentState, claim intent.OutcomeClaim) int {
	return slices.IndexFunc(state.Gaps, func(gap intent.ChunkGap) bool {
		return gap.Milestone == claim.Milestone && gap.Outcome == claim.Outcome
	})
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
