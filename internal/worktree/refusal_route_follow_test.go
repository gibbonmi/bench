// Refusal route follow tests for the landing: each face's producing fixture carries out
// the route the landing printed, re-runs the landing, and the landing completes. The route
// walk helpers here serve the merge walk too, with the merge fixture sets it builds on and
// the reset fixtures it walks.
package worktree

import (
	"bytes"
	"context"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/env"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/gate/greenmarker"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/shellcommand"
)

// TestLandingFacesFollowTheirRoutes is RR15, with RR18 and RR19 folded in. Each producing
// fixture's printed route is carried out step by step: the fixture carries out the steps it
// declares by its own means, and every other step, each agent command step included, runs
// verbatim with only its operator slots filled. The last step re-runs the landing, and the
// landing completes, so the route finishes the recovery rather than trading the face for
// another refusal. A reviewer face's route opens with the reviewer marker and an agent
// face's does not.
func TestLandingFacesFollowTheirRoutes(t *testing.T) {
	t.Parallel()
	wrapper := installedWrapper(t, testRunBinary(t))
	faces := map[string]refusalroute.Face{}
	for _, face := range refusalroute.Faces(refusalroute.Land) {
		faces[face.Name] = face
	}
	for _, fixture := range landingRefusalFixtures() {
		t.Run(fixture.face, func(t *testing.T) {
			t.Parallel()
			face := faces[fixture.face]
			produced := produceLandingFace(t, fixture)
			var names []string
			if fixture.names != nil {
				names = fixture.names(produced.f)
			}
			steps := printedSteps(t, face, produced.next, laterProofsSkipped, names)
			if fixture.clear != nil {
				fixture.clear(t, produced.f)
			}
			fill := operatorFill(produced)
			carry := func(index int) (func(), bool) {
				step, carried := fixture.carry[index]
				return func() { step(t, produced.f) }, carried
			}
			last := followRoute(t, face, steps, carry, func(step string) verbResult {
				return runPrintedStep(t, wrapper, produced.f.repoHome, fill(t, step))
			})
			if _, landed := landedField(last.stdout, "worktree"); last.exit != 0 || !landed {
				t.Fatalf("%s re-run = (%d, %q, %q), want exit 0 and a landed record", fixture.face, last.exit, last.stdout, last.stderr)
			}
		})
	}
}

// printedSteps checks one printed route against the face that printed it and returns its
// steps. The reviewer marker opens the route exactly when the face is the reviewer's, the
// route names each of names, and after the preface it prints the steps the face declares.
func printedSteps(t *testing.T, face refusalroute.Face, next, preface string, names []string) []string {
	t.Helper()
	route, reviewer := strings.CutPrefix(next, reviewerRoute)
	if reviewer != (face.Authority == refusalroute.Reviewer) {
		t.Fatalf("%s next = %q, want the reviewer marker exactly when the face is the reviewer's", face.Name, next)
	}
	for _, name := range names {
		if !strings.Contains(route, name) {
			t.Fatalf("%s next = %q, want it to name %q", face.Name, next, name)
		}
	}
	steps := refusalroute.Steps(strings.TrimPrefix(route, preface+"; "))
	if len(steps) != len(face.Route) {
		t.Fatalf("%s next = %q prints %d steps, want the %d the face declares", face.Name, next, len(steps), len(face.Route))
	}
	return steps
}

// followRoute carries out a printed route step by step and returns the result of the last
// step it ran. carry answers the fixture's own means for one step, by the step's index in
// the face's route: an instruction, or a step of a reviewer route. run runs every other
// step verbatim. An agent command step always runs verbatim, so the route the agent reads
// is the route the walk proves.
func followRoute(t *testing.T, face refusalroute.Face, steps []string, carry func(index int) (func(), bool), run func(step string) verbResult) verbResult {
	t.Helper()
	var last verbResult
	for index, step := range face.Route {
		carried, ok := carry(index)
		switch {
		case ok && face.Authority == refusalroute.Agent && step.IsCommand():
			t.Fatalf("%s step %d is an agent command, which the walk runs verbatim; the fixture may not carry it", face.Name, index+1)
		case ok:
			carried()
		case step.IsCommand():
			last = run(steps[index])
		default:
			t.Fatalf("%s step %d %q is an instruction that the fixture does not carry out", face.Name, index+1, steps[index])
		}
	}
	return last
}

