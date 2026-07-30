// Package gate is the oracle's selection logic in one Go home: the ordered
// resolution chain (`.bench/gate.sh` beats `$BENCH_GATE` beats auto-detect), the gate
// run from the repo root, and the verdict-cache record keyed to git.TreeHash. Both the
// standalone `bench gate` (via the shell's one-glance run_gate → `bench gate-run`) and
// the in-process shift loop read this package, so gate resolution and the cache-write
// format each live in exactly one place — a second live resolver, or a second cache
// writer, is the worst class of bug in a kit whose premise is "the gate is the oracle".
package gate

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

var gateTimeout = bounds.GateTimeout
var errGateTimeout = errors.New("gate deadline exceeded")

// Kind names the resolved gate. The zero value None is the no-gate case (exit 3,
// nothing recorded); the rest map to a command run from the repo root.
type Kind int

const (
	None Kind = iota
	GateSh
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
	return [...]string{"none", "gate-script", "bench-gate", "pnpm", "npm", "python", "cargo"}[kind]
}

// treeHashRE is the shape a real git tree hash must match before it is written to the
// verdict cache. Anything else (notably git.TreeHash's "none" on failure) is refused,
// so Record never forges a tree — the no-forged-verdict guarantee shared with the Stop
// hook, whose recordGate delegates here.
var treeHashRE = regexp.MustCompile(`^[0-9a-f]+$`)

// FS injects the two filesystem probes Resolve needs — `-x` for the executable
// `.bench/gate.sh` and `-f` for the auto-detect lockfiles — so the resolution
// precedence is a pure function unit-testable without a real tree.
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

// Resolve is the ordered chain as a pure function: an executable `.bench/gate.sh`
// wins, then a non-empty `$BENCH_GATE`, then the first auto-detect lockfile in the
// fixed order pnpm → npm → pyproject → cargo, then None. A reordered chain would
// silently run the wrong oracle; this is the precedence the NEW resolution-order
// contract and the table test both pin.
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
// BENCH_KIT/BENCH_WRAPPER so the binary can find its assets. Those are not part of the
// project gate's contract; leaking them into the gate makes fixture wrappers resolve the
// live kit instead of their own fabricated layout.
//
// The capability skip log goes with them: a run owns the log its own phases append to,
// so an inherited path — the outer gate's, reaching a canary's inner run — must never
// survive into a child. A run that collects sets its own value back on each phase.
func gateEnv() []string {
	var env []string
	for _, kv := range capability.WithoutEnvironment(os.Environ(), capability.LogEnv) {
		if strings.HasPrefix(kv, "BENCH_KIT=") || strings.HasPrefix(kv, "BENCH_WRAPPER=") {
			continue
		}
		env = append(env, kv)
	}
	return env
}

// Run executes the resolved gate from the repo root and returns its exit code, with
// the gate's own output streamed to stdout/stderr. The gate is run from the working
// tree by design (an agent can edit the file it is graded by; the canary tripwire, not
// this call site, keeps that safe). None must not reach here — the caller handles the
// no-gate exit-3-nothing-recorded case.
func Run(root string, res Resolution, stdout, stderr io.Writer) int {
	return runResolved(context.Background(), root, res, gateEnv(), stdout, stderr, false).Code
}

// RunContext executes the resolved gate like Run, but puts the gate in its own process
// group and kills that group before returning when ctx is canceled. Shift uses this
// path so an interrupt cannot release the pooled worktree while a gate child keeps
// running and writing into it. Standalone `bench gate` uses Run, preserving normal
// foreground-process signal delivery.
func RunContext(ctx context.Context, root string, res Resolution, stdout, stderr io.Writer) int {
	result := runResolved(ctx, root, res, gateEnv(), stdout, stderr, true)
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

// RunCommand is the `bench gate-run [--fresh] [root]` plumbing subcommand: the shell's
// one-glance run_gate forwards here so gate resolution lives in exactly one place. Root
// is the first non-flag argument when the shell passes the resolved repo root, else the
// cwd's repo — resolved so the gate always runs from the top level even when invoked from
// a subdirectory. `--fresh` may sit on either side of it.
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
	return executeWithEngineAfterAcquire(context.Background(), root, stdout, stderr, productionGateEngine{}, notifyGateSignals, mode).ActionExit
}

type Result struct {
	GateExit   int
	ActionExit int
	Inspection Inspection
}

func Execute(ctx context.Context, root string, stdout, stderr io.Writer) Result {
	return executeWithEngine(ctx, root, stdout, stderr, productionGateEngine{})
}

// ExecuteReusingFreshGreen answers for root's tree like Execute, but a verdict already
// reusable for this subject answers before the execution lock is touched. A gate run in
// progress elsewhere therefore neither refuses the caller nor demotes the green it would
// have reused, which is what makes the gated commit safe to run beside one. Everything else
// falls through to Execute and pays a real run under the lock.
func ExecuteReusingFreshGreen(ctx context.Context, root string, stdout, stderr io.Writer) Result {
	if reuse := Inspect(root); reuse.ReusableGreen {
		return reusedGreenResult(stdout, reuse)
	}
	return Execute(ctx, root, stdout, stderr)
}

