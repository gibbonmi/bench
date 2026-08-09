package gate

// A component whose phase the run never launches. An optional phase whose tool this host
// lacks skips green without grading anything, and the failure class these rows exist for is
// the green tail authoring that component's slot anyway: the slot's whole meaning is "some
// run graded this component at this identity", and a run that never started the phase graded
// nothing. The declared inputs of an absent tool's component do not move when the tool is
// installed, so a slot authored here would answer forever and the check would never run.
//
// Evidence is read from the slot bytes on disk and execution from the durable phase marker,
// never from a return value that only claims what happened.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/canary"
)

// absentToolComponent is the component the fixture declares as an optional phase pointing at
// a tool this host does not have. shellcheck is the kit table's own optional phase, so the
// fixture reproduces the shape the real gate carries rather than inventing a component.
const absentToolComponent = "shellcheck"

// absentToolPath is where the fixture's optional phase looks for its tool. It sits inside the
// root but outside every component's declared input set, so a test can install the tool
// without moving any component's identity — which is exactly the position a developer is in
// when they install shellcheck after a green run.
const absentToolPath = ".bench/tools/" + absentToolComponent

// absentToolScript is the body the installed tool carries. It appends the phase marker every
// other fixture phase appends, so an installed tool that ran is observable the same way.
const absentToolScript = "#!/usr/bin/env bash\necho " + absentToolComponent + " >> .git/phase-runs\n"

// gradedControlComponent is a component the fixture genuinely runs. It anchors the control
// row: without one, withholding every slot would satisfy the two rows above it.
const gradedControlComponent = canary.PhaseVet

// newAbsentToolFixture is a kit-shaped root whose absentToolComponent phase is declared
// optional with an argv naming a tool that is not installed. Its marker script goes with the
// declaration: the resolved gate script discovers its phases by globbing those scripts, so a
// leftover script would leave a marker for a phase the table says the host cannot run.
func newAbsentToolFixture(t *testing.T) kitShapedFixture {
	t.Helper()
	fixture := newKitShapedFixture(t)
	var doc manifestDoc
	for _, phase := range fixture.phases {
		declaration := manifestPhase{Name: phase.Name, Argv: phase.Argv, Needs: phase.Needs}
		if phase.Name == absentToolComponent {
			declaration.Argv = []string{filepath.Join(fixture.root, filepath.FromSlash(absentToolPath))}
			declaration.Optional = true
		}
		doc.Phases = append(doc.Phases, declaration)
	}
	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	writeGateTestFile(t, fixture.root, canary.PhaseManifestPath, string(data)+"\n", 0o644)
	if err := os.Remove(filepath.Join(fixture.root, ".bench", "phase-"+absentToolComponent+".sh")); err != nil {
		t.Fatal(err)
	}
	phases, err := phaseTable(fixture.root, fixture.root)
	if err != nil {
		t.Fatalf("re-resolve the fixture phase table: %v", err)
	}
	fixture.phases = phases
	if !carriesPhase(fixture.phases, absentToolComponent) || !carriesPhase(fixture.phases, gradedControlComponent) {
		t.Fatalf("resolved table = %v, want both %q and %q", fixture.phaseNames(), absentToolComponent, gradedControlComponent)
	}
	return fixture
}

// componentHasSlot reports whether the store answers for component at its current identity.
func componentHasSlot(t *testing.T, root, component string) bool {
	t.Helper()
	identity := mustResolveComponentIdentities(t, root)[component]
	if identity == "" {
		t.Fatalf("no identity for %q; the row below would observe nothing", component)
	}
	return resolveComponentSlot(root, component, identity, time.Now().UTC()).Skippable
}

// [PS48] A green run authors no slot for a component whose optional phase found no tool to
// run, and stays green. The gate's optional-phase idiom is untouched — the absent tool is
// still not a red — but it leaves no evidence behind, because nothing graded the component.
func TestAbsentOptionalToolLeavesNoSlot(t *testing.T) {
	t.Parallel()
	fixture := newAbsentToolFixture(t)
	mustExecuteGreen(t, fixture.root, productionGateEngine{})

	if componentHasSlot(t, fixture.root, absentToolComponent) {
		t.Fatalf("a slot answers for %q, but its tool is absent so the run never launched the phase", absentToolComponent)
	}
}

