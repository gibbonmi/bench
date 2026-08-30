package gate

// The composed change list: the raw diff between the base tree and the tree a commit
// composes. Each row here drives the real derivation against a real repository, because
// the framing and the mode bytes are Git's own and no fake reproduces them.

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

// composedFixture is a repository whose base commit holds the given files, and it
// answers that commit. A row commits its graded side after it.
func composedFixture(t *testing.T, files map[string]string) (root, base string) {
	t.Helper()
	root = gittest.RepoOnBranch(t, "main")
	commitComposed(t, root, "base", files)
	return root, composedRev(t, root, "HEAD^{commit}")
}

// commitComposed writes and commits one set of files, so a row states the two trees its
// derivation is measured between.
func commitComposed(t *testing.T, root, message string, files map[string]string) {
	t.Helper()
	for name, body := range files {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	stageAndCommitComposed(t, root, message)
}

func stageAndCommitComposed(t *testing.T, root, message string) {
	t.Helper()
	if out, err := benchgit.Raw("-C", root, "add", "-A"); err != nil {
		t.Fatalf("stage %s: %v\n%s", message, err, out)
	}
	if out, err := benchgit.Raw("-C", root, "commit", "-q", "-m", message); err != nil {
		t.Fatalf("commit %s: %v\n%s", message, err, out)
	}
}

func composedRev(t *testing.T, root, revision string) string {
	t.Helper()
	resolved, err := benchgit.Output("-C", root, "rev-parse", revision)
	if err != nil {
		t.Fatal(err)
	}
	return resolved
}

// PL1: the list is the tree-to-tree diff, so a commit that names the directory `docs`
// reaches its changed files. A list built from the named paths holds `docs` and no file.
func TestComposedChangesExpandsANamedDirectory(t *testing.T) {
	root, base := composedFixture(t, map[string]string{"docs/a.md": "before\n", "docs/b.go": "package docs\n"})
	commitComposed(t, root, "graded", map[string]string{"docs/a.md": "after\n", "docs/b.go": "package docs // after\n"})
	tree := composedRev(t, root, "HEAD^{tree}")

	changes, err := ComposedChanges(root, base, tree)
	if err != nil {
		t.Fatal(err)
	}
	want := []ComposedChange{
		{Status: "M", SrcMode: "100644", DstMode: "100644", Path: "docs/a.md"},
		{Status: "M", SrcMode: "100644", DstMode: "100644", Path: "docs/b.go"},
	}
	if !reflect.DeepEqual(changes, want) {
		t.Fatalf("changes = %+v, want the two files under the directory %+v", changes, want)
	}
}

// PL2: rename detection is off, so a rename is one deletion and one addition. An `R`
// entry names two paths in one row, and no class reader expects that shape.
func TestComposedChangesRepresentsARenameAsDeletionAndAddition(t *testing.T) {
	root, base := composedFixture(t, map[string]string{"docs/old.md": "prose that is long enough to match\n"})
	if out, err := benchgit.Raw("-C", root, "mv", "docs/old.md", "docs/new.md"); err != nil {
		t.Fatalf("rename the file: %v\n%s", err, out)
	}
	stageAndCommitComposed(t, root, "graded")
	tree := composedRev(t, root, "HEAD^{tree}")

	changes, err := ComposedChanges(root, base, tree)
	if err != nil {
		t.Fatal(err)
	}
	want := []ComposedChange{
		{Status: "A", SrcMode: "000000", DstMode: "100644", Path: "docs/new.md"},
		{Status: "D", SrcMode: "100644", DstMode: "000000", Path: "docs/old.md"},
	}
	if !reflect.DeepEqual(changes, want) {
		t.Fatalf("changes = %+v, want the deletion and the addition %+v", changes, want)
	}
}

// PL3: the destination mode reaches the reader, so a symbolic link named `link.go` is
// not read as Go source. A name-only list carries the suffix and loses the mode.
func TestComposedChangesCarriesTheSymlinkMode(t *testing.T) {
	root, base := composedFixture(t, map[string]string{"target.go": "package fixture\n"})
	if err := os.Symlink("target.go", filepath.Join(root, "link.go")); err != nil {
		capability.Capability(t, capability.Symlink, "symlinks unavailable: "+err.Error())
	}
	stageAndCommitComposed(t, root, "graded")
	tree := composedRev(t, root, "HEAD^{tree}")

	changes, err := ComposedChanges(root, base, tree)
	if err != nil {
		t.Fatal(err)
	}
	want := []ComposedChange{{Status: "A", SrcMode: "000000", DstMode: "120000", Path: "link.go"}}
	if !reflect.DeepEqual(changes, want) {
		t.Fatalf("changes = %+v, want the link with its own mode %+v", changes, want)
	}
}
