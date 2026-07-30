# FT129 — attribute inner test aborts

Status: staged

## Problem

A canary run currently reduces every nonzero inner run that omits its fixture's
`EXPECT` to “did not bite.” That diagnosis is wrong when a Go test binary panics
or otherwise aborts before the owning check can report: it names the silent
fixture instead of the test or process failure that silenced it. FT122
demonstrated the failure when a runtime contract test panicked on an empty
subject-reported hash and an unrelated, correct behavior-owned fixture was
reported as no longer biting.

This spec is compiled from FT129's reviewed roadmap finding under the
reviewer-directed batch-drain override. No closed decision map exists. Every
choice supplied below is therefore marked `[defaulted]` for post-hoc veto or
`[reviewer-directed]` when it follows the reviewer's explicit Codex routing.

## Solution

Classify an inner test abort before applying the ordinary red-plus-`EXPECT`
bite rule. A panic names the deepest failing Go test or subtest the output
identifies; a process-level abort names the scoped owner or package when one is
known. The diagnostic is bounded and control-safe, the sweep continues to
aggregate other fixture failures, and ordinary missing-`EXPECT` runs keep the
existing “did not bite” result.

## User stories

1. As a Bench maintainer, I want a panicking inner Go test reported as an inner
   test abort with the responsible test or subtest, so that diagnosis starts at
   the cause instead of an unrelated fixture. **[reviewer-directed]** Line:
   `gpt-5.6-terra` / medium.
   This is oracle code at a known seam whose real-output grammar and precedence
   need careful implementation.
2. As a Bench maintainer, I want a signaled, unstartable, or otherwise
   process-level inner abort reported separately from an ordinary test failure,
   so that a run that never completed cannot masquerade as a weakened canary.
   **[reviewer-directed]** Line: `gpt-5.6-terra` / medium. This changes the
   runner result contract and must preserve cross-platform process behavior.
3. As a Bench maintainer, I want normal green runs and red runs that merely omit
   `EXPECT` to keep their current bite classification, so that sharper abort
   attribution does not weaken the canary or relabel ordinary failures.
   **[reviewer-directed]** Line: `gpt-5.6-terra` / medium. The gate fully
   observes the result, but the fail-safe precedence belongs to the oracle.
4. As a Bench maintainer, I want mixed abort and did-not-bite failures aggregated
   in the sweep's existing stable fixture order, so that one abort neither hides
   another failure nor makes concurrent diagnostics nondeterministic.
   **[reviewer-directed]** Line: `gpt-5.6-terra` / medium. The existing
   concurrency seam is known, while completeness of oracle output warrants the
   mid line.

## Implementation decisions

- **[defaulted] Existing module owner.** The canary package remains the single
  owner of fixture execution, Go-output classification, bite precedence, and
  aggregate diagnostics. No public package or fixture marker is added.
- **[reviewer-directed] Build routing.** The reviewer requires Sol for Codex
  work except code authorship. Every story is code authorship, and the project
  profile routes gate and conformance logic to the mid tier, so all four stories
  use `gpt-5.6-terra` / medium.
- **[defaulted] Existing test seam.** `Sweep` plus its injected runner is the
  acceptance seam. The runner result grows only the structured termination fact
  needed to distinguish a normal numeric exit from a spawn or signal abort;
  callers do not infer that fact a second time from prose.
- **[defaulted] Abort precedence.** A recognized panic, runtime-fatal marker, or
  process-level abort wins before both successful bite and did-not-bite
  classification. An aborted run cannot count as a bite merely because its
  partial output happened to contain `EXPECT`.
- **[defaulted] Go-output grammar.** Panic classification recognizes Go runner
  lines, after an optional outer gate-phase prefix, rather than searching for
  the substring `panic:` anywhere. It reports the deepest valid `Test...`
  failure header associated with the abort. A scoped behavior-owned bite falls
  back to its declared owning test; an unattributable process abort names its
  package or inner-gate scope and never invents a test.
