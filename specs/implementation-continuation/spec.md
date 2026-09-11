# Continue the retained implementation session while it advances

Status: staged
Decision source: `docs/adr/0021-benchmark-workflow-orchestration.md`
Verification log: 2 iteration(s) to accept — Sol/high approved the repaired spec and ticket graph. Attempt counting, stop causes, and diagnostic boundaries were corrected.

## Problem

A numeric iteration cap can interrupt useful approved work. Repeating the same failed approach also wastes work. The retained author needs a clear distinction between continued progress, diagnosis, and a real stopping condition.

## Solution

Support an explicit uncapped declaration for the retained implementation session. Keep progress observable and use a small diagnostic escalation after repeated nonprogress. Preserve the user budget and the existing external-blocker boundary.

Build dependency: retained-implementation-workflow.

## User stories

Line: gpt-5.6-sol / high. One retained implementation session. Starting effort is a recommendation for approval.

Sol/high is recommended because the continuation policy has closed behavioral decisions and known guidance owners. Chunk 2 needs more care because it defines the diagnostic exception.

1. As an operator, I want to choose uncapped continuation explicitly, so that I can assess and continue approved work accurately.
2. As an operator, I want to continue while evidence advances the work, so that I can assess and continue approved work accurately.
3. As an operator, I want to count meaningful attempts, so that I can assess and continue approved work accurately.
4. As an operator, I want to reassess repeated nonprogress, so that I can assess and continue approved work accurately.
5. As an operator, I want to use focused debugging, so that I can assess and continue approved work accurately.
6. As an operator, I want to obtain brief expert advice, so that I can assess and continue approved work accurately.
7. As an operator, I want to respect explicit budgets, so that I can assess and continue approved work accurately.
8. As an operator, I want to surface required decisions and blockers, so that I can assess and continue approved work accurately.
9. As an operator, I want to preserve control of authorship, so that I can assess and continue approved work accurately.
10. As an operator, I want to keep the separate runner behavior stable, so that I can assess and continue approved work accurately.

11. As an operator, I want to continue from useful evidence, so that the required evidence remains complete.
12. As an operator, I want to count complete attempts, so that the required evidence remains complete.
13. As an operator, I want to choose a discriminating check, so that the required evidence remains complete.
14. As an operator, I want to honor a selected numeric cap, so that the required evidence remains complete.
15. As an operator, I want to wait for required decisions, so that the required evidence remains complete.
16. As an operator, I want to surface external blockage, so that the required evidence remains complete.
17. As an operator, I want to adjust effort while retaining the model, so that the required evidence remains complete.
18. As an operator, I want to consult an expert without a new approval, so that the required evidence remains complete.
19. As an operator, I want to keep consultation context small, so that the required evidence remains complete.
20. As an operator, I want to obtain a useful short response, so that the required evidence remains complete.
21. As an operator, I want to retain the implementation model, so that the required evidence remains complete.
22. As an operator, I want to stop when cancelled, so that the required evidence remains complete.
23. As an operator, I want to keep diagnosis out of attempt counts, so that the required evidence remains complete.

## Implementation decisions

This is portable guidance for the retained implementation session. It does not change `bench shift`, its adapter process, its iteration limits, or its wall-time controls. The reviewer confirmed that boundary on 2026-09-11.

The line declaration names model, starting effort, and either a numeric cap or an explicit uncapped policy. Exhaustion of a selected numeric cap stops the run. Uncapped means no artificial iteration stop within the approved spec. It does not imply a paid comparison trial, a background task, an automatic restart, or an unlimited user budget. Keep the declaration in the normal session handoff and progress updates.

Progress means a verified acceptance improvement or new useful evidence that changes the next action. An attempt is one coherent hypothesis, implementation change, and verification. A diagnostic-only action is not a completed implementation-and-verification attempt. Expected TDD reds and individual tool calls are not attempts. Two completed attempts with neither progress nor new useful evidence trigger reassessment before another implementation attempt.

Reassessment states the changed hypothesis and the next discriminating check. The author may use `$bench-debug`, change effort, or ask a higher-tier model for brief diagnostic help. Native consultation is pre-approved through top tier. Give one question, the relevant error, minimal code, and attempted hypotheses.

Request a short diagnosis and next check. The helper does not write implementation or repairs. Record the actual line used. If that model is unavailable, keep the failure explicit and use an available authorized diagnostic route without substituting an undeclared model.

