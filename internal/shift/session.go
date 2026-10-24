package shift

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/env"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/otelrecord"
	"github.com/gibbonmi/bench/internal/structure"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
	"github.com/gibbonmi/bench/internal/worktree"
)

// iterationPrompt is the text a shift iteration writes to its adapter's stdin. It is
// immutable reviewer-facing content compiled into the executable. %s is the objective.
// The prompt travels on stdin, never argv, so it never appears in a process listing.
const iterationPrompt = `You are one iteration of a Bench shift. Objective: %s
First read .bench-notes.md for what prior iterations learned, did, and left
unfinished. Then make ONE small, self-contained change toward the objective, at
the pre-agreed seams. Read the spec under specs/ and projects/ if present. Do not
try to finish everything; advance it by one honest step. Do not weaken or skip any
gate check. Before you stop, append 2–4 lines to .bench-notes.md: what you changed,
what you learned, and the next step you'd take. Then stop — the gate, not you,
decides if it counts.
`

// refactorPrompt scopes the refactor phase to the files this shift flagged. %s is the
// structure report naming those files, never repo-wide debt.
const refactorPrompt = `The implementation is complete and tests are green, but the structure budget is
exceeded. These are the flagged files and directories this shift touched — fix only
these, nothing else:

%s

Fix them by splitting along responsibility, using the deletion test from the craft-seams
skill: lift a cluster out only if extracting it *concentrates* complexity behind a real
interface rather than just moving it. Never fragment a cohesive file to beat the line
count — if a file is genuinely one deep module, leave it and say so. Group a crowded
directory into a package with a clear entry point. Keep every test green; change
structure, not behavior. Make one split, then stop — the loop re-checks and continues.
`

var timeNow = time.Now

// session carries a running shift's state. The signal handler and the loop share one
// view of the worktree root, the adapter lifetime, and the once-only teardown through it.
// It also carries the fields the interrupt checkpoint needs to emit an honest
// shift_result and record the intent outcome itself. os.Exit skips deferred cleanup, so
// nothing at that path can rely on a defer.
type session struct {
	agent  string
	root   string
	stdout io.Writer
	stderr io.Writer

	mainRoot string
	branch   string
	entry    *intent.Entry
	record   *shiftRecord
	released bool // set once teardown has released the worktree to the pool

	mu sync.Mutex
	// adapter is the lifetime of the in-flight adapter, or nil between runs. runAdapter
	// stores it before Start, so a cancellation never waits on a published process.
	adapter *subprocess.Lifetime
	// cancel is the first shutdown request of the session. A request that arrives between
	// adapter runs reaches the next lifetime before its Start.
	cancel      subprocess.Cancellation
	cancelGate  context.CancelFunc
	interrupted atomic.Bool // set by the signal handler; parks the loop at its checkpoints
	deadline    atomic.Bool // set by the wall timer; parks the loop like interrupted, but reads as incomplete
	teardownOne sync.Once
	// preserve marks that this session's charged worktree has already been finalized,
	// either snapshotted and released or retained and locked, by preserveAndRecover. Every
	// post-mutation exit path calls preserveAndRecover explicitly and sets this itself. The
	// deferred cleanup in Loop only runs teardown when this is still false. A retained
	// worktree is never re-cleaned, and a released one is never released twice.
	preserve atomic.Bool

	committed      int // main-loop commits this shift; the "committed" evidence for the taxonomy
	iterationsUsed int // the highest iteration number this shift entered
}

// touchedViolations reads this shift's touched-scope structure result, the flagged-files
// string and the violation count, tolerating a git-query failure as an empty scope with
// zero violations. The tolerance is deliberate and single-sourced here for all three
// refactor-gate reads. The shift loop's own `bench gate` run is this worktree's loud
// oracle, so a broken diff degrades the refactor gate rather than crashing the loop.
// `bench structure` is the loud-error path for the same query.
func (s *session) touchedViolations(base string) (flagged string, violations int) {
	flagged, violations, _ = structure.Touched(s.root, base)
	return flagged, violations
}