- **[defaulted] Bounded diagnostics.** The failure class and attributable owner
  are rendered through the existing control-safe preview convention. Raw stack
  traces, arbitrary child output, and unbounded subprocess errors are not copied
  into the aggregate canary error.
- **[defaulted] Fail-safe unknowns.** A recognized abort with a malformed or
  missing test header remains an abort and reports an unknown test within the
  known package/scope. Losing attribution must not demote it to did-not-bite.
- **[defaulted] No semantic widening.** Fixture discovery, `EXPECT`, `TEST`,
  baselines, per-test narrowing, worker budgeting, phase routing, and the
  requirement for red plus `EXPECT` on a completed run are unchanged.
- **[defaulted] Static unsafe-slice arm rejected.** A conformance rule forbidding
  slices of “subject-reported output” is not included. Syntax alone cannot know
  value provenance or whether a dominating length guard makes a slice safe, so
  such a rule would be both incomplete and false-positive-prone. The generic
  runtime abort classifier covers the demonstrated failure without encoding a
  one-package heuristic.

## Testing decisions

- **[defaulted]** A good acceptance test drives `Sweep` with real fixture
  selection and an injected runner result, then observes the aggregate error.
  It does not call an output-parser helper directly.
- **[defaulted]** A helper subprocess re-executes the package test binary to
  produce authentic Go panic and subtest-panic output. This is the compatibility
  probe against the official producer's grammar; table fixtures cover malformed,
  prefixed, signaled, and hostile-output edges.
- **[defaulted] Historical red record.** On 2026-07-25 FT122's first gate
  produced `canary 'worktree-lifecycle-safety-bypassed' did not bite` after
  another runtime contract test panicked on an empty hash. The correct fixture
  was untouched. This is the observed red for the missing failure class; build
  entry promotes the recorded output shape into a focused test before changing
  production code.
- **[defaulted] Focused red-probe record.** On 2026-07-30 a throwaway test file
  in an isolated worktree drove `Sweep` with representative Go panic,
  runtime-fatal, signal, partial-output, hostile-output, and mixed aggregate
  outcomes. Every named FT129 target below exited nonzero on its desired
  assertion against the current generic classifier, with no setup or compile
  failure. The temporary file was removed and the worktree released after the
  red log was recorded.
- **[defaulted] Cheapest wrong implementation.** Searching for any `panic:`
  substring and replacing “did not bite” passes the happy case but fails on an
  incidental diagnostic containing that text, subtest attribution, partial
  output that already contains `EXPECT`, a signal with no panic text, and mixed
  aggregate ordering.
- **[defaulted] Focused verification:**
  `go test -count=1 ./internal/canary ./internal/subprocess`.
- **[defaulted] Gate:** `.bench/gate.sh` through `bench gate --fresh`.

