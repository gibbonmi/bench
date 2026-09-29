# Return task-shaped reads and selected query results

Status: implemented

Decision source: `specs/session-context-efficiency/decisions/session-context-efficiency.md` (ready compiled map), with the reviewer-confirmed changes in the dated Amendment of 2026-09-28 section.

Verification log: 1 iteration(s) to accept — one Opus/high review accepted with two fence and fixture folds and four accounting nits. F1 removed the census fence, and F2 fixed the QU24 fixture kind and hash width. The author folded them under coordinator verification.

Coordinator: [Session context efficiency](../session-context-efficiency/spec.md)

## Problem

Agents repeat path and history queries, then carry unrelated fields and full bodies into context.
General call aggregation does not solve the oversized-result problem.

## Solution

Extend the existing worktree and spec-history owners with selected multi-target views.
Add focused raw-read guidance with explicit recovery paths.
The response bound that FT336 landed applies to both selected views, and the views add no numeric default.
Current domain-owned truncation policies remain authoritative.

## User stories

Line: opus / high.
Implementation-line reason: Shared history is the hardest chunk. The selected contracts are precise, but producer preservation needs new tests. The mid tier binds `opus` in this harness and `gpt-5.6-sol` in Codex.
Harder chunks: QU-C2.
The source fixes the outcome, but the seams require careful evidence and compatibility work.

1. As an agent, I want selected worktree identity, path, and state, so that one read supplies my next operation.
2. As an agent, I want per-target worktree errors, so that one missing target does not erase useful facts.
3. As an agent, I want duplicate identities collapsed, so that aliases do not repeat the same result.
4. As an agent, I want several exact spec histories, so that one owner answers the repeated intent.
5. As an agent, I want bounded history events with true omission measures, so that a selected query stays useful and auditable.
6. As an agent, I want a complete-history route per spec, so that I can recover omitted evidence.
7. As an agent, I want per-spec failures, so that one failed history does not erase another.
8. As an existing caller, I want additive selectors first, so that current invocations retain their behavior.
9. As an agent, I want invalid query grammar identified, so that a typo cannot change the selected scope.
10. As an agent, I want relevant raw-read defaults, so that file and command evidence starts with the needed material.
11. As an agent, I want explicit full reads to remain available, so that projection cannot hide needed evidence.
12. As a reviewer, I want authority boundaries retained, so that fewer calls cannot bypass approval or verification.
13. As an agent, I want hostile target text handled per target, so that output remains usable.
14. As an agent, I want unrequested output omitted, so that selected queries actually reduce context.

15. As an agent, I want a complete worktree-detail route, so that the selected view cannot hide additional worktree facts.

## Implementation decisions

Proposed worktree grammar: `bench worktree list --view paths --target <target> [--target <target>]...`.
The owner uses the existing assignment selector for paths, labels, IDs, and supported unique prefixes.
It does not use the authority-taking active-path resolver.
The fixed view emits `target`, `id`, `path`, `state`, and `error`.
Successful identities appear once in first-request order.
Unresolved operands retain one result per distinct operand.

The selected view includes `bench worktree list` in its contextual help as the complete-detail action.
Its output is one `worktrees` table and one help action.

The selected worktree route starts only when `--view` or `--target` is present.
That check comes before the bare list grammar, `usage.Parse(worktreeListGrammar)` in `ListCommand`.
Every other argument list keeps the bare grammar and its checked-in responses.

Proposed history grammar: `bench spec history --spec <slug-or-path> [--spec <slug-or-path>]... --limit <positive-count>`.
The initial selected mode requires an explicit event limit.
This caller-selected limit creates no new numeric default.
Normalized duplicate slugs appear once in first-request order.
The existing history producer retains exact slug matching, retire/delete classification, commit deduplication, and newest-first order.
No additional Git parser or batch command language joins the owner.

The selected history route starts only when `--spec` or `--limit` is present.
That check comes before the positional operand guard, `specArg` in `internal/spec/spec.go`.
Every other argument list keeps the `specArg` help, missing-operand, and unknown-flag responses.

The shared `History` producer keeps its complete ordered fact contract for roadmap context and positional history.
Selection and limits apply only after that producer returns.

