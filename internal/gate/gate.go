// Package gate holds the oracle's selection logic in one Go home. It holds the
// ordered resolution chain (`.bench/gate.sh` beats `$BENCH_GATE` beats auto-detect),
// the gate run from the repo root, and the verdict-cache record keyed to
// git.TreeHash.
//
// Both the standalone `bench gate` (via the shell's one-glance run_gate →
// `bench gate-run`) and the in-process shift loop read this package. So gate
// resolution and the cache-write format each live in exactly one place. A second live
// resolver, or a second cache writer, is the worst class of bug in this kit. This kit's
// premise is "the gate is the oracle".
package gate

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/env"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gocache"
	"github.com/gibbonmi/bench/internal/otelrecord"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/toon"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

var gateTimeout = bounds.GateTimeout
var errGateTimeout = errors.New("gate deadline exceeded")

// Kind names the resolved gate. The zero value None is the no-gate case (exit 3,
// nothing recorded); the rest map to a command run from the repo root.
type Kind int

const (
	None Kind = iota
	GateSh
	ProspectiveGateSh
	BenchGate
	Pnpm
	Npm
	Pyproject
	Cargo
)

// Resolution is the chosen gate: its Kind and, for BenchGate, the command string the
// `$BENCH_GATE` env var carried.
type Resolution struct {
	Kind    Kind
	Command string
}

func resolutionName(kind Kind) string {
	return [...]string{"none", "gate-script", "prospective-gate-script", "bench-gate", "pnpm", "npm", "python", "cargo"}[kind]
}

// treeHashRE is the shape a real git tree hash must match before it is written to the
// verdict cache. treeHashRE refuses anything else, notably git.TreeHash's "none" on
// failure, so Record never forges a tree. The Stop hook shares this no-forged-verdict
// guarantee, because its recordGate delegates here.
var treeHashRE = regexp.MustCompile(`^[0-9a-f]+$`)

// FS injects the two filesystem probes Resolve needs: `-x` for the executable
// `.bench/gate.sh` and `-f` for the auto-detect lockfiles. This keeps the resolution
// precedence a pure function, unit-testable without a real tree.
type FS struct {
	Executable func(path string) bool
	Exists     func(path string) bool
}

// RealFS is the production probe set: a regular executable file for Executable, a
// regular (non-directory) file for Exists.
func RealFS() FS {
	return FS{
		Executable: func(p string) bool {
			info, err := os.Stat(p)
			return err == nil && !info.IsDir() && info.Mode()&0o111 != 0
		},
		Exists: func(p string) bool {
			info, err := os.Stat(p)
			return err == nil && !info.IsDir()
		},
	}
}

// Resolve is the ordered chain as a pure function. An executable `.bench/gate.sh`
// wins first, then a non-empty `$BENCH_GATE` wins next. After that, the first
// auto-detect lockfile wins, in the fixed order pnpm → npm → pyproject → cargo, then
// None wins. A reordered chain would silently run the wrong oracle. This is the
// precedence the NEW resolution-order contract and the table test both pin.
func Resolve(root, benchGate string, fs FS) Resolution {
	if fs.Executable(filepath.Join(root, ".bench", "gate.sh")) {
		return Resolution{Kind: GateSh}
	}
	if benchGate != "" {
		return Resolution{Kind: BenchGate, Command: benchGate}
	}
	for _, d := range []struct {
		file string
		kind Kind
	}{
		{"pnpm-lock.yaml", Pnpm},
		{"package.json", Npm},
		{"pyproject.toml", Pyproject},
		{"Cargo.toml", Cargo},
	} {
		if fs.Exists(filepath.Join(root, d.file)) {
			return Resolution{Kind: d.kind}
		}
	}
	return Resolution{Kind: None}
}

// command builds the shell command a resolution runs from the repo root. None has no
// command (handled by the caller). The auto-detect strings mirror bench.sh's
// best-effort defaults byte-for-byte.
func (r Resolution) command(root string) *exec.Cmd {
	switch r.Kind {
	case GateSh:
		return exec.Command(filepath.Join(root, ".bench", "gate.sh"))
	case ProspectiveGateSh:
		return exec.Command(filepath.Join(root, ".bench", "gate-prospective.sh"), root)
	case BenchGate:
		return exec.Command("bash", "-c", r.Command)
	case Pnpm:
		return exec.Command("bash", "-c", "pnpm -s typecheck && pnpm -s test && pnpm -s lint")
	case Npm:
		return exec.Command("bash", "-c", "npm run -s typecheck && npm test --silent && npm run -s lint")
	case Pyproject:
		return exec.Command("bash", "-c", "mypy . && pytest -q && ruff check .")
	case Cargo:
		return exec.Command("bash", "-c", "cargo test --quiet && cargo clippy -q -- -D warnings")
	default:
		return nil
	}
}

