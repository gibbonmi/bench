# Require source-bound review evidence at completion checkpoints

Status: staged
Decision source: `docs/adr/0021-benchmark-workflow-orchestration.md`
Verification log: 2 iteration(s) to accept — Sol/high verified the design repairs. The author folded its final three E21 prose corrections and checked them against the coverage row.

## Problem

The current preflight collects review inputs. It does not record terminal reviewer results. Clean reviews leave no durable pickup artifact, so a later session cannot distinguish a completed clean review from one that never ran.

## Solution

Extend the existing review artifact with one machine-readable result record. The gate checks that record at chunk advancement and final landing. Keep author verification distinct from independent review and retain the source identity examined by each performer.

Build dependency: retained-implementation-workflow.

## User stories

Line: gpt-6-astra / high. One retained implementation session. Starting effort is a recommendation for approval.

Astra/high is recommended because the hardest chunks change source identity, gate reuse, and the publication boundary. These guarantees require careful composition across existing owners.

1. As an operator, I want to retain clean review outcomes, so that I can assess and continue approved work accurately.
2. As an operator, I want to separate author and reviewer evidence, so that I can assess and continue approved work accurately.
3. As an operator, I want to identify examined source, so that I can assess and continue approved work accurately.
4. As an operator, I want to distinguish unfinished review, so that I can assess and continue approved work accurately.
5. As an operator, I want to distinguish failed review, so that I can assess and continue approved work accurately.
6. As an operator, I want to distinguish skipped review, so that I can assess and continue approved work accurately.
7. As an operator, I want to detect missing review, so that I can assess and continue approved work accurately.
8. As an operator, I want to detect stale review, so that I can assess and continue approved work accurately.
9. As an operator, I want to cover repairs, so that I can assess and continue approved work accurately.
10. As an operator, I want to retain each review obligation, so that I can assess and continue approved work accurately.
11. As an operator, I want to keep ordinary work possible, so that I can assess and continue approved work accurately.
12. As an operator, I want to prevent weak verdict reuse, so that I can assess and continue approved work accurately.
13. As an operator, I want to avoid self-staling review records, so that I can assess and continue approved work accurately.
14. As an operator, I want to bind gate evidence to record contents, so that I can assess and continue approved work accurately.
15. As an operator, I want to cover the complete spec, so that I can assess and continue approved work accurately.
16. As an operator, I want to block publication before proof, so that I can assess and continue approved work accurately.
17. As an operator, I want to cover destination composition, so that I can assess and continue approved work accurately.
18. As an operator, I want to keep trusted completion possible, so that I can assess and continue approved work accurately.
19. As an operator, I want to reject malformed record input, so that I can assess and continue approved work accurately.
20. As an operator, I want to read hostile paths safely, so that I can assess and continue approved work accurately.
21. As an operator, I want to retain inspectable native results, so that I can assess and continue approved work accurately.
22. As an operator, I want to report evidence without overstating it, so that I can assess and continue approved work accurately.
23. As an operator, I want to retain valid earlier chunk evidence, so that I can assess and continue approved work accurately.
24. As an operator, I want to review plan amendments, so that I can assess and continue approved work accurately.

25. As an operator, I want to require author verification, so that the required evidence remains complete.
26. As an operator, I want to reject pending verification, so that the required evidence remains complete.
27. As an operator, I want to reject failed verification, so that the required evidence remains complete.
28. As an operator, I want to reject stale verification, so that the required evidence remains complete.
29. As an operator, I want to identify the verifier, so that the required evidence remains complete.
30. As an operator, I want to require meaningful probes, so that the required evidence remains complete.
31. As an operator, I want to require probe restoration, so that the required evidence remains complete.
32. As an operator, I want to separate reviews from test execution, so that the required evidence remains complete.
33. As an operator, I want to execute final integration checks, so that the required evidence remains complete.
34. As an operator, I want to reject failed final verification, so that the required evidence remains complete.
35. As an operator, I want to reject stale final verification, so that the required evidence remains complete.
36. As an operator, I want to limit the trusted status transform, so that the required evidence remains complete.
37. As an operator, I want to reach the public checkpoint owner, so that the required evidence remains complete.
38. As an operator, I want to retain occurrence history, so that the required evidence remains complete.

