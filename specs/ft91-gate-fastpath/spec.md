# ft91-gate-fastpath — per-test canary bites + gate-level verdict reuse

Status: implemented

Compiled from `decisions/gate-critical-path.md` (the gate assessment, adopted
as the working map) under the reviewer-directed batch-drain override: levers 1
and 2 were explicitly pushed into this spec ("write the spec for lever 1";
"include lever 2 in that spec"). Every decision this spec defaulted rather
than received from the map is flagged **[defaulted]** for post-hoc veto. The
map stays at `decisions/gate-critical-path.md` until the in-flight
`ft91-artifact-hoist` retire commit lands; this spec's staging commit then
moves it to `specs/ft91-gate-fastpath/decisions/` per the map-move rule. A
mid-tier falsification pass reviewed the draft; its four blocking findings
(R17 fault-family breakage, lock-demotion deletion, wrapper-forwarding
misstatements, ungraded `-test.list` argv) are folded in below.

## Problem

The dev gate is ~135–172 s and one package wide: `internal/contract/surface/
artifact` (~109 s solo) floors the contract phase once and the canary sweep
five more times, because each behavior-owned fixture's compiled bite runs
every `Test*` in its owning package when its EXPECT is, by stage 2's own
ruling, one owning test's failure message. Separately, a green verdict for an
identical closed subject is re-judged from scratch everywhere except `bench
commit`: a final-check → commit → push chain, a post-review re-gate on an
untouched tree, or a re-typed `bench gate` each pay the full wall for a
subject the oracle already answered, and the 10-minute freshness window
expires across any real review gap.

## Solution

**Lever 1 — per-test bites.** Each behavior-owned fixture names its owning
contract test in a `TEST` file (mirroring `CHECK`); the sweep invokes the
compiled package binary with `-test.run` anchored to exactly that test.
Binding defects — missing, empty, stray, or unknown names — are loud reds
reported together, with names graded against the compiled binary's own
`-test.list`. Vacuity baselines stay package-wide (the
wider-baseline-fails-safe principle already documented on `scopeBaselines`).
Bite and vacuity semantics are unchanged; execution narrows to the semantics
stage 2 already declared. Expected: canary ~172 s → ~30–60 s.

**Lever 2 — verdict reuse at the gate.** Gate execution short-circuits on
`ReusableGreen` — the existing closed-subject key (working tree including
untracked files, oracle closure, policy, declared env/paths/tools) —
**under the execution lock, after the under-lock subject re-check and before
the pending replace**. Everything up to and including lock acquisition is
byte-identical to today, so the lock-contention demotion
(`persistInterruptedIfGreen`) and every pre-acquire contract stay untouched;
what a reuse skips is the pending write, the run, the post-run re-check, and
the final write. The reused result prints an explicit line, writes nothing,
and returns the full green tuple. `bench gate --fresh` forces a real run.
Freshness extends 10 → 60 minutes after the tool closure is completed (`go`,
`node`, `npm` declared in `.bench/gate-inputs.json`). The R17 fault-injection
contract family, which today seeds a reusable green and relies on `Execute`
always running, is re-staged in the same change to seed non-reusable
subjects — a declared, reviewer-signed contract co-move, not a quiet test
edit. ADR 0002 posture 5 is revised to name gate execution as a reusing
consumer. Expected: unchanged-tree gate ~2–5 s.

## User stories

1. As the canary sweep, I want a bite invocation (`RunKind` `RunBite`) to
   carry `-test.run` anchored `^<owner>$` with the owner regexp-quoted, each
   fixture resolving its own owner, so a fixture pays one owning test
   instead of its whole package suite.
   Line: opus (mid alias) / medium. This is oracle-harness logic, which the
   profile's cached routing sends to the mid tier at medium effort.
2. As a fixture author, I want each behavior-owned fixture to declare its
   owning test in a `TEST` file — read with the same trim discipline as
   `CHECK` — so the binding lives in the fixture the way every other canary
   fact does. **[defaulted: the `TEST`-file convention itself; the map named
   the need, not the vehicle]**
   Line: opus (mid alias) / medium. The file convention touches sweep
   discovery, which is oracle infrastructure on the cached mid routing.
