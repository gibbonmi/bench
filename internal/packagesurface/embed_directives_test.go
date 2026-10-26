package packagesurface

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

// TestScanEmbedDirectivesReadsEveryGoFile proves that the scan reads a test file and a file
// outside cmd/ and internal/, resolves each pattern against its file's directory, drops the
// all: prefix, and skips the Git entry of the root.
func TestScanEmbedDirectivesReadsEveryGoFile(t *testing.T) {
	root := t.TempDir()
	for name, body := range map[string]string{
		"docs/gen/gen_test.go": "package gen\n\n//go:embed fixture.txt all:assets\nvar fixture string\n",
		"tools/tool.go":        "package tools\n\n//go:embed \"*.txt\"\nvar notes string\n",
		".git/hooks/hook.go":   "package hooks\n\n//go:embed hidden.txt\nvar hidden string\n",
	} {
		file := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(file, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	got, err := ScanEmbedDirectives(root)
	slices.Sort(got)
	if want := []string{"docs/gen/assets", "docs/gen/fixture.txt", "tools/*.txt"}; err != nil || !slices.Equal(got, want) {
		t.Fatalf("ScanEmbedDirectives = %q, %v; want %q", got, err, want)
	}
}

// TestEmbedDirectivesMatchTheirPaths proves that a directive matches its exact path, a path
// below it, a glob match, and a path below a glob match, and nothing else.
func TestEmbedDirectivesMatchTheirPaths(t *testing.T) {
	directives := EmbedDirectives{"docs/x.txt", "docs/assets", "notes/*.txt", "tree/*"}
	for rel, want := range map[string]bool{
		"docs/x.txt":         true,
		"docs/assets/a.txt":  true,
		"notes/a.txt":        true,
		"tree/sub/deep.md":   true,
		"docs/x.txt.bak":     false,
		"docs/assets2/a.txt": false,
		"notes/sub/a.txt":    false,
		"docs":               false,
	} {
		if got := directives.Matches(rel); got != want {
			t.Errorf("Matches(%q) = %t; want %t", rel, got, want)
		}
	}
}
