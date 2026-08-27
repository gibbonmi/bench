package gocache

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

// fixtureDir builds a cache-shaped directory: two regular files at two depths, one
// `-d` executable directory with a regular file inside it, one FIFO, and one symlink.
// It returns the directory and the byte sum of its regular files.
func fixtureDir(t *testing.T) (string, int64) {
	t.Helper()
	dir := t.TempDir()
	write := func(path, body string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write(filepath.Join(dir, "aa", "one-a"), "abc")
	write(filepath.Join(dir, "bb", "two-a"), "defgh")
	write(filepath.Join(dir, "cc", "three-d", "exe"), "ijklmno")
	if err := syscall.Mkfifo(filepath.Join(dir, "aa", "pipe"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(dir, "aa", "one-a"), filepath.Join(dir, "aa", "link")); err != nil {
		t.Fatal(err)
	}
	return dir, int64(len("abc") + len("defgh") + len("ijklmno"))
}

// R02: the walk sums the regular files at every depth, and it recurses into a `-d`
// directory. A FIFO and a symlink are not regular, so neither adds a byte or a file.
func TestMeasureSumsRegularFilesAndRecursesIntoADirectory(t *testing.T) {
	t.Parallel()
	dir, want := fixtureDir(t)
	footprint := Measure(dir)
	if footprint.Bytes != want || footprint.Files != 3 {
		t.Fatalf("Measure = %d bytes in %d files, want %d bytes in 3 files", footprint.Bytes, footprint.Files, want)
	}
	if footprint.Dir != dir {
		t.Fatalf("Measure dir = %q, want %q", footprint.Dir, dir)
	}
	if footprint.OverBound() {
		t.Fatalf("OverBound = true at %d bytes, want false below the %d bound", footprint.Bytes, Bound)
	}
}

// R06: the trim file's unix seconds render as UTC RFC 3339, with the trailing newline a
// hand edit leaves behind trimmed first.
func TestLastTrimRendersUnixSecondsAsUTC(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, trimFile), []byte("  1700000000\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := Measure(dir).LastTrim; got != "2023-11-14T22:13:20Z" {
		t.Fatalf("LastTrim = %q, want %q", got, "2023-11-14T22:13:20Z")
	}
}

// R07: an absent, symlinked, FIFO, or unparsable trim file leaves the last trim empty.
// The lstat check comes first, so a FIFO is never opened and never blocks the read.
func TestLastTrimIsEmptyWithoutARegularFile(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name  string
		build func(*testing.T, string)
	}{
		{name: "absent", build: func(*testing.T, string) {}},
		{name: "symlink", build: func(t *testing.T, dir string) {
			target := filepath.Join(dir, "elsewhere")
			if err := os.WriteFile(target, []byte("1700000000"), 0o644); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(target, filepath.Join(dir, trimFile)); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "fifo", build: func(t *testing.T, dir string) {
			if err := syscall.Mkfifo(filepath.Join(dir, trimFile), 0o644); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "unparsable", build: func(t *testing.T, dir string) {
			if err := os.WriteFile(filepath.Join(dir, trimFile), []byte("yesterday"), 0o644); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			tc.build(t, dir)
			if got := Measure(dir).LastTrim; got != "" {
				t.Fatalf("LastTrim = %q, want it empty", got)
			}
		})
	}
}