Stop for a required user decision, an external blocker with no remaining independent work, or an explicit user budget. Stop immediately when the user cancels. A stalled hypothesis alone calls for reassessment. It is not authority to switch implementation models or sessions. If no useful next check remains, report the unresolved blocker and the smallest decision or evidence needed. A handoff preserves completed work, outstanding acceptance, attempts, and the next check.

## Implementation chunks

Implement these tickets in the retained session. Each ticket is one initial review chunk. Commit the ticket on its lane pass, run the three delegated axes, and repair findings before the next chunk. Review line: decision #13, resolved for this implementation as gpt-6-astra / high. Final acceptance reconciliation stays with the author.

| chunk / ticket | blocked by | delivered outcome | harder chunk |
| --- | --- | --- | --- |
| 1.md — Support explicit continued progress | none | Update the declaration, progress, attempt, and stop rules for the retained session | no |
| 2.md — Provide bounded diagnostic escalation | 1.md | Connect nonprogress to debug, effort adjustment, and brief higher-tier advice | yes |

## Testing decisions

Use the named existing owner and new tests below. A new test name is a planned seam, not a claim that the test exists. Read its nearest fixture before implementation. Demonstrate each required omission or behavioral mutation as a diagnostic red, restore it, and show green. A compile failure is not the required red. Semantic review judges prose quality and the sufficiency of the test.

### Seam diagram

```text
retained author -> progress / attempt policy -> next check or stop
                               |
                       diagnostic helper (read only)
                               |
                    anchors + negative fixtures -> docs gate
```

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| C1 | 1 | The line declaration supports an explicit uncapped policy | New test: TestImplementationContinuation in internal/conformance/implementation_continuation_test.go | Restore a mandatory numeric cap for the retained session and observe red. |
| C2 | 2 | Verified acceptance improvement permits continuation | New test: TestImplementationContinuation in internal/conformance/implementation_continuation_test.go | Remove the verified-improvement continuation arm and observe red. |
| C3 | 3 | Expected TDD reds do not consume attempts | New test: TestImplementationContinuation in internal/conformance/implementation_continuation_test.go | Count each expected red as an attempt in the fixture and observe red. |
| C4 | 4 | Two nonprogress attempts require a changed hypothesis | New test: TestImplementationContinuation in internal/conformance/implementation_continuation_test.go | Permit a third identical attempt in the fixture and observe red. |
| C5 | 5 | The retained author can invoke the existing debug skill | New test: TestImplementationContinuation in internal/conformance/implementation_continuation_test.go | Remove the debug route and observe red. |
| C6 | 6 | Higher-tier consultation remains diagnostic only | New test: TestImplementationContinuation in internal/conformance/implementation_continuation_test.go | Permit the helper to repair code and observe red. |
| C7 | 7 | A user budget stops uncapped continuation | New test: TestImplementationContinuation in internal/conformance/implementation_continuation_test.go | Remove the budget stop condition and observe red. |
| C8 | 8 | A blocked run reports the next needed decision or evidence | New test: TestImplementationContinuation in internal/conformance/implementation_continuation_test.go | Remove the actionable handoff requirement and observe red. |
| C9 | 9 | A diagnostic escalation does not switch the implementation session | New test: TestImplementationContinuation in internal/conformance/implementation_continuation_test.go | Conflate the helper model with the implementation line and observe red. |
| C10 | 10 | The continuation change leaves shift runtime limits unchanged | review-owned: Standards compares the fence and final diff | Reject a shift source or default change under this guidance scope. |

