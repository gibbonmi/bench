# Test cost signal in bench test

Status: staged

Roadmap: FT409

Decision source: the named reviewed artifact `roadmap/FT409.md` (reviewer approval 2026-10-09), with its late forks closed on 2026-10-21 by the reviewer-delegated consultants (fable medium, opus high, and sonnet high as the tiebreaker)

Verification log: 3 iteration(s) to accept — the spec round took 2 iterations, because three rows had fixtures that a cheap wrong build passed; the ticket round took 1

## Problem

No Bench output shows the cost of one test. A ticket can double the run time of a package, and the first signal is a gate that times out.

In the FT217 build, the adopt package grew from about 180 to 468 seconds inside the gate. The first signal was the Go per-binary timeout of 600 seconds.

`bench test` prints the elapsed time of each package today. It discards the elapsed time of each test, and it gives no warning before a package reaches the limit.

No Bench source states the limit. The gate passes no `-timeout` flag, so the limit is the implicit Go default, and a `GOFLAGS` value can change it silently.

FT369 wants an author charge that names a package time budget. That contract needs a cost signal that a charge can cite.

## Solution

One constant in `internal/bounds` owns the per-binary go test timeout. Every Bench-owned go test command passes that constant as an explicit `-timeout` flag. So the value that `bench test` warns against is the value that stops the test binary.

`bench test` prints a new `slowest` table. The table names the three slowest top-level tests of each package, with the elapsed time of each test.

`bench test` prints a `timeout_warnings` table when a package uses half or more of the timeout. The table is first in the output, so it stays in the visible head of a bounded response. A warning never changes the exit code.

## User stories

Line: opus / medium.
Implementation-line reason: TC-C2 is the hardest chunk, because it changes the decoder and the render order and it updates the full-equality goldens. The decisions are exact, the seams are known in-process seams with canned-event precedents, and full-equality tests turn a wrong shape red cheaply.
Harder chunks: TC-C2.

### One owner for the timeout

1. As a kit maintainer, I want the release core test step to pass the per-binary timeout from one constant, so that the step obeys the stated limit.
2. As a gate operator, I want the gate `test` phase to pass the timeout explicitly, so that a `GOFLAGS` value cannot change the limit silently.
3. As a gate operator, I want the gate `system` phase to pass the same timeout, so that every phase obeys one limit.
4. As a gate operator, I want the gate `race` phase to pass the same timeout, so that the slowest phase obeys the same limit.
5. As a release operator, I want the ship conformance step to pass the same timeout, so that a release run obeys the gate limit.
6. As an agent, I want the focused `bench test` run to pass the same timeout, so that a focused run and the gate stop a binary at one limit.
7. As an agent, I want each named check run of `bench test` to pass the same timeout, so that a check run obeys the gate limit.
8. As a kit maintainer, I want the bounds policy to turn red when the gate stops reading the constant, so that a second copy of the value cannot return.
9. As a teammate who reads the profile, I want the profile phase table to show the timeout flag, so that the document states the real argv.

### The slowest tests of each package

10. As an agent, I want `bench test` to show the elapsed time of each slow test, so that I can see which test costs the most.
11. As an agent, I want at most three tests for each package, so that the table stays short on a large run.
12. As an agent, I want the slowest test first in each package, so that I read the most costly test first.
13. As an agent, I want tests with equal elapsed times in name order, so that two runs give the same row order.
14. As an agent, I want the rows in package order, so that I can find one package quickly.
15. As an agent, I want subtests left out, so that the time of a subtest does not count twice.
16. As an agent, I want skipped tests left out, so that a skip does not look like a cost.
17. As an agent, I want failed tests included, so that a slow failing test stays visible.
18. As an agent, I want a test that did not finish left out, so that the table shows only measured times.
19. As an agent, I want a test event without an elapsed value to show 0, so that the row stays honest about the missing value.
20. As an agent, I want the package event kept out of the slowest table, so that a package total does not look like one test.
21. As an agent, I want the slowest table present when no test ran, so that an empty table is a definitive answer.
22. As an agent, I want the slowest table after the skips table, so that the packages and failures tables stay in the visible head.
23. As an agent, I want the slowest table before the root conformance table, so that the conformance route stays last.
24. As an author of an FT369 charge, I want the package elapsed time unchanged, so that a charge can compare base and tip with the same cell.

### The timeout warning

