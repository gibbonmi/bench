package worktree

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

// requireLandedReaders runs the three landed readers over one repository and holds each to
// the same verdict. The list cell, the selector of clean --landed, and the session-start
// count must treat every creation in landed as landed and every creation in unlanded as
// unlanded work. The summary is the count line that the session start must print.
func requireLandedReaders(t *testing.T, fixture repoHome, landed, unlanded []Creation, summary string) {
	t.Helper()
	list := runVerb(t, verbList, fixture.call())
	requireTest(t, list.exit == 0, "ListCommand exit=%d output=%q", list.exit, list.stdout)
	cells := map[string]any{}
	for _, row := range list.mustRows(t, listTable) {
		fields := row.(map[string]any)
		cells[fields["id"].(string)] = fields["landed"]
	}
	for _, creation := range landed {
		requireTest(t, cells[creation.Assignment.ID] == true, "list landed cell of %s = %v, want true", creation.Assignment.ID, cells[creation.Assignment.ID])
	}
	for _, creation := range unlanded {
		requireTest(t, cells[creation.Assignment.ID] == false, "list landed cell of %s = %v, want false", creation.Assignment.ID, cells[creation.Assignment.ID])
	}

	plan := runVerb(t, verbClean, fixture.call("--landed"))
	requireTest(t, plan.exit == 0, "clean --landed exit=%d stdout=%q stderr=%q", plan.exit, plan.stdout, plan.stderr)
	for _, creation := range landed {
		requireTest(t, strings.Contains(plan.stdout, creation.Path), "clean --landed output=%q, want %s", plan.stdout, creation.Path)
	}
	for _, creation := range unlanded {
		requireTest(t, !strings.Contains(plan.stdout, creation.Path), "clean --landed selected the unlanded worktree: %q", plan.stdout)
	}

	resumed := runVerb(t, verbResumeClean, fixture.call())
	requireTest(t, resumed.exit == 0, "ResumeCleanCommand exit=%d stdout=%q stderr=%q", resumed.exit, resumed.stdout, resumed.stderr)
	requireTest(t, strings.Contains(resumed.stdout, summary), "summary=%q, want %q", resumed.stdout, summary)
}

// TestDirtyAssignmentWithoutOwnCommitIsNotLanded pins FT408. A branch tip that sits at the
// recorded start is a main commit, but an assignment with no commit of its own and
// uncommitted output holds unlanded work. Each landed reader must agree on that verdict,
// and the two neighbors that keep a landed verdict must stay landed.
func TestDirtyAssignmentWithoutOwnCommitIsNotLanded(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	dirtyNoCommit := mustCreate(t, root, home, "no-commit-dirty", "no commit dirty")
	cleanNoCommit := mustCreate(t, root, home, "no-commit-clean", "no commit clean")
	dirtyOwnCommit := mustCreate(t, root, home, "own-commit-dirty", "own commit dirty")
	landAssignment(t, root, dirtyOwnCommit, "own-commit.txt")
	mustWrite(t, filepath.Join(dirtyNoCommit.Path, "research.txt"), []byte("output\n"), 0o644)
	mustWrite(t, filepath.Join(dirtyOwnCommit.Path, "extra.txt"), []byte("output\n"), 0o644)

	requireLandedReaders(t, repoHome{root, home}, []Creation{cleanNoCommit, dirtyOwnCommit}, []Creation{dirtyNoCommit}, "retained active=1 landed=2")
}

// TestFastForwardOverDefaultCommitsGivesNoOwnCommit pins the own-commit rule of FT408. A
// fast-forward merge of main moves the branch tip past the recorded start, but the commits
// it takes in belong to main. One of them is a landing merge, whose second-parent commit
// sits in the range from the start to the tip and yet is no commit of this assignment. The
// dirty ones therefore stay unlanded work for each landed reader, and the clean one stays
// landed.
func TestFastForwardOverDefaultCommitsGivesNoOwnCommit(t *testing.T) {
	t.Parallel()
	f := mergeFixture(t, "plain", "landing", "clean")
	plain, landing, clean := f.created[0], f.created[1], f.created[2]
	commitOnDefault(t, f.root, "plain.txt", "plain\n")
	gitRun(t, f.root, "checkout", "-q", "-b", "side")
	commitOnDefault(t, f.root, "side.txt", "side\n")
	gitRun(t, f.root, "checkout", "-q", "-")
	gitRun(t, f.root, "merge", "--no-ff", "-q", "-m", "land side", "side")
	for _, creation := range f.created {
		merged := runVerb(t, verbMerge, f.merge("--from", "main", creation.Assignment.ID))
		requireTest(t, merged.exit == 0, "merge exit=%d stdout=%q stderr=%q", merged.exit, merged.stdout, merged.stderr)
		requireTest(t, strings.Contains(mergedRecord(t, merged.stdout), "kind=fast-forward"), "merge stdout=%q, want a fast-forward", merged.stdout)
	}
	mustWrite(t, filepath.Join(plain.Path, "research.txt"), []byte("output\n"), 0o644)
	mustWrite(t, filepath.Join(landing.Path, "research.txt"), []byte("output\n"), 0o644)

	requireLandedReaders(t, f.repoHome, []Creation{clean}, []Creation{plain, landing}, "retained active=2 landed=1")
}

// TestOrphanSweepReportsDirtyAssignmentWithoutOwnCommit pins the sweep's use of the landed
// verdict. An orphaned record whose tree holds output and no own commit is no landed
// record, so the sweep reports it. An orphaned clean one stays out of the report.
func TestOrphanSweepReportsDirtyAssignmentWithoutOwnCommit(t *testing.T) {
	t.Parallel()
	root, home := newSweepRepo(t)
	dirty := mustCreate(t, root, home, "sweep-dirty-no-commit", "aged, dirty, no own commit")
	clean := mustCreate(t, root, home, "sweep-clean-no-commit", "aged, clean, no own commit")
	mustWrite(t, filepath.Join(dirty.Path, "research.txt"), []byte("output\n"), 0o644)
	backdate(t, root, dirty.Assignment, 8*24*time.Hour)
	backdate(t, root, clean.Assignment, 8*24*time.Hour)

	var reported []string
	for _, orphan := range mustSweep(t, root, home).Orphans {
		reported = append(reported, orphan.ID)
	}
	requireTest(t, slices.Equal(reported, []string{dirty.Assignment.ID}), "sweep reported %v, want only the dirty assignment %s", reported, dirty.Assignment.ID)
}
