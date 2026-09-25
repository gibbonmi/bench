package responsebound

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/benchhome"
)

// storeLevels answers each spill store directory below home, from the store root to the
// scope directory of a process outside any repository.
func storeLevels(home string) []string {
	root := filepath.Join(home, storeDirName)
	repository := filepath.Join(root, noRepository)
	return []string{root, repository, filepath.Join(repository, primaryScope)}
}

func storeDir(home string) string {
	levels := storeLevels(home)
	return levels[len(levels)-1]
}

// BO14: a new spill directory is private to the owner, and so is the file.
func TestSpillStorePrivateModes(t *testing.T) {
	home := privateHome(t)
	respond(t, stdoutLines(numbered(1, 11)...))
	for _, dir := range storeLevels(home) {
		info, err := os.Lstat(dir)
		if err != nil || !info.IsDir() || info.Mode().Perm() != 0o700 {
			t.Errorf("directory %s = (%v, %v), want a directory of mode 0700", dir, info, err)
		}
	}
	info, err := os.Lstat(onlySpill(t, home))
	if err != nil || !info.Mode().IsRegular() || info.Mode().Perm() != 0o600 {
		t.Fatalf("spill file = (%v, %v), want a regular file of mode 0600", info, err)
	}
}

// assertCreateFailure grades the create-failure route: the complete output, then one
// line-safe spill-failed line.
func assertCreateFailure(t *testing.T, got, output string) {
	t.Helper()
	rest, found := strings.CutPrefix(got, output)
	if !found {
		t.Fatalf("response = %q, want the complete output %q first", got, output)
	}
	if !strings.HasPrefix(rest, "spill-failed{reason=") || !strings.HasSuffix(rest, "}\n") || strings.Count(rest, "\n") != 1 {
		t.Fatalf("response tail = %q, want one spill-failed{reason=<reason>} line", rest)
	}
}

// atEachStoreLevel runs check once for each store level, in a fresh Bench home whose
// parent directories exist. check plants its entry at the level path.
func atEachStoreLevel(t *testing.T, check func(t *testing.T, level string)) {
	for i, name := range storeLevels("") {
		t.Run(name, func(t *testing.T) {
			level := storeLevels(privateHome(t))[i]
			if err := os.MkdirAll(filepath.Dir(level), 0o700); err != nil {
				t.Fatal(err)
			}
			check(t, level)
		})
	}
}

// BO15: a symlink at any store level takes the create-failure route, and nothing is
// written where the symlink points.
func TestSpillStoreRefusesSymlink(t *testing.T) {
	atEachStoreLevel(t, func(t *testing.T, level string) {
		elsewhere := t.TempDir()
		if err := os.Symlink(elsewhere, level); err != nil {
			t.Fatal(err)
		}
		writes := stdoutLines(numbered(1, 11)...)
		assertCreateFailure(t, respond(t, writes), joined(writes))
		if entries, err := os.ReadDir(elsewhere); err != nil || len(entries) != 0 {
			t.Fatalf("symlink target entries = %v (%v), want none", entries, err)
		}
	})
}

// BO18: a regular file at any store level gives the complete output and then the
// spill-failed line.
func TestOwnerCreateFailureKeepsOutput(t *testing.T) {
	atEachStoreLevel(t, func(t *testing.T, level string) {
		if err := os.WriteFile(level, []byte("occupied\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		writes := append(stdoutLines(numbered(1, 8)...), tagged{stderr: true, data: "err 09\n"}, tagged{data: "line 10\n"}, tagged{data: "line 11\n"}, tagged{data: "line 12\n"})
		assertCreateFailure(t, respond(t, writes), joined(writes))
	})
}

// BO18: on the create-failure route, the spill-failed line is whole in the stdout view,
// the stderr view, and the combined view. Each stream whose own last line is open takes
// one newline before the line, so neither stream can hide or merge the other's line.
func TestOwnerCreateFailureSeparatesStdoutLine(t *testing.T) {
	privateHome(t)
	head := strings.Join(numbered(1, 10), "")
	line := "spill-failed{reason=" + errInjected.Error() + "}\n"
	for _, row := range []struct {
		name                     string
		writes                   []tagged
		stdout, stderr, combined string
	}{
		{"open stdout", append(stdoutLines(head, "partial"), tagged{stderr: true, data: "err line\n"}), head + "partial\n", "err line\n", head + "partialerr line\n\n"},
		{"open stderr", append(stdoutLines(head, "line 11\n"), tagged{stderr: true, data: "err partial"}), head + "line 11\n", "err partial\n", head + "line 11\nerr partial\n"},
		{"both open", append(stdoutLines(head, "partial"), tagged{stderr: true, data: "err partial"}), head + "partial\n", "err partial\n", head + "partialerr partial\n\n"},
	} {
		t.Run(row.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			owner := New(benchhome.Dir(), &stdout, &stderr, outsideRepository)
			owner.create = func(string) (io.WriteCloser, error) { return nil, errInjected }
			respondWith(t, owner, &stdout, row.writes)
			if stdout.String() != row.stdout+line || stderr.String() != row.stderr {
				t.Fatalf("streams = (%q, %q), want (%q, %q)", stdout.String(), stderr.String(), row.stdout+line, row.stderr)
			}
			var sink bytes.Buffer
			owner = New(benchhome.Dir(), &sink, &sink, outsideRepository)
			owner.create = func(string) (io.WriteCloser, error) { return nil, errInjected }
			if got := respondWith(t, owner, &sink, row.writes); got != row.combined+line {
				t.Fatalf("combined response = %q, want %q", got, row.combined+line)
			}
		})
	}
}

// BO25: a Bench home that holds a newline takes the create-failure route, so no raw path
// can forge a second response line.
func TestOwnerRefusesUnsafeSpillPath(t *testing.T) {
	home := filepath.Join(t.TempDir(), "unsafe\nhome")
	t.Setenv(benchhome.Env, home)
	writes := stdoutLines(numbered(1, 11)...)
	assertCreateFailure(t, respond(t, writes), joined(writes))
	if _, err := os.Stat(home); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("unsafe home was created: %v", err)
	}
}
