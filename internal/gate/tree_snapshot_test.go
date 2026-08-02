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
