# Gate pipeline (FT91, pipeline arm)

## Destination

The dev gate becomes a local pipeline of first-class phases with declared
dependencies: `checkGoCore`'s seven serial toolchain steps overlap instead of
running inside one test function, conformance grades structure only, and every
canary fixture pays for only the check its `CHECK` file names. Phase
definitions become project-owned data (the runner ships in the kit, shaped
against regroup-app), cutting dev wall-clock — today 4m36s, long pole
`package-core-guard` ~86 s — without changing what green means.

## #1: What does the phase manifest declare, and who owns its failure modes?

Type: Grill

### Question
The phase list must move from Go code (`BenchkitPhases`) to project-owned
declarative data. Open: the schema (name, command, dependencies, tier,
timeout, serial-vs-concurrent, capability requirements — which of these are
v1 fields and which stay runner defaults); the format and location (a file
beside `.bench/gate.sh`, and whether `gate.sh` stays the entry that invokes
the runner or becomes a generated shim); and the fail posture for the
absent / empty / malformed manifest family — absent presumably falls back to
the built-in kit table, but empty-versus-malformed is exactly the class this
repo keeps meeting and must be decided, not defaulted.

### Answer
— (open)

## #2: How much of `internal/gate`'s existing runner survives as the pipeline runner?

Type: Research

### Question
The runner already owns a `Phase` struct, a serial/concurrent split
(`splitSerialPhases`), capability skips, per-phase timing, and subject
binding — the pipeline may be an extension (add dependency edges and manifest
loading) rather than a rewrite. Read `internal/gate` (`phases.go`,
`runner.go`, `capability_skips.go`, `subject.go`) and map: which parts already
implement the phase/verdict/timing model the FT91 row calls universal, what
assumes the phase table is compiled-in, what assumes exactly one
serial-then-concurrent stage shape, and where dependency edges would attach.
Asset: a short summary with per-claim file citations, so #3 and #4 grill
against facts rather than recollection.

### Answer
— (open)

## #3: Where exactly does `checkGoCore` split?

Blocked by: #2
Type: Grill

### Question
Which of the composite's steps become first-class phases
(gofmt / build / vet / test / race / cross-compile — the last is ship-tier),
and what the residual conformance check still asserts once the toolchain work
leaves it. The complication that stays: `checkGoCore` grades arbitrary roots —
canary fixtures and linked repos — so the split must yield a host-repo phase
set plus a narrower structural check that still works against any root. What
that narrower check owns is the decision; the 19+ package-core-guard fixtures
asserting the composite's diagnostics hang off it (#6).

### Answer
— (open)

## #4: What are the DAG's execution semantics?

Blocked by: #1
Type: Grill

### Question
Width budget: the gate-concurrency arm settled a product-budget model
(`runtime.GOMAXPROCS(0)`, no knob, inner gates pinned) — do overlapping
toolchain phases share that budget by declared weight, or run unweighted and
rely on the OS, and how do nested canary inner gates inherit it? Failure
posture: fail-fast on first red phase versus run-all-and-report (today's
phases already differ from a fix-loop's needs — a red that cancels siblings
saves wall-clock but destroys the shrinking-red-set signal FT107's fifth
clause wants). Timeout ownership: per-phase deadlines in the manifest versus
runner defaults, given the 600 s go-test package default has twice presented
as a gate hang.

### Answer
— (open)

## #5: How does a canary fixture's inner run scope to its named check, and what shape is the dedup?

Type: Grill

### Question
The prerequisite slice, buildable on today's gate: inside the conformance
phase every fixture still runs all 16 registry checks to observe the one its
`CHECK` file names — cheap at dev, ruinous at ship (~12 min sweep estimate,
rejected). Two levers, both needing mechanism decisions: (a) scope the inner
run to the named check — an env the registry already honours, a derived
`-run` filter, or a registry-level filter argument, and what happens when a
fixture names no check; (b) dedup by check — N fixtures grading one check
plus a shared vacuity baseline currently pay N+1 runs, and both ship fixtures
grade the same check, so ship's cost is irreducible by (a) alone. Open
whether the vacuity baseline can be shared per-check without weakening the
did-not-bite verdict.

### Answer
— (open)

## #6: What do the package-core-guard fixtures migrate to?

Blocked by: #3
Type: Research

### Question
19+ fixtures assert the composite check's diagnostics; after the split they
must grade the phase that owns each diagnostic. Inventory what each fixture's
EXPECT strings actually assert, bucket them by destination phase versus
residual structural check, and name any fixture whose expectation no seam
will emit after the split — those need a decision, not a mechanical move.
Asset: the bucketed inventory.

### Answer
— (open)

## #7: What proves parity — that restaging lost no check?

Blocked by: #3
Type: Grill

### Question
The oracle constraint: green must keep meaning the same thing. The tier split
proved membership with a registry both directions could interrogate
(`TestRegistryBindsEveryCheck`); the pipeline needs the analogue — every
sub-step of the old composite maps to exactly one phase or to the residual
check, asserted mechanically rather than by review. Decide the mechanism
(registry-style parity test, a recorded red per moved step per `craft-gate`,
or both) and what the migration's demonstrated-bite evidence is.

### Answer
— (open)

## #8: Does the manifest hold against regroup-app?

Blocked by: #1
Type: Task

### Question
The FT91 row's decided shaping rule: design the manifest against regroup-app
(the first external, non-Go validation target) rather than generalising from
this repo. The reviewer makes regroup-app's toolchain shape readable (what
its build / test / lint phases actually are); the worker then walks the draft
schema against it on paper and reports where the schema bends. No code — the
deliverable is the list of schema changes regroup-app forces.

### Answer
— (open)

## #9: What happens to the gate's timing and output contract?

Type: Grill

### Answer
Resolved at bootstrap — continuity, decided by the tier split: per-check
timing is permanent gate observability owned by the conformance driver, and
the runner keeps emitting per-phase timing in the same stable format
(ordering byte-stable, values free to vary). New overlapping phases each get
a timing line like any phase today. No output-contract redesign rides this
map; any change to the format is its own reviewer decision later.

## Not yet specified

- Ship-tier phases joining the manifest (`prep-release` as a manifest
  consumer) — dim until #1 fixes the schema.
- How linked repos receive runner upgrades once phase definitions are theirs
  (`bench upgrade` semantics for a project-owned manifest).

## Out of scope

- Cross-language incrementality (phases declare input globs, runner hashes
  them) — a small build system; deferred behind FT91's standing revive
  trigger, never built speculatively on this map.
- Removing `-count=1` and cache infrastructure — same trigger;
  oracle-semantics decisions, reviewer-led.
- Diff-scoped gating in any form — ruled unsound; the ruling stands.
- Weakening or dropping any check to buy wall-clock — green keeps meaning
  the same thing.
- Outer conformance/contract width capping — dormant unless contention
  flakes persist.
- Removing canary nesting — running the real gate against broken trees is
  its premise.
- The two interim defects on the FT91 row (`BENCH_CONFORMANCE_TIER` scrub
  symmetry, probe-output spill) — shape already decided there; build work,
  not map fog.
