package conformance

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/capability"
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

			diags := RunConformance(root, h.KitRoot, registry.Dev, "")

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

			diags := RunConformance(root, h.KitRoot, registry.Dev, "")

			if !containsDiagnostic(diags, expect) {
				t.Fatalf("%s did not bite under Go conformance; want %q in diagnostics:\n%s", fixture, expect, strings.Join(diags, "\n"))
			}
		})
	}
}

func TestDocsCurrencyTokenDietAndWorkflowFixturesBite(t *testing.T) {
	h := NewHarness(t)
	all, err := canary.Fixtures(filepath.Join(h.KitRoot, "tests", "canary"))
	requireFixtureNoError(t, err)
	var fixtures []string
	for name, fixture := range all {
		if fixture.Family == "docs-currency-token-diet" || fixture.Family == "workflow-guidance-anchors" {
			fixtures = append(fixtures, name)
		}
	}
	slices.Sort(fixtures)
	for _, fixture := range fixtures {
		t.Run(fixture, func(t *testing.T) {
			root := materializeConformanceFixture(t, fixture)
			h := NewHarness(t)
			expect := readExpectation(t, filepath.Join(canaryFixturePath(t, h.KitRoot, fixture), "EXPECT"))

			diags := RunConformance(root, h.KitRoot, registry.Dev, "")

			if !containsDiagnostic(diags, expect) {
				t.Fatalf("%s did not bite under Go conformance; want %q in diagnostics:\n%s", fixture, expect, strings.Join(diags, "\n"))
			}
		})
	}
}

