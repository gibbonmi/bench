package visibletree_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/visibletree"
)

func writeTreeFixture(t *testing.T, root string, rels ...string) {
	t.Helper()
	for _, rel := range rels {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(rel+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestFilesListsTrackedAndUntrackedButNotIgnored(t *testing.T) {
	root := gittest.RepoIgnoring(t, ".logs/", "local-skill/")
	writeTreeFixture(t, root, "tracked.go", "gone.go")
	gittest.Output(t, root, "add", "tracked.go", "gone.go")
	gittest.Output(t, root, "commit", "-qm", "tracked")
	if err := os.Remove(filepath.Join(root, "gone.go")); err != nil {
		t.Fatal(err)
	}
	writeTreeFixture(t, root, "new/untracked.go", ".logs/copy/tracked.go", "local-skill/SKILL.md")

	got, err := visibletree.Files(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{".gitignore", "new/untracked.go", "tracked.go"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Files = %v, want %v (no ignored path, no deleted tracked file)", got, want)
	}

	ignored, err := visibletree.Ignored(root, ".logs/copy", "local-skill", "tracked.go", "new/untracked.go")
	if err != nil {
		t.Fatal(err)
	}
	if want := map[string]bool{".logs/copy": true, "local-skill": true}; !reflect.DeepEqual(ignored, want) {
		t.Fatalf("Ignored = %v, want %v", ignored, want)
	}
}

func TestFilesOutsideARepositoryIgnoresNothing(t *testing.T) {
	root := t.TempDir()
	writeTreeFixture(t, root, "a.go", ".logs/b.go", ".git/config")

	got, err := visibletree.Files(root)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{".logs/b.go", "a.go"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Files = %v, want %v", got, want)
	}
	ignored, err := visibletree.Ignored(root, ".logs")
	if err != nil || len(ignored) != 0 {
		t.Fatalf("Ignored = %v, %v, want none outside a repository", ignored, err)
	}
}