3. As the sweep, I want a missing or empty `TEST` on a behavior-owned
   fixture, a `TEST` outside the behavior-owned family, and a `TEST` at any
   non-fixture level inside it (a family or package directory, where nothing
   would read it) each to be a loud red naming the offending path — all
   binding defects reported together, beside the existing
   `assertContractScopes` refusals — so there is no silent fallback to a
   full-package run (stage-1 posture) and no unread binding lying anywhere
   in the tree.
   Line: opus (mid alias) / medium. Fail-closed refusal shape must match the
   existing structural refusals exactly, a correctness-over-speed edit.
4. As the sweep, I want every declared owner graded against its compiled
   binary's `-test.list` output after the per-package compile and before any
   baseline or fixture runs — unknown names reported together, one
   diagnostic per fixture — so a renamed or mistyped owner is a named
   refusal rather than a did-not-bite archaeology session. **[defaulted
   twice: `-test.list` against the real binary, chosen over a second
   Go-parser derivation, per the run-the-real-path standard; and the
   departure from the map's "at sweep start" shape — name validation cannot
   precede the compile that produces the binary, so structural defects red
   pre-compile (story 3) and name defects red post-compile, both before any
   graded run]**
   Line: opus (mid alias) / medium. The validation is the one guard between
   a typo and a misgraded oracle, mid per the cached routing.
5. As the vacuity check, I want contract-group baselines to keep running the
   whole package binary with no `-test.run`, so a baseline is never narrower
   than the runs it grades and a vacuous EXPECT cannot clear the screen in
   silence — the same wide-baseline direction `scopeBaselines` already
   documents for phase pins. **[defaulted: baseline stays package-wide
   rather than gaining a per-test twin — deliberately reversing the map's
   scoped-baseline sketch, which would have narrowed the screen in the
   fail-dangerous direction]**
   Line: opus (mid alias) / medium. Deliberate non-change at the oracle's
   vacuity screen, graded at the cached mid routing.
6. As the marker-file reader, I want one guarded reader shared by `CHECK`
   and `TEST` that rejects non-regular files before opening, so a FIFO
   planted at either path cannot wedge the sweep — closing the existing gap
   where `fixtureCheck` opens `CHECK` unguarded (profile hostile-input
   class: special files in discovered paths). **[defaulted: fixing the
   pre-existing `CHECK` gap in this spec rather than parking it — it is the
   same reader]**
   Line: opus (mid alias) / medium. A fail-closed guard on the oracle's
   discovery path, mid per the cached routing.
7. As the migration, I want all 33 behavior-owned fixtures to carry `TEST`
   files naming their EXPECT's owning test, the full sweep green, and one
   deliberately wrong owner demonstrated red (did-not-bite) during the
   build and reverted, so the scoping is proven to bite rather than assumed.
   Line: opus (mid alias) / low — a declared deviation from the cached
   mid-effort routing, because the story is mechanical transcription of
   EXPECT owners that the sweep itself grades on every run.
8. As gate execution, I want the resolve-run-record path to short-circuit
   under the held lock, after the under-lock subject re-check and before the
   pending replace, when the subject's verdict is `ReusableGreen`: it prints
   `gate: green (fresh verdict reused for this tree)` on stdout, writes
   nothing — the verdict record, its `RecordedAt`, and the cache file stay
   byte-identical, so reuse never slides the freshness window — and returns
   the full green tuple (`GateExit` 0, `ActionExit` 0, an `Inspection`
   projecting the reusable green), so shift, commit, status, and the
   projection contracts all read a truthful result.
   Line: opus (mid alias) / medium. This is the oracle's authority surface;
   the cached gate-logic routing applies at full effort.
9. As the fault-injection contracts, I want the R17 family
   (`runner_test.go` / `runtime_gate_partial_proof_test.go`) re-seeded to
   non-reusable subjects (an expired `RecordedAt` via the fake engine's
   clock) for every faulted op the short-circuit now precedes, with their
   registry tuples updated in the same change, so they keep grading fault
   durability on the real run path — a declared contract co-move under this
   spec's sign-off, never a quiet edit to make a test pass. Pre-acquire
   contracts (lock-open, lock-acquisition, the `interrupted-pending`
   demotion in `owner_record_test.go`) need no edit: the placement leaves
   their path byte-identical.
   Line: opus (mid alias) / medium. Re-staging oracle contracts is gate
   logic at the cached mid routing.
