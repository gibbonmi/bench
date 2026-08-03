package gate

// The per-component scoping decision. One gate run, one verdict, one phase ordering —
// inside it, each component whose declared inputs the tree has not moved since some run
// graded it skips on its own ancestor slot, and every other phase runs.
//
// This is the one site where a skip is decided, and that is the whole point of the file.
// Every path out of scopeComponents that is not a validated slot at the component's exact
// current identity answers run-the-component: a root that is not the kit, a gate that does
// not route the phase table, a derivation that failed, an identity that could not be
// computed, a slot that could not be read or that answers for something else. Fail-closed
// is then a property of one function rather than a claim repeated at seven call sites, and
// the quantifier "every error runs the component" is checkable by reading one file.
//
// Nothing outside the run selects. No flag, caller, agent, or session names components to
// run or skip; the decision reads the tree's content identities against the retained
// evidence and nothing else. The one operator lever is `bench gate --fresh`, which arrives
// here as forceRun and buys more grading, never less.

import (
	"path/filepath"
	"sort"
	"time"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	// The package's own freshness constant owns the bare name, so the seal package is
	// reached through an alias.
	benchfreshness "github.com/gibbonmi/bench/internal/freshness"
)

// componentScoping is one execution's per-component answer: the identity of every
// component this root may scope, the phases a narrowed run executes, and the components it
// skips with the evidence covering each.
//
// identities is populated whenever the root scopes at all, including under forceRun, where
// no component skips: a green run has to author a slot for every component it executed,
// and that authorship is what a later run inherits from. It names only the components this
// run can grade — a phase the run cannot launch is absent from it, so nothing downstream
// can offer that component a skip or credit it with evidence.
type componentScoping struct {
	// eligible reports that the per-component declarations reach this root at all, which is
	// a question about the root's shape and never about whether a derivation succeeded on
	// it. It stays true through every refusal below, so a caller can tell "these
	// declarations are not this root's" from "they are, and something failed".
	eligible   bool
	runnerRoot string
	identities map[string]string
	checks     conformanceCheckPartition
	phases     []Phase
	skipped    []ComponentSkip
	// buildArtifact is the binary this root's build phase publishes, set only when the
	// resolved table carries that phase. It is the one spelling both halves of the
	// attestation use — the skip decision reads it and the green tail authors against it —
	// because the record is addressed by the artifact's path, and two spellings would have
	// the author write where the reader never looks.
	buildArtifact string
}

// partial reports whether this run graded less than every component. The announcement, the
// recorded partition, and the evidence a green run withholds all key off this one answer,
// so a run is narrow in all three senses or in none.
func (s componentScoping) partial() bool {
	return len(s.skipped) > 0 || len(s.checks.Inherited) > 0
}

// executedScopedComponents are the components this run graded for itself and whose evidence
// is an ancestor slot — the set a green run authors slots for and a red run retires them
// for. build is excluded because its evidence is the attested seal beside the artifact
// rather than a slot, and the unconditional phases are absent from identities entirely.
func (s componentScoping) executedScopedComponents() []string {
	skipped := make(map[string]bool, len(s.skipped))
	for _, skip := range s.skipped {
		skipped[skip.Component] = true
	}
	executed := make([]string, 0, len(s.identities))
	for component := range s.identities {
		if !skipped[component] && componentSkipsOnAncestorEvidence(component) {
			executed = append(executed, component)
		}
	}
	sort.Strings(executed)
	return executed
}

// executedBuild reports that this run's own build phase produced the binary now beside the
// root — the condition under which a green run attests it. A root whose table carries no
// build phase has no artifact to answer for, and a run that skipped the build produced
// nothing, so re-attesting there would credit the gate with a binary it did not make.
func (s componentScoping) executedBuild() bool {
	if s.buildArtifact == "" {
		return false
	}
	for _, skip := range s.skipped {
		if skip.Component == canary.PhaseBuild {
			return false
		}
	}
	return true
}

// executedPhaseNames are the names of the phases a narrowed run executes, sorted and
// deduplicated as the partial record class requires.
func (s componentScoping) executedPhaseNames() []string {
	names := make([]string, 0, len(s.phases))
	for _, phase := range s.phases {
		names = append(names, phase.Name)
	}
	sort.Strings(names)
	return names
}

