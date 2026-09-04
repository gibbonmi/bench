package status

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/handoffdoc"
	"github.com/gibbonmi/bench/internal/intent"
)

// commitFile makes one commit touching a uniquely-named file and returns its full sha.
func commitFile(t *testing.T, root, name string) string {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, name), []byte(name+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", name)
	return gitRun(t, root, "rev-parse", "HEAD")
}

// commitHandoff writes and commits the handoff with distinct body text, returning the
// commit that wrote it. The body varies so a rewrite is a real commit, not an empty one.
func commitHandoff(t *testing.T, root, body string) string {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, "capture/session-handoff.md"), []byte("# Session handoff\n\n"+body+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "handoff")
	return gitRun(t, root, "rev-parse", "HEAD")
}

// A repo with no handoff produces no row: not every repo keeps one, and absence is a
// choice rather than a defect.
func TestAppendHandoffAbsentIsSilent(t *testing.T) {
	root := initRepo(t)
	commitFile(t, root, "base.txt")

	if rows := appendHandoff(nil, root); len(rows) != 0 {
		t.Fatalf("absent handoff produced rows: %#v", rows)
	}
}

// The quiet case that matters: a cold session picking up a tree whose last commit wrote
// the handoff. Nothing has landed since, so there is nothing to report.
func TestAppendHandoffWrittenAtHeadIsSilent(t *testing.T) {
	root := initRepo(t)
	commitFile(t, root, "base.txt")
	commitHandoff(t, root, "first")

	if rows := appendHandoff(nil, root); len(rows) != 0 {
		t.Fatalf("handoff written at HEAD produced rows: %#v", rows)
	}
}

// The failure this signal exists for is the one this repo actually hit. Commits landed
// after the handoff was written and nobody rewrote it, so a cold session would trust a
// document the tree has moved past.
func TestAppendHandoffBehindHeadReportsDistance(t *testing.T) {
	root := initRepo(t)
	commitFile(t, root, "base.txt")
	written := commitHandoff(t, root, "first")
	commitFile(t, root, "one.txt")
	commitFile(t, root, "two.txt")

	rows := appendHandoff(nil, root)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one handoff row", rows)
	}
	if rows[0].signal != "handoff" {
		t.Errorf("signal = %q, want handoff", rows[0].signal)
	}
	if !strings.Contains(rows[0].detail, "2 commits behind") {
		t.Errorf("detail = %q, want it to name the 2-commit distance", rows[0].detail)
	}
	if !strings.Contains(rows[0].detail, Short(written)) {
		t.Errorf("detail = %q, want it to name the commit that wrote the handoff", rows[0].detail)
	}
	// The action is the command that rewrites the handoff, not prose describing the job.
	// A row naming the work in words can never be selected as a next command by a reader
	// that requires an invocation. That left this signal invisible to the handoff's own
	// Next-command field.
	if rows[0].action.render() != "bench handoff" {
		t.Errorf("action = %q, want the command that rewrites it", rows[0].action.render())
	}
}

// The distance is measured from the handoff's own last write, not from whatever the tree
// most recently committed. A later commit elsewhere must not reset the handoff's age.
func TestAppendHandoffDistanceIgnoresUnrelatedCommits(t *testing.T) {
	root := initRepo(t)
	commitFile(t, root, "base.txt")
	commitHandoff(t, root, "first")
	commitFile(t, root, "one.txt")

	// Rewriting the handoff resets the age; the unrelated commit above must not.
	rows := appendHandoff(nil, root)
	if len(rows) != 1 || !strings.Contains(rows[0].detail, "1 commit behind") {
		t.Fatalf("rows = %#v, want a 1-commit distance from the handoff's own write", rows)
	}
	commitHandoff(t, root, "rewritten at the later commit")
	if rows := appendHandoff(nil, root); len(rows) != 0 {
		t.Fatalf("rewriting the handoff did not reset its age: %#v", rows)
	}
}

