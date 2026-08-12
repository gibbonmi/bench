package worktree

import (
	"errors"
	"os"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

var worktreeListFields = []string{"id", "label", "state", "source", "tree", "lease", "landed", "ignored"}

var worktreeListGrammar = usage.Grammar{Cmd: usage.WorktreeList, Help: "usage: " + usage.WorktreeList}

// ListCommand implements the read-only AXI worktree population query.
func ListCommand(args []string) (string, int) {
	parsed, line, code := usage.Parse(worktreeListGrammar, args)
	if line != "" {
		return line + "\n", code
	}
	if parsed.EndedFlags {
		return toon.Usage(usage.WorktreeList, "--") + "\n", 2
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
	orphanPaths := make([]string, 0, len(registrations))
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
		tree := listTree(registration.Path)
		rows = append(rows, []any{"foreign", label, "foreign", "foreign", tree, listLease(registration.Path), listLanded(root, registration.Branch, def, defaultResolved), listIgnored(registration.Path)})
		if tree == "missing" {
			orphanPaths = append(orphanPaths, registration.Path)
		}
	}
	out, err := toon.TableTyped("worktrees", worktreeListFields, rows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	help, err := axi.RenderHelp(actionsForRows(rows, orphanPaths))
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return out + help, 0
}

func actionsForRows(rows [][]any, orphanPaths []string) []axi.Action {
	actions := make([]axi.Action, 0, len(rows))
	orphanIndex := 0
	for _, row := range rows {
		if len(row) < 3 {
			continue
		}
		if row[2] == string(intent.StateActive) {
			id, ok := row[0].(string)
			if !ok || id == "" {
				continue
			}
			actions = append(actions,
				axi.ExecutableInvocation("inspect active worktree", axi.KnownArgument("worktree"), axi.KnownArgument("path"), axi.KnownArgument(id)),
				axi.ExecutableInvocation("run a command in the active worktree", axi.KnownArgument("worktree"), axi.KnownArgument("exec"), axi.KnownArgument(id), axi.KnownArgument("--"), axi.FutureInput("command")))
			continue
		}
		if row[2] == "foreign" && len(row) > 4 && row[4] == "missing" && orphanIndex < len(orphanPaths) {
			actions = append(actions, axi.ExecutableInvocation("clean the orphaned worktree", axi.KnownArgument("worktree"), axi.KnownArgument("clean"), axi.KnownArgument(orphanPaths[orphanIndex])))
			orphanIndex++
		}
	}
	return actions
}

func listAssignmentRow(root string, assignment intent.Assignment, def string, defaultResolved bool) []any {
	return []any{assignment.ID, assignment.Label, string(assignment.State), "assignment", listTree(assignment.Worktree), listLease(assignment.Worktree), listLanded(root, assignment.Branch, def, defaultResolved), listIgnored(assignment.Worktree)}
}

func listTree(path string) string {
	if _, err := os.Stat(path); err != nil {
		return "missing"
	}
	return "present"
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
