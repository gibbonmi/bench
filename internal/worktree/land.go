package worktree

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/gate/greenmarker"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/spec"
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
	path, err := canonicalPath(parsed.Positionals[0])
	if err != nil {
		return landRefusal(stdout, "worktree path is not canonical")
	}
	// A stale executable enforces whatever landing contract it was built with, so its own
	// freshness is proven before any repository proof. A stale one is rebuilt through the
	// sanctioned build and the landing re-runs under the fresh binary; only a rebuild that
	// still fails the proof surfaces the owner's message, which carries the remedy.
	if freshness.DeclaresBuildInputs(root) {
		if err := verifyLandingExecutable(root, executable); err != nil {
			if os.Getenv(rebuiltLandingEnv) != "" {
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
		SpecPath: source.specPath, SpecBytes: source.specBytes, SpecMode: source.specMode,
		Message: parsed.Flags["-m"], Stdout: stdout, Stderr: stderr,
	})
	if err != nil {
		var conflict landing.ConflictError
		if errors.As(err, &conflict) {
			return landRefusalError(stdout, refusalError{refusal{detail: conflict.Error(), paths: conflict.Paths}})
		}
		return landRefusal(stdout, err.Error())
	}
	// The destination CAS above is the commit point. Later errors name the durable
	// commit and retain the source. first-run never attempts to publish again.
	if err := advanceLandingMarker(context.Background(), root, branch, result.Commit, priorMarker); err != nil {
		return landedIncomplete(stdout, result, parsed.Flags["--spec"], path, assignment.ID, "marker")
	}
	if err := reconcileLanding(root, result.Commit); err != nil {
		return landedIncomplete(stdout, result, parsed.Flags["--spec"], path, assignment.ID, "reconcile")
	}
	var releaseDiagnostic bytes.Buffer
	if release := releaseLandingAssignment(root, []string{"--request", parsed.Flags["--request"], path}, io.Discard, &releaseDiagnostic); release != 0 {
		if releaseDiagnostic.Len() > 0 {
			fmt.Fprintln(stderr, sanitize.Controls(strings.TrimSuffix(releaseDiagnostic.String(), "\n")))
		}
		return landedIncomplete(stdout, result, parsed.Flags["--spec"], path, assignment.ID, "release")
	}
	fmt.Fprintf(stdout, "landed{source_base=%s,source_tip=%s,destination_base=%s,published_commit=%s,tree=%s,worktree=released}\n", result.SourceBase, result.SourceTip, result.DestinationBase, result.Commit, result.Tree)
	return 0
}

func hasResumeFlag(args []string) bool {
	return usage.FlagPresent(landGrammar, args, "--resume")
}

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
	if destination == published {
		if err := advanceLandingMarker(context.Background(), root, branch, published, marker); err != nil {
			return landedIncomplete(stdout, result, parsed.Flags["--spec"], path, assignmentID, "marker")
		}
	} else if marker == "" || !git.OK("-C", root, "merge-base", "--is-ancestor", published, marker) {
		return landRefusal(stdout, "project-green marker is absent, behind, or divergent from the published landing")
	}
	if err := reconcileLanding(root, destination); err != nil {
		return landedIncomplete(stdout, result, parsed.Flags["--spec"], path, assignmentID, "reconcile")
	}
	if !active {
		fmt.Fprintf(stdout, "landed{source_base=%s,source_tip=%s,destination_base=%s,published_commit=%s,tree=%s,worktree=already-complete}\n", sourceBase, parsed.Flags["--source-tip"], destinationBase, published, tree)
		return 0
	}
	if releaseLandingAssignment(root, []string{"--request", parsed.Flags["--request"], assignment.Worktree}, io.Discard, stderr) != 0 {
		return landedIncomplete(stdout, result, parsed.Flags["--spec"], path, assignmentID, "release")
	}
	fmt.Fprintf(stdout, "landed{source_base=%s,source_tip=%s,destination_base=%s,published_commit=%s,tree=%s,worktree=released}\n", sourceBase, parsed.Flags["--source-tip"], destinationBase, published, tree)
	return 0
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
	nested, err := classifyNestedState(root)
	if err != nil || nested != nestedClean {
		return errors.New("landing destination has nested repositories")
	}
	raw, err := git.Raw("-C", root, "status", "--porcelain=v1", "-z", "--no-renames", "--untracked-files=all", "--ignored")
	if err != nil {
		return errors.New("landing destination status is unreadable")
	}
	entries, err := git.ParsePorcelainZStrict(raw)
	if err != nil {
		return errors.New("landing destination status is malformed")
	}
	staged := false
	allowedIgnored, allowanceKnown := false, false
	for _, entry := range entries {
		switch entry.Status {
		case "":
			continue
		case "!!":
			if !allowanceKnown {
				allowedIgnored, allowanceKnown = ignoredResidueDeclared(root), true
			}
			if allowedIgnored {
				continue
			}
			return errors.New("landing destination has ignored residue")
		case "??":
			return errors.New("landing destination has untracked collisions")
		}
		if entry.Status[1] != ' ' {
			return errors.New("landing destination has tracked-worktree changes")
		}
		staged = staged || entry.Status[0] != ' '
	}
	if staged && (destination != published || !git.OK("-C", root, "diff", "--cached", "--quiet", destinationBase, "--")) {
		return errors.New("landing destination has staged changes")
	}
	return nil
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
			detail: assignmentMismatchDetail,
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
	if (a.State != intent.StateActive && a.State != intent.StateCleanupPending) || a.Worktree != path {
		return intent.Assignment{}, false, errors.New(assignmentMismatchDetail)
	}
	evidence, markerErr := validateOwnerMarker(root, path)
	if markerErr != nil || evidence.marker.OwnerID != a.OwnerID || evidence.marker.Path != a.Worktree || evidence.registration.BranchRef != a.Branch || !evidence.registration.Locked || evidence.registration.LockReason != lockReason(a) {
		return intent.Assignment{}, false, errors.New("owner marker or assignment branch mismatch")
	}
	if _, err := landingSource(root, a, base, tip, slug); err != nil {
		return intent.Assignment{}, false, err
	}
	return a, true, nil
}

