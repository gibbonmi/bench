# Gate budget (FT171)

Status: ready

## Destination

One machine-wide resource budget the whole gate draws from, at every nesting
depth, so the run cannot oversubscribe the box — and less duplicated work
inside the sweep it schedules. Two levers under one destination: bound the
demand, and reduce it.

The outer gate saturates the machine on its own: load average ~123 on 16 cores
at the 2026-07-22 baseline with nothing else running. There is no quiet box
during a run, so every symptom here is self-contention and no external load is
needed to reproduce it. Reserved interactive headroom is a first-class
requirement, not a tuning preference. Every width is a formula over the
resolved budget, never a literal, because core counts differ per machine.

Supersedes `gate-concurrency`'s decision #1 (budget model) and #3 (budget =
`runtime.GOMAXPROCS(0)`, the whole box). That map's landed canary arm stands as
built; what changes is that its arithmetic stops computing from the box.

Measured state has moved twice since the #20 census: the #23 concurrency route
cut the focused `internal/gate` median from 150.85 s to 56.72 s, then the
single-build serial gate and the 2026-08-09 branch-native rebuild (`3701c4a0`)
replaced the fixture-driven workload wholesale — one host binary per top-level
run, one phase process at a time, direct mutation-to-check canaries, no
stripped-subject reruns. The target of a full gate under 2 minutes is now
exceeded threefold: #26's census on `a3b599ea` measured 38 s wall, ~51 s
CPU, and 25 peak descendants, with no Bench-owned fan-out left. #27 closed
the map 2026-08-13: the destination is met by other means, no spec follows,
and `ready` here marks terminal closure rather than spec-readiness.

## #1: What is bounded — processes, or cores?

Blocked by: none
Type: Grill

### Question

Outer phases set no `GOMAXPROCS` at all, so contract, test, and conformance
each run `go test` at full box width; `fixtureWorkers` divides
`runtime.GOMAXPROCS(0)` by the inner width, sizing canary's pool as though
canary were alone. Each layer independently claims the whole machine. Does the
fix bound concurrent phases, per-phase width, or something else?

### Answer

Bound total live demand, denominated in cores rather than processes. Process
count is the wrong unit: the demand is each `go test`'s internal parallelism,
governed by `GOMAXPROCS`/`-p`, so a bound on process count leaves each process
free to claim the box. A spawn acquires weight `w` and is launched pinned to
width `w`. This generalizes the landed inner-width pin to every layer.

## #2: How does a token cross the process boundary?

Blocked by: #1
Type: Grill

### Question

Canary's inner gates are separate OS processes, so an in-memory semaphore in
the top gate does not reach them. Hierarchical env allowance, a cross-process
broker, or an inherited descriptor?

### Answer

An inherited file descriptor holding tokens, jobserver-style. The pool owner
opens a pipe filled with N tokens; every spawn inherits the fd, reads bytes to
acquire and writes them back to release, at any nesting depth. No broker
process and no failure mode where a coordinator dies. Work-conserving: a quiet
layer's tokens are available to a busy one, which is what makes reserved
headroom cost less than a static split would. This is the established answer to
recursive-make explosion.

A process holding the fd and finding zero tokens free blocks on the read until
someone releases. That is the mechanism, not a failure: the pool is a queue.

## #3: What shape does the reserved-headroom formula take?

Blocked by: #1
Type: Grill

### Question

Fixed core reserve, proportional reserve, or both?

### Answer

Proportional and floored: `budget = max(1, resolved * (1 - r))`, with `r` a
single named constant priced by #8. Scales across a 4-core laptop and a 32-core
box, and holds on cgroup-limited CI where a fixed count could consume the whole
allocation. `GOMAXPROCS` remains the operator lever, applied before the
reserve, per `gate-concurrency` #3's no-knob ruling.

The pool owner resolves that operator width in the outer process before the
gate constructs its closed subject environment. Children receive the derived
grant through the inherited pool and the explicit width pin; ambient
`GOMAXPROCS` or `GOFLAGS` do not become undeclared subject inputs. The current
unbudgeted gate does not yet have that bridge: the memory profile observed a
12-core canary pool despite invoking the gate with `GOMAXPROCS=2`, which is an
implementation gap this destination closes rather than a second operator knob.

## #4: What decides whether a process needs a pool?

Blocked by: #2
Type: Grill

### Question

A gate-launched marker plus fd absence could red as a plumbing defect, but no
general marker exists — `BENCH_CANARY_PHASE` marks canary's sweep only. And a
hand-run of a suite that itself fans out reproduces the explosion with no
marker in sight.

### Answer

Fan-out owns a pool. Any process about to spawn concurrent children inherits a
pool or creates one; leaf processes never ask and run at their default width.
The distinction is structural rather than contextual: no marker to set, none to
delete, and no ambient condition to misread — which is what makes it safe with
an agent as the implementer.

Consequences: a hand-run of the canary package's own tests opens its own pool
and its children inherit it, closing the hole a gate-launched marker would have
left. A focused `go test ./internal/foo` never fans out, so it keeps full width
and the agent's triage loop stays fast — one process at 1x the box, which is an
ordinary test run, not the N-times-over the gate produces.

Fan-out inventory, walked 2026-08-04: exactly two Bench-owned CPU fan-out sites
exist — `internal/gate/runner.go`'s phase launch loop and
`internal/canary/canary.go`'s `eachIndex`. `internal/models` fans out
network calls, not CPU work. `internal/guards`'s scan spawns one candidate at a
time, and the artifact fixture harness's goroutine is a readiness poller;
neither is a fan-out. Every other concurrency site found is a single
`cmd.Wait()` await. Parallelism inside the tests themselves — `go test -p` and
`t.Parallel()` — is governed by #11's pinning rather than by the pool.