## Implementation decisions

The existing `reviews/<slug>.md` is the durable review and completion record. It retains the human findings and one fenced `bench-review-record` JSON payload. `internal/reviewrecord` is the single parser and source-coverage owner used by preflight and the gate. It is a deep data module with two actual consumers. It does not execute reviewers, grant approval, or decide semantic correctness.

Version 1 fields are `version`, `spec`, `plan_digest`, `implementation_session`, `chunks`, and `completion`. Each chunk has `id`, `base`, `tip`, `source_digest`, `acceptance_rows`, `verification`, and `reviews`. Each evidence item has `id`, `performer`, `role`, `model`, `effort`, `source_digest`, `state`, `outcome`, and `native_ref`. Verification also names `command`, `exit_code`, and any required probe mutation and restore result.

Review adds `axis`, `finding_ids`, and `supersedes`. Completion records `state`, the final source digest, the author's acceptance reconciliation, and final `verification` items. Final verification names the executed overall acceptance and integration commands, performer, source, outcome, and exit code. Unknown model or effort is explicit unknown, never an inferred value.

Stored occurrence states are pending, completed, failed, and skipped. Validation states are missing, stale, invalid, and current. Missing, stale, and invalid are computed and are never imported as occurrence states. Missing is derived when a required entry is absent. Stale is derived when its identity or review chain does not cover the requested subject. Completed with zero findings is a positive terminal result.

Failed transport, skipped review, and a request sent to a reviewer cannot become completed. Native references name the harness result or command log and its digest. Import or attach the minimal result excerpt needed for cold inspection in the record. Local-only links can supplement it, but are not the sole durable evidence.

The gate checks fields and source coverage. A record cannot mechanically prove the review judgment correct or authenticate an invented harness transcript.

The source digest is a canonical Git tree identity over the checked source, excluding only the exact review record for this spec. Include the spec, ticket plan, tests, and gate code. Bind a completed review to frozen base and tip as well. The full gate subject still includes the review-record bytes.

Thus a record update does not invalidate its own reviewed source, but it does invalidate prior gate evidence. The record path cannot exclude a directory or another source file. Use the existing Git tree and gate subject owners, not a second repository walker.

Each chunk starts from the prior covered source. The three axis results cover its entire delta against the whole spec. Later deltas require their own review. Earlier records stay historical evidence in the continuous chain.

A gap, wrong base, unrelated tip, plan change, or uncovered repair is stale. A revised chunk plan requires an explicit ID mapping and review of the plan delta. A repair can use a scoped follow-up, but every axis must have a current result or a reviewer-authored reaffirmation bound to the repaired tip. The author cannot self-exempt an axis.

The final chain ends at the source being completed. A later cross-chunk concern reopens the affected evidence and blocks completion until resolved.

Extend the public gate grammar with `bench gate --checkpoint <spec-path> --chunk <id>` and `bench gate --checkpoint <spec-path> --complete`. These are gate-owned review obligations. A successful chunk checkpoint permits the phase to advance. It does not create a separate approval record or run a second lifecycle engine.

Checkpoint purpose and record identity participate in gate evidence reuse. An ordinary green verdict cannot answer a stronger checkpoint request. Reuse existing exact verification only through the gate owner; no phase-side cache shortcuts.

Ordinary ticket lane checks tolerate an absent record, an active chunk, pending review, repair in progress, and a staged sibling spec. They do not claim chunk completion. A checkpoint requires the completed current chunk's author verification and all three review axes. Final completion also requires every planned row, a continuous review chain, and final integration reconciliation.

Missing or stale evidence fails closed at those boundaries with the chunk, axis or verification item, reason, and next record or review action. Syntax errors fail before source traversal. Evidence failures precede any publication.

Spec-backed landing passes the complete obligation into the existing prospective gate owner. It binds review coverage to the reviewed source and the exact landing composition. A destination change that alters the reviewed content requires review of that delta. The existing broker remains the independent publication authority. Its trusted status flip can be recognized only as the exact existing completion transform, not an arbitrary path exclusion. The final whole-project gate still grades the prospective published tree.

