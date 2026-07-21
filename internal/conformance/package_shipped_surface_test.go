package conformance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestRootPackageProjectMetadata(t *testing.T) {
	root := NewHarness(t).Root
	var pkg struct {
		Repository string `json:"repository"`
		Homepage   string `json:"homepage"`
		Bugs       string `json:"bugs"`
		Author     string `json:"author"`
	}
	data, err := os.ReadFile(filepath.Join(root, "package.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		t.Fatal(err)
	}
	if pkg.Repository != "git+https://github.com/gibbonmi/bench.git" || pkg.Homepage != "https://github.com/gibbonmi/bench#readme" || pkg.Bugs != "https://github.com/gibbonmi/bench/issues" || pkg.Author != "gibbonmi" {
		t.Fatalf("root project metadata = %+v", pkg)
	}
}

func TestRepairScriptPolicyAndManifestPathParity(t *testing.T) {
	root := NewHarness(t).Root
	script := readIfExists(filepath.Join(root, "bin", "bench-repair-binary.mjs"))
	for _, fact := range []string{
		"const FETCH_DEADLINE_MS = 60_000;",
		"const DOWNLOAD_LIMIT = 100 * 1024 * 1024;",
		"const DECOMPRESSED_LIMIT = 200 * 1024 * 1024;",
	} {
		if !strings.Contains(script, fact) {
			t.Fatalf("repair policy omits %q", fact)
		}
	}
	var requirements struct {
		BinaryPinManifest struct {
			Path string `json:"path"`
		} `json:"binary_pin_manifest"`
	}
	data, err := os.ReadFile(filepath.Join(root, "internal", "releaseevidence", "requirements.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &requirements); err != nil {
		t.Fatal(err)
	}
	want := `const PIN_MANIFEST_PATH = "` + requirements.BinaryPinManifest.Path + `";`
	if requirements.BinaryPinManifest.Path == "" || !strings.Contains(script, want) {
		t.Fatalf("repair pin path does not match requirement %q", requirements.BinaryPinManifest.Path)
	}
}

func TestRepairAppearsInColdPickupInventory(t *testing.T) {
	root := NewHarness(t).Root
	guide := readIfExists(filepath.Join(root, ".bench", "BENCH.md"))
	if !strings.Contains(guide, "`bench repair") {
		t.Fatal("CLI inventory omits bench repair")
	}
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
	var out, excluded []string
	for _, entry := range pkg.Files {
		if strings.HasPrefix(entry, "!") {
			excluded = append(excluded, strings.TrimPrefix(entry, "!"))
			continue
		}
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
	filtered := out[:0]
	for _, file := range out {
		rel, err := filepath.Rel(root, file)
		if err != nil {
			continue
		}
		if !npmFilesExcluded(excluded, filepath.ToSlash(rel)) {
			filtered = append(filtered, file)
		}
	}
	return uniqueSorted(filtered)
}

func npmFilesExcluded(patterns []string, value string) bool {
	for _, pattern := range patterns {
		if npmFilesPatternMatches(pattern, value) {
			return true
		}
	}
	return false
}

func npmFilesPatternMatches(pattern, value string) bool {
	quoted := regexp.QuoteMeta(pattern)
	quoted = strings.ReplaceAll(quoted, `\*\*/`, `(?:.*/)?`)
	quoted = strings.ReplaceAll(quoted, `\*\*`, `.*`)
	quoted = strings.ReplaceAll(quoted, `\*`, `[^/]*`)
	return regexp.MustCompile("^" + quoted + "$").MatchString(value)
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

// checkShippedIdentityStrings sweeps the shipped surface plus the identity sources
// (package.json, the release workflow) for npm strings naming a package the project
// does not own. A half-done redbench rename leaves one behind; each survivor becomes a
// diag rather than shipping. dist/ is a compiled artifact, not authored surface.
func checkShippedIdentityStrings(root string) []string {
	// The spec enumerates three examples (@benchkit, npx benchkit, npm i -g benchkit);
	// the profile's "npm package `benchkit`" claim and the README/doctor uninstall line
	// are the same class of unowned-identity string on surfaces this rename touched, so
	// they join the net (reviewer-approved).
	stale := []string{
		"@benchkit",
		"npx benchkit",
		"npm i -g benchkit",
		"npm uninstall -g benchkit",
		"npm package `benchkit`",
	}
	surface := append(shippedFiles(root),
		filepath.Join(root, "package.json"),
		filepath.Join(root, ".github", "workflows", "release.yml"),
	)
	var diags []string
	for _, file := range surface {
		rel := slashRel(root, file)
		if strings.HasPrefix(rel, "dist/") {
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

func checkUserFacingBenchkitStrings(root string) []string {
	var diags []string
	for _, rel := range []string{"cmd/bench/main.go", "bin/bench.sh", "scripts/build-release-evidence.mjs"} {
		file := filepath.Join(root, filepath.FromSlash(rel))
		for lineNo, line := range strings.Split(readIfExists(file), "\n") {
			if strings.Contains(line, "benchkit") {
				diags = append(diags, fmt.Sprintf("%s:%d exposes internal identity benchkit in user-facing output", rel, lineNo+1))
			}
		}
	}
	return diags
}

func TestUserFacingBenchkitSweepBites(t *testing.T) {
	root := t.TempDir()
	for _, rel := range []string{"cmd/bench/main.go", "bin/bench.sh", "scripts/build-release-evidence.mjs"} {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte("clean\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if got := checkUserFacingBenchkitStrings(root); len(got) != 0 {
		t.Fatalf("clean sweep = %v", got)
	}
	if err := os.WriteFile(filepath.Join(root, "bin", "bench.sh"), []byte("echo benchkit\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := checkUserFacingBenchkitStrings(root); len(got) != 1 {
		t.Fatalf("seeded sweep = %v", got)
	}
}

// TestShippedIdentityStringSweepBites is the recorded bite proof for
// checkShippedIdentityStrings (per craft-gate): a clean shipped surface passes, a stale
// string confined under dist/ is ignored, and each matched pattern —
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

	// Clean surface: a stale string only under dist/ (a compiled artifact, not
	// authored surface) must NOT fire.
	write("package.json", `{"files":["README.md","projects/","CHANGELOG.md","dist/"]}`+"\n")
	write("README.md", "install with npx redbench link\n")
	write("projects/benchkit.md", "the npm package `redbench`.\n")
	write("CHANGELOG.md", "# Changelog\n")
	write("dist/bench", "stub with @benchkit/linux-x64 bytes\n")
	if diags := checkShippedIdentityStrings(root); len(diags) != 0 {
		t.Fatalf("clean surface (stale strings only in dist/): want no diagnostics, got %v", diags)
	}
	write("package.json", `{"files":["README.md","projects/","CHANGELOG.md","internal/","!internal/**/*_test.go"]}`+"\n")
	write("internal/fixture_test.go", "npx benchkit\n")
	if diags := checkShippedIdentityStrings(root); len(diags) != 0 {
		t.Fatalf("negated npm test source should not be shipped, got %v", diags)
	}

	// Each pattern must fire, and each swept location (shipped doc, the profile,
	// package.json, the release workflow) must be reachable.
	cases := []struct{ file, body, want string }{
		{"README.md", "run npx benchkit link\n", "npx benchkit"},
		{"CHANGELOG.md", "changed npx benchkit link\n", "npx benchkit"},
		{"README.md", "or npm i -g benchkit today\n", "npm i -g benchkit"},
		{"README.md", "remove with npm uninstall -g benchkit\n", "npm uninstall -g benchkit"},
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
		write2("package.json", `{"files":["README.md","projects/","CHANGELOG.md"]}`+"\n")
		write2("README.md", "clean\n")
		write2("projects/benchkit.md", "clean\n")
		write2(c.file, c.body)
		diags := checkShippedIdentityStrings(root)
		if len(diags) != 1 || !strings.Contains(diags[0], slashRel(root, filepath.Join(root, filepath.FromSlash(c.file)))) || !strings.Contains(diags[0], c.want) {
			t.Fatalf("seeded %q in %s: want one diagnostic naming %q, got %v", c.want, c.file, c.want, diags)
		}
	}
}
