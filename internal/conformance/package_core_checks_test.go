package conformance

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	kitpayload "github.com/gibbonmi/bench"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/guards"
	"github.com/gibbonmi/bench/internal/packagesurface"
)

func checkPackageCoreAndGuards(root string, tier registry.Tier) []string {
	var diags []string
	diags = append(diags, checkPackageFiles(root, tier)...)
	diags = append(diags, checkGoToolchain(root)...)
	diags = append(diags, checkReleaseWorkflow(root)...)
	diags = append(diags, checkNativeRuntimeWorkflow(root)...)
	diags = append(diags, checkReleasePreflight(root)...)
	diags = append(diags, checkShippedIdentityStrings(root)...)
	diags = append(diags, checkUserFacingBenchkitStrings(root)...)
	diags = append(diags, checkGuardHeaderManifests(root)...)
	diags = append(diags, checkGuardResolverOrderDrift(root)...)
	return diags
}

func checkPackageFiles(root string, tier registry.Tier) []string {
	path := filepath.Join(root, "package.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var pkg struct {
		Files []string `json:"files"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil
	}
	var diags []string
	for _, file := range pkg.Files {
		if strings.HasPrefix(file, "!") {
			continue
		}
		if exists(filepath.Join(root, filepath.FromSlash(file))) {
			continue
		}
		// A gitignored files[] entry (dist/, the repo-local Go build) exists only once
		// something has built it — a bare or prospective checkout never carries one, so
		// its completeness is a built-payload concern below ship tier, not a tree-shape
		// one. Ship-tier grading keeps the strict check: the release rehearsal build
		// precedes it.
		if tier != registry.Ship && benchgit.OK("-C", root, "check-ignore", "-q", "--", file) {
			continue
		}
		diags = append(diags, "package.json files[] missing "+file)
	}
	if len(pkg.Files) == 0 {
		return diags
	}

	// --ignore-scripts: inspect files[] membership only. The prepare build (npx-from-git
	// enablement) is a lifecycle side effect the git-install probe exercises for real;
	// running it here would rebuild dist/bench and defeat the built/unbuilt determinism
	// this shape check is meant to hold.
	probe := runAtCleanEnv(root, "npm", "pack", "--dry-run", "--json", "--ignore-scripts")
	if probe != nil && probe.ExitCode != 0 {
		diags = append(diags, formatProbeFailure("npm pack --dry-run failed", probe, root))
	} else if probe != nil {
		buildAssets, err := packagesurface.RequiredBuildPackAssets(root)
		if err != nil {
			diags = append(diags, "npm package build inputs are unreadable")
		} else {
			required := append([]string{}, packagesurface.RequiredPackAssets...)
			required = append(required, buildAssets...)
			diags = append(diags, checkNpmPackAssets(probe.Stdout, required)...)
		}
		// Independent of whether the build-input inventory resolved: the kit-only
		// guard only needs the packed file list and the allowlist, both already in
		// hand, so a fixture too minimal to carry cmd/ and internal/ still exercises it.
		diags = append(diags, checkNoKitOnlyPackedAssets(root, probe.Stdout)...)
	}
	return append(diags, checkRepoOnlyPackageClaims(root)...)
}

// TestCheckPackageFilesExemptsGitignoredEntryBelowShipTier pins the prospective-checkout
// repro: a files[] entry that is both absent and gitignored (dist/, in the real tree) is
// exactly what any bare `git worktree add` + `read-tree` checkout produces, since a build
// artifact never rides along with tracked content. Dev tier must not flag it — it is not a
// tree-shape defect — while ship tier, the release rehearsal, keeps the strict check. A
// genuinely missing tracked entry still flags at both tiers, so the exemption stays narrow.
func TestCheckPackageFilesExemptsGitignoredEntryBelowShipTier(t *testing.T) {
	root := t.TempDir()
	if out, err := exec.Command("git", "-C", root, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte("dist/\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "package.json"), []byte(`{"files":["dist/","README.md"]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	// README.md is deliberately left absent and untracked (not gitignored): a genuinely
	// missing tracked entry, the case the check must keep catching at every tier.

	dev := checkPackageFiles(root, registry.Dev)
	if containsDiagnostic(dev, "missing dist/") {
		t.Fatalf("dev tier flagged the gitignored, unbuilt dist/ entry: %v", dev)
	}
	if !containsDiagnostic(dev, "missing README.md") {
		t.Fatalf("dev tier did not flag the genuinely missing tracked entry: %v", dev)
	}

	ship := checkPackageFiles(root, registry.Ship)
	if !containsDiagnostic(ship, "missing dist/") {
		t.Fatalf("ship tier did not flag the missing built entry: %v", ship)
	}
	if !containsDiagnostic(ship, "missing README.md") {
		t.Fatalf("ship tier did not flag the genuinely missing tracked entry: %v", ship)
	}
}

// checkNoKitOnlyPackedAssets is the FT85 story 3 forbidden-asset guard: it derives the
// kit-only prefix set from the same canonical allowlist buildLinkPlan and
// build-release-evidence.mjs read, and grades the real npm pack --dry-run output
// against it, so a kit-only path readmitted anywhere in package.json's files[] (not
// just the wholesale .agents/ case the allowlist itself replaced) still turns the gate
// red. It reads .bench/consumer-payload.json from the graded root, not the running
// binary's own copy, so a canary fixture's mutated allowlist is what gets graded. An
// allowlist present but declaring no kit-only rows is itself a diagnostic: with an
// empty prefix set every packed path passes vacuously, so the emptied-allowlist case
// must be caught here rather than silently skipped alongside the missing-file case.
func checkNoKitOnlyPackedAssets(root, packJSON string) []string {
	rows, absent, err := kitpayload.PayloadRowsAt(filepath.Join(root, ".bench", "consumer-payload.json"))
	if absent {
		return nil
	}
	if err != nil {
		// A present allowlist that does not resolve cannot be silently skipped the way
		// an absent one is: the prefix set would be empty and every packed path would
		// pass vacuously, which is the exact failure this guard exists to catch.
		return []string{err.Error()}
	}
	prefixes := kitpayload.PayloadKitOnlyPrefixes(rows)
	var diags []string
	if len(prefixes) == 0 {
		diags = append(diags, "consumer payload allowlist declares no kit-only rows; the packed-asset guard would pass vacuously")
	}
	var packs []struct {
		Files []struct {
			Path string `json:"path"`
		} `json:"files"`
	}
	if err := json.Unmarshal([]byte(packJSON), &packs); err != nil || len(packs) == 0 {
		return diags
	}
	for _, file := range packs[0].Files {
		rel := strings.TrimPrefix(file.Path, "package/")
		if kitpayload.PayloadExcluded(rel, prefixes) {
			diags = append(diags, "npm package includes kit-only allowlist asset "+rel)
		}
	}
	return diags
}

func checkNpmPackAssets(packJSON string, required []string) []string {
	var packs []struct {
		Files []struct {
			Path string `json:"path"`
		} `json:"files"`
	}
	if err := json.Unmarshal([]byte(packJSON), &packs); err != nil {
		return []string{"npm pack --dry-run JSON unreadable: " + err.Error()}
	}
	files := map[string]bool{}
	if len(packs) > 0 {
		for _, file := range packs[0].Files {
			files[file.Path] = true
		}
	}
	var diags []string
	for _, asset := range required {
		if !files[asset] {
			diags = append(diags, "npm package missing "+asset)
		}
	}
	for _, forbidden := range packagesurface.ForbiddenPackAssets {
		if files[forbidden] {
			diags = append(diags, "npm package includes local-only file "+forbidden)
		}
	}
	return diags
}

func checkRepoOnlyPackageClaims(root string) []string {
	// Mirrors the package fragment's lightweight prose sweep over shipped markdown.
	var diags []string
	files := packageMarkdownFiles(root)
	claimRe := regexp.MustCompile(`(?i)\b(ship|ships|shipped|shipping|package|packaged|tarball|installable|included|includes)\b`)
	headingRe := regexp.MustCompile(`(?i)^#{1,6}\s+.*\b(ship|ships|shipped|shipping|package|packaged|tarball|installable|included|includes|surfaces?)\b`)
	repoOnlyRe := regexp.MustCompile(`(?i)\b(repo-only|development context|local development|not shipped|not in the npm package|not part of the npm package)\b`)
	for _, file := range files {
		lines := strings.Split(readIfExists(file), "\n")
		inClaimSection := false
		for i, line := range lines {
			if strings.HasPrefix(line, "#") {
				inClaimSection = headingRe.MatchString(line)
			}
			if !(inClaimSection || claimRe.MatchString(line)) || repoOnlyRe.MatchString(line) {
				continue
			}
			for _, repoOnlyPath := range []string{"specs/", "decisions/", "tests/"} {
				if strings.Contains(line, repoOnlyPath) {
					diags = append(diags, fmt.Sprintf("%s:%d claims repo-only path '%s' is shipped/package content; label it repo-only development context", slashRel(root, file), i+1, repoOnlyPath))
				}
			}
		}
	}
	return diags
}

// checkGoToolchain grades the two Go facts no gate phase can own. A toolchain absent
// from PATH is the condition under which every probed phase silently declines to
// materialize, so the diagnostic naming it has to come from a check that runs without
// Go. The cross-compile matrix is ship-tier, which the dev phase table does not reach.
// Both grade any root cheaply.
func checkGoToolchain(root string) []string {
	if !exists(filepath.Join(root, "go.mod")) {
		return nil
	}
	if _, err := exec.LookPath("go"); err != nil {
		return []string{"go.mod present but no Go toolchain on PATH — the compiled core is load-bearing; install Go"}
	}
	return crossCompileMatrix(root, filepath.Join(root, "scripts", "go-build.sh"))
}

// TestResidualCheckBuildsNothing keeps the residual check from competing with the run
// owner's single build. An assertion about dist/bench alone would miss a build aimed at
// a throwaway path.
func TestResidualCheckBuildsNothing(t *testing.T) {
	root := t.TempDir()
	writeFixtureFile(t, filepath.Join(root, "go.mod"), "module fixture\n\ngo 1.25\n")
	writeFixtureFile(t, filepath.Join(root, "cmd", "bench", "main.go"), "package main\n\nfunc main() {}\n")
	// Records the invocation instead of building, so the assertion reads whether the
	// check reached the helper at all rather than what it asked it to produce.
	writeFixtureFile(t, filepath.Join(root, "scripts", "go-build.sh"),
		"#!/usr/bin/env bash\nprintf '%s\\n' \"$2\" > \"$1/recorded-out\"\n")

	checkGoToolchain(root)

	if recorded := strings.TrimSpace(readIfExists(filepath.Join(root, "recorded-out"))); recorded != "" {
		t.Fatalf("residual check invoked the build helper (output path %q); the run owner owns the only ordinary build", recorded)
	}
	if exists(filepath.Join(root, "dist", "bench")) {
		t.Fatal("residual check wrote the graded root's dist/bench")
	}
}

// TestResidualCheckReportsAbsentToolchain keeps the only diagnostic a host with no Go
// can produce. Without it such a host grades green on a tree whose compiled core is
// load-bearing, because every phase that would have noticed is gated on the same absent
// toolchain.
func TestResidualCheckReportsAbsentToolchain(t *testing.T) {
	root := t.TempDir()
	writeFixtureFile(t, filepath.Join(root, "go.mod"), "module fixture\n\ngo 1.25\n")
	t.Setenv("PATH", "")

	diags := checkGoToolchain(root)

	if !containsDiagnostic(diags, "go.mod present but no Go toolchain on PATH") {
		t.Fatalf("residual check lost the absent-toolchain diagnostic: %#v", diags)
	}
}

func writeFixtureFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func offlineSmokeRecoversInterruptedStages(smoke string) bool {
	for _, anchor := range []string{
		"trap interrupt INT TERM HUP",
		"offline_process comparison",
		"offline_stage installation",
		"offline_stage registry-service",
		"offline_stage removal",
	} {
		if strings.Count(smoke, anchor) != 1 {
			return false
		}
	}
	return strings.Contains(smoke, "exit 130")
}

func containsKey(records []requirementRecord, want string) bool {
	for _, record := range records {
		if record.Key == want {
			return true
		}
	}
	return false
}

func checkGuardHeaderManifests(root string) []string {
	var diags []string
	for _, guard := range []struct {
		name string
		path string
	}{
		{"block-dangerous-git", filepath.Join(root, ".bench", "hooks", "block-dangerous-git.sh")},
		{"check-agent-line", filepath.Join(root, ".bench", "hooks", "check-agent-line.sh")},
		{"stop", filepath.Join(root, ".bench", "hooks", "stop.sh")},
		{"session-start", filepath.Join(root, ".bench", "hooks", "session-start.sh")},
		{"worktree-lifecycle", filepath.Join(root, ".bench", "hooks", "worktree-lifecycle.sh")},
		{"pre-push", filepath.Join(root, "internal", "adopt", "prepush.sh")},
	} {
		path := guard.path
		if !exists(path) {
			continue
		}
		fields, err := guards.HeaderFields(path)
		if err != nil {
			diags = append(diags, fmt.Sprintf("guard %s manifest unreadable", guard.name))
			continue
		}
		for _, key := range fields.MissingRequired() {
			diags = append(diags, fmt.Sprintf("guard %s manifest missing %s key", guard.name, key))
		}
		if guard.name == "session-start" && fields["denies"] != "nothing (informational)" {
			diags = append(diags, "session-start is not classified informational (denies: nothing)")
		}
	}
	return diags
}

type platformTarget struct {
	Goos   string `json:"goos"`
	Goarch string `json:"goarch"`
}

func platformMatrix(path string) ([]platformTarget, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var plan struct {
		Targets []platformTarget `json:"targets"`
	}
	if err := json.Unmarshal(data, &plan); err != nil {
		return nil, err
	}
	return plan.Targets, nil
}

// hostilePayloadPlanters extends the shared producer partition to the allowlist: the
// byte-shape half reuses hostileSkillPlanters (one owner for FIFO, both link forms,
// oversized, and invalid UTF-8), and the semantic half adds the row defects only the
// canonical parser can see. A JSON-decode-only reader survives the first group and
// admits the second.
func hostilePayloadPlanters(t *testing.T) map[string]func(*testing.T, string) {
	t.Helper()
	planters := map[string]func(*testing.T, string){}
	for kind, plant := range hostileSkillPlanters {
		planters[kind] = plant
	}
	for kind, content := range map[string]string{
		"empty":            "",
		"invalid JSON":     "{not json",
		"unknown audience": `[{"source":"AGENTS.md","audience":"everyone"}]`,
		"empty source":     `[{"source":"","audience":"kit-only"}]`,
		"unsafe source":    `[{"source":"../escape","audience":"kit-only"}]`,
		"duplicate source": `[{"source":"AGENTS.md","audience":"consumer"},{"source":"AGENTS.md","audience":"kit-only"}]`,
	} {
		planters[kind] = func(t *testing.T, path string) {
			t.Helper()
			if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	return planters
}

// TestRegisteredPayloadConsumersRefuseHostilePayload is the allowlist composition row.
// The package-core guard, the registered package-shipped-surface check, and the
// packagesurface contract-document inventory all read the same tracked file, so each
// has to complete over a hostile or invalid one and name the path it refused —
// otherwise skills-index can go green while a package reader hangs in open(2) or ships
// an inventory derived from rows the allowlist forbids.
func TestRegisteredPayloadConsumersRefuseHostilePayload(t *testing.T) {
	const rel = ".bench/consumer-payload.json"
	// One packed file list, enough for the guard to have something to grade: the
	// refusal under test happens before it is consulted.
	const packJSON = `[{"files":[{"path":"package/AGENTS.md"}]}]`
	for kind, plant := range hostilePayloadPlanters(t) {
		t.Run(kind, func(t *testing.T) {
			root := t.TempDir()
			if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
				t.Fatal(err)
			}
			// package.json is the shipped-surface check's own entry condition; without
			// it that consumer returns before it reaches the allowlist.
			if err := os.WriteFile(filepath.Join(root, "package.json"), []byte(`{"files":[".agents/"]}`), 0o644); err != nil {
				t.Fatal(err)
			}
			plant(t, filepath.Join(root, ".bench", "consumer-payload.json"))

			consumers := map[string]func() []string{
				// The guard function itself, not the registered package-core-guard
				// binding: that binding reaches this reader only through a real
				// `npm pack`, which a fixture root cannot satisfy.
				"package-core-guard": func() []string { return checkNoKitOnlyPackedAssets(root, packJSON) },
				"package-shipped-surface": func() []string {
					binding, bound := conformanceChecks["package-shipped-surface"]
					if !bound {
						t.Fatal("package-shipped-surface conformance owner is not bound")
					}
					return binding.run(root, root, registry.Dev)
				},
				"packagesurface.ContractDocumentInputs": func() []string {
					_, err := packagesurface.ContractDocumentInputs(root)
					if err == nil {
						return nil
					}
					return []string{err.Error()}
				},
			}
			for name, run := range consumers {
				t.Run(name, func(t *testing.T) {
					done := make(chan []string, 1)
					go func() { done <- run() }()
					select {
					case diags := <-done:
						if !containsDiagnostic(diags, rel) {
							t.Fatalf("%s over a %s %s produced %q, want a diagnostic naming the refused path", name, kind, rel, diags)
						}
					case <-time.After(bounds.TestDeadline(0)):
						t.Fatalf("%s blocked on a %s %s, so it opened the path before classifying it", name, kind, rel)
					}
				})
			}
		})
	}
}
