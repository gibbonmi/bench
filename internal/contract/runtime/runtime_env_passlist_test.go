package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
	"github.com/gibbonmi/bench/internal/gate"
)

// Seam-A sentinel contracts: the environment a Bench-launched subprocess actually
// receives, observed at the built binary rather than at a Go-level unit. Every
// assertion here plants a marker in the parent process and looks for it in the
// child's own dump of its exported names, so an implementation that silently
// failed to launch anything cannot pass by finding nothing — readEnvDump refuses
// a missing or empty dump before any absence is asserted.
func TestRuntimeEnvPasslistContracts(t *testing.T) {
	t.Parallel()
	contract.RunParallel(t, "adapter drops a non-passlisted marker", testAdapterDropsMarker)
	contract.RunParallel(t, "adapter keeps passlisted names", testAdapterKeepsPasslistedNames)
	contract.RunParallel(t, "env.allow [agent] entry reaches the adapter", testAllowAgentEntryReachesAdapter)
	contract.RunParallel(t, "env.allow glob admits every matching name", testAllowGlobAdmitsMatchingNames)
	contract.RunParallel(t, "malformed env.allow refuses the adapter launch", testMalformedAllowRefusesAdapterLaunch)
	contract.RunParallel(t, "stale [gate] section refuses the adapter launch", testStaleGateSectionRefusesLaunch)
	contract.RunParallel(t, "environment sentinel runs in the gate's contract phase", testEnvSentinelIsGateAttached)
}

// envMarkers is the parent environment every seam-A shift run plants. FT88_MARKER
// matches no default passlist pattern and no env.allow entry any fixture writes,
// so its presence in a child dump is a leak by construction; the remaining names
// are chosen to match exactly one documented passlist entry each.
func envMarkers(agent, home string) map[string]string {
	return map[string]string{
		"FT88_MARKER":          "ft88-marker-must-not-reach-a-subprocess",
		"ANTHROPIC_FT88_TOKEN": "ft88-anthropic",
		"FT88_OPT_IN":          "ft88-opt-in",
		"FT88P_ONE":            "one",
		"FT88P_TWO":            "two",
		"BENCH_AGENT":          agent,
		"BENCH_MAX_ITERS":      "1",
		"BENCH_HOME":           home,
	}
}

// envDumpFixture builds a shift-capable repo whose stub adapter writes its own
// exported-variable names to an absolute path outside the repo, and commits the
// optional .bench/env.allow contents so the shift worktree carries them. The dump
// path is baked into the adapter script rather than passed through the
// environment: a path read from the environment would itself be subject to the
// passlist under test, and a filter defect would then look like a test bug.
func envDumpFixture(t *testing.T, allow string) (contract.Fixture, string) {
	t.Helper()
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	dump := filepath.Join(t.TempDir(), "adapter-env")
	f.WriteExecutable("agent", fmt.Sprintf("#!/usr/bin/env bash\n{ compgen -e || true; } > %q\nexit 0\n", dump))
	if allow != "" {
		f.WriteFile(".bench/env.allow", allow)
	}
	f.CommitAll("seam-A adapter")
	return f, dump
}

// runShiftDump drives one capped shift iteration with the markers planted and
// returns the adapter's dumped names. A no-op adapter makes no change, so the
// loop's honest taxonomy is no-op/4 — pinned here so a run that failed for an
// unrelated reason is not mistaken for a clean observation.
func runShiftDump(t *testing.T, f contract.Fixture, dump, objective string) []string {
	t.Helper()
	home := t.TempDir()
	probe := f.BenchEnv(envMarkers(filepath.Join(f.Root, "agent"), home), "shift", objective)
	probe.RequireExit(4)
	return readEnvDump(t, dump)
}

