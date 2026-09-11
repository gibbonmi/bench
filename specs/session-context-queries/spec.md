# Return task-shaped reads and selected query results

Status: staged

Decision source: `specs/session-context-efficiency/decisions/session-context-efficiency.md` (ready compiled map).

Verification log: 0 iteration(s) to accept — pending the requested Sol/high review.

Coordinator: [Session context efficiency](../session-context-efficiency/spec.md)

## Problem

Agents repeat path and history queries, then carry unrelated fields and full bodies into context.
General call aggregation does not solve the oversized-result problem.

## Solution

Extend the existing worktree and spec-history owners with selected multi-target views.
Add focused raw-read guidance with explicit recovery paths.
Apply new numeric defaults only after the measurement report and reviewer budget decision.
Current domain-owned truncation policies remain authoritative.

## User stories

Line: gpt-5.6-terra / high.
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
13. As a reviewer, I want measured numeric policies, so that a build cannot invent new limits.
14. As an agent, I want recoverable bounded defaults, so that approved limits preserve the route to complete evidence.
15. As an agent, I want hostile target text handled per target, so that output remains usable.
16. As an agent, I want unrequested output omitted, so that selected queries actually reduce context.

## Implementation decisions

Proposed worktree grammar: `bench worktree list --view paths --target <target> [--target <target>]...`.
The owner uses the existing assignment selector for paths, labels, IDs, and supported unique prefixes.
It does not use the authority-taking active-path resolver.
The fixed view emits `target`, `id`, `path`, `state`, and `error`.
Successful identities appear once in first-request order.
Unresolved operands retain one result per distinct operand.

Proposed history grammar: `bench spec history --spec <slug-or-path> [--spec <slug-or-path>]... --limit <positive-count>`.
The initial selected mode requires an explicit event limit.
This caller-selected limit creates no new numeric default.
Normalized duplicate slugs appear once in first-request order.
The existing history producer retains exact slug matching, retire/delete classification, commit deduplication, and newest-first order.
No additional Git parser or batch command language joins the owner.

The selected history output has one summary row per requested spec.
Its fields are `target`, `slug`, `total_events`, `total_bytes`, `omitted_events`, `detail`, and `error`.
A separate event table carries `slug`, `hash`, `date`, `kind`, and `subject`.
The summary counts complete serialized history before projection.
The `detail` cell contains the correctly quoted `bench spec history <slug>` command.
A successful empty history has zero events and no error.

Successful selected queries return exit 0.
A valid selection with any per-target failure returns exit 1 after all results render.
A malformed whole-command grammar returns exit 2 before the query starts.
An unsafe operand uses a stable ordinal in the result and a bounded explanatory error.
It never enters a shell command unquoted.
The existing TOON owner retains escaping policy.

Focused guidance covers raw file, Git, test, and shell reads.
It selects relevant sections, changed paths, failures, or task rows and gives the exact available full-detail route.
Independent archive and log reads can share one discovery step.
Polling, mutation, verification, approval, and publication keep their existing boundaries.
Build-then-run consolidation needs separate evidence and remains outside these tickets.

The final ticket waits for the measurement child's budget evidence and an explicit reviewer decision.
Before that ticket starts, spec authoring records each numeric value, surface, unit, omission metadata, and full-detail route.
That authoring pass also closes any additional owner fence required by the selected surfaces.
The build cannot amend its own policy or expand its own fence.
Until that checkpoint, the final ticket is not on the executable frontier.
Existing budgets remain active and do not require reapproval.

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
| QU1 | 1 | The selected worktree view emits identity, path, and state for each resolved target | planned TestSelectedWorktreeFacts in internal/worktree | Calling the active-only path command loses non-active state |
| QU2 | 2 | A failed worktree target retains its own result beside successful targets | planned TestSelectedWorktreePartialFailure in internal/worktree | A whole-command early return drops a requested result |
| QU3 | 3 | Worktree aliases resolving to one identity emit one row at their first occurrence | planned TestSelectedWorktreeAliases in internal/worktree | Repeated labels and IDs cannot inflate the result set |
| QU4 | 4 | Each selected spec receives newest-first history from the existing history owner | planned TestSelectedSpecHistories in internal/spec | A prefix match includes another spec or changes retire/delete classification |
| QU5 | 5 | The history view emits at most the explicit per-spec event limit | planned TestSelectedSpecLimit in internal/spec | Concatenating complete histories violates the requested bound |
| QU6 | 5 | Each omitted history reports its complete event count and complete serialized byte count | planned TestSelectedSpecOmissions in internal/spec | A clipped count cannot describe the omitted evidence |
| QU7 | 6 | Each selected spec result includes an exact complete-history command | planned TestSelectedSpecDetailRoute in internal/spec | A generic help action cannot recover the omitted target |
| QU8 | 7 | A failed spec target retains its own result beside successful targets | planned TestSelectedSpecPartialFailure in internal/spec | One Git failure cannot erase another target result |
| QU9 | 8 | Existing bare worktree and positional history views match the baseline before new budgets | planned TestSelectedViewsPreserveDefaults in the respective command packages | A differential input matrix catches accidental default changes |
| QU10 | 9 | Malformed selection grammar returns usage at exit 2 | planned TestSelectedQueryGrammar in the respective command packages | Mixed modes and invalid limits cannot silently choose a different query |
| QU11 | 10 | Raw-read guidance selects relevant sections, paths, failures, or rows before full detail | review-owned: Standards axis reads the three guidance files | Advice to concatenate complete results defeats projection |
| QU12 | 11 | Raw-read guidance retains an explicit complete-detail route | review-owned: Standards axis reads the examples | A bounded example without recovery hides required evidence |
| QU13 | 12 | Guidance keeps polling and distinct authority boundaries as separate operations | review-owned: Spec axis compares ticket 4 and ticket 6 | A call-saving example cannot combine mutation with later approval |
| QU14 | 13 | New numeric defaults require an approved per-surface byte policy before their build ticket starts | review-owned: budget record and ticket-entry inspection | A build cannot turn the diagnostic cut into its own cap |
| QU15 | 14 | Approved defaults retain their complete-detail routes and owner metadata | planned TestApprovedQueryBudget in the respective command packages | Boundary fixtures catch lost detail routes or missing omission metadata |
| QU16 | 15 | Control-bearing target failures remain representable without hiding other results | planned TestSelectedQueryHostileTarget in the respective command packages | An unsafe TOON cell cannot collapse the complete result into one render error |
| QU17 | 16 | Selected query results omit unrequested worktree rows and history bodies | planned TestSelectedQueriesExcludeOldOutput in the respective command packages | Presence-only assertions would let the old full output survive |

