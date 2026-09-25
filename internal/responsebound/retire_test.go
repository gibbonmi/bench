package responsebound

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/poolkey"
	"github.com/gibbonmi/bench/internal/responsebound/responseboundtest"
)

// spillOnce sends one over-bound response through a new owner and answers its spill path.
func spillOnce(t *testing.T, home string, root func() string) string {
	t.Helper()
	var sink bytes.Buffer
	return responseboundtest.Path(t, respondWith(t, New(home, &sink, &sink, root, false), &sink, stdoutLines(numbered(1, 11)...)))
}

// BO17: the 65th spill in a `primary` scope leaves the newest 64 files. The counts are
// authored apart from the owner, because the spec fixes them. The scope below the `none`
// repository key and the scope below a repository key both keep the limit.
func TestPrimarySpillsKeepNewest(t *testing.T) {
	repository := t.TempDir()
	for _, row := range []struct {
		name string
		root func() string
		key  string
	}{
		{"none", outsideRepository, noRepository},
		{"repository", func() string { return repository }, poolkey.Key(repository)},
	} {
		t.Run(row.name, func(t *testing.T) {
			home := t.TempDir()
			var paths []string
			for range 65 {
				paths = append(paths, spillOnce(t, home, row.root))
			}
			scope := filepath.Join(home, storeDirName, row.key, primaryScope)
			entries, err := os.ReadDir(scope)
			if err != nil || len(entries) != 64 {
				t.Fatalf("primary scope entries = %d (%v), want 64", len(entries), err)
			}
			if _, err := os.Stat(paths[0]); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("oldest spill %s = %v, want removed", paths[0], err)
			}
			for _, path := range paths[1:] {
				if filepath.Dir(path) != scope {
					t.Fatalf("spill %s, want it in the primary scope %s", path, scope)
				}
				if _, err := os.Stat(path); err != nil {
					t.Fatalf("newer spill %s = %v, want kept", path, err)
				}
			}
		})
	}
}

// The spill drop refuses an identifier that is not an assignment id, as the census drop
// does, so no operand composes the path that the drop removes.
func TestDropRefusesNonAssignmentID(t *testing.T) {
	home := t.TempDir()
	root := t.TempDir()
	scope := filepath.Join(home, storeDirName, poolkey.Key(root), primaryScope)
	if err := os.MkdirAll(scope, 0o700); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"..", primaryScope, "", strings.Repeat("g", 32)} {
		if err := Drop(home, root, id); err == nil {
			t.Fatalf("Drop(%q) = nil, want a refusal", id)
		}
	}
	if _, err := os.Stat(scope); err != nil {
		t.Fatalf("a refused drop removed the primary scope: %v", err)
	}
}

// The spill drop never follows a symlink at a store level. A symlink at `responses/` or at
// `responses/<repo-key>/` that names a directory outside the home would otherwise send the
// removal to that directory.
func TestDropRefusesSymlinkedStore(t *testing.T) {
	root := t.TempDir()
	key, assignment := poolkey.Key(root), strings.Repeat("a", 32)
	for _, row := range []struct {
		name  string
		level string
		held  string
	}{
		{"store", storeDirName, key},
		{"repository key", filepath.Join(storeDirName, key), ""},
	} {
		t.Run(row.name, func(t *testing.T) {
			home, outside := t.TempDir(), t.TempDir()
			target := filepath.Join(outside, row.held, assignment)
			if err := os.MkdirAll(target, 0o700); err != nil {
				t.Fatal(err)
			}
			link := filepath.Join(home, row.level)
			if err := os.MkdirAll(filepath.Dir(link), 0o700); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(outside, link); err != nil {
				t.Fatal(err)
			}
			_ = Drop(home, root, assignment)
			if _, err := os.Stat(target); err != nil {
				t.Fatalf("a drop through a symlink at %s removed %s: %v", row.level, target, err)
			}
		})
	}
}
