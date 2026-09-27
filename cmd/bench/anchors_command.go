// The anchors query renders the registered anchors that pin one repo-relative path.
package main

import (
	"os"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/anchors"
	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/canonicalpath"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

var anchorsGrammar = usage.Grammar{
	Cmd:     "bench anchors",
	Help:    "usage: bench anchors <file|dir>",
	MinArgs: 1,
	MaxArgs: 1,
}

func anchorsCommand(args []string) (string, int) {
	parsed, line, code := usage.Parse(anchorsGrammar, args)
	if line != "" {
		return line + "\n", code
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	path, dir := anchorQueryPath(root, parsed.Positionals[0])
	if dir {
		return anchorsDirectory(root, path)
	}
	evaluation := anchors.EvaluatePath(root, path)
	if evaluation.State.Failed() {
		return toon.RecordError(path, evaluation.State, evaluation.Reason) + "\n", 1
	}
	var rows [][]any
	for _, anchor := range evaluation.Locations {
		rows = append(rows, []any{anchorKindName(anchor.Kind), anchor.Section, anchor.Step, anchor.Needle, anchor.Line})
	}
	out, err := toon.TableTyped("anchors", []string{"kind", "section", "step", "needle", "line"}, rows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	for _, diagnostic := range evaluation.Diagnostics {
		out += toon.Errorf("anchor", diagnostic) + "\n"
		code = 1
	}
	help, err := axi.RenderHelp(nil)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	out += help
	return out, code
}

// anchorsDirectory grades each registered file below dir: one row per file with its anchor
// count, its diagnostic count, and its verdict. The registry bounds the rows, and the file
// form of the verb carries each file's own anchor rows and diagnostics.
func anchorsDirectory(root, dir string) (string, int) {
	code := 0
	var rows [][]any
	for _, file := range anchors.FilesBelow(dir) {
		evaluation := anchors.EvaluatePath(root, file)
		verdict := "pass"
		switch {
		case evaluation.State.Failed():
			verdict = "refused"
		case len(evaluation.Diagnostics) > 0:
			verdict = "fail"
		}
		if verdict != "pass" {
			code = 1
		}
		rows = append(rows, []any{file, len(evaluation.Locations), len(evaluation.Diagnostics), verdict})
	}
	out, err := toon.TableTyped("files", []string{"file", "anchors", "diagnostics", "verdict"}, rows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	var actions []axi.Action
	if code != 0 {
		actions = append(actions, axi.ExecutableInvocation("read the anchor rows and diagnostics of a file that did not pass", axi.KnownArgument("anchors"), axi.FutureInput("file")))
	}
	help, err := axi.RenderHelp(actions)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return out + help, code
}

// anchorQueryPath spells the operand the way the anchor registry keys its files, and it
// reports whether the operand is a directory. A path that does not exist has no file to
// key, so the query keeps the operand as typed. A link is never a directory here, so the
// file form refuses it.
func anchorQueryPath(root, arg string) (string, bool) {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = ""
	}
	candidate, display := canonicalpath.Operand(root, cwd, arg)
	if info, err := os.Lstat(candidate); err == nil {
		return display, info.IsDir()
	}
	return filepath.ToSlash(filepath.Clean(arg)), false
}

func anchorKindName(kind anchors.Kind) string {
	switch kind {
	case anchors.Require:
		return "require"
	case anchors.Forbid:
		return "forbid"
	case anchors.RequireInSection:
		return "require-in-section"
	case anchors.ForbidInSection:
		return "forbid-in-section"
	case anchors.ForbidCaseFoldedEmphasis:
		return "forbid-case-folded-emphasis"
	case anchors.RequireInStep:
		return "require-in-step"
	default:
		return "unknown"
	}
}
