package runtime

// The per-component evidence store across a process boundary. Ancestor slots and the build
// attestation are serialized by one gate process and re-read by the next, so a defect that
// only appears on reload — a record no reader parses, an address the author and the reader
// spell differently, a partition answered from memory rather than from the store — is
// invisible to any test that stays inside one process. Every assertion here is therefore on
// what a second, fresh `bench gate-run` printed, which components it executed, and what the
// store held afterwards.
//
// The fixture is the kit-shaped root: a real Go module with a ./cmd/bench main and a sealed
// dist/bench, so the components the input registry declares resolve identities and the build
// phase has an artifact to be attested. Its phases are marker scripts — the executed set is
// read from a durable marker rather than from a return value — and the resolved gate script
// hands off to gate-phases, which is what lets the root be narrowed at all.

import (
	"encoding/json"
	"io/fs"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/contract"
	benchfreshness "github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/gate"
)

// boundaryCapturePath is the edit every test here makes between two gate runs. No component's
// declared inputs cover it, so it moves the whole tree — retiring the whole-tree green the
// first run recorded — while leaving every component's identity where the first run authored
// it. That is exactly the changeset per-component scoping exists to narrow.
const boundaryCapturePath = "ROADMAP.md"

// boundaryGateScript is the fixture's resolved gate. The closing exec is the gate-phases
// hand-off the narrowing decision reads before it will narrow this root at all; the loop
// above it stands in for the phase table on the unnarrowed path, discovering the phase
// scripts rather than listing them so the manifest stays the only place the table is written.
const boundaryGateScript = `#!/usr/bin/env bash
set -uo pipefail
echo full >> .git/full-runs
for script in .bench/phase-*.sh; do
  bash "$script" || exit 1
done
exec true gate-phases "$PWD"
`

// [PC21a] Slots and an attestation authored by one CLI process are honored by a fresh one:
// the second process skips every component some evidence covers, announces each with the
// evidence that covered it, executes only the components no evidence can cover, and leaves
// the store it read exactly as it found it. The last clause is not decoration — a reader that
// re-stamped what it consumed would let evidence age forever while every announcement stayed
// current, which is the same wrong skip a re-timestamped whole-tree verdict would buy.
func TestPartialEvidenceSurvivesAProcessBoundary(t *testing.T) {
	t.Parallel()
	contract.NoteContractFailure(t, "per-component evidence did not survive a process boundary")
	fixture := seededBoundaryFixture(t)
	authored := fixture.readEvidence(t)
	seeded := fixture.markers(t, "phase-runs")

	fixture.captureOnlyEdit(t)
	probe := fixture.gateRun(t)
	probe.RequireExit(0)

	output := probe.Stdout + probe.Stderr
	announced := announcedSkips(t, output)
	skippable, unconditional := boundaryComponents(fixture.phases)
	if got := slices.Sorted(maps.Keys(announced)); !slices.Equal(got, skippable) {
		t.Fatalf("second process announced skips for %v, want every evidence-covered component %v:\n%s", got, skippable, output)
	}
	if got := fixture.markers(t, "phase-runs")[len(seeded):]; !slices.Equal(slices.Sorted(slices.Values(got)), unconditional) {
		t.Fatalf("second process executed %v, want the components no evidence covers %v:\n%s", got, unconditional, output)
	}
	if got := len(fixture.markers(t, "full-runs")); got != 1 {
		t.Fatalf("resolved gate runs = %d, want the seed run only — the narrowed run paid the whole gate", got)
	}

	// Each announced skip names the evidence class its component actually rests on: build
	// reuses the artifact its seal describes, and every other component an ancestor slot the
	// first process wrote.
	for component, skip := range announced {
		if component == canary.PhaseBuild {
			if skip.seal == "" || skip.seal != boundarySealDigest(t, fixture.root) {
				t.Fatalf("build skip announced %+v, want the seal source digest %s", skip, boundarySealDigest(t, fixture.root))
			}
			continue
		}
		if got := authored.slots[component].record.Identity; skip.identity != got {
			t.Fatalf("%s skip announced identity %q, want the slot the first process authored at %q", component, skip.identity, got)
		}
	}

	reread := fixture.readEvidence(t)
	for _, name := range slices.Sorted(maps.Keys(authored.entries)) {
		before, after := authored.entries[name], reread.entries[name]
		if after.data == nil {
			t.Fatalf("store entry %s is gone after the second process; it skipped on that record", name)
		}
		if !os.SameFile(before.info, after.info) || string(before.data) != string(after.data) {
			t.Fatalf("store entry %s was re-authored by the process that only read it:\nbefore %s\nafter  %s",
				name, before.data, after.data)
		}
	}
}