25. As an agent, I want a warning when a package uses half of the timeout, so that I see the risk long before the limit.
26. As an agent, I want no warning below half of the timeout, so that a healthy package stays quiet.
27. As an agent, I want the warning rows in package order, so that two runs give the same output.
28. As an agent, I want the warning first in the output, so that a bounded response always shows it.
29. As an agent, I want the warning after the check row of a named check, so that the check row still leads the result.
30. As an agent, I want a passing run with a warning to exit 0, so that the warning does not fail my run.
31. As an agent, I want a failing run with a warning to exit 1, so that the warning does not hide a failure.
32. As an agent, I want no warning table on a healthy run, so that the existing output stays the same.
33. As a teammate, I want the changelog to state the new flag and the two new tables, so that I know the output changed.
34. As a kit maintainer, I want the tests to derive every timeout number from the constant, except the one policy pin, so that a later change of the limit moves the tests too.

### Reviewed exclusions

35. As a gate operator, I want the gate phase output unchanged, so that the gate grades as before.
36. As an agent, I want the total wall time of a run left to FT415, so that two specs do not own one fact.
37. As an agent, I want a package with no elapsed value at its timeout left to the failures table, so that no parser of panic text is necessary.

### The policy value pin

38. As a kit maintainer, I want `TestProductionPolicyValues` in `internal/bounds` to pin the timeout as ten minutes, so that an unnoticed change of the limit turns the gate red.

## Implementation decisions

### The timeout owner

`internal/bounds` declares `GoTestTimeout` beside `GateTimeout`, with the value of ten minutes. The value equals the current Go default, so the enforced limit does not change on a host without a `GOFLAGS` timeout. `GateTimeout` stays the deadline of the whole gate, which is a different fact.

`gate.BaseTestArgv` appends one token, `-timeout=` plus `bounds.GoTestTimeout.String()`, right after `-count=1`. The token comes before the caller operands. A command-line flag takes precedence over a `GOFLAGS` value, so the limit stops depending on the environment.

`toolchainPhases` also runs in a linked repository. So a linked repository that sets a `-timeout` in `GOFLAGS` now gets the ten-minute limit in its gate `test` phase. The changelog entry names this audience.

`bench consumers gate.BaseTestArgv` names six callers, and each caller inherits the flag:

- `toolchainPhases` for the gate `test` phase
- `BenchkitPhases` for the gate `system` phase
- `raceDriverArgv` for the gate `race` phase
- `coreTestArgv` for the release core step
- `preprelease.Steps` for the ship conformance step
- `focusedTestArgv` for every `bench test` run, named checks included

The bounds policy owners map in `internal/conformance/bounds_policy_test.go` gains the entry `internal/gate/phases.go` with the token `bounds.GoTestTimeout`. It also gains the entry `internal/testreport/test_cost.go`, which reads the constant. A dropped consumer then turns the `bounds-policy` check red.

### The per-test elapsed time

`decode` keeps the `Elapsed` value of each test terminal event. A test terminal event is an event with a test name and the action `pass`, `fail`, or `skip`. The value goes into a new field of `testResult`. The last terminal event of a test wins, as the `failed` and `skipped` fields do today.

One helper converts seconds to whole milliseconds with `math.Round`. The package cell and the test cell both use that helper, so the conversion has one source.

### The slowest table

The new table is `slowest[N]{package,test,elapsed_ms}`. It holds at most three rows for each package. The row limit is one named constant in `internal/testreport`.

A test is a candidate when both conditions are true:

- its name contains no `/`
- it has a terminal event, and its last terminal action is `pass` or `fail`

A test with only `run` and `output` events has a `testResult` entry but no terminal event, so it is not a candidate.

Rows sort by package name, then by elapsed time from high to low, then by test name. The table renders on every run that renders the packages table, and an empty table prints `slowest[0]{package,test,elapsed_ms}:`.

The render order becomes `packages`, `failures`, `skips`, `slowest`, and then the conditional `root_conformance`.

### The timeout warning table

The new table is `timeout_warnings[N]{package,elapsed_ms,timeout_ms}`. It renders only when it has a row. It is the first block of the `render` result. Under `--check`, the check row comes before the `render` result, so the warning follows the check row and comes before `packages`.

