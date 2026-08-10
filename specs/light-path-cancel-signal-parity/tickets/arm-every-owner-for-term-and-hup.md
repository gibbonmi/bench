# Arm every signal owner for TERM and HUP from one signal source

Blocked by: none
Ownership fence: `internal/subprocess/cancel.go`, `internal/subprocess/cancel_test.go`, `internal/runbinary/runbinary.go`, `internal/testreport/testreport.go`, `internal/testreport/cancel_test.go`, `internal/worktree/exec.go`, `internal/worktree/lifecycle.go`, `internal/gate/gate.go`, `internal/gate/phases.go`, `internal/preflight/command.go`, `internal/preprelease/preprelease.go`, `internal/publication/command.go`, `internal/shift/loop.go`, `internal/freshness/freshness.go`, `cmd/bench/specbuild.go`
Integration surfaces: cancel signal set→`internal/subprocess/cancel.go`; drained builder group→`internal/testreport/cancel_test.go` + CS1; private selection directory→`internal/runbinary/runbinary.go` `Selection.Close` exercised unchanged + CS4
Contracts: the cancel signal set crosses `internal/subprocess`→every converted owner as `subprocess.CancelSignals`, membership exactly `SIGINT, SIGTERM, SIGHUP` in that order and nothing else, asserted by CS5 against the exported value and by CS1–CS3 against a real detached builder group; the cancelled-builder grace crosses `internal/runbinary`→`internal/testreport/cancel_test.go` as `runbinary.BuilderCancelGrace`, so the test deadline derives from the window it must outlast rather than a literal
Closure: CS1/a SIGTERMed owner leaves no process in its builder group, CS2/a SIGHUPed owner leaves no process in its builder group, CS3/a SIGINTed owner still leaves no process in its builder group, CS4/a signalled owner removes its private selection directory, CS5/`CancelSignals` holds exactly the three signals in order
Assumptions: SIGKILL is uncatchable and stays out of this ticket — closing that half needs `Pdeathsig` on the group leader, which changes when a child may outlive its owner and is a reviewer decision, not a parity repair; every converted site already means "cancel this run and drain what it started" by its signal, so widening the set changes which signals reach that path, never what the path does

## What to build

`bench test` loses its builder group on SIGTERM, and the signal set that decides this is
written out thirteen times in three spellings.

`runbinary.canonicalBuild` (`internal/runbinary/runbinary.go:189`) starts
`scripts/go-build.sh` with `Setpgid: true` and drains it only through `drainBuilderGroup`
after `Wait` or on context cancellation. Detaching the group also removes the kernel's
fallback: an orphan in the owner's own group would still receive the terminal's SIGINT or
SIGHUP, and a detached one receives nothing. When the owner dies without reaching its
handler, the builder has no cleanup path left at all.

`bench test` arms that context for `os.Interrupt` alone
(`internal/testreport/testreport.go:43`). Reproduced 2026-08-09: with the builder live,
SIGINT drains the group, while SIGTERM leaves `go-build.sh`, its `go build`, and a
toolchain `compile` alive and still consuming CPU, with the `/tmp/bench-run-*` selection
directory never reclaimed because `Own`'s incomplete-path `Close` never runs. Three such
orphans were found alive after three days. SIGTERM is the signal that matters here: it is
what session and harness teardown sends.

The thirteen production registrations spell the set three ways — `{INT}` at
`internal/testreport/testreport.go:43`, `internal/worktree/exec.go:44`, and
`internal/worktree/lifecycle.go:474`; `{INT, TERM}` at `internal/gate/phases.go:208`,
`internal/preflight/command.go:45`, `internal/preprelease/preprelease.go:169`,
`internal/shift/loop.go:218`, and the three in `internal/publication/command.go`;
`{INT, TERM, HUP}` at `internal/gate/gate.go:356`, `internal/freshness/freshness.go:208`,
and `cmd/bench/specbuild.go:30`. Thirteen derivations of one fact is how three spellings
drifted apart.

Land the set once in `internal/subprocess` — the package that already owns running an
external command — as `CancelSignals` plus a `NotifyCancel` context helper, and move every
production site onto it. `internal/shift/loop.go` and `internal/freshness/freshness.go`
register a channel rather than a context, so they take
`signal.Notify(ch, subprocess.CancelSignals...)`; the set still has one source.

Not every converted site owns a detached group — `internal/worktree/lifecycle.go` guards
plain `git worktree` calls, and the ticket converts it for one-source reasons, not because
it leaks. The group-drain acceptance below is graded at the `bench test` seam, which is the
site the repro actually reds.

## Acceptance

- [ ] [CS1] A `bench test` owner sent SIGTERM while its builder group is live leaves no surviving process in that group.
- [ ] [CS2] The same owner sent SIGHUP leaves no surviving process in that group.
- [ ] [CS3] The same owner sent SIGINT leaves no surviving process in that group, as it did before this ticket.
- [ ] [CS4] A `bench test` owner cancelled by any of the three signals removes the private selection directory it created.
- [ ] [CS5] `subprocess.CancelSignals` holds exactly SIGINT, SIGTERM, SIGHUP in that order, and every production `signal.Notify` and `signal.NotifyContext` call site takes its signals from it.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CS1 | restore `signal.NotifyContext(context.Background(), os.Interrupt)` at `internal/testreport/testreport.go:43` | the SIGTERM leg of `TestCommandDrainsBuilderGroupOnCancelSignal` | revert the call site, run `go test ./internal/testreport -run TestCommandDrainsBuilderGroupOnCancelSignal`, expect the SIGTERM subtest red on the surviving builder pid and the SIGINT subtest still green |
| CS2 | drop `syscall.SIGHUP` from `subprocess.CancelSignals` | the SIGHUP leg of the same test | narrow the exported set, run the same test, expect the SIGHUP subtest red on the surviving builder pid and the SIGTERM subtest still green |
| CS3 | drop `os.Interrupt` from `subprocess.CancelSignals` | the SIGINT leg of the same test | narrow the exported set, run the same test, expect the SIGINT subtest red — the leg that proves widening the set did not silently drop the signal every site already had |
| CS4 | set `complete = true` before `build` returns in `runbinary.Factory.Own` | the selection-directory leg of the same test | promote the flag early, run the same test, expect the leg red on the surviving `bench-run-*` directory while the three drain legs stay green |
| CS5 | append a fourth signal to `CancelSignals` | `TestCancelSignalsMembershipIsExact` | extend the exported set, run `go test ./internal/subprocess -run TestCancelSignalsMembershipIsExact`, expect the exact-membership failure naming the extra signal |

## Out of scope

Two halves stay open and are named here rather than silently dropped.

SIGKILL, and the `Pdeathsig` design change that would cover it: after this ticket a
SIGKILLed owner still leaks its builder group.

The drift guard. Nothing in the gate stops a fourteenth site from being written with a
signal literal, which is exactly how the three spellings drifted. A conformance check —
production Go only, `internal/subprocess` excluded, `_test.go` excluded because
`internal/systemtest/owner_test.go:250` legitimately registers its own `os.Interrupt` as a
fixture — is the guard, and it is a second ticket rather than this one, because it turns a
single green step into a migration plus a new gate rule. Captured for the roadmap.

Adjacent roadmap rows FT197 (the Go core owning gate invocation and process lifetime) and
FT178 (the worktree bare verb's signal leak) touch neighbouring surfaces; neither owns the
signal-set fact, so this ticket is not folded into either.
