package runtime

// The reduced run through the built binary: when a change is confined to the declared
// allowlist and a full-green ancestor's stripped identity still answers for the tree,
// `bench gate-run` executes only the included phases — and says so. The announcement is
// the contract here: silent reduction reads as a gate that never ran, matching the
// failure the existing reused-verdict announcement exists to prevent, and the operator
// has no other signal that the verdict just recorded is narrower than a full run's.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/contract"
	"github.com/gibbonmi/bench/internal/gate"
)

// [R24] A reduced run announces the phases it skipped and the full-green run whose
// evidence it inherited, and the announcement travels with an actual reduction: the
// resolved gate script must not run again, the excludable phase must not run at all,
// and the included phase must still be graded against the real root.
func TestReducedRunAnnouncesSkippedPhases(t *testing.T) {
	t.Parallel()
	contract.RequireFreshBench(t)
	contract.NoteContractFailure(t, "reduced-run announcement contract failed")
	root, f, env, bench := seededReducedGateFixture(t, "")

	writeReducedFixtureFile(t, root, "decisions/probe.md", "decision-only edit\n", 0o644)
	probe := f.RunEnvSpec(env, bench, "gate-run", root)
	probe.RequireExit(0)

	output := probe.Stdout + probe.Stderr
	for _, want := range []string{
		"gate: reduced run",
		"skipping " + canary.PhaseTest,
		"evidence inherited from full green",
		"phase conformance: green",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("reduced run output missing %q:\nstdout:\n%s\nstderr:\n%s", want, probe.Stdout, probe.Stderr)
		}
	}
	if strings.Contains(output, "phase "+canary.PhaseTest+":") {
		t.Fatalf("reduced run still reported the excludable phase:\n%s", output)
	}
	if got := reducedFixtureMarker(t, root, "full-runs"); got != "full\n" {
		t.Fatalf("gate marker after the reduced run = %q, want the seed run only — the announced reduction paid the resolved gate", got)
	}
	if got := reducedFixtureMarker(t, root, "phase-runs"); got != "conformance\n" {
		t.Fatalf("phase marker after the reduced run = %q, want the included phase only", got)
	}
}

// [RB1] A root that is not the kit runs unreduced and pays a second full run for a
// capture-only edit. The reduced scope is the kit's own declaration; a linked repo's
// wrapper names the kit checkout through BENCH_KIT while grading its own tree, and that
// repo never declared the allowlist the reduction would inherit against — so the fixture
// here differs from the announcement fixture in exactly one fact: BENCH_KIT names a tree
// other than the graded root.
func TestForeignRootRunsUnreduced(t *testing.T) {
	t.Parallel()
	contract.RequireFreshBench(t)
	contract.NoteContractFailure(t, "foreign-root reduction contract failed")
	root, f, env, bench := seededReducedGateFixture(t, t.TempDir())

	writeReducedFixtureFile(t, root, "ROADMAP.md", "capture-only edit\n", 0o644)
	probe := f.RunEnvSpec(env, bench, "gate-run", root)
	probe.RequireExit(0)

	output := probe.Stdout + probe.Stderr
	if strings.Contains(output, "gate: reduced run") {
		t.Fatalf("foreign root reduced against a scope it never declared:\nstdout:\n%s\nstderr:\n%s", probe.Stdout, probe.Stderr)
	}
	if got := reducedFixtureMarker(t, root, "phase-runs"); got != "" {
		t.Fatalf("phase marker on the foreign root = %q, want none — the phase table ran instead of the resolved gate", got)
	}
	if got := reducedFixtureMarker(t, root, "full-runs"); got != "full\nfull\n" {
		t.Fatalf("gate marker after the capture-only edit = %q, want a second full run", got)
	}
}

// seededReducedGateFixture builds the reduced-path repository — one included and one
// excludable phase, both observable through durable markers — and seeds the full-green
// ancestor a reduction would inherit from. kit chooses whose declaration the graded root
// claims: empty means the fixture is its own kit (the reduction-eligible shape), while a
// non-empty path plays a linked repo's wrapper naming a kit checkout elsewhere.
func seededReducedGateFixture(t *testing.T, kit string) (string, contract.Fixture, contract.Env, string) {
	t.Helper()
	scope := gate.ReducedScope()
	if !scope.Member("ROADMAP.md") {
		t.Fatal("fixture capture path is no longer declared; repoint the fixture")
	}
	if !scope.Excludable(canary.PhaseTest) || scope.Excludable("conformance") {
		t.Fatal("fixture phase names no longer match the declaration; repoint the fixture")
	}

	root := t.TempDir()
	writeReducedFixtureFile(t, root, ".bench/gate.sh", "#!/usr/bin/env bash\necho full >> .git/full-runs\nexec true gate-phases \"$PWD\"\n", 0o755)
	writeReducedFixtureFile(t, root, ".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`, 0o644)
	writeReducedFixtureFile(t, root, ".bench/phase-conformance.sh", "echo conformance >> .git/phase-runs\n", 0o644)
	writeReducedFixtureFile(t, root, ".bench/phase-test.sh", "echo test >> .git/phase-runs\n", 0o644)
	writeReducedFixtureFile(t, root, canary.PhaseManifestPath, `{"phases":[`+
		`{"name":"conformance","argv":["bash",".bench/phase-conformance.sh"]},`+
		fmt.Sprintf(`{"name":%q,"argv":["bash",".bench/phase-test.sh"]}]}`, canary.PhaseTest)+"\n", 0o644)
	writeReducedFixtureFile(t, root, "ROADMAP.md", "roadmap\n", 0o644)
	writeReducedFixtureFile(t, root, "graded.txt", "graded content\n", 0o644)
	fixtureGit(t, root, "init", "-q")
	fixtureGit(t, root, "add", "-A")
	fixtureGit(t, root, "commit", "-q", "-m", "fixture")

	f := contract.NewExecFixtureAt(t, root)
	if kit == "" {
		kit = root
	}
	env := contract.Env{
		"BENCH_KIT":                  &kit,
		"BENCH_GATE":                 nil,
		"BENCH_CANARY_INNER":         nil,
		"BENCH_REQUIRE_CAPABILITIES": nil,
		capability.LogEnv:            nil,
	}
	bench := filepath.Join(contract.SubjectRoot(t), "dist", "bench")

	seed := f.RunEnvSpec(env, bench, "gate-run", root)
	seed.RequireExit(0)
	if got := reducedFixtureMarker(t, root, "full-runs"); got != "full\n" {
		t.Fatalf("seed run gate marker = %q, want one full run", got)
	}
	return root, f, env, bench
}

func writeReducedFixtureFile(t *testing.T, root, rel, body string, mode os.FileMode) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	contract.WriteFileAbs(t, path, body)
	if err := os.Chmod(path, mode); err != nil {
		t.Fatal(err)
	}
}

func reducedFixtureMarker(t *testing.T, root, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, ".git", name))
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	return string(data)
}
