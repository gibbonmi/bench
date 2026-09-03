package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/runbinary"
)

func TestRootConformance(t *testing.T) {
	root := os.Getenv(registry.ConformanceRootEnv)
	if root == "" {
		capability.Environment(t, registry.ConformanceRootEnv+" not set")
	}
	h := NewHarness(t)
	selected, selectedSet := os.LookupEnv(registry.ConformanceChecksEnv)
	var selectedValue *string
	if selectedSet {
		selectedValue = &selected
	}
	inherited, inheritedSet := os.LookupEnv(registry.ConformanceInheritedEnv)
	var inheritedValue *string
	if inheritedSet {
		inheritedValue = &inherited
	}
	for _, diag := range RunConformanceSelection(root, h.KitRoot, registry.TierFor(os.Getenv(registry.ConformanceTierEnv)), os.Getenv(registry.ConformanceScopeEnv), selectedValue, inheritedValue) {
		t.Errorf("gate: %s", diag)
	}
}

func checkGateEntryContract(root string) []string {
	path := filepath.Join(root, ".bench", "gate.sh")
	if !exists(path) {
		return nil
	}
	gate := readIfExists(path)
	var diags []string
	needle := `exec env BENCH_KIT="$kit" BENCH_RUN_BINARY="$bench" "$bench" gate-phases "$root"`
	if !strings.Contains(gate, needle) {
		diags = append(diags, fmt.Sprintf(".bench/gate.sh missing %q", needle))
	}
	check := `"$bench" freshness-check "$kit"`
	if strings.Index(gate, check) < 0 || strings.Index(gate, check) > strings.Index(gate, needle) {
		diags = append(diags, fmt.Sprintf(".bench/gate.sh does not run current-source verification %q before %q", check, needle))
	}

	for _, retired := range []string{
		`BENCH_CONFORMANCE_ROOT="$root" go test -count=1 ./internal/conformance -run '^TestRootConformance$'`,
		`BENCH_CONTRACT_ROOT="$root" go test -count=1 ./internal/contract/...`,
		`bin/bench.sh" canary "$root"`,
		"gate-docs-contracts.sh",
		"gate-line-contracts.sh",
		"gate-package-contracts.sh",
		"gate-runtime-shift-contracts.sh",
		"gate-contract-runner.sh",
		"gate-go-contracts.sh",
		"gate-link-contracts.sh",
		"gate-runtime-contracts.sh",
		"gate-runtime-git-contracts.sh",
		"gate-doctor-contracts.sh",
		"gate-axi-contracts.sh",
		"gate-axi-guards-contracts.sh",
		"gate-axi-wave2-contracts.sh",
	} {
		if strings.Contains(gate, retired) {
			diags = append(diags, ".bench/gate.sh still references retired conformance fragment "+retired)
		}
	}
	return diags
}

// gateEntryWrapperAction is the invocation the wrapper-less refusal must name. The run-
// binary variable is wrapper-owned. Naming the variable alone leaves a session that meets
// the refusal with no next action.
const gateEntryWrapperAction = "bash bin/bench.sh gate"

func TestGateEntryRefusesUnverifiedBinaryBeforeGatePhases(t *testing.T) {
	h := NewHarness(t)
	_, _, nested, gatePath := newGateEntryFixture(t, h)

	stripped := capability.WithoutEnvironment(os.Environ(), runbinary.Env)
	absent := gateEntryRefusal(t, nested, stripped, gatePath)
	if !strings.Contains(absent, runbinary.Env) {
		t.Fatalf("gate entry missing selected-path refusal:\n%s", absent)
	}
	if !strings.Contains(absent, gateEntryWrapperAction) {
		t.Fatalf("gate entry refusal names no next action, want %q:\n%s", gateEntryWrapperAction, absent)
	}
	if strings.Contains(strings.ToLower(absent), "worktree") {
		t.Fatalf("gate entry refusal names a worktree, but it fires for every wrapper-less caller:\n%s", absent)
	}
	if strings.Contains(absent, "phase ") {
		t.Fatalf("gate entry resolved phases before refusing its unverified binary:\n%s", absent)
	}

	relative := gateEntryRefusal(t, nested, runbinary.WithEnv(stripped, "dist/bench"), gatePath)
	if relative != absent {
		t.Fatalf("relative run-binary refusal = %q, want the absent-value refusal %q", relative, absent)
	}
}

