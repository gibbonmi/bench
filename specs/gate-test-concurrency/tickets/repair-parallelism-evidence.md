# Repair the parallelism evidence and overlap test

Blocked by: repair-public-entry-kit-witness.md, repair-ambient-kit-migrations.md, repair-dead-kit-wrappers.md, repair-unused-prospective-kit-seam.md
Ownership fence: `internal/gate/runner_test.go`, `specs/gate-test-concurrency/tickets/adopt-gate-test-parallelism.md`, `specs/gate-test-concurrency/implementation-evidence.md`
Integration surfaces: `internal/gate/runner_test.go` overlap assertion→`specs/gate-test-concurrency/implementation-evidence.md` + RPE1; `specs/gate-test-concurrency/tickets/adopt-gate-test-parallelism.md` timeout contract→`specs/gate-test-concurrency/implementation-evidence.md` + RPE1; repair-public-entry-kit-witness.md, repair-ambient-kit-migrations.md, repair-dead-kit-wrappers.md, and repair-unused-prospective-kit-seam.md→`specs/gate-test-concurrency/implementation-evidence.md` + RPE1
Contracts: the census in `specs/gate-test-concurrency/implementation-evidence.md` is derived from the exact repaired candidate; serial reasons map one-for-one to every structurally serial top-level test; a baseline command is green before its mutation; all timings name command, width, exit, wall time, and exact subject
Closure: RPE1/overlap-serialization, RPE1/census, RPE1/serial-inventory, RPE1/timing, RPE1/mutation-ledger, RPE1/narrow-timeout

## What to build

Keep `TestRunnerRunsPhasesConcurrently` serial and state its load-bearing
reason: it measures the runner's own phase overlap and must not depend on the
Go test scheduler providing overlap. Correct TP1's narrow-width mutation
timeout from 120 seconds to 600 seconds so the unmutated width-one baseline is
green before the proposed deadlock mutation.

Create `implementation-evidence.md` for the exact repaired candidate. Derive
and record the top-level `internal/gate` test count, first-statement
`t.Parallel` count, and complete serial inventory with one substantive reason
per test. Record the original three-run before timings, then three new
same-command after timings on this exact candidate and their median. Record
green focused race, `GOMAXPROCS=2`, `-parallel 1`, and `-parallel 2` runs. For
the width runs use the 600-second timeout and collect filesystem output and
maximum RSS with `/usr/bin/time -v` so local-development cost stays visible.

The ledger also preserves the mutation-red and restored-green facts already
observed for IK, RP, and TP, plus the two repair tickets. Each entry names the
mutation, focused command, red symptom/exit, restoration check, and restored
green result. Re-run a mutation when the earlier fact was not attributable to
a green baseline; never describe the 120-second narrow-width timeout as a
valid mutation red.

## Acceptance

- [ ] [RPE1] (covers local) the overlap assertion is serial; the exact repaired census and complete reasoned serial list are durable; three exact-candidate runs retain a median at or below 90 seconds; race, two-core, width-one, width-two, and hostile-ambient runs are green; and every claimed mutation red is attributable to a green baseline and exact restoration, repairing TP1/P3/C1/C2/C3/STD1.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RPE1/overlap-serialization | add `t.Parallel` back to `TestRunnerRunsPhasesConcurrently` | the structural serial-list audit | apply, run the audit, expect the test to be named as wrongly parallel |
| RPE1/census | remove one first-statement marker from an eligible test | the exact census audit | apply, run the audit, expect the named missing marker red |
| RPE1/serial-inventory | omit one structurally serial test from the evidence list | the exact list comparison | apply, run the comparison, expect the omitted test name red |
| RPE1/timing | remove enough eligible markers to return the measured median above 90 seconds | the same-command three-run measurement | apply, run three repetitions, expect the median budget red |
| RPE1/mutation-ledger | omit a required mutation field or record a red whose baseline is not green | the ledger-shape and baseline audit | apply, run the audit, expect the incomplete or unattributable entry to be named |
| RPE1/narrow-timeout | set the baseline width-one timeout below its observed green duration | baseline-before-mutation validation | apply, run width one, expect rejection because the baseline itself is red, not acceptance as mutation evidence |