// refactorPhase pays down structural debt this shift touched, but only once the touched
// scope is over budget, never pre-existing debt, and never mid-implementation. It runs
// within the rcap budget, scopes each prompt to the flagged files, and stops on a no-op
// pass. A refactor commit failure is returned as an error rather than swallowed, so the
// caller can map it onto the outcome taxonomy honestly.
func (s *session) refactorPhase(base string, rcap int) error {
	if _, violations := s.touchedViolations(base); violations == 0 {
		return nil
	}
	fmt.Fprintln(s.stdout, "▶ structure over budget — refactor phase (split at green, not before)")
	defer s.record.endPass()
	attempted := 0
	for r := 1; r <= rcap; r++ {
		s.checkpoint()
		flagged, violations := s.touchedViolations(base)
		if violations == 0 {
			break
		}
		attempted = r
		passCtx := s.record.beginPass(refactorSeam)
		fmt.Fprintf(s.stdout, "── refactor %d/%d ──\n", r, rcap)
		pre := dirtyPaths(s.root)
		s.record.adapterRan(s.runAdapter(passCtx, fmt.Sprintf(refactorPrompt, flagged)))
		s.checkpoint()
		post := dirtyPaths(s.root)
		if s.runGate(passCtx) == 0 {
			s.checkpoint()
			if err := stageTouched(s.root, pre, post); err != nil {
				fmt.Fprintf(s.stderr, "could not stage refactor %d: %v\n", r, err)
				return fmt.Errorf("could not stage refactor pass %d", r)
			}
			if nothingStaged(s.root) {
				fmt.Fprintf(s.stdout, "  gate green, refactor %d made no staged change - stopping refactor phase\n", r)
				break
			}
			if err := exec.Command("git", "-C", s.root, "commit", "-q", "-m", "refactor: reduce structural debt").Run(); err != nil {
				fmt.Fprintf(s.stderr, "could not commit refactor %d: %v\n", r, err)
				return fmt.Errorf("could not commit refactor pass %d", r)
			}
			s.record.passCommitted(s.root)
			fmt.Fprintf(s.stdout, "  ✓ tests green - refactor %d committed\n", r)
		} else {
			s.checkpoint()
			fmt.Fprintln(s.stdout, "  ✗ refactor broke the gate — rolling back")
			if err := rollback(s.root); err != nil {
				fmt.Fprintf(s.stderr, "could not roll back refactor %d: %v\n", r, err)
				return fmt.Errorf("could not roll back refactor pass %d", r)
			}
		}
	}
	if _, violations := s.touchedViolations(base); violations == 0 {
		fmt.Fprintln(s.stdout, "  structure back under budget.")
	} else {
		n := attempted
		if n == 0 {
			n = rcap
		}
		fmt.Fprintf(s.stdout, "  ⚠ still over budget after %d refactor pass(es) - review manually, or run a Bench deep pass with bench structure and craft-seams.\n", n)
	}
	return nil
}

// stepAdapterStart is the startup transition of an adapter run: its lifetime is published
// to the cancellation paths, and its Start has not run. A test pauses it there.
const stepAdapterStart shiftStep = "adapter-start"

// adapterPolicy keeps the adapter's result after a normal exit. A cancellation sends the
// caller's signal, and KILL follows the shutdown grace.
var adapterPolicy = subprocess.Policy{
	Normal: subprocess.Preserve,
	Grace:  bounds.ProcessWindow(bounds.FixedWindow(bounds.ProcessShutdownGrace)),
	Final:  bounds.ProcessWindow(bounds.FixedWindow(bounds.ProcessFinalWindow)),
	Poll:   bounds.ProcessPoll,
}

