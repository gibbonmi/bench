// First-run flow for the worktree landing command: grammars, injectable seams, and the stable-owner proofs.
package worktree

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gibbonmi/bench/internal/census"
	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
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

// joins is this package's injectable seam set. A verb resolves one value at its boundary
// with defaultJoins and passes it down; nothing below the boundary reads a package
// variable. One value carries every family because the families travel together: a
// landing releases its assignment, a release cleans the worktree, and that cleanup reads
// the ignored inventory and the live-binary identity. A test builds the value, replaces
// one field, and calls the internal form, so two tests never share one stub point.
//
// A field whose default has to reach another seam takes the value as its first argument.
// The default then reads the caller's joins rather than a captured copy.
type joins struct {
	landReviewed             func(context.Context, landing.ReviewedRequest) (landing.ReviewedResult, error)
	advanceLandingMarker     func(context.Context, string, string, string, string) error
	reconcileLanding         func(joins, string, string, string, string) error
	releaseLandingAssignment func(joins, string, string, []string, io.Writer, io.Writer) int
	authorizeLandingSource   func(string, string, string) (diff.SourceRange, error)
	// cleanupBoundary is the deterministic transaction fault seam. A nil value carries no
	// fault, exactly as hit reads it, so it is also the default.
	cleanupBoundary      Fault
	cleanupLockAttempt   func(string)
	creationLockAttempt  func(string)
	claimTakeoverGap     func(string)
	claimStealGap        func(string)
	restoreClean         func(string)
	chmodPool            func(string, os.FileMode) error
	ignoredLstat         func(string) (os.FileInfo, error)
	resolveRunningBinary func() (string, error)
	liveBinaryWarnings   io.Writer
	// planLandedExplicit exposes the target-path boundary so a hostile-path test can stand
	// in for the planner without the fixture the real one needs.
	planLandedExplicit   func(joins, string, string, CleanupOptions) (CleanupPlan, error)
	reauthorizeUnlock    func(string, string) error
	reauthorizeLock      func(string, string, string) error
	reauthorizeBeforeCAS func(*intent.Assignment)
	// mergeLane resolves the fast lane the merge's composed tree is graded under. It is a
	// seam because the resolution reads the kit root out of the process environment, and a
	// fixture that bound that environment would leave the package's parallel set.
	mergeLane func(string) ([]gate.Phase, string, error)
	// mergeReconcile is the merge verb's publication boundary: the checkout catch-up that
	// runs after the branch ref moved. It is a seam because its failure is the one
	// outcome that reads apart from a refusal, and no fixture can make a bare reset fail.
	mergeReconcile func(string, string) error
	// home is the Bench home the verb's own boundary resolved. The retirement path
	// needs it to drop the retired assignment's census records, and it travels in the
	// seam set because every verb that retires an assignment already carries the set
	// down to that path. The default names the operator's home, so an in-package entry
	// point that resolves no home of its own still drops the right records.
	home string
}

// defaultJoins names the real function behind every seam. It is the one place a default
// lives, so a verb's boundary and a test's starting value cannot disagree.
func defaultJoins() joins {
	return joins{
		landReviewed: func(ctx context.Context, request landing.ReviewedRequest) (landing.ReviewedResult, error) {
			return landing.New().LandReviewed(ctx, request)
		},
		advanceLandingMarker:     authorization.AdvanceMarker,
		reconcileLanding:         reconcileLandingDestination,
		releaseLandingAssignment: releaseCommandWith,
		authorizeLandingSource:   preflight.AuthorizeReviewedSource,
		cleanupLockAttempt:       func(string) {},
		creationLockAttempt:      func(string) {},
		claimTakeoverGap:         func(string) {},
		claimStealGap:            func(string) {},
		restoreClean:             restoreCleanCheckout,
		chmodPool:                os.Chmod,
		ignoredLstat:             os.Lstat,
		resolveRunningBinary:     os.Executable,
		liveBinaryWarnings:       os.Stderr,
		planLandedExplicit:       planExplicitWith,
		home:                     Home(),
		reauthorizeUnlock:        unlockWorktree,
		reauthorizeLock:          lockWorktree,
		mergeLane:                gate.LaneForCommit,
		mergeReconcile:           reconcileMergeCheckout,
	}
}