Use bounded regular-file reads, reject path escape and symlink traversal, and reject unknown versions and duplicate IDs. The parser is read-only. The phase session embeds terminal native results in the review artifact through its native file tools. No record writer runs a shell command found in a record.

A missing legacy record remains visible. Unchanged historical implemented specs do not retroactively block ordinary work. A newly requested completion checkpoint has no legacy bypass.

## Implementation chunks

Implement these tickets in the retained session. Each ticket is one initial review chunk. Commit the ticket on its lane pass, run the three delegated axes, and repair findings before the next chunk. Review line: decision #13, resolved for this implementation as gpt-5.6-sol / high. Final acceptance reconciliation stays with the author.

| chunk / ticket | blocked by | delivered outcome | harder chunk |
| --- | --- | --- | --- |
| 1.md — Retain source-bound review outcomes | none | Extend the existing review artifact and preflight instructions with the canonical result format | no |
| 2.md — Check evidence before chunk advancement | 1.md | Add gate checkpoint obligations and connect the implementation phase | yes |
| 3.md — Require complete evidence before landing | 2.md | Pass the final obligation from spec-backed landing to the prospective gate | yes |

## Verification inventory

The fenced plan names required commands and probes. Each chunk derives its rows and dependencies from its named tickets.
The parser hashes the plan, spec, and ticket bytes. A plan amendment must map old chunk IDs to new IDs.

A review occurrence also records its frozen `base` and `tip`. Each chunk records its own `plan_digest`.
These fields preserve earlier evidence when a later reviewed delta changes the plan.

```bench-completion-plan
{
  "version": 1,
  "chunks": [
    {
      "id": "1",
      "tickets": [
        "1.md",
        "r1.md"
      ],
      "verification": [
        {
          "id": "record-tests",
          "command": "bench test --package ./internal/reviewrecord --run TestReviewRecord",
          "probe": "omit missing-axis rejection"
        },
        {
          "id": "preflight-tests",
          "command": "bench test --package ./internal/preflight --run TestReviewCharge"
        }
      ]
    },
    {
      "id": "2",
      "tickets": [
        "2.md",
        "r2.md"
      ],
      "verification": [
        {
          "id": "checkpoint-tests",
          "command": "bench test --package ./internal/gate --run TestReviewCheckpoint",
          "probe": "omit checkpoint evidence validation"
        },
        {
          "id": "axis-tests",
          "command": "bench test --package ./internal/gate --run TestReviewCheckpointCanonicalAxes",
          "probe": "omit canonical Coverage axis"
        },
        {
          "id": "route-tests",
          "command": "bench test --package ./cmd/bench --run TestGateCheckpointRoute"
        }
      ]
    },
    {
      "id": "3",
      "tickets": [
        "3.md"
      ],
      "verification": [
        {
          "id": "landing-tests",
          "command": "bench test --package ./internal/landing --run TestLandingCompletionEvidence",
          "probe": "omit landing completion obligation"
        }
      ]
    }
  ],
  "final_verification": [
    {
      "id": "acceptance",
      "command": "bench test --changed --base de1447b31903b170679b66bd200431cc15ddc73f"
    },
    {
      "id": "integration",
      "command": "bench test --check system"
    }
  ]
}
```

## Testing decisions

Use the named existing owner and new tests below. A new test name is a planned seam, not a claim that the test exists. Read its nearest fixture before implementation. Demonstrate each required omission or behavioral mutation as a diagnostic red, restore it, and show green. A compile failure is not the required red. Semantic review judges prose quality and the sufficiency of the test.

### Seam diagram

