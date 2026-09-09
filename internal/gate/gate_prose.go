package gate

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/prose"
	"github.com/gibbonmi/bench/internal/toon"
)

// gateProseUsage carries the grammar line and one example per form. The root operand is a
// directory, so the single-file form names the file after the `--` separator; the examples
// show both forms, so a caller does not learn either by tripping the refusal.
const gateProseUsage = "usage: bench gate-prose <root> (--staged | [--] [path...])\n" +
	"example: bench gate-prose . -- <path>\n" +
	"example: bench gate-prose . --staged"

// gateProseFields is the pass table's schema, the shape `roadmap/FT270.md` decides for
// this verb.
var gateProseFields = []string{"path", "verdict"}

// GateProseCommand is the `bench gate-prose <root> (--staged | [--] [path...])` plumbing
// command. It grades the named paths through the same per-subject grader the whole-tree
// prose check composes, so the lane and the gate agree on one rule. `--staged` selects the
// staged Markdown instead of a path list and grades the index bytes, so a pre-commit
// caller grades what the commit would carry. A sole `--help` writes usage to
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
	root, paths, staged, ok := parseGateProseArgs(args)
	if !ok {
		fmt.Fprintln(stderr, gateProseUsage)
		return 2
	}
	if staged {
		return gateProseStaged(root, stdout)
	}
	if !rootIsDirectory(root) {
		// The usage text's example carries the single-file form, so this sentence states
		// the refusal alone rather than repeating that form.
		fmt.Fprintf(stderr, "gate-prose: root %q is not a directory: the root operand must be a directory\n", root)
		fmt.Fprintln(stderr, gateProseUsage)
		return 2
	}
	return renderProseVerdict(paths, prose.GradeNamedResults(root, paths), stdout)
}

// renderProseVerdict states the verdict both forms owe: the pass table over the graded
// subjects on a clean list, and one diagnostic line per finding otherwise. The two forms
// select their subjects differently and report them the same way.
func renderProseVerdict(subjects []string, findings []prose.NamedResult, stdout io.Writer) int {
	if len(findings) > 0 {
		for _, f := range findings {
			fmt.Fprintln(stdout, prose.RenderNamedResult(f))
		}
		return 1
	}
	rows := make([][]string, 0, len(subjects))
	for _, path := range subjects {
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

// gateProseStaged grades the Markdown the index holds rather than the working files. The
// subjects are the staged regular-file `.md` entries, the bytes come from each index blob,
// and the exclusion policy comes from the index blob of the policy file. So the verb
// grades what a commit would carry, which is what a pre-commit caller asks about.
//
// The refusal for a root that is not a working-tree top reaches stdout at exit 1 rather
// than stderr at exit 2, because the operand is well formed and only the repository state
// refuses it.
func gateProseStaged(root string, stdout io.Writer) int {
	index, err := git.ReadStagedIndex(root)
	if err != nil {
		fmt.Fprintf(stdout, "prose: %q: refused staged grade: %s\n", root, err)
		return 1
	}
	entries := make([]string, 0, len(index.Entries))
	policy := prose.IndexSource{}
	for _, entry := range index.Entries {
		entries = append(entries, entry.Path)
		if entry.Path == prose.ExclusionFile && entry.IsRegularFile() {
			policy.PolicyPresent = true
		}
	}
	policy.Entries = entries
	if policy.PolicyPresent {
		body, err := git.IndexBlob(root, prose.ExclusionFile)
		if err != nil {
			fmt.Fprintf(stdout, "prose: %q: refused unreadable exclusion file: %s\n", prose.ExclusionFile, err)
			return 1
		}
		policy.Policy = body
	}
	grader, diags := prose.NewGraderFromIndex(policy)
	if len(diags) > 0 {
		for _, diagnostic := range diags {
			fmt.Fprintln(stdout, diagnostic)
		}
		return 1
	}
	var subjects []string
	var findings []prose.NamedResult
	for _, entry := range index.Staged {
		if !entry.IsRegularFile() || !strings.HasSuffix(entry.Path, ".md") {
			continue
		}
		subjects = append(subjects, entry.Path)
		body, err := git.IndexBlob(root, entry.Path)
		if err != nil {
			fmt.Fprintf(stdout, "prose: %q: refused unreadable subject: %s\n", entry.Path, err)
			return 1
		}
		findings = append(findings, grader.GradeBytes(entry.Path, body)...)
	}
	return renderProseVerdict(subjects, findings, stdout)
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

// parseGateProseArgs splits args into the root, the named path list, and the staged form.
// `--staged` is the verb's one flag, and it selects its subjects itself: a path list or a
// `--` separator beside it names two subject sources at once and is a usage error. Any
// other argument that starts with `-` before a `--` separator is an unrecognized flag and
// a usage error. An empty path list is valid: it grades nothing and passes.
func parseGateProseArgs(args []string) (root string, paths []string, staged, ok bool) {
	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		return "", nil, false, false
	}
	root = args[0]
	sawSep := false
	for _, a := range args[1:] {
		if a == "--staged" {
			staged = true
			continue
		}
		if !sawSep && a == "--" {
			sawSep = true
			continue
		}
		if !sawSep && strings.HasPrefix(a, "-") {
			return "", nil, false, false
		}
		paths = append(paths, a)
	}
	if staged && (sawSep || len(paths) > 0) {
		return "", nil, false, false
	}
	return root, paths, staged, true
}
