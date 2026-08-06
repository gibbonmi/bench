// Package commit owns the public command grammar and adapts it to exact landing.
package commit

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// Command runs a path-attributed prospective landing. Help exits 0, grammar errors exit
// 2, and operational refusals exit 1; the landing owner alone composes, authorizes, and
// publishes the prospective tree.
func Command(args []string, stdout, stderr io.Writer) int {
	msg, specSlug, paths, help, usageErr := parseArgs(args)
	if help != "" {
		fmt.Fprintln(stdout, help)
		return 0
	}
	if usageErr != "" {
		fmt.Fprintln(stderr, grammar.Help+" (--spec marks the spec implemented; "+usageErr+")")
		return 2
	}
	root, err := git.Root()
	if err != nil {
		fmt.Fprintln(stderr, toon.NotInRepo())
		return 1
	}

	// Capture publication identity before reading attributed content. A detached checkout
	// updates literal HEAD; an attached checkout updates its full branch ref.
	destination := "HEAD"
	if out, symbolicErr := git.Raw("-C", root, "symbolic-ref", "-q", "HEAD"); symbolicErr == nil {
		destination = strings.TrimSpace(string(out))
	}
	expectedBytes, expectedErr := git.Raw("-C", root, "rev-parse", "--verify", "HEAD^{commit}")
	if expectedErr != nil {
		fmt.Fprintln(stderr, "error: destination has no commit base")
		return 1
	}

	named := make([]string, 0, len(paths))
	for _, path := range paths {
		rel, relErr := rootRel(root, path)
		if relErr != nil {
			fmt.Fprintf(stderr, "error: cannot resolve path %q relative to repo root: %v\n", path, relErr)
			return 1
		}
		named = append(named, rel)
	}
	if _, err := landing.New().Land(context.Background(), landing.Request{
		Root: root, Destination: destination, Expected: strings.TrimSpace(string(expectedBytes)),
		Message: msg, Paths: named, Spec: specSlug, Stdout: stdout, Stderr: stderr,
	}); err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "committed %d path(s)\n", len(named))
	return 0
}

var grammar = usage.Grammar{
	Cmd:     "bench commit",
	Help:    "usage: bench commit -m <msg> [--spec <slug>] [--] <path>...",
	Flags:   []usage.Flag{{Name: "-m", HasValue: true}, {Name: "--spec", HasValue: true}},
	MaxArgs: -1,
}

func parseArgs(args []string) (msg string, specSlug string, paths []string, help string, usageErr string) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		if code == 0 {
			return "", "", nil, line, ""
		}
		return "", "", nil, "", line
	}
	msg, msgSet := parsed.Flags["-m"]
	if !msgSet {
		return "", "", nil, "", "-m <msg> is required"
	}
	if strings.TrimSpace(msg) == "" {
		return "", "", nil, "", "-m <msg> must not be empty"
	}
	if len(parsed.Positionals) == 0 {
		return "", "", nil, "", "at least one <path> is required"
	}
	return msg, parsed.Flags["--spec"], parsed.Positionals, "", ""
}

func rootRel(root, arg string) (string, error) {
	abs, err := filepath.Abs(arg)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(root, abs)
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(rel), nil
}
