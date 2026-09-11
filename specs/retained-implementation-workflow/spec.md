# Retain implementation authorship and review each chunk

Status: staged
Decision source: `specs/retained-implementation-workflow/decisions/benchmark-orchestration.md`
Verification log: 2 iteration(s) to accept — Sol/high approved the repaired spec and ticket graph. Source predicates, tier cases, and reader closure were corrected.

## Problem

The current kit assigns a fresh writer to each ticket and waits until the whole spec is built for its routine implementation review. This loses implementation context and delays review feedback. The reviewer has selected one retained author with independent review after each coherent chunk.

## Solution

Make retained authorship the portable default. The spec author defines review chunks and recommends an implementation line. The implementation session writes, verifies, repairs, and commits each ticket. Independent reviewers examine each chunk before work starts on its successor.

Build dependency: none.

## User stories

Line: gpt-5.6-sol / high. One retained implementation session. Starting effort is a recommendation for approval.

Sol/high is recommended because this build changes guidance and its conformance fixtures at known owners. Chunk 1 is harder because it reconciles many readers. Gate and landing implementation belongs to the separate completion-evidence spec.

1. As an operator, I want to retain one author across tickets, so that I can assess and continue approved work accurately.
2. As an operator, I want to define coherent chunks before implementation, so that I can assess and continue approved work accurately.
3. As an operator, I want to keep ticket verification checkpoints, so that I can assess and continue approved work accurately.
4. As an operator, I want to review work before the next chunk, so that I can assess and continue approved work accurately.
5. As an operator, I want to keep the author responsible for repairs, so that I can assess and continue approved work accurately.
6. As an operator, I want to reconcile the entire spec at completion, so that I can assess and continue approved work accurately.
7. As an operator, I want to avoid redundant full reviews, so that I can assess and continue approved work accurately.
8. As an operator, I want to revise chunks without approval delays, so that I can assess and continue approved work accurately.
9. As an operator, I want to expand necessary write authority, so that I can assess and continue approved work accurately.
10. As an operator, I want to expand gate coverage safely, so that I can assess and continue approved work accurately.
11. As an operator, I want to retain a record of autonomous changes, so that I can assess and continue approved work accurately.
12. As an operator, I want to receive a reasoned implementation recommendation, so that I can assess and continue approved work accurately.
13. As an operator, I want to use the agreed GPT tiers, so that I can assess and continue approved work accurately.
14. As an operator, I want to use consistent independent reviewers, so that I can assess and continue approved work accurately.
15. As an operator, I want to keep explicit user control, so that I can assess and continue approved work accurately.
16. As an operator, I want to remove conflicting instructions, so that I can assess and continue approved work accurately.

17. As an operator, I want to map chunk acceptance, so that the required evidence remains complete.
18. As an operator, I want to plan chunk tests, so that the required evidence remains complete.
19. As an operator, I want to keep chunks vertical, so that the required evidence remains complete.
20. As an operator, I want to retain diagnostic access, so that the required evidence remains complete.
21. As an operator, I want to retain implementation ownership, so that the required evidence remains complete.
22. As an operator, I want to bind the top tier, so that the required evidence remains complete.
23. As an operator, I want to bind the low tier, so that the required evidence remains complete.

## Implementation decisions

The canonical authorship and authority policy lives in `.bench/BENCH.md`. Craft skills and phase commands apply it by reference. `craft-line` owns implementation-line judgment and the single review-line policy. `.bench/lines.env` owns executable bindings, with the project table checked against it. These owners serve every linked repository.

An implementation chunk is a coherent behavior outcome with acceptance rows, tests, and a review checkpoint. A ticket remains a green commit checkpoint. A chunk can contain several tickets. The initial plan uses one ticket per chunk unless splitting the commits makes the same outcome safer to deliver. Chunk IDs survive retitles. A changed plan records old-to-new IDs and preserves all acceptance coverage and dependencies.

The retained author writes production changes, tests, probes, and repairs. Review subagents do not become repair writers. A brief diagnostic consultation follows the continuation policy when that capability lands. The user can explicitly choose another authorship arrangement. The default alone does not grant a model switch.