// [PC21b] A forged slot and a forged attestation planted between the two runs are refused,
// and the components they claimed to cover run. The forgeries are the two shapes the store
// cannot rule out by construction: a record filed at one component's address that answers for
// another, and an attestation naming a binary other than the one the seal describes.
func TestForgedEvidenceIsRefusedOnReload(t *testing.T) {
	t.Parallel()
	contract.NoteContractFailure(t, "forged per-component evidence was honored on reload")
	fixture := seededBoundaryFixture(t)
	authored := fixture.readEvidence(t)
	seeded := fixture.markers(t, "phase-runs")

	// The slot is overwritten with a valid record for a different component. It stays a
	// well-formed slot, so nothing but the reader's own comparison of the record against the
	// address it was found at can refuse it.
	forged := authored.slots[canary.PhaseVet]
	victim := authored.slots[canary.PhaseGofmt].record.Component
	writeBoundaryRecord(t, forged.path, storeRecord{
		Schema:     forged.record.Schema,
		Component:  victim,
		Identity:   forged.record.Identity,
		AuthoredAt: forged.record.AuthoredAt,
	})
	attestation := authored.attestation
	writeBoundaryRecord(t, attestation.path, storeRecord{
		Schema:     attestation.record.Schema,
		Executable: strings.Repeat("b", len(attestation.record.Executable)),
		AuthoredAt: attestation.record.AuthoredAt,
	})

	fixture.captureOnlyEdit(t)
	probe := fixture.gateRun(t)
	probe.RequireExit(0)

	output := probe.Stdout + probe.Stderr
	announced := announcedSkips(t, output)
	skippable, unconditional := boundaryComponents(fixture.phases)
	refused := []string{canary.PhaseBuild, canary.PhaseVet}
	if got := slices.Sorted(maps.Keys(announced)); !slices.Equal(got, without(skippable, refused)) {
		t.Fatalf("second process announced skips for %v, want the forged components %v refused:\n%s", got, refused, output)
	}
	want := slices.Sorted(slices.Values(slices.Concat(unconditional, refused)))
	if got := fixture.markers(t, "phase-runs")[len(seeded):]; !slices.Equal(slices.Sorted(slices.Values(got)), want) {
		t.Fatalf("second process executed %v, want the unconditional components and the forged ones %v:\n%s", got, want, output)
	}
}

