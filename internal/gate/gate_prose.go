package gate

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/prose"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/toon"
)

// gateProseUsage carries the grammar line and one example per form. The root operand is a
// directory, so the single-file form names the file after the `--` separator; the examples
// show both forms, so a caller does not learn either by tripping the refusal.
const gateProseUsage = "usage: bench gate-prose <root> (--staged | [--] [path...])\n" +
	"example: bench gate-prose . -- <path>\n" +
	"example: bench gate-prose . --staged"

// gateProseFields is the verdict table's schema.
var gateProseFields = []string{"path", "verdict"}

// The verdict cells of a clean list. A subject that the exclusion policy keeps out of the
// grade reads as excluded, so no row claims a grade that did not run.
const (
	proseVerdictPass     = "pass"
	proseVerdictExcluded = "excluded"
)

// GateProseCommand is the `bench gate-prose <root> (--staged | [--] [path...])` plumbing
// command. It grades the named paths through the same per-subject grader the whole-tree
// prose check composes, so the lane and the gate agree on one rule. `--staged` selects the
// staged Markdown instead of a path list and grades the index bytes, so a pre-commit
// caller grades what the commit would carry. A sole `--help` writes usage to
// stdout and exits 0. Exit 0 is otherwise a clean list, 1 is a list with findings or a
// missing named path printed to stdout, and 2 is a usage error: an unknown flag or an omitted root. A pass states its
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
	if refusedOperands(root, paths, stdout) {
		return 1
	}
	findings, excluded := prose.GradeNamedResults(root, paths)
	return renderProseVerdict(paths, excluded, findings, stdout)
}

// refusedOperands writes one refusal line for each named path the verb cannot find under
// root, and it reports whether it wrote any. The prose grader skips an absent path, so
// without this check a path on the wrong checkout would read as a pass row. The lane names
// only the Markdown the composed tree holds, so a deleted file never reaches this check.
func refusedOperands(root string, paths []string, stdout io.Writer) bool {
	refused := false
	for _, path := range paths {
		_, err := os.Lstat(filepath.Join(root, filepath.FromSlash(path)))
		if err == nil {
			continue
		}
		reason := err
		if inner := errors.Unwrap(err); inner != nil {
			// The path error repeats the joined path; the refusal names the operand itself.
			reason = inner
		}
		fmt.Fprintln(stdout, prose.UnreadableSubjectDiagnostic(path, reason.Error()))
		refused = true
	}
	return refused
}

