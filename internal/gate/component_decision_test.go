package gate

// The per-component decision observed from outside: which phases a run executed, which it
// skipped, what evidence each skip named, and what the store held afterwards. The failure
// class here is a verdict that credits work nobody graded, so execution is read from the
// durable phase marker and evidence from the slot bytes on disk — never from a return value
// that only claims what happened.

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/canary"
)

// gateObservation is one execution's externally visible account: the phases whose marker
// scripts ran during it, everything it printed, and the exit it reported.
type gateObservation struct {
	executed []string
	stdout   string
	exit     int
}

func (o gateObservation) ran(phase string) bool { return slicesContains(o.executed, phase) }

func slicesContains(names []string, want string) bool {
	for _, name := range names {
		if name == want {
			return true
		}
	}
	return false
}

// observeGate runs one gate execution against root and reports what it did. The marker delta
// is taken across the call, so an observation names this run's phases and not the seed's.
func observeGate(t *testing.T, root string) gateObservation {
	t.Helper()
	before := len(phaseRunNames(t, root))
	var stdout bytes.Buffer
	result := executeWithEngineAtKit(context.Background(), root, root, &stdout, io.Discard, productionGateEngine{})
	executed := append([]string(nil), phaseRunNames(t, root)[before:]...)
	sort.Strings(executed)
	return gateObservation{executed: executed, stdout: stdout.String(), exit: result.ActionExit}
}

func observeGreenGate(t *testing.T, root string) gateObservation {
	t.Helper()
	observation := observeGate(t, root)
	if observation.exit != 0 {
		t.Fatalf("execution = %+v, want green", observation)
	}
	return observation
}

// seededScopingFixture is a kit-shaped root whose first full green has authored a slot for
// every scoped component. Every row below starts here: with no slots, nothing can skip, and
// a scoping assertion written against that root would observe an ordinary full run.
func seededScopingFixture(t *testing.T) kitShapedFixture {
	t.Helper()
	fixture := newKitShapedFixture(t)
	mustExecuteGreen(t, fixture.root, productionGateEngine{})
	for component, identity := range mustResolveComponentIdentities(t, fixture.root) {
		if !componentSkipsOnAncestorEvidence(component) {
			continue
		}
		if !resolveComponentSlot(fixture.root, component, identity, time.Now().UTC()).Skippable {
			t.Fatalf("the seed run authored no slot for %q; nothing below can observe a skip", component)
		}
	}
	return fixture
}

// seededDecisionFixture authors the prerequisite evidence without executing the phases
// whose fail-closed decisions the caller observes.
func seededDecisionFixture(t *testing.T) kitShapedFixture {
	t.Helper()
	fixture := newKitShapedFixture(t)
	authoredAt := time.Now().UTC()
	for component, identity := range mustResolveComponentIdentities(t, fixture.root) {
		if !componentSkipsOnAncestorEvidence(component) {
			continue
		}
		if err := authorComponentSlot(fixture.root, component, identity, authoredAt); err != nil {
			t.Fatalf("author the %q decision fixture slot: %v", component, err)
		}
		if inspection := resolveComponentSlot(fixture.root, component, identity, authoredAt); !inspection.Skippable {
			t.Fatalf("read the %q decision fixture slot = %+v, want skippable", component, inspection)
		}
	}

	return fixture
}

// slotBytes maps every scoped component to the bytes its slot holds at that component's
// current identity, and to nil where the store holds none. Comparing these maps across a run
// is what "left byte-identical" means: a re-stamped slot and a retired one both show up.
func slotBytes(t *testing.T, root string) map[string][]byte {
	t.Helper()
	dir, err := componentSlotDir(root)
	if err != nil {
		t.Fatal(err)
	}
	held := map[string][]byte{}
	for component, identity := range mustResolveComponentIdentities(t, root) {
		if !componentSkipsOnAncestorEvidence(component) {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, componentSlotName(component, identity)))
		if err != nil {
			if !os.IsNotExist(err) {
				t.Fatal(err)
			}
			held[component] = nil
			continue
		}
		held[component] = data
	}
	return held
}

// forcePhaseRed repoints a fixture phase's marker script at a red body. The script is in no
// component's declared inputs, so the red arrives without moving any identity — which is the
// only way to red a component the decision would otherwise let skip.
func forcePhaseRed(t *testing.T, root, phase string) {
	t.Helper()
	writeGateTestFile(t, root, ".bench/phase-"+phase+".sh", "echo "+phase+" >> .git/phase-runs\nexit 1\n", 0o644)
}

