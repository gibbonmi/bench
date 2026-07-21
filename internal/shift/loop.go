package shift

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/worktree"
	refreshop "github.com/gibbonmi/bench/internal/worktree/refresh"
)

// finish emits the shift_result block, records the outcome on the intent entry when one
// exists (a validation failure exits before an entry is created, and that is correct —
// there is nothing yet to enrich), and resolves res to its process exit code. Every
// Loop return, and checkpoint's os.Exit path (via exitPreserving), funnels through this
// one path — the single source for the emit → record → exit-code sequence. A failed
// Upsert does not change the outcome or the exit code (the gate's verdict already
// happened; the ledger record is enrichment, not the oracle) but is not silently
// swallowed either: it is reported to stderr so an operator can see the ledger fell out
// of sync.
func finish(stdout, stderr io.Writer, mainRoot string, entry *intent.Entry, res Result) int {
	res.Emit(stdout)
	if entry != nil {
		entry.Outcome = string(res.Outcome)
		entry.Recovery = res.Recovery
		err := intent.Upsert(mainRoot, *entry)
		if err == nil {
			err = hitShift(shiftFault, stepIntentUpsert)
		}
		if err != nil {
			fmt.Fprintf(stderr, "warning: could not record shift outcome: %v\n", err)
		}
	}
	return res.ExitCode()
}

// usage is the exit-2 shorthand for every setup failure before the first adapter run —
// there is no intent entry yet, so nothing is enriched.
func usage(stdout, stderr io.Writer, detail string) int {
	return finish(stdout, stderr, "", nil, Result{Outcome: OutcomeUsage, Detail: detail})
}

// evidenceResult is the one place a post-mutation failure both preserves the dirty tree
// (snapshot-and-release, or retain-and-lock on a snapshot failure — preserveAndRecover's
// uniform rule) and builds the Result, split by the session's committed count per the
// evidence rule. A teardown failure on the release side still resolves to failed/1,
// regardless of the evidence split, per teardownFailureResult.
func evidenceResult(s *session, detail string) Result {
	recovery, teardownErr := s.preserveAndRecover(detail)
	if teardownErr != nil {
		return teardownFailureResult(s, recovery, teardownErr)
	}
	return Result{
		Outcome:        evidenceOutcome(s.committed),
		Branch:         s.branch,
		Committed:      s.committed,
		IterationsUsed: s.iterationsUsed,
		Recovery:       recovery,
		Detail:         detail,
	}
}

// branchCollisionRetries bounds how many disambiguating suffixes createShiftBranch will
// try before giving up and reporting the creation failure — ten total attempts (the
// bare per-second name, then -2 through -10) is generous headroom for concurrent
// same-second shifts while still failing fast on a genuinely broken repo.
const branchCollisionRetries = 10

// createShiftBranch derives the bench/shift-<timestamp> branch name and switches wt
// onto a freshly created branch of that name. A same-second collision (two shifts
// deriving the same timestamp) is not fatal: it retries with a disambiguating "-2",
// "-3", … suffix — appended to the per-second name, so the recovery ref path, which is
// built from the resolved branch name, gets a fresh, non-colliding pair too — until
// creation succeeds or branchCollisionRetries is exhausted, at which point it reports
// the same creation failure this always reported for an unresolvable collision.
func createShiftBranch(wt, timestamp string) (string, error) {
	base := "bench/shift-" + timestamp
	var lastErr error
	for attempt := 1; attempt <= branchCollisionRetries; attempt++ {
		candidate := base
		if attempt > 1 {
			candidate = fmt.Sprintf("%s-%d", base, attempt)
		}
		if err := exec.Command("git", "-C", wt, "switch", "-q", "-c", candidate).Run(); err != nil {
			lastErr = fmt.Errorf("could not create shift branch %s: %w", candidate, err)
			continue
		}
		return candidate, nil
	}
	return "", lastErr
}

