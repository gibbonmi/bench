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
)

func TestRootConformance(t *testing.T) {
	root := os.Getenv("BENCH_CONFORMANCE_ROOT")
	if root == "" {
		capability.Environment(t, "BENCH_CONFORMANCE_ROOT not set")
	}
	h := NewHarness(t)
	// The scope env is passed through verbatim: any normalising here would let a stale
	// or misspelled value slide into a silent full run instead of the driver's red.
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
	for _, diag := range RunConformanceSelection(root, h.KitRoot, registry.TierFor(os.Getenv(registry.ConformanceTierEnv)), os.Getenv(registry.ConformanceCheckEnv), selectedValue, inheritedValue) {
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
	needle := `exec env BENCH_KIT="$kit" "$bench" gate-phases "$root"`
	if !strings.Contains(gate, needle) {
		diags = append(diags, fmt.Sprintf(".bench/gate.sh missing %q", needle))
	}
	check := `go run ./internal/freshness/check "$kit" "$bench"`
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

	probe := runAt(nested, "bash", gatePath)
	if probe == nil || probe.ExitCode == 0 {
		t.Fatalf("gate entry exit = %+v, want an untrusted-binary refusal", probe)
	}
	output := probe.Stdout + probe.Stderr
	rebuild := freshness.RebuildAction(kit)
	if !strings.Contains(output, rebuild) {
		t.Fatalf("gate entry missing freshness rebuild action %q:\n%s", rebuild, output)
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

	legacy := runAt(nested, "bash", gatePath)
	if legacy == nil || legacy.ExitCode == 0 {
		t.Fatalf("legacy gate entry exit = %+v, want refusal", legacy)
	}
	legacyOutput := legacy.Stdout + legacy.Stderr
	if rebuild := freshness.RebuildAction(kit); !strings.Contains(legacyOutput, rebuild) {
		t.Fatalf("legacy gate entry missing freshness rebuild action %q:\n%s", rebuild, legacyOutput)
	}
	if strings.Contains(legacyOutput, "old phase") {
		t.Fatalf("legacy gate entry ran the old table:\n%s", legacyOutput)
	}
	if _, err := os.Stat(filepath.Join(root, ".git", "gate-phases-ran")); !os.IsNotExist(err) {
		t.Fatalf("legacy gate entry scheduled the old table: %v", err)
	}

	publishGateFixtureBench(t, kit, fmt.Sprintf("#!/usr/bin/env bash\ncase \"$1\" in\nfreshness-check) test \"$2\" = %q ;;\ngate-phases) test \"$BENCH_KIT\" = %q || exit 98; printf 'replacement phase\\n'; printf 'replacement\\n' >> \"$2/.git/gate-phases-ran\" ;;\nesac\n", kit, kit))
	replacement := runAt(nested, "bash", gatePath)
	if replacement == nil || replacement.ExitCode != 0 {
		t.Fatalf("replacement gate entry exit = %+v, want green", replacement)
	}
	if output := replacement.Stdout + replacement.Stderr; strings.Count(output, "replacement phase") != 1 {
		t.Fatalf("replacement gate output = %q, want replacement phase exactly once", output)
	}
	repeated := runAt(root, "bash", gatePath)
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

func TestGateEntryNormalizesIndeterminateFreshnessFailures(t *testing.T) {
	h := NewHarness(t)
	root := t.TempDir()
	if probe := runAt(root, "git", "init", "-q"); probe == nil || probe.ExitCode != 0 {
		t.Fatalf("git init failed: %+v", probe)
	}
	nested := filepath.Join(root, "nested [*] path")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir nested gate cwd: %v", err)
	}
	for _, failure := range []struct {
		name   string
		mutate func(*testing.T, string)
	}{
		{name: "missing executable", mutate: func(t *testing.T, bench string) { removeGateArtifact(t, bench) }},
		{name: "missing seal", mutate: func(t *testing.T, bench string) { removeGateArtifact(t, bench+".seal") }},
		{name: "unreadable seal", mutate: func(t *testing.T, bench string) {
			seal := bench + ".seal"
			t.Cleanup(func() { _ = os.Chmod(seal, 0o644) })
			if err := os.Chmod(seal, 0); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "malformed seal", mutate: func(t *testing.T, bench string) { writeFixtureFile(t, bench+".seal", "{}\n") }},
		{name: "partial seal", mutate: func(t *testing.T, bench string) { writeFixtureFile(t, bench+".seal", `{"schema":`) }},
		{name: "digest mismatch", mutate: func(t *testing.T, bench string) {
			writeFixtureFile(t, bench, "#!/usr/bin/env bash\nexit 0\n")
			if err := os.Chmod(bench, 0o755); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "legacy binary"},
		{name: "always-green altered executable", mutate: func(t *testing.T, bench string) {
			writeFixtureFile(t, bench, "#!/usr/bin/env bash\nexit 0\n")
			if err := os.Chmod(bench, 0o755); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(failure.name, func(t *testing.T) {
			kit := filepath.Join(t.TempDir(), "kit [*] path")
			gatePath := writeGateEntryFixture(t, h, kit)
			bench := filepath.Join(kit, "dist", "bench")
			publishGateFixtureBench(t, kit, "#!/usr/bin/env bash\ncase \"$1\" in\nfreshness-check) exit 2 ;;\ngate-phases) touch \"$2/.git/gate-phases-ran\" ;;\nesac\n")
			if failure.mutate != nil {
				failure.mutate(t, bench)
			}

			probe := runAt(nested, "bash", gatePath)
			if probe == nil || probe.ExitCode == 0 {
				t.Fatalf("gate entry exit = %+v, want refusal", probe)
			}
			output := probe.Stdout + probe.Stderr
			rebuild := freshness.RebuildAction(kit)
			if !strings.Contains(output, rebuild) || strings.Count(output, "rebuild with ") != 1 {
				t.Fatalf("gate output = %q, want one stable rebuild action %q", output, rebuild)
			}
			if strings.Contains(output, "unknown subcommand") {
				t.Fatalf("gate leaked output from an unverified binary:\n%s", output)
			}
			if _, err := os.Stat(filepath.Join(root, ".git", "gate-phases-ran")); !os.IsNotExist(err) {
				t.Fatalf("gate entry scheduled phases after %s: %v", failure.name, err)
			}
		})
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
