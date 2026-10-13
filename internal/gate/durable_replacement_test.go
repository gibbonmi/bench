package gate

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

// gateWriter is one gate record writer: the entry it replaces in its directory, and the
// encoded record that entry holds after the rename.
type gateWriter struct {
	name, entry, want string
	write             func(durablefile.Replacer, string) error
}

// replacementVerdict is the record that the verdict and evidence writers encode, and
// replacementVerdictJSON is its encoding.
var replacementVerdict = verdictRecord{Schema: 1, State: Ready, Status: "green", Tree: "abc123", Oracle: "feed", RecordedAt: "2026-10-12T00:00:00Z"}

const replacementVerdictJSON = `{"schema":1,"state":"ready","status":"green","tree":"abc123","oracle":"feed","recorded_at":"2026-10-12T00:00:00Z"}` + "\n"

var (
	verdictWriter = gateWriter{
		name:  "verdict",
		entry: benchgit.GateCacheFile,
		want:  replacementVerdictJSON,
		write: func(replacer durablefile.Replacer, dir string) error {
			return durableReplace(replacer, dir, replacementVerdict)
		},
	}
	evidenceWriter = gateWriter{
		name:  "evidence",
		entry: "entry",
		want:  replacementVerdictJSON,
		write: func(replacer durablefile.Replacer, dir string) error {
			return durableReplaceAt(replacer, dir, "entry", replacementVerdict)
		},
	}
	laneWriter = gateWriter{
		name:  "lane",
		entry: laneRecordFile,
		want:  `{"schema":1,"tree":"abc123","recorded_at":"2026-10-12T00:00:00Z","lane":"fast","outcome":"pass","run_binary":"beef"}` + "\n",
		write: func(replacer durablefile.Replacer, dir string) error {
			return writeLaneRecord(replacer, dir, LaneResult{Outcome: lanePass, Tree: "abc123", Lane: "fast", RunBinary: "beef", RecordedAt: time.Date(2026, 10, 12, 0, 0, 0, 0, time.UTC)})
		},
	}
)

// D16: the verdict writer returns the injected directory-sync failure of its replacement
// after the rename, and the cache holds the encoded record at mode 0600.
func TestVerdictReplacement(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	cause := errors.New("injected verdict directory-sync failure")
	rec := &durablefiletest.Recorder{Fail: map[string]error{"directory-sync": cause}}

	err := verdictWriter.write(rec.Replacer(), dir)

	durablefiletest.RequireFailure(t, err, "directory-sync", true, cause)
	durablefiletest.RequireFile(t, filepath.Join(dir, verdictWriter.entry), 0o600, verdictWriter.want)
}

// D17: the evidence writer and the lane writer return the injected directory-sync failure
// of their shared store replacement, and each entry holds its encoded record at mode 0600.
func TestEvidenceReplacement(t *testing.T) {
	t.Parallel()
	for _, tc := range []gateWriter{evidenceWriter, laneWriter} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			cause := errors.New("injected " + tc.name + " directory-sync failure")
			rec := &durablefiletest.Recorder{Fail: map[string]error{"directory-sync": cause}}

			durablefiletest.RequireFailure(t, tc.write(rec.Replacer(), dir), "directory-sync", true, cause)
			durablefiletest.RequireFile(t, filepath.Join(dir, tc.entry), 0o600, tc.want)
		})
	}
}

// D36: a rename failure of each gate record writer returns an unpublished failure that
// wraps its cause and leaves no entry.
func TestGatePreRenameFailure(t *testing.T) {
	t.Parallel()
	for _, tc := range []gateWriter{verdictWriter, evidenceWriter, laneWriter} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			cause := errors.New("injected " + tc.name + " rename failure")
			rec := &durablefiletest.Recorder{Fail: map[string]error{"rename": cause}}

			durablefiletest.RequireFailure(t, tc.write(rec.Replacer(), dir), "rename", false, cause)
			durablefiletest.RequireAbsent(t, filepath.Join(dir, tc.entry))
		})
	}
}

// D30: an evidence directory with a mode other than private refuses the green evidence
// before the replacement starts, and no evidence entry exists.
func TestGateEvidenceDirectoryRefusal(t *testing.T) {
	t.Parallel()
	root := gittest.RepoOnBranch(t, "main")
	dir := filepath.Join(root, ".git", evidenceDestination)
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	plan := subject{Tree: replacementVerdict.Tree, Oracle: replacementVerdict.Oracle}
	rec := &durablefiletest.Recorder{}

	err := retainGreen(rec.Replacer(), root, plan, time.Date(2026, 10, 12, 0, 0, 0, 0, time.UTC))

	durablefiletest.RequireRefusal(t, err, rec.Trace)
	durablefiletest.RequireAbsent(t, filepath.Join(dir, evidenceName(plan)))
}

