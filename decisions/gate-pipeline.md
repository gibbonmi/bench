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
Resolved 2026-07-26.

**v1 fields — six:** `name`; `argv` as an exec array (no shell strings, the
existing Phase contract); `env` pairs; `needs` as dependency edges by phase
name; `optional` for skip-if-binary-absent; `dir` — phase working directory
relative to the graded root, default the root, validated to stay inside it
(added by #8: regroup-app's frontend phases forced it). `Serial` disappears —
it is the degenerate edge build→everything, expressed via `needs`. Timeout stays a
runner concern pending #4; tier stays out (ship-tier manifest membership is
declared fog); no weight field without a scheduler that reads it.

**Format and location:** JSON at `.bench/phases.json`, parsed with
`encoding/json` — no new dependency, and `.bench/` is already the
project-owned seam `bench upgrade` respects.

**Entry:** `.bench/gate.sh` stays the hand-ownable oracle entry (it already
just execs `bench gate-phases "$root"`); `gate-phases` gains manifest loading
behind it. The resolution chain is untouched. Generated shim rejected — a
write-owner problem for zero gain.

**Fail posture:** absent file → the built-in kit table (today's behavior;
linked repos keep working through the migration). Empty — blank or zero
phases → red: deleting the file is what asks for the default, so empty is
likelier a truncated write than an intent (the empty-CHECK-file precedent).
Malformed — parse error, duplicate name, dangling or cyclic `needs`, empty
argv → red naming the defect; falling back to the built-in table would
silently grade with the wrong oracle.

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
Extension, not rewrite. Survives as-is: the `Phase` model
(`Name/Argv/Env/Optional/Serial`, `internal/gate/phases.go:40`) with
process-group launch, cancel/timeout codes 130/124, optional-binary skip, and
per-phase prefixed output + `phase X: green/red` summaries (`runner.go`); the
subject/lock/verdict machinery wraps the phase run as a black box
(`gate.go:234`, phase-agnostic); capability-skip collection is already
per-phase env injection (`capability_skips.go:57`); the conformance timing
print keys on phase name only (`runner.go:249,287`).

Compiled-in assumptions the manifest replaces: `BenchkitPhases(root, kit)`
builds the table in Go with kit knowledge — go-build.sh probe, shellcheck file
list, canary argv (`phases.go:58–149`); the only injection seam is the
`benchkitPhasesForCommand` var (`phases.go:48`). `.bench/gate.sh` already just
execs `bench gate-phases "$root"` — the entry survives; the manifest loads
behind `gate-phases`.

Single-stage-shape assumption a DAG replaces: `splitSerialPhases`
(`runner.go:201`) is the one ordering source for both runners — serial phases
fail-fast in table order, then all concurrent phases run-all-and-report.
Dependency edges attach by replacing that split with a scheduler; `Serial:
true` is today's degenerate edge (build → everything). Inner-mode filtering
(`phasesForMode`, `phases.go:271`) selects phases by name via
`BENCH_CANARY_PHASE`, so phases must stay name-addressable.

