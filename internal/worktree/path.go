package worktree

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/usage"
)

// PathCommand resolves one active Bench-owned assignment and prints its portable path.
func PathCommand(root string, args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 || strings.TrimSpace(args[0]) == "" {
		fmt.Fprintln(stderr, "usage: "+usage.WorktreePath)
		return 2
	}
	path, err := resolvePath(root, args[0])
	if err != nil {
		fmt.Fprintln(stderr, "bench worktree path: target is not one active Bench-owned worktree")
		return 1
	}
	if !lineSafe(path) {
		fmt.Fprintln(stderr, "bench worktree path: resolved path is not safe for line output")
		return 1
	}
	fmt.Fprintln(stdout, path)
	return 0
}

func resolvePath(root, target string) (string, error) {
	path, err := resolveWorktree(root, target)
	if err != nil {
		return "", err
	}
	return compactHomePath(path)
}

func resolveWorktree(root, target string) (string, error) {
	if !lineSafe(target) {
		return "", errors.New("target contains control characters")
	}
	assignments, err := intent.Assignments(root)
	if err != nil {
		return "", err
	}
	selected, err := selectAssignment(assignments, target)
	if err != nil {
		return "", err
	}
	if selected.State != intent.StateActive {
		return "", errors.New("assignment is not active")
	}
	if err := validateCreationBundle(root, selected); err != nil {
		return "", err
	}
	return selected.Worktree, nil
}

func planAbandon(root, request, path string) (CleanupPlan, error) {
	repo, target, err := cleanupIdentity(root, path)
	if err != nil {
		return CleanupPlan{}, err
	}
	if plan, found, err := abandonReceipt(root, repo, target, requestDigest(request)); err != nil || found {
		return plan, err
	}
	if _, err := os.Lstat(target); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return CleanupPlan{}, errors.New("abandon request, assignment, or path mismatch; checkout retained")
		}
		return CleanupPlan{}, err
	}
	plan, err := PlanExplicit(root, target)
	if err != nil {
		return CleanupPlan{}, err
	}
	if !plan.owned || plan.assignment == nil || plan.assignment.Request != requestDigest(request) {
		return CleanupPlan{}, errors.New("abandon request, assignment, or path mismatch; checkout retained")
	}
	return plan, nil
}

func abandonReceipt(root, repo, target, request string) (CleanupPlan, bool, error) {
	receipt, found, err := intent.CleanupReceiptForRequest(root, repo, cleanupOperation, target, request)
	if err != nil || !found {
		return CleanupPlan{}, found, err
	}
	if receipt.State != intent.ReceiptInFlight && receipt.State != intent.ReceiptComplete {
		return CleanupPlan{}, true, errors.New("abandon request, assignment, or path mismatch; checkout retained")
	}
	assigned, active, err := intent.FindAssignmentByRequest(root, request)
	if err != nil {
		return CleanupPlan{}, true, err
	}
	if !active {
		if receipt.State == intent.ReceiptComplete {
			return planFromReceipt(receipt), true, nil
		}
		return CleanupPlan{}, true, errors.New("abandon request, assignment, or path mismatch; checkout retained")
	}
	assignmentPath, err := canonicalPath(assigned.Worktree)
	if err != nil || assigned.ID != receipt.Assignment || assigned.OwnerID != receipt.Owner || assigned.Request != request || assignmentPath != target {
		return CleanupPlan{}, true, errors.New("abandon request, assignment, or path mismatch; checkout retained")
	}
	return planFromReceipt(receipt), true, nil
}

// PlanAbandon returns the exact recovery-aware cleanup fingerprint for one owned assignment.
func PlanAbandon(root, request, path string) (string, error) {
	plan, err := planAbandon(root, request, path)
	return plan.Fingerprint, err
}

// ApplyAbandon preserves Git-visible work before removing one exact owned assignment.
func ApplyAbandon(root, request, path, fingerprint string) (CleanupPlan, error) {
	plan, err := planAbandon(root, request, path)
	if err != nil {
		return plan, err
	}
	if plan.Fingerprint != fingerprint {
		return plan, errStaleFingerprint
	}
	applied, err := ApplyExplicit(root, path, plan.Fingerprint)
	if err != nil || applied.Action != ActionRetain {
		return applied, err
	}
	return applied, retainedReleaseError(applied, request, path)
}

func selectAssignment(assignments []intent.Assignment, target string) (intent.Assignment, error) {
	path, isPath, err := targetPath(target)
	if err != nil {
		return intent.Assignment{}, err
	}
	var selected *intent.Assignment
	for i := range assignments {
		matches := assignments[i].Label == target
		if isPath {
			worktree, worktreeErr := canonicalPath(assignments[i].Worktree)
			matches = worktreeErr == nil && worktree == path
		}
		if !matches {
			continue
		}
		if selected != nil {
			return intent.Assignment{}, errors.New("target is ambiguous")
		}
		selected = &assignments[i]
	}
	if selected == nil {
		return intent.Assignment{}, errors.New("target is unassigned")
	}
	return *selected, nil
}

func targetPath(target string) (string, bool, error) {
	if target == "~" || strings.HasPrefix(target, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", false, err
		}
		if target == "~" {
			path, err := canonicalPath(home)
			return path, true, err
		}
		path, err := canonicalPath(filepath.Join(home, target[2:]))
		return path, true, err
	}
	if strings.HasPrefix(target, "~") {
		return "", false, errors.New("unsupported home target")
	}
	if filepath.IsAbs(target) {
		path, err := canonicalPath(target)
		return path, true, err
	}
	if target == "." || strings.HasPrefix(target, "."+string(filepath.Separator)) || strings.HasPrefix(target, ".."+string(filepath.Separator)) || strings.Contains(target, string(filepath.Separator)) {
		return "", false, errors.New("relative path targets are unsupported")
	}
	return "", false, nil
}

func compactHomePath(path string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	home, err = canonicalPath(home)
	if err != nil {
		return "", err
	}
	path, err = canonicalPath(path)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(home, path)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return path, nil
	}
	if rel == "." {
		return "~", nil
	}
	return filepath.Join("~", rel), nil
}
