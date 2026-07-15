package worktree

import (
	"errors"
	"os"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

var worktreeListFields = []string{"id", "label", "state", "source", "tree", "lease", "landed", "ignored"}

// ListCommand implements the read-only AXI worktree population query.
func ListCommand(args []string) (string, int) {
	if len(args) > 0 {
		if len(args) == 1 && (args[0] == "-h" || args[0] == "--help") {
			return "usage: " + usage.WorktreeList + "\n", 0
		}
		return toon.Usage(usage.WorktreeList, args[0]) + "\n", 2
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	assignments, err := intent.Assignments(root)
	if err != nil {
		return toon.Errorf("cannot read worktree assignments", "repair the Bench intent ledger and retry") + "\n", 1
	}
	registrations, err := git.Worktrees(root)
	if err != nil {
		return toon.Errorf("cannot read registered worktrees", "run git worktree list and retry") + "\n", 1
	}

	def, defaultResolved := git.ResolvedDefault(root)
	rows := make([][]any, 0, len(assignments)+len(registrations))
	assignedPaths := make(map[string]bool, len(assignments))
	for _, assignment := range assignments {
		assignedPaths[assignment.Worktree] = true
		rows = append(rows, listAssignmentRow(root, assignment, def, defaultResolved))
	}
	mainRoot := canonicalRoot(root)
	for _, registration := range registrations {
		if samePath(registration.Path, mainRoot) || assignedPaths[registration.Path] {
			continue
		}
		label := registration.Branch
		if label == "" {
			label = registration.Path
		}
		rows = append(rows, []any{"foreign", label, "foreign", "foreign", "present", listLease(registration.Path), listLanded(root, registration.Branch, def, defaultResolved), listIgnored(registration.Path)})
	}
	out, err := toon.TableTyped("worktrees", worktreeListFields, rows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return out, 0
}

func listAssignmentRow(root string, assignment intent.Assignment, def string, defaultResolved bool) []any {
	tree := "present"
	if _, err := os.Stat(assignment.Worktree); err != nil {
		tree = "missing"
	}
	return []any{assignment.ID, assignment.Label, string(assignment.State), "assignment", tree, listLease(assignment.Worktree), listLanded(root, assignment.Branch, def, defaultResolved), listIgnored(assignment.Worktree)}
}

func listLease(path string) string {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "none"
		}
		return "unknown"
	}
	lease, err := LeaseFile(path)
	if err != nil {
		return "unknown"
	}
	if _, err := os.Stat(lease); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "none"
		}
		return "unknown"
	}
	return string(ProbeLease(lease))
}

func listLanded(root, branch, def string, defaultResolved bool) any {
	if !defaultResolved || branch == "" {
		return "unknown"
	}
	landed, _, err := git.LandedInDefault(root, branch, def)
	if err != nil {
		return "unknown"
	}
	return landed
}

func listIgnored(path string) any {
	if _, err := os.Stat(path); err != nil {
		return "unknown"
	}
	inventory, _, err := inventoryIgnored(path, true)
	if err != nil || inventory.Uncertain {
		return "unknown"
	}
	return inventory.Count
}
