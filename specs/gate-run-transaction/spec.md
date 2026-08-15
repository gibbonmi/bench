# gate-run-transaction

Status: staged

Decision source: specs/gate-run-transaction/decisions/deepening-2026-08.md (compiled ready map, reviewer-resolved 2026-08-15; tickets #4, #6, #7 govern this spec)

Verification log: spec 1 + tickets 1 iteration(s) to accept — reviewer-capped at one round per loop (Codex `gpt-5.6-sol`, high, read-only, `--reviewer` override); loop 1 returned 8 blocking findings, all folded (GC5 lock holder isolated from any in-flight run, GC6 seeded green + `--fresh`, GT4/GT5 persistence-failure rows and GT6 structural predicate added, GR2 enumerated per narrow class, GR5 validator rows, GC11 abandoned-pending row replacing a false already-covered cite); loop 2 returned 5 blocking, all folded except the per-row split of the two characterization tickets, kept fixture-grouped as a flagged reviewer call. Acceptance was by cap, not by a clean re-review.

## Problem

The gate run — the oracle every Bench verdict rests on — is one 151-line function
(`executeSubjectWithRunBinary`) that owns subject acceptance, two reuse checks, lock
acquisition and the owner file, the under-lock drift check, run-binary selection,
dispatch, and four persistence branches (pending, timeout, green-retain, red-invalidate)
each with hand-rolled rollback. The injection seam that once made it testable
(`gateEngine`, twelve methods) has had exactly one implementation since the
2026-08-09 branch-native test commit deleted its fault engine, and four wrappers plus
`engineEvaluation` now have no caller. Verdict record-class validity is spelled three
ways in `verdict.go`: an unread map (`readyFieldClasses`, whose comment cites a file, a
function, and a test that do not exist), the live `switch` in `validateRecordBytes`,
and an inline Pending field list. Nothing in the tree today observes the executor's
persistence branches: no test at any seam drives lock refusal, the pending record, the
timeout record, drift during execution, or retain/invalidate. Understanding what
survives a crash mid-run means reading two files end to end.

## Solution

Behind the unchanged public seam (`Execute`, `ExecuteReusingFreshGreen`, `ExecuteTree`,
`RunCommand`, `Command`, `RunAndRecord`, `RunAndRecordContext`, `Inspect`,
`InspectTree`, `ValidateProjectGreen`, `ComposedGreen`, `PinCommand`) the run becomes
one deep module — "grade this subject and leave exactly one durable record" — that owns
the lock, the owner file, the under-lock drift check, the pending→terminal record pair,
retain/invalidate pairing, and the crash-safe replace. Callers learn the outcome and the
inspection, nothing about locks or temp-renames. The `gateEngine` interface,
`productionGateEngine`, and the four dead wrappers are deleted (map ticket #6);
`engineEvaluation` goes with them as a consequence — it exists only to wrap an engine.
The public seam is enumerated above because "unchanged" quantifies over it (`craft-spec`);
the list adds no function the tree does not already export.
Verdict record classes become one registry — one row per class carrying its name, its
exact field set, and its validator — read by the loader, by `narrowVerdictReason`, and
by the reuse refusal, with an independently authored guard test. Before anything moves,
characterization tests pin today's observable behavior at the public seam so the
refactor's exit test has teeth: the pre-existing suite and the characterization rows
pass with test logic unmodified.

## User stories

1. **Characterize the gate run at its public seam.** As the reviewer I can see, from
   tests that drive `gate.Execute` / `gate.RunCommand` / `gate.Inspect` against a fixture
   repository whose `.bench/gate.sh` is a script, every outcome the run promises today:
   green retains and records, red invalidates and records, a fresh green is reused
   without re-recording, `--fresh` forces the run, a held lock refuses with the owner
   diagnostic and demotes a reusable green to pending, a timeout records `timeout` and
   invalidates, the oracle sees the pending record while it runs, an abandoned pending
   record reloads as pending and does not block the next run, a subject that moves
   during execution is not recorded as a verdict, and a missing gate exits 3.
   Line: opus / medium. Gate logic routes mid effort per the profile; the behaviors are
   fixed by the current code, so the work is careful observation, not design.
2. **Run the gate as one transaction.** As a maintainer I read one production file in
   `internal/gate` that is the only one referencing the lock path, the owner path, the
   pending record constructor, `retainGreen`, `invalidateEvidence`, and the terminal
   record construction — the run's whole lock → owner file → drift check → pending →
   oracle → terminal ± evidence sequence — while the public entry points contain none
   of those calls; I find no `gateEngine`, no `productionGateEngine`, and no
   `executeWithEngine`/`executeWithEngineAfterAcquire`/`executeSubjectWithEngine`/
   `newWorkingTreeEvaluation`; every exported function keeps its signature and every
   pre-existing test plus story 1's rows pass unmodified.
   Line: opus / medium. A behavior-preserving move at the oracle; the exit test grades
   it, and the profile routes gate logic mid.
3. **Enumerate verdict record classes once.** As a maintainer I find one registry with
   exactly five rows — full, partial, check-partial, combined-partial, pending — each
   carrying its name, its exact field set, and its validator (the per-class check the
   loader runs after the field set matches); the loader selects a class through the
   registry and runs the row's validator, `narrowVerdictReason` derives today's
   `partial verdict` reason from the row rather than from ad-hoc predicates, the reuse
   refusal reads that reason, no `*ReadyFields` package variable or inline field list remains, the
   `readyFieldClasses` comment's phantom references are gone, and an independently
   authored guard test names the five classes and their field sets so a dropped or
   mutated row is red. `partialTestRecord`/`fullTestRecord` are deleted.
   Line: sonnet / medium. Fully gate-observable through the loader at the public
   `Inspect` seam; mechanical once the shape is fixed.

## Implementation decisions

- **One deep run module, one durable record.** The transaction module owns, in
  order: subject acceptance and refusal, the pre-lock reuse answer, git-dir and lock
  acquisition, the owner file, the under-lock drift re-validation, the under-lock reuse
  answer, run-binary selection, the pending record, the oracle dispatch, and exactly one
  terminal outcome that pairs evidence retain-or-invalidate with the terminal record and
  restores the pending record if the terminal replace fails. Public entry points become
  thin adapters over it. Its interface is what a caller must know today: the outcome
  (`Result`) and the inspection — no lock, owner file, temp-rename, or rollback is
  visible.
- **Delete the hypothetical seam.** `gateEngine` and `productionGateEngine` go; the
  clock, filesystem, lock, and subject calls they wrapped become direct calls. One
  adapter is a hypothetical seam (`craft-seams`); a future fake is a future decision.
  `engineEvaluation`, `newEngineEvaluationAtKit`, and the four zero-caller wrappers go
  with it. `durableReplaceWithEngine`, `inspectWithEngine`, `operationalWithEngine`,
  and `persistInterruptedIfGreen` lose their engine parameter.
- **Record-class registry.** One ordered registry, one row per class: name, exact
  sorted field set, validator (map ticket #7). The loader selects the row whose field
  set the record carries, then runs that row's validator (`validatePartition`,
  `validateCheckPartition`, both, the pending checks, or the full class's none); the
  four `*ReadyFields` variables and the inline Pending field list are replaced by rows.
  How `narrowVerdictReason` derives the existing `partial verdict` reason from the
  selected row (a per-row attribute or a name predicate) is spec-writer discretion —
  the observable reason strings do not change. Narrow classes remain readable exactly as
  today — no class is dropped or renamed. Adding a class is one row.
- **Guard test under ADR 0006.** The guard test's expectation is an independently
  authored literal — five names with their field sets — not derived from the registry;
  its independence is what makes a dropped row or a mutated field set red, and that red
  is demonstrated and recorded in the ticket.
- **Exit-test rule (map ticket #4).** Every ticket lands with the pre-existing test
  suite and story 1's characterization tests passing with test logic unmodified;
  mechanical renames are the only permitted test edit. A changed assertion reverts the
  move and reroutes as a behavior change. A defect discovered mid-move is parked with
  `bench idea` and fixed as its own commit, never bundled.
- **No new behavior.** Stories 2 and 3 add no observable behavior. Story 1 adds tests
  only. The one in-package knob the characterization uses is the existing `gateTimeout`
  package variable (a time boundary), overridden in-test and restored.

## Testing decisions

- A good test drives a real fixture repository (`internal/gittest`) with a scripted,
  executable `.bench/gate.sh`, calls the public seam, and observes exit codes, stderr
  lines, `Inspect` state, and the file effects the script itself can witness (its own
  run marker; the pending record it can read mid-run; permissions it flips before
  exiting). Behavior switches between runs go through `.gitignore`d sentinels so the
  subject stays byte-identical; the lock is held from outside a run by a re-exec'd
  helper process (prior art `internal/freshness/freshness_test.go`). It never reaches
  into the transaction module.
- Seams receiving tests: the public gate seam (`gate.Execute`, `gate.RunCommand`,
  `gate.Inspect`) — prior art `internal/landing/landing_test.go` (drives a scripted
  gate through `authorization.Authorize`) and `internal/status/status_test.go` (writes
  and reads verdict records); the loader through `gate.Inspect` over a hand-written
  record file — prior art `internal/status/status_test.go:177`.
- Gate seam observing the feature: the ordinary `test` phase (in-package tests in
  `internal/gate`) plus every existing consumer test that crosses the public seam
  (`internal/landing`, `internal/worktree/land_test.go`, `internal/preprelease`,
  `internal/status`, `internal/stophook`, `internal/shift`, `cmd/bench`).

### Seam diagram

    trigger: bench gate | bench gate-run | gated commit | shift iteration | landing
        │
        ▼
    root, mode, ctx  ──▶  [ gate run transaction: accept → reuse? → lock+owner → drift → ]  ──▶  Result{GateExit, ActionExit, Inspection}
                          [ pending record → oracle → terminal record ± evidence         ]       + <git-dir>/bench-last-gate
                          [ ─ registry: class row → loader / narrowVerdictReason ─       ]       + bench-gate-evidence/<subject>
                              ◀ tests attach here: fixture repo + scripted .bench/gate.sh; call Execute/RunCommand/Inspect;
                                assert exit, stderr, Inspect state, script-witnessed files, record bytes

### Acceptance coverage map
| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| GC1 | 1 | a green script gate: `Execute` returns `GateExit 0`/`ActionExit 0`; `Inspect` reports `Ready`/`green` with `ReusableGreen` true for the unchanged tree; a second `Execute` prints `gate: green (fresh verdict reused for this tree)`, does not run the script again (script-owned run counter stays 1), and leaves the record's `recorded_at` byte-identical | public gate seam | characterization: green on the current tree; red demonstrated by mutating the reuse check to skip `reusableEvidence` (script runs twice) | pins retain + record + reuse-without-rerecord, the property the transaction move is most likely to reorder |
| GC2 | 1 | a red script gate (exit 7): `Execute` returns `GateExit 7`/`ActionExit 7`; `Inspect` reports `Ready`/`red`, `ReusableGreen` false; a prior green's evidence for the same subject no longer authorizes reuse (a following `Execute` runs the script again) | public gate seam | characterization; red demonstrated by mutating the red branch to skip `invalidateEvidence` (the following run reuses a stale green) | pins the invalidate half of the retain/invalidate pair |
| GC3 | 1 | `RunCommand([--fresh, root])` after a fresh green runs the script again (counter 2) | public gate seam | characterization; red demonstrated by mutating `forceRun` handling to fall through to reuse | pins the operator's only escape from a reusable green |
| GC4 | 1 | with the lock held by an in-flight run (the script polls for a fixture-owned sentinel file before exiting), a concurrent `Execute` returns `ActionExit 1`, stderr contains `gate execution already in progress` and `gate owner: pid <n> (alive)`, and the script's run counter shows the second call never ran it | public gate seam | characterization; red demonstrated by mutating the acquire-failure branch to fall through | pins refusal-not-queueing under contention and the owner diagnostic |
| GC5 | 1 | with a reusable green retained by a completed run, and the lock then held by a separate fixture-owned holder process (re-exec of the test binary, prior art `internal/freshness/freshness_test.go`) that has written no record, a refused `Execute` leaves `Inspect` in `Pending` and never runs the script; before the refused call `Inspect` was `Ready`/`green` reusable | public gate seam | characterization; red demonstrated by deleting the `persistInterruptedIfGreen` call on the acquire-failure branch (`Inspect` stays `Ready`/`green`) | pins the demotion the doc comment on `ExecuteReusingFreshGreen` contrasts against, with the holder isolated so no in-flight run writes the pending record itself |
| GC6 | 1 | after a green run retained evidence, an ignored (`.gitignore`d, so the subject is unchanged) sentinel switches the script to sleep; with `gateTimeout` overridden to a sub-second value, `RunCommand([--fresh, root])` returns `124`, stderr contains `gate: timeout`, `Inspect` reports `Ready`/`timeout`; with the sentinel removed, an ordinary `Execute` runs the script again (counter 3) instead of reusing the earlier green | public gate seam | characterization; red demonstrated by deleting `invalidateEvidence` from the timeout branch (the ordinary run reuses the retained green) and, separately, by mutating the timeout branch to record `green` | pins the timeout record class and the invalidation that keeps a timed-out subject from inheriting its earlier green |
| GC7 | 1 | while the script runs it can read `<git-dir>/bench-last-gate` and the bytes it copies aside parse as `state: pending` with `owner_pid` equal to the test process pid and `tree` equal to the subject `Inspect` later reports | public gate seam (script-witnessed) | characterization; red demonstrated by moving the pending replace after the oracle dispatch | pins the pending-before-oracle order that makes a crash mid-run visible to `bench status` |
| GC8 | 1 | a script that appends to a tracked file: `Execute` returns `ActionExit 1` with the script's own exit as `GateExit`, stderr contains `gate subject changed during execution`, and `Inspect` does not report `ReusableGreen` | public gate seam | characterization; red demonstrated by mutating `sameSubject(plan, after)` to true | pins the drift refusal that keeps a moved tree from inheriting a verdict |
| GC9 | 1 | a fixture with no `.bench/gate.sh` and no `BENCH_GATE`: `Execute` returns `GateExit 3`/`ActionExit 3` and stderr contains `no gate found` | public gate seam | characterization; red demonstrated by mutating the `None` branch's exit | pins the absent-oracle exit the shell wrapper relies on |
| GC10 | 1 | a green prospective tree run through `authorization.Authorize` is `Green` and `InspectTree` reports `ReusableGreen` | authorization seam | already covered — `internal/worktree/land_test.go:65-73`, `internal/landing/landing_test.go` | the prospective path is the second caller of the transaction; its existing tests are part of the exit test |
| GC11 | 1 | a hand-written valid `state: pending` record whose `owner_pid` names an exited process: `Inspect` reports `Pending`; a following `Execute` runs the script (counter 1) and ends `Ready`/`green` | public gate seam | characterization; red demonstrated by mutating the loader to report `Invalid` for a pending record whose owner is not alive | pins the abandoned-run reload path — the process-boundary case a fresh `bench status` meets after a crash |
| GT1 | 2 | after the move, every GC row and every pre-existing test in the repository passes with test logic unmodified (renames only); a change to any test assertion is not permitted — the move reverts and the change reroutes as a behavior change through `/bench-write-spec` or `/bench-debug` (map ticket #4) | ordinary `test` phase + `git diff -- '*_test.go'` shows only rename hunks | exit test — a red here reverts the move | the refactor's only oracle is that nothing observable changed |
| GT2 | 2 | `gateEngine`, `productionGateEngine`, `engineEvaluation`, `newEngineEvaluationAtKit`, `executeWithEngine`, `executeWithEngineAfterAcquire`, `executeSubjectWithEngine`, `newWorkingTreeEvaluation` no longer exist in `internal/gate` | source (`rg` in review) | not TDD-able — structural deletion; verified by review with `rg` and by the build | a surviving identifier means the seam was kept, contradicting map ticket #6 |
| GT3 | 2 | after a green run retained evidence, a sentinel-switched `RunCommand([--fresh, root])` whose script strips the evidence directory to `0500` before exiting 0 does not write a green terminal record: `ActionExit 1`, stderr `gate evidence persistence failed`, `Inspect` not `ReusableGreen` | public gate seam (`ensureEvidenceDir` requires exactly `0700`; capability-skipped when privileged; fixture restores mode in cleanup) | characterization; red demonstrated by writing the terminal record before the retain | pins the retain half of the pairing the deep module exists to own |
| GT4 | 2 | after a green run retained evidence, a sentinel-switched `RunCommand([--fresh, root])` whose script makes `bench-gate-evidence` unwritable before exiting 7 exits operationally: `ActionExit 1`, stderr `gate evidence invalidation failed`, `Inspect` not `ReusableGreen` | public gate seam (`os.Remove` of the retained file fails; capability-skipped when privileged; fixture restores mode in cleanup) | characterization; red demonstrated by dropping the invalidation error check on the red branch | pins the invalidate half of the pairing |
| GT5 | 2 | after a green run retained evidence, a sentinel-switched `RunCommand([--fresh, root])` whose script makes `<git-dir>` itself unwritable before exiting 0 (the evidence subdirectory stays writable, so retain succeeds and only the terminal replace fails) leaves the pending record in place: `ActionExit 1`, stderr `gate final persistence failed`, `Inspect` reports `Pending` | public gate seam (capability-skipped when privileged; the fixture restores permissions in cleanup) | characterization; red demonstrated by dropping the pending restore after the failed terminal replace | pins crash-safe replace with rollback to pending — the property that keeps a half-written verdict from reading as ready |
| GT6 | 2 | exactly one production file in `internal/gate` references `bench-gate.lock`, `bench-gate-owner`, `interruptedRecord(`, `retainGreen(`, `invalidateEvidence(`, and the `Status: "timeout"` / terminal `verdictRecord` construction; the exported entry points (`Execute`, `ExecuteReusingFreshGreen`, `ExecuteTree`, `RunCommand`, `RunAndRecord`, `RunAndRecordContext`) reference none of them and each is at most a thin adapter that resolves roots/mode/log and calls the run module | source (`rg` in review) | not TDD-able — structural predicate; verified by review with `rg` over `internal/gate/*.go` excluding tests | the cheapest wrong implementation deletes the symbols and leaves the orchestration split across `gate.go` and `verdict.go`; this predicate makes that red |
| GR1 | 3 | a record file carrying `state: ready` with fields from two classes (full + `executed` without `skipped`/`skip_evidence`, or partial + `check_executed` alone) makes `Inspect` report `Invalid`; a record with exactly the full set, exactly the partial set, exactly the check-partial set, and exactly the combined set is `Ready` (all five fixtures asserted) | loader via `gate.Inspect` | characterization; red demonstrated by mutating the registry lookup to accept any superset of the full set | pins the alternatives-not-spectrum rule the registry must preserve |
| GR2 | 3 | each of the three narrow classes — a valid partial record, a valid check-partial record, and a valid combined-partial record — inspected for its own tree yields a non-reusable inspection whose reason is `partial verdict` (three fixtures asserted); a full-class record yields `ReusableGreen` | loader + `narrowVerdictReason` via `gate.Inspect` | characterization; red demonstrated three times, once per row, by making that row report the full class's reason (each single mutation must be red on its own fixture) | pins the reporter's and the reuse refusal's read of the registry for every narrow class, not just one |
| GR3 | 3 | the guard test's independently authored literal — five class names each with its exact field set — matches the registry; removing a row or adding a field to one row's set is red | guard test in `internal/gate` | first-write red: registry stubbed empty, the literal comparison fails; then the omission mutation (drop the `check-partial verdict` row) is applied and observed red, recorded in the ticket | the registry is the one enumeration; without an independent expectation nothing detects a lost class |
| GR4 | 3 | `readyFieldClasses`, `fullReadyFields`, `partialReadyFields`, `checkPartialReadyFields`, `combinedPartialReadyFields`, `partialTestRecord`, `fullTestRecord`, and the strings `storeRecordClasses`, `record_classes.go`, `TestVerdictReadyFieldsAreAllRegistered` no longer occur in `internal/gate` | source (`rg` in review) | not TDD-able — structural; verified by review with `rg` | a surviving spelling is a second enumeration or a phantom reference |
| GR5 | 3 | a partial-class record whose `executed` and `skipped` lists share a component is `Invalid` (`validatePartition`); a check-partial record whose `check_inherited` names a component with no `check_evidence` entry is `Invalid` (`validateCheckPartition`); a pending record with `owner_pid: 0` is `Invalid` | loader via `gate.Inspect` | characterization; red demonstrated per class by replacing that row's validator with a no-op | proves the loader runs the row's validator — a registry that carries validators nothing calls passes GR1 and GR3 |

### Edge inventory
- **error path** — lock refusal (GC4/GC5), retain failure (GT3), invalidate failure (GT4), terminal replace failure (GT5), red oracle (GC2), absent oracle (GC9).
- **empty/absent input** — no gate found (GC9); absent record file is `Unavailable` (already covered, `internal/status/status_test.go`).
- **boundary values** — timeout (GC6); exact field sets (GR1).
- **malformed input** — mixed-class record (GR1); duplicate JSON names and bad times remain covered by the loader's existing checks, unchanged by this spec (`validateRecordBytes` before the class switch).
- **interrupted / partial state** — pending record visible mid-run (GC7); an abandoned pending record reloads as `Pending` and does not block the next run (GC11); a failed terminal replace rolls back to pending (GT5).
- **re-run idempotency** — second run reuses without re-recording (GC1); `--fresh` re-runs (GC3).
- **process-boundary lifecycle** — GC7 crosses it (the oracle script, a child process, reads the record the parent wrote); GC1's second run in the same process reads what the first wrote through the store.
- **hostile environment** — evidence directory or git-dir made unwritable mid-run (GT3/GT4/GT5, capability-skipped when privileged).
- **Won't handle:** subject drift between the pre-lock acceptance and the under-lock re-validation — reproducing it needs an in-process race hook the public seam does not offer; the moved code keeps its `gate subject changed before execution` exit and review re-derives it from the diff.
- **Won't handle:** lock file unopenable (git-dir unwritable) — privilege-dependent, the same capability class GT3 skips under; the operational exit stays.
- **Won't handle:** run-binary selection failures (`gate Bench executable unavailable`, selection timeout) — they need the phase-table gate (`"$bench" gate-phases`) and a broken build; the branch moves verbatim and the kit's own gate exercises the success path on every run.
- **Won't handle:** a class registered without a row in the guard literal — that is the guard test's job (GR3), not a runtime edge.

## Ownership fences

One writer at a time; the tickets are serial (story 1 → story 2 → story 3), so no fence
has two writers concurrently. Reviewer disposition: approve as listed.

- Story 1: `internal/gate/` (new characterization test file(s) only; no non-test edit).
- Story 2: `internal/gate/` (production files; test edits limited to mechanical renames).
- Story 3: `internal/gate/` (`verdict.go`, `verdict_test_records_test.go`, new registry
  and guard-test files).
- Every ticket: `specs/gate-run-transaction/` for status and ticket checkboxes.

## Out of scope

- **One green-marker reader** (map candidate 6) — a light-path ticket after this spec:
  3 edits, 2 gate runs.
- **Dropping the narrow (partial) record classes** now that no writer produces them —
  a behavior change to the loader and `bench status`; separate decision. 4 edits, 2 gate
  runs once decided.
- **A fake gate engine / injectable clock** — the deleted seam returns only when a second
  adapter exists (`craft-seams`).
- **FT162/FT185 evidence-subject unification** across status/handoff/landing — roadmap
  work with its own rows.
- **Structure budget** for `internal/gate/` (21 files, `gate.go` over budget) — a
  reviewer grant or split decision, not this spec's; the deepening should not add net
  files beyond the tests and one run module.