func resumePublished(root, destination, value, base, source, slug string) (published, sourceBase, destinationBase, tree string, err error) {
	published, err = git.Output("-C", root, "rev-parse", "--verify", value+"^{commit}")
	if err != nil {
		return "", "", "", "", errors.New("published commit is not an exact commit identity")
	}
	if published != value {
		return "", "", "", "", identityRefusal(value, published, "published commit is not an exact commit identity")
	}
	if !git.OK("-C", root, "merge-base", "--is-ancestor", published, destination) {
		return "", "", "", "", errors.New("published commit is not reachable from the destination")
	}
	parents, err := git.Output("-C", root, "rev-list", "--parents", "-n", "1", published)
	if err != nil {
		return "", "", "", "", errors.New("published commit parents are unreadable")
	}
	parts := strings.Fields(parents)
	if len(parts) != 3 {
		return "", "", "", "", errors.New("published commit does not authenticate the reviewed source parent")
	}
	if parts[2] != source {
		return "", "", "", "", identityRefusal(source, parts[2], "published commit does not authenticate the reviewed source parent")
	}
	destinationBase = parts[1]
	sourceRange, kind, _ := diff.ResolveSourceRange(root, base, source)
	if kind != "" {
		return "", "", "", "", errors.New("review base does not authenticate the published source")
	}
	sourceBase = sourceRange.Base
	// A spec-less first run published no transition, so a resume without a spec has no
	// transition to authenticate. The parents and the range already bind the same
	// published commit either way, so no second publication can follow from the skip.
	if slug != "" {
		specPath := spec.LiveSpecPath(slug)
		staged, stagedErr := git.Raw("-C", root, "show", source+":"+specPath)
		implemented, implementedErr := spec.Implemented(staged)
		publishedSpec, publishedErr := git.Raw("-C", root, "show", published+":"+specPath)
		if stagedErr != nil || implementedErr != nil || publishedErr != nil || !bytes.Equal(implemented, publishedSpec) {
			return "", "", "", "", errors.New("published commit does not carry the source staged spec transition")
		}
	}
	tree, err = git.Output("-C", root, "rev-parse", published+"^{tree}")
	if err != nil {
		return "", "", "", "", errors.New("published tree is unreadable")
	}
	return published, sourceBase, destinationBase, tree, nil
}

type landingSourceFact struct {
	base, tip, specPath string
	fingerprint         string
	specBytes           []byte
	specMode            os.FileMode
}