### Edge inventory

QU1–QU3 cover active, complete, cleanup-pending, absent, and ambiguous worktree selections.
QU4–QU8 cover empty histories, duplicate slugs, retire-only, delete-only, mixed history, and per-target Git failures.
QU5 covers a limit of one, exact-boundary output, and omitted output.
QU10 covers zero, negative, nonnumeric, and overflowing limits, missing values, and mixed positional modes.

QU16 covers spaces, control bytes, leading dashes, and command-shaped operands for kit and linked-repository callers.
QU9 preserves existing no-repository and default-view behavior through a differential matrix.
No query acquires mutation authority or substitutes a package variable across subprocesses.

Won't handle: arbitrary command batches — the two existing query owners remain the callers.
Won't handle: a general Git query layer — the existing history producer and diff owner retain their responsibilities.
Won't handle: a general guidance rewrite — the named read examples remain the only guidance edits.

## Ownership fences

- `.agents/commands/bench-debug.md`
- `.agents/commands/bench-drain.md`
- `.agents/skills/bench-craft-cli/SKILL.md`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `cmd/bench/main.go`
- `cmd/bench/worktree_leaves.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/spec/history.go`
- `internal/spec/history_command_test.go`
- `internal/spec/history_test.go`
- `internal/spec/spec.go`
- `internal/usage/worktree.go`
- `internal/worktree/list.go`
- `internal/worktree/list_actions_test.go`
- `internal/worktree/list_selected_test.go`
- `internal/worktree/path.go`
- `reviews/session-context-queries.md`

Reviewer disposition: pending sign-off.
The fence is the union of ticket writes and the review pickup.
A build cannot change this spec, its acceptance rows, or its tickets.

## Ticket graph

| Ticket | Blocked by | Delivered coverage |
| --- | --- | --- |
| [1. Select worktree path facts](tickets/1-select-worktrees.md) | none | QU1, QU2, QU3, QU9, QU10, QU16, QU17 |
| [2. Select bounded spec histories](tickets/2-select-histories.md) | none | QU4, QU5, QU6, QU7, QU8, QU9, QU10, QU16, QU17 |
| [3. Guide relevant raw reads](tickets/3-guide-relevant-reads.md) | 1-select-worktrees.md, 2-select-histories.md | QU11, QU12, QU13 |
| [4. Apply reviewed owner budgets](tickets/4-apply-reviewed-budgets.md) | 1-select-worktrees.md, 2-select-histories.md, 3-guide-relevant-reads.md | QU14, QU15 |

## Out of scope

A general command batch language is a separate capability: approximately 8 edits, 2 gate runs.
A general guidance rewrite belongs to FT100: approximately 10 edits, 2 gate runs.
Build-then-run consolidation needs separate evidence: approximately 4 edits, 1 gate run.

## Further notes

### Source trace

| Source clause | Coverage |
| --- | --- |
| Ticket 3: bounded raw reads and Bench queries | QU5–QU7, QU11, QU12, QU14, QU15 |
| Ticket 4: selected multi-target queries and archive guidance | QU1–QU4, QU8, QU11, QU13 |
| Ticket 8: identity, path, state, and bounded histories | QU1–QU8, QU16, QU17 |
| Ticket 9: reviewer-approved numeric budgets | QU14, QU15 |
| FT173: existing owners and full-detail routes | QU4, QU7, QU9, QU12, QU15 |

### Reader sweep and proof checklist

[Seam evidence](../session-context-efficiency/assets/seam-evidence.md) records producer and enforcement reads.
Cited symbols: `selectAssignment`, `History`, `historyCommand`, and `ListCommand` resolve to the existing owners.
Import edges: no new cross-package producer is required.
Source-row clauses and occurrences: the sole compiled map owns the clauses above.

Promised field labels: the two selected schemas appear in Implementation decisions.
Changed-function callers: worktree leaf dispatch calls `ListCommand`; spec dispatch calls `historyCommand`.
Their local tests, public help tests, and command registry tests consume their grammar and output.
Copy survival: QU17 fails if the old full output survives beside the new projection.

The bare worktree fixture family remains unchanged under QU9.
The debug and drain history examples receive QU11–QU13 and exact ownership fences.
The repository-wide sweep found no JavaScript or workflow schema reader.
The AXI approved-query set remains unchanged.
The selected history view retains the spec owner's operational contract.

### Flagged additions

Repeated selectors, first-occurrence deduplication, and the explicit history limit are proposed grammar decisions for this sign-off.
QU1–QU10 and QU16–QU17 grade these additions.
The numeric-default ticket remains staged until its external checkpoint is complete.
The checkpoint retains the full bounded-default scope rather than silently declaring the early selectors complete.