func partialRecord(t *testing.T, root string) verdictRecord {
	t.Helper()
	rec := slotRecord(t, root, time.Now().UTC())
	if rec.partition() == nil {
		t.Fatalf("recorded verdict = %+v, want a partial verdict", rec)
	}
	return rec
}

// [PC1] A capture-only changeset pays the unconditional phases and skips every
// evidence-covered component on that component's own evidence. The recorded evidence is
// checked against each component's own identity, so a decision answering one question of the
// whole changeset — the shape the reduced run already has — cannot satisfy this row by
// naming one shared ancestor for all seven.
func TestCaptureOnlyChangesetExecutesConformanceOnly(t *testing.T) {
	t.Parallel()
	fixture := seededScopingFixture(t)
	for _, path := range captureSurfacePaths(ReducedScope()) {
		writeGateTestFile(t, fixture.root, path, "capture-only edit\n", 0o644)
	}
	identities := mustResolveComponentIdentities(t, fixture.root)
	observation := observeGreenGate(t, fixture.root)

	unconditional, skippable := unconditionalPhaseNames(fixture.phases)
	sort.Strings(unconditional)
	sort.Strings(skippable)
	if !reflect.DeepEqual(observation.executed, unconditional) {
		t.Fatalf("executed %v, want exactly the unconditional phases %v", observation.executed, unconditional)
	}
	rec := partialRecord(t, fixture.root)
	if !reflect.DeepEqual(rec.Skipped, skippable) {
		t.Fatalf("recorded skips = %v, want every evidence-covered component %v", rec.Skipped, skippable)
	}
	seen := map[string]string{}
	for _, component := range rec.Skipped {
		evidence := rec.SkipEvidence[component]
		if evidence.Identity != identities[component] {
			t.Fatalf("%s skipped on identity %q, want its own identity %q", component, evidence.Identity, identities[component])
		}
		if other, shared := seen[evidence.Identity]; shared {
			t.Fatalf("%s and %s skipped on one shared identity %q, want per-component evidence", component, other, evidence.Identity)
		}
		seen[evidence.Identity] = component
	}
}

// [PC2c] Every skip is announced on its own line naming the component, the ancestor identity
// it inherited, and when that ancestor was recorded, and the record carries one entry per
// skip. A single summary line tells an operator that something was skipped without letting
// them check what stood in for it.
func TestEverySkipIsAnnouncedWithItsEvidence(t *testing.T) {
	t.Parallel()
	fixture := seededScopingFixture(t)
	writeGateTestFile(t, fixture.root, "ROADMAP.md", "capture-only edit\n", 0o644)
	observation := observeGreenGate(t, fixture.root)
	rec := partialRecord(t, fixture.root)

	var announcements []string
	for _, line := range strings.Split(observation.stdout, "\n") {
		if strings.HasPrefix(line, "gate: skipping ") {
			announcements = append(announcements, line)
		}
	}
	if len(announcements) != len(rec.Skipped) {
		t.Fatalf("announced %d skips in\n%s\nwant one line for each of %v", len(announcements), observation.stdout, rec.Skipped)
	}
	for _, component := range rec.Skipped {
		evidence := rec.SkipEvidence[component]
		// Each evidence form gets its own line, because each asks the operator to trust a
		// different thing: a run that graded the component, or a build the gate itself ran.
		want := "gate: skipping " + component + " (evidence inherited from ancestor " + evidence.Identity +
			" recorded " + evidence.AuthoredAt + ")"
		if evidence.Seal != "" {
			want = "gate: skipping " + component + " (artifact reused from the gate-attested build sealed at " + evidence.Seal + ")"
		}
		if !slicesContains(announcements, want) {
			t.Fatalf("skip announcements %v are missing %q", announcements, want)
		}
	}
	if len(rec.SkipEvidence) != len(rec.Skipped) {
		t.Fatalf("recorded %d evidence entries for %d skips: %+v", len(rec.SkipEvidence), len(rec.Skipped), rec)
	}
}

