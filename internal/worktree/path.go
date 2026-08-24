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

// PathCommand resolves one active Bench-owned assignment and prints its resolved
// absolute path. A quoted `~` never expands in a shell, so the verb emits the form
// every caller can paste; the path-taking verbs still accept the `~` form.
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
	return canonicalPath(path)
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

func selectAssignment(assignments []intent.Assignment, target string) (intent.Assignment, error) {
	path, isPath, err := targetPath(target)
	if err != nil {
		return intent.Assignment{}, err
	}
	var selected *intent.Assignment
	for i := range assignments {
		// The id is the address `bench worktree list` advertises in its executable help
		// rows, and it is the only unique one. Labels can collide, and the ambiguity
		// guard below is what answers when they do.
		matches := assignments[i].Label == target || assignments[i].ID == target
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
	if selected == nil && !isPath && len(target) >= minOperandPrefix && len(target) <= maxIdentifierPrefix {
		// An unambiguous 8-12 character prefix of the label or the id also resolves.
		// Shorter prefixes stay unresolved so a short word cannot grab a worktree.
		for i := range assignments {
			if !strings.HasPrefix(assignments[i].ID, target) && !strings.HasPrefix(assignments[i].Label, target) {
				continue
			}
			if selected != nil {
				return intent.Assignment{}, errors.New("target is ambiguous")
			}
			selected = &assignments[i]
		}
	}
	if selected == nil {
		return intent.Assignment{}, errors.New("target is unassigned")
	}
	return *selected, nil
}

// minOperandPrefix and maxIdentifierPrefix bound the unambiguous-prefix window. The
// floor is one policy fact shared by the identifier form and the fingerprint form, so
// it has one source; only the identifier form carries a ceiling.
const (
	minOperandPrefix    = 8
	maxIdentifierPrefix = 12
)

// resolveVerbOperand widens every path-taking worktree verb's operand: a label, an
// assignment id, or an unambiguous 8-12 character prefix of either resolves to the
// assignment's worktree path. A path-shaped or unresolvable operand returns unchanged,
// so each verb keeps its own refusal for it.
func resolveVerbOperand(root, operand string) string {
	if path, err := resolveWorktree(root, operand); err == nil {
		return path
	}
	return operand
}

// expandHomeTarget resolves the portable `~`-prefixed form that every worktree command
// prints, and owns that grammar for all of them. The second result reports whether the
// target used the home form at all, so a caller can leave an ordinary path untouched.
// `~user` stays unsupported, because the address this expands is one Bench itself
// printed, and Bench never prints another account's home.
func expandHomeTarget(target string) (string, bool, error) {
	if target != "~" && !strings.HasPrefix(target, "~/") {
		if strings.HasPrefix(target, "~") {
			return "", true, errors.New("unsupported home target")
		}
		return "", false, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", true, err
	}
	if target == "~" {
		return home, true, nil
	}
	return filepath.Join(home, target[2:]), true, nil
}

func targetPath(target string) (string, bool, error) {
	if expanded, isHome, err := expandHomeTarget(target); isHome {
		if err != nil {
			return "", false, err
		}
		path, err := canonicalPath(expanded)
		return path, true, err
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
