package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// treeSnapshotSourceFile is the one file allowed to run and parse `ls-tree`.
const treeSnapshotSourceFile = "tree_snapshot.go"

// [PS41] One parser. The stripped whole-tree identity and the per-component identities are
// two readings of one snapshot, and a second `ls-tree` invocation beside this one would be a
// second derivation of the same fact — free to drift in its flags, its separator, or its
// handling of a malformed entry, so that one identity sees a path the other does not. The
// claim is about the package's shape rather than about any one call's result, so it is
// graded over the package's own sources; test files are out of scope, where listing paths
// for an assertion is not an identity's input.
func TestTreeSnapshotIsTheOnlyListingParser(t *testing.T) {
	sources, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	var parsers []string
	for _, source := range sources {
		if strings.HasSuffix(source, "_test.go") {
			continue
		}
		body, err := os.ReadFile(source)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(body), `"ls-tree"`) {
			parsers = append(parsers, source)
		}
	}
	if len(parsers) != 1 || parsers[0] != treeSnapshotSourceFile {
		t.Fatalf("files running ls-tree = %v, want only %s — every identity reads one parsed snapshot", parsers, treeSnapshotSourceFile)
	}
}

// literalTreeSource feeds a fabricated listing to generation capture: malformed metadata
// cannot be induced through a real Git tree without replacing Git itself.
type literalTreeSource struct{ listing string }

func (s literalTreeSource) tree() (string, error)       { return strings.Repeat("a", 40), nil }
func (s literalTreeSource) list(string) (string, error) { return s.listing, nil }
func (s literalTreeSource) blob(string) ([]byte, error) { return nil, nil }

// [RM1] Metadata that is not git's mode-type-object shape refuses the snapshot at parse,
// so a malformed listing never authorizes a generation that identities would hash.
func TestTreeSnapshotRefusesMalformedListingMetadata(t *testing.T) {
	object := strings.Repeat("a", 40)
	for name, listing := range map[string]string{
		"missing object field": "100644 blob\tREADME.md\x00",
		"extra metadata field": "100644 blob " + object + " extra\tREADME.md\x00",
		"non-octal mode":       "10064x blob " + object + "\tREADME.md\x00",
		"unknown object type":  "100644 glob " + object + "\tREADME.md\x00",
		"mismatched type pair": "120000 commit " + object + "\tREADME.md\x00",
		"non-hex object":       "100644 blob " + strings.Repeat("z", 40) + "\tREADME.md\x00",
		"truncated object":     "100644 blob " + object[:12] + "\tREADME.md\x00",
		"empty metadata":       "\tREADME.md\x00",
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := parseTreeSnapshot(listing); err == nil {
				t.Fatalf("malformed metadata %q parsed into a snapshot", listing)
			}
			if generation, err := captureTreeGeneration(literalTreeSource{listing: listing}); err == nil || generation != nil {
				t.Fatalf("malformed metadata %q authorized a generation: %v, %v", listing, generation, err)
			}
		})
	}
}

// [RM2] The shapes git actually records — regular, executable, and symlink entries — keep
// parsing through a real working-tree capture, and blob reads keep their non-blob refusal.
func TestTreeSnapshotCaptureAcceptsRealGitObjectShapes(t *testing.T) {
	root := t.TempDir()
	gitRun(t, root, "init", "-q")
	writeGateTestFile(t, root, "plain.txt", "content\n", 0o644)
	writeGateTestFile(t, root, "tool.sh", "#!/bin/sh\n", 0o755)
	if err := os.Symlink("plain.txt", filepath.Join(root, "link")); err != nil {
		t.Fatal(err)
	}
	generation, err := captureWorkingTree(root)
	if err != nil {
		t.Fatalf("real capture refused: %v", err)
	}
	for path, mode := range map[string]string{"plain.txt": "100644", "tool.sh": "100755", "link": "120000"} {
		entry, tracked := generation.entry(path)
		if !tracked || !strings.HasPrefix(entry.Metadata, mode+" blob ") {
			t.Fatalf("entry %s = %+v/tracked=%v, want mode %s blob", path, entry, tracked, mode)
		}
	}
	entry, _ := generation.entry("plain.txt")
	if data, err := generation.blob(entry); err != nil || string(data) != "content\n" {
		t.Fatalf("blob read = %q, %v; want file content", data, err)
	}
	if _, err := generation.blob(treeEntry{Path: "plain.txt", Metadata: "absent"}); err == nil {
		t.Fatal("non-blob metadata read a blob")
	}
}

// Identity families receive generations from gateEvaluation. A family-local capture would
// preserve its hash while making source work grow with the number of families.
func TestProductionIdentityCapturesBelongToEvaluation(t *testing.T) {
	sources, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, source := range sources {
		if strings.HasSuffix(source, "_test.go") || source == treeSnapshotSourceFile {
			continue
		}
		body, err := os.ReadFile(source)
		if err != nil {
			t.Fatal(err)
		}
		if source != "evaluation.go" {
			for _, forbidden := range []string{"captureWorkingTree(", "captureProspectiveTree(", "workingTreeSource{", "prospectiveTreeSource{", "readTreeSnapshot("} {
				if strings.Contains(string(body), forbidden) {
					t.Fatalf("%s contains independent capture route %q", source, forbidden)
				}
			}
		}
		if strings.Contains(string(body), `"cat-file"`) && strings.Contains(string(body), `"blob"`) {
			t.Fatalf("%s reads a snapshot blob outside %s", source, treeSnapshotSourceFile)
		}
	}
}
