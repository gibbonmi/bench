// Package commit owns the public command grammar and adapts it to exact landing.
package commit

import (
	"bytes"
	"cmp"
	"context"
	"errors"
	"fmt"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/otelrecord"
	"github.com/gibbonmi/bench/internal/poolkey"
	"github.com/gibbonmi/bench/internal/preflight/evidencecmd"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
	"go.opentelemetry.io/otel/attribute"
)

// The commit seam's record. The verb boundary builds the provider once the repository is
// known and closes the span with the verb's own exit, the way the gate's boundary does.
const otelCommitSeam = "commit"

// commitMeasures is what the record states about the work one commit did: the subject
// the commit published, and the size of the write set it composed. A refusal reaches
// neither, so each carries its own presence.
type commitMeasures struct {
	subject   string
	pathCount int
	counted   bool
}

// beginCommitSpan starts the commit's span and returns the closer that ends it. The span
// carries the subject digest and never the commit subject text: the message holds
// objective text by design, and DATA_HANDLING.md keeps it out of a third durable place.
func beginCommitSpan(root string) func(int, commitMeasures) {
	_, span, finish := otelrecord.Begin("", root, otelCommitSeam)
	if id, ok := poolkey.SplitAssignmentSegment(filepath.Base(root)); ok {
		span.SetAttributes(attribute.String(otelrecord.AttrAssignmentID, id))
	}
	return func(exit int, measures commitMeasures) {
		if measures.subject != "" {
			span.SetAttributes(attribute.String(otelrecord.AttrSubjectID, measures.subject))
		}
		// A refusal before attribution composed no path set, and a zero would read as an
		// empty commit rather than as a commit that never got that far.
		if measures.counted {
			span.SetAttributes(attribute.String(otelrecord.AttrMeasurePathCount, strconv.Itoa(measures.pathCount)))
		}
		span.SetAttributes(attribute.String(otelrecord.AttrOutcome, otelrecord.PublishedExitOutcome(exit)))
		finish()
	}
}

// Command runs a path-attributed prospective landing. Help exits 0, grammar errors exit
// 2, operational refusals exit 1, and a commit that published without reconciling its
// checkout exits 3; the landing owner alone composes, authorizes, and publishes the
// prospective tree.
func Command(args []string, stdout, stderr io.Writer) int {
	_, exit := Run(args, stdout, stderr)
	return exit
}

// Outcome is what one commit call hands the command layer: the checkout the commit ran
// in, the commit it published, and the spec slug that --preflight-build named. A call
// that published nothing leaves Published empty.
type Outcome struct {
	Root, Published, PreflightBuild string
}

// Run is Command with its Outcome, so the commit chain reads the published commit from
// the call that published it.
func Run(args []string, stdout, stderr io.Writer) (Outcome, int) {
	req, help, usageErr := parseRequest(args)
	if help != "" {
		fmt.Fprintln(stdout, helpText)
		return Outcome{}, 0
	}
	if usageErr != "" {
		if strings.HasPrefix(usageErr, "usage: ") {
			fmt.Fprintln(stderr, usageErr)
			return Outcome{}, 2
		}
		fmt.Fprintln(stderr, grammar.Help+" ("+usageErr+")")
		return Outcome{}, 2
	}
	outcome := Outcome{PreflightBuild: req.preflightBuild}
	root, err := git.Root()
	if err != nil {
		fmt.Fprintln(stderr, toon.NotInRepo())
		return outcome, 1
	}
	outcome.Root = root
	// The seam record opens once the repository is known, because the record is addressed
	// by repository. A grammar answer reaches no repository and records nothing.
	var measures commitMeasures
	finishSpan := beginCommitSpan(root)
	exit := commitAttributed(&measures, root, req.msg, req.paths, req.dryRun, stdout, stderr)
	finishSpan(exit, measures)
	outcome.Published = measures.subject
	return outcome, exit
}

// ChainSteps are the steps of `bench commit --preflight-build`, each with its owner's own
// signature. The command layer binds them, so this package imports neither the worktree
// build nor the preflight command.
type ChainSteps struct {
	Commit    func(args []string, stdout, stderr io.Writer) (Outcome, int)
	Build     func(root, home string, args []string, stdout, stderr io.Writer) int
	Home      func() string
	Preflight func(args []string) (string, int)
}

// Chain runs the commit and, when --preflight-build names a slug, the worktree build and
// the build preflight at the published commit. Each step prints its own response, then
// one commit-chain line states every step. The exit is the first non-zero step exit.
func Chain(steps ChainSteps, args []string, stdout, stderr io.Writer) int {
	outcome, exit := steps.Commit(args, stdout, stderr)
	if outcome.PreflightBuild == "" {
		return exit
	}
	build, preflight := "skipped", "skipped"
	// Exit 3 published a commit that the checkout does not match, so a build would grade
	// an unreconciled tree. It stops the chain the way a refusal does.
	if exit == 0 {
		exit = steps.Build(outcome.Root, steps.Home(), []string{outcome.Root}, stdout, stderr)
		build = stepState[exit == 0]
		if exit == 0 {
			out, code := steps.Preflight([]string{evidencecmd.ModeBuild, outcome.PreflightBuild, evidencecmd.FlagTip, outcome.Published})
			fmt.Fprint(stdout, out)
			exit, preflight = code, stepState[code == 0]
		}
	}
	fmt.Fprintf(stdout, "commit-chain{commit=%s,build=%s,preflight=%s}\n", cmp.Or(outcome.Published, "none"), build, preflight)
	return exit
}

