# Gate critical path (FT91, eighth arm)

## Destination

Find where the dev gate's remaining wall-clock actually lives, decide the lever
that takes it — or decide FT91 is done — and settle whether
`decisions/gate-pipeline.md` reopens now that slice C falsified its premise.
Green keeps meaning the same thing throughout. The stop condition for FT91 is
itself decided here (#5), against evidence rather than picked blind.

## #1: Why did the gate absorb only 24 s of a 131 s suite win — what is the critical path now?

Type: Research

### Question
The seventh arm cut the contract suites 249 s → 118 s measured solo, but the
whole gate moved only 4m51s → 4m27s: under gate parallelism the artifact suite
inflates ~106 s → ~152 s, and the contract phase is no longer clearly the
critical path. Nobody has diagnosed the gap. Run the gate with its permanent
per-phase timing, capture a phase timeline (start/stop overlap, not just
durations) plus load, and name: the current critical path, the mechanism of the
artifact-suite inflation (CPU oversubscription against canary, ambient
build-cache contention, disk, the output-directory lock), and the realistic
wall-clock floor if the contract phase were free. Deliverable: a numbers asset
with per-claim citations; #4 and #5 grill against it.

### Answer
— (open)

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

### Question
The pipeline map closed on the premise that conformance was the long pole;
slice C's measurement falsified it (`package-core-guard` 1m52.8s → 3.3s,
whole gate unchanged). One ruling covers both faces (reviewer-folded
2026-07-28): whether the map's remaining fog — the kit-owned `.bench/phases.json`
(story 9, unbuilt, its acceptance row unsatisfiable as specced) and the probed
phases that shipped instead of it (stories 4, 5) — is still worth building now
that the premise is gone, and the retirement of the deliberately-unretired
spec rides the same answer.

### Answer
— (open)

## #5: What ends FT91 — and what, if anything, does the eighth arm build?

Blocked by: #1, #3, #4
Type: Grill

### Question
The closing ticket. Against #1's numbers: a target wall-clock, a
diminishing-returns rule, or close the row now. And the arm selection: the
prepared-artifact hoist (#3), whatever mechanism fix #1 surfaces, both, or
nothing. Decided here so the stop condition is evidence-based; seven arms in,
the last one bought 24 s.

### Answer
— (open)

## Not yet specified

- Reviving the outer conformance/contract width cap — dormant per
  `gate-concurrency.md`; #1's diagnosis may fire its contention trigger.
- Gate-verdict caching keyed on the pinned subject — parked behind
  re-measurement; #1 is that re-measurement, graduates only if the diagnosis
  makes it the lever.

## Out of scope

- Diff-scoped gating in any form — ruled unsound; the ruling stands.
- Weakening or dropping any check to buy wall-clock — green keeps meaning the
  same thing.
- `-count=1` freshness semantics — parked in
  `decisions/cost-follows-project-size.md`, reviewer-led, not this map's.
- Cross-language incrementality — separate later capability behind its
  existing revive trigger, shaped against regroup-app.
- Byte-reproducibility tiering — decided and shipped (seventh arm); dev opt-in
  and ship-tier hermeticity stay as ruled.