func TestSpecBuildCadenceAnchorsRejectDeletionSwapAndRawGitRouting(t *testing.T) {
	h := NewHarness(t)
	tests := []struct{ name, rel, old, replacement, diag string }{
		{"lifecycle deletion", ".agents/commands/bench-implement-spec.md", "`start` → `assign` → `checkpoint` →\n`integrate` → `review` → `promote`; `status` inspects the run and `abandon`\nplans or applies cleanup.", "", "bench-implement-spec dropped or reordered the eight-operation spec-build lifecycle"},
		{"initial capacity deletion", ".agents/commands/bench-implement-spec.md", "Re-derive the complete ready frontier and the harness's live capacity before\ndispatch. Assign every ownership-safe ticket up to the smaller of frontier size\nand available capacity.", "", "bench-implement-spec dropped initial frontier capacity dispatch"},
		{"additive generic unused-slot reason", ".agents/commands/bench-implement-spec.md", "For every unused harness slot, record exactly one\nreason: dependency, overlapping ownership fence, unavailable harness capacity,\nor measured resource constraint.", "For every unused harness slot, record exactly one reason: dependency, overlapping ownership fence, unavailable harness capacity, or measured resource constraint. An unused slot may instead be `NOT\n  PARALLELIZABLE`.", "bench-implement-spec permits a generic unused-slot reason outside the closed set"},
		{"exact candidate input deletion", ".agents/commands/bench-review-implementation.md", "For an active spec build, read `bench spec build status\n   <slug> --full` and bind the review inputs to the exact candidate subject and\n   recorded run base it reports. Confirm that subject is unchanged immediately\n   before receipt submission; a changed candidate invalidates the review rather\n   than letting a delta review authorize a new composition. ", "", "bench-review-implementation dropped exact-candidate review input binding"},
		{"frontier swap", ".agents/commands/bench-implement-spec.md", "or measured resource constraint. Refill the ownership-safe frontier after every\nintegration or assignment release while another delegate remains active.", "or measured resource constraint. Wait for the ownership-safe frontier to drain before refill after every integration or assignment release.", "bench-implement-spec replaced continuous frontier refill with drain-then-refill cadence"},
		{"repair deletion", ".agents/commands/bench-implement-spec.md", "Accepted findings become new ownership-fenced repair tickets and re-enter\n  `assign`, `checkpoint`, and `integrate` before a fresh composed review.", "", "bench-implement-spec routes an accepted repair outside the provisional lifecycle"},
		{"recomposition purpose deletion", ".bench/BENCH-reference.md", "; a moved tip recomposes through `promote`, discarding the review |", " |", "BENCH-reference dropped promote's moved-tip recomposition from the lifecycle lookup"},
		{"recomposition discard deletion", ".agents/commands/bench-implement-spec.md", "When the branch tip moves, `promote` is the operation that recomposes the run\n  onto the new tip, and recomposition discards the review.\n  ", "", "bench-implement-spec dropped moved-tip recomposition through promote or its review discard"},
		{"repair round deletion", ".agents/commands/bench-implement-spec.md", "\n  The repair round is therefore\n  repair → `promote` → `review` → `assign` … `integrate` → `review` → `promote`.", "", "bench-implement-spec dropped the ordered moved-tip repair round"},
		{"probe kind", ".agents/skills/bench-craft-delegate/SKILL.md", "The\ncoordinator probe's mutation kind differs from the delegate author's mutation\nkind.", "The\ncoordinator probe's mutation kind matches the delegate author's mutation\nkind.", "craft-delegate allows the coordinator probe to repeat the author's mutation kind"},
		{"ordinary commit route", ".bench/BENCH.md", "Provisional\ncadence is exclusive to reviewed spec-backed builds; light-path work, `bench\nshift`, and ordinary `bench commit` remain commit-on-green.", "Provisional cadence covers reviewed spec-backed builds, light-path work, `bench shift`, and ordinary `bench commit`.", ".bench/BENCH.md broadened provisional cadence beyond reviewed spec-backed builds"},
		{"purpose swap", ".bench/BENCH-reference.md", "| `assign` | lease one ownership-fenced ticket worktree |", "| `assign` | validate focused evidence and bind a provisional commit |", "BENCH-reference misroutes spec build assign"},
		{"flag positional", "bin/bench.sh", "bench spec build assign <slug> --ticket <ticket> --request <id>", "bench spec build assign <slug> <ticket> --ticket --request <id>", "bench help dropped or malformed spec build assign grammar"},
		{"line replacement", "projects/benchkit.md", "Spec-build guidance cadence** → **`gpt-5.6-sol / high`", "Spec-build guidance cadence** → **`gpt-5.6-terra / high`", "benchkit profile replaced the approved spec-build guidance line"},
		{"control deletion", "CHANGELOG.md", "Light-path changes, `bench shift`, and ordinary `bench commit` keep\n  commit-on-green cadence.", "", "CHANGELOG dropped the unchanged-path control for provisional spec builds"},
		{"raw git route", ".agents/commands/bench-implement-spec.md", "Submit focused delegate evidence plus the coordinator-owned, different-kind\n  probe through `checkpoint`.", "Create the checkpoint with `git commit`; the public `checkpoint` token remains documented.", "bench-implement-spec synthesizes lifecycle Git plumbing outside the eight public operations"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(h.KitRoot, filepath.FromSlash(tc.rel)))
			if err != nil || strings.Count(string(data), tc.old) != 1 {
				t.Fatalf("mutation anchor count for %s = %d, %v", tc.rel, strings.Count(string(data), tc.old), err)
			}
			root := t.TempDir()
			path := filepath.Join(root, filepath.FromSlash(tc.rel))
			requireFixtureNoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
			requireFixtureNoError(t, os.WriteFile(path, []byte(strings.Replace(string(data), tc.old, tc.replacement, 1)), 0o644))
			if diags := checkWorkflowAnchors(root); !containsDiagnostic(diags, tc.diag) {
				t.Fatalf("mutation did not bite with %q:\n%s", tc.diag, strings.Join(diags, "\n"))
			}
		})
	}
	t.Run("additive direct working branch permission", func(t *testing.T) {
		const rel = ".agents/commands/bench-implement-spec.md"
		const anchor = "Accepted findings become new ownership-fenced repair tickets and re-enter\n  `assign`, `checkpoint`, and `integrate` before a fresh composed review."
		const diag = "bench-implement-spec permits an accepted repair to bypass provisional assignment and write directly to the working branch"
		data, err := os.ReadFile(filepath.Join(h.KitRoot, filepath.FromSlash(rel)))
		requireFixtureNoError(t, err)
		root := t.TempDir()
		path := filepath.Join(root, filepath.FromSlash(rel))
		requireFixtureNoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		if strings.Count(string(data), anchor) != 1 {
			t.Fatalf("accepted-repair paragraph anchor count = %d", strings.Count(string(data), anchor))
		}
		for _, contradiction := range []string{
			"For an accepted repair finding, the coordinator may instead write the repair directly to the working branch before `promote`.",
			"For an accepted `repair` finding, the coordinator may instead write the repair directly to the\n  `working branch` before `promote`.",
		} {
			mutated := strings.Replace(string(data), anchor, anchor+"\n  "+contradiction, 1)
			requireFixtureNoError(t, os.WriteFile(path, []byte(mutated), 0o644))
			if diags := checkWorkflowAnchors(root); !containsDiagnostic(diags, diag) {
				t.Fatalf("additive contradiction did not bite with %q:\n%s", diag, strings.Join(diags, "\n"))
			}
			requireFixtureNoError(t, os.WriteFile(path, data, 0o644))
			if diags := checkWorkflowAnchors(root); containsDiagnostic(diags, diag) {
				t.Fatalf("additive contradiction remained red after removal:\n%s", strings.Join(diags, "\n"))
			}
		}
	})
}

func requireFixtureNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCoverageMapValidationFixtureBite(t *testing.T) {
	fixtures := []string{
		"broken-coverage-map",
		"no-map-not-historical",
		"stray-flat-live-spec",
	}
	for _, fixture := range fixtures {
		t.Run(fixture, func(t *testing.T) {
			root := materializeConformanceFixture(t, fixture)
			h := NewHarness(t)
			expect := readExpectation(t, filepath.Join(canaryFixturePath(t, h.KitRoot, fixture), "EXPECT"))

			diags := RunConformance(root, h.KitRoot, registry.Dev, "")

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

			diags := RunConformance(root, h.KitRoot, registry.Dev, "")

			if !containsDiagnostic(diags, expect) {
				t.Fatalf("%s did not bite under Go conformance; want %q in diagnostics:\n%s", fixture, expect, strings.Join(diags, "\n"))
			}
		})
	}
}

func TestPackageCoreAndGuardFixturesBite(t *testing.T) {
	// Only fixtures a conformance check grades belong here. A fixture whose failure a
	// gate phase owns is proved by the canary sweep at that phase instead; running
	// conformance over its tree compiles and tests nothing, so it would report
	// did-not-bite forever.
	fixtures := []string{
		"missing-files-entry",
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

			diags := RunConformance(root, h.KitRoot, registry.Dev, "")

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
	for _, family := range registry.Families() {
		if family == "coverage-map-validation" {
			continue
		}
		familyDir := filepath.Join(kitRoot, "tests", "canary", family)
		requireFixtureNoError(t, os.MkdirAll(familyDir, 0o755))
		requireFixtureNoError(t, os.MkdirAll(filepath.Join(familyDir, "sentinel"), 0o755))
	}

	diags := RunConformance(root, kitRoot, registry.Dev, "")

	want := `canary conformance family "coverage-map-validation" has no fixture directories under tests/canary/coverage-map-validation`
	if !containsDiagnostic(diags, want) {
		t.Fatalf("absent canary family did not produce diagnostic %q:\n%s", want, strings.Join(diags, "\n"))
	}
}

func TestRunConformanceReportsEmptyCanaryFamily(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")
	kitRoot := t.TempDir()
	for _, family := range registry.Families() {
		familyDir := filepath.Join(kitRoot, "tests", "canary", family)
		requireFixtureNoError(t, os.MkdirAll(familyDir, 0o755))
		if family == "coverage-map-validation" {
			continue
		}
		requireFixtureNoError(t, os.MkdirAll(filepath.Join(familyDir, "sentinel"), 0o755))
	}

	diags := RunConformance(root, kitRoot, registry.Dev, "")

	want := `canary conformance family "coverage-map-validation" has no fixture directories under tests/canary/coverage-map-validation`
	if !containsDiagnostic(diags, want) {
		t.Fatalf("empty canary family did not produce diagnostic %q:\n%s", want, strings.Join(diags, "\n"))
	}
}

// TestRunConformanceReportsUnboundCanaryFamily grades the direction the derived
// family list cannot see: a family directory on disk that the table does not bind.
// Its fixtures would each silently run a full inner gate — the cost the scoping
// exists to remove — so the kit's tree and its table have to agree in both
// directions. The behavior-owned directory and the legacy flat fixture beside it are
// not conformance families and must stay unreported.
func TestRunConformanceReportsUnboundCanaryFamily(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")
	kitRoot := t.TempDir()
	canaryDir := filepath.Join(kitRoot, "tests", "canary")
	families := append(registry.Families(), "unbound-family", "behavior-owned")
	for _, family := range families {
		requireFixtureNoError(t, os.MkdirAll(filepath.Join(canaryDir, family, "sentinel"), 0o755))
	}
	flat := filepath.Join(canaryDir, "legacy-flat")
	requireFixtureNoError(t, os.MkdirAll(filepath.Join(flat, "files"), 0o755))
	requireFixtureNoError(t, os.WriteFile(filepath.Join(flat, "EXPECT"), []byte("target\n"), 0o644))

	diags := RunConformance(root, kitRoot, registry.Dev, "")

	want := `canary conformance family "unbound-family" is bound to no conformance check; add it to the registry family table so its fixtures run scoped`
	if !containsDiagnostic(diags, want) {
		t.Fatalf("unbound canary family did not produce diagnostic %q:\n%s", want, strings.Join(diags, "\n"))
	}
	joined := strings.Join(diags, "\n")
	for _, excluded := range []string{"behavior-owned", "legacy-flat"} {
		if strings.Contains(joined, excluded) {
			t.Errorf("%q is not a conformance family but was reported:\n%s", excluded, joined)
		}
	}
}

// TestSymlinkedCanaryFamilyIsInvisibleToTreeAndSweep pins the agreement that makes
// skipping a symlinked family directory the right answer rather than a hole. os.ReadDir
// reports a symlink by its own type, so neither the unbound-family read nor the canary
// package's fixture walk descends into one — a symlinked family therefore contributes
// no fixture to the sweep, has no inner run to scope, and cannot be unbound. Reporting
// it would demand a table binding no fixture ever uses. The two sides share one reading
// of the tree; changing either alone reds a family with no fixtures or leaves a real
// family's fixtures silently unscoped.
func TestSymlinkedCanaryFamilyIsInvisibleToTreeAndSweep(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")
	kitRoot := t.TempDir()
	canaryDir := filepath.Join(kitRoot, "tests", "canary")
	for _, family := range registry.Families() {
		writeCanaryFixture(t, filepath.Join(canaryDir, family, family+"-fx"))
	}
	// The target sits outside tests/canary, so only the link can make it read as a
	// family — and it holds a real fixture, so a walk that followed the link would both
	// report the family unbound and sweep the fixture under it.
	target := filepath.Join(kitRoot, "outside", "linked-family")
	writeCanaryFixture(t, filepath.Join(target, "linked-fixture"))
	if err := os.Symlink(target, filepath.Join(canaryDir, "symlinked-family")); err != nil {
		capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable on this filesystem: %v", err))
	}

	diags := RunConformance(root, kitRoot, registry.Dev, "")
	if joined := strings.Join(diags, "\n"); strings.Contains(joined, "symlinked-family") {
		t.Errorf("a symlinked family contributes no fixtures but was reported:\n%s", joined)
	}

	var mu sync.Mutex
	var swept []string
	err := canary.Sweep(kitRoot, func(call canary.RunCall) canary.RunResult {
		if call.FixtureDir == "" {
			return canary.RunResult{ExitCode: 1, Output: "baseline noise\n"}
		}
		mu.Lock()
		swept = append(swept, filepath.Base(call.FixtureDir))
		mu.Unlock()
		return canary.RunResult{ExitCode: 1, Output: "target-" + filepath.Base(call.FixtureDir) + "\n"}
	})
	if err != nil {
		t.Fatalf("Sweep err = %v", err)
	}
	if slices.Contains(swept, "linked-fixture") {
		t.Errorf("the sweep graded %v, which reaches through the symlinked family the tree check skips", swept)
	}
}

// TestRunConformanceReportsEveryFamilyWhenCanaryTreeIsUnreadable pins the direction the
// unbound-family read cannot cover. It returns nothing when tests/canary will not open,
// which is safe only because the family-presence loop iterates the registry table and
// reports every family it cannot find fixtures for — an absent or unreadable tree is
// therefore the loudest red the check has, not a silent skip.
func TestRunConformanceReportsEveryFamilyWhenCanaryTreeIsUnreadable(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T, canaryDir string)
	}{
		{"absent", func(*testing.T, string) {}},
		{"unreadable", func(t *testing.T, canaryDir string) {
			for _, family := range registry.Families() {
				writeCanaryFixture(t, filepath.Join(canaryDir, family, family+"-fx"))
			}
			// The restore is registered before the strip so it runs ahead of TempDir's
			// own removal, which cannot descend into a directory it cannot enter.
			t.Cleanup(func() { _ = os.Chmod(canaryDir, 0o700) })
			if err := os.Chmod(canaryDir, 0o000); err != nil {
				capability.Capability(t, capability.Privilege, fmt.Sprintf("cannot strip directory permissions: %v", err))
			}
			if _, err := os.ReadDir(canaryDir); err == nil {
				capability.Capability(t, capability.Privilege, "mode 0o000 directory is still readable by this user")
			}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			runGit(t, root, "init")
			kitRoot := t.TempDir()
			tt.setup(t, filepath.Join(kitRoot, "tests", "canary"))

			diags := RunConformance(root, kitRoot, registry.Dev, "")

			for _, family := range registry.Families() {
				want := fmt.Sprintf("canary conformance family %q has no fixture directories", family)
				if !containsDiagnostic(diags, want) {
					t.Errorf("no diagnostic %q:\n%s", want, strings.Join(diags, "\n"))
				}
			}
		})
	}
}

