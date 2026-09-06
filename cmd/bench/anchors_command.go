// The anchors query renders the registered anchors that pin one repo-relative path.
package main

import (
	"os"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/anchors"
	"github.com/gibbonmi/bench/internal/axi"
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
	var rows [][]string
	for _, anchor := range anchors.Entries() {
		if anchor.File == path {
			rows = append(rows, []string{anchorKindName(anchor.Kind), anchor.Section, anchor.Needle})
		}
	}
	out, err := toon.Table("anchors", []string{"kind", "section", "needle"}, rows)
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

func anchorQueryPath(root, arg string) string {
	candidate := arg
	if !filepath.IsAbs(candidate) {
		if cwd, err := os.Getwd(); err == nil {
			candidate = filepath.Join(cwd, candidate)
		}
	}
	if _, err := os.Lstat(candidate); err == nil {
		if relative, err := filepath.Rel(root, candidate); err == nil {
			return filepath.ToSlash(filepath.Clean(relative))
		}
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