// reusedGreenResult is the one place a reused verdict is announced and shaped into a result.
// The announcement is not optional: a skipped run that says nothing reads as a gate that
// never ran, and the operator has no way to tell the difference.
func reusedGreenResult(stdout io.Writer, reuse Inspection) Result {
	fmt.Fprintln(stdout, "gate: green (fresh verdict reused for this tree)")
	return Result{Inspection: reuse}
}

// runMode says whether a fresh green already recorded for this subject may answer the
// execution. `bench gate --fresh` picks forceRun: it is the operator's only escape from a
// green the closure still calls current but the oracle would no longer stand behind.
type runMode int

const (
	reuseFreshGreen runMode = iota
	forceRun
)

type postAcquireContextArm func(context.Context) (context.Context, func())

func notifyGateSignals(ctx context.Context) (context.Context, func()) {
	return signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
}

func executeWithEngine(ctx context.Context, root string, stdout, stderr io.Writer, engine gateEngine) Result {
	return executeWithEngineAfterAcquire(ctx, root, stdout, stderr, engine, nil, reuseFreshGreen)
}

func executeWithEngineAfterAcquire(ctx context.Context, root string, stdout, stderr io.Writer, engine gateEngine, arm postAcquireContextArm, mode runMode) Result {
	plan, err := engine.BuildSubject(root)
	if err != nil {
		return operationalWithEngine(engine, root, 0, stderr, "gate subject unavailable")
	}
	if plan.Resolution.Kind == None {
		fmt.Fprintln(stderr, "no gate found: add an executable .bench/gate.sh or set BENCH_GATE")
		return Result{GateExit: 3, ActionExit: 3, Inspection: inspectAt(root, engine.Now())}
	}
	gitdir, err := engine.GitDir(root)
	if err != nil {
		return operationalWithEngine(engine, root, 0, stderr, "git directory unavailable")
	}
	lock, err := engine.OpenLock(filepath.Join(gitdir, "bench-gate.lock"))
	if err != nil {
		persistInterruptedIfGreen(engine, root, gitdir, plan)
		return operationalWithEngine(engine, root, 0, stderr, "gate lock unavailable")
	}
	defer lock.Close()
	if err := engine.Acquire(lock); err != nil {
		persistInterruptedIfGreen(engine, root, gitdir, plan)
		fmt.Fprintln(stderr, "gate execution already in progress")
		writeOwnerDiagnostic(stderr, filepath.Join(gitdir, "bench-gate-owner"))
		inspection := inspectAt(root, engine.Now())
		inspection.ReusableGreen = false
		return Result{ActionExit: 1, Inspection: inspection}
	}
	defer engine.Unlock(lock)
	if arm != nil {
		var stop func()
		ctx, stop = arm(ctx)
		defer stop()
	}
	ownerPath := filepath.Join(gitdir, "bench-gate-owner")
	defer func() { _ = engine.Remove(ownerPath) }()
	if err := engine.WriteFile(ownerPath, ownerRecord(engine.Now()), 0o600); err != nil {
		return operationalWithEngine(engine, root, 0, stderr, "gate owner persistence failed")
	}
	underLock, err := engine.BuildSubject(root)
	if err != nil || !sameSubject(plan, underLock) {
		return operationalWithEngine(engine, root, 0, stderr, "gate subject changed before execution")
	}
	// A reusable green is answered from the record without touching it: re-recording the
	// verdict would push RecordedAt forward on every read and make the freshness window
	// unbounded. The check sits ahead of the pending replace so a reuse returns with nothing
	// written — no pending record to leave behind, no verdict to restore.
	if mode == reuseFreshGreen {
		if reuse := inspectAt(root, engine.Now()); reuse.ReusableGreen {
			return reusedGreenResult(stdout, reuse)
		}
	}
	pending := interruptedRecord(plan, engine.Now())
	if err := durableReplaceWithEngine(engine, gitdir, pending); err != nil {
		_ = durableReplaceWithEngine(engine, gitdir, pending)
		return operationalWithEngine(engine, root, 0, stderr, "gate pending persistence failed")
	}
	runCtx, cancelRun := bounds.ContextCause(ctx, gateTimeout, errGateTimeout)
	defer cancelRun()
	rc := runCaptured(runCtx, root, plan, stdout, stderr)
	if ctx.Err() != nil {
		return Result{GateExit: rc, ActionExit: rc, Inspection: inspectAt(root, engine.Now())}
	}
	if errors.Is(context.Cause(runCtx), errGateTimeout) {
		fmt.Fprintln(stderr, "gate: timeout")
		ready := verdictRecord{Schema: 1, State: Ready, Status: "timeout", Tree: plan.Tree, Oracle: plan.Oracle, RecordedAt: engine.Now().UTC().Truncate(time.Second).Format(time.RFC3339)}
		if err := durableReplaceWithEngine(engine, gitdir, ready); err != nil {
			_ = durableReplaceWithEngine(engine, gitdir, pending)
			return operationalWithEngine(engine, root, 124, stderr, "gate timeout persistence failed")
		}
		return Result{GateExit: 124, ActionExit: 124, Inspection: inspectAt(root, engine.Now())}
	}
	after, err := engine.PostRunSubject(root)
	if err != nil || !sameSubject(plan, after) {
		fmt.Fprintln(stderr, "gate subject changed during execution")
		return Result{GateExit: rc, ActionExit: 1, Inspection: inspectAt(root, engine.Now())}
	}
	status := "red"
	if rc == 0 {
		status = "green"
	}
	ready := verdictRecord{Schema: 1, State: Ready, Status: status, Tree: plan.Tree, Oracle: plan.Oracle, RecordedAt: engine.Now().UTC().Truncate(time.Second).Format(time.RFC3339)}
	if err := durableReplaceWithEngine(engine, gitdir, ready); err != nil {
		_ = durableReplaceWithEngine(engine, gitdir, pending)
		fmt.Fprintln(stderr, "gate final persistence failed")
		return Result{GateExit: rc, ActionExit: 1, Inspection: inspectAt(root, engine.Now())}
	}
	return Result{GateExit: rc, ActionExit: rc, Inspection: inspectAt(root, engine.Now())}
}