| C11 | 11 | Useful new evidence permits continuation | New test: TestImplementationContinuation in internal/conformance/implementation_continuation_test.go | Remove the new-evidence continuation arm. |
| C12 | 12 | Individual tool calls do not consume attempts | New test: TestImplementationContinuation in internal/conformance/implementation_continuation_test.go | Count two tool calls as two implementation attempts. |
| C13 | 13 | Reassessment names the next discriminating check | New test: TestImplementationContinuation in internal/conformance/implementation_continuation_test.go | Retain a changed hypothesis but omit the next check. |
| C14 | 14 | Exhaustion of a selected numeric cap stops continuation | New test: TestImplementationContinuation in internal/conformance/implementation_continuation_test.go | Retain uncapped support but delete numeric-cap exhaustion. |
| C15 | 15 | A required user decision stops dependent implementation | New test: TestImplementationContinuation in internal/conformance/implementation_continuation_test.go | Continue past a required scope decision. |
| C16 | 16 | An external blocker with no independent work stops the run | New test: TestImplementationContinuation in internal/conformance/implementation_continuation_test.go | Repeat work after the only required service remains unavailable. |
| C17 | 17 | The retained author may adjust effort | New test: TestImplementationContinuation in internal/conformance/implementation_continuation_test.go | Require user approval for every effort adjustment. |
| C18 | 18 | Brief higher-tier diagnostic help is pre-approved through top tier | New test: TestImplementationContinuation in internal/conformance/implementation_continuation_test.go | Require a new approval before top-tier diagnostic advice. |
| C19 | 19 | A diagnostic request contains only the bounded question context | New test: TestImplementationContinuation in internal/conformance/implementation_continuation_test.go | Request the whole session history instead of the question error code and attempted hypotheses. |
| C20 | 20 | The helper returns a short diagnosis and next check | New test: TestImplementationContinuation in internal/conformance/implementation_continuation_test.go | Replace the response contract with a general report. |
| C21 | 21 | Diagnostic escalation does not switch the implementation model | New test: TestImplementationContinuation in internal/conformance/implementation_continuation_test.go | Treat the helper line as permission to switch the author model. |
| C22 | 22 | User cancellation stops continuation | New test: TestImplementationContinuation in internal/conformance/implementation_continuation_test.go | Continue an uncapped run after explicit cancellation. |
| C23 | 23 | A diagnostic-only action does not consume an implementation attempt | New test: TestImplementationContinuation in internal/conformance/implementation_continuation_test.go | Count a read-only diagnosis as a completed implementation attempt. |

### Edge inventory

The implementation phase and craft-line are the progress-policy readers. The operating guide references that policy. Craft-delegate and its reference consume the diagnostic exception. The existing anchors registry and docs conformance root run the new negative fixtures. The shift loop and session adapter were inspected to establish the excluded runner boundary.

The check enforces the written contract. It cannot infer the usefulness of an agent's new hypothesis. Spec and Coverage reviewers judge recorded attempts in ordinary work. No telemetry threshold becomes an additional stopping oracle.

## Ownership fences

These paths are the union of ticket expectations. A directory entry is an exact prefix for that existing owner or fixture family. Expansion follows decision #5, with the plan updated before use. It cannot weaken existing guarantees.

- `.bench/BENCH.md`
- `.agents/commands/bench-implement-spec.md`
- `.agents/skills/bench-craft-line/SKILL.md`
- `internal/anchors/registry_data.go`
- `internal/conformance/implementation_continuation_test.go` (new)
- `tests/canary/workflow-guidance-anchors`
- `.agents/skills/bench-craft-delegate/SKILL.md`
- `.agents/skills/bench-craft-delegate/references/delegation-discipline.md`
- `specs/implementation-continuation/spec.md`
- `specs/implementation-continuation/tickets`
- `reviews/implementation-continuation.md` (new)

- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `tests/canary/docs-currency-token-diet/benchref-imported`
- `tests/canary/docs-currency-token-diet/benchref-pointer-dropped`
- `tests/canary/docs-currency-token-diet/benchref-section-duplicated`
- `tests/canary/docs-currency-token-diet/dogfood-referent-shipped`
- `tests/canary/docs-currency-token-diet/missing-cli-inventory`
- `tests/canary/docs-currency-token-diet/stale-cli-doc-reference`
- `tests/canary/load-validity-metadata/readme-shared-rule-drift`
- `tests/canary/load-validity-metadata/shared-rule-drift`
- `tests/canary/skills-index-command-adapters/adapter-inert-invocation-key`
- `tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy`
- `tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted`

## Out of scope

Changing bench shift, adding restart automation, choosing an implementation model automatically, and inventing a paid-trial budget are outside this spec.

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
| #3: "Work continues within approved scope while it makes progress" | C1, C2, C7, C8 |
| #6: reassess after two attempts and permit brief diagnostic help | C3, C4, C5, C6, C9 |
| #3 dated clarification: retained session only | C10 |

Plan updates follow source decision #5 and preserve approved behavior. Bulk rewriting excludes this spec and its tickets.

Review repair coverage: C11, C12, C13, C14, C15, C16, C17, C18, C19, C20, C21, C22, C23. These rows refine the source clauses already mapped above.

Authoring close: this spec and its tickets are staged for user sign-off. The reviews above assess the proposed build. No implementation or implementation test result is claimed.

Reviewer amendment on 2026-09-11: apply decision #13 to the declared implementation model. Sol implementations use Astra/high review axes. This amendment follows the spec-authoring reviews recorded above.
