package gate

// The fast lane: the declared check list a worktree commit runs in place of the
// whole-project gate. A lane is an ordered []Phase with the phase manifest's entry
// schema, and it runs through the same scheduler the gate's phase table runs through.
//
// A lane is not a gate run. It writes its own record class into its own file, and it
// touches neither the verdict cache nor the evidence store. So no reader can mistake a
// lane pass for a graded green.

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gocache"
	"github.com/gibbonmi/bench/internal/otelrecord"
	"github.com/gibbonmi/bench/internal/runbinary"
	"go.opentelemetry.io/otel/attribute"
)

// The two lane outcomes. A lane never reports "green": that word belongs to the
// landing's whole-project gate, and a lane pass authorizes nothing the gate decides.
const (
	lanePass = "pass"
	laneFail = "fail"
)

// LaneNamedMarkdownToken is the placeholder a prose check carries for the Markdown
// paths a commit names. RunLane replaces the one argv element holding it with the
// caller's path list, which may be empty. The check is declared, so it still runs on
// an empty list; only its subject follows the named paths.
const LaneNamedMarkdownToken = "<bench-named-markdown>"

// laneRecordFile is the lane record's name inside the worktree's own Git dir. It is
// distinct from the gate cache file, because a lane record is a distinct record class
// that no gate reader may read as a verdict.
const laneRecordFile = "bench-last-lane"

// defaultLaneName identifies a lane whose caller named none.
const defaultLaneName = "lane"

// laneRunBinary is the factory the lane selects its Bench executable through. It is the
// gate's own selection chain, so the Bench-owned checks grade with the tree's own code.
// It is a variable for the same reason the phase table is: a test supplies its build and
// verification seams rather than paying a real build.
var laneRunBinary = runbinary.Factory{}

// BenchkitLane is the built-in fast lane for the kit root: gofmt, prose, vet, build, and
// then one check per document family the registry binds. A Bench-owned check carries the
// run-binary token, which the run replaces with the executable it selected, and the two
// toolchain checks run the tool directly. The build disables VCS stamping, as every gate
// Go argv does: the private checkout is a linked worktree under a temporary directory,
// and Go's own discovery skips its `.git` file and adopts whatever `.git` directory sits
// above it. The whole-project test phase is deliberately absent. The lane is the worktree
// commit's check, and the landing's gate stays the one full grade.
func BenchkitLane(root, kit string) []Phase {
	_ = kit
	return append([]Phase{
		{Name: "gofmt", Argv: []string{runBinaryArgvToken, "gate-go", "gofmt"}},
		{Name: "prose", Argv: []string{runBinaryArgvToken, "gate-prose", root, "--", LaneNamedMarkdownToken}},
		{Name: "vet", Argv: []string{"go", "vet", trimPath, "./..."}},
		{Name: "build", Argv: []string{"go", "build", trimPath, DisableBuildVCS, "./..."}},
	}, documentLaneChecks()...)
}

// RunLane grades the composed tree against the declared checks and records one lane
// record. It materializes the tree as a private checkout, so it grades what the commit
// composes rather than the working tree beside it. The run is bounded by the gate
// timeout.
//
// An interrupted or timed-out run returns an error and writes nothing. A record of a
// partial run would carry an outcome no check ever reached.
func RunLane(ctx context.Context, req LaneRequest) (LaneResult, error) {
	ctx, finishSpan := beginLaneSpan(ctx, req)
	result, err := runLane(ctx, req)
	finishSpan(result, err)
	return result, err
}

