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

	"github.com/gibbonmi/bench/internal/env"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/structure"
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
// view of the worktree root, the adapter child, and the once-only teardown through it.
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

	mu          sync.Mutex
	adapter     *exec.Cmd // the in-flight adapter child, or nil between runs
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
	attempted := 0
	for r := 1; r <= rcap; r++ {
		s.checkpoint()
		flagged, violations := s.touchedViolations(base)
		if violations == 0 {
			break
		}
		attempted = r
		fmt.Fprintf(s.stdout, "── refactor %d/%d ──\n", r, rcap)
		pre := dirtyPaths(s.root)
		s.runAdapter(fmt.Sprintf(refactorPrompt, flagged))
		s.checkpoint()
		post := dirtyPaths(s.root)
		if s.runGate() == 0 {
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
			fmt.Fprintf(s.stdout, "  ✓ tests green - refactor %d committed\n", r)
		} else {
			s.checkpoint()
			fmt.Fprintln(s.stdout, "  ✗ refactor broke the gate — rolling back")
			rollback(s.root)
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

// runAdapter invokes the harness adapter with the prompt written to its stdin and no
// positional argument, with BENCH_SHIFT=1 armed, from the worktree root. BENCH_SHIFT=1
// arms the Stop hook so the agent cannot declare done on red. Stdin transport keeps the
// prompt out of the machine's process listing on the hop Bench controls. The child runs
// in its own process group, so a pulled line can tear down the whole adapter tree, not
// just the immediate child. The returned error, a spawn failure or a nonzero exit, is
// evidence for progress, not the oracle. The gate still decides whether an iteration's
// work counts.
func (s *session) runAdapter(prompt string) error {
	adapterEnv, err := env.Build(s.root)
	if err != nil {
		fmt.Fprintln(s.stderr, err)
		return err
	}
	cmd := exec.Command(s.agent)
	cmd.Dir = s.root
	cmd.Env = append(adapterEnv, "BENCH_SHIFT=1")
	cmd.Stdin = strings.NewReader(prompt)
	cmd.Stdout, cmd.Stderr = s.stdout, s.stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	s.mu.Lock()
	s.adapter = cmd
	s.mu.Unlock()
	startErr := cmd.Start()
	var runErr error
	if startErr != nil {
		runErr = startErr
	} else {
		runErr = cmd.Wait()
	}
	s.mu.Lock()
	s.adapter = nil
	s.mu.Unlock()
	return runErr
}

// killAdapter signals the in-flight adapter's whole process group, so a pulled line
// reaches the adapter and everything it spawned. It is a no-op when no adapter is running.
func (s *session) killAdapter(sig syscall.Signal) {
	s.mu.Lock()
	cmd := s.adapter
	s.mu.Unlock()
	if cmd != nil && cmd.Process != nil {
		_ = syscall.Kill(-cmd.Process.Pid, sig)
	}
}

// runGate runs the gate against the session's worktree and reports its exit code
// straight through, doing nothing to the session itself. Both call sites share this one
// implementation. The main loop's call propagates a red result through the evidence-split
// preservation path (evidenceResult). The refactor probe's call instead rolls a red
// result back, by design. Preservation itself happens once, explicitly, at each caller's
// own return site, never implied by a flag this method sets on its way out.
func (s *session) runGate() int {
	ctx, cancel := context.WithCancel(context.Background())
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
	os.Exit(finish(s.stdout, s.stderr, s.mainRoot, s.entry, res))
}

// teardown removes the shift scratch and releases the pool lease, exactly once whether
// reached by the normal path, the deferred cleanup, or preserveAndRecover. It returns
// the first error hit, today only the injectable teardown fault, since worktree.Release
// does not itself report git failures. A real teardown problem is reported rather than
// silently swallowed. A second call, the deferred safety net once preserve is already
// set, is a no-op via sync.Once and returns nil.
func (s *session) teardown() error {
	var err error
	s.teardownOne.Do(func() {
		cleanupScratch(s.root)
		if ferr := hitShift(shiftFault, stepTeardown); ferr != nil {
			err = ferr
			return
		}
		worktree.Release(s.root)
	})
	return err
}

// preserveAndRecover is the one place a post-mutation failure preserves work before this
// session's charged worktree leaves the process's hands. When nothing beyond scratch is
// dirty, it releases through teardown and reports RecoveryNone. Otherwise it retains and
// locks the worktree, leaving the dirty tree exactly where it is and never running
// teardown/Release. The work stays visible at a path the operator can read, rather than
// in a ref no command hands back. It always prints the resulting location and marks the
// session so the deferred cleanup never re-finalizes this worktree.
func (s *session) preserveAndRecover(reason string) (recovery string, teardownErr error) {
	defer s.preserve.Store(true)
	if len(dirtyPaths(s.root)) == 0 {
		return RecoveryNone, s.teardown()
	}
	if lockErr := worktree.RetainAndLock(s.root, "bench shift recovery: "+reason); lockErr != nil {
		fmt.Fprintf(s.stderr, "could not retain worktree %s: %v\n", s.root, lockErr)
	}
	fmt.Fprintf(s.stdout, "  recovery: %s\n", recoveryWorktree(s.root))
	return recoveryWorktree(s.root), nil
}

// teardownFailureResult is the one Result a teardown error resolves to, regardless of
// what post-mutation path reached it. Its outcome is failed/1, with a detail naming
// what is already safe, the branch, and the recovery pointer when one exists. The
// teardown failure is real even when the work is not lost.
func teardownFailureResult(s *session, recovery string, err error) Result {
	detail := fmt.Sprintf("teardown failed: %v; branch %s is safe", err, s.branch)
	if recovery != "" && recovery != RecoveryNone {
		detail += "; recovery " + recovery
	}
	return Result{Outcome: OutcomeFailed, Branch: s.branch, Committed: s.committed, IterationsUsed: s.iterationsUsed, Recovery: recovery, Detail: detail}
}

// nothingStaged reports whether the index has no staged changes, the "gate green, no
// change this iteration" signal.
func nothingStaged(root string) bool {
	return exec.Command("git", "-C", root, "diff", "--cached", "--quiet").Run() == nil
}

// rollback discards a red iteration's work while preserving the shift scratch, so the
// next iteration still reads what the last one learned.
func rollback(root string) {
	_ = exec.Command("git", "-C", root, "reset", "-q", "--hard").Run()
	_ = exec.Command("git", "-C", root, "clean", "-qfdx", "-e", ".bench-objective", "-e", ".bench-notes.md").Run()
}