func landingDestination(root string) (string, string, string, string, error) {
	branch, ok := git.ResolvedDefault(root)
	if !ok {
		return "", "", "", "", errors.New("default branch is unresolved")
	}
	current, err := git.CheckedOutBranch(root)
	if err != nil || current != branch {
		return "", "", "", "", errors.New("landing checkout is not attached to the default branch")
	}
	if raw, err := git.Raw("-C", root, "status", "--porcelain=v1", "-z", "--no-renames", "--untracked-files=all"); err != nil {
		return "", "", "", "", errors.New("landing destination is not clean")
	} else if entries, err := git.ParsePorcelainZStrict(raw); err != nil || len(entries) > 0 {
		paths := make([]string, 0, len(entries))
		for _, entry := range entries {
			paths = append(paths, entry.Path)
		}
		return "", "", "", "", refusalError{refusal{detail: "landing destination is not clean", paths: paths}}
	}
	ignored, _, ignoredErr := inventoryIgnored(root, false)
	declared, _, declarationErr := loadBuildOutputs(root)
	if ignoredErr != nil || declarationErr != nil || (ignored.Count > 0 && !ignoredWithinLandingAllowance(ignored, declared)) {
		if ignoredErr == nil && declarationErr == nil {
			return "", "", "", "", refusalError{refusal{detail: "landing destination has undeclared ignored residue", paths: ignored.Paths}}
		}
		return "", "", "", "", errors.New("landing destination has undeclared ignored residue")
	}
	tip, err := git.Output("-C", root, "rev-parse", "HEAD^{commit}")
	if err != nil {
		return "", "", "", "", errors.New("landing destination has no commit")
	}
	marker, err := landingMarker(root, branch, tip)
	if err != nil {
		return "", "", "", "", err
	}
	fingerprint, err := landing.CheckoutFingerprint(root)
	if err != nil {
		return "", "", "", "", errors.New("landing destination fingerprint is unavailable")
	}
	return tip, branch, marker, fingerprint, nil
}

func landingMarker(root, branch, destination string) (string, error) {
	marker, _, err := greenmarker.Read(root, branch)
	if err != nil {
		return "", err
	}
	if err := authorization.CheckMarker(context.Background(), root, branch, destination, marker); err != nil {
		return "", err
	}
	return marker, nil
}

func landingAssignment(root, path, request, base, requestedTip string) (intent.Assignment, error) {
	a, err := assignmentForRequest(root, request, assignmentRecoveryContext{
		target: path,
		detail: assignmentMismatchDetail,
		base:   base,
		tip:    requestedTip,
	})
	if err != nil || a.State != intent.StateActive || a.Worktree != path {
		if err != nil {
			return intent.Assignment{}, err
		}
		return intent.Assignment{}, errors.New(assignmentMismatchDetail)
	}
	evidence, err := validateOwnerMarker(root, path)
	if err != nil || evidence.marker.OwnerID != a.OwnerID || evidence.marker.Path != a.Worktree || evidence.registration.BranchRef != a.Branch || !evidence.registration.Locked || evidence.registration.LockReason != lockReason(a) {
		return intent.Assignment{}, errors.New("owner marker or assignment branch mismatch")
	}
	if requestedTip == "" {
		return intent.Assignment{}, errors.New("source tip is required")
	}
	return a, nil
}

func landingSource(root string, a intent.Assignment, base, requestedTip, slug string) (landingSourceFact, error) {
	branch, err := git.Output("-C", a.Worktree, "symbolic-ref", "--quiet", "HEAD")
	if err != nil || branch != a.Branch {
		return landingSourceFact{}, errors.New("assignment branch is not checked out")
	}
	head, err := git.Output("-C", a.Worktree, "rev-parse", "HEAD^{commit}")
	if err != nil {
		return landingSourceFact{}, errors.New("worktree source tip mismatch")
	}
	if head != requestedTip {
		return landingSourceFact{}, identityRefusal(requestedTip, head, "worktree source tip mismatch")
	}
	branchTip, err := git.Output("-C", root, "rev-parse", "--verify", a.Branch+"^{commit}")
	if err != nil {
		return landingSourceFact{}, errors.New("assignment branch source tip mismatch")
	}
	if branchTip != requestedTip {
		return landingSourceFact{}, identityRefusal(requestedTip, branchTip, "assignment branch source tip mismatch")
	}
	if dirty, err := git.Output("-C", a.Worktree, "status", "--porcelain=v1", "--untracked-files=all"); err != nil || dirty != "" {
		return landingSourceFact{}, errors.New("reviewed source is not clean")
	}
	rangeFact, detail, err := landingSourceRange(a.Worktree, slug, base, head)
	if err != nil {
		return landingSourceFact{}, err
	}
	if rangeFact.Tip != requestedTip {
		return landingSourceFact{}, identityRefusal(requestedTip, rangeFact.Tip, detail)
	}
	fact := landingSourceFact{base: rangeFact.Base, tip: rangeFact.Tip}
	if slug != "" {
		bytes, resolved, _, ok, err := spec.Resolve(a.Worktree, slug)
		if err != nil || !ok {
			return landingSourceFact{}, errors.New("staged spec is unreadable")
		}
		if _, err := spec.Implemented(bytes); err != nil {
			return landingSourceFact{}, err
		}
		rel, err := filepath.Rel(a.Worktree, resolved)
		if err != nil || rel == "." {
			return landingSourceFact{}, errors.New("staged spec path is invalid")
		}
		info, err := os.Lstat(resolved)
		if err != nil || !info.Mode().IsRegular() {
			return landingSourceFact{}, errors.New("staged spec is not a regular file")
		}
		fact.specPath, fact.specBytes, fact.specMode = filepath.ToSlash(rel), bytes, info.Mode().Perm()
	}
	fingerprint, err := landing.CheckoutFingerprint(a.Worktree)
	if err != nil {
		return landingSourceFact{}, errors.New("reviewed source fingerprint is unavailable")
	}
	fact.fingerprint = fingerprint
	return fact, nil
}

