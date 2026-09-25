package responsebound

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/poolkey"
)

// spillPath answers the path that the spill line of one response names.
func spillPath(t *testing.T, response string) string {
	t.Helper()
	for _, line := range strings.SplitAfter(response, "\n") {
		if _, path, found := strings.Cut(strings.TrimSuffix(line, "}\n"), ",path="); found && strings.HasPrefix(line, "spilled{") {
			return path
		}
	}
	t.Fatalf("response = %q, want a spill line", response)
	return ""
}

// spillOnce sends one over-bound response through a new owner and answers its spill path.
func spillOnce(t *testing.T, root func() string) string {
	t.Helper()
	var sink bytes.Buffer
	return spillPath(t, respondWith(t, New(&sink, &sink, root), &sink, stdoutLines(numbered(1, 11)...)))
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
			home := privateHome(t)
			var paths []string
			for range 65 {
				paths = append(paths, spillOnce(t, row.root))
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
