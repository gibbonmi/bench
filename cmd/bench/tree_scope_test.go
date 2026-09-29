package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/benchhome"
)

// TT6 and TT7: a repository-scoped verb refuses a tree target as its first argument, before
// the verb runs. The expected usage line is authored apart from the dispatcher, so a
// refusal that names another command or another argument fails the exact match.
func TestRepositoryVerbRefusesTreeTarget(t *testing.T) {
	t.Setenv(benchhome.Env, t.TempDir())
	root := newAXIEnvelopeRepo(t)
	// An ignored inbox lets `bench idea` park its line in this primary checkout, so a verb
	// that runs before the refusal leaves the file behind.
	writeAXIFixture(t, filepath.Join(root, ".gitignore"), "capture/\n")

	t.Run("version", func(t *testing.T) {
		result := runAXICommandAt(t, root, []string{"version", "--in", "primary"})
		if want := "usage: bench version (unknown argument: --in)\n"; result.code != 2 || result.stdout != want {
			t.Fatalf("version --in primary = (%d, %q, %q), want (2, %q)", result.code, result.stdout, result.stderr, want)
		}
	})
	t.Run("idea", func(t *testing.T) {
		result := runAXICommandAt(t, root, []string{"idea", "--in", "primary", "x"})
		if result.code != 2 {
			t.Fatalf("idea --in primary x = (%d, %q, %q), want exit 2", result.code, result.stdout, result.stderr)
		}
		if _, err := os.Stat(filepath.Join(root, "capture", "IDEAS.md")); !errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("idea --in primary x left capture/IDEAS.md: %v", err)
		}
	})
}
