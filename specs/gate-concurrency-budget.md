# gate-concurrency-budget

Status: implemented

## Problem

FT91's first arm, compiled from the closed `decisions/gate-concurrency.md` map.

The kit's gate takes 10–15 minutes on a 16-core box, and most of that wall clock is
contention rather than work. Measured 2026-07-22: load average ~123 during a gate run.

The mechanism is nesting nobody budgets. The canary phase fans 144 fixtures over
`runtime.NumCPU()` workers, and each worker spawns a full inner gate whose own `go
build` and `go test` default to machine width again. Sixteen workers times sixteen
inner threads is roughly a sixteen-fold oversubscription of the box, and it runs
concurrently with the outer conformance and contract phases, which are also at machine
width.

The cost is not only time. The marker stalls and cleanup flakes the retired
gate-trustworthiness arm fixed were contention symptoms, so an unbudgeted gate is also
a less trustworthy one. Nothing in the tree states an aggregate concurrency budget:
`internal/canary` decides its own worker count, and each inner gate decides its own
width, and neither knows about the other.

## Solution

Budget the product, not either factor alone.

**One lever per inner gate.** Each inner gate is invoked with an explicit
`GOMAXPROCS=k` in its environment. That single entry caps both `go build` parallelism
and the test binary it produces, so an inner gate can no longer expand to machine
width on its own.

**Workers derived from the budget.** The canary worker count becomes
`runtime.GOMAXPROCS(0) / k`, floored at one and capped at the fixture count, so workers
times inner width lands at roughly the cores actually available.

**k is 2, and it is measured.** The map's ticket #2 ran the full gate at k ∈ {1, 2, 4}
on this repo with a warm compile cache on an idle box: k=1 took 330 s, k=2 took 332 s,
k=4 took 432 s. k=1 and k=2 tie on wall clock, and k=2 halves the number of concurrent
inner gates — halving memory and temp-directory pressure at no measured cost. k=4 pays
30% more wall for under-parallelizing 144 fixtures. After the change the gate's long
pole is the conformance phase at 319 s, not canary.

**No Bench-specific knob.** The budget is `runtime.GOMAXPROCS(0)` of the outer process.
In Go 1.25 that value is already cgroup- and container-aware and already honors an
explicit `GOMAXPROCS` environment variable, so `GOMAXPROCS=8 bench gate` is the escape
hatch and there is no second concurrency vocabulary to learn. `runtime.NumCPU()` stops
being a budget source anywhere in this arm.

**The pin cannot be leaked past.** Because Go's exec environment gives no guaranteed
precedence to a duplicate key, the inner environment strips any inherited `GOMAXPROCS`
before appending its own, so exactly one entry reaches the inner gate.

Nothing else moves. No CLI surface, exit code, or output text changes; green keeps
meaning exactly what it meant before.

## User stories

1. As a reviewer running `bench gate` on this repo, I want the canary phase's nested
   inner gates to fit the machine's core budget in aggregate instead of demanding
   roughly sixteen times it, so that the gate's wall clock measures work rather than
   contention.
   Line: `gpt-5.6-terra` / medium. This story's deliverable is the post-change
   measurement judged against the decision map's k table, which is a reading of
   evidence rather than code the gate can grade.

2. As a reviewer, I want each inner gate invoked with an explicit `GOMAXPROCS` width, so
   that one lever caps both its build parallelism and its test binary instead of each
   inner gate fanning out to machine width on its own.
   Line: `gpt-5.6-luna` / medium. The width is a decided constant and the environment
   entry it produces is directly assertable through the injected runner.

3. As a reviewer, I want the canary worker count derived as the outer budget divided by
   that width, floored at one and capped at the fixture count, so that workers times
   inner width stays at roughly the cores available.
   Line: `gpt-5.6-luna` / medium. The derivation is a total function of two integers and
   this spec enumerates its answers, so the cheap tier is building to a fully pinned
   target.

4. As an operator on a container or a smaller box, I want `GOMAXPROCS=8 bench gate` to
   shrink the whole canary budget with no Bench-specific knob, so that the standard Go
   lever stays the one escape hatch.
   Line: `gpt-5.6-luna` / medium. Reading the budget from `runtime.GOMAXPROCS(0)` rather
   than `runtime.NumCPU()` is a one-call change whose effect is observable in the
   derived worker count.