func runLane(ctx context.Context, req LaneRequest) (LaneResult, error) {
	if len(req.Checks) == 0 {
		return LaneResult{}, errors.New("gate: lane declares no checks")
	}
	gitdir, err := benchgit.Output("-C", req.Root, "rev-parse", "--absolute-git-dir")
	if err != nil || gitdir == "" {
		return LaneResult{}, errors.New("gate: git directory unavailable")
	}
	artifacts, err := openProspectiveArtifacts(req.Root, req.Tree)
	if err != nil {
		return LaneResult{}, err
	}
	defer artifacts.Close()
	checkout := artifacts.Checkout()

	// The lane holds the shared cache lock for its span, so one rule covers every holder
	// and a clean cannot remove an archive a lane check is reading. A lock the lane cannot
	// take never fails the lane: the checks, not the cache, decide the outcome.
	if holder, err := gocache.Hold(os.Environ()); err == nil {
		defer holder.Release()
	}

	selected, names, classes, err := selectLaneChecks(req, checkout)
	if err != nil {
		return LaneResult{}, err
	}
	checks := resolveLane(selected, req.Root, checkout, proseSubject(req.Changes))
	runBinary, checks, closeSelection, err := selectLaneRunBinary(ctx, req, checkout, artifacts.Root(), checks)
	if err != nil {
		return LaneResult{}, err
	}
	defer closeSelection()

	runCtx, cancel := bounds.ContextCause(ctx, gateTimeout, errGateTimeout)
	defer cancel()
	diagnostics := &laneDiagnostics{first: map[string]string{}}
	results, cancelled := schedule(runCtx, checkout, checks, laneWriters(req.Stdout, req.Stderr, diagnostics))
	if cancelled {
		return LaneResult{}, errors.New("gate: lane interrupted")
	}
	if errors.Is(context.Cause(runCtx), errGateTimeout) {
		return LaneResult{}, errGateTimeout
	}

	result := LaneResult{
		Outcome:    lanePass,
		Checks:     names,
		Classes:    classes,
		Tree:       req.Tree,
		Lane:       laneName(req.Lane),
		RunBinary:  runBinary,
		RecordedAt: time.Now().UTC().Truncate(time.Second),
	}
	for _, settled := range results {
		if settled.Code == 0 {
			continue
		}
		result.Outcome, result.Check = laneFail, settled.Name
		result.Diagnostic = diagnostics.firstLine(settled.Name, settled.StartErr)
		break
	}
	if err := writeLaneRecord(gitdir, result); err != nil {
		return LaneResult{}, fmt.Errorf("gate: lane record persistence failed: %w", err)
	}
	return result, nil
}

// The lane's seam record. The lane runs in this process, so the boundary that opens the
// span is the lane itself rather than a verb above it.
const otelLaneSeam = "lane"

// beginLaneSpan starts the lane's span and returns the closer that ends it. A red lane
// carries the first failing check and that check's first diagnostic line, because the
// FT232 tripwire reads the lane's red from that one pair. The lane name rides in the
// span name, the way the gate's mode does, since the declared attribute set carries no
// lane identity.
func beginLaneSpan(ctx context.Context, req LaneRequest) (context.Context, func(LaneResult, error)) {
	ctx, span, finish := otelrecord.BeginIn(ctx, "", req.Root, otelLaneSeam, otelLaneSeam+"."+laneName(req.Lane))
	return ctx, func(result LaneResult, err error) {
		if req.Tree != "" {
			span.SetAttributes(attribute.String(otelrecord.AttrSubjectID, req.Tree))
		}
		// An interrupted or timed-out run settles no check, so it reaches no outcome to
		// record. The lane writes no lane record for that run either.
		if err == nil {
			span.SetAttributes(attribute.String(otelrecord.AttrOutcome, result.Outcome))
			// A pass names no failing check, and an empty attribute would read as one.
			if result.Check != "" {
				span.SetAttributes(attribute.String(otelrecord.AttrOutcomeCheck, result.Check))
			}
			// The first line alone. No further diagnostic output enters the record.
			if result.Diagnostic != "" {
				span.SetAttributes(attribute.String(otelrecord.AttrOutcomeDiagnostic, result.Diagnostic))
			}
		}
		finish()
	}
}

func laneName(declared string) string {
	if declared == "" {
		return defaultLaneName
	}
	return declared
}

// resolveLane rewrites a declared check list for the checkout that will run it. Two
// substitutions apply. The prose placeholder becomes the named Markdown paths. Any argv
// element or working directory anchored to the graded root is re-anchored to the private
// checkout, because that is where the tree the lane grades actually is.
func resolveLane(checks []Phase, root, checkout string, namedMarkdown []string) []Phase {
	resolved := make([]Phase, len(checks))
	for i, check := range checks {
		resolved[i] = check
		resolved[i].Argv = nil
		for _, arg := range check.Argv {
			if arg == LaneNamedMarkdownToken {
				resolved[i].Argv = append(resolved[i].Argv, namedMarkdown...)
				continue
			}
			resolved[i].Argv = append(resolved[i].Argv, reanchor(arg, root, checkout))
		}
		resolved[i].Dir = reanchor(check.Dir, root, checkout)
	}
	return resolved
}

// reanchor moves a path under root to the same place under checkout, and leaves every
// other value alone. A declared value that is not a path under root is not a path the
// checkout has a counterpart for.
func reanchor(value, root, checkout string) string {
	if value == "" || root == "" {
		return value
	}
	if value == root {
		return checkout
	}
	if strings.HasPrefix(value, root+string(filepath.Separator)) {
		return checkout + value[len(root):]
	}
	return value
}