The selected history output has one summary table, `histories`, with one row per requested spec.
Its fields are `target`, `slug`, `total_events`, `total_bytes`, `omitted_events`, `detail`, and `error`.
A separate event table, `history`, carries `slug`, `hash`, `date`, `kind`, and `subject`.
The summary counts complete serialized history before projection.
The `detail` cell contains the correctly quoted `bench spec history <slug>` command.
A successful empty history has zero events and no error.

`total_bytes` is the UTF-8 byte length of the complete positional history table for that spec.
That table is the exact `bench spec history <slug>` stdout: its `history[N]{hash,date,kind,subject}:` header, every event row, and each trailing newline.
One shared renderer produces that table for the positional path and for the count.

A failed history uses unknown counts and a per-spec error.
The selected view validates each history before it joins the combined event table.
A recovery command uses the end-of-options marker when the slug would otherwise select help or a flag.
If a second slug normalization changes the slug, the recovery command keeps the original operand.

Successful selected queries return exit 0.
A valid selection with any per-target failure returns exit 1 after all results render.
A malformed whole-command grammar returns exit 2 before the query starts.
An unsafe operand uses a stable ordinal in the result and a bounded explanatory error.
It never enters a shell command unquoted.
The existing TOON owner retains escaping policy.

The FT336 response bound already covers both routes, and this spec does not change it.
The `list` leaf in `cmd/bench/worktree_leaves.go` and the `spec` root in `cmd/bench/main.go` both declare `Bound: boundResponse`.
The bound values are `ResponseLines` 10 and `ResponseBytes` 4096 in `internal/bounds/bounds.go`.
A response within both values prints unchanged.
A larger response prints 4 head lines, 1 spill line, and 5 tail lines, and a private file keeps the complete output.
The selected views add no cap of their own.

Focused guidance covers raw file, Git, test, and shell reads.
It selects relevant sections, changed paths, failures, or task rows and gives the exact available full-detail route.
Independent archive and log reads can share one discovery step.
Polling, mutation, verification, approval, and publication keep their existing boundaries.
Build-then-run consolidation needs separate evidence and remains outside these tickets.

### Stream reference source

The 2026-09-22 wave built chunks QU-C1 to QU-C3 on branch `bench/assign/3458eb8d3cac8810295a65f9fdf866f2/910100df1a45f972d086be3f35818c63` at tip `194b7dba932a2881a085650789fa4c98d6478feb`.
That branch does not merge into `main`, because `main` moved 645 commits after the wave.
Each ticket names its stream reference files as `194b7dba:<path>`.
The fresh author reads those files as a starting point, not as an approved diff.
The ticket's `Writes:` line, its acceptance rows, and the current tree govern each port.

## Implementation chunks

Each ticket goes to a fresh author session under the authorship rule of `.bench/BENCH.md`.
Each ticket forms one named review chunk and one serial commit checkpoint.
The table orders independent tickets that share command inventory writes.
After each chunk, freeze its predecessor and current tips for Standards, Spec, and Coverage review.
The successor starts after accepted findings have current repair coverage.

| chunk / ticket | blocked by | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- | --- |
| QU-C1 / `1-select-worktrees.md` | none | Selected worktree facts | QU1, QU2, QU3, QU9, QU10, QU16, QU17, QU18, QU26 | TestSelectedWorktreeFacts and the remaining owned-row tests | no |
| QU-C2 / `2-select-histories.md` | none | Selected bounded histories | QU4, QU5, QU6, QU7, QU8, QU19, QU20, QU21, QU22, QU23, QU24, QU25, QU27 | TestSelectedHistoryPreservesProducer and the remaining owned-row tests | yes |
| QU-C3 / `3-guide-relevant-reads.md` | 1-select-worktrees.md, 2-select-histories.md | Focused raw-read guidance | QU11, QU12, QU13 | Standards and Spec inspection at QU11–QU13 | no |

The coverage map supplies the complete test inventory for each chunk's owned rows.
The final reconciliation checks every acceptance row and the integrated result.
The existing evidence checkpoints continue to block their implementation chunks.
The rules for execution-plan changes are in the platform guide.

### Claim verification

Before each ticket charge, the orchestrator verifies the current-code claims of that ticket against the charge source tip.
The claims are the cited symbols, the operand guards, the bound disposition, and the caller lists in Further notes.
The cost is about 12 targeted reads for each ticket and one `bench anchors` run for each guidance path of ticket 3.
A false claim stops the dispatch and returns to `/bench-write-spec`.

### Measurement after the first slice

