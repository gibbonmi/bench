# Plan and apply one complete cleanup target set

Status: staged

Decision source: `specs/session-context-efficiency/decisions/session-context-efficiency.md` (ready compiled map).

Verification log: 2 iteration(s) to accept — Sol/high accepted the suite. Trace-only partials and review-state bookkeeping are folded.

Coordinator: [Session context efficiency](../session-context-efficiency/spec.md)

## Problem

Agents repeat cleanup plans for several targets, then invalidate later plans through their own earlier removals.
Existing set selectors do not yet report every unstarted outcome after a partial apply.

## Solution

Extend the lifecycle owner with one explicit target-set operation alongside its existing selectors.
Qualify the complete selection before removal and retain each target’s existing locked lifecycle transaction.
Report every selected outcome and an exact re-plan action.
This capability requires no new numeric output budget.

## User stories

Line: gpt-5.6-sol / high.
Implementation-line reason: Complete-set preflight is the hardest chunk. The source fixes the lifecycle outcome, but concurrent drift needs tests at existing transaction boundaries.
Harder chunks: CL-C2.
The work crosses uncertain runtime or lifecycle boundaries.

1. As an agent, I want one fingerprint for my selected cleanup set, so that one apply represents the whole intent.
2. As an agent, I want cleanup aliases collapsed, so that one target cannot be removed twice.
3. As a reviewer, I want selection failures to prevent removal, so that an unresolved target cannot shrink my approved set.
4. As a reviewer, I want stale plans refused before removal, so that known drift cannot authorize partial cleanup.
5. As a reviewer, I want complete-set preflight, so that a later stale row cannot be discovered after an avoidable deletion.
6. As a reviewer, I want each locked recheck retained, so that preflight cannot authorize a later race.
7. As an agent, I want completed outcomes after failure, so that I know which effects already occurred.
8. As an agent, I want unstarted targets reported, so that a partial result accounts for the whole selection.
9. As an agent, I want an exact re-plan action, so that recovery keeps the intended scope and modifiers.
10. As a reviewer, I want current cleanup authority retained, so that a set operation cannot remove active work.
11. As an existing caller, I want current lifecycle effects preserved, so that set support cannot alter single-target semantics.
12. As a reviewer, I want spent plans rejected, so that retries cannot repeat completed cleanup.
13. As an agent, I want an honest empty plan, so that absence cannot look like a failed mutation.
14. As an agent, I want conflicting selection modes rejected, so that the command cannot guess my scope.
15. As an agent, I want cleanup independent of budget work, so that a separate measurement decision cannot block lifecycle improvements.
16. As an agent, I want hostile operands kept as data, so that selection and recovery cannot execute injected commands.

## Implementation decisions

Proposed grammar adds repeated `--target <target>` operands to `bench worktree clean`.
The existing positional path, `--landed`, and `--unclaimed` forms remain available.
Exactly one selection mode is valid per call.
Existing discard modifiers retain their current eligibility and authority rules.
The new mode does not imply either discard modifier.

Explicit targets use the existing assignment identity resolver.
Aliases collapse by assignment identity, and canonical identity order makes the fingerprint deterministic.
Any unresolved or ambiguous explicit target prevents an applicable fingerprint and all removal.
The plan still reports each selection outcome so the agent can correct the request.
An unsafe but resolved target retains its existing refused or retained plan disposition.
No new path resolver or deletion engine joins the command.

The fingerprint covers complete membership and the existing removal-relevant state tuple for each selected target.
Apply recomputes the complete set before the first transaction.
It also requalifies every removable row before the first removal.
A mismatched membership or tuple refuses without removal.
The same check remains inside each target's existing locked transaction.
Preflight does not promise an atomic filesystem snapshot across later concurrent mutations.

Each target transaction keeps the existing receipt, recovery, branch, assignment, census, and handoff behavior.
If a race or failure occurs after an earlier removal, the earlier result remains completed.
The failing target reports its actual retained or failed outcome.
Each remaining selected target reports `not-attempted`.
The command returns nonzero and does not claim rollback of completed effects.

A stale result preserves the selected mode and modifiers in its exact re-plan command.
Explicit targets are safely quoted in that command.
The command never substitutes a fresh fingerprint into the failed apply automatically.
Existing `--apply-current` behavior remains limited to its existing unclaimed-branch mode.
No new efficiency denial or numeric result cap is introduced.
The landing effect retains its existing scoped plan-and-apply call without a user-facing fingerprint round trip.

The existing late-drift test remains valid.
It changes the second target after the first target's terminal receipt.
That fixture proves the required per-target race posture, not stale state present before set preflight.
A separate fixture introduces later-target drift before apply entry and requires zero removals.

## Implementation chunks

One retained implementation session owns this child after approval.
Each existing ticket forms one named review chunk and one serial commit checkpoint.
The table orders independent tickets that share command inventory writes.
After each chunk, freeze its predecessor and current tips for Standards, Spec, and Coverage review.
The successor starts after accepted findings have current repair coverage.

