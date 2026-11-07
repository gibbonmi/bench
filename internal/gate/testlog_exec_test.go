package gate

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// gateGoShimEnv makes this test binary serve as `bench gate-go`, so a phase runs the real
// trace shim with no built Bench.
const gateGoShimEnv = "BENCH_TEST_GATE_GO"

// init runs before TestMain, so a test binary that serves as `bench gate-go` runs no test.
func init() {
	if os.Getenv(gateGoShimEnv) != "" && len(os.Args) > 1 && os.Args[1] == "gate-go" {
		os.Exit(GateGoCommand(os.Args[2:], os.Stdout, os.Stderr))
	}
}

// shimPackage is the package directory, below the fixture root, where runShim runs the shim.
// shimOperands are the arguments that the shim passes on to its probe binary.
const shimPackage = "pkg/sub"

var shimOperands = []string{"-a", "b c"}

// recordArgs is a probe body that writes its arguments, one to a line, beside the probe.
const recordArgs = `printf '%s\n' "$@" > "$0.args"` + "\n"

// TestGateGoExecAddsTheTestLog requires that, with a log directory, the shim runs the binary
// with the log flag of its package before its own arguments.
func TestGateGoExecAddsTheTestLog(t *testing.T) {
	logs := t.TempDir()
	probe := probeBinary(t, recordArgs)
	root := runShim(t, probe, map[string]string{testLogEnv: logs})
	log := shimLog(t, logs, root)
	if want := append([]string{"-test.testlogfile=" + log}, shimOperands...); !slices.Equal(probeArgs(t, probe), want) {
		t.Fatalf("probe arguments = %q, want %q", probeArgs(t, probe), want)
	}
}

// TestGateGoExecWithoutALogDirectory requires that, with no log directory, the shim runs the
// binary with its arguments unchanged.
func TestGateGoExecWithoutALogDirectory(t *testing.T) {
	probe := probeBinary(t, recordArgs)
	runShim(t, probe, map[string]string{})
	if got := probeArgs(t, probe); !slices.Equal(got, shimOperands) {
		t.Fatalf("probe arguments = %q, want %q", got, shimOperands)
	}
}

// TestGateGoExecCreatesTheLogFirst requires that a binary that exits at once leaves an empty log.
func TestGateGoExecCreatesTheLogFirst(t *testing.T) {
	logs := t.TempDir()
	root := runShim(t, probeBinary(t, "exit 0\n"), map[string]string{testLogEnv: logs})
	if info, err := os.Stat(shimLog(t, logs, root)); err != nil || info.Size() != 0 {
		t.Fatalf("log of a binary that exits at once = %v, %v; want an empty file", info, err)
	}
}

