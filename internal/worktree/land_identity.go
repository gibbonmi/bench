// Landing identity resolution: destination, marker, assignment, source range, and identity reconciliation.
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

func landingAssignment(root, path, request, base, requestedTip string) (intent.Assignment, error) {
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
	rangeFact, detail, err := landingSourceRange(a.Worktree, rangeSlug, base, head)
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

func abbreviatedIdentity(value, full string) bool {
	if len(value) < 4 || len(value) > 39 || !fullCommitIdentity(full) {
		return false
	}
	if !hexIdentity(value) {
		return false
	}
	return strings.HasPrefix(strings.ToLower(full), strings.ToLower(value))
}

func reconcileLandingDestination(root, destination, published, destinationBase string) error {
	if err := resumeDestructiveDestinationState(root, destination, published, destinationBase); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", root, "reset", "--merge", destination).Run(); err != nil {
		return err
	}
	if err := resumeDestructiveDestinationState(root, destination, published, destinationBase); err != nil {
		return err
	}
	return nil
}