// stepState names a chained step that ran by whether it exited 0.
var stepState = map[bool]string{true: "green", false: "red"}

// commitAttributed is the verb's own work, with the span's measures written to measures
// as each becomes known. The exit code it returns is the verb's, so the record and the
// shell agree about the same commit.
func commitAttributed(measures *commitMeasures, root, msg string, paths []string, dryRun bool, stdout, stderr io.Writer) int {
	primary, err := git.IsPrimaryCheckout(root)
	if err != nil {
		fmt.Fprintln(stderr, toon.Errorf("checkout identity is unknown", "repair Git metadata, then retry from a Bench worktree"))
		return 1
	}
	if primary {
		fmt.Fprintln(stderr, usage.PrimaryCheckoutRefusal())
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

	named, err := landing.ResolveAttributedPaths(root, strings.TrimSpace(string(expectedBytes)), paths)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	// The composed path count is the attributed set the commit publishes. It is counted
	// where the list already exists, so no second counter derives the same fact.
	measures.pathCount = len(named)
	measures.counted = true
	// The lane is resolved before anything is graded. A declared lane replaces the
	// whole-project gate for this commit; a malformed declaration refuses the run and
	// names the defect, because a lane nobody can read grades nothing.
	lane, laneErr := gate.LaneForCommit(root)
	if laneErr != nil {
		fmt.Fprintf(stderr, "error: %v\n", laneErr)
		return 1
	}
	if !dryRun {
		formatted, formatErr := formatNamedGoFiles(root, named)
		if formatErr != nil {
			fmt.Fprintf(stderr, "error: format named Go files: %v\n", formatErr)
			return 1
		}
		if len(formatted) > 0 {
			shown := make([]string, len(formatted))
			for i, path := range formatted {
				shown[i] = sanitize.Controls(path)
			}
			fmt.Fprintf(stdout, "formatted Go paths: %s\n", strings.Join(shown, " "))
		}
	}
	owner := landing.NewForLane(lane, strings.TrimSpace(string(expectedBytes)))
	if dryRun {
		if err := owner.DryRun(context.Background(), landing.Request{
			Root: root, Destination: destination, Expected: strings.TrimSpace(string(expectedBytes)),
			Message: msg, Paths: named, Stdout: stdout, Stderr: stderr,
		}); err != nil {
			fmt.Fprintf(stderr, "error: %v\n", err)
			return 1
		}
		// A lane pass is not green, and the lane already stated its own outcome, so the
		// summary borrows neither the word nor a second verdict.
		if lane != nil {
			fmt.Fprintf(stdout, "dry run: composed %d path(s); nothing committed\n", len(named))
		} else {
			fmt.Fprintf(stdout, "dry run: composed %d path(s) authorized green; nothing committed\n", len(named))
		}
		return 0
	}
	result, err := owner.Land(context.Background(), landing.Request{
		Root: root, Destination: destination, Expected: strings.TrimSpace(string(expectedBytes)),
		Message: msg, Paths: named, Stdout: stdout, Stderr: stderr,
	})
	// A published commit identifies the subject, whether or not the checkout reconciled
	// after it. A refusal published nothing and names no subject.
	measures.subject = result.Commit
	if err != nil {
		var remainder *landing.PublishedUnreconciledError
		if errors.As(err, &remainder) {
			return publicationRemainder(stdout, remainder)
		}
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "committed %d path(s)\n", len(named))
	return 0
}

// publicationRemainder reports the publication boundary the landing owner reached: the
// commit exists and the checkout does not match it. The record uses the landing verb's
// name{key=value,...} grammar, and its exit code separates this outcome from a refusal
// that published nothing.
func publicationRemainder(stdout io.Writer, remainder *landing.PublishedUnreconciledError) int {
	fmt.Fprintf(stdout, "committed{published_commit=%s,path=%s,next=%s}\n",
		remainder.Commit, sanitize.Controls(remainder.Path), restoreNext(remainder.Commit, remainder.Paths))
	return 3
}

// restoreNext names the one restore that reconciles every named path against the
// published commit. The restore is idempotent, so it covers the paths that already
// reconciled as well as the remainder. A path that is not line-safe takes the landing
// verb's pointer form: quoting would still emit the raw byte into a line-structured
// record, and escaping would name a path that does not exist.
//
// The value is line-safe by construction, so it reaches the record unescaped; the
// sanitizer's backslash escaping would break the quoting a reader pastes.
func restoreNext(commit string, paths []string) string {
	command := "git restore --source=" + commit + " --staged --worktree --"
	quoted := make([]string, 0, len(paths))
	for _, path := range paths {
		if !sanitize.LineSafe(path) {
			return command + " <named-paths>"
		}
		quoted = append(quoted, sanitize.ShellQuote(path))
	}
	return command + " " + strings.Join(quoted, " ")
}

// LaneClause names what a commit grades: the declared lane, or the gate when the
// project declares no lane. The `bench help` row and the --dry-run help line read it.
const LaneClause = "run the declared lane (or the gate when no lane is declared)"

// helpText adds one concrete example and the exit-code meanings the grammar line
// cannot carry. A usage error prints the grammar line alone, so only a help request
// pays for them. The example shows the trailing `-- <path>...` form, so a caller does
// not learn the argument shape by tripping the usage line.
var helpText = grammar.Help + "\n" +
	"example: bench commit -m \"fix: tighten the guard\" -- internal/gitguard/scan.go docs/adr/0007.md\n" +
	"--dry-run: " + LaneClause + " on the exact composed snapshot and report the outcome; commit nothing\n" +
	preflightBuildFlag + " <slug>: after a commit that publishes, run bench worktree build on this worktree, then bench preflight build <slug> at the published commit; end with commit-chain{commit,build,preflight} and exit with the first non-zero step exit\n" +
	"exit 1: refused before publication; nothing was committed\n" +
	"exit 2: grammar error\n" +
	"exit 3: published; the checkout did not reconcile — paste next= to repair"

var grammar = usage.Grammar{
	Cmd:  "bench commit",
	Help: "usage: bench commit [--dry-run | " + preflightBuildFlag + " <slug>] -m <msg> [--] <path>...",
	Flags: []usage.Flag{{Name: "-m", HasValue: true}, {Name: "--dry-run"},
		{Name: preflightBuildFlag, HasValue: true, NoEmptyValue: true}},
	MaxArgs: -1,
}

const preflightBuildFlag = "--preflight-build"

// request is one parsed commit call.
type request struct {
	msg            string
	paths          []string
	dryRun         bool
	preflightBuild string
}

func parseRequest(args []string) (req request, help string, usageErr string) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		if code == 0 {
			return request{}, line, ""
		}
		return request{}, "", line
	}
	_, req.dryRun = parsed.Flags["--dry-run"]
	req.preflightBuild = parsed.Flags[preflightBuildFlag]
	// A dry run publishes nothing, so the chain would have no commit to build.
	if req.dryRun && req.preflightBuild != "" {
		return request{}, "", "--dry-run excludes " + preflightBuildFlag
	}
	msg, msgSet := parsed.Flags["-m"]
	if !msgSet {
		return request{}, "", "-m <msg> is required"
	}
	if strings.TrimSpace(msg) == "" {
		return request{}, "", "-m <msg> must not be empty"
	}
	if len(parsed.Positionals) == 0 {
		return request{}, "", "at least one <path> is required"
	}
	req.msg, req.paths = msg, parsed.Positionals
	return req, "", ""
}