// componentSkipsOnEvidence reports whether any evidence at all can cover this component.
//
// conformance and the conformance suite enforce the input declarations, the slot class, and
// the profile bindings this scoping rests on, so they are never skippable: a surface that
// could skip its own enforcement would let a drifted declaration certify itself.
func componentSkipsOnEvidence(component string) bool {
	switch component {
	case conformancePhaseName:
		return false
	}
	return true
}

// componentSkipsOnAncestorEvidence reports whether a skippable component's evidence is its own
// ancestor slot. build is the one that skips on something else — it produces the binary the
// other phases exec, and its evidence is the attested seal beside that artifact — so it is
// also the one component this run neither authors nor retires a slot for.
func componentSkipsOnAncestorEvidence(component string) bool {
	return component != canary.PhaseBuild && componentSkipsOnEvidence(component)
}

// scopeComponents decides, for root's resolved table, which components may skip on their own
// evidence. The returned scoping is empty for a root that does not scope at all.
//
// Eligibility comes before any per-component question, and all three of its parts are about
// whose declaration is being applied rather than about how this run went. The input
// declarations are the kit's own — compiled from the kit's source, bound to its profile, and
// conformance-checked only in its tree — so only the kit's own root may inherit against them;
// identity rather than spelling, so a symlinked path still counts and any stat failure
// refuses. The resolved gate must provably route through the phase table, because a narrowed
// run executes that table directly: narrowing under a hand-written gate would swap the
// repository's oracle for one it never chose. And every declaration is derived from a Go
// module, so a root without one has nothing here to be scoped by.
func scopeComponents(root string, res Resolution, mode runMode, now time.Time) componentScoping {
	kit := kitRoot(root)
	if !sameDirectory(root, kit) || !phaseTableGate(root, res) || !isRegularFile(filepath.Join(root, "go.mod")) {
		return componentScoping{}
	}
	eligible := componentScoping{eligible: true, runnerRoot: kit}
	table, err := phaseTable(root, kit)
	if err != nil {
		return eligible
	}
	checkIdentities, checkIdentityErr := ResolveConformanceCheckIdentities(root, registry.Dev)
	canaryIdentities, canaryIdentityErr := resolveConformanceCanaryIdentities(root, registry.Dev)
	if mode == reuseFreshGreen {
		eligible.checks = partitionConformanceChecks(root, registry.Dev, checkIdentities, checkIdentityErr, now)
	} else {
		eligible.checks = executeAllConformanceChecks(registry.Dev, checkIdentities)
		if checkIdentityErr != nil {
			eligible.checks.Identities = nil
		}
	}
	decorateConformanceCanarySelection(root, &eligible.checks, canaryIdentities, canaryIdentityErr)
	// One resolution for the whole family, and a failure anywhere in it — a listing that
	// would not run, a seal that could not be read, a declared input the snapshot has no
	// entry for — takes every component with it. A partial set of identities names fewer
	// inputs than the components read, which is exactly the shape that buys a wrong skip.
	identities, err := ResolveComponentIdentities(root)
	if err != nil {
		eligible.checks = executeAllConformanceChecks(registry.Dev, checkIdentities)
		if checkIdentityErr != nil {
			eligible.checks.Identities = nil
		}
		eligible.checks.Canary = canaryIdentities
		return eligible
	}
	scoping := componentScoping{eligible: true, runnerRoot: kit, identities: gradableIdentities(table, identities), checks: eligible.checks}
	// The artifact is named from the resolved table rather than from the declarations: a
	// component the registry declares but this tree materializes no phase for has nothing
	// here to skip and nothing to attest.
	if carriesPhase(table, canary.PhaseBuild) {
		scoping.buildArtifact = buildArtifactPath(root)
	}
	if mode != reuseFreshGreen {
		return scoping
	}
	phases := make([]Phase, 0, len(table))
	var skipped []ComponentSkip
	var selectedCanaryFamilies []string
	for _, phase := range table {
		identity, scoped := identities[phase.Name]
		if !scoped {
			phases = append(phases, phase)
			continue
		}
		skip, skippable := componentSkip(root, scoping, phase.Name, identity, now)
		if phase.Name == "canary" && skippable && (scoping.checks.CanaryFull || len(scoping.checks.CanaryFamilies) > 0) {
			if !scoping.checks.CanaryFull {
				selectedCanaryFamilies = scoping.checks.CanaryFamilies
			}
			phases = append(phases, phase)
			continue
		}
		if !skippable {
			phases = append(phases, phase)
			continue
		}
		skipped = append(skipped, skip)
	}
	// Nothing skipped makes a narrowed run the same work under a narrower record, and
	// nothing left to run would green on a run of no phases. Both fall back, keeping the
	// identities so the green tail still authors what it executed.
	if (len(skipped) == 0 && len(scoping.checks.Inherited) == 0) || len(phases) == 0 {
		return scoping
	}
	sort.Slice(skipped, func(i, j int) bool { return skipped[i].Component < skipped[j].Component })
	scoping.phases = withConformanceCheckSelection(phases, registry.Dev, scoping.checks.Executed, scoping.checks.verdictInherited())
	scoping.phases = withCanaryFamilySelection(scoping.phases, selectedCanaryFamilies)
	scoping.skipped = skipped
	return scoping
}

