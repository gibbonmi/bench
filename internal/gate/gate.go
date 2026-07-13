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
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

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
func gateEnv() []string {
	var env []string
	for _, kv := range os.Environ() {
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
	cmd := res.command(root)
	if cmd == nil {
		return 3
	}
	cmd.Dir = root
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Env = gateEnv()
	if err := cmd.Run(); err != nil {
		if cmd.ProcessState != nil {
			if code := cmd.ProcessState.ExitCode(); code > 0 {
				return code
			}
		}
		return 1 // failed to start, or a signal death: treat as red
	}
	return 0
}

// RunContext executes the resolved gate like Run, but puts the gate in its own process
// group and kills that group before returning when ctx is canceled. Shift uses this
// path so an interrupt cannot release the pooled worktree while a gate child keeps
// running and writing into it. Standalone `bench gate` uses Run, preserving normal
// foreground-process signal delivery.
func RunContext(ctx context.Context, root string, res Resolution, stdout, stderr io.Writer) int {
	cmd := res.command(root)
	if cmd == nil {
		return 3
	}
	cmd.Dir = root
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Env = gateEnv()
	result := runProcessGroupCommand(ctx, cmd)
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

// RunCommand is the `bench gate-run [root]` plumbing subcommand: the shell's one-glance
// run_gate forwards here so gate resolution lives in exactly one place. Root is args[0]
// when the shell passes the resolved repo root, else the cwd's repo — resolved so the
// gate always runs from the top level even when invoked from a subdirectory.
func RunCommand(args []string, stdout, stderr io.Writer) int {
	var root string
	if len(args) > 0 && args[0] != "" {
		root = args[0]
	} else {
		r, err := git.Root()
		if err != nil {
			fmt.Fprintln(stderr, toon.NotInRepo())
			return 1
		}
		root = r
	}
	return RunAndRecord(root, stdout, stderr)
}

type Result struct {
	GateExit   int
	ActionExit int
	Inspection Inspection
}

func Execute(ctx context.Context, root string, stdout, stderr io.Writer) Result {
	return executeWithEngine(ctx, root, stdout, stderr, productionGateEngine{})
}

func executeWithEngine(ctx context.Context, root string, stdout, stderr io.Writer, engine gateEngine) Result {
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
		return operationalWithEngine(engine, root, 0, stderr, "gate lock unavailable")
	}
	defer lock.Close()
	if err := engine.Acquire(lock); err != nil {
		fmt.Fprintln(stderr, "gate execution already in progress")
		inspection := inspectAt(root, engine.Now())
		inspection.ReusableGreen = false
		return Result{ActionExit: 1, Inspection: inspection}
	}
	defer engine.Unlock(lock)
	underLock, err := engine.BuildSubject(root)
	if err != nil || !sameSubject(plan, underLock) {
		return operationalWithEngine(engine, root, 0, stderr, "gate subject changed before execution")
	}
	now := engine.Now().UTC().Truncate(time.Second)
	pending := verdictRecord{Schema: 1, State: Pending, Tree: plan.Tree, Oracle: plan.Oracle, StartedAt: now.Format(time.RFC3339), OwnerPID: os.Getpid()}
	if err := durableReplaceWithEngine(engine, gitdir, pending); err != nil {
		_ = durableReplaceWithEngine(engine, gitdir, pending)
		return operationalWithEngine(engine, root, 0, stderr, "gate pending persistence failed")
	}
	rc := runCaptured(ctx, root, plan, stdout, stderr)
	if ctx.Err() != nil {
		return Result{GateExit: rc, ActionExit: rc, Inspection: inspectAt(root, engine.Now())}
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
	cmd := s.Resolution.command(root)
	if cmd == nil {
		return 3
	}
	cmd.Dir, cmd.Stdout, cmd.Stderr = root, controlSafeWriter{stdout}, controlSafeWriter{stderr}
	cmd.Env = append([]string(nil), s.Env...)
	return runProcessGroupCommand(ctx, cmd).Code
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

const processGroupCancelGrace = 2 * time.Second

type processGroupResult struct {
	Code     int
	StartErr error
}

func runProcessGroupCommand(ctx context.Context, cmd *exec.Cmd) processGroupResult {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setpgid = true
	if err := cmd.Start(); err != nil {
		if ctx.Err() != nil {
			return processGroupResult{Code: 130, StartErr: err}
		}
		return processGroupResult{Code: 1, StartErr: err}
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		return processGroupResult{Code: processExitCode(cmd, err)}
	case <-ctx.Done():
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGINT)
		select {
		case <-done:
			// The command leader may exit on INT while a descendant in the same
			// process group ignores it. A final group kill is still required: Wait
			// only proves the leader and its output pipes are done.
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		case <-time.After(processGroupCancelGrace):
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			<-done
		}
		return processGroupResult{Code: 130}
	}
}

func processExitCode(cmd *exec.Cmd, err error) int {
	if err == nil {
		return 0
	}
	if cmd.ProcessState != nil {
		if code := cmd.ProcessState.ExitCode(); code > 0 {
			return code
		}
	}
	return 1
}
