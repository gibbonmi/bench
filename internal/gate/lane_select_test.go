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

// laneChange is one raw-diff entry with the two regular-file modes, which is the shape
// an ordinary modification carries.
func laneChange(path string) ComposedChange {
	return ComposedChange{Status: "M", SrcMode: "100644", DstMode: "100644", Path: path}
}

// laneDeletion is the entry a removed file carries: a source mode and no destination.
func laneDeletion(path string) ComposedChange {
	return ComposedChange{Status: "D", SrcMode: "100644", DstMode: "000000", Path: path}
}

// laneSelectableFixture is an outcome fixture whose committed tree carries the two
// source directories the embed derivation walks. The kit's own checkout always carries
// them, and a tree that omits one makes the derivation report an absent directory
// instead of an empty embed list.
func laneSelectableFixture(t *testing.T) string {
	t.Helper()
	root := outcomeFixture(t)
	outcomeWrite(t, root, "cmd/bench/main.go", "package main\n\nfunc main() {}\n", 0o644)
	outcomeWrite(t, root, "internal/x/x.go", "package x\n", 0o644)
	outcomeGit(t, root, "add", "-A")
	outcomeGit(t, root, "commit", "-q", "-m", "sources")
	return root
}

// TestSelectLaneByClass is PL9 to PL13, PL15 to PL18, PL42, and PL47. Each row states
// the checks the kit lane runs for one change set, and the classes that selected them.
func TestSelectLaneByClass(t *testing.T) {
	every := []string{"gofmt", "prose", "vet", "build"}
	embed := []string{"internal/adopt/prepush.sh"}
	for _, tc := range []struct {
		name    string
		changes []ComposedChange
		embeds  []string
		checks  []string
		classes []string
	}{
		{
			name:    "PL9 a Go source file",
			changes: []ComposedChange{laneChange("internal/x/y.go")},
			checks:  []string{"gofmt", "vet", "build"},
			classes: []string{"go-source"},
		},
		{
			name:    "PL10 the module checksum file",
			changes: []ComposedChange{laneChange("go.sum")},
			checks:  []string{"vet", "build"},
			classes: []string{"go-build-input"},
		},
		{
			name:    "PL11 an embed target",
			changes: []ComposedChange{laneChange("internal/adopt/prepush.sh")},
			embeds:  embed,
			checks:  []string{"vet", "build"},
			classes: []string{"go-build-input"},
		},
		{
			name:    "PL12 a Markdown file",
			changes: []ComposedChange{laneChange("docs/note.md")},
			checks:  []string{"prose"},
			classes: []string{"markdown"},
		},
		{
			name:    "PL13 the prose exclusion list",
			changes: []ComposedChange{laneChange(".bench/prose-exclusions")},
			checks:  []string{"prose"},
			classes: []string{"prose-policy"},
		},
		{
			name:    "PL15 two classes take the union in declared order",
			changes: []ComposedChange{laneChange("a.go"), laneChange("b.md")},
			checks:  every,
			classes: []string{"go-source", "markdown"},
		},
		{
			name:    "PL16 a path no class claims",
			changes: []ComposedChange{laneChange("bin/bench.sh")},
			checks:  every,
			classes: []string{"unknown"},
		},
		{
			name:    "PL17 a link on the source side",
			changes: []ComposedChange{{Status: "T", SrcMode: "120000", DstMode: "100644", Path: "x.go"}},
			checks:  every,
			classes: []string{"unknown"},
		},
		{
			name:    "PL17 a link on the destination side",
			changes: []ComposedChange{{Status: "T", SrcMode: "100644", DstMode: "120000", Path: "x.go"}},
			checks:  every,
			classes: []string{"unknown"},
		},
		{
			name:    "PL18 a gitlink",
			changes: []ComposedChange{{Status: "M", SrcMode: "160000", DstMode: "160000", Path: "vendor/kit"}},
			checks:  every,
			classes: []string{"unknown"},
		},
		{
			name:    "PL42 a file under a glob embed pattern",
			changes: []ComposedChange{laneChange("templates/a.txt")},
			embeds:  []string{"templates/*"},
			checks:  every,
			classes: []string{"unknown"},
		},
		{
			name:    "PL47 a deleted Go source file",
			changes: []ComposedChange{laneDeletion("internal/x/y.go")},
			checks:  []string{"gofmt", "vet", "build"},
			classes: []string{"go-source"},
		},
		{
			name:    "PL47 a deleted embed target",
			changes: []ComposedChange{laneDeletion("internal/adopt/prepush.sh")},
			embeds:  embed,
			checks:  []string{"vet", "build"},
			classes: []string{"go-build-input"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			selected, classes := SelectLane(BenchkitLane("/repo", "/repo"), tc.changes, tc.embeds)
			if got := laneCheckNames(selected); !reflect.DeepEqual(got, tc.checks) {
				t.Errorf("selected checks = %v, want %v", got, tc.checks)
			}
			if !reflect.DeepEqual(classes, tc.classes) {
				t.Errorf("classes = %v, want %v", classes, tc.classes)
			}
		})
	}
}

// TestSelectLaneReturnsADeclaredSubsequence is PL38. The lane below declares two of the
// checks the classes name, so a selection that added a class's own check name would
// hand the run a check the lane never declared.
func TestSelectLaneReturnsADeclaredSubsequence(t *testing.T) {
	declared := []Phase{
		{Name: "prose", Argv: []string{"true"}},
		{Name: "build", Argv: []string{"true"}},
	}
	for _, changes := range [][]ComposedChange{
		nil,
		{laneChange("a.go")},
		{laneChange("go.mod")},
		{laneChange("a.md")},
		{laneChange("a.go"), laneChange("a.md"), laneChange(".bench/prose-exclusions")},
		{laneChange("bin/x.sh")},
		{laneChange("a.go"), laneChange("bin/x.sh")},
	} {
		selected, _ := SelectLane(declared, changes, nil)
		next := 0
		for _, check := range selected {
			for next < len(declared) && declared[next].Name != check.Name {
				next++
			}
			if next == len(declared) {
				t.Fatalf("selection %v for %v is no subsequence of %v",
					laneCheckNames(selected), changes, laneCheckNames(declared))
			}
			next++
		}
	}
}
