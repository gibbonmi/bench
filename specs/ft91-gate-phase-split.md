# FT91 slice C — the checkGoCore split

Status: implemented

Compiled from `decisions/gate-pipeline.md` (#3 the split, #6 fixture migration,
#7 parity, #9 output contract; Handoff items 1–9), which the reviewer confirmed
closed on 2026-07-27. Slice A (canary check-scoping) and slice B (manifest + DAG
runner) shipped and are retired; this slice is the first consumer of both. The
map's dependency note bundles FT143's kit-root family→check binding assertion
into this spec (reviewer-bundled 2026-07-27), and the slice-B spec sequenced the
kit's own `.bench/phases.json` and the consumer-facing manifest docs here.

## Problem

`checkGoCore` runs seven toolchain steps — gofmt, a throwaway build, vet, the
core `go test`, a filtered `internal/conformance` run, the worktree race test,
and (at ship) cross-compile — serially inside one conformance check. Three costs
follow.

**Duration.** `package-core-guard` is now the gate's long pole at ~86 s against a
4m36s whole gate, and every second of it is serial work inside a single Go test
function that the DAG scheduler cannot overlap with anything.

**Recursion and stale state.** The compile-and-test work lives inside the test
binary the gate is running, which is the surface that produced both 600 s
package-timeout incidents and the stale-state class the FT91 build paid ~12
messages of git archaeology to diagnose.

**Overreach.** The race step runs `go test -race ./internal/worktree` against any
root carrying a `go.mod`, so a linked Go repo without that package reds on a
check the kit wrote for itself.

Meanwhile the gate's own restaging cannot be trusted by review alone: gofmt, vet,
and the race step have **zero canary fixtures today**, and the two fixtures that
do grade moved steps both die on a mechanical move, so almost nothing proves a
moved step still bites at its new owner.

## Solution

Hoist the toolchain work out of the conformance test binary into first-class gate
phases with declared dependencies, leaving conformance to grade structure only.
gofmt, vet, and the core `go test` become phases the built-in fallback table
probes for (so canary fixtures and linked Go repos keep toolchain grading); the
worktree race test and the filtered `internal/conformance` suite become phases
the kit declares in its own `.bench/phases.json` (so they stop grading trees that
never asked for them); the throwaway build validation collapses into the gate's
existing build phase. The new phases declare no build edge, so they overlap with
the build and with each other instead of running end to end inside one check.

Because the restaged phases are only as good as their tripwires, the same pass
migrates the fixture inventory: the strays gain `CHECK` files, the two movers get
phase-native EXPECTs at their destination families, four new fixtures cover the
steps that never had one, and an inventory test asserts every moved step owns a
fixture bound to its new owner. Green keeps meaning the same thing, in both
directions — the dev tier and `prep-release` each keep every check they had.

## User stories

1. As a worker running the gate, I want `gofmt` to be a first-class phase that
   reds when `gofmt -l` names any file, so that formatting is graded in parallel
   with the rest of the gate instead of serially inside a conformance check.
   Line: opus (mid) / medium. The step and its diagnostic already exist; the work
   is relocating them behind a phase whose exit code is honest.
2. As a worker running the gate, I want `go vet` to be a first-class phase, so
   that vet overlaps the build and the other toolchain phases.
   Line: opus (mid) / low. `go vet` exits nonzero on its own, so the phase is a
   direct argv with no policy behind it.
3. As a worker running the gate, I want the core `go test` run to be a
   first-class phase whose package set comes from one place, so that the
   exclusion policy that keeps contract, conformance, and release-only packages
   off this run cannot drift from the surfaces that own them.
   Line: opus (mid) / medium. The policy is decided and already written; moving
   it out of a test file into importable code is the whole task.
4. As a linked-repo owner, I want the worktree race test declared only in the
   kit's own manifest, so that my Go repo without `internal/worktree` stops
   reding on a check the kit wrote for itself.
   Line: opus (mid) / medium. Closing a fail-open-for-the-kit, fail-wrong-for-
   consumers overreach, which is the audience question this slice keeps meeting.
5. As a maintainer, I want the filtered `internal/conformance` suite run declared
   only in the kit's own manifest, so that the one step whose non-recursion
   contract is kit-specific stops being probed for on arbitrary roots.
   Line: opus (mid) / medium. Same shape as story 4, with the recursion contract
   as the thing that must survive the move intact.
6. As a maintainer, I want the throwaway build-validation step gone and the
   gate's existing build phase to be the single build, so that one owner writes
   one binary and the throwaway temp-dir path disappears with the reason it
   existed.
   Line: opus (mid) / medium. A deletion plus a fixture-tree change; the risk is
   in what the deletion silently stops grading, which story 15 canaries.
7. As a reviewer, I want the residual `package-core-guard` check to keep its
   structural siblings — package files, release and native workflow, preflight,
   identity strings, guard manifests, resolver drift — plus the "go.mod present
   but no Go toolchain" diagnostic, and to grade any root cheaply, so that the
   28 fixtures that grade those diagnostics stay exactly where they are.
   Line: opus (mid) / medium. Mostly subtraction, but it is the check whose
   meaning must not change, so it is graded by the fixtures that already exist.
8. As a canary fixture and as a linked Go repo, I want the built-in fallback
   table to gain go.mod-probed gofmt, vet, and test phases, so that a root with
   no manifest keeps toolchain grading and the fixture EXPECTs still bite.
   Line: opus (mid) / medium. The build phase's probe pattern is the prior art;
   the judgment is what the probe requires before a phase materializes.
9. As the kit, I want to ship my own `.bench/phases.json` declaring the whole kit
   phase table, so that the kit-specific phases exist as data and the manifest
   loader slice B shipped has a real consumer.
   Line: opus (mid) / medium. The schema is decided and the loader is built; the
   care is in the table being byte-equivalent in effect to today's.
10. As a worker waiting on the gate, I want the new toolchain phases to declare
    no `needs` edge on the build phase, so that they overlap the build and each
    other — which is where this slice's duration win actually comes from.
    Line: opus (mid) / medium. The edge set is a correctness claim about what
    reads `dist/bench`, and it is asserted rather than assumed.
11. As a maintainer, I want one `bench gate-go <step>` plumbing subcommand to own
    the four toolchain steps whose argv is policy rather than a plain command —
    gofmt (`gofmt -l` exits 0 while naming files), test (the package set is
    policy), race (a `-run` filter that matches nothing exits 0), and
    conformance-suite (its skip pattern has one source in the registry) — so
    that the manifest stays declarative and no policy is retyped into JSON.
    Line: opus (mid) / high. Three fail-open shapes and one single-source
    obligation collapse here; getting any of them wrong makes a phase that can
    never red or a second copy of a contract.
12. As a canary fixture author, I want `FixturePhase` to route phase-named
    top-level families (`build`, `gofmt`, `vet`, `test`, `race`) and
    `phasesForMode`'s owner allowlist widened to the phase table's names, so that
    a fixture in a phase-named family runs only the phase that owns its failure.
    Line: opus (mid) / medium. Two coupled edits the map's watch-out says must
    land together, since a name absent from either runs every non-canary phase.
13. As a maintainer, I want the ship-tier cross-compile matrix to stay reachable
    through the residual check after the split, so that ship green does not lose
    the four-platform matrix the same way story 18 keeps it from losing the
    release-only suites.
    Line: opus (mid) / medium. A subtraction risk with no fixture behind it
    today, which is exactly why it needs an assertion rather than care.
14. As a reviewer, I want the two fixtures whose EXPECTs no seam will emit after
    the split — `go-build-broken` and `go-test-failing` — moved to their
    destination families with phase-native EXPECTs, and `go-build-broken`'s tree
    given what the build phase's probe requires, so that neither fixture silently
    goes green on a tree the gate no longer grades the same way.
    Line: opus (mid) / high. Both fixtures die on a mechanical move; each needs a
    recorded red at its new owner before it counts.
15. As a reviewer, I want new canary fixtures authored in this same pass for
    every moved step that has none today — gofmt, vet, race, and the filtered
    conformance suite — so that the split does not ship four uncanaried seams.
    Line: opus (mid) / high. Authoring a fixture that bites for the intended
    reason, and proving it, is the slice's real evidence.
16. As a maintainer, I want an inventory test asserting every moved step has a
    canary fixture bound to its destination phase, so that a later edit that
    orphans a step turns the gate red instead of quietly leaving it ungraded.
    Line: opus (mid) / medium. The `TestRegistryBindsEveryCheck` analogue, with a
    settled shape to copy.
17. As a maintainer, I want `TestCoreSubprocessFailuresUseProbeFormatter`'s label
    inventory to follow the split, so that the text scan keeps guarding the
    labels that remain instead of failing on the ones that left.
    Line: opus (mid) / low. A known mechanical casualty named in the map, but it
    edits a conformance guard, which the profile routes to the mid tier
    regardless of how small the edit is.
18. As a maintainer, I want `bench prep-release` to keep running the release-only
    package suites after the core test step leaves the conformance check, so that
    ship green does not silently lose three suites in the restaging.
    Line: opus (mid) / high. Parity in the ship direction, which the map's
    fixture-based parity requirement does not reach; flagged below.
19. As a reviewer, I want an unbound canary conformance family to red before the
    sweep pays for it, on both `bench gate` and standalone `bench canary`, so
    that the family→check binding cannot be reintroduced late on one path and
    silently skipped on the other (FT143).
    Line: opus (mid) / high. The scoping predicate is the subtle part: it must
    fire for the kit's own tree and stay silent for every adopting repo.
20. As a linked-repo owner, I want `.bench/phases.json` documented in the
    reference file that ships to me, so that the manifest format is discoverable
    without reading the kit's Go source or a retired spec.
    Line: fable (top) / high. The profile's cached leverage override for kit
    guidance prose that ships to every consumer; flagged in the table for veto.
21. As a canary and contract consumer, I want the gate's output and timing
    contract unchanged — summary lines, per-check conformance timing, prefix
    framing, inner-mode byte shape, exit codes — so that substring EXPECTs and
    downstream parsers keep biting across the restaging.
    Line: opus (mid) / medium. Continuity pinned by existing suites; the work is
    not breaking them.
22. As the reviewer, I want the post-change whole-gate and per-phase wall clock
    measured against the 4m36s / `package-core-guard` ~86 s baseline and reported
    in the close, so that the arm's premise is evidenced rather than asserted.
    Line: opus (mid) / low. One measured run and one line of reporting; the map
    names this as the slice's only non-gate-assertable outcome.

## Implementation decisions

**The split.** `checkGoCore` loses gofmt, the throwaway build, vet, the core
`go test`, the filtered `internal/conformance` run, and the race test. What
remains is the "go.mod present but no Go toolchain on PATH" diagnostic and the
ship-tier `crossCompileMatrix` call, both of which grade any root cheaply and
neither of which the phase table can own (cross-compile stays ship-tier and
untouched, per the map). The registry check name `package-core-guard` does not
change, so the family binding and the 28 residual fixtures stay put.

**Phase set and edges.** The built-in fallback table gains `gofmt`, `vet`, and
`test`, each materializing only when the graded root carries a `go.mod` and `go`
resolves on PATH — the build phase's probe pattern, extended. The kit's
`.bench/phases.json` declares the same table plus `race` (the worktree race test)
and `conformance-suite` (the filtered `internal/conformance` run). None of the
five new phases declares a `needs` edge: none of them execs `dist/bench`, which
is the only artifact the build edge exists to sequence. That claim is asserted,
not assumed — story 10's row proves the enumerated test package set excludes
`internal/contract`, the one core-adjacent consumer of the built binary.

**The `bench gate-go` subcommand** *(map-silent extension — flagged)*. The map's
Handoff names one new plumbing subcommand, for test-package enumeration. Four
steps need it, not one. Three of the moved commands cannot red on their own:
`gofmt -l` exits 0 while printing the offending files; the test step's package
set is exclusion policy; and the race step's `-run` filter exits 0 when it
matches nothing, which is exactly why today's check asserts `=== RUN` in the
output. The fourth, the filtered conformance suite, could red on its own but
carries `registry.InnerSkipPattern()` in its argv — retyping that regex into
JSON would give the non-recursion contract a second source, which this repo's
code standard forbids. `bench gate-go <gofmt|test|race|conformance-suite> [root]`
owns all four. `vet` stays a direct `go` argv, since it exits nonzero honestly
and carries no policy. An unknown step name is a usage error (exit 2), never a
silent success.

**Where the exclusion policy lives.** `isExcludedTestPackage` and
`releaseOnlyPackages` sit in `package_core_checks_test.go` today, so the compiled
binary cannot import them. They move to `internal/conformance/registry`, which
already owns the tier vocabulary, the scoping env names, and `InnerSkipTests`,
and which imports nothing from `internal/conformance`. The subcommand reads the
tier from the existing `BENCH_CONFORMANCE_TIER`, so dev and ship keep resolving
the same package sets they resolve today.

**Diagnostic continuity** *(map deviation — flagged)*. The map's #6 answer gives
the migrated fixtures "tool diagnostics plus the stable `phase <name>: red`
summary" as their EXPECTs. The summary half is not available: `runPhasesSequential`
— the inner-mode runner every canary fixture drives — emits no per-phase summary
lines at all, and `TestRunnerInnerModeByteShape` pins that byte shape while #9
closes against redesigning the output contract. So the new and migrated fixtures
EXPECT the tool's own diagnostic only. To keep the vocabulary continuous, the
gofmt step keeps today's exact label (`gofmt: unformatted Go files:` followed by
the file list) rather than emitting bare `gofmt -l` output; vet, test, and race
stream the tool's own output, and each fixture's EXPECT is the most stable stem
that output offers, settled at build time within that rule (the map's Handoff
item 7 already delegates the exact strings).

