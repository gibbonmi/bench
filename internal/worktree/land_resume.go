// Resume flow for an interrupted landing: assignment recovery, publication facts, and residue policy.
package worktree

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/usage"
	"github.com/gibbonmi/bench/internal/worktree/landingpolicy"
)

// ResumeLandCommand finishes only a published landing's marker, destination checkout,
// and release work. Its proofs keep a retry from becoming a second publication attempt.
func ResumeLandCommand(root string, args []string, stdout, stderr io.Writer) int {
	parsed, line, code := usage.Parse(resumeLandGrammar, args)
	if line != "" {
		fmt.Fprintln(stderr, line)
		return code
	}
	path, err := canonicalPath(parsed.Positionals[0])
	if err != nil {
		return landRefusal(stdout, "worktree path is not canonical")
	}
	for _, flag := range []string{"--resume", "--base", "--source-tip"} {
		parsed.Flags[flag] = expandIdentity(root, parsed.Flags[flag])
	}
	destination, branch, marker, err := resumeLandingDestination(root)
	if err != nil {
		return landRefusal(stdout, err.Error())
	}
	published, sourceBase, destinationBase, tree, err := resumePublished(root, destination, parsed.Flags["--resume"], parsed.Flags["--base"], parsed.Flags["--source-tip"], parsed.Flags["--spec"])
	if err != nil {
		return landRefusalError(stdout, err)
	}
	result := landing.ReviewedResult{SourceBase: sourceBase, SourceTip: parsed.Flags["--source-tip"], DestinationBase: destinationBase, Commit: published, Tree: tree}
	assignment, active, err := resumeAssignment(root, path, parsed.Flags["--request"], parsed.Flags["--source-tip"], parsed.Flags["--base"], landingSlug(parsed.Flags["--spec"]))
	if err != nil {
		return landRefusalError(stdout, err)
	}
	assignmentID := assignment.ID
	if !active {
		receipt, err := terminalResumeReceipt(root, path, parsed.Flags["--request"], parsed.Flags["--source-tip"])
		if err != nil {
			return landRefusalError(stdout, err)
		}
		assignmentID = receipt.Tracked
	}
	if err := resumeDestructiveDestinationState(root, destination, published, destinationBase); err != nil {
		return landRefusal(stdout, err.Error())
	}
	switch landingpolicy.ResumeMarker(resumeMarkerFacts(root, destination, published, marker)) {
	case landingpolicy.MarkerAdvance:
		if err := advanceLandingMarker(context.Background(), root, branch, published, marker); err != nil {
			return landedIncomplete(stdout, result, parsed.Flags["--spec"], path, assignmentID, "marker")
		}
	case landingpolicy.MarkerRefuse:
		return landRefusal(stdout, landingpolicy.MarkerRefusalDetail)
	}
	if err := reconcileLanding(root, destination, published, destinationBase); err != nil {
		return landedIncomplete(stdout, result, parsed.Flags["--spec"], path, assignmentID, "reconcile")
	}
	if !active {
		return landedComplete(stdout, result, false)
	}
	if releaseLandingAssignment(root, []string{"--request", parsed.Flags["--request"], assignment.Worktree}, io.Discard, stderr) != 0 {
		return landedIncomplete(stdout, result, parsed.Flags["--spec"], path, assignmentID, "release")
	}
	return landedComplete(stdout, result, true)
}

func terminalResumeReceipt(root, path, request, sourceTip string) (intent.CleanupReceipt, error) {
	repo, target, err := cleanupIdentity(root, path)
	if err != nil {
		return intent.CleanupReceipt{}, errors.New("missing-terminal-receipt")
	}
	receipt, found, err := intent.CleanupReceiptFor(root, repo, releaseOperation, target, intent.RequestDigest(request))
	if err != nil || !found || receipt.State != intent.ReceiptComplete || receipt.Phase != intent.ReceiptPhaseTerminal || !receipt.Owned || receipt.Action != string(ActionRemoved) || !intent.ValidIdentity(receipt.Tracked) || receipt.Branch == "" || !strings.HasSuffix(receipt.Branch, "/"+receipt.Tracked) {
		return intent.CleanupReceipt{}, errors.New("missing-terminal-receipt")
	}
	if receipt.BranchOID != sourceTip {
		return intent.CleanupReceipt{}, identityRefusal(sourceTip, receipt.BranchOID, "terminal receipt source tip mismatch")
	}
	return receipt, nil
}