| chunk / ticket | blocked by | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- | --- |
| CL-C1 / `1-plan-explicit-sets.md` | none | Explicit cleanup sets | CL1, CL2, CL3, CL10, CL11, CL13, CL14, CL15, CL16, CL17 | TestCleanExplicitSetPlan and the remaining owned-row tests | no |
| CL-C2 / `2-complete-preflight-outcomes.md` | 1-plan-explicit-sets.md | Complete preflight and outcomes | CL4, CL5, CL6, CL7, CL8, CL9, CL12, CL18, CL19 | TestCleanSetPreflightAllRows and the existing landing-effect tests | yes |

The coverage map supplies the complete test inventory for each chunk's owned rows.
The final reconciliation checks every acceptance row and the integrated result.
The existing evidence checkpoints continue to block their implementation chunks.
Execution-plan changes follow `.bench/BENCH.md`.

## Testing decisions

Tests drive the production owner through controlled inputs and its existing injected boundaries.
Ordinary package tests execute through the existing Go test phase.
Hook integration tests execute through `bench test --check system` where the row requires a real subprocess.
Live harness evidence remains review-owned and cannot be replaced by a simulated callback.
Planned test names below identify required tests, not tests that already exist.

### Seam diagram

```text
caller -> domain entrypoint -> verification -> existing operation -> complete result
                ^                  ^                  ^
          command tests       stale/fault cases   durable evidence
```

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| CL1 | 1 | One explicit target set produces one fingerprint for the complete resolved selection | planned TestCleanExplicitSetPlan in internal/worktree | Independent per-target plans cannot identify one selected intent |
| CL2 | 2 | Aliases of one cleanup identity appear once in deterministic order | planned TestCleanExplicitSetAliases in internal/worktree | Duplicate identities cannot trigger repeated removal |
| CL3 | 3 | A target-resolution failure prevents every removal | planned TestCleanExplicitSetSelectionFailure in internal/worktree | Partial selection must not silently authorize a smaller mutation set |
| CL4 | 4 | A stale selected-set fingerprint refuses before the first removal | planned TestCleanSetPreexistingDrift in internal/worktree | A later target already dirty at apply entry cannot permit an earlier deletion |
| CL5 | 5 | Every removable row is requalified before the first transaction begins | planned TestCleanSetPreflightAllRows in internal/worktree | A row-by-row-only preflight lets an early target disappear before known later drift |
| CL6 | 6 | Each target retains the existing under-lock lifecycle recheck | planned TestCleanSetLateDrift in internal/worktree | A race after preflight cannot use an obsolete removal plan |
| CL7 | 7 | A partial apply reports completed outcomes without claiming rollback | planned TestCleanSetPartialApply in internal/worktree | A failed later transaction cannot erase the earlier completed result |
| CL8 | 8 | A partial apply reports every unstarted target as not attempted | planned TestCleanSetUnstartedOutcomes in internal/worktree | Omitting remaining rows hides part of the selected intent |
| CL9 | 9 | A stale result names the exact selector-preserving re-plan command | planned TestCleanSetStaleReplanAction in internal/worktree | A generic clean command loses the selection or discard modifiers |
| CL10 | 10 | Active or unsafe targets retain the existing cleanup refusal | planned TestCleanSetRetainsAuthority in internal/worktree | Set selection cannot widen deletion authority |
| CL11 | 11 | Existing single-target and selector success cases match the baseline lifecycle effects | planned TestCleanSetCompatibility in internal/worktree | A differential state comparison catches changed branch or receipt behavior |
| CL12 | 12 | A completed target cannot be removed twice through a repeated apply | planned TestCleanSetSpentPlan in internal/worktree | Reusing a spent plan cannot replay cleanup side effects |
| CL13 | 13 | A present empty landed-assignment inventory reports an empty plan without mutation | planned TestCleanSetPresentEmptyInventory in internal/worktree | A present empty inventory cannot produce fabricated removable rows |
| CL14 | 14 | Invalid mixed selection grammar returns usage before plan creation | planned TestCleanSetGrammar in internal/worktree | A selector plus explicit targets cannot silently widen the selected set |
| CL15 | 15 | Cleanup starts without a new numeric output-budget dependency | review-owned: ticket graph and entry checks | A measurement prerequisite on cleanup contradicts the approved independent capability |
| CL16 | 16 | Hostile operand text remains data throughout selection and re-plan output | planned TestCleanSetHostileOperand in internal/worktree | A command-shaped path cannot execute through the recovery action renderer |
| CL17 | 13 | An absent landed-assignment inventory reports an empty plan without mutation | planned TestCleanSetAbsentInventory in internal/worktree | Treating a missing assignment store as a command failure violates the existing empty state |
| CL18 | 11 | Landing retains automatic cleanup of the sibling carried by that landing | `internal/worktree/land_effects_cleanup_test.go` (`TestLandCleansTheFoldedSibling`) | A command-only fingerprint requirement would prevent the existing landing effect |
| CL19 | 11 | Landing retains a previously landed assignment outside its cleanup scope | `internal/worktree/land_effects_cleanup_test.go` (`TestLandLeavesAPriorLandedAssignment`) | Discarding the scope argument would turn one landing into repository-wide cleanup |