// [PC11] A red run retires the slot of what it executed and leaves what it skipped alone. The
// components that skipped are the ones holding real slots at their current identities, so an
// implementation retiring every slot on any red re-charges the next run for a red none of
// those components own.
func TestRedComponentInvalidatesOnlyItsOwnSlot(t *testing.T) {
	t.Parallel()
	fixture := seededScopingFixture(t)
	// Moving the canary fixture surface moves canary's identity and no other's, so canary
	// is the one component with no slot to skip on when the run starts.
	writeGateTestFile(t, fixture.root, "tests/canary/fixture.txt", "moved canary fixture\n", 0o644)
	forcePhaseRed(t, fixture.root, "canary")
	before := slotBytes(t, fixture.root)

	observation := observeGate(t, fixture.root)
	if observation.exit == 0 {
		t.Fatalf("execution = %+v, want red", observation)
	}
	if !observation.ran("canary") {
		t.Fatalf("executed %v, want the moved canary component to run", observation.executed)
	}
	after := slotBytes(t, fixture.root)
	if after["canary"] != nil {
		t.Fatalf("canary's slot survived its own red: %q", after["canary"])
	}
	for component, data := range before {
		if component == "canary" {
			continue
		}
		if data == nil {
			t.Fatalf("%s held no slot before the red run; the row observes nothing about it", component)
		}
		if !bytes.Equal(data, after[component]) {
			t.Fatalf("%s's slot moved across another component's red:\nbefore %q\nafter  %q", component, data, after[component])
		}
	}
}

// [PC11b] A forced red retires the slots of the components it just reddened. This is the one
// shape where a component executes while a green slot already stands at its exact current
// identity — an ordinary run reaches a component only when its identity moved past whatever
// slot it had — so leaving those slots up lets the very next run skip a component `--fresh`
// just proved red.
func TestForcedRedRetiresTheSlotsItExecuted(t *testing.T) {
	t.Parallel()
	fixture := seededScopingFixture(t)
	forcePhaseRed(t, fixture.root, canary.PhaseVet)
	before := slotBytes(t, fixture.root)
	for component, data := range before {
		if data == nil {
			t.Fatalf("%s held no slot before the forced run; the row observes nothing", component)
		}
	}

	if got := executeWithEngineAfterAcquireAtKit(context.Background(), fixture.root, fixture.root, io.Discard, io.Discard, productionGateEngine{}, notifyGateSignals, forceRun); got.ActionExit == 0 {
		t.Fatalf("forced run exit = %d, want red", got.ActionExit)
	}
	for component, data := range slotBytes(t, fixture.root) {
		if data != nil {
			t.Fatalf("%s kept its slot through a red forced run: %q", component, data)
		}
	}
	// The point of the retirement, stated as behavior: the next ordinary run re-grades
	// everything rather than skipping on evidence the forced red contradicted.
	writeGateTestFile(t, fixture.root, ".bench/phase-"+canary.PhaseVet+".sh",
		"# repaired\necho "+canary.PhaseVet+" >> .git/phase-runs\n", 0o644)
	observation := observeGreenGate(t, fixture.root)
	if !observation.ran(canary.PhaseVet) {
		t.Fatalf("executed %v after the forced red, want %s re-graded", observation.executed, canary.PhaseVet)
	}
}

// [PC12] Editing a canary fixture runs the canary and skips the toolchain components. This is
// the decision's worked example: one changeset the components disagree about, which no
// whole-changeset predicate can express.
func TestCanaryFixtureEditRunsCanary(t *testing.T) {
	t.Parallel()
	fixture := seededScopingFixture(t)
	writeGateTestFile(t, fixture.root, "tests/canary/fixture.txt", "moved canary fixture\n", 0o644)
	observation := observeGreenGate(t, fixture.root)

	if !observation.ran("canary") {
		t.Fatalf("executed %v, want canary run for an edit under tests/canary/", observation.executed)
	}
	for _, component := range []string{canary.PhaseGofmt, canary.PhaseVet, canary.PhaseTest} {
		if observation.ran(component) {
			t.Fatalf("executed %v, want %s skipped — the canary fixture is not its input", observation.executed, component)
		}
	}
}

// [PC13] Editing the wrapper script the canary phase execs runs the canary. Canary's inputs
// are hand-declared, and the cheapest declaration — its two directories — is blind to the
// wiring that decides what the phase actually runs.
func TestWrapperEditRunsCanary(t *testing.T) {
	t.Parallel()
	fixture := seededScopingFixture(t)
	writeGateTestFile(t, fixture.root, "bin/bench.sh", "#!/usr/bin/env bash\n# moved wrapper\nexit 0\n", 0o755)
	observation := observeGreenGate(t, fixture.root)

	if !observation.ran("canary") {
		t.Fatalf("executed %v, want canary run for an edit to bin/bench.sh", observation.executed)
	}
}

