package git

import (
	"os"
	"path/filepath"
	"testing"
)

// noteWorktree adds a linked worktree of an existing repo and returns its resolved
// toplevel. It reuses newRepo as the package's one repository constructor.
func noteWorktree(t *testing.T, primary string) string {
	t.Helper()
	linked := filepath.Join(t.TempDir(), "linked")
	runGit(t, primary, "worktree", "add", "-q", "-b", "topic", linked)
	out, err := Output("-C", linked, "rev-parse", "--show-toplevel")
	if err != nil {
		t.Fatalf("rev-parse: %v", err)
	}
	return out
}

// TestLocalNoteRootTrackedKeepsCaller covers a path git does not ignore: the caller's
// checkout stays authoritative, from the primary and from a linked worktree alike.
func TestLocalNoteRootTrackedKeepsCaller(t *testing.T) {
	primary := newRepo(t)
	linked := noteWorktree(t, primary)
	for _, root := range []string{primary, linked} {
		got, ignored, err := LocalNoteRoot(root, "capture/note.md")
		if err != nil || ignored || got != root {
			t.Fatalf("root %s = (%q, %v, %v), want the caller's root, not ignored", root, got, ignored, err)
		}
	}
}

// TestLocalNoteRootIgnoredResolvesPrimary covers an ignored path: both checkouts
// resolve to the primary checkout's copy, so the note survives a worktree release.
func TestLocalNoteRootIgnoredResolvesPrimary(t *testing.T) {
	primary := newRepo(t)
	linked := noteWorktree(t, primary)
	for _, root := range []string{primary, linked} {
		if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte("capture/note.md\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	resolvedPrimary, err := Output("-C", primary, "rev-parse", "--show-toplevel")
	if err != nil {
		t.Fatalf("rev-parse: %v", err)
	}
	for _, root := range []string{resolvedPrimary, linked} {
		got, ignored, err := LocalNoteRoot(root, "capture/note.md")
		if err != nil || !ignored || got != resolvedPrimary {
			t.Fatalf("root %s = (%q, %v, %v), want the primary checkout, ignored", root, got, ignored, err)
		}
	}
}
