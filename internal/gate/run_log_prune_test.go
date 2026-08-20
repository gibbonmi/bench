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
		t.Fatalf("surviving records\n got %d: %v\nwant %d: %v", len(got), got, len(want), want)
	}
}

// H27: a gate run over a log directory holding more than 20 records leaves exactly
// the newest 20 — asserted by the identity of the survivors, so a pruner that
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
	if len(want) != gateLogRetainedRecords {
		t.Fatalf("expectation names %d records, retention is %d", len(want), gateLogRetainedRecords)
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

// H29: a file in .logs outside the gate's name shape survives pruning, and a
// candidate that is not a regular file is left alone rather than removed — the
// pruner stats before it removes.
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

	// Special files carrying the gate's own name shape: age 0 and 1 would both be
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
// verdict and exit code unchanged — the run's own finish record still lands with
// the Result's exit, and finish returns rather than reporting an error upward.
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

// H27/story 39: pruning is a side effect of a gate run, not of a separate chore, so this
// drives the closure beginGateRunLog hands back rather than calling the pruner. A build
// that keeps the pruner correct but never wires it into the run leaves .logs unbounded
// exactly as before, and every direct-call assertion above still passes.
func TestGateRunPrunesThroughTheClosureBeginHandsBack(t *testing.T) {
	root := newLoggingPruneRoot(t)
	seeded := seedRecords(t, root, ages(0, 24)...)

	var stderr bytes.Buffer
	_, finish := beginGateRunLog(context.Background(), root, &stderr, "dev")
	if strings.Contains(stderr.String(), "progress logging unavailable") {
		t.Fatalf("the run never opened a record, so this asserts nothing: %q", stderr.String())
	}
	finish(Result{GateExit: 0, ActionExit: 0})

	survivors := logDirEntries(t, root)
	if len(survivors) != gateLogRetainedRecords {
		t.Fatalf("after one gate run: %d records retained, want %d\n%v", len(survivors), gateLogRetainedRecords, survivors)
	}
	// The record this run wrote is the newest of all, so it must have survived its own
	// pruning: a pruner that counts before the write would drop it.
	seededNames := strings.Join(seeded, "\n")
	own := 0
	for _, name := range survivors {
		if !strings.Contains(seededNames, name) {
			own++
		}
	}
	if own != 1 {
		t.Errorf("survivors hold %d records this run wrote, want exactly 1\n%v", own, survivors)
	}
}

// H27's ordering rule when two runs share a start instant. seededRun spaces records a
// minute apart, so the tiebreak below sort.Slice's comparator is unreached there; without
// it the sort is unstable on ties and two gates starting in the same instant get
// nondeterministic retention.
func TestPruneGateRunLogsBreaksTimestampTiesDeterministically(t *testing.T) {
	root := newPruneRoot(t)
	// The tie sits at the retention cut: 19 newer records plus 3 sharing one instant is
	// 22, so exactly two of the tied three are dropped. Which two is a fact about the
	// name, not about readdir order.
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

	if got := len(logDirEntries(t, root)); got != gateLogRetainedRecords {
		t.Fatalf("retained %d, want %d", got, gateLogRetainedRecords)
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
	previous := gateLogPathIgnored
	gateLogPathIgnored = func(string) bool { return true }
	t.Cleanup(func() { gateLogPathIgnored = previous })
	return root
}
