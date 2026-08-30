// The worktree build verb: one executable at <worktree>/dist/bench through the tree's own build script.
package worktree

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// BuildCommand builds one active Bench-owned worktree's tree into that worktree's
// `dist/bench`, so the assignment's own grammar has an executable. The build runs the
// tree's sanctioned build script, which writes the executable and its seal, so the verb
// authors no artifact of its own.
func BuildCommand(root, _ string, args []string, stdout, stderr io.Writer) int {
	return buildWith(defaultJoins(), root, args, stdout, stderr)
}

// buildWith is BuildCommand with the seam set resolved at the caller's boundary.
func buildWith(j joins, root string, args []string, stdout, stderr io.Writer) int {
	parsed, line, code := usage.Parse(buildGrammar, args)
	if line != "" {
		if code == 0 {
			fmt.Fprintln(stdout, line)
			return 0
		}
		fmt.Fprintln(stderr, line)
		return code
	}
	path, assignment, err := resolveBuildTarget(root, parsed.Positionals[0])
	if err != nil {
		return printTargetRefusal(stderr, buildVerb, err)
	}
	output := filepath.Join(path, "dist", "bench")
	ctx, stop := subprocess.NotifyCancel(context.Background())
	defer stop()
	if buildErr := j.build(ctx, path, output); buildErr != nil {
		fmt.Fprintf(stderr, "%s: %v\n", buildVerb, buildErr)
		return nameWorktree(stderr, path, buildExitCode(buildErr))
	}
	table, err := toon.Table("worktree_build", []string{"worktree", "executable"}, [][]string{{assignment.ID, output}})
	if err != nil {
		fmt.Fprintln(stderr, toon.RenderError(err))
		return 1
	}
	fmt.Fprint(stdout, table)
	fmt.Fprintf(stdout, "next[1]:\n%s\n", buildExecNext(assignment))
	return 0
}

// buildVerb is the name every refusal this verb prints carries.
const buildVerb = "bench worktree build"

var buildGrammar = usage.Grammar{
	Cmd:     buildVerb,
	Help:    "usage: " + usage.WorktreeBuild,
	MinArgs: 1,
	MaxArgs: 1,
}

// resolveBuildTarget returns the authorized worktree path and the record the output
// names. resolveAssignment is the whole authority, and the record it answers carries both.
func resolveBuildTarget(root, target string) (string, intent.Assignment, error) {
	assignment, err := resolveAssignment(root, target)
	if err != nil {
		return "", intent.Assignment{}, err
	}
	return assignment.Worktree, assignment, nil
}

// buildExitCode maps an interrupted build apart from a failed one, so a reader can tell
// a signal from a broken tree.
func buildExitCode(err error) int {
	if errors.Is(err, context.Canceled) {
		return 130
	}
	return 1
}

// buildExecNext names the command that runs what the build produced. A label that is not
// line-safe gives way to the assignment id, by the mergeReconcileNext precedent, because
// no quoting makes a control byte pasteable.
func buildExecNext(assignment intent.Assignment) string {
	address := assignment.ID
	if lineSafe(assignment.Label) {
		address = sanitize.ShellQuote(assignment.Label)
	}
	return "  bench worktree exec " + address + " -- ./dist/bench <verb>"
}