// selectLaneRunBinary answers the run binary the lane records, and hands back the checks
// bound to it. A lane naming the token gets the gate's own selection chain. A lane naming
// no token runs no Bench-owned check, so there is nothing to select and the record
// carries the executable that ran the lane. An authored selection is written inside
// artifactRoot, so the bundle owns every executable it authors; an inherited selection
// stays where its own owner put it.
func selectLaneRunBinary(ctx context.Context, req LaneRequest, checkout, artifactRoot string, checks []Phase) (string, []Phase, func(), error) {
	if !laneUsesRunBinary(checks) {
		digest, err := ownExecutableDigest()
		return digest, checks, func() {}, err
	}
	source := req.Kit
	if source == "" {
		source = checkout
	}
	factory := laneRunBinary
	factory.TempRoot = artifactRoot
	selection, err := factory.ReuseOrOwn(ctx, source)
	if err != nil {
		return "", nil, func() {}, fmt.Errorf("gate: lane Bench executable unavailable: %w", err)
	}
	digest, err := fileDigest(selection.Path)
	if err != nil {
		_ = selection.Close()
		return "", nil, func() {}, err
	}
	return digest, withRunBinary(checks, selection), func() { _ = selection.Close() }, nil
}

func laneUsesRunBinary(checks []Phase) bool {
	for _, check := range checks {
		if len(check.Argv) > 0 && check.Argv[0] == runBinaryArgvToken {
			return true
		}
	}
	return false
}

func ownExecutableDigest() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("gate: lane executable unavailable: %w", err)
	}
	return fileDigest(path)
}

func fileDigest(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("gate: lane executable unreadable: %w", err)
	}
	defer file.Close()
	digest := sha256.New()
	if _, err := io.Copy(digest, file); err != nil {
		return "", fmt.Errorf("gate: lane executable unreadable: %w", err)
	}
	return hex.EncodeToString(digest.Sum(nil)), nil
}

// writeLaneRecord installs the record through the store's durable temp-then-rename path,
// the same one every other record class publishes through. It writes only the lane file.
// The verdict cache and the evidence store stay untouched, because a lane grades a check
// list and not the oracle.
func writeLaneRecord(gitdir string, result LaneResult) error {
	data, err := json.Marshal(verdictRecord{
		Schema:     verdictSchema,
		Tree:       result.Tree,
		Lane:       result.Lane,
		Outcome:    result.Outcome,
		RunBinary:  result.RunBinary,
		RecordedAt: result.RecordedAt.UTC().Truncate(time.Second).Format(time.RFC3339),
	})
	if err != nil {
		return err
	}
	return durableReplaceRecordAt(gitdir, laneRecordFile, data)
}

// laneDiagnostics keeps each check's first output line while the run streams. A caller
// naming a failing check needs one line of why, and re-reading the stream afterwards
// would need the whole stream retained.
type laneDiagnostics struct {
	mu    sync.Mutex
	first map[string]string
}

func (d *laneDiagnostics) record(name, line string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, held := d.first[name]; !held {
		d.first[name] = line
	}
}

// firstLine answers the check's first output line, and falls back to the start error for
// a check that never ran and so wrote nothing.
func (d *laneDiagnostics) firstLine(name string, startErr error) string {
	d.mu.Lock()
	defer d.mu.Unlock()
	if line, held := d.first[name]; held {
		return line
	}
	if startErr != nil {
		return startErr.Error()
	}
	return ""
}

// laneWriters is the lane's output plumbing. It keeps the gate's prefixed streams and
// taps each of them on the way through. The tap holds only the current partial line, so
// a chatty check costs the lane no memory.
func laneWriters(stdout, stderr io.Writer, diagnostics *laneDiagnostics) func(Phase) (io.Writer, io.Writer, func()) {
	inner := prefixedPhaseWriters(discardIfNil(stdout), discardIfNil(stderr))
	return func(check Phase) (io.Writer, io.Writer, func()) {
		out, errOut, closeWriters := inner(check)
		tapOut := &laneTapWriter{name: check.Name, dst: out, diagnostics: diagnostics}
		tapErr := &laneTapWriter{name: check.Name, dst: errOut, diagnostics: diagnostics}
		return tapOut, tapErr, func() {
			// A final line with no newline after it still names the failure, so it is
			// offered before the underlying writers close.
			tapOut.flush()
			tapErr.flush()
			closeWriters()
		}
	}
}

func discardIfNil(w io.Writer) io.Writer {
	if w == nil {
		return io.Discard
	}
	return w
}

type laneTapWriter struct {
	name        string
	dst         io.Writer
	diagnostics *laneDiagnostics
	pending     []byte
}

func (w *laneTapWriter) Write(p []byte) (int, error) {
	w.pending = append(w.pending, p...)
	for {
		idx := bytes.IndexByte(w.pending, '\n')
		if idx < 0 {
			break
		}
		w.diagnostics.record(w.name, string(w.pending[:idx]))
		w.pending = w.pending[idx+1:]
	}
	return w.dst.Write(p)
}

func (w *laneTapWriter) flush() {
	if len(w.pending) == 0 {
		return
	}
	w.diagnostics.record(w.name, string(w.pending))
	w.pending = nil
}