// TestCreateTestLogKeysThePackageByTheConformanceRoot requires that the root-conformance root
// of the process names the log: the working package is keyed by its path below that root.
func TestCreateTestLogKeysThePackageByTheConformanceRoot(t *testing.T) {
	root, logs := t.TempDir(), t.TempDir()
	pkg := filepath.Join(root, "internal", "sample")
	if err := os.MkdirAll(pkg, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(pkg)
	t.Setenv(registry.ConformanceRootEnv, root)
	path, err := createTestLog(logs)
	if want := filepath.Join(logs, "internal%2Fsample.log"); err != nil || path != want {
		t.Fatalf("createTestLog = %q, %v; want %q", path, err, want)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("the keyed log was not created: %v", err)
	}
}

// TestTestLogKeysAreDistinct requires that two directories that differ only by a separator,
// or by an escaped separator, map to two different log names, each one path element.
func TestTestLogKeysAreDistinct(t *testing.T) {
	for _, pair := range [][2]string{{"a/b", "a__b"}, {"a/b", "a%2Fb"}} {
		first, second := testLogKey(pair[0]), testLogKey(pair[1])
		if first == second || strings.Contains(first, "/") || strings.Contains(second, "/") {
			t.Errorf("keys of %q and %q = %q and %q, want two different file names", pair[0], pair[1], first, second)
		}
	}
}

// TestWithRunBinaryReplacesTheExecToken requires that the selection path, quoted, replaces the
// token inside the `-exec=` element. A path that holds a single quote takes double quotes.
func TestWithRunBinaryReplacesTheExecToken(t *testing.T) {
	for path, want := range map[string]string{
		"/run binary/bench": "-exec='/run binary/bench' gate-go exec",
		"/it's/bench":       `-exec="/it's/bench" gate-go exec`,
	} {
		phases := withRunBinary([]Phase{{Name: canary.PhaseTest, Argv: []string{"go", "test", traceExecOperand, "./..."}}}, &runbinary.Selection{Path: path, SourceRoot: t.TempDir()})
		if got := phases[0].Argv[2]; got != want {
			t.Errorf("exec element for %q = %q, want %q", path, got, want)
		}
	}
}

// TestLinkedTestPhaseCarriesNoTrace requires that a linked root's test phase carries no trace.
func TestLinkedTestPhaseCarriesNoTrace(t *testing.T) {
	kit, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	writeLaneFile(t, filepath.Join(root, "go.mod"), "module linked\n\ngo 1.24\n")
	for _, phase := range BenchkitPhases(root, kit) {
		if phase.Name != canary.PhaseTest {
			continue
		}
		if slices.ContainsFunc(phase.Argv, func(arg string) bool { return strings.HasPrefix(arg, "-exec=") }) {
			t.Fatalf("linked-root test argv = %q, want no -exec= element", phase.Argv)
		}
		return
	}
	t.Fatal("linked root produced no test phase")
}

// TestWholeTestPhaseWritesOneLogPerPackage requires that a run of the kit table with a log
// directory leaves one Go test log for each test package.
func TestWholeTestPhaseWritesOneLogPerPackage(t *testing.T) {
	root := tracedKit(t, "")
	logs := t.TempDir()
	t.Setenv(testLogEnv, logs)
	stdout, stderr := runTracedTable(t, root, testExecutable(t))
	entries, err := os.ReadDir(logs)
	if err != nil || len(entries) != len(tracedPackages) {
		t.Fatalf("logs = %v, %v; want one for each of %q\nstdout:\n%s\nstderr:\n%s", entries, err, tracedPackages, stdout, stderr)
	}
	for _, pkg := range tracedPackages {
		log, err := testLogFile(logs, root, filepath.Join(root, filepath.FromSlash(pkg)))
		if err != nil {
			t.Fatal(err)
		}
		if data, err := os.ReadFile(log); err != nil || !strings.HasPrefix(string(data), "# test log\n") {
			t.Errorf("log of %s = %q, %v; want a Go test log", pkg, data, err)
		}
	}
}

// TestTraceExecSurvivesASpaceInThePath requires that the test phase reaches the planted red
// through a shim whose path holds a space. It clears an inherited log directory, so the run
// writes into no outer run's logs.
func TestTraceExecSurvivesASpaceInThePath(t *testing.T) {
	t.Setenv(testLogEnv, "")
	root := tracedKit(t, `t.Fatal("reached through the shim")`)
	spaced := filepath.Join(t.TempDir(), "run binary", "bench")
	if err := os.Mkdir(filepath.Dir(spaced), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(testExecutable(t), spaced); err != nil {
		t.Fatal(err)
	}
	stdout, stderr := runTracedTable(t, root, spaced)
	reached := func(row string) bool { return strings.Contains(row, "reached through the shim") }
	if !slices.ContainsFunc(rowsForPhase(t, stdout, canary.PhaseTest), reached) {
		t.Fatalf("test phase through %q did not reach the planted red\nstdout:\n%s\nstderr:\n%s", spaced, stdout, stderr)
	}
}

// TestPhaseChildrenDropTheTestLogDirectory requires that no phase child inherits a test log
// directory.
func TestPhaseChildrenDropTheTestLogDirectory(t *testing.T) {
	t.Setenv(testLogEnv, t.TempDir())
	env, err := gateEnv()
	if err != nil || slices.ContainsFunc(env, func(entry string) bool { return strings.HasPrefix(entry, testLogEnv+"=") }) {
		t.Fatalf("phase child environment = %q, %v; want no test log directory", env, err)
	}
}

// probeBinary writes a shell script with body outside every fixture root and answers its path.
func probeBinary(t *testing.T, body string) string {
	t.Helper()
	probe := filepath.Join(t.TempDir(), "probe")
	if err := os.WriteFile(probe, []byte("#!/bin/sh\n"+body), 0o755); err != nil {
		t.Fatal(err)
	}
	return probe
}

// probeArgs answers the arguments that a recordArgs probe wrote.
func probeArgs(t *testing.T, probe string) []string {
	t.Helper()
	return strings.Split(strings.TrimSuffix(string(outcomeRead(t, probe+".args")), "\n"), "\n")
}

// runShim runs `bench gate-go exec <probe>` with shimOperands as the one phase of a fixture
// root, in shimPackage, with env added to the phase. It answers the root.
func runShim(t *testing.T, probe string, env map[string]string) string {
	t.Helper()
	phase := map[string]any{"name": "shim", "argv": append([]string{testExecutable(t), "gate-go", gateGoExecStep, probe}, shimOperands...), "dir": shimPackage, "env": env}
	manifest, err := json.Marshal(map[string]any{"phases": []any{phase}})
	if err != nil {
		t.Fatal(err)
	}
	root := fixturePhaseRoot(t, string(manifest))
	if err := os.MkdirAll(filepath.Join(root, shimPackage), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv(gateGoShimEnv, "1")
	t.Setenv(registry.ConformanceRootEnv, root)
	if code, stdout, stderr := runFixturePhases(context.Background(), t, root); code != 0 {
		t.Fatalf("shim phase = %d; stdout=%q stderr=%q", code, stdout, stderr)
	}
	return root
}

// shimLog answers the log path of shimPackage under logs.
func shimLog(t *testing.T, logs, root string) string {
	t.Helper()
	log, err := testLogFile(logs, root, filepath.Join(root, filepath.FromSlash(shimPackage)))
	if err != nil {
		t.Fatal(err)
	}
	return log
}

// tracedPackages are the test packages of tracedKit. Two of them differ only by a separator.
var tracedPackages = []string{"a", "a/b", "a__b", conformancePackagePath}

// tracedKit is a fixture kit whose whole test phase is traced, because it declares the
// root-conformance entry. The test of package a runs body.
func tracedKit(t *testing.T, body string) string {
	t.Helper()
	t.Setenv(baselinePolicyEnv, "")
	root := gittest.Repo(t)
	outcomeWrite(t, root, "go.mod", "module fixture\n\ngo 1.21\n", 0o644)
	for _, pkg := range tracedPackages {
		name, testBody := "TestTraced", ""
		switch pkg {
		case conformancePackagePath:
			name = registry.RootConformanceTest
		case "a":
			testBody = body
		}
		outcomeWrite(t, root, pkg+"/traced_test.go", goTest(filepath.Base(pkg), name, testBody), 0o644)
	}
	return root
}

// runTracedTable runs the phase table of the kit root with runBinary as the selected
// executable, which serves as `bench gate-go`. It answers stdout and stderr.
func runTracedTable(t *testing.T, root, runBinary string) (string, string) {
	t.Helper()
	t.Setenv(gateGoShimEnv, "1")
	var stdout, stderr bytes.Buffer
	phasesCommandAtKitWithSelection(context.Background(), root, root, &runbinary.Selection{Path: runBinary, SourceRoot: root}, &stdout, &stderr)
	return stdout.String(), stderr.String()
}

// testExecutable answers the path of this test binary.
func testExecutable(t *testing.T) string {
	t.Helper()
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	return exe
}