```text
native results -> existing review artifact -> reviewrecord parser
                                                 |
phase checkpoint / landing broker -> gate obligation -> publish or refuse
                                                 |
                               exact source + record + purpose identity
```

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| E1 | 1 | A zero-finding axis has a durable completed result | `internal/reviewrecord/record_test.go` (`TestReviewRecord`) | Delete the clean result before loading and observe a missing axis. |
| E2 | 2 | Author verification cannot satisfy a review axis | `internal/reviewrecord/source_test.go` (`TestReviewRecordTerminal`) | Relabel a test command as independent review and require refusal. |
| E3 | 3 | Every terminal result binds its source and performer | `internal/reviewrecord/source_test.go` (`TestReviewRecordTerminal`) | Remove the source digest from a completed item and require refusal. |
| E4 | 4 | A pending axis blocks a checkpoint | `internal/gate/review_checkpoint_test.go` (`TestReviewCheckpoint`) | Supply two complete axes and one pending axis. |
| E5 | 5 | A failed axis blocks a checkpoint | `internal/gate/review_checkpoint_test.go` (`TestReviewCheckpoint`) | Supply a failed transport result with an empty findings list. |
| E6 | 6 | A skipped axis blocks a checkpoint | `internal/gate/review_checkpoint_test.go` (`TestReviewCheckpoint`) | Supply a skipped axis with a success-looking summary. |
| E7 | 7 | An absent required axis blocks a checkpoint | `internal/gate/review_checkpoint_test.go` (`TestReviewCheckpoint`) | Omit Coverage from an otherwise valid record. |
| E8 | 8 | An uncovered source delta blocks a checkpoint | `internal/gate/review_checkpoint_test.go` (`TestReviewCheckpointLaterSource`) | Edit one source file after recording the reviewed tip. |
| E9 | 9 | An unreviewed repair blocks advancement | `internal/gate/review_checkpoint_test.go` (`TestReviewCheckpointLaterSource`) | Append a repair commit without axis coverage for its tip. |
| E10 | 10 | A completed checkpoint requires all three canonical axes | `internal/gate/review_checkpoint_test.go` (`TestReviewCheckpointCanonicalAxes`) | Remove one axis from the required inventory and demonstrate the independent omission red. |
| E11 | 11 | An ordinary lane tolerates incomplete review state | `internal/gate/review_checkpoint_test.go` (`TestReviewCheckpointOrdinaryWork`) | Exercise absent active pending repair and sibling-spec fixtures through the normal lane. |
| E12 | 12 | Ordinary green cannot satisfy a completion checkpoint | `internal/gate/review_checkpoint_test.go` (`TestReviewCheckpointReuse`) | Seed ordinary green before requesting a stronger obligation. |
| E13 | 13 | A record-only update preserves the reviewed source identity | `internal/reviewrecord/source_test.go` (`TestReviewRecordSource`) | Append a native result while leaving source bytes unchanged. |
| E14 | 14 | Changed record bytes invalidate checkpoint verdict reuse | `internal/gate/review_checkpoint_test.go` (`TestReviewCheckpointReuse`) | Remove a completed axis after a green checkpoint. |
| E15 | 15 | Completion requires final acceptance reconciliation | New test: TestLandingCompletionEvidence in internal/landing/completion_evidence_test.go | Leave one planned row without a final disposition. |
| E16 | 16 | A landing with missing completion evidence publishes no ref | New test: TestLandingCompletionEvidence in internal/landing/completion_evidence_test.go | Pass a valid source pair with an incomplete review record. |
| E17 | 17 | A new destination delta cannot inherit source-only review | New test: TestLandingCompletionEvidence in internal/landing/completion_evidence_test.go | Compose a destination change that modifies reviewed content. |
| E18 | 18 | The exact broker status transform preserves review coverage | New test: TestLandingCompletionEvidence in internal/landing/completion_evidence_test.go | Run the existing status flip through the complete gate obligation. |
| E19 | 19 | Invalid records fail before source traversal | `internal/reviewrecord/source_test.go` (`TestReviewRecordTerminal`) | Compete a malformed version with an escaping native reference. |
| E20 | 20 | A nonregular or escaping record path is refused | `internal/reviewrecord/source_test.go` (`TestReviewRecordPaths`, `TestReviewRecordControlPath`, `TestReviewRecordLiteralTickets`) | Use FIFO symlink traversal and control-byte fixtures without opening the target. |
| E21 | 21 | A retained native result remains inspectable without its local log | `internal/reviewrecord/source_test.go` (`TestReviewRecordSource`) | Remove the supplemental local log and load the embedded terminal result. |
| E22 | 22 | The result distinguishes review occurrence from judgment correctness | review-owned: Spec checks phase output and field descriptions | Reject any claim that record validation proves semantic correctness. |
| E23 | 23 | A continuous reviewed chain covers completed chunks | `internal/reviewrecord/source_test.go` (`TestReviewRecordSource`) | Cover two chunk deltas and then mutate the second base to create a gap. |
| E24 | 24 | An unreviewed plan amendment makes completion stale | `internal/reviewrecord/source_test.go` (`TestReviewRecordPlanAmendment`) | Change chunk rows or dependency mappings after the last review. |

