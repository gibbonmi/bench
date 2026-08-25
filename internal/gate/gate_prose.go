package gate

import (
	"fmt"
	"io"
	"strings"

	"github.com/gibbonmi/bench/internal/prose"
)

const gateProseUsage = "usage: bench gate-prose <root> [--] [path...]"

// GateProseCommand is the `bench gate-prose <root> [--] [path...]` plumbing command. It
// grades the named paths through the same per-subject grader the whole-tree prose check
// composes, so the lane and the gate agree on one rule. Exit 0 is a clean list, 1 is a
// list with findings printed to stdout, and 2 is a usage error: an unknown flag or an
// omitted root.
func GateProseCommand(args []string, stdout, stderr io.Writer) int {
	root, paths, ok := parseGateProseArgs(args)
	if !ok {
		fmt.Fprintln(stderr, gateProseUsage)
		return 2
	}
	findings := prose.GradeNamed(root, paths)
	if len(findings) == 0 {
		return 0
	}
	for _, f := range findings {
		fmt.Fprintln(stdout, f)
	}
	return 1
}

// parseGateProseArgs splits args into the root and the named path list. The verb takes
// no flags of its own, so an argument that starts with `-` before a `--` separator is an
// unrecognized flag and a usage error. An empty path list is valid: it grades nothing and
// passes.
func parseGateProseArgs(args []string) (root string, paths []string, ok bool) {
	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		return "", nil, false
	}
	root = args[0]
	sawSep := false
	for _, a := range args[1:] {
		if !sawSep && a == "--" {
			sawSep = true
			continue
		}
		if !sawSep && strings.HasPrefix(a, "-") {
			return "", nil, false
		}
		paths = append(paths, a)
	}
	return root, paths, true
}
