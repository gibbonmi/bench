package testrepo_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/testrepo"
)

// TestCopyTreeKeepsModes uses modes that os.CopyFS does not keep under the usual umask, and
// a symbolic link that sorts after its target, so a chmod through the link changes the target.
func TestCopyTreeKeepsModes(t *testing.T) {
	src, dst := t.TempDir(), filepath.Join(t.TempDir(), "copy")
	for _, file := range []struct {
		name string
		mode os.FileMode
	}{{"a-private", 0o600}, {"dir/read-only", 0o444}, {"dir/tool", 0o750}} {
		path := filepath.Join(src, filepath.FromSlash(file.name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(file.name), file.mode); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Chmod(filepath.Join(src, "dir"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("a-private", filepath.Join(src, "z-link")); err != nil {
		t.Fatal(err)
	}
	if err := testrepo.CopyTree(dst, src); err != nil {
		t.Fatal(err)
	}
	err := filepath.WalkDir(src, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		want, err := os.Lstat(path)
		if err != nil {
			return err
		}
		got, err := os.Lstat(filepath.Join(dst, rel))
		if err != nil {
			return err
		}
		if got.Mode() != want.Mode() {
			t.Errorf("%s mode = %v, want %v", rel, got.Mode(), want.Mode())
		}
		if want.Mode().IsRegular() {
			body, err := os.ReadFile(filepath.Join(dst, rel))
			if err != nil {
				return err
			}
			if string(body) != filepath.ToSlash(rel) {
				t.Errorf("%s bytes = %q, want %q", rel, body, filepath.ToSlash(rel))
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if target, err := os.Readlink(filepath.Join(dst, "z-link")); err != nil || target != "a-private" {
		t.Fatalf("z-link target = %q, %v, want a-private", target, err)
	}
}