// gradableIdentities drops every component whose phase this run cannot launch. An optional
// phase whose tool the host lacks settles green having graded nothing, so its identity is not
// one this run may answer for: leaving it in would have the green tail author a slot claiming
// the component was graded, and because installing the tool moves none of that component's
// declared inputs, the slot would answer every run afterwards and the check would never run
// again. That is the failure a slot exists to make impossible, and it costs a red nothing to
// prevent — the phase is free either way.
//
// The drop is here rather than at authorship so one answer covers the whole decision: such a
// component cannot skip on evidence, authors none, and retires none. It also reaches the run
// shapes that never touch this package's runner, because the resolved gate executes the phase
// table in a child process whose results the green tail cannot see.
//
// A capability- or environment-skipped case inside a phase is deliberately not covered. That
// phase launched and graded its component; individual cases the host could not populate are
// reported on their own `capability-skips` rows and can be made fatal, so the component's
// grade stands as the run's honest answer for it. Withholding evidence there would strip the
// slot of every component on any host with an incomplete capability population, which is
// every developer's, and the tally is per-run rather than per-component, so there is no sound
// component to withhold it from.
func gradableIdentities(table []Phase, identities map[string]string) map[string]string {
	gradable := make(map[string]string, len(identities))
	for component, identity := range identities {
		gradable[component] = identity
	}
	for _, phase := range table {
		if _, absent := phaseToolAbsent(phase); absent {
			delete(gradable, phase.Name)
		}
	}
	return gradable
}

// componentSkip answers one component's skip question and, where the answer is yes, the
// evidence entry that covered it. Both evidence forms are reached from here, so the loop
// above stays one pass over the table: every scoped component but build resolves its own
// ancestor slot, while build reuses the artifact beside the root and is answered by the seal
// and the attestation instead. A component that skips on neither runs.
func componentSkip(root string, scoping componentScoping, component, identity string, now time.Time) (ComponentSkip, bool) {
	if !componentSkipsOnEvidence(component) {
		return ComponentSkip{}, false
	}
	if component == canary.PhaseBuild {
		return attestedBuildSkip(root, scoping.buildArtifact, now)
	}
	inspection := resolveComponentSlot(root, component, identity, now)
	if !inspection.Skippable {
		return ComponentSkip{}, false
	}
	return ComponentSkip{Component: component, Identity: identity, AuthoredAt: inspection.AuthoredAt}, true
}

// attestedBuildSkip is build's whole skip decision, and the only one here answered by an
// artifact rather than by a slot.
//
// Both conjuncts are required and neither implies the other. freshness.Check says the binary
// on disk is the one its own seal describes and that the seal was taken from the sources
// present now; the attestation says a gate build produced exactly those bytes. Check alone
// skips on a planted binary republished with its own recomputed seal, which is the whole case
// the attestation exists for; the attestation alone skips on a gate-built binary whose sources
// have since moved, which is a stale artifact every reader would then exec. No clock and no
// mtime enters either half — the question is only ever about content.
//
// Every refusal runs the build, which republishes the seal and re-authors the attestation
// together, so nothing here can leave the artifact permanently unskippable.
func attestedBuildSkip(root, executable string, now time.Time) (ComponentSkip, bool) {
	if executable == "" {
		return ComponentSkip{}, false
	}
	if err := benchfreshness.Check(root, executable); err != nil {
		return ComponentSkip{}, false
	}
	if !verifyBuildAttestation(root, executable, now).Attested {
		return ComponentSkip{}, false
	}
	// The recorded evidence is the source digest rather than the executable one: what a
	// reader of the verdict needs to know is which build inputs the reused artifact answers
	// for, and the executable digest names bytes that say nothing about the tree.
	sources, _, err := benchfreshness.SealDigests(executable)
	if err != nil || !isContentAddress(sources) {
		return ComponentSkip{}, false
	}
	return ComponentSkip{Component: canary.PhaseBuild, Seal: sources}, true
}

