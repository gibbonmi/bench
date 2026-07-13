package shift

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/structure"
	"github.com/gibbonmi/bench/internal/toon"
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
// view of the worktree root, the adapter child, and the once-only teardown.
type session struct {
	agent  string
	root   string
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

	mu          sync.Mutex
	adapter     *exec.Cmd // the in-flight adapter child, or nil between runs
	cancelGate  context.CancelFunc
	interrupted atomic.Bool // set by the signal handler; parks the loop at its checkpoints
	teardownOne sync.Once
	preserve    bool // a failed oracle transaction retains the charged worktree verbatim
}

// Loop runs the gated shift: preflight the adapter, acquire a pooled worktree, branch,
// iterate (commit on green, preserve on any oracle failure) to the objective or the iteration cap,
// pay down touched-scope structural debt at green, then release. Acquire → loop →
// release run in one process because lease ownership is this process's pid. Returns 0
// on a completed shift, 1 on a preflight/setup failure; a SIGINT/SIGTERM cancels the
// running child, tears the run down, and exits 130.
func Loop(objective string, stdin io.Reader, stdout, stderr io.Writer) int {
	if err := requireAdapter(os.Getenv("BENCH_AGENT")); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	mainRoot, err := git.Root()
	if err != nil {
		fmt.Fprintln(stderr, toon.NotInRepo())
		return 1
	}
	// Audit #10 — tolerate: an empty parse reads as a clean tree, but the very next
	// `rev-parse HEAD` fails the loop loudly on a broken repo, so no broken repo slips past.
	if dirty, _ := git.Output("-C", mainRoot, "status", "--porcelain"); dirty != "" {
		fmt.Fprintln(stderr, "working tree not clean; commit or stash first")
		return 1
	}
	base, err := git.Output("-C", mainRoot, "rev-parse", "HEAD")
	if err != nil {
		fmt.Fprintln(stderr, "could not resolve HEAD")
		return 1
	}
	intentEntry := intent.NewEntry(intent.KindShift, objective)
	if err := intent.Upsert(mainRoot, intentEntry); err != nil {
		fmt.Fprintf(stderr, "could not persist shift intent: %v\n", err)
		return 1
	}
	wt, err := worktree.Acquire(mainRoot, base, "hard")
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	s := &session{agent: os.Getenv("BENCH_AGENT"), root: wt, stdin: stdin, stdout: stdout, stderr: stderr}
	branch := "bench/shift-" + timeNow().Format("20060102-150405")
	if err := exec.Command("git", "-C", wt, "switch", "-q", "-c", branch).Run(); err != nil {
		fmt.Fprintf(stderr, "could not create shift branch %s: %v\n", branch, err)
		s.teardown()
		return 1
	}
	intentEntry.Worktree = wt
	intentEntry.Branch = branch
	if err := intent.Upsert(mainRoot, intentEntry); err != nil {
		fmt.Fprintf(stderr, "could not enrich shift intent: %v\n", err)
		s.teardown()
		return 1
	}
	// The true review base for this branch: `bench diff` resolves it from here, and
	// worktrees share repo config so the key is visible wherever review runs.
	if err := exec.Command("git", "-C", wt, "config", "branch."+branch+".benchBase", base).Run(); err != nil {
		fmt.Fprintf(stderr, "could not configure shift branch %s: %v\n", branch, err)
		s.teardown()
		return 1
	}
	if err := os.WriteFile(wt+"/.bench-objective", []byte(objective+"\n"), 0o644); err != nil {
		fmt.Fprintf(stderr, "could not write shift objective: %v\n", err)
		s.teardown()
		return 1
	}
	if err := os.WriteFile(wt+"/.bench-notes.md", nil, 0o644); err != nil {
		fmt.Fprintf(stderr, "could not write shift notes: %v\n", err)
		s.teardown()
		return 1
	}

	// Signal handling: a pulled line cancels the running child. The loop exits at its
	// checkpoints so cleanup never races git work or an in-flight gate.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)
	go func() {
		sig := <-sigCh
		s.interrupted.Store(true)
		if sys, ok := sig.(syscall.Signal); ok {
			s.killAdapter(sys)
		} else {
			s.killAdapter(syscall.SIGTERM)
		}
		s.cancelRunningGate()
	}()
	defer func() {
		if !s.preserve {
			s.teardown()
		}
	}()

	maxIters := envInt("BENCH_MAX_ITERS", 12)
	fmt.Fprintf(stdout, "▶ shift on %s — objective: %s\n", branch, objective)
	fmt.Fprintf(stdout, "  worktree: %s\n", wt)
	fmt.Fprintf(stdout, "  cap: %d iterations. Ctrl-C to pull the line.\n", maxIters)
	started := time.Now()

	committed := 0
	for i := 1; i <= maxIters; i++ {
		s.checkpoint()
		fmt.Fprintf(stdout, "── iteration %d/%d ──\n", i, maxIters)
		pre := dirtyPaths(wt)
		s.runAdapter(fmt.Sprintf(iterationPrompt, objective))
		s.checkpoint()
		post := dirtyPaths(wt)
		if s.runGate() == 0 {
			s.checkpoint()
			stageTouched(wt, pre, post)
			if nothingStaged(wt) {
				fmt.Fprintln(stdout, "  gate green, no change this iteration — objective likely met.")
				break
			}
			if err := exec.Command("git", "-C", wt, "commit", "-q", "-m", fmt.Sprintf("shift: iteration %d — %s", i, objective)).Run(); err != nil {
				fmt.Fprintf(stderr, "could not commit iteration %d: %v\n", i, err)
				return 1
			}
			committed++
			fmt.Fprintf(stdout, "  ✓ green — committed iteration %d\n", i)
			if objectiveMet(wt, objective) {
				fmt.Fprintln(stdout, "  objective met.")
				break
			}
		} else {
			s.checkpoint()
			s.preserve = true
			fmt.Fprintf(stdout, "  ✗ gate failed — preserving iteration %d in %s\n", i, wt)
			return 1
		}
	}

	if rc := s.refactorPhase(base); rc != 0 {
		return rc
	}

	s.teardown()
	fmt.Fprintf(stdout, "■ shift done: %s, %d committed iteration(s), %dm elapsed\n", branch, committed, int(time.Since(started).Minutes()))
	fmt.Fprintf(stdout, "  review: git -C %s log --oneline %s..%s\n", mainRoot, base, branch)
	fmt.Fprintln(stdout, "  the merge is yours.")
	return 0
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
// runs within the BENCH_REFACTOR_ITERS budget, scopes each prompt to the flagged files,
// and stops on a no-op pass.
func (s *session) refactorPhase(base string) int {
	if _, violations := s.touchedViolations(base); violations == 0 {
		return 0
	}
	fmt.Fprintln(s.stdout, "▶ structure over budget — refactor phase (split at green, not before)")
	rcap := envInt("BENCH_REFACTOR_ITERS", 4)
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
				return 1
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
	return 0
}

// runAdapter invokes the harness adapter with the prompt as its single positional
// argument, BENCH_SHIFT=1 armed (which arms the Stop hook so the agent cannot declare
// done on red), from the worktree root. The child runs in its own process group so a
// pulled line can tear down the whole adapter tree, not just the immediate child. The
// adapter's exit status is ignored — the gate, not the adapter, decides an iteration.
func (s *session) runAdapter(prompt string) {
	cmd := exec.Command(s.agent, prompt)
	cmd.Dir = s.root
	cmd.Env = append(os.Environ(), "BENCH_SHIFT=1")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = s.stdin, s.stdout, s.stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	s.mu.Lock()
	s.adapter = cmd
	s.mu.Unlock()
	if err := cmd.Start(); err == nil {
		_ = cmd.Wait()
	}
	s.mu.Lock()
	s.adapter = nil
	s.mu.Unlock()
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

// checkpoint exits when a signal has been caught, after the running adapter or gate
// has been signaled. This is the well-defined point (mirroring bash's
// trap-between-commands) at which an interrupt takes effect.
func (s *session) checkpoint() {
	if s.interrupted.Load() {
		s.teardown()
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
