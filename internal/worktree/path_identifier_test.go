package worktree

import (
	"bytes"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var listPathAction = regexp.MustCompile(`bench worktree path (\S+),inspect active worktree`)

// Every help row `bench worktree list` emits is an exact executable invocation, so the
// identifier it prints has to be the one the resolver accepts. The row is authored from
// the id column rather than the label, because ids are unique and labels can collide.
// A follow-up that fails is worse than no follow-up at all.
// It is the only address an agent has for a worktree it did not create.
func TestListPathActionRunsAsAdvertised(t *testing.T) {
	root, creation, home := newOwnedAssignment(t, "advertised")
	chdir(t, root)
	listed, code := ListCommand(root, home, nil)
	if code != 0 {
		t.Fatalf("list code=%d out=%q", code, listed)
	}
	match := listPathAction.FindStringSubmatch(listed)
	if match == nil {
		t.Fatalf("list emitted no path action:\n%s", listed)
	}
	if match[1] != creation.Assignment.ID {
		t.Fatalf("path action addresses %q, want the id column %q", match[1], creation.Assignment.ID)
	}
	var stdout, stderr bytes.Buffer
	if code := PathCommand(root, home, []string{match[1]}, &stdout, &stderr); code != 0 {
		t.Fatalf("advertised %q exited %d: %s", match[1], code, stderr.String())
	}
}

// [PB34] The path serves the file tools, and a delegate that pastes it into a shell step
// leaves the worktree boundary. The note says so on stderr, because stdout stays the one
// path line a `$(...)` capture reads. The note quotes the target as typed, so the
// suggested command is the one the caller can run.
func TestPathNotesTheFileToolRouteOnStderr(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "filetools")
	target := creation.Assignment.Label
	var stdout, stderr bytes.Buffer
	if code := PathCommand(root, home, []string{target}, &stdout, &stderr); code != 0 {
		t.Fatalf("path exited %d: %s", code, stderr.String())
	}
	printed := strings.TrimSuffix(stdout.String(), "\n")
	if strings.Contains(printed, "\n") || !filepath.IsAbs(printed) {
		t.Fatalf("stdout = %q, want one absolute path line alone", stdout.String())
	}
	want := "note: the path serves the file tools; run a shell step through bench worktree exec " + target + " -- <command>\n"
	if stderr.String() != want {
		t.Fatalf("stderr = %q, want %q", stderr.String(), want)
	}
}

// The label stays a valid address; accepting the id widens the grammar, it does not
// replace it.
func TestPathResolvesTheLabelAndTheIdAlike(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "both")
	for _, target := range []string{creation.Assignment.ID, creation.Assignment.Label} {
		var byTarget, stderr bytes.Buffer
		if code := PathCommand(root, home, []string{target}, &byTarget, &stderr); code != 0 {
			t.Fatalf("target %q exited %d: %s", target, code, stderr.String())
		}
		if byTarget.Len() == 0 {
			t.Fatalf("target %q resolved to no path", target)
		}
	}
}