10. As an operator, I want `bench gate --fresh` to force a real run past a
    reusable green: `gate_command` in `bin/bench.sh` accepts and forwards
    the flag (today it arity-rejects any second token except `pin`/help),
    `run_gate` forwards arguments, `RunCommand` parses the flag without
    colliding with its optional root positional, and any other unknown flag
    still exits 2 with usage on stderr and an untouched oracle — the
    existing wrapper contract in `runtime_gate_test.go` extended, not
    bypassed. **[defaulted: flag spelling `--fresh`; `bench commit` gains no
    such flag — its contract is reuse-when-valid]**
    Line: opus (mid alias) / medium. The flag crosses the wrapper contract
    surface, which is oracle plumbing at the cached routing.
11. As `bench commit`, I want my own pre-gate `ReusableGreen` check
    collapsed into the single gate-package home, so verdict-reuse policy has
    exactly one source (the package's stated charter) and commit's observable
    behavior — reuse on fresh green, refusal on red, exactly one gate run
    tallied by the runtime commit contract — is unchanged, with the
    reuse-line emitter moving from `commit.go` to the gate home and the
    contract expectation updated to the new emitter in the same change.
    **[defaulted: collapse now rather than leaving the duplicate]**
    Line: opus (mid alias) / medium. Touches the gated-commit authority
    path, mid per the cached routing.
12. As the oracle policy, I want freshness extended from 10 to 60 minutes in
    the one shared constant — editing the existing hard-coded boundary pin
    in `gate_test.go` in the same change — accepting that `bench status`'s
    stale signal and `bench prep-release`'s dev-green entry check inherit
    the same window, and with `policyVersion` left unbumped, accepting the
    stated consequence that verdicts recorded under the 10-minute policy
    become retroactively reusable for 60 minutes on the new binary's first
    run. **[defaulted: 60 min; single shared window; no policy bump — all
    three are reviewer calls the map left open]**
    Line: opus (mid alias) / medium. The freshness window is the single
    highest-stakes semantic the reviewer signs here; cached routing, full
    effort.
13. As the subject closure, I want `go`, `node`, and `npm` declared in
    `.bench/gate-inputs.json` tools so a toolchain upgrade changes the
    subject and can never be papered over by a reused green, while
    `shellcheck` stays undeclared because its legitimate absence (an
    optional phase) would open the subject and silently disable reuse on
    every host without it — that in-place-upgrade residual is documented in
    the ADR revision. **[defaulted: the exact tool list and the shellcheck
    exclusion]**
    Line: opus (mid alias) / medium. Closure membership decides what a
    reused green can hide; cached gate routing.
14. As the decision record, I want ADR 0002 posture 5 revised to name gate
    execution (not only the gated commit) as the consumer that reuses a
    fresh green for the identical closed subject, with the same reopen
    trigger, and to record the new accepted residual (an in-place
    `shellcheck` upgrade inside the freshness window under a reused green).
    Line: fable (top alias) / high — the profile's doc-authoring leverage
    override: the revision authors a new accepted-risk clause, not a
    transcription. This is a top-tier bump; approving this spec approves the
    bump.
15. As the map, I want the post-build re-measurement (solo canary, full gate
    changed-tree, and one unchanged-tree reuse timing) recorded in
    `decisions/gate-critical-path.md` against the ≤60 s stop rule, so the
    evidence lives where the levers were decided.
    Line: opus (mid alias) / low — a declared deviation from the cached
    mid-effort routing, because the story transcribes measurements into the
    map with no oracle logic touched.

## Implementation decisions

- `internal/canary`: `RunCall` gains a `Test` field and a `RunList` kind;
  `runnerCommand` gains a real-argv case for both — the bite case appends
  `-test.run ^<QuoteMeta(owner)>$`, the list case invokes the compiled
  binary with `-test.list '.*'` — and the existing real-runner dispatch
  test (`TestDefaultRunnerDispatchesOnCallKind`) grows both cases, because
  the fake-runner seam alone is the exact insufficiency its own comment
  names. `subjectCall` sets `Test` for fixtures only — baselines never
  carry it. The anchored-regexp prior art is `raceTestFilter` in
  `internal/gate`.
- `-test.list` semantics: one list call per bound package, in the compile
  stage's sequencing (after its package's compile, before baselines),
  sharing the sweep's worker budget; environment is the sweep base plus the
  width pin, no subject root (nothing is graded); a nonzero exit is a sweep
  error naming the package; membership accepts only lines beginning `Test`
  (the flag also lists benchmarks, fuzz targets, and examples, none of
  which can own an EXPECT).