// diagnosticStep reports whether a printed step is `bench doctor`. Its exit reports the
// health of the whole install, and a fixture repository is never a linked one, so a walk
// runs the step verbatim and does not grade its exit. Every other Bench verb keeps the
// exit-0 rule.
func diagnosticStep(words []string) bool {
	return len(words) == 2 && words[0] == "bench" && words[1] == "doctor"
}

// operatorFill fills the slots a printed route leaves to the operator, with the values the
// operator holds: the paths it changed in the source, its own request and message, and the
// source tip after its repair. A slot the walk has no value for stays, and the step then
// refuses to run.
func operatorFill(produced producedFace) func(t *testing.T, step string) string {
	return func(t *testing.T, step string) string {
		t.Helper()
		source := produced.f.creation.Path
		tip := gitOutput(t, source, "rev-parse", "HEAD")
		return strings.NewReplacer(append(commitSlots(t, source),
			"<message>", "'land the followed route'",
			"<request>", sanitize.ShellQuote(produced.request),
			repairedTipArg, sanitize.ShellQuote(tip),
		)...).Replace(step)
	}
}

// commitSlots are the replacement pairs for the operator slots of a printed commit step:
// the message, and the uncommitted paths of the checkout the step commits in.
func commitSlots(t *testing.T, checkout string) []string {
	t.Helper()
	var paths []string
	for _, path := range strings.Fields(gitOutput(t, checkout, "ls-files", "--modified", "--others", "--exclude-standard")) {
		paths = append(paths, sanitize.ShellQuote(path))
	}
	return []string{"<msg>", "'follow the refusal route'", "<path>...", strings.Join(paths, " ")}
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

// wholeGateMergeFixture is the merge fixture whose root declares no lane, so the whole gate
// grades each fold. script is the gate's body, and the set holds one assignment.
func wholeGateMergeFixture(t *testing.T, script string) mergeSet {
	t.Helper()
	f := mergeSetOver(t, newWorktreeRepo(t), filepath.Join(t.TempDir(), "bench-home"))
	commitWholeGate(t, f.root, script)
	f.created = append(f.created, mustCreate(t, f.root, f.home, "merge-integration", "integration"))
	return f
}

// The red-source fixtures grade one file: a tree that holds redFile is red, under the lane
// and under the whole gate alike. It is Markdown, so a lane over the named Markdown names it.
const (
	redFile       = "red.md"
	redFileAbsent = "[ ! -e " + redFile + " ]"
)

var redFileLane = gate.Phase{Name: "unit", Argv: []string{"sh", "-c", redFileAbsent + " || { echo " + redFile + " is present; exit 1; }"}}

// markTargetGreen grades the target's own tip with the whole gate and advances the target
// branch's project-green marker to it, as a landing advances main's.
func markTargetGreen(t *testing.T, target Creation) {
	t.Helper()
	if r := gate.Execute(context.Background(), target.Path, io.Discard, io.Discard); r.ActionExit != 0 {
		t.Fatalf("target gate = %+v, want green", r)
	}
	branch := strings.TrimPrefix(target.Assignment.Branch, "refs/heads/")
	mustNoError(t, greenmarker.Advance(target.Path, branch, gitOutput(t, target.Path, "rev-parse", "HEAD"), ""))
	if !gate.ValidateProjectGreen(target.Path, branch).ReusableGreen {
		t.Fatal("the target tip carries no reusable project green")
	}
}

// producingFixtures requires a producing fixture for each cause of each face of verb, and a
// registered face for each fixture. keys holds each fixture's face and cause. It returns the
// verb's faces by name.
func producingFixtures(t *testing.T, verb refusalroute.Verb, keys [][2]string) map[string]refusalroute.Face {
	t.Helper()
	faces := map[string]refusalroute.Face{}
	for _, face := range refusalroute.Faces(verb) {
		faces[face.Name] = face
	}
	produced := map[string]bool{}
	for _, key := range keys {
		if produced[key[0]+"/"+key[1]] {
			t.Fatalf("%s face %q has two producing fixtures for cause %q", verb, key[0], key[1])
		}
		if _, ok := faces[key[0]]; !ok {
			t.Fatalf("fixture %q produces no registered %s face", key[0], verb)
		}
		produced[key[0]], produced[key[0]+"/"+key[1]] = true, true
	}
	for name := range faces {
		if !produced[name] {
			t.Errorf("registry %s face %q has no producing fixture", verb, name)
		}
	}
	return faces
}

// resetRefusalFixture produces exactly one reset face. build makes an owned assignment and
// breaks it so the reset prints that face; it returns the assignment and the reset's
// arguments. cause, carry, and after are the merge fixture's fields over the assignment.
type resetRefusalFixture struct {
	face  string
	cause string
	build func(t *testing.T) (ownedAssignment, []string)
	carry map[int]func(t *testing.T, f ownedAssignment)
	after func(t *testing.T, wrapper string, f ownedAssignment, last verbResult)
}

// followResetFace drives one reset fixture: the reset prints the face, the walk carries out
// the printed route, and the reset plan reruns out of the face. A refused apply reruns as
// its plan, because the plan prints the apply that the repaired checkout takes.
func followResetFace(t *testing.T, wrapper string, face refusalroute.Face, fixture resetRefusalFixture) {
	f, args := fixture.build(t)
	r := runVerb(t, verbReset, f.call(args...))
	record := "refused{detail=" + face.Sentence
	next, printed := recordField(r.stdout, record, refusalroute.NextField)
	if r.exit != 1 || !printed || next == "" {
		t.Fatalf("face %s = (%d, %q, %q), want exit 1 and a non-empty next= field", face.Name, r.exit, r.stdout, r.stderr)
	}
	carry := func(index int) (func(), bool) {
		step, carried := fixture.carry[index]
		return func() { step(t, f) }, carried
	}
	last := followRoute(t, face, printedSteps(t, face, next, "", nil), carry, func(step string) verbResult {
		return runPrintedStep(t, wrapper, f.repoHome, step)
	})
	if fixture.after != nil {
		fixture.after(t, wrapper, f, last)
	}
	if rerun := runVerb(t, verbReset, f.call(args[:3]...)); strings.Contains(rerun.stdout, record) {
		t.Fatalf("%s reset re-run = (%d, %q, %q), want it out of the face", face.Name, rerun.exit, rerun.stdout, rerun.stderr)
	}
}

// ownedResetFixture is the owned assignment of a reset fixture and the reset of its
// checkout to the assignment start.
func ownedResetFixture(t *testing.T, request string) (ownedAssignment, []string) {
	t.Helper()
	f := newOwnedAssignment(t, "reset-route-"+request)
	return f, resetToStart(f)
}

// resetToStart is the reset of the assignment's checkout to the assignment start.
func resetToStart(f ownedAssignment) []string {
	return []string{"--to", f.creation.Assignment.Start, f.creation.Assignment.ID}
}

func resetRefusalFixtures() []resetRefusalFixture {
	return append([]resetRefusalFixture{
		{
			face: faceResetCheckoutConflicted,
			build: func(t *testing.T) (ownedAssignment, []string) {
				f, args := ownedResetFixture(t, "conflicted")
				setupConflict(t, f.creation.Path)
				return f, args
			},
			// The printed clean plans the checkout's retirement, and the operator applies it.
			after: func(t *testing.T, wrapper string, f ownedAssignment, plan verbResult) {
				runPrintedStep(t, wrapper, f.repoHome, "bench worktree clean "+f.creation.Assignment.ID+" --apply "+plan.mustFingerprint(t))
			},
		},
		{
			// The checkout changes between the plan and its apply, so the apply's fingerprint
			// is stale. The printed plan prints the apply of the changed checkout.
			face: faceResetPlanStale,
			build: func(t *testing.T) (ownedAssignment, []string) {
				f, args := ownedResetFixture(t, "stale")
				mustWrite(t, filepath.Join(f.creation.Path, "work.txt"), []byte("planned\n"), 0o644)
				fingerprint := runVerb(t, verbReset, f.call(args...)).mustFingerprint(t)
				mustWrite(t, filepath.Join(f.creation.Path, "work.txt"), []byte("changed\n"), 0o644)
				return f, append(args, "--apply", fingerprint)
			},
			after: func(t *testing.T, wrapper string, f ownedAssignment, plan verbResult) {
				applyResetPlan(t, wrapper, f.repoHome, plan)
			},
		},
		{
			// A hidden index flag is a cause that the reset names no repair for, so it hands
			// back. The reviewer clears the flag, and the printed plan then runs.
			face: faceResetHandback,
			build: func(t *testing.T) (ownedAssignment, []string) {
				f, args := ownedResetFixture(t, "handback")
				gitRun(t, f.creation.Path, "update-index", "--assume-unchanged", "README.md")
				mustWrite(t, filepath.Join(f.creation.Path, "README.md"), []byte("hidden edit\n"), 0o644)
				return f, args
			},
			carry: map[int]func(*testing.T, ownedAssignment){0: func(t *testing.T, f ownedAssignment) {
				gitRun(t, f.creation.Path, "update-index", "--no-assume-unchanged", "README.md")
			}},
		},
	}, missingTreeResetFixtures()...)
}

// mergeTargetArgs are the arguments of a merge of incoming into the set's first
// assignment.
func mergeTargetArgs(f mergeSet, incoming string) []string {
	return []string{"--from", incoming, f.created[0].Assignment.ID}
}

// installedWrapper writes the wrapper an installed kit puts in front of the executable. It
// names a kit other than the fixture, so the fixture grades as a linked project, and it
// names itself as the wrapper, so a tree-scoped call's child runs behind it too.
//
// A child process writes the wrapper, so this process never holds the file open for write.
// A sibling test's fork inherits each descriptor this process holds until that child execs,
// and an exec of a file that a descriptor still holds for write fails as text file busy.
func installedWrapper(t *testing.T, binary string) string {
	t.Helper()
	wrapper := filepath.Join(t.TempDir(), "bench.sh")
	script := "#!/bin/sh\nBENCH_KIT=" + sanitize.ShellQuote(t.TempDir()) + " " + env.WrapperEnv + "=" + sanitize.ShellQuote(wrapper) +
		" exec " + sanitize.ShellQuote(binary) + " \"$@\"\n"
	write := descendant(t, "sh", "-c", `cat >"$1" && chmod 755 "$1"`, "sh", wrapper)
	write.Stdin = strings.NewReader(script)
	if out, err := write.CombinedOutput(); err != nil {
		t.Fatalf("write the installed wrapper: %v, %s", err, out)
	}
	return wrapper
}

// runPrintedStep runs one printed command step as the operator would paste it. A landing
// runs in process at the fixture's home, and any other Bench verb runs from the fixture's
// root checkout behind the installed wrapper. The step must read as one simple command with
// no placeholder left, and a Bench verb other than the landing must exit 0, unless the step
// is the diagnostic one.
func runPrintedStep(t *testing.T, wrapper string, f repoHome, step string) verbResult {
	t.Helper()
	stream := shellcommand.Parse(step)
	if stream.Unlexed || len(stream.Commands) != 1 || len(stream.Tokens) != stream.Commands[0].End-stream.Commands[0].Start {
		t.Fatalf("printed step %q is not one simple command with every slot filled", step)
	}
	words := shellcommand.ProjectCommandWords(stream.Tokens)
	if len(words) < 2 || words[0] != "bench" {
		t.Fatalf("printed step %q runs no Bench verb", step)
	}
	if len(words) > 3 && words[1] == "worktree" && words[2] == "land" {
		return runVerb(t, verbLand, f.call(words[3:]...))
	}
	var stdout, stderr bytes.Buffer
	cmd := descendant(t, wrapper, words[1:]...)
	cmd.Dir, cmd.Env, cmd.Stdout, cmd.Stderr = f.root, journeyChildEnv(t, f.home), &stdout, &stderr
	err := cmd.Run()
	r := verbResult{exit: exitCode(err), stdout: stdout.String(), stderr: stderr.String()}
	if r.exit != 0 && !diagnosticStep(words) {
		t.Fatalf("printed step %q = (%d, %q, %q), want exit 0; run error: %v", step, r.exit, r.stdout, r.stderr, err)
	}
	return r
}