After the QU-C1 checkpoint, the orchestrator measures the path task on the integration source.
The task asks for the path and state of one existing assignment.
It records the printed bytes and the follow-on calls for two routes: the selected view with one `--target`, and the bare `bench worktree list`.
The premise holds when the selected route needs no follow-on call and prints fewer bytes than the bare route.
If either condition fails, the build stops and returns to `/bench-write-spec` with the measurement.
The 2026-09-23 wave measurement is prior evidence: 300 bytes for the selected view against 16938 bytes for a 55-row inventory (`docs/research/parallel-implementation-wave.md:69-70`).

The same wave showed that a selected history of small specs can print more bytes than the positional histories: 640 against 499 (`docs/research/parallel-implementation-wave.md:78`).
The history premise is one call for several specs with true omission counts, not fewer bytes for small histories.
The measurement therefore grades only the worktree slice.

## Testing decisions

The command tests exercise the real owner with controlled files and records.
The ordinary Go test phase executes new package tests through its existing package census.
The conformance phase checks command inventories and ticket ownership.
No test-only helper crosses a package boundary.
Planned test names below identify required tests, not tests that already exist.

### Seam diagram

```text
named input -> existing command owner -> typed producer -> projected result
                       ^                       ^
                 command tests          controlled records
```

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| QU1 | 1 | The selected worktree view emits identity, path, and state for each resolved target | planned TestSelectedWorktreeFacts in internal/worktree | The valid state matrix includes recovered records, which an active-only path command loses |
| QU2 | 2 | A failed worktree target retains its own result beside successful targets | planned TestSelectedWorktreePartialFailure in internal/worktree | A whole-command early return drops a requested result |
| QU3 | 3 | Worktree aliases resolving to one identity emit one row at their first occurrence | planned TestSelectedWorktreeAliases in internal/worktree | Repeated labels and IDs cannot inflate successful rows or stored-path refusal rows |
| QU4 | 4 | Each selected spec receives newest-first history from the existing history owner | planned TestSelectedSpecHistories in internal/spec | A prefix match includes another spec or changes retire/delete classification |
| QU5 | 5 | The history view emits at most the explicit per-spec event limit | planned TestSelectedSpecLimit in internal/spec | For a 3-event history, limit 1 emits 1 row, limit 2 emits 2 rows, and limit 3 emits 3 rows, so a concatenated complete history of 3 rows fails limits 1 and 2 |
| QU6 | 5 | Each omitted history reports its complete event count | planned TestSelectedSpecOmissions in internal/spec | A clipped count cannot describe the omitted evidence |
| QU7 | 6 | Each selected spec result includes an exact complete-history command | planned TestSelectedSpecDetailRoute in internal/spec | A generic help action cannot recover the omitted target |
| QU8 | 7 | A failed spec target retains its own result beside successful targets | planned TestSelectedSpecPartialFailure in internal/spec | One Git failure cannot erase another target result |
| QU9 | 8 | The bare worktree view matches the baseline | planned TestSelectedWorktreesPreserveDefault in internal/worktree | Captured facts rendered through the TOON owner catch default changes, including numeric-looking identities |
| QU10 | 9 | Malformed worktree selection grammar returns usage at exit 2 | planned TestSelectedWorktreeGrammar in internal/worktree | Mixed modes and missing values cannot silently choose a different query |
| QU11 | 10 | Raw-read guidance selects relevant sections, paths, failures, or rows before full detail | review-owned: Standards axis reads the three guidance files | Advice to concatenate complete results defeats projection |
| QU12 | 11 | Raw-read guidance retains an explicit complete-detail route | review-owned: Standards axis reads the examples | A bounded example without recovery hides required evidence |
| QU13 | 12 | Guidance keeps polling and distinct authority boundaries as separate operations | review-owned: Spec axis compares `specs/session-context-efficiency/decisions/session-context-efficiency/tickets/4.md` and `specs/session-context-efficiency/decisions/session-context-efficiency/tickets/6.md` | A call-saving example cannot combine mutation with later approval |
| QU16 | 13 | A control-bearing worktree target failure does not hide other target results | planned TestSelectedWorktreeHostileTarget in internal/worktree | Later and repeated unsafe operands retain their request ordinals, and permitted stored-path controls round-trip without erasing other results |
| QU17 | 14 | Selected worktree output omits unrequested worktree rows | planned TestSelectedWorktreesExcludeOldOutput in internal/worktree | Presence-only assertions would let the old full output survive |
| QU18 | 15 | The selected worktree view names `bench worktree list` as its complete-detail action | planned TestSelectedWorktreeDetailRoute in internal/worktree | A result without the exact full inventory command fails the recovery contract |
| QU19 | 8 | The positional history view matches the baseline | planned TestSelectedHistoriesPreserveDefault in internal/spec | A differential history matrix catches accidental changes to complete history |
| QU20 | 9 | Malformed spec selection grammar returns usage at exit 2 | planned TestSelectedHistoryGrammar in internal/spec | An invalid history limit cannot silently choose another query |
| QU21 | 13 | A control-bearing spec target failure does not hide other target results | planned TestSelectedHistoryHostileTarget in internal/spec | One unsafe target cannot collapse the whole result into a render error |
| QU22 | 14 | Selected history output omits unrequested history bodies | planned TestSelectedHistoriesExcludeOldOutput in internal/spec | Presence-only assertions would let full histories survive beside selected output |
| QU23 | 4, 8 | The shared history producer preserves the baseline complete ordered fact sequence | planned TestSelectedHistoryPreservesProducer in internal/spec | Applying selection limits inside History would truncate roadmap context |
| QU24 | 5 | Each omitted history reports its complete serialized UTF-8 byte count | planned TestSelectedHistoryTrueBytes in internal/spec | For 3 events of 31 bytes each under a 36-byte header, the complete count is 36 + 3 × 31 = 129, and a count after limit 1 is 36 + 31 = 67 |
| QU25 | 7, 13 | An unrepresentable history subject fails only its selected spec result | planned TestSelectedHistoryHostileSubject in internal/spec | A combined table failure must not erase valid histories for other specs |
| QU26 | 14 | A selected worktree view of three resolved targets prints exactly 6 lines with no spill line through the bounded dispatcher | planned TestSelectedWorktreeWithinResponseBound in cmd/bench | The expected count is 1 table header + 3 rows + 1 help header + 1 help row = 6, which is under 10, so one unrequested row gives 7 and fails the exact count |
| QU27 | 14 | A selected history of two 3-event specs at limit 2 prints exactly 8 lines with no spill line through the bounded dispatcher | planned TestSelectedHistoryWithinResponseBound in cmd/bench | The expected count is 1 summary header + 2 summary rows + 1 event header + 2 × 2 event rows = 8, and an ignored limit gives 1 + 2 + 1 + 2 × 3 = 10 |