// D37: a failed pending write prints its classified cause and starts no oracle.
func TestGatePendingDurabilityFailure(t *testing.T) {
	root := outcomeFixture(t)
	cause := errors.New("injected pending directory-sync failure")
	faults := newFaultPlan(map[durablefiletest.Fault]map[string]error{{Destination: benchgit.GateCacheFile, Ordinal: 1}: {"directory-sync": cause}})

	result, stderr := executeWithFaults(root, faults)

	if _, err := os.Stat(runCountWitness.path(t, root)); !os.IsNotExist(err) {
		t.Fatalf("run counter stat = %v, want no oracle start after the pending failure", err)
	}
	requireUnpersisted(t, result, 0, stderr)
	want := "gate pending persistence failed: " + faults.Failure(benchgit.GateCacheFile, "directory-sync", true, cause)
	if !strings.Contains(stderr, want) {
		t.Fatalf("stderr = %q, want %q", stderr, want)
	}
}

// D38: a failed terminal write after a real green, red, or timed-out oracle result keeps
// its classified cause at ActionExit 1, and the recovery leaves the pending record.
func TestGateTerminalDurabilityFailure(t *testing.T) {
	for _, tc := range []struct {
		name, marker, category string
		timeout                time.Duration
		gateExit               int
	}{
		{name: "green", category: "gate final persistence failed"},
		{name: "red", marker: ".gate-red", category: "gate final persistence failed", gateExit: 7},
		{name: "timeout", marker: ".gate-sleep", category: "gate timeout persistence failed", timeout: 100 * time.Millisecond, gateExit: 124},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := failureOutcomeFixture(t)
			if tc.marker != "" {
				outcomeWrite(t, root, tc.marker, "\n", 0o644)
			}
			if tc.timeout != 0 {
				previousTimeout := gateTimeout
				gateTimeout = tc.timeout
				t.Cleanup(func() { gateTimeout = previousTimeout })
			}
			cause := errors.New("injected terminal directory-sync failure")
			faults := newFaultPlan(map[durablefiletest.Fault]map[string]error{{Destination: benchgit.GateCacheFile, Ordinal: 2}: {"directory-sync": cause}})

			result, stderr := executeWithFaults(root, faults)

			requireUnpersisted(t, result, tc.gateExit, stderr)
			want := tc.category + ": " + faults.Failure(benchgit.GateCacheFile, "directory-sync", true, cause)
			if !strings.Contains(stderr, want) || strings.Contains(stderr, "pending recovery failed") {
				t.Fatalf("stderr = %q, want %q and no recovery failure", stderr, want)
			}
			if got := outcomeRuns(t, root); got != 1 {
				t.Fatalf("runs = %d, want 1", got)
			}
			if inspection := Inspect(root); inspection.State != Pending || inspection.ReusableGreen {
				t.Fatalf("inspection after the recovery = %#v, want pending and non-reusable", inspection)
			}
		})
	}
}

// D39: a failed green-evidence write prints its classified cause at ActionExit 1 and
// writes no final verdict.
func TestGateEvidenceDurabilityFailure(t *testing.T) {
	root := outcomeFixture(t)
	cause := errors.New("injected evidence directory-sync failure")
	faults := newFaultPlan(map[durablefiletest.Fault]map[string]error{{Destination: evidenceDestination, Ordinal: 1}: {"directory-sync": cause}})

	result, stderr := executeWithFaults(root, faults)

	requireUnpersisted(t, result, 0, stderr)
	want := "gate evidence persistence failed: " + faults.Failure(evidenceDestination, "directory-sync", true, cause)
	if !strings.Contains(stderr, want) {
		t.Fatalf("stderr = %q, want %q", stderr, want)
	}
	if got := faults.Count(benchgit.GateCacheFile); got != 1 {
		t.Fatalf("cache replacements = %d, want only the pending record", got)
	}
}