**Fixture migration.** *Map-versus-tree correction:* #6 decided that the seven
stray fixtures without `CHECK` files would gain them, but that work shipped with
slice A — all seven files exist today, naming five distinct checks (bounds-policy
three times, plus default-branch-single-source, marker-wait-deadlines,
skip-ownership, and subcommand-routing). This slice does not restage it, and no
story claims it; the map's count of "six other checks" is stale by one.

What does move: the two movers relocate to phase-named families —
`tests/canary/build/go-build-broken/`, whose tree gains the `scripts/go-build.sh`
the build phase's probe requires, and `tests/canary/test/go-test-failing/`. Four
new fixtures cover the steps with none today: gofmt, vet, race, and the filtered
conformance suite. `race` and `conformance-suite` are kit-manifest-only phases,
so a fixture grading either ships its own minimal `.bench/phases.json` declaring
that phase — the loader slice B built is what makes this possible, and it is the
only way a fixture root can carry a phase the built-in table does not probe for.
Fixtures in phase-named families resolve to the empty scope, so they join the
existing unscoped baseline group; that baseline runs a superset of what the
scoped fixture runs, which is conservative for the vacuity comparison.
Sequencing: `unboundCanaryFamilies` reds on any new top-level family until
`FixturePhase` and the allowlist widen, so story 12's two edits and the new
family directories land in one commit.