### Edge inventory

QU1–QU3 cover active, complete, cleanup-pending, recovered, absent, and ambiguous worktree selections.
Recovered fixtures carry valid recovery metadata from the assignment-state producer.
QU3 also selects two aliases of one assignment whose stored path is unrepresentable.

QU4–QU8 cover empty histories, duplicate slugs, retire-only, delete-only, mixed history, and per-target Git failures.
QU5 covers a limit of one, exact-boundary output, and omitted output.
QU7 also covers the operands `x.md.md` and `.md`, whose slug normalization is not stable.
QU10 covers missing worktree values and mixed worktree modes.
QU20 covers zero, negative, nonnumeric, and overflowing limits, missing values, and mixed history modes.

QU16 and QU21 cover spaces, control bytes, leading dashes, and command-shaped operands for kit and linked-repository callers.
QU16 varies unsafe operands across first, middle, and repeated positions, and it keeps the first occurrence ordinal.
Its stored-path matrix covers refused ESC and permitted tab, newline, and return beside valid results.

QU9 and QU19 preserve existing no-repository and default-view behavior through separate differential matrices.
QU9 keeps the captured baseline facts and renders expected tables through the TOON owner.
A deterministic numeric-looking identity makes quoting changes observable.

QU24 fixture rows use the kind `delete`, because a `retire` event needs a subject that ends with `spec-retire: <slug>`.
The fixture pins `core.abbrev 8`, a hash that does not look numeric, a 10-character date, and a one-byte subject, so each row has 31 bytes.
The test takes its expected count from the positional stdout of the same fixture and also asserts the row arithmetic.
No query acquires mutation authority or substitutes a package variable across subprocesses.

Won't handle: arbitrary command batches — the two existing query owners remain the callers.
Won't handle: a general Git query layer — the existing history producer and diff owner retain their responsibilities.
Won't handle: a general guidance rewrite — the named read examples remain the only guidance edits.
Won't handle: a selected-view cap — a selection past 10 lines or 4096 bytes takes the existing response-bound projection and spill file.

