package runtime

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

// Seam-A2 sentinel contracts: the environment the project gate script actually
// receives, observed at the built binary. FT78 already closes this subject —
// internal/gate launches .bench/gate.sh with PATH plus only the names declared
// under `environment` in .bench/gate-inputs.json — so these are GREEN on the
// shipped tree. Their whole value is a regression guard: a future change
// widening the gate subject turns them red instead of green-by-silence. Each
// plants a marker in the parent process and looks for it in the gate script's
// own dump of its exported names, and reuses readEnvDump so an implementation
// that silently failed to launch the gate cannot pass by finding nothing.
//
// Red capability was proven once by a targeted mutation, recorded on
// testGateDropsMarker below.
func TestRuntimeGateEnvSentinelContracts(t *testing.T) {
	t.Parallel()
	contract.RunParallel(t, "gate drops a non-declared marker", testGateDropsMarker)
	contract.RunParallel(t, "gate keeps PATH and a manifest-declared name", testGateKeepsPathAndDeclared)
}

// gateEnvMarkers is the parent environment every seam-A2 gate run plants.
// FT88_GATE_MARKER is undeclared in every fixture's manifest, so its presence in
// the gate script's dump is a leak by construction; FT88_GATE_DECLARED is the one
// name the fixture declares under `environment`, so it must survive.
func gateEnvMarkers() map[string]string {
	return map[string]string{
		"FT88_GATE_MARKER":   "ft88-gate-marker-must-not-reach-the-gate-script",
		"FT88_GATE_DECLARED": "ft88-gate-declared",
	}
}

// gateEnvDumpFixture builds a gate-capable repo whose stub .bench/gate.sh dumps
// its own exported-variable names to an absolute path outside the repo, and
// declares exactly one marker name (FT88_GATE_DECLARED) under the manifest's
// `environment`. The dump path is baked into the gate script rather than passed
// through the environment: a path read from the environment would itself be
// subject to the closed subject under test, and a leak would then look like a
// test bug.
func gateEnvDumpFixture(t *testing.T) (contract.Fixture, string) {
	t.Helper()
	f := contract.NewFixture(t)
	dump := filepath.Join(t.TempDir(), "gate-env")
	f.WriteExecutable(".bench/gate.sh", fmt.Sprintf("#!/usr/bin/env bash\n{ compgen -e || true; } > %q\nexit 0\n", dump))
	f.WriteFile(".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":["FT88_GATE_DECLARED"],"paths":[],"tools":[]}`+"\n")
	f.CommitAll("seam-A2 gate")
	return f, dump
}

// runGateDump runs one `bench gate` with the markers planted and returns the gate
// script's dumped names, after readEnvDump has proved the script actually ran.
func runGateDump(t *testing.T, f contract.Fixture, dump string) []string {
	t.Helper()
	f.BenchEnv(gateEnvMarkers(), "gate").RequireExit(0)
	return readEnvDump(t, dump)
}

// testGateDropsMarker is story 2's marker-absence claim: a name the reviewer
// happens to have exported reaches the gate script only if the subject fails to
// filter by declaration. Green on the shipped tree because FT78 closes the
// subject to PATH-plus-declared.
//
// Demonstrated red: in internal/gate/subject.go, buildSubject constructs the
// gate script's environment as s.Env, seeded with PATH and extended only by
// manifest-declared names. Temporarily appending the parent marker to that
// subject —
//
//	s.Env = append(s.Env, "FT88_GATE_MARKER="+os.Getenv("FT88_GATE_MARKER"))
//
// inserted after s.Env is seeded in buildSubject (and adding the "os" import) —
// widened the subject exactly as a careless future change would. This test then
// failed with:
//
//	subprocess environment carries non-passlisted name "FT88_GATE_MARKER": [...]
//
// The mutation was reverted after the red was observed; the shipped tree is
// green.
func testGateDropsMarker(t *testing.T) {
	f, dump := gateEnvDumpFixture(t)
	names := runGateDump(t, f, dump)
	requireNoNames(t, names, "FT88_GATE_MARKER")
}

// testGateKeepsPathAndDeclared pins the subject as PATH-plus-declared rather than
// empty, which the marker-absence row alone would accept: PATH and the fixture's
// one manifest-declared name both survive into the gate script's environment.
func testGateKeepsPathAndDeclared(t *testing.T) {
	f, dump := gateEnvDumpFixture(t)
	names := runGateDump(t, f, dump)
	requireNames(t, names, "PATH", "FT88_GATE_DECLARED")
	requireNoNames(t, names, "FT88_GATE_MARKER")
}
