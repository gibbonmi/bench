package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/responsebound/responseboundtest"
	"github.com/gibbonmi/bench/internal/usage"
)

// putSelectedAssignment records one complete assignment whose tree is absent, so the
// selected view answers a non-active record that the active-path resolver refuses.
func putSelectedAssignment(t *testing.T, root, id, owner string) intent.Assignment {
	t.Helper()
	assignment := intent.Assignment{
		Schema: intent.AssignmentRecordSchema, ID: id, OwnerID: owner,
		Request: intent.RequestDigest("selected-route-" + id), RequestToken: "selected-route-" + id,
		Label: "selected route " + id[:1], Start: strings.TrimSpace(runAXIGit(t, "-C", root, "rev-parse", "HEAD")),
		Branch: intent.AssignmentBranchRef(owner, id), Worktree: filepath.Join(root, "missing tree "+id[:1]), State: intent.StateComplete,
	}
	if err := intent.PutAssignment(root, assignment); err != nil {
		t.Fatal(err)
	}
	return assignment
}

func TestSelectedWorktreeCommandRoutes(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	assignment := putSelectedAssignment(t, root, strings.Repeat("a", 32), strings.Repeat("b", 32))
	for _, executable := range []string{"bench", filepath.Join(root, ".bench", "bin", "bench")} {
		for _, target := range []string{"bad\x1btarget", "space target", "--target", "$(touch sentinel)"} {
			result := runAXICommandAsAt(t, filepath.Join(root, "nested", "deep"), executable, []string{"worktree", "list", "--view", "paths", "--target", target, "--target", assignment.ID})
			if result.code != 1 || result.stderr != "" {
				t.Fatalf("selected route=%#v", result)
			}
			document, err := axitest.DecodeDocument(result.stdout)
			if err != nil {
				t.Fatal(err)
			}
			rows, err := document.Rows("worktrees")
			if err != nil {
				t.Fatal(err)
			}
			if len(rows) != 2 {
				t.Fatalf("selected rows=%#v", rows)
			}
			row := rows[1].(map[string]any)
			if row["id"] != assignment.ID || row["path"] != assignment.Worktree || row["state"] != "complete" || row["error"] != "" {
				t.Fatalf("non-active route=%#v", row)
			}
		}
	}
}

func TestSelectedWorktreeHelpDiscovery(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	result := runAXICommandAt(t, root, []string{"worktree", "--help"})
	if result.code != 0 || !strings.Contains(result.stdout, usage.WorktreeListPaths) {
		t.Fatalf("worktree help=%#v", result)
	}
	result = runAXICommandAt(t, root, []string{"worktree", "list", "--view", "paths", "--help"})
	if result.code != 0 || result.stdout != "usage: "+usage.WorktreeListPaths+"\n" {
		t.Fatalf("selected help=%#v", result)
	}
}

// TestSelectedWorktreeWithinResponseBound reads the printed response, not the spill file.
// Three resolved targets give 1 table header, 3 rows, 1 help header, and 1 help row: 6 lines.
func TestSelectedWorktreeWithinResponseBound(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	owner := strings.Repeat("b", 32)
	argv := []string{"worktree", "list", "--view", "paths"}
	for _, id := range []string{strings.Repeat("a", 32), strings.Repeat("c", 32), strings.Repeat("d", 32)} {
		argv = append(argv, "--target", putSelectedAssignment(t, root, id, owner).ID)
	}
	putSelectedAssignment(t, root, strings.Repeat("e", 32), owner)
	t.Chdir(root)
	var stdout, stderr bytes.Buffer
	code := Command{Stdout: &stdout, Stderr: &stderr, Executable: "bench"}.Run(argv)
	printed := stdout.String()
	if code != 0 || stderr.Len() != 0 {
		t.Fatalf("selected response exit=%d stdout=%q stderr=%q", code, printed, stderr.String())
	}
	if _, spilled := responseboundtest.Find(printed); spilled {
		t.Fatalf("selected response spilled: %q", printed)
	}
	if lines := strings.Count(printed, "\n"); lines != 6 || !strings.HasSuffix(printed, "\n") {
		t.Fatalf("selected response has %d lines, want exactly 6: %q", lines, printed)
	}
}
