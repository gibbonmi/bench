package census

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/benchguard"
	"github.com/gibbonmi/bench/internal/poolkey"
)

// TestDirSitsBesideThePool proves the records live under <home>/census and never
// under the worktree pool, which `bench worktree reclaim` enumerates.
// (Coverage row EC26.)
func TestDirSitsBesideThePool(t *testing.T) {
	t.Parallel()
	home := filepath.Join(string(filepath.Separator), "home-a", ".bench")
	root := filepath.Join(string(filepath.Separator), "repos", "example")
	got := Dir(home, root)
	if want := filepath.Join(home, "census", poolkey.Key(root)); got != want {
		t.Fatalf("Dir = %q, want %q", got, want)
	}
	if !strings.HasPrefix(got, filepath.Join(home, "census")+string(filepath.Separator)) {
		t.Fatalf("Dir %q is not under the census home", got)
	}
	if strings.Contains(got, filepath.Dir(poolkey.Pool(home, root))) {
		t.Fatalf("Dir %q is inside the worktree pool", got)
	}
}

// TestDirDependsOnTheInjectedHome proves the home arrives explicitly, with no
// environment read below the effect boundary.
func TestDirDependsOnTheInjectedHome(t *testing.T) {
	t.Parallel()
	root := filepath.Join(string(filepath.Separator), "repos", "example")
	a := Dir(filepath.Join(string(filepath.Separator), "home-a"), root)
	b := Dir(filepath.Join(string(filepath.Separator), "home-b"), root)
	if filepath.Base(a) != filepath.Base(b) {
		t.Fatalf("the census key must depend only on the root: %s vs %s", a, b)
	}
	if a == b {
		t.Fatalf("Dir ignored the injected home: %s", a)
	}
}

const (
	ownerID   = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	knownID   = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	unknownID = "cccccccccccccccccccccccccccccccc"
)

var fixedTime = time.Date(2026, 8, 25, 9, 30, 0, 0, time.UTC)

// fixtureHome returns a temporary Bench home, a repository root, and the pool
// directory the two make. The root is not a repository, so the key derivation
// falls back to the root itself and the pool path stays stable.
func fixtureHome(t *testing.T) (home, root, pool string) {
	t.Helper()
	home, root = t.TempDir(), t.TempDir()
	return home, root, poolkey.Pool(home, root)
}

// records returns the lines of one assignment's record file. An absent file reads
// as no records.
func records(t *testing.T, home, root, id string) []string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(Dir(home, root), id))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	text := strings.TrimSuffix(string(data), "\n")
	if text == "" {
		return nil
	}
	return strings.Split(text, "\n")
}

// TestRecordCountsOneRawCall proves a plain edit against a pool path is counted.
// (Coverage row EC01.)
func TestRecordCountsOneRawCall(t *testing.T) {
	t.Parallel()
	home, root, pool := fixtureHome(t)
	command := "sed -i s/a/b/ " + filepath.Join(pool, ownerID+"-"+knownID, "x")
	if err := Record(command, root, home, fixedTime); err != nil {
		t.Fatal(err)
	}
	if got := records(t, home, root, knownID); len(got) != 1 {
		t.Fatalf("records = %v, want one", got)
	}
}

// TestRecordCountsAChainOnce proves one Bash call is one raw call, whatever the
// number of simple commands that name the pool. (Coverage row EC02.)
func TestRecordCountsAChainOnce(t *testing.T) {
	t.Parallel()
	home, root, pool := fixtureHome(t)
	dir := filepath.Join(pool, ownerID+"-"+knownID)
	command := "cat " + filepath.Join(dir, "a") + " && sed -i x " + filepath.Join(dir, "b")
	if err := Record(command, root, home, fixedTime); err != nil {
		t.Fatal(err)
	}
	if got := records(t, home, root, knownID); len(got) != 1 {
		t.Fatalf("records = %v, want one", got)
	}
}

