package compatibilitytest

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// HomeEntry retains the observable bytes and mode of one personal-home entry.
type HomeEntry struct {
	Mode    os.FileMode
	Content string
}

// SnapshotHome captures one private fixture home without following symbolic links.
func SnapshotHome(t testing.TB, home string) map[string]HomeEntry {
	t.Helper()
	entries := map[string]HomeEntry{}
	err := filepath.WalkDir(home, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(home, path)
		if err != nil {
			return err
		}
		var content string
		switch {
		case info.Mode().IsRegular():
			data, readErr := os.ReadFile(path)
			content, err = string(data), readErr
		case info.Mode()&os.ModeSymlink != 0:
			content, err = os.Readlink(path)
		}
		entries[rel] = HomeEntry{info.Mode(), content}
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return entries
}

// RemovePreserving removes one regular fixture file and restores its bytes and mode at cleanup.
func RemovePreserving(t testing.TB, path string) {
	t.Helper()
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return
	}
	if err != nil {
		t.Fatal(err)
	}
	if !info.Mode().IsRegular() {
		t.Fatalf("fixture preservation requires a regular file: %s", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.WriteFile(path, data, info.Mode()); err != nil {
			t.Error(err)
		}
		if err := os.Chmod(path, info.Mode()); err != nil {
			t.Error(err)
		}
	})
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
}
