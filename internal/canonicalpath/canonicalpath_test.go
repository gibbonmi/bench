package canonicalpath

import (
	"os"
	"path/filepath"
	"testing"
)

// TestOperandJoinsTheWorkingDirectoryAndSpellsTheDisplay pins the one operand rule both
// `bench anchors` and `bench probe` read. The rows are built from filepath.Join rather
// than from literals, so the case runs on a platform whose separator is not "/".
func TestOperandJoinsTheWorkingDirectoryAndSpellsTheDisplay(t *testing.T) {
	root := filepath.Join(string(filepath.Separator), "repo")
	cwd := filepath.Join(root, "internal", "probe")
	cases := []struct {
		name        string
		root        string
		cwd         string
		arg         string
		wantPath    string
		wantDisplay string
	}{
		{
			name:        "an absolute operand stays",
			root:        root,
			cwd:         cwd,
			arg:         filepath.Join(root, "cmd", "bench", "main.go"),
			wantPath:    filepath.Join(root, "cmd", "bench", "main.go"),
			wantDisplay: "cmd/bench/main.go",
		},
		{
			name:        "a relative operand joins onto the working directory",
			root:        root,
			cwd:         cwd,
			arg:         filepath.Join("subject.go"),
			wantPath:    filepath.Join(cwd, "subject.go"),
			wantDisplay: "internal/probe/subject.go",
		},
		{
			name:        "a dot-dot operand cleans",
			root:        root,
			cwd:         cwd,
			arg:         filepath.Join("..", "toon", "toon.go"),
			wantPath:    filepath.Join(root, "internal", "toon", "toon.go"),
			wantDisplay: "internal/toon/toon.go",
		},
		{
			name:        "an operand outside the root spells its way out",
			root:        root,
			cwd:         cwd,
			arg:         filepath.Join(string(filepath.Separator), "elsewhere", "file.go"),
			wantPath:    filepath.Join(string(filepath.Separator), "elsewhere", "file.go"),
			wantDisplay: "../elsewhere/file.go",
		},
		{
			name:        "a root that cannot relativize keeps the cleaned operand",
			root:        filepath.Join("relative", "root"),
			cwd:         cwd,
			arg:         filepath.Join("..", "probe", "subject.go"),
			wantPath:    filepath.Join(root, "internal", "probe", "subject.go"),
			wantDisplay: "../probe/subject.go",
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			path, display := Operand(testCase.root, testCase.cwd, testCase.arg)
			if path != testCase.wantPath {
				t.Fatalf("Operand path = %q, want %q", path, testCase.wantPath)
			}
			if display != testCase.wantDisplay {
				t.Fatalf("Operand display = %q, want %q", display, testCase.wantDisplay)
			}
		})
	}
}

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