| E25 | 25 | Missing required author verification blocks a checkpoint | `internal/gate/review_checkpoint_test.go` (`TestReviewCheckpoint`) | Provide all review axes and omit one required test result. |
| E26 | 26 | Pending author verification blocks a checkpoint | `internal/gate/review_checkpoint_test.go` (`TestReviewCheckpoint`) | Provide a pending test with complete review axes. |
| E27 | 27 | A nonzero required verification exit blocks a checkpoint | `internal/gate/review_checkpoint_test.go` (`TestReviewCheckpoint`) | Set exit code one with an otherwise complete record. |
| E28 | 28 | A stale verification source blocks a checkpoint | `internal/gate/review_checkpoint_test.go` (`TestReviewCheckpoint`) | Bind the test to the pre-repair source. |
| E29 | 29 | Missing verification performer blocks a checkpoint | `internal/gate/review_checkpoint_test.go` (`TestReviewCheckpoint`) | Omit the performer from a passing command result. |
| E30 | 30 | A missing required mutation result blocks a checkpoint | `internal/gate/review_checkpoint_test.go` (`TestReviewCheckpoint`) | Keep tests green but omit the planned probe evidence. |
| E31 | 31 | A failed required probe restore blocks a checkpoint | `internal/gate/review_checkpoint_test.go` (`TestReviewCheckpoint`) | Record a biting mutation with failed restoration. |
| E32 | 32 | Review evidence cannot substitute for author verification | `internal/gate/review_checkpoint_test.go` (`TestReviewCheckpoint`) | Populate only review items for a required command. |
| E33 | 33 | Missing final integration execution blocks completion | New test: TestLandingCompletionEvidence in internal/landing/completion_evidence_test.go | Supply reconciliation prose with no final command result. |
| E34 | 34 | Failed final integration execution blocks completion | New test: TestLandingCompletionEvidence in internal/landing/completion_evidence_test.go | Record a final integration command with a nonzero exit. |
| E35 | 35 | Stale final integration execution blocks completion | New test: TestLandingCompletionEvidence in internal/landing/completion_evidence_test.go | Change source after the final command ran. |
| E36 | 36 | An extra spec-byte change beside the status flip blocks publication | New test: TestLandingCompletionEvidence in internal/landing/completion_evidence_test.go | Change one acceptance byte beside the valid status transform before updateRef. |
| E37 | 37 | The public gate wrapper forwards checkpoint arguments | `cmd/bench/gate_route_test.go` (`TestGateCheckpointRoute`) | Drive chunk and complete forms through bin/bench.sh. |
| E38 | 38 | Repairs append new evidence without erasing earlier outcomes | `internal/reviewrecord/source_test.go` (`TestReviewRecordTerminal`) | Load initial findings and their superseding repair results. |

### Edge inventory

Read owners: preflight collects diff, consumers, and coverage inputs in `internal/preflight/review.go`. The review phase currently omits a clean artifact. The gate command and subject builder own execution and reuse. Prospective grading belongs to `internal/gate/prospective.go`; the landing broker remains the publisher. New tests must reach these public paths rather than a second test-only oracle.

The parser fixtures include every state in E4 through E9. They also cover empty and absent files, bad versions, duplicate IDs, and wrong source. Include plan changes, chain gaps, oversized data, unsafe paths, embedded native results, and legacy records. A missing `reviews/` directory and an empty directory both mean no evidence at a requested checkpoint. In ordinary work both are permitted. Required native result data must remain inspectable after a local log disappears.

