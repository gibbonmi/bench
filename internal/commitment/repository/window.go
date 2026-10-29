package repository

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/brokermanifest"
	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/git"
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

// DeclaredBuild reports whether the ledger at the store root holds a declared build, and
// answers the integration target of the first one.
func (store Store) DeclaredBuild() (intent.Assignment, bool, error) {
	builds, err := store.declaredBuilds()
	if err != nil || len(builds) == 0 {
		return intent.Assignment{}, false, err
	}
	return builds[0].target, true, nil
}

// declaredBuilds lists the declared builds of the ledger at the store root, in claim order.
// A declared build is an open build whose bound spec declares a persisted field at the
// default-branch tip. A declaration that does not parse counts as none, because the install
// refuses that spec, so no candidate can be installed for it.
func (store Store) declaredBuilds() ([]openBuild, error) {
	builds, err := store.readOpenBuilds()
	if err != nil {
		return nil, err
	}
	var declared []openBuild
	for _, build := range builds {
		fields, err := store.DefaultDeclaration(build.deliverable)
		if errors.Is(err, spec.ErrPersistedFieldsInvalid) {
			continue
		}
		if err != nil {
			return nil, err
		}
		if len(fields) > 0 {
			declared = append(declared, build)
		}
	}
	return declared, nil
}

// CheckEntry is the candidate check of a guarded entry in the checkout at root, over the
// declaration of the spec at path in the commit that HEAD names and the running executable.
// Where the repository declares Bench build inputs, or where its primary checkout does not
// resolve, a spec that does not read and a declaration that does not parse refuse with their
// cause. A declaration with a field then refuses with candidate-not-installed when the primary
// checkout does not resolve.
func CheckEntry(root, path string) error {
	// An executable that does not resolve leaves no path, which no digest reads, so the
	// check refuses it.
	executable, _ := os.Executable()
	return checkCandidate(executable, root, func() ([]spec.PersistedField, error) {
		data, err := git.ReadTreeFile(root, "HEAD", path)
		if err != nil {
			return nil, fmt.Errorf("read the declaration of %s at HEAD: %w", path, err)
		}
		return spec.PersistedFields(data)
	})
}

// CheckCandidate refuses a guarded write of a declared build unless executable is the
// installed candidate: its digest must equal the digest that the primary checkout's broker
// manifest binds. The primary checkout comes from the Git common directory of root, never
// from root itself, because a worktree's own wrapper directory can hold a manifest of the
// worktree executable. An undeclared spec, and a repository that declares no Bench build
// inputs, pass. The check compares digests only: an executable that knows a declared field
// but is not the installed bytes still refuses. While more than one declared build is open,
// no executable is the installed candidate of each, so the check refuses before the digest
// test.
func CheckCandidate(declaration []spec.PersistedField, executable, root string) error {
	return checkCandidate(executable, root, func() ([]spec.PersistedField, error) { return declaration, nil })
}

// checkCandidate is CheckCandidate over the declaration that declare reads. A repository
// without Bench build inputs keeps the declaration inert, so the check calls declare only
// after the build-inputs test.
func checkCandidate(executable, root string, declare func() ([]spec.PersistedField, error)) error {
	primary, located := git.PrimaryCheckout(root)
	if located == nil && !freshness.DeclaresBuildInputs(primary) {
		return nil
	}
	declaration, err := declare()
	if err != nil {
		return err
	}
	if len(declaration) == 0 {
		return nil
	}
	// Declared builds that do not read, whatever the cause, list no build: the ledger, the
	// default-branch tip, or a bound spec can fail. A refusal with no build prints the target
	// slot. With no build, the overlap count passes, so the digest test runs when the primary
	// checkout resolves.
	builds, _ := Store{Root: root}.declaredBuilds()
	if located != nil {
		return candidateNotInstalled(builds, fmt.Errorf("the primary checkout does not resolve: %w", located))
	}
	if len(builds) > 1 {
		return declaredBuildsOverlap(builds)
	}
	manifest := brokermanifest.Path(primary)
	fields, err := brokermanifest.Read(manifest)
	if err != nil {
		return candidateNotInstalled(builds, fmt.Errorf("the primary broker manifest does not read: %w", err))
	}
	digest, err := brokermanifest.Digest(executable)
	if err != nil {
		return candidateNotInstalled(builds, fmt.Errorf("the running executable %q does not digest: %w", executable, err))
	}
	if digest != fields[brokermanifest.DigestField] {
		return candidateNotInstalled(builds, fmt.Errorf("the running executable %q has digest %s, and the primary broker manifest %s binds %s", executable, digest, manifest, fields[brokermanifest.DigestField]))
	}
	return nil
}

// candidateNotInstalled raises the candidate refusal for cause. The route names the
// integration target of the first declared build of builds, whatever checkout the caller
// runs in. With no build, the route prints the target slot. A ledger that holds no open
// declared build leaves no build, and so does a declared-build read that fails.
func candidateNotInstalled(builds []openBuild, cause error) error {
	values := map[string]string{}
	if len(builds) > 0 {
		values[refusalroute.FactTarget] = builds[0].target.Label
	}
	return refusalroute.Raised{Name: refusalroute.CandidateNotInstalled, Err: fmt.Errorf("%s: %w", refusalroute.CandidateNotInstalled, cause), Values: values}
}

// declaredBuildsOverlap raises the overlap refusal, which names the outcome and the spec of
// each open declared build of builds.
func declaredBuildsOverlap(builds []openBuild) error {
	open := make([]string, 0, len(builds))
	for _, build := range builds {
		open = append(open, build.claim.Outcome+" ("+build.deliverable+")")
	}
	return refusalroute.Raised{Name: refusalroute.DeclaredBuildsOverlap, Err: fmt.Errorf("%s: the declared builds %s are open, and one primary executable serves one installed candidate", refusalroute.DeclaredBuildsOverlap, strings.Join(open, ", "))}
}

// readOpenBuilds derives the open builds from one strict read of the ledger at the store
// root.
func (store Store) readOpenBuilds() ([]openBuild, error) {
	ledger, err := store.ledger().Read()
	if err != nil {
		return nil, err
	}
	return store.openBuilds(ledger)
}

// AdmitWindow decides the build window for source from one ledger read. The landing asks it
// before the lease, and PublishAdmitted decides again under the intent lock.
func (store Store) AdmitWindow(source Publication) error {
	ledger, err := store.ledger().Read()
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
