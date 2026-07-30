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
	return compactHomePath(selected.Worktree)
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
