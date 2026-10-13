package probe

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// deletionRestore is the restore of a --delete probe. A deletion replaces nothing, so that
// restore is the first replacement of the subject.
var deletionRestore = durablefiletest.Fault{Destination: "subject", Ordinal: 1}

// deleteArgs is the --delete invocation every row below runs.
var deleteArgs = []string{"clamp.go", "--delete", "--package", "./", "--run", "^TestClampNegative$"}

// executableSubject gives the subject an executable mode, so a restore that recreates the
// file at a default mode shows.
func executableSubject(t *testing.T, f *fixture) {
	t.Helper()
	if err := os.Chmod(f.subject, 0o750); err != nil {
		t.Fatal(err)
	}
}

// PS9, PS10: the subject is absent while the mutated run runs, and the restore brings back
// its start bytes and its start mode.
func TestDeleteRemovesTheSubjectForTheRun(t *testing.T) {
	f := newFixture(t)
	executableSubject(t, f)
	installStubGo(t, f, cannedFailure("caught"), 1, "[ ! -e "+sanitize.ShellQuote(f.subject)+" ] || exit 99\n")
	out, code := runProbe(t, deleteArgs...)
	requireRow(t, out, code, 0, probeRow(t, "bit", "clamp.go", "delete", "failed", 1, "yes"))
	durablefiletest.RequireFile(t, f.subject, 0o750, clampSource)
	requireHomeEmpty(t, f)
}

// modeless is a temporary file whose mode change does nothing, so a restore through it
// leaves the bytes right and the mode wrong.
type modeless struct{ durablefile.File }

func (modeless) Chmod(fs.FileMode) error { return nil }

// failRemoval makes the subject's removal fail after it removed the file, which is the
// failure that can leave the subject absent.
func failRemoval(f *fixture, replacer durablefile.Replacer, cause error) durablefile.Replacer {
	remove := replacer.Remove
	replacer.Remove = func(name string) error {
		if name != f.subject {
			return remove(name)
		}
		_ = remove(name)
		return cause
	}
	return replacer
}

// PS11: a failed restore after --delete reads restore-failed and names the kept copy, and a
// failed removal takes the recovery of a published replace failure: the restore runs, and
// a failed restore keeps the copy at exit 2.
func TestDeleteReportsARestoreFailure(t *testing.T) {
	removalCause := errors.New("injected removal failure")
	t.Run("restore proves the mode", func(t *testing.T) {
		f := newFixture(t)
		executableSubject(t, f)
		replacer := durablefile.Native()
		create := replacer.CreateTemp
		replacer.CreateTemp = func(dir, pattern string) (durablefile.File, error) {
			file, err := create(dir, pattern)
			if err != nil {
				return nil, err
			}
			return modeless{file}, nil
		}
		out, code := probeWith(t, f, replacer, deleteArgs...)
		requireRow(t, out, code, 2, probeRow(t, "restore-failed", "clamp.go", "delete", "failed", 1, "no"))
		if kept := preservedCopy(t, f); !strings.Contains(out, kept) {
			t.Fatalf("stdout = %q, want it to name the copy %s", out, kept)
		}
	})
	t.Run("failed removal and restore", func(t *testing.T) {
		f := newFixture(t)
		router := subjectFaults()
		restoreCause := errors.New("injected restore rename failure")
		router.Fail[deletionRestore] = map[string]error{"rename": restoreCause}
		out, code := probeWith(t, f, failRemoval(f, router.Replacer(), removalCause), deleteArgs...)
		if code != 2 {
			t.Fatalf("exit = %d, want 2\n%s", code, out)
		}
		kept := preservedCopy(t, f)
		for _, want := range []string{kept, removalCause.Error(), restoreCause.Error()} {
			if !strings.Contains(out, want) {
				t.Fatalf("stdout = %q, want it to name %q", out, want)
			}
		}
		durablefiletest.RequireBytes(t, kept, clampSource)
		requireBaselineOnly(t, f)
	})
	t.Run("failed removal restored", func(t *testing.T) {
		f := newFixture(t)
		executableSubject(t, f)
		out, code := probeWith(t, f, failRemoval(f, durablefile.Native(), removalCause), deleteArgs...)
		if want := "error: probe mutation failed — " + removalCause.Error() + "\n"; out != want || code != 1 {
			t.Fatalf("probe = (%q, %d), want (%q, 1)", out, code, want)
		}
		durablefiletest.RequireFile(t, f.subject, 0o750, clampSource)
		requireHomeEmpty(t, f)
		requireBaselineOnly(t, f)
	})
}

// PS12: an omission whose result has zero bytes refuses with the line that names --delete,
// for --omit-file and for a combined mutant alike.
func TestOmissionThatEmptiesTheSubjectRefuses(t *testing.T) {
	head, tail, _ := strings.Cut(clampSource, "func ")
	for _, tc := range []struct {
		name string
		args func(t *testing.T) []string
	}{
		{"omit-file", func(t *testing.T) []string {
			path := filepath.Join(scratchDir(t), "whole")
			writeFixtureFile(t, path, clampSource, 0o644)
			return []string{"--omit-file", path}
		}},
		{"combined", func(*testing.T) []string { return []string{"--omit", head, "--omit", "func " + tail} }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := refusalFixture(t)
			out, code := runProbe(t, append(append([]string{"clamp.go"}, tc.args(t)...), "--package", "./")...)
			want := "error: probe mutation empties the subject — use --delete to remove clamp.go for the run\n"
			requireRefused(t, f, out, code, want, 1)
		})
	}
}
