// Refusal route follow tests for the landing: each face's producing fixture carries out
// the route the landing printed, re-runs the landing, and the landing completes. The route
// walk helpers here serve the merge walk too.
package worktree

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/env"
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
