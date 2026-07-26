package conformance

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/conformance/registry"
)

func TestLoadValidityMetadataFixturesBite(t *testing.T) {
	fixtures := []string{
		"invalid-json",
		"codex-hooks-broken",
		"codex-hooks-timeout",
		"codex-hooks-timeout-typed",
		"bad-frontmatter",
		"claude-skills-unmirrored",
		"extensionless-gate-ref",
		"gate-input-gitignored",
		"shared-rule-drift",
		"readme-shared-rule-drift",
	}
	for _, fixture := range fixtures {
		t.Run(fixture, func(t *testing.T) {
			root := materializeConformanceFixture(t, fixture)
			h := NewHarness(t)
			expect := readExpectation(t, filepath.Join(canaryFixturePath(t, h.KitRoot, fixture), "EXPECT"))

			diags := RunConformance(root, h.KitRoot, registry.Dev)

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
		"roadmap-promotion-persistence",
	}
	for _, fixture := range fixtures {
		t.Run(fixture, func(t *testing.T) {
			root := materializeConformanceFixture(t, fixture)
			h := NewHarness(t)
			expect := readExpectation(t, filepath.Join(canaryFixturePath(t, h.KitRoot, fixture), "EXPECT"))

			diags := RunConformance(root, h.KitRoot, registry.Dev)

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
		"retired-command-reference",
		"stale-cli-doc-reference",
		"stale-skill-cli-reference",
		"missing-cli-inventory",
		"historical-marker-prose",
		"benchref-missing",
		"benchref-pointer-dropped",
		"benchref-imported",
		"benchref-section-duplicated",
		"dogfood-referent-shipped",
		"readme-command-first",
		"signal-vocabulary-drift",
		"structured-phase-progress-anchor",
		"acceptance-coverage-anchor",
		"coverage-axis-anchor",
		"command-handoff-anchor",
		"debug-archaeology-anchor",
		"debug-red-commit",
		"readme-shaping-skip",
		"implement-spec-inline-exception",
		"implement-spec-landing-commit",
		"edge-inventory-anchor",
		"fix-pass-sentinel-anchor",
		"implement-spec-mandatory-delegation-anchor",
		"implement-spec-status-flip-anchor",
		"implement-spec-structure-pointer",
		"review-persistence-anchor",
		"shared-worktree-path-pin",
		"delegate-parallel-route-anchor",
		"delegate-stash-refusal-anchor",
		"shape-idea-bypass",
		"shape-idea-bypass-wrapped",
		"shape-idea-handoff-anchor",
		"shape-idea-grill-continuation",
		"what-next-anchor",
		"what-next-spec-history-anchor",
		"what-next-roadmap-context-anchor",
		"spec-retire-roadmap-row",
		"staged-command-sweep-anchor",
		"capture-sink-anchor",
		"craft-seams-structure-headroom",
		"story-line-anchor-missing",
		"write-spec-handoff-anchor",
		"write-spec-map-required",
		"write-spec-reviewer-closed-comment-spoof",
		"write-spec-reviewer-closed-fast-path",
		"write-spec-review-trigger-dropped",
		"write-spec-review-tier-escalated",
		"write-spec-review-made-conditional",
		"write-spec-open-fork-fallback",
		"shape-idea-write-spec-entry-contract-pointer",
		"line-anchor-missing",
	}
	for _, fixture := range fixtures {
		t.Run(fixture, func(t *testing.T) {
			root := materializeConformanceFixture(t, fixture)
			h := NewHarness(t)
			expect := readExpectation(t, filepath.Join(canaryFixturePath(t, h.KitRoot, fixture), "EXPECT"))

			diags := RunConformance(root, h.KitRoot, registry.Dev)

			if !containsDiagnostic(diags, expect) {
				t.Fatalf("%s did not bite under Go conformance; want %q in diagnostics:\n%s", fixture, expect, strings.Join(diags, "\n"))
			}
		})
	}
}

func TestCoverageMapValidationFixtureBite(t *testing.T) {
	fixtures := []string{
		"broken-coverage-map",
		"no-map-not-historical",
	}
	for _, fixture := range fixtures {
		t.Run(fixture, func(t *testing.T) {
			root := materializeConformanceFixture(t, fixture)
			h := NewHarness(t)
			expect := readExpectation(t, filepath.Join(canaryFixturePath(t, h.KitRoot, fixture), "EXPECT"))

			diags := RunConformance(root, h.KitRoot, registry.Dev)

			if !containsDiagnostic(diags, expect) {
				t.Fatalf("%s did not bite under Go conformance; want %q in diagnostics:\n%s", fixture, expect, strings.Join(diags, "\n"))
			}
		})
	}
}

func TestLineRoutingFixturesBite(t *testing.T) {
	fixtures := []string{
		"line-binding-prose-drift",
		"agent-hook-unwired",
		"agent-hook-broken",
		"stop-hook-unwired",
		"adapter-line-broken",
	}
	for _, fixture := range fixtures {
		t.Run(fixture, func(t *testing.T) {
			root := materializeConformanceFixture(t, fixture)
			h := NewHarness(t)
			expect := readExpectation(t, filepath.Join(canaryFixturePath(t, h.KitRoot, fixture), "EXPECT"))

			diags := RunConformance(root, h.KitRoot, registry.Dev)

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
		"guard-resolver-order-drift",
		"default-branch-refabricated",
		"kit-only-asset-admitted",
		"kit-only-allowlist-emptied",
	}
	for _, fixture := range fixtures {
		t.Run(fixture, func(t *testing.T) {
			root := materializeConformanceFixture(t, fixture)
			h := NewHarness(t)
			expect := readExpectation(t, filepath.Join(canaryFixturePath(t, h.KitRoot, fixture), "EXPECT"))

			diags := RunConformance(root, h.KitRoot, registry.Dev)

			if !containsDiagnostic(diags, expect) {
				t.Fatalf("%s did not bite under Go conformance; want %q in diagnostics:\n%s", fixture, expect, strings.Join(diags, "\n"))
			}
		})
	}
}

func TestRunConformanceReportsAbsentCanaryFamily(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")
	kitRoot := t.TempDir()
	for _, family := range conformanceFamilies {
		if family == "coverage-map-validation" {
			continue
		}
		familyDir := filepath.Join(kitRoot, "tests", "canary", family)
		if err := os.MkdirAll(familyDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Join(familyDir, "sentinel"), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	diags := RunConformance(root, kitRoot, registry.Dev)

	want := `canary conformance family "coverage-map-validation" has no fixture directories under tests/canary/coverage-map-validation`
	if !containsDiagnostic(diags, want) {
		t.Fatalf("absent canary family did not produce diagnostic %q:\n%s", want, strings.Join(diags, "\n"))
	}
}

func TestRunConformanceReportsEmptyCanaryFamily(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")
	kitRoot := t.TempDir()
	for _, family := range conformanceFamilies {
		familyDir := filepath.Join(kitRoot, "tests", "canary", family)
		if err := os.MkdirAll(familyDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if family == "coverage-map-validation" {
			continue
		}
		if err := os.MkdirAll(filepath.Join(familyDir, "sentinel"), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	diags := RunConformance(root, kitRoot, registry.Dev)

	want := `canary conformance family "coverage-map-validation" has no fixture directories under tests/canary/coverage-map-validation`
	if !containsDiagnostic(diags, want) {
		t.Fatalf("empty canary family did not produce diagnostic %q:\n%s", want, strings.Join(diags, "\n"))
	}
}

func TestRunConformanceAcceptsHostileRootPath(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "root with spaces [glob]")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	h := NewHarness(t)
	if err := canary.MaterializeFixture(filepath.Join(canaryFixturePath(t, h.KitRoot, "invalid-json"), "files"), root); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "init")

	diags := RunConformance(root, h.KitRoot, registry.Dev)

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

	diags := RunConformance(root, NewHarness(t).KitRoot, registry.Dev)

	if !containsDiagnostic(diags, "bin/bench.sh is not executable in git") {
		t.Fatalf("non-executable tracked command path was not diagnosed:\n%s", strings.Join(diags, "\n"))
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

func TestConformanceSubprocessEnvProvidesWritableNpmCache(t *testing.T) {
	oldCache, hadCache := os.LookupEnv("NPM_CONFIG_CACHE")
	if err := os.Unsetenv("NPM_CONFIG_CACHE"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if hadCache {
			_ = os.Setenv("NPM_CONFIG_CACHE", oldCache)
			return
		}
		_ = os.Unsetenv("NPM_CONFIG_CACHE")
	})

	var cache string
	for _, kv := range conformanceSubprocessEnv() {
		if strings.HasPrefix(kv, "NPM_CONFIG_CACHE=") {
			cache = strings.TrimPrefix(kv, "NPM_CONFIG_CACHE=")
			break
		}
	}
	if cache == "" {
		t.Fatal("NPM_CONFIG_CACHE missing from conformance subprocess env")
	}
	if !strings.HasPrefix(filepath.Clean(cache), filepath.Clean(os.TempDir())+string(os.PathSeparator)) {
		t.Fatalf("NPM_CONFIG_CACHE = %q, want temp-backed cache", cache)
	}
}

func TestCheckPackageFilesToleratesNpmStderrNotice(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "bin"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "bin", "bench.sh"), []byte("#!/usr/bin/env bash\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "package.json"), []byte(`{"files":["bin/bench.sh"]}`), 0o644); err != nil {
		t.Fatal(err)
	}

	// npm's update notifier intermittently writes "npm notice ..." to stderr;
	// the pack JSON on stdout must survive that chatter, so the stub replays
	// both streams.
	stub := t.TempDir()
	script := "#!/usr/bin/env bash\n" +
		"printf '[{\"files\":[{\"path\":\"bin/bench.sh\"}]}]\\n'\n" +
		"echo 'npm notice New major version of npm available!' >&2\n"
	if err := os.WriteFile(filepath.Join(stub, "npm"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", stub+string(os.PathListSeparator)+os.Getenv("PATH"))

	for _, diag := range checkPackageFiles(root) {
		if strings.Contains(diag, "JSON unreadable") {
			t.Fatalf("npm stderr notice corrupted the pack JSON parse: %s", diag)
		}
	}
}

func materializeConformanceFixture(t *testing.T, fixture string) string {
	t.Helper()
	h := NewHarness(t)
	root := t.TempDir()
	src := filepath.Join(canaryFixturePath(t, h.KitRoot, fixture), "files")
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

func canaryFixturePath(t *testing.T, kitRoot, fixture string) string {
	t.Helper()
	path, ok := canaryFixturePaths(t, filepath.Join(kitRoot, "tests", "canary"))[fixture]
	if !ok {
		t.Fatalf("canary fixture %q not found", fixture)
	}
	return path
}

func readExpectation(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read EXPECT: %v", err)
	}
	return strings.TrimRight(string(data), "\n")
}