**FT143 — the binding assertion.** Both entry points already funnel through
`canary.SweepTier` (`bench canary` calls it directly; the gate's canary phase
execs `bench canary <root>`), so one assertion at the top of `SweepTier` reaches
both. It is scoped to the kit by comparing the swept root against `BENCH_KIT`,
which `bin/bench.sh` always exports to the binary — equal means the kit is
grading itself and the kit-owned family table is authoritative; unequal means an
adopting repo, whose families the table will never carry, and the assertion stays
silent. Both sides are resolved through `filepath.EvalSymlinks` before comparison
*(map-silent — flagged)*: `bin/bench.sh` derives its default with a physical
`cd -P` while `SweepTier` normalizes with `filepath.Abs` alone, so a symlinked
checkout path would otherwise make the kit look like an adopting repo and skip
its own assertion. The residual case — a globally installed wrapper run inside
the kit checkout, where `BENCH_KIT` honestly points at the install directory — is
a won't-handle below: there the sweep really is grading a tree the running kit
does not own. The scan itself is single-sourced: `unboundCanaryFamilies` moves from
`internal/conformance` into `internal/canary` (which already owns the fixture-tree
layout) as an exported function, and the existing conformance check calls it
instead of rederiving it. One derivation, two callers, two entry points.

**Ship parity** *(map-silent — flagged)*. Today `prep-release`'s
`conformance-ship` step reaches `checkGoCore`, whose ship-tier package
enumeration pulls in `internal/preflight`, `internal/releaseevidence`, and
`internal/publication`. Once the test step leaves the check, that path is gone
and ship green would silently stop covering three suites. `Steps` therefore gains
a step invoking `bench gate-go test` with `BENCH_CONFORMANCE_TIER=ship`. This
also closes FT142's flagged item — the release-only `go test` step the decisions
promised, folded silently into ship-tier `goCoreTestPackages` — by making it an
explicit step. It cannot be demonstrated green end to end, because `prep-release`
stays blocked on FT116's `enumerateGuards` leak; the coverage row attaches at the
`Steps` seam, which is where the existing `conformance-ship` assertion lives.