// gateEnv returns the caller's environment with wrapper-routing internals removed.
// `bench gate` reaches this package through bin/bench.sh -> route_binary, which sets
// them so the binary can find its assets. Those variables are not part of the project
// gate's contract. Leaking them into the gate makes fixture wrappers resolve the live
// kit instead of their own fabricated layout. The selected executable is not among
// them, because the gate's phase children are contracted to inherit it.
//
// The capability skip log goes too. A run owns the log its own phases append to, so an
// inherited path must never survive into a child. A collecting run sets its own value
// back on each phase.
//
// The Bench build cache entry is applied last, so a gate child writes to the one
// Bench-owned directory and an ambient GOCACHE never reaches it. An environment with
// no absolute HOME cannot name that directory, so gateEnv refuses rather than let a
// child fall back to Go's default.
func gateEnv() ([]string, error) {
	// The baseline phase-schedule selector goes too. It addresses the one gate-phases
	// process the owner launched; a phase child that inherited it would resolve its own
	// schedule against a root it was never handed.
	//
	// The seam record's two entries go too, for the same reason. They address the gate
	// run that started this process. A phase child that inherited them would attach its
	// own record lines to a run it is not part of. The parent sets them back on the one
	// child that runs the phase table.
	base := env.WithoutWrapperRouting(capability.WithoutEnvironment(capability.WithoutEnvironment(os.Environ(), capability.LogEnv), baselinePolicyEnv))
	for _, name := range otelGateEnv {
		base = capability.WithoutEnvironment(base, name)
	}
	return gocache.Apply(base)
}

// Run executes the resolved gate from the repo root and returns its exit code, with
// the gate's own output streamed to stdout/stderr. Run executes the gate from the
// working tree by design. An agent can edit the file it is graded by. The canary
// tripwire, not this call site, keeps that safe. None must not reach here; the caller
// handles the no-gate exit-3-nothing-recorded case.
func Run(root string, res Resolution, stdout, stderr io.Writer) int {
	childEnv, err := gateEnv()
	if err != nil {
		fmt.Fprintf(stderr, "gate environment unavailable: %v\n", err)
		return 1
	}
	return runResolved(context.Background(), root, res, childEnv, stdout, stderr, false).Code
}

// RunContext executes the resolved gate like Run, but puts the gate in its own process
// group. RunContext kills that group before returning when ctx is canceled. Shift uses
// this path so an interrupt cannot release the pooled worktree while a gate child keeps
// running and writing into it. Standalone `bench gate` uses Run, which preserves normal
// foreground-process signal delivery.
func RunContext(ctx context.Context, root string, res Resolution, stdout, stderr io.Writer) int {
	childEnv, err := gateEnv()
	if err != nil {
		fmt.Fprintf(stderr, "gate environment unavailable: %v\n", err)
		return 1
	}
	result := runResolved(ctx, root, res, childEnv, stdout, stderr, true)
	if result.StartErr != nil {
		return 1
	}
	return result.Code
}

// RunAndRecord resolves, runs, and records the gate for root, returning its exit code.
// The no-gate case exits 3 and records nothing (the chain resolved to None); every
// resolved gate runs and then records its verdict. Shared by the `gate-run` subcommand
// and the in-process shift loop, so neither carries its own resolve-run-record chain.
func RunAndRecord(root string, stdout, stderr io.Writer) int {
	return Execute(context.Background(), root, stdout, stderr).ActionExit
}

// RunAndRecordContext is RunAndRecord with cancellation for in-process callers that
// own teardown. A canceled gate is not recorded as red because the oracle did not
// finish judging the tree.
func RunAndRecordContext(ctx context.Context, root string, stdout, stderr io.Writer) int {
	return Execute(ctx, root, stdout, stderr).ActionExit
}

