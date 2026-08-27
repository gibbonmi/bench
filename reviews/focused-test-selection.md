# Focused test selection repair review

Frozen base: `0dae83df5372727d608c6fb2b29fa118d732562f`

Prior reviewed tip: `38074d760fc7204a8cf17995d1af69fe527f1bef`

Reviewed repair tip: `dd049e235aa7f9643105c5d5d975da3edfe14683`

Raw findings: 8. De-duplicated repair targets: 6. The three cancellation
findings share one production-and-oracle repair.

## Standards

Count: 2. Worst issue: P2 duplicated parking fixture.

- **P2 — auto-fix — single-source the parking script.** The one-source rule in
  `AGENTS.md` applies to fixture harnesses. `internal/testreport/cancel_test.go:125-139`
  and lines 238-253 both build the same background sleep, atomic PID publication,
  and wait script. Extract one script source for both target writers.

- **P2 — auto-fix — map the legacy terminator to F01.** Ticket 06 maps the
  legacy positional route to F04 at
  `specs/focused-test-selection/tickets/06-complete-hostile-input-matrices.md:22`.
  The acceptance map defines legacy positional compatibility in F01 at
  `spec.md:240`; F04 at line 243 concerns zero run events.

## Spec

Count: 2. Worst issue: P1 changed graph loading leaks ambient controls.

- **P1 — auto-fix — select once before the first Go hop and close its
  environment.** `spec.md:157-160` requires every mode to remove ambient
  controls and derive the selected binary and `BENCH_KIT` from one selection.
  `internal/testreport/command.go:116-130` resolves changed packages before the
  selection. `internal/testreport/selection.go:86-92` starts `go list` without
  an owned environment. Pass the exact selected environment to graph loading and
  retain the empty-subject posture.

- **P1 — auto-fix — make `go test` cancellation observable after decode
  completes.** N02 at `spec.md:257` requires process-group cancellation.
  `internal/testreport/command.go:209-218` can receive the decoded report first
  and then block in `cmd.Wait`; context cancellation cannot reach that branch.
  Wait for process completion and cancellation concurrently, then drain the
  group and return the common interruption result.

## Coverage

Count: 4. Worst issue: P1 the cancellation oracle misses the decoded-before-wait
state.

- **P2 — auto-fix — make both rename halves independently red-capable.**
  `TestChangedRenameFeedsDeletionAndAdditionToSelection` at
  `selection_test.go:193-220` uses two Go paths in one surviving package. Either
  half can disappear while the final package result stays green. Give each half
  a distinct observable classification or result.

- **P2 — auto-fix — exercise space and glob paths through the focused seam.**
  `spec.md:264-265` requires package and embed paths to remain typed.
  `selection_test.go:134-145` calls `selectCurrentPackages` with handcrafted
  paths, so Git-subject or inspection corruption stays green. Use a real subject
  through `Command` or `resolveChangedPackages` for both path classes.

- **P2 — auto-fix — require the final drain in both Go-hop oracles.** The fake
  Go at `cancel_test.go:125-139` uses an ordinary sleep. The initial group SIGINT
  kills it, so removal of `drainGoProcessGroup` stays green. Make the descendant
  survive the graceful signal and prove the final drain removes it for both hops.

- **P1 — auto-fix — cover decode completion before cancellation.** The same fake
  Go holds stdout open and always takes the current context branch. Add a test
  mode that emits valid terminal JSON, closes stdout, and then parks a
  SIGINT-resistant descendant. It must interrupt without a hang or partial table.
