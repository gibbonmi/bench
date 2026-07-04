package conformance

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
)

func TestLoadValidityMetadataFixturesBite(t *testing.T) {
	fixtures := []string{
		"invalid-json",
		"codex-hooks-broken",
		"bad-frontmatter",
		"claude-skills-unmirrored",
		"extensionless-gate-ref",
		"shared-rule-drift",
	}
	for _, fixture := range fixtures {
		t.Run(fixture, func(t *testing.T) {
			root := materializeConformanceFixture(t, fixture)
			h := NewHarness(t)
			expect := readExpectation(t, h.KitPath("tests", "canary", fixture, "EXPECT"))

			diags := RunConformance(root, h.KitRoot)

			if !containsDiagnostic(diags, expect) {
				t.Fatalf("%s did not bite under Go conformance; want %q in diagnostics:\n%s", fixture, expect, strings.Join(diags, "\n"))
			}
		})
	}
}

func TestSkillsIndexAndCommandAdapterFixturesBite(t *testing.T) {
	fixtures := []string{
		"dangling-index",
		"missing-index-field",
		"stale-index-wording",
		"unindexed-skill",
	}
	for _, fixture := range fixtures {
		t.Run(fixture, func(t *testing.T) {
			root := materializeConformanceFixture(t, fixture)
			h := NewHarness(t)
			expect := readExpectation(t, h.KitPath("tests", "canary", fixture, "EXPECT"))

			diags := RunConformance(root, h.KitRoot)

			if !containsDiagnostic(diags, expect) {
				t.Fatalf("%s did not bite under Go conformance; want %q in diagnostics:\n%s", fixture, expect, strings.Join(diags, "\n"))
			}
		})
	}
}

func TestDocsCurrencyTokenDietAndWorkflowFixturesBite(t *testing.T) {
	fixtures := []string{
		"stale-command-reference",
		"stale-codex-adapter-reference",
		"stale-cli-doc-reference",
		"historical-marker-prose",
		"benchref-missing",
		"benchref-pointer-dropped",
		"benchref-imported",
		"benchref-section-duplicated",
		"readme-command-first",
		"acceptance-coverage-anchor",
		"coverage-axis-anchor",
		"command-handoff-anchor",
		"debug-archaeology-anchor",
		"edge-inventory-anchor",
		"implement-spec-status-flip-anchor",
		"shape-idea-bypass",
		"shape-idea-bypass-wrapped",
		"shape-idea-handoff-anchor",
		"story-line-anchor-missing",
		"write-spec-handoff-anchor",
		"write-spec-map-required",
		"line-anchor-missing",
	}
	for _, fixture := range fixtures {
		t.Run(fixture, func(t *testing.T) {
			root := materializeConformanceFixture(t, fixture)
			h := NewHarness(t)
			expect := readExpectation(t, h.KitPath("tests", "canary", fixture, "EXPECT"))

			diags := RunConformance(root, h.KitRoot)

			if !containsDiagnostic(diags, expect) {
				t.Fatalf("%s did not bite under Go conformance; want %q in diagnostics:\n%s", fixture, expect, strings.Join(diags, "\n"))
			}
		})
	}
}

func TestCoverageMapValidationFixtureBite(t *testing.T) {
	fixture := "broken-coverage-map"
	root := materializeConformanceFixture(t, fixture)
	h := NewHarness(t)
	expect := readExpectation(t, h.KitPath("tests", "canary", fixture, "EXPECT"))

	diags := RunConformance(root, h.KitRoot)

	if !containsDiagnostic(diags, expect) {
		t.Fatalf("%s did not bite under Go conformance; want %q in diagnostics:\n%s", fixture, expect, strings.Join(diags, "\n"))
	}
}

func TestLineRoutingFixturesBite(t *testing.T) {
	fixtures := []string{
		"line-binding-prose-drift",
		"agent-hook-unwired",
		"agent-hook-broken",
		"adapter-line-broken",
	}
	for _, fixture := range fixtures {
		t.Run(fixture, func(t *testing.T) {
			root := materializeConformanceFixture(t, fixture)
			h := NewHarness(t)
			expect := readExpectation(t, h.KitPath("tests", "canary", fixture, "EXPECT"))

			diags := RunConformance(root, h.KitRoot)

			if !containsDiagnostic(diags, expect) {
				t.Fatalf("%s did not bite under Go conformance; want %q in diagnostics:\n%s", fixture, expect, strings.Join(diags, "\n"))
			}
		})
	}
}

