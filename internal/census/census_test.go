package census

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/benchguard"
	"github.com/gibbonmi/bench/internal/capability"
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

// TestRecordKeysOnTheFirstAssignmentID proves a text naming two assignment ids
// records once, under the first id in the text. (Coverage row EC02.)
func TestRecordKeysOnTheFirstAssignmentID(t *testing.T) {
	t.Parallel()
	home, root, pool := fixtureHome(t)
	command := "cat " + filepath.Join(pool, ownerID+"-"+knownID, "x") + " " + filepath.Join(pool, ownerID+"-"+unknownID, "y")
	if err := Record(command, root, home, fixedTime); err != nil {
		t.Fatal(err)
	}
	if got := records(t, home, root, knownID); len(got) != 1 {
		t.Fatalf("records under the first id = %v, want one", got)
	}
	if got := records(t, home, root, unknownID); len(got) != 0 {
		t.Fatalf("records under the second id = %v, want none", got)
	}
}

// TestRecordSkipsABenchCall proves the verb is free: a text whose pool path sits in a
// `bench` command records nothing, at the top level and one wrapper level deep. An
// allowed exec follow-on also records nothing, because the verb reads the call's own
// head and never re-classifies the follow-on. (Coverage rows EC03, EC11, and G13.)
func TestRecordSkipsABenchCall(t *testing.T) {
	t.Parallel()
	inner := "bench worktree exec a -- sed -i x "
	for _, tc := range []struct {
		name string
		text func(path string) string
	}{
		{"top level", func(path string) string { return inner + path }},
		{"inside a bash string", func(path string) string { return "bash -c '" + inner + path + "'" }},
		{"with an allowed follow-on", func(path string) string { return inner + path + "; cp a b" }},
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
		capability.Capability(t, capability.Symlink, fmt.Sprintf("the filesystem refuses symlinks: %v", err))
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

// TestRecordNamesTheHeadBehindAHeredocBody proves a scripted edit still counts with a
// real head: the pool path sits in the heredoc body the parser strips, so the head
// falls back to the first command's resolved head. (Coverage row EC05.)
func TestRecordNamesTheHeadBehindAHeredocBody(t *testing.T) {
	t.Parallel()
	home, root, pool := fixtureHome(t)
	path := filepath.Join(pool, ownerID+"-"+knownID, "x")
	command := "python3 - <<EOF\nopen(\"" + path + "\")\nEOF"
	if err := Record(command, root, home, fixedTime); err != nil {
		t.Fatal(err)
	}
	got := records(t, home, root, knownID)
	if len(got) != 1 {
		t.Fatalf("records = %v, want one", got)
	}
	if head := strings.Split(got[0], "\t")[1]; head != "python3" {
		t.Fatalf("record head = %q, want %q", head, "python3")
	}
}

// TestRecordNamesTheGitSubcommand proves a bare `git` head never hides the missing
// verb: the head carries the first subcommand word. (Coverage row EC07.)
func TestRecordNamesTheGitSubcommand(t *testing.T) {
	t.Parallel()
	home, root, pool := fixtureHome(t)
	path := filepath.Join(pool, ownerID+"-"+knownID, "x")
	command := "git add " + path
	if err := Record(command, root, home, fixedTime); err != nil {
		t.Fatal(err)
	}
	got := records(t, home, root, knownID)
	if len(got) != 1 {
		t.Fatalf("records = %v, want one", got)
	}
	if head := strings.Split(got[0], "\t")[1]; head != "git add" {
		t.Fatalf("record head = %q, want %q", head, "git add")
	}
}

// TestRecordSkipsARoutinePrefix proves an unresolved prefix never gets recorded in
// the verb's place: an assignment and `timeout` are stepped over. (Coverage row EC08.)
func TestRecordSkipsARoutinePrefix(t *testing.T) {
	t.Parallel()
	home, root, pool := fixtureHome(t)
	path := filepath.Join(pool, ownerID+"-"+knownID, "y")
	command := "FOO=1 timeout 5 sed -i x " + path
	if err := Record(command, root, home, fixedTime); err != nil {
		t.Fatal(err)
	}
	got := records(t, home, root, knownID)
	if len(got) != 1 {
		t.Fatalf("records = %v, want one", got)
	}
	if head := strings.Split(got[0], "\t")[1]; head != "sed" {
		t.Fatalf("record head = %q, want %q", head, "sed")
	}
}

// TestRecordNamesCdOverACompoundEdit proves a `cd` into the assignment directory
// records the real head `cd`, not the trailing `make`.
func TestRecordNamesCdOverACompoundEdit(t *testing.T) {
	t.Parallel()
	home, root, pool := fixtureHome(t)
	dir := filepath.Join(pool, ownerID+"-"+knownID)
	command := "cd " + dir + " && make"
	if err := Record(command, root, home, fixedTime); err != nil {
		t.Fatal(err)
	}
	got := records(t, home, root, knownID)
	if len(got) != 1 {
		t.Fatalf("records = %v, want one", got)
	}
	if head := strings.Split(got[0], "\t")[1]; head != "cd" {
		t.Fatalf("record head = %q, want %q", head, "cd")
	}
}

// TestRecordSkipsAGitGlobalFlagBeforeTheSubcommand proves a flag before the
// subcommand, such as `-C`, never becomes the subcommand.
func TestRecordSkipsAGitGlobalFlagBeforeTheSubcommand(t *testing.T) {
	t.Parallel()
	home, root, pool := fixtureHome(t)
	dir := filepath.Join(pool, ownerID+"-"+knownID)
	command := "git -C " + dir + " add x"
	if err := Record(command, root, home, fixedTime); err != nil {
		t.Fatal(err)
	}
	got := records(t, home, root, knownID)
	if len(got) != 1 {
		t.Fatalf("records = %v, want one", got)
	}
	if head := strings.Split(got[0], "\t")[1]; head != "git add" {
		t.Fatalf("record head = %q, want %q", head, "git add")
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

// writeRecordFile puts raw bytes at one name in the census directory, so a test can
// state a file shape the writer never makes.
func writeRecordFile(t *testing.T, home, root, name, body string) {
	t.Helper()
	dir := Dir(home, root)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

// TestCountsReadsAnAbsentDirectoryAsEmpty proves the board gets a map, not a failure,
// before any raw call is recorded.
func TestCountsReadsAnAbsentDirectoryAsEmpty(t *testing.T) {
	t.Parallel()
	home, root, _ := fixtureHome(t)
	counts, err := Counts(home, root)
	if err != nil || len(counts) != 0 {
		t.Fatalf("Counts on an absent directory = %v, %v; want an empty map and no error", counts, err)
	}
}

// TestCountsReadsEachRecordFileShape proves an empty file counts zero and a last line
// with no newline still counts as one record.
func TestCountsReadsEachRecordFileShape(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, body string
		want       int
	}{
		{"empty file", "", 0},
		{"one closed line", "t\tsed\n", 1},
		{"one unterminated line", "t\tsed", 1},
		{"two lines, the last unterminated", "t\tsed\nt\tpython3", 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			home, root, _ := fixtureHome(t)
			writeRecordFile(t, home, root, knownID, tc.body)
			counts, err := Counts(home, root)
			if err != nil || counts[knownID] != tc.want {
				t.Fatalf("Counts = %v, %v; want %s = %d", counts, err, knownID, tc.want)
			}
		})
	}
}

// TestCountsIgnoresAForeignName proves only a 32-hex assignment id is read, so a
// stray file in the census directory never becomes a count.
func TestCountsIgnoresAForeignName(t *testing.T) {
	t.Parallel()
	home, root, _ := fixtureHome(t)
	writeRecordFile(t, home, root, "README", "t\tsed\n")
	writeRecordFile(t, home, root, ownerID+"-"+knownID, "t\tsed\n")
	counts, err := Counts(home, root)
	if err != nil || len(counts) != 0 {
		t.Fatalf("Counts = %v, %v; want no foreign name counted", counts, err)
	}
}

// TestCountsRefusesAFifoWithoutBlocking proves a refused file type reads as zero and
// never holds the board open on a reader that has no writer.
func TestCountsRefusesAFifoWithoutBlocking(t *testing.T) {
	t.Parallel()
	home, root, _ := fixtureHome(t)
	dir := Dir(home, root)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(filepath.Join(dir, knownID), 0o600); err != nil {
		capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
	}
	done := make(chan map[string]int, 1)
	go func() {
		counts, _ := Counts(home, root)
		done <- counts
	}()
	select {
	case counts := <-done:
		if len(counts) != 0 {
			t.Fatalf("Counts on a FIFO record = %v, want none", counts)
		}
	case <-time.After(time.Second):
		t.Fatal("Counts blocked on a FIFO record file")
	}
}

// TestCountsCountsTheRecordedCalls proves the reader and the writer agree: the count
// is the number of lines Record appended.
func TestCountsCountsTheRecordedCalls(t *testing.T) {
	t.Parallel()
	home, root, pool := fixtureHome(t)
	command := "sed -i s/a/b/ " + filepath.Join(pool, ownerID+"-"+knownID, "x")
	for range 3 {
		if err := Record(command, root, home, fixedTime); err != nil {
			t.Fatal(err)
		}
	}
	counts, err := Counts(home, root)
	if err != nil || counts[knownID] != 3 {
		t.Fatalf("Counts = %v, %v; want %s = 3", counts, err, knownID)
	}
}

// TestDropRemovesOneAssignmentsRecords proves the drop is exact: the retired
// assignment's file goes, and every other assignment keeps its records.
func TestDropRemovesOneAssignmentsRecords(t *testing.T) {
	t.Parallel()
	home, root, _ := fixtureHome(t)
	writeRecordFile(t, home, root, knownID, "t\tsed\n")
	writeRecordFile(t, home, root, unknownID, "t\tsed\n")
	if err := Drop(home, root, knownID); err != nil {
		t.Fatalf("Drop = %v, want no error", err)
	}
	counts, err := Counts(home, root)
	if err != nil || len(counts) != 1 || counts[unknownID] != 1 {
		t.Fatalf("Counts after the drop = %v, %v; want only %s", counts, err, unknownID)
	}
}

// TestDropReadsAnAbsentFileAsDone proves a retirement of an assignment that made no
// raw call, and a retirement that runs twice, both complete.
func TestDropReadsAnAbsentFileAsDone(t *testing.T) {
	t.Parallel()
	home, root, _ := fixtureHome(t)
	if err := Drop(home, root, knownID); err != nil {
		t.Fatalf("Drop of an absent file = %v, want no error", err)
	}
	writeRecordFile(t, home, root, knownID, "t\tsed\n")
	if err := Drop(home, root, knownID); err != nil {
		t.Fatal(err)
	}
	if err := Drop(home, root, knownID); err != nil {
		t.Fatalf("second Drop = %v, want no error", err)
	}
}

// TestDropRefusesAnIdentifierThatIsNotAnAssignment proves the drop never removes a
// path a malformed identifier composes, because the removal is unrecoverable.
func TestDropRefusesAnIdentifierThatIsNotAnAssignment(t *testing.T) {
	t.Parallel()
	home, root, _ := fixtureHome(t)
	writeRecordFile(t, home, root, knownID, "t\tsed\n")
	if err := Drop(home, root, filepath.Join("..", "census")); err == nil {
		t.Fatal("Drop accepted an identifier that is not an assignment")
	}
	if counts, err := Counts(home, root); err != nil || counts[knownID] != 1 {
		t.Fatalf("the refused drop changed the records: %v, %v", counts, err)
	}
}

// TestHeadBreakdownCountsEachHeadAndSortsByCount proves the landing's evidence line
// states one count for each verb head, largest first, and settles a tie by the head
// name so the text is stable between runs.
func TestHeadBreakdownCountsEachHeadAndSortsByCount(t *testing.T) {
	t.Parallel()
	home, root, _ := fixtureHome(t)
	writeRecordFile(t, home, root, knownID, "t\tsed\nt\tpython3\nt\tsed\nt\tsed\nt\tawk\nt\tpython3\n")
	if got := HeadBreakdown(home, root, knownID); got != "sed=3,python3=2,awk=1" {
		t.Fatalf("HeadBreakdown = %q, want %q", got, "sed=3,python3=2,awk=1")
	}
}

// TestHeadBreakdownSortsATieByHeadName proves two heads with the same count print in
// name order, not in the order the file happens to hold them.
func TestHeadBreakdownSortsATieByHeadName(t *testing.T) {
	t.Parallel()
	home, root, _ := fixtureHome(t)
	writeRecordFile(t, home, root, knownID, "t\tsed\nt\tawk\nt\tsed\nt\tawk\n")
	if got := HeadBreakdown(home, root, knownID); got != "awk=2,sed=2" {
		t.Fatalf("HeadBreakdown = %q, want %q", got, "awk=2,sed=2")
	}
}

// TestHeadBreakdownReadsTheSecondTabField proves the head comes from the field the
// writer puts it in, and that a last line with no newline is still one record.
func TestHeadBreakdownReadsTheSecondTabField(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, body, want string
	}{
		{"one closed line", "t\tsed\n", "sed=1"},
		{"one unterminated line", "t\tsed", "sed=1"},
		{"a third field is not the head", "t\tsed\textra\n", "sed=1"},
		{"a line with no head counts under none", "t\n\tsed\n", "sed=1"},
		{"an empty head counts under none", "t\t\n", ""},
		{"an empty head beside a real head", "t\t\nt\tsed\n", "sed=1"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			home, root, _ := fixtureHome(t)
			writeRecordFile(t, home, root, knownID, tc.body)
			if got := HeadBreakdown(home, root, knownID); got != tc.want {
				t.Fatalf("HeadBreakdown = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestHeadBreakdownReadsAnUnreadableCensusAsEmpty proves the breakdown is evidence
// beside the landing and never a condition on it: an absent directory, an absent
// file, and an empty file each render no text.
func TestHeadBreakdownReadsAnUnreadableCensusAsEmpty(t *testing.T) {
	t.Parallel()
	home, root, _ := fixtureHome(t)
	if got := HeadBreakdown(home, root, knownID); got != "" {
		t.Fatalf("HeadBreakdown on an absent directory = %q, want no text", got)
	}
	writeRecordFile(t, home, root, unknownID, "t\tsed\n")
	if got := HeadBreakdown(home, root, knownID); got != "" {
		t.Fatalf("HeadBreakdown on an absent file = %q, want no text", got)
	}
	writeRecordFile(t, home, root, knownID, "")
	if got := HeadBreakdown(home, root, knownID); got != "" {
		t.Fatalf("HeadBreakdown on an empty file = %q, want no text", got)
	}
}

// TestHeadBreakdownRefusesAFifoWithoutBlocking proves the breakdown has the same
// file-type posture as Counts: a refused file type renders no text and never holds
// the landing open on a reader that has no writer.
func TestHeadBreakdownRefusesAFifoWithoutBlocking(t *testing.T) {
	t.Parallel()
	home, root, _ := fixtureHome(t)
	dir := Dir(home, root)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(filepath.Join(dir, knownID), 0o600); err != nil {
		capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
	}
	done := make(chan string, 1)
	go func() {
		done <- HeadBreakdown(home, root, knownID)
	}()
	select {
	case got := <-done:
		if got != "" {
			t.Fatalf("HeadBreakdown on a FIFO record = %q, want no text", got)
		}
	case <-time.After(time.Second):
		t.Fatal("HeadBreakdown blocked on a FIFO record file")
	}
}

// TestHeadBreakdownSanitizesAForeignHead proves a head that a foreign writer put in
// the record file cannot move the terminal cursor when the landing prints the line.
func TestHeadBreakdownSanitizesAForeignHead(t *testing.T) {
	t.Parallel()
	home, root, _ := fixtureHome(t)
	writeRecordFile(t, home, root, knownID, "t\tse\x1b[2Kd\n")
	if got := HeadBreakdown(home, root, knownID); strings.ContainsRune(got, '\x1b') {
		t.Fatalf("HeadBreakdown = %q, want the control character removed", got)
	}
}
