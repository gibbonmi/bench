package worktree

import (
	"fmt"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// selectedWorktreeGrammar is the selected view's grammar. ListCommand takes this route only
// when a selector flag is present, so every other argument list keeps the bare grammar.
var selectedWorktreeGrammar = usage.Grammar{
	Cmd:  usage.WorktreeList,
	Help: "usage: " + usage.WorktreeListPaths,
	Flags: []usage.Flag{
		{Name: "--view", HasValue: true, NoEmptyValue: true, Required: true},
		{Name: "--target", HasValue: true, NoEmptyValue: true, Required: true, Repeatable: true},
	},
}

// selectedPathUnrepresentable is the error cell of a resolved record whose stored path the
// TOON table cannot carry.
const selectedPathUnrepresentable = "assignment path is not representable"

// selectsWorktrees reports whether args choose the selected view. The check runs before
// the bare grammar parses, because that grammar refuses every flag.
func selectsWorktrees(args []string) bool {
	return usage.FlagPresent(selectedWorktreeGrammar, args, "--view") || usage.FlagPresent(selectedWorktreeGrammar, args, "--target")
}

// listSelectedWorktrees answers identity, path, and state for each requested target. It
// reads the ledger through the assignment selector and not through the active-path
// resolver, so a non-active record answers and no target takes worktree authority. Each
// distinct failed operand keeps its own row, and a resolved identity appears once, at
// the first request that names it.
func listSelectedWorktrees(root string, args []string) (string, int) {
	parsed, line, code := usage.Parse(selectedWorktreeGrammar, args)
	if line != "" {
		return line + "\n", code
	}
	if parsed.Flags["--view"] != "paths" || parsed.EndedFlags {
		return selectedWorktreeGrammar.Help + "\n", 2
	}
	if !inRepository(root) {
		return toon.NotInRepo() + "\n", 1
	}
	assignments, err := intent.Assignments(root)
	if err != nil {
		return assignmentsReadRefusal(), 1
	}
	var rows [][]any
	exit := 0
	seenTargets := map[string]bool{}
	seenIDs := map[string]bool{}
	for i, target := range parsed.Repeated["--target"] {
		if seenTargets[target] {
			continue
		}
		seenTargets[target] = true
		// An unsafe operand never reaches a cell: its row names the request ordinal.
		if !lineSafe(target) {
			rows = append(rows, []any{fmt.Sprintf("target-%d", i+1), "", "", "", errTargetControls.Error()})
			exit = 1
			continue
		}
		selected, err := selectAssignment(assignments, target)
		if err != nil {
			rows = append(rows, []any{target, "", "", "", err.Error()})
			exit = 1
			continue
		}
		// The identity check runs before the stored-path refusal, so two aliases of one
		// unrepresentable record give one refusal row.
		if seenIDs[selected.ID] {
			continue
		}
		seenIDs[selected.ID] = true
		if !toon.Representable(selected.Worktree) {
			rows = append(rows, []any{target, "", "", "", selectedPathUnrepresentable})
			exit = 1
			continue
		}
		rows = append(rows, []any{target, selected.ID, selected.Worktree, string(selected.State), ""})
	}
	out, err := toon.TableTyped("worktrees", []string{"target", "id", "path", "state", "error"}, rows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	help, err := axi.RenderHelp([]axi.Action{axi.ExecutableInvocation("inspect the complete worktree inventory", axi.KnownArgument("worktree"), axi.KnownArgument("list"))})
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return out + help, exit
}
