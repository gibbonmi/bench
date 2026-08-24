// First-run flow for the worktree landing command: grammars, injectable seams, and the stable-owner proofs.
package worktree

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/usage"
)

var landGrammar = usage.Grammar{
	Cmd:     "bench worktree land",
	Help:    "usage: " + usage.WorktreeLand,
	MinArgs: 1,
	MaxArgs: 1,
	Flags: []usage.Flag{
		{Name: "--resume", HasValue: true, NoEmptyValue: true},
		{Name: "--request", HasValue: true, NoEmptyValue: true, Required: true},
		{Name: "--base", HasValue: true, NoEmptyValue: true, Required: true},
		{Name: "--source-tip", HasValue: true, NoEmptyValue: true, Required: true},
		{Name: "--spec", HasValue: true, NoEmptyValue: true},
		{Name: "-m", HasValue: true, NoEmptyValue: true, Required: true},
	},
}

var resumeLandGrammar = usage.Grammar{
	Cmd:     "bench worktree land",
	Help:    "usage: " + usage.WorktreeLandResume,
	MinArgs: 1,
	MaxArgs: 1,
	Flags: []usage.Flag{
		{Name: "--resume", HasValue: true, NoEmptyValue: true, Required: true},
		{Name: "--request", HasValue: true, NoEmptyValue: true, Required: true},
		{Name: "--base", HasValue: true, NoEmptyValue: true, Required: true},
		{Name: "--source-tip", HasValue: true, NoEmptyValue: true, Required: true},
		{Name: "--spec", HasValue: true, NoEmptyValue: true},
	},
}

var landReviewed = func(ctx context.Context, request landing.ReviewedRequest) (landing.ReviewedResult, error) {
	return landing.New().LandReviewed(ctx, request)
}

var advanceLandingMarker = authorization.AdvanceMarker

var reconcileLanding = reconcileLandingDestination

var releaseLandingAssignment = ReleaseCommand

var authorizeLandingSource = preflight.AuthorizeReviewedSource

// rebuiltLandingEnv survives only for the retired rebuild guard in effects.go; the
// stable owner never rebuilds or re-executes a landing, so nothing sets it.
const rebuiltLandingEnv = "BENCH_LANDING_REBUILT"