// D40: a failed final write and an independently failed pending recovery report both
// causes. The final rename completed, so the cache holds the new verdict and not the
// pending record.
func TestGatePendingRecoveryFailure(t *testing.T) {
	root := outcomeFixture(t)
	final := errors.New("injected final directory-sync failure")
	recovery := errors.New("injected recovery rename failure")
	faults := newFaultPlan(map[durablefiletest.Fault]map[string]error{
		{Destination: benchgit.GateCacheFile, Ordinal: 2}: {"directory-sync": final},
		{Destination: benchgit.GateCacheFile, Ordinal: 3}: {"rename": recovery},
	})

	result, stderr := executeWithFaults(root, faults)

	requireUnpersisted(t, result, 0, stderr)
	want := "gate final persistence failed: " + faults.Failure(benchgit.GateCacheFile, "directory-sync", true, final) +
		"; pending recovery failed: " + faults.Failure(benchgit.GateCacheFile, "rename", false, recovery)
	if !strings.Contains(stderr, want) {
		t.Fatalf("stderr = %q, want %q", stderr, want)
	}
	if inspection := Inspect(root); inspection.State != Ready || inspection.Status != "green" {
		t.Fatalf("inspection after the failed recovery = %#v, want the visible green verdict", inspection)
	}
}

// The interrupted-green demotion on an unavailable lock and on a contended lock reports
// its failed persistence beside the lock refusal.
func TestGateDemotionDurabilityFailure(t *testing.T) {
	for _, tc := range []struct {
		name, refusal string
		block         func(t *testing.T, root string)
	}{
		{
			name:    "lock unavailable",
			refusal: "gate lock unavailable",
			block: func(t *testing.T, root string) {
				lock := gateLockPath(outcomeGit(t, root, "rev-parse", "--absolute-git-dir"))
				if err := os.Remove(lock); err != nil {
					t.Fatal(err)
				}
				if err := os.Mkdir(lock, 0o700); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name:    "lock contention",
			refusal: "gate execution already in progress",
			block:   func(t *testing.T, root string) { startGateLockHolder(t, root) },
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := outcomeFixture(t)
			if result := Execute(context.Background(), root, &bytes.Buffer{}, &bytes.Buffer{}); result.ActionExit != 0 {
				t.Fatalf("green result = %#v", result)
			}
			tc.block(t, root)
			demotion := errors.New("injected demotion rename failure")
			recovery := errors.New("injected demotion recovery rename failure")
			faults := newFaultPlan(map[durablefiletest.Fault]map[string]error{
				{Destination: benchgit.GateCacheFile, Ordinal: 1}: {"rename": demotion},
				{Destination: benchgit.GateCacheFile, Ordinal: 2}: {"rename": recovery},
			})

			result, stderr := executeWithFaults(root, faults)

			requireUnpersisted(t, result, 0, stderr)
			want := tc.refusal + "; interrupted-green demotion failed: " + faults.Failure(benchgit.GateCacheFile, "rename", false, demotion) +
				"; pending recovery failed: " + faults.Failure(benchgit.GateCacheFile, "rename", false, recovery)
			if !strings.Contains(stderr, want) {
				t.Fatalf("stderr = %q, want %q", stderr, want)
			}
		})
	}
}

// evidenceDestination names every entry of the evidence store, whose file names are
// content addresses.
const evidenceDestination = "bench-gate-evidence"

// newFaultPlan routes each replacement of one execution by its destination: the file name
// of a cache record, or evidenceDestination for every evidence entry.
func newFaultPlan(fail map[durablefiletest.Fault]map[string]error) *durablefiletest.Router {
	return &durablefiletest.Router{Fail: fail, Destination: func(path string) string {
		if filepath.Base(filepath.Dir(path)) == evidenceDestination {
			return evidenceDestination
		}
		return filepath.Base(path)
	}}
}

// executeWithFaults runs one fresh execution of root through faults and returns its
// result and stderr.
func executeWithFaults(root string, faults *durablefiletest.Router) (Result, string) {
	var stdout, stderr bytes.Buffer
	result := executeSubjectReplacing(context.Background(), root, root, &stdout, &stderr, nil, forceRun, newGateEvaluation(root), nil, "", faults.Replacer())
	return result, stderr.String()
}

func requireUnpersisted(t *testing.T, result Result, gateExit int, stderr string) {
	t.Helper()
	if result.ActionExit != 1 || result.GateExit != gateExit || result.Inspection.ReusableGreen {
		t.Fatalf("result = %#v, stderr=%q, want ActionExit 1, GateExit %d, and no reusable green", result, stderr, gateExit)
	}
}