// TestGateEntryKeepsTheLaterRefusalWording pins the two refusals that follow the wrapper-
// less one. Each refusal already names a condition an operator can act on. A reword of
// the first refusal must not reach these two.
func TestGateEntryKeepsTheLaterRefusalWording(t *testing.T) {
	h := NewHarness(t)
	_, kit, nested, gatePath := newGateEntryFixture(t, h)
	stripped := capability.WithoutEnvironment(os.Environ(), runbinary.Env)

	plain := filepath.Join(kit, "dist", "not [*] executable")
	writeFixtureFile(t, plain, "#!/usr/bin/env bash\nexit 0\n")
	if got, want := gateEntryRefusal(t, nested, runbinary.WithEnv(stripped, plain), gatePath), "error: "+runbinary.Env+" is not a regular executable"; !strings.Contains(got, want) {
		t.Fatalf("non-executable refusal = %q, want it to contain %q", got, want)
	}

	physical := filepath.Join(kit, "physical [*] dir")
	linked := filepath.Join(kit, "linked [*] dir")
	if err := os.MkdirAll(physical, 0o755); err != nil {
		t.Fatalf("mkdir physical gate dir: %v", err)
	}
	if err := os.Symlink(physical, linked); err != nil {
		capability.Environment(t, "symbolic links unavailable: "+err.Error())
	}
	staged := filepath.Join(physical, "bench")
	writeFixtureFile(t, staged, "#!/usr/bin/env bash\nexit 0\n")
	if err := os.Chmod(staged, 0o755); err != nil {
		t.Fatalf("chmod staged gate binary: %v", err)
	}
	viaLink := filepath.Join(linked, "bench")
	if got, want := gateEntryRefusal(t, nested, runbinary.WithEnv(stripped, viaLink), gatePath), "error: "+runbinary.Env+" must be a cleaned physical path"; !strings.Contains(got, want) {
		t.Fatalf("physical-path refusal = %q, want it to contain %q", got, want)
	}
}

func TestGateEntryRejectsLegacyBeforeRunningOldTableAndRunsReplacementOnce(t *testing.T) {
	h := NewHarness(t)
	root, kit, nested, gatePath := newGateEntryFixture(t, h)
	publishGateFixtureBench(t, kit, fmt.Sprintf("#!/usr/bin/env bash\ncase \"$1\" in\nfreshness-check) test \"$2\" = %q || exit 97; exit 2 ;;\ngate-phases) test \"$BENCH_KIT\" = %q || exit 98; printf 'old phase\\n'; printf 'old\\n' >> \"$2/.git/gate-phases-ran\" ;;\nesac\n", kit, kit))
	selected := filepath.Join(kit, "dist", "bench")
	env := runbinary.WithEnv(os.Environ(), selected)
	env = append(env, "BENCH_KIT="+kit)

	legacy := runAtEnv(nested, env, "bash", gatePath)
	if legacy == nil || legacy.ExitCode == 0 {
		t.Fatalf("legacy gate entry exit = %+v, want refusal", legacy)
	}
	legacyOutput := legacy.Stdout + legacy.Stderr
	if strings.Contains(legacyOutput, "old phase") {
		t.Fatalf("legacy gate entry ran the old table:\n%s", legacyOutput)
	}
	if _, err := os.Stat(filepath.Join(root, ".git", "gate-phases-ran")); !os.IsNotExist(err) {
		t.Fatalf("legacy gate entry scheduled the old table: %v", err)
	}

	publishGateFixtureBench(t, kit, fmt.Sprintf("#!/usr/bin/env bash\ncase \"$1\" in\nfreshness-check) test \"$2\" = %q ;;\ngate-phases) test \"$BENCH_KIT\" = %q || exit 98; printf 'replacement phase\\n'; printf 'replacement\\n' >> \"$2/.git/gate-phases-ran\" ;;\nesac\n", kit, kit))
	replacement := runAtEnv(nested, env, "bash", gatePath)
	if replacement == nil || replacement.ExitCode != 0 {
		t.Fatalf("replacement gate entry exit = %+v, want green", replacement)
	}
	if output := replacement.Stdout + replacement.Stderr; strings.Count(output, "replacement phase") != 1 {
		t.Fatalf("replacement gate output = %q, want replacement phase exactly once", output)
	}
	repeated := runAtEnv(root, env, "bash", gatePath)
	if repeated == nil || repeated.ExitCode != 0 {
		t.Fatalf("repeated fresh gate entry exit = %+v, want green", repeated)
	}
	if output := repeated.Stdout + repeated.Stderr; strings.Count(output, "replacement phase") != 1 {
		t.Fatalf("repeated fresh gate output = %q, want replacement phase exactly once", output)
	}
	runs, err := os.ReadFile(filepath.Join(root, ".git", "gate-phases-ran"))
	if err != nil {
		t.Fatalf("read replacement phase record: %v", err)
	}
	if got, want := string(runs), "replacement\nreplacement\n"; got != want {
		t.Fatalf("phase table runs = %q, want %q", got, want)
	}
}

