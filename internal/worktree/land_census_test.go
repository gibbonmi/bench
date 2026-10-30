package worktree

import (
	"github.com/gibbonmi/bench/internal/census"
	"github.com/gibbonmi/bench/internal/handoffdoc"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func recordRawCalls(t *testing.T, home, root, path string, n int) {
	t.Helper()
	recordRawCallsWithHead(t, home, root, path, "sed -i s/a/b/", n)
}

// recordRawCallsWithHead appends n raw-call records that one command text makes, which
// lets a test state a breakdown over more than one verb head.
func recordRawCallsWithHead(t *testing.T, home, root, path, command string, n int) {
	t.Helper()
	for range n {
		if err := census.Record(command+" "+filepath.Join(path, "owned.txt"), root, home, time.Now()); err != nil {
			t.Fatal(err)
		}
	}
}

// censusRecordPath names one assignment's record file.
func censusRecordPath(home, root, assignment string) string {
	return filepath.Join(census.Dir(home, root), assignment)
}

// TestLandCommandStatesTheCensusCountAndDropsTheRecords is EC20, the landing half of
// EC24, and BO60. Before the release step drops the records, the landing prints the
// raw-call count for each verb head, heaviest first, and the output breakdown on stderr,
// so the retro reads both from the run. The landed record carries the raw-call count as
// its last key, which the output records leave unchanged, and the release step leaves no
// record file for the retired assignment.
func TestLandCommandStatesTheCensusCountAndDropsTheRecords(t *testing.T) {
	t.Parallel()
	request := "census-landed-count"
	f := publicLandingFixture(t, request, "", "")
	recordRawCallsWithHead(t, f.home, f.root, f.creation.Path, "sed -i s/a/b/", 2)
	recordRawCallsWithHead(t, f.home, f.root, f.creation.Path, "awk -f x", 1)
	for _, size := range []int64{16938, 12040} {
		output := census.Output{Head: "bench worktree list", Lines: 55, Bytes: size, Spilled: true}
		if err := census.RecordOutput(f.home, f.root, f.creation.Assignment.ID, output, time.Now()); err != nil {
			t.Fatal(err)
		}
	}
	survivor := seedHandoffSections(t, f.root, f.creation.Assignment)
	r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...))
	if r.exit != 0 || !strings.HasSuffix(r.stdout, ",census=3}\n") {
		t.Fatalf("landed record = (%d, %q, %q), want the raw-call count census=3 as the last key", r.exit, r.stdout, r.stderr)
	}
	if !strings.Contains(r.stderr, "census heads{sed=2,awk=1}\n") || !strings.Contains(r.stderr, "census output{bench worktree list=2/28978}\n") {
		t.Fatalf("landing evidence = %q, want the head and output breakdowns on stderr", r.stderr)
	}
	if _, err := os.Stat(censusRecordPath(f.home, f.root, f.creation.Assignment.ID)); !os.IsNotExist(err) {
		t.Fatalf("the released landing kept the census record: %v", err)
	}
	requireHandoffSections(t, f.root, handoffdoc.MainKey, survivor)
}

// TestLandCommandPrintsNoHeadsLineWithoutRecords proves an assignment that made no raw
// call prints no breakdown at all, and still states the zero count in its record.
func TestLandCommandPrintsNoHeadsLineWithoutRecords(t *testing.T) {
	t.Parallel()
	request := "census-landed-no-heads"
	f := publicLandingFixture(t, request, "", "")
	r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...))
	if r.exit != 0 || strings.Contains(r.stderr, "census heads{") || !strings.HasSuffix(r.stdout, ",census=0}\n") {
		t.Fatalf("empty census landing = (%d, %q, %q), want no heads line and census=0", r.exit, r.stdout, r.stderr)
	}
}

// TestLandCommandRefusalKeepsTheCensusRecords proves a landing that refuses before
// its gate prints no landed record and drops nothing, so the operator can repair the
// invocation and land with the evidence intact.
func TestLandCommandRefusalKeepsTheCensusRecords(t *testing.T) {
	t.Parallel()
	request := "census-landed-refusal"
	f := publicLandingFixture(t, request, "", "")
	recordRawCalls(t, f.home, f.root, f.creation.Path, 2)
	r := runVerb(t, verbLand, f.call(landArgs("no-such-request", f.base, f.tip, f.creation.Path)...))
	if r.exit == 0 || !strings.Contains(r.stdout, "refused{") || strings.Contains(r.stdout, "landed{") {
		t.Fatalf("refused landing = (%d, %q, %q), want a refusal and no landed record", r.exit, r.stdout, r.stderr)
	}
	if _, err := os.Stat(f.tally); !os.IsNotExist(err) {
		t.Fatalf("the refusal ran the gate: %v", err)
	}
	if _, err := os.Stat(censusRecordPath(f.home, f.root, f.creation.Assignment.ID)); err != nil {
		t.Fatalf("the refusal dropped the census records: %v", err)
	}
}
