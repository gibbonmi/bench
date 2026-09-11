package worktree

import (
	"bytes"
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

// TestLandCommandStatesTheCensusCountAndDropsTheRecords is EC20 and the landing half
// of EC24. The landed record carries the count as its last key, and the release step
// the landing runs leaves no record file for the retired assignment.
func TestLandCommandStatesTheCensusCountAndDropsTheRecords(t *testing.T) {
	t.Parallel()
	request := "census-landed-count"
	root, creation, base, tip, _, home := publicLandingFixture(t, request, "", "")
	recordRawCalls(t, home, root, creation.Path, 3)
	survivor := seedHandoffSections(t, root, creation.Assignment)
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 0 || !strings.HasSuffix(stdout.String(), ",census=3}\n") {
		t.Fatalf("landed record = (%d, %q, %q), want census=3 as the last key", code, stdout.String(), stderr.String())
	}
	if _, err := os.Stat(censusRecordPath(home, root, creation.Assignment.ID)); !os.IsNotExist(err) {
		t.Fatalf("the released landing kept the census record: %v", err)
	}
	requireHandoffSections(t, root, handoffdoc.MainKey, survivor)
}

// TestLandCommandStatesZeroForAnAssignmentWithNoRecords is EC21. Zero is a stated
// fact, and an absent record file is not a landing failure.
func TestLandCommandStatesZeroForAnAssignmentWithNoRecords(t *testing.T) {
	t.Parallel()
	request := "census-landed-zero"
	root, creation, base, tip, _, home := publicLandingFixture(t, request, "", "")
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 0 || !strings.HasSuffix(stdout.String(), ",census=0}\n") {
		t.Fatalf("landed record = (%d, %q, %q), want census=0", code, stdout.String(), stderr.String())
	}
}

// TestLandCommandPrintsTheCensusHeadBreakdown proves the landing states the raw-call
// count for each verb head before the release step drops the records, so the retro
// reads the breakdown from the run. The heaviest head prints first.
func TestLandCommandPrintsTheCensusHeadBreakdown(t *testing.T) {
	t.Parallel()
	request := "census-landed-heads"
	root, creation, base, tip, _, home := publicLandingFixture(t, request, "", "")
	recordRawCallsWithHead(t, home, root, creation.Path, "sed -i s/a/b/", 2)
	recordRawCallsWithHead(t, home, root, creation.Path, "awk -f x", 1)
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 0 || !strings.Contains(stderr.String(), "census heads{sed=2,awk=1}\n") {
		t.Fatalf("landing evidence = (%d, %q, %q), want the head breakdown on stderr", code, stdout.String(), stderr.String())
	}
	if !strings.HasSuffix(stdout.String(), ",census=3}\n") {
		t.Fatalf("landed record = %q, want census=3 beside the breakdown", stdout.String())
	}
}

// TestLandCommandPrintsNoHeadsLineWithoutRecords proves an assignment that made no raw
// call prints no breakdown at all, and still states the zero count in its record.
func TestLandCommandPrintsNoHeadsLineWithoutRecords(t *testing.T) {
	t.Parallel()
	request := "census-landed-no-heads"
	root, creation, base, tip, _, home := publicLandingFixture(t, request, "", "")
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 0 || strings.Contains(stderr.String(), "census heads{") || !strings.HasSuffix(stdout.String(), ",census=0}\n") {
		t.Fatalf("empty census landing = (%d, %q, %q), want no heads line and census=0", code, stdout.String(), stderr.String())
	}
}

// TestLandCommandRefusalKeepsTheCensusRecords proves a landing that refuses before
// its gate prints no landed record and drops nothing, so the operator can repair the
// invocation and land with the evidence intact.
func TestLandCommandRefusalKeepsTheCensusRecords(t *testing.T) {
	t.Parallel()
	request := "census-landed-refusal"
	root, creation, base, tip, tally, home := publicLandingFixture(t, request, "", "")
	recordRawCalls(t, home, root, creation.Path, 2)
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs("no-such-request", base, tip, creation.Path), &stdout, &stderr)
	if code == 0 || !strings.Contains(stdout.String(), "refused{") || strings.Contains(stdout.String(), "landed{") {
		t.Fatalf("refused landing = (%d, %q, %q), want a refusal and no landed record", code, stdout.String(), stderr.String())
	}
	if _, err := os.Stat(tally); !os.IsNotExist(err) {
		t.Fatalf("the refusal ran the gate: %v", err)
	}
	if _, err := os.Stat(censusRecordPath(home, root, creation.Assignment.ID)); err != nil {
		t.Fatalf("the refusal dropped the census records: %v", err)
	}
}