Conformance side (for #3/#5): `RunConformance` loops `registry.Checks`
filtered only by tier (`checks_test.go:62`) — no per-check filter exists;
narrowing today is tier env + phase env only. `checkGoCore`
(`package_core_checks_test.go:180`) runs gofmt → build (throwaway path) → vet
→ test (excluding contract/conformance/release-only) → filtered inner
conformance test → worktree race test → cross-compile serially inside the
`package-core-guard` check.

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
Resolved 2026-07-26.

**Phase set:** gofmt, vet, test, and the worktree race test leave the
composite as first-class phases. The throwaway build-validation step
collapses into the gate's existing serial build phase — one build, one owner;
the throwaway copy existed only because conformance couldn't own
`dist/bench`. Cross-compile stays ship-tier, untouched.

**Residual:** `package-core-guard` keeps its structural siblings — package
files, release/native workflow, preflight, identity strings, guard manifests,
resolver drift — plus the "go.mod present but no toolchain" diagnostic. All
grade any root cheaply.

**Layering:** the built-in fallback table gains go.mod-probed gofmt/vet/test
phases (the build-phase probe pattern), so manifest-less roots — canary
fixtures, linked Go repos — keep toolchain grading and fixture EXPECTs still
bite. The kit-specific steps — the worktree race test and the filtered
`internal/conformance` suite run — are declared only in the kit's
`.bench/phases.json`. This also fixes an overreach: today's race test runs
against any go.mod root and would spuriously red a linked repo without
`internal/worktree`.

**Test-phase enumeration:** the phase's argv invokes a bench plumbing
subcommand that runs `go list`, applies the exclusion policy (today's
`isExcludedTestPackage`), and execs `go test` — policy single-sourced in Go,
manifest stays declarative, works against arbitrary roots.

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
Resolved 2026-07-26.

**Width budget:** unweighted — new phases rely on the OS scheduler exactly
like today's concurrent phases; nested canary inner gates keep their
`GOMAXPROCS=2` pin and worker derivation unchanged. The concurrency arm's
evidence rules: after the canary cap, load peaks ~2× cores only in bursts and
wall is conformance-bound, so outer capping buys nothing and could cost. No
weight field (per #1); revive trigger stays persistent contention flakes.

**Red posture:** a red phase marks its downstream dependents
skipped-with-cause — reported distinctly, never as red — since they would
grade an artifact that failed to materialize (today's serial fail-fast
rationale). Phases with no path from the red one run to completion and
report, preserving the full red-set signal a fix loop needs. Today's two
behaviors, generalized by the edges. Full fail-fast and run-everything both
rejected: the first destroys the shrinking-red-set signal, the second emits
cascade reds that read as independent defects.

**Timeouts:** the single gate-level deadline (`bounds.GateTimeout`, exit 124)
stays the only timeout; when it fires, the runner names the phases still
running so a hang is attributed. Both 600 s go-test incidents hurt through
diagnosis, not the wait. Per-phase deadlines — runner default or manifest
field — stay behind evidence they'd catch something the gate deadline
doesn't.

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
Resolved 2026-07-26, buildable on today's gate as the prerequisite slice.

**Scope source:** a registry-owned family→check table, extending the
`FixturePhase` family→phase convention; a fixture-level `CHECK` file
overrides, as the two ship fixtures already do. Family membership stays the
binding and `CHECK` stays written only where it changes the answer. #6
verified the mapping per emitting code: it holds many-to-one, and the nine
stray fixtures carry CHECK overrides (see #6).

**Transport:** a new `BENCH_CONFORMANCE_CHECK` env read at the conformance
entry point, filtering `RunConformance`'s registry loop; the canary sets it
per fixture scrub-then-set (the `ConformanceTierEnv` precedent, owned by
`innerEnv`). Absent env = full-tier run, so behavior-owned, legacy flat, and
unmapped-by-design fixtures stay correct by default.

**Fail posture — all loud:** unknown check name → red diagnostic at the
conformance entry; a scope the running tier does not execute → red (a
zero-check scoped run would otherwise surface as a baffling "did not bite");
a conformance family missing from the family→check map → sweep error,
matching `fixtureTier`'s unknown-CHECK posture. No silent fallback to a full
run — that hides drift and re-pays the cost this slice removes.

**Dedup:** one shared scoped empty-tree baseline per check group, run only
for groups present in the tier's sweep; each EXPECT is vacuity-checked
against its own group's baseline. Unscoped fixtures keep today's full
baseline. Sound because the guard's premise is "would EXPECT match with no
mutation under the same run shape the fixture pays" — scoped-vs-scoped is
the consistent comparison. Ship drops from three full inner gates to three
scoped probe runs. Merging same-check fixture runs was rejected: each
fixture is a distinct mutated tree and must be graded alone.

## #6: What do the package-core-guard fixtures migrate to?

Blocked by: #3
Type: Research

### Question
39 fixtures assert the composite check's diagnostics; after the split they
must grade the phase that owns each diagnostic. Inventory what each fixture's
EXPECT strings actually assert, bucket them by destination phase versus
residual structural check, and name any fixture whose expectation no seam
will emit after the split — those need a decision, not a mechanical move.
Also carries #5's assumption to verify: for each of the nine conformance
families, confirm every fixture actually grades the family's mapped check, so
the family→check table holds; name any fixture that doesn't. Asset: the
bucketed inventory. The #7 parity requirement reads off it: every moved step
must end with ≥1 fixture biting at its destination phase.

### Answer
Resolved 2026-07-26; full inventory in
`decisions/assets/gate-pipeline-fixture-inventory.md` (per-claim citations,
spot-checked). Headlines: of 39 fixtures, 28 stay with the structural
residual, 2 are ship, 7 are strays belonging to six other checks, and only 2
move to new phases — and both die on the move (their EXPECTs are
`formatProbeFailure` framing no standalone phase emits;
`go-build-broken`'s tree also lacks the build-phase predicate,
`scripts/go-build.sh`). gofmt, vet, and race have zero fixtures today, so
the split creates uncanaried seams unless fixtures are authored in the same
pass. `TestCoreSubprocessFailuresUseProbeFormatter` text-scans the label
inventory and breaks the moment a label leaves — mechanical migration
casualty. Family→check verification: mapping holds per emitter, but
many-to-one (three families share `docs-currency-workflow`) and
`compliance-hardening` binds to `canary-inner-compliance`, not the
same-named kit check.

Decisions taken on the findings: **binding is many-to-one** — scoping needs
no partition (that is tier's job) — and the 7 stray fixtures without CHECK
files gain them (the decided override mechanism; no directory moves).
**Fixtures reach new phases via phase-named top-level families**
(`tests/canary/build|gofmt|vet|test|race/`), extending the directory
convention `FixturePhase` already implements, with `phasesForMode`'s owner
allowlist widened to the phase table's names; the two movers get
phase-native EXPECTs (tool diagnostics plus the stable `phase <name>: red`
summary) and `go-build-broken`'s tree gains what the build predicate needs;
the missing gofmt/vet/race/test fixtures are **authored in the same pass** —
#7's parity requirement already demands one biting fixture per moved step.
Rejected: a per-fixture PHASE file (second binding file where the directory
carries the fact); deferring the new fixtures (ships uncanaried phases).

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
Resolved 2026-07-26. Fixture-backed parity: the #6 migration must end with
every moved step owning at least one canary fixture that bites at its
destination phase, and an inventory test asserts each moved step has a
fixture bound to its new owner — the `TestRegistryBindsEveryCheck` analogue,
grounded in fixtures that prove behavior rather than bookkeeping. The
migration's ship evidence is the recorded red per moved step (per
`craft-gate`), demonstrated during the #6 fixture migration. Rejected: a
static step→owner mapping test alone (stays green while a migrated EXPECT
matches nothing any seam emits) and recorded reds alone (no tripwire after
the build).

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
Resolved 2026-07-26 against `~/workspace/regroup-app` (Python backend via uv
+ TS frontend). Its gate phases: `uv run ruff check`, `uv run mypy`,
`uv run pytest` at the root, all independent; `npm run build`
(`tsc --noEmit && vite build`), `npm run test` (vitest), `npm run test:e2e`
(Playwright, needs the built dist — a real `needs` edge — plus env like
`REGROUP_E2E_REUSE_DIST=1`), `npm run lint:design` — all run from
`frontend/`. Composite package scripts stay composite; the manifest points
at them. Service lifecycle (the e2e backend server) stays inside the phase
command — the runner never manages services.

**One schema change forced: a `dir` field** — phase working directory
relative to the graded root, default the root, validated to stay inside it.
The runner hardcodes cwd=root and argv is an exec array (no shell, no `cd`);
the alternatives are tool-specific luck (`npm --prefix`) or a wrapper script
per phase. Adopted into v1 (see #1). **One gap recorded as fog, not v1:**
only the kit's Go helpers can write capability-skip rows, so a non-Go phase
(Playwright browsers missing) has no structured skipped-not-green surface;
`optional` covers absent binaries, everything else is red-on-missing, and
regroup's own `dev-preflight.sh` diagnoses the host. Everything else in the
five-field schema held.

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
  consumer).
- How linked repos receive runner upgrades once phase definitions are theirs
  (`bench upgrade` semantics for a project-owned manifest).
- A cross-language capability-skip surface — a documented skip-line contract
  or a `bench skip` helper a non-Go phase can exec to report
  skipped-not-green (#8's gap; `optional` covers absent binaries today).

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
- Removing canary nesting — clause transferred to
  `decisions/gate-critical-path.md`, where it was
  reopened and ruled
  (2026-07-28): behavior-owned bites move to the owning contract test.
- The two interim defects on the FT91 row (`BENCH_CONFORMANCE_TIER` scrub
  symmetry, probe-output spill) — shape already decided there; build work,
  not map fog.

## Handoff

1. **Module boundaries.** `internal/gate` owns the manifest loader/validator,
   the DAG scheduler that replaces `splitSerialPhases`, the built-in fallback
   table (gaining go.mod-probed gofmt/vet/test phases), and
   straggler-naming at the gate deadline. `internal/conformance` owns the
   residual `package-core-guard` check and the check-filter in
   `RunConformance`. `internal/conformance/registry` owns the family→check
   table and the filter-env constant. `internal/canary` owns scope
   resolution (family/CHECK → env), per-check shared vacuity baselines, and
   the widened family→phase routing. A new bench plumbing subcommand owns
   test-phase package enumeration. The kit's `.bench/phases.json` declares
   the kit-specific phases (worktree race test, filtered conformance suite).
   `internal/bounds` and the shift/verdict machinery are outside — untouched.
2. **Contracts.** Manifest: JSON at `.bench/phases.json`, six fields
   (`name`, `argv` exec-array, `env`, `needs`, `optional`, `dir`
   default-root inside-root); absent → built-in table, empty or malformed
   (parse error, duplicate name, dangling/cyclic needs, empty argv,
   escaping dir) → red naming the defect. Scoping env
   `BENCH_CONFORMANCE_CHECK`: absent → full tier; unknown name or
   tier-mismatch → red. Runner output: today's `phase <name>: green|red`
   summary lines and timing formats unchanged (#9); dependents of a red
   phase report skipped-with-cause, distinct from red; exit codes
   0/1/3/124/130 unchanged; on deadline the runner names still-running
   phases. `.bench/gate.sh` stays the one-line exec entry.
3. **Deep vs thin.** Deep: the scheduler (hides DAG execution and
   cancellation), the manifest loader (hides the validation family), the
   canary scope resolver + baseline grouping, the test-enumeration
   subcommand (owns exclusion policy). Thin: `gate.sh`, the manifest itself
   (data), the phase-named family convention.
4. **Black-box assertables.** Loader: exit code + stderr diagnostic per
   malformed-manifest class from fixture files. Scheduler (fake runner):
   needs-respecting order, dependents-cancelled-with-cause, independents
   complete, straggler naming. Canary (fake `Runner`): each fixture's
   `RunCall.Env` carries the resolved check var; baseline run count equals
   the number of check groups present. Parity: the inventory test asserting
   every moved step has a fixture bound to its destination phase; each new
   fixture's recorded red. Conformance filter: scoped run executes exactly
   the named check (timing lines observable).
5. **Gate attachment.** The kit's own suites run inside the gate's
   conformance/test phases, so the gate sees every seam above; the canary
   sweep itself proves fixture routing and scoping bite. Not
   gate-assertable: the wall-clock outcome — ship evidence is the manual
   post-change measurement against the 4m36s / `package-core-guard` ~86 s
   baseline.
6. **Hostile-input owners.** Malformed-manifest family → loader.
   Hostile `BENCH_CONFORMANCE_CHECK` values → conformance entry (unknown /
   tier-mismatch reds). Empty or unknown CHECK file → existing
   `fixtureTier` posture; unmapped conformance family → sweep error.
   `dir` escaping the root → loader validation. Duplicate env keys →
   strip-then-set everywhere (see watch-outs).
7. **Uncertainty flags.** None blocking. The exact EXPECT strings for the
   new phase fixtures are settled at build time within the decided rule
   (tool diagnostic + `phase <name>: red`).
8. **Rejected alternatives.** Manifest: weight field, per-phase timeout
   field, tier field, YAML/TOON, generated gate.sh shim, generator-written
   package lists, static `go test ./...`. Execution: weighted budget,
   full fail-fast, run-everything-on-red, runner per-phase timeout default.
   Canary: single full baseline for scoped runs, merging same-check fixture
   runs, per-fixture PHASE file, relocating stray fixtures, deferring the
   new-phase fixtures. Silent fallback on any malformed input.
9. **Domain watch-outs.** The kit's own test suite runs nested inside the
   gate's conformance phase with the phase env passed through, and Go's exec
   env has no duplicate-key precedence — every inner env pin must
   strip-then-set, never append. The gate's build phase materializes only
   when both `scripts/go-build.sh` and `go.mod` exist in the graded tree — a
   fixture targeting it must ship both. Conformance-family fixtures are
   routed by `FixturePhase` and `phasesForMode`'s owner allowlist; a phase
   name absent from both runs every non-canary phase, so fixture migration
   and allowlist widening must land together.
   `TestCoreSubprocessFailuresUseProbeFormatter` text-scans the composite's
   label inventory and goes red when any label leaves it.

Dependency order: confirmed by the reviewer 2026-07-27. Slice A (the canary
check-scoping prerequisite, #5) shipped 2026-07-26 (spec
`ft91-canary-check-scoping`, retired). Next is slice B, the manifest + DAG
runner (#1, #4, #9); then slice C, the `checkGoCore` split with the fixture
migration and parity test (#3, #6, #7) — C's spec also carries FT143's cheap
kit-root family→check binding assertion (reviewer-bundled 2026-07-27), so the
registry rework cannot reintroduce the late red on either entry point. The
two interim FT91 defects (conformance env-scrub symmetry, probe-output spill)
land ahead of both slices as a reviewer-approved direct fix-and-gate pass.