// TestRecordSkipsABenchCall proves the verb is free: a text whose pool path sits in a
// `bench` command records nothing, at the top level and one wrapper level deep.
// (Coverage rows EC03 and EC11.)
func TestRecordSkipsABenchCall(t *testing.T) {
	t.Parallel()
	inner := "bench worktree exec a -- sed -i x "
	for _, tc := range []struct {
		name string
		text func(path string) string
	}{
		{"top level", func(path string) string { return inner + path }},
		{"inside a bash string", func(path string) string { return "bash -c '" + inner + path + "'" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			home, root, pool := fixtureHome(t)
			path := filepath.Join(pool, ownerID+"-"+knownID, "b")
			if err := Record(tc.text(path), root, home, fixedTime); err != nil {
				t.Fatal(err)
			}
			if got := records(t, home, root, knownID); got != nil {
				t.Fatalf("records = %v, want none", got)
			}
		})
	}
}

// TestRecordSkipsAWrapperPath proves the `bench` test reads the resolved executable,
// not the word: a wrapper invoked by an absolute path is a Bench call. The stub
// resolver keeps the test off the machine's own PATH. (Coverage row EC04.)
func TestRecordSkipsAWrapperPath(t *testing.T) {
	t.Parallel()
	resolver := benchguard.Resolver{
		Getwd:    func() (string, error) { return string(filepath.Separator), nil },
		LookPath: func(string) (string, error) { return "", errors.New("not on PATH") },
		EvalSymlinks: func(path string) (string, error) {
			if path == "/tmp/r/bin/bx" {
				return "/opt/bench/bin/bench", nil
			}
			return path, nil
		},
	}
	for _, wrapper := range []string{"/tmp/r/bin/bench.sh", "/tmp/r/bin/bx"} {
		t.Run(wrapper, func(t *testing.T) {
			t.Parallel()
			home, root, pool := fixtureHome(t)
			command := wrapper + " status " + filepath.Join(pool, ownerID+"-"+knownID)
			if err := recordWith(command, root, home, fixedTime, resolver); err != nil {
				t.Fatal(err)
			}
			if got := records(t, home, root, knownID); got != nil {
				t.Fatalf("records = %v, want none", got)
			}
		})
	}
}

// TestRecordLineHoldsTimeAndHeadOnly proves the record is not a command log.
// (Coverage row EC06.)
func TestRecordLineHoldsTimeAndHeadOnly(t *testing.T) {
	t.Parallel()
	home, root, pool := fixtureHome(t)
	path := filepath.Join(pool, ownerID+"-"+knownID, "secret-operand")
	if err := Record("sed -i s/a/b/ "+path, root, home, fixedTime); err != nil {
		t.Fatal(err)
	}
	got := records(t, home, root, knownID)
	if len(got) != 1 {
		t.Fatalf("records = %v, want one", got)
	}
	field := strings.Split(got[0], "\t")
	if len(field) != 2 {
		t.Fatalf("record %q, want a time and a head", got[0])
	}
	if _, err := time.Parse(time.RFC3339, field[0]); err != nil {
		t.Fatalf("record time %q: %v", field[0], err)
	}
	if field[1] != "sed" {
		t.Fatalf("record head = %q, want %q", field[1], "sed")
	}
	for _, fragment := range []string{"secret-operand", "s/a/b/", path} {
		if strings.Contains(got[0], fragment) {
			t.Fatalf("record %q holds the command text %q", got[0], fragment)
		}
	}
}

// TestRecordKeysAnUnknownAssignment proves the writer never reads the ledger, so a
// path of a dead or unknown assignment still records. (Coverage row EC14.)
func TestRecordKeysAnUnknownAssignment(t *testing.T) {
	t.Parallel()
	home, root, pool := fixtureHome(t)
	command := "sed -i x " + filepath.Join(pool, unknownID+"-"+unknownID, "x")
	if err := Record(command, root, home, fixedTime); err != nil {
		t.Fatal(err)
	}
	if got := records(t, home, root, unknownID); len(got) != 1 {
		t.Fatalf("records = %v, want one", got)
	}
}

// TestRecordKeysOnTheAssignmentID proves the key is the assignment id alone, which is
// the name the ledger carries. (Coverage row EC16.)
func TestRecordKeysOnTheAssignmentID(t *testing.T) {
	t.Parallel()
	home, root, pool := fixtureHome(t)
	segment := ownerID + "-" + knownID
	if err := Record("sed -i x "+filepath.Join(pool, segment, "x"), root, home, fixedTime); err != nil {
		t.Fatal(err)
	}
	entry, err := os.ReadDir(Dir(home, root))
	if err != nil {
		t.Fatal(err)
	}
	if len(entry) != 1 || entry[0].Name() != knownID {
		t.Fatalf("census directory holds %v, want one file named %q", entry, knownID)
	}
}

