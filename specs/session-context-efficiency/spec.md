# Session context efficiency coordination

Status: staged

Decision source: `specs/session-context-efficiency/decisions/session-context-efficiency.md` (ready compiled map).

Verification log: 0 iteration(s) to accept — pending the requested Sol/high review.

## Problem

A small number of large tool results dominate session context.
Repeated domain operations also require repeated calls.
Uncoordinated changes can duplicate measurement owners or hide required evidence.

## Solution

One coordinating spec links four child specs.
Each child owns its implementation stories, acceptance rows, and ticket graph.
This document owns their dependencies and reviewer checkpoints.
It adds no executable orchestration layer and has no implementation tickets.
The child ticket folders provide the complete build breakdown.

| Child spec | Delivered outcome | Prerequisite |
| --- | --- | --- |
| [Measurement](../session-context-measurement/spec.md) | Reliable measures and comparative budget evidence | Spec and ticket approval |
| [Queries](../session-context-queries/spec.md) | Task-shaped reads and both selected query families | Measurement and budget review before new defaults |
| [Overflow](../session-context-overflow/spec.md) | Verified complete-output preservation and bounded replacement | Measurement, budget review, and runtime capability evidence |
| [Cleanup](../session-context-cleanup/spec.md) | One selected-set cleanup plan and apply | Spec and ticket approval |

## User stories

### Coordinate the delivered capabilities

Line: gpt-5.6-terra / high.
The work coordinates uncertain semantics across existing owners.

1. As a reviewer, I want one source for each child requirement, so that the program cannot drift between copies.
2. As an agent, I want measurement work to remain independently useful, so that unavailable harness data does not block known facts.
3. As an agent, I want cleanup to proceed independently, so that output-budget decisions do not delay lifecycle improvements.
4. As a reviewer, I want measured budget proposals, so that new defaults need my explicit approval.
5. As a reviewer, I want runtime evidence before overflow enablement, so that a documented hook does not become an unsupported promise.
6. As an agent, I want the complete program scope retained, so that a useful first child does not silently close the rest.
7. As a reviewer, I want separate child approvals, so that this coordinating record does not authorize an unchecked build.
8. As a reviewer, I want one compiled decision source, so that subsequent sessions use the same closed decisions.

## Implementation decisions

This is a coordination document, not a fifth product capability.
The four sibling spec folders are independently reviewable build sources.
Their ticket graphs express local dependencies with sibling ticket basenames.
This document expresses dependencies between specs.
A child build reads both its local prerequisites and this coordination table.

The first build frontier contains measurement and cleanup after their approvals.
Query selectors can precede new default budgets only where their child spec permits that slice.
No build fills numeric budget placeholders or enables an unverified harness path.
Those changes return to spec authoring with the required evidence.
The reviewer owns the resulting budget and enablement decisions.

The measurement owner supplies facts to the other children.
The query and lifecycle owners retain command policy.
The overflow owner supplies capability evidence without claiming provider-token attribution.
The approved map remains the sole record of reviewer decisions.

## Testing decisions

The coordinating record has review-owned acceptance.
Each child names executable tests at its own seam.
The spec checker validates all five coverage maps and the child ticket grammar.
The project gate retains its existing six-phase architecture.
This phase does not change the gate or manufacture executable coverage for a document-only checkpoint.

### Seam diagram

```text
reviewed map -> coordinating record -> four child specs -> child ticket graphs
                        |
                        +-> reviewer budget and capability checkpoints
```

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| SC1 | 1, 6 | The coordinator links each approved child spec without copying its acceptance rows | review-owned: Spec axis compares ticket 7 with the child table | Omitting one child leaves an approved outcome without an owner |
| SC2 | 2 | Measurement has an independent build entrypoint | review-owned: ticket graph inspection | A dependency on query or overflow implementation creates a cycle |
| SC3 | 3 | Cleanup has no dependency on a new output budget | review-owned: prerequisite inspection | A measurement blocker on cleanup contradicts the approved split |
| SC4 | 4 | New numeric defaults wait for reviewer approval of measured proposals | review-owned: child readiness and budget record inspection | A numeric default chosen during a build bypasses the retained decision |
| SC5 | 5 | Overflow enablement waits for evidence from the installed harness path | review-owned: capability acceptance in the overflow child | A documentation-only capability claim cannot meet the checkpoint |
| SC6 | 7 | Each child requires its own spec and ticket approval | review-owned: approval record inspection | Treating this record as blanket build approval fails the predicate |
| SC7 | 8 | The map and topic folder exist only in the compiled location | review-owned: repository path census | A surviving top-level copy leaves two decision sources |

### Edge inventory

The audience includes the kit repository and repositories that link Bench.
A child with unknown harness measures preserves explicit absence.
A child waiting for budget evidence remains staged and cannot select its own numeric defaults.
A completed child does not retire an incomplete sibling.
These dispositions attach to SC2, SC4, SC5, and SC1 respectively.

Won't handle: automatic cross-spec execution — the reviewer invokes the existing build phase for each approved child.
Won't handle: a new orchestration service — the existing phase and worktree owners remain the callers.

## Ownership fences

- `reviews/session-context-efficiency.md`

The reviewer disposition is pending sign-off.
This coordination source authorizes no product-code writes.
Each child declares its own complete fence.

## Out of scope

There is no fifth implementation capability: 0 product edits, 0 separate product gate runs.
The children preserve the compiled map's reviewed exclusions.
No provider-dollar score, general guidance rewrite, or new efficiency denial joins this program.

## Further notes

### Reader sweep and source trace

The compiled map's ticket 7 supplies SC1, SC2, SC3, and SC6.
Ticket 9 supplies SC4; tickets 3 and 6 supply SC5.
The phase's map-move contract supplies SC7.
The existing decision-map integrity, coverage, and ticket-grammar readers consume the staged artifacts.
Their enforcement references are in [seam evidence](assets/seam-evidence.md).

### Pre-review proof checklist

- Cited symbols: none in acceptance claims.
- Import edges: none.
- Source-row clauses and occurrences: the trace above names the sole compiled decision source.
- Promised field labels: none in product output.
- Changed-function callers: none.
- Copy survival: SC7 checks the moved source.

### Flagged additions

The coordination format introduces no behavior beyond the approved coordinating spec and four-child split.
The review-owned dependency record does not add a new CLI or gate check.