- Binding defects split by what they need: structural refusals (missing,
  empty, stray, non-fixture-level `TEST`) red pre-compile beside
  `assertContractScopes`; name refusals red post-compile via the list.
  Both stages report every defect they can see.
- One regular-file-guarded marker reader serves `CHECK` and `TEST`
  (story 6) — two open-coded reads would be the knowledge duplication the
  code standard names.
- `internal/gate`: the reuse short-circuit lives in
  `executeWithEngineAfterAcquire` between the under-lock subject re-check
  and the pending replace. Pre-acquire behavior (lock contention,
  `persistInterruptedIfGreen` demotion, owner-record ordering) is
  byte-identical. R17 subtests for post-short-circuit ops re-seed via the
  fake engine's clock so every faulted op is still reached on a real run.
- `--fresh` plumbing: `gate_command` and `run_gate` in `bin/bench.sh`
  forward the flag; `RunCommand` distinguishes it from the root positional;
  the wrapper's unknown-flag usage contract is extended to accept exactly
  this token.
- `bench commit` calls the gate home unconditionally and loses its private
  check; the reuse line is printed by the gate. Inner canary gates route
  through `PhasesCommand`, which records and consults nothing — untouched.
- `freshness` moves to 60 min where it is declared; no second constant; no
  `policyVersion` bump.
- `.bench/gate-inputs.json` gains `"go", "node", "npm"` in `tools` — the
  identity collector already hashes declared tools and their shebang
  chains, and a declared-but-absent tool only opens the subject (disables
  reuse, fail-safe), never fakes a closure.
- ADR 0002 posture 5 wording change; no other posture moves.

## Testing decisions

- A good test here drives the real surface: the sweep through an injected
  `Runner` asserting the `RunCall`s it emits (prior art:
  `compiled_bite_test.go`, `canary_test.go`) **plus** the real-argv
  dispatch test for every new call shape, and gate reuse through `Execute`
  against a throwaway repo whose gate script leaves a run marker, so "did
  not re-run" is an observable fact rather than a mock's opinion.
- Seams: the canary `Runner` injection seam and the real-argv dispatch
  seam; the gate execute seam (`executeWithEngine` + throwaway repos,
  prior art `gate_test.go`, `fault_engine_test.go`, `owner_record_test.go`);
  the runtime contract suites for wrapper and commit behavior; the live
  gate itself for migration proof.
- Gate command: `bench gate` (dev tier), green required to commit.

### Seam diagram — canary bite scoping

    trigger: canary phase of `bench gate` / `bench canary <root>`
        │
        ▼
    fixture dir (EXPECT, TEST, files/) ──▶ [ sweep: discover → refuse →    ] ──▶ RunCall{RunList, Binary}
    compiled pkg binary  ──────────────▶  [ compile → -test.list validate  ]     RunCall{RunBite, Binary,
                                          [ → baseline → per-fixture run   ] ──▶   Test: owner}
                                                                            ──▶ bite / did-not-bite / vacuous
                      ◀ tests attach here: injected Runner records every RunCall
                        (scoping, refusals); the real-argv dispatch test grades
                        what each RunCall kind actually executes