// LandCommand is the first-run reviewed-source landing operation. It performs every
// reversible proof before the exact-tree owner receives authority to publish. The
// invoked process is the one promotion owner for the complete landing: it never
// consults, rebuilds, or re-executes a repository executable, so candidate landing
// code cannot run during its own promotion.
func LandCommand(root, _ string, args []string, stdout, stderr io.Writer) int {
	if hasResumeFlag(args) {
		return ResumeLandCommand(root, args, stdout, stderr)
	}
	parsed, line, code := usage.Parse(landGrammar, args)
	if line != "" {
		fmt.Fprintln(stderr, line)
		return code
	}
	path, err := canonicalPath(resolveVerbOperand(root, parsed.Positionals[0]))
	if err != nil {
		return landRefusal(stdout, "worktree path is not canonical")
	}
	base := expandIdentity(root, parsed.Flags["--base"])
	tip := expandIdentity(root, parsed.Flags["--source-tip"])
	// Every reversible proof runs before the first refusal prints, so one preflight
	// names every refusal the caller must clear. The destination proofs and the
	// assignment proofs are independent; the source proofs need the assignment.
	var refusals []error
	destination, branch, priorMarker, destinationFingerprint, err := landingDestination(root)
	if err != nil {
		refusals = append(refusals, err)
	}
	var source landingSourceFact
	assignment, err := landingAssignment(root, path, parsed.Flags["--request"], base, tip)
	if err != nil {
		refusals = append(refusals, err)
	} else if source, err = landingSource(root, assignment, base, tip, parsed.Flags["--spec"]); err != nil {
		refusals = append(refusals, err)
	} else if destination != "" && !git.OK("-C", root, "merge-base", "--is-ancestor", source.base, destination) {
		// The review base binds before composition: a base outside the destination's
		// history grades a range the destination never reviewed against.
		refusals = append(refusals, identityRefusal(source.base, destination, "review base is not an ancestor of the landing destination"))
	}
	if len(refusals) > 0 {
		for _, err := range refusals {
			landRefusalError(stdout, err)
		}
		return 1
	}
	fmt.Fprintf(stderr, "landing source{review_base=%s,assignment_start=%s}\n", source.base, assignment.Start)
	if notice := brokerChangeNotice(assignment.Worktree, source.base, source.tip); notice != "" {
		fmt.Fprintln(stderr, notice)
	}
	result, err := landReviewed(context.Background(), landing.ReviewedRequest{
		Root: root, Destination: "refs/heads/" + branch, DestinationBase: destination,
		Source: assignment.Branch, SourceTip: source.tip, ReviewBase: source.base,
		SourceWorktree: assignment.Worktree, SourceFingerprint: source.fingerprint, DestinationFingerprint: destinationFingerprint,
		SpecPath: source.specPath, SpecBytes: source.specBytes, SpecMode: source.specMode, ClosePath: source.closePath,
		Message: parsed.Flags["-m"], Stdout: stdout, Stderr: stderr,
	})
	if err != nil {
		var conflict landing.ConflictError
		if errors.As(err, &conflict) {
			return landRefusalError(stdout, refusalError{refusal{
				detail: conflict.Error(),
				paths:  conflict.Paths,
				next:   landingConflictNext(destination, assignment.ID, parsed.Flags["--spec"], path),
			}})
		}
		return landRefusal(stdout, err.Error())
	}
	// The destination CAS above is the commit point. Later errors name the durable
	// commit and retain the source. first-run never attempts to publish again.
	if err := advanceLandingMarker(context.Background(), root, branch, result.Commit, priorMarker); err != nil {
		return landedIncomplete(stdout, result, parsed.Flags["--spec"], path, assignment.ID, "marker")
	}
	if err := reconcileLanding(root, result.Commit, result.Commit, result.DestinationBase); err != nil {
		return landedIncomplete(stdout, result, parsed.Flags["--spec"], path, assignment.ID, "reconcile")
	}
	var releaseDiagnostic bytes.Buffer
	if release := releaseLandingAssignment(root, []string{"--request", parsed.Flags["--request"], path}, io.Discard, &releaseDiagnostic); release != 0 {
		if releaseDiagnostic.Len() > 0 {
			fmt.Fprintln(stderr, sanitize.Controls(strings.TrimSuffix(releaseDiagnostic.String(), "\n")))
		}
		return landedIncomplete(stdout, result, parsed.Flags["--spec"], path, assignment.ID, "release")
	}
	return landedComplete(stdout, result, true)
}

func hasResumeFlag(args []string) bool {
	return usage.FlagPresent(landGrammar, args, "--resume")
}

// brokerChangeNotice names the install step when the reviewed diff changes the
// promotion broker's own build inputs. Source publication cannot replace the broker's
// authority: the installed broker keeps landing until the release or repair path
// installs the new one. An unresolvable input set reports nothing; the landing itself
// stays under the installed owner either way.
func brokerChangeNotice(worktree, base, tip string) string {
	if !freshness.DeclaresBuildInputs(worktree) {
		return ""
	}
	inputs, err := freshness.BuildInputs(worktree)
	if err != nil {
		return ""
	}
	names, err := git.Output("-C", worktree, "diff", "--name-only", base, tip)
	if err != nil {
		return ""
	}
	changed := map[string]struct{}{}
	for _, name := range strings.Split(names, "\n") {
		changed[name] = struct{}{}
	}
	for _, input := range inputs {
		if _, ok := changed[input]; ok {
			return "landing changes the promotion broker source; the installed broker keeps authority until 'bench repair' or the release install publishes the new broker"
		}
	}
	return ""
}
