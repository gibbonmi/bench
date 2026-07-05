# Test-suite structure

Status: implemented

## Problem

`bench structure` is red: `internal/conformance` holds 13 source files and
`internal/contract` holds 23, both over the 12-file directory cap. The two caps
squeeze jointly — the 400-line file cap forces splits, the splits feed the
dir cap — so the suite cannot grow without tripping the oracle. Separately,
`tests/canary/` is 59 flat fixture directories whose check-family classification
lives in a 59-entry hand-tended registry map, and the one-canary-per-check rule
is manual, not enforced. This is one debt-reduction program, not three bugs.

## Solution

Three independent slices, all behavior-preserving where they touch tested code:

- **Contract** (`internal/contract`, 23 files) splits into a root
  fixture-harness library plus three family subpackages (`runtime/`, `axi/`,
  `surface/`). The gate's contract phase and the Go-core exclusion widen to the
  subtree so the suite still runs exactly once.
- **Conformance** (`internal/conformance`, 13 files) stays one package but merges
  two helper-beside-consumer file pairs down to 11 files. It does *not*
  subpackage — its single check registry and single pinned gate entry are one
  deep module by design.
- **Canary** (`tests/canary/`) regroups into `<family>/<fixture>/`, derives the
  family from the path instead of a registry column, walks two levels, and gains
  a meta-check that fails conformance if any check family owns zero fixtures.

When all three land, `bench structure` reports zero issues and `bench gate` stays
green with the suite running exactly as before.

## User stories

1. **As the kit maintainer, I want `internal/conformance` merged from 13 files to
   11**, so that the directory clears the 12-file cap without subpackaging a
   registry whose point is to be one list. Merge `harness_helpers_test.go` (145
   lines) into `harness_test.go` (88) and fold `root_conformance_test.go` (17)
   into `gate_entry_test.go`; every merged file stays under the 400-line cap and
   `docs_workflow_helpers_test.go` stays separate (two consumers).
   Line: claude-sonnet-5 / low. Pure helper-beside-consumer concatenation at a
   known seam that the structure check and the conformance compile+run fully
   observe, so the cheapest tier is honest here.

2. **As the kit maintainer, I want the shared contract harness promoted to an
   exported surface**, so that family subpackages can import the fixture harness
   and the generic assertions instead of re-declaring them. `helper.go` already
   exports the fixture harness in `package contract`; promote the generic
   test-scoped helpers in `runtime_helpers_test.go` (`runAt`, `requireContains`,
   `requireNotContains`, `requireIntEqual`, `lineCount`, `nonEmptyLines`, and
   peers used across families) into an exported non-test assert surface in the
   root package. Family-specific helpers (`copyRuntimeFile`, `cksum`,
   `countStatusRows`, …) stay unexported in their family package.
   Line: claude-opus-4-8 / medium. The generic-vs-family cut and the call-site
   rewrite across every consumer is a judgment the gate cannot grade, and it
   carries the non-gate-observable behavior-preservation claim.