## Ownership fences

- `.agents/commands/bench-debug.md`
- `.agents/commands/bench-drain.md`
- `.agents/skills/bench-craft-cli/SKILL.md`
- `CHANGELOG.md`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `cmd/bench/main.go`
- `cmd/bench/selected_queries_test.go`
- `cmd/bench/worktree_leaves.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_debug_loop.go`
- `internal/anchors/registry_ft311_review_dispatch.go`
- `internal/anchors/registry_retained_workflow.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/spec/history.go`
- `internal/spec/history_command_test.go`
- `internal/spec/history_selected.go`
- `internal/spec/history_selected_test.go`
- `internal/spec/history_test.go`
- `internal/spec/spec.go`
- `internal/spec/spec_test.go`
- `internal/spec/testdata/selected-defaults.json`
- `internal/sanitize/sanitize.go`
- `internal/sanitize/sanitize_test.go`
- `internal/usage/worktree.go`
- `internal/worktree/list.go`
- `internal/worktree/list_actions_test.go`
- `internal/worktree/list_selected.go`
- `internal/worktree/list_selected_test.go`
- `internal/worktree/path.go`
- `internal/worktree/testdata/selected-defaults.json`
- `internal/worktree/unlanded_route_test.go`
- `reviews/session-context-queries.md`
- `tests/canary/package-core-guard/unrouted-subcommand`
- `tests/canary/row-next-grammar/token-table-lacks-kit-edit`
- `tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted`
- `tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns`
- `tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary`
- `tests/canary/workflow-guidance-anchors/debug-archaeology-anchor`
- `tests/canary/workflow-guidance-anchors/debug-phase1-stop-gate-softened`
- `tests/canary/workflow-guidance-anchors/debug-red-commit`
- `tests/canary/workflow-guidance-anchors/debug-reproduction-economics-deleted`
- `tests/canary/workflow-guidance-anchors/dg-1`
- `tests/canary/workflow-guidance-anchors/dg-1-contradiction`
- `tests/canary/workflow-guidance-anchors/dg-2`
- `tests/canary/workflow-guidance-anchors/dg-2-contradiction`
- `tests/canary/workflow-guidance-anchors/dg-3`
- `tests/canary/workflow-guidance-anchors/dg-3-blanket-ban`
- `tests/canary/workflow-guidance-anchors/dg-3-contradiction`
- `tests/canary/workflow-guidance-anchors/dg-4`
- `tests/canary/workflow-guidance-anchors/dg-4-contradiction`
- `tests/canary/workflow-guidance-anchors/dg-5`
- `tests/canary/workflow-guidance-anchors/dg-5-contradiction`
- `tests/canary/workflow-guidance-anchors/dg-5-dirty`
- `tests/canary/workflow-guidance-anchors/dg-5-dirty-contradiction`
- `tests/canary/workflow-guidance-anchors/dg-6`
- `tests/canary/workflow-guidance-anchors/dg-6-command`
- `tests/canary/workflow-guidance-anchors/dg-6-digest`
- `tests/canary/workflow-guidance-anchors/dg-6-dirty`
- `tests/canary/workflow-guidance-anchors/dg-6-surface`
- `tests/canary/workflow-guidance-anchors/drain-anchor`
- `tests/canary/workflow-guidance-anchors/drain-implement-now-commit`
- `tests/canary/workflow-guidance-anchors/drain-implement-now-per-spec-exception`
- `tests/canary/workflow-guidance-anchors/drain-implement-now-route`
- `tests/canary/workflow-guidance-anchors/drain-implement-now-row-fallback`
- `tests/canary/workflow-guidance-anchors/drain-implement-now-second-exception`
- `tests/canary/workflow-guidance-anchors/drain-roadmap-context-anchor`
- `tests/canary/workflow-guidance-anchors/drain-spec-history-anchor`
- `tests/canary/workflow-guidance-anchors/drain-split-board-detail-owner`
- `tests/canary/workflow-guidance-anchors/drain-split-board-retirement-pair`
- `tests/canary/workflow-guidance-anchors/drain-split-board-row-detail-owner`
- `tests/canary/workflow-guidance-anchors/implementation-retro-drain-anchor`

