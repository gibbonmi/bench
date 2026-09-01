package gate

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gibbonmi/bench/internal/prose"
	"github.com/gibbonmi/bench/internal/toon"
)

// gateProseUsage carries the grammar line and one example. The root operand is a
// directory, so the single-file form names the file after the `--` separator; the
// example shows that form, so a caller does not learn it by tripping the refusal.
const gateProseUsage = "usage: bench gate-prose <root> [--] [path...]\n" +
	"example: bench gate-prose . -- <path>"

// gateProseFields is the pass table's schema, the shape `roadmap/FT270.md` decides for
// this verb.
var gateProseFields = []string{"path", "verdict"}

// GateProseCommand is the `bench gate-prose <root> [--] [path...]` plumbing command. It
// grades the named paths through the same per-subject grader the whole-tree prose check
// composes, so the lane and the gate agree on one rule. A sole `--help` writes usage to
// stdout and exits 0. Exit 0 is otherwise a clean list, 1 is a list with findings printed
// to stdout, and 2 is a usage error: an unknown flag or an omitted root. A pass states its
// verdict as a `prose[N]{path,verdict}` table, so a caller tells a clean list from a list
// that graded nothing. The word `green` stays out of that table: the lane composes this
// verb, and a lane pass is not a graded green.
func GateProseCommand(args []string, stdout, stderr io.Writer) int {
	if len(args) == 1 && args[0] == "--help" {
		fmt.Fprintln(stdout, gateProseUsage)
		return 0
	}
	root, paths, ok := parseGateProseArgs(args)
	if !ok {
		fmt.Fprintln(stderr, gateProseUsage)
		return 2
	}
	if !rootIsDirectory(root) {
		// The usage text's example carries the single-file form, so this sentence states
		// the refusal alone rather than repeating that form.
		fmt.Fprintf(stderr, "gate-prose: root %q is not a directory: the root operand must be a directory\n", root)
		fmt.Fprintln(stderr, gateProseUsage)
		return 2
	}
	findings := prose.GradeNamedResults(root, paths)
	if len(findings) == 0 {
		rows := make([][]string, 0, len(paths))
		for _, path := range paths {
			rows = append(rows, []string{path, "pass"})
		}
		out, err := toon.Table("prose", gateProseFields, rows)
		if err != nil {
			// A path the encoder cannot carry leaves the verb no honest pass block, so it
			// reports the refusal and exits red rather than forging one.
			fmt.Fprintln(stdout, toon.RenderError(err))
			return 1
		}
		fmt.Fprint(stdout, out)
		return 0
	}
	for _, f := range findings {
		fmt.Fprintln(stdout, prose.RenderNamedResult(f))
	}
	return 1
}

// rootIsDirectory reports whether the root operand names an existing directory. A path
// that does not exist reads as true, so the grader keeps its own diagnostic for a missing
// root; only an existing non-directory is the malformed argument this guard refuses.
func rootIsDirectory(root string) bool {
	info, err := os.Stat(root)
	if err != nil {
		return true
	}
	return info.IsDir()
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
