package worktree

// The resume-clean verb reads the home and the instant of the ambient value that its entry
// builds. Each test runs the verb through the verb runner and binds no process environment.

import (
	"os"
	"strings"
	"testing"
	"time"
)

// TestResumeCleanDropsTheCensusRecordUnderTheCallerHome retires a cleanup-pending
// assignment under a fixture home that is not the process home. A retirement that drops
// the record under the process home leaves this record behind.
func TestResumeCleanDropsTheCensusRecordUnderTheCallerHome(t *testing.T) {
	t.Parallel()
	f := newPendingAssignment(t, "resume-census-home")
	recordRawCalls(t, f.home, f.root, f.creation.Path, 2)
	record := censusRecordPath(f.home, f.root, f.creation.Assignment.ID)
	if _, err := os.Stat(record); err != nil {
		t.Fatalf("the fixture wrote no census record under the caller home: %v", err)
	}
	r := runVerb(t, verbResumeClean, f.call())
	if r.exit != 0 {
		t.Fatalf("resume-clean = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	if _, err := os.Stat(f.creation.Path); !os.IsNotExist(err) {
		t.Fatalf("resume-clean kept the cleanup-pending checkout: %v", err)
	}
	if _, err := os.Stat(record); !os.IsNotExist(err) {
		t.Fatalf("resume-clean kept the census record under the caller home: %v", err)
	}
}

// TestResumeCleanJudgesTheClockValue runs resume-clean eight days after an assignment's
// creation. The assignment is new by the real clock, so only a run that judges the call's
// clock value names the clean command for its path.
func TestResumeCleanJudgesTheClockValue(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "resume-clock-value")
	makeUnlandedAssignment(t, f.creation)
	created, err := time.Parse(time.RFC3339, *f.creation.Assignment.CreatedAt)
	mustNoError(t, err)
	call := f.call()
	call.clock = func() time.Time { return created.Add(8 * 24 * time.Hour) }
	r := runVerb(t, verbResumeClean, call)
	if r.exit != 0 || !strings.Contains(r.stdout, cleanLineFor(f.creation.Path)) {
		t.Fatalf("resume-clean = (%d, %q, %q), want the clean command for %s", r.exit, r.stdout, r.stderr, f.creation.Path)
	}
}
