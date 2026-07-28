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
Resolved 2026-07-28 (reviewer): the nesting premise is formally narrowed —
a behavior-owned fixture's bite is proven at its owning contract test, not by
a nested gate. The sweep stops spawning nested gates for this family; the
owning package's test binary runs once per package group with each fixture
root as a subtest (same test code, same `BENCH_CONTRACT_ROOT` seam, one
process instead of a process tree). Bite and vacuity semantics unchanged:
EXPECT observed as a real failure per fixture, EXPECT absent on the unmutated
baseline under the same run shape. The gate-level plumbing (red test → red
phase → red gate) is asserted once by the gate's own phase tests, not 34
times. Stage 1's subfamily binding is stage 2's input — nothing is throwaway.
The migration must verify per fixture that every EXPECT is observable in
test-level output rather than gate-level framing (the slice-C migration
casualty class). Expected canary → ~20–40 s; the inner-width pin becomes
moot for this family. Transport details (env contract for the fixture list,
baseline invocation shape) are spec-time within this decided shape.

## Handoff

Written 2026-07-28 for the two canary slices (#6, #7), same-session with the
reviewer's rulings; #2/#3 stay open and gate only the artifact-hoist slice.

1. **Module boundaries.** `internal/canary` owns the fixture walk extension
   (package subpaths under `behavior-owned`, fixture = directory holding
   EXPECT), contract scope resolution, per-package vacuity baseline groups,
   and the per-fixture inner env. `internal/gate` owns the inner-mode
   contract-phase argv narrowing and the outer-mode scrub of the new env.
   Stage 2 adds the fixture-parameterized bite mode to the owning
   `internal/contract` packages. `tests/canary/behavior-owned/<package
   path>/<fixture>/` is the binding — data, no second file.
   `internal/bounds` untouched.
2. **Contracts.** One new env var carrying the bound package path relative to
   `internal/contract` (scrub-then-set everywhere; exact name is spec-time).
   Absent → today's full `./internal/contract/...` argv, so adopting repos
   and non-contract fixtures are untouched. Set to a package the kit tree
   lacks → red naming it. A behavior-owned fixture directly under the family
   (no package dir) → red, missing binding. EXPECT, bite, did-not-bite, and
   vacuity semantics unchanged; gate exit codes unchanged.
3. **Deep vs thin.** The sweep (walk, scope resolution, baseline grouping)
   stays the deep unit behind the injected `Runner` seam; `scopedEnv`-style
   plumbing and the subfamily convention are thin.
4. **Black-box assertables.** Via a fake `Runner`: each behavior-owned
   fixture's `RunCall.Env` carries phase=contract plus its package env;
   baseline run count equals the scope groups present; unbound fixture and
   unknown package each produce their red. Via the gate's phase-table tests:
   inner-mode argv narrows to the one package; outer mode strips an ambient
   export of the env.
5. **Gate attachment.** The canary sweep itself enforces bite and vacuity
   per fixture; the kit's own suites run in the gate's test and
   conformance-suite phases. The wall-clock outcome is not gate-assertable —
   ship evidence is the post-change measurement against
   `assets/gate-critical-path-timeline.md`.
6. **Hostile-input owners.** Garbage package values (traversal, absolute,
   empty) → the scope resolver / phase-table reader, red. Ambient export of
   the new env → `innerEnv` strip list plus the outer-mode scrub. Empty or
   missing binding directory → the walk, red (empty-CHECK precedent).
7. **Uncertainty flags.** None blocking stage 1. Stage 2's transport (fixture
   list env, baseline invocation) is spec-time within #7's decided shape.
8. **Rejected alternatives.** Per-fixture binding file; registry
   fixture→package table; test-level `-run` binding at stage 1; keeping
   full-suite nesting; any silent fallback for a kit-owned unbound fixture.
9. **Domain watch-outs.** Go's exec env has no duplicate-key precedence —
   every inner pin strips-then-sets. Fixture base names are globally unique
   across all families; subfamily moves must not collide them. The empty-tree
   vacuity baseline runs whatever phases its env selects, so a scoped group's
   baseline must carry exactly its fixtures' phase+package env or the
   scoped-vs-scoped comparison silently breaks. The walk today treats any
   directory under a family as a fixture; EXPECT presence is the only safe
   discriminator once package segments nest.

Dependency order: stage 1 (#6, package-scoped nesting) → stage 2 (#7,
in-process bites) → artifact hoist (pends #2/#3) → re-measure against the
≤60 s stop rule; oracle-semantics tickets only if the re-measure stays above.

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