5. As a reviewer, I want any inherited `GOMAXPROCS` stripped before the inner one is
   set, so that an operator override cannot leak past the cap and hand an inner gate
   machine width.
   Line: `gpt-5.6-luna` / medium. Go's exec environment has no guaranteed precedence for
   a duplicate key, so the assertion is that exactly one entry survives, which is a
   count a test reads straight off the recorded call.

6. As a Bench maintainer, I want the inner width to live in `internal/bounds` as one
   named policy constant that `internal/canary` consumes and the bounds-policy
   conformance check enforces, so that this tunable has a single owner like every other
   resource bound in the kit.
   Line: `gpt-5.6-terra` / medium. The change extends the check that decides whether
   policy values stay single-sourced, and the profile keeps oracle-affecting edits off
   the cheap tier.

7. As an agent whose canary package tests run nested inside a fixture's own inner gate,
   I want every concurrency expectation in those tests keyed to the derived worker bound
   rather than to machine width, so that a nested run at a bound of one asserts
   truthfully or skips by capability instead of deadlocking until the phase timeout
   turns conformance red.
   Line: `gpt-5.6-luna` / high. The seam and the three affected tests are named here, but
   the failure mode is a nested-environment deadlock rather than a wrong value, so the
   effort buys the care that separates an honest capability skip from a weakened
   assertion.

8. As a reviewer, I want the canary sweep's observable contract unchanged — the same CLI
   surface, exit codes, output text, baseline-before-fixtures ordering, sorted error
   ordering, and temp-directory cleanup — so that this arm stays a scheduling change.
   Line: `gpt-5.6-luna` / low. The existing canary tests already state every one of these
   properties, so the story is satisfied by leaving them passing and unedited.

## Implementation decisions

**`internal/canary` owns everything that changes; `internal/gate` is untouched.** The
outer phase table, the phase runner, and the outer phase environment stay exactly as
they are. The map's evidence says so: after the canary cap, load still peaks at roughly
twice cores in bursts, but wall clock is conformance-bound, so capping outer phase width
buys nothing and could cost. That question stays dormant unless contention flakes
persist.

**`bounds.CanaryInnerWidth = 2` is the one source of k.** It joins the const block in
`internal/bounds/bounds.go` beside the other tunables. `internal/canary/canary.go` is
its only consumer. Three registrations in `checkBoundsPolicy` follow the pattern the
existing entries already set: the name joins the `required` list, `internal/canary/canary.go`
joins the `owners` map against `bounds.CanaryInnerWidth`, and `boundLikeName` gains a
words entry so a redeclaration under a bound-like name is caught the way
`ProviderTimeout` already is.

**The derivation is a pure function of two integers.** `runFixtures` calls a small
in-package helper taking the outer budget and the fixture count and returning the worker
count: divide, floor at one, cap at the fixture count. Passing the budget in rather than
reading `runtime.GOMAXPROCS(0)` inside the helper is what makes the policy testable as a
total function instead of as a property of whichever machine runs the test, and it keeps
one source — the tests that need the bound call the same helper the sweep does.

**The inner environment strips, then appends.** `innerEnv` adds `GOMAXPROCS=` to the
prefixes it already filters out, then appends exactly one `GOMAXPROCS=<k>` alongside the
existing `BENCH_CANARY_INNER=1`. Never append without stripping: Go's exec environment
has no guaranteed duplicate-key precedence, so a duplicate is a coin flip rather than an
override. The per-fixture `BENCH_CANARY_PHASE` append rides on top of this environment
and must not displace it.

**k stays constant even when the budget is smaller than k.** At an outer budget of 1 the
derivation yields one worker at inner width 2, a bounded two-fold overshoot on the
smallest possible box. Special-casing it would add a second rule to a policy whose value
is being one rule, and the overshoot at that size is two threads.

**`Sweep(root, Runner)` keeps its signature, and no new abstraction appears.**
`runFixtures` stays the deep unit that hides scheduling; the injected `Runner` stays the
seam. The knobless constant replaces the prototype's environment variable rather than
adding to it.