Each chunk review uses Standards, Spec, and Coverage against the whole approved spec, focused on the frozen chunk delta. Resolve findings and obtain current repair coverage before starting the next chunk. After the last chunk, the author reconciles overall acceptance and integration. Repeat delegated review only for a later delta or a cross-chunk concern that invalidates prior evidence. Retain the current native dispatch and incomplete-axis rules.

The spec author recommends one implementation tier and starting effort. The reason considers the hardest material chunk, spec precision, uncertainty, and tests. Mark harder chunks in the chunk table. The author can adjust effort and report it.

A different implementation model or session requires user direction. A recommendation for a review exception names its reason. The canonical default is one fresh mid/high subagent per review axis per chunk. A Sol implementation instead uses Astra/high for each axis.

Bind GPT top to `gpt-6-astra`, mid to `gpt-5.6-sol`, and cheap to `gpt-5.6-terra`. Cheap is the existing configuration key for the user's low tier. Preserve other harness columns and their resolution behavior.

The approved spec owns overall write authority. Ticket Writes entries predict the touched files. The author can split, combine, or reorder chunks and expand writes or gate coverage within approved behavior. Update the affected spec and tickets before using the expansion.

Preserve existing checks, pass criteria, and required behavior. Record each change with `bench learning`, including what changed, why, and verification. The drain reviews these records later. Unrelated scope and weakened guarantees still require a user decision.

Deliver guidance with honest enforcement limits. Anchors detect missing or contradictory instructions. They do not prove that an agent followed them. The completion-evidence spec adds the gate-owned checkpoint check. Until it lands, this phase follows the review checkpoint in its instructions and must not claim mechanical enforcement.

## Implementation chunks

Implement these tickets in the retained session. Each ticket is one initial review chunk. Commit the ticket on its lane pass, run the three delegated axes, and repair findings before the next chunk. Review line: decision #13, resolved for this implementation as gpt-6-astra / high. Final acceptance reconciliation stays with the author.

| chunk / ticket | blocked by | delivered outcome | harder chunk |
| --- | --- | --- | --- |
| 3.md — Bind the model tiers and review recommendation | none | Update GPT bindings and the checked project table | no |
| 1.md — Retain the author through chunk review | 3.md | Apply the authorship and review cadence end to end, including the drain and field guide | yes |
| 2.md — Permit logged within-scope plan expansion | 1.md | Apply standing authority to chunk amendments, Writes expectations, and gate additions | no |

## Testing decisions

Use the named existing owner and new tests below. A new test name is a planned seam, not a claim that the test exists. Read its nearest fixture before implementation. Demonstrate each required omission or behavioral mutation as a diagnostic red, restore it, and show green. A compile failure is not the required red. Semantic review judges prose quality and the sufficiency of the test.

### Seam diagram

