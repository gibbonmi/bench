package shift

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/structure"
	"github.com/gibbonmi/bench/internal/worktree"
)

// iterationPrompt is the text a shift iteration hands its adapter as the single
// positional argument. It has always lived inside the executable (a heredoc in the
// shell); it is reviewer-facing content only through the running loop, never a tunable
// file. %s is the objective.
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
// structure report naming those files — never repo-wide debt.
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

// session carries a running shift's state so the signal handler and the loop share one
// view of the worktree root, the adapter child, and the once-only teardown. It also
// carries the fields the interrupt checkpoint needs to emit an honest shift_result and
// record the intent outcome itself — os.Exit skips deferred cleanup, so nothing at that
// path can rely on a defer.
type session struct {
	agent  string
	root   string
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

	mainRoot string
	branch   string
	entry    *intent.Entry

	mu          sync.Mutex
	adapter     *exec.Cmd // the in-flight adapter child, or nil between runs
	cancelGate  context.CancelFunc
	interrupted atomic.Bool // set by the signal handler; parks the loop at its checkpoints
	teardownOne sync.Once
	preserve    atomic.Bool // a failed oracle transaction retains the charged worktree verbatim

	committed      int // main-loop commits this shift; the "committed" evidence for the taxonomy
	iterationsUsed int // the highest iteration number this shift entered
}

// touchedViolations reads this shift's touched-scope structure result — the flagged-files
// string and the violation count — tolerating a git-query failure as an empty scope (zero
// violations). The tolerance is deliberate and single-sourced here for all three refactor-
// gate reads: the shift loop's own `bench gate` run is this worktree's loud oracle, so a
// broken diff degrades the refactor gate rather than crashing the loop; `bench structure`
// is the loud-error path for the same query.
func (s *session) touchedViolations(base string) (flagged string, violations int) {
	flagged, violations, _ = structure.Touched(s.root, base)
	return flagged, violations
}

// refactorPhase pays down structural debt this shift touched, but only once the touched
// scope is over budget — never pre-existing debt, and never mid-implementation. It
// runs within the rcap budget, scopes each prompt to the flagged files, and stops on a
// no-op pass. A refactor commit failure is returned as an error rather than swallowed,
// so the caller can map it onto the outcome taxonomy honestly.
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
			stageTouched(s.root, pre, post)
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

// runAdapter invokes the harness adapter with the prompt as its single positional
// argument, BENCH_SHIFT=1 armed (which arms the Stop hook so the agent cannot declare
// done on red), from the worktree root. The child runs in its own process group so a
// pulled line can tear down the whole adapter tree, not just the immediate child. The
// returned error — a spawn failure or a nonzero exit — is evidence for progress, not the
// oracle: the gate still decides whether an iteration's work counts.
func (s *session) runAdapter(prompt string) error {
	cmd := exec.Command(s.agent, prompt)
	cmd.Dir = s.root
	cmd.Env = append(os.Environ(), "BENCH_SHIFT=1")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = s.stdin, s.stdout, s.stderr
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
// reaches the adapter and everything it spawned. A no-op when no adapter is running.
func (s *session) killAdapter(sig syscall.Signal) {
	s.mu.Lock()
	cmd := s.adapter
	s.mu.Unlock()
	if cmd != nil && cmd.Process != nil {
		_ = syscall.Kill(-cmd.Process.Pid, sig)
	}
}

func (s *session) runGate() int {
	ctx, cancel := context.WithCancel(context.Background())
	s.mu.Lock()
	s.cancelGate = cancel
	s.mu.Unlock()
	rc := gate.RunAndRecordContext(ctx, s.root, s.stdout, s.stderr)
	if ctx.Err() != nil {
		s.preserve.Store(true)
	}
	s.mu.Lock()
	s.cancelGate = nil
	s.mu.Unlock()
	cancel()
	return rc
}

// runPreservingGate records a failed iteration's ownership before returning control
// to the caller, so an interrupt at the next checkpoint cannot release its worktree.
// Refactor probes use runGate directly because their red result is rolled back.
func (s *session) runPreservingGate() int {
	rc := s.runGate()
	if rc != 0 {
		s.preserve.Store(true)
	}
	return rc
}

func (s *session) cancelRunningGate() {
	s.mu.Lock()
	cancel := s.cancelGate
	s.mu.Unlock()
	if cancel != nil {
		s.preserve.Store(true)
		cancel()
	}
}

// checkpoint exits when a signal has been caught, after the running adapter or gate
// has been signaled. This is the well-defined point (mirroring bash's
// trap-between-commands) at which an interrupt takes effect. os.Exit skips deferred
// cleanup, so the shift_result block and the intent outcome are emitted and recorded
// explicitly here, never left to a defer.
func (s *session) checkpoint() {
	if s.interrupted.Load() {
		if !s.preserve.Load() {
			s.teardown()
		}
		res := Result{
			Outcome:        OutcomeInterrupted,
			Branch:         s.branch,
			Committed:      s.committed,
			IterationsUsed: s.iterationsUsed,
			Recovery:       "none",
			Detail:         "interrupted by signal",
		}
		res.Emit(s.stdout)
		if s.entry != nil {
			s.entry.Outcome = string(res.Outcome)
			_ = intent.Upsert(s.mainRoot, *s.entry)
		}
		os.Exit(130)
	}
}

// teardown removes the shift scratch and releases the pool lease, exactly once whether
// reached by the normal path, the deferred cleanup, or the signal handler.
func (s *session) teardown() {
	s.teardownOne.Do(func() {
		cleanupScratch(s.root)
		worktree.Release(s.root)
	})
}

// nothingStaged reports whether the index has no staged changes — the "gate green, no
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
