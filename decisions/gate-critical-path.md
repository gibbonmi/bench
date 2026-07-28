# Gate critical path (FT91, eighth arm)

## Destination

Find where the dev gate's remaining wall-clock actually lives, decide the lever
that takes it, and settle whether `decisions/gate-pipeline.md` reopens now that
slice C falsified its premise. Green keeps meaning the same thing throughout.
Decided stop condition (#5): FT91 ends at a measured full dev gate **≤60 s** on
this repo, accepting that the last stretch may pull oracle-semantics decisions
into this row.

## #1: Why did the gate absorb only 24 s of a 131 s suite win — what is the critical path now?

Type: Research

### Question
The seventh arm cut the contract suites 249 s → 118 s measured solo, but the
whole gate moved only 4m51s → 4m27s: under gate parallelism the artifact suite
inflates ~106 s → ~152 s, and the contract phase is no longer clearly the
critical path. Nobody has diagnosed the gap.

### Answer
Measured 2026-07-28; full timeline, decomposition, and citations in
`decisions/assets/gate-critical-path-timeline.md`. The gate is canary-bound end
to end: gate wall 267 s, canary solo 250 s, contract phase done at t+162. The
34 `behavior-owned` fixtures route to the contract phase, so each inner gate
re-runs the entire kit contract suite (`go test -count=1
./internal/contract/...`) at inner width 2 — ~1,700 of canary's ~2,000
worker-seconds. The artifact-suite inflation (106 s solo → 158 s in-gate) is
CPU oversubscription while contract, test, and canary's nested runs overlap
(load 15–23 on 16 cores). Floor if contract were free: unchanged — canary is
underneath. That is the whole 24 s-of-131 s gap.

## #2: Which artifact-suite tests build, and which only inspect?

Type: Research

### Question
The suite is ~20 host-only generator runs at ~3.7 s each; the
`BENCH_TEST_PREPARED_ARTIFACTS` seam exists but `artifactPreparedGeneration`
is per-test, so even inspection-only tests pay a full build. Read
`internal/contract/surface/artifact` and classify every test: mutates the
artifact set or its environment, asserts on the act of generation itself
(atomicity, promotion, refusals), or inspects prepared output only. State what
per-test scoping actually guarantees each class, so #3's independence ruling
is made against facts. Deliverable: a bucketed inventory asset with per-claim
citations.

### Answer
— (open)

## #3: Which tests may share one package-scoped prepared artifact set?

Blocked by: #2
Type: Grill

### Question
Hoisting `artifactPreparedGeneration` to package scope lets one build serve
every inspection-only test, but sharing an artifact set is a test-independence
ruling, not a build: which #2 classes may share, what shape the hoisted seam
takes, and the fail posture when a test in the sharing group mutates the
shared set.

### Answer
— (open)

## #4: Does `decisions/gate-pipeline.md` reopen, and what happens to `ft91-gate-phase-split` stories 4, 5, and 9?

Blocked by: #1
Type: Grill

### Answer
No reopen (reviewer, 2026-07-28). The shipped structure — manifest loader,
DAG scheduler, phase split — stands on its own merits; only its wall-clock
rationale was falsified, and this map now owns wall-clock. Stories 4 and 5
are accepted as the probed phases that shipped; story 9 is dropped as
unsatisfiable as specced; the spec is retired. One transfer: gate-pipeline's
"removing canary nesting" out-of-scope clause moves here, reopened and ruled
in #7.

## #5: What ends FT91 — and what does the eighth arm build?

Blocked by: #1, #3, #4
Type: Grill

### Answer
Resolved 2026-07-28 (reviewer). **Stop condition: a measured full dev gate
≤60 s on this repo retires the row** — a hard target, chosen over ≤90 s
knowing it may force the oracle-semantics decisions (verdict caching,
`-count=1`) into FT91 once the structural work lands and is re-measured.
**The arm builds, in order:** stage 1, package-scoped behavior-owned nesting
(#6); stage 2, in-process bite checks (#7); the prepared-artifact hoist
(#3, pending its independence ruling). Expected end state ~60–90 s with
oracle semantics untouched; the gap to ≤60 s, if the re-measurement shows
one, graduates the parked oracle-semantics fog below into tickets here.

## #6: How is a behavior-owned fixture's inner run scoped to one contract package (stage 1)?

Type: Grill

### Answer
Resolved 2026-07-28 (reviewer). Each fixture's inner gate runs only the one
contract package that owns its EXPECT. Binding lives in **subfamily
directories** — `tests/canary/behavior-owned/<package>/<fixture>/` — the
directory carries the fact (slice-C precedent; per-fixture binding files and
a registry table rejected). Postures follow the fifth arm: vacuity baselines
are scoped-vs-scoped per package group, missing or unknown binding is a loud
red, did-not-bite stays red, no silent fallback to a full run. Expected
canary ~250 s → ~60–80 s (floor: slowest bound package at inner width 2).

## #7: What replaces canary nesting for behavior-owned fixtures (stage 2)?

Blocked by: #6
Type: Grill

### Answer
Resolved 2026-07-28 (reviewer), refined same-day after stage 1 landed. The
nesting premise is formally narrowed: a behavior-owned fixture's bite is
proven at its owning contract test, not by a nested gate. The sweep stops
spawning nested gates for this family.

**Run shape.** `go test -c` compiles the owning contract package once per
package group; that binary is then invoked once per fixture root with
`BENCH_CONTRACT_ROOT` set. The originally-worded "each fixture root as a
subtest in one process" is not buildable — `SubjectRoot` reads the env per
call, and Go cannot re-enter a package's `Test*` functions as parameterized
subtests without a per-package registry, which would duplicate the package's
own test list and collide with the `t.Parallel` those tests already call.
Compile-once delivers the same intent: no nested gate, no process tree, no
inner-width throttle, compile paid once per group, and each fixture keeps the
clean process the root swap needs anyway.

**Ownership.** The canary sweep still invokes these runs and still owns bite,
did-not-bite, and the vacuity baseline, so `bench canary .` stays
standalone-provable and vacuity stays a canary concept. Bite and vacuity
semantics unchanged: EXPECT observed as a real failure per fixture, EXPECT
absent on the unmutated baseline under the same run shape. The gate-level
plumbing (red test → red phase → red gate) is asserted once by the gate's own
phase tests, not 33 times. Stage 1's subfamily binding is stage 2's input —
nothing is throwaway, and the inner-width pin becomes moot for this family.

**Migration casualties.** Every EXPECT must be observable in test-level output
rather than gate-level framing (the slice-C casualty class). One that is not is
an ordinary did-not-bite red, fixed by making the owning test's own failure
message carry the fact — no hybrid family, no fallback to nesting.

**Stage 1's inner-mode narrowing is removed by stage 2.** Contract scope is
behavior-owned-only, so once no such fixture spawns an inner gate,
`narrowContractScope` and the package env's inner-gate plumbing have no caller;
they go, with their now-unreachable tests. The subfamily directory binding
stays — stage 2 reads it.

**Acceptance.** No wall-clock threshold on this slice; ship evidence is a
recorded before/after in `assets/gate-critical-path-timeline.md`. #5's ≤60 s
gate rule stays the only hard bar and is judged after the artifact hoist.
Stage 1 showed a per-slice threshold did not change the decision, and
degenerate migration is caught by bite and vacuity semantics, not by a clock.
The artifact multiplier stage 1 exposed — five `surface/artifact` fixtures each
paying the artifact suite — is not stage 2's to solve; it belongs to the
prepared-artifact hoist (#2/#3), which stays the next slice.

Remaining spec-time detail: where the compiled test binaries live, and their
cleanup.

## Handoff

Rewritten 2026-07-28 for stage 2 (#7) after stage 1 shipped and was measured;
#2/#3 stay open and gate only the artifact-hoist slice.

1. **Module boundaries.** `internal/canary` owns the whole stage-2 change:
   the behavior-owned family's run mode becomes compile-then-invoke instead of
   nested-gate spawn, keeping the existing walk, the subfamily package binding,
   and per-package vacuity baseline grouping. `internal/gate` **loses** the
   inner-mode contract-phase argv narrowing and its env plumbing, which stage 2
   deletes. The `internal/contract` packages are untouched except where an
   EXPECT proves unobservable in test-level output, where the fix is that
   owning test's failure message. `tests/canary/behavior-owned/<package
   path>/<fixture>/` stays the binding — data, no second file. `internal/bounds`
   untouched.
2. **Contracts.** Per package group: one `go test -c` producing a binary, then
   one invocation of that binary per fixture root with `BENCH_CONTRACT_ROOT`
   set to the mutated tree, plus one baseline invocation on the unmutated tree.
   A compile failure for a group is a red naming the package, never a skipped
   group. EXPECT, bite, did-not-bite, and vacuity semantics unchanged; canary
   and gate exit codes unchanged. The contract-package env disappears as an
   external contract along with the inner-gate narrowing.
3. **Deep vs thin.** The sweep (walk, group resolution, compile, per-fixture
   invocation, baseline grouping) stays the deep unit behind the injected
   `Runner` seam. The compile step and per-root env are thin plumbing inside it.
4. **Black-box assertables.** Via a fake `Runner`: one compile call per package
   group present; one invocation per fixture carrying that group's binary and
   its own `BENCH_CONTRACT_ROOT`; no call carrying the canary inner-gate mode
   for this family; baseline invocation count equals the group count; a failed
   compile surfaces as that group's red. Via the gate's phase-table tests: the
   contract phase's argv is the unnarrowed subtree in every mode, and an ambient
   export of the retired env changes nothing.
5. **Gate attachment.** The canary sweep still enforces bite and vacuity per
   fixture; the kit's own suites run in the gate's test and conformance-suite
   phases. The red-test → red-phase → red-gate path is asserted once by the
   gate's own phase tests. The wall-clock outcome is not gate-assertable — ship
   evidence is the post-change measurement recorded in
   `assets/gate-critical-path-timeline.md`, with no threshold gating the slice.
6. **Hostile-input owners.** Missing or empty subfamily binding directory → the
   walk, red (empty-CHECK precedent). A bound package the kit tree lacks → the
   group resolver at compile time, red naming it. Ambient export of the retired
   contract-package env → inert, asserted by the phase-table test. A fixture
   root that fails to materialize → the sweep, red, never a silent skip.
7. **Uncertainty flags.** None blocking stage 2. Compiled-binary location and
   cleanup are spec-time within the decided shape.
8. **Rejected alternatives.** Per-package registry of parameterized subtests
   (duplicates the test list, forces `t.Setenv`, kills `t.Parallel`); plain
   `go test -run` per fixture with no pre-compile (re-pays link per fixture);
   folding the bites into the contract phase and dropping the family from
   canary; a hybrid family keeping nesting for gate-framed EXPECTs; keeping the
   now-callerless inner-mode narrowing; a wall-clock threshold gating the slice;
   resequencing the artifact hoist ahead of stage 2.
9. **Domain watch-outs.** `SubjectRoot` reads `BENCH_CONTRACT_ROOT` per call, so
   one process cannot serve two roots — the per-fixture process is load-bearing,
   not incidental. Go's exec env has no duplicate-key precedence: every pin
   strips-then-sets. The vacuity baseline must run the identical shape as its
   group's fixtures or the scoped-vs-scoped comparison silently breaks. Fixture
   base names are globally unique across all families. The five
   `surface/artifact` fixtures each pay a full artifact suite; that multiplier
   survives stage 2 by design and is the hoist's target.

Dependency order: stage 1 (#6, package-scoped nesting) — **shipped** → stage 2
(#7, compile-once bites) → artifact hoist (pends #2/#3) → re-measure against
the ≤60 s stop rule; oracle-semantics tickets only if the re-measure stays
above.

## Not yet specified

- Oracle-semantics levers — gate-verdict caching keyed on the pinned subject,
  and `-count=1` freshness (parked in
  `decisions/cost-follows-project-size.md`) — graduate into tickets here only
  if the post-stage-2 re-measurement leaves the gate above 60 s (#5's ruling).
- Reviving the outer conformance/contract width cap — dormant per
  `gate-concurrency.md`; the oversubscription #1 measured may shrink to
  nothing once nesting is gone, so re-check only after stage 2.

## Out of scope

- Diff-scoped gating in any form — ruled unsound; the ruling stands.
- Weakening or dropping any check to buy wall-clock — green keeps meaning the
  same thing; scoping a fixture to the check its EXPECT names is not
  weakening (fifth-arm precedent, extended by #6/#7).
- Cross-language incrementality — separate later capability behind its
  existing revive trigger, shaped against regroup-app.
- Byte-reproducibility tiering — decided and shipped (seventh arm); dev
  opt-in and ship-tier hermeticity stay as ruled.
- Reopening `decisions/gate-pipeline.md` — ruled closed (#4); its canary
  clause lives here now.