// buildArtifactPath is where root's build phase publishes its binary. The phase's own argv,
// the skip decision, and the attestation's address are all taken from here. The attestation is
// keyed by this path's cleaned absolute spelling, so a second derivation that normalized
// differently would have the author write to an address the reader never looks at — which does
// not fail loudly, it just silently pays for a build phase on every run forever.
func buildArtifactPath(root string) string { return filepath.Join(root, "dist", "bench") }

// announcement is how one skip reports itself to an operator. The two evidence forms read
// differently because they are different claims: an ancestor slot names a run that graded this
// component, while build's names an artifact a gate build produced and a seal that ties it to
// the current sources. Which one a skip rests on is what an operator would check.
func (s ComponentSkip) announcement() string {
	if s.Seal != "" {
		return "gate: skipping " + s.Component + " (artifact reused from the gate-attested build sealed at " + s.Seal + ")"
	}
	return "gate: skipping " + s.Component + " (evidence inherited from ancestor " + s.Identity +
		" recorded " + s.AuthoredAt.UTC().Format(time.RFC3339) + ")"
}

// evidence is the same skip as the verdict record carries it. It sits beside announcement so
// the printed line and the recorded entry are one derivation of one skip's evidence form:
// a verdict naming evidence the announcement did not is exactly the discrepancy an operator
// checking a skip has no way to resolve.
func (s ComponentSkip) evidence() skipEvidence {
	if s.Seal != "" {
		return skipEvidence{Seal: s.Seal}
	}
	return skipEvidence{
		Identity:   s.Identity,
		AuthoredAt: s.AuthoredAt.UTC().Truncate(time.Second).Format(time.RFC3339),
	}
}

// attestExecutedBuild records that this green run's own build phase produced the binary now
// beside the root.
//
// The digest is taken from the artifact's bytes, never from the seal the build just wrote: a
// digest read back from a seal would only restate what the build script put on disk, and being
// the one half no seal writer can supply is the entire point of the record.
func attestExecutedBuild(scoping componentScoping, root string, authoredAt time.Time) error {
	if !scoping.executedBuild() {
		return nil
	}
	digest, err := benchfreshness.ExecutableDigest(scoping.buildArtifact)
	if err != nil {
		return err
	}
	return authorBuildAttestation(root, scoping.buildArtifact, digest, authoredAt)
}

// authorExecutedComponentSlots records that this green run graded each component it executed
// at that component's current identity. A component this run skipped is left alone: the run
// graded nothing about it, and re-stamping its slot would dress older evidence as newer.
func authorExecutedComponentSlots(root string, scoping componentScoping, authoredAt time.Time) error {
	for _, component := range scoping.executedScopedComponents() {
		if err := authorComponentSlot(root, component, scoping.identities[component], authoredAt); err != nil {
			return err
		}
	}
	return nil
}

// invalidateExecutedComponentSlots retires the slot of every component this red run executed,
// at that component's current identity. A component the run skipped keeps its slot: none of
// what ran grades what that slot answers for, and retiring it would re-charge the next run
// for a red it does not own.
//
// The set is the executed one rather than the single red component because the run's own
// verdict is the finest attribution available: the resolved gate runs the phase table in a
// child process, so a forced whole-tree run — the one case where a component executes while a
// green slot already stands at its current identity — cannot say which phase was red. Leaving
// those slots up would let the very next run skip a component this run just reddened.
func invalidateExecutedComponentSlots(root string, scoping componentScoping) error {
	for _, component := range scoping.executedScopedComponents() {
		if err := invalidateComponentSlot(root, component, scoping.identities[component]); err != nil {
			return err
		}
	}
	return nil
}
