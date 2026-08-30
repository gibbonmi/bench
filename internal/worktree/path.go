package worktree

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/usage"
)

// PathCommand resolves one active Bench-owned assignment and prints its resolved
// absolute path. A quoted `~` never expands in a shell, so the verb emits the form
// every caller can paste; the path-taking verbs still accept the `~` form.
func PathCommand(root, _ string, args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 || strings.TrimSpace(args[0]) == "" {
		fmt.Fprintln(stderr, "usage: "+usage.WorktreePath)
		return 2
	}
	path, err := resolvePath(root, args[0])
	if err != nil {
		return printTargetRefusal(stderr, "bench worktree path", err)
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
	selected, err := resolveAssignment(root, target)
	if err != nil {
		return "", err
	}
	return selected.Worktree, nil
}

// resolveAssignment is the one authority every target-taking verb resolves through: it
// answers the whole record, so a verb that names the assignment in its output reads the
// ledger once. A verb that needs the path alone composes resolveWorktree.
func resolveAssignment(root, target string) (intent.Assignment, error) {
	if !lineSafe(target) {
		return intent.Assignment{}, errors.New("target contains control characters")
	}
	assignments, err := intent.Assignments(root)
	if err != nil {
		return intent.Assignment{}, err
	}
	selected, err := selectAssignment(assignments, target)
	if err != nil {
		return intent.Assignment{}, err
	}
	if !landingActiveState(selected.State) {
		return intent.Assignment{}, componentRefusal(componentAssignmentState, selected.ID, string(selected.State), string(intent.StateActive))
	}
	// The tree check runs before the creation bundle, because every bundle component reads
	// evidence inside the tree. A tree that is gone is the fact the operator acts on. Any
	// other stat error leaves the bundle to describe what it can read.
	if _, statErr := os.Stat(selected.Worktree); statErr != nil && errors.Is(statErr, os.ErrNotExist) {
		return intent.Assignment{}, missingTreeRefusal(root, selected)
	}
	if err := validateCreationBundle(root, selected); err != nil {
		return intent.Assignment{}, err
	}
	return selected, nil
}

// missingTreeRefusal names the fact the operator can act on and the one verb that clears
// the record, so a retained assignment whose tree is gone is never a dead end. The
// landedness comes from the branch, because a tree that is gone proves nothing.
func missingTreeRefusal(root string, assignment intent.Assignment) error {
	def, defaultResolved := git.ResolvedDefault(root)
	landed := listLanded(root, assignment.Branch, def, defaultResolved) == true
	recovery := recoverMissingTree(landed, assignment.RequestToken, assignment.Worktree)
	return refusalError{refusal{detail: "worktree tree is missing", next: recovery.line()}}
}

// errTargetUnassigned is the selector outcome a caller may act on rather than report: the
// spelling names no assignment, so a caller with a second lookup falls through to it. The
// ambiguity outcome carries the colliding ids, so it stays a typed error a caller can
// tell apart from the grammar and lookup outcomes it may fall through.
var errTargetUnassigned = errors.New("target is unassigned")

// ambiguousTargetError is the one selector outcome no second lookup can resolve: the
// spelling names more than one assignment. Callers match it with errors.As.
type ambiguousTargetError struct{ IDs []string }

func (e ambiguousTargetError) Error() string {
	return "target is ambiguous: " + strings.Join(e.IDs, ", ")
}

func selectAssignment(assignments []intent.Assignment, target string) (intent.Assignment, error) {
	path, isPath, err := targetPath(target)
	if err != nil {
		return intent.Assignment{}, err
	}
	matched := matchingAssignments(assignments, func(a intent.Assignment) bool {
		if isPath {
			worktree, worktreeErr := canonicalPath(a.Worktree)
			return worktreeErr == nil && worktree == path
		}
		// The id is the address `bench worktree list` advertises in its executable help
		// rows, and it is the only unique one. Labels can collide, and the ambiguity
		// refusal below is what answers when they do.
		return a.Label == target || a.ID == target
	})
	if len(matched) == 0 && !isPath && len(target) >= minOperandPrefix && len(target) <= maxIdentifierPrefix {
		// An unambiguous 8-12 character prefix of the label or the id also resolves.
		// Shorter prefixes stay unresolved so a short word cannot grab a worktree.
		matched = matchingAssignments(assignments, func(a intent.Assignment) bool {
			return strings.HasPrefix(a.ID, target) || strings.HasPrefix(a.Label, target)
		})
	}
	if len(matched) == 0 {
		return intent.Assignment{}, errTargetUnassigned
	}
	if len(matched) > 1 {
		// The refusal names every colliding id, because the id is the address that
		// resolves the collision the operator just hit.
		ids := make([]string, 0, len(matched))
		for _, a := range matched {
			ids = append(ids, a.ID)
		}
		return intent.Assignment{}, ambiguousTargetError{IDs: ids}
	}
	return matched[0], nil
}

func matchingAssignments(assignments []intent.Assignment, matches func(intent.Assignment) bool) []intent.Assignment {
	var selected []intent.Assignment
	for _, a := range assignments {
		if matches(a) {
			selected = append(selected, a)
		}
	}
	return selected
}

// printTargetRefusal is the one printer every target-taking verb uses, so `worktree path`,
// `worktree exec`, `worktree show`, `worktree build`, and `worktree create` cannot each
// describe one failure their own way. A component refusal prints
// its detail sentence: the operator reads the named check, not the refused record. The
// second line names the verb that answers the refusal; a target that never resolved is
// answered by the lookup, so that is the default.
func printTargetRefusal(stderr io.Writer, verb string, err error) int {
	reason, next := err.Error(), "bench worktree list"
	var refused refusalError
	if errors.As(err, &refused) {
		reason = refused.detail
		if refused.next != "" {
			next = refused.next
		}
	}
	fmt.Fprintln(stderr, verb+": "+sanitize.Controls(reason))
	fmt.Fprintln(stderr, "next="+sanitize.Controls(next))
	return 1
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
