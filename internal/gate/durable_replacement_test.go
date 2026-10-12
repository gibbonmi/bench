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
)

// D16: the verdict writer returns the injected directory-sync failure of its replacement
// after the rename, and the cache holds the encoded record at mode 0600.
func TestVerdictReplacement(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	cause := errors.New("injected verdict directory-sync failure")
	rec := &durablefiletest.Recorder{Fail: map[string]error{"directory-sync": cause}}
	record := verdictRecord{Schema: 1, State: Ready, Status: "green", Tree: "abc123", Oracle: "feed", RecordedAt: "2026-10-12T00:00:00Z"}

	err := durableReplace(rec.Replacer(), dir, record)

	requirePublishedFailure(t, err, cause)
	requireRecordFile(t, filepath.Join(dir, benchgit.GateCacheFile),
		`{"schema":1,"state":"ready","status":"green","tree":"abc123","oracle":"feed","recorded_at":"2026-10-12T00:00:00Z"}`+"\n")
}

// D17: the evidence writer and the lane writer return the injected directory-sync failure
// of their shared store replacement, and each entry holds its encoded record at mode 0600.
func TestEvidenceReplacement(t *testing.T) {
	t.Parallel()
	recordedAt := time.Date(2026, 10, 12, 0, 0, 0, 0, time.UTC)
	for _, tc := range []struct {
		name, entry, want string
		write             func(durablefile.Replacer, string) error
	}{
		{
			name:  "evidence",
			entry: "entry",
			want:  `{"schema":1,"state":"ready","status":"green","tree":"abc123","oracle":"feed","recorded_at":"2026-10-12T00:00:00Z"}` + "\n",
			write: func(replacer durablefile.Replacer, dir string) error {
				return durableReplaceAt(replacer, dir, "entry", verdictRecord{Schema: 1, State: Ready, Status: "green", Tree: "abc123", Oracle: "feed", RecordedAt: "2026-10-12T00:00:00Z"})
			},
		},
		{
			name:  "lane",
			entry: laneRecordFile,
			want:  `{"schema":1,"tree":"abc123","recorded_at":"2026-10-12T00:00:00Z","lane":"fast","outcome":"pass","run_binary":"beef"}` + "\n",
			write: func(replacer durablefile.Replacer, dir string) error {
				return writeLaneRecord(replacer, dir, LaneResult{Outcome: lanePass, Tree: "abc123", Lane: "fast", RunBinary: "beef", RecordedAt: recordedAt})
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			cause := errors.New("injected " + tc.name + " directory-sync failure")
			rec := &durablefiletest.Recorder{Fail: map[string]error{"directory-sync": cause}}

			requirePublishedFailure(t, tc.write(rec.Replacer(), dir), cause)
			requireRecordFile(t, filepath.Join(dir, tc.entry), tc.want)
		})
	}
}

// D37: a failed pending write prints its classified cause and starts no oracle.
func TestGatePendingDurabilityFailure(t *testing.T) {
	root := outcomeFixture(t)
	cause := errors.New("injected pending directory-sync failure")
	faults := newFaultPlan(map[fault]map[string]error{{benchgit.GateCacheFile, 1}: {"directory-sync": cause}})

	result, stderr := executeWithFaults(root, faults)

	if _, err := os.Stat(runCountWitness.path(t, root)); !os.IsNotExist(err) {
		t.Fatalf("run counter stat = %v, want no oracle start after the pending failure", err)
	}
	requireUnpersisted(t, result, 0, stderr)
	want := "gate pending persistence failed: " + faults.failure(benchgit.GateCacheFile, "directory-sync", true, cause)
	if !strings.Contains(stderr, want) {
		t.Fatalf("stderr = %q, want %q", stderr, want)
	}
}