Reviewer disposition: the 2026-09-28 amendment awaits user spec and ticket sign-off.
The fence is the union of ticket writes and the review pickup.
The existing tests in `internal/spec/spec_test.go` and `internal/worktree/unlanded_route_test.go` are in the fence so that an author can add a case; their current expectations stay unchanged.

## Ticket graph

| Ticket | Blocked by | Delivered coverage |
| --- | --- | --- |
| [1. Select worktree path facts](tickets/1-select-worktrees.md) | none | QU1, QU2, QU3, QU9, QU10, QU16, QU17, QU18, QU26 |
| [2. Select bounded spec histories](tickets/2-select-histories.md) | none | QU4, QU5, QU6, QU7, QU8, QU19, QU20, QU21, QU22, QU23, QU24, QU25, QU27 |
| [3. Guide relevant raw reads](tickets/3-guide-relevant-reads.md) | 1-select-worktrees.md, 2-select-histories.md | QU11, QU12, QU13 |

## Out of scope

A general command batch language is a separate capability: approximately 8 edits, 2 gate runs.
A general guidance rewrite belongs to FT100: approximately 10 edits, 2 gate runs.
Build-then-run consolidation needs separate evidence: approximately 4 edits, 1 gate run.

## Further notes

### Source trace

| Source clause | Coverage |
| --- | --- |
| Ticket 3: bounded raw reads and Bench queries | QU5–QU7, QU11, QU12, QU18, QU24 |
| Ticket 4: selected multi-target queries and archive guidance | QU1–QU4, QU8, QU11, QU13 |
| Ticket 8: identity, path, state, and bounded histories | QU1–QU8, QU16–QU18, QU21, QU22, QU24–QU27 |
| Ticket 9: reviewer-approved numeric budgets | retired; FT336 landed the response bound in commit `5ec03981`, and its spec retired in commit `1952c9ce` |
| FT173: existing owners and full-detail routes | QU4, QU7, QU9, QU12, QU18–QU20, QU23 |

### Reader sweep and proof checklist

[Seam evidence](../session-context-efficiency/assets/seam-evidence.md) records producer and enforcement reads.
Cited symbols, verified at `4bdffd04`: `selectAssignment` at `internal/worktree/path.go:117`, `ListCommand` at `internal/worktree/list.go:30`, `History` at `internal/spec/history.go:37`, `historyCommand` at `internal/spec/history.go:150`, and `specArg` at `internal/spec/spec.go:224`.
Import edges: no new cross-package producer is required.
Source-row clauses and occurrences: the sole compiled map owns the clauses above.

Promised field labels: the two selected schemas appear in Implementation decisions.
Changed-function callers: the worktree `list` leaf calls `ListCommand`, and spec dispatch calls `historyCommand`.

The roadmap context parser also consumes every fact from the shared `History` producer at `internal/roadmap/context_parse.go:185`, under QU23.
Their local tests, public help tests, and command registry tests consume their grammar and output.
Copy survival: QU17 and QU22 fail if old full output survives beside the respective selected projection.

The bare worktree fixture family remains unchanged under QU9.
The file `internal/worktree/testdata/pre-disclosure-argv-pairs.json` stays unchanged, because no argument list in it carries `--view` or `--target`.
`TestListCommandCheckedInOldNewArgvCompatibility` and `TestListCommandHelpAndArgumentMatrix` in `internal/worktree/list_actions_test.go` grade that claim.
The bare help stays `usage: bench worktree list`.

The test callers of `ListCommand` are in `list_actions`, `landed`, `path_identifier`, `request_token`, `clean_landed_hostile`, and `unlanded_route` in the worktree package.
Every one of them calls the bare form, and QU9 keeps their responses.
`TestSpecSubcommandHelpUsesDeclaredUsage` and `TestSpecSubcommandUsageRefusals` in `internal/spec/spec_test.go` pin the positional history help and refusals, and QU19 keeps them.

The debug and drain history examples receive QU11–QU13 and exact ownership fences.
`bench anchors` reports 34 anchors on `bench-debug.md`, 39 on `bench-drain.md`, and 3 on `bench-craft-cli/SKILL.md`.
Ticket 3 keeps the drain needle ``use `bench spec history <slug>` for the shipped-row check`` from `internal/anchors/registry_data.go:156`.

The repository-wide sweep found no JavaScript or workflow schema reader.
The roadmap context reader remains unchanged because QU23 preserves its producer contract.
The AXI approved-query set remains unchanged.
The selected history view retains the spec owner's operational contract.