// runAdapter invokes the harness adapter with the prompt written to its stdin and no
// positional argument, with BENCH_SHIFT=1 armed, from the worktree root. BENCH_SHIFT=1
// arms the Stop hook so the agent cannot declare done on red. Stdin transport keeps the
// prompt out of the machine's process listing on the hop Bench controls. The child runs
// in its own process group under one lifetime, so a pulled line can tear down the whole
// adapter tree, not just the immediate child. The lifetime is published before its Start,
// so a cancellation during startup reaches the adapter. The launch registers on the
// inherited resources and the worktree, so an adapter group that outlives the run keeps the
// worktree protected. The returned error, a spawn failure or a nonzero exit, is evidence
// for progress, not the oracle. The gate still decides whether an iteration's work counts.
// The adapter also gets the trace handoff for the pass span on ctx, after every inherited
// value, so its line resolution records under this pass.
func (s *session) runAdapter(ctx context.Context, prompt string) error {
	adapterEnv, err := env.Build(s.root)
	if err != nil {
		fmt.Fprintln(s.stderr, err)
		return err
	}
	adapterEnv = otelrecord.WithHandoff(ctx, s.mainRoot, append(adapterEnv, "BENCH_SHIFT=1"))
	resources, err := worktree.ChildResources(s.root, adapterEnv)
	if err != nil {
		fmt.Fprintln(s.stderr, err)
		return err
	}
	cmd := exec.Command(s.agent)
	cmd.Dir = s.root
	cmd.Env = protection.Environ(adapterEnv, resources)
	streams := subprocess.Streams{Stdin: strings.NewReader(prompt), Stdout: subprocess.Output{Writer: s.stdout}, Stderr: subprocess.Output{Writer: s.stderr}}
	lifetime := subprocess.NewLifetime(cmd, adapterPolicy, streams, subprocess.Registration{Protocol: protection.Installed(), Resources: resources})
	s.mu.Lock()
	s.adapter = lifetime
	pending := s.cancel
	s.mu.Unlock()
	defer func() {
		s.mu.Lock()
		s.adapter = nil
		s.mu.Unlock()
	}()
	if pending.Cause != subprocess.CauseNone {
		lifetime.Cancel(pending.Cause, pending.Signal)
	}
	if err := hitShift(shiftFault, stepAdapterStart); err != nil {
		return err
	}
	if err := lifetime.Start(); err != nil {
		return err
	}
	outcome, _ := lifetime.Wait()
	if unresolved := outcome.Uncertified(); unresolved != nil {
		// The worktree's protection keeps what the adapter group can still change.
		fmt.Fprintf(s.stderr, "adapter: %v\n", unresolved)
	}
	if outcome.Process.StartErr != nil {
		return outcome.Process.StartErr
	}
	return outcome.Process.WaitErr
}

// cancelAdapter requests shutdown of the in-flight adapter's whole process group with
// sig, so a pulled line reaches the adapter and everything it spawned. The session keeps
// the first request, so an adapter run that has not published its lifetime yet receives
// it before its Start.
func (s *session) cancelAdapter(cause subprocess.Cause, sig syscall.Signal) {
	s.mu.Lock()
	if s.cancel.Cause == subprocess.CauseNone {
		s.cancel = subprocess.Cancellation{Cause: cause, Signal: sig}
	}
	adapter := s.adapter
	s.mu.Unlock()
	if adapter != nil {
		adapter.Cancel(cause, sig)
	}
}

// interrupt is the signal path: a pulled line forwards its own signal to the adapter and
// cancels a running gate. The loop exits at its next checkpoint.
func (s *session) interrupt(sig os.Signal) {
	s.interrupted.Store(true)
	sys, ok := sig.(syscall.Signal)
	if !ok {
		sys = syscall.SIGTERM
	}
	s.cancelAdapter(subprocess.CauseInterruption, sys)
	s.cancelRunningGate()
}