// readEnvDump reads the child's dump and refuses to return until it has proved
// the child actually ran: the file must exist, be non-empty, and carry PATH,
// which every passlist admits. Without this guard an absence assertion would pass
// against a subprocess that never launched.
func readEnvDump(t *testing.T, path string) []string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("subprocess wrote no environment dump at %s: %v — the sentinel proves nothing unless the child ran", path, err)
	}
	names := contract.NonEmptyLines(string(data))
	if len(names) == 0 {
		t.Fatal("subprocess environment dump is empty — the sentinel proves nothing unless the child ran")
	}
	if !hasName(names, "PATH") {
		t.Fatalf("subprocess environment dump carries no PATH: %#v — the child did not launch with a real environment", names)
	}
	return names
}

func hasName(names []string, want string) bool {
	for _, name := range names {
		if name == want {
			return true
		}
	}
	return false
}

func requireNames(t *testing.T, names []string, want ...string) {
	t.Helper()
	for _, name := range want {
		if !hasName(names, name) {
			t.Fatalf("subprocess environment is missing passlisted name %q: %#v", name, names)
		}
	}
}

func requireNoNames(t *testing.T, names []string, unwanted ...string) {
	t.Helper()
	for _, name := range unwanted {
		if hasName(names, name) {
			t.Fatalf("subprocess environment carries non-passlisted name %q: %#v", name, names)
		}
	}
}

// testAdapterDropsMarker is the story-1 claim itself: a name the reviewer happens
// to have exported reaches the harness adapter only if construction fails to
// filter by name.
func testAdapterDropsMarker(t *testing.T) {
	f, dump := envDumpFixture(t, "")
	names := runShiftDump(t, f, dump, "marker-absence")
	requireNoNames(t, names, "FT88_MARKER")
}

// testAdapterKeepsPasslistedNames pins the passlist as a filter rather than a
// wipe. Marker-absence alone would be greened by an adapter launched with an
// empty environment, which breaks every real harness.
func testAdapterKeepsPasslistedNames(t *testing.T) {
	f, dump := envDumpFixture(t, "")
	names := runShiftDump(t, f, dump, "passlist-survival")
	requireNames(t, names, "PATH", "HOME", "ANTHROPIC_FT88_TOKEN", "BENCH_MAX_ITERS")
	requireNoNames(t, names, "FT88_MARKER")
}

// testAllowAgentEntryReachesAdapter proves the committed opt-in is read at all:
// an exact name under [agent] admits a variable no default pattern matches, while
// the marker beside it still does not survive.
func testAllowAgentEntryReachesAdapter(t *testing.T) {
	f, dump := envDumpFixture(t, "# seam-A opt-in\n[agent]\nFT88_OPT_IN\n")
	names := runShiftDump(t, f, dump, "allow-agent-entry")
	requireNames(t, names, "FT88_OPT_IN")
	requireNoNames(t, names, "FT88_MARKER")
}

// testAllowGlobAdmitsMatchingNames pins glob support in the opt-in parser.
// Exact-match-only is the likeliest cheap implementation, and it would admit
// neither of the two matching names while still greening the exact-entry row.
func testAllowGlobAdmitsMatchingNames(t *testing.T) {
	f, dump := envDumpFixture(t, "[agent]\nFT88P_*\n")
	names := runShiftDump(t, f, dump, "allow-glob")
	requireNames(t, names, "FT88P_ONE", "FT88P_TWO")
	requireNoNames(t, names, "FT88_MARKER", "FT88_OPT_IN")
}