// [PS34] The second process's announcements name what the first process recorded rather than
// anything it could recompute for itself.
//
// The recorded time is what makes that falsifiable. A component's identity is a content
// address, and a component whose identity moved cannot skip at all, so the announced identity
// and the identity this run computes are equal by construction — an announcement built from
// the run's own state is only distinguishable in the other field. Back-dating the authored
// slots is what opens that gap: the records stay valid (no clock retires a slot), and an
// announcement sourced from anywhere but the record then names the wrong instant.
func TestAnnouncedAncestorsMatchWhatWasAuthored(t *testing.T) {
	t.Parallel()
	contract.NoteContractFailure(t, "announced ancestors did not match the authored records")
	fixture := seededBoundaryFixture(t)
	backdated := time.Now().UTC().Add(-30 * time.Hour).Truncate(time.Second)
	authored := fixture.backdateSlots(t, backdated)

	fixture.captureOnlyEdit(t)
	probe := fixture.gateRun(t)
	probe.RequireExit(0)

	output := probe.Stdout + probe.Stderr
	announced := announcedSkips(t, output)
	if len(announced) == 0 {
		t.Fatalf("second process announced no skips at all:\n%s", output)
	}
	for component, skip := range announced {
		if component == canary.PhaseBuild {
			continue
		}
		entry, held := authored.slots[component]
		if !held {
			t.Fatalf("%s skipped on evidence the store holds no slot for:\n%s", component, output)
		}
		if skip.identity != entry.record.Identity || skip.recordedAt != entry.record.AuthoredAt {
			t.Fatalf("%s announced ancestor %s recorded %s, want the authored slot %s recorded %s",
				component, skip.identity, skip.recordedAt, entry.record.Identity, entry.record.AuthoredAt)
		}
	}
}

// [PS35] A `--fresh` third run buys the whole tree back: it announces no skip, pays the
// resolved gate, executes every phase in the table, and re-authors every slot and the
// attestation it could have inherited from.
//
// Re-authorship is observed as the store entry being replaced rather than as changed bytes.
// A forced run over an unmoved tree legitimately records the same identity for the same
// component within the same second, so equal bytes prove nothing either way; the store
// installs every record by rename, so a record that was written again is a different file.
func TestFreshReauthorsAcrossProcesses(t *testing.T) {
	t.Parallel()
	contract.NoteContractFailure(t, "--fresh did not re-author per-component evidence")
	fixture := seededBoundaryFixture(t)

	fixture.captureOnlyEdit(t)
	fixture.gateRun(t).RequireExit(0)
	inherited := fixture.readEvidence(t)
	executed := fixture.markers(t, "phase-runs")

	probe := fixture.gateRun(t, "--fresh")
	probe.RequireExit(0)

	output := probe.Stdout + probe.Stderr
	if got := announcedSkips(t, output); len(got) != 0 {
		t.Fatalf("--fresh announced skips for %v, want a whole-tree run:\n%s", slices.Sorted(maps.Keys(got)), output)
	}
	if got := fixture.markers(t, "phase-runs")[len(executed):]; !slices.Equal(slices.Sorted(slices.Values(got)), slices.Sorted(slices.Values(boundaryPhaseNames(fixture.phases)))) {
		t.Fatalf("--fresh executed %v, want the whole resolved table %v:\n%s", got, boundaryPhaseNames(fixture.phases), output)
	}
	if got := len(fixture.markers(t, "full-runs")); got != 2 {
		t.Fatalf("resolved gate runs = %d, want 2 — --fresh did not pay the whole gate", got)
	}

	reauthored := fixture.readEvidence(t)
	skippable, _ := boundaryComponents(fixture.phases)
	for _, component := range skippable {
		if component == canary.PhaseBuild {
			continue
		}
		before, after := inherited.slots[component], reauthored.slots[component]
		if after.path == "" || os.SameFile(before.info, after.info) {
			t.Fatalf("%s slot was not re-authored by --fresh; it still carries the record the forced run inherited from", component)
		}
	}
	if before, after := inherited.attestation, reauthored.attestation; after.path == "" || os.SameFile(before.info, after.info) {
		t.Fatalf("the build attestation was not re-authored by --fresh; the forced run reused the record it was meant to replace")
	}
}

// boundaryFixture is a seeded kit-shaped root together with the phase table it resolves and
// the CLI invocation that grades it. Tests take their expected component sets from phases; a
// second literal list would disagree with the table the runs are actually made of.
type boundaryFixture struct {
	root   string
	f      contract.Fixture
	env    contract.Env
	bench  string
	phases []gate.Phase
}

