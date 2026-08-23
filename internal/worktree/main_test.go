package worktree

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// inheritedBenchHome records the operator's BENCH_HOME before TestMain replaces it.
// Package variables initialize before TestMain runs, so this is the pre-swap value.
var inheritedBenchHome = os.Getenv("BENCH_HOME")

// privateBenchHome is the process-private BENCH_HOME every test in this package runs
// under. A fixture that forgets to bind its own home lands here instead of the
// operator's pool. The post-run residue check turns that into a red.
var privateBenchHome string

// TestMain gives the package its own BENCH_HOME, created empty in the OS temp
// directory and exported before the first test. It checks the home for residue after
// the last test, then removes it. Creation failure is fail-closed, because running the
// package against the inherited home is the leak this exists to prevent.
func TestMain(m *testing.M) {
	home, err := os.MkdirTemp("", "bench-worktree-home-")
	if err != nil {
		fmt.Fprintln(os.Stderr, "private BENCH_HOME:", err)
		os.Exit(1)
	}
	privateBenchHome = home
	if err := os.Setenv("BENCH_HOME", home); err != nil {
		fmt.Fprintln(os.Stderr, "private BENCH_HOME:", err)
		os.RemoveAll(home)
		os.Exit(1)
	}
	code := m.Run()
	// The one test-run executable owner outlives every journey; release its private
	// directory only after the last child has returned.
	packageRunBinary.close()
	entries, err := homeResidue(home)
	if err != nil {
		fmt.Fprintln(os.Stderr, "private BENCH_HOME residue:", err)
		code = 1
	}
	if len(entries) > 0 {
		fmt.Fprintf(os.Stderr, "%d residue entries under private BENCH_HOME %s — a test created a worktree under a home it did not bind itself\n", len(entries), home)
		for _, entry := range entries {
			for _, line := range entry.render() {
				fmt.Fprintln(os.Stderr, line)
			}
		}
		code = 1
	}
	if err := os.RemoveAll(home); err != nil {
		fmt.Fprintln(os.Stderr, "private BENCH_HOME removal:", err)
		code = 1
	}
	os.Exit(code)
}

// residueEntry is one top-level entry left under a home after a run: a stray file, or
// a directory such as worktrees/, empty or not. Nothing legitimate materializes any of
// them. origins carries the gitdir: line of every leaked pool worktree beneath it,
// which names the test that created it.
type residueEntry struct {
	path    string
	origins []string
}

func (e residueEntry) render() []string {
	return append([]string{e.path}, e.origins...)
}

// homeResidue computes the residue for one home. A directory it cannot walk costs
// that entry its origin lines, never the entry itself. The report has to name the
// leak even when it cannot name the leaker.
func homeResidue(home string) ([]residueEntry, error) {
	names, err := os.ReadDir(home)
	if err != nil {
		return nil, err
	}
	var entries []residueEntry
	for _, name := range names {
		path := filepath.Join(home, name.Name())
		origins, err := residueOrigins(path)
		if err != nil {
			origins = append(origins, "  origin scan: "+err.Error())
		}
		entries = append(entries, residueEntry{path: path, origins: origins})
	}
	return entries, nil
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
	// An empty inherited value means benchHome falls back to the operator's ~/.bench,
	// which the next assertion covers unconditionally, so the pair is never both vacuous.
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

	// A .git that is a full repository rather than a pointer, and a pointer carrying no
	// gitdir: line — both are still residue. Both simply have no origin to name.
	malformed := t.TempDir()
	mustMkdirAll(t, filepath.Join(malformed, "worktrees", "001-2", "repo", ".git"), 0o755)
	mustMkdirAll(t, filepath.Join(malformed, "worktrees", "001-3", "blank"), 0o755)
	mustWrite(t, filepath.Join(malformed, "worktrees", "001-3", "blank", ".git"), []byte("not a pointer\n"), 0o644)

	for _, tc := range []struct {
		name    string
		home    string
		entries int
		origins int
	}{
		{"leaked worktree and stray file", leaked, 2, 1},
		{"empty pool directory", onlyPool, 1, 0},
		{"clean home", empty, 0, 0},
		{"malformed git pointers", malformed, 1, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			entries, err := homeResidue(tc.home)
			requireTest(t, err == nil, "homeResidue: %v", err)
			requireTest(t, len(entries) == tc.entries, "entries = %#v, want %d", entries, tc.entries)
			origins := 0
			for _, entry := range entries {
				origins += len(entry.origins)
			}
			requireTest(t, origins == tc.origins, "origin lines = %d, want %d", origins, tc.origins)
		})
	}

	entries, err := homeResidue(leaked)
	requireTest(t, err == nil, "homeResidue: %v", err)
	var lines []string
	for _, entry := range entries {
		lines = append(lines, entry.render()...)
	}
	joined := strings.Join(lines, "\n")
	requireTest(t, strings.Contains(joined, "gitdir: "+origin), "report %q does not name the origin %q", joined, origin)
}

// withinDir reports whether path is parent itself or lies beneath it. It stays
// independent of the production insidePool so that the Pool assertion above cannot
// be satisfied by a bug the two share.
func withinDir(parent, path string) bool {
	rel, err := filepath.Rel(parent, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func TestWithinDir(t *testing.T) {
	for _, tc := range []struct {
		parent, path string
		want         bool
	}{
		{"/a/b", "/a/b", true},
		{"/a/b", "/a/b/c", true},
		{"/a/b", "/a/bc", false},
		{"/a/b", "/a", false},
		{"/a/b", "/", false},
	} {
		got := withinDir(tc.parent, tc.path)
		requireTest(t, got == tc.want, "withinDir(%q, %q) = %t, want %t", tc.parent, tc.path, got, tc.want)
	}
}
