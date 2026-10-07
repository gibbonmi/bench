// Refusal route tests for the merge: each merge face's producing fixture makes the merge
// print the face, follows the route the merge printed, and reruns the merge to exit 0.
package worktree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// mergeRefusalFixture produces exactly one merge face. build makes the merge set and
// breaks it so the merge prints that face; it returns the set and the merge's arguments.
// The registry walk requires a fixture for each merge face, so a face added with no
// fixture turns TestMergeFacesFollowTheirRoutes red.
type mergeRefusalFixture struct {
	face string
	// cause names the cause of a face that more than one cause raises, so each cause has
	// its own producing fixture.
	cause string
	build func(t *testing.T) (mergeSet, []string)
	// exit is the exit the face prints at: 1 for a refusal, 3 for a published merge whose
	// checkout did not reconcile.
	exit int
	// clear removes the fixture's own fault scaffolding once the face printed, which no
	// operator's tree holds, so the route repairs only the cause the face names.
	clear func(t *testing.T, f mergeSet)
	// carry carries out a printed route step by the fixture's own means, keyed by the
	// step's index in the face's route: an instruction, or a step of a reviewer route.
	carry map[int]func(t *testing.T, f mergeSet)
	// after carries out the route that the last printed step's own verb printed, such as
	// the apply of a reset plan.
	after func(t *testing.T, wrapper string, f mergeSet, last verbResult)
	// names states what the printed route must name, and absent what it must not.
	names, absent func(f mergeSet) []string
}

// producedMerge is what the merge printed for one face's producing fixture: the set, the
// merge's arguments, the run, and the next= value of the face's record.
type producedMerge struct {
	f    mergeSet
	args []string
	r    verbResult
	next string
}

// produceMergeFace drives one producing fixture through the merge and returns the face's
// next= value. It fails t unless the merge exits as the fixture declares and the face
// printed a non-empty route.
func produceMergeFace(t *testing.T, fixture mergeRefusalFixture) producedMerge {
	t.Helper()
	f, args := fixture.build(t)
	r := runVerb(t, verbMerge, f.merge(args...))
	record := "refused{detail=" + refusalroute.Sentence(fixture.face)
	if fixture.exit == 3 {
		record = "merged{"
	}
	next, printed := recordField(r.stdout, record, refusalroute.NextField)
	if r.exit != fixture.exit || !printed || next == "" {
		t.Fatalf("face %s = (%d, %q, %q), want exit %d and a non-empty next= field", fixture.face, r.exit, r.stdout, r.stderr, fixture.exit)
	}
	return producedMerge{f: f, args: args, r: r, next: next}
}

// mergeFixtureFor returns the one fixture that produces the named merge face, so a proof
// that raises one face composes its fixture rather than copying it.
func mergeFixtureFor(t *testing.T, face string) mergeRefusalFixture {
	t.Helper()
	for _, fixture := range mergeRefusalFixtures() {
		if fixture.face == face {
			return fixture
		}
	}
	t.Fatalf("no producing fixture for merge face %q", face)
	return mergeRefusalFixture{}
}

// TestMergeFacesFollowTheirRoutes is RR28 and RR29 for the merge faces, and RR26, RR56, and
// RR57 for the fold red. The registry is the source of the merge's face set, so each merge
// face needs a producing fixture for each of its causes. Each fixture's printed route is
// carried out step by step, as the landing walk carries out its own, and the merge then
// reruns to exit 0.
func TestMergeFacesFollowTheirRoutes(t *testing.T) {
	t.Parallel()
	faces := map[string]refusalroute.Face{}
	for _, face := range refusalroute.Faces(refusalroute.Merge) {
		faces[face.Name] = face
	}
	produced := map[string]bool{}
	for _, fixture := range mergeRefusalFixtures() {
		if produced[fixture.face+"/"+fixture.cause] {
			t.Fatalf("merge face %q has two producing fixtures for cause %q", fixture.face, fixture.cause)
		}
		if _, ok := faces[fixture.face]; !ok {
			t.Fatalf("fixture %q produces no registered merge face", fixture.face)
		}
		produced[fixture.face], produced[fixture.face+"/"+fixture.cause] = true, true
	}
	for name := range faces {
		if !produced[name] {
			t.Errorf("registry merge face %q has no producing fixture", name)
		}
	}
	wrapper := installedWrapper(t, testRunBinary(t))
	for _, fixture := range mergeRefusalFixtures() {
		t.Run(strings.TrimSuffix(fixture.face+"/"+fixture.cause, "/"), func(t *testing.T) {
			t.Parallel()
			face := faces[fixture.face]
			p := produceMergeFace(t, fixture)
			var names []string
			if fixture.names != nil {
				names = fixture.names(p.f)
			}
			if fixture.absent != nil {
				for _, name := range fixture.absent(p.f) {
					if strings.Contains(p.next, name) {
						t.Fatalf("%s next = %q, want it not to name %q", fixture.face, p.next, name)
					}
				}
			}
			steps := printedSteps(t, face, p.next, "", names)
			if fixture.clear != nil {
				fixture.clear(t, p.f)
			}
			carry := func(index int) (func(), bool) {
				step, carried := fixture.carry[index]
				return func() { step(t, p.f) }, carried
			}
			last := followRoute(t, face, steps, carry, func(step string) verbResult {
				return runPrintedStep(t, wrapper, p.f.repoHome, mergeOperatorFill(t, p.f, step))
			})
			if fixture.after != nil {
				fixture.after(t, wrapper, p.f, last)
			}
			r := runVerb(t, verbMerge, p.f.merge(p.args...))
			if r.exit != 0 || !strings.Contains(r.stdout, "merged{") {
				t.Fatalf("%s merge re-run = (%d, %q, %q), want exit 0 and a merged record", fixture.face, r.exit, r.stdout, r.stderr)
			}
		})
	}
}