// seededBoundaryFixture builds the root and pays the first CLI process: a whole-tree green
// that authors an ancestor slot for every component it graded and an attestation for the
// artifact beside it. Everything each test asserts is about what a later process makes of
// what this run published.
func seededBoundaryFixture(t *testing.T) boundaryFixture {
	t.Helper()
	contract.RequireFreshBench(t)
	if _, err := exec.LookPath("go"); err != nil {
		capability.Environment(t, "go toolchain absent; the per-component fixture needs a real module")
	}
	root := t.TempDir()
	writeBoundaryTree(t, root)
	fixtureGit(t, root, "init", "-q")
	// The manifest is generated from the production table's own answer for this tree, so the
	// executed table carries the phases the kit table materializes here and nothing else: a
	// tree that stops satisfying a phase's shape stops declaring that phase.
	phases := gate.BenchkitPhases(root, root)
	writeBoundaryManifest(t, root, phases)
	sealBoundaryBinary(t, root)

	kit := root
	fixture := boundaryFixture{
		root: root,
		f:    contract.NewExecFixtureAt(t, root),
		env: contract.Env{
			"BENCH_KIT":                  &kit,
			"BENCH_GATE":                 nil,
			"BENCH_CANARY_INNER":         nil,
			"BENCH_REQUIRE_CAPABILITIES": nil,
			capability.LogEnv:            nil,
		},
		bench:  filepath.Join(contract.SubjectRoot(t), "dist", "bench"),
		phases: phases,
	}
	seed := fixture.gateRun(t)
	seed.RequireExit(0)
	if got := len(fixture.markers(t, "full-runs")); got != 1 {
		t.Fatalf("seed run paid the resolved gate %d times, want 1:\n%s%s", got, seed.Stdout, seed.Stderr)
	}
	if got := slices.Sorted(slices.Values(fixture.markers(t, "phase-runs"))); !slices.Equal(got, slices.Sorted(slices.Values(boundaryPhaseNames(phases)))) {
		t.Fatalf("seed run executed %v, want the whole resolved table %v", got, boundaryPhaseNames(phases))
	}
	skippable, _ := boundaryComponents(phases)
	if !slices.Equal(slices.Sorted(maps.Keys(fixture.readEvidence(t).slots)), without(skippable, []string{canary.PhaseBuild})) {
		t.Fatalf("seed run authored slots for %v, want one per evidence-covered component in %v",
			slices.Sorted(maps.Keys(fixture.readEvidence(t).slots)), skippable)
	}
	return fixture
}

// gateRun execs the built CLI over the fixture in a process of its own. Every observation in
// this file is of what one such process printed, executed, or left in the store.
func (b boundaryFixture) gateRun(t *testing.T, flags ...string) contract.Probe {
	t.Helper()
	args := append([]string{"gate-run"}, flags...)
	return b.f.RunEnvSpec(b.env, b.bench, append(args, b.root)...)
}

// captureOnlyEdit moves the tree without moving any component's declared inputs.
func (b boundaryFixture) captureOnlyEdit(t *testing.T) {
	t.Helper()
	writeReducedFixtureFile(t, b.root, boundaryCapturePath, "capture-only edit\n", 0o644)
}

// markers are the lines a durable fixture marker has collected. The gate script appends to
// full-runs and each phase script appends its own name to phase-runs, so a run's executed set
// outlives the process that produced it.
func (b boundaryFixture) markers(t *testing.T, name string) []string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(b.root, ".git", name))
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return nil
	}
	return lines
}

// backdateSlots rewrites every authored slot to record an instant well in the past, and
// returns the store as it then stands. The records stay valid — a slot is retired by its
// component's identity moving and never by a clock — so what changes is only whether an
// announcement built from the record and one built from the running process agree.
func (b boundaryFixture) backdateSlots(t *testing.T, at time.Time) boundaryEvidence {
	t.Helper()
	stored := b.readEvidence(t)
	for _, component := range slices.Sorted(maps.Keys(stored.slots)) {
		entry := stored.slots[component]
		entry.record.AuthoredAt = at.Format(time.RFC3339)
		writeBoundaryRecord(t, entry.path, entry.record)
	}
	return b.readEvidence(t)
}

