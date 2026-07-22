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

	"github.com/gibbonmi/bench/internal/guards"
	"github.com/gibbonmi/bench/internal/packagesurface"
)

func checkPackageCoreAndGuards(root string) []string {
	var diags []string
	diags = append(diags, checkPackageFiles(root)...)
	diags = append(diags, checkGoCore(root)...)
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
		diags = append(diags, "npm pack --dry-run failed")
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
// binary's own copy, so a canary fixture's mutated allowlist is what gets graded.
func checkNoKitOnlyPackedAssets(root, packJSON string) []string {
	prefixes, err := kitOnlyAllowlistPrefixes(root)
	if err != nil || len(prefixes) == 0 {
		return nil
	}
	var packs []struct {
		Files []struct {
			Path string `json:"path"`
		} `json:"files"`
	}
	if err := json.Unmarshal([]byte(packJSON), &packs); err != nil || len(packs) == 0 {
		return nil
	}
	var diags []string
	for _, file := range packs[0].Files {
		rel := strings.TrimPrefix(file.Path, "package/")
		if excludedByAllowlistPrefix(rel, prefixes) {
			diags = append(diags, "npm package includes kit-only allowlist asset "+rel)
		}
	}
	return diags
}

func kitOnlyAllowlistPrefixes(root string) ([]string, error) {
	data, err := os.ReadFile(filepath.Join(root, ".bench", "consumer-payload.json"))
	if err != nil {
		return nil, err
	}
	var rows []struct {
		Source   string `json:"source"`
		Audience string `json:"audience"`
	}
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil, err
	}
	var out []string
	for _, r := range rows {
		if r.Audience == "kit-only" {
			out = append(out, r.Source)
		}
	}
	return out, nil
}