// expire is the wall-deadline path. It acts like a pulled line, but it sets deadline
// rather than interrupted, so the next checkpoint resolves incomplete/3.
func (s *session) expire() {
	s.deadline.Store(true)
	s.cancelAdapter(subprocess.CauseDeadline, syscall.SIGTERM)
	s.cancelRunningGate()
}

// runGate runs the gate against the session's worktree and reports its exit code
// straight through, doing nothing to the session itself. Both call sites share this one
// implementation. The main loop's call propagates a red result through the evidence-split
// preservation path (evidenceResult). The refactor probe's call instead rolls a red
// result back, by design. Preservation itself happens once, explicitly, at each caller's
// own return site, never implied by a flag this method sets on its way out. The gate runs
// below ctx, so its span is the child of the pass span on ctx.
func (s *session) runGate(ctx context.Context) int {
	ctx, cancel := context.WithCancel(ctx)
	s.mu.Lock()
	s.cancelGate = cancel
	s.mu.Unlock()
	rc := gate.RunAndRecordContext(ctx, s.root, s.stdout, s.stderr)
	s.mu.Lock()
	s.cancelGate = nil
	s.mu.Unlock()
	cancel()
	return rc
}

func (s *session) cancelRunningGate() {
	s.mu.Lock()
	cancel := s.cancelGate
	s.mu.Unlock()
	if cancel != nil {
		cancel()
	}
}

// checkpointOutcome resolves which exit, if either, a checkpoint takes, tested directly
// since checkpoint itself exits the process. Decided precedence: the wall deadline wins
// over an interrupt when both flags are set. A deadline that fires while a pulled line
// is already in flight still resolves incomplete/3, not interrupted/130. The deadline is
// this shift's own bound, not an external signal. The taxonomy should not depend on
// which of two concurrent cancellations happened to set its flag first.
func (s *session) checkpointOutcome() (outcome Outcome, detail string, ok bool) {
	switch {
	case s.deadline.Load():
		return OutcomeIncomplete, "wall deadline exceeded", true
	case s.interrupted.Load():
		return OutcomeInterrupted, "interrupted by signal", true
	}
	return "", "", false
}

// checkpoint exits when a signal or the wall deadline has fired, after the running
// adapter or gate has already been signaled or cancelled. This is the well-defined point,
// mirroring bash's trap-between-commands, at which an interrupt or deadline takes effect.
// os.Exit skips deferred cleanup, so preservation, the shift_result block, and the intent
// outcome are all run and recorded explicitly here, never left to a defer.
func (s *session) checkpoint() {
	if outcome, detail, ok := s.checkpointOutcome(); ok {
		s.exitPreserving(outcome, detail)
	}
}

// exitPreserving is checkpoint's shared exit path for both a signal and a wall-deadline
// trip. It preserves any dirty work and resolves the outcome; a teardown failure here
// still resolves to failed/1, same as every other exit path. It then hands off to finish
// for the emit-record-exit-code sequence, the same single-sourced path every ordinary
// Loop return uses, before exiting with its resolved code. os.Exit here, rather than a
// return, is deliberate. This path is reached from inside the loop's own call stack via
// checkpoint, and skips the deferred cleanup a normal return would trigger. See the
// session doc comment.
func (s *session) exitPreserving(outcome Outcome, detail string) {
	recovery, teardownErr := s.preserveAndRecover(detail)
	res := Result{Outcome: outcome, Branch: s.branch, Committed: s.committed, IterationsUsed: s.iterationsUsed, Recovery: recovery, Detail: detail}
	if teardownErr != nil {
		res = teardownFailureResult(s, recovery, teardownErr)
	}
	os.Exit(finish(s.stdout, s.stderr, s.mainRoot, s.entry, s.record, res))
}

// nothingStaged reports whether the index has no staged changes, the "gate green, no
// change this iteration" signal.
func nothingStaged(root string) bool {
	return exec.Command("git", "-C", root, "diff", "--cached", "--quiet").Run() == nil
}
