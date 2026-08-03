package runtime

// The retired reduced run, observed through the built binary: the whole-changeset
// reduced path no longer exists, so a capture-only edit on a root component scoping
// cannot reach pays a full run — never a narrowed one it would announce. The durable
// markers are the contract: the resolved gate must run again, and the phase table must
// not run directly.

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

// [R24] A capture-only edit on the kit-root shape pays a second full run: the retired
// whole-changeset reduction must not resurface as an announcement, a skipped phase, or
// an inherited-evidence claim — the resolved gate runs again and the phase table never
// runs directly.
func TestKitRootCaptureEditPaysFullRun(t *testing.T) {
	t.Parallel()
	contract.RequireFreshBench(t)
	contract.NoteContractFailure(t, "kit-root full-run contract failed")
	root, f, env, bench := seededReducedGateFixture(t, "")

	writeReducedFixtureFile(t, root, "decisions/probe.md", "decision-only edit\n", 0o644)
	probe := f.RunEnvSpec(env, bench, "gate-run", root)
	probe.RequireExit(0)

	output := probe.Stdout + probe.Stderr
	for _, retired := range []string{
		"gate: reduced run",
		"skipping " + canary.PhaseTest,
		"evidence inherited from full green",
	} {
		if strings.Contains(output, retired) {
			t.Fatalf("kit-root run resurfaced the retired reduction marker %q:\nstdout:\n%s\nstderr:\n%s", retired, probe.Stdout, probe.Stderr)
		}
	}
	if got := reducedFixtureMarker(t, root, "phase-runs"); got != "" {
		t.Fatalf("phase marker after the capture-only edit = %q, want none — the phase table ran instead of the resolved gate", got)
	}
	if got := reducedFixtureMarker(t, root, "full-runs"); got != "full\nfull\n" {
		t.Fatalf("gate marker after the capture-only edit = %q, want a second full run", got)
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

// seededReducedGateFixture builds the retired reduced path's repository shape — one
// included and one excludable phase, both observable through durable markers — and seeds
// a full green. kit chooses whose declaration the graded root claims: empty means the
// fixture is its own kit (the shape that once reduced), while a non-empty path plays a
// linked repo's wrapper naming a kit checkout elsewhere.
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