A package gets a row when its `elapsed_ms` is at least `bounds.GoTestTimeout.Milliseconds()` divided by 2. The divisor is one named constant in `internal/testreport`. The threshold derives from the timeout and is never a second duration. The `timeout_ms` cell is `bounds.GoTestTimeout.Milliseconds()`. Rows sort by package name.

A warning changes no `Outcome` field and no exit code. `bench probe` and the landing fold read the outcome kind and the exit code, so their decisions do not change.

### Structure limits

`internal/testreport/testreport.go` has 329 lines, and the limit is 400. New table helpers go into a new file in `internal/testreport`. `check_test.go` has 481 lines and `selection_test.go` has 602 lines, so no ticket adds lines to them. New tests go into new test files: TC-C2 writes `test_cost_test.go`, and TC-C3 writes `timeout_warning_test.go`. So no new file passes 400 lines.

`check_test.go` does not import `internal/bounds`, and an added import line would turn `bench structure --growth` red. So `testreport_test.go`, which has 382 lines, declares one test helper that returns the timeout token from `bounds.GoTestTimeout`. The argv lines of `check_test.go` call that helper in place, so the line count of `check_test.go` does not change.

## Implementation chunks

Each chunk has one ticket. Each ticket fits one fresh author session of about 150k tokens of context. TC-C1 and TC-C2 write disjoint paths. TC-C3 waits for both, because it reads the constant of TC-C1 and the table file of TC-C2.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| TC-C1 / 1-timeout-owner.md | Every Bench-owned go test command passes the one timeout constant | TC1, TC2, TC3, TC4, TC5, TC6, TC7, TC8, TC9, TC10, TC11, TC38 | TestProductionPolicyValues, TestCoreTestArgvCarriesTrimPathAndCountOne, TestBenchkitPhasesKeepWholeProjectTestArgv, TestBenchkitPhasesUseBranchNativeDrivers, TestBranchNativeArchitectureCensus, TestShipConformanceStepArgvPrefix, TestFocusedTestArgvCarriesTrimPathAndCountOne, TestNamedCheckOwnsConformanceEnvironment, TestSystemCheckOwnsTheGateEnvironment, TestRootConformanceCheckRunsUnscoped, bounds-policy | no |
| TC-C2 / 2-slowest-tests.md | `bench test` prints the three slowest top-level tests of each package | TC12, TC13, TC14, TC15, TC16, TC17, TC18, TC19, TC20, TC21, TC22, TC23, TC24, TC25, TC26 | TestSlowestCarriesTestElapsed, TestSlowestKeepsThreePerPackage, TestSlowestOrdersByElapsed, TestSlowestBreaksTiesByName, TestSlowestOrdersByPackage, TestSlowestExcludesSubtests, TestSlowestExcludesSkippedTests, TestSlowestIncludesFailedTests, TestSlowestExcludesUnfinishedTests, TestSlowestMissingElapsedIsZero, TestSlowestExcludesPackageEvent, TestCommandKeepsItsBaseOutput, TestSlowestPrecedesRootConformance, TestPackageElapsedIgnoresTestElapsed | yes |
| TC-C3 / 3-timeout-warnings.md | `bench test` warns first when a package uses half of the timeout | TC27, TC28, TC29, TC30, TC31, TC32, TC33, TC34, TC35, TC36, TC37 | TestTimeoutWarningAtThreshold, TestTimeoutWarningBelowThreshold, TestTimeoutWarningsOrderByPackage, TestTimeoutWarningsRenderFirst, TestTimeoutWarningsFollowCheckRow, TestTimeoutWarningKeepsPassingExit, TestTimeoutWarningKeepsFailingExit, TestCommandKeepsItsBaseOutput, bounds-policy | no |

### Completion plan

