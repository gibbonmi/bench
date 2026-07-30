package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/conformance/registry"
)

func TestRootConformance(t *testing.T) {
	root := os.Getenv("BENCH_CONFORMANCE_ROOT")
	if root == "" {
		capability.Environment(t, "BENCH_CONFORMANCE_ROOT not set")
	}
	h := NewHarness(t)
	// The scope env is passed through verbatim: any normalising here would let a stale
	// or misspelled value slide into a silent full run instead of the driver's red.
	for _, diag := range RunConformance(root, h.KitRoot, registry.TierFor(os.Getenv(registry.ConformanceTierEnv)), os.Getenv(registry.ConformanceCheckEnv)) {
		t.Errorf("gate: %s", diag)
	}
}

func TestGateEntryRunsGoConformanceAndBehaviorContracts(t *testing.T) {
	h := NewHarness(t)
	gate := h.ReadRootFile(".bench", "gate.sh")

	needle := `exec env BENCH_KIT="$kit" "$bench" gate-phases "$root"`
	if !strings.Contains(gate, needle) {
		t.Fatalf(".bench/gate.sh missing %q", needle)
	}
	freshness := `"$bench" freshness-check "$kit"`
	if strings.Index(gate, freshness) < 0 || strings.Index(gate, freshness) > strings.Index(gate, needle) {
		t.Fatalf(".bench/gate.sh does not verify %q before %q", freshness, needle)
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
			t.Fatalf(".bench/gate.sh still references retired conformance fragment %s", retired)
		}
	}
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
	rebuild := freshnessRebuildAction(t, h, kit, nested)
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
	bench := filepath.Join(kit, "dist", "bench")
	writeFixtureFile(t, bench, fmt.Sprintf("#!/usr/bin/env bash\ncase \"$1\" in\nfreshness-check) test \"$2\" = %q || exit 97; exit 2 ;;\ngate-phases) test \"$BENCH_KIT\" = %q || exit 98; printf 'old phase\\n'; printf 'old\\n' >> \"$2/.git/gate-phases-ran\" ;;\nesac\n", kit, kit))
	if err := os.Chmod(bench, 0o755); err != nil {
		t.Fatalf("chmod legacy bench: %v", err)
	}

	legacy := runAt(nested, "bash", gatePath)
	if legacy == nil || legacy.ExitCode == 0 {
		t.Fatalf("legacy gate entry exit = %+v, want refusal", legacy)
	}
	legacyOutput := legacy.Stdout + legacy.Stderr
	if rebuild := freshnessRebuildAction(t, h, kit, nested); !strings.Contains(legacyOutput, rebuild) {
		t.Fatalf("legacy gate entry missing freshness rebuild action %q:\n%s", rebuild, legacyOutput)
	}
	if strings.Contains(legacyOutput, "old phase") {
		t.Fatalf("legacy gate entry ran the old table:\n%s", legacyOutput)
	}
	if _, err := os.Stat(filepath.Join(root, ".git", "gate-phases-ran")); !os.IsNotExist(err) {
		t.Fatalf("legacy gate entry scheduled the old table: %v", err)
	}

	writeFixtureFile(t, bench, fmt.Sprintf("#!/usr/bin/env bash\ncase \"$1\" in\nfreshness-check) test \"$2\" = %q ;;\ngate-phases) test \"$BENCH_KIT\" = %q || exit 98; printf 'replacement phase\\n'; printf 'replacement\\n' >> \"$2/.git/gate-phases-ran\" ;;\nesac\n", kit, kit))
	if err := os.Chmod(bench, 0o755); err != nil {
		t.Fatalf("chmod replacement bench: %v", err)
	}
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
		name, program string
	}{
		{name: "missing executable"},
		{name: "missing seal", program: "printf 'untrusted missing seal\\n' >&2; exit 1"},
		{name: "unreadable seal", program: "printf 'untrusted unreadable seal\\n' >&2; exit 1"},
		{name: "malformed seal", program: "printf 'untrusted malformed seal\\n' >&2; exit 1"},
		{name: "partial seal", program: "printf 'untrusted partial seal\\n' >&2; exit 1"},
		{name: "digest mismatch", program: "printf 'untrusted digest mismatch\\n' >&2; exit 1"},
		{name: "legacy binary", program: "printf 'unknown subcommand freshness-check\\n' >&2; exit 2"},
	} {
		t.Run(failure.name, func(t *testing.T) {
			kit := filepath.Join(t.TempDir(), "kit [*] path")
			gatePath := writeGateEntryFixture(t, h, kit)
			if failure.program != "" {
				bench := filepath.Join(kit, "dist", "bench")
				writeFixtureFile(t, bench, "#!/usr/bin/env bash\ncase \"$1\" in\nfreshness-check) "+failure.program+" ;;\ngate-phases) touch \"$2/.git/gate-phases-ran\" ;;\nesac\n")
				if err := os.Chmod(bench, 0o755); err != nil {
					t.Fatalf("chmod fixture bench: %v", err)
				}
			}

			probe := runAt(nested, "bash", gatePath)
			if probe == nil || probe.ExitCode == 0 {
				t.Fatalf("gate entry exit = %+v, want refusal", probe)
			}
			output := probe.Stdout + probe.Stderr
			rebuild := freshnessRebuildAction(t, h, kit, nested)
			if !strings.Contains(output, rebuild) || strings.Count(output, "rebuild with ") != 1 {
				t.Fatalf("gate output = %q, want one stable rebuild action %q", output, rebuild)
			}
			if strings.Contains(output, "untrusted ") || strings.Contains(output, "unknown subcommand") {
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
	gate, err := os.ReadFile(h.KitPath(".bench", "gate.sh"))
	if err != nil {
		t.Fatalf("read gate entry: %v", err)
	}
	path := filepath.Join(kit, ".bench", "gate.sh")
	writeFixtureFile(t, path, string(gate))
	if err := os.Chmod(path, 0o755); err != nil {
		t.Fatalf("chmod gate entry: %v", err)
	}
	return path
}

func freshnessRebuildAction(t *testing.T, h Harness, source, cwd string) string {
	t.Helper()
	probe := runAt(cwd, h.KitPath("dist", "bench"), "freshness-check", source)
	if probe == nil || probe.ExitCode == 0 {
		t.Fatalf("freshness-check exit = %+v, want refusal for fixture root", probe)
	}
	output := probe.Stdout + probe.Stderr
	start := strings.Index(output, "rebuild with ")
	if start < 0 {
		t.Fatalf("freshness-check missing rebuild action:\n%s", output)
	}
	return strings.TrimSpace(output[start:])
}