Gate fixtures cover verdict persistence failure before execution, interruption during execution, and failure to persist the terminal verdict. None may produce reusable checkpoint green. Existing exact-verdict and broker tests remain in force. Introduce focused tests beside these owners, with a second test file if command-level scenarios exceed the local file budget. Ordinary Go test phases execute them; any new package joins the existing package and injected-port registries. Prove omission reds by removing the checkpoint call at advancement and at landing separately.

No signed-attestation system is promised. The threat model covers omitted, malformed, stale, and mismatched evidence under the existing trusted harness and broker boundary. Review-owned acceptance E22 forbids claiming protection against a fabricated transcript.

## Ownership fences

These paths are the union of ticket expectations. A directory entry is an exact prefix for that existing owner or fixture family. Expansion follows decision #5, with the plan updated before use. It cannot weaken existing guarantees.

- `internal/reviewrecord` (new)
- `internal/git/tree.go`
- `CHANGELOG.md`
- `tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns`
- `tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary`
- `internal/preflight/review.go`
- `internal/preflight/review_charge_test.go`
- `.agents/commands/bench-review-implementation.md`
- `internal/anchors/registry_ft311_review_dispatch.go`
- `internal/anchors/registry_data.go`
- `tests/canary/workflow-guidance-anchors/review-clean-terminal-result`
- `internal/conformance/registry/packages.go`
- `internal/conformance/injected_ports_registry_test.go`
- `internal/gate`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `.agents/commands/bench-implement-spec.md`
- `.bench/BENCH-reference.md`
- `internal/preflight/charge.go`
- `internal/preflight/charge_cells_test.go`
- `internal/landing`
- `internal/worktree`
- `.agents/commands/bench-final-check.md`
- `projects/benchkit.md`
- `specs/completion-evidence/spec.md`
- `specs/completion-evidence/tickets`
- `reviews/completion-evidence.md` (new)

- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `tests/canary/docs-currency-token-diet/benchref-imported`
- `tests/canary/docs-currency-token-diet/benchref-pointer-dropped`
- `tests/canary/docs-currency-token-diet/benchref-section-duplicated`
- `tests/canary/docs-currency-token-diet/introduces-undeclared-command`
- `tests/canary/docs-currency-token-diet/stale-command-reference`
- `tests/canary/guidance-prose-budgets/over-budget-skill`
- `tests/canary/line-routing/line-binding-prose-drift`
- `tests/canary/skills-index-command-adapters/adapter-inert-invocation-key`
- `tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy`
- `tests/canary/skills-index-command-adapters/dangling-index`
- `tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted`
- `tests/canary/skills-index-command-adapters/missing-index-field`
- `tests/canary/skills-index-command-adapters/stale-index-wording`
- `tests/canary/skills-index-command-adapters/unindexed-skill`
- `tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading`
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-owner`
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-routing`
- `tests/canary/workflow-guidance-anchors/benchkit-spec-ownership`
- `tests/canary/workflow-guidance-anchors/benchkit-system-suite-route`
- `tests/canary/workflow-guidance-anchors/coverage-axis-anchor`
- `tests/canary/workflow-guidance-anchors/final-check-bare-leftover-clean-retired`
- `tests/canary/workflow-guidance-anchors/final-check-census-read-before-land`
- `tests/canary/workflow-guidance-anchors/final-check-landed-worktree-sweep`
- `tests/canary/workflow-guidance-anchors/final-check-light-path-changelog-heading`
- `tests/canary/workflow-guidance-anchors/final-check-scratch-branch-clean`
- `tests/canary/workflow-guidance-anchors/implement-spec-coverage-task-seeding`
- `tests/canary/workflow-guidance-anchors/implement-spec-entry-validation`
- `tests/canary/workflow-guidance-anchors/implement-spec-incapable-harness`
- `tests/canary/workflow-guidance-anchors/implement-spec-inline-exception`
- `tests/canary/workflow-guidance-anchors/implement-spec-mandatory-delegation-anchor`
- `tests/canary/workflow-guidance-anchors/implement-spec-offer-retired`
- `tests/canary/workflow-guidance-anchors/implement-spec-offer-scope`
- `tests/canary/workflow-guidance-anchors/implement-spec-read-only-helper`
- `tests/canary/workflow-guidance-anchors/implement-spec-status-flip-anchor`
- `tests/canary/workflow-guidance-anchors/implement-spec-worktree-before-preflight`
- `tests/canary/workflow-guidance-anchors/implement-spec-write-delegation`
- `tests/canary/workflow-guidance-anchors/line-anchor-missing`
- `tests/canary/workflow-guidance-anchors/prepared-build-approval`
- `tests/canary/workflow-guidance-anchors/prepared-build-freshness`
- `tests/canary/workflow-guidance-anchors/prepared-review-axis-returns`
- `tests/canary/workflow-guidance-anchors/prepared-review-blast-evidence`
- `tests/canary/workflow-guidance-anchors/prepared-review-capable-handoff`
- `tests/canary/workflow-guidance-anchors/prepared-review-inline-axis-route`
- `tests/canary/workflow-guidance-anchors/prepared-review-legacy-entry-points`
- `tests/canary/workflow-guidance-anchors/prepared-review-native-dispatch`
- `tests/canary/workflow-guidance-anchors/prepared-review-runtime-capability`
- `tests/canary/workflow-guidance-anchors/prepared-review-shared-evidence`
- `tests/canary/workflow-guidance-anchors/reference-agent-push-rule`
- `tests/canary/workflow-guidance-anchors/reference-bench-operational-layer`
- `tests/canary/workflow-guidance-anchors/reference-category-context`
- `tests/canary/workflow-guidance-anchors/reference-category-oracle`
- `tests/canary/workflow-guidance-anchors/reference-category-setup`
- `tests/canary/workflow-guidance-anchors/reference-category-work`
- `tests/canary/workflow-guidance-anchors/reference-gate-authority`
- `tests/canary/workflow-guidance-anchors/reference-kit-only-ship`
- `tests/canary/workflow-guidance-anchors/reference-no-path-fallback`
- `tests/canary/workflow-guidance-anchors/reference-progressive-loading-term`
- `tests/canary/workflow-guidance-anchors/reference-refusal-route-shape`
- `tests/canary/workflow-guidance-anchors/reference-retro-capture-owner`
- `tests/canary/workflow-guidance-anchors/reference-retro-drain-owner`
- `tests/canary/workflow-guidance-anchors/reference-skills-guidance`
- `tests/canary/workflow-guidance-anchors/reference-upgrade-route`
- `tests/canary/workflow-guidance-anchors/review-base-merged-main-tip`
- `tests/canary/workflow-guidance-anchors/review-falsification-accept-routing`
- `tests/canary/workflow-guidance-anchors/review-falsification-dispositions`
- `tests/canary/workflow-guidance-anchors/review-kit-guidance-set`
- `tests/canary/workflow-guidance-anchors/review-persistence-anchor`
- `tests/canary/workflow-guidance-anchors/review-preflight-explicit-base`
- `tests/canary/workflow-guidance-anchors/review-repair-ticket-covers`
- `tests/canary/workflow-guidance-anchors/review-repair-ticket-owner`
- `tests/canary/workflow-guidance-anchors/review-standing-falsification`
- `tests/canary/workflow-guidance-anchors/review-universal-claim-bar`

- `bin/bench.sh`
- `cmd/bench/gate_route_test.go`
- `cmd/bench/main_test.go`
- `cmd/bench/main.go`
- `internal/worktree/parallel_census_test.go`

- `tests/canary/docs-currency-token-diet/missing-cli-inventory`
- `tests/canary/docs-currency-token-diet/stale-cli-doc-reference`
- `tests/canary/docs-currency-token-diet/stale-skill-cli-reference`
- `tests/canary/load-validity-metadata/extensionless-gate-ref`
- `tests/canary/package-core-guard/bounds-duplicate-owner`
- `tests/canary/package-core-guard/reintroduced-bare-skip`
- `tests/canary/package-core-guard/unrouted-subcommand`