func excludedByAllowlistPrefix(rel string, prefixes []string) bool {
	for _, p := range prefixes {
		if rel == p || strings.HasPrefix(rel, p+"/") {
			return true
		}
	}
	return false
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

func checkGoCore(root string) []string {
	if !exists(filepath.Join(root, "go.mod")) {
		return nil
	}
	if _, err := exec.LookPath("go"); err != nil {
		return []string{"go.mod present but no Go toolchain on PATH — the compiled core is load-bearing; install Go"}
	}
	var diags []string
	if probe := runAtCleanEnv(root, "gofmt", "-l", "."); probe != nil && strings.TrimSpace(probe.Stdout) != "" {
		diags = append(diags, "gofmt: unformatted Go files: "+strings.Join(strings.Fields(probe.Stdout), " "))
	}
	buildHelper := filepath.Join(root, "scripts", "go-build.sh")
	if exists(buildHelper) {
		// Throwaway output path, never root's dist/bench — the gate's serial
		// build phase owns the one real write (rationale: gate.BenchkitPhases).
		tmp, err := os.MkdirTemp("", "bench-build-*")
		if err != nil {
			diags = append(diags, "go build setup failed: "+err.Error())
		} else {
			defer os.RemoveAll(tmp)
			if probe := runAtCleanEnv(root, "bash", buildHelper, root, filepath.Join(tmp, "bench")); probe == nil || probe.ExitCode != 0 {
				diags = append(diags, "go build failed")
			}
		}
	} else if probe := runAtCleanEnv(root, "go", "build", "./..."); probe == nil || probe.ExitCode != 0 {
		diags = append(diags, "go build failed")
	}
	if probe := runAtCleanEnv(root, "go", "vet", "./..."); probe == nil || probe.ExitCode != 0 {
		diags = append(diags, "go vet failed")
	}
	testPackages, ok := goCoreTestPackages(root)
	if !ok {
		diags = append(diags, "go list failed")
	} else if len(testPackages) > 0 {
		args := append([]string{"go", "test"}, testPackages...)
		if probe := runAtCleanEnv(root, args...); probe == nil || probe.ExitCode != 0 {
			diags = append(diags, "go test failed")
		}
	}
	const cleanupRaceTest = "TestConcurrentCleanupRecordsOneTransaction"
	race := runAtCleanEnv(root, "go", "test", "-race", "-count=1", "-v", "./internal/worktree", "-run", "^"+cleanupRaceTest+"$")
	if race == nil || race.ExitCode != 0 {
		diags = append(diags, "worktree cleanup race test failed")
	} else if !strings.Contains(race.Stdout, "=== RUN   "+cleanupRaceTest) {
		diags = append(diags, "worktree cleanup race test did not run")
	}
	diags = append(diags, crossCompileMatrix(root, buildHelper)...)
	return diags
}

func goCoreTestPackages(root string) ([]string, bool) {
	probe := runAtCleanEnv(root, "go", "list", "./...")
	if probe == nil || probe.ExitCode != 0 {
		return nil, false
	}
	var packages []string
	for _, pkg := range strings.Fields(probe.Stdout) {
		// The gate runs internal/contract separately with BENCH_CONTRACT_ROOT set.
		if isContractPackage(pkg) {
			continue
		}
		packages = append(packages, pkg)
	}
	return packages, true
}

func isContractPackage(pkg string) bool {
	return pkg == "internal/contract" ||
		strings.HasPrefix(pkg, "internal/contract/") ||
		strings.HasSuffix(pkg, "/internal/contract") ||
		strings.Contains(pkg, "/internal/contract/")
}

func TestCheckGoCoreDoesNotWriteRootDistBench(t *testing.T) {
	root := t.TempDir()
	writeFixtureFile(t, filepath.Join(root, "go.mod"), "module fixture\n\ngo 1.25\n")
	writeFixtureFile(t, filepath.Join(root, "cmd", "bench", "main.go"), "package main\n\nfunc main() {}\n")
	// Records the output path it was handed instead of building, so the assertion
	// reads the real argv checkGoCore passes to the helper.
	writeFixtureFile(t, filepath.Join(root, "scripts", "go-build.sh"),
		"#!/usr/bin/env bash\nprintf '%s\\n' \"$2\" > \"$1/recorded-out\"\n")

	diags := checkGoCore(root)
	for _, diag := range diags {
		if strings.Contains(diag, "go build failed") {
			t.Fatalf("build helper probe went red: %#v", diags)
		}
	}
	recorded := strings.TrimSpace(readIfExists(filepath.Join(root, "recorded-out")))
	if recorded == "" {
		t.Fatal("build helper was not invoked")
	}
	if strings.HasPrefix(recorded, root) {
		t.Fatalf("build output path %q is inside the tree under grade; the gate's serialized build phase owns the real dist/bench write", recorded)
	}
	if exists(filepath.Join(root, "dist", "bench")) {
		t.Fatal("checkGoCore wrote the graded root's dist/bench")
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

func TestGoCoreTestPackagesExcludesContractSubtreeOnly(t *testing.T) {
	h := NewHarness(t)
	packages, ok := goCoreTestPackages(h.KitRoot)
	if !ok {
		t.Fatal("goCoreTestPackages failed")
	}
	if containsPackageSuffix(packages, "/internal/contract") {
		t.Fatalf("goCoreTestPackages included internal/contract:\n%s", strings.Join(packages, "\n"))
	}
	if containsPackageSuffix(packages, "/internal/contract/axi") {
		t.Fatalf("goCoreTestPackages included internal/contract/axi:\n%s", strings.Join(packages, "\n"))
	}
	if !containsPackageSuffix(packages, "/internal/conformance") {
		t.Fatalf("goCoreTestPackages excluded internal/conformance:\n%s", strings.Join(packages, "\n"))
	}
}

func containsPackageSuffix(packages []string, suffix string) bool {
	for _, pkg := range packages {
		if strings.HasSuffix(pkg, suffix) {
			return true
		}
	}
	return false
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
