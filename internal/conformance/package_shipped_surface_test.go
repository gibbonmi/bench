package conformance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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

// checkShippedIdentityStrings sweeps the shipped surface plus the identity sources
// (package.json, the release workflow) for npm strings naming a package the project
// does not own. A half-done redbench rename leaves one behind; each survivor becomes a
// diag rather than shipping. CHANGELOG.md is excluded by design — history names the old
// identity and is not rewritten; dist/ is a compiled artifact, not authored surface.
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
