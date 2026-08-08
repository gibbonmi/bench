# Adopt t.Parallel across eligible gate tests

Blocked by: retire-fixture-kit-pins.md
Ownership fence: `internal/gate`
Integration surfaces: pin-free fixture construction→retire-fixture-kit-pins.md; serial-list enumeration→this ticket's build evidence, graded by the lifecycle's review phase
Contracts: pin-free fixture construction crosses retire-fixture-kit-pins.md→`internal/gate`, with the fixture kit carried as an explicit path string before parallel adoption, absence represented by no ambient kit override, and TP1 asserting the boundary through the hostile-environment package run and serial-list review
Closure: TP1/serial-list, TP1/global-swap, TP1/lock-registry, TP1/race, TP1/narrow-width, TP1/wall-median

## What to build

Every `internal/gate` test satisfying the spec's structural eligibility
predicate calls `t.Parallel`. Ineligible — and enumerated in the build
evidence with each test's pinning reason — is any test that mutates
process-global state: the environment including PATH (`t.Setenv`,
`os.Setenv`/`os.Unsetenv`), the working directory (`t.Chdir`/`os.Chdir`), a
package-level production variable (the gate timeout swap, the builtin-table
stub), the shared execution-lock owners map, or process lifecycle (both
`os.Exit` re-exec families — the phases-command helpers and the prospective
helpers). The representative pinned entry tests from the blocker ticket join
the list with their reason carried over. Go's runner fences only the env and
chdir cases; the variable swaps restore via `t.Cleanup` and corrupt concurrent
tests silently, so the predicate — not the runner — decides the list. The
enumeration is graded against the predicate by the spec build's review phase,
never self-certified. No adopted test may require real overlap: all stay
correct at `-parallel=1` and `GOMAXPROCS=2`.

Build evidence recorded with the ticket: the enumerated serial list; three
focused `go test -count=1 ./internal/gate` repetitions before and after on the
same host with the after-median at or below 90 s; one focused
`go test -race -count=1 ./internal/gate` green; one
`GOMAXPROCS=2 go test -count=1 -timeout 600s ./internal/gate` green.

## Acceptance

- [ ] [TP1] (covers KC5) every eligible test calls `t.Parallel`, the ineligible tests are enumerated with reasons for the review phase to grade, and the package is green under repeated `-count=1`, focused `-race`, and `GOMAXPROCS=2`, with the three-run focused median ≤ 90 s on the reference host.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| TP1/serial-list | add `t.Parallel` to a `t.Setenv`-pinning test | the Go runner's own fence | apply, run `go test -count=1 -run <that test> ./internal/gate`, expect the deterministic panic red |
| TP1/global-swap | add `t.Parallel` to the gate-timeout-swapping test | the package run | apply, run `go test -count=1 -timeout 600s ./internal/gate`, expect concurrent real-gate tests red under the 50 ms timeout |
| TP1/lock-registry | add `t.Parallel` to the test that mutates the execution-lock owners map | the package run plus the focused race run | apply, run `go test -race -count=1 -timeout 600s ./internal/gate`, expect the lock-acquisition tests red or the detector red |
| TP1/race | share one mutable buffer between two adopted tests | the focused race run | apply, run `go test -race -count=1 -timeout 600s ./internal/gate`, expect the detector red |
| TP1/narrow-width | make one adopted test block on a signal another parallel test must send | the narrow-width run | apply, run `GOMAXPROCS=2 go test -count=1 -parallel 1 -timeout 600s ./internal/gate`, expect the bounded timeout red |
| TP1/wall-median | revert adoption on every test at or above one second in the ticket's recorded before-measurement | the recorded three-run measurement | apply, run the three focused repetitions, expect the median above 90 s |
