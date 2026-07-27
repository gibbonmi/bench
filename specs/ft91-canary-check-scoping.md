# ft91-canary-check-scoping

Status: implemented

Source map: `decisions/gate-pipeline.md` (#5, #6; Handoff slice A). Inventory
evidence: `decisions/assets/gate-pipeline-fixture-inventory.md`.

## Problem

Inside the gate's conformance phase, every canary fixture's inner run executes
every check its tier runs (16 at dev, all 17 at ship) to observe the one check
its family or `CHECK` file targets. Cheap per check at dev, this is ruinous at ship (the rejected ~12 min
sweep estimate), and N fixtures grading one check plus the shared vacuity
baseline pay N+1 full runs. Nine fixtures also sit in a family whose check
they do not grade, so no scoping mechanism can be honest until the binding is.

## Solution

Each conformance-family fixture's inner run is scoped to the single check it
grades: a registry-owned family→check table (many-to-one) provides the
default, a fixture's `CHECK` file overrides it, and the scope crosses the
process boundary as a `BENCH_CONFORMANCE_CHECK` env the conformance entry
honours. Vacuity baselines are shared per scope group instead of one full
gate. Every mismatch — unknown check, tier-mismatched scope, unmapped family —
is a loud red, never a silent full run. Green keeps meaning the same thing:
no check is weakened or dropped; only redundant re-execution is removed.

## User stories

1. As the kit maintainer, I want a registry-owned family→check table
   (many-to-one, `compliance-hardening` → `canary-inner-compliance` by
   emitting code), with the conformance layer's family list derived from it,
   so the fixture→check binding has exactly one source.
   Line: opus (`gpt-5.6-terra`) / medium. This is gate/conformance logic at a
   settled seam, which the profile's cached routings place on the mid model at
   medium effort.

2. As the kit maintainer, I want the conformance driver to honour a check
   scope read at the entry point — executing exactly the named check, red on
   an unknown name, red on a scope the active tier does not run, and the full
   tier when no scope is set — so a scoped inner run pays for one check and a
   wrong scope can never pass as green silence.
   Line: opus (`gpt-5.6-terra`) / medium. The behavior is precisely specified
   with gate-observable outcomes, so the mid model at medium effort fits the
   cached gate-logic routing.

3. As the kit maintainer, I want the canary sweep to resolve each fixture's
   scope (CHECK override, else family table) and inject it strip-then-set into
   the inner env — behavior-owned and legacy flat fixtures staying unscoped,
   and an operator's inherited `BENCH_CONFORMANCE_CHECK` never leaking
   through — so fixtures stop paying 16 checks without any change to what a
   bite proves.
   Line: opus (`gpt-5.6-terra`) / medium. The env-hygiene pattern already
   exists in `innerEnv` and the seam is the injected Runner, so the mid model
   at medium effort is enough.

4. As the kit maintainer, I want the gate to fail loudly when one of the kit's
   own conformance families is missing from the family→check table, so adding a
   family without binding it is an error, not a silent cost regression.
   Line: opus (`gpt-5.6-terra`) / medium. A small, precisely specified fail
   posture at the same seam as story 3 keeps the mid model at medium effort.
   **Amended during implementation (reviewer veto surface):** this story
   originally placed the error in the sweep, during fixture selection. That
   seam is wrong — `internal/canary` sweeps every adopting repo, and `bench
   init` scaffolds a seed canary family a kit-owned table can never bind, so
   the sweep-side error reddened every linked repo on its first sweep. The
   sweep now resolves an unbound family to no scope (a full inner run, exactly
   today's behavior), and the conformance layer's family check raises the red
   against the kit root alone, where the table and the tree it describes are
   both in scope.

5. As the kit maintainer, I want one shared empty-tree vacuity baseline per
   scope group present in the tier's sweep — each EXPECT vacuity-checked
   against its own group's baseline, unscoped fixtures keeping today's full
   baseline — so N same-check fixtures pay N scoped runs plus one shared
   baseline while the did-not-bite comparison stays run-shape-consistent.
   Line: opus (`gpt-5.6-terra`) / medium. The grouping rule is fully settled
   in the implementation decisions (group by resolved check name, unscoped
   fixtures keep the single full baseline), so the story matches the cached
   mid-model gate-logic routing with no residual uncertainty.

6. As the kit maintainer, I want the seven stray fixtures to carry `CHECK`
   files naming their real checks (three `bounds-policy`, one each
   `default-branch-single-source`, `marker-wait-deadlines`, `skip-ownership`,
   `subcommand-routing`), so the family binding stays honest without moving
   directories.
   Line: sonnet (`gpt-5.6-luna`) / low. Seven one-line tracked files whose
   correctness the existing `fixtureTier` validation already grades makes this
   mechanical work for the cheap model.

7. As the release maintainer, I want the ship sweep scoped the same way, so
   `bench prep-release`'s canary step pays two scoped fixture runs plus one
   shared scoped baseline instead of three full inner gates.
   Line: opus (`gpt-5.6-terra`) / medium. Ship reuses the dev mechanism
   through the existing `SweepTier` path, so the mid model at medium effort
   covers the tier-specific assertions.

## Implementation decisions

- The family→check table and a `BENCH_CONFORMANCE_CHECK` env constant live in
  `internal/conformance/registry`, which stays imported-by (never importing)
  the conformance and canary layers. Many-to-one is legal; the mapping follows
  the emitting code, per the inventory asset.
- The conformance layer's family-presence check derives its family list from
  the registry table — the current literal list is deleted (one source per
  fact).
- The conformance entry point reads the scope env (as it already reads the
  tier env) and passes it into the driver; the driver owns all three fail
  postures — unknown name, tier mismatch, absent-scope-runs-full-tier — so
  posture has one owner. Diagnostics are ordinary check diagnostics: the run
  goes red through the existing path, no new output format.
- The canary widens its one `CHECK` read (today `fixtureTier`) to yield both
  tier and check name, so tier and scope cannot disagree about what the file
  says. Scope resolution: `CHECK` name if present; else the family table;
  behavior-owned and legacy flat fixtures resolve to no scope. A conformance
  family absent from the table resolves to no scope too — the sweep grades
  adopting repos whose families a kit-owned table never carries — and the
  conformance layer's family check reads the kit's own tree and reds on any
  family the table leaves unbound.
- `innerEnv` strips the inherited scope env (strip-then-set discipline);
  `runFixture` appends the fixture's resolved scope beside the existing
  per-fixture phase env.
- Vacuity baselines group by the resolved check name alone — the phase env is
  **not** part of the key. Every unscoped fixture (behavior-owned, legacy
  flat) shares one full baseline exactly as today, preserving the existing
  asymmetry in which phase-filtered fixture runs are vacuity-checked against
  the full baseline. Each check group present in the tier's sweep gets one
  scoped empty git-init baseline, executed before that group's fixtures under
  the existing worker budget (`fixtureWorkers` unchanged, no new bounds).
- No CLI surface, exit code, or output-format change anywhere. `Sweep`,
  `SweepShip`, `SweepTier` signatures unchanged.

## Testing decisions

- A good test drives the public seam and asserts observable output — env
  content on `RunCall`, diagnostics returned, timing lines written — never
  scheduler internals.
- Gate command: `bench gate` (the kit's own suites run inside its
  conformance phase, so every new test is gate-observed).
- Prior art: `internal/conformance/tier_test.go` (driver + timing-line
  assertions), `internal/canary/canary_test.go` and
  `canary_concurrency_test.go` (fake-`Runner` sweeps),
  `fixture_tier_test.go` (CHECK-file postures).

### Seam diagram

Seam 1 — the conformance driver (in-process):

    trigger: gate conformance phase → TestRootConformance (entry reads tier + scope envs)
        │
        ▼
    root, kitRoot, tier, scope ──▶ [ RunConformance: registry loop ] ──▶ diags []string
                                   [ + scope filter, fail postures ] ──▶ timing lines (git dir)
                      ◀ tests attach here: call the driver with a scope; assert
                        diags content and which checks the timing file records

Seam 2 — the canary sweep behind the injected Runner:

    trigger: gate canary phase → Sweep/SweepTier(root, tier, runner)
        │
        ▼
    fixtures (families, CHECK) ──▶ [ scope resolution → grouping →   ] ──▶ RunCall{Cwd, Env, FixtureDir}
    registry family→check table ──▶ [ baselines → workers → vacuity ] ──▶ sweep error strings
                      ◀ tests attach here: fake Runner records every RunCall;
                        assert env pairs, baseline count per group, error text

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | every family-table value names a registered check | registry unit test | new `TestFamilyCheckTableBindsRegisteredChecks` red before the table exists | a table row naming a phantom check would scope fixtures to a check that never runs, reporting did-not-bite forever |
| 1 | dropping a family from the table cannot silently shrink the derived family list | seam 1 | covered by story 4's row — the on-disk family becomes unbound and the conformance family check reds | the disk tree is the independent anchor: the family-presence check iterates the derived list, so only a check that reads the tree can red a table omission |
| 1, 3 | a wrong-but-registered binding, or a conformance fixture left unscoped, cannot survive | gate canary phase (live sweep) | already covered — a mis-scoped fixture's inner run never emits its EXPECT and the sweep reports did-not-bite | scoping a fixture away from its emitter turns the strongest existing signal (every fixture must bite) into the binding's own enforcement, fixture by fixture |
| 2 | scoped run executes exactly the named check | seam 1 | new driver test asserting the timing file records only the named check — red today (driver has no scope parameter) | an ignored filter leaves all 16 checks in the timing record |
| 2 | unknown scope name → red diagnostic | seam 1 | new driver test — red today | silent fallback would re-pay the full run and hide binding drift |
| 2 | scope naming a check the tier does not run → red diagnostic | seam 1 | new driver test — red today | a zero-check scoped run would otherwise read as green silence and the fixture as a baffling did-not-bite |
| 2 | absent scope runs the full tier | seam 1 | already covered — existing tier tests assert dev runs every dev check with no scope set | scoping regressions on the default path break the existing assertions |
| 3 | scoped fixture's `RunCall.Env` carries exactly one scope var, inherited outer value stripped | seam 2 | new fake-Runner test with `t.Setenv` outer scope — red today | a leaked or duplicated env pair hands ambient exports control of what the sweep grades (no duplicate-key precedence in exec env) |
| 3 | behavior-owned and legacy flat fixtures carry no scope var | seam 2 | not TDD-able as a pre-red — green today by construction; lands with story 3's tests as a regression guard on the new scoping path | wrongly scoping a contract fixture would skip the phase its EXPECT needs and kill the bite |
| 3 | every selected fixture yields exactly one inner run carrying its `FixtureDir` | seam 2 | not TDD-able as a pre-red — green today by construction; the fake-Runner scoping test enumerates the `RunCall`s per fixture as a regression guard | an implementation that merged or skipped same-check fixtures would grade fewer mutated trees than were selected — the map's rejected shortcut |
| 4 | the kit's own unbound conformance family → red from the conformance family check | seam 1 | new driver test with an unbound family in a synthetic kit tree — red today | an unbound family silently running unscoped re-pays the cost this slice removes and hides the missing binding; reading the tree is the only direction that sees a table omission |
| 4 | an unbound family still sweeps, unscoped, so an adopting repo is never reddened by the kit's table | seam 2 | new fake-Runner test asserting the sweep succeeds and injects no scope var — red today | `bench init` scaffolds a seed canary family a kit-owned table can never bind, so a sweep-side error breaks every linked repo on first run |
| 5 | one baseline per check group present, not one per fixture and not one full gate | seam 2 | new fake-Runner test counting empty-`FixtureDir` calls per check group — red today | the N+1 cost and the run-shape mismatch both reappear if grouping is wrong |
| 5 | vacuity graded against the fixture's own group baseline | seam 2 | new fake-Runner test: EXPECT matching its own group's baseline output → vacuous error; matching only another group's → no error — red today | cross-group comparison either false-flags scoped EXPECTs or lets a genuinely vacuous EXPECT pass |
| 5 | an all-unscoped sweep runs exactly one full baseline | seam 2 | not TDD-able as a pre-red — green today by construction; lands with story 5's tests as a regression guard proving the check-name grouping key never splits the unscoped group | phase-keyed grouping would split contract fixtures from flat ones into two baselines, changing today's semantics and doubling flat-sweep cost |
| 6 | the seven stray CHECK files name registry checks the fixtures' EXPECTs are emitted by | seam 2 + gate canary phase | already covered — `fixtureTier`'s unknown-name error grades the files; the live sweep proves the bites survive scoping | a misspelled CHECK errors the sweep; a wrong-but-registered CHECK scopes the run away from the EXPECT's emitter and the fixture reports did-not-bite |
| 7 | ship sweep = two scoped fixture runs + one shared scoped baseline | seam 2 | new fake-Runner `SweepTier(Ship)` test — red today | ship's cost reduction is the slice's headline; without the assertion a regression to full inner gates would land silently |

### Edge inventory

Walked against the profile's hostile-input checklist:

- absent vs empty `CHECK`: absent → family table; empty → existing error —
  already covered (`fixtureTier` tests), unchanged by the widened read.
- `CHECK` without trailing newline: existing `TrimSpace` read — already
  covered.
- hostile scope env value (control bytes, renamed check): lands in the
  unknown-name red (story 2 row); no render surface is added.
- duplicate env keys: strip-then-set asserted in the story 3 row.
- family with zero fixtures: existing family-presence check — already
  covered.
- interrupt mid-sweep / temp cleanup: existing temp-cleanup and overlap tests
  unchanged — already covered.
- fixture counts 0/1 and worker clamps: `fixtureWorkers` untouched — already
  covered by existing concurrency tests.
- **Won't handle:** a FIFO/device planted as a `CHECK` file — fixtures are
  kit-authored tracked trees, not user input, and the read surface predates
  this spec.
- **Won't handle:** scoping the outer conformance/contract phases — dormant
  per the map unless contention flakes persist.

## Out of scope

- **Slice B — manifest + DAG runner** (`.bench/phases.json`, scheduler,
  straggler naming): a separate capability with its own map decisions;
  ~40 edits, ~7 gate runs.
- **Slice C — `checkGoCore` split + fixture migration + parity test**:
  depends on B's phases and this slice's scoping; ~50 edits, ~9 gate runs.
- **Cross-language capability-skip surface** (`bench skip` helper or
  documented log contract): recorded fog on the map; ~15 edits, ~3 gate runs.
- The two FT91 interim defects (`BENCH_CONFORMANCE_TIER` scrub symmetry,
  probe-output spill) stay on the roadmap row; slice C absorbs their
  surfaces.
