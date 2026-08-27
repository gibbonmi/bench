package gocache

import (
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// cacheFixture is the shape Go leaves in a build cache: two-hex shard directories holding
// archives, beside the lock file, the trim file, and Go's README. shardBytes and
// shardFiles are what a clean removes, because Go's own clean takes the shards alone.
type cacheFixture struct {
	home       string
	dir        string
	shardBytes int64
	shardFiles int64
}

// newCacheFixture writes the fixture under a fresh HOME and returns its measurements.
func newCacheFixture(t *testing.T) cacheFixture {
	t.Helper()
	home := t.TempDir()
	dir := cacheDir(home)
	fixture := cacheFixture{home: home, dir: dir}
	for _, shard := range []struct{ name, body string }{
		{"aa", "archive-a"},
		{"bf", "archive-bf-longer"},
	} {
		if err := os.MkdirAll(filepath.Join(dir, shard.name), 0o755); err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(dir, shard.name, shard.name+"0123456789-a")
		if err := os.WriteFile(path, []byte(shard.body), 0o644); err != nil {
			t.Fatal(err)
		}
		fixture.shardBytes += int64(len(shard.body))
		fixture.shardFiles++
	}
	write(t, filepath.Join(dir, trimFile), "1700000000\n")
	write(t, filepath.Join(dir, "README"), "This is the Go build cache.\n")
	write(t, filepath.Join(dir, LockFile), "")
	return fixture
}

func write(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// entries lists the fixture directory's immediate names, sorted.
func entries(t *testing.T, dir string) []string {
	t.Helper()
	read, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	names := make([]string, 0, len(read))
	for _, entry := range read {
		names = append(names, entry.Name())
	}
	sort.Strings(names)
	return names
}

// cleanEnv is the environment slice the clean reads: the fixture home and the PATH that
// resolves the real toolchain.
func cleanEnv(home string) []string {
	return []string{"HOME=" + home, "PATH=" + os.Getenv("PATH")}
}

// While a holder holds the cache lock, the clean exits 1. This row grades the Hold-to-clean
// contract itself. L01, L02, and L03 bind that contract to each production holder, and they
// live beside those call sites.
func TestCleanRefusesWhileAHolderRuns(t *testing.T) {
	fixture := newCacheFixture(t)
	holder, err := Hold([]string{"HOME=" + fixture.home})
	if err != nil {
		t.Fatalf("hold = %v", err)
	}
	defer holder.Release()
	out, code := runCleanProcess(t, fixture.home)
	if code != 1 {
		t.Fatalf("clean under a holder = exit %d, want 1; out=%q", code, out)
	}
	if !strings.HasPrefix(out, "error: cache in use — ") {
		t.Fatalf("clean under a holder = %q, want the cache-in-use refusal", out)
	}
}

// L04: a refused clean removes nothing. A clean that removed before it locked would empty
// the directory under a live holder.
func TestRefusedCleanRemovesNoFile(t *testing.T) {
	fixture := newCacheFixture(t)
	holder, err := Hold([]string{"HOME=" + fixture.home})
	if err != nil {
		t.Fatalf("hold = %v", err)
	}
	defer holder.Release()
	out, code := runCleanProcess(t, fixture.home)
	if code != 1 {
		t.Fatalf("clean = exit %d, want 1; out=%q", code, out)
	}
	if got, want := entries(t, fixture.dir), []string{"README", "aa", LockFile, "bf", trimFile}; !slices.Equal(got, want) {
		t.Fatalf("directory after a refused clean = %q, want %q", got, want)
	}
	// The shards plus the three files a clean would keep: README, trim.txt, and the lock.
	if got, want := Measure(fixture.dir).Files, fixture.shardFiles+3; got != want {
		t.Fatalf("files after a refused clean = %d, want %d", got, want)
	}
}

// L11: the refusal names the holder's pid, so the operator can find the run that is
// compiling instead of being told only that something is.
func TestCleanRefusalNamesTheHolderPID(t *testing.T) {
	fixture := newCacheFixture(t)
	holder, err := Hold([]string{"HOME=" + fixture.home})
	if err != nil {
		t.Fatalf("hold = %v", err)
	}
	defer holder.Release()
	out, code := runCleanProcess(t, fixture.home)
	if code != 1 {
		t.Fatalf("clean = exit %d, want 1; out=%q", code, out)
	}
	if want := "pid " + strconv.Itoa(os.Getpid()); !strings.Contains(out, want) {
		t.Fatalf("clean refusal = %q, want it to name %q", out, want)
	}
}

// L06: with no holder the clean removes every two-hex subdirectory, so no bytes stay behind.
// L07: it leaves bench.lock, trim.txt, and README, so the file the next holder opens survives.
func TestCleanRemovesTheShardsAndKeepsTheRest(t *testing.T) {
	fixture := newCacheFixture(t)
	out, code := clean(cleanEnv(fixture.home))
	if code != 0 {
		t.Fatalf("clean = exit %d, want 0; out=%q", code, out)
	}
	if got, want := entries(t, fixture.dir), []string{"README", LockFile, trimFile}; !slices.Equal(got, want) {
		t.Fatalf("directory after a clean = %q, want %q", got, want)
	}
}

// L12: the clean reports the bytes and the files the removal took, measured before it ran.
// A clean that measured after the removal reports zero.
func TestCleanReportsThePreRemovalMeasurement(t *testing.T) {
	fixture := newCacheFixture(t)
	out, code := clean(cleanEnv(fixture.home))
	if code != 0 {
		t.Fatalf("clean = exit %d, want 0; out=%q", code, out)
	}
	want := "go_build_cache_clean[1]{dir,bytes_removed,files_removed}:\n" +
		"  " + fixture.dir + "," + strconv.FormatInt(fixture.shardBytes, 10) + "," + strconv.FormatInt(fixture.shardFiles, 10) + "\n"
	if out != want {
		t.Fatalf("clean = %q, want %q", out, want)
	}
}

// L08: an absent directory reports zero removed at exit 0 and creates nothing, so a fresh
// machine passes. A stat error that became a refusal would fail here.
func TestCleanOnAnAbsentDirectory(t *testing.T) {
	home := t.TempDir()
	dir := cacheDir(home)
	out, code := clean(cleanEnv(home))
	if code != 0 {
		t.Fatalf("clean = exit %d, want 0; out=%q", code, out)
	}
	want := "go_build_cache_clean[1]{dir,bytes_removed,files_removed}:\n  " + dir + ",0,0\n"
	if out != want {
		t.Fatalf("clean = %q, want %q", out, want)
	}
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatalf("stat after a clean on an absent directory = %v, want it still absent", err)
	}
}

// L10: no `go` on PATH is a refusal that names `go`, so the operator reads the missing tool
// rather than a bare exec error.
func TestCleanWithoutGoOnPath(t *testing.T) {
	fixture := newCacheFixture(t)
	out, code := clean([]string{"HOME=" + fixture.home, "PATH=" + t.TempDir()})
	if code != 1 {
		t.Fatalf("clean = exit %d, want 1; out=%q", code, out)
	}
	if !strings.HasPrefix(out, "error: go not found — ") || !strings.Contains(out, "`go`") {
		t.Fatalf("clean = %q, want a refusal that names go", out)
	}
	if got, want := entries(t, fixture.dir), []string{"README", "aa", LockFile, "bf", trimFile}; !slices.Equal(got, want) {
		t.Fatalf("directory after the refusal = %q, want %q", got, want)
	}
}

// A derivation refusal reaches the operator as the AXI error contract, the same way the
// report's does.
func TestCleanRefusesWithoutAnAbsoluteHome(t *testing.T) {
	out, code := clean([]string{"PATH=/usr/bin"})
	if code != 1 || !strings.HasPrefix(out, "error: cache directory not derived — ") || !strings.Contains(out, "HOME") {
		t.Fatalf("clean = %q exit %d, want a HOME refusal at 1", out, code)
	}
}

// A control byte in the derived path refuses before any lock or removal.
func TestCleanRefusesAControlByteInThePath(t *testing.T) {
	out, code := clean([]string{"HOME=/home/ag\x1bent", "PATH=/usr/bin"})
	if code != 1 || !strings.HasPrefix(out, "error: unrepresentable cache directory — ") {
		t.Fatalf("clean = %q exit %d, want the unrepresentable refusal at 1", out, code)
	}
}
