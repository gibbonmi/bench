# Gate pipeline (FT91, pipeline arm)

Status: shaping

## Destination

The dev gate becomes a local pipeline of first-class phases with declared
dependencies. `checkGoCore`'s seven serial toolchain steps overlap instead of
running inside one test function. Conformance grades structure only. Every
canary fixture pays for only the check its `CHECK` file names. Phase
definitions become project-owned data (the runner ships in the kit, shaped
against regroup-app). This cuts dev wall-clock — today 4m36s, long pole
`package-core-guard` ~86 s — without changing what green means.

## #1: What does the phase manifest declare, and who owns its failure modes?

Blocked by: none
Type: Grill

### Question
The phase list must move from Go code (`BenchkitPhases`) to project-owned
declarative data. Open questions cover three things. First, the schema: name, command,
dependencies, tier, timeout, serial-vs-concurrent, and capability
requirements — which of these are v1 fields, and which stay runner
defaults? Second, the format and location: a file beside `.bench/gate.sh`,
and whether `gate.sh` stays the entry that invokes the runner or becomes a
generated shim. Third, the fail posture for the absent, empty, or malformed
manifest family. Absent falls back to the built-in kit table, but
empty-versus-malformed is exactly the class this repo keeps meeting, and it
must be decided, not defaulted.

### Answer
Resolved 2026-07-26.

**v1 fields — six total:**

- `name`.
- `argv` as an exec array (no shell strings, the existing Phase contract).
- `env` pairs.
- `needs` as dependency edges by phase name.
- `optional` for skip-if-binary-absent.
- `dir` — phase working directory relative to the graded root, default the
  root, validated to stay inside it (added by #8: regroup-app's frontend
  phases forced it).

`Serial` disappears. It is the degenerate edge build→everything, expressed
via `needs`. Timeout stays a runner concern pending #4. Tier stays out
(ship-tier manifest membership is declared fog). No weight field exists
without a scheduler that reads it.

**Format and location:** JSON at `.bench/phases.json`, parsed with
`encoding/json` — no new dependency, and `.bench/` is already the
project-owned seam `bench upgrade` respects.

**Entry:** `.bench/gate.sh` stays the hand-ownable oracle entry (it already
just execs `bench gate-phases "$root"`); `gate-phases` gains manifest loading
behind it. The resolution chain is untouched. Generated shim rejected — a
write-owner problem for zero gain.

**Fail posture:** absent file → the built-in kit table (today's behavior;
linked repos keep working through the migration). Empty — blank or zero
phases — reds. Deleting the file is what asks for the default, so empty is
likelier a truncated write than an intent (the empty-CHECK-file precedent).
Malformed — parse error, duplicate name, dangling or cyclic `needs`, empty
argv — reds naming the defect. Falling back to the built-in table would
silently grade with the wrong oracle.

## #2: How much of `internal/gate`'s existing runner survives as the pipeline runner?

Blocked by: none
Type: Research

### Question
The runner already owns a `Phase` struct, a serial/concurrent split
(`splitSerialPhases`), capability skips, per-phase timing, and subject
binding. The pipeline may be an extension (add dependency edges and manifest
loading) rather than a rewrite.

Read `internal/gate` (`phases.go`, `runner.go`, `capability_skips.go`,
`subject.go`) and map four things. Which parts already implement the
phase/verdict/timing model the FT91 row calls universal? What assumes the
phase table is compiled-in? What assumes exactly one serial-then-concurrent
stage shape? Where would dependency edges attach? Asset: a short summary
with per-claim file citations, so #3 and #4 grill against facts rather than
recollection.

### Answer
Extension, not rewrite. Survives as-is: the `Phase` model
(`Name/Argv/Env/Optional/Serial`, `internal/gate/phases.go:40`) with
process-group launch, cancel/timeout codes 130/124, optional-binary skip, and
per-phase prefixed output plus `phase X: green/red` summaries (`runner.go`).
The subject/lock/verdict machinery wraps the phase run as a black box
(`gate.go:234`, phase-agnostic). Capability-skip collection is already
per-phase env injection (`capability_skips.go:57`). The conformance timing
print keys on phase name only (`runner.go:249,287`).

