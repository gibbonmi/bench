package gocache

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// homeEnv builds the one-entry environment slice a report reads, with home as HOME.
func homeSlice(home string) []string { return []string{"HOME=" + home} }

// R01: the report is one go_build_cache table with the six declared columns, and the
// row carries the measured directory, bytes, files, last trim, bound, and flag.
func TestReportPrintsTheGoBuildCacheTable(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	dir := filepath.Join(home, ".cache", "bench", "go-build")
	if err := os.MkdirAll(filepath.Join(dir, "aa"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "aa", "one-a"), []byte("abcde"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, trimFile), []byte("1700000000"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code := report(homeSlice(home))
	if code != 0 {
		t.Fatalf("report exit = %d, want 0; out=%q", code, out)
	}
	// trim.txt is itself a regular file under the directory, so the count is the archive
	// plus the trim file: 5 bytes and 10 bytes in two files.
	want := "go_build_cache[1]{dir,bytes,files,last_trim,bound,over_bound}:\n" +
		"  " + dir + ",15,2,\"2023-11-14T22:13:20Z\",10737418240,false\n"
	if out != want {
		t.Fatalf("report = %q, want %q", out, want)
	}
}

// R03: an absent directory is zeros at exit 0, so a machine with no Bench cache yet
// reports rather than refuses.
func TestReportOnAnAbsentDirectory(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	out, code := report(homeSlice(home))
	assertZeroRow(t, filepath.Join(home, ".cache", "bench", "go-build"), out, code)
}

// R04: an empty directory is the same zero row at exit 0. A walk that needed a trim file
// would refuse here.
func TestReportOnAnEmptyDirectory(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	dir := filepath.Join(home, ".cache", "bench", "go-build")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	out, code := report(homeSlice(home))
	assertZeroRow(t, dir, out, code)
}

func assertZeroRow(t *testing.T, dir, out string, code int) {
	t.Helper()
	if code != 0 {
		t.Fatalf("report exit = %d, want 0; out=%q", code, out)
	}
	want := "go_build_cache[1]{dir,bytes,files,last_trim,bound,over_bound}:\n" +
		"  " + dir + ",0,0,\"\",10737418240,false\n"
	if out != want {
		t.Fatalf("report = %q, want %q", out, want)
	}
}

// R05: the report reads the environment alone, so it answers from a working directory
// outside every git repository. A git-root lookup would refuse here instead.
func TestReportRunsOutsideAGitRepository(t *testing.T) {
	outside := t.TempDir()
	if _, err := os.Stat(filepath.Join(outside, ".git")); !os.IsNotExist(err) {
		t.Fatalf("fixture %q is not outside a repository", outside)
	}
	t.Chdir(outside)
	out, code := report(homeSlice(outside))
	if code != 0 {
		t.Fatalf("report exit = %d, want 0 outside a repository; out=%q", code, out)
	}
	if !strings.HasPrefix(out, "go_build_cache[1]{") {
		t.Fatalf("report = %q, want the go_build_cache table", out)
	}
}

// R14: a control byte in the derived path refuses with a named reason before any table
// renders it.
func TestReportRefusesAControlByteInThePath(t *testing.T) {
	t.Parallel()
	out, code := report(homeSlice("/home/ag\x1bent"))
	if code != 1 {
		t.Fatalf("report exit = %d, want 1; out=%q", code, out)
	}
	if !strings.HasPrefix(out, "error: unrepresentable cache directory — ") || !strings.Contains(out, "control byte") {
		t.Fatalf("report = %q, want a refusal that names the control byte", out)
	}
	if strings.Contains(out, "\x1b") {
		t.Fatalf("report = %q, want no control byte in the refusal", out)
	}
}

// A derivation refusal reaches the operator as the AXI error contract, not as a panic or
// a zero row that looks measured.
func TestReportRefusesWithoutAnAbsoluteHome(t *testing.T) {
	t.Parallel()
	out, code := report([]string{"PATH=/usr/bin"})
	if code != 1 || !strings.HasPrefix(out, "error: cache directory not derived — ") || !strings.Contains(out, "HOME") {
		t.Fatalf("report = %q exit %d, want a HOME refusal at 1", out, code)
	}
}
