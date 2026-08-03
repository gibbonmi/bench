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
	if plan, err := buildSubject(root); err == nil {
		if reuse := reusableEvidence(root, plan, time.Now()); reuse.ReusableGreen {
			return reusedGreenResult(stdout, reuse)
		}
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

// sameDirectory reports whether two paths name the same directory, by file identity
// rather than by string equality, so a symlinked or differently-spelled path to one tree
// still matches. Either stat failing answers no — for the reduction guard that is the
// fail-closed direction.
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
// phase table this package runs. Only such a gate may be reduced: running the table
// under a gate that never execs it would swap the repository's oracle — a hand-written
// script — for one it never chose. A declared phase manifest chooses *which* table a
// routed run resolves, never *whether* the gate routes through one, so the proof is the
// same with or without a manifest: the resolved gate must be the script, and the script
// must carry the gate-phases hand-off the kit's own entry uses (the exec line the
// gate-entry conformance check pins). Anything else pays the full run — fail closed,
// never fail cheap.
func phaseTableGate(root string, res Resolution) bool {
	if res.Kind != GateSh {
		return false
	}
	script, err := os.ReadFile(filepath.Join(root, ".bench", "gate.sh"))
	return err == nil && strings.Contains(string(script), "gate-phases")
}

type postAcquireContextArm func(context.Context) (context.Context, func())

func notifyGateSignals(ctx context.Context) (context.Context, func()) {
	return signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
}

func executeWithEngine(ctx context.Context, root string, stdout, stderr io.Writer, engine gateEngine) Result {
	return executeSubjectWithEngine(ctx, root, root, stdout, stderr, engine, nil, reuseFreshGreen,
		func() (subject, error) { return engine.BuildSubject(root) },
		func() (subject, error) { return engine.PostRunSubject(root) })
}

func executeWithEngineAfterAcquire(ctx context.Context, root string, stdout, stderr io.Writer, engine gateEngine, arm postAcquireContextArm, mode runMode) Result {
	return executeSubjectWithEngine(ctx, root, root, stdout, stderr, engine, arm, mode,
		func() (subject, error) { return engine.BuildSubject(root) },
		func() (subject, error) { return engine.PostRunSubject(root) })
}

func executeSubjectWithEngine(ctx context.Context, runtimeRoot, storageRoot string, stdout, stderr io.Writer, engine gateEngine, arm postAcquireContextArm, mode runMode, build, postRun func() (subject, error)) Result {
	plan, err := build()
	if err != nil {
		return operationalWithEngine(engine, storageRoot, 0, stderr, "gate subject unavailable")
	}
	if plan.Resolution.Kind == None {
		fmt.Fprintln(stderr, "no gate found: add an executable .bench/gate.sh or set BENCH_GATE")
		return Result{GateExit: 3, ActionExit: 3, Inspection: inspectAt(storageRoot, engine.Now())}
	}
	gitdir, err := engine.GitDir(storageRoot)
	if err != nil {
		return operationalWithEngine(engine, storageRoot, 0, stderr, "git directory unavailable")
	}
	lock, err := engine.OpenLock(filepath.Join(gitdir, "bench-gate.lock"))
	if err != nil {
		persistInterruptedIfGreen(engine, storageRoot, gitdir, plan)
		return operationalWithEngine(engine, storageRoot, 0, stderr, "gate lock unavailable")
	}
	defer lock.Close()
	if err := engine.Acquire(lock); err != nil {
		persistInterruptedIfGreen(engine, storageRoot, gitdir, plan)
		fmt.Fprintln(stderr, "gate execution already in progress")
		writeOwnerDiagnostic(stderr, filepath.Join(gitdir, "bench-gate-owner"))
		inspection := inspectAt(storageRoot, engine.Now())
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
		return operationalWithEngine(engine, storageRoot, 0, stderr, "gate owner persistence failed")
	}
	underLock, err := build()
	if err != nil || !sameSubject(plan, underLock) {
		return operationalWithEngine(engine, storageRoot, 0, stderr, "gate subject changed before execution")
	}
	// A reusable green is answered from the record without touching it: re-recording the
	// verdict would push RecordedAt forward on every read and make the freshness window
	// unbounded. The check sits ahead of the pending replace so a reuse returns with nothing
	// written — no pending record to leave behind, no verdict to restore.
	if mode == reuseFreshGreen {
		if reuse := reusableEvidence(storageRoot, plan, engine.Now()); reuse.ReusableGreen {
			return reusedGreenResult(stdout, reuse)
		}
	}
	// The narrowing decision sits under the execution lock, after the whole-tree reuse
	// answer — reuse costs nothing and grades everything, so it stays the first answer.
	// Only the ordinary path (runtime and storage root the same tree) narrows; a
	// prospective execution grades a tree about to become HEAD and keeps the full gate.
	// An ineligible scoping falls back to the full run, never to a coarser skip: a
	// fail-closed refusal stays a full run.
	//
	// forceRun consults no partition: `bench gate --fresh` is the operator's one escape
	// to a real whole-tree run. It still resolves identities, because a forced green has to
	// re-author the slot of every component it just graded.
	ordinary := runtimeRoot == storageRoot
	var scoping componentScoping
	if ordinary {
		scoping = scopeComponents(runtimeRoot, plan.Resolution, mode, engine.Now())
	}
	pending := interruptedRecord(plan, engine.Now())
	if err := durableReplaceWithEngine(engine, gitdir, pending); err != nil {
		_ = durableReplaceWithEngine(engine, gitdir, pending)
		return operationalWithEngine(engine, storageRoot, 0, stderr, "gate pending persistence failed")
	}
	runCtx, cancelRun := bounds.ContextCause(ctx, gateTimeout, errGateTimeout)
	defer cancelRun()
	var rc int
	switch {
	case scoping.partial():
		// The announcement is not optional: a skipped grading surface that says
		// nothing reads as a gate that never ran (the reused-verdict line above is
		// the same rule). One line per component, each naming the evidence that
		// covered that component — a single summary line would tell an operator that
		// something was skipped without letting them check what stood in for it.
		for _, skip := range scoping.skipped {
			fmt.Fprintln(stdout, skip.announcement())
		}
		rc = runPhases(runCtx, scoping.runnerRoot, scoping.phases, outerMode, controlSafeWriter{stdout}, controlSafeWriter{stderr})
	default:
		rc = runCaptured(runCtx, runtimeRoot, plan, stdout, stderr)
	}
	if ctx.Err() != nil {
		return Result{GateExit: rc, ActionExit: rc, Inspection: inspectAt(storageRoot, engine.Now())}
	}
	if errors.Is(context.Cause(runCtx), errGateTimeout) {
		fmt.Fprintln(stderr, "gate: timeout")
		if err := invalidateEvidence(storageRoot, plan); err != nil {
			return operationalWithEngine(engine, storageRoot, 124, stderr, "gate evidence invalidation failed")
		}
		ready := verdictRecord{Schema: 1, State: Ready, Status: "timeout", Tree: plan.Tree, Oracle: plan.Oracle, RecordedAt: engine.Now().UTC().Truncate(time.Second).Format(time.RFC3339)}
		if err := durableReplaceWithEngine(engine, gitdir, ready); err != nil {
			_ = durableReplaceWithEngine(engine, gitdir, pending)
			return operationalWithEngine(engine, storageRoot, 124, stderr, "gate timeout persistence failed")
		}
		return Result{GateExit: 124, ActionExit: 124, Inspection: inspectAt(storageRoot, engine.Now())}
	}
	after, err := postRun()
	if err != nil || !sameSubject(plan, after) {
		fmt.Fprintln(stderr, "gate subject changed during execution")
		return Result{GateExit: rc, ActionExit: 1, Inspection: inspectAt(storageRoot, engine.Now())}
	}
	status := "red"
	if rc == 0 {
		status = "green"
	}
	recordedAt := engine.Now()
	ready := verdictRecord{Schema: 1, State: Ready, Status: status, Tree: plan.Tree, Oracle: plan.Oracle, RecordedAt: recordedAt.UTC().Truncate(time.Second).Format(time.RFC3339)}
	if scoping.partial() {
		// The record names what this run graded and, per component it did not, the
		// evidence that covered it. The identity and the recorded time are carried
		// rather than referenced: the cache is a single slot, so a reference would
		// point at whichever record replaced this one, and a consumer refusing a
		// release needs to name the evidence a skip rested on.
		ready.Executed = scoping.executedPhaseNames()
		ready.SkipEvidence = make(map[string]skipEvidence, len(scoping.skipped))
		for _, skip := range scoping.skipped {
			ready.Skipped = append(ready.Skipped, skip.Component)
			ready.SkipEvidence[skip.Component] = skip.evidence()
		}
	}
	if status == "green" {
		// The slots are this run's own account of what it graded, so they are written
		// whatever shape the run had: a full run authors every component's, a partial
		// run authors the ones it executed, and a forced run re-authors all of them.
		// Failing to author only costs a future run its skip, but failing silently
		// would hide that cost, so it shares the persistence-failure posture below.
		if err := authorExecutedComponentSlots(storageRoot, scoping, recordedAt); err != nil {
			fmt.Fprintln(stderr, "gate evidence persistence failed")
			return Result{GateExit: rc, ActionExit: 1, Inspection: inspectAt(storageRoot, engine.Now())}
		}
		// The build's evidence is the artifact rather than a slot, so it is authored here
		// beside them: this run executed the build phase, and the attestation is the record
		// saying the binary now on disk is the one it produced. A run that skipped the build
		// authors nothing, leaving the seal and the attestation it inherited untouched.
		if err := attestExecutedBuild(scoping, storageRoot, recordedAt); err != nil {
			fmt.Fprintln(stderr, "gate evidence persistence failed")
			return Result{GateExit: rc, ActionExit: 1, Inspection: inspectAt(storageRoot, engine.Now())}
		}
		// A narrowed run retains nothing here, which is why only the full
		// branch below writes evidence. The verdict cache refuses to reuse a narrow
		// record by its class, but this store has no class at the reuse call site —
		// reusableEvidence asks only whether a green is retained for this tree and
		// oracle — so a whole-tree record written by a run that skipped components
		// would answer a release-path reuse for the next hour with phases nobody ran.
		// Stripped evidence is the same door one step over: writing it from a narrow
		// run would re-stamp an ever-older ancestor as recent.
		if !scoping.partial() {
			if err := retainGreen(storageRoot, plan, recordedAt); err != nil {
				fmt.Fprintln(stderr, "gate evidence persistence failed")
				return Result{GateExit: rc, ActionExit: 1, Inspection: inspectAt(storageRoot, engine.Now())}
			}
			// A full green is the sole author of the ancestor slot: the same run's
			// evidence retained under the stripped identity, which is what a later
			// allowlist-confined tree resolves. Failing to author it only costs a
			// future full run, but failing silently would hide that cost, so it
			// shares the persistence-failure posture above.
			if ordinary {
				strippedPlan, err := buildStrippedSubject(runtimeRoot)
				if err == nil {
					err = retainGreen(storageRoot, strippedPlan, recordedAt)
				}
				if err != nil {
					fmt.Fprintln(stderr, "gate evidence persistence failed")
					return Result{GateExit: rc, ActionExit: 1, Inspection: inspectAt(storageRoot, engine.Now())}
				}
			}
		}
	} else {
		if err := invalidateEvidence(storageRoot, plan); err != nil {
			return operationalWithEngine(engine, storageRoot, rc, stderr, "gate evidence invalidation failed")
		}
		if err := invalidateExecutedComponentSlots(storageRoot, scoping); err != nil {
			return operationalWithEngine(engine, storageRoot, rc, stderr, "gate evidence invalidation failed")
		}
		// A full red also contradicts any ancestor sharing this tree's stripped
		// identity — the `--fresh` red after a capture-only edit — so that slot goes
		// with it.
		if ordinary {
			if strippedPlan, err := buildStrippedSubject(runtimeRoot); err == nil {
				if err := invalidateEvidence(storageRoot, strippedPlan); err != nil {
					return operationalWithEngine(engine, storageRoot, rc, stderr, "gate evidence invalidation failed")
				}
			}
		}
	}
	if err := durableReplaceWithEngine(engine, gitdir, ready); err != nil {
		_ = durableReplaceWithEngine(engine, gitdir, pending)
		fmt.Fprintln(stderr, "gate final persistence failed")
		return Result{GateExit: rc, ActionExit: 1, Inspection: inspectAt(storageRoot, engine.Now())}
	}
	return Result{GateExit: rc, ActionExit: rc, Inspection: inspectAt(storageRoot, engine.Now())}
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
