package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/benchguard"
	"github.com/gibbonmi/bench/internal/census"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gitguard"
	"github.com/gibbonmi/bench/internal/lines"
	"github.com/gibbonmi/bench/internal/modelid"
	"github.com/gibbonmi/bench/internal/otelrecord"
	"github.com/gibbonmi/bench/internal/poolkey"
	"github.com/gibbonmi/bench/internal/usage"
	"github.com/gibbonmi/bench/internal/worktree"
	"github.com/gibbonmi/bench/internal/writeguard"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var resolveModelGrammar = usage.Grammar{
	Cmd:   "bench resolve-model",
	Help:  "usage: bench resolve-model --harness <" + strings.Join(lines.Harnesses, "|") + ">",
	Flags: []usage.Flag{harnessFlag},
}

// resolveModel is the `bench resolve-model` plumbing subcommand for the shift adapters.
// It prints the model to pass via the harness --model flag, empty for passthrough, to
// stdout and returns an exit code. Any warning or error goes to os.Stderr directly: the
// map signature carries only stdout, and the adapter captures stdout as the model, so a
// warning must never ride there. BENCH_MODEL names a tier and --harness names the column.
// In a routed repo an unset or unbound tier exits 1 and the adapter refuses to launch.
// The verdict lives in internal/lines, so it is unit-tested without a repo. Each
// resolution past the grammar records one line.resolve span.
func resolveModel(args []string) (string, int) {
	harness, line, code := parseHarness(resolveModelGrammar, args)
	if line != "" {
		fmt.Fprintln(os.Stderr, line)
		return "", code
	}
	benchModel, set := os.LookupEnv("BENCH_MODEL")
	model, code, stderr := lines.ResolveModelVerdict(harness, benchModel, set, linesEnv())
	recordResolve(harness, benchModel, model, code)
	if stderr != "" {
		fmt.Fprintln(os.Stderr, stderr)
	}
	if model == "" {
		return "", code
	}
	return model + "\n", code
}

// The line resolution's record seam.
const otelResolveSeam = "line.resolve"

// recordResolve writes the line.resolve span of one resolution. Each value is a known
// word or a safe token, so an operator's harness, tier, or model text never reaches the
// record.
func recordResolve(harness, tier, model string, exit int) {
	attrs := []attribute.KeyValue{attribute.String(otelrecord.AttrOutcome, otelrecord.ExitOutcome(exit))}
	if slices.Contains(lines.Harnesses, harness) {
		attrs = append(attrs, attribute.String(otelrecord.AttrLineHarness, harness))
	}
	if slices.Contains(lines.Tiers, tier) {
		attrs = append(attrs, attribute.String(otelrecord.AttrLineTier, tier))
	}
	if model != "" && modelid.SafeToken(model) {
		attrs = append(attrs, attribute.String(otelrecord.AttrLineModel, model))
	}
	span, end := beginResolveSpan()
	span.SetAttributes(attrs...)
	end()
}

// beginResolveSpan starts the line.resolve span under the handed-off trace when a shift
// handed one off, and in a trace of its own in the current repository otherwise.
func beginResolveSpan() (trace.Span, func()) {
	ctx, detach := otelrecord.AttachHandoff(context.Background())
	if trace.SpanContextFromContext(ctx).IsValid() {
		_, span := otelrecord.TracerFrom(ctx).Start(ctx, otelResolveSeam, trace.WithAttributes(attribute.String(otelrecord.AttrSeam, otelResolveSeam)))
		return span, func() { span.End(); detach() }
	}
	detach()
	root, err := git.Root()
	if err != nil {
		return trace.SpanFromContext(ctx), func() {}
	}
	_, span, end := otelrecord.Begin("", root, otelResolveSeam)
	return span, end
}

// guardGit is the destructive-git guard subcommand. It reads the PreToolUse envelope on
// stdin, classifies through internal/gitguard, and yields the verdict as an exit code:
// 0 allow, 2 block, with the `BLOCKED:` message on stderr, or 3 a genuine failure to run.
// The deferred recover maps any panic to 3, not Go's default exit-2, so exit 2 means
// only an intentional block, and the shim can trust it.
func guardGit(_ []string, stdin io.Reader, _ io.Writer, stderr io.Writer) (code int) {
	defer func() {
		if r := recover(); r != nil {
			code = 3
		}
	}()
	data, err := io.ReadAll(stdin)
	if err != nil {
		return 3
	}
	command := gitguard.CommandFromEnvelope(data)
	if command == "" {
		return 0
	}
	chk := gitguard.Checker{
		RefResolves:     git.RefResolves,
		BranchExists:    git.BranchExists,
		DefaultBranch:   guardDefaultBranch,
		CheckedOut:      guardCheckedOut,
		BareDestination: guardBareDestination,
	}
	label := gitguard.Classify(command, chk)
	if label == "" {
		return 0
	}
	fmt.Fprintln(stderr, gitguard.BlockMessage(label))
	return 2
}

