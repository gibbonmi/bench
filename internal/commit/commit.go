// Package commit owns the public command grammar and adapts it to exact landing.
package commit

import (
	"bytes"
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

	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/otelrecord"
	"github.com/gibbonmi/bench/internal/poolkey"
	"github.com/gibbonmi/bench/internal/refusalroute"
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

// Outcome is what one commit call hands the command layer: the checkout the commit ran
// in, the commit it published, the spec slug that --preflight-build named, and whether
// the call answered a help request. A call that published nothing leaves Published empty.
type Outcome struct {
	Root, Published, PreflightBuild string
	Help                            bool
}

// Run runs a path-attributed prospective landing and returns its Outcome, so the commit
// chain reads the published commit from the call that published it. Help exits 0,
// grammar errors exit 2, operational refusals exit 1, and a commit that published without
// reconciling its checkout exits 3; the landing owner alone composes, authorizes, and
// publishes the prospective tree.
func Run(args []string, stdout, stderr io.Writer) (Outcome, int) {
	req, help, usageErr := parseRequest(args)
	if help != "" {
		fmt.Fprintln(stdout, helpText)
		return Outcome{Help: true}, 0
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
	exit := commitAttributed(&measures, root, req, stdout, stderr)
	finishSpan(exit, measures)
	outcome.Published = measures.subject
	return outcome, exit
}

// commitAttributed is the verb's own work, with the span's measures written to measures
// as each becomes known. The exit code it returns is the verb's, so the record and the
// shell agree about the same commit.
func commitAttributed(measures *commitMeasures, root string, req request, stdout, stderr io.Writer) int {
	// handback refuses a cause that no face of its own names.
	handback := func(err error) int { return refuse(stderr, root, req, faceHandback, err.Error()) }
	primary, err := git.IsPrimaryCheckout(root)
	if err != nil {
		return handback(errors.New("checkout identity is unknown"))
	}
	if primary {
		fmt.Fprintln(stderr, usage.PrimaryCheckoutRefusal())
		return refuse(stderr, root, req, facePrimaryCheckout, "")
	}

	// Capture publication identity before reading attributed content. A detached checkout
	// updates literal HEAD; an attached checkout updates its full branch ref.
	destination := "HEAD"
	if out, symbolicErr := git.Raw("-C", root, "symbolic-ref", "-q", "HEAD"); symbolicErr == nil {
		destination = strings.TrimSpace(string(out))
	}
	expectedBytes, expectedErr := git.Raw("-C", root, "rev-parse", "--verify", "HEAD^{commit}")
	if expectedErr != nil {
		return handback(errors.New("destination has no commit base"))
	}

	named, err := landing.ResolveAttributedPaths(root, strings.TrimSpace(string(expectedBytes)), req.paths)
	if err != nil {
		return handback(err)
	}
	candidate, err := landing.CandidateTree(root, strings.TrimSpace(string(expectedBytes)), named)
	if err == nil {
		err = (commitrepo.Store{Root: root}).AuthorizeCandidate(candidate)
	}
	if err != nil {
		fmt.Fprintf(stderr, "error: commitment: %v\n", err)
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
		return handback(laneErr)
	}
	if !req.dryRun {
		formatted, formatErr := formatNamedGoFiles(root, named)
		if formatErr != nil {
			return handback(fmt.Errorf("format named Go files: %w", formatErr))
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
	landingRequest := landing.Request{
		Root: root, Destination: destination, Expected: strings.TrimSpace(string(expectedBytes)),
		Message: req.msg, Paths: named, Stdout: stdout, Stderr: stderr,
	}
	if req.dryRun {
		if err := owner.DryRun(context.Background(), landingRequest); err != nil {
			return refuse(stderr, root, req, landingFace(err), err.Error())
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
	result, err := owner.Land(context.Background(), landingRequest)
	// A published commit identifies the subject, whether or not the checkout reconciled
	// after it. A refusal published nothing and names no subject.
	measures.subject = result.Commit
	if err != nil {
		var remainder *landing.PublishedUnreconciledError
		if errors.As(err, &remainder) {
			return publicationRemainder(stdout, root, remainder)
		}
		return refuse(stderr, root, req, landingFace(err), err.Error())
	}
	fmt.Fprintf(stdout, "committed %d path(s)\n", len(named))
	return 0
}

// The commit's refusal faces. The shared registry declares each face, and a raising site
// names the face it raises.
const (
	facePublishedUnreconciled = "commit-published-unreconciled"
	facePrimaryCheckout       = "commit-primary-checkout"
	faceRed                   = "commit-red"
	faceInfrastructure        = "commit-infrastructure"
	faceHandback              = "commit-handback"
)

// landingFace picks the face of a refusal that the landing owner returned. The kind that
// the gate attributed picks the face of an authorization refusal, and every other refusal
// hands back.
func landingFace(err error) string {
	var refused landing.AuthorizationRefusal
	switch {
	case !errors.As(err, &refused):
	case refused.Result.Kind == authorization.Infrastructure:
		return faceInfrastructure
	case landing.RedKind(refused.Result.Kind):
		return faceRed
	}
	return faceHandback
}

// refuse prints a refusal that published nothing, sentence and then the route of face on
// its own next= line, and answers exit 1. The route re-runs the caller's own commit in the
// worktree of the active assignment that owns root. With no owner, the zero owner's empty
// label prints the label slot.
func refuse(stderr io.Writer, root string, req request, face, sentence string) int {
	owner, _ := intent.AssignmentForWorktree(root)
	args := []string{}
	if req.dryRun {
		args = append(args, "--dry-run")
	}
	if req.preflightBuild != "" {
		args = append(args, PreflightBuildFlag, refusalroute.Arg("slug", req.preflightBuild))
	}
	args = append(args, "-m", refusalroute.Arg("msg", req.msg), "--")
	for _, path := range req.paths {
		args = append(args, refusalroute.Arg("path", path))
	}
	refusal := refusalroute.New(face, refusalroute.Facts{Sentence: sentence, Values: map[string]string{
		refusalroute.FactLabel: owner.Label, refusalroute.FactArguments: strings.Join(args, " "),
	}})
	if refusal.Sentence != "" {
		fmt.Fprintln(stderr, "error: "+refusal.Sentence)
	}
	fmt.Fprintln(stderr, refusalroute.NextField+"="+refusal.Route)
	return 1
}

// publicationRemainder reports the publication boundary the landing owner reached: the
// commit exists and the checkout at root does not match it. The record uses the landing
// verb's name{key=value,...} grammar, and its exit code separates this outcome from a
// refusal that published nothing.
//
// The route is line-safe by construction, so it reaches the record unescaped; the
// sanitizer's backslash escaping would break the quoting a reader pastes.
func publicationRemainder(stdout io.Writer, root string, remainder *landing.PublishedUnreconciledError) int {
	next := refusalroute.New(facePublishedUnreconciled, refusalroute.Facts{Values: map[string]string{
		refusalroute.FactPublishedCommit: remainder.Commit, refusalroute.FactCheckout: root,
	}}).Route
	fmt.Fprintf(stdout, "committed{published_commit=%s,path=%s,%s=%s}\n",
		remainder.Commit, sanitize.Controls(remainder.Path), refusalroute.NextField, next)
	return 3
}

// LaneClause names what a commit grades: the declared lane, or the gate when the
// project declares no lane. The `bench help` row and the --dry-run help line read it.
const LaneClause = "run the declared lane (or the gate when no lane is declared)"

// helpText adds one concrete example, the deletion route, and the exit-code meanings
// the grammar line cannot carry. A usage error prints the grammar line alone, so only a
// help request pays for them. The example shows the trailing `-- <path>...` form, so a
// caller does not learn the argument shape by tripping the usage line. The command
// layer that runs the chain adds the --preflight-build line.
var helpText = grammar.Help + "\n" +
	"example: bench commit -m \"fix: tighten the guard\" -- internal/gitguard/scan.go docs/adr/0007.md\n" +
	"deleted path: name a file or folder deleted from the worktree and the commit publishes its deletion; no git rm is needed\n" +
	"--dry-run: " + LaneClause + " on the exact composed snapshot and report the outcome; commit nothing\n" +
	"exit 1: refused before publication; nothing was committed\n" +
	"exit 2: grammar error\n" +
	"exit 3: published; the checkout did not reconcile — paste " + refusalroute.NextField + "= to repair"

var grammar = usage.Grammar{
	Cmd:  "bench commit",
	Help: "usage: bench commit [--dry-run | " + PreflightBuildFlag + " <slug>] -m <msg> [--] <path>...",
	Flags: []usage.Flag{{Name: "-m", HasValue: true}, {Name: "--dry-run"},
		{Name: PreflightBuildFlag, HasValue: true, NoEmptyValue: true}},
	MaxArgs: -1,
}

// PreflightBuildFlag names the spec slug whose build the command layer chains after a
// commit that publishes.
const PreflightBuildFlag = "--preflight-build"

// HelpRowSuffix is the argument shape of the `bench help` row.
const HelpRowSuffix = " -m <msg> [" + PreflightBuildFlag + " <slug>] <path>..."

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
	req.preflightBuild = parsed.Flags[PreflightBuildFlag]
	// A dry run publishes nothing, so the chain would have no commit to build.
	if req.dryRun && req.preflightBuild != "" {
		return request{}, "", "--dry-run excludes " + PreflightBuildFlag
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
