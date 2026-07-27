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

	kitpayload "github.com/gibbonmi/bench"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/guards"
	"github.com/gibbonmi/bench/internal/packagesurface"
)

func checkPackageCoreAndGuards(root string, tier registry.Tier) []string {
	var diags []string
	diags = append(diags, checkPackageFiles(root)...)
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

func checkPackageFiles(root string) []string {
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
		if !exists(filepath.Join(root, filepath.FromSlash(file))) {
			diags = append(diags, "package.json files[] missing "+file)
		}
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
	data, err := os.ReadFile(filepath.Join(root, ".bench", "consumer-payload.json"))
	if err != nil {
		return nil
	}
	var rows []kitpayload.PayloadRow
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil
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

// checkGoToolchain is what survives of the old checkGoCore after gofmt, the throwaway
// build validation, vet, the core `go test`, the filtered conformance suite, and the
// worktree race test became gate phases of their own. The two things left are the ones
// no phase can own: a toolchain absent from PATH is the condition under which every
// probed phase silently declines to materialize, so the diagnostic naming it has to
// come from a check that runs without Go; and the cross-compile matrix stays ship-tier,
// which the dev phase table does not reach. Both grade any root cheaply.
func checkGoToolchain(root string) []string {
	if !exists(filepath.Join(root, "go.mod")) {
		return nil
	}
	if _, err := exec.LookPath("go"); err != nil {
		return []string{"go.mod present but no Go toolchain on PATH — the compiled core is load-bearing; install Go"}
	}
	return crossCompileMatrix(root, filepath.Join(root, "scripts", "go-build.sh"))
}

// TestResidualCheckBuildsNothing pins the deletion of the throwaway build validation.
// The gate's build phase is now the single build, so the residual check must reach no
// build helper at all — a deletion that left the probe behind passes every other test
// in the package, and a probe writing to a throwaway path passes an assertion about
// dist/bench alone.
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
		t.Fatalf("residual check invoked the build helper (output path %q); the gate's build phase owns the only build", recorded)
	}
	if exists(filepath.Join(root, "dist", "bench")) {
		t.Fatal("residual check wrote the graded root's dist/bench")
	}
}

// TestResidualCheckReportsAbsentToolchain keeps the one diagnostic that outlives the
// steps it used to introduce. Without it a host with no Go grades green on a tree whose
// compiled core is load-bearing, because every phase that would have noticed is gated
// on the same absent toolchain.
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
