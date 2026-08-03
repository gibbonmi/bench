package runtime

// `bench commit` with the reduced path retired: every staged set pays the full gate —
// confined to the old allowlist or mixing an unlisted path, the shape of the staged set
// no longer narrows the run. Both assert on durable markers rather than on elapsed time,
// so the evidence is deterministic.

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/contract"
	"github.com/gibbonmi/bench/internal/gate"
)

// [R22] The reduced path is retired: an allowlist-confined staged set pays the full gate
// through the resolved script, announces no reduction, and the commit still lands.
func TestCommitPaysFullGateForConfinedStagedSet(t *testing.T) {
	t.Parallel()
	contract.RequireFreshBench(t)
	contract.NoteContractFailure(t, "retired-reduction commit contract failed")
	root, f, env, bench := commitReducedFixture(t)

	writeReducedFixtureFile(t, root, "ROADMAP.md", "capture-only edit\n", 0o644)
	probe := f.RunEnvSpec(env, bench, "commit", "-m", "capture edit", "ROADMAP.md")
	probe.RequireExit(0)

	output := probe.Stdout + probe.Stderr
	if strings.Contains(output, "gate: reduced run") {
		t.Fatalf("confined commit announced a reduced run:\nstdout:\n%s\nstderr:\n%s", probe.Stdout, probe.Stderr)
	}
	if got := reducedFixtureMarker(t, root, "full-runs"); got != "full\nfull\n" {
		t.Fatalf("gate marker after the confined commit = %q, want a second full run", got)
	}
	if got := reducedFixtureMarker(t, root, "phase-runs"); got != "" {
		t.Fatalf("phase marker after the confined commit = %q, want none — the phase table ran instead of the resolved gate", got)
	}
	if names := committedNames(f); !strings.Contains(names, "ROADMAP.md") {
		t.Fatalf("confined commit did not land ROADMAP.md; committed:\n%s", names)
	}
}

// [R23] A staged set mixing an allowlisted path with an unlisted one runs the full gate.
// Confinement is every path or none: an "any" rule reduces here and lets the unlisted
// file land ungraded.
func TestCommitMixedStagedSetRunsFullGate(t *testing.T) {
	t.Parallel()
	contract.RequireFreshBench(t)
	contract.NoteContractFailure(t, "mixed staged set contract failed")
	root, f, env, bench := commitReducedFixture(t)

	writeReducedFixtureFile(t, root, "ROADMAP.md", "capture-only edit\n", 0o644)
	writeReducedFixtureFile(t, root, "graded.txt", "graded edit\n", 0o644)
	probe := f.RunEnvSpec(env, bench, "commit", "-m", "mixed edit", "ROADMAP.md", "graded.txt")
	probe.RequireExit(0)

	output := probe.Stdout + probe.Stderr
	if strings.Contains(output, "gate: reduced run") {
		t.Fatalf("mixed staged set reduced; the unlisted path would land ungraded:\nstdout:\n%s\nstderr:\n%s", probe.Stdout, probe.Stderr)
	}
	if got := reducedFixtureMarker(t, root, "full-runs"); got != "full\nfull\n" {
		t.Fatalf("gate marker after the mixed commit = %q, want a second full run", got)
	}
	if names := committedNames(f); !strings.Contains(names, "graded.txt") {
		t.Fatalf("mixed commit did not land graded.txt; committed:\n%s", names)
	}
}

// commitReducedFixture builds a repository whose declared phase table has one excludable
// and one included phase, seeds a full green, and
// returns everything a `bench commit` probe needs. The repository carries a git identity
// because `bench commit` runs `git commit` itself, with no inline identity to lend it.
func commitReducedFixture(t *testing.T) (string, contract.Fixture, contract.Env, string) {
	t.Helper()
	scope := gate.ReducedScope()
	if !scope.Member("ROADMAP.md") || scope.Member("graded.txt") {
		t.Fatal("fixture paths no longer straddle the declaration; repoint the fixture")
	}
	if !scope.Excludable(canary.PhaseTest) || scope.Excludable("conformance") {
		t.Fatal("fixture phase names no longer match the declaration; repoint the fixture")
	}

	root := t.TempDir()
	// The gate script carries the gate-phases hand-off the kit's own entry uses — a
	// stand-in binary keeps the exec inert, and the marker line is this fixture's
	// observable. The hand-off is what makes the root reducible: routing, not the
	// manifest beside it, is the proof the reduction requires.
	writeReducedFixtureFile(t, root, ".bench/gate.sh",
		"#!/usr/bin/env bash\necho full >> .git/full-runs\nexec true gate-phases \"$PWD\"\n", 0o755)
	writeReducedFixtureFile(t, root, ".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`, 0o644)
	writeReducedFixtureFile(t, root, ".bench/phase-conformance.sh", "echo conformance >> .git/phase-runs\n", 0o644)
	writeReducedFixtureFile(t, root, ".bench/phase-test.sh", "echo test >> .git/phase-runs\n", 0o644)
	writeReducedFixtureFile(t, root, canary.PhaseManifestPath, `{"phases":[`+
		`{"name":"conformance","argv":["bash",".bench/phase-conformance.sh"]},`+
		fmt.Sprintf(`{"name":%q,"argv":["bash",".bench/phase-test.sh"]}]}`, canary.PhaseTest)+"\n", 0o644)
	writeReducedFixtureFile(t, root, "ROADMAP.md", "roadmap\n", 0o644)
	writeReducedFixtureFile(t, root, "graded.txt", "graded content\n", 0o644)
	fixtureGit(t, root, "init", "-q")
	fixtureGit(t, root, "config", "user.email", "fixture@bench.invalid")
	fixtureGit(t, root, "config", "user.name", "bench-fixture")
	fixtureGit(t, root, "add", "-A")
	fixtureGit(t, root, "commit", "-q", "-m", "fixture")

	f := contract.NewExecFixtureAt(t, root)
	kit := root
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