The plan uses version 2 for the delegated serial author checkpoints. The build records one fresh implementation author for each ticket on the accepted line. Review closes each frozen chunk before the next author starts.

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"ft409-test-cost-signal-full-20261021","orchestrator_session":"claude:tcs-orchestrator-20261021","author_limit":1,"assignments":{"1-timeout-owner.md":[{"session":"claude:tcs-t1-author-20261021","assignment":"585a61f066a815ffb7118b8bc83840be","model":"opus","effort":"medium","source":"3cf8494a079e87a65401e669736e85175eba3131","native_ref":"claude-agent:bench-writer-medium:tcs-t1-author-20261021"}],"2-slowest-tests.md":[],"3-timeout-warnings.md":[]}},"chunks":[{"id":"TC-C1","tickets":["1-timeout-owner.md"],"verification":[{"id":"t1-bounds","command":"bench test --package ./internal/bounds","ticket":"1-timeout-owner.md"},{"id":"t1-gate","command":"bench test --package ./internal/gate","ticket":"1-timeout-owner.md"},{"id":"t1-preprelease","command":"bench test --package ./internal/preprelease","ticket":"1-timeout-owner.md"},{"id":"t1-testreport","command":"bench test --package ./internal/testreport","ticket":"1-timeout-owner.md"},{"id":"t1-conformance","command":"bench test --package ./internal/conformance","ticket":"1-timeout-owner.md"},{"id":"t1-bounds-policy","command":"bench test --check bounds-policy","ticket":"1-timeout-owner.md"},{"id":"t1-root-conformance","command":"bench test --check root-conformance","ticket":"1-timeout-owner.md"},{"id":"t1-policy-pin","command":"bench test --package ./internal/bounds --run '^TestProductionPolicyValues$'","probe":"Swap GoTestTimeout = 10 * time.Minute for GoTestTimeout = 5 * time.Minute in internal/bounds/bounds.go; TestProductionPolicyValues must red.","ticket":"1-timeout-owner.md"}]},{"id":"TC-C2","tickets":["2-slowest-tests.md"],"verification":[{"id":"t2-testreport","command":"bench test --package ./internal/testreport","ticket":"2-slowest-tests.md"},{"id":"t2-root-conformance","command":"bench test --check root-conformance","ticket":"2-slowest-tests.md"}]},{"id":"TC-C3","tickets":["3-timeout-warnings.md"],"verification":[{"id":"t3-testreport","command":"bench test --package ./internal/testreport","ticket":"3-timeout-warnings.md"},{"id":"t3-conformance","command":"bench test --package ./internal/conformance","ticket":"3-timeout-warnings.md"},{"id":"t3-bounds-policy","command":"bench test --check bounds-policy","ticket":"3-timeout-warnings.md"},{"id":"t3-root-conformance","command":"bench test --check root-conformance","ticket":"3-timeout-warnings.md"}]}],"final_verification":[{"id":"testreport","command":"bench test --package ./internal/testreport"},{"id":"gate-package","command":"bench test --package ./internal/gate"},{"id":"whole-project","command":"bench gate"},{"id":"coverage","command":"bench coverage --check specs/test-cost-signal/spec.md"}]}
```

## Testing decisions

A good test drives `bench test` behavior through canned `go test -json` events and reads the rendered text. No test compares an elapsed value from a real run. Every test builds each timeout number from `bounds.GoTestTimeout`, with one exception.

The exception is `TestProductionPolicyValues` in `internal/bounds`, which pins each sibling policy value as a literal. That test also pins `GoTestTimeout` as ten minutes. Every other test derives from the constant, so only this literal turns red when a change of the policy value goes unnoticed. The ticket records the mutation `GoTestTimeout = 5 * time.Minute` and the observed red of that test, as the independent-expectation exception of `AGENTS.md` requires.

The argv rows use the existing exact-argv tests. Each test builds the flag token from `bounds.GoTestTimeout.String()` by its own expression, not by a call into the code under test.

The table rows use two precedents. The decode precedent is `TestPackageTableCarriesElapsedMilliseconds`, which calls `decode` over inline events and then `render`. The command precedent is `commandOverEvents`, which runs `Command` over a stub `go` that replays the events.

A threshold fixture sets the package `Elapsed` from the constant. The fixture at the threshold uses `GoTestTimeout.Milliseconds() / 2` milliseconds. The fixture below it uses one millisecond less. Both sides of the `>=` predicate then have a row.

The gate observes the feature through the ordinary `test` phase, which runs the `internal/testreport`, `internal/gate`, and `internal/conformance` tests. Root conformance runs the `bounds-policy` check inside that phase.

### Seam diagram

    trigger: an agent runs `bench test`, or the gate runs a phase
        │
        ▼
    bounds.GoTestTimeout ──▶ [ gate.BaseTestArgv ] ──▶ argv with -timeout
                                ◀ tests attach here: exact-argv tests compare each caller's argv

    canned go test -json events
        │
        ▼
    [ decode ] ──▶ [ render ] ──▶ timeout_warnings, packages, failures, skips, slowest, root_conformance
                       ◀ tests attach here: decode-level and commandOverEvents tests read the text and the exit

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| TC1 | 1 | The release core argv puts the timeout token right after `-count=1` | `internal/gate/branch_native_phases_test.go` (`TestCoreTestArgvCarriesTrimPathAndCountOne`) | Exact equality reds an argv without the token |
| TC2 | 2 | The gate `test` phase argv equals the base flags, the timeout token, and `./...` | `internal/gate/lane_test.go` (`TestBenchkitPhasesKeepWholeProjectTestArgv`) | Exact equality reds a missing or moved token |
| TC3 | 3 | The gate `system` phase argv carries the timeout token before `-tags=system` | `internal/gate/branch_native_phases_test.go` (`TestBenchkitPhasesUseBranchNativeDrivers`) | Exact equality reds a system argv that skips the owner |
| TC4 | 4 | The race argv starts with the base flags, the timeout token, `-race`, and `-v` | `internal/conformance/ordinary_build_census_test.go` (`TestBranchNativeArchitectureCensus`) | The census prefix compare reds a race driver without the token |
| TC5 | 5 | The ship conformance argv starts with `go -C <kit> test`, the base flags, and the timeout token | `internal/preprelease/preprelease_test.go` (`TestShipConformanceStepArgvPrefix`) | The prefix compare reds a ship step without the token |
| TC6 | 6 | The focused argv equals the base flags, the timeout token, `-json`, and the package | `internal/testreport/testreport_test.go` (`TestFocusedTestArgvCarriesTrimPathAndCountOne`) | Exact equality reds a focused run that the gate limit does not cover |
| TC7 | 7 | A named conformance check child receives the timeout token before `-json` | `internal/testreport/check_test.go` (`TestNamedCheckOwnsConformanceEnvironment`) | The recorded argv line reds a check run without the token |
| TC8 | 7 | The system check child receives the timeout token before `-json` | `internal/testreport/check_test.go` (`TestSystemCheckOwnsTheGateEnvironment`) | The recorded argv line reds a system check without the token |
| TC9 | 7 | The root conformance check child receives the timeout token before `-json` | `internal/testreport/root_conformance_test.go` (`TestRootConformanceCheckRunsUnscoped`) | The recorded argv line reds a root run without the token |
| TC10 | 8 | The bounds policy reds when `internal/gate/phases.go` does not contain `bounds.GoTestTimeout` | the `bounds-policy` check through `bench test --check bounds-policy` | The owners entry turns a dropped consumer red |
| TC11 | 9 | The profile phase table shows the timeout token in the `test`, `race`, and `system` rows | review-owned | No check grades the phase table cells |
| TC12 | 10 | A top-level test with terminal `Elapsed` 1.5 renders the row `canned,TestA,1500` in `slowest` | planned TestSlowestCarriesTestElapsed in internal/testreport/test_cost_test.go | The current decoder drops the test value |
| TC13 | 11 | Four top-level tests in one package render three rows and omit the fastest test | planned TestSlowestKeepsThreePerPackage in internal/testreport/test_cost_test.go | A missing limit prints the fourth row |
| TC14 | 12 | Rows of one package order from high to low elapsed when the events arrive fastest first | planned TestSlowestOrdersByElapsed in internal/testreport/test_cost_test.go | Event order or an ascending sort gives the wrong first row |
| TC15 | 13 | Two tests with equal elapsed times, whose events arrive in reverse name order, render in test name order | planned TestSlowestBreaksTiesByName in internal/testreport/test_cost_test.go | Event order or map iteration order puts the later name first |
| TC16 | 11, 14 | Two packages with four top-level tests each render six rows, three for each package, in package name order, when the later-named package holds the highest elapsed time | planned TestSlowestOrdersByPackage in internal/testreport/test_cost_test.go | A global limit of three rows or a sort by elapsed alone gives the wrong rows |
| TC17 | 15 | A subtest `TestA/sub` with the largest elapsed time renders no row | planned TestSlowestExcludesSubtests in internal/testreport/test_cost_test.go | A filter that keeps subtests counts the time twice |
| TC18 | 16 | A test whose last terminal action is `skip` renders no row | planned TestSlowestExcludesSkippedTests in internal/testreport/test_cost_test.go | A filter on terminal events alone keeps the skip |
| TC19 | 17 | A test whose last terminal action is `fail` renders a row | planned TestSlowestIncludesFailedTests in internal/testreport/test_cost_test.go | A filter on `pass` alone drops the failing test |
| TC20 | 18 | A test with a `run` event and an `output` event and no terminal event renders no row | planned TestSlowestExcludesUnfinishedTests in internal/testreport/test_cost_test.go | A candidate set built from the test results without a terminal marker keeps the test at 0 |
| TC21 | 19 | A test terminal event without `Elapsed` renders the value 0 | planned TestSlowestMissingElapsedIsZero in internal/testreport/test_cost_test.go | A filter on positive values drops the test |
| TC22 | 20 | The package terminal event renders no `slowest` row with an empty test name | planned TestSlowestExcludesPackageEvent in internal/testreport/test_cost_test.go | A candidate set that reads every terminal event adds the package total |
| TC23 | 21 | A run with no test events prints `slowest[0]{package,test,elapsed_ms}:` | `internal/testreport/outcome_test.go` (`TestCommandKeepsItsBaseOutput`) | Full equality reds a table that is absent when empty |
| TC24 | 22 | The `slowest` block comes right after the `skips` block | `internal/testreport/outcome_test.go` (`TestCommandKeepsItsBaseOutput`) | Full equality reds any other position |
| TC25 | 23 | With a skipped root conformance test, the `slowest` block comes before `root_conformance` | planned TestSlowestPrecedesRootConformance in internal/testreport/test_cost_test.go | A block appended last moves the route out of the tail |
| TC26 | 24 | A top-level test terminal with `Elapsed` 2.0, then a package terminal with `Elapsed` 0.5, render `canned,pass,500,1` in `packages` and `canned,TestA,2000` in `slowest` | planned TestPackageElapsedIgnoresTestElapsed in internal/testreport/test_cost_test.go | A packages cell read from the test value prints 2000 |
| TC27 | 25 | A package at the threshold renders `timeout_warnings[1]{package,elapsed_ms,timeout_ms}` with the row `canned,<GoTestTimeout.Milliseconds()/2>,<GoTestTimeout.Milliseconds()>` | planned TestTimeoutWarningAtThreshold in internal/testreport/timeout_warning_test.go | A strict `>` or a missing table drops the row |
| TC28 | 26 | A package one millisecond below the threshold renders no `timeout_warnings` block | planned TestTimeoutWarningBelowThreshold in internal/testreport/timeout_warning_test.go | A lower threshold or an always-present table prints the block |
| TC29 | 27 | Two warned packages, whose events arrive in reverse name order, render in package name order | planned TestTimeoutWarningsOrderByPackage in internal/testreport/timeout_warning_test.go | Event order or map iteration order puts the later name first |
| TC30 | 28 | A warned plain run starts with the `timeout_warnings` header | planned TestTimeoutWarningsRenderFirst in internal/testreport/timeout_warning_test.go | A block after `packages` falls out of the visible head |
| TC31 | 29 | A warned named check run renders the check row, then `timeout_warnings`, then `packages` | planned TestTimeoutWarningsFollowCheckRow in internal/testreport/timeout_warning_test.go | A warning ahead of the check row breaks the check contract |
| TC32 | 30 | A passing warned run exits 0 | planned TestTimeoutWarningKeepsPassingExit in internal/testreport/timeout_warning_test.go | A warning that sets a failure exits 1 |
| TC33 | 31 | A failing warned run exits 1 and renders its failure row | planned TestTimeoutWarningKeepsFailingExit in internal/testreport/timeout_warning_test.go | A warning that replaces the outcome hides the failure |
| TC34 | 32 | A healthy run renders no `timeout_warnings` block | `internal/testreport/outcome_test.go` (`TestCommandKeepsItsBaseOutput`) | Full equality reds an empty warning table |
| TC35 | 33 | The changelog Unreleased section states the timeout flag and its effect on a linked repository, the `slowest` table, and the `timeout_warnings` table | review-owned | No check grades changelog content |
| TC36 | 34 | No test of this spec, except the TC38 pin, spells the timeout, its half, or its milliseconds as a literal | review-owned | A literal survives a later change of the constant |
| TC37 | 8 | The bounds policy reds when `internal/testreport/test_cost.go` does not contain `bounds.GoTestTimeout` | the `bounds-policy` check through `bench test --check bounds-policy` | The owners entry turns a dropped warning consumer red |
| TC38 | 38 | `TestProductionPolicyValues` reds when `GoTestTimeout` is not ten minutes, and the ticket records the red of the mutation `GoTestTimeout = 5 * time.Minute` | `internal/bounds/bounds_test.go` (`TestProductionPolicyValues`) | Every other test derives from the constant, so only this literal reds a changed value |

Not covered: story 35 — the gate phase argv rows TC2 to TC4 already pin the only gate change, and the gate report code stays untouched.
Not covered: story 36 — FT415 owns the total wall time, and this spec adds no total.
Not covered: story 37 — the existing failures table shows the timeout panic line.

### Edge inventory

The canonical edge classes walk the decode and render seams. The hostile-input classes of the profile apply as follows:

- Control bytes in a test name: a top-level test name is a Go identifier, so it holds no control byte. `toon.Table` still refuses one.
- A cell that looks numeric: test names start with `Test`, and the fixtures use the package name `canned`. So no fixture expectation needs the encoder to quote a cell.
- An absent value against a zero value: TC21 asserts the absent `Elapsed` case. The empty table of TC23 asserts the run with no test events.
- A command whose own write changes a reported fact: `bench test` writes no artifact that it also reports.
- Unquoted arguments and spaces: the flag is one argv token, so no shell splits it.

The threshold edge has both sides: TC27 asserts the exact value, and TC28 asserts one millisecond less. The tie edge has TC15. The limit edge has TC13.

**Won't handle:** a timed-out package whose terminal event has no `Elapsed` value — the failures table still prints the Go timeout panic line for that package.
**Won't handle:** a package with no terminal event, such as an interrupted run — the interrupted outcome still reports the incomplete packages.
**Won't handle:** per-test or per-package timing in the gate phase output — the gate `test` phase runs without `-json`, and its exact argv rows still pin the shared flag.
**Won't handle:** the total wall time of a run — FT415 owns that fact, and the packages table still prints each package time.
**Won't handle:** the race phase of the release evidence registry — the gate does not run it, and its argv stays part of the release evidence inputs.
**Won't handle:** a configurable timeout for each package — the one constant stays the only limit, and every caller still passes it.

## Out of scope

- Per-package timing and warnings in the gate phase output: 6 edits, 2 gate runs. The gate needs `-json` and a decoder first.
- The total wall time of a red run: FT415 owns it. 4 edits, 1 gate run.
- A configurable timeout through an environment value or a profile field: 5 edits, 1 gate run.
- The author charge that compares base and tip cost: FT369 owns it. 3 edits, 1 gate run.

## Further notes

### Derived Writes scope

The slicer verifies this scope against the reader sweep. The scope below is the prediction that the approval table carries.

- TC-C1: `internal/bounds/bounds.go`, `internal/bounds/bounds_test.go`, `internal/gate/phases.go`, `internal/gate/gate_go.go`, `internal/gate/lane_test.go`, `internal/gate/branch_native_phases_test.go`, `internal/conformance/ordinary_build_census_test.go`, `internal/conformance/bounds_policy_test.go`, `internal/preprelease/preprelease_test.go`, `internal/testreport/testreport_test.go`, `internal/testreport/check_test.go`, `internal/testreport/root_conformance_test.go`, `projects/benchkit.md`
- TC-C2: `internal/testreport/testreport.go`, `internal/testreport/test_cost.go`, `internal/testreport/test_cost_test.go`, `internal/testreport/outcome_test.go`, `internal/testreport/tests_run_test.go`, `internal/testreport/failure_lines_test.go`
- TC-C3: `internal/testreport/testreport.go`, `internal/testreport/test_cost.go`, `internal/testreport/test_cost_test.go`, `internal/testreport/timeout_warning_test.go`, `internal/testreport/check_row_test.go`, `internal/testreport/outcome_test.go`, `internal/testreport/canned_go_test.go`, `internal/testreport/tests_run_test.go`, `internal/conformance/bounds_policy_test.go`, `CHANGELOG.md`

The slicing research added `internal/gate/gate_go.go` to TC-C1 and four `internal/testreport` test files to TC-C3. Each addition came from one of two readers, and the keep rule holds each path.

TC-C1 and TC-C3 both write `internal/conformance/bounds_policy_test.go`. TC-C1 adds the gate entry, and TC-C3 adds the testreport entry. The serial order keeps the two edits apart.

### Pre-review proof checklist

- Cited symbols: `BaseTestArgv`, `coreTestArgv`, `raceDriverArgv`, `focusedTestArgv`, `decode`, `render`, `testResult`, `checkRow`, `checkBoundsPolicy`, and `GateTimeout` resolve in the files named above
- Import edges: `internal/gate` already imports `internal/bounds`, and `internal/testreport` already imports `internal/bounds` and `internal/gate`, as `go list` shows
- Source-row clauses and occurrences: the source-to-row table below
- Promised field labels: `slowest[N]{package,test,elapsed_ms}` and `timeout_warnings[N]{package,elapsed_ms,timeout_ms}`
- Changed-function callers: `bench consumers gate.BaseTestArgv testreport.decode` lists six production callers of `BaseTestArgv`, one production caller of `decode`, and four test callers of `decode`
- Copy survival: TC10 and TC37 red when a consumer stops reading the constant, TC1 to TC9 red a caller without the token, and TC38 reds a changed value
- Rendered-shape readers: the needle `skips[0]{package,test,reason}:\n"` hits `outcome_test.go` lines 258 to 282, `tests_run_test.go` line 31, and `failure_lines_test.go` lines 24, 36, and 47, all in TC-C2
- Pin operators: TC1 to TC3 and TC6 use `reflect.DeepEqual`, TC4 and TC5 use a prefix compare, TC7 to TC9 use `strings.Contains`, and TC23, TC24, and TC34 use `!=` on the whole output
- Entry reads: none
- Derived expectations: each argv test builds the token from `bounds.GoTestTimeout`, and each threshold fixture builds its value from the same constant
- Consolidated rules: none
- Quantified obligations: "every Bench-owned go test command" covers the six callers that `bench consumers` lists, and TC1 to TC9 cover each one
- Workflow-step writes: none