**Deviation from the map, flagged for veto.** Handoff item 4 says the existing overlap,
baseline-order, error-order, and temp-cleanup tests are unchanged. That holds for
baseline-order and temp-cleanup, but not for the two tests that assume at least two
concurrent workers. The conformance phase runs the kit's own `go test` over every core
package as a subprocess inheriting the phase environment, so inside a fixture's inner
gate the canary package tests run at `GOMAXPROCS=2`, where the derived bound is one.
`TestSweepRunsFixturesConcurrently` and `TestSweepReportsErrorsInSortedFixtureOrder` both
require two workers to make progress, and `TestSweepBoundsFixtureConcurrencyAtNumCPU`
waits for machine-width in-flight before releasing, which is the deadlock the map's
watch-out #9 records observing. All three retarget to the derived bound: the two that
need overlap gate on the bound through the existing `capability.CPU` helper, and the
bounds test releases and asserts at the derived bound. This is story 7, and it is the one
thing in this spec that departs from the Handoff.

## Testing decisions

A good test here drives `Sweep` with a fake `Runner` and reads what the sweep did: how
many calls were in flight at once, and what environment each recorded call carried. No
test reaches into worker goroutines or scheduling state, and no test asserts a wall-clock
number — the wall-clock outcome is measured evidence, not a gate assertion.

Prior art in this repo, followed rather than reinvented:

- `internal/canary/canary_concurrency_test.go` — the existing fake-`Runner` tests for
  overlap, bounds, ordering, and cleanup. Every behavioral row extends this file.
- `internal/bounds/bounds_test.go` — table-driven unit tests over a pure policy helper.
  The worker derivation's table takes the same shape.
- `tests/canary/package-core-guard/bounds-duplicate-owner` — the mutation fixture that
  proves `checkBoundsPolicy` bites. The new ownership fixture is its sibling.

Gate command: `.bench/gate.sh` (the project gate).

The canary package's tests run inside the gate's conformance phase, as part of the
`go test` over core packages that check runs as a subprocess, so the gate sees this seam
directly.

### Seam diagram

**Seam 1 — the sweep, through the injected `Runner`.** The behavioral seam for
concurrency and environment.

    trigger: Sweep(root, runner) — the gate's canary phase, or a test
        │
        ▼
    fixture tree   ──▶  [ Sweep                                 ]  ──▶  RunCall per fixture
    process env    ──▶  [   innerEnv: strip + pin GOMAXPROCS    ]  ──▶    .Env
    GOMAXPROCS(0)  ──▶  [   runFixtures: derived worker pool    ]  ──▶  in-flight call count
                        [   runFixture: per-fixture verdict     ]  ──▶  error strings
                            ◀ tests attach here: pass a fake Runner that records each
                              call's Env and tracks in-flight high-water, then assert
                              both against the derived bound

**Seam 2 — the worker derivation.** Pure; no repo, no process, no goroutines.

    trigger: runFixtures, once per sweep
        │
        ▼
    outer budget   ──▶  [ fixtureWorkers(budget, fixtures) ]  ──▶  worker count
    fixture count  ──▶  [   divide · floor 1 · cap         ]
                            ◀ tests attach here: call it with a table of integer pairs
                              and assert the returned count