// reviewedRangeDetail names the refusal a `--base` outside the assignment's reviewed
// range carries.
const reviewedRangeDetail = "review base is outside the assignment's reviewed range"

// rebuiltLandingEnv survives only for the retired rebuild guard in effects.go; the
// stable owner never rebuilds or re-executes a landing, so nothing sets it.
const rebuiltLandingEnv = "BENCH_LANDING_REBUILT"

// LandCommand is the first-run reviewed-source landing operation. It performs every
// reversible proof before the exact-tree owner receives authority to publish. The
// invoked process is the one promotion owner for the complete landing: it never
// consults, rebuilds, or re-executes a repository executable, so candidate landing
// code cannot run during its own promotion.
func LandCommand(root, home, executable string, args []string, stdout, stderr io.Writer) int {
	return landWith(defaultJoins(), root, home, executable, args, stdout, stderr)
}

// landWith is LandCommand with the seam set resolved explicitly at the caller's boundary.
func landWith(j joins, root, home, _ string, args []string, stdout, stderr io.Writer) int {
	j.home = home
	if hasResumeFlag(args) {
		return resumeLandWith(j, root, home, args, stdout, stderr)
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
	destination, branch, priorMarker, destinationFingerprint, err := landingDestination(j, root)
	if err != nil {
		refusals = append(refusals, err)
	}
	var source landingSourceFact
	assignment, err := landingAssignment(j, root, path, parsed.Flags["--request"], base, tip)
	if err != nil {
		refusals = append(refusals, err)
	} else if source, err = landingSource(j, root, assignment, base, tip, parsed.Flags["--spec"]); err != nil {
		refusals = append(refusals, err)
	} else if !git.OK("-C", root, "merge-base", "--is-ancestor", assignment.Start, source.base) {
		// The review base binds to the assignment's recorded start or to a descendant
		// of it. The destination advances while an assignment is open, so a landing
		// rebases forward and names the moved base; a base behind the recorded start
		// instead grades a range the assignment never authorized. `--is-ancestor`
		// accepts the recorded start itself, which is the unmoved case.
		refusals = append(refusals, identityRefusal(source.base, assignment.Start, reviewedRangeDetail))
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
	// The count is read before the release step, because that step drops the records.
	// A landing that stops at an earlier step states the same count, and its resume
	// reads the file the release never removed.
	records := censusCount(home, root, assignment.ID)
	fmt.Fprintf(stderr, "landing source{review_base=%s,assignment_start=%s}\n", source.base, assignment.Start)
	if notice := brokerChangeNotice(assignment.Worktree, source.base, source.tip); notice != "" {
		fmt.Fprintln(stderr, notice)
	}
	result, err := j.landReviewed(context.Background(), landing.ReviewedRequest{
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
	if err := j.advanceLandingMarker(context.Background(), root, branch, result.Commit, priorMarker); err != nil {
		return landedIncomplete(stdout, result, parsed.Flags["--spec"], path, assignment.ID, "marker", records)
	}
	if err := j.reconcileLanding(j, root, result.Commit, result.Commit, result.DestinationBase); err != nil {
		return landedIncomplete(stdout, result, parsed.Flags["--spec"], path, assignment.ID, "reconcile", records)
	}
	var releaseDiagnostic bytes.Buffer
	if release := j.releaseLandingAssignment(j, root, home, []string{"--request", parsed.Flags["--request"], path}, io.Discard, &releaseDiagnostic); release != 0 {
		if releaseDiagnostic.Len() > 0 {
			fmt.Fprintln(stderr, sanitize.Controls(strings.TrimSuffix(releaseDiagnostic.String(), "\n")))
		}
		return landedIncomplete(stdout, result, parsed.Flags["--spec"], path, assignment.ID, "release", records)
	}
	return landedComplete(stdout, result, true, records)
}

// censusCount is the assignment's raw-call count for the landed record. An unreadable
// census reads as zero: the count is evidence beside the landing, never a condition
// on it.
func censusCount(home, root, assignment string) int {
	counts, _ := census.Counts(home, root)
	return counts[assignment]
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