// landingSourceRange resolves the reviewed range and returns the refusal detail its
// identity mismatch carries. A spec names an ownership fence, so its range comes from
// the authorization owner. Without a spec there is no fence to authorize, so only the
// range resolves; every other source proof is unchanged either way.
func landingSourceRange(worktree, slug, base, head string) (diff.SourceRange, string, error) {
	if slug == "" {
		const detail = "reviewed source range is invalid"
		resolved, kind, hint := diff.ResolveSourceRange(worktree, base, head)
		if kind != "" {
			return diff.SourceRange{}, detail, fmt.Errorf("%s: %s: %s", detail, kind, hint)
		}
		return resolved, detail, nil
	}
	const detail = "reviewed source range or ownership fence is invalid"
	resolved, err := authorizeLandingSource(worktree, slug, base)
	if err != nil {
		return diff.SourceRange{}, detail, fmt.Errorf("%s: %s", detail, err)
	}
	return resolved, detail, nil
}

func identityRefusal(observed, wanted, detail string) error {
	return refusalError{refusal{detail: detail, observed: observed, wanted: wanted}}
}

// expandIdentity resolves an abbreviated commit identity to the full one the
// repository knows, so every later proof and every printed value pins the exact
// commit. Git refuses an ambiguous prefix, and a value that is not a hex prefix of
// a commit passes through unchanged for the proof that owns it to refuse.
func expandIdentity(repository, value string) string {
	full, err := git.Output("-C", repository, "rev-parse", "--verify", "--quiet", value+"^{commit}")
	if err == nil && abbreviatedIdentity(value, full) {
		return full
	}
	return value
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

func abbreviatedIdentity(value, full string) bool {
	if len(value) < 4 || len(value) > 39 || !fullCommitIdentity(full) {
		return false
	}
	if !hexIdentity(value) {
		return false
	}
	return strings.HasPrefix(strings.ToLower(full), strings.ToLower(value))
}

func reconcileLandingDestination(root, commit string) error {
	if err := exec.Command("git", "-C", root, "reset", "--hard", commit).Run(); err != nil {
		return err
	}
	return nil
}

func landedIncomplete(stdout io.Writer, result landing.ReviewedResult, specArg, path, assignment, step string) int {
	next := landingResumeNext(result, specArg, path, assignment)
	fmt.Fprintf(stdout, "landed{source_base=%s,source_tip=%s,destination_base=%s,published_commit=%s,tree=%s,worktree=incomplete:%s,next=%s}\n", result.SourceBase, result.SourceTip, result.DestinationBase, result.Commit, result.Tree, step, sanitize.Controls(next))
	return 3
}

// landingSlug is the spec slug a landing argument names, and the empty slug for the
// spec-less landing that named no argument.
func landingSlug(arg string) string {
	if arg == "" {
		return ""
	}
	return spec.LiveSpecSlug(arg)
}

func landingResumeNext(result landing.ReviewedResult, specArg, path, assignment string) string {
	values := []string{result.Commit, result.SourceBase, result.SourceTip, specArg}
	for _, value := range values {
		if !lineSafe(value) {
			return "bench worktree exec " + assignment + " -- bench worktree land --resume <full-published-commit> --request <request> --base <full-review-base> --source-tip <full-source-tip> --spec <spec> ."
		}
	}
	// A spec-less landing resumes spec-less, so the resume command it names carries no
	// --spec at all rather than an empty value the grammar refuses.
	specFlag := ""
	if specArg != "" {
		specFlag = " --spec " + sanitize.ShellQuote(specArg)
	}
	command := "bench worktree land --resume " + sanitize.ShellQuote(result.Commit) + " --request <request> --base " + sanitize.ShellQuote(result.SourceBase) + " --source-tip " + sanitize.ShellQuote(result.SourceTip) + specFlag
	if lineSafe(path) {
		return command + " " + sanitize.ShellQuote(path)
	}
	return "bench worktree exec " + assignment + " -- " + command + " ."
}

func landRefusal(stdout io.Writer, detail string) int {
	fmt.Fprintln(stdout, "refused{detail="+sanitize.Controls(detail)+"}")
	return 1
}

func landRefusalError(stdout io.Writer, err error) int {
	var typed refusalError
	if !errors.As(err, &typed) {
		return landRefusal(stdout, err.Error())
	}
	fmt.Fprintln(stdout, "refused{"+typed.fields()+"}")
	fmt.Fprint(stdout, typed.table())
	return 1
}