// [PS49][PS52] Installing the tool is enough to make the component run, and from there it is
// graded like any other. Its declared inputs do not move when the tool appears, so the only
// thing that can put the phase back on the gate is the absence of evidence to inherit; the run
// that grades it then authors the slot, and the run after that skips on it.
//
// Those last two states are the positive half of the rule, and they are what stops the rule
// from reading optionality itself as absence. A component held permanently ungradable is
// fail-closed, so it would ship quietly having deleted its whole share of the cost win with no
// diagnostic anywhere. The row follows one optional component across all three states rather
// than borrowing a required component's slot to stand in for the graded case, which is the
// substitution that would leave "optional" and "absent" indistinguishable.
func TestInstalledOptionalToolIsGradedLikeAnyOtherComponent(t *testing.T) {
	t.Parallel()
	fixture := newAbsentToolFixture(t)
	mustExecuteGreen(t, fixture.root, productionGateEngine{})
	before := mustResolveComponentIdentities(t, fixture.root)[absentToolComponent]

	writeGateTestFile(t, fixture.root, absentToolPath, absentToolScript, 0o755)
	graded := observeGreenGate(t, fixture.root)

	if after := mustResolveComponentIdentities(t, fixture.root)[absentToolComponent]; after != before {
		t.Fatalf("%s identity moved from %s to %s; installing the tool must not move its declared inputs", absentToolComponent, before, after)
	}
	if !graded.ran(absentToolComponent) {
		t.Fatalf("%s did not run; the executed set was %v and the gate said:\n%s", absentToolComponent, graded.executed, graded.stdout)
	}
	if rec := partialRecord(t, fixture.root); slicesContains(rec.Skipped, absentToolComponent) {
		t.Fatalf("recorded skips = %v, want %s absent — it had no evidence to inherit", rec.Skipped, absentToolComponent)
	}
	if !componentHasSlot(t, fixture.root, absentToolComponent) {
		t.Fatalf("the run that graded %s authored no slot; withholding evidence has become disabling the component", absentToolComponent)
	}

	for _, path := range captureSurfacePaths(ReducedScope()) {
		writeGateTestFile(t, fixture.root, path, "capture-only edit\n", 0o644)
	}
	skipping := observeGreenGate(t, fixture.root)

	if skipping.ran(absentToolComponent) {
		t.Fatalf("%s ran again over a capture-only changeset; the slot it earned covers it, so an installed optional tool is paying for every run forever", absentToolComponent)
	}
	if rec := partialRecord(t, fixture.root); !slicesContains(rec.Skipped, absentToolComponent) {
		t.Fatalf("recorded skips = %v, want %s among them", rec.Skipped, absentToolComponent)
	}
}

// [PS50] A component whose phase genuinely ran green still gets its slot, and still skips on
// the next run. This is the control: withholding every component's evidence would satisfy the
// two rows above while silently taking per-component scoping off the gate entirely.
func TestGradedComponentKeepsItsSlotAndSkips(t *testing.T) {
	t.Parallel()
	fixture := newAbsentToolFixture(t)
	mustExecuteGreen(t, fixture.root, productionGateEngine{})

	if !componentHasSlot(t, fixture.root, gradedControlComponent) {
		t.Fatalf("the green run authored no slot for %q, which it executed", gradedControlComponent)
	}
	for _, path := range captureSurfacePaths(ReducedScope()) {
		writeGateTestFile(t, fixture.root, path, "capture-only edit\n", 0o644)
	}
	observation := observeGreenGate(t, fixture.root)

	if observation.ran(gradedControlComponent) {
		t.Fatalf("%s ran again over a capture-only changeset; its own slot covers it", gradedControlComponent)
	}
	if rec := partialRecord(t, fixture.root); !slicesContains(rec.Skipped, gradedControlComponent) {
		t.Fatalf("recorded skips = %v, want %s among them", rec.Skipped, gradedControlComponent)
	}
}

// buildArtifactSourceFile is the one file allowed to derive the published binary's path.
const buildArtifactSourceFile = "component_decision.go"

// [PS51] buildArtifactPath is the package's only derivation of the alternate-package
// proof artifact. The attestation is addressed by that path, so a second derivation
// that normalized differently would have the author write where the reader never looks — and
// a second spelling of the same fact is the class this repository's code standard refuses
// whether or not today's spellings happen to agree.
func TestBuildArtifactPathIsTheOnlyArtifactDerivation(t *testing.T) {
	t.Parallel()
	sources, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	var derivations []string
	for _, source := range sources {
		if strings.HasSuffix(source, "_test.go") {
			continue
		}
		body, err := os.ReadFile(source)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(body), `"dist", "bench"`) {
			derivations = append(derivations, source)
		}
	}
	if len(derivations) != 1 || derivations[0] != buildArtifactSourceFile {
		t.Fatalf("files deriving the artifact path = %v, want only %s — the attestation is addressed by one spelling", derivations, buildArtifactSourceFile)
	}
}