// objectiveBanner formats the shift-start line. The objective is reviewer-authored and
// already rejected at intake if it carries a control byte, but the banner renders it
// through the shared sanitizer anyway — one policy for every terminal render of
// operator-influenced text, no render path with a raw escape sequence behind it. The
// " — objective:" delimiter is load-bearing: the shift-start parser splits the branch on it.
func objectiveBanner(branch, objective string) string {
	return fmt.Sprintf("▶ shift on %s — objective: %s", branch, sanitize.Preview(objective))
}

// Loop runs the gated shift: validate the objective and env, preflight the adapter,
// acquire a pooled worktree, branch, iterate (commit on green, preserve on a red gate)
// to the objective or the iteration cap, pay down touched-scope structural debt at
// green, then release. Acquire → loop → release run in one process because lease
// ownership is this process's pid. Every exit path resolves through finish, which emits
// the shift_result TOON block and records the outcome on the intent entry.
func Loop(objective string, stdout, stderr io.Writer) int {
	return loop(objective, false, stdout, stderr)
}

func loop(objective string, refresh bool, stdout, stderr io.Writer) int {
	if err := validateObjective(objective); err != nil {
		fmt.Fprintln(stderr, err)
		return usage(stdout, stderr, err.Error())
	}
	maxIters, err := parseBoundedInt("BENCH_MAX_ITERS", bounds.MainIterationsDefault, bounds.IterationMin, bounds.IterationMax)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return usage(stdout, stderr, err.Error())
	}
	refactorIters, err := parseBoundedInt("BENCH_REFACTOR_ITERS", bounds.RefactorIterationsDefault, bounds.IterationMin, bounds.IterationMax)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return usage(stdout, stderr, err.Error())
	}
	wallDur, err := parseWallDuration("BENCH_MAX_WALL", 0, bounds.MaxWall)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return usage(stdout, stderr, err.Error())
	}
	if err := requireAdapter(os.Getenv("BENCH_AGENT")); err != nil {
		fmt.Fprintln(stderr, err)
		return usage(stdout, stderr, err.Error())
	}
	mainRoot, err := git.Root()
	if err != nil {
		fmt.Fprintln(stderr, toon.NotInRepo())
		return usage(stdout, stderr, "not in a git repository")
	}
	// Audit #10 — tolerate: an empty parse reads as a clean tree, but the very next
	// `rev-parse HEAD` fails the loop loudly on a broken repo, so no broken repo slips past.
	if dirty, _ := git.Output("-C", mainRoot, "status", "--porcelain"); dirty != "" {
		fmt.Fprintln(stderr, "working tree not clean; commit or stash first")
		return usage(stdout, stderr, "working tree not clean")
	}
	startRef := "HEAD"
	if refresh {
		result := refreshop.Refresh(mainRoot)
		fmt.Fprint(stdout, refreshop.RenderRefresh(result))
		if result.Status == "refreshed" {
			startRef = refreshop.RefreshedStartRef(mainRoot)
		}
	}
	base, err := git.Output("-C", mainRoot, "rev-parse", startRef+"^{commit}")
	if err != nil {
		fmt.Fprintln(stderr, "could not resolve HEAD")
		return usage(stdout, stderr, "could not resolve HEAD")
	}
	intentEntry := intent.NewEntry(intent.KindShift)
	if err := intent.Upsert(mainRoot, intentEntry); err != nil {
		fmt.Fprintf(stderr, "could not persist shift intent: %v\n", err)
		return usage(stdout, stderr, "could not persist shift intent")
	}
	wt, err := worktree.Acquire(mainRoot, base, "hard")
	if err != nil {
		fmt.Fprintln(stderr, err)
		return finish(stdout, stderr, mainRoot, &intentEntry, Result{Outcome: OutcomeUsage, Detail: "could not acquire a worktree"})
	}

	s := &session{agent: os.Getenv("BENCH_AGENT"), root: wt, stdout: stdout, stderr: stderr, mainRoot: mainRoot, entry: &intentEntry}
	branch, err := createShiftBranch(wt, timeNow().Format("20060102-150405"))
	if err != nil {
		fmt.Fprintf(stderr, "%v\n", err)
		s.teardown()
		return finish(stdout, stderr, mainRoot, &intentEntry, Result{Outcome: OutcomeUsage, Detail: err.Error()})
	}
	s.branch = branch
	intentEntry.Worktree = wt
	intentEntry.Branch = branch
	if err := intent.Upsert(mainRoot, intentEntry); err != nil {
		fmt.Fprintf(stderr, "could not enrich shift intent: %v\n", err)
		s.teardown()
		return finish(stdout, stderr, mainRoot, &intentEntry, Result{Outcome: OutcomeUsage, Branch: branch, Detail: "could not enrich shift intent"})
	}
	// The true review base for this branch: `bench diff` resolves it from here, and
	// worktrees share repo config so the key is visible wherever review runs.
	if err := exec.Command("git", "-C", wt, "config", "branch."+branch+".benchBase", base).Run(); err != nil {
		fmt.Fprintf(stderr, "could not configure shift branch %s: %v\n", branch, err)
		s.teardown()
		return finish(stdout, stderr, mainRoot, &intentEntry, Result{Outcome: OutcomeUsage, Branch: branch, Detail: "could not configure shift branch " + branch})
	}
	// 0600: the worktree scratch file is the one place the full objective text persists,
	// so it is readable only by the user who started the shift.
	if err := os.WriteFile(wt+"/.bench-objective", []byte(objective+"\n"), 0o600); err != nil {
		fmt.Fprintf(stderr, "could not write shift objective: %v\n", err)
		s.teardown()
		return finish(stdout, stderr, mainRoot, &intentEntry, Result{Outcome: OutcomeUsage, Branch: branch, Detail: "could not write shift objective"})
	}
	if err := os.WriteFile(wt+"/.bench-notes.md", nil, 0o644); err != nil {
		fmt.Fprintf(stderr, "could not write shift notes: %v\n", err)
		s.teardown()
		return finish(stdout, stderr, mainRoot, &intentEntry, Result{Outcome: OutcomeUsage, Branch: branch, Detail: "could not write shift notes"})
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
		if !s.preserve.Load() {
			_ = s.teardown()
		}
	}()

	// The wall deadline: on expiry it acts like a pulled line — kill the adapter process
	// group and cancel a running gate — but sets deadline rather than interrupted, so the
	// next checkpoint resolves incomplete/3 with a deadline detail, not interrupted/130.
	if wallDur > 0 {
		wallTimer := time.AfterFunc(wallDur, func() {
			s.deadline.Store(true)
			s.killAdapter(syscall.SIGTERM)
			s.cancelRunningGate()
		})
		defer wallTimer.Stop()
	}

	fmt.Fprintln(stdout, objectiveBanner(branch, objective))
	fmt.Fprintf(stdout, "  worktree: %s\n", wt)
	fmt.Fprintf(stdout, "  cap: %d iterations. Ctrl-C to pull the line.\n", maxIters)
	started := time.Now()

	// stopReason distinguishes how the loop left its for-range: "" means it ran to the
	// iteration cap (incomplete), "stopped" means it broke clean (complete, or no-op if
	// nothing landed), "adapter-failed" means an adapter spawn/exit failure ended it.
	var stopReason, stopDetail string
	for i := 1; i <= maxIters; i++ {
		s.checkpoint()
		s.iterationsUsed = i
		fmt.Fprintf(stdout, "── iteration %d/%d ──\n", i, maxIters)
		pre := dirtyPaths(wt)
		adapterErr := s.runAdapter(fmt.Sprintf(iterationPrompt, objective))
		s.checkpoint()
		post := dirtyPaths(wt)
		if s.runGate() == 0 {
			s.checkpoint()
			if err := stageTouched(wt, pre, post); err != nil {
				fmt.Fprintf(stderr, "could not stage iteration %d: %v\n", i, err)
				return finish(stdout, stderr, mainRoot, &intentEntry, evidenceResult(s, fmt.Sprintf("could not stage iteration %d", i)))
			}
			if nothingStaged(wt) {
				if adapterErr == nil {
					if objectiveMet(wt, objective) {
						stopDetail = "completion predicate already satisfied"
					} else {
						stopDetail = "gate green, no further change"
					}
					fmt.Fprintln(stdout, "  gate green, no change this iteration — objective likely met.")
					stopReason = "stopped"
				} else {
					stopDetail = fmt.Sprintf("adapter exited nonzero on iteration %d with no change", i)
					fmt.Fprintf(stdout, "  gate green, no change this iteration, but adapter exited nonzero — stopping.\n")
					stopReason = "adapter-failed"
				}
				break
			}
			if err := exec.Command("git", "-C", wt, "commit", "-q", "-m", fmt.Sprintf("shift: iteration %d — %s", i, objective)).Run(); err != nil {
				fmt.Fprintf(stderr, "could not commit iteration %d: %v\n", i, err)
				return finish(stdout, stderr, mainRoot, &intentEntry, evidenceResult(s, fmt.Sprintf("could not commit iteration %d", i)))
			}
			s.committed++
			fmt.Fprintf(stdout, "  ✓ green — committed iteration %d\n", i)
			if objectiveMet(wt, objective) {
				fmt.Fprintln(stdout, "  objective met.")
				stopDetail = "objective met"
				stopReason = "stopped"
				break
			}
			if adapterErr != nil {
				fmt.Fprintf(stdout, "  adapter exited nonzero after committing iteration %d — stopping.\n", i)
				stopDetail = fmt.Sprintf("adapter exited nonzero after committing iteration %d", i)
				stopReason = "adapter-failed"
				break
			}
		} else {
			s.checkpoint()
			// Neither the human line nor the detail may name the pool worktree as the
			// preservation site: on the snapshot path it is released and cleaned right
			// after. The location is the "recovery:" line preserveAndRecover prints
			// (the ref, or the retained worktree path on the fallback) plus the
			// shift_result recovery cell — never this message.
			fmt.Fprintf(stdout, "  ✗ gate failed — snapshotting iteration %d\n", i)
			return finish(stdout, stderr, mainRoot, &intentEntry, evidenceResult(s, fmt.Sprintf("gate failed on iteration %d", i)))
		}
	}

	if stopReason == "adapter-failed" {
		return finish(stdout, stderr, mainRoot, &intentEntry, evidenceResult(s, stopDetail))
	}

	if err := s.refactorPhase(base, refactorIters); err != nil {
		return finish(stdout, stderr, mainRoot, &intentEntry, evidenceResult(s, err.Error()))
	}

	recovery, teardownErr := s.preserveAndRecover("shift teardown")
	if teardownErr != nil {
		return finish(stdout, stderr, mainRoot, &intentEntry, teardownFailureResult(s, recovery, teardownErr))
	}
	fmt.Fprintf(stdout, "■ shift done: %s, %d committed iteration(s), %dm elapsed\n", branch, s.committed, int(time.Since(started).Minutes()))
	fmt.Fprintf(stdout, "  review: git -C %s log --oneline %s..%s\n", mainRoot, base, branch)
	fmt.Fprintln(stdout, "  the merge is yours.")

	outcome := OutcomeComplete
	detail := stopDetail
	switch {
	case stopReason == "":
		outcome = OutcomeIncomplete
		detail = "iteration cap exhausted"
	case s.committed == 0:
		outcome = OutcomeNoOp
	}
	return finish(stdout, stderr, mainRoot, &intentEntry, Result{
		Outcome: outcome, Branch: branch, Committed: s.committed, IterationsUsed: s.iterationsUsed, Recovery: recovery, Detail: detail,
	})
}