**Seam 3 — bounds ownership.** The repository-wide static check.

    trigger: the conformance phase, with BENCH_CONFORMANCE_ROOT set
        │
        ▼
    repo tree  ──▶  [ checkBoundsPolicy: registry names,   ]  ──▶  violations
                    [ owner consumption, redeclaration     ]  ──▶  exit code
                        ◀ tests attach here: a mutation fixture that stops canary.go
                          consuming the constant, asserted red by the canary phase

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 3 | With more fixtures than workers, the sweep's in-flight high-water equals the derived worker bound rather than machine width | seam 1, fake `Runner` | retarget `TestSweepBoundsFixtureConcurrencyAtNumCPU`: release at the derived bound, settle, assert high-water equals it — red today, `workers = runtime.NumCPU()` reaches machine width whenever the bound is lower | Asserting equality rather than a ceiling rejects both directions at once: an uncapped sweep overshoots, and a sweep that serialized everything undershoots and would pass a `<=` assertion |
| 3 | The derived worker count is the budget divided by the width, floored at one | seam 2, unit table | new table over `fixtureWorkers` with budgets 0, 1, 2, 3, 16 against a large fixture count — red against a derivation without the floor, which returns zero workers | A zero-worker pool does not fail loudly, it hangs on an unbuffered job channel, so the value has to be asserted where it is a returned number rather than a deadlock |
| 3 | The derived worker count never exceeds the fixture count | seam 2, unit table | same table with fixture counts 1 and 3 against a 16-core budget — red against a derivation that drops the existing cap while rewriting the line above it | The cap exists in the tree today and this story rewrites the exact statement that carries it, so its survival is asserted rather than assumed |
| 4 | The budget comes from `runtime.GOMAXPROCS(0)`, so a process started under a lowered `GOMAXPROCS` derives fewer workers | seam 2, unit table | the same table is the assertion: `fixtureWorkers` takes the budget as an argument, and `runFixtures` passing `runtime.NumCPU()` instead is a one-line diff visible at the single call site | Making the budget a parameter is what turns an untestable machine property into a graded one; the remaining risk is which value the one call site passes, which is a reading, not a behavior |
| 2, 5 | Every `RunCall.Env` — baseline and fixtures — carries exactly one `GOMAXPROCS` entry, set to the width constant, when the process environment already carries `GOMAXPROCS=32` | seam 1, fake `Runner` | new test setting the outer variable and counting `GOMAXPROCS=` prefixes on each recorded call — red today, no entry is set at all, and red against an append-without-strip implementation, which yields two | Counting rather than searching is the whole assertion: a duplicate key satisfies a contains-check while leaving which value wins to the exec implementation |
| 2, 5 | A fixture call that also carries `BENCH_CANARY_PHASE` still carries the single pinned `GOMAXPROCS` | seam 1, fake `Runner` | same test asserting on a fixture call rather than only the baseline — red against an implementation that builds the pin into the baseline environment only | The per-fixture environment is derived by appending to the shared one, so an implementation that pins in the wrong place passes on the baseline call and fails every fixture |
| 6 | `internal/canary/canary.go` failing to consume `bounds.CanaryInnerWidth` is reported by the conformance check | seam 3, mutation fixture | new canary fixture replacing the constant reference with a bare literal, expecting `internal/canary/canary.go does not consume bounds.CanaryInnerWidth` — red against a check whose `owners` map was never extended | Registering the name in the `required` list alone would leave the ownership unenforced, and the fixture is what distinguishes a registered constant from an enforced one |
| 6 | A production redeclaration of the width under a bound-like name is reported | seam 3, `internal/bounds` policy check | the existing `bounds-duplicate-owner` fixture family's assertion, extended with the new words entry — red against a words map left without the entry | Without the words entry the value `2` is owned in name only, and the redeclaration path that already caught `ProviderTimeout` would pass silently for this constant |
| 7 | The whole canary package passes when run under a budget where the derived bound is one | seam 1, package run | observed red: with story 3 landed and the tests unretargeted, `GOMAXPROCS=2 go test -timeout 120s ./internal/canary` fails the overlap and error-order tests and times out on the bounds test | This is the exact nested condition the gate creates, and it is the failure the map recorded observing, so it is graded by reproducing the environment rather than by reasoning about it |
| 7 | At a derived bound of one, the two tests that require overlap skip through the capability helper rather than failing or hanging | seam 1, package run | the same `GOMAXPROCS=2` run, asserted green with two `bench-skip kind=capability class=cpu` lines on stdout — red against retargeting that deleted the assertions instead of gating them | A deleted assertion and a capability skip both make the run green; only the emitted evidence line distinguishes honest unavailability from a weakened test |
| 7 | At full machine width the same tests still run and still assert overlap | seam 1, package run | `go test ./internal/canary` with no budget override, asserted green with no capability skip line — red against a gate condition that always fires | Gating on the bound is only correct if the gate is false on the machine that can run the test; an always-true condition silently retires the overlap assertion everywhere |
| 8 | Baseline completes before any fixture starts, errors report in sorted fixture order, and temp work directories are removed on both green and red paths | seam 1, fake `Runner` | already covered — `TestSweepCompletesBaselineBeforeStartingFixtures`, `TestSweepReportsErrorsInSortedFixtureOrder`, and the two temp-cleanup tests, which must pass unedited except for story 7's capability gate | These are the properties a scheduling rewrite is most likely to break, and they are already stated as tests, so the signal is that they keep passing without their assertions being touched |
| 8 | `bench canary` keeps its exit codes and output text | `internal/contract` | already covered — the existing canary command contract tests, unedited | This arm changes no argument, no message, and no verdict, so an unedited contract suite passing is the correct evidence rather than new assertions |
| 1 | Gate wall clock and load improve against the 2026-07-22 baseline, and the k=2 row of the map's table reproduces | none — measured evidence | not TDD-able: a full-gate wall-clock and load measurement, compared to the map's ticket #2 table and the 10–15 minute, load ~123 baseline | Wall clock is a property of the machine and its load at the moment of the run, so asserting it in the gate would make the oracle report the box's mood rather than the code's correctness |
| edge: hostile budget | An outer `GOMAXPROCS` that is empty, non-numeric, zero, or negative still yields exactly one pinned inner entry and at least one worker | seam 1 and seam 2 | the environment test parameterized over `GOMAXPROCS=`, `GOMAXPROCS=abc`, and `GOMAXPROCS=0`, plus the derivation table's zero and one budgets — red against a strip keyed on a parsed value rather than the key prefix | Go's runtime ignores an unparseable value, so the hostile input never reaches the derivation; the risk is entirely in the strip, which must key on the name and not on whether the value made sense |
| edge: single fixture | A sweep of exactly one fixture runs it and returns its verdict unchanged | seam 1 and seam 2 | the derivation table's fixture-count-1 rows, plus the existing single-fixture cleanup test — red against a derivation that returns zero workers for one fixture | One fixture is the boundary where the cap and the floor disagree, and it is the shape every temp-cleanup test already uses, so a regression there is broad |