func TestPackageCoreAndGuardFixturesBite(t *testing.T) {
	fixtures := []string{
		"missing-files-entry",
		"go-build-broken",
		"go-test-failing",
		"guard-describe-boundary-dropped",
	}
	for _, fixture := range fixtures {
		t.Run(fixture, func(t *testing.T) {
			root := materializeConformanceFixture(t, fixture)
			h := NewHarness(t)
			expect := readExpectation(t, h.KitPath("tests", "canary", fixture, "EXPECT"))

			diags := RunConformance(root, h.KitRoot)

			if !containsDiagnostic(diags, expect) {
				t.Fatalf("%s did not bite under Go conformance; want %q in diagnostics:\n%s", fixture, expect, strings.Join(diags, "\n"))
			}
		})
	}
}

func TestRunConformanceAcceptsHostileRootPath(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "root with spaces [glob]")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	h := NewHarness(t)
	if err := canary.MaterializeFixture(h.KitPath("tests", "canary", "invalid-json", "files"), root); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "init")

	diags := RunConformance(root, h.KitRoot)

	if !containsDiagnostic(diags, "invalid JSON in package.json") {
		t.Fatalf("hostile root path did not produce expected diagnostic:\n%s", strings.Join(diags, "\n"))
	}
}

func TestRunConformanceChecksExecutableGitMode(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")
	if err := os.MkdirAll(filepath.Join(root, "bin"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "bin", "bench.sh"), []byte("#!/usr/bin/env bash\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "add", "bin/bench.sh")

	diags := RunConformance(root, NewHarness(t).KitRoot)

	if !containsDiagnostic(diags, "bin/bench.sh is not executable in git") {
		t.Fatalf("non-executable tracked command path was not diagnosed:\n%s", strings.Join(diags, "\n"))
	}
}

func TestRunConformanceDistinguishesAbsentAndEmptyInputs(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")

	absent := RunConformance(root, NewHarness(t).KitRoot)
	if !containsDiagnostic(absent, "JSON file missing: package.json") {
		t.Fatalf("absent package.json diagnostic missing:\n%s", strings.Join(absent, "\n"))
	}
	if !containsDiagnostic(absent, "lines.env missing: .bench/lines.env") {
		t.Fatalf("absent lines.env diagnostic missing:\n%s", strings.Join(absent, "\n"))
	}

	if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".bench", "lines.env"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	empty := RunConformance(root, NewHarness(t).KitRoot)
	if !containsDiagnostic(empty, "lines.env tier unset: BENCH_TIER_TOP has no value") {
		t.Fatalf("empty lines.env diagnostic missing:\n%s", strings.Join(empty, "\n"))
	}
}

func TestConformanceSubprocessEnvStripsRootOverride(t *testing.T) {
	t.Setenv("BENCH_CONFORMANCE_ROOT", "/tmp/outer-root")

	for _, kv := range conformanceSubprocessEnv() {
		if strings.HasPrefix(kv, "BENCH_CONFORMANCE_ROOT=") {
			t.Fatalf("BENCH_CONFORMANCE_ROOT leaked into subprocess env: %q", kv)
		}
	}
}

func materializeConformanceFixture(t *testing.T, fixture string) string {
	t.Helper()
	h := NewHarness(t)
	root := t.TempDir()
	src := h.KitPath("tests", "canary", fixture, "files")
	if err := canary.MaterializeFixture(src, root); err != nil {
		t.Fatalf("materialize %s: %v", fixture, err)
	}
	cmd := exec.Command("git", "init", "-q")
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init fixture %s: %v\n%s", fixture, err, out)
	}
	return root
}

func readExpectation(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read EXPECT: %v", err)
	}
	return strings.TrimRight(string(data), "\n")
}