### Edge inventory

CL3 and CL14 cover unresolved, ambiguous, missing-value, and mixed-mode selection.
CL4 and CL5 cover membership drift, changed head, dirty tracked state, ignored-file changes, and changed leases before removal.
CL6–CL8 cover under-lock refusal, transaction faults, terminal receipt faults, and remaining unstarted targets.
CL10 and CL11 retain the existing lifecycle owner's unsafe, active, dirty, ignored, and branch-preservation cases.

CL12 covers repeated apply after a completed target.
CL13 and CL17 distinguish present-empty and absent landed-assignment inventories for kit and linked-repository callers.
CL16 covers spaces, newlines, leading dashes, and shell-shaped target text.
Tests use the existing cleanup boundary in-process; no package-variable swap is expected to affect a separate executable.

Won't handle: rollback of completed removals — the existing lifecycle transaction reports durable per-target effects.
Won't handle: cross-repository target sets — each existing cleanup invocation retains one repository boundary.
Won't handle: automatic stale-plan approval — the agent runs the rendered re-plan operation explicitly.

## Ownership fences

- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/usage/worktree.go`
- `internal/worktree/clean_landed.go`
- `internal/worktree/clean_landed_apply_test.go`
- `internal/worktree/clean_set.go`
- `internal/worktree/clean_set_apply_test.go`
- `internal/worktree/clean_set_command_test.go`
- `internal/worktree/clean_set_test.go`
- `internal/worktree/clean_unclaimed.go`
- `internal/worktree/clean_unclaimed_test.go`
- `internal/worktree/land_effects.go`
- `internal/worktree/land_effects_cleanup_test.go`
- `internal/worktree/worktree.go`
- `reviews/session-context-cleanup.md`

Reviewer disposition: Sol/high review accepted; user spec and ticket sign-off remains pending.
The fence is the union of ticket writes and the review pickup.
`.bench/BENCH.md` governs execution-plan changes.

## Ticket graph

| Ticket | Blocked by | Delivered coverage |
| --- | --- | --- |
| [1. Plan and apply explicit cleanup sets](tickets/1-plan-explicit-sets.md) | none | CL1, CL2, CL3, CL10, CL11, CL13, CL14, CL15, CL16, CL17 |
| [2. Preflight and report the complete set](tickets/2-complete-preflight-outcomes.md) | 1-plan-explicit-sets.md | CL4, CL5, CL6, CL7, CL8, CL9, CL12, CL18, CL19 |

## Out of scope

Cross-repository cleanup is a separate capability: approximately 7 edits, 2 gate runs.
Rollback across completed removals is a separate lifecycle capability: approximately 9 edits, 2 gate runs.

## Further notes

### Source trace

| Source clause | Coverage |
| --- | --- |
| Ticket 4: one selected-set plan, fingerprint, apply, and per-target outcomes | CL1–CL9, CL12, CL13, CL17 |
| Ticket 4: stale checks before removal and a re-plan operation | CL4–CL6, CL9 |
| Ticket 6: existing safety guards remain | CL10, CL11, CL14, CL16, CL18, CL19 |
| Ticket 7: independently useful cleanup | CL15 |

### Reader sweep and proof checklist

[Seam evidence](../session-context-efficiency/assets/seam-evidence.md) records the lifecycle and enforcement reads.
Cited symbols: `CleanCommand`, `applyLandedSet`, `applyCleanupTransaction`, and `executeCleanup` remain the command and lifecycle owners.
Import edges: no new cross-package lifecycle owner is required.
Source-row clauses and occurrences: the sole compiled map owns the clauses above.

Promised field labels: `not-attempted` is an apply outcome, never a planned removal claim.
Changed-function callers: the worktree leaf calls `CleanCommand`; landed apply calls the existing cleanup transaction.

The landing effect also calls `planLandedSet` and `applyLandedSet` with its destination-base scope.

CL18 and CL19 protect that existing caller.
Local cleanup tests and public grammar tests consume the changed behavior.
Copy survival: no lifecycle copy is authorized.

The existing post-first-removal drift fixture remains an in-scope success under CL6 and CL7.
The new pre-existing-drift fixture belongs to CL4 and CL5.
The unclaimed selector retains its current whole-set replan and branch transaction authority.
The landing effect retains its automatic scoped cleanup through the existing caller.
The final ticket carries the cleanup package's complete-set invariant.

### Flagged additions

Repeated explicit targets are a proposed scope clarification for reviewer sign-off.
The earlier optional clarification has no recorded reviewer answer.
CL1–CL3, CL14, and CL16 grade the proposed selection grammar.
The existing selector improvements remain independently valuable if the reviewer declines explicit targets.