func resumeLandingDestination(root string) (string, string, string, error) {
	branch, ok := git.ResolvedDefault(root)
	if !ok {
		return "", "", "", errors.New("default branch is unresolved")
	}
	current, err := git.CheckedOutBranch(root)
	if err != nil || current != branch {
		return "", "", "", errors.New("landing checkout is not attached to the default branch")
	}
	destination, err := git.Output("-C", root, "rev-parse", "refs/heads/"+branch+"^{commit}")
	if err != nil {
		return "", "", "", errors.New("landing destination has no commit")
	}
	marker, err := landingMarker(root, branch, destination)
	if err != nil {
		return "", "", "", err
	}
	return destination, branch, marker, nil
}

func resumeDestructiveDestinationState(root, destination, published, destinationBase string) error {
	if detail := landingpolicy.Residue(destinationResidueFacts(root, destination, published, destinationBase)); detail != "" {
		return errors.New(detail)
	}
	return nil
}

// destinationResidueFacts translates the destination's Git and filesystem state
// into the typed residue facts once at the boundary. The expensive allowance
// and staged-content facts bind as lazy suppliers the policy consults on demand.
func destinationResidueFacts(root, destination, published, destinationBase string) landingpolicy.ResidueFacts {
	facts := landingpolicy.ResidueFacts{
		DestinationAtPublished: destination == published,
		IgnoredDeclared:        func() bool { return ignoredResidueDeclared(root) },
		StagedMatchesPublished: func() bool { return git.OK("-C", root, "diff", "--cached", "--quiet", destinationBase, "--") },
	}
	nested, err := classifyNestedState(root)
	facts.NestedClean = err == nil && nested == nestedClean
	if !facts.NestedClean {
		return facts
	}
	raw, err := git.Raw("-C", root, "status", "--porcelain=v1", "-z", "--no-renames", "--untracked-files=all", "--ignored")
	facts.StatusReadable = err == nil
	if !facts.StatusReadable {
		return facts
	}
	entries, err := git.ParsePorcelainZStrict(raw)
	facts.StatusWellFormed = err == nil
	for _, entry := range entries {
		facts.Entries = append(facts.Entries, landingpolicy.StatusEntry{Status: entry.Status, Path: entry.Path})
	}
	return facts
}

// ignoredResidueDeclared reports whether the destination's ignored residue sits
// entirely inside its declared build outputs. Resume applies the same allowance the
// first run did. A landing that validly proceeded with declared outputs and then
// failed at release can still be completed, rather than being permanently stuck.
func ignoredResidueDeclared(root string) bool {
	ignored, _, ignoredErr := inventoryIgnored(root, false)
	declared, _, declarationErr := loadBuildOutputs(root)
	if ignoredErr != nil || declarationErr != nil {
		return false
	}
	return ignoredWithinLandingAllowance(ignored, declared)
}

func resumeAssignment(root, path, request, tip, base, slug string) (intent.Assignment, bool, error) {
	a, found, err := intent.FindAssignmentForRequest(root, request)
	if err != nil {
		return intent.Assignment{}, false, err
	}
	if !found {
		recovery, recoverable, err := unmatchedRequestRecovery(root, assignmentRecoveryContext{
			target: path,
			base:   base,
			tip:    tip,
		})
		if err != nil {
			return intent.Assignment{}, false, err
		}
		if recoverable {
			return intent.Assignment{}, false, refusalError{recovery}
		}
		return intent.Assignment{}, false, nil
	}
	if err := identityBundleRefusal(root, path, a, resumeActiveState); err != nil {
		return intent.Assignment{}, false, err
	}
	if _, err := landingSource(root, a, base, tip, slug); err != nil {
		return intent.Assignment{}, false, err
	}
	return a, true, nil
}