Compiled-in assumptions the manifest replaces: `BenchkitPhases(root, kit)`
builds the table in Go with kit knowledge — go-build.sh probe, shellcheck file
list, canary argv (`phases.go:58–149`). The only injection seam is the
`benchkitPhasesForCommand` var (`phases.go:48`). `.bench/gate.sh` already just
execs `bench gate-phases "$root"` — the entry survives; the manifest loads
behind `gate-phases`.

A DAG replaces the single-stage-shape assumption. `splitSerialPhases`
(`runner.go:201`) is the one ordering source for both runners: serial phases
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
Which of the composite's steps become first-class phases (gofmt / build /
vet / test / race / cross-compile — the last is ship-tier)? What does the
residual conformance check still assert once the toolchain work leaves it?
The complication that stays: `checkGoCore` grades arbitrary roots — canary
fixtures and linked repos. The split must therefore yield a host-repo phase
set plus a narrower structural check that still works against any root. What
that narrower check owns is the decision; the 19+ package-core-guard fixtures
asserting the composite's diagnostics hang off it (#6).

### Answer
Resolved 2026-07-26.

**Phase set:** gofmt, vet, test, and the worktree race test leave the
composite as first-class phases. The throwaway build-validation step
collapses into the gate's existing serial build phase — one build, one
owner. The throwaway copy existed only because conformance couldn't own
`dist/bench`. Cross-compile stays ship-tier, untouched.

**Residual:** `package-core-guard` keeps its structural siblings — package
files, release/native workflow, preflight, identity strings, guard manifests,
resolver drift — plus the "go.mod present but no toolchain" diagnostic. All
grade any root cheaply.

**Layering:** the built-in fallback table gains go.mod-probed gofmt/vet/test
phases (the build-phase probe pattern). This keeps toolchain grading working
for manifest-less roots — canary fixtures, linked Go repos — so fixture
EXPECTs still bite. The kit-specific steps — the worktree race test and the
filtered `internal/conformance` suite run — are declared only in the kit's
`.bench/phases.json`. This also fixes an overreach: today's race test runs
against any go.mod root and would spuriously red a linked repo without
`internal/worktree`.

**Test-phase enumeration:** the phase's argv invokes a bench plumbing
subcommand that runs `go list`, applies the exclusion policy (today's
`isExcludedTestPackage`), and execs `go test`. The policy stays
single-sourced in Go; the manifest stays declarative and works against
arbitrary roots.

## #4: What are the DAG's execution semantics?

Blocked by: #1
Type: Grill

### Question
Width budget: the gate-concurrency arm settled a product-budget model
(`runtime.GOMAXPROCS(0)`, no knob, inner gates pinned). Do overlapping
toolchain phases share that budget by declared weight, or run unweighted and
rely on the OS? How do nested canary inner gates inherit it? Failure
posture: fail-fast on first red phase versus run-all-and-report. Today's
phases already differ from a fix-loop's needs — a red that cancels siblings
saves wall-clock but destroys the shrinking-red-set signal FT107's fifth
clause wants. Timeout ownership: per-phase deadlines in the manifest versus
runner defaults, given the 600 s go-test package default has twice presented
as a gate hang.

### Answer
Resolved 2026-07-26.