### Seam diagram — gate verdict reuse

    trigger: `bench gate` [--fresh] / shift loop / stop hook / `bench commit`
        │
        ▼
    working tree + oracle closure ──▶ [ gate execute: subject → lock →     ] ──▶ exit 0 + reuse line +
    verdict cache (git dir)      ──▶ [ re-check → ReusableGreen? reuse :   ]     green tuple (no write), or
                                     [ pending → run → record              ] ──▶ real run + recorded verdict
                      ◀ tests attach here: throwaway repo, gate script appends
                        to a marker file; tests assert marker count, stdout
                        line, returned tuple, and byte-identical cache record

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | two fixtures with distinct TEST owners produce bite RunCalls carrying each its own anchored `-test.run` value, no clash | injected-Runner sweep test (two-fixture shape, prior art `TestBiteCarriesItsOwnFixtureTreeAsSubjectRoot`) | observed red required before green: the assertion is written against `RunCall.Test` while the field does not exist | a per-package or first-fixture-wins owner resolution — the cheapest wrong scoping — passes any single-fixture assertion and fails this one |
| 1 | an owner name containing regexp metacharacters matches only itself | unit test on the argv builder | observed red required before green, then one recorded mutation during the build: the unquoted-interpolation variant observed red and reverted | unquoted interpolation would let a name match a superset and pass for the wrong test; the recorded mutation separates quoting from the filter's mere existence |
| 1, 4 | `runnerCommand` builds the real argv for both new shapes: bite with `-test.run`, list with `-test.list '.*'` | the real-argv dispatch test (`TestDefaultRunnerDispatchesOnCallKind` extended) | observed red required before green: both cases asserted while `runnerCommand` has no list case and no filter append | the fake-runner seam alone is explicitly insufficient per that test's own comment — a sweep emitting correct RunCall metadata while executing unfiltered binaries would pass every fake-seam row and grade nothing in production |
| 3 | missing TEST, empty TEST, stray TEST outside the family, and TEST at a family or package level inside it are distinct loud reds naming the path, reported together | injected-Runner sweep test | observed red required before green: all four refusal messages asserted while the sweep still runs such fixtures unscoped | absent-vs-empty is a profile edge class; a package-level TEST is an unread binding — a lie in the tree — and silent fallback is the stage-1 posture violation this spec exists to prevent |
| 4 | an owner absent from the compiled binary's `-test.list` is a named refusal before any baseline or fixture run, all such defects together | injected-Runner sweep test (fake runner answers the list call) | observed red required before green: unknown-owner case red while validation does not exist | a renamed owning test would otherwise surface as did-not-bite archaeology per fixture instead of one named refusal |
| 5 | contract-group baselines carry no `-test.run` | injected-Runner sweep test | observed red required before green: asserted against a variant that scopes baselines | a scoped baseline passes every vacuous EXPECT in its group in silence — the exact failure the wide-baseline principle names |
| 6 | a FIFO at CHECK or TEST is rejected before opening, with a named refusal | unit test at the marker reader (FIFO planted in a temp fixture), bounded by the test's own deadline | observed red required before green: while the reader opens unguarded, the read blocks and the bounded deadline reds the test rather than hanging the suite | an unguarded open of a FIFO blocks forever; the deadline converts the wedge into an observable red at the unit seam, where the sweep's 10-minute panic would otherwise be the only signal |
| 7 | all 33 migrated fixtures bite non-vacuously under per-test scoping | the live gate's canary phase | already covered — the sweep's bite/vacuity semantics are unchanged and grade every fixture on every run | a fixture whose TEST names the wrong owner reports did-not-bite; the sweep is the standing proof |
| 7 | the scoping bites: a deliberately wrong owner reds as did-not-bite | one recorded break-it run during the build (wrong TEST, observe red, revert) | observed red required once during the build | proves the filter grades the named test rather than accidentally passing on any package failure |
| 8 | a second `Execute` on an unchanged tree within freshness runs nothing: marker count stays 1, stdout carries the reuse line, cache record byte-identical, and the returned tuple is `GateExit` 0 / `ActionExit` 0 / an `Inspection` projecting the reusable green | gate execute seam, throwaway repo | observed red required before green: test red while `Execute` always runs | a reuse that re-runs is the lever missing; one that rewrites `RecordedAt` builds a sliding window; one that returns a zero `Inspection` makes status and the projection contracts report `absent` on a green tree — the tuple assertion catches the third |
| 8 | a red, expired, pending, subject-open, or tree-changed state pays a real run | gate execute seam (fake engine `Now()` for expiry; edited tree; red gate script) | observed red required before green: each non-reuse state asserted to re-run while the short-circuit is written | the cheapest wrong implementation short-circuits on any cached green; these rows pin reuse to exactly the ReusableGreen predicate |
| 9 | every R17 faulted op is still reached twice (or once for pre-acquire ops) on re-seeded non-reusable subjects, and the pre-acquire demotion contracts pass unedited | existing R17 + owner-record suites, re-seeded | observed red guaranteed at build time: the unrevised R17 seeds red the moment the short-circuit lands, and stay red until re-seeded — the co-move is driven by real reds, not rewritten expectations | if the short-circuit accidentally reordered the pre-acquire path, the unedited demotion contracts red; if a faulted op became unreachable, its re-seeded subtest reds on `opCounts` |
| 10 | `bench gate --fresh` forces a real run (marker count 2 past a reusable green); any other unknown flag still exits 2 with usage on stderr and an untouched oracle | gate execute seam + the wrapper contract in the runtime suite (extended) | observed red required before green: marker count asserted at 2 while the flag does not exist; the extended wrapper case red while `gate_command` arity-rejects the token | without the opt-out, a poisoned-but-green cached verdict has no operator escape; the wrapper case pins that the flag arrives at the binary instead of dying at arity checking |
| 11 | `bench commit` still reuses a fresh green (exactly one gate run tallied across gate-then-commit) and still refuses on red, through the collapsed path | runtime commit contract suite (`runtime_commit_test.go` gate-run tally; reuse-line expectation moved to the gate emitter in the same change) | already covered for the run-count behavior — the tally contract predates this spec; the emitter move is graded by updating that contract's expected output source, red until the collapse lands | the tally reds if the collapse double-runs or stops reusing; the updated expectation reds if the reuse line disappears from the surface the operator reads |
| 12 | freshness is 60 minutes, pinned at both edges | unit test on the inspect path with fake `Now()` at 59 and 61 minutes, replacing the existing 10-minute boundary pin in `gate_test.go` | observed red required before green: the 59-minute reuse case red while the constant is 10 min | the window is the semantic the reviewer signs; an unpinned constant drifts silently, and the stale 10-minute pin would red the moment the constant moves — naming it here makes the edit a declared co-move |
| 13 | a declared-but-absent tool opens the subject (no reuse, gate still runs) | the `R1/manifest-missing-tool` contract (`runtime_gate_proof_test.go`) | already covered — that contract proves the open-on-missing-tool mechanism; the new entries are manifest data riding it, and the build verifies the case exercises a declared-tool name rather than a fixture-only manifest | fail-safe direction is the whole argument for declaring tools; the contract keeps it observed |
| 14 | ADR 0002 posture 5 names gate execution as reusing consumer and records the shellcheck residual | — | not TDD-able — prose; graded on review against the ADR standard | — |
| 15 | re-measurement recorded in the map | — | not TDD-able — wall-clock is not gate-assertable (map precedent); ship evidence is the recorded measurement | — |