**Test-cache posture.** The core `go test` run carries no `-count=1` today, and
the phase keeps it that way. Adding one would change cost without changing what
green means, and the gate's `-count=1` uses are deliberate arms elsewhere.

**Untouched.** `.bench/gate.sh` stays the one-line exec entry. The DAG scheduler,
the manifest loader, the verdict and subject machinery, `internal/bounds`, the
canary width budget and its `GOMAXPROCS` pin, the conformance per-check timing
driver, the `BENCH_CONFORMANCE_CHECK` scoping mechanism, and the exit-code set
all stay as slice A and B left them.

**Map-silent default inventory (veto list).** Everything decided here that the
map does not carry, in one place: (a) `bench gate-go` covering four steps rather
than one; (b) `vet` as a direct argv with no wrapper; (c) the exclusion policy
relocating into `internal/conformance/registry`; (d) fixture EXPECTs dropping the
`phase <name>: red` half, with the gofmt label preserved to compensate; (e)
kit-manifest-only phases canaried through fixture-local manifests; (f) the
`prep-release` ship-parity step and its FT142 overlap; (g) symlink-resolved
`BENCH_KIT` equality as the kit-root discriminator for the FT143 assertion, and
the installed-wrapper case left as a won't-handle; (h) no `-count=1` on the test
phase; (i) the new phases declaring no build edge; (j) the residual check keeping
the name `package-core-guard`; (k) the conformance suite's skip pattern staying
in Go behind `gate-go` rather than being retyped into the manifest; (l) the phase
name `conformance-suite` itself; (m) cross-compile continuity getting an assertion
of its own (story 13), which the map assumed rather than required.

## Testing decisions

- A good test here drives the real seam with cheap subprocesses and fixture trees
  and asserts observable output, exit codes, and which phases ran — never
  scheduler or registry internals. Phase-table tests assert the resolved table
  for a probed root; subcommand tests assert exit code and diagnostic against
  temp trees; canary tests drive `SweepTier` with a fake `Runner` and assert the
  `RunCall`s; the fixtures themselves are the end-to-end bite.
- Prior art to follow: `TestPhaseTableBuildPhase` (probe-conditional phases),
  `TestPhasesCommandManifestFieldsEndToEnd` (manifest end-to-end),
  `TestRegistryBindsEveryCheck` (inventory parity), the `scope_test.go` fake-runner
  pattern (sweep assertions), `TestRunnerInnerModeByteShape` and
  `TestRunnerSummaryLineByteShape` (output contract), and the `Steps` unit test
  in `internal/preprelease` (ship step list).