// WM40: the conflict refusal names the hand repair in the order the operator runs it, up
// to the commit that records the resolution. An empty next= leaves the operator to invent
// the repair the landing already spells. The guard reserves the raw merge to the reviewer,
// so the route opens with the reviewer marker.
func TestMergeConflictRefusalNamesTheHandRepair(t *testing.T) {
	t.Parallel()
	f := mergeFixture(t, "integration")
	target := f.created[0]
	commitInWorktree(t, target.Path, "tracked.txt", "target edit\n", "target edit")
	incoming := commitOnDefault(t, f.root, "tracked.txt", "incoming edit\n")

	r := runVerb(t, verbMerge, f.merge("--from", incoming, target.Assignment.ID))
	wantNext := "next=reviewer: git -C '" + target.Assignment.Worktree + "' merge '" + incoming +
		"' (bench worktree merge refuses this conflict; resolve it by hand); then bench commit"
	if r.exit != 1 || !strings.Contains(r.stdout, wantNext) {
		t.Fatalf("conflict repair next = (%d, %q, %q), want %q", r.exit, r.stdout, r.stderr, wantNext)
	}
}

// TestRedSourceFoldNamesAnExit is collision 5a: a red source folds a moved main. RR21-RR25
// and RR55 follow the lane fail out: the target tip alone fails its lane, so the fold names
// the target's own repair, and a fold into the target that the repair left dirty names the
// commit. After the red file goes and the printed commit runs, the printed fold publishes.
// RR59 is the same red under a whole gate, whose inherited kind attributes no red either.
func TestRedSourceFoldNamesAnExit(t *testing.T) {
	t.Parallel()
	// requireTargetRed returns the steps of a fold refusal's route, which the agent runs:
	// the commit in the target, and last the fold of main into the target again.
	requireTargetRed := func(t *testing.T, r verbResult, target Creation) []string {
		t.Helper()
		next, printed := recordField(r.stdout, "refused{", refusalroute.NextField)
		commit := "bench commit --in " + sanitize.ShellQuote(target.Assignment.Label) + " "
		fold := "bench worktree merge --from " + sanitize.ShellQuote("main") + " " + sanitize.ShellQuote(target.Assignment.ID)
		if r.exit != 1 || !printed || strings.HasPrefix(next, reviewerRoute) || !strings.Contains(next, commit) || !strings.HasSuffix(next, fold) {
			t.Fatalf("red-source fold = (%d, %q), want an agent next= that commits with %q and ends with %q", r.exit, r.stdout, commit, fold)
		}
		return refusalroute.Steps(next)
	}
	t.Run("lane fail", func(t *testing.T) {
		t.Parallel()
		f, args := mergeFixtureFor(t, faceMergeTargetRed).build(t)
		target := f.created[0]
		fold := runVerb(t, verbMerge, f.merge(args...))
		requireMergeLaneRefusal(t, fold, "unit")
		steps := requireTargetRed(t, fold, target)
		// The operator starts the repair, so the target checkout is dirty.
		mustRemove(t, filepath.Join(target.Path, redFile))
		dirty := runVerb(t, verbMerge, f.merge(args...))
		requireMergeRefusal(t, dirty, "not clean")
		if next, _ := recordField(dirty.stdout, "refused{", refusalroute.NextField); !strings.Contains(next, "bench commit --in ") {
			t.Fatalf("dirty fold next = %q, want the commit route", next)
		}
		wrapper := installedWrapper(t, testRunBinary(t))
		runPrintedStep(t, wrapper, f.repoHome, mergeOperatorFill(t, f, steps[len(steps)-2]))
		if r := runPrintedStep(t, wrapper, f.repoHome, steps[len(steps)-1]); !strings.Contains(r.stdout, "merged{") {
			t.Fatalf("printed fold = %q, want a merged record", r.stdout)
		}
	})
	t.Run("inherited", func(t *testing.T) {
		t.Parallel()
		f := wholeGateMergeFixture(t, redFileAbsent)
		target := f.created[0]
		commitInWorktree(t, target.Path, redFile, "red\n", "inherited red")
		commitOnDefault(t, f.root, "incoming.txt", "incoming\n")
		requireTargetRed(t, runVerb(t, verbMerge, f.merge(mergeTargetArgs(f, "main")...)), target)
	})
}

