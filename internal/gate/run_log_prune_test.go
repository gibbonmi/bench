package gate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/capability"
)

var pruneEpoch = time.Date(2026, 3, 4, 5, 6, 7, 0, time.UTC)

// seededRun names the record a seeded run would have written. Ages are one minute
// apart so the retained set is identifiable by name, not merely countable.
func seededRun(age int) string {
	return gateLogRunToken(pruneEpoch.Add(time.Duration(age)*time.Minute), 4000+age)
}

func newPruneRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(gateLogDir(root), 0o700); err != nil {
		t.Fatal(err)
	}
	return root
}

func seedRecords(t *testing.T, root string, ages ...int) []string {
	t.Helper()
	names := make([]string, 0, len(ages))
	for _, age := range ages {
		run := seededRun(age)
		if err := os.WriteFile(gateLogRecordPath(root, run), []byte("{}\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		names = append(names, gateLogRecordName(run))
	}
	return names
}

// seedRunPairs seeds what a run that opened its stream leaves behind: a record and a
// stream beside it, under one run token.
func seedRunPairs(t *testing.T, root string, ages ...int) []string {
	t.Helper()
	names := seedRecords(t, root, ages...)
	for _, age := range ages {
		run := seededRun(age)
		if err := os.WriteFile(gateLogStreamPath(root, run), []byte("[vet] a line\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		names = append(names, gateLogStreamName(run))
	}
	return names
}

func ages(from, to int) []int {
	out := make([]int, 0, to-from+1)
	for age := from; age <= to; age++ {
		out = append(out, age)
	}
	return out
}

func logDirEntries(t *testing.T, root string) []string {
	t.Helper()
	entries, err := os.ReadDir(gateLogDir(root))
	if err != nil {
		t.Fatal(err)
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	sort.Strings(names)
	return names
}

func assertLogDirHolds(t *testing.T, root string, want []string) {
	t.Helper()
	sort.Strings(want)
	got := logDirEntries(t, root)
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("surviving files\n got %d: %v\nwant %d: %v", len(got), got, len(want), want)
	}
}

// H27: a gate run over a log directory holding more than 20 records leaves exactly
// the newest 20. This test asserts the identity of the survivors, so a pruner that
// retains a bounded-but-wrong set cannot pass on the count alone.
func TestPruneGateRunLogsRetainsExactlyTheNewestRecords(t *testing.T) {
	root := newPruneRoot(t)
	seedRecords(t, root, ages(0, 24)...)
	current := seededRun(25)
	seedRecords(t, root, 25)

	pruneGateRunLogs(root, current, io.Discard)

	want := []string{gateLogRecordName(current)}
	for _, age := range ages(6, 24) {
		want = append(want, gateLogRecordName(seededRun(age)))
	}
	if len(want) != gateLogRetainedRuns {
		t.Fatalf("expectation names %d runs, retention is %d", len(want), gateLogRetainedRuns)
	}
	assertLogDirHolds(t, root, want)
}

// BG25: the retention counts runs, not files. Twenty-one runs that each opened a stream
// leave forty-two files, and a pruner counting files would keep ten runs. The oldest
// run's two files go together, so no surviving table points at a stream that is gone.
func TestPruneGateRunLogsCountsRunsAndRemovesEachRunsPairTogether(t *testing.T) {
	root := newPruneRoot(t)
	seedRunPairs(t, root, ages(0, 19)...)
	current := seededRun(20)
	seedRunPairs(t, root, 20)

	pruneGateRunLogs(root, current, io.Discard)

	want := make([]string, 0, 2*gateLogRetainedRuns)
	for _, age := range ages(1, 20) {
		run := seededRun(age)
		want = append(want, gateLogRecordName(run), gateLogStreamName(run))
	}
	if len(want) != 2*gateLogRetainedRuns {
		t.Fatalf("expectation names %d files, retention is %d runs", len(want), gateLogRetainedRuns)
	}
	assertLogDirHolds(t, root, want)
}

// A run that opened no stream prunes on its record alone, so a directory mixing the two
// shapes still retains twenty runs. The oldest run here has a pair and the run above the
// cut has only its record; neither shape changes what the count keeps.
func TestPruneGateRunLogsCountsARunWithNoStreamAsOneRun(t *testing.T) {
	root := newPruneRoot(t)
	seedRunPairs(t, root, ages(0, 9)...)
	seedRecords(t, root, ages(10, 20)...)

	pruneGateRunLogs(root, seededRun(20), io.Discard)

	want := make([]string, 0)
	for _, age := range ages(1, 9) {
		run := seededRun(age)
		want = append(want, gateLogRecordName(run), gateLogStreamName(run))
	}
	for _, age := range ages(10, 20) {
		want = append(want, gateLogRecordName(seededRun(age)))
	}
	assertLogDirHolds(t, root, want)
}

// H28: the record the current run is writing survives its own pruning. The current
// run is seeded as the OLDEST record here, so only an explicit exclusion saves it.
func TestPruneGateRunLogsNeverRemovesTheCurrentRun(t *testing.T) {
	root := newPruneRoot(t)
	current := seededRun(0)
	seedRecords(t, root, 0)
	seedRecords(t, root, ages(1, 25)...)

	pruneGateRunLogs(root, current, io.Discard)

	want := []string{gateLogRecordName(current)}
	for _, age := range ages(6, 25) {
		want = append(want, gateLogRecordName(seededRun(age)))
	}
	assertLogDirHolds(t, root, want)
}

// H29: a file in .logs outside the gate's name shape survives pruning. A candidate
// that is not a regular file is also left alone rather than removed. The pruner
// stats before it removes.
func TestPruneGateRunLogsLeavesForeignEntriesAlone(t *testing.T) {
	root := newPruneRoot(t)
	logs := gateLogDir(root)
	seedRecords(t, root, ages(6, 25)...)

	foreign := []string{
		"notes.txt",
		"gate.jsonl",
		gateLogRecordName(""),
		gateLogRecordName("not-a-run"),
		gateLogRecordName("20260304T050607.000000000Z"),
		gateLogRecordName("20260304T050607.000000000Z-abc"),
		gateLogRecordName("20260304T050607.000000000Z-004000"),
		gateLogRecordName("20260304T050607Z-4000"),
	}
	for _, name := range foreign {
		if err := os.WriteFile(filepath.Join(logs, name), []byte("keep me\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	// These special files carry the gate's own name shape. Age 0 and 1 would both be
	// pruned on count, so removing them is exactly the failure the stat prevents.
	dangling := gateLogRecordName(seededRun(0))
	if err := os.Symlink(filepath.Join(logs, "gone.jsonl"), filepath.Join(logs, dangling)); err != nil {
		t.Fatal(err)
	}
	fifo := gateLogRecordName(seededRun(1))
	if err := syscall.Mkfifo(filepath.Join(logs, fifo), 0o600); err != nil {
		t.Fatal(err)
	}

	current := seededRun(26)
	seedRecords(t, root, 26)

	pruneGateRunLogs(root, current, io.Discard)

	want := append([]string{gateLogRecordName(current), dangling, fifo}, foreign...)
	for _, age := range ages(7, 25) {
		want = append(want, gateLogRecordName(seededRun(age)))
	}
	assertLogDirHolds(t, root, want)
}

// H34: a pruning failure writes exactly one stderr warning and leaves the gate's
// verdict and exit code unchanged. The run's own finish record still lands with the
// Result's exit, and finish returns rather than reporting an error upward.
func TestGateRunLogFinishSurvivesPruningFailure(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T, root string)
	}{
		{"unreadable-directory", func(t *testing.T, root string) {
			logs := gateLogDir(root)
			t.Cleanup(func() { _ = os.Chmod(logs, 0o700) })
			if err := os.Chmod(logs, 0o500); err != nil {
				capability.Capability(t, capability.Privilege, fmt.Sprintf("cannot strip directory permissions: %v", err))
			}
			if err := os.Remove(gateLogRecordPath(root, seededRun(0))); err == nil {
				capability.Capability(t, capability.Privilege, "mode 0o500 directory is still writable by this user")
			}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := newPruneRoot(t)
			seedRecords(t, root, ages(0, 24)...)
			current := gateLogRunToken(pruneEpoch.Add(30*time.Minute), 9999)
			file, err := os.OpenFile(gateLogRecordPath(root, current), os.O_WRONLY|os.O_CREATE|os.O_EXCL|os.O_APPEND, 0o600)
			if err != nil {
				t.Fatal(err)
			}
			tt.setup(t, root)

			var stderr bytes.Buffer
			log := &gateRunLog{file: file, stderr: &stderr, run: current, root: root, started: time.Now().UTC()}
			log.finish(Result{GateExit: 3, ActionExit: 7})

			warnings := 0
			for _, line := range strings.Split(strings.TrimSuffix(stderr.String(), "\n"), "\n") {
				if strings.Contains(line, "pruning") {
					warnings++
				}
			}
			if warnings != 1 {
				t.Fatalf("pruning warnings = %d, want 1; stderr:\n%s", warnings, stderr.String())
			}

			data, err := os.ReadFile(gateLogRecordPath(root, current))
			if err != nil {
				t.Fatal(err)
			}
			var record gateLogRecord
			last := strings.TrimSuffix(string(data), "\n")
			if err := json.Unmarshal([]byte(last[strings.LastIndex(last, "\n")+1:]), &record); err != nil {
				t.Fatal(err)
			}
			if record.Event != "gate.finish" || record.Exit == nil || *record.Exit != 7 {
				t.Fatalf("finish record = %+v, want gate.finish with exit 7", record)
			}
		})
	}
}

// H27/story 39: pruning is a side effect of a gate run, not of a separate chore. So
// this test drives the closure beginGateRunLog hands back, rather than calling the
// pruner. A build that keeps the pruner correct but never wires it into the run
// leaves .logs unbounded exactly as before. Every direct-call assertion above still
// passes.
func TestGateRunPrunesThroughTheClosureBeginHandsBack(t *testing.T) {
	root := newLoggingPruneRoot(t)
	seeded := seedRecords(t, root, ages(0, 24)...)

	var stderr bytes.Buffer
	_, finish := beginGateRunLog(context.Background(), root, &stderr, "dev")
	if strings.Contains(stderr.String(), "progress logging unavailable") {
		t.Fatalf("the run never opened a record, so this asserts nothing: %q", stderr.String())
	}
	finish(Result{GateExit: 0, ActionExit: 0})

	// The retention is a count of runs, and this run left both of its files, so the
	// survivors are counted by their run token rather than by name.
	survivors := logDirEntries(t, root)
	seededNames := strings.Join(seeded, "\n")
	runs, own := map[string]bool{}, map[string]bool{}
	for _, name := range survivors {
		run, _, ok := gateLogRunFromRecordName(name)
		if !ok {
			t.Fatalf("survivor %q is not a name the gate writes\n%v", name, survivors)
		}
		runs[run] = true
		if !strings.Contains(seededNames, name) {
			own[run] = true
		}
	}
	if len(runs) != gateLogRetainedRuns {
		t.Fatalf("after one gate run: %d runs retained, want %d\n%v", len(runs), gateLogRetainedRuns, survivors)
	}
	// The files this run wrote are the newest of all, so they must have survived their
	// own pruning. A pruner that counts before the write would drop them.
	if len(own) != 1 {
		t.Errorf("survivors hold %d runs this run wrote, want exactly 1\n%v", len(own), survivors)
	}
}

// H27's ordering rule when two runs share a start instant. seededRun spaces records a
// minute apart, so the tiebreak below sort.Slice's comparator is unreached there.
// Without it, the sort is unstable on ties, and two gates starting in the same
// instant get nondeterministic retention.
func TestPruneGateRunLogsBreaksTimestampTiesDeterministically(t *testing.T) {
	root := newPruneRoot(t)
	// The tie sits at the retention cut. 19 newer records plus 3 sharing one instant
	// is 22, so exactly two of the tied three are dropped. Which two is a fact about
	// the name, not about readdir order.
	instant := pruneEpoch.Add(-10 * time.Minute)
	tied := make([]string, 0, 3)
	for _, pid := range []int{7001, 7002, 7003} {
		run := gateLogRunToken(instant, pid)
		if err := os.WriteFile(gateLogRecordPath(root, run), []byte("{}\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		tied = append(tied, gateLogRecordName(run))
	}
	sort.Strings(tied)
	seedRecords(t, root, ages(0, 18)...)

	pruneGateRunLogs(root, "", io.Discard)
	survivors := strings.Join(logDirEntries(t, root), "\n")

	if got := len(logDirEntries(t, root)); got != gateLogRetainedRuns {
		t.Fatalf("retained %d runs, want %d", got, gateLogRetainedRuns)
	}
	if !strings.Contains(survivors, tied[2]) {
		t.Errorf("the greatest-named tied record was dropped; the comparator must keep it\n%s", survivors)
	}
	for _, dropped := range tied[:2] {
		if strings.Contains(survivors, dropped) {
			t.Errorf("tied record %s survived; only the greatest name may\n%s", dropped, survivors)
		}
	}
}

// newLoggingPruneRoot is newPruneRoot with the ignore precondition satisfied, stubbed
// rather than built — see gateLogPathIgnored for why this package's tests stand up no
// repository.
func newLoggingPruneRoot(t *testing.T) string {
	t.Helper()
	root := newPruneRoot(t)
	stubGateLogPathIgnored(t)
	return root
}

// stubGateLogPathIgnored satisfies the ignore precondition for this test's run.
func stubGateLogPathIgnored(t *testing.T) {
	t.Helper()
	previous := gateLogPathIgnored
	gateLogPathIgnored = func(string) bool { return true }
	t.Cleanup(func() { gateLogPathIgnored = previous })
}