// RunCommand is the `bench gate-run [--fresh] [root]` plumbing subcommand. The shell's
// one-glance run_gate forwards here so gate resolution lives in exactly one place. Root
// is the first non-flag argument when the shell passes the resolved repo root, else the
// cwd's repo. RunCommand resolves root this way so the gate always runs from the top
// level, even when invoked from a subdirectory. `--fresh` may sit on either side of it.
func RunCommand(args []string, stdout, stderr io.Writer) int {
	var root string
	mode := reuseFreshGreen
	for _, arg := range args {
		switch {
		case arg == "--fresh":
			mode = forceRun
		case root == "":
			root = arg
		}
	}
	if root == "" {
		r, err := git.Root()
		if err != nil {
			fmt.Fprintln(stderr, toon.NotInRepo())
			return 1
		}
		root = r
	}
	ctx, finishSpan := beginGateSpan(context.Background(), root, mode.String())
	ctx, finishLog := beginGateRunLog(ctx, root, stderr, mode.String())
	result := executeAfterAcquire(ctx, root, stdout, stderr, notifyGateSignals, mode)
	finishLog(result)
	finishSpan(result)
	return result.ActionExit
}

const commandUsage = "usage: bench gate [--fresh|pin]"

func Command(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		return RunCommand(nil, stdout, stderr)
	}
	switch args[0] {
	case "pin":
		return PinCommand(args[1:], stdin, stdout, stderr)
	case "--fresh":
		if len(args) == 1 {
			return RunCommand(args, stdout, stderr)
		}
	case "--help", "-h", "help":
		if len(args) == 1 {
			fmt.Fprintln(stdout, commandUsage)
			return 0
		}
	}
	fmt.Fprintln(stderr, commandUsage)
	return 2
}

type Result struct {
	GateExit   int
	ActionExit int
	Inspection Inspection
}

func Execute(ctx context.Context, root string, stdout, stderr io.Writer) Result {
	ctx, finishSpan := beginGateSpan(ctx, root, "ordinary")
	ctx, finishLog := beginGateRunLog(ctx, root, stderr, "ordinary")
	result := execute(ctx, root, stdout, stderr)
	finishLog(result)
	finishSpan(result)
	return result
}

// ExecuteReusingFreshGreen answers for root's tree like Execute, but a verdict already
// reusable for this subject answers before the execution lock is touched. A gate run in
// progress elsewhere therefore neither refuses the caller nor demotes the green it would
// have reused. This is what makes the gated commit safe to run beside one.
//
// The optimistic subject comes from an evaluation-owned pre generation. The reuse
// decision authorizes skipping the gate, so it answers from the same snapshot contract
// a real run accepts, never from an independent capture. Everything else falls through
// to Execute and pays for a real run under the lock.
func ExecuteReusingFreshGreen(ctx context.Context, root string, stdout, stderr io.Writer) Result {
	if plan, err := newGateEvaluation(root).acceptPre(); err == nil {
		if reuse := reusableEvidence(root, plan, time.Now()); reuse.ReusableGreen {
			return reusedGreenResult(stdout, reuse)
		}
	}
	return execute(ctx, root, stdout, stderr)
}

// reusedGreenResult is the one place a reused verdict is announced and shaped into a
// result. The announcement is not optional. A skipped run that says nothing reads as a
// gate that never ran, and the operator has no way to tell the difference.
func reusedGreenResult(stdout io.Writer, reuse Inspection) Result {
	fmt.Fprintln(stdout, "gate: green (fresh verdict reused for this tree)")
	return Result{Inspection: reuse}
}

// runMode says whether a fresh green already recorded for this subject may answer the
// execution. `bench gate --fresh` picks forceRun. It is the operator's only escape from
// a green the closure still calls current, but the oracle would no longer stand behind.
type runMode int

const (
	reuseFreshGreen runMode = iota
	forceRun
)

func (m runMode) String() string {
	if m == forceRun {
		return "fresh"
	}
	return "ordinary"
}

func runBinarySource(runtimeRoot, storageRoot string, plan subject) string {
	if plan.Resolution.Kind == ProspectiveGateSh &&
		isRegularFile(filepath.Join(runtimeRoot, "scripts", "go-build.sh")) &&
		isRegularFile(filepath.Join(runtimeRoot, "go.mod")) {
		return runtimeRoot
	}
	return kitRoot(storageRoot)
}

// sameDirectory reports whether two paths name the same directory, by file identity
// rather than by string equality. This way a symlinked or differently-spelled path to
// one tree still matches. Either stat failing answers no, which is the fail-closed
// direction for the reduction guard.
func sameDirectory(a, b string) bool {
	ai, err := os.Stat(a)
	if err != nil {
		return false
	}
	bi, err := os.Stat(b)
	if err != nil {
		return false
	}
	return os.SameFile(ai, bi)
}

