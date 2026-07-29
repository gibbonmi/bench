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
`assets/gate-critical-path-timeline.md`. The gate is canary-bound end
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
Resolved 2026-07-29; full inventory with per-claim citations in
`assets/artifact-test-inventory.md`. Load-bearing facts: the
prepared seam copies rather than consumes (`build-artifacts.sh:110` `cp -a`),
so a prepared set is read-only input; no test in the package calls
`t.Parallel`, so sharing faces ordering hazards only; `TestMain` puts the
package on the shared-cache posture and posture tests strip it per row; the
hermetic default is a double generation; the missing-pin refusal fires after
the build matrix, so that test pays a full generation before refusing. Up to
7 of ~17 generation passes are host-only builds of the identical recipe and
collapse to one shared set; non-collapsible by subject: 4 hermetic
double-build passes, 3 GOPROXY-off passes, 2 stale-record passes, and
`TestDistributableArtifactContracts`' unique host+non-host shape.

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
Resolved 2026-07-29 (reviewer). **Who shares:** the inspection/consumption
tests (`TestArtifactSourceStagesCommittedHostPlan`,
`TestSharedCacheBuildPromotesNoRecord`, `TestOfflineArchiveProjection`,
`TestPackedArtifactRunsSetupOfflineFromASpacedPrefix`) plus the two
prepared-seam promotion tests (`TestArtifactPromotionIsAtomicAndExclusive`,
`TestSharedCacheBuildRestoresRecordOnInterruptedPromotion`) share one
package-scoped host-only set. Posture tests (hermetic, GOPROXY-off,
stale-record), refusal tests, and `TestDistributableArtifactContracts`'
host+non-host build stay private — their build is the subject.
`TestArtifactBuilderRefusesMissingBinaryPinManifest` joins only if the seam
still reaches the pin check; otherwise it stays private. **Seam shape:** a
lazy in-package singleton — first requester stages the host-narrowed
committed clone and runs one real shared-cache build, exposing the staged
source and promoted output read-only; prepared-seam tests use the output as
their prepared input; nothing builds when no sharer is selected; fresh per
`go test` invocation, so no cross-run staleness and no contact with the
parked verdict-caching levers. **Fail posture:** both belts — the set is
chmod'd read-only after generation so a mutator fails at write time with
attribution, and a digest fingerprint recorded at build is re-verified before
each consumer use, failing loud as "shared artifact set mutated" (the digest
helper exists: `promotedArtifactDigests`).

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

Rewritten 2026-07-29 for the artifact-hoist slice after stages 1 and 2
shipped (their decisions stay recorded in #6 and #7); #2/#3 are closed and
this is the map's last open slice before the #5 re-measurement.

1. **Module boundaries.** The whole change is test-harness work inside
   `internal/contract/surface/artifact`: a package-scoped lazy singleton
   helper owning the staged host-narrowed clone and its one real build, plus
   the sharer migrations. `scripts/build-artifacts.sh` untouched — the
   existing prepared seam (`cp -a`, copy-not-consume) already serves.
   Production packages untouched.
2. **Contracts.** The helper: first requester stages one host-narrowed
   committed clone and runs one real shared-cache build; it exposes the
   staged source root, the promoted output dir, and the build fingerprint,
   all read-only; later requesters reuse; no requester selected → no build.
   The six ruled sharers (#3) consume it; a verify-before-use check reds with
   a failure naming the shared set as mutated, never a downstream assertion.
3. **Deep vs thin.** The singleton (staging, build, chmod, fingerprint,
   verify) is the deep unit; per-test verify-and-consume calls are thin.
4. **Black-box assertables.** All six sharers pass against the shared set; a
   deliberate probe that writes into the shared set produces the named
   mutation red, not a downstream failure; a `-run` selecting only
   fabricated-fixture tests invokes no build; suite time drop recorded in
   `assets/gate-critical-path-timeline.md`.
5. **Gate attachment.** The gate's contract phase runs the package as today;
   no new gate check. Wall-clock is not gate-assertable — ship evidence is
   the post-change measurement, and #5's ≤60 s rule is judged on the full
   re-measure after this slice lands.
6. **Hostile-input owners.** A sharer mutating the set → write-time chmod red
   or the fingerprint verify, attributed to the set. A chmod or fingerprint
   failure at build → the helper, red, fail closed. A test needing a
   different plan shape → private build by ruling, never a widened shared
   set. The pin-manifest refusal test → private unless the spec shows the
   seam reaches the pin check.
7. **Uncertainty flags.** Only the `TestArtifactBuilderRefusesMissingBinaryPinManifest`
   membership: the prepared seam today skips packing, so the pin check is
   unreachable through it; joining would need a narrower prepared-binaries
   seam in the production script, which crosses boundary 1 — recommend it
   stays private, spec decides with that constraint stated.
8. **Rejected alternatives.** Seam-tests-only minimal sharing; a host+non-host
   shared shape folding `TestDistributableArtifactContracts` (drifts every
   host-only consumer's assertions); a `TestMain` up-front build (pays on
   every `-run`); an on-disk tree-hash cache (cross-run staleness — the
   parked oracle-semantics territory); convention-only mutation posture.
9. **Domain watch-outs.** The prepared seam copies (`cp -a`), so a prepared
   set is reusable input. No test in the package calls `t.Parallel`;
   introducing it changes the hazard analysis from ordering to races.
   `TestMain` sets the shared-cache posture package-wide — the shared build
   inherits it, and posture tests strip it and must never share. The five
   `surface/artifact` canary fixtures each pay the whole suite, so the
   per-suite saving multiplies about six-fold.

Dependency order: single spec — the hoist — then the full-gate re-measure
against the ≤60 s stop rule (#5); oracle-semantics tickets only if the
re-measure stays above.

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
