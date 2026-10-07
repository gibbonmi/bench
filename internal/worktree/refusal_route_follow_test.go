// Refusal route follow tests for the landing: each face's producing fixture carries out
// the route the landing printed, re-runs the landing, and the face no longer prints.
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

// routeStepJoiner is the joiner the spec pins between two printed route steps.
const routeStepJoiner = "; then "

// TestLandingFacesFollowTheirRoutes is RR15, with RR18 and RR19 folded in. Each producing
// fixture's printed route is carried out step by step: the fixture carries out the steps it
// declares by its own means, and every other step, each agent command step included, runs
// verbatim with only its operator slots filled. The last step re-runs the landing, and the
// face's refused record no longer prints. A reviewer face's route opens with the reviewer
// marker and an agent face's does not.
func TestLandingFacesFollowTheirRoutes(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	faces := map[string]refusalroute.Face{}
	for _, face := range refusalroute.Faces(refusalroute.Land) {
		faces[face.Name] = face
	}
	for _, fixture := range landingRefusalFixtures() {
		t.Run(fixture.face, func(t *testing.T) {
			t.Parallel()
			face := faces[fixture.face]
			produced := produceLandingFace(t, fixture)
			route, reviewer := strings.CutPrefix(produced.next, reviewerRoute)
			if reviewer != (face.Authority == refusalroute.Reviewer) {
				t.Fatalf("%s next = %q, want the reviewer marker exactly when the face is the reviewer's", fixture.face, produced.next)
			}
			for _, name := range fixtureRouteNames(fixture, produced.f) {
				if !strings.Contains(route, name) {
					t.Fatalf("%s next = %q, want it to name %q", fixture.face, produced.next, name)
				}
			}
			route = strings.TrimPrefix(route, laterProofsSkipped+"; ")
			steps := strings.Split(route, routeStepJoiner)
			if len(steps) != len(face.Route) {
				t.Fatalf("%s next = %q prints %d steps, want the %d the face declares", fixture.face, produced.next, len(steps), len(face.Route))
			}
			detail := refusedDetail(produced.r.stdout, refusalroute.Sentence(fixture.face))
			fill := operatorFill(produced)
			var last verbResult
			for index, step := range face.Route {
				carry, carried := fixture.carry[index]
				switch {
				case carried && face.Authority == refusalroute.Agent && step.IsCommand():
					t.Fatalf("%s step %d is an agent command, which the walk runs verbatim; the fixture may not carry it", fixture.face, index+1)
				case carried:
					carry(t, produced.f)
				case step.IsCommand():
					last = runPrintedStep(t, binary, produced.f, fill(t, steps[index]))
				default:
					t.Fatalf("%s step %d %q is an instruction that the fixture does not carry out", fixture.face, index+1, steps[index])
				}
			}
			if fixture.stage == stageIncomplete {
				if last.exit != 0 || strings.Contains(last.stdout, "worktree=incomplete") {
					t.Fatalf("%s resume = (%d, %q, %q), want the landing complete", fixture.face, last.exit, last.stdout, last.stderr)
				}
				return
			}
			if strings.Contains(last.stdout, "refused{detail="+detail) {
				t.Fatalf("%s re-run = (%d, %q, %q), want the face's refusal %q gone", fixture.face, last.exit, last.stdout, last.stderr, detail)
			}
		})
	}
}

// fixtureRouteNames is what a fixture's printed route must name, read from the fixture.
func fixtureRouteNames(fixture landingRefusalFixture, f landingFixture) []string {
	if fixture.names == nil {
		return nil
	}
	return fixture.names(f)
}

// refusedDetail is the sentence that marks a face's refused record: the face's declared
// sentence, or, for a face whose sentence a policy owns, the detail the record printed. That
// detail runs to the first field after it.
func refusedDetail(stdout, sentence string) string {
	if sentence != "" {
		return sentence
	}
	for _, line := range strings.Split(stdout, "\n") {
		detail, found := strings.CutPrefix(line, "refused{detail=")
		if !found {
			continue
		}
		for _, field := range []string{",observed=", ",wanted=", ",next=", "}"} {
			detail, _, _ = strings.Cut(detail, field)
		}
		return detail
	}
	return ""
}

// operatorFill fills the slots a printed route leaves to the operator, with the values the
// operator holds: the paths it changed in the source, its own request and message, and the
// source tip after its repair. A slot the walk has no value for stays, and the step then
// refuses to run.
func operatorFill(produced producedFace) func(t *testing.T, step string) string {
	return func(t *testing.T, step string) string {
		t.Helper()
		source := produced.f.creation.Path
		var paths []string
		for _, path := range strings.Fields(gitOutput(t, source, "ls-files", "--modified", "--others", "--exclude-standard")) {
			paths = append(paths, sanitize.ShellQuote(path))
		}
		tip := gitOutput(t, source, "rev-parse", "HEAD")
		return strings.NewReplacer(
			"<msg>", "'follow the refusal route'",
			"<path>...", strings.Join(paths, " "),
			"<message>", "'land the followed route'",
			"<request>", sanitize.ShellQuote(produced.request),
			"<repaired-source-tip>", sanitize.ShellQuote(tip),
		).Replace(step)
	}
}

// installedWrapper writes the wrapper an installed kit puts in front of the executable. It
// names a kit other than the fixture, so the fixture grades as a linked project, and it
// names itself as the wrapper, so a tree-scoped call's child runs behind it too.
func installedWrapper(t *testing.T, binary string) string {
	t.Helper()
	wrapper := filepath.Join(t.TempDir(), "bench.sh")
	script := "#!/bin/sh\nBENCH_KIT=" + sanitize.ShellQuote(t.TempDir()) + " " + env.WrapperEnv + "=" + sanitize.ShellQuote(wrapper) +
		" exec " + sanitize.ShellQuote(binary) + " \"$@\"\n"
	mustWrite(t, wrapper, []byte(script), 0o755)
	return wrapper
}

// runPrintedStep runs one printed command step as the operator would paste it. A landing
// runs in process at the fixture's home, and any other Bench verb runs through the sealed
// test-run executable from the destination checkout, behind a wrapper as an installed kit
// runs it. The step must read as one simple command with no placeholder left, and a Bench
// verb other than the landing must exit 0.
func runPrintedStep(t *testing.T, binary string, f landingFixture, step string) verbResult {
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
	cmd := descendant(t, installedWrapper(t, binary), words[1:]...)
	cmd.Dir, cmd.Env, cmd.Stdout, cmd.Stderr = f.root, journeyChildEnv(t, f.home), &stdout, &stderr
	r := verbResult{exit: exitCode(cmd.Run()), stdout: stdout.String(), stderr: stderr.String()}
	if r.exit != 0 {
		t.Fatalf("printed step %q = (%d, %q, %q), want exit 0", step, r.exit, r.stdout, r.stderr)
	}
	return r
}