3. **As the kit maintainer, I want `internal/contract` split into `runtime/`,
   `axi/`, and `surface/` subpackages**, so that the directory clears the cap by
   grouping into real Go packages (the check's own remedy). `runtime/`: gate,
   git, shift, shift adapters, shift helpers, status, structure, worktree.
   `axi/`: axi, wave2, guards, fail_closed, axi asserts. `surface/`: doctor,
   doctor_shim, link, link_marker, package, go_routing, binary_repair. The
   relocation is move-only and behavior-preserving — same tests, same subtest
   counts, `BENCH_CONTRACT_ROOT` reaching every subpackage unchanged.
   Line: claude-opus-4-8 / medium. Behavior-preservation across a large
   multi-file move needs care the gate can't fully see (subtest-count equality).

4. **As the kit maintainer, I want the gate's contract phase to run the whole
   subtree**, so that the split subpackages still execute and stay pinned. The
   `.bench/gate.sh` contract line becomes
   `BENCH_CONTRACT_ROOT="$root" go test -count=1 ./internal/contract/...` — still
   one phase, one run — and the `gate_entry` conformance anchor that pins that
   line verbatim updates in the same change.
   Line: claude-opus-4-8 / medium. This edits what the oracle covers;
   craft-gate discipline and the anchor re-pin apply.

5. **As the kit maintainer, I want the Go-core exclusion widened to the contract
   subtree**, so that subpackages like `internal/contract/axi` cannot escape the
   exclusion and silently resurrect the duplicate contract run. `goCoreTestPackages`
   excludes the whole `internal/contract/` subtree (prefix match) instead of the
   current exact/suffix match on `internal/contract`.
   Line: claude-opus-4-8 / medium. A coverage-losing edit to the oracle with a
   TDD bite proof; the too-wide direction is the only one that loses coverage.

6. **As the kit maintainer, I want `tests/canary/` regrouped into
   `<family>/<fixture>/`**, so that the family classification lives in the tree
   instead of a hand-tended registry column. Move all 59 fixtures under a family
   directory (44 conformance fixtures under their 7 check families, 15 behavior
   fixtures under `behavior-owned/`); `Sweep` and `fixtures()` walk two levels;
   `bench canary` behavior (exit codes, diagnostics, vacuous-EXPECT rejection,
   bite semantics) is unchanged.
   Line: claude-opus-4-8 / medium. The two-level walk and name-uniqueness
   handling is logic, not a pure move.

7. **As the kit maintainer, I want the family derived from the path**, so that the
   59-entry classification map collapses into the tree. Delete the `Family` field
   from `canaryFixtureRegistry`; the registry keeps only the owner and the
   retired-shell-source columns; the classification test derives each fixture's
   family from its parent directory.
   Line: claude-opus-4-8 / low. Mechanical once the layout lands, but it rewrites
   the classification test's source of truth.

8. **As the kit maintainer, I want a meta-check that every check family owns at
   least one fixture directory**, so that one-canary-per-check is enforced instead
   of hand-tracked. A new check inside `RunConformance` reads the canonical list
   of conformance check families and emits a conformance diagnostic for any family
   with zero fixture directories under `tests/canary/`; the rename-coupling is
   intended (a family rename must move its fixtures).
   Line: claude-opus-4-8 / medium. A new oracle check needs a real bite proof;
   craft-gate discipline applies.

## Implementation decisions

- **Contract root becomes the exported harness library.** `helper.go` already
  holds the exported `Fixture`/`Probe`/`NewFixture`/`KitRoot` harness in
  `package contract`. Add a non-test `asserts.go` (also `package contract`) with
  the promoted generic assertions, exported (`RunAt`, `RequireContains`,
  `RequireNotContains`, `RequireIntEqual`, `LineCount`, `NonEmptyLines`, …).
  Subpackages `package runtime` / `package axi` / `package surface` import
  `github.com/gibbonmi/bench/internal/contract` and call `contract.RequireContains`
  etc. The shared surface must live in non-test `.go` files (not `_test.go`),
  because a package's `_test.go` symbols are not importable by other packages;
  `helper.go` is the existing precedent.

- **Family subpackage assignment** (from map #2, verified against the current 23
  files): root lib = `helper.go`, `helper_test.go`, `asserts.go`. `runtime/` =
  `runtime_gate`, `runtime_git`, `runtime_shift`, `runtime_shift_adapters`,
  `runtime_shift_helpers`, `runtime_status`, `runtime_structure`,
  `runtime_worktree`, plus the runtime-specific remainder of
  `runtime_helpers_test.go`. `axi/` = `axi`, `axi_wave2`, `axi_guards`,
  `axi_fail_closed`, `axi_helpers`. `surface/` = `doctor`, `doctor_shim`, `link`,
  `link_marker`, `package`, `go_routing`, `binary_repair`. Each subpackage lands
  ≤ 12 files; the root lands at 3.

- **Gate rewiring is one change with three touchpoints** (map #3): the contract
  phase becomes `./internal/contract/...`; `goCoreTestPackages` widens to the
  whole `internal/contract/` subtree; the `gate_entry` anchor re-pins the new
  gate line verbatim. No-weakening: the contract suite still runs exactly once,
  and the unit sweep loses only the packages the contract phase itself owns.

- **Conformance stays one package** (map #4). Two helper-beside-single-consumer
  merges, not responsibility consolidation. This is not a reopen of the #1
  consolidation rejection.

- **Canary layout and registry.** `tests/canary/<family>/<fixture>/`. `Sweep`
  walks two levels; `runFixture` still names a fixture by `filepath.Base`, and all
  59 current fixture names are globally unique, so base-name diagnostics stay
  unambiguous after nesting. The `Family` field is removed from
  `fixtureRegistration`; the `conformanceFixture` and `behaviorFixture`
  constructors and the classification test change accordingly. The two-level walk
  must enumerate directory entries, never glob a pattern.

- **Canonical family list for the meta-check (seam the map left silent — flagged
  for sign-off).** To catch a check family with *zero* fixture directories, the
  meta-check needs an independent statement of which families should exist —
  deriving the family set from the directory tree alone is circular. Decision:
  add a canonical `conformanceFamilies` inventory (a single ordered `[]string` of
  the 7 check-family names) in the conformance package; the meta-check asserts
  each owns ≥ 1 fixture directory, and the classification test rejects a fixture
  filed under an unknown family. This is the sanctioned rename-coupling from map
  #5 (a family rename touches the inventory and moves the dir), not a duplicated
  fact — the inventory states *what checks exist*, the tree states *where each
  canary lives*.

## Testing decisions

- **What a good test is here:** exercise the black-box behavior of the oracle
  surfaces — `bench structure`, `bench gate`, `bench canary`, and
  `RunConformance` diagnostics — not the internal file layout. The contract split
  is graded by the existing contract suite staying green plus a one-time
  subtest-count diff, because move-only relocation has no new behavior to assert.
- **Seams tested** (prior art in parentheses): `bench structure` output
  (`runtime_structure_test.go`); the `gate_entry` anchor and `goCoreTestPackages`
  unit (`gate_entry_test.go`, `package_core_checks_test.go`); the canary sweep and
  registry classification (`canary_test.go`, `registry_test.go`).
- **Gate command:** the project gate, `bench gate`.

### Seam diagram

**Seam 1 — `bench structure` (whole-program clearance):**

    trigger: bench structure   /   bench gate → structure check
        │
        ▼
    internal/contract/ (23)  ──▶  [ dir-cap walk over source tree ]  ──▶  "N issue(s)"
    internal/conformance/(13) ──▶ [                              ]
                                      ◀ tests attach here: assert zero dir-cap
                                        issues after the restructure

**Seam 2 — gate contract phase + Go-core exclusion:**

    trigger: bench gate  (and go test ./internal/conformance)
        │
        ▼
    .bench/gate.sh line ──▶ [ go test -count=1 ./internal/contract/... ] ──▶ pass once
    go list ./...        ──▶ [ goCoreTestPackages: drop internal/contract/ ] ──▶ unit set
                                 ◀ tests attach here:
                                   - gate_entry anchor pins the new gate line verbatim
                                   - unit: contract-subtree pkg excluded,
                                     internal/conformance sibling survives

**Seam 3 — canary regroup + per-family meta-check:**

    trigger: bench canary        /        RunConformance → TestRootConformance
        │                                        │
        ▼                                        ▼
    tests/canary/<family>/<fixture>/ ──▶ [ Sweep: 2-level walk, materialize ] ──▶ bite/no-bite
    conformanceFamilies (7 names)    ──▶ [ meta-check: ≥1 fixture dir/family ] ──▶ diagnostic
                                             ◀ tests attach here:
                                               - canary green with nested fixtures
                                               - classification derives family from path
                                               - meta-check red when a family dir is empty

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | conformance dir ≤ 12 files | `bench structure` | `bench structure` today: `13 source files (max 12) internal/conformance/` (observed RED) | the check counts files per dir; an un-merged conformance dir stays over cap |
| 2 | subpackages import shared asserts | contract compile via `go test ./internal/contract/...` | before promotion, `package axi` referencing `requireContains` fails to compile (unexported, cross-package) | the compiler rejects test-scoped helpers used across packages, forcing the exported surface |
| 3 | contract dirs ≤ 12; suite unchanged | `bench structure` + contract suite | `bench structure` today: `23 source files (max 12) internal/contract/` (observed RED); subtest-count diff not gate-observable — one-time `go test -v` inventory before/after | file-count red proves the split; the inventory diff proves the move added/dropped no subtest |
| 4 | gate runs the whole contract subtree | `gate_entry` anchor test | flipping the gate line to `./internal/contract/...` without updating the anchor turns `TestRootConformance` RED (anchor pins the literal) | the anchor asserts the verbatim gate line, so phase text and anchor must move together |
| 5 | exclusion covers the subtree | unit on `goCoreTestPackages` | new unit asserting `internal/contract/axi` is excluded and `internal/conformance` survives fails RED under the current exact/suffix match | the too-wide direction is the only one that loses coverage; the unit pins the boundary |
| 6 | canary green with nested fixtures | `bench canary` | after moving fixtures under `<family>/` with the one-level walk, `fixtures()` returns family dirs with no `EXPECT` → `bench canary` reports "no EXPECT" (RED) | the walk depth must match the layout depth or every fixture reads as malformed |
| 7 | family derived from path | `TestCanaryFixtureRegistryClassifiesEveryFixture` | after nesting, the flat-dir classification test finds family dirs, not fixture names → mismatch RED | the test's source of truth must move from the registry column to the parent dir |
| 8 | meta-check bites on an empty family | `RunConformance` via a bite test | a test that points the meta-check at a tree with one family dir emptied expects the zero-fixtures diagnostic; absent the check, no diagnostic is emitted → RED | proves the check actually fires when one-canary-per-check regresses |

### Edge inventory

Walked per behavior; each resolved as a coverage row above or a **Won't handle**
line here.

- **Empty / absent input** — family dir with zero fixtures → meta-check (story 8,
  row above). Fixture missing `EXPECT` or `files/` → existing `runFixture`
  diagnostics, unchanged.
- **Boundary — fixture name collision across families.** `runFixture` names a
  fixture by `filepath.Base`; two families with the same fixture base name would
  produce ambiguous diagnostics. Covered: a test (or the classification test's
  uniqueness assertion) asserting fixture base names stay globally unique across
  families, since all 59 current names already are.
- **Malformed input — paths with spaces/globs in fixture names.** The two-level
  `Sweep`/`fixtures()` walk must enumerate `os.ReadDir` entries, never expand a
  glob pattern. Covered by the existing space-path fixture harness plus the
  canary green assertion (story 6).
- **Re-run idempotency** — the relocation is move-only; `bench gate` green is
  idempotent across re-runs. Covered by the gate itself.
- **Interrupted/partial state — a subpackage that fails to compile mid-split.**
  *Won't handle* — a partial split is a build-time state the gate catches
  wholesale (compile failure); no per-partial recovery seam is warranted.
- **Hostile environment** — `BENCH_CONTRACT_ROOT` must reach every subpackage
  test unchanged. Covered: the contract suite green under `./internal/contract/...`
  (story 4) exercises it, since the tests read that env.
- **Boundary — meta-check family set drift.** A check family named in
  `conformanceFamilies` with a typo (no matching dir) → meta-check RED, which is
  the intended coupling. *Won't handle beyond that* — reconciling the inventory
  against `RunConformance`'s actual dispatch is out of scope (see Out of scope).

## Out of scope

- **Reconcile `conformanceFamilies` against `RunConformance`'s dispatch.** A
  second meta-check proving the canonical family list matches the checks actually
  run (not just the fixture tree) is a distinct enforcement capability with its
  own bite proof. Deferred estimate: ~3 edits, ~2 gate runs. Parked rather than
  built now because the family set is small, stable, and reviewer-visible in one
  file.
- **Unify the subprocess-capture seams** (conformance `runProbe`, `Harness.Run`,
  canary `defaultRunner`). A separate parked change with a bigger blast radius —
  it spans conformance and canary and the exit-code-extraction idiom re-derived
  3× (per the #6 bloat audit). Deferred estimate: ~8 edits, ~4 gate runs. It is
  its own dedup feature, not the rest of this one.