## #5: How are tokens recovered when a process dies holding them?

Blocked by: #2
Type: Grill

### Question

Nothing writes a dead child's tokens back, so the pool shrinks for the rest of
the run and enough of that ends in a gate that blocks forever — the worst
outcome for an agent, since a hang produces no diagnostic to act on.

### Answer

The fan-out owner reclaims on child exit. It already waits on every child it
spawns, so it knows the exact grant to write back whether the child exited
cleanly, crashed, or was killed; the common case is leak-free by construction.
A grandchild dying below a surviving parent remains possible, so an acquisition
deadline backstops it: on expiry, proceed with a one-core grant and emit a
diagnostic naming the shortfall. Never unbounded, never a hang.

## #6: What rule picks `r` from the prototype's numbers, and what counts as a symptom?

Blocked by: #3
Type: Grill

### Question

Reviewer picks live from a numbers asset (the `gate-concurrency` #2
precedent), or a mechanical rule? And measured by which observable?

### Answer

Mechanical: the smallest `r` with zero contention symptoms.

Symptom is span inflation only — a phase's full-gate span over its focused
span, against a threshold priced by #8. Load average was rejected as too noisy
and lagging to discriminate between adjacent `r` values. Deadline exhaustion
was rejected as a separate hard-fail axis: a phase inflated enough to blow a
subprocess deadline fails the inflation threshold first, so the case is
covered by one continuous measure rather than two constants.

## #7: What is measured, on which trees and under what conditions?

Blocked by: #6
Type: Grill

### Question

`go-build-cache-footprint` story 3 adds `-count=1` at `coreTestStep` and
`ConformanceSuiteArgv`, turning a mostly-cached test phase into real execution
that overlaps the critical path. Numbers taken before it lands may not survive
it. And the landed arm measured on a box running only the gate.

### Answer

Both trees: today's tree, then a throwaway branch adding `-count=1` at the two
owners (2 lines, no spec). `r` is picked against the post-spec profile; the
today's-tree runs also price what story 3 costs before it is approved.

Conditions: a box running only the gate. External background load is out of
scope — the gate saturates the machine unaided, so the symptom reproduces
without it, and every existing figure in this map was taken that way and stays
comparable. Consequence recorded honestly: the prototype cannot demonstrate the
ten-minute case is gone. That claim rests on the ceiling holding by
construction, not on measurement.

## #8: What are `r` and the span-inflation threshold?

Blocked by: #7, #26, #27
Type: Prototype

### Question

Throwaway pool implementation, then measure wall-clock and worst-phase span
inflation across candidate `r` values on both trees, and across #11's two grant
splits. How many repetitions state the variance honestly? Deliverable is a short
numbers asset. Discard a sample whose subject does not match the build under
test rather than reporting it; record the exact commit, worktree, and run time
beside every sample. Also measure what #9's narrow baselines actually save,
which this map declined to assert.

### Answer

