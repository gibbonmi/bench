package worktree

// The resume-clean verb reads the home and the instant of the ambient value that its entry
// builds. Each test runs the verb through the verb runner and binds no process environment.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestResumeCleanDropsTheCensusRecordUnderTheCallerHome retires a cleanup-pending
// assignment under a fixture home that is not the process home. A retirement that drops
// the record under the process home leaves this record behind. The home is a sibling
// directory of the repository, and its name holds a space, so the retired checkout path
// holds a space.
func TestResumeCleanDropsTheCensusRecordUnderTheCallerHome(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	f := ownedAssignment{repoHome: repoHome{root, filepath.Join(filepath.Dir(root), "auto home")}}
	f.creation = mustCreate(t, f.root, f.home, "resume-census-home", "census home")
	markPending(t, f.root, f.creation.Assignment)
	if !strings.Contains(f.creation.Path, " ") {
		t.Fatalf("the fixture checkout path %q holds no space", f.creation.Path)
	}
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
// clock value names the clean command for its path. The per-checkout plan and the orphan
// sweep each judge the age, so the test reads the retained stale-active count and the clean
// command.
func TestResumeCleanJudgesTheClockValue(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "resume-clock-value")
	makeUnlandedAssignment(t, f.creation)
	created, err := time.Parse(time.RFC3339, *f.creation.Assignment.CreatedAt)
	mustNoError(t, err)
	call := f.call()
	call.clock = func() time.Time { return created.Add(8 * 24 * time.Hour) }
	r := runVerb(t, verbResumeClean, call)
	if r.exit != 0 || !strings.Contains(r.stdout, "retained stale-active=1;") || !strings.Contains(r.stdout, cleanLineFor(f.creation.Path)) {
		t.Fatalf("resume-clean = (%d, %q, %q), want the stale-active count and the clean command for %s", r.exit, r.stdout, r.stderr, f.creation.Path)
	}
}
