// Merge refusal fixtures for the grade of the target tip alone, which picks the face of a
// fold red that names no cause: a lane whose subject follows its base, so the grade must
// measure the target against the incoming commit, and the two infrastructure arms that
// stop the grade.
package worktree

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// namedRedLane is redFileLane over a subject: it grades only when the graded tree names
// changed Markdown against its base. The target tip's own redFile reaches it only when the
// grade measures the target against the incoming commit.
var namedRedLane = gate.Phase{Name: "unit", Argv: []string{"sh", "-c", "[ $# -eq 0 ] || " + redFileLane.Argv[2], "lane", gate.LaneNamedMarkdownToken}}

// commitRouteNames states what a route that commits in the set's assignment at index names.
func commitRouteNames(index int) func(f mergeSet) []string {
	return func(f mergeSet) []string { return []string{"bench commit --in ", f.created[index].Assignment.Label} }
}

// commitWholeGate commits the whole gate at dir, with script as its body.
func commitWholeGate(t *testing.T, dir, script string) {
	t.Helper()
	landingGateFixture(t).MustWrite(t, dir, script+"\n", script+"\n")
	gitRun(t, dir, "add", ".bench")
	gitRun(t, dir, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "declare the whole gate")
}

// namedLaneTarget gives the set's first target namedRedLane, and a whole gate whose body is
// script. The whole gate reads the other verdict, so a grade that skips the lane picks the
// other face.
func namedLaneTarget(t *testing.T, f mergeSet, script string) mergeSet {
	t.Helper()
	commitWholeGate(t, f.created[0].Path, script)
	commitLaneManifest(t, f.created[0].Path, namedRedLane)
	return f
}

// targetGradeFault produces the infrastructure face of one cause that stops the grade of the
// target tip alone. The target's lane fails once: the fold's run plants the fault and fails,
// and every later run passes. fault answers the shell step that plants the fault, where $1
// is the base that the lane measures against. clear removes the fault.
func targetGradeFault(cause string, fault func(f mergeSet) string, clear func(t *testing.T, f mergeSet)) mergeRefusalFixture {
	return mergeRefusalFixture{
		face:  faceMergeInfrastructure,
		cause: cause,
		exit:  1,
		build: func(t *testing.T) (mergeSet, []string) {
			f := mergeFixture(t, "integration")
			armed := filepath.Join(t.TempDir(), "armed")
			mustWrite(t, armed, nil, 0o644)
			script := "[ -e " + sanitize.ShellQuote(armed) + " ] || exit 0; rm " + sanitize.ShellQuote(armed) + "; " + fault(f) + "; exit 1"
			commitLaneManifest(t, f.created[0].Path, gate.Phase{Name: "unit", Argv: []string{"sh", "-c", script, "lane", gate.LaneBaseToken}})
			return f, mergeTargetArgs(f, commitOnDefault(t, f.root, "incoming.txt", "incoming\n"))
		},
		clear: clear,
	}
}

// The unreadable-tree fault moves the target tip's loose commit object out of the set's
// object store to hiddenTip, and its clear moves the object back.
func setObjects(f mergeSet) string { return filepath.Join(f.root, ".git", "objects") }
func hiddenTip(f mergeSet) string  { return filepath.Join(f.root, ".git", "hidden-tip") }

func targetGradeFixtures() []mergeRefusalFixture {
	return []mergeRefusalFixture{
		{
			// The target tip alone holds the red file, which its lane grades red and its whole
			// gate grades green. Main adds Markdown, so the fold's lane grades too.
			face:  faceMergeTargetRed,
			cause: "named lane",
			exit:  1,
			build: func(t *testing.T) (mergeSet, []string) {
				f := namedLaneTarget(t, admittedMergeFixture(t, "integration"), "exit 0")
				commitInWorktree(t, f.created[0].Path, redFile, "red\n", "inherited red")
				return f, mergeTargetArgs(f, commitOnDefault(t, f.root, "incoming.md", "incoming\n"))
			},
			names: commitRouteNames(0),
			carry: map[int]func(*testing.T, mergeSet){0: func(t *testing.T, f mergeSet) {
				mustRemove(t, filepath.Join(f.created[0].Path, redFile))
			}},
		},
		// The target tip alone passes its lane and fails its whole gate.
		foldRedFixture("named lane", func(t *testing.T) mergeSet {
			return namedLaneTarget(t, mergeFixture(t, "integration"), "exit 1")
		}),
		// The fold's lane breaks the target's manifest, so the grade cannot resolve its lane.
		targetGradeFault("target lane unreadable", func(f mergeSet) string {
			return "printf '{' > " + sanitize.ShellQuote(filepath.Join(f.created[0].Path, ".bench", "phases.json"))
		}, func(t *testing.T, f mergeSet) {
			gitRun(t, f.created[0].Path, "checkout", "-q", "--", ".bench/phases.json")
		}),
		// The fold's lane hides the target tip's commit, so the grade cannot read its tree.
		targetGradeFault("target tree unreadable", func(f mergeSet) string {
			return `mv ` + sanitize.ShellQuote(setObjects(f)) + `/"${1%"${1#??}"}/${1#??}" ` + sanitize.ShellQuote(hiddenTip(f))
		}, func(t *testing.T, f mergeSet) {
			tip := gitOutput(t, f.created[0].Path, "rev-parse", "HEAD")
			mustNoError(t, os.Rename(hiddenTip(f), filepath.Join(setObjects(f), tip[:2], tip[2:])))
		}),
	}
}
