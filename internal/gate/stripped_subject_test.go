package gate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStrippedIdentityIgnoresAllowlistedEdit(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func(*testing.T, string)
	}{
		{"declared root file edited", func(t *testing.T, root string) { writeFixture(t, root, "ROADMAP.md", "edited\n") }},
		{"shift scratch edited", func(t *testing.T, root string) { writeFixture(t, root, ".bench-notes.md", "edited\n") }},
		{"capture surface edited", func(t *testing.T, root string) { writeFixture(t, root, "capture/learnings.md", "edited\n") }},
		{"capture surface added", func(t *testing.T, root string) { writeFixture(t, root, "capture/retros/new.md", "retro\n") }},
		{"spec surface edited", func(t *testing.T, root string) { writeFixture(t, root, "specs/thing/spec.md", "edited\n") }},
		{"spec surface added", func(t *testing.T, root string) { writeFixture(t, root, "specs/thing/tickets/one.md", "ticket\n") }},
		{"capture surface deleted", func(t *testing.T, root string) { removeFixture(t, root, "capture/session-handoff.md") }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := newIdentityFixture(t)
			whole, stripped := fixtureIdentities(t, root)
			tc.mutate(t, root)
			mutatedWhole, mutatedStripped := fixtureIdentities(t, root)

			if mutatedStripped != stripped {
				t.Errorf("stripped identity moved on an allowlisted edit: %s -> %s", stripped, mutatedStripped)
			}
			if mutatedWhole == whole {
				t.Errorf("whole-tree identity held at %s, so the fixture edit changed nothing", whole)
			}
		})
	}
}

// TestStrippedIdentityMovesOnUnlistedEdit covers the dangerous direction: an identity that
// holds still for a tree that really changed reuses evidence for ungraded work. The cases
// that reach beyond the declaration — the root the declared directories sit in, and the
// names a prefix match swallows — are the strips this rules out.
func TestStrippedIdentityMovesOnUnlistedEdit(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func(*testing.T, string)
	}{
		{"root file beside the declared entries", func(t *testing.T, root string) { writeFixture(t, root, "README.md", "edited\n") }},
		{"source below the root", func(t *testing.T, root string) {
			writeFixture(t, root, "internal/gate/engine.go", "package gate // edited\n")
		}},
		{"directory sharing a declared prefix", func(t *testing.T, root string) { writeFixture(t, root, "capture-old/note.md", "edited\n") }},
		{"file whose name extends a declared one", func(t *testing.T, root string) { writeFixture(t, root, "ROADMAP.md.bak", "edited\n") }},
		{"root file whose name extends a declared directory", func(t *testing.T, root string) { writeFixture(t, root, "specsheet.md", "edited\n") }},
		{"file added outside the declaration", func(t *testing.T, root string) { writeFixture(t, root, "internal/gate/added.go", "package gate\n") }},
		{"file deleted outside the declaration", func(t *testing.T, root string) { removeFixture(t, root, "NOTES.md") }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := newIdentityFixture(t)
			_, stripped := fixtureIdentities(t, root)
			tc.mutate(t, root)
			if _, mutatedStripped := fixtureIdentities(t, root); mutatedStripped == stripped {
				t.Errorf("stripped identity held at %s across an edit outside the allowlist", stripped)
			}
		})
	}
}

// newIdentityFixture builds a repository carrying one file of every shape the identities
// have to tell apart: declared files and directories, the near-miss names a prefix match
// would swallow, and ordinary sources at the root and below it.
func newIdentityFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	gitRun(t, root, "init", "-q")
	for path, content := range map[string]string{
		"README.md":                  "readme\n",
		"NOTES.md":                   "notes\n",
		"internal/gate/engine.go":    "package gate\n",
		"ROADMAP.md":                 "roadmap\n",
		"ROADMAP.md.bak":             "backup\n",
		".bench-notes.md":            "scratch\n",
		"capture/learnings.md":       "learnings\n",
		"capture/session-handoff.md": "handoff\n",
		"capture-old/note.md":        "old note\n",
		"specs/thing/spec.md":        "spec\n",
		"specsheet.md":               "sheet\n",
	} {
		writeFixture(t, root, path, content)
	}
	return root
}

func writeFixture(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir for %s: %v", rel, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", rel, err)
	}
}

func removeFixture(t *testing.T, root, rel string) {
	t.Helper()
	if err := os.Remove(filepath.Join(root, filepath.FromSlash(rel))); err != nil {
		t.Fatalf("remove %s: %v", rel, err)
	}
}

func fixtureIdentities(t *testing.T, root string) (whole, stripped string) {
	t.Helper()
	full, err := buildSubject(root)
	if err != nil {
		t.Fatalf("whole-tree subject: %v", err)
	}
	reduced, err := buildStrippedSubject(root)
	if err != nil {
		t.Fatalf("stripped subject: %v", err)
	}
	if full.Oracle == "" || reduced.Oracle == "" {
		t.Fatal("empty identity")
	}
	return full.Oracle, reduced.Oracle
}
