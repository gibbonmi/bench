package canary

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRestoreMutationFixtureReinstatesBaseAndRemovesOverlay(t *testing.T) {
	source := t.TempDir()
	if err := os.WriteFile(filepath.Join(source, "owned.txt"), []byte("clean\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	fixture := t.TempDir()
	files := filepath.Join(fixture, filesDirName)
	if err := os.MkdirAll(files, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range map[string]string{"owned.txt": "overlay\n", "added.txt": "added\n"} {
		if err := os.WriteFile(filepath.Join(files, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(fixture, "BASE"), []byte("owned.txt\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	dst := t.TempDir()
	if err := materializeMutationFixture(source, fixture, dst); err != nil {
		t.Fatalf("materialize fixture: %v", err)
	}
	if err := RestoreMutationFixture(source, fixture, dst); err != nil {
		t.Fatalf("restore fixture: %v", err)
	}
	if got, err := os.ReadFile(filepath.Join(dst, "owned.txt")); err != nil || string(got) != "clean\n" {
		t.Fatalf("restored base = %q, %v", got, err)
	}
	if _, err := os.Stat(filepath.Join(dst, "added.txt")); !os.IsNotExist(err) {
		t.Fatalf("overlay remains after restore: %v", err)
	}
}

func TestRestoreMutationFixtureUsesSourceForOverlayPaths(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct{ fixture, rel string }{
		{filepath.Join(root, "tests", "canary", "load-validity-metadata", "codex-hooks-broken"), filepath.Join(".codex", "hooks.json")},
	} {
		dst := t.TempDir()
		if err := MaterializeMutationFixture(root, test.fixture, dst); err != nil {
			t.Fatal(err)
		}
		if err := RestoreMutationFixture(root, test.fixture, dst); err != nil {
			t.Fatal(err)
		}
		want, sourceErr := os.ReadFile(filepath.Join(root, test.rel))
		got, restoredErr := os.ReadFile(filepath.Join(dst, test.rel))
		if sourceErr != nil || restoredErr != nil || string(got) != string(want) {
			t.Fatalf("%s restored %s differs from source: source=%v restored=%v", test.fixture, test.rel, sourceErr, restoredErr)
		}
	}
}