// writeCanaryFixture creates the minimum a canary fixture needs to be swept: a files/
// tree and an EXPECT the sweep helpers echo back.
func writeCanaryFixture(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, "files"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "EXPECT"), []byte("target-"+filepath.Base(dir)+"\n"), 0o644); err != nil {
		t.Fatal(err)
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

	diags := RunConformance(root, h.KitRoot, registry.Dev, "")

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

	diags := RunConformance(root, NewHarness(t).KitRoot, registry.Dev, "")

	if !containsDiagnostic(diags, "bin/bench.sh is not executable in git") {
		t.Fatalf("non-executable tracked command path was not diagnosed:\n%s", strings.Join(diags, "\n"))
	}
}

func TestConformanceSubprocessEnvStripsConformanceControlVars(t *testing.T) {
	t.Setenv("BENCH_CONFORMANCE_ROOT", "/tmp/outer-root")
	t.Setenv(registry.ConformanceTierEnv, "ship")
	t.Setenv(registry.ConformanceCheckEnv, "package-core-guard")

	for _, kv := range conformanceSubprocessEnv() {
		for _, name := range []string{"BENCH_CONFORMANCE_ROOT", registry.ConformanceTierEnv, registry.ConformanceCheckEnv} {
			if strings.HasPrefix(kv, name+"=") {
				t.Fatalf("%s leaked into the probe subprocess env: %q", name, kv)
			}
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
	// npm's update notifier writes stderr chatter; the pack JSON must survive while the stub replays both streams.
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
	fixturePath := canaryFixturePath(t, h.KitRoot, fixture)
	if err := canary.MaterializeMutationFixture(h.KitRoot, fixturePath, root); err != nil {
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
	found, ok := canaryFixturePaths(t, filepath.Join(kitRoot, "tests", "canary"))[fixture]
	if !ok {
		t.Fatalf("canary fixture %q not found", fixture)
	}
	return found.Dir
}

func readExpectation(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read EXPECT: %v", err)
	}
	return strings.TrimRight(string(data), "\n")
}
