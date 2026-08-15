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

var worktreeListGrammar = usage.Grammar{Cmd: usage.WorktreeList, Help: "usage: " + usage.WorktreeList, HelpOnlyWhenSole: true, UnquotedEmptyPositional: true}

type listRow struct {
	values     []any
	orphanPath string
}

// ListCommand implements the read-only AXI worktree population query.
func ListCommand(args []string) (string, int) {
	_, line, code := usage.Parse(worktreeListGrammar, args)
	if line != "" {
		return line + "\n", code
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	registrations, err := git.Worktrees(root)
	if err != nil {
		var typed git.WorktreeFailure
		if errors.As(err, &typed) {
			return toon.Errorf(typed.Error(), typed.WorktreeAction()) + "\n", 1
		}
		return toon.Errorf("cannot read registered worktrees", "run git worktree list and retry") + "\n", 1
	}
	assignments, err := intent.Assignments(root)
	if err != nil {
		return toon.Errorf("cannot read worktree assignments", "repair the Bench intent ledger and retry") + "\n", 1
	}

	def, defaultResolved := git.ResolvedDefault(root)
	rows := make([]listRow, 0, len(assignments)+len(registrations))
	assignedPaths := make(map[string]bool, len(assignments))
	for _, assignment := range assignments {
		assignedPaths[assignment.Worktree] = true
		rows = append(rows, listRow{values: listAssignmentRow(root, assignment, def, defaultResolved)})
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
		row := listRow{values: []any{"foreign", label, "foreign", "foreign", tree, listLease(registration.Path), listLanded(root, registration.Branch, def, defaultResolved)}}
		row.values = append(row.values, listIgnored(registration.Path))
		if tree == "missing" {
			row.orphanPath = registration.Path
		}
		rows = append(rows, row)
	}
	values := make([][]any, 0, len(rows))
	for _, row := range rows {
		values = append(values, row.values)
	}
	out, err := toon.TableTyped("worktrees", worktreeListFields, values)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	help, err := axi.RenderHelp(actionsForRows(rows))
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return out + help, 0
}

func actionsForRows(rows []listRow) []axi.Action {
	actions := make([]axi.Action, 0, len(rows))
	for _, row := range rows {
		if len(row.values) < 3 {
			continue
		}
		if row.values[2] == string(intent.StateActive) {
			id, ok := row.values[0].(string)
			if !ok || id == "" {
				continue
			}
			actions = append(actions,
				axi.ExecutableInvocation("inspect active worktree", axi.KnownArgument("worktree"), axi.KnownArgument("path"), axi.KnownArgument(id)),
				axi.ExecutableInvocation("run a command in the active worktree", axi.KnownArgument("worktree"), axi.KnownArgument("exec"), axi.KnownArgument(id), axi.KnownArgument("--"), axi.FutureInput("command")))
			continue
		}
		if row.values[2] == "foreign" && row.orphanPath != "" {
			actions = append(actions, axi.ExecutableInvocation("clean the orphaned worktree", axi.KnownArgument("worktree"), axi.KnownArgument("clean"), axi.KnownArgument(row.orphanPath)))
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