// storeRecord is as much of a retained-evidence record as this file reads: a slot names the
// component and identity it answers for, an attestation names the digest of the binary a gate
// build produced, and the whole-tree verdict classes carry neither name. Decoding every store
// entry into this one shape is what lets the classes be told apart without the test knowing
// how any of their addresses are computed.
type storeRecord struct {
	Schema     int    `json:"schema"`
	Component  string `json:"component,omitempty"`
	Identity   string `json:"identity,omitempty"`
	Executable string `json:"executable,omitempty"`
	AuthoredAt string `json:"authored_at"`
}

// boundaryEntry is one store entry: where it sits, the bytes it holds, what those bytes say,
// and the file identity a later stat is compared against to tell a record that was written
// again from one that was only read.
type boundaryEntry struct {
	path   string
	data   []byte
	info   os.FileInfo
	record storeRecord
}

// boundaryEvidence is the retained-evidence store as one process left it, indexed both by
// store name — for the entries a test compares wholesale — and by the claim each record makes.
type boundaryEvidence struct {
	entries     map[string]boundaryEntry
	slots       map[string]boundaryEntry
	attestation boundaryEntry
}

func (b boundaryFixture) readEvidence(t *testing.T) boundaryEvidence {
	t.Helper()
	dir := filepath.Join(b.root, ".git", "bench-gate-evidence")
	names, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read the retained-evidence store: %v", err)
	}
	stored := boundaryEvidence{entries: map[string]boundaryEntry{}, slots: map[string]boundaryEntry{}}
	for _, name := range names {
		path := filepath.Join(dir, name.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read the store entry %s: %v", name.Name(), err)
		}
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("stat the store entry %s: %v", name.Name(), err)
		}
		entry := boundaryEntry{path: path, data: data, info: info}
		if err := json.Unmarshal(data, &entry.record); err != nil {
			t.Fatalf("decode the store entry %s: %v", name.Name(), err)
		}
		switch {
		case entry.record.Component != "":
			stored.slots[entry.record.Component] = entry
		case entry.record.Executable != "":
			stored.attestation = entry
		default:
			// A whole-tree verdict shares the store with the per-component records and is
			// neither of the claims this file reads, so it is indexed by name alone.
		}
		stored.entries[name.Name()] = entry
	}
	return stored
}

func writeBoundaryRecord(t *testing.T, path string, record storeRecord) {
	t.Helper()
	data, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}

// announcedSkip is one skip line as an operator reads it. The two evidence forms are
// different claims — an ancestor slot names a run that graded the component, build's names the
// artifact a gate build produced — so they are parsed apart rather than into one field.
type announcedSkip struct {
	identity   string
	recordedAt string
	seal       string
}

var (
	ancestorSkipLine = regexp.MustCompile(`^gate: skipping (\S+) \(evidence inherited from ancestor ([0-9a-f]{64}) recorded (\S+)\)$`)
	attestedSkipLine = regexp.MustCompile(`^gate: skipping (\S+) \(artifact reused from the gate-attested build sealed at ([0-9a-f]{64})\)$`)
)

// announcedSkips parses a gate process's per-component announcements. A line that says a
// component was skipped and matches neither form is a failure rather than a line to ignore:
// the announcement is the whole operator-facing contract of a narrowed run, and a test that
// skipped over an unrecognized one would report a silent reduction as a clean pass.
func announcedSkips(t *testing.T, output string) map[string]announcedSkip {
	t.Helper()
	skips := map[string]announcedSkip{}
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "gate: skipping ") {
			continue
		}
		if match := ancestorSkipLine.FindStringSubmatch(line); match != nil {
			skips[match[1]] = announcedSkip{identity: match[2], recordedAt: match[3]}
			continue
		}
		if match := attestedSkipLine.FindStringSubmatch(line); match != nil {
			skips[match[1]] = announcedSkip{seal: match[2]}
			continue
		}
		t.Fatalf("unparsable skip announcement %q; the operator-facing form moved", line)
	}
	return skips
}

