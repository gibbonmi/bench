# ft91-canary-contract-scoping

Status: implemented

Map: `decisions/gate-critical-path.md` (#1, #6; Handoff written same-session
2026-07-28 under the entry contract's reviewer-closed exception). Flagged for
reviewer veto — decided here, not carried by the map: the env var name
(`BENCH_CANARY_CONTRACT_PACKAGE`); the prefixed baseline-group key; the
validation split (structure and manifest-declaring checks sweep-side,
value-shape and package-existence checks gate-side against the kit tree); the
relocation of the two non-contract-owned fixtures to legacy flat fixtures;
landing the 34-fixture migration in the same change as the machinery; the
kit-only flat-set guard test (an independently authored expectation under the
code standard's named-mutation exception) and the ≤100 s solo-canary
acceptance threshold on the ship evidence.

## Problem

The dev gate is canary-bound: canary solo is 250 s of a 267 s gate wall, and
~1,700 of its ~2,000 worker-seconds are the 34 `behavior-owned` fixtures, each
of whose nested gates re-runs the entire kit contract suite
(`go test -count=1 ./internal/contract/...`) at inner width 2 — seven packages
re-proven per fixture when each fixture's EXPECT names exactly one failure.
Evidence: `decisions/assets/gate-critical-path-timeline.md`.

## Solution

Scope each behavior-owned fixture's inner gate to the one contract package
that owns its EXPECT, the way the fifth arm scoped conformance fixtures to
their one check. The binding is the directory tree itself —
`tests/canary/behavior-owned/<package path>/<fixture>/`, where `<package
path>` is the path relative to `internal/contract` (nested for nested
packages, e.g. `surface/artifact/`) — and every failure mode is loud. The two
fixtures whose EXPECT no contract package emits relocate to the existing
legacy flat class (full inner gate). Expected canary ~250 s → ~60–80 s with
bite and vacuity semantics unchanged per fixture.
`decisions/assets/behavior-owned-package-bindings.md` is the traced binding
inventory the migration starts from.

## User stories

1. As the canary sweep, I discover behavior-owned fixtures under package
   subpaths — a directory holding an EXPECT file is a fixture, any directory
   above it (below the family) is a package path segment — and each
   `selected` fixture carries its family and package path explicitly, so no
   later code re-derives the family from the fixture's parent basename.
   Fixture base names stay globally unique across all families. Line:
   gpt-5.6-terra / medium. The walk change touches shared discovery code
   whose regressions the whole sweep inherits, and the two parent-basename
   call sites it retires are silent-misroute hazards.
2. As the canary sweep, I resolve each behavior-owned fixture's scope to its
   package path and red loudly on the structural defects I own: a fixture
   directly under `behavior-owned/` (no package directory), a package
   directory that does not exist under the kit tree's `internal/contract/`,
   and a fixture whose `files/` tree declares a phase manifest
   (`dot-bench/phases.json`) — a manifest-bearing root bypasses the built-in
   table the narrowing lives in, so such a fixture must be flat, never
   silently unscoped. Line: gpt-5.6-terra / medium. Scope resolution carries
   the fail-posture family and must match the fifth arm's loud-red precedent.
3. As a behavior-owned fixture's inner gate, my `RunCall.Env` carries exactly
   `BENCH_CANARY_PHASE=contract` and
   `BENCH_CANARY_CONTRACT_PACKAGE=<package path>` (multi-segment values
   included, e.g. `surface/artifact`), with no `BENCH_CONFORMANCE_CHECK`
   entry, and with any ambient export of either variable stripped-then-set.
   Line: gpt-5.6-terra / medium. Env plumbing is small but the
   strip-then-set watch-out has bitten this repo twice, and the exact-set
   assertion is what catches a conformance scope leaking in from the old
   family derivation.
4. As the built-in phase table (`BenchkitPhases`), when
   `BENCH_CANARY_CONTRACT_PACKAGE` is set in inner mode I narrow the contract
   phase argv from `./internal/contract/...` to
   `./internal/contract/<package path>`, and I red naming the value when it
   is empty, absolute, `..`-carrying, or names a directory the kit tree
   lacks. Absent the variable, the argv is byte-identical to today.
   Manifest-declared tables are never rewritten — story 2's sweep-side red is
   what keeps that path unreachable for scoped fixtures. Line: gpt-5.6-terra
   / medium. This is the seam where a silent mistake weakens what green
   means, and the manifest bypass is exactly the trap a builder working from
   the old draft would have fallen into.
5. As the outer gate, `gateEnv` strips `BENCH_CANARY_CONTRACT_PACKAGE` (as it
   scrubs the other canary control variables), so an operator export can
   never scope a real gate's contract phase. Line: gpt-5.6-terra / medium.
   The scrub point is the one place a removal can be expressed — `Phase.Env`
   merging can only override, never delete.
6. As the vacuity guard, each contract scope group gets its own empty-tree
   baseline run whose env carries the same
   `BENCH_CANARY_PHASE`/`BENCH_CANARY_CONTRACT_PACKAGE` pair its fixtures
   use — a deliberate change from today's phase-env-less baselines, per the
   map's scoped-vs-scoped watch-out — and each EXPECT is graded against its
   own group's baseline, with a fixture whose EXPECT appears in that baseline
   reported vacuous. The single unscoped full baseline remains for legacy
   flat and toolchain-phase fixtures. Line: gpt-5.6-terra / medium. Baseline
   grouping is the soundness argument made executable, and the group key is
   written and read in different functions — the desync silently un-guards
   vacuity, so it needs its own red.
7. As the kit maintainer, the 34 fixtures migrate per
   `decisions/assets/behavior-owned-package-bindings.md`: 32 move into their
   package subdirectories, and the two whose EXPECT no contract package emits
   (`phase-manifest-defect-admitted`, `native-proof-digest-binding-bypassed`)
   relocate to legacy flat fixtures, in the same change as the machinery, so
   the sweep itself (bite + vacuity, both enforced red) proves every binding.
   The one ambiguous binding (`roadmap-regressed`) and the one script-emitted
   EXPECT (`publication-unpublish-attempt`) are verified by bite before the
   change is called done. A kit-only guard test pins the flat set to exactly
   the two named fixtures, so flattening is never an escape from scoping —
   its independent expectation is justified by the named mutation "relocate a
   package-bound fixture to flat", which must red. Line: gpt-5.6-terra /
   medium. The inventory is
   traced but a wrong binding surfaces only as a did-not-bite red at gate
   time, so verification is part of the story, not an afterthought.
8. As the reviewer, I get ship evidence: a post-change instrumented gate and
   solo canary measurement against the timeline asset's baselines (267 s /
   250 s), recorded in that asset — and the story is satisfied only if solo
   canary lands at or below 100 s, so an unchanged measurement (the signature
   of a degenerate migration) stops the build and reports instead of shipping.
   Line: gpt-5.6-luna / low. Running and recording a measurement is
   mechanical once the build is green; the threshold is what makes the
   evidence an acceptance, not a note.

## Implementation decisions

- The binding is the directory tree, no second source: `internal/canary`'s
  walk learns that under `behavior-owned` a fixture is the directory holding
  EXPECT and everything between family and fixture is the package path. The
  `selected` record carries family and package explicitly; the two call
  sites that today re-derive family via `filepath.Base(filepath.Dir(fx))`
  (the phase-env pin in `runFixture` and the scope lookup in `fixtureScope`)
  are retired onto it. The existing global-uniqueness check on fixture base
  names is kept as-is.
- Contract scope groups key the baseline map as `contract:<package path>`,
  so a package path can never collide with a conformance check name; the
  same `selected` field feeds the write and the read so the key cannot
  desync by construction, and story 6's vacuity row proves it.
- Transport is `BENCH_CANARY_CONTRACT_PACKAGE`, owned by `internal/canary`:
  added to `innerEnv`'s strip list, set per fixture and per baseline
  scrub-then-set. `internal/gate` reads it only in `BenchkitPhases` under
  inner mode; `gateEnv` strips it in the outer gate. Manifest-parsed tables
  never see it (story 2 makes that path a sweep-side red for this family).
- Validation is split by owner: the sweep's scope resolver reds structure
  (unbound fixture, nonexistent package under the kit tree,
  manifest-declaring fixture) naming the fixture; the phase-table reader
  reds value shape (empty, absolute, `..`, nonexistent under the kit tree it
  builds argv for) naming the value, as backstop.
- Bite, did-not-bite, vacuity, exit codes, and every non-behavior-owned
  fixture's behavior are unchanged. Baseline env for contract groups is the
  one deliberate behavior change, stated in story 6. The
  `bounds.CanaryInnerWidth` pin is untouched (moot at stage 2, not here).
- The migration lands in the same change as the machinery — this spec's own
  call, flagged above: the sweep's enforcement is what proves the bindings,
  so splitting the change would ship machinery no fixture exercises.

## Testing decisions

- A good test drives the sweep or the phase table through its public seam and
  asserts observable outputs — `RunCall` contents via a fake `Runner`, argv
  via the phase-table constructor, reds via returned errors — never internals.
- Seams: the injected `Runner` in `internal/canary` (prior art:
  `TestSweepBoundsFixtureConcurrencyAtNumCPU` and the fifth arm's scoping
  tests), the phase-table functions in `internal/gate` (prior art: the
  inner-mode tests around `phasesForMode`), and `gateEnv` (prior art: the
  slice-B env-scrub symmetry tests).
- Gate: `bench gate` (dev tier) must be green; the canary phase within it is
  the end-to-end proof that all migrated fixtures still bite.

### Seam diagram

Sweep seam (stories 1–3, 6, 7):

    trigger: gate canary phase → Sweep(root, runner)
        │
        ▼
    tests/canary tree ──▶ [ walk + scope resolution ] ──▶ selected{dir, family, package}
                               │
                               ▼
    scope groups      ──▶ [ scopeBaselines ]  ──▶ baseline output per contract:<pkg> group
    selected fixtures ──▶ [ runFixtures    ]  ──▶ RunCall{Cwd, Env} per fixture ──▶ errs
                      ◀ tests attach here: fake Runner records every RunCall;
                        asserts exact env sets, baseline count and env, red messages

Phase-table seam (story 4) and scrub seam (story 5):

    trigger: bench gate-phases (inner mode, via .bench/gate.sh in the fixture tree)
        │
        ▼
    root, kit ──▶ [ BenchkitPhases (reads scope env, inner mode) ] ──▶ []Phase (argv)
              ◀ tests attach here: build the table with the env set/unset/garbage;
                assert contract argv equality and reds

    trigger: bench gate (outer mode)
        │
        ▼
    os.Environ ──▶ [ gateEnv ] ──▶ phase env without BENCH_CANARY_CONTRACT_PACKAGE
               ◀ tests attach here: export the var, assert it is stripped

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | fixture under `behavior-owned/surface/artifact/x/` discovered once as fixture `x` with package `surface/artifact` | sweep walk via fake Runner | new test red before the walk learns package segments | today's walk treats `surface` as the fixture and never finds `x`; the package assertion also kills the cheapest wrong walk that keeps only the last segment |
| 1 | name collision across subfamilies still errors | sweep walk | already covered (`fixtures()` uniqueness test), extended red for the nested case | the extended case would pass silently if uniqueness ignored nested dirs |
| 2 | fixture directly under `behavior-owned/` reds naming it | sweep via fake Runner | new test red while the walk treats it as an unscoped full run | an unbound fixture silently paying the full suite is the cost this spec removes |
| 2 | package dir absent under kit's `internal/contract/` reds naming fixture and package | sweep via fake Runner | new test red before the resolver checks the kit tree | a rotted binding would otherwise surface only as a did-not-bite red naming the wrong defect |
| 2 | fixture whose `files/` declares `dot-bench/phases.json` reds telling it to be flat | sweep via fake Runner | new test red before the resolver stats the manifest | a manifest-bearing root bypasses `BenchkitPhases`, so its inner run would silently skip the narrowing |
| 3 | behavior-owned `RunCall.Env` carries exactly the phase pin and `BENCH_CANARY_CONTRACT_PACKAGE=surface/artifact`, no conformance-check entry, ambient exports stripped | fake Runner exact-env assertion | new test red before env plumbing exists | an exact-set assertion with a multi-segment value kills both the single-segment shortcut and a conformance scope leaking from the old family derivation |
| 4 | inner-mode contract argv narrows to `./internal/contract/surface/artifact` | phase-table test | new test red before the reader exists | argv stays `./internal/contract/...`, failing the equality assertion; the nested value kills the last-segment shortcut |
| 4 | empty, absolute, `..`, and unknown-dir values each red naming the value | phase-table test | new test red per value class | a garbage value that silently ran the full suite would hide binding rot |
| 4 | absent variable → argv byte-identical to today | phase-table test | new test red if the reader changes the unset path | adopting repos and unscoped fixtures must be untouched |
| 5 | `gateEnv` strips an exported `BENCH_CANARY_CONTRACT_PACKAGE` | gateEnv test | new test red before the strip lands | an operator export scoping a real gate would weaken what green means; the argv-side test alone cannot catch it because outer tables never read the var |
| 6 | fixtures across two distinct packages produce two baseline calls, independently keyed `contract:<pkg>`, each env matching its own group's fixtures | fake Runner call recording, two-package case | new test red before grouping exists | pinning two groups in one test kills the degenerate single-shared-baseline: call count, key set, and per-call env all fail if one baseline serves every key |
| 6 | with two groups whose baselines differ, a fixture reds as vacuous only against its own group's output | fake Runner: group B's baseline emits group B's EXPECT, group A's does not | new test red if any fixture is graded against another group's baseline or an empty lookup | a write/read key desync yields an empty baseline, `strings.Contains("", x)` is always false, and vacuity is silently un-guarded; the cross-group case also catches grading against the wrong group |
| 7 | all migrated fixtures still bite and none is vacuous | gate canary phase | already covered — the sweep reds on did-not-bite and vacuity by construction | a wrong binding makes the bound package not emit the EXPECT, a did-not-bite red at the next gate |
| 7 | the flat-fixture set is exactly the two named relocations | kit-only guard test over the tests/canary tree | new test red, demonstrated by temporarily flattening a package-bound fixture | flattening every fixture is the cheapest wrong migration — full runs keep bite and vacuity green while removing the entire win; this guard reds the first fixture that escapes scoping |
| 8 | measured post-change solo canary ≤100 s recorded in the asset | manual, thresholded | not TDD-able — wall-clock is not gate-assertable (map Handoff item 5) | an unchanged ~250 s measurement is the degenerate-migration signature; the threshold turns the evidence row into an acceptance that fails it |

### Edge inventory

- Empty scope value, absolute path, `..` traversal, unknown package → coverage
  rows (story 4); nonexistent package also sweep-side (story 2).
- Ambient export of the scope variable → coverage rows (stories 3, 5).
- Fixture tree declaring a phase manifest → coverage row (story 2).
- Scoped empty-tree baseline producing no output → **Won't handle** beyond
  existing behavior: `scopeBaselines` already reds any group whose baseline
  prints nothing, and a scoped contract run against an empty
  `BENCH_CONTRACT_ROOT` reds with output; if a group ever trips the
  no-output red, that is the correct loud outcome.
- Bound package exists but no test emits the EXPECT → **Won't handle** beyond
  existing behavior: did-not-bite red at the gate is the correct loud outcome
  for a rotted binding (story 7 verifies the two known-fuzzy bindings by
  bite before done).
- Behavior-owned fixture carrying a CHECK file → **Won't handle**: CHECK
  continues to govern tier only; its conformance-scope meaning stays inert
  for non-conformance families, unchanged by this spec.
- Adopting repo with its own family names → **Won't handle**: absent the
  variable the argv is unchanged; an adopting repo's `behavior-owned` family
  gets the same scoping against its own tree, and its structural reds are the
  same loud postures.
- Concurrent sweeps / worker interleaving → **Won't handle**: scheduling is
  untouched; the existing concurrency tests keep their authority.

## Out of scope

- **Stage 2 — in-process bite checks** (map #7): a separate capability — it
  replaces the nesting transport rather than narrowing argv — and specs the
  moment this lands. ~25 edits, ~6 gate runs.
- **Prepared-artifact hoist** (map #3): separate capability on the contract
  suite itself, pending the #2 inventory and the reviewer's independence
  ruling. ~15 edits, ~4 gate runs.
- **Oracle-semantics levers** (verdict caching, `-count=1`): enter FT91 only
  if the post-stage-2 re-measure stays above 60 s (map #5). Unpriced until
  their ruling.