## Out of scope

This spec excludes semantic proof of judgment, a new approval service, automatic reviewers, and automatic repair. It also excludes probe classification changes and a parallel implementation engine. General preflight enforcement remains under FT200.

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
| #7: source-bound terminal verification and review results | E1, E2, E3, E13, E19, E20, E21, E22 |
| #10: "Missing or stale required evidence blocks chunk advancement and spec completion" | E4, E5, E6, E7, E8, E9, E10, E11, E12, E14, E15, E16 |
| #12: chunk review and final integration reconciliation | E17, E18, E23, E24 |

The retained phase session is the artifact producer. It imports the native reviewer return without converting transport failure into success. This spec adds no record-writing CLI or persistence port. Every required verification item comes from the approved acceptance and testing plan. A review item cannot satisfy that inventory. Completed verification requires outcome pass, exit code zero, current source, and successful required mutation and restore evidence.

| transition | producer and durable boundary | gate disposition |
| --- | --- | --- |
| Review request | The retained phase adds a pending occurrence after freezing source. An ordinary lane can commit it. | Pending is permitted during work and blocks a checkpoint. |
| Terminal review | The phase retains the native return and its performer. It writes completed, failed, or skipped as reported. | Only a completed current result can satisfy an axis. |
| Finding repair | The author retains prior results and appends the repair delta with new results or native reaffirmations. | Uncovered repair makes the requested checkpoint stale. |
| Chunk checkpoint | Commit the review artifact on its ordinary lane, then ask the gate to check that chunk. | Check the continuous source chain and required verification before advancement. |
| Final verification | The author runs the overall acceptance and integration checks and records their native outcomes before the complete request. | Missing, failed, or stale final verification blocks completion. |
| Landing | The existing broker grades the exact prospective subject and performs its defined status transform. | No ref update precedes successful complete evidence. |
| Retirement | The existing retirement path can remove the live pickup after its committed results remain reachable in history. | Historical source references use commit and path. No uncommitted result may be discarded. |

Occurrence states record what ran. Validation states compare that occurrence with the requested source and obligation. An artifact-only commit can follow the frozen review tip only when its delta changes the exact record path alone. Arbitrary spec changes cannot hide in that exception.

Checkpoint enforcement applies to named completion requests in the Bench workflow. It does not claim to prevent arbitrary raw edits. The phase must obtain checkpoint green before starting a successor chunk. Ordinary ticket lane checks keep their existing purpose. The checkpoint path delegates to the gate owner and never writes an independent approval ledger.

Review repair coverage: E25, E26, E27, E28, E29, E30, E31, E32, E33, E34, E35, E36, E37, E38. These rows refine the source clauses already mapped above.

Plan-edit authority is limited to chunk boundaries, dependencies, row assignments, Writes expectations, and necessary ownership expansion under decision #5. Preserve the approved behavior and pass criteria. A material acceptance change still needs the user's decision. Bulk rewrites exclude this spec and its tickets. Completion evidence must review a changed plan digest before it can satisfy the next checkpoint.

Authoring close: this spec and its tickets are staged for user sign-off. The reviews above assess the proposed build. No implementation or implementation test result is claimed.

Reviewer amendment on 2026-09-11: apply decision #13 to the declared implementation model. Sol implementations use Astra/high review axes. This amendment follows the spec-authoring reviews recorded above.

Implementation plan repair: chunk 1 includes `r1.md` for the accepted review findings.
The final acceptance command uses the existing changed-package selector from the pinned landing base.
This repair changes test citations and ownership only; the approved behavior remains the same.

Chunk 1 guard repair replaces the retired clean-review no-artifact requirement with a terminal-result requirement.
The omission fixture retains the gate’s refusal behavior and the existing persistence constraints.

Chunk 2 keeps public gate grammar in the Go gate owner. The shell forwards its arguments unchanged.
The help inventory derives its suffix from that owner, and wrapper tests execute that same owner.

Chunk 2 repair keeps verification cases in the shared checkpoint fixture.
It adds partial-inventory, completion-purpose reuse, and hostile-path wrapper coverage before advancement.