// boundaryComponents splits a resolved table into the components some evidence can cover and
// the phases every run executes. The split reads the production input registry: a component it
// declares is a component that resolves an identity and so can skip on its own evidence, and a
// phase it does not declare has no identity to be covered by anything. A component the registry
// gains therefore lands on the skippable side here without an edit.
func boundaryComponents(table []gate.Phase) (skippable, unconditional []string) {
	declared := map[string]bool{}
	for _, source := range gate.ComponentInputSources() {
		declared[source.Component] = true
	}
	for _, phase := range table {
		if declared[phase.Name] {
			skippable = append(skippable, phase.Name)
			continue
		}
		unconditional = append(unconditional, phase.Name)
	}
	return slices.Sorted(slices.Values(skippable)), slices.Sorted(slices.Values(unconditional))
}

func boundaryPhaseNames(table []gate.Phase) []string {
	names := make([]string, 0, len(table))
	for _, phase := range table {
		names = append(names, phase.Name)
	}
	return names
}

func boundarySealDigest(t *testing.T, root string) string {
	t.Helper()
	sources, _, err := benchfreshness.SealDigests(filepath.Join(root, "dist", "bench"))
	if err != nil {
		t.Fatalf("read the fixture seal: %v", err)
	}
	return sources
}

// writeBoundaryTree lays down the tree shape the production phase table reads: a module with a
// ./cmd/bench main, the build helper beside its auxiliary input manifest, the wrapper script
// the canary phase execs, and the canary sources and fixtures. dist/ stays out of the tree
// identity as it does in the kit, so republishing the binary cannot move the subject being
// graded. No script here is ever executed — the manifest routes every phase through a marker
// script — so each carries an inert body.
func writeBoundaryTree(t *testing.T, root string) {
	t.Helper()
	writeReducedFixtureFile(t, root, ".gitignore", "dist/\n", 0o644)
	writeReducedFixtureFile(t, root, "go.mod", "module benchboundaryfixture\n\n"+subjectGoDirective(t)+"\n", 0o644)
	writeReducedFixtureFile(t, root, "cmd/bench/main.go", "package main\n\nfunc main() {}\n", 0o644)
	writeReducedFixtureFile(t, root, "internal/canary/canary.go",
		"package canary\n\n// Name is the surface the canary phase grades.\nfunc Name() string { return \"canary\" }\n", 0o644)
	writeReducedFixtureFile(t, root, "tests/canary/fixture.txt", "canary fixture\n", 0o644)
	writeReducedFixtureFile(t, root, "scripts/go-build.sh", "#!/usr/bin/env bash\nexit 0\n", 0o755)
	writeReducedFixtureFile(t, root, "scripts/go-build.inputs", "build_script=scripts/go-build.sh\n", 0o644)
	writeReducedFixtureFile(t, root, "bin/bench.sh", "#!/usr/bin/env bash\nexit 0\n", 0o755)
	writeReducedFixtureFile(t, root, ".bench/gate-inputs.json",
		`{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`, 0o644)
	writeReducedFixtureFile(t, root, ".bench/gate.sh", boundaryGateScript, 0o755)
	writeReducedFixtureFile(t, root, boundaryCapturePath, "roadmap\n", 0o644)
}