// A handoff being written right now is not stale: the session is mid-rewrite and its age
// is about to be reset. A row here would fire on the very act that fixes it.
func TestAppendHandoffInFlightEditIsSilent(t *testing.T) {
	root := initRepo(t)
	commitFile(t, root, "base.txt")
	commitHandoff(t, root, "first")
	commitFile(t, root, "one.txt")
	if err := os.WriteFile(filepath.Join(root, "capture/session-handoff.md"), []byte("# Session handoff\n\nrewritten\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if rows := appendHandoff(nil, root); len(rows) != 0 {
		t.Fatalf("in-flight handoff edit produced rows: %#v", rows)
	}
}

// An untracked handoff has never been handed off to anyone, so it has no age to report.
func TestAppendHandoffUntrackedIsSilent(t *testing.T) {
	root := initRepo(t)
	commitFile(t, root, "base.txt")
	if err := os.WriteFile(filepath.Join(root, "capture/session-handoff.md"), []byte("# Session handoff\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if rows := appendHandoff(nil, root); len(rows) != 0 {
		t.Fatalf("untracked handoff produced rows: %#v", rows)
	}
}

// commitFileAt makes one commit whose committer date is pinned, so a test can place it
// before or after a file's write time deterministically. The committer date is the one
// `rev-list --since` filters on, and only the environment can set it.
func commitFileAt(t *testing.T, root, name string, at time.Time) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, name), []byte(name+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	stamp := at.Format(time.RFC3339)
	t.Setenv("GIT_AUTHOR_DATE", stamp)
	t.Setenv("GIT_COMMITTER_DATE", stamp)
	defer os.Unsetenv("GIT_AUTHOR_DATE")
	defer os.Unsetenv("GIT_COMMITTER_DATE")
	gitRun(t, root, "commit", "-m", name)
}

// An ignored handoff is a local file with no commit of its own, so its age comes from
// the file's write time. Commits after that write still report a distance, and touching
// the file resets it.
func TestAppendHandoffIgnoredUsesFileTime(t *testing.T) {
	root := initRepo(t)
	commitFileAt(t, root, "base.txt", time.Now().Add(-2*time.Hour))
	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte("capture/session-handoff.md\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "capture"), 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "capture/session-handoff.md")
	if err := os.WriteFile(path, []byte("# Session handoff\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// The write sits one hour in the past, between the base commit (two hours back)
	// and the fresh commit below, so exactly one commit follows it.
	past := time.Now().Add(-time.Hour)
	if err := os.Chtimes(path, past, past); err != nil {
		t.Fatal(err)
	}
	commitFile(t, root, "one.txt")

	rows := appendHandoff(nil, root)
	if len(rows) != 1 || rows[0].signal != "handoff" {
		t.Fatalf("rows = %#v, want one handoff row", rows)
	}
	if !strings.Contains(rows[0].detail, "1 commit behind") {
		t.Errorf("detail = %q, want the 1-commit distance", rows[0].detail)
	}

	// A fresh write resets the age: nothing has landed since.
	now := time.Now()
	if err := os.Chtimes(path, now, now); err != nil {
		t.Fatal(err)
	}
	if rows := appendHandoff(nil, root); len(rows) != 0 {
		t.Fatalf("fresh ignored handoff produced rows: %#v", rows)
	}
}

// An ignored repo that keeps no handoff at all reports nothing: absence stays a choice.
func TestAppendHandoffIgnoredAbsentIsSilent(t *testing.T) {
	root := initRepo(t)
	commitFile(t, root, "base.txt")
	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte("capture/session-handoff.md\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if rows := appendHandoff(nil, root); len(rows) != 0 {
		t.Fatalf("absent ignored handoff produced rows: %#v", rows)
	}
}

// ignoreHandoff makes the handoff a local file, which is the state this repository keeps
// it in and the one the file-age clock covers.
func ignoreHandoff(t *testing.T, root string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte("capture/session-handoff.md\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

// seedBranchAssignment writes one active assignment and points its branch ref at tip. It
// returns the request digest, which is the key of the section that assignment owns.
func seedBranchAssignment(t *testing.T, root, id, tip string) string {
	t.Helper()
	seedAssignment(t, root, id, id, intent.StateActive)
	gitRun(t, root, "update-ref", intent.AssignmentBranchRef(censusOwnerID, id), tip)
	return intent.RequestDigest("request-" + id)
}

// seedHandoffSection writes one request section that records tip as its worktree tip. The
// State body varies so a rewrite changes the document rather than replacing it byte for
// byte.
func seedHandoffSection(t *testing.T, root, key, tip, state string) {
	t.Helper()
	section := handoffdoc.Section{
		Key:    key,
		Fields: []handoffdoc.Field{{Label: handoffdoc.LabelWorktreeTip, Value: tip}},
		State:  state,
	}
	if err := handoffdoc.WriteSection(filepath.Join(root, "capture/session-handoff.md"), section); err != nil {
		t.Fatal(err)
	}
}

// twoSectionRepo builds the fixture both section tests read: a document with one section
// three commits behind its assignment branch and one section at that branch's tip. It
// returns the root and the two request digests.
func twoSectionRepo(t *testing.T) (root, behind, fresh string) {
	t.Helper()
	root = initRepo(t)
	commitFile(t, root, "base.txt")
	ignoreHandoff(t, root)
	base := gitRun(t, root, "rev-parse", "HEAD")
	commitFile(t, root, "one.txt")
	commitFile(t, root, "two.txt")
	head := commitFile(t, root, "three.txt")

	behind = seedBranchAssignment(t, root, strings.Repeat("d", 32), head)
	fresh = seedBranchAssignment(t, root, strings.Repeat("e", 32), head)
	seedHandoffSection(t, root, behind, base, "behind")
	seedHandoffSection(t, root, fresh, head, "fresh")
	return root, behind, fresh
}

// behindRow returns the one row that states a section's distance, and fails when the rows
// hold none or more than one.
func behindRow(t *testing.T, rows []row) row {
	t.Helper()
	var found []row
	for _, r := range rows {
		if strings.Contains(r.detail, "behind") {
			found = append(found, r)
		}
	}
	if len(found) != 1 {
		t.Fatalf("rows = %#v, want exactly one row stating a distance", rows)
	}
	return found[0]
}

// A section is dated against its own assignment branch, not against the document's write
// time. A file-level clock reads both sections current, so the stale sibling hides.
// (Coverage row HS21.)
func TestAppendHandoffSectionNamesTheBehindSection(t *testing.T) {
	root, behind, fresh := twoSectionRepo(t)

	got := behindRow(t, appendHandoff(nil, root))
	if got.signal != "handoff" {
		t.Errorf("signal = %q, want handoff", got.signal)
	}
	if !strings.Contains(got.detail, "3 commits behind") {
		t.Errorf("detail = %q, want the 3-commit distance", got.detail)
	}
	if !strings.Contains(got.detail, Short(behind)) {
		t.Errorf("detail = %q, want it to name the behind section", got.detail)
	}
	if strings.Contains(got.detail, Short(fresh)) {
		t.Errorf("detail = %q, want it to leave the fresh section unnamed", got.detail)
	}
}

// Rewriting the fresh section resets the document's write time. The behind section's own
// distance is unmoved by that, so the row it produces must not change. (Coverage row HS22.)
func TestAppendHandoffSectionRewriteLeavesTheDistance(t *testing.T) {
	root, _, fresh := twoSectionRepo(t)
	before := behindRow(t, appendHandoff(nil, root))

	head := gitRun(t, root, "rev-parse", "HEAD")
	seedHandoffSection(t, root, fresh, head, "rewritten")

	if after := behindRow(t, appendHandoff(nil, root)); after.detail != before.detail {
		t.Errorf("detail = %q after the fresh rewrite, want %q", after.detail, before.detail)
	}
}

// The row has to be visible where it matters: on the otherwise-clean tree a cold session
// picks up. A stale handoff leads the board rather than being budgeted off it.
func TestRenderStaleHandoffLeadsCleanBoard(t *testing.T) {
	root := initRepo(t)
	commitFile(t, root, "base.txt")
	commitHandoff(t, root, "second")
	commitFile(t, root, "one.txt")

	out := render(root, false)
	if !strings.HasPrefix(out, "▶ bench handoff  (handoff)") {
		t.Fatalf("board did not lead with the handoff signal:\n%s", out)
	}
}