The reader sweep found one parser of the output, `bench probe`, which reads the outcome kind and passes the bytes through. The landing fold prints the `--changed` output and reads the exit code. No workflow file, hook, or shell script reads `bench test` output. The `bench test` verb is outside the AXI contract fragments, so no fragment grades its tables.

The `bounds-policy` check reds a production value spec outside `internal/bounds` whose value text equals a registry value. So the constant must live in `internal/bounds`, and a copy of `10 * time.Minute` elsewhere stays red-capable.

### Source-to-row table

| reviewed source clause | occurrence | rows or disposition |
|---|---|---|
| `bench test` reports the slowest tests of each package | `roadmap/FT409.md` body, criterion FT409.slowest | TC12 to TC25 |
| and the package elapsed time | `roadmap/FT409.md` body, criterion FT409.slowest | TC26; the packages table already prints it |
| It warns when a package nears the go test timeout | `roadmap/FT409.md` body, criterion FT409.timeout | TC27 to TC34 |
| that the gate uses | `roadmap/FT409.md` body | TC1 to TC10 make the warned value the enforced value |
| This signal makes the test-cost contract in FT369 checkable | `roadmap/FT409.md` body | TC26 keeps the cell a charge compares; FT369 is out of scope |

### Closed decisions

The consultants closed these forks on 2026-10-21, under the reviewer delegation:

- The timeout is one constant in `internal/bounds`, which `BaseTestArgv` passes to every caller.
- The warning threshold is half of the timeout, derived through one divisor constant.
- The slowest table holds three top-level tests for each package and is always present after `skips`.
- The warning table is conditional, is named `timeout_warnings`, and renders first.
- A warning never changes the exit code.
- The signal appears in `bench test` only.
- A timed-out package with no elapsed value gets no warning.
- No test spells a timeout number as a literal, except the one pin in `TestProductionPolicyValues`.
- The slowest table admits only a test whose last terminal action is `pass` or `fail`.

### Flagged additions and residuals

The `bounds-policy` owners entries of TC10 and TC37 are an addition beyond the source. They make the single owner of the timeout red-capable under the code standard.

The profile phase table edit of TC11 is a currency repair. The table states the gate argv, and the new token changes that argv.

The `TestProductionPolicyValues` pin of TC38 obeys the tree convention for policy values. Both consultants chose it on 2026-10-21. It is the one literal timeout in the tests of this spec.

This planning pass runs the prose and coverage checks only. It supplies no implementation, runtime-test, mutation, or gate evidence.