// foldRedFixture produces the fold red of one red cause: main adds the red file to a target
// whose tip alone grades green. target makes the target set, and it returns the set.
func foldRedFixture(cause string, target func(t *testing.T) mergeSet) mergeRefusalFixture {
	return mergeRefusalFixture{
		face:  faceMergeFoldRed,
		cause: cause,
		exit:  1,
		build: func(t *testing.T) (mergeSet, []string) {
			f := target(t)
			commitOnDefault(t, f.root, redFile, "red\n")
			return f, mergeTargetArgs(f, "main")
		},
		names: func(f mergeSet) []string { return []string{"FT342", f.created[0].Assignment.Label} },
		carry: map[int]func(*testing.T, mergeSet){
			// The reviewer reads the cause: the target tip alone holds no red.
			0: func(t *testing.T, f mergeSet) {
				if _, err := os.Stat(filepath.Join(f.created[0].Path, redFile)); !os.IsNotExist(err) {
					t.Fatalf("target holds %s: %v, want the red from main alone", redFile, err)
				}
			},
			// The reviewer decides that main repairs the red it added.
			1: func(t *testing.T, f mergeSet) {
				gitRun(t, f.root, "rm", "-q", redFile)
				gitRun(t, f.root, "commit", "-q", "-m", "repair the red main added")
			},
		},
	}
}