Retired unpriced by #27's ruling: no saturating class remains to certify `r`
or a grant split against, and the fixtures this prototype would have swept
(#9's baselines included) were retired with the fixture-driven architecture.
A future phase-overlap reintroduction re-prices from a fresh census rather
than from this ticket.

## #9: Can behavior-owned baselines narrow soundly?

Blocked by: none
Type: Grill

### Question

Behavior-owned fixtures are already scoped: each carries a `TEST` file naming
its owning contract test, `biteArgs` turns it into an anchored `-test.run`, and
the package binary compiles once and is shared. What is not narrowed is the
baseline. `group()` keys these fixtures as `contract:<pkg>`, one baseline per
contract package, run `wideBaseline` — the full package binary over an empty
tree. That is roughly 133 + 52 + 39 + 20 + 4.5 ~= 250 worker-seconds, now the
dominant behavior-owned cost.

The wideness is deliberate: a scoped baseline "prints a fraction of what the
empty tree can produce, so an EXPECT the wide run already emits goes
unflagged." But `gate-pipeline` #5 made the opposite call on the conformance
axis — one shared scoped baseline per check group, sound because the guard's
premise is "would EXPECT match with no mutation under the same run shape the
fixture pays", making scoped-vs-scoped the consistent comparison. Does that
argument carry to a per-`(package, test)` baseline here, or do the two axes
differ in a way the current wideness is protecting against?

### Answer

It carries. Narrow the baseline to the test its group names and re-key
`group()` from `contract:<pkg>` to `(package, test)`; without the re-key,
fixtures naming different tests would be graded against a yardstick from a test
they never run.

The asymmetry currently runs the wrong way, which is why this is a correctness
fix and not only a cost one. `runFixture` grades a narrowed run
(`-test.run ^TEST$`, output from one test) against `strings.Contains` over a
wide baseline carrying every test's empty-tree output. Any other test emitting
a string containing this fixture's EXPECT flags it vacuous, though the
fixture's own run could never produce that output. The wide yardstick catches
nothing a narrow one misses: if the named test emits EXPECT unmutated, a
baseline running that test sees it.

Widening the fixture run to restore symmetry the other way was rejected — it
undoes lever 1, the change that took the gate from 135–172 s to 128 s, and buys
no soundness.

Consistent with the documented purpose in `projects/benchkit.md`: a
behavior-owned EXPECT is never checked for mutation-specificity, and its
baseline is only a collision screen against infrastructure noise. Narrowing
preserves that screen — noise emitted while the named test runs still appears —
and drops only noise from tests the fixture run could never execute.

Unscoped conformance fixtures keep today's full baseline, per `gate-pipeline`
#5. The saving is real but unmeasured: a narrow empty-tree baseline may still
be most of a wide one if a single test dominates the package. Measure it with
#8 rather than asserting it here.

## #10: Does the Go toolchain participate in an inherited token pool?

Blocked by: #2
Type: Research

### Question

If `go build`/`go test` honour a jobserver-style descriptor, tokens govern the
toolchain's own internal parallelism directly. If they do not, our spawns hold
the tokens and pin each child's width instead — sufficient, but a different
mechanism to build and to test. Answer from primary Go documentation and a
runnable probe, not from recollection.

### Answer

No. Probed 2026-08-04 on go1.25.0 linux/amd64: `go help build` documents `-p n`
("the number of programs, such as build commands or test binaries, that can be
run in parallel", default `GOMAXPROCS`) as the only parallelism lever, and
grepping `$(go env GOROOT)/src/cmd/` for `jobserver` and for `MAKEFLAGS`
returns zero hits.

So tokens govern Bench's own spawns and each spawn is pinned to its grant; the
`go` command neither draws from nor returns to the pipe. This is the sufficient
branch #2 anticipated, not a blocker.

The probe also found that one `go test` invocation carries two parallelism
knobs whose product is its real demand: `-p` for concurrent package test
binaries, `GOMAXPROCS` for threads inside each. Outer phases set neither today.
Pinning an invocation to weight `w` therefore means constraining both, which
#11 rules on.

## #11: A spawn is granted weight `w`. How is `w` split between `-p` and `GOMAXPROCS`?

Blocked by: #10
Type: Prototype

### Question

One `go test` invocation's demand is the product of `-p` (concurrent package
test binaries) and `GOMAXPROCS` (threads inside each). Pinning to `w` means
constraining both, and the split is a real behavioral choice: package-level
parallelism is what the artifact split exposed to remove 38.09 s, while
in-binary threads are what a phase bounded by one heavy package needs.

### Answer

Priced by #8, against exactly two candidates:

- `-p = w`, `GOMAXPROCS = 1` — the whole grant spent on concurrent package
  binaries, each single-threaded, so one token means one runnable thread.
- `-p = 1`, `GOMAXPROCS = w` — one package binary at a time with `w` threads
  inside, matching the landed inner-width pin's shape.

A balanced or square-root split is rejected: it is optimal for neither failure
mode and no evidence justifies the arithmetic.

## #12: Do the cheap phases draw from the pool?

Blocked by: #4
Type: Grill

### Question

`gofmt`, `vet`, and `shellcheck` are single short processes. Exempting them
preserves today's fast fail on formatting and lint; requiring them to acquire
keeps the ceiling exact.

### Answer

Every gate-launched phase acquires; the trivial ones take a single token. No
exemption list to maintain and nothing runs outside the ceiling. Accepted cost:
a fast-failing `gofmt` can queue behind heavy phases, so the cheapest red
arrives later than it does today.

## #13: What actually makes the expensive packages expensive?

Blocked by: none
Type: Research

### Question

The 2026-08-04 census ranked the cost but does not explain it, and it moved the
map's premise: `internal/preflight` (168 s wall, 369 s CPU) and `internal/gate`
(147 s wall) carry 40% of all serial wall and are neither contract nor canary,
so they appear in no figure in `gate-critical-path.md`. The workload is three
shapes — saturating, idle, serial — and a token pool addresses only the first.

Probes are listed in the asset. The decisive one: `SharedBuildCacheEnv` exists
in `internal/contract/command.go` and is honoured by `scripts/build-artifacts.sh`,
but nothing under `.bench/`, `bin/`, or the gate sets it, so dev contract tests
build against private disposable caches and recompile from cold on every build.
Run `posture` with and without the opt-in and compare CPU. If CPU collapses,
repeated cold compilation — not contention — is what the saturating class is
paying, and the fix is cache posture rather than scheduling.

### Answer

The expensive rows do not describe one dev-gate workload.

`posture` already opts every package run into shared build caches through its
`TestMain`; the gate does not need to set the environment itself. A controlled
two-build A/B still proved the mechanism matters: shared posture used 8.71
seconds CPU and 8.18 seconds wall, while the hermetic counterfactual used 164.65
seconds CPU and 43.31 seconds wall before failing the expected shared-posture
assertion. Cold recompilation is expensive, but it is not an unimplemented dev
cache fix. `posture` remains saturating because the package deliberately proves
the hermetic default as well as shared posture.

`internal/preflight` is not in the dev gate at all. The conformance registry
excludes it from the dev package set and `bench prep-release` runs it as a
release-only suite. Its 169.93-second green timing is ship-tier cost. Per-test
timing attributes 66.06 seconds (38.9%) to the bootstrap test whose full verify
stubs every phase except a real, hermetic artifact build; vulnerability is
stubbed there and the direct vulnerability tests are negligible.

So #8 prices only current dev-tier saturating subjects. Idle subjects are
timeouts, serial and mixed subjects cannot spend a wider grant, and ship-only
subjects cannot certify a dev reserve. Whole-gate wall remains the global
outcome, but the zero-contention acceptance signal must come from the class the
pool can affect. The measurements and ranked per-test table are in the census
asset.

## #14: Which workload class is allowed to certify `r`?

Blocked by: #13
Type: Grill

### Question

#6 currently says the smallest `r` with zero span-inflation symptoms, without
restricting the acceptance subjects. Should #8 apply that rule only to the
current dev-tier saturating class, while reporting whole-gate wall globally and
reporting idle, mixed, serial, and ship-only costs without letting them certify
the reserve?

### Answer

Yes. #8's zero-contention acceptance subjects are the current dev-tier
saturating class. Whole-gate wall remains the global outcome and every class
is still reported, but idle, mixed, serial, and ship-only work cannot certify
the core reserve because the pool cannot improve it.

## #15: Where does one gate run duplicate expensive work?

Blocked by: #14
Type: Research

### Question

Trace the current dev and ship workflows from their public entry points through
every nested gate, package suite, artifact build, and canary bite. For each
expensive child, identify whether the same job already ran elsewhere in the
same workflow, whether the repetition is an independent mutation proof or
accidental orchestration, and which owner should collapse accidental
duplication. Include current `internal/gate` per-test timings: a repository this
small should not need minute-scale test packages without a named external or
nested workload explaining them.

### Answer

The public dev workflow does not launch the same whole gate twice by accident.
`bench commit` asks the gate once and reuses an exact green verdict when one
exists; the final-check workflow explicitly forbids a pre-commit gate; and
release preflight's `gate` phase reaches the same reuse path. The dev phase
table also partitions ordinary package tests: core excludes contract,
conformance, and release-only packages; contract and conformance own their own
surfaces. Targeted `-race` tests and canary mutations repeat selected behavior
under a different oracle or tree deliberately.

The duplication is below that public boundary, in two forms.

First, one gate decision repeatedly derives the same source fact. A focused
`TestPublicDocumentClassesProjectTheirExactCheckPartition/projects/benchkit.md`
run took 2.90 seconds and launched 415 Git processes. The same uncommitted
fixture tree was materialized 47 times: 118 `rev-parse`, 94 `read-tree`, 84
`cat-file`, 47 `add`, 47 `write-tree`, 24 `ls-tree`, and one `init`. The
production path has the same ownership defect even though a committed tree
avoids the fixture's failed-`read-tree HEAD` fallback: `scopeComponents`
independently asks the conformance-check, canary, and component resolvers for
snapshots; subject and stripped-subject construction ask again; pre/post-run
inspection rebuilds them again; repeated blob reads are not memoized. Each
local helper says it snapshots once, but no gate-decision-wide snapshot owns
the fact.

The test then multiplies that defect. The 35.14-second public-document-class
test has 21 rows. Each row runs the real gate engine for a mutation and a
deletion and runs a restoring green gate after each: about 84 gate executions
to prove one mapping table. The complete `internal/gate` package was green at
153.40 seconds wall and 116.33 seconds CPU. Eighteen top-level tests at or
above two seconds carried 96.12 seconds; only the 35.14-second matrix exceeded
ten. Six cancellation/deadline cases also pay the real two-second termination
grace. This is orchestration and test-seam cost, not slow Go computation.

Second, package tests repeatedly cross a seam wider than the behavior they
assert. `internal/conformance` was green at 85.11 seconds wall and 138.49
seconds CPU. Its 83 documentation/workflow fixture rows each call the complete
dev conformance table, totaling 35.13 seconds. Gate-entry tests spend another
25 seconds creating fresh temporary Go modules and invoking the shell entry,
including 19.95 seconds across eight freshness-failure variants. The existing
single-check `RunConformance` interface is the correct seam for fixture rows;
one representative entry proof can keep the shell path while the failure
matrix lives at the freshness interface.

The remaining long packages confirm the shape rather than hiding another
whole-gate recursion. `internal/specbuild` is 145 Git-backed lifecycle tests
spread over 43.53 seconds; `internal/worktree` is 134 such tests spread over
29.65 seconds. `internal/contract/runtime` is 69.78 seconds wall and 197.73
seconds CPU, led by the 58-case FT78 action proof ledger (13.68 seconds), the
four-route freshness matrix (12.33 seconds), and real cross-process gate/spec
tests. `internal/publication` spends 30.03 of 30.49 seconds waiting on
`http://127.0.0.1:1` through `http.DefaultClient` with `context.Background()`;
the dev-tier publication contract carries the same unreachable-port posture.
That is an uncontrolled network timeout, not repository work, and remains
owned by FT87.

Ship orchestration does contain literal duplicate jobs:

- A current dev-green verdict is a precondition, but `core-tests-ship` uses the
  ship package enumeration and therefore reruns all 54 dev-core packages just
  to add the three release-only packages. Its comment describes a release-only
  step; its command is not one.
- Ship conformance reruns the dev checks as a superset before adding the one
  release-evidence probe and cross-compile assertion. This is currently an
  explicit lifecycle-final policy, not an accidental call graph, but it is
  inherited reproof that must be priced as such.
- Release preflight's `vet ./...` repeats the dev gate's identical vet job.
  Its `race ./...` executes all 70 packages under the detector after the dev
  gate's three-test authoritative race registry; the detector is a distinct
  oracle, but broad integration suites whose meaningful work is in ordinary
  child binaries gain little race coverage from being rerun this way.
- `prep-release` builds `dist/artifacts` before invoking release preflight,
  whose registered `artifacts` phase rebuilds the same root into the same
  output. Conformance ship also builds the same snapshot in an authenticated
  clone. Two ship canaries build two distinct mutated trees; those are
  independent bite proofs, not duplicate inputs, although each repays the full
  compile-and-package pipeline.

Owners follow the highest observable seam. Gate evaluation should own one
immutable parsed tree/blob snapshot for all identity families and pre/post
inspection. Mapping matrices should grade the resolver interface exhaustively
and retain a small end-to-end bite set. Conformance fixture rows should pass
their registered single-check scope. Pre-release should enumerate only the
three release-only suites, consume one artifact build, and either inherit
exact dev vet evidence or stop claiming it as a separately executed phase.
Race coverage should come from one authoritative registry, extended with
ship-only cases when needed.

The ordering consequence is decisive: with `-count=1`, the 153.40-second
focused `internal/gate` package alone exceeds the destination's two-minute
whole-gate target. A token pool cannot make a mixed/serial package faster and
may make it slower. Accidental demand must be collapsed before #8 prices the
reserve; otherwise the prototype is optimizing contention around an impossible
critical path. Full measurements and the workflow graph are in the census
asset.

## #16: Must demand reduction land before the pool is priced?

Blocked by: #15
Type: Grill

### Question

Should FT171 order the dev-path snapshot and over-wide-test-seam reductions
ahead of #8, then re-run the census and price the pool against the reduced
workload, while sending ship-only duplicate-proof policy to a separate
follow-up? The recommendation is yes: it is the only order in which the
two-minute target is achievable and #8 measures the workload the finished gate
will actually schedule.

### Answer

Yes. Demand reduction lands before #8, as three independently green specs in
dependency order:

1. `gate-evaluation-snapshot` makes one gate evaluation own one immutable,
   generation-scoped parsed tree/blob snapshot. Pre- and post-evaluation remain
   separate generations, prospective trees use the same snapshot contract
   through their own source adapter, and all current identity derivations move
   onto it without changing oracle behavior.
2. `gate-decision-test-seam` moves the exhaustive component/check mapping matrix
   to the decision seam, retaining representative full-engine tests that prove
   wiring and failure behavior without rematerializing trees for every policy
   row.
3. `conformance-harness-scope` gives each conformance fixture the registered
   check it owns and moves freshness assertions to the lower seam, retaining
   representative shell-level integration coverage.

After those three land, #20 re-runs the cold workload census. Only that reduced
workload may feed #8's pool and reserve pricing. The first ordered spec is
`gate-evaluation-snapshot`; the other two are downstream and out of its scope.

The three landed reductions were necessary but not sufficient: #20 found the
focused package floor still above the destination. #21 therefore continues
demand reduction before #8 without reopening this ordering decision.

Ship-tier duplicate-proof policy is separate shaping because it changes the
meaning and composition of publication evidence, not merely the dev gate's
implementation. FT87 continues to own publication timeout behavior.

## #17: Does one gate evaluation own one snapshot generation?

Blocked by: #16
Type: Task

### Question

Write and land the first ordered spec, `gate-evaluation-snapshot`. One pre-gate
generation and one post-gate generation each parse their tree and blobs once;
working-tree and prospective-tree inputs satisfy the same snapshot contract.
Migrate component, check, canary, subject, and stripped-subject identity reads
while preserving drift detection, prospective-tree semantics, fail-closed
errors, and emitted evidence. Prove both semantic equivalence and the bounded
process/materialization count. Moving the mapping matrix is not part of this
spec.

### Answer

Resolved and landed. `gate-evaluation-snapshot` promoted candidate `60b44de` as
`c3d8e47` and retired at `60842e0`. One evaluation now owns its accepted pre
generation and a distinct post generation; working and prospective sources
share the generation contract, all five identity families consume it, and the
operation-count controls bound materializations, listings, and repeated blob
reads without changing drift detection, fail-closed posture, or evidence.

## #18: Does the exhaustive gate matrix run at the decision seam?

Blocked by: #17
Type: Task

### Question

Write and land `gate-decision-test-seam`. Exercise the exhaustive mapping table
at the pure decision boundary, then keep representative full-engine tests that
prove source wiring, materialization, and fail-closed behavior. Production
oracle behavior does not change. Owner: a fresh spec lifecycle beginning with
`/bench-write-spec`, followed through build, review, promotion, and retirement.

### Answer

Resolved and landed. `gate-decision-test-seam` promoted candidate `f7f0ea8c`
as `3cce8aab` and retired at `acf5c9da`. The exhaustive public-document
partition now runs at the decision seam; its seeded full-engine control and
failure-path representatives retain the composition proof.

## #19: Does each conformance fixture run only its registered check?

Blocked by: #18
Type: Task

### Question

Write and land `conformance-harness-scope`. Route each fixture through its
registered check, test freshness at the lower seam that owns it, and retain
representative shell-level integration coverage. The executable check registry
remains the single source of policy. Owner: a fresh spec lifecycle beginning
with `/bench-write-spec`, after #18 retires, followed through build, review,
promotion, and retirement.

### Answer

Resolved and landed. `conformance-harness-scope` promoted candidate
`cd0b7300` as `6cf1816e` and retired at `c8ea95ea`. Its 83 direct fixture bites
now execute the canary-resolved ordinary check, including fixture-level
`CHECK` precedence, while broad controls retain the full table.

## #20: What workload remains after demand reduction?

Blocked by: #19
Type: Research

### Question

Re-run the cold package census and the focused-versus-in-gate span probes after
#17–#19 land. Update the census asset with exact commits, subjects, repetitions,
CPU, wall, and process/materialization counts. Report the reduced workload that
#8 must price; do not carry the pre-refactor demand forward as an assumption.

### Answer

Resolved on exact commit `eb6845f` (tree `cd2ece9`), 2026-08-07 local time.
The 71-package serial census was green at 767.25 s wall and 1110.99 s CPU.
`internal/conformance` fell to a 25.83 s three-run mean, but the reduced
workload still cannot feed #8: `internal/gate` was 147.34, 159.70, and 152.11
s focused, so its minimum alone exceeds the whole-gate target. One green fresh
gate was 246 s, with `internal/gate` at 230.663 s, `specbuild` at 177.258 s,
and contract `runtime` at 152.126 s.

The same gate peaked at 97 concurrent descendants while seven outer phase
roots overlapped. A two-second ancestry sample saw at least 1,682 distinct
child PIDs across the run; that is a lower bound because no process-event
tracer was installed. The public-document matrix now performs one generation
capture per changed state and no per-row full-engine run, but its 57 captures
still cost 20.11–21.71 s. Removing that cost alone would not put the focused
package safely below 120 s. Full method, CPU figures, span inflation, process
classes, and the cache ruling are in the census asset.

Do not clear the ambient Go cache for this census. `-count=1` already disables
test-result reuse, while the compile-only warm-up separates test work from
build work. Clearing the shared build cache would add a different cold-compile
workload, perturb other worktrees, and invalidate comparison with the existing
warm-cache gate; the sampled gate still grew that cache by 133.0 MB.

## #21: Which remaining demand must land before outer-width pricing?

Blocked by: #20
Type: Research

### Question

Trace the current `internal/gate` and `internal/specbuild` focused workloads
from the post-reduction census. Identify independently useful demand reductions
that put the longest fresh focused dev package below the two-minute whole-gate
target with enough margin for gate setup, without weakening an oracle or
reclassifying a test as cached evidence. Separate work a token pool can improve
from mixed, serial, fixed-deadline, and process-materialization cost. Retain
per-test spans and process/materialization counts, and turn any policy choice
into a new Grill ticket rather than choosing it in research.

### Answer

Resolved on #20's exact commit `eb6845f`, 2026-08-07. Both packages are
strictly serial test chains: per-test elapsed sums equal package walls
(`internal/gate` 139.67 s across 241 tests against 140.80 s wall;
`internal/specbuild` 57.38–59.81 s across 192 tests against 57.80–60.05 s
walls), so the focused floor is the sum of every test's subprocess waits and a
token pool cannot shorten it — the pool only stops sibling phases from
inflating it.

`internal/gate`'s demand is distributed: the largest test is 18.65 s and 40
tests carry 106.8 s. The sized serial cuts — observing the 2 s termination
grace instead of paying it live in four of five cancellation tests (~8 s),
synthesizing the document matrix's 57 generation captures at the decision seam
(~15 s), and memoizing the per-resolution `go list` module-closure derivation
(~10–20 s, estimated) — project a 105–125 s floor. That fails the target with
margin, so no set of single-test fixes reaches the destination and the
remaining lever is concurrency inside the package run. That lever is blocked
by process-global state, not test isolation: production `kitRoot` reads
ambient `BENCH_KIT`, 51 fixture constructions pin it with `t.Setenv` (which
excludes `t.Parallel`), and two helpers chdir. `internal/specbuild` has no
such coupling — its cost is pure process churn (35,952–35,957 Git spawns per
run, 60% `rev-parse` re-deriving repository facts) and test-only `t.Parallel`
is structurally available.

Gate Git-spawn counts are deterministic: exactly 13,156 per package run,
dominated by 5,343 `rev-parse` and 1,253 `write-tree`. Choosing the
concurrency route is a workload-shape and seam decision, so it moved to #22
rather than being selected here. Per-test tables, spawn histograms, probe
timings, and method are in the census asset.

## #22: Does the dev gate adopt intra-package test concurrency, and through which mechanism?

Blocked by: #21
Type: Grill

### Question

#21 found both long dev packages are strictly serial chains and that the
enumerated serial cuts project only a 105–125 s `internal/gate` floor — short
of the 120 s whole-gate target with margin. The remaining lever is concurrency
inside the package run, with three routes:

- `t.Parallel` inside `internal/gate`, which first needs a kit-root injection
  seam: production `kitRoot` reads ambient `BENCH_KIT`, 51 fixture
  constructions pin it with `t.Setenv` (mutually exclusive with
  `t.Parallel`), and two helpers chdir.
- Splitting `internal/gate`'s tests into several test packages, letting the
  existing `go test -p` package-level parallelism and per-process env
  isolation do the work without a production seam change — at the cost of
  exporting or relocating the internals those tests reach.
- Test-only `t.Parallel` in `internal/specbuild` (its two `t.Setenv` tests
  stay serial), independent of either gate route.

The choice changes the workload shape #8 prices — a parallelized package
becomes saturating, and its width is exactly what the pool's grant pin
governs — so it precedes any further census. Its answer also decides whether
the sized serial cuts (grace observation, synthesized matrix generations,
closure memoization) ride the same spec or land separately, and orders the
resulting Task tickets.

### Answer

Route one, plus the specbuild lever. `internal/gate` gains a kit-root
injection seam — production `kitRoot` takes its root as an explicit input
instead of reading ambient `BENCH_KIT`, retiring the fixture constructors'
`t.Setenv` pins (two constructor pins reached from roughly 52 construction
sites). The chdir-helper tests keep their chdir and stay serial — cwd is the
input they exercise, so that pin cannot retire (spec-time amendment,
2026-08-07) — and the package's tests then adopt `t.Parallel`. The
package-split route was rejected: it leaves the ambient read in place and
works around it, and every extra test package repays binary build and fixture
setup, the process-materialization cost #21 measured. `internal/specbuild`
adopts test-only `t.Parallel` as well; its two `t.Setenv` tests stay serial.

The three sized serial cuts land separately, after the concurrency spec: the
seam spec decides the workload shape #8 prices and is not bloated with
independently green cuts, which still pay under parallelism because they
remove CPU and spawns rather than wall alone.

Routing and order: #23 carries the seam and gate `t.Parallel` as a full spec
lifecycle; #24 lands specbuild `t.Parallel` light-path and may land
immediately; #25 lands the three cuts as light-path tickets after #23, each
stopping for a reviewer decision if it turns out to cross a declared seam;
#26 re-runs the census after #23–#25 land and replaces this ticket in #8's
blockers.

## #23: Land the gate kit-root seam and `t.Parallel` adoption

Blocked by: #22
Type: Task

### Question

Write and land `gate-test-concurrency`: production `kitRoot` becomes an
injected input rather than an ambient `BENCH_KIT` read, the 51 `t.Setenv`
fixture constructions and two chdir helpers retire, and `internal/gate`'s
tests run under `t.Parallel`. Oracle behavior does not change; the closed
subject environment keeps ambient env out of subject inputs per #3. Owner: a
fresh spec lifecycle beginning with `/bench-write-spec`, followed through
build, review, promotion, and retirement.

### Answer

Resolved and landed. `gate-test-concurrency` staged at `f2a925ee` and retired
at `bc2d1aff` (2026-08-08): production kit-root became an injected input, the
fixture `t.Setenv` pins retired, and 192 of 245 top-level gate tests adopted
`t.Parallel` with 53 reasoned serial holdouts. The exact-candidate focused
package median fell from 150.85 s to 56.72 s; the measured shared-fixture
follow-up then took width one to 111.67 s / 699,904 output blocks and width
two to 78.79 s / 699,712 blocks — flat write volume across widths, so
concurrency compresses wall without removing intrinsic fixture work. The
2026-08-09 branch-native rebuild (`3701c4a0`) has since replaced that
fixture-driven suite wholesale — `internal/gate` now carries five in-package
top-level tests — so these figures are the route's landed evidence, not the
current workload.

## #24: Land specbuild test-only `t.Parallel`

Blocked by: #22
Type: Task

### Question

Mark `internal/specbuild`'s tests parallel, keeping its two `t.Setenv` tests
serial. Test-only, one independently green ticket, no declared seam — rides
the standing light-path approval and may land before or alongside #23.

### Answer

Retired as moot. The spec-build lifecycle core was removed wholesale
(`dae240df`) and `internal/specbuild` no longer exists, so there is no package
to parallelize. Removed from #26's blockers per the reviewed roadmap drain.

## #25: Land the three sized serial cuts

Blocked by: #23
Type: Task

### Question

After #23 retires, land as light-path tickets: observe the 2 s termination
grace through an injected duration keeping at least one live-cascade proof
(~8 s), synthesize the document matrix's 57 generation captures at the
decision seam (~15 s), and memoize the per-resolution `go list` module-closure
derivation (~10–20 s, estimated). Each ticket stops for a reviewer decision
if its change turns out to cross a declared seam instead of riding light path.

### Answer

Moot as sized, flagged for reviewer veto against the roadmap's "land #25's
cuts". All three estimates priced a workload the branch-native rebuild
(`3701c4a0`) deleted: the document matrix and the in-package cancellation
tests are gone, and the grace is already injectable through context
(`processGroupGrace`, landed with `af23f587`) with the 2 s default intact.
The `go list` module-closure derivation survives unmemoized in
`internal/gate/gate_go.go`; whether it is still worth memoizing is priced by
#26's census on the current workload rather than landed on a stale estimate,
per #20's own rule against carrying pre-refactor demand forward.

## #26: What workload remains after the concurrency route lands?

Blocked by: #23, #25
Type: Research

### Question

Re-run the cold package census and focused-versus-in-gate span probes on the
current baseline — the branch-native, single-build serial gate: one host
binary per top-level run, one phase process at a time, direct
mutation-to-check canaries, no stripped-subject reruns — updating the census
asset with exact commits, repetitions, CPU, wall, and the new workload-shape
classification. Re-walk #4's fan-out inventory while there: the rebuild
retired `eachIndex` and serialized the phase table, so the pool's client
sites must be re-enumerated before #8 prices reintroduced width. Include the
surviving `go list` module-closure derivation from #25. Only this reduced,
reshaped workload may feed #8's reserve and split pricing.

### Answer

Resolved 2026-08-13 on exact commit `a3b599ea`, same 12-online-CPU host. The
serial package census is 64 packages at 105.2 s wall / 46.7 s CPU (#20:
767.25 s / 1110.99 s); two fresh gates were 38.26 s (green) and 38.00 s
(red) wall at ~51 s CPU each, peaking at 25 concurrent descendants against
97 pre-rebuild, ~1.3 of 12 cores busy on average. Phase spans: test 31.9 s
(floored by `internal/publication`'s FT87 30 s idle wait), race 2.4 s,
system 1.4 s, everything else sub-second, ~1.6 s setup. The fan-out re-walk
found zero Bench-owned CPU fan-out sites — the only remaining width owner is
`go test`'s own `-p` inside the single test phase — and the `go list`
closure derivation is 0.10–0.31 s warm, so the #25 residual is dead. The red
was `TestListCommandCheckedInCompletedAssignmentTerminalPair` in
`internal/worktree` (2 reds in 6 package runs this session, FT203's rate but
not its named test — the flake family is wider than the roadmap row). No
dev-tier saturating class remains for #8 to certify against; whether the
pool destination survives at all is #27. Full tables and method are in the
census asset.

## #27: Does the pool destination survive the serial baseline?

Blocked by: #26
Type: Grill

### Question

The destination exists to stop the gate oversubscribing the box, on evidence
of load ~123 and a ten-minute worst case. #26 measured the current baseline
at 38 s wall, ~51 s CPU, 25 peak descendants, ~1.3 of 12 cores busy: the
symptom cannot reproduce, the two-minute target is exceeded threefold, and
zero Bench-owned fan-out sites remain to hold a token. Does #8 still price
`r` and the grant split — against what symptom, on what class? — or do #8
and the token-pool destination retire as achieved by the single-build serial
gate, with any future reintroduction of phase overlap reopening shaping
first? Recommendation: retire #8, mark the map's destination met by other
means, and close the map; the mechanism decisions (#1–#5, #10–#12) remain
recorded as settled design if overlap ever returns.

### Answer

Retire and close, reviewer-accepted 2026-08-13. The destination is met by
other means: the single-build serial gate and the branch-native rebuild
bound demand by construction, and no workload class remains for a reserve to
protect. #8 retires unpriced. The mechanism decisions (#1–#5, #10–#12) stand
as settled design; any future reintroduction of phase overlap reopens
shaping and re-censuses before adopting them.

## Not yet specified

## Spec-writer discretion

- The token encoding in the pipe, and how the acquisition deadline in #5 is
  expressed, provided the fail-toward-less-concurrency posture holds.
- Where the pool owner's bookkeeping lives, provided reclaim-on-child-exit is
  observable from one place.
- #8's sweep order — comparing #11's two splits at one fixed `r` before sweeping
  `r` keeps the matrix from multiplying, provided the chosen split is re-checked
  at the selected `r`.

## Out of scope

- Scoped or diff-based gating of any kind: FT91 ruled it unsound here, contract
  and canary being behavior contracts with no file-to-test map, and that ruling
  stands. Scope is not a speed lever.
- Weakening any check to buy wall-clock. Green keeps meaning the same thing.
- Merging same-check fixture runs: `gate-pipeline` #5 rejected it because each
  fixture is a distinct mutated tree and must be graded alone. #9 asks about
  baselines only.
- A cache quota, automatic eviction, or the Go build-cache footprint itself —
  owned by `go-build-cache-footprint`.
- Ship-tier duplicate-proof cleanup, which needs separate shaping because it
  changes publication evidence rather than the dev gate's implementation.
- Publication timeout behavior, which remains owned by FT87.

## Sources

- Path: `decisions/gate-critical-path.md`
  Supports: the 89.91 s entry state and 85.415 s vs 50.917 s `posture` span behind #6's choice of span inflation.
  Drift: historical trigger evidence on the artifact-split landed tree; #20 is the current workload account.
- Path: `decisions/gate-concurrency.md`
  Supports: the landed canary arm, its ~123 load baseline, and decisions #1 and #3 this map supersedes.
  Drift: describes shipped code; re-read if the canary arithmetic changes.
- Path: `decisions/gate-pipeline.md`
  Supports: #5's scoped-baseline soundness argument that #9 asks to extend.
  Drift: none expected while that decision stands.
- Path: `ROADMAP.md`
  Supports: FT171's contention evidence, including the 12-core sample where both `TestSetupConflictContracts` FIFO cases exhausted 15 s deadlines under overlap and then passed 3/3 focused at ~0.43 s.
  Drift: a working prioritization document; the row moves as the work lands.
- Path: `decisions/assets/gate-budget-cpu-wall-census.md`
  Supports: #13's three-shape finding, cache A/B, preflight ownership and per-test timing, #20's post-reduction 71-package census, focused repetitions, exact-subject gate span, process fan-out, and cache ruling, plus #21's serial-chain finding, per-test attribution, deterministic Git-spawn histograms, and concurrency constraints.
  Drift: #13–#21's figures are pre-rebuild and historical; the decision #26 section is current as of `a3b599ea`, and #27 closed the map with no width pricing to follow.
- Path: `decisions/assets/gate-budget-memory-profile.md`
  Supports: #1–#5's machine-wide process-boundary budget, the non-recursive primary/stripped/canary overlap, the current operator-width plumbing gap, and #20's required memory/process/I/O observables.
  Drift: one green run on `6607236` and one 12-core host before #18–#19; mechanism evidence only, never authority to price #8.
- Path: `decisions/assets/ft171-shared-fixture-staged-binary.md`
  Supports: #23's measured residual — the shared immutable fixture binary and narrowed setup-only work behind the 111.67 s / 78.79 s width figures and the flat ~700k output-block volume.
  Drift: measured on the pre-rebuild fixture-driven suite; landed-route evidence, never the current workload.
- Path: `internal/conformance/registry/packages.go`
  Supports: #13's finding that `internal/preflight` is excluded from the dev package set and owned by the ship tier.
  Drift: code-derived; re-verify if tier ownership changes.
- Path: `internal/gate/evaluation.go`
  Supports: #17's current evaluation-owned accepted pre generation, distinct post generation, and common working/prospective source contract.
  Drift: code-derived; re-verify with the operation ceilings in `internal/gate/evaluation_test.go` if evaluation or snapshot ownership moves.
- Path: `internal/preprelease/preprelease.go`
  Supports: #15's ship step order, repeated core package enumeration, first artifact build, conformance probe, preflight, and ship canary.
  Drift: code-derived; re-verify if the ship sequence changes.
- Path: `internal/releaseevidence/registry.json`
  Supports: #15's release-preflight gate, full-race, vet, artifact, and smoke commands.
  Drift: code-derived; re-verify if release evidence phase ownership changes.
- Path: `internal/gate/runner.go`
  Supports: the current serialized phase table behind #26's baseline and the context-injectable `processGroupGrace` behind #25's mootness.
  Drift: code-derived; the pre-rebuild `canary.go`/`check_slots_test.go`/artifact-fixture sources behind #1, #9, #13, and #18's findings were deleted by `3701c4a0`, and those answers now rest on the census and memory-profile assets plus spec-retirement history.
