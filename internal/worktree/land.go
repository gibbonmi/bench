// First-run flow for the worktree landing command: grammars, injectable seams, and the rebuild-and-rerun path.
package worktree

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/runbinary"
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

var verifyLandingExecutable = freshness.Verify

var rebuildLandingExecutable = runbinary.Build

var reexecLanding = syscall.Exec

// rebuiltLandingEnv marks a landing process that a stale executable already rebuilt
// and re-ran once. A second stale verdict under it is the owner's refusal, never
// another rebuild, so a build that cannot reach freshness cannot loop.
const rebuiltLandingEnv = "BENCH_LANDING_REBUILT"

// LandCommand is the first-run reviewed-source landing operation. It performs every
// reversible proof before the exact-tree owner receives authority to publish.
func LandCommand(root, executable string, args []string, stdout, stderr io.Writer) int {
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
	// A stale executable enforces whatever landing contract it was built with, so its own
	// freshness is proven before any repository proof. A stale one is rebuilt through the
	// sanctioned build and the landing re-runs under the fresh binary; only a rebuild that
	// still fails the proof surfaces the owner's message, which carries the remedy.
	if freshness.DeclaresBuildInputs(root) {
		if err := verifyLandingExecutable(root, executable); err != nil {
			if landingAlreadyRebuilt() {
				return landRefusal(stdout, err.Error())
			}
			return rebuildAndRerunLanding(root, args, err, stdout, stderr)
		}
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
	}
	if len(refusals) > 0 {
		for _, err := range refusals {
			landRefusalError(stdout, err)
		}
		return 1
	}
	fmt.Fprintf(stderr, "landing source{review_base=%s,assignment_start=%s}\n", source.base, assignment.Start)
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

// rebuildAndRerunLanding republishes the repository's own Bench executable through the
// sanctioned build and replaces this process with the fresh binary running the same
// landing. The marker environment bounds this to one rebuild per landing attempt.
func rebuildAndRerunLanding(root string, args []string, cause error, stdout, stderr io.Writer) int {
	fresh := filepath.Join(root, "dist", "bench")
	if err := rebuildLandingExecutable(context.Background(), root, fresh); err != nil {
		return landRefusal(stdout, cause.Error()+"; rebuild failed: "+err.Error())
	}
	fmt.Fprintf(stderr, "landing executable was stale; rebuilt %s and re-ran the landing under it\n", fresh)
	argv := append([]string{fresh, "worktree", "land"}, args...)
	if err := reexecLanding(fresh, argv, append(os.Environ(), rebuiltLandingEnv+"=1")); err != nil {
		return landRefusal(stdout, cause.Error()+"; re-run failed: "+err.Error())
	}
	return landRefusal(stdout, "rebuilt landing executable did not take over the process")
}