```text
approved source -> canonical workflow and line policy -> phase and skill readers
                                      |
                         anchors + negative fixtures -> docs gate
```

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| W1 | 1 | Authorship remains in one implementation session | New test: TestRetainedWorkflow in internal/conformance/retained_workflow_test.go | Replace the retained-author instruction with fresh writers and observe its named diagnostic. |
| W2 | 2 | Each planned chunk declares a behavior outcome | New test: TestRetainedWorkflow in internal/conformance/retained_workflow_test.go | Omit the outcome from the chunk contract and observe red. |
| W3 | 3 | Each ticket keeps a serial green commit checkpoint | New test: TestRetainedWorkflow in internal/conformance/retained_workflow_test.go | Remove the ticket checkpoint while retaining chunk review and observe red. |
| W4 | 4 | The phase orders three-axis chunk review before advancement | New test: TestRetainedWorkflow in internal/conformance/retained_workflow_test.go | Move the review after advancement in the phase fixture and observe red. |
| W5 | 5 | Review findings return to the retained author | New test: TestRetainedWorkflow in internal/conformance/retained_workflow_test.go | Restore writer delegation for repairs and observe red. |
| W6 | 6 | Final acceptance reconciliation covers integration | New test: TestRetainedWorkflow in internal/conformance/retained_workflow_test.go | Delete the final reconciliation instruction and observe red. |
| W7 | 7 | Routine final review is replaced by evidence-triggered review | New test: TestRetainedWorkflow in internal/conformance/retained_workflow_test.go | Restore unconditional final full review in any live owner and observe red. |
| W8 | 8 | Chunk amendments preserve acceptance and review checkpoints | New test: TestRetainedWorkflow in internal/conformance/retained_workflow_test.go | Omit preservation from the amendment instruction and observe red. |
| W9 | 9 | Within-scope write expansion requires an updated plan before use | New test: TestRetainedWorkflow in internal/conformance/retained_workflow_test.go | Restore the chunk fence as an approval boundary and observe red. |
| W10 | 10 | Gate additions preserve existing pass criteria | New test: TestRetainedWorkflow in internal/conformance/retained_workflow_test.go | Permit dropping an existing check in the expansion fixture and observe red. |
| W11 | 11 | Plan and gate expansions enter the learning drain | New test: TestRetainedWorkflow in internal/conformance/retained_workflow_test.go | Remove the learning step from the canonical policy and observe red. |
| W12 | 12 | The spec template requires one implementation line with a reason | New test: TestRetainedWorkflow in internal/conformance/retained_workflow_test.go | Omit the hardest-chunk factor and observe red. |
| W13 | 13 | The GPT mid binding is gpt-5.6-sol | New test: TestRetainedWorkflow in internal/conformance/retained_workflow_test.go | Replace only the mid binding and observe red. |
| W14 | 14 | Each axis resolves the canonical conditional review line | New test: TestRetainedWorkflow in internal/conformance/retained_workflow_test.go | Test Sol implementation with Astra/high review and other implementation lines with mid/high review. Mutate each branch independently. |
| W15 | 15 | A model switch requires user direction | New test: TestRetainedWorkflow in internal/conformance/retained_workflow_test.go | Remove the switch boundary while retaining effort changes and observe red. |
| W16 | 16 | No live workflow reader retains mandatory fresh implementation writers | New test: TestRetainedWorkflow in internal/conformance/retained_workflow_test.go | Leave an old writer instruction in the drain or field guide and observe red. |

| W17 | 17 | Each planned chunk names its acceptance rows | New test: TestRetainedWorkflow in internal/conformance/retained_workflow_test.go | Omit row references while retaining the outcome. |
| W18 | 18 | Each planned chunk names its tests | New test: TestRetainedWorkflow in internal/conformance/retained_workflow_test.go | Omit tests while retaining outcomes and acceptance. |
| W19 | 19 | A chunk delivers coherent behavior with its tests | review-owned: Spec reviews the chunk graph | Reject a file-only function-only or tests-only split. |
| W20 | 20 | The retained-author policy permits brief diagnostic consultation | New test: TestRetainedWorkflow in internal/conformance/retained_workflow_test.go | Remove the diagnostic exception from the author policy. |
| W21 | 21 | A diagnostic helper receives no implementation or repair assignment | New test: TestRetainedWorkflow in internal/conformance/retained_workflow_test.go | Allow the helper to write a repair and observe red. |
| W22 | 22 | The GPT top binding is gpt-6-astra | New test: TestRetainedWorkflow in internal/conformance/retained_workflow_test.go | Replace only the top binding and observe red. |
| W23 | 23 | The GPT cheap binding is gpt-5.6-terra | New test: TestRetainedWorkflow in internal/conformance/retained_workflow_test.go | Replace only the cheap binding and observe red. |

### Edge inventory

Reader sweep: the operating guide, implement, write-spec, review, final-check, drain, craft-spec, craft-tickets, craft-delegate, its delegation reference, craft-line, and field guide consume these facts. The anchors registry and workflow conformance helpers enforce the prose. The native-review registry retains its three-axis protections. Binding readers use `internal/lines`; their generic tier fixtures need no model-literal rewrite. Historical audit reports remain historical.

