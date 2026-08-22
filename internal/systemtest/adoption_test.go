//go:build system

package systemtest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/adopt"
)

// TestAdoptionSmokeJourney adopts one disposable repository with `bench setup --yes` and
// drives its scaffolded gate through the installed wrapper. The kit's own oracle observes
// the adopter's side of adoption, not an audit. Every launch binds one private BENCH_HOME
// under the test's temporary directory. runSelected and runWrapper carry fixed
// environments with no BENCH_HOME override. This journey composes observeSelected plus
// runAt the way TestWorktreeReauthorizeJourney does. It asserts after every leg that the
// private home stayed empty.
func TestAdoptionSmokeJourney(t *testing.T) {
	repo := owner.repos[1]
	home := filepath.Join(t.TempDir(), "bench-home")
	environment := []string{"BENCH_HOME=" + home, "BENCH_RUN_BINARY=" + owner.selected.path, "BENCH_KIT=" + owner.kit}
	launch := func(program string, args ...string) processResult {
		if err := owner.observeSelected(); err != nil {
			t.Fatal(err)
		}
		return owner.runAt(repo, environment, program, args...)
	}

	setup := launch(owner.selected.path, "setup", "--yes")
	sentinelRow := ".bench/gate.sh is still the unconfigured fail-closed stub (replace the " + adopt.SentinelMarker + " sentinel with real checks)"
	if setup.code != 3 || !strings.Contains(setup.stdout+setup.stderr, sentinelRow) {
		t.Fatalf("bench setup --yes = (%d, %q, %q)", setup.code, setup.stdout, setup.stderr)
	}
	for _, rel := range []string{".bench/gate.sh", ".bench/gate-inputs.json", ".bench/bin/bench.sh", ".bench/dist/bench"} {
		if _, err := os.Stat(filepath.Join(repo, filepath.FromSlash(rel))); err != nil {
			t.Fatalf("bench setup left %s unwritten: %v", rel, err)
		}
	}
	assertPrivateHomeEmpty(t, home)

	wrapper := filepath.Join(repo, ".bench", "bin", "bench.sh")
	gate := func(args ...string) processResult {
		return launch("bash", append([]string{wrapper}, args...)...)
	}

	stub := gate("gate")
	if stub.code != 1 || !strings.Contains(stub.stderr, "configure .bench/gate.sh - replace this sentinel with real checks") {
		t.Fatalf("untouched stub gate = (%d, %q, %q)", stub.code, stub.stdout, stub.stderr)
	}
	assertPrivateHomeEmpty(t, home)

	gateScript := filepath.Join(repo, ".bench", "gate.sh")
	retireSentinelLine(t, gateScript)
	canaryDir := filepath.Join(repo, "tests", "canary")
	if _, err := os.Stat(canaryDir); !os.IsNotExist(err) {
		t.Fatalf("tests/canary exists before the fixture leg: %v", err)
	}
	retired := gate("gate", "--fresh")
	if retired.code != 0 || !hasExactLine(retired.stdout, "gate: green") {
		t.Fatalf("gate with the sentinel retired = (%d, %q, %q)", retired.code, retired.stdout, retired.stderr)
	}
	assertPrivateHomeEmpty(t, home)

	fixture := filepath.Join(canaryDir, "adoption-smoke", "seeded-fixture")
	if err := os.MkdirAll(fixture, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(fixture, "input.txt"), []byte("seeded\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	withFixture := gate("gate", "--fresh")
	if withFixture.code != 0 || !strings.Contains(withFixture.stdout, "canary inventory ok (1 fixture bindings)") || !hasExactLine(withFixture.stdout, "gate: green") {
		t.Fatalf("gate with one project fixture = (%d, %q, %q)", withFixture.code, withFixture.stdout, withFixture.stderr)
	}
	assertPrivateHomeEmpty(t, home)

	manifest := filepath.Join(repo, ".bench", "gate-inputs.json")
	declaration, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(manifest); err != nil {
		t.Fatal(err)
	}
	// The test removes the manifest so HOME cannot reach the gate's environment. The
	// wrapper's own pool-home refusal fires in response. The conformance wrapper-pool-home
	// test pins the exact wording once. This leg asserts only that an undeclared input
	// reds the gate on HOME.
	unbound := gate("gate", "--fresh")
	if unbound.code != 1 || !strings.Contains(unbound.stderr, "HOME:") {
		t.Fatalf("gate without the seeded manifest = (%d, %q, %q)", unbound.code, unbound.stdout, unbound.stderr)
	}
	assertPrivateHomeEmpty(t, home)

	if err := os.WriteFile(manifest, declaration, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(filepath.Join(canaryDir, "adoption-smoke")); err != nil {
		t.Fatal(err)
	}
	empty := gate("gate", "--fresh")
	if empty.code != 1 || !strings.Contains(empty.stderr, "canary fixture inventory is empty") {
		t.Fatalf("gate with an empty tests/canary = (%d, %q, %q)", empty.code, empty.stdout, empty.stderr)
	}
	assertPrivateHomeEmpty(t, home)
}

// retireSentinelLine performs the one documented operator step. It removes exactly the
// line that carries the sentinel marker and leaves the rest of the scaffolded gate intact.
func retireSentinelLine(t *testing.T, path string) {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(string(content), "\n")
	kept := make([]string, 0, len(lines))
	removed := 0
	for _, line := range lines {
		if strings.Contains(line, adopt.SentinelMarker) {
			removed++
			continue
		}
		kept = append(kept, line)
	}
	if removed != 1 {
		t.Fatalf("sentinel lines in %s = %d, want 1", path, removed)
	}
	if err := os.WriteFile(path, []byte(strings.Join(kept, "\n")), 0o755); err != nil {
		t.Fatal(err)
	}
}

// hasExactLine matches a whole output line, never a substring. A reused verdict prints
// "gate: green (fresh verdict reused for this tree)". A substring check on "gate: green"
// would wrongly accept that line as a real green run.
func hasExactLine(output, want string) bool {
	for _, line := range strings.Split(output, "\n") {
		if line == want {
			return true
		}
	}
	return false
}

func assertPrivateHomeEmpty(t *testing.T, home string) {
	t.Helper()
	entries, err := os.ReadDir(home)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("private BENCH_HOME %s is not empty: %v", home, entries)
	}
}