- Gate: `.bench/gate.sh` (the project gate), green required to commit. Note that
  adding `.bench/phases.json` changes the `.bench` tree, so the pre-push hook
  wants a fresh `bench gate pin`.

### Seam diagram — phase table resolution (fronted by `gate-phases`)

    trigger: bench gate → .bench/gate.sh → bench gate-phases <root>
        │
        ▼
    <root>/.bench/phases.json  ──▶  [ phaseTable: manifest or built-in    ] ──▶ []Phase (kit: build, gofmt,
    <root>/go.mod              ──▶  [ BenchkitPhases: probe go.mod + go   ]      vet, test, race,
    `go` on PATH               ──▶  [ for gofmt/vet/test, go-build.sh for ]      conformance-suite, …)
    <root>/scripts/go-build.sh ──▶  [ build; edges only where dist/bench  ] ──▶ []Phase (probed subset)
                                    [ is actually consumed                ] ──▶ no toolchain phases
                      ◀ tests attach here: BenchkitPhases against temp roots with
                        each probe input present/absent; assert phase names, argv,
                        and Needs. PhasesCommand with a stubbed table for end-to-end.

### Seam diagram — the toolchain step subcommand (`bench gate-go`)

    trigger: a phase's argv, from either the built-in table or a manifest
        │
        ▼
    step name + root  ──▶  [ gate-go gofmt: gofmt -l, red on any named file  ] ──▶ tool diagnostic + exit 0/1
    BENCH_CONFORMANCE ──▶  [ gate-go test:  go list → exclusion policy →     ] ──▶ usage error + exit 2
      _TIER           ──▶  [                go test <packages>               ]      (unknown step)
                           [ gate-go race:  -race run + did-it-actually-run  ]
                      ◀ tests attach here: run the subcommand against temp trees
                        (unformatted file, failing test, absent package) and assert
                        exit code plus diagnostic; assert the package set directly.

### Seam diagram — canary routing, scoping, and the binding assertion

    trigger: bench canary <root>  ──and──  bench gate → canary phase → bench canary <root>
        │
        ▼
    tests/canary/<family>/<fx>  ──▶  [ SweepTier: kit-root binding assert → ] ──▶ RunCall per fixture
    CHECK file (stray override) ──▶  [ selectTier → fixtureScope →          ]      (Env: PHASE + CHECK)
    BENCH_KIT vs root           ──▶  [ FixturePhase(family) → owner phase   ] ──▶ unbound-family error
                                     [ → phasesForMode filters inner table  ]      before any inner run
                      ◀ tests attach here: SweepTier with a fake Runner over a temp
                        fixture tree; assert each RunCall's Env, the baseline count,
                        and that an unbound family errors before any RunCall happens.

### Seam diagram — residual check and parity inventory (conformance)

    trigger: gate conformance phase → TestRootConformance → package-core-guard
        │
        ▼
    graded root  ──▶  [ residual: package files, workflows, preflight,   ] ──▶ structural diagnostics
                      [ identity, guard manifests, resolver drift,       ]
                      [ go.mod-but-no-toolchain, ship cross-compile      ]
    phase table  ──▶  [ parity inventory: every moved step has a fixture ] ──▶ orphaned-step diagnostic
    tests/canary ──▶  [ bound to its destination phase                   ]
                      ◀ tests attach here: existing 28 fixtures for the residual;
                        the inventory test reads the phase table and the fixture
                        tree and reds when a moved step owns no fixture.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | `gofmt` phase reds when a file is unformatted, green when none are | gate-go subcommand against temp trees | `go test ./internal/gate -run TestGateGoGofmt` observed red before build | a phase that shells `gofmt -l` directly exits 0 with output and can never red; the row demands exit 1 plus the named file |