// phaseTableGate reports whether root's resolved gate provably routes through the
// phase table this package runs. Only such a gate may be reduced. Running the table
// under a gate that never execs it would swap the repository's oracle, a hand-written
// script, for one it never chose.
//
// A declared phase manifest chooses *which* table a routed run resolves, never
// *whether* the gate routes through one. So the proof is the same with or without a
// manifest. The resolved gate must be the script. The script must also carry the
// gate-phases hand-off the kit's own entry uses (the exec line the gate-entry
// conformance check pins). Anything else pays the full run; phaseTableGate fails
// closed, never fails cheap.
func phaseTableGate(root string, res Resolution) bool {
	if res.Kind != GateSh && res.Kind != ProspectiveGateSh {
		return false
	}
	script, err := os.ReadFile(filepath.Join(root, ".bench", "gate.sh"))
	return err == nil && strings.Contains(string(script), `"$bench" gate-phases`)
}

type postAcquireContextArm func(context.Context) (context.Context, func())

func notifyGateSignals(ctx context.Context) (context.Context, func()) {
	return subprocess.NotifyCancel(ctx)
}

func execute(ctx context.Context, root string, stdout, stderr io.Writer) Result {
	return executeSubjectWithRunBinary(ctx, root, root, stdout, stderr, nil, reuseFreshGreen, newGateEvaluation(root), productionRunBinaryOwner(), "")
}

func executeAfterAcquire(ctx context.Context, root string, stdout, stderr io.Writer, arm postAcquireContextArm, mode runMode) Result {
	return executeSubjectWithRunBinary(ctx, root, root, stdout, stderr, arm, mode, newGateEvaluation(root), productionRunBinaryOwner(), "")
}

func operational(root string, gateExit int, stderr io.Writer, msg string) Result {
	fmt.Fprintln(stderr, msg)
	inspection := inspectAt(root, time.Now().UTC())
	inspection.ReusableGreen = false
	return Result{GateExit: gateExit, ActionExit: 1, Inspection: inspection}
}

// The gate's seam record. The verb boundary resolves the Bench home once, builds the
// provider, and puts its tracer on the context for the whole run, the way the run log
// threads its own file. The phase table runs in a second process, so the record's
// repository and this run's trace reach that process through the environment.
const (
	otelGateSeam      = "gate"
	otelGatePhaseSeam = "gate.phase"

	otelRootEnv        = "BENCH_OTEL_ROOT"
	otelTraceparentEnv = "BENCH_OTEL_TRACEPARENT"
)

// otelGateEnv is the whole set the phases process inherits, so the stripper and the
// composer read one list rather than two that can disagree.
var otelGateEnv = []string{otelRootEnv, otelTraceparentEnv}

// otelRecordRootKey addresses the repository the run records under. The child that runs
// the phases records under the same repository, so the run's spans land in one file.
type otelRecordRootKey struct{}

// beginGateSpan starts the run's root span and returns the closer that ends it. The
// mode rides in the span name: story 19's declared attribute set names the seam, the
// subject, the outcome, and the measures, and none of those carries a run mode.
func beginGateSpan(ctx context.Context, root, mode string) (context.Context, func(Result)) {
	ctx, span, finish := otelrecord.BeginIn(ctx, "", root, otelGateSeam, otelGateSeam+"."+mode)
	ctx = context.WithValue(ctx, otelRecordRootKey{}, root)
	return ctx, func(result Result) {
		// A subject the run never resolved has no digest to group its iterations by, and
		// an empty attribute would read as one.
		if subject := result.Inspection.CurrentTree; subject != "" {
			span.SetAttributes(attribute.String(otelrecord.AttrSubjectID, subject))
		}
		// The action exit is the one the operator sees, so the record and the shell
		// agree about the same run.
		span.SetAttributes(attribute.String(otelrecord.AttrOutcome, otelrecord.ExitOutcome(result.ActionExit)))
		finish()
	}
}

// withGateSpanEnv hands the phases child the repository it records under and this run's
// trace, so its phase spans join the root span rather than start a trace of their own.
// A context outside a recorded run adds nothing.
func withGateSpanEnv(ctx context.Context, base []string) []string {
	root, _ := ctx.Value(otelRecordRootKey{}).(string)
	if root == "" {
		return base
	}
	env := append(base, otelRootEnv+"="+root)
	carrier := propagation.MapCarrier{}
	propagation.TraceContext{}.Inject(ctx, carrier)
	if parent := carrier.Get("traceparent"); parent != "" {
		env = append(env, otelTraceparentEnv+"="+parent)
	}
	return env
}