// TestConcurrentRecordsKeepEveryLine proves the append leaves whole lines when two
// writers record to one assignment. `internal/racetests` runs it under `-race`,
// because the ordinary suite cannot observe the loss. (Coverage row EC15.)
func TestConcurrentRecordsKeepEveryLine(t *testing.T) {
	t.Parallel()
	home, root, pool := fixtureHome(t)
	command := "sed -i x " + filepath.Join(pool, ownerID+"-"+knownID, "x")
	var group sync.WaitGroup
	for range 2 {
		group.Add(1)
		go func() {
			defer group.Done()
			if err := Record(command, root, home, fixedTime); err != nil {
				t.Error(err)
			}
		}()
	}
	group.Wait()
	got := records(t, home, root, knownID)
	if len(got) != 2 {
		t.Fatalf("records = %v, want two", got)
	}
	for _, line := range got {
		if len(strings.Split(line, "\t")) != 2 {
			t.Fatalf("line %q is not intact", line)
		}
	}
}

// TestRecordOnAnUnwritableHomeFails proves a home the process cannot create returns an
// error and leaves no file. The caller ignores the error, so the call still passes.
func TestRecordOnAnUnwritableHomeFails(t *testing.T) {
	t.Parallel()
	base, root := t.TempDir(), t.TempDir()
	blocker := filepath.Join(base, "blocker")
	if err := os.WriteFile(blocker, []byte("not a directory\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	home := filepath.Join(blocker, "home")
	pool := poolkey.Pool(home, root)
	err := Record("sed -i x "+filepath.Join(pool, ownerID+"-"+knownID, "x"), root, home, fixedTime)
	if err == nil {
		t.Fatal("Record returned no error for an unwritable home")
	}
	if _, statErr := os.Stat(Dir(home, root)); statErr == nil {
		t.Fatal("the census directory exists after the failure")
	}
}

// TestRecordRefusesASymlinkedCensusDirectory proves the writer never follows a
// redirected census directory.
func TestRecordRefusesASymlinkedCensusDirectory(t *testing.T) {
	t.Parallel()
	home, root, pool := fixtureHome(t)
	target := t.TempDir()
	if err := os.MkdirAll(filepath.Join(home, "census"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, Dir(home, root)); err != nil {
		t.Skipf("the filesystem refuses symlinks: %v", err)
	}
	err := Record("sed -i x "+filepath.Join(pool, ownerID+"-"+knownID, "x"), root, home, fixedTime)
	if err == nil {
		t.Fatal("Record returned no error for a symlinked census directory")
	}
	entry, readErr := os.ReadDir(target)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if len(entry) != 0 {
		t.Fatalf("the writer followed the symlink and wrote %v", entry)
	}
}

// TestRecordMatchesOnlyAPoolPath walks the edges of the raw-text match: the separator
// after the pool, the absolute form, the folded quoted word, and the assignment shape
// of the segment.
func TestRecordMatchesOnlyAPoolPath(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		command func(pool string) string
		want    int
	}{
		{"a sibling of the pool", func(pool string) string {
			return "sed -i x " + pool + "x" + string(filepath.Separator) + ownerID + "-" + knownID + "/y"
		}, 0},
		{"a relative path", func(pool string) string {
			return "sed -i x " + strings.TrimPrefix(pool, string(filepath.Separator)) + "/" + ownerID + "-" + knownID + "/y"
		}, 0},
		{"a segment without the assignment shape", func(pool string) string {
			return "sed -i x " + filepath.Join(pool, "scratch", "y")
		}, 0},
		{"a quoted path with a space", func(pool string) string {
			return "sed -i x '" + filepath.Join(pool, ownerID+"-"+knownID, "a b") + "'"
		}, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			home, root, pool := fixtureHome(t)
			if err := Record(tc.command(pool), root, home, fixedTime); err != nil {
				t.Fatal(err)
			}
			if got := records(t, home, root, knownID); len(got) != tc.want {
				t.Fatalf("records = %v, want %d", got, tc.want)
			}
		})
	}
}