The posture change reds the current writer anchors in `internal/anchors/registry_data.go`, including the ticket-stage, ticket-light-path, repair-writer, drain, and field-guide cases. Inspect the matching fixtures under `tests/canary/workflow-guidance-anchors` and replace their mutation and expected diagnostic together. The existing line-binding prose-drift canary remains a drift test. Exclude staged specs and their tickets from bulk rewriting.

The new workflow test belongs to the existing docs conformance entry. Keep the operating guide as policy owner. The anchors registry holds enforcement metadata that pins references and rejects conflicting instructions. Independent negative expectations must demonstrate the named omission red and restored green. Semantic quality of a chunk boundary remains review-owned. Anchor presence cannot establish that quality.

Before this spec is retired, promote its compiled map to a durable reviewed decision artifact. Update surviving spec references in that same change. Do not delete a decision source still used by another staged spec.

## Ownership fences

These paths are the union of ticket expectations. A directory entry is an exact prefix for that existing owner or fixture family. Expansion follows decision #5, with the plan updated before use. It cannot weaken existing guarantees.

- `.bench/BENCH.md`
- `.bench/BENCH-reference.md`
- `CONTEXT.md`
- `.agents/commands/bench-implement-spec.md`
- `.agents/commands/bench-write-spec.md`
- `.agents/commands/bench-review-implementation.md`
- `.agents/commands/bench-final-check.md`
- `.agents/commands/bench-drain.md`
- `.agents/skills/bench-craft-spec/SKILL.md`
- `.agents/skills/bench-craft-tickets/SKILL.md`
- `.agents/skills/bench-craft-delegate/SKILL.md`
- `.agents/skills/bench-craft-delegate/references/delegation-discipline.md`
- `.agents/skills/bench-craft-line/SKILL.md`
- `docs/field-guide.html`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_ft311_preparation.go`
- `internal/anchors/registry_ft311_review_dispatch.go`
- `internal/conformance/docs_workflow_helpers_test.go`
- `internal/conformance/docs_workflow_checks_test.go`
- `internal/conformance/ft311_preparation_test.go`
- `tests/canary/workflow-guidance-anchors`
- `internal/conformance/retained_workflow_test.go` (new)
- `.agents/skills/bench-craft-gate/SKILL.md`
- `.bench/lines.env`
- `projects/benchkit.md`
- `internal/conformance/line_routing_static_test.go`
- `tests/canary/line-routing/line-binding-prose-drift/files/dot-bench/lines.env`
- `specs/retained-implementation-workflow/spec.md`
- `specs/retained-implementation-workflow/tickets`
- `reviews/retained-implementation-workflow.md` (new)

- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `tests/canary/docs-currency-token-diet/benchref-imported`
- `tests/canary/docs-currency-token-diet/benchref-pointer-dropped`
- `tests/canary/docs-currency-token-diet/benchref-section-duplicated`
- `tests/canary/docs-currency-token-diet/dogfood-referent-shipped`
- `tests/canary/docs-currency-token-diet/introduces-undeclared-command`
- `tests/canary/docs-currency-token-diet/missing-cli-inventory`
- `tests/canary/docs-currency-token-diet/signal-vocabulary-drift`
- `tests/canary/docs-currency-token-diet/stale-cli-doc-reference`
- `tests/canary/docs-currency-token-diet/stale-command-reference`
- `tests/canary/guidance-prose-budgets/over-budget-skill`
- `tests/canary/line-routing/line-binding-prose-drift`
- `tests/canary/load-validity-metadata/readme-shared-rule-drift`
- `tests/canary/load-validity-metadata/shared-rule-drift`
- `tests/canary/row-next-grammar/token-table-lacks-kit-edit`
- `tests/canary/skills-index-command-adapters/adapter-inert-invocation-key`
- `tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy`
- `tests/canary/skills-index-command-adapters/dangling-index`
- `tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted`
- `tests/canary/skills-index-command-adapters/missing-index-field`
- `tests/canary/skills-index-command-adapters/stale-index-wording`
- `tests/canary/skills-index-command-adapters/unindexed-skill`

- `internal/conformance/recurrence_maintenance_contract_test.go`

## Out of scope

A new implementation scheduler, automatic model switching, changes to other harness tier values, and measurement-based default adoption are outside this spec. The completion-evidence and continuation specs own their stronger contracts.

## Further notes

Approval is pending for this spec and its ticket graph. Authoring these documents does not authorize their implementation. The source decisions remain closed. The current kit instructions have not yet been changed by these specs.

Flagged additions: the engineering mechanisms named under Implementation decisions, the new test seams, and the ticket graph. They implement the source outcomes. The review must reject an additional behavior with no coverage row.

Pre-review proof checklist:

- Cited symbols: existing owner names were read in this session. New test names are explicitly planned and must be created by the ticket that owns their row.
- Import edges: no new forbidden-import claim. New data owners list their actual consumers in the seam diagram.
- Source-row clauses and occurrences: the table below maps the closed source clauses. The index and resolved tickets are the authoritative occurrences. Historical recommendation assets do not override them.
- Promised field labels: the exact record fields, command arguments, and plan fields appear under Implementation decisions where this spec adds them. Guidance-only specs add no machine record.
- Changed-function callers: the reader and enforcement inventory appears under Edge inventory. Recheck function-level callers before changing an existing signature. No signature change is authorized by assumption.
- Copy survival: W16 owns live workflow-copy removal. Other specs add one owner and require consumers to call it rather than parse the same facts independently. Historical evidence remains readable.

| source clause | acceptance rows |
| --- | --- |
| #2: "Implementation stays in one model session" | W1, W3, W5, W16 |
| #12: "Each chunk receives delegated review before the next starts" | W2, W4, W6, W7 |
| #4 and #5: amend chunks and use standing expansion authority | W8, W9, W10, W11 |
| #13: recommend implementation line and resolve its review policy | W12, W13, W14, W15 |

Exact reader closure for W16: `.bench/BENCH.md`, `.bench/BENCH-reference.md`, `.agents/commands/bench-implement-spec.md`, `.agents/commands/bench-write-spec.md`, `.agents/commands/bench-review-implementation.md`, `.agents/commands/bench-final-check.md`, and `.agents/commands/bench-drain.md`.

Skill readers are `.agents/skills/bench-craft-spec/SKILL.md`, `.agents/skills/bench-craft-tickets/SKILL.md`, and `.agents/skills/bench-craft-line/SKILL.md`. Include `.agents/skills/bench-craft-delegate/SKILL.md` and its `references/delegation-discipline.md`. The rendered reader is `docs/field-guide.html`.

The docs check calls `checkWorkflowAnchors`, which calls `anchors.EvaluateGroup`, `checkReviewConvergenceContract`, `checkIntegrationSourceWorkflowCurrency`, and `checkSpecAuthorizationContract`. These helpers live in the authorized workflow conformance files. Also update the independent drain expectations in `internal/conformance/recurrence_maintenance_contract_test.go`. Preserve its recurrence, trust, and ignored-capture guarantees while changing the writer route.

Review repair coverage: W17, W18, W19, W20, W21, W22, W23. These rows refine the source clauses already mapped above.

Plan-edit authority is limited to chunk boundaries, dependencies, row assignments, Writes expectations, and necessary ownership expansion under decision #5. Preserve the approved behavior and pass criteria. A material acceptance change still needs the user's decision. Bulk rewrites exclude this spec and its tickets. Completion evidence must review a changed plan digest before it can satisfy the next checkpoint.

Authoring close: this spec and its tickets are staged for user sign-off. The reviews above assess the proposed build. No implementation or implementation test result is claimed.

Reviewer amendment on 2026-09-11: apply decision #13 to the declared implementation model. Sol implementations use Astra/high review axes. This amendment follows the spec-authoring reviews recorded above.

Execute tickets in dependency order: 3.md, 1.md, then 2.md. The binding ticket enables Astra before the first review. This order preserves each ticket checkpoint and each chunk review.
