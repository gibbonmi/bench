package commit

// The commit verb and another checkout's whole green of its composed tree. A commit with no
// lane looks up that green before it runs the gate, and a commit with a declared lane runs
// its lane whatever the store holds.

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// gradeComposition writes a.txt in root, grades the tree that a commit of a.txt composes, and
// answers that tree. With fromPrimary, the primary checkout runs the gate, and otherwise root.
func gradeComposition(t *testing.T, root, before string, fromPrimary bool) string {
	t.Helper()
	gittest.MustRun(t, root, "reset", "-q", "--hard", "HEAD")
	mustWrite(t, filepath.Join(root, "a.txt"), "landed\n", 0o644)
	tree, err := landing.CandidateTree(root, before, []string{"a.txt"})
	if err != nil {
		t.Fatal(err)
	}
	grader := root
	if fromPrimary {
		grader = filepath.Dir(gittest.Output(t, root, "rev-parse", "--path-format=absolute", "--git-common-dir"))
	}
	var stdout, stderr bytes.Buffer
	if result := gate.ExecuteTree(context.Background(), grader, tree, &stdout, &stderr); result.ActionExit != 0 {
		t.Fatalf("run of the composed tree = %+v, stderr=%q", result, stderr.String())
	}
	if inspection := gate.InspectTree(root, tree); !inspection.ReusableGreen {
		t.Fatalf("inspection from the commit checkout = %+v, want a reusable green", inspection)
	}
	return tree
}

// GR49: in a repository with no lane, a commit of a composition that another checkout graded
// green adds no whole-gate run.
func TestCommitReusesAGreenOfTheComposedTree(t *testing.T) {
	tally := filepath.Join(t.TempDir(), "gate-tally")
	root, before := landingRepo(t, 0, func(t *testing.T, root string) {
		mustWrite(t, filepath.Join(root, ".bench", "gate.sh"),
			"#!/bin/sh\nprintf g >> "+sanitize.ShellQuote(tally)+"\nexit 0\n", 0o755)
	})
	tree := gradeComposition(t, root, before, true)

	code, stdout, stderr := runCommand(t, root, "-m", "m", "a.txt")
	if code != 0 || headTree(t, root) != tree {
		t.Fatalf("commit = (%d, %q, %q), want the graded tree %s published", code, stdout, stderr, tree)
	}
	if recorded, err := os.ReadFile(tally); err != nil || string(recorded) != "g" {
		t.Fatalf("gate tally = %q, %v; want the one run of the other checkout", recorded, err)
	}
}

// GR51: a commit with a declared lane runs its lane when a whole green of the composed tree
// exists, so the failing lane check refuses the commit.
func TestLaneCommitKeepsItsLane(t *testing.T) {
	root, before := laneRepo(t, 1, noWrite)
	gradeComposition(t, root, before, false)

	code, stdout, stderr := runCommand(t, root, "-m", "m", "a.txt")
	if code != 1 || !strings.Contains(stdout, "lane{outcome=fail,check=check}") || head(t, root) != before {
		t.Fatalf("lane commit = (%d, %q, %q), want the failing lane check to refuse", code, stdout, stderr)
	}
}