// gateEntryRefusal runs the gate entry as a real subprocess. It returns the combined
// output of a run that must refuse before it reaches the phase table.
func gateEntryRefusal(t *testing.T, dir string, env []string, gatePath string) string {
	t.Helper()
	probe := runAtEnv(dir, env, "bash", gatePath)
	if probe == nil || probe.ExitCode == 0 {
		t.Fatalf("gate entry exit = %+v, want a refusal", probe)
	}
	return probe.Stdout + probe.Stderr
}

// newGateEntryFixture stages a gate entry under a git root. It returns the root, the kit
// it reads, a nested working directory whose glob characters must survive quoting, and
// the script path.
func newGateEntryFixture(t *testing.T, h Harness) (root, kit, nested, gatePath string) {
	t.Helper()
	root = t.TempDir()
	if probe := runAt(root, "git", "init", "-q"); probe == nil || probe.ExitCode != 0 {
		t.Fatalf("git init failed: %+v", probe)
	}
	nested = filepath.Join(root, "nested [*] path")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir nested gate cwd: %v", err)
	}
	kit = filepath.Join(t.TempDir(), "kit [*] path")
	return root, kit, nested, writeGateEntryFixture(t, h, kit)
}

func writeGateEntryFixture(t *testing.T, h Harness, kit string) string {
	t.Helper()
	for _, rel := range []string{
		".bench/gate.sh",
		"internal/freshness/freshness.go",
		"internal/freshness/check/main.go",
		"scripts/go-build.sh",
		"scripts/go-build.inputs",
		"package.json",
		"internal/releaseevidence/requirements.json",
	} {
		data, err := os.ReadFile(h.KitPath(filepath.FromSlash(rel)))
		if err != nil {
			t.Fatalf("read gate fixture source %s: %v", rel, err)
		}
		writeFixtureFile(t, filepath.Join(kit, filepath.FromSlash(rel)), string(data))
	}
	writeFixtureFile(t, filepath.Join(kit, "go.mod"), "module github.com/gibbonmi/bench\n\ngo 1.25\n")
	writeFixtureFile(t, filepath.Join(kit, "cmd", "bench", "main.go"), "package main\n\nfunc main() {}\n")
	path := filepath.Join(kit, ".bench", "gate.sh")
	if err := os.Chmod(path, 0o755); err != nil {
		t.Fatalf("chmod gate entry: %v", err)
	}
	return path
}

func publishGateFixtureBench(t *testing.T, kit, program string) {
	t.Helper()
	staged := filepath.Join(kit, "dist", "staged")
	writeFixtureFile(t, staged, program)
	if err := os.Chmod(staged, 0o755); err != nil {
		t.Fatalf("chmod staged fixture bench: %v", err)
	}
	if err := freshness.Publish(kit, staged, filepath.Join(kit, "dist", "bench"), "1.2.3"); err != nil {
		t.Fatalf("publish fixture bench: %v", err)
	}
}

func removeGateArtifact(t *testing.T, path string) {
	t.Helper()
	if err := os.Remove(path); err != nil {
		t.Fatalf("remove gate fixture artifact: %v", err)
	}
}