func formatNamedGoFiles(root string, named []string) ([]string, error) {
	args := []string{"-C", root, "status", "--porcelain=v1", "-z", "--no-renames", "--untracked-files=all", "--"}
	for _, path := range named {
		args = append(args, ":(literal)"+path)
	}
	raw, err := git.Raw(args...)
	if err != nil {
		return nil, err
	}
	entries, err := git.ParsePorcelainZStrict(raw)
	if err != nil {
		return nil, err
	}
	type edit struct {
		path string
		body []byte
		mode os.FileMode
	}
	seen := map[string]bool{}
	var edits []edit
	for _, entry := range entries {
		path := entry.Path
		if seen[path] || !strings.HasSuffix(path, ".go") {
			continue
		}
		seen[path] = true
		full := filepath.Join(root, filepath.FromSlash(path))
		info, statErr := os.Lstat(full)
		if errors.Is(statErr, os.ErrNotExist) {
			continue
		}
		if statErr != nil {
			return nil, fmt.Errorf("inspect %q: %w", path, statErr)
		}
		if !info.Mode().IsRegular() {
			continue
		}
		body, readErr := os.ReadFile(full)
		if readErr != nil {
			return nil, fmt.Errorf("read %q: %w", path, readErr)
		}
		formatted, formatErr := format.Source(body)
		if formatErr != nil {
			return nil, fmt.Errorf("%q: %w", path, formatErr)
		}
		if !bytes.Equal(body, formatted) {
			edits = append(edits, edit{path: path, body: formatted, mode: info.Mode().Perm()})
		}
	}
	sort.Slice(edits, func(i, j int) bool { return edits[i].path < edits[j].path })
	formatted := make([]string, 0, len(edits))
	for _, edit := range edits {
		if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(edit.path)), edit.body, edit.mode); err != nil {
			return formatted, fmt.Errorf("write %q: %w", edit.path, err)
		}
		formatted = append(formatted, edit.path)
	}
	return formatted, nil
}
