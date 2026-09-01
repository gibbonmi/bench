// Landing identity resolution: destination, marker, assignment, source range, and identity
// reconciliation, with the checkout predicates every identity proof in the package shares.
package worktree

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/gate/greenmarker"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/spec"
)

type landingSourceFact struct {
	base, tip, specPath string
	// closePath is the tickets-only folder the landing consumes, empty otherwise.
	closePath   string
	fingerprint string
	specBytes   []byte
	specMode    os.FileMode
}

func landingDestination(j joins, root string) (string, string, string, string, error) {
	branch, ok := git.ResolvedDefault(root)
	if !ok {
		return "", "", "", "", errors.New("default branch is unresolved")
	}
	current, err := git.CheckedOutBranch(root)
	if err != nil || current != branch {
		return "", "", "", "", errors.New("landing checkout is not attached to the default branch")
	}
	if err := checkoutClean(root, "landing destination is not clean", ""); err != nil {
		return "", "", "", "", err
	}
	ignored, _, ignoredErr := inventoryIgnored(j, root, false)
	declared, _, declarationErr := loadBuildOutputs(root)
	if ignoredErr != nil || declarationErr != nil || (ignored.Count > 0 && !ignoredWithinLandingAllowance(ignored, declared)) {
		if ignoredErr == nil && declarationErr == nil {
			return "", "", "", "", refusalError{refusal{detail: "landing destination has undeclared ignored residue", paths: undeclaredLandingIgnoredPaths(ignored, declared)}}
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

func landingAssignment(j joins, root, path, request, base, requestedTip string) (intent.Assignment, error) {
	a, err := assignmentForRequest(root, request, assignmentRecoveryContext{
		target: path,
		base:   base,
		tip:    requestedTip,
	})
	if err != nil {
		return intent.Assignment{}, err
	}
	if err := identityBundleRefusal(root, path, a, landingActiveState); err != nil {
		return intent.Assignment{}, err
	}
	if requestedTip == "" {
		return intent.Assignment{}, errors.New("source tip is required")
	}
	return a, nil
}

func landingSource(j joins, root string, a intent.Assignment, base, requestedTip, slug string) (landingSourceFact, error) {
	if _, ok := assignmentBranchCheckedOut(a); !ok {
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
	if err := checkoutClean(a.Worktree, "reviewed source is not clean", ""); err != nil {
		// The hostile-source surface pins this refusal to one line, so the source reads the
		// shared fact and states it without the predicate's path table.
		return landingSourceFact{}, errors.New("reviewed source is not clean")
	}
	// A --spec naming a tickets-only folder is a close, not a staged spec. It names no
	// ownership fence and carries no transition, so its range resolves and its proof
	// runs exactly as the spec-less landing's do.
	closeSlug := ""
	if slug != "" && landing.TicketsOnlyFolder(a.Worktree, landingSlug(slug)) {
		closeSlug = landingSlug(slug)
	}
	rangeSlug := slug
	if closeSlug != "" {
		rangeSlug = ""
	}
	rangeFact, detail, err := landingSourceRange(j, a.Worktree, rangeSlug, base, head)
	if err != nil {
		return landingSourceFact{}, err
	}
	if rangeFact.Tip != requestedTip {
		return landingSourceFact{}, identityRefusal(requestedTip, rangeFact.Tip, detail)
	}
	fact := landingSourceFact{base: rangeFact.Base, tip: rangeFact.Tip}
	if closeSlug != "" {
		fact.closePath = landing.ClosedFolderPath(closeSlug)
	} else if slug != "" {
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
func landingSourceRange(j joins, worktree, slug, base, head string) (diff.SourceRange, string, error) {
	if slug == "" {
		const detail = "reviewed source range is invalid"
		resolved, kind, hint := diff.ResolveSourceRange(worktree, base, head)
		if kind != "" {
			return diff.SourceRange{}, detail, fmt.Errorf("%s: %s: %s", detail, kind, hint)
		}
		return resolved, detail, nil
	}
	detail := landingRefusalFaceByName(faceSourceNotFenced).detail
	resolved, err := j.authorizeLandingSource(worktree, slug, base)
	if err != nil {
		// The unfenced paths arrive typed, so they print as the refusal's own path table
		// rather than inside the sentence. The assembler holds the caller's flags, so it
		// attaches this face's route there.
		var unfenced preflight.UnauthorizedPathsError
		if errors.As(err, &unfenced) && len(unfenced.Paths) > 0 {
			return diff.SourceRange{}, detail, landingFaceRefusal(faceSourceNotFenced, "", unfenced.Paths)
		}
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

func abbreviatedIdentity(value, full string) bool {
	if len(value) < 4 || len(value) > 39 || !fullCommitIdentity(full) {
		return false
	}
	if !hexIdentity(value) {
		return false
	}
	return strings.HasPrefix(strings.ToLower(full), strings.ToLower(value))
}

func reconcileLandingDestination(j joins, root, destination, published, destinationBase string) error {
	if err := resumeDestructiveDestinationState(j, root, destination, published, destinationBase); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", root, "reset", "--merge", destination).Run(); err != nil {
		return err
	}
	if err := resumeDestructiveDestinationState(j, root, destination, published, destinationBase); err != nil {
		return err
	}
	return nil
}

// checkoutClean proves one checkout holds no uncommitted work. It refuses any status
// entry, untracked included; ignored residue is excluded, because a worktree's build
// output is not uncommitted work. Each caller states the detail and the repair its own
// refusal carries, because the fact reads the same but the repair does not.
func checkoutClean(path, detail, next string) error {
	raw, err := git.Raw("-C", path, "status", "--porcelain=v1", "-z", "--no-renames", "--untracked-files=all")
	if err != nil {
		return refusalError{refusal{detail: "checkout status is unreadable"}}
	}
	entries, err := git.ParsePorcelainZStrict(raw)
	if err != nil {
		return refusalError{refusal{detail: "checkout status is unreadable"}}
	}
	if len(entries) == 0 {
		return nil
	}
	paths := make([]string, 0, len(entries))
	for _, entry := range entries {
		paths = append(paths, entry.Path)
	}
	return refusalError{refusal{detail: detail, next: next, paths: paths}}
}

// assignmentBranchCheckedOut reports the ref one assignment's checkout has attached, and
// whether that ref is the assignment's own branch. The observed value is returned even
// when it disagrees, because a refusal that names it tells the operator what to restore.
func assignmentBranchCheckedOut(a intent.Assignment) (string, bool) {
	branch, err := git.Output("-C", a.Worktree, "symbolic-ref", "--quiet", "HEAD")
	return branch, err == nil && branch == a.Branch
}