**Width budget:** unweighted. New phases rely on the OS scheduler exactly
like today's concurrent phases; nested canary inner gates keep their
`GOMAXPROCS=2` pin and worker derivation unchanged. The concurrency arm's
evidence rules: after the canary cap, load peaks ~2× cores only in bursts,
and wall is conformance-bound. Outer capping buys nothing and could cost. No
weight field (per #1); revive trigger stays persistent contention flakes.

**Red posture:** a red phase marks its downstream dependents
skipped-with-cause, reported distinctly, never as red. They would otherwise
grade an artifact that failed to materialize (today's serial fail-fast
rationale). Phases with no path from the red one run to completion and
report, preserving the full red-set signal a fix loop needs. Today's two
behaviors, generalized by the edges. Full fail-fast and run-everything both
rejected: the first destroys the shrinking-red-set signal, the second emits
cascade reds that read as independent defects.

**Timeouts:** the single gate-level deadline (`bounds.GateTimeout`, exit 124)
stays the only timeout. When it fires, the runner names the phases still
running so a hang is attributed. Both 600 s go-test incidents hurt through
diagnosis, not the wait. Per-phase deadlines — runner default or manifest
field — stay behind evidence they'd catch something the gate deadline
doesn't.

## #5: How does a canary fixture's inner run scope to its named check, and what shape is the dedup?

Blocked by: none
Type: Grill

### Question
This is the prerequisite slice, buildable on today's gate. Inside the
conformance phase, every fixture runs all 16 registry checks to observe the
one its `CHECK` file names. This is cheap at dev, but ruinous at ship (~12
min sweep estimate, rejected).

Two levers need mechanism decisions. Lever (a) scopes the inner run to the
named check: an env the registry already honours, a derived `-run` filter,
or a registry-level filter argument. It also decides what happens when a
fixture names no check.

Lever (b) dedups by check. N fixtures grading one check, plus a shared
vacuity baseline, currently pay N+1 runs. Both ship fixtures grade the same
check, so ship's cost is irreducible by (a) alone. Open question: can the
vacuity baseline be shared per-check without weakening the did-not-bite
verdict?

### Answer
Resolved 2026-07-26, buildable on today's gate as the prerequisite slice.

**Scope source:** a registry-owned family→check table, extending the
`FixturePhase` family→phase convention; a fixture-level `CHECK` file
overrides, as the two ship fixtures already do. Family membership stays the
binding and `CHECK` stays written only where it changes the answer. #6
verified the mapping per emitting code: it holds many-to-one, and the nine
stray fixtures carry CHECK overrides (see #6).

**Transport:** a new `BENCH_CONFORMANCE_CHECK` env read at the conformance
entry point, filtering `RunConformance`'s registry loop. The canary sets it
per fixture scrub-then-set (the `ConformanceTierEnv` precedent, owned by
`innerEnv`). Absent env = full-tier run, so behavior-owned, legacy flat, and
unmapped-by-design fixtures stay correct by default.

**Fail posture — all loud:** an unknown check name reds as a diagnostic at
the conformance entry. A scope the running tier does not execute also reds
— a zero-check scoped run would otherwise surface as a baffling "did not
bite". A conformance family missing from the family→check map reds as a
sweep error, matching `fixtureTier`'s unknown-CHECK posture. No silent
fallback to a full run exists — that would hide drift and re-pay the cost
this slice removes.

**Dedup:** one shared scoped empty-tree baseline per check group, run only
for groups present in the tier's sweep. Each EXPECT is vacuity-checked
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
EXPECT strings actually assert. Bucket them by destination phase versus
residual structural check. Name any fixture whose expectation no seam will
emit after the split — those need a decision, not a mechanical move.

This also carries #5's assumption to verify. For each of the nine
conformance families, confirm every fixture actually grades the family's
mapped check, so the family→check table holds. Name any fixture that
doesn't. Asset: the bucketed inventory. The #7 parity requirement reads off
it: every moved step must end with ≥1 fixture biting at its destination
phase.

### Answer
Resolved 2026-07-26; full inventory in
`decisions/assets/gate-pipeline-fixture-inventory.md` (per-claim citations,
spot-checked). Headlines: of 39 fixtures, 28 stay with the structural
residual, 2 are ship, and 7 are strays belonging to six other checks. Only 2
move to new phases, and both die on the move. Their EXPECTs are
`formatProbeFailure` framing no standalone phase emits, and
`go-build-broken`'s tree also lacks the build-phase predicate,
`scripts/go-build.sh`.

gofmt, vet, and race have zero fixtures today, so the split creates
uncanaried seams unless fixtures are authored in the same pass.
`TestCoreSubprocessFailuresUseProbeFormatter` text-scans the label inventory
and breaks the moment a label leaves — mechanical migration casualty.
Family→check verification: mapping holds per emitter, but many-to-one
(three families share `docs-currency-workflow`) and `compliance-hardening`
binds to `canary-inner-compliance`, not the same-named kit check.

Decisions taken on the findings: **binding is many-to-one**, so scoping
needs no partition (that is tier's job). The 7 stray fixtures without CHECK
files gain them (the decided override mechanism; no directory moves).

**Fixtures reach new phases via phase-named top-level families**
(`tests/canary/build|gofmt|vet|test|race/`), extending the directory
convention `FixturePhase` already implements. `phasesForMode`'s owner
allowlist widens to the phase table's names. The two movers get
phase-native EXPECTs (tool diagnostics plus the stable `phase <name>: red`
summary), and `go-build-broken`'s tree gains what the build predicate needs.
The missing gofmt/vet/race/test fixtures are **authored in the same pass**
— #7's parity requirement already demands one biting fixture per moved
step. Rejected: a per-fixture PHASE file (second binding file where the
directory carries the fact); deferring the new fixtures (ships uncanaried
phases).

## #7: What proves parity — that restaging lost no check?

Blocked by: #3
Type: Grill

### Question
The oracle constraint: green must keep meaning the same thing. The tier split
proved membership with a registry both directions could interrogate
(`TestRegistryBindsEveryCheck`). The pipeline needs the analogue: every
sub-step of the old composite maps to exactly one phase or to the residual
check. This mapping is asserted mechanically, not by review. Decide the
mechanism: a registry-style parity test, a recorded red per moved
step per `craft-gate`, or both. Also decide what the migration's
demonstrated-bite evidence is.

### Answer
Resolved 2026-07-26. Fixture-backed parity: the #6 migration must end with
every moved step owning at least one canary fixture that bites at its
destination phase. An inventory test asserts each moved step has a fixture
bound to its new owner — the `TestRegistryBindsEveryCheck` analogue. It is
grounded in fixtures that prove behavior rather than bookkeeping.

The migration's ship evidence is the recorded red per moved step (per
`craft-gate`), demonstrated during the #6 fixture migration. Rejected: a
static step→owner mapping test alone, since it stays green while a migrated
EXPECT matches nothing any seam emits. Recorded reds alone are rejected too,
since they leave no tripwire after the build.

## #8: Does the manifest hold against regroup-app?

Blocked by: #1
Type: Task

### Question
The FT91 row's decided shaping rule: design the manifest against regroup-app
(the first external, non-Go validation target) rather than generalising from
this repo. The reviewer makes regroup-app's toolchain shape readable (what
its build / test / lint phases actually are). The worker then walks the draft
schema against it on paper and reports where the schema bends. No code — the
deliverable is the list of schema changes regroup-app forces.

### Answer
Resolved 2026-07-26 against `~/workspace/regroup-app` (Python backend via uv
+ TS frontend). Its gate phases at the root are: `uv run ruff check`,
`uv run mypy`, `uv run pytest` — all independent. From `frontend/`:
`npm run build` (`tsc --noEmit && vite build`), `npm run test` (vitest),
`npm run test:e2e` (Playwright, needs the built dist — a real `needs` edge
— plus env like `REGROUP_E2E_REUSE_DIST=1`), and `npm run lint:design`.
Composite package scripts stay composite; the manifest points at them.
Service lifecycle (the e2e backend server) stays inside the phase
command — the runner never manages services.

**One schema change forced: a `dir` field** — phase working directory
relative to the graded root, default the root, validated to stay inside it.
The runner hardcodes cwd=root, and argv is an exec array (no shell, no
`cd`). The alternatives are tool-specific luck (`npm --prefix`) or a
wrapper script per phase. Adopted into v1 (see #1).

**One gap recorded as fog, not v1:** only the kit's Go helpers can write
capability-skip rows. A non-Go phase (Playwright browsers missing)
therefore has no structured skipped-not-green surface. `optional` covers
absent binaries, and everything else is red-on-missing. regroup's own
`dev-preflight.sh` diagnoses the host. Everything else in the five-field
schema held.

## #9: What happens to the gate's timing and output contract?

Blocked by: none
Type: Grill

### Question

Which timing and output behaviors remain stable as the gate pipeline gains
overlapping phases?

### Answer
Resolved at bootstrap — continuity, decided by the tier split. Per-check
timing is permanent gate observability owned by the conformance driver. The
runner keeps emitting per-phase timing in the same stable format (ordering
byte-stable, values free to vary). New overlapping phases each get
a timing line like any phase today. No output-contract redesign rides this
map; any change to the format is its own reviewer decision later.

## Not yet specified

- Ship-tier phases joining the manifest (`prep-release` as a manifest
  consumer).
- How linked repos receive runner upgrades once phase definitions are theirs
  (`bench upgrade` semantics for a project-owned manifest).
- A cross-language capability-skip surface: a documented skip-line contract,
  or a `bench skip` helper a non-Go phase can exec to report
  skipped-not-green. This is #8's gap; `optional` covers absent binaries
  today.

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

## Spec-writer discretion

## Sources
