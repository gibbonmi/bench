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
	for _, diag := range RunConformanceSelection(root, h.KitRoot, registry.TierFor(os.Getenv(registry.ConformanceTierEnv)), "", selectedValue, inheritedValue) {
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

func TestGateEntryRefusesUnverifiedBinaryBeforeGatePhases(t *testing.T) {
	h := NewHarness(t)
	root := t.TempDir()
	if probe := runAt(root, "git", "init", "-q"); probe == nil || probe.ExitCode != 0 {
		t.Fatalf("git init failed: %+v", probe)
	}
	nested := filepath.Join(root, "nested [*] path")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir nested gate cwd: %v", err)
	}
	kit := filepath.Join(t.TempDir(), "kit [*] path")
	gatePath := writeGateEntryFixture(t, h, kit)

	env := capability.WithoutEnvironment(os.Environ(), runbinary.Env)
	probe := runAtEnv(nested, env, "bash", gatePath)
	if probe == nil || probe.ExitCode == 0 {
		t.Fatalf("gate entry exit = %+v, want an untrusted-binary refusal", probe)
	}
	output := probe.Stdout + probe.Stderr
	if !strings.Contains(output, runbinary.Env) {
		t.Fatalf("gate entry missing selected-path refusal:\n%s", output)
	}
	if strings.Contains(output, "phase ") {
		t.Fatalf("gate entry resolved phases before refusing its unverified binary:\n%s", output)
	}
}

func TestGateEntryRejectsLegacyBeforeRunningOldTableAndRunsReplacementOnce(t *testing.T) {
	h := NewHarness(t)
	root := t.TempDir()
	if probe := runAt(root, "git", "init", "-q"); probe == nil || probe.ExitCode != 0 {
		t.Fatalf("git init failed: %+v", probe)
	}
	nested := filepath.Join(root, "nested [*] path")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir nested gate cwd: %v", err)
	}
	kit := filepath.Join(t.TempDir(), "kit [*] path")
	gatePath := writeGateEntryFixture(t, h, kit)
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
	if err := freshness.Publish(kit, staged, filepath.Join(kit, "dist", "bench")); err != nil {
		t.Fatalf("publish fixture bench: %v", err)
	}
}

func removeGateArtifact(t *testing.T, path string) {
	t.Helper()
	if err := os.Remove(path); err != nil {
		t.Fatalf("remove gate fixture artifact: %v", err)
	}
}