// guardProbeRoot is the directory the guard's three push facts read. The guarded Bash
// command runs in the agent's cwd, and git.RefResolves and git.BranchExists already probe
// that directory with no root operand, so the push facts name it explicitly to reach the
// same repository.
const guardProbeRoot = "."

// guardDefaultBranch reports the repository's default branch, from the one Go owner of
// that fact. No answer denies the push, so the guard never guesses a protected name.
func guardDefaultBranch() (string, bool) { return git.ResolvedDefault(guardProbeRoot) }

// guardCheckedOut reports the checked-out branch, or no branch, from the one Go owner of
// that mapping. No answer denies a `HEAD` refspec with the unresolved class.
func guardCheckedOut() (string, bool) { return git.CheckedOutName(guardProbeRoot) }

// guardBareDestination reports the branch a bare `git push` targets, from the one Go
// owner of that fact.
func guardBareDestination() (string, bool) { return git.BarePushDestination(guardProbeRoot) }

func guardBenchFollowOn(_ []string, stdin io.Reader, _ io.Writer, stderr io.Writer) (code int) {
	defer func() {
		if recover() != nil {
			code = 3
		}
	}()
	data, err := io.ReadAll(stdin)
	if err != nil {
		return 3
	}
	command, err := benchguard.CommandFromEnvelope(data)
	if err != nil {
		fmt.Fprintln(stderr, "WARNING: block-bench-follow-on: unreadable command field — allowing Bash.")
		return 0
	}
	recordFollowOn(command)
	// The pool denial runs first. A pool reference with a Bench call after it has two
	// faults, and the pool reference is the cause the reader repairs.
	if target := benchguard.PoolReference(command, poolkey.Pools(worktree.Home())); target != "" {
		fmt.Fprintln(stderr, benchguard.PoolReferenceMessage(target))
		return 2
	}
	verdict := benchguard.Classify(command, benchguard.DefaultResolver())
	if !verdict.Blocked {
		return 0
	}
	fmt.Fprintln(stderr, verdict.Message())
	return 2
}

// guardFileWrite is the primary-checkout file-write guard subcommand. It reads the
// PreToolUse envelope on stdin, classifies through internal/writeguard, and yields the
// verdict as an exit code: 0 allow, 2 block with the `BLOCKED:` message on stderr, or 3 a
// genuine failure to run. The deferred recover maps any panic to 3, not Go's default
// exit-2, so exit 2 means only an intentional block, and the shim can trust it.
func guardFileWrite(_ []string, stdin io.Reader, _ io.Writer, stderr io.Writer) (code int) {
	defer func() {
		if recover() != nil {
			code = 3
		}
	}()
	data, err := io.ReadAll(stdin)
	if err != nil {
		return 3
	}
	path, err := writeguard.PathFromEnvelope(data)
	if err != nil {
		fmt.Fprintln(stderr, "WARNING: block-primary-file-write: unreadable file_path field — allowing the write.")
		return 0
	}
	verdict := writeguard.Classify(path, writeguard.Checker{
		RootAt:    git.RootAt,
		IsPrimary: git.IsPrimaryCheckout,
		IsTracked: guardPathTracked,
		IsIgnored: guardPathIgnored,
	})
	if !verdict.Blocked {
		return 0
	}
	fmt.Fprintln(stderr, verdict.Message())
	return 2
}

// guardPathTracked reports whether git tracks path in root. An absolute pathspec is what
// the envelope carries, and git resolves it against the repository itself.
func guardPathTracked(root, path string) bool {
	return git.OK("-C", root, "ls-files", "--error-unmatch", "--", path)
}

// guardPathIgnored reports whether root's ignore rules cover path.
func guardPathIgnored(root, path string) bool {
	return git.OK("-C", root, "check-ignore", "-q", "--", path)
}

// recordFollowOn records a raw call through the exec census. It tests the command
// text for the pool prefix before it resolves any root, so an ordinary call outside
// a Bench worktree spawns no git process. Its own failure is silent and never reaches
// the verdict: this call sits before the verdict, so no later return can skip it.
func recordFollowOn(command string) {
	home := worktree.Home()
	if !strings.Contains(command, poolkey.Pools(home)+string(filepath.Separator)) {
		return
	}
	root, err := git.Root()
	if err != nil {
		return
	}
	_ = census.Record(command, root, home, time.Now())
}
