package gate

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/prose"
	"github.com/gibbonmi/bench/internal/toon"
)

const gateProseUsage = "usage: bench gate-prose <root> [--] [path...]"

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
	findings := prose.GradeNamed(root, paths)
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
		fmt.Fprintln(stdout, proseFindingWithSentence(root, paths, f))
	}
	return 1
}

func proseFindingWithSentence(root string, paths []string, finding string) string {
	if !strings.Contains(finding, ": sentence of ") {
		return finding
	}
	for _, path := range paths {
		prefix := fmt.Sprintf("prose: %q line ", path)
		rest, ok := strings.CutPrefix(finding, prefix)
		if !ok {
			continue
		}
		lineText, _, ok := strings.Cut(rest, ":")
		if !ok {
			return finding
		}
		line, err := strconv.Atoi(lineText)
		if err != nil || line < 1 {
			return finding
		}
		body, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil {
			return finding
		}
		lines := strings.Split(string(body), "\n")
		if line > len(lines) {
			return finding
		}
		return finding + ": " + strconv.Quote(strings.TrimSpace(lines[line-1]))
	}
	return finding
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
