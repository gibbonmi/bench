// Refusal route tests for the merge: each merge face's producing fixture makes the merge
// print the face, follows the route the merge printed, and reruns the merge to exit 0.
package worktree

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// mergeRefusalFixture produces exactly one merge face. build makes the merge set and
// breaks it so the merge prints that face; it returns the set and the merge's arguments.
// The registry walk requires one fixture per merge face, so a face added with no fixture
// turns TestMergeFacesFollowTheirRoutes red.
type mergeRefusalFixture struct {
	face  string
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

// TestMergeFacesFollowTheirRoutes is RR28 and RR29 for the merge faces. The registry is the
// source of the merge's face set, so each merge face needs exactly one producing fixture.
// Each fixture's printed route is carried out step by step, as the landing walk carries
// out its own, and the merge then reruns to exit 0.
func TestMergeFacesFollowTheirRoutes(t *testing.T) {
	t.Parallel()
	faces := map[string]refusalroute.Face{}
	for _, face := range refusalroute.Faces(refusalroute.Merge) {
		faces[face.Name] = face
	}
	produced := map[string]bool{}
	for _, fixture := range mergeRefusalFixtures() {
		if produced[fixture.face] {
			t.Fatalf("merge face %q has two producing fixtures", fixture.face)
		}
		if _, ok := faces[fixture.face]; !ok {
			t.Fatalf("fixture %q produces no registered merge face", fixture.face)
		}
		produced[fixture.face] = true
	}
	for name := range faces {
		if !produced[name] {
			t.Errorf("registry merge face %q has no producing fixture", name)
		}
	}
	wrapper := installedWrapper(t, testRunBinary(t))
	for _, fixture := range mergeRefusalFixtures() {
		t.Run(fixture.face, func(t *testing.T) {
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

// mergeOperatorFill fills the slots a printed merge route leaves to the operator: the
// commit slots of the assignment checkout that a commit step addresses.
func mergeOperatorFill(t *testing.T, f mergeSet, step string) string {
	t.Helper()
	for _, created := range f.created {
		if strings.Contains(step, " --in "+sanitize.ShellQuote(created.Assignment.Label)+" ") {
			return strings.NewReplacer(commitSlots(t, created.Path)...).Replace(step)
		}
	}
	return step
}

// admittedMergeFixture is the merge fixture whose assignments each hold a delivery
// binding, so a printed `bench commit` in one of them passes the commitment policy.
func admittedMergeFixture(t *testing.T, labels ...string) mergeSet {
	t.Helper()
	root := newWorktreeRepo(t)
	mustMkdirAll(t, filepath.Join(root, filepath.Dir(closureSpec)), 0o755)
	mustWrite(t, filepath.Join(root, closureSpec), []byte("# x\n\nStatus: staged\n"), 0o644)
	commitmenttest.SeedAdmission(t, root, closureSpec)
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "approve the merge fixture delivery")
	f := mergeSetAt(t, root, filepath.Join(t.TempDir(), "bench-home"))
	for _, label := range labels {
		created := mustCreate(t, f.root, f.home, "merge-"+label, label)
		commitmenttest.Admit(t, created.Path, "merge-"+label, closureSpec)
		f.created = append(f.created, created)
	}
	return f
}

// mergeTargetArgs are the arguments of a merge of incoming into the set's first
// assignment.
func mergeTargetArgs(f mergeSet, incoming string) []string {
	return []string{"--from", incoming, f.created[0].Assignment.ID}
}

func mergeRefusalFixtures() []mergeRefusalFixture {
	label := func(index int) func(f mergeSet) []string {
		return func(f mergeSet) []string { return []string{"bench commit --in ", f.created[index].Assignment.Label} }
	}
	return []mergeRefusalFixture{
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
				root := newWorktreeRepo(t)
				landingGateFixture(t).MustWrite(t, root, "set -eu\n", "set -eu\n")
				gitRun(t, root, "add", ".bench")
				gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "declare the whole gate")
				f := mergeSet{repoHome: repoHome{root: root, home: filepath.Join(t.TempDir(), "bench-home")}, joins: defaultJoins(), kit: t.TempDir()}
				f.created = append(f.created, mustCreate(t, f.root, f.home, "merge-infrastructure", "integration"))
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
