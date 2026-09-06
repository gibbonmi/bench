package canonicalpath

import (
	"os"
	"path/filepath"
	"testing"
)

// TestResolveFollowsALinkAndKeepsAnAbsentPath pins the two halves of the derivation the
// six former copies each carried. The link row compares two spellings of one directory
// rather than a literal, because the temp root itself can sit under a link. The absent row
// compares a literal, because nothing under it resolves.
func TestResolveFollowsALinkAndKeepsAnAbsentPath(t *testing.T) {
	t.Run("a link resolves to its target's spelling", func(t *testing.T) {
		root := t.TempDir()
		target := filepath.Join(root, "target")
		if err := os.Mkdir(target, 0o755); err != nil {
			t.Fatal(err)
		}
		link := filepath.Join(root, "link")
		if err := os.Symlink(target, link); err != nil {
			t.Fatal(err)
		}
		viaLink, err := Resolve(link)
		if err != nil {
			t.Fatal(err)
		}
		direct, err := Resolve(target)
		if err != nil {
			t.Fatal(err)
		}
		if viaLink != direct {
			t.Fatalf("Resolve through the link = %q, want the target's spelling %q", viaLink, direct)
		}
	})

	t.Run("a dot-dot after a symlink resolves physically", func(t *testing.T) {
		base := t.TempDir()
		physical := filepath.Join(base, "physical")
		child := filepath.Join(physical, "child")
		if err := os.MkdirAll(child, 0o755); err != nil {
			t.Fatal(err)
		}
		jump := filepath.Join(base, "jump")
		if err := os.Symlink(child, jump); err != nil {
			t.Fatal(err)
		}
		// filepath.Join would clean "jump/.." away before Resolve ever saw it, so the
		// root is built by concatenation to keep the symlinked component in place.
		root := jump + string(filepath.Separator) + ".."
		want, err := Resolve(physical)
		if err != nil {
			t.Fatal(err)
		}
		got, err := Resolve(root)
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("Resolve(%q) = %q, want the physical directory %q", root, got, want)
		}
	})

	t.Run("an absent path keeps its absolute spelling", func(t *testing.T) {
		absent := filepath.Join(t.TempDir(), "absent", "child.txt")
		got, err := Resolve(absent)
		if err != nil {
			t.Fatalf("Resolve of an absent path refused: %v", err)
		}
		if got != absent {
			t.Fatalf("Resolve of an absent path = %q, want %q", got, absent)
		}
	})
}