### Flagged additions

Repeated selectors, first-occurrence deduplication, and the explicit history limit are proposed grammar decisions for this sign-off.
QU1–QU10 and QU16–QU22 grade these additions.
This spec has no numeric-default ticket.
QU26 and QU27 are new in the 2026-09-28 amendment, and they grade the landed response bound.

### Amendment of 2026-09-28

Decision source for this amendment: the reviewer-confirmed conversation of 2026-09-28.
The reviewer chose to re-author the queries stream on `main` for FT348 and not to merge the stream.
The amendment keeps the stream's valid spec changes, which are these:

- the QU1, QU3, QU9, and QU16 edge work;
- the output table names and the byte definition;
- the recovery-operand rule;
- the added fence paths.

Stream ticket 5 was a review repair of QU-C1, and its acceptance now belongs to ticket 1.
The amendment drops stream ticket 4 and chunk QU-C4, because `b7667d8f` retired them.
It does not carry the stream review record; the new build writes its own.

The build plan amendment moves the plan to version 2 with one verification for each ticket.
It records one delegated run with an author limit of 1, and each author and repair session enters the plan before its dispatch.
The QU-C2 review expanded the ticket 2 fence to `internal/sanitize` and the two worktree selected files. One owner next to `LineSafe` then supplies the unsafe-operand ordinal for both selected views.

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"session-context-queries-full-20260928","orchestrator_session":"claude:session_01HvaChf55KS1vG4A5DV4mwW","author_limit":1,"assignments":{"1-select-worktrees.md":[{"session":"claude:bench-writer/scq-t1-author","assignment":"scq-t1-author","model":"opus","effort":"high","source":"12c7d857c6fe14e03c618c1584d3374371068277","native_ref":"claude:agent/scq-t1-author-20260928@12c7d857c6fe14e03c618c1584d3374371068277"},{"session":"claude:bench-writer/scq-t1-repair-1","assignment":"scq-t1-repair-1","model":"opus","effort":"high","source":"0e94692d3f2c73e7415a5a3d49e13b463ed440e2","native_ref":"claude:agent/scq-t1-repair-1-20260928@0e94692d3f2c73e7415a5a3d49e13b463ed440e2","predecessor":"claude:bench-writer/scq-t1-author","trigger":"user-directed","stopped":"The author returned its final report after record commit a5a527c8 and holds no write.","preserved":"0e94692d3f2c73e7415a5a3d49e13b463ed440e2"}],"2-select-histories.md":[{"session":"claude:bench-writer/scq-t2-author","assignment":"scq-t2-author","model":"opus","effort":"high","source":"a2c71f7a3165191bad3b2db13619f7888066d1a2","native_ref":"claude:agent/scq-t2-author-20260928@a2c71f7a3165191bad3b2db13619f7888066d1a2"},{"session":"claude:bench-writer/scq-t2-repair-1","assignment":"scq-t2-repair-1","model":"opus","effort":"high","source":"98aaa0e98e1bb4d863fb2786253e1cc036d7cd02","native_ref":"claude:agent/scq-t2-repair-1-20260928@98aaa0e98e1bb4d863fb2786253e1cc036d7cd02","predecessor":"claude:bench-writer/scq-t2-author","trigger":"user-directed","stopped":"The author returned its final report after record commit f4927b60 and holds no write.","preserved":"98aaa0e98e1bb4d863fb2786253e1cc036d7cd02"},{"session":"claude:bench-writer/scq-t2-repair-2","assignment":"scq-t2-repair-2","model":"opus","effort":"high","source":"7e51cb0acd4bff12d71d658e6bf5c8955aa87402","native_ref":"claude:agent/scq-t2-repair-2-20260928@7e51cb0acd4bff12d71d658e6bf5c8955aa87402","predecessor":"claude:bench-writer/scq-t2-repair-1","trigger":"user-directed","stopped":"The first repair session returned its final report after record commit 85286d10 and holds no write.","preserved":"7e51cb0acd4bff12d71d658e6bf5c8955aa87402"}],"3-guide-relevant-reads.md":[{"session":"claude:bench-writer/scq-t3-author","assignment":"scq-t3-author","model":"opus","effort":"high","source":"5fa6563815c5e4955c2f193bce62b651a8169c6e","native_ref":"claude:agent/scq-t3-author-20260928@5fa6563815c5e4955c2f193bce62b651a8169c6e"},{"session":"claude:bench-writer/scq-t3-repair-1","assignment":"scq-t3-repair-1","model":"opus","effort":"high","source":"6a88e82d4f2975a65df695f9323fa088fa32a444","native_ref":"claude:agent/scq-t3-repair-1-20260929@6a88e82d4f2975a65df695f9323fa088fa32a444","predecessor":"claude:bench-writer/scq-t3-author","trigger":"user-directed","stopped":"The author returned its final report after record commit ee78005b and holds no write.","preserved":"6a88e82d4f2975a65df695f9323fa088fa32a444"},{"session":"claude:bench-writer/scq-t3-repair-2","assignment":"scq-t3-repair-2","model":"opus","effort":"high","source":"78a50ee52784f1780c5ade1411807f4f50b878c5","native_ref":"claude:agent/scq-t3-repair-2-20260929@78a50ee52784f1780c5ade1411807f4f50b878c5","predecessor":"claude:bench-writer/scq-t3-repair-1","trigger":"user-directed","stopped":"The first repair session returned its final report after record commit 654e3d2f and holds no write.","preserved":"78a50ee52784f1780c5ade1411807f4f50b878c5"},{"session":"claude:bench-writer/scq-t3-reverify-1","assignment":"scq-t3-reverify-1","model":"opus","effort":"high","source":"551637e540a3a67536a6e55a5e4b0285c01beef3","native_ref":"claude:agent/scq-t3-reverify-1-20260929@551637e540a3a67536a6e55a5e4b0285c01beef3","predecessor":"claude:bench-writer/scq-t3-repair-2","trigger":"user-directed","stopped":"The second repair session returned its final report after record commit 0e77d29a and holds no write.","preserved":"551637e540a3a67536a6e55a5e4b0285c01beef3"}]}},"chunks":[{"id":"QU-C1","tickets":["1-select-worktrees.md"],"verification":[{"id":"1-worktree","command":"bench test --package ./internal/worktree --run TestSelected","ticket":"1-select-worktrees.md"},{"id":"1-bare-matrix","command":"bench test --package ./internal/worktree --run 'TestList|TestPath|TestCleanLanded|TestLanded|TestUnlanded|TestParallelCensusOnTheLiveTree|TestSerialSetStaysBelowTheCeiling|TestPackage'","ticket":"1-select-worktrees.md"},{"id":"1-command-route","command":"bench test --package ./cmd/bench --run 'TestSelected|TestCommandRegistryAXI|TestAXIRegistry|TestHelp|TestKept|TestWorktree'","ticket":"1-select-worktrees.md"},{"id":"1-state-probe","command":"bench probe internal/worktree/list_selected.go --swap 'string(selected.State)' --with 'string(intent.StateActive)' --package ./internal/worktree --run TestSelectedWorktreeFacts","probe":"swap","ticket":"1-select-worktrees.md"}]},{"id":"QU-C2","tickets":["2-select-histories.md"],"verification":[{"id":"2-history","command":"bench test --package ./internal/spec","ticket":"2-select-histories.md"},{"id":"2-command-route","command":"bench test --package ./cmd/bench --run 'TestSelected|TestHelp|TestAXIRegistry'","ticket":"2-select-histories.md"},{"id":"2-limit-probe","command":"bench probe internal/spec/history_selected.go --swap 'events[:limit]' --with 'events' --package ./internal/spec --run TestSelectedSpecLimit","probe":"swap","ticket":"2-select-histories.md"}]},{"id":"QU-C3","tickets":["3-guide-relevant-reads.md"],"verification":[{"id":"3-workflow","command":"bench test --check docs-currency-workflow","ticket":"3-guide-relevant-reads.md"},{"id":"3-skills","command":"bench test --check skills-index-command-adapters","ticket":"3-guide-relevant-reads.md"},{"id":"3-budgets","command":"bench test --check guidance-prose-budgets","ticket":"3-guide-relevant-reads.md"},{"id":"3-prose","command":"bench test --check prose-mechanics","ticket":"3-guide-relevant-reads.md"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/session-context-queries/spec.md"},{"id":"worktree","command":"bench test --package ./internal/worktree"},{"id":"history","command":"bench test --package ./internal/spec"},{"id":"commands","command":"bench test --package ./cmd/bench"},{"id":"workflow","command":"bench test --check docs-currency-workflow"}]}
```
