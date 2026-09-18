package anchors

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"
)

// writeRegistryFile plants one file in a temporary root's anchor registry directory.
func writeRegistryFile(t *testing.T, root, name, body string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(RegistryDir), filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
}

func scan(t *testing.T, root string) map[string][]string {
	t.Helper()
	refs, err := ReferencingFiles(root)
	if err != nil {
		t.Fatalf("ReferencingFiles: %v", err)
	}
	return refs
}

// TestReferencingFilesIncludeTestFiles covers SC4. A test file that states the guidance
// path independently joins the files that reference it.
func TestReferencingFilesIncludeTestFiles(t *testing.T) {
	root := t.TempDir()
	writeRegistryFile(t, root, "registry_a.go", "package anchors\n\nvar a = Anchor{File: \".agents/x/SKILL.md\"}\n")
	writeRegistryFile(t, root, "registry_a_test.go", "package anchors\n\nvar want = \".agents/x/SKILL.md\"\n")
	got := scan(t, root)[".agents/x/SKILL.md"]
	want := []string{"internal/anchors/registry_a.go", "internal/anchors/registry_a_test.go"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("referencing files = %q, want %q", got, want)
	}
}

// TestReferencingFilesReadEveryStringLiteral covers SC5. A constant and a positional
// value name the path without a File: key.
func TestReferencingFilesReadEveryStringLiteral(t *testing.T) {
	root := t.TempDir()
	writeRegistryFile(t, root, "constant.go", "package anchors\n\nconst guide = \".agents/const guide.md\"\n")
	writeRegistryFile(t, root, "positional.go", "package anchors\n\nvar rows = [][2]string{{\"needle\", \".agents/positional.md\"}}\n")
	refs := scan(t, root)
	for path, file := range map[string]string{".agents/const guide.md": "internal/anchors/constant.go", ".agents/positional.md": "internal/anchors/positional.go"} {
		if got := refs[path]; !reflect.DeepEqual(got, []string{file}) {
			t.Errorf("referencing files of %q = %q, want [%s]", path, got, file)
		}
	}
}

// TestReferencingFilesUnquoteRawLiterals covers SC6. A backquoted literal matches by its
// value, not by its quoted source text.
func TestReferencingFilesUnquoteRawLiterals(t *testing.T) {
	root := t.TempDir()
	writeRegistryFile(t, root, "raw.go", "package anchors\n\nvar raw = `.agents/raw.md`\nvar escaped = \"\\x2eagents/escaped.md\"\n")
	refs := scan(t, root)
	for _, path := range []string{".agents/raw.md", ".agents/escaped.md"} {
		if got := refs[path]; !reflect.DeepEqual(got, []string{"internal/anchors/raw.go"}) {
			t.Errorf("referencing files of %q = %q, want the raw.go file", path, got)
		}
	}
}

// TestReferencingFilesAbsentDirectory covers SC8. A root with no anchor registry
// directory, or with an empty one, has no referencing files and no refusal.
func TestReferencingFilesAbsentDirectory(t *testing.T) {
	root := t.TempDir()
	if refs := scan(t, root); len(refs) != 0 {
		t.Fatalf("absent directory = %v, want empty", refs)
	}
	if err := os.MkdirAll(filepath.Join(root, filepath.FromSlash(RegistryDir)), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if refs := scan(t, root); len(refs) != 0 {
		t.Fatalf("empty directory = %v, want empty", refs)
	}
}

// TestReferencingFilesSkipsNestedAndOtherFiles pins the top-level `.go` scope: the scan
// does not recurse, and it reads no file of another kind.
func TestReferencingFilesSkipsNestedAndOtherFiles(t *testing.T) {
	root := t.TempDir()
	writeRegistryFile(t, root, "nested/deep.go", "package nested\n\nvar x = \".agents/nested.md\"\n")
	writeRegistryFile(t, root, "notes.txt", "\".agents/notes.md\"\n")
	if refs := scan(t, root); len(refs) != 0 {
		t.Fatalf("nested and non-Go files = %v, want empty", refs)
	}
}

func assertRefused(t *testing.T, root, path string) {
	t.Helper()
	refs, err := ReferencingFiles(root)
	if err == nil || !strings.Contains(err.Error(), path) {
		t.Fatalf("ReferencingFiles = (%v, %v), want a refusal naming %s", refs, err, path)
	}
}

// TestReferencingFilesRefuseSymlink covers SC10. A link to bytes outside the tree
// refuses the scan by path and returns none of those bytes.
func TestReferencingFilesRefuseSymlink(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.go")
	if err := os.WriteFile(outside, []byte("package outside\n\nvar x = \".agents/outside.md\"\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	writeRegistryFile(t, root, "registry_a.go", "package anchors\n")
	if err := os.Symlink(outside, filepath.Join(root, filepath.FromSlash(RegistryDir), "link.go")); err != nil {
		t.Fatalf("Symlink: %v", err)
	}
	assertRefused(t, root, "internal/anchors/link.go")
}

// TestReferencingFilesRefuseFIFO covers SC11. The scan classifies a FIFO before it
// opens it, so the refusal returns without a blocking read.
func TestReferencingFilesRefuseFIFO(t *testing.T) {
	root := t.TempDir()
	writeRegistryFile(t, root, "registry_a.go", "package anchors\n")
	if err := syscall.Mkfifo(filepath.Join(root, filepath.FromSlash(RegistryDir), "pipe.go"), 0o644); err != nil {
		t.Fatalf("Mkfifo: %v", err)
	}
	done := make(chan error, 1)
	go func() {
		_, err := ReferencingFiles(root)
		done <- err
	}()
	select {
	case err := <-done:
		if err == nil || !strings.Contains(err.Error(), "internal/anchors/pipe.go") {
			t.Fatalf("ReferencingFiles error = %v, want a refusal naming internal/anchors/pipe.go", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("ReferencingFiles blocked on the FIFO")
	}
}

// TestReferencingFilesRefuseTokenizerError names the file whose source does not
// tokenize, so a broken registry file never reads as one with no literals.
func TestReferencingFilesRefuseTokenizerError(t *testing.T) {
	root := t.TempDir()
	writeRegistryFile(t, root, "broken.go", "package anchors\n\nvar x = \"unterminated\n")
	assertRefused(t, root, "internal/anchors/broken.go")
}
