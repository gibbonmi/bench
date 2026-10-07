// First-run flow for the worktree landing command: grammars and stable-owner proofs.
package worktree

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/census"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/otelrecord"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/usage"
	"go.opentelemetry.io/otel/attribute"
)

var landGrammar = usage.Grammar{
	Cmd:     usage.CommandName(usage.WorktreeLand),
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
	Cmd:     usage.CommandName(usage.WorktreeLandResume),
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

// reviewedRangeDetail names the refusal a `--base` outside the assignment's reviewed
// range carries.
const reviewedRangeDetail = "review base is outside the assignment's reviewed range"

// LandCommand is the first-run reviewed-source landing operation. It performs every
// reversible proof before the exact-tree owner receives authority to publish. The
// invoked process is the one promotion owner for the complete landing: it never
// consults, rebuilds, or re-executes a repository executable, so candidate landing
// code cannot run during its own promotion.
func LandCommand(root, home string, args []string, stdout, stderr io.Writer) int {
	return landWith(defaultJoins(), newAmbient(home, stderr), root, args, stdout, stderr)
}

// landWith is LandCommand with the seam set and the ambient value resolved explicitly at
// the caller's boundary. It is also the landing's record boundary: one span covers the
// composition and the publication, and the resume path runs inside it, so a resumed
// landing records the same seam as the first run.
func landWith(j joins, a ambient, root string, args []string, stdout, stderr io.Writer) int {
	var measures landingMeasures
	ctx, finishSpan := beginLandingSpan(a.home, root)
	exit := landAttributed(ctx, &measures, j, a, root, args, stdout, stderr)
	finishSpan(exit, measures)
	return exit
}

// The landing seam's record. The span opens at the verb's own boundary, so it spans the
// composition and the publication as one seam rather than as two.
const otelLandingSeam = "worktree.land"

// landingMeasures is what the record states about one landing: the published subject,
// the size of the reviewed write set, and the assignment's census raw-call count. The
// census count rides here because the release step drops the records, so the record is
// where the count survives.
type landingMeasures struct {
	assignment         string
	subject            string
	pathCount          int
	counted            bool
	censusRawCalls     int
	censusRawCallsRead bool
}

// beginLandingSpan starts the landing's span and returns its context with the closer
// that ends it. The home is the one the verb boundary already resolved, and the context
// carries the span, so the authorization gate joins this trace rather than opening one.
func beginLandingSpan(home, root string) (context.Context, func(int, landingMeasures)) {
	ctx, span, finish := otelrecord.Begin(home, root, otelLandingSeam)
	return ctx, func(exit int, measures landingMeasures) {
		if measures.assignment != "" {
			span.SetAttributes(attribute.String(otelrecord.AttrAssignmentID, measures.assignment))
		}
		if measures.subject != "" {
			span.SetAttributes(attribute.String(otelrecord.AttrSubjectID, measures.subject))
		}
		// A landing refused before its preflight cleared reached neither measure, and a
		// zero would read as a landing that composed nothing and called nothing.
		if measures.counted {
			span.SetAttributes(attribute.String(otelrecord.AttrMeasurePathCount, strconv.Itoa(measures.pathCount)))
		}
		if measures.censusRawCallsRead {
			span.SetAttributes(attribute.String(otelrecord.AttrMeasureCensusRawCalls, strconv.Itoa(measures.censusRawCalls)))
		}
		span.SetAttributes(attribute.String(otelrecord.AttrOutcome, otelrecord.PublishedExitOutcome(exit)))
		finish()
	}
}

// landAttributed is the first-run landing itself, with the span's measures written to
// measures as each becomes known.
func landAttributed(ctx context.Context, measures *landingMeasures, j joins, a ambient, root string, args []string, stdout, stderr io.Writer) int {
	if hasResumeFlag(args) {
		return resumeLandWith(j, a, root, args, stdout, stderr)
	}
	parsed, line, code := usage.Parse(landGrammar, args)
	if line != "" {
		fmt.Fprintln(stderr, line)
		return code
	}
	operand := resolveVerbOperand(root, parsed.Positionals[0])
	path, err := canonicalPath(operand)
	if err != nil {
		rerun := landingRerun(parsed.Flags["--request"], parsed.Flags["--base"], parsed.Flags["--source-tip"], parsed.Flags["--spec"], operand, "")
		return landRefusalError(stdout, landingFaceRoute(errors.New("worktree path is not canonical"), rerun, false))
	}
	base := expandIdentity(root, parsed.Flags["--base"])
	tip := expandIdentity(root, parsed.Flags["--source-tip"])
	// Every reversible proof runs before the first refusal prints, so one preflight
	// names every refusal the caller must clear. The destination proofs and the
	// assignment proofs are independent; the source proofs need the assignment.
	//
	// The groups are the destination group and the assignment group, and the assignment
	// group runs the assignment proofs ahead of the source proofs that read their result.
	// A fault that stops a group's later proofs says so, and the operator then expects a
	// second refusal from the same group after the repair. The destination group holds one
	// stage, and the source stage is the assignment group's last, so neither states one.
	var refusals []error
	// The destination proof runs before the assignment resolves, so its route has no
	// assignment id to address and names the operator's own worktree path instead.
	unassignedRerun := landingRerun(parsed.Flags["--request"], base, tip, parsed.Flags["--spec"], path, "")
	destination, branch, priorMarker, destinationFingerprint, err := landingDestination(root)
	if err != nil {
		refusals = append(refusals, landingFaceRoute(err, unassignedRerun, false))
	}
	var source landingSourceFact
	assignment, err := landingAssignment(j, root, path, parsed.Flags["--request"], base, tip)
	if err != nil {
		refusals = append(refusals, landingFaceRoute(err, unassignedRerun, true))
	} else if source, err = landingSourceProofs(j, root, assignment, base, tip, parsed.Flags["--spec"], destination); err != nil {
		// The source proofs read the resolved assignment, so their route addresses it.
		refusals = append(refusals, landingSourceRoute(err, parsed.Flags["--request"], base, tip, parsed.Flags["--spec"], path, assignment.ID))
	}
	if len(refusals) > 0 {
		for _, err := range refusals {
			landRefusalError(stdout, err)
		}
		return 1
	}
	measures.assignment = assignment.ID
	// The count is read before the release step, because that step drops the records.
	// A landing that stops at an earlier step states the same count, and its resume
	// reads the file the release never removed.
	records := censusCount(a.home, root, assignment.ID)
	measures.censusRawCalls = records
	measures.censusRawCallsRead = true
	// The write-set size is the reviewed name-only diff, read from the same range
	// resolver the source proof above accepted. An unreadable range states no size.
	if reviewed, kind, _ := diff.ResolveSourceRange(assignment.Worktree, source.base, source.tip); kind == "" {
		measures.pathCount = len(reviewed.CommittedPaths)
		measures.counted = true
	}
	fmt.Fprintf(stderr, "landing source{review_base=%s,assignment_start=%s}\n", source.base, assignment.Start)
	printCensusHeads(stderr, a.home, root, assignment.ID)
	// The release removes the source worktree, so the broker check reads it now and the
	// notice prints after the effects report the refresh.
	brokerChanged := brokerSourceChanged(assignment.Worktree, source.base, source.tip)
	request := landing.ReviewedRequest{
		Root: root, Destination: "refs/heads/" + branch, DestinationBase: destination,
		Source: assignment.Branch, SourceTip: source.tip, ReviewBase: source.base,
		SourceWorktree: assignment.Worktree, SourceFingerprint: source.fingerprint, DestinationFingerprint: destinationFingerprint,
		SpecPath: source.specPath, SpecBytes: source.specBytes, SpecMode: source.specMode, ClosePath: source.closePath,
		Message: parsed.Flags["-m"], Stdout: stdout, Stderr: stderr,
	}
	result, err := j.landReviewed(ctx, request, commitmentAdmission{
		store:  commitrepo.Store{Root: root},
		source: commitrepo.Publication{Assignment: assignment.ID, Request: intent.RequestDigest(parsed.Flags["--request"]), Worktree: assignment.Worktree, Source: source.tip, Deliverable: request.Deliverable()},
		gap:    j.publicationGap,
	})
	if err != nil {
		var conflict landing.ConflictError
		if errors.As(err, &conflict) {
			return landRefusalError(stdout, landingConflictRefusal(conflict, destination, assignment.ID, parsed.Flags["--spec"], path, assignment.Worktree))
		}
		return landRefusalError(stdout, landingCompositionRoute(err, parsed.Flags["--request"], base, tip, parsed.Flags["--spec"], path, assignment))
	}
	// The destination CAS above is the commit point. Later errors name the durable
	// commit and retain the source. first-run never attempts to publish again.
	measures.subject = result.Commit
	if err := authorization.AdvanceMarker(context.Background(), root, branch, result.Commit, priorMarker); err != nil {
		return landedIncomplete(stdout, result, parsed.Flags["--spec"], path, assignment.ID, "marker", records)
	}
	if err := j.reconcileCommitment(root); err != nil {
		return landedIncomplete(stdout, result, parsed.Flags["--spec"], path, assignment.ID, "commitment", records)
	}
	if err := reconcileLandingDestination(j, root, result.Commit, result.Commit, result.DestinationBase); err != nil {
		return landedIncomplete(stdout, result, parsed.Flags["--spec"], path, assignment.ID, "reconcile", records)
	}
	if _, err := intent.PruneUnclaimedLandedBranches(root); err != nil {
		return landedIncomplete(stdout, result, parsed.Flags["--spec"], path, assignment.ID, "prune", records)
	}
	var releaseDiagnostic bytes.Buffer
	if release := j.releaseLandingAssignment(j, a, root, []string{"--request", parsed.Flags["--request"], path}, io.Discard, &releaseDiagnostic); release != 0 {
		if releaseDiagnostic.Len() > 0 {
			fmt.Fprintln(stderr, sanitize.Controls(strings.TrimSuffix(releaseDiagnostic.String(), "\n")))
		}
		return landedIncomplete(stdout, result, parsed.Flags["--spec"], path, assignment.ID, "release", records)
	}
	return landedAfterEffects(j, a, root, result, parsed.Flags["--spec"], path, assignment.ID, true, brokerChanged, records, stdout, stderr)
}

// landingCompositionRoute attaches a route to an error the composition and its
// authorization returned. The kind the gate attributed picks the face of an authorization
// refusal: the repair of a red commits in the source, so its re-run names the repaired tip,
// and an infrastructure outcome re-runs the caller's own command after the diagnosis. Every
// other error has no face of its own, so it hands back to the reviewer.
func landingCompositionRoute(err error, request, base, tip, specArg, path string, a intent.Assignment) error {
	rerun := landingRerun(request, base, tip, specArg, path, a.ID)
	var refused landing.AuthorizationRefusal
	if !errors.As(err, &refused) {
		return landingFaceRoute(err, rerun, false)
	}
	raised := refusal{detail: err.Error()}
	switch {
	case refused.Result.Kind == authorization.Infrastructure:
		return landingFaceRefusal(faceLandInfrastructure, raised, rerun, "")
	case landing.RedKind(refused.Result.Kind):
		raised.values = map[string]string{refusalroute.FactLabel: a.Label}
		return landingFaceRefusal(faceLandRed, raised, landingRerunAt(request, base, repairedSourceTipFlag, specArg, path, a.ID), "")
	}
	return landingFaceRoute(err, rerun, false)
}

// landingSourceProofs runs the first run's source proofs over a resolved assignment: the
// source itself, the review base's binding, and the destination's collisions with the
// source tree. The first fault stops the rest.
func landingSourceProofs(j joins, root string, assignment intent.Assignment, base, tip, specArg, destination string) (landingSourceFact, error) {
	source, err := landingSource(j, root, assignment, base, tip, specArg)
	if err != nil {
		return landingSourceFact{}, err
	}
	// The review base binds to the assignment's recorded start or to a descendant of it.
	// The destination advances while an assignment is open, so a landing rebases forward
	// and names the moved base; a base behind the recorded start instead grades a range the
	// assignment never authorized. `--is-ancestor` accepts the recorded start itself, which
	// is the unmoved case.
	if !git.OK("-C", root, "merge-base", "--is-ancestor", assignment.Start, source.base) {
		return landingSourceFact{}, identityRefusal(source.base, assignment.Start, reviewedRangeDetail)
	}
	// The review base binds before composition: a base outside the destination's history
	// grades a range the destination never reviewed against.
	if destination != "" && !git.OK("-C", root, "merge-base", "--is-ancestor", source.base, destination) {
		return landingSourceFact{}, identityRefusal(source.base, destination, landingBaseNotAncestorDetail)
	}
	if err := landingDestinationCollisions(root, source.tip); err != nil {
		return landingSourceFact{}, err
	}
	return source, nil
}

// commitmentAdmission is the installed broker's publication decision for one frozen
// source assignment. It reads authority from the destination repository, never from the
// candidate tree that it grades.
type commitmentAdmission struct {
	store  commitrepo.Store
	source commitrepo.Publication
	gap    func(string)
}

func (a commitmentAdmission) Check(tree string) error {
	if err := a.store.AdmitPublication(a.source, tree); err != nil {
		return fmt.Errorf("commitment: %w", err)
	}
	return nil
}

// Publish leaves a publication failure unwrapped, so compare-and-swap and recovery
// refusals keep their own text. Only an admission refusal carries the commitment prefix.
func (a commitmentAdmission) Publish(tree string, publish func() error) error {
	ran := false
	err := a.store.PublishAdmitted(a.source, tree, func() error {
		a.gap(a.store.Root)
		ran = true
		return publish()
	})
	if ran || err == nil {
		return err
	}
	return fmt.Errorf("commitment: %w", err)
}

// censusCount is the assignment's raw-call count for the landed record. An unreadable
// census reads as zero: the count is evidence beside the landing, never a condition
// on it.
func censusCount(home, root, assignment string) int {
	counts, _ := census.Counts(home, root)
	return counts[assignment]
}

// printCensusHeads states the raw-call count for each verb head the assignment used, and
// the response calls and bytes for each Bench verb head, beside the landing's other
// evidence. The lines print before the release step, which drops the records, so the
// retro reads the breakdown from the run instead of from memory. A line with no records
// behind it prints nothing.
func printCensusHeads(stderr io.Writer, home, root, assignment string) {
	if breakdown := census.HeadBreakdown(home, root, assignment); breakdown != "" {
		fmt.Fprintf(stderr, "census heads{%s}\n", breakdown)
	}
	if breakdown := census.OutputBreakdown(home, root, assignment); breakdown != "" {
		fmt.Fprintf(stderr, "census output{%s}\n", breakdown)
	}
}

func hasResumeFlag(args []string) bool {
	return usage.FlagPresent(landGrammar, args, "--resume")
}

// brokerSourceChanged reports whether the reviewed diff in worktree changes the promotion
// broker's own build inputs. An unresolvable input set reports false; the landing itself
// stays under the installed owner either way.
func brokerSourceChanged(worktree, base, tip string) bool {
	if !freshness.DeclaresBuildInputs(worktree) {
		return false
	}
	inputs, err := freshness.BuildInputs(worktree)
	if err != nil {
		return false
	}
	names, err := git.Output("-C", worktree, "diff", "--name-only", base, tip)
	if err != nil {
		return false
	}
	changed := map[string]struct{}{}
	for _, name := range strings.Split(names, "\n") {
		changed[name] = struct{}{}
	}
	for _, input := range inputs {
		if _, ok := changed[input]; ok {
			return true
		}
	}
	return false
}

// brokerChangeNotice names the install step that a broker-changing landing at root still
// owes after its refresh effect reported refresh, where kit is the kit value that the verb
// entry read. Source publication cannot replace the installed broker's authority.
//
// The kit's own source checkout carries no pin manifest, so 'bench repair' refuses there.
// Its route is the stamped rebuild and 'bench doctor --fix', and the refresh effect is
// that route, so the notice names it only after a failed refresh. Elsewhere the refresh
// republishes the destination's executable but not the installed broker, so the notice
// names 'bench repair' or the release install after every refresh result. The rebuild
// command comes from the one rebuild owner, never from a second copy here.
func brokerChangeNotice(kit, root, refresh string) string {
	step := "'bench repair' or the release install"
	if gate.KitSourceCheckoutAtKit(root, kit) {
		if refresh != effectFailed {
			return ""
		}
		step = freshness.RebuildAction(root) + " with 'bench doctor --fix'"
	}
	return "landing changes the promotion broker source; the installed broker keeps authority until " + step + " publishes the new broker"
}
