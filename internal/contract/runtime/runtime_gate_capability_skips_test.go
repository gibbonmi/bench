package runtime

// Strict capability evidence through the built binary: `bench gate-phases` against a
// fixture kit whose phase tests append real skip lines to the log the gate hands them.

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/contract"
)

func TestRuntimeGateCapabilitySkipContracts(t *testing.T) {
	contract.SkipIfSubjectFileMissing(t, "dist/bench")
	t.Parallel()
	contract.RunParallel(t, "strict capability skip is red naming the class", testRuntimeStrictCapabilitySkipIsRed)
	contract.RunParallel(t, "unset strict flag keeps skips informational", testRuntimeCapabilitySkipRowsWithoutStrict)
	contract.RunParallel(t, "environment skip never trips strict mode", testRuntimeEnvironmentSkipIsNotStrict)
}

func testRuntimeStrictCapabilitySkipIsRed(t *testing.T) {
	contract.NoteContractFailure(t, "strict capability-skip verdict contract failed")
	probe := runCapabilitySkipGate(t, map[string]string{"BENCH_REQUIRE_CAPABILITIES": "1"}, capability.Skip{
		Kind: capability.KindCapability, Class: capability.Symlink, Reason: "host cannot create unprivileged symlinks",
	})
	if probe.ExitCode == 0 {
		t.Fatalf("strict gate stayed green with a capability skip:\n%s\n%s", probe.Stdout, probe.Stderr)
	}
	output := probe.Stdout + probe.Stderr
	for _, want := range []string{"BENCH_REQUIRE_CAPABILITIES=1", "symlink"} {
		if !strings.Contains(output, want) {
			t.Fatalf("strict verdict does not name %q:\n%s", want, output)
		}
	}
}

func testRuntimeCapabilitySkipRowsWithoutStrict(t *testing.T) {
	contract.NoteContractFailure(t, "informational capability-skip row contract failed")
	probe := runCapabilitySkipGate(t, nil, capability.Skip{
		Kind: capability.KindCapability, Class: capability.Symlink, Reason: "host cannot create unprivileged symlinks",
	})
	probe.RequireExit(0)
	for _, want := range []string{"capability-skips: 1 (capability=1 environment=0)", "capability-skips class=symlink: 1"} {
		if !strings.Contains(probe.Stdout, want) {
			t.Fatalf("gate output missing row %q:\n%s", want, probe.Stdout)
		}
	}
}

func testRuntimeEnvironmentSkipIsNotStrict(t *testing.T) {
	contract.NoteContractFailure(t, "environment-skip strict-exclusion contract failed")
	probe := runCapabilitySkipGate(t, map[string]string{"BENCH_REQUIRE_CAPABILITIES": "1"}, capability.Skip{
		Kind: capability.KindEnvironment, Reason: "fixture binary was never staged",
	})
	probe.RequireExit(0)
	if want := "capability-skips: 1 (capability=0 environment=1)"; !strings.Contains(probe.Stdout, want) {
		t.Fatalf("gate output missing row %q:\n%s", want, probe.Stdout)
	}
}

// runCapabilitySkipGate drives the built binary's phase runner over a fixture kit
// whose conformance phase reports skip. The gate reaches the skip through the log it
// sets for every phase it launches, which is the only transport that survives:
// `go test` without -v discards a passing package's output entirely.
func runCapabilitySkipGate(t *testing.T, env map[string]string, skip capability.Skip) contract.Probe {
	t.Helper()
	kit := capabilitySkipKit(t, skip)
	graded := t.TempDir()
	f := contract.NewExecFixtureAt(t, kit)
	run := map[string]string{"BENCH_KIT": kit}
	for key, value := range env {
		run[key] = value
	}
	return contract.RunAt(t, f, kit, run, filepath.Join(contract.SubjectRoot(t), "dist", "bench"), "gate-phases", graded)
}

// capabilitySkipKit builds the smallest tree BenchkitPhases resolves to green phases:
// a module holding the conformance and contract test packages the gate names, and a
// wrapper stub for the canary and shellcheck phases. The conformance test appends the
// rendered skip line itself rather than restating the line shape.
func capabilitySkipKit(t *testing.T, skip capability.Skip) string {
	t.Helper()
	line, err := capability.Render(skip)
	if err != nil {
		t.Fatalf("render fixture skip: %v", err)
	}
	kit := t.TempDir()
	contract.WriteFileAbs(t, filepath.Join(kit, "go.mod"), "module benchcapabilityfixture\n\n"+subjectGoDirective(t)+"\n")
	contract.WriteFileAbs(t, filepath.Join(kit, "internal", "conformance", "root_test.go"), skipEmittingSource("conformance", "TestRootConformance", line))
	contract.WriteFileAbs(t, filepath.Join(kit, "internal", "contract", "probe", "probe_test.go"), skipEmittingSource("probe", "TestProbe", ""))
	contract.WriteExecutableAbs(t, filepath.Join(kit, "bin", "bench.sh"), "#!/usr/bin/env bash\nexit 0\n")
	return kit
}

// skipEmittingSource writes a test that appends line to the gate's skip log, or only
// asserts the log was handed to it when line is empty. Failing on an unset variable
// makes a gate that never points its phases at a log red here rather than silently
// reporting zero skips forever.
func skipEmittingSource(pkg, name, line string) string {
	return fmt.Sprintf(`package %s

import (
	"os"
	"testing"
)

func %s(t *testing.T) {
	path := os.Getenv(%q)
	if path == "" {
		t.Fatal("gate did not point the phase at a capability skip log")
	}
	line := %q
	if line == "" {
		return
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatalf("open skip log: %%v", err)
	}
	defer f.Close()
	if _, err := f.WriteString(line); err != nil {
		t.Fatalf("append skip log: %%v", err)
	}
}
`, pkg, name, capability.LogEnv, line)
}

// subjectGoDirective reuses the subject's own language version so the fixture module
// compiles under exactly the toolchain the gate runs.
func subjectGoDirective(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(contract.SubjectRoot(t), "go.mod"))
	if err != nil {
		t.Fatalf("read subject go.mod: %v", err)
	}
	match := regexp.MustCompile(`(?m)^go .*$`).FindString(string(data))
	if match == "" {
		t.Fatal("subject go.mod carries no go directive")
	}
	return match
}
