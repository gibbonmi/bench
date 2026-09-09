// The anchors query renders the registered anchors that pin one repo-relative path.
package main

import (
	"os"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/anchors"
	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/canonicalpath"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

var anchorsGrammar = usage.Grammar{
	Cmd:     "bench anchors",
	Help:    "usage: bench anchors <path>",
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
	path := anchorQueryPath(root, parsed.Positionals[0])
	// Classification precedes the read: a link, a special file, or an unreadable file
	// answers a structured refusal instead of blocking or lying with a row of zeros.
	classified := bounds.ClassifyNoFollow(filepath.Join(root, filepath.FromSlash(path)))
	if classified.State.Failed() {
		return toon.RecordError(path, classified.State, classified.Reason) + "\n", 1
	}
	data := string(classified.Data)
	var rows [][]any
	for _, anchor := range anchors.Entries() {
		if anchor.File == path {
			rows = append(rows, []any{anchorKindName(anchor.Kind), anchor.Section, anchor.Needle, anchors.Locate(anchor.Kind, anchor.Section, anchor.Needle, data)})
		}
	}
	out, err := toon.TableTyped("anchors", []string{"kind", "section", "needle", "line"}, rows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	help, err := axi.RenderHelp(nil)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	out += help
	return out, 0
}

// anchorQueryPath spells the operand the way the anchor registry keys its files. A path
// that does not exist has no file to key, so the query keeps the operand as typed.
func anchorQueryPath(root, arg string) string {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = ""
	}
	candidate, display := canonicalpath.Operand(root, cwd, arg)
	if _, err := os.Lstat(candidate); err == nil {
		return display
	}
	return filepath.ToSlash(filepath.Clean(arg))
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
	default:
		return "unknown"
	}
}
