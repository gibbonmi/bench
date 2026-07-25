package status

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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
	if err := os.WriteFile(filepath.Join(root, "session-handoff.md"), []byte("# Session handoff\n\n"+body+"\n"), 0o644); err != nil {
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

// The failure this signal exists for, and the one this repo actually hit: commits landed
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
	// that requires an invocation, which left this signal invisible to the handoff's own
	// Next-command field.
	if rows[0].action != "bench handoff" {
		t.Errorf("action = %q, want the command that rewrites it", rows[0].action)
	}
}

// The distance is measured from the handoff's own last write, not from whatever the tree
// most recently committed: a later commit elsewhere must not reset the handoff's age.
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
// is about to be reset, so a row here would fire on the very act that fixes it.
func TestAppendHandoffInFlightEditIsSilent(t *testing.T) {
	root := initRepo(t)
	commitFile(t, root, "base.txt")
	commitHandoff(t, root, "first")
	commitFile(t, root, "one.txt")
	if err := os.WriteFile(filepath.Join(root, "session-handoff.md"), []byte("# Session handoff\n\nrewritten\n"), 0o644); err != nil {
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
	if err := os.WriteFile(filepath.Join(root, "session-handoff.md"), []byte("# Session handoff\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if rows := appendHandoff(nil, root); len(rows) != 0 {
		t.Fatalf("untracked handoff produced rows: %#v", rows)
	}
}

// The row has to be visible where it matters: on the otherwise-clean tree a cold session
// picks up, a stale handoff leads the board rather than being budgeted off it.
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