func mergeRefusalFixtures() []mergeRefusalFixture {
	label := func(index int) func(f mergeSet) []string {
		return func(f mergeSet) []string { return []string{"bench commit --in ", f.created[index].Assignment.Label} }
	}
	return []mergeRefusalFixture{
		{
			// The target tip alone fails its lane, and main moves by a commit that adds no red.
			face: faceMergeTargetRed,
			exit: 1,
			build: func(t *testing.T) (mergeSet, []string) {
				f := admittedMergeFixture(t, "integration")
				commitLaneManifest(t, f.created[0].Path, redFileLane)
				commitInWorktree(t, f.created[0].Path, redFile, "red\n", "inherited red")
				commitOnDefault(t, f.root, "incoming.txt", "incoming\n")
				return f, mergeTargetArgs(f, "main")
			},
			names: label(0),
			// The operator's repair: the red file goes, and the printed commit records it.
			carry: map[int]func(*testing.T, mergeSet){0: func(t *testing.T, f mergeSet) {
				mustRemove(t, filepath.Join(f.created[0].Path, redFile))
			}},
		},
		foldRedFixture("lane fail", func(t *testing.T) mergeSet {
			f := mergeFixture(t, "integration")
			commitLaneManifest(t, f.created[0].Path, redFileLane)
			return f
		}),
		foldRedFixture("inherited", func(t *testing.T) mergeSet { return wholeGateMergeFixture(t, redFileAbsent) }),
		// The target's own tip carries the project-green marker and its green evidence, so
		// the gate attributes the fold's red to the fold.
		foldRedFixture("candidate", func(t *testing.T) mergeSet {
			f := wholeGateMergeFixture(t, redFileAbsent)
			markTargetGreen(t, f.created[0])
			return f
		}),
		{
			face: faceMergeTargetNotClean,
			exit: 1,
			build: func(t *testing.T) (mergeSet, []string) {
				f := admittedMergeFixture(t, "integration")
				incoming := commitOnDefault(t, f.root, "incoming.txt", "incoming\n")
				mustWrite(t, filepath.Join(f.created[0].Path, "target-work.txt"), []byte("uncommitted\n"), 0o644)
				return f, mergeTargetArgs(f, incoming)
			},
			names: label(0),
		},
		{
			face: faceMergeSiblingNotClean,
			exit: 1,
			build: func(t *testing.T) (mergeSet, []string) {
				f := admittedMergeFixture(t, "integration", "delegate")
				sibling := f.created[1]
				commitInWorktree(t, sibling.Path, "sibling.txt", "sibling\n", "sibling work")
				mustWrite(t, filepath.Join(sibling.Path, "sibling.txt"), []byte("uncommitted\n"), 0o644)
				return f, mergeTargetArgs(f, sibling.Assignment.Label)
			},
			names:  label(1),
			absent: func(mergeSet) []string { return []string{"bench worktree exec"} },
		},
		{
			face: faceMergeConflict,
			exit: 1,
			build: func(t *testing.T) (mergeSet, []string) {
				f := mergeFixture(t, "integration")
				commitInWorktree(t, f.created[0].Path, "tracked.txt", "target edit\n", "target edit")
				incoming := commitOnDefault(t, f.root, "tracked.txt", "incoming edit\n")
				return f, mergeTargetArgs(f, incoming)
			},
			// The reviewer's hand merge of the incoming commit into the target, its
			// resolution, and the commit that records it.
			carry: map[int]func(*testing.T, mergeSet){
				0: func(t *testing.T, f mergeSet) {
					target := f.created[0].Path
					merge := descendant(t, "git", "-C", target, "merge", "--no-commit", "main")
					if out, err := merge.CombinedOutput(); err == nil || !strings.Contains(string(out), "CONFLICT") {
						t.Fatalf("hand merge = %v, %s; want the conflict", err, out)
					}
					mustWrite(t, filepath.Join(target, "tracked.txt"), []byte("resolved edit\n"), 0o644)
					gitRun(t, target, "add", "tracked.txt")
				},
				1: func(t *testing.T, f mergeSet) {
					gitRun(t, f.created[0].Path, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "resolve the composition")
				},
			},
		},
		{
			// The target declares no lane, so the whole gate grades the composed tree. The
			// gate cannot open its lock until the fixture frees it, so the fault outlasts the
			// merge's single retry and the face prints.
			face: faceMergeInfrastructure,
			exit: 1,
			build: func(t *testing.T) (mergeSet, []string) {
				f := wholeGateMergeFixture(t, "set -eu")
				commitInWorktree(t, f.created[0].Path, "target.txt", "target\n", "target work")
				incoming := commitOnDefault(t, f.root, "incoming.txt", "incoming\n")
				blockGateLock(t, f.created[0].Path)
				return f, mergeTargetArgs(f, incoming)
			},
			clear: func(t *testing.T, f mergeSet) { mustRemove(t, gateLock(t, f.created[0].Path)) },
		},
		{
			face:  faceMergePublishedUnreconciled,
			exit:  3,
			build: reconcileFailingMerge,
			// The lane's stale index lock is the fixture's fault, and the reset needs the
			// index lock too.
			clear: func(t *testing.T, f mergeSet) { mustRemove(t, mustAdminPath(t, f.created[0].Path, "index.lock")) },
			// The reset plan prints its own apply, which the operator runs next.
			after: func(t *testing.T, wrapper string, f mergeSet, plan verbResult) {
				apply, printed := recordField(plan.stdout, "reset_plan{", refusalroute.NextField, "envelope", "fingerprint")
				if !printed || apply == "" {
					t.Fatalf("reset plan = %q, want a next= apply", plan.stdout)
				}
				runPrintedStep(t, wrapper, f.repoHome, apply)
			},
		},
		{
			// A detached target checkout is a state the merge's own bookkeeping cannot
			// name a repair for, so it hands back.
			face: faceMergeHandback,
			exit: 1,
			build: func(t *testing.T) (mergeSet, []string) {
				f := mergeFixture(t, "integration")
				incoming := commitOnDefault(t, f.root, "incoming.txt", "incoming\n")
				gitRun(t, f.created[0].Path, "checkout", "-q", "--detach", "HEAD")
				return f, mergeTargetArgs(f, incoming)
			},
			carry: map[int]func(*testing.T, mergeSet){0: func(t *testing.T, f mergeSet) {
				gitRun(t, f.created[0].Path, "checkout", "-q", strings.TrimPrefix(f.created[0].Assignment.Branch, "refs/heads/"))
			}},
		},
	}
}

// reconcileFailingMerge is the merge whose publication outruns its reconcile: the target's
// lane leaves a stale index lock, so the merge publishes and the checkout cannot follow.
func reconcileFailingMerge(t *testing.T) (mergeSet, []string) {
	t.Helper()
	f := mergeFixture(t, "integration")
	commitInWorktree(t, f.created[0].Path, "target.txt", "target\n", "target work")
	commitIndexLockLane(t, f.created[0].Path)
	incoming := commitOnDefault(t, f.root, "incoming.txt", "incoming\n")
	return f, mergeTargetArgs(f, incoming)
}
