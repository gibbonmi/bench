# FT91 slice B — phase manifest + DAG runner

Status: staged

Compiled from `decisions/gate-pipeline.md` (#1 manifest, #4 execution semantics,
#9 output contract; Handoff items 1–9). Slice order is reviewer-confirmed
(2026-07-27): this slice ships the mechanism; slice C (`checkGoCore` split,
fixture migration, parity test, probed fallback phases) consumes it.

## Problem

The gate's phase table is compiled into Go (`BenchkitPhases`), so a linked
repo cannot declare its own phases, and the only ordering the runner knows is
one serial-then-concurrent split — no way to express real dependency edges
(regroup-app's e2e-needs-build), no per-phase attribution when the gate
deadline fires, and a red serial phase throws away the run-everything signal a
fix loop needs from unrelated phases.

## Solution

Phase definitions become project-owned data: an optional JSON manifest at
`.bench/phases.json` in the graded root, loaded by `gate-phases` behind the
unchanged `.bench/gate.sh` entry. A DAG scheduler replaces
`splitSerialPhases`: phases start when their `needs` complete, a red phase
marks its downstream dependents skipped-with-cause while independents run to
completion, and at the gate deadline the runner names the phases still
running. Absent manifest keeps today's built-in kit table, so linked repos
work unchanged through the migration. Green keeps meaning the same thing.

## User stories

1. As a linked-repo owner, I want `gate-phases` to load `.bench/phases.json`
   from the graded root when it exists, so that my project declares its own
   gate phases without kit code changes.
   Line: opus (mid) / medium. The profile caches gate logic at the mid tier
   because the loader is oracle code with fully decided semantics.
2. As a linked-repo owner, I want an absent manifest to fall back to the
   built-in kit table, so that every existing repo keeps working through the
   migration with zero action.
   Line: opus (mid) / medium. The fallback is decided behavior at a known
   seam, but it guards the oracle so it stays at the mid tier.
3. As a reviewer, I want a present-but-empty manifest (blank file or zero
   phases) to red the gate naming the defect, so that a truncated write is
   never silently graded with the default oracle.
   Line: opus (mid) / medium. The empty-versus-absent posture is decided in
   the map and the test is mechanical, but it is fail-closed oracle behavior.
4. As a reviewer, I want every malformed-manifest class — parse error,
   duplicate phase name, invalid phase name, dangling `needs`, cyclic
   `needs`, empty `argv`, `dir` escaping the root, unknown field — to red the
   gate with a diagnostic naming the defect and the class, so that a wrong
   manifest never silently grades with the wrong oracle.
   Line: opus (mid) / medium. The class list is decided; the work is a
   table-driven validator with one diagnostic per class.
5. As a reviewer, I want a manifest path that exists but cannot be read as a
   regular file — a dangling symlink or a special file such as a FIFO — to
   red the gate rather than fall back or block, so that a broken link is
   never classified as an authoritative absent state.
   Line: opus (mid) / medium. The hostile-input checklist names both classes
   and the stat-before-read shape is settled kit practice.
6. As a project owner, I want each manifest phase to support the six decided
   fields — `name`, `argv` exec-array, `env` pairs, `needs`, `optional`,
   `dir` — with `dir` resolved relative to the graded root, defaulting to
   the root, and validated to stay inside it, so that regroup-app's
   `frontend/` phases run where their tools expect.
   Line: opus (mid) / medium. The schema is fully decided in the map
   including the regroup-forced `dir` field.
7. As a project owner, I want a phase's `env` pairs applied strip-then-set
   over the gate environment — and the existing built-in phase env
   application converged on the same rule — so that a duplicate key can
   never reach a child process with two values (the map's domain watch-out).
   Line: opus (mid) / medium. One merge helper with a decided rule, applied
   at one call site.
8. As a worker running the gate, I want the scheduler to start a phase only
   after all its `needs` are green and to run phases with no unmet edges
   concurrently and unweighted, so that independent toolchain work overlaps
   exactly as the width-budget decision settled.
   Line: opus (mid) / high. The DAG scheduling and its cancellation paths are
   the one genuinely subtle concurrent piece of this slice.
9. As a worker reading a red gate, I want a red phase's downstream dependents
   reported as skipped-with-cause — distinct from red — while phases with no
   path from the red one run to completion and report, so that the red set I
   see is real defects, not cascade noise, and the shrinking-red-set signal
   survives.
   Line: opus (mid) / high. The red posture is decided but its propagation
   through a graph under concurrency shares story 8's subtlety.
10. As a worker on a host missing an optional binary, I want an optional
    phase's skip to propagate to its dependents as skipped-with-cause, so
    that a phase never grades an artifact its skipped dependency did not
    produce.
    Line: opus (mid) / medium. Small decided extension of the skip posture;
    flagged below as a map-silent default for veto.
11. As a maintainer, I want the `Serial` field gone and the built-in kit
    table to express today's ordering as `needs` edges (every other phase
    needs `build` when the build phase materializes), so that one scheduler
    is the only ordering source and today's gate behavior is preserved
    byte-for-byte where it is pinned.
    Line: opus (mid) / medium. Mechanical conversion guarded by the existing
    runner, timing, and canary suites.
12. As a worker diagnosing a hang, I want the runner to name the phases
    still running when the gate deadline fires — the single gate-level
    deadline staying `gate-run`'s `bounds.GateTimeout`, with `gate-phases`
    printing stragglers on the termination signal before it dies — so that
    a hang is attributed to a phase instead of ending as a silent SIGKILL.
    Line: opus (mid) / high. The signal-cascade change and the straggler
    report are concurrent cancellation code, the subtlest part of the slice.
13. As a canary and contract consumer, I want the runner's output contract
    unchanged — `phase <name>: green|red` summaries, timing line formats,
    prefix framing, inner-mode byte shape, exit codes 0/1/3/124/130 — with
    skipped-with-cause as the only additive line class, so that substring
    EXPECTs and downstream parsers keep biting.
    Line: opus (mid) / medium. Continuity is pinned by existing suites; the
    work is not breaking them.
14. As a reviewer, I want a canary fixture whose tree carries a malformed
    `.bench/phases.json` with an EXPECT on the loader's diagnostic, so that
    the loader's fail-closed posture stays bitten by the sweep permanently,
    not just at build time.
    Line: opus (mid) / medium. One fixture following the existing canary
    conventions; the inner gate reds at load so the fixture is cheap.
15. As the inner-mode canary path, I want manifest-declared tables to flow
    through the existing `phasesForMode` filtering and sequential inner
    runner unchanged (topological order, no summaries, same byte shape), so
    that nested grading semantics do not fork by table source.
    Line: opus (mid) / medium. The filter and inner byte shape are pinned by
    existing tests; the scheduler must simply preserve them.

## Implementation decisions

- **Loader.** New loader in `internal/gate` (the Handoff's module owner):
  `gate-phases` stats `.bench/phases.json` under the graded root before
  reading; a non-regular or unreadable-but-present path is malformed, not
  absent. Absent → `benchkitPhasesForCommand` exactly as today. All loader
  reds print one diagnostic to stderr naming the file, the defect class, and
  the offending element, then exit 1 before any phase runs.
- **Manifest shape** *(map-silent default — flagged for veto)*: top-level
  object `{"phases": [...]}` decoded strictly (unknown fields at either level
  are the unknown-field malformed class). Strict decoding follows the map's
  posture that silent tolerance grades with the wrong oracle; the object
  wrapper leaves room for later top-level fields without a format break.
- **Field semantics.** `name`: required, unique, non-empty, no whitespace or
  control bytes *(charset rule is a map-silent default — flagged)* — names
  are addressable in `BENCH_CANARY_PHASE`, summary lines, and `[name]`
  prefixes, so a splitting byte would corrupt three contracts. `argv`:
  required non-empty exec array, no element interpolation, argv[0] non-empty.
  `env`: JSON object of string pairs, applied in sorted key order
  *(shape flagged)*. `needs`: array of phase names; dangling and cyclic
  (including self) edges are malformed; duplicate entries are deduplicated.
  `optional`: bool, today's skip-if-binary-absent semantics.
  `dir`: relative path, default the root, validated lexically to stay inside
  the root (absolute paths and `..` escapes are the escaping-dir class);
  the phase runs with its working directory at root/dir.
- **`Phase` struct.** `Serial` is removed; `Needs []string` and `Dir string`
  are added. `BenchkitPhases` emits `needs: [build]` on conformance,
  contract, shellcheck, and canary when the build phase materializes
  (`scripts/go-build.sh` + `go.mod` present), preserving today's ordering as
  the degenerate edge the map names.
- **Scheduler.** Replaces `splitSerialPhases` in both runners. Outer mode:
  a phase starts when every need is green; phases with unmet-able needs
  (a need red, skipped, or itself blocked) resolve skipped-with-cause without
  launching; everything runnable runs concurrently, unweighted, under the
  existing prefix writers and capability-skip env injection. Inner mode:
  the same graph executed sequentially in topological order (declaration
  order among ready phases), preserving today's no-summary byte shape and
  the fail-fast-equivalent behavior via skip propagation.
- **Skip/red propagation** *(skip propagation is a map-silent default —
  flagged)*: a dependent of a red phase reports
  `phase <name>: skipped (needs <failed-phase>)`; a dependent of a skipped
  phase propagates the same way. Skipped-with-cause is never red on its own;
  the gate's verdict comes from the red phase. Summary lines print in table
  (declaration) order after all phases settle *(ordering unification
  flagged: today a red serial phase prints early and short-circuits; under
  the DAG its dependents are skipped so the observable effect for the
  current kit table is unchanged)*.
- **Deadline and stragglers** *(mechanism is a map-silent default — flagged;
  the falsification pass killed the draft's first design)*: the single
  gate-level deadline stays exactly where it is — `gate-run`'s
  `bounds.GateTimeout`, exit 124, `internal/bounds` untouched. What changes
  is the kill cascade on that one path: `gate-run`'s deadline branch sends
  SIGTERM to the gate process group, waits the existing
  `processGroupCancelGrace`, then SIGKILLs — mirroring the SIGINT path —
  and `gate-phases`, on its termination signal, prints
  `gate: cancelled; still running: <name>[, <name>…]` to stderr (phases
  already settled are excluded) before exiting. `gate-run` still maps the
  deadline to exit 124; operator SIGTERM/SIGINT gains the same straggler
  report for free. Rejected: a second, earlier `GateTimeout − lead` timer
  inside `gate-phases` — it would add a bounds constant and a second
  deadline, both of which the map closes against.
- **Env merge.** One merge helper applies phase env strip-then-set over
  `gateEnv()`; the existing `append(gateEnv(), phase.Env...)` call site
  converges on it.
- **Timing (#9 reading — flagged for veto).** The tree emits no per-phase
  timing today: the only runner-emitted timing is the conformance driver's
  per-check lines, keyed on the phase name `conformance`. This spec reads
  #9's "each new overlapping phase gets a timing line like any phase today"
  as exactly that status quo — manifest phases named `conformance` inherit
  the per-check timing hook, and no new timing format is introduced,
  because #9 also closes against any output-contract redesign. If the
  reviewer intended a per-phase wall-clock line, that is a new output
  format and needs its own decision first.
- **Complete map-silent default inventory (veto list).** Everything this
  spec decided that the map does not carry, in one place: (a) top-level
  `{"phases": [...]}` object + strict decoding; (b) `env` as a JSON object
  applied in sorted key order; (c) the name charset rule; (d) skip
  propagation to dependents; (e) summary lines printed in declaration order
  after all phases settle; (f) the SIGTERM-grace straggler mechanism above;
  (g) the #9 timing reading above; (h) dangling-symlink and special-file
  manifest paths red rather than absent (derived from the profile's
  hostile-input checklist, not the map); (i) duplicate entries within one
  `needs` array deduplicated silently; (j) `dir` containment validated
  lexically only; (k) `canary`/`conformance` stay name-keyed contracts,
  documented not reserved; (l) inner-mode tie-break among ready phases is
  declaration order; (m) invalid UTF-8 inside JSON strings is coerced to
  U+FFFD by `encoding/json`, not refused — refusal would need a bespoke
  byte scan and the manifest is the project's own declaration.
- **Untouched.** `.bench/gate.sh` stays the one-line exec entry; the
  resolution chain, subject/verdict machinery, `internal/bounds` timeout
  value, capability-skip tally, conformance timing print (keyed on the
  phase name `conformance`), and `phasesForMode`'s owner allowlist all stay
  as they are. The kit ships no `.bench/phases.json` in this slice — the
  kit-specific manifest lands in slice C with the phases it declares.

## Testing decisions

- A good test here drives `runPhases`/`PhasesCommand` with real cheap
  subprocesses (`fakePhase`, marker files, `pwd`/`env`-printing scripts) and
  asserts observable output, exit codes, and filesystem effects — never
  scheduler internals. Loader tests are table-driven over manifest byte
  fixtures asserting exit code and diagnostic substring per class.
- Seams: the loader behind `PhasesCommand` (with `benchkitPhasesForCommand`
  stubbed as prior art does), and the scheduler behind `runPhases` (prior
  art: `TestRunnerRunsPhasesConcurrently`, `TestRunnerAggregatesAllPhases`,
  `TestRunnerInnerModeByteShape`, `TestPhasesCommandSignalCancelsRunningPhaseGroups`).
  The deadline injects via the same package-var pattern as `gateTimeout`.
- End-to-end bite: one canary fixture with a malformed manifest EXPECTing
  the loader diagnostic, following existing fixture conventions.
- Gate: `.bench/gate.sh` (the project gate), green required to commit.

### Seam diagram — manifest loader (fronted by `gate-phases`)

    trigger: bench gate → .bench/gate.sh → bench gate-phases <root>
        │
        ▼
    <root>/.bench/phases.json ──▶ [ loader: stat → read → decode → validate ] ──▶ []Phase (manifest)
    absent file               ──▶ [                                         ] ──▶ []Phase (built-in table)
    malformed/empty/special   ──▶ [                                         ] ──▶ stderr diagnostic + exit 1
                  ◀ tests attach here: PhasesCommand against a temp root with a
                    manifest byte-fixture; assert exit code, stderr substring,
                    and (via stubbed built-in table) which table ran

### Seam diagram — DAG scheduler (fronted by `runPhases`)

    trigger: gate-phases after load (outer or inner mode)
        │
        ▼
    []Phase{Name,Argv,Env,Needs,Optional,Dir} ──▶ [ scheduler: ready-set launch,   ] ──▶ per-phase prefixed output
    ctx (termination signal)                  ──▶ [ red/skip propagation, settle,  ] ──▶ summary + timing lines
                                                  [ straggler naming on signal     ] ──▶ exit 0/1/124/130
                  ◀ tests attach here: runPhases with fakePhase subprocesses that
                    write marker files / sleep / exit nonzero; assert ordering by
                    observed effects, summary lines, exit code

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | valid manifest replaces built-in table | PhasesCommand + stubbed built-in table | `go test ./internal/gate -run TestPhasesCommandLoadsManifest` observed red before build | if loading is missing the stub table's marker phase runs and the manifest phase's marker file never appears |
| 2 | absent manifest falls back to built-in table | PhasesCommand + stubbed built-in table | already covered (`TestPhasesCommandRoutesCanaryToOwningPhase` and siblings run rootless of any manifest) plus new explicit `TestPhasesCommandAbsentManifestFallsBack` observed red | existing suites red if fallback breaks; the explicit test reds if absence is misclassified as malformed |
| 3 | blank file and `{"phases":[]}` both red naming empty | loader table test | `go test ./internal/gate -run TestManifestEmptyIsRed` observed red before build | a fallback-on-empty implementation exits 0 with the built-in table and the asserted diagnostic never prints |
| 4 | parse error reds with a class-naming diagnostic | loader table test | `TestManifestMalformed/parse-error` observed red before build | tolerant decoding would exit 0; the row demands exit 1 plus the parse-class substring |
| 4 | duplicate phase name reds naming the class and the duplicated name | loader table test | `TestManifestMalformed/duplicate-name` observed red before build | a one-generic-diagnostic loader fails the class+element assertion; without the check the last duplicate silently wins |
| 4 | invalid phase name (empty, whitespace, control byte) reds naming the class and the offending name | loader table test | `TestManifestMalformed/invalid-name` observed red before build | an unvalidated name reaches summary/prefix/env contracts; a generic diagnostic fails the class+element assertion |
| 4 | dangling `needs` reds naming the class and the unknown edge target | loader table test | `TestManifestMalformed/dangling-needs` observed red before build | an unknown edge target either panics or is silently dropped; the class+element assertion rejects a generic red |
| 4 | cyclic `needs` (incl. self) reds naming the class and a phase on the cycle | loader table test | `TestManifestMalformed/cyclic-needs` observed red before build | a cycle deadlocks the scheduler or drops phases; the loader must refuse first, attributably |
| 4 | empty `argv` / empty `argv[0]` reds at load naming the class and the phase | loader table test | `TestManifestMalformed/empty-argv` observed red before build | today empty argv only reds at phase run time; the row pins attributable load-time refusal |
| 4 | `dir` absolute or escaping root reds naming the class and the dir value | loader table test | `TestManifestMalformed/escaping-dir` observed red before build | unvalidated dir runs a phase outside the graded root and exits 0 |
| 4 | unknown field reds at both levels (top-level key and phase-level key, two subcases) naming the key | loader table test | `TestManifestMalformed/unknown-field-top` and `/unknown-field-phase` observed red before build | a decoder strict at only one level accepts the other; a typo like `"need"` silently drops an edge and exits 0 |
| 5 | dangling-symlink manifest reds, no fallback | loader test with symlink fixture | `TestManifestDanglingSymlinkIsRed` observed red before build | a plain-read implementation sees ENOENT and silently grades with the built-in table |
| 5 | FIFO at manifest path reds without blocking | loader test with mkfifo | `TestManifestSpecialFileIsRed` observed red before build (with test timeout as the hang tripwire) | a read-first implementation blocks forever on the FIFO; stat-first exits 1 promptly |
| 6 | phase runs with cwd root/`dir`, default root | runPhases with pwd-printing fakePhase | `TestRunnerPhaseDirIsRelativeToRoot` observed red before build | if `dir` is ignored the printed cwd is the root and the assertion on the subdir fails |
| 6 | a valid manifest's six fields all survive decode end-to-end | PhasesCommand against a temp root whose manifest exercises `needs` ordering, an `env` value, `optional` skip, `dir` cwd, and names in summaries | `TestPhasesCommandManifestFieldsEndToEnd` observed red before build | a loader that decodes only `name`/`argv` and discards the rest passes every direct-`Phase` runner test but fails this one, because the discarded fields' observable effects never appear |
| 7 | phase env applied strip-then-set over gate env | runPhases with env-printing fakePhase carrying a key gateEnv also sets | `TestRunnerPhaseEnvStripsThenSets` observed red before build | append semantics hand the child two values for one key; the test asserts exactly one, with the phase's value |
| 8 | dependent starts only after its need completes | runPhases, marker-file ordering | `TestSchedulerRespectsNeeds` observed red before build | with edges ignored both phases overlap and the dependent observes the need's marker missing |
| 8 | phases with no unmet edges overlap unweighted, no internal width cap | runPhases: four independent fakePhases that each wait for all four start-markers to exist before exiting | `TestSchedulerOverlapsIndependents` observed red before build (flat two-way overlap already covered by `TestRunnerRunsPhasesConcurrently`) | a serializing or fixed-two-worker scheduler never lets all four subprocesses exist at once, so the barrier times out; the subprocess barrier is GOMAXPROCS-safe because waiting goroutines are syscall-blocked, not P-holding |
| 9 | dependents of a red phase report skipped-with-cause, distinct from red | runPhases with red fakePhase + dependent | `TestSchedulerSkipsDependentsOfRed` observed red before build | fail-fast cancels and reports nothing; run-everything reports the dependent red — the asserted `skipped (needs …)` line rejects both |
| 9 | independents of a red phase run and report | runPhases, marker files | same test as above, sibling assertion, observed red before build | full fail-fast leaves the independent's marker missing and its summary absent |
| 10 | dependents of a skipped optional phase skip with cause | runPhases with optional-absent phase + dependent | `TestSchedulerPropagatesOptionalSkip` observed red before build | treating skip as green runs the dependent against an artifact never produced |
| 11 | built-in table expresses build→all as needs; `Serial` gone | updated `TestPhaseTableBuildPhase` | updated test asserting `Needs: ["build"]` on all four downstream phases observed red before the table conversion (the current test asserts `Serial`, the opposite, so it cannot stand as coverage); removing the field makes every `Serial` reference a compile red | a conversion that deletes `Serial` without adding the edges passes compilation and every name-pinning test, but fails the explicit `Needs` assertion |
| 12 | on termination the runner names only phases still running, excluding settled ones | PhasesCommand seam (stubbed table, prior art `TestPhasesCommandSignalCancelsRunningPhaseGroups`): one phase completes and writes its marker, then the test signals while a second sleeps | `TestPhasesCommandNamesStragglersOnTermination` observed red before build | a runPhases-only implementation that never wires the signal path at the command seam, and a lazy report that names every phase, both fail — the assertion demands the sleeper named and the completed phase absent |
| 12 | `gate-run`'s deadline path delivers SIGTERM-grace-SIGKILL so the straggler report lands, still exit 124 | gate-run deadline branch with injected short `gateTimeout` (prior-art package-var injection) against a gate script that traps TERM | `TestGateRunDeadlineTermGraceThenKill` observed red before build | today's branch SIGKILLs immediately, so the child's trap output never appears and the test's report assertion fails; exit 124 is asserted unchanged |
| 13 | exit 130 on interrupt, process groups killed | PhasesCommand signal path | already covered — `TestPhasesCommandSignalCancelsRunningPhaseGroups` | the scheduler rewrite sits under the same entry; a broken 130 path turns this existing test red |
| 13 | final `gate:` lines, exit codes, conformance timing, inner byte shape unchanged | existing runner, timing, canary suites | already covered — `TestRunnerFinalLineAndExitCodes` (final lines, 0/1/3), `TestRunnerPrintsConformanceTiming`, `TestRunnerInnerModeByteShape` (green inner path), plus the canary sweep's substring EXPECTs | these pin those exact strings; any drift under the new scheduler reds them |
| 13 | outer `phase <name>: green`/`red (exit N)`/`skipped (…)` summary lines pinned under the DAG scheduler | new runPhases outer-mode test asserting all three exact summary forms (nothing pins the green/red forms today — only the skip verdict is asserted, in `runner_serial_test.go`) | `TestRunnerSummaryLineByteShape` observed red before build (asserted against the current strings first, then kept green through the rewrite) | a scheduler rewrite that drops or reformats summaries passes every existing test; this row makes the summary contract a tripwire before the rewrite starts |
| 14 | loader red bites through the real gate entry permanently | canary sweep fixture with malformed manifest | fixture EXPECT observed red in the sweep during build (craft-gate recorded red) | if the loader posture rots to silent fallback the inner gate goes green and the sweep reports the vacuous EXPECT |
| 15 | inner mode loads the manifest, filters via `phasesForMode`, and runs topologically, byte shape unchanged | PhasesCommand with `BENCH_CANARY_INNER=1` against a temp root whose manifest carries needs edges | `TestPhasesCommandInnerModeManifestDagOrder` observed red before build | an implementation that loads manifests only in outer mode, skips the mode filter, or ignores edges sequentially each produces the wrong phase set or order at this seam; a direct-runPhases test would see none of that |

### Edge inventory

Walked classes landing as rows above: absent vs present-but-empty (rows 2/3),
malformed family per class (row group 4), dangling symlink (row 5), special
file (row 5), duplicate env keys (row 7), paths with spaces (existing
`TestRunnerRootWithSpace` continues to pin the runner; `dir` row 6 uses a
plain subdir), interrupt mid-run (row 13), control bytes in names (row 4
invalid-name).

Won't handle:
- manifest JSON without trailing newline — `encoding/json` has no line
  semantics; the parse rows exercise byte-level tolerance already.
- `dir` escaping via a symlink inside the root — the manifest already
  executes arbitrary `argv` from the graded tree, so containment is a
  correctness guard, not a security boundary; lexical validation is the
  decided level.
- a manifest phase named `canary` or `conformance` — name-keyed contracts
  (inner-mode exclusion, timing print) apply to those names by design; the
  loader documents rather than reserves them.
- duplicate entries inside one `needs` array — deduplicated silently; no
  ambiguity, nothing graded wrong.
- concurrent mutation of `phases.json` during a run — the gate already
  rejects a subject that changes mid-run; the manifest is read once at load.
- non-UTF-8 manifest bytes — structural garbage surfaces as the parse-error
  class; invalid bytes inside a JSON string are coerced to U+FFFD by
  `encoding/json` rather than refused (veto item m) — the result is an odd
  but harmless name or value in the project's own declaration.

## Out of scope

- **The kit's own `.bench/phases.json` and consumer-facing manifest docs** —
  a separate capability: the kit has no manifest consumer until slice C
  splits `checkGoCore` and declares the kit-specific phases; documenting an
  unconsumed format would drift. Lands in slice C (reviewer-sequenced).
  Estimate if pulled forward: 4 edits, 2 gate runs.
- **`checkGoCore` split, fixture migration, parity inventory test, go.mod-
  probed fallback phases, `phasesForMode` allowlist widening** — slice C,
  already sequenced by the reviewer on the map; distinct capability
  (restaging checks) from this slice's mechanism. Estimate: the slice-C spec.
- **Ship-tier manifest membership, `bench upgrade` semantics for a
  project-owned manifest, cross-language capability-skip surface** — the
  map's "Not yet specified" list; each needs its own decision before spec.
- Per-phase timeouts, weight fields, YAML/TOON formats, generated `gate.sh`
  shims — rejected alternatives on the map; not re-arguable here.