// writeBoundaryManifest declares one marker phase per phase in declared, and writes the script
// each of them runs. The declared edges are carried through: they order the writer of
// dist/bench ahead of its readers, and a fixture that dropped them would leave the build
// phase's skip untested against the edges it is claimed to satisfy.
func writeBoundaryManifest(t *testing.T, root string, declared []gate.Phase) {
	t.Helper()
	type manifestPhase struct {
		Name  string   `json:"name"`
		Argv  []string `json:"argv"`
		Needs []string `json:"needs,omitempty"`
	}
	var doc struct {
		Phases []manifestPhase `json:"phases"`
	}
	for _, phase := range declared {
		script := ".bench/phase-" + phase.Name + ".sh"
		writeReducedFixtureFile(t, root, script, "echo "+phase.Name+" >> .git/phase-runs\n", 0o644)
		doc.Phases = append(doc.Phases, manifestPhase{Name: phase.Name, Argv: []string{"bash", script}, Needs: phase.Needs})
	}
	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	writeReducedFixtureFile(t, root, canary.PhaseManifestPath, string(data)+"\n", 0o644)
}

// sealBoundaryBinary publishes dist/bench through the seal package's own Publish, the only
// writer that produces a seal answering for the tree the binary was built from — moving the
// bytes into place by hand leaves an artifact no reader can verify, and the build phase would
// then run for a reason no test here is about.
//
// It publishes twice because the two halves are circular: the declarations name surfaces the
// tree has to carry, and resolving them reads the seal of the binary published from that same
// tree. The first publish is what lets the declarations resolve at all; the second answers for
// the surfaces materializing them added.
func sealBoundaryBinary(t *testing.T, root string) {
	t.Helper()
	staged := filepath.Join(root, "dist", "bench.staged")
	build := exec.Command("go", "build", "-buildvcs=false", "-o", staged, "./cmd/bench")
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build the fixture binary: %v\n%s", err, out)
	}
	spare := filepath.Join(root, "dist", "bench.spare")
	copyRuntimeFile(t, staged, spare, 0o755)
	executable := filepath.Join(root, "dist", "bench")
	if err := benchfreshness.Publish(root, staged, executable); err != nil {
		t.Fatalf("publish the fixture binary: %v", err)
	}
	materializeDeclaredInputs(t, root)
	if err := benchfreshness.Publish(root, spare, executable); err != nil {
		t.Fatalf("republish the fixture binary: %v", err)
	}
}

// materializeDeclaredInputs writes whatever the production input registry declares and the
// tree does not already carry. It is derived from the registry rather than listed again here,
// so a declaration that gains a surface lands on this root without an edit — and a component
// whose declared input is missing resolves no identity at all, which would leave every skip
// assertion in this file passing over a component nothing ever graded.
func materializeDeclaredInputs(t *testing.T, root string) {
	t.Helper()
	sets, err := gate.ResolveComponentInputs(root)
	if err != nil {
		t.Fatalf("resolve the declared component inputs: %v", err)
	}
	for _, component := range slices.Sorted(maps.Keys(sets)) {
		for _, path := range sets[component].Paths() {
			full := filepath.Join(root, filepath.FromSlash(path))
			if strings.HasSuffix(path, "/") {
				// A directory entry is satisfied by any file beneath it.
				if !holdsBoundaryFile(full) {
					writeReducedFixtureFile(t, root, path+"declared.txt", "declared surface\n", 0o644)
				}
				continue
			}
			if info, err := os.Stat(full); err == nil && info.Mode().IsRegular() {
				continue
			}
			// The declared file entries a derived closure cannot have named are the wrapper
			// scripts a phase's wiring execs. Nothing here runs them.
			writeReducedFixtureFile(t, root, path, "#!/usr/bin/env bash\nexit 0\n", 0o755)
		}
	}
}

func holdsBoundaryFile(dir string) bool {
	found := false
	_ = filepath.WalkDir(dir, func(_ string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type().IsRegular() {
			found = true
			return filepath.SkipAll
		}
		return nil
	})
	return found
}

func without(values, excluded []string) []string {
	kept := make([]string, 0, len(values))
	for _, value := range values {
		if !slices.Contains(excluded, value) {
			kept = append(kept, value)
		}
	}
	return slices.Sorted(slices.Values(kept))
}