### Seam diagram — canary fixture result classification

    trigger: dev or ship canary sweep grades one fixture
        │
        ▼
    selected fixture + RunCall ──▶ [ injected Runner ] ──▶ exit/output/termination
                                             │
                                             ▼
                                   [ abort then bite classifier ]
                                             │
                     ┌───────────────────────┼──────────────────────┐
                     ▼                       ▼                      ▼
              attributed abort       did not bite            fixture bit
                     └───────────────────────┬──────────────────────┘
                                             ▼
                                  stable aggregate sweep error
                    ◀ tests attach here: Sweep with real fixtures and
                      fake or helper-subprocess runner outcomes

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | a root-test panic is reported as an inner test abort naming that test, not as did-not-bite | `Sweep` over a behavior-owned fixture with authentic Go panic output | observed red 2026-07-30: `TestSweepAttributesPanickingTest` received generic did-not-bite instead of an abort naming `TestSubjectHash`; the 2026-07-25 FT122 gate is the production reproduction | the current implementation reaches its generic line, while a classifier that cannot name the test fails the exact diagnostic |
| 1 | a subtest panic reports the deepest failing `Test.../...` name | `Sweep` with helper-subprocess output from a panicking subtest | observed red 2026-07-30: `TestSweepAttributesPanickingSubtest` received generic did-not-bite instead of `TestLifecycle/empty_subject/hash` | taking the first parent `--- FAIL` header cannot satisfy the deepest-name assertion |
| edge of 1 | optional gate-phase prefixes are removed only for grammar recognition, and an incidental diagnostic containing `panic:` remains an ordinary did-not-bite | `Sweep` with prefixed and hostile runner-output table rows | observed red 2026-07-30: `TestSweepAbortGrammar` left the phase-prefixed panic generic while its incidental-text control correctly remained did-not-bite | a raw substring search misclassifies the hostile row, while a parser that ignores phase prefixes misses the real panic |
| edge of 1 | malformed or truncated panic output remains an abort with unknown test/package attribution and a bounded diagnostic | `Sweep` with a truncated panic block | observed red 2026-07-30: `TestSweepAbortGrammar/truncated` received generic did-not-bite | a parser requiring a complete stack or valid test header falls through to did-not-bite |
| edge of 1 | a Go runtime `fatal error:` marker is an inner test abort even when no panic marker exists | `Sweep` with authentic runtime-fatal output grammar | observed red 2026-07-30: `TestSweepAbortGrammar/runtime-fatal` received generic did-not-bite | an implementation that recognizes only `panic:` falls through to did-not-bite |
| edge of 1 | the attributed test is the deepest failing header associated with the abort, not an unrelated earlier failure or merely the last header in the stream | `Sweep` with parent/subtest and competing failure headers | observed red 2026-07-30: `TestSweepAbortGrammar/competing-headers` received generic did-not-bite instead of the panic-associated owner | a first-header or global-last-header heuristic names the wrong test while still appearing specific |
| 2 | a spawn failure and a signaled process are reported as process-level aborts using structured termination state, even with empty output | `Sweep` through the runner-result contract | observed red 2026-07-30: `TestSweepAttributesProcessAbort` received generic did-not-bite for an empty-output abnormal exit | output-only detection cannot classify the empty-output rows and therefore emits did-not-bite |
| edge of 2 | an ordinary numeric test exit is not a process-level abort | `Sweep` through the runner-result contract | already covered in part by `TestSweepReportsDidNotBiteWhenARedRunOmitsItsExpect`; the new process-abort table adds numeric exits 1 and 2 beside signal/spawn rows | treating every nonzero runner error as an abort relabels the existing ordinary failure |
| 3 | a completed red run without `EXPECT` retains the exact did-not-bite diagnostic | existing `Sweep` fixture contract | already covered by `TestSweepReportsDidNotBiteWhenARedRunOmitsItsExpect` | any precedence change that absorbs ordinary failures breaks the exact existing assertion |
| 3 | exit zero is did-not-bite even when output contains `EXPECT` | existing `Sweep` fixture contract | already covered by `TestSweepReportsDidNotBite` | an implementation that treats text match alone as a bite fails the current exact assertion |
| 3 | a completed red run containing `EXPECT` remains a successful bite with no sweep error | existing `Sweep` fixture contract | already covered by `TestSweepMaterializesFixtureAndRequiresTargetedBite` and the green concurrent-sweep fixtures | an implementation that attributes aborts but relabels every completed run as did-not-bite fails the existing positive assertion |
| edge of 3 | panic classification wins when partial output already contains `EXPECT` | `Sweep` with expected output followed by a panic | observed red 2026-07-30: `TestSweepAbortPrecedesBite/panic` returned no error because the partial output contained `EXPECT` | checking red plus substring first accepts the partial run and hides the panic |
| edge of 3 | process-abort classification wins when partial output already contains `EXPECT` | `Sweep` with expected output plus structured abnormal termination | observed red 2026-07-30: `TestSweepAbortPrecedesBite/process` returned no error because the partial output contained `EXPECT` | an output-only classifier accepts the partial run because no panic text exists |
| 4 | one panic, one process abort, and one ordinary did-not-bite are all returned in fixture order under concurrent execution | `Sweep` aggregate error | observed red 2026-07-30: `TestSweepOrdersMixedAbortFailures` returned three generic did-not-bite lines instead of the two attributed abort classes followed by did-not-bite | returning early hides rows, while append-on-completion makes output scheduler-dependent |
| edge of 4 | repeated classification of identical outcomes is byte-stable and carries no raw stack or unbounded child error | repeated `Sweep` runs over fixed runner results | observed red 2026-07-30: `TestSweepBoundsAbortDiagnostic` received generic did-not-bite instead of the bounded abort diagnostic its assertion also constrained | embedding raw output changes with stacks, temp paths, or control bytes and breaks the exact bounded result |