### Edge inventory

Walked per the profile's hostile-input checklist; every class lands as a row
above or a **Won't handle** line here.

- Trailing-newline absence and absent-vs-empty TEST → story 3's row plus a
  trim assertion at the marker reader (story 2/6 build); special files
  (FIFO), regexp metacharacters, unknown/renamed owner, unread
  package-level TEST → rows above.
- **Won't handle:** a dangling symlink at TEST — reads as missing and reds
  loudly either way; fail-closed in both branches, so the stat-first
  distinction buys no different outcome.
- **Won't handle:** multiline TEST content — after trimming, the embedded
  newline can never match a `-test.list` name, so it reds as an unknown
  owner with the offending bytes in the message.
- **Won't handle:** an owning test that capability-skips on this host — the
  EXPECT never appears and the fixture reds did-not-bite, byte-for-byte the
  behavior a full-package run produces today; no new exposure.
- **Won't handle:** an inner canary gate consulting the reuse cache — a
  fixture's mutated tree is a different subject by construction, and
  `PhasesCommand` neither records nor consults the cache; structurally
  unreachable rather than guarded.
- **Won't handle:** concurrent gate executions — the reuse check sits under
  the same execution lock every real run takes, so contention, the
  `interrupted-pending` demotion, and the owner-record ordering are
  byte-identical to today; story 9's unedited pre-acquire contracts are the
  standing proof.
- **Won't handle:** `shellcheck` upgraded in place within the freshness
  window under a reused green — accepted residual, recorded by story 14;
  declaring it would disable reuse on every host that legitimately lacks
  it.

## Out of scope

- **Artifact-suite restructure (assessment lever 4)** — a separate
  capability (test-package architecture, not oracle semantics); ~20 edits,
  ~8 gate runs. Next slice if the re-measure stays above 60 s.
- **`-count=1` removal (assessment lever 3)** — rejected in the map: Go's
  test-cache key cannot see subprocess effects, and lever 2 subsumes it
  with a complete key. Not deferred — refused.
- **Outer conformance/contract width-cap revival** — dormant per
  `decisions/gate-concurrency.md`; re-check rides story 15's measurement,
  not this build.
- **Per-test scoping for conformance-family fixtures** — already scoped to
  one check via `CHECK`/family binding; no further narrowing exists.