| 1 | gofmt materializes as a phase in the probed table | BenchkitPhases against a probed temp root | `TestPhaseTableProbedToolchainPhases` observed red before build | without the table entry no gofmt phase exists and the asserted name is absent from the resolved table |
| 2 | `vet` phase reds on a vet diagnostic | gate-phases against a temp root with a vet-triggering file | `TestPhasesCommandVetPhaseReds` observed red before build | a table that omits vet, or names it with wrong argv, leaves the run green on a tree `go vet` rejects |
| 3 | the test phase's package set excludes contract, conformance, and release-only at dev, and includes release-only at ship | gate-go test package enumeration, called directly | `TestGateGoTestPackageSet` observed red before build (both tier subcases) | a naive `go test ./...` includes the excluded packages, restaging the recursion and the 600 s timeout; the ship subcase catches an exclusion that forgot the tier |
| 3 | the test phase reds when a core package's test fails | gate-go subcommand against a temp tree with a failing test | `TestGateGoTestReds` observed red before build | an enumeration-only implementation that never execs `go test` exits 0 on a tree whose tests fail |
| 4 | the race step is absent from the probed table for a go.mod root | BenchkitPhases against a probed temp root | `TestPhaseTableProbedToolchainPhases` sibling assertion, observed red before build | today's behavior is the overreach; the row fails while the race step is still probed for on arbitrary roots |
| 4 | the race step is present in the kit's manifest and reds when its target test does not run | gate-go race against a tree with no `internal/worktree` | `TestGateGoRaceRequiresTheTestToRun` observed red before build | a `-run` filter matching nothing exits 0; the did-it-run guard is the only thing that separates "passed" from "never executed" |
| 5 | the filtered conformance suite runs as its own kit phase, with the registry's skip pattern reaching it from Go rather than from the manifest text | `gate-go conformance-suite` argv assertion plus the kit manifest read through `phaseTable` | `TestGateGoConformanceSuiteUsesRegistrySkipPattern` observed red before build, asserting the argv against `registry.InnerSkipPattern()` and asserting the manifest carries no regex literal | a manifest with the pattern retyped into JSON satisfies any phase-presence test while giving the non-recursion contract a second source that drifts silently on the next rename |
| 6 | no temp build output path is created and `dist/bench` is written only by the build phase | the existing `TestCheckGoCoreDoesNotWriteRootDistBench`, rewritten for the residual check | the existing test asserts the throwaway path is used, so it is inverted to assert no build probe runs at all; observed red before the deletion | a deletion that leaves the throwaway build behind passes every other test in the package; the inverted assertion is what notices the step is still there |
| 7 | the residual check still emits every structural diagnostic it owns | the 28 existing package-core-guard fixtures | already covered — the sweep's existing EXPECTs | a subtraction that removes too much turns those fixtures' EXPECTs vacuous and the sweep reports them |
| 7 | go.mod present with no `go` on PATH still reds with the toolchain diagnostic | residual check with a temp root and a stripped PATH | `TestResidualCheckReportsAbsentToolchain` observed red before build | if the diagnostic leaves with the steps that used it, a host without Go grades green on a tree whose compiled core is load-bearing |
| 8 | a root with `go.mod` but no `go` on PATH materializes no toolchain phases | BenchkitPhases with a stripped PATH | `TestPhaseTableNoToolchainNoPhases` observed red before build | an unprobed table emits phases that fail with exec errors, turning one attributable diagnostic into three unattributable phase reds |
| 8 | a root with no `go.mod` materializes no toolchain phases | BenchkitPhases against a non-Go temp root | same test, sibling assertion, observed red before build | a table that probes only PATH adds Go phases to a Python repo, reding a linked repo the gate never used to touch |
| 8 | a `go.mod` that is a dangling symlink or a special file is not treated as a present Go root | BenchkitPhases against temp roots with a broken link and a FIFO | `TestPhaseTableProbeRejectsNonRegularGoMod` observed red before build (test timeout is the FIFO hang tripwire) | a plain `os.Stat`-and-open probe blocks on the FIFO and a plain read misclassifies the broken link; the regular-file predicate is what separates them |
| 9 | the kit's manifest resolves to a table equal, field by field, to the `BenchkitPhases`-derived table for the kit root plus exactly the two kit-only phases | `phaseTable` against the kit root, compared against `BenchkitPhases(kit, kit)` | `TestKitManifestMatchesDerivedTable` observed red before build, diffing name, argv, env, needs, optional, and dir | a name-presence test passes while the manifest drops an env pin or mangles an argv, so the gate runs a phase that grades nothing; the field-by-field diff is what makes a silent weakening visible |
| 10 | the new toolchain phases carry no `needs` edge | BenchkitPhases and kit manifest table | `TestPhaseTableProbedToolchainPhases` Needs assertion, observed red before build | copying the build edge from the existing phases costs the whole overlap win with no test noticing |
| 10 | no package in the test phase's set consumes `dist/bench` | gate-go test package enumeration, asserting `internal/contract` and its subpackages are absent | covered by the story-3 package-set row's dev subcase | if `internal/contract` re-enters the set, a build-edge-free test phase races the build phase for the binary; this is the assertion that makes the no-edge decision safe |
| 11 | an unknown `gate-go` step is a usage error, not a silent success | gate-go subcommand | `TestGateGoUnknownStep` observed red before build | a default-case fallthrough returns 0, so a typo'd step in a manifest reds nothing and grades nothing |
| 11 | `gate-go` is reachable through the shipped `bin/bench.sh` route and recorded in the routing registry | subcommand-routing conformance check | `TestRootConformance` subcommand-routing check observed red before the registry entry is added | the check reads dispatch names from the file itself, so a new subcommand that is not registered reds the gate — this row is a tripwire the build must satisfy, not optional |
| 11 | `gate-go` works against a root whose path contains a space | gate-go subcommand with a spaced temp root | `TestGateGoSpacedRoot` observed red before build | an argv assembled by string concatenation splits the path and the step grades the wrong tree, or nothing |
| 12 | a fixture in a phase-named family routes to that phase, not to `conformance` | SweepTier with a fake Runner over a temp fixture tree, asserting `RunCall.Env`'s phase value | `TestSweepRoutesPhaseNamedFamilies` observed red before build | `FixturePhase`'s default arm returns `conformance` for any unrecognized family, so an unwidened build routes a gofmt fixture into the conformance phase, which never emits its EXPECT — the assertion is on the routed phase name, not merely on the fixture failing |
| 12 | `phasesForMode` keeps a phase-named owner instead of clearing it | phasesForMode directly, with each new phase name | `TestPhasesForModeAcceptsPhaseTableNames` observed red before build | the current allowlist clears any owner outside conformance/contract, so the filter runs every non-canary phase and the fixture's targeted red drowns in unrelated work |
| 13 | the ship tier still reaches the cross-compile matrix after the split | residual check under `-tags stress`, asserting a cross-compile diagnostic for a root that cannot cross-compile | `go test -tags stress ./internal/conformance -run TestResidualCheckKeepsCrossCompile` observed red with the call removed, before build | cross-compile is a dev-tier no-op with zero fixtures, so a split that drops the call from the residual leaves every other row green while ship silently loses the four-platform matrix |
| 14 | `go-build-broken` bites at the build phase | the real canary sweep | fixture EXPECT observed red in the sweep during build (craft-gate recorded red) | the fixture's tree lacks the build probe's inputs today, so an unmodified move makes it green — the recorded red is the only thing that proves otherwise |
| 14 | `go-test-failing` bites at the test phase | the real canary sweep | fixture EXPECT observed red in the sweep during build (craft-gate recorded red) | its EXPECT is `formatProbeFailure` framing no standalone phase emits; without a new EXPECT the fixture passes vacuously |
| 15 | gofmt, vet, race, and the filtered conformance suite each own a fixture that bites at its destination phase | the real canary sweep | each new fixture's EXPECT observed red in the sweep during build (craft-gate recorded red) | these four steps have zero fixtures today, so nothing else in the gate would notice if a restaged step stopped grading anything |
| 15 | each kit-manifest-only phase is canaried through a fixture-local manifest | the real canary sweep, race and conformance-suite fixtures | both fixtures' EXPECTs observed red in the sweep during build | a fixture root has no kit manifest, so without its own neither phase materializes and the fixture grades the empty set — the failure mode is a green fixture, which only a recorded red rules out |
| 16 | every moved step has a fixture bound to its destination phase, the filtered conformance suite included | parity inventory test in conformance | `TestEveryMovedStepOwnsAFixture` observed red before the fixtures are authored, grading all five stepless phases red | a static step→owner mapping test stays green while a migrated EXPECT matches nothing; reading the fixture tree is what makes an orphaned step visible, and enumerating every moved step is what stops the inventory quietly exempting the one phase nobody wrote a fixture for |
| 17 | the probe-formatter label scan matches the labels that remain | `TestCoreSubprocessFailuresUseProbeFormatter` | the existing test observed red the moment a label leaves `checkGoCore` | this is the coupled assertion the map names; it fails immediately and mechanically, which is the intended behavior, not a defect to route around |
| 18 | `prep-release` still runs the release-only package suites | `Steps(root, kit)` unit test | `TestStepsRunReleaseOnlyPackages` observed red before build | after the split nothing else reaches those three suites; without this step ship green loses them and no dev-tier test notices |
| 19 | an unbound conformance family errors before any inner gate runs, on both entry points | SweepTier with a fake Runner, asserting zero `RunCall`s | `TestSweepRefusesUnboundFamilyBeforeAnyRun` observed red before build | today's red lands in the conformance layer after the sweep already paid for full unscoped runs, and standalone `bench canary` never reaches it at all — asserting zero RunCalls is what pins "before the cost" |
| 19 | an adopting repo's unbound families do not error | SweepTier with root ≠ `BENCH_KIT` | same test, sibling assertion, observed red before build | an unscoped assertion breaks every linked repo, whose families a kit-owned table can never carry — this is the fail-posture-per-audience trap FT144 was opened for |
| 19 | the swept root is compared to `BENCH_KIT` through resolved symlinks, so a symlinked kit checkout still asserts | SweepTier with root reached via a symlink and `BENCH_KIT` set to the physical path | `TestSweepBindingAssertionResolvesSymlinks` observed red before build | a raw string compare makes the kit look like an adopting repo and skips its own assertion silently — the one failure mode that leaves FT143 open while looking fixed |
| 19 | the unbound-family scan has one derivation serving both callers | conformance check delegating to the canary package | the existing `internal/conformance/fixture_bite_test.go` unbound-family unit test (which runs in the filtered conformance suite, not under `TestRootConformance`) stays red for a synthetic unbound family through the move, observed | two copies drift, and the copy that drifts is the one nobody runs; keeping that existing test red across the refactor is what proves the conformance caller still reaches the relocated scan |
| 20 | the manifest format is documented in the file that ships to consumers | docs-currency conformance check over `.bench/BENCH-reference.md` | `TestRootConformance` docs-currency check observed red if the new plumbing subcommand is absent from the reference enumeration | the reference's plumbing list is a single-source enumeration the gate grades, so an undocumented `gate-go` reds; the prose describing the manifest schema is graded by review, which this row states honestly rather than claiming coverage it lacks |
| 21 | summary lines, inner byte shape, timing lines, and exit codes are unchanged | existing runner, timing, and canary suites | already covered — `TestRunnerSummaryLineByteShape`, `TestRunnerInnerModeByteShape`, `TestRunnerPrintsConformanceTiming`, plus every existing fixture EXPECT | these pin the exact strings; adding five phases to the table cannot drift them without reding at least one |
| 22 | whole-gate and per-phase wall clock improve against the 4m36s baseline | not TDD-able — the gate cannot grade its own duration without becoming flaky on a loaded host | not TDD-able: ship evidence is the manual post-change measurement the map names | a timing assertion would red on an unrelated busy host and get weakened; the map records this as the one outcome outside the gate |