### Edge inventory

Walked against the profile's hostile-input checklist for the surfaces this spec touches.
This arm parses no new input and adds no new CLI surface, so most classes resolve as
won't-handle with a reason rather than as coverage.

- Absent vs present-but-empty value — the hostile-budget row: `GOMAXPROCS` unset,
  `GOMAXPROCS=`, and `GOMAXPROCS=abc` are three distinct inputs to the strip.
- Duplicate-key precedence in the exec environment — the exactly-one-entry rows under
  stories 2 and 5. This is the class watch-out #9 names.
- Boundary fixture counts — the single-fixture row. A zero-fixture tree never reaches
  the derivation: `fixtures()` already returns the absent-harness error first.
- Re-run idempotency — **Won't handle**: the derivation is a pure function of two
  integers, so a second run on the same machine derives the same bound by construction.
- Host-backed filesystems under I/O pressure — **Won't handle** as an assertion. Reducing
  concurrent inner gates is precisely how this arm reduces that pressure, but the effect
  is the story 1 measurement, and the profile's own note records that WSL2 `fsync` stalls
  are not gate-observable.
- Control bytes in environment values — **Won't handle**: the environment is passed
  through to `exec` verbatim exactly as it is today, and this arm neither renders it nor
  routes it into a `toon.Table`.
- Paths with spaces or glob characters, special files in discovery paths, symlinked
  invocation, cwd deeper than the repo root, invocation through every shipped surface —
  **Won't handle**: no path, no discovery, and no invocation surface changes in this arm.
- Interrupt mid-loop — **Won't handle**: the worker pool's shape, its job channel, and
  its cleanup are unchanged; only the worker count moves, and cancellation was never
  keyed to it.
- Required tool missing from PATH, non-TTY stdin, destructive worktree state —
  **Won't handle**: no tool lookup, no prompt, and no worktree operation exists in this
  arm.

## Out of scope

- **Capping the outer conformance and contract phases' `go test` width.** A separate
  capability: it changes the gate's own phase environment in `internal/gate` rather than
  canary's nesting, and the map's ticket #2 evidence says it buys nothing today because
  wall clock is now conformance-bound. Dormant unless contention flakes persist after
  this arm ships. Estimate: 3 edits, 2 gate runs.
- **Removing the hardcoded `-count=1` and letting Go's test result cache work.** A
  separate FT91 arm and a separate capability: it changes what a green gate proves about
  freshness, which is an oracle-semantics decision rather than a scheduling one.
  Estimate: 4 edits, 3 gate runs.
- **A shared hermetic build cache across inner gates, and caching keyed on the pinned
  gate subject.** Separate FT91 arms: both add cache infrastructure and its invalidation
  rules, which is new machinery rather than a budget for machinery that already exists.
  Estimate: 12 edits, 6 gate runs.