// [PC14] An ordinary Go edit runs the toolchain components and skips the canary. The binary
// is deliberately absent from canary's declaration, so this pins the ruled narrowing against
// a later edit widening canary's inputs back to the artifact it execs.
func TestOrdinaryGoEditSkipsCanary(t *testing.T) {
	t.Parallel()
	fixture := seededScopingFixture(t)
	writeGateTestFile(t, fixture.root, "cmd/bench/main.go", "package main\n\n// moved source\nfunc main() {}\n", 0o644)
	observation := observeGreenGate(t, fixture.root)

	if observation.ran("canary") {
		t.Fatalf("executed %v, want canary skipped for an ordinary Go edit", observation.executed)
	}
	for _, component := range []string{canary.PhaseGofmt, canary.PhaseVet, canary.PhaseTest} {
		if !observation.ran(component) {
			t.Fatalf("executed %v, want %s run for an ordinary Go edit", observation.executed, component)
		}
	}
}

func TestAgentMarkdownEditRunsConsumersAndSkipsToolchain(t *testing.T) {
	t.Parallel()
	fixture := seededScopingFixture(t)
	writeGateTestFile(t, fixture.root, ".agents/skills/new-skill/SKILL.md", "# New skill\n", 0o644)
	observation := observeGreenGate(t, fixture.root)

	for _, component := range []string{canary.PhaseContract, "canary"} {
		if !observation.ran(component) {
			t.Fatalf("executed %v, want %s run for an agent Markdown edit", observation.executed, component)
		}
	}
	for _, component := range []string{
		canary.PhaseGofmt, canary.PhaseVet, canary.PhaseTest, canary.PhaseRace,
	} {
		if observation.ran(component) {
			t.Fatalf("executed %v, want %s skipped for an agent Markdown edit", observation.executed, component)
		}
	}
}

// [PC19] Every error class at the decision site answers run-the-component. One fail-open path
// anywhere here silently credits ungraded work, and the classes are exercised one per subtest
// so a refusal that stops refusing is named rather than absorbed by a sibling.
func TestDecisionSiteFailsClosed(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		affects string
		inject  func(t *testing.T, fixture kitShapedFixture)
	}{
		{"slot unreadable", canary.PhaseVet, func(t *testing.T, fixture kitShapedFixture) {
			// The store's file discipline refuses anything but a private regular file, so a
			// widened mode leaves bytes the reader will not answer from.
			if err := os.Chmod(scopedSlotPath(t, fixture.root, canary.PhaseVet), 0o644); err != nil {
				t.Fatal(err)
			}
			writeGateTestFile(t, fixture.root, "ROADMAP.md", "capture-only edit\n", 0o644)
		}},
		{"derivation failure", canary.PhaseVet, func(t *testing.T, fixture kitShapedFixture) {
			writeGateTestFile(t, fixture.root, "internal/broken/broken.go", "package broken\n\nfunc (\n", 0o644)
		}},
		{"identity failure", "canary", func(t *testing.T, fixture kitShapedFixture) {
			// A declared file entry the snapshot has no entry for is the declaration and the
			// tree disagreeing at a named point.
			if err := os.Remove(filepath.Join(fixture.root, "bin", "bench.sh")); err != nil {
				t.Fatal(err)
			}
		}},
		{"domain mismatch", canary.PhaseVet, func(t *testing.T, fixture kitShapedFixture) {
			plantForeignSlot(t, fixture.root, canary.PhaseVet, canary.PhaseTest)
			writeGateTestFile(t, fixture.root, "ROADMAP.md", "capture-only edit\n", 0o644)
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fixture := seededDecisionFixture(t)
			tc.inject(t, fixture)
			scoping := mustScopeComponents(t, fixture.root, Resolve(fixture.root, "", RealFS()), reuseFreshGreen, time.Now().UTC())
			executed := fixture.phaseNames()
			if scoping.partial() {
				executed = scoping.executedPhaseNames()
			}
			if !slicesContains(executed, tc.affects) {
				t.Fatalf("decision executed %v and skipped %v, want %s run past a %s", executed, scoping.skipped, tc.affects, tc.name)
			}
		})
	}
}

