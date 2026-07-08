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

	"github.com/gibbonmi/bench/internal/packagesurface"
)

func checkPackageCoreAndGuards(root string) []string {
	var diags []string
	diags = append(diags, checkPackageFiles(root)...)
	diags = append(diags, checkGoCore(root)...)
	diags = append(diags, checkReleaseWorkflow(root)...)
	diags = append(diags, checkShippedIdentityStrings(root)...)
	diags = append(diags, checkGuardDescribeManifests(root)...)
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
		diags = append(diags, checkNpmPackAssets(probe.Stdout)...)
	}
	return append(diags, checkRepoOnlyPackageClaims(root)...)
}

func checkNpmPackAssets(packJSON string) []string {
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
	for _, required := range packagesurface.RequiredPackAssets {
		if !files[required] {
			diags = append(diags, "npm package missing "+required)
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
		if probe := runAtCleanEnv(root, "bash", buildHelper, root, filepath.Join(root, "dist", "bench")); probe == nil || probe.ExitCode != 0 {
			diags = append(diags, "go build failed")
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
	if exists(filepath.Join(root, "scripts", "platforms.json")) && exists(buildHelper) {
		matrix, err := platformMatrix(filepath.Join(root, "scripts", "platforms.json"))
		if err != nil {
			diags = append(diags, "platform matrix unreadable: "+err.Error())
		}
		tmp, err := os.MkdirTemp("", "bench-cross-*")
		if err != nil {
			diags = append(diags, "cross-compile setup failed: "+err.Error())
		} else {
			defer os.RemoveAll(tmp)
			for _, target := range matrix {
				env := append(conformanceSubprocessEnv(), "GOOS="+target.Goos, "GOARCH="+target.Goarch)
				probe := runAtEnv(root, env, "bash", buildHelper, root, filepath.Join(tmp, "bench-"+target.Goos+"-"+target.Goarch))
				if probe == nil || probe.ExitCode != 0 {
					diags = append(diags, fmt.Sprintf("cross-compile failed: %s/%s", target.Goos, target.Goarch))
				}
			}
		}
	}
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

func checkReleaseWorkflow(root string) []string {
	if !exists(filepath.Join(root, "scripts", "platforms.json")) {
		return nil
	}
	wf := filepath.Join(root, ".github", "workflows", "release.yml")
	if !exists(wf) {
		return []string{"release workflow missing (.github/workflows/release.yml)"}
	}
	text := readIfExists(wf)
	var diags []string
	if !regexp.MustCompile(`(?m)^\s*tags:`).MatchString(text) {
		diags = append(diags, "release workflow does not trigger on tags")
	}
	if !strings.Contains(text, "scripts/platforms.json") {
		diags = append(diags, "release workflow does not derive targets from the matrix (scripts/platforms.json)")
	}
	if !strings.Contains(text, "scripts/gen-platform-packages.sh") {
		diags = append(diags, "release workflow does not run the platform-package generator")
	}
	if !strings.Contains(text, "npm publish") {
		diags = append(diags, "release workflow does not publish to npm")
	}
	if !strings.Contains(text, "provenance") {
		diags = append(diags, "release workflow does not publish with provenance")
	}
	return diags
}

// checkShippedIdentityStrings sweeps the shipped surface plus the identity sources
// (package.json, the release workflow) for npm strings naming a package the project
// does not own. A half-done redbench rename leaves one behind; each survivor becomes a
// diag rather than shipping. CHANGELOG.md is excluded by design — history names the old
// identity and is not rewritten; dist/ is a compiled artifact, not authored surface.
func checkShippedIdentityStrings(root string) []string {
	stale := []string{
		"@benchkit",
		"npx benchkit",
		"npm i -g benchkit",
		"npm package `benchkit`",
	}
	surface := append(shippedFiles(root),
		filepath.Join(root, "package.json"),
		filepath.Join(root, ".github", "workflows", "release.yml"),
	)
	var diags []string
	for _, file := range surface {
		rel := slashRel(root, file)
		if rel == "CHANGELOG.md" || strings.HasPrefix(rel, "dist/") {
			continue
		}
		for i, line := range strings.Split(readIfExists(file), "\n") {
			for _, s := range stale {
				if strings.Contains(line, s) {
					diags = append(diags, fmt.Sprintf("%s:%d ships stale npm identity %q; the project publishes as redbench", rel, i+1, s))
				}
			}
		}
	}
	return diags
}

// TestShippedIdentityStringSweepBites is the recorded bite proof for
// checkShippedIdentityStrings (per craft-gate): a clean shipped surface passes, a stale
// string confined to CHANGELOG.md or under dist/ is ignored, and each matched pattern —
// in any of the swept locations (shipped doc, package.json, the release workflow) —
// makes the sweep fire.
func TestShippedIdentityStringSweepBites(t *testing.T) {
	root := t.TempDir()
	write := func(rel, body string) {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	// Clean surface: a stale string only in CHANGELOG.md or under dist/ (a compiled
	// artifact, not authored surface) must NOT fire — both are excluded by design.
	write("package.json", `{"files":["README.md","projects/","CHANGELOG.md","dist/"]}`+"\n")
	write("README.md", "install with npx redbench link\n")
	write("projects/benchkit.md", "the npm package `redbench`.\n")
	write("CHANGELOG.md", "0.1.0 advertised npx benchkit link\n")
	write("dist/bench", "stub with @benchkit/linux-x64 bytes\n")
	if diags := checkShippedIdentityStrings(root); len(diags) != 0 {
		t.Fatalf("clean surface (stale strings only in CHANGELOG and dist/): want no diagnostics, got %v", diags)
	}

	// Each pattern must fire, and each swept location (shipped doc, the profile,
	// package.json, the release workflow) must be reachable.
	cases := []struct{ file, body, want string }{
		{"README.md", "run npx benchkit link\n", "npx benchkit"},
		{"README.md", "or npm i -g benchkit today\n", "npm i -g benchkit"},
		{"projects/benchkit.md", "the npm package `benchkit`.\n", "npm package `benchkit`"},
		{"package.json", `{"files":["README.md"],"optionalDependencies":{"@benchkit/linux-x64":"1"}}` + "\n", "@benchkit"},
		{".github/workflows/release.yml", "  for d in dist/packages/@benchkit/*/; do\n", "@benchkit"},
	}
	for _, c := range cases {
		root := t.TempDir()
		// A minimal files[] so shippedFiles resolves; package.json and release.yml are
		// swept whether or not files[] lists them.
		write2 := func(rel, body string) {
			full := filepath.Join(root, filepath.FromSlash(rel))
			if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		write2("package.json", `{"files":["README.md","projects/"]}`+"\n")
		write2("README.md", "clean\n")
		write2("projects/benchkit.md", "clean\n")
		write2(c.file, c.body)
		diags := checkShippedIdentityStrings(root)
		if len(diags) != 1 || !strings.Contains(diags[0], slashRel(root, filepath.Join(root, filepath.FromSlash(c.file)))) || !strings.Contains(diags[0], c.want) {
			t.Fatalf("seeded %q in %s: want one diagnostic naming %q, got %v", c.want, c.file, c.want, diags)
		}
	}
}

func checkGuardDescribeManifests(root string) []string {
	var diags []string
	for _, guard := range []string{"block-dangerous-git", "check-agent-line", "stop", "session-start"} {
		path := filepath.Join(root, ".bench", "hooks", guard+".sh")
		if !exists(path) {
			continue
		}
		probe := runAtCleanEnv(root, "bash", path, "--describe")
		if probe == nil || probe.ExitCode != 0 {
			exit := 1
			if probe != nil {
				exit = probe.ExitCode
			}
			diags = append(diags, fmt.Sprintf("guard %s --describe did not exit 0 (exit %d)", guard, exit))
			continue
		}
		for _, key := range []string{"name", "boundary", "denies", "why"} {
			if !regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(key) + `: `).MatchString(probe.Stdout) {
				diags = append(diags, fmt.Sprintf("guard %s --describe manifest missing %s key", guard, key))
			}
		}
		if guard == "session-start" && !regexp.MustCompile(`(?m)^denies: nothing \(informational\)$`).MatchString(probe.Stdout) {
			diags = append(diags, "session-start --describe is not classified informational (denies: nothing)")
		}
	}
	return diags
}

// shippedFiles returns every file the npm tarball would carry, expanding package.json
// files[] directory entries recursively. It is the single source of the shipped-file
// set; callers filter it (markdown prose sweep, identity-string sweep).
func shippedFiles(root string) []string {
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
	var out []string
	for _, entry := range pkg.Files {
		full := filepath.Join(root, filepath.FromSlash(entry))
		info, err := os.Stat(full)
		if err != nil {
			continue
		}
		if info.IsDir() {
			_ = filepath.WalkDir(full, func(path string, d os.DirEntry, err error) error {
				if err == nil && !d.IsDir() {
					out = append(out, path)
				}
				return nil
			})
			continue
		}
		out = append(out, full)
	}
	return uniqueSorted(out)
}

func packageMarkdownFiles(root string) []string {
	var out []string
	for _, file := range shippedFiles(root) {
		if strings.HasSuffix(file, ".md") {
			out = append(out, file)
		}
	}
	return out
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
	var targets []platformTarget
	if err := json.Unmarshal(data, &targets); err != nil {
		return nil, err
	}
	return targets, nil
}