// testMalformedAllowRefusesAdapterLaunch pins the fail-closed posture end to end:
// the run exits non-zero, the diagnostic names the offending line number, and —
// the part a parser that skipped bad lines would fail — the adapter never ran, so
// nothing was launched with a silently-widened environment.
func testMalformedAllowRefusesAdapterLaunch(t *testing.T) {
	f, dump := envDumpFixture(t, "[agent]\nFT88_OPT_IN\nBAD ENTRY!\n")
	home := t.TempDir()
	probe := f.BenchEnv(envMarkers(filepath.Join(f.Root, "agent"), home), "shift", "malformed-allow")
	if probe.ExitCode == 0 {
		t.Fatalf("malformed .bench/env.allow exited 0:\n%s\n%s", probe.Stdout, probe.Stderr)
	}
	output := probe.Stdout + probe.Stderr
	if !strings.Contains(output, ".bench/env.allow:3:") {
		t.Fatalf("malformed .bench/env.allow diagnostic does not name the offending line 3:\n%s", output)
	}
	if _, err := os.Stat(dump); !os.IsNotExist(err) {
		t.Fatalf("adapter launched despite a malformed .bench/env.allow (dump stat error = %v)", err)
	}
}

// testStaleGateSectionRefusesLaunch is story 3's unknown-section row for the one
// unknown section the retired gate-class draft would most plausibly leave behind:
// a committed .bench/env.allow carrying a stale [gate] header refuses the shift
// launch, exits non-zero, and names the offending line — the gate opt-in lives in
// the manifest, not this file, so a parser that silently skipped the section would
// be indistinguishable from a working opt-in. The comment above the header keeps
// the header off line 1 so the line number is a real assertion rather than a
// constant.
func testStaleGateSectionRefusesLaunch(t *testing.T) {
	f, dump := envDumpFixture(t, "# stale from the retired gate-class draft\n[gate]\nFT88_OPT_IN\n")
	home := t.TempDir()
	probe := f.BenchEnv(envMarkers(filepath.Join(f.Root, "agent"), home), "shift", "stale-gate-section")
	if probe.ExitCode == 0 {
		t.Fatalf("stale [gate] section exited 0:\n%s\n%s", probe.Stdout, probe.Stderr)
	}
	output := probe.Stdout + probe.Stderr
	if !strings.Contains(output, ".bench/env.allow:2:") {
		t.Fatalf("stale [gate] diagnostic does not name the offending line 2:\n%s", output)
	}
	if _, err := os.Stat(dump); !os.IsNotExist(err) {
		t.Fatalf("adapter launched despite a stale [gate] section (dump stat error = %v)", err)
	}
}

// testEnvSentinelIsGateAttached is the story-5 claim that this proof is the
// gate's, not a human's: the file it lives in must sit under a package the gate's
// contract phase actually executes. It resolves this file's own package rather
// than naming it, so moving the sentinel out of the phase's reach turns it red.
func testEnvSentinelIsGateAttached(t *testing.T) {
	kit := contract.KitRoot(t)
	pkg, err := filepath.Rel(kit, sentinelPackageDir(t))
	if err != nil {
		t.Fatalf("locate the sentinel package under the kit root: %v", err)
	}
	pkg = filepath.ToSlash(pkg)
	for _, phase := range gate.BenchkitPhases(kit, kit) {
		if phase.Name != "contract" {
			continue
		}
		for _, arg := range phase.Argv {
			if packagePatternCovers(arg, pkg) {
				return
			}
		}
		t.Fatalf("the gate's contract phase does not run the environment sentinel's package %q: argv %#v", pkg, phase.Argv)
	}
	t.Fatal("the gate has no contract phase, so the environment sentinel is not gate-attached")
}

// packagePatternCovers reports whether a `go test` package pattern argument
// selects pkg — either the package itself or, for a `/...` pattern, any package
// beneath it.
func packagePatternCovers(arg, pkg string) bool {
	pattern := strings.TrimPrefix(arg, "./")
	if pattern == arg || pattern == "" {
		return false
	}
	if strings.HasSuffix(pattern, "/...") {
		root := strings.TrimSuffix(pattern, "/...")
		return pkg == root || strings.HasPrefix(pkg, root+"/")
	}
	return pkg == pattern
}

func sentinelPackageDir(t *testing.T) string {
	t.Helper()
	dir, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("resolve the sentinel package directory: %v", err)
	}
	return dir
}
