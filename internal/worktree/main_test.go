package worktree

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// inheritedBenchHome records the operator's BENCH_HOME before TestMain replaces it, so
// TestPackageHomeIsPrivate can prove the tests never resolve a pool there. Package
// variables initialize before TestMain runs, which is what makes the record honest.
var inheritedBenchHome = os.Getenv("BENCH_HOME")

// privateBenchHome is the process-private BENCH_HOME every test in this package runs
// under. A fixture that forgets to bind its own home lands here instead of the
// operator's pool, and the post-run residue check turns that into a red.
var privateBenchHome string

// TestMain gives the package a BENCH_HOME of its own: created empty in the OS temp
// directory, exported before the first test, checked for residue after the last, and
// removed. The exit code combines the three verdicts — a failing test, residue, or a
// home that will not remove — so neither can mask the others. Creation failure is
// fail-closed: running the package against the inherited home is the leak this exists
// to prevent.
func TestMain(m *testing.M) {
	home, err := os.MkdirTemp("", "bench-worktree-home-")
	if err != nil {
		fmt.Fprintln(os.Stderr, "private BENCH_HOME:", err)
		os.Exit(1)
	}
	privateBenchHome = home
	if err := os.Setenv("BENCH_HOME", home); err != nil {
		fmt.Fprintln(os.Stderr, "private BENCH_HOME:", err)
		os.Exit(1)
	}
	code := m.Run()
	report, err := homeResidue(home)
	if err != nil {
		fmt.Fprintln(os.Stderr, "private BENCH_HOME residue:", err)
		code = 1
	}
	if len(report.entries) > 0 {
		fmt.Fprintf(os.Stderr, "%d residue entries under private BENCH_HOME %s — a test created a worktree under a home it did not bind itself\n", len(report.entries), home)
		for _, line := range report.lines {
			fmt.Fprintln(os.Stderr, line)
		}
		code = 1
	}
	if err := os.RemoveAll(home); err != nil {
		fmt.Fprintln(os.Stderr, "private BENCH_HOME removal:", err)
		code = 1
	}
	os.Exit(code)
}

// residueReport is what remains under a home after a run. The unit is the top-level
// entry — a stray file, or a directory such as worktrees/ whether or not it holds
// anything — because nothing legitimate materializes any of them.
type residueReport struct {
	entries []string
	lines   []string
}

// homeResidue computes the report for one home. Each leaked pool worktree beneath a
// top-level entry contributes its .git pointer's gitdir: line, which carries the
// temporary root of the test that created it and so names the offender.
func homeResidue(home string) (residueReport, error) {
	names, err := os.ReadDir(home)
	if err != nil {
		return residueReport{}, err
	}
	var report residueReport
	for _, name := range names {
		entry := filepath.Join(home, name.Name())
		report.entries = append(report.entries, entry)
		report.lines = append(report.lines, entry)
		origins, err := residueOrigins(entry)
		if err != nil {
			return residueReport{}, err
		}
		report.lines = append(report.lines, origins...)
	}
	return report, nil
}

func residueOrigins(entry string) ([]string, error) {
	var origins []string
	err := filepath.WalkDir(entry, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || d.Name() != ".git" {
			return err
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			// A residue worktree without a readable pointer is still listed by path;
			// its origin line is simply absent.
			return nil
		}
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "gitdir:") {
				origins = append(origins, "  "+strings.TrimSpace(line))
			}
		}
		return nil
	})
	return origins, err
}

func TestPackageHomeIsPrivate(t *testing.T) {
	home := os.Getenv("BENCH_HOME")
	requireTest(t, home == privateBenchHome, "BENCH_HOME = %q, want the private home %q", home, privateBenchHome)
	info, err := os.Stat(home)
	requireTest(t, err == nil && info.IsDir(), "private home stat = %v, %v", info, err)
	entries, err := os.ReadDir(home)
	requireTest(t, err == nil && len(entries) == 0, "private home entries = %#v, %v, want empty", entries, err)
	requireTest(t, home != inheritedBenchHome, "private home equals the inherited BENCH_HOME %q", inheritedBenchHome)
	if inheritedBenchHome != "" {
		requireTest(t, !withinDir(inheritedBenchHome, home), "private home %q is under the inherited BENCH_HOME %q", home, inheritedBenchHome)
	}
	userHome, err := os.UserHomeDir()
	requireTest(t, err == nil, "user home: %v", err)
	fallback := filepath.Join(userHome, ".bench")
	requireTest(t, !withinDir(fallback, home), "private home %q is under the operator's %q", home, fallback)
	pool := Pool(t.TempDir())
	requireTest(t, withinDir(home, pool), "Pool of a fresh root = %q, want under %q", pool, home)
}

func TestHomeResidueListsLeakedWorktreesWithOrigin(t *testing.T) {
	origin := filepath.Join(t.TempDir(), "TestReauthorizeSomething1234", "001")
	leaked := t.TempDir()
	mustMkdirAll(t, filepath.Join(leaked, "worktrees", "001-1", "abcdef"), 0o755)
	mustWrite(t, filepath.Join(leaked, "worktrees", "001-1", "abcdef", ".git"), []byte("gitdir: "+origin+"/.git/worktrees/abcdef\n"), 0o644)
	mustWrite(t, filepath.Join(leaked, "stray.txt"), []byte("stray\n"), 0o644)

	empty := t.TempDir()
	onlyPool := t.TempDir()
	mustMkdirAll(t, filepath.Join(onlyPool, "worktrees"), 0o755)

	for _, tc := range []struct {
		name    string
		home    string
		entries int
	}{
		{"leaked worktree and stray file", leaked, 2},
		{"empty pool directory", onlyPool, 1},
		{"clean home", empty, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			report, err := homeResidue(tc.home)
			requireTest(t, err == nil, "homeResidue: %v", err)
			requireTest(t, len(report.entries) == tc.entries, "entries = %#v, want %d", report.entries, tc.entries)
		})
	}

	report, err := homeResidue(leaked)
	requireTest(t, err == nil, "homeResidue: %v", err)
	joined := strings.Join(report.lines, "\n")
	requireTest(t, strings.Contains(joined, "gitdir: "+origin), "report %q does not name the origin %q", joined, origin)
}

// withinDir reports whether path is parent itself or lies beneath it.
func withinDir(parent, path string) bool {
	rel, err := filepath.Rel(parent, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