func resumePublished(root, destination, value, base, source, slug string) (published, sourceBase, destinationBase, tree string, err error) {
	facts, sourceBase, destinationBase := publicationFacts(root, destination, value, base, source, slug)
	if refusal := landingpolicy.Publication(facts); refusal.Detail != "" {
		if refusal.Observed != "" || refusal.Wanted != "" {
			return "", "", "", "", identityRefusal(refusal.Observed, refusal.Wanted, refusal.Detail)
		}
		return "", "", "", "", errors.New(refusal.Detail)
	}
	tree, err = git.Output("-C", root, "rev-parse", facts.Published+"^{tree}")
	if err != nil {
		return "", "", "", "", errors.New("published tree is unreadable")
	}
	return facts.Published, sourceBase, destinationBase, tree, nil
}

// publicationFacts translates the repository's view of a claimed published
// landing into the typed publication facts once at the boundary. Gathering
// stops at the first fact the policy will refuse on, so no later repository
// read depends on state an earlier failure left unresolved. It also returns
// the resolved review-source base and the published commit's destination base.
func publicationFacts(root, destination, value, base, source, slug string) (landingpolicy.PublicationFacts, string, string) {
	facts := landingpolicy.PublicationFacts{Requested: value, RequestedSource: source}
	published, err := git.Output("-C", root, "rev-parse", "--verify", value+"^{commit}")
	if err != nil {
		return facts, "", ""
	}
	facts.Resolved, facts.Published = true, published
	if published != value {
		return facts, "", ""
	}
	facts.ReachableFromDestination = git.OK("-C", root, "merge-base", "--is-ancestor", published, destination)
	if !facts.ReachableFromDestination {
		return facts, "", ""
	}
	parents, err := git.Output("-C", root, "rev-list", "--parents", "-n", "1", published)
	if err != nil {
		return facts, "", ""
	}
	facts.ParentsReadable, facts.Parents = true, strings.Fields(parents)
	if len(facts.Parents) != 3 || facts.Parents[2] != source {
		return facts, "", ""
	}
	destinationBase := facts.Parents[1]
	sourceRange, kind, _ := diff.ResolveSourceRange(root, base, source)
	facts.RangeAuthenticates = kind == ""
	if !facts.RangeAuthenticates {
		return facts, "", destinationBase
	}
	facts.Spec = specTransitionFacts(root, published, source, slug)
	return facts, sourceRange.Base, destinationBase
}

// specTransitionFacts translates what the source and published trees carry for
// the landing's named spec. A spec-less first run published no transition, so
// a resume without a spec has no transition to authenticate. A source folder
// with no spec.md published a close, not a transition; its authentication is
// the folder's absence from the published tree.
func specTransitionFacts(root, published, source, slug string) landingpolicy.SpecTransition {
	transition := landingpolicy.SpecTransition{Named: slug != ""}
	if !transition.Named {
		return transition
	}
	specPath := spec.LiveSpecPath(slug)
	folder := landing.ClosedFolderPath(landingSlug(slug))
	staged, stagedErr := git.Raw("-C", root, "show", source+":"+specPath)
	if stagedErr != nil && git.OK("-C", root, "cat-file", "-e", source+":"+folder) {
		transition.TicketsOnlyClose = true
		transition.PublishedHasFolder = git.OK("-C", root, "cat-file", "-e", published+":"+folder)
		return transition
	}
	implemented, implementedErr := spec.Implemented(staged)
	publishedSpec, publishedErr := git.Raw("-C", root, "show", published+":"+specPath)
	transition.TransitionMatches = stagedErr == nil && implementedErr == nil && publishedErr == nil && bytes.Equal(implemented, publishedSpec)
	return transition
}

// resumeMarkerFacts translates the destination, published commit, and recorded
// green marker into the typed resume-marker facts once at the boundary.
func resumeMarkerFacts(root, destination, published, marker string) landingpolicy.MarkerFacts {
	return landingpolicy.MarkerFacts{
		DestinationAtPublished: destination == published,
		MarkerPresent:          marker != "",
		MarkerReachesPublished: func() bool { return git.OK("-C", root, "merge-base", "--is-ancestor", published, marker) },
	}
}