### Edge inventory

**[defaulted]** The following resolutions and **Won't handle** lines are the
no-map edge decisions for reviewer veto.

- **Error path:** root-test panic, subtest panic, runtime-fatal marker, spawn
  failure, signal termination, ordinary numeric failure, and a red run missing
  `EXPECT` each have rows above.
- **Empty or absent input:** empty output plus structured process abort is an
  abort; empty output plus an ordinary numeric exit is did-not-bite. An absent
  fixture marker remains owned by existing discovery checks.
- **Boundary values:** exits 0, 1, 2, and signal/spawn failure; zero/one/multiple
  failure headers; parent and deepest subtest names; `EXPECT` absent, present
  before abort, and present on a completed red are enumerated.
- **Malformed input:** a truncated panic block, an invalid test header, an
  incidental `panic:` diagnostic, and an optional phase prefix have coverage
  rows. Invalid owner text is bounded and control-sanitized before rendering.
- **Interrupted or partial state:** a signal with no output and a panic after
  partial expected output are both classified before bite evaluation.
- **Re-run idempotency:** classification is pure over the captured outcome and
  aggregate order remains fixture order rather than completion order.
- **Paths and names with spaces or glob characters:** no new path parser is
  introduced. Existing fixture materialization and argv tests remain the owner.
- **Control bytes in line sinks:** parsed child text never reaches a diagnostic
  raw; only a validated owner or bounded shared preview is rendered.
- **Missing trailing newline:** terminal panic and runtime-fatal markers are
  recognized at EOF.
- **Absent, empty, special, or dangling control records:** no new control record
  is added; existing `EXPECT` and `TEST` marker classification remains the owner.
- **Required tools missing:** a binary that cannot start is a process abort.
  Contract-package compile failures remain the earlier compile-stage diagnostic.
- **Invocation through every shipped surface:** `bench canary`, the dev gate,
  and ship-tier canary all call the same sweep classifier.
- **Deep cwd and non-TTY stdin:** the classifier consumes a captured runner
  outcome and neither resolves cwd nor prompts.
- **Host-backed filesystem pressure:** no deadline or timing assertion is added.
- **Compatibility probe:** the helper subprocess uses the installed Go test
  binary as the official producer of panic and subtest output.
- **Won't handle:** arbitrary application output that merely contains
  `panic:` — without Go runner grammar or structured process termination it is
  not evidence that the test binary aborted.
- **Won't handle:** recovering or continuing the aborted inner test binary —
  Go process recovery is outside the canary's authority; this feature attributes
  the abort and keeps the outer sweep running.

## Out of scope

- **[defaulted] Static detection of subject-derived slicing without a dominating length
  guard** — a separate source-analysis policy requiring value-provenance and
  control-flow decisions; a sound implementation is roughly 8 edits and 3 gate
  runs. The demonstrated diagnosis gap is covered generically here, and a
  textual or shallow-AST approximation is rejected rather than parked.