func ownerRecord(now time.Time) []byte {
	return []byte(strconv.Itoa(os.Getpid()) + " " + now.UTC().Truncate(time.Second).Format(time.RFC3339) + "\n")
}

func writeOwnerDiagnostic(stderr io.Writer, path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	fields := strings.Fields(string(data))
	if len(fields) != 2 {
		return
	}
	pid, err := strconv.Atoi(fields[0])
	if err != nil || pid <= 0 {
		return
	}
	if _, err := time.Parse(time.RFC3339, fields[1]); err != nil {
		return
	}
	liveness := "alive"
	if err := syscall.Kill(pid, 0); err != nil && err != syscall.EPERM {
		liveness = "not alive"
	}
	fmt.Fprintf(stderr, "gate owner: pid %d (%s)\n", pid, liveness)
}

func interruptedRecord(plan subject, now time.Time) verdictRecord {
	return verdictRecord{Schema: 1, State: Pending, Tree: plan.Tree, Oracle: plan.Oracle, StartedAt: now.UTC().Truncate(time.Second).Format(time.RFC3339), OwnerPID: os.Getpid()}
}

func operational(root string, gateExit int, stderr io.Writer, msg string) Result {
	return operationalWithEngine(productionGateEngine{}, root, gateExit, stderr, msg)
}

func operationalWithEngine(engine gateEngine, root string, gateExit int, stderr io.Writer, msg string) Result {
	fmt.Fprintln(stderr, msg)
	inspection := inspectAt(root, engine.Now())
	inspection.ReusableGreen = false
	return Result{GateExit: gateExit, ActionExit: 1, Inspection: inspection}
}

func sameSubject(a, b subject) bool {
	return a.Tree == b.Tree && a.Oracle == b.Oracle && a.Resolution == b.Resolution && a.Closed == b.Closed && a.Reason == b.Reason
}

func runCaptured(ctx context.Context, root string, s subject, stdout, stderr io.Writer) int {
	return runResolved(ctx, root, s.Resolution, s.Env, controlSafeWriter{stdout}, controlSafeWriter{stderr}, true).Code
}

func runResolved(ctx context.Context, root string, res Resolution, env []string, stdout, stderr io.Writer, processGroup bool) processGroupResult {
	cmd := res.command(root)
	if cmd == nil {
		return processGroupResult{Code: 3}
	}
	cmd.Dir, cmd.Stdout, cmd.Stderr = root, stdout, stderr
	cmd.Env = append([]string(nil), env...)
	if processGroup {
		return runProcessGroupCommand(ctx, cmd)
	}
	if err := cmd.Run(); err != nil {
		if cmd.ProcessState != nil {
			if code := cmd.ProcessState.ExitCode(); code > 0 {
				return processGroupResult{Code: code}
			}
		}
		return processGroupResult{Code: 1, StartErr: err}
	}
	return processGroupResult{}
}

// controlSafeWriter preserves gate output while removing C0 bytes that can execute
// terminal controls. Newline, carriage return, and tab remain ordinary formatting.
type controlSafeWriter struct{ io.Writer }

func (w controlSafeWriter) Write(p []byte) (int, error) {
	safe := make([]byte, 0, len(p))
	for _, b := range p {
		if (b >= 0x20 && b != 0x7f) || b == '\n' || b == '\r' || b == '\t' {
			safe = append(safe, b)
		}
	}
	if _, err := w.Writer.Write(safe); err != nil {
		return 0, err
	}
	return len(p), nil
}
