# Bounded repair policy enforcement evidence

## Observed source

Revision: 80e438a0a309f26724e0fe6f44f246950c8dd52a
Read date: 2026-09-13
Drift: The named policy consumer, conformance family, or review-record validator changes.

## Policy seam and precedent

The implementation phase charges the line skill's continuation policy throughout the ticket graph.
Its post-review path retains author verification and current review results.
Source: .agents/commands/bench-implement-spec.md:36 and .agents/commands/bench-implement-spec.md:48.

The line skill owns progress, attempt identity, numeric caps, and the pre-review no-progress reassessment.
It currently has 128 physical lines, with its guidance budget held in the project profile.
Source: .agents/skills/bench-craft-line/SKILL.md:98 and projects/benchkit.md, Guidance prose budgets.

The review phase already restricts repeat review to invalidated evidence.
It separates repair routing from authority to make an edit.
Source: .agents/commands/bench-review-implementation.md:28 and .agents/commands/bench-review-implementation.md:140.

The existing implementation-continuation family protects these instructions.
Its test independently declares expected predicates and drives omission cases through the shared anchor harness.
Sources: internal/anchors/registry_retained_workflow.go:53 and internal/conformance/implementation_continuation_test.go:11.

## Executed conformance root

The executable check registry binds `docs-currency-workflow` to `checkDocsCurrencyAndWorkflow`.
That function calls `checkWorkflowAnchors`, which evaluates `AfterImplementSpec`.
The anchor registry already includes `implementationContinuationAnchors` in that group.

- internal/conformance/checks_test.go:55
- internal/conformance/docs_workflow_checks_test.go:20
- internal/conformance/docs_workflow_helpers_test.go:30
- internal/anchors/registry_data.go:26
- internal/anchors/registry.go:55

`runAnchorBites` writes each conformant fixture, substitutes a contradictory value, and requires that anchor's diagnostic.
It is shared fixture infrastructure, not a new fixture owner for this spec.
Source: internal/conformance/docs_workflow_helpers_test.go:452.

The current review-convergence checker pins required review coverage and successor ordering.
The change can preserve its exact requirements while clarifying that optional advice is separate from blocking findings.
Sources: internal/conformance/docs_workflow_helpers_test.go:86 and internal/conformance/docs_workflow_helpers_test.go:123.

## Named reader sweep

| Reader | Fact consumed | Treatment |
| --- | --- | --- |
| .agents/skills/bench-craft-line/SKILL.md | Attempt, progress, and cap policy | Edit the consumer pointer and qualify the post-review exception |
| .agents/commands/bench-implement-spec.md | Repair continuation and short-stop routing | Edit the policy reference and bounded handoff route |
| .agents/commands/bench-review-implementation.md | Finding classification, repeat review, and repair evidence | Edit classification and reference the one allowance owner |
| .agents/skills/bench-craft-review/SKILL.md | Axis findings and their cited basis | Add the policy reference without copying its limit |
| .bench/BENCH.md | Ownership of the continuation policy | Preserve the existing craft-line pointer and review guarantees |
| projects/benchkit.md | Model binding and prose budgets | Preserve the current bindings and budgets |
| internal/anchors/registry_retained_workflow.go | Required continuation instructions | Extend the existing family |
| internal/conformance/implementation_continuation_test.go | Independent required-predicate expectations | Extend the existing test and demonstrate anchor omission red |
| internal/anchors/registry_data.go | Registry membership | Existing family membership already reaches new entries |
| internal/conformance/docs_workflow_helpers_test.go | Review convergence and anchor execution | Preserve the checker and reuse its existing fixture harness |
| internal/conformance/retained_workflow_test.go | Existing authorship and review guarantees | Preserve the pinned workflow predicates |
| internal/anchors/registry_ft311_review_dispatch.go | Existing repeat-review and evidence requirements | Preserve the pinned workflow predicates |

Each edited reader has an exact spec fence.
The deeper anchor evaluator, file classifier, and conformance dispatcher retain their existing behavior.
No script or workflow consumer of a repair count exists in the current tree.
The count is new workflow state in existing prose artifacts, not a new machine-readable schema.

## Native review-record readers

The JSON parser requires a passing completed review to have no finding IDs.
The completion validator also requires an independent current review for each axis.

- internal/reviewrecord/parse.go:84
- internal/reviewrecord/record.go:79
- internal/reviewrecord/check.go:27
- internal/reviewrecord/check.go:58

The policy therefore records optional advice outside unresolved blocking finding IDs.
These validators stay unchanged, and no optional classification can waive failed required verification.

`TestReviewCheckpointFindingAndReviewIdentity` covers unresolved findings, stale source identity, and a stale reviewed pair.
`TestReviewRecordPaths` retains the record's hostile-file behavior.
The historical supersession test retains old findings after a later passing result.

- internal/gate/review_checkpoint_test.go:224
- internal/reviewrecord/source_test.go:126
- internal/reviewrecord/source_test.go:107

These files are read-only regression consumers, not implementation fences.

## Source movement

The ready map index and its complete topic folder move into this spec's decisions directory.
The asset retains the trace snapshot digest and the evidence limits.
The move updates the index's asset path, research ticket #6's report path, and the report's consuming-ticket path.
The session handoff receives the compiled path at phase close.
No top-level map copy remains.
The collection specification can use the compiled source as its named reviewed artifact.

## Scope and verification limits

The gate can detect missing instructions, consumer pointers, and independent expectations.
It cannot prove agent compliance or classify arbitrary review prose.
Those properties remain review-owned, with explicit scenarios in the acceptance map.
The future implementation must demonstrate one omitted policy instruction and one omitted production anchor as distinct reds.