### Edge inventory

Walked classes landing as rows above: absent vs present-but-empty (`go.mod`
absent, present-but-no-toolchain, and empty gofmt output — story 7/8 rows);
dangling symlink and special file where a file is expected (story 8's probe row,
with the FIFO hang tripwire); required tool missing from PATH (story 7/8 rows);
paths containing spaces (story 11 row); a flag value read as a positional
(story 11's unknown-step row, which pins exit 2 rather than a fallthrough);
invocation through every shipped surface (story 11's routing-registry row);
invocation through a symlink (story 19's symlink-resolution row, which is where
this class actually bites in this slice); and
the kit-versus-linked-repo audience split, which is walked explicitly for every
new fail posture in this slice — the race phase (story 4), the probed fallback
table (story 8), and the FT143 assertion's scoping (story 19). That audience
question is FT144's proposed `craft-spec` prompt, applied here ahead of the kit
edit that would make it standing.

Won't handle:
- **Control bytes in tool output** — `gofmt`, `go vet`, and `go test` output
  reaches the phase's stdout unfiltered today via the runner's prefix writers, and
  this slice does not change that path; the conformance driver's sanitizer went
  with the steps it framed.
- **`gofmt -l` output whose last line lacks a trailing newline** — the file list
  is parsed with `strings.Fields`, which has no line semantics.
- **`.bench/phases.json` leaking into the consumer payload** — the payload
  allowlist is a complete positive enumeration, so an unlisted file is
  undistributed by construction; `.bench/gate.sh` is the standing precedent for a
  kit-only `.bench` file.
- **Concurrent `gate-go` invocations against one root** — the gate already
  refuses a subject that changes mid-run, and the steps are read-only apart from
  the Go build cache, which is concurrency-safe by design.
- **A manifest phase named `gofmt`/`vet`/`test` in an adopting repo colliding
  with the probed fallback** — a present manifest replaces the built-in table
  entirely; there is no merge, so there is nothing to collide.
- **Interrupt mid-phase for the new phases** — the scheduler's cancellation and
  process-group teardown are slice B's, unchanged and already covered.
- **A globally installed `bench` run inside the kit checkout** — `BENCH_KIT`
  honestly names the install directory, so the FT143 assertion stays silent.
  That is the correct answer, not a miss: the running kit does not own the swept
  tree's table, and the kit's own gate always runs its own `bin/bench.sh`.
- **Glob characters in the graded root's path** — the phases pass root through
  argv arrays with no shell between them, so the spaced-root row is the whole of
  the exposure; a glob character is inert in an exec array.

## Out of scope

- **Ship-tier phases joining the manifest** (`prep-release` as a manifest
  consumer) — a separate capability: this slice keeps `prep-release`'s own step
  list as its ordering source and only preserves parity. On the map's "Not yet
  specified" list, so it needs a decision before a spec. Estimate: 6 edits, 3
  gate runs plus a ship run.
- **`bench upgrade` semantics for a project-owned manifest** — how linked repos
  receive runner upgrades once phase definitions are theirs. A distinct
  capability (asset lifecycle, not gate execution) and undecided on the map.
  Estimate: 8 edits, 3 gate runs.
- **A cross-language capability-skip surface** — a documented skip-line contract
  or a `bench skip` helper a non-Go phase can exec. #8's recorded gap; a new
  contract, not the rest of this feature. Estimate: 6 edits, 3 gate runs.
- **Cross-compile as a phase** — stays ship-tier inside the residual check, per
  the map's #3 answer. Not a cut this slice is choosing; a decision already
  closed.
- **FT142's remaining residuals and FT116's `enumerateGuards` leak** — board rows
  with their own shapes. This slice touches FT142 only where ship parity forces
  it (story 18), and that overlap is flagged above.
- **FT144's workflow decision** — whether "intent stands, seam moves" gets a
  named lighter path. The reviewer's call, unmade; this spec applies the
  edge-inventory half of FT144 without waiting on it.