// renderProseVerdict states the verdict both forms owe: the verdict table over the
// subjects on a clean list, and one diagnostic line per finding otherwise. The excluded
// set is the grader's own decision, so the table marks a subject excluded only where the
// grade skipped it. The two forms select their subjects differently and report them the
// same way.
func renderProseVerdict(subjects []string, excluded map[string]bool, findings []prose.NamedResult, stdout io.Writer) int {
	if len(findings) > 0 {
		for _, f := range findings {
			fmt.Fprintln(stdout, prose.RenderNamedResult(f))
		}
		return 1
	}
	rows := make([][]string, 0, len(subjects))
	for _, path := range subjects {
		verdict := proseVerdictPass
		if excluded[path] {
			verdict = proseVerdictExcluded
		}
		rows = append(rows, []string{path, verdict})
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
// grades what a commit would carry, which is what a pre-commit caller asks about. Every
// blob is bounded and classified before it is graded, so an index holds no bytes the named
// form would refuse and this form would parse.
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
		body, reason := readIndexBlob(root, prose.ExclusionFile)
		if reason != "" {
			fmt.Fprintln(stdout, prose.UnreadableExclusionDiagnostic(reason))
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
	excluded := map[string]bool{}
	for _, entry := range index.Staged {
		if !entry.IsRegularFile() || !prose.IsSubjectName(entry.Path) {
			continue
		}
		subjects = append(subjects, entry.Path)
		if grader.Excludes(entry.Path) {
			excluded[entry.Path] = true
		}
		body, reason := readIndexBlob(root, entry.Path)
		if reason != "" {
			fmt.Fprintln(stdout, prose.UnreadableSubjectDiagnostic(entry.Path, reason))
			return 1
		}
		findings = append(findings, grader.GradeBytes(entry.Path, body)...)
	}
	return renderProseVerdict(subjects, excluded, findings, stdout)
}

// readIndexBlob reads one index blob through classifiedBlob.
func readIndexBlob(root, path string) ([]byte, string) {
	stream, err := git.IndexBlobReader(root, path)
	if err != nil {
		return nil, err.Error()
	}
	return classifiedBlob(stream)
}

// classifiedBlob reads one blob stream under the control-record bound, closes it, and grades
// the bytes through the classifier the named form's file read uses. A blob over the limit or
// holding invalid UTF-8 is refused rather than parsed, so the staged form, the checkpoint
// record grade, and the named form dispose of the same bytes the same way. The bound applies
// to the stream, so an oversized blob stops the read at the limit rather than reaching memory
// whole. The returned reason is empty for bytes a grade may trust, and it names the fault
// otherwise.
func classifiedBlob(stream io.ReadCloser) ([]byte, string) {
	read := bounds.Read(stream, bounds.ControlRecordLimit)
	if err := stream.Close(); err != nil {
		return nil, err.Error()
	}
	classified := bounds.ClassifyBytes(read)
	if classified.State.Failed() {
		return nil, classified.Reason
	}
	return classified.Data, ""
}

// gradeCheckpointRecord grades the review record prose of the chunk checkpoint that ctx
// carries. The run grades the tree without the record, and a reuse runs nothing, so a caller
// grades the record before its reuse return. It answers whether it refused, and the refusal.
func gradeCheckpointRecord(ctx context.Context, root string, stderr io.Writer, mode runMode) (Result, bool) {
	source, chunk := ctx.Value(chunkSourceKey{}).(chunkSource)
	if !chunk {
		return Result{}, false
	}
	checkpoint, _ := ctx.Value(checkpointKey{}).(Checkpoint)
	findings := recordProse(root, source.tree, checkpoint.Spec)
	if len(findings) == 0 {
		return Result{}, false
	}
	for _, finding := range findings {
		fmt.Fprintln(stderr, finding)
	}
	return refuse(ctx, root, stderr, mode, faceRecordProse, "checkpoint review record prose refused"), true
}

// recordProse answers one diagnostic for each prose finding of the review record of spec in
// tree, as the staged form grades a staged blob. The record blob and the exclusion-policy blob
// come from tree. The prose package holds a tree with no exclusion policy ungradeable
// (prose.NewGraderFromIndex), so the record grade skips such a tree and answers none.
func recordProse(root, tree, spec string) []string {
	record, err := reviewrecord.RecordPath(spec)
	if err != nil {
		return []string{err.Error()}
	}
	generation, err := captureProspectiveTree(root, tree)
	if err != nil {
		return []string{prose.UnreadableSubjectDiagnostic(record, err.Error())}
	}
	if _, present := generation.entry(prose.ExclusionFile); !present {
		return nil
	}
	policy := prose.IndexSource{PolicyPresent: true}
	for _, entry := range generation.snapshot.entries {
		policy.Entries = append(policy.Entries, entry.Path)
	}
	var reason string
	if policy.Policy, reason = readTreeBlob(root, tree, prose.ExclusionFile); reason != "" {
		return []string{prose.UnreadableExclusionDiagnostic(reason)}
	}
	grader, diagnostics := prose.NewGraderFromIndex(policy)
	if len(diagnostics) > 0 {
		return diagnostics
	}
	body, reason := readTreeBlob(root, tree, record)
	if reason != "" {
		return []string{prose.UnreadableSubjectDiagnostic(record, reason)}
	}
	var findings []string
	for _, finding := range grader.GradeBytes(record, body) {
		findings = append(findings, prose.RenderNamedResult(finding))
	}
	return findings
}

// readTreeBlob reads the regular file at path of tree through classifiedBlob.
func readTreeBlob(root, tree, path string) ([]byte, string) {
	data, err := git.ReadTreeFile(root, tree, path)
	if err != nil {
		return nil, err.Error()
	}
	return classifiedBlob(io.NopCloser(bytes.NewReader(data)))
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
