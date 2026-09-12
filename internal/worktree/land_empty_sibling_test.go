package worktree

import (
	"bytes"
	"errors"
	"io"
	"path/filepath"
	"strings"
	"testing"
)

func TestLandRetainsSiblingBornAtPublishedTip(t *testing.T) {
	t.Parallel()
	for _, resume := range []bool{false, true} {
		name := "ordinary"
		if resume {
			name = "resumed"
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := "land-newborn-" + name
			root, source, base, tip, _, home := publicLandingFixture(t, request, "", "")
			j, _ := refreshJoins(nil)
			release := j.releaseLandingAssignment
			var sibling Creation
			j.releaseLandingAssignment = func(inner joins, root, home string, args []string, stdout, stderr io.Writer) int {
				sibling = mustCreate(t, root, home, request+"-sibling", "new sibling")
				if resume {
					return 1
				}
				return release(inner, root, home, args, stdout, stderr)
			}
			var stdout, stderr bytes.Buffer
			code := landWith(j, root, home, "", landArgs(request, base, tip, source.Path), &stdout, &stderr)
			if resume {
				if code != 3 || !strings.Contains(stdout.String(), "worktree=incomplete:release") {
					t.Fatalf("interrupted landing = (%d, %q, %q)", code, stdout.String(), stderr.String())
				}
				published := gitOutput(t, root, "rev-parse", "main")
				j.releaseLandingAssignment = release
				stdout.Reset()
				stderr.Reset()
				code = resumeLandWith(j, root, home, resumeLandArgs(published, request, base, tip, source.Path), &stdout, &stderr)
				if got := gitOutput(t, root, "rev-parse", "main"); got != published {
					t.Fatalf("resume republished: %s", got)
				}
			}
			if code != 0 {
				t.Fatalf("landing = (%d, %q, %q)", code, stdout.String(), stderr.String())
			}
			requirePresent(t, sibling.Path, "new sibling worktree")
			if !assignmentActive(t, root, sibling.Assignment.ID) {
				t.Fatal("new sibling assignment was retired")
			}
			if got := gitOutput(t, root, "rev-parse", sibling.Assignment.Branch); got != sibling.Assignment.Start {
				t.Fatalf("new sibling tip = %q, want its start %q", got, sibling.Assignment.Start)
			}
			row, selected := selectLandedCleanupRow(j, root, sibling.Assignment, "main", "none", CleanupOptions{}, "")
			if !selected || !row.plan.Action.Removes() {
				t.Fatalf("explicit cleanup changed: selected=%t, plan=%#v", selected, row.plan)
			}
		})
	}
}

func TestScopedCleanupRequiresProvenContribution(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	sibling := mustCreate(t, root, filepath.Join(root, ".bench-home"), "contribution-proof", "sibling")
	scope := sibling.Assignment.Start
	landAssignment(t, root, sibling, "contribution.txt")
	head := gitOutput(t, root, "rev-parse", sibling.Assignment.Branch)
	commitInWorktree(t, root, "later.txt", "later\n", "later main commit")
	later := gitOutput(t, root, "rev-parse", "main")
	j := defaultJoins()
	planned, selected := selectLandedCleanupRow(j, root, sibling.Assignment, "main", "none", CleanupOptions{}, scope)
	if !selected || !planned.plan.Action.Removes() {
		t.Fatalf("contributed sibling excluded: selected=%t, plan=%#v", selected, planned.plan)
	}
	for _, tc := range []struct{ name, start string }{
		{"no contribution", head},
		{"unresolved start", strings.Repeat("0", len(head))},
		{"inconsistent start", later},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assignment := sibling.Assignment
			assignment.Start = tc.start
			writeAssignmentLedger(t, root, assignment)
			if _, selected := selectLandedCleanupRow(j, root, assignment, "main", "none", CleanupOptions{}, scope); selected {
				t.Fatal("automatic cleanup selected unproven contribution")
			}
			if _, err := requalifyLandedRow(j, root, planned, CleanupOptions{}, scope); !errors.Is(err, errStaleFingerprint) {
				t.Fatalf("requalification = %v, want stale refusal", err)
			}
			requirePresent(t, sibling.Path, "unproven sibling worktree")
			if !assignmentActive(t, root, assignment.ID) {
				t.Fatal("unproven assignment retired")
			}
		})
	}
}