// D38: a failed final write after a real green or red oracle result keeps its classified
// cause at ActionExit 1, and the recovery leaves the pending record.
func TestGateTerminalDurabilityFailure(t *testing.T) {
	for _, tc := range []struct {
		name     string
		red      bool
		gateExit int
	}{{"green", false, 0}, {"red", true, 7}} {
		t.Run(tc.name, func(t *testing.T) {
			root := outcomeFixture(t)
			if tc.red {
				outcomeWrite(t, root, ".gate-red", "\n", 0o644)
			}
			cause := errors.New("injected final directory-sync failure")
			faults := newFaultPlan(map[fault]map[string]error{{benchgit.GateCacheFile, 2}: {"directory-sync": cause}})

			result, stderr := executeWithFaults(root, faults)

			requireUnpersisted(t, result, tc.gateExit, stderr)
			want := "gate final persistence failed: " + faults.failure(benchgit.GateCacheFile, "directory-sync", true, cause)
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
	faults := newFaultPlan(map[fault]map[string]error{{evidenceDestination, 1}: {"directory-sync": cause}})

	result, stderr := executeWithFaults(root, faults)

	requireUnpersisted(t, result, 0, stderr)
	want := "gate evidence persistence failed: " + faults.failure(evidenceDestination, "directory-sync", true, cause)
	if !strings.Contains(stderr, want) {
		t.Fatalf("stderr = %q, want %q", stderr, want)
	}
	if got := faults.count[benchgit.GateCacheFile]; got != 1 {
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
	faults := newFaultPlan(map[fault]map[string]error{
		{benchgit.GateCacheFile, 2}: {"directory-sync": final},
		{benchgit.GateCacheFile, 3}: {"rename": recovery},
	})

	result, stderr := executeWithFaults(root, faults)

	requireUnpersisted(t, result, 0, stderr)
	want := "gate final persistence failed: " + faults.failure(benchgit.GateCacheFile, "directory-sync", true, final) +
		"; pending recovery failed: " + faults.failure(benchgit.GateCacheFile, "rename", false, recovery)
	if !strings.Contains(stderr, want) {
		t.Fatalf("stderr = %q, want %q", stderr, want)
	}
	if inspection := Inspect(root); inspection.State != Ready || inspection.Status != "green" {
		t.Fatalf("inspection after the failed recovery = %#v, want the visible green verdict", inspection)
	}
}

// The interrupted-green demotion on a lock failure reports its failed persistence beside
// the lock refusal.
func TestGateDemotionDurabilityFailure(t *testing.T) {
	root := outcomeFixture(t)
	if result := Execute(context.Background(), root, &bytes.Buffer{}, &bytes.Buffer{}); result.ActionExit != 0 {
		t.Fatalf("green result = %#v", result)
	}
	lock := gateLockPath(outcomeGit(t, root, "rev-parse", "--absolute-git-dir"))
	if err := os.Remove(lock); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(lock, 0o700); err != nil {
		t.Fatal(err)
	}
	demotion := errors.New("injected demotion rename failure")
	recovery := errors.New("injected demotion recovery rename failure")
	faults := newFaultPlan(map[fault]map[string]error{
		{benchgit.GateCacheFile, 1}: {"rename": demotion},
		{benchgit.GateCacheFile, 2}: {"rename": recovery},
	})

	result, stderr := executeWithFaults(root, faults)

	requireUnpersisted(t, result, 0, stderr)
	want := "gate lock unavailable; interrupted-green demotion failed: " + faults.failure(benchgit.GateCacheFile, "rename", false, demotion) +
		"; pending recovery failed: " + faults.failure(benchgit.GateCacheFile, "rename", false, recovery)
	if !strings.Contains(stderr, want) {
		t.Fatalf("stderr = %q, want %q", stderr, want)
	}
}

// evidenceDestination names every entry of the evidence store, whose file names are
// content addresses.
const evidenceDestination = "bench-gate-evidence"

// fault names one replacement by its destination and by its ordinal among the renames of
// that destination, counting from 1.
type fault struct {
	destination string
	ordinal     int
}

// faultPlan routes each replacement of one execution through a durablefiletest.Recorder
// that fails the stages the plan names for that replacement. The replacement is known
// from its rename, so a plan fails the rename and the stages after it.
type faultPlan struct {
	fail  map[fault]map[string]error
	count map[string]int
	paths map[string]string
}

func newFaultPlan(fail map[fault]map[string]error) *faultPlan {
	return &faultPlan{fail: fail, count: map[string]int{}, paths: map[string]string{}}
}

func (p *faultPlan) replacer() durablefile.Replacer {
	native := durablefile.Native()
	current := native
	return durablefile.Replacer{
		CreateTemp: native.CreateTemp,
		Rename: func(oldpath, newpath string) error {
			name := filepath.Base(newpath)
			if filepath.Base(filepath.Dir(newpath)) == evidenceDestination {
				name = evidenceDestination
			}
			p.count[name]++
			p.paths[name] = newpath
			current = (&durablefiletest.Recorder{Fail: p.fail[fault{name, p.count[name]}]}).Replacer()
			return current.Rename(oldpath, newpath)
		},
		OpenDirectory: func(name string) (durablefile.Directory, error) { return current.OpenDirectory(name) },
		Remove:        native.Remove,
	}
}

// failure is the message of the injected cause at stage of the last replacement of
// destination.
func (p *faultPlan) failure(destination, stage string, published bool, cause error) string {
	return (&durablefile.Error{Path: p.paths[destination], Stage: stage, Published: published, Err: cause}).Error()
}

// executeWithFaults runs one fresh execution of root through faults and returns its
// result and stderr.
func executeWithFaults(root string, faults *faultPlan) (Result, string) {
	var stdout, stderr bytes.Buffer
	result := executeSubjectReplacing(context.Background(), root, root, &stdout, &stderr, nil, forceRun, newGateEvaluation(root), nil, "", faults.replacer())
	return result, stderr.String()
}

func requireUnpersisted(t *testing.T, result Result, gateExit int, stderr string) {
	t.Helper()
	if result.ActionExit != 1 || result.GateExit != gateExit || result.Inspection.ReusableGreen {
		t.Fatalf("result = %#v, stderr=%q, want ActionExit 1, GateExit %d, and no reusable green", result, stderr, gateExit)
	}
}

func requirePublishedFailure(t *testing.T, err, cause error) {
	t.Helper()
	var replaceErr *durablefile.Error
	if !errors.As(err, &replaceErr) || !replaceErr.Published || replaceErr.Stage != "directory-sync" || !errors.Is(err, cause) {
		t.Fatalf("err = %v, want a published directory-sync *durablefile.Error that wraps %v", err, cause)
	}
}

func requireRecordFile(t *testing.T, path, want string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("mode = %v, want 0600", got)
	}
	if got := string(outcomeRead(t, path)); got != want {
		t.Fatalf("record = %q, want %q", got, want)
	}
}