// [PS29] `bench gate --fresh` grades everything and re-authors what it graded, and an
// unedited tree is still answered by the whole-tree reuse before any component is asked.
func TestFreshExecutesEveryComponent(t *testing.T) {
	t.Parallel()
	t.Run("the whole-tree reuse answers ahead of the decision", func(t *testing.T) {
		fixture := seededScopingFixture(t)
		observation := observeGreenGate(t, fixture.root)
		if len(observation.executed) != 0 {
			t.Fatalf("unedited tree executed %v, want the retained whole-tree green reused", observation.executed)
		}
		if !strings.Contains(observation.stdout, "fresh verdict reused") {
			t.Fatalf("unedited tree said:\n%s\nwant the reused-verdict announcement", observation.stdout)
		}
	})

	t.Run("the forced run executes every phase past every slot", func(t *testing.T) {
		fixture := seededScopingFixture(t)
		writeGateTestFile(t, fixture.root, "ROADMAP.md", "capture-only edit\n", 0o644)
		before := len(phaseRunNames(t, fixture.root))
		if got := executeWithEngineAfterAcquireAtKit(context.Background(), fixture.root, fixture.root, io.Discard, io.Discard, productionGateEngine{}, notifyGateSignals, forceRun); got.ActionExit != 0 {
			t.Fatalf("forced run exit = %d, want 0", got.ActionExit)
		}
		executed := append([]string(nil), phaseRunNames(t, fixture.root)[before:]...)
		want := fixture.phaseNames()
		sort.Strings(executed)
		sort.Strings(want)
		if !reflect.DeepEqual(executed, want) {
			t.Fatalf("forced run executed %v, want the whole resolved table %v", executed, want)
		}
		if rec := slotRecord(t, fixture.root, time.Now().UTC()); rec.partition() != nil {
			t.Fatalf("forced record = %+v, want a whole-tree verdict", rec)
		}
	})

	t.Run("the forced green re-authors every slot", func(t *testing.T) {
		fixture := seededScopingFixture(t)
		before := slotBytes(t, fixture.root)
		writeGateTestFile(t, fixture.root, "ROADMAP.md", "capture-only edit\n", 0o644)
		later := time.Now().UTC().Add(2 * time.Hour).Truncate(time.Second)
		if got := executeWithEngineAfterAcquire(context.Background(), fixture.root, io.Discard, io.Discard,
			&faultEngine{now: later}, nil, forceRun); got.ActionExit != 0 {
			t.Fatalf("forced execution = %+v, want green", got)
		}
		for component, data := range slotBytes(t, fixture.root) {
			if data == nil {
				t.Fatalf("%s holds no slot after a forced green", component)
			}
			if bytes.Equal(data, before[component]) {
				t.Fatalf("%s's slot was not re-authored by the forced green: %q", component, data)
			}
		}
	})
}

// [PS46] A partial green retains nothing a reuse can answer from. The verdict cache refuses a
// partial record by its class, but the retained store is asked only whether a green is held
// for this tree and oracle — so a whole-tree record written here would let `bench commit`
// credit skipped components for the length of the freshness window, and a stripped record
// would hand a later reduced run an ancestor that graded none of what it inherits.
func TestPartialGreenRetainsNoWholeTreeEvidence(t *testing.T) {
	t.Parallel()
	fixture := seededScopingFixture(t)
	gitdir := commonGitDirOf(t, fixture.root)
	retained := evidenceFiles(t, gitdir)
	sort.Strings(retained)

	writeGateTestFile(t, fixture.root, "ROADMAP.md", "capture-only edit\n", 0o644)
	observeGreenGate(t, fixture.root)
	partialRecord(t, fixture.root)

	after := evidenceFiles(t, gitdir)
	sort.Strings(after)
	if !reflect.DeepEqual(after, retained) {
		t.Fatalf("retained evidence after the partial green = %v, want the seed's %v unchanged", after, retained)
	}
	plan := mustSubject(t, fixture.root)
	if reuse := reusableEvidence(fixture.root, plan, time.Now().UTC()); reuse.ReusableGreen {
		t.Fatalf("partial green is reusable as a whole-tree green: %+v", reuse)
	}
	// The same answer through the surface a gated commit actually asks: it must pay a run
	// rather than be handed the partial green.
	var stdout bytes.Buffer
	before := len(phaseRunNames(t, fixture.root))
	if got := executeReusingFreshGreenAtKit(context.Background(), fixture.root, fixture.root, &stdout, io.Discard); got.ActionExit != 0 {
		t.Fatalf("commit-path execution = %+v, want green", got)
	}
	if len(phaseRunNames(t, fixture.root)) == before {
		t.Fatalf("the commit path executed nothing and said:\n%s\nwant the partial green refused as a reuse", stdout.String())
	}
}

func TestComposedGreenAcceptsOnlyCompleteExactTipEvidence(t *testing.T) {
	t.Run("no verdict", func(t *testing.T) {
		root := routedKitFixture(t)
		if composedGreenAtKit(root, root) {
			t.Fatal("absent verdict composed to whole-tree green")
		}
	})

	t.Run("invalid verdict", func(t *testing.T) {
		root := routedKitFixture(t)
		writeCache(t, cachePath(t, root), "{", 0o600)
		if composedGreenAtKit(root, root) {
			t.Fatal("invalid verdict composed to whole-tree green")
		}
	})

	t.Run("full", func(t *testing.T) {
		root := routedKitFixture(t)
		mustExecuteGreen(t, root, productionGateEngine{})
		t.Setenv("BENCH_KIT", root)
		if !ComposedGreen(root) {
			t.Fatal("full green did not compose to whole-tree green")
		}
	})

	t.Run("capture-only edit runs full", func(t *testing.T) {
		// The whole-changeset reduced path is retired: a capture-only edit on a root
		// component scoping cannot reach (no go.mod) pays the full run, and that full
		// green composes.
		root := routedKitFixture(t)
		mustExecuteGreen(t, root, productionGateEngine{})
		writeGateTestFile(t, root, "ROADMAP.md", "capture-only edit\n", 0o644)
		mustExecuteGreen(t, root, productionGateEngine{})
		if got := fullRunCount(t, root); got != 2 {
			t.Fatalf("gate runs = %d, want 2 full runs — the capture-only edit must not narrow", got)
		}
		if !composedGreenAtKit(root, root) {
			t.Fatal("full green after a capture-only edit did not compose")
		}
	})

	t.Run("partial", func(t *testing.T) {
		fixture := seededScopingFixture(t)
		writeGateTestFile(t, fixture.root, "ROADMAP.md", "capture-only edit\n", 0o644)
		observeGreenGate(t, fixture.root)
		if !composedGreenAtKit(fixture.root, fixture.root) {
			t.Fatal("partial green with every retained component did not compose")
		}
	})

	t.Run("missing skipped component evidence", func(t *testing.T) {
		fixture := seededScopingFixture(t)
		writeGateTestFile(t, fixture.root, "ROADMAP.md", "capture-only edit\n", 0o644)
		observeGreenGate(t, fixture.root)
		record := partialRecord(t, fixture.root)
		var skipped string
		for _, component := range record.Skipped {
			if componentSkipsOnAncestorEvidence(component) {
				skipped = component
				break
			}
		}
		if skipped == "" {
			t.Fatal("partial fixture has no slot-backed skipped component")
		}
		if err := os.Remove(scopedSlotPath(t, fixture.root, skipped)); err != nil {
			t.Fatal(err)
		}
		if composedGreenAtKit(fixture.root, fixture.root) {
			t.Fatal("partial green composed after a skipped component lost its evidence")
		}
	})

	t.Run("red", func(t *testing.T) {
		root := routedKitFixture(t)
		mustExecuteGreen(t, root, productionGateEngine{})
		writeGateTestFile(t, root, ".bench/gate.sh", "#!/usr/bin/env bash\nexit 1\n", 0o755)
		if got := executeWithEngineAtKit(context.Background(), root, root, io.Discard, io.Discard, productionGateEngine{}); got.ActionExit == 0 {
			t.Fatal("red fixture gate unexpectedly passed")
		}
		if composedGreenAtKit(root, root) {
			t.Fatal("red verdict composed to whole-tree green")
		}
	})
}

// scopedSlotPath is where component's slot at its current identity lives.
func scopedSlotPath(t *testing.T, root, component string) string {
	t.Helper()
	dir, err := componentSlotDir(root)
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(dir, componentSlotName(component, mustResolveComponentIdentities(t, root)[component]))
}

// plantForeignSlot files an otherwise valid slot record for owner at the address host's slot
// is looked up under — the domain mismatch a per-component policy domain exists to refuse.
func plantForeignSlot(t *testing.T, root, host, owner string) {
	t.Helper()
	identities := mustResolveComponentIdentities(t, root)
	record := componentSlotRecord{
		Schema:     componentSlotSchema,
		Component:  owner,
		Identity:   identities[host],
		AuthoredAt: time.Now().UTC().Truncate(time.Second).Format(time.RFC3339),
	}
	data, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	dir, err := componentSlotDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := durableReplaceRecordAt(dir, componentSlotName(host, identities[host]), data); err != nil {
		t.Fatal(err)
	}
}
