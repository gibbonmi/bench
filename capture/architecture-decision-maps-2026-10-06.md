# Architecture review: decision maps

The maps cover all eleven deepening candidates from the 2026-10-06 architecture review.
They link the reviewed specifications and preserve earlier approved decisions.
They authorize no implementation.

The reviewer approved the FT362 owner, migration scope, and failure posture.
The reviewer then delegated the remaining architecture recommendations.
Each decision ticket holds its answer; this index provides navigation and sequence.

Source commit: `37ac80b8adeff4ecbd57f3d813ff8937f2bcbc68`

Repository: `/home/mgibs/workspace/bench`

Assignment: `ft362-process-lifetime-map`

The [full-roadmap assessment](bench-maintenance-assessment-2026-10-07.md) gives all 99 rows a proposed disposition.
This map set covers its eleven architecture candidates, not every feature in that ranking.
The active FT290 commitment and all existing roadmap dependencies remain unchanged.

## Candidate coverage

| Review candidate | Map | Existing owner | Next planning action |
|---|---|---|---|
| C01: authored Markdown | [Markdown](../decisions/architecture-markdown.md) | FT358; markdown-block-reader | Accepted spec and ticket amendment |
| C02: process lifetime | [Process lifetime](../specs/process-lifetime/decisions/ft362-process-lifetime.md) | FT362 | Accepted spec and thirteen tickets; integrated |
| C03: source observations | [Observations](../decisions/architecture-observations.md) | landing-test-efficiency; residual FT365 | Three reviewed plans are integrated |
| C04: generic Git fixtures | [Fixtures](../decisions/architecture-fixtures.md) | FT360; shared-test-fixtures | Accepted spec and ticket amendment |
| C05: cleanup traversal | [Cleanup](../specs/cleanup-member-traversal/decisions/architecture-cleanup.md) | FT363 | Accepted spec and four implementation tickets |
| C06: lifecycle decision consumption | [Adoption](../specs/adoption-lifecycle-consumption/decisions/architecture-adoption.md) | FT217 | Accepted spec and four implementation tickets |
| C07: shared guard grammar | [Guards](../decisions/architecture-guards.md) | FT366 | Accepted projection and scan plans; integrated |
| C08: shared primitives | [Primitives](../decisions/architecture-primitives.md) | FT354 | Four reviewed plans are integrated |
| C09: production test seams | [Test seams](../specs/production-test-seam-policy/decisions/architecture-test-seams.md) | FT343 | Accepted spec and four tickets; integrated |
| C10: planning operations | [Planning](../decisions/architecture-planning.md) | FT293, FT375, FT318, FT125 | Four reviewed plans are integrated |
| C11: release cuts | [Release](../specs/release-evidence-forwards/decisions/architecture-release.md) | FT364; partial FT368 | Accepted forwarding plan; integrated |

## Recommended sequence

Start with the narrow false-success and cancellation repairs in C08, C07, and C02.
The completed process prototype establishes the stream strategy for its specification.
The retained-resource design still needs production proof during implementation.
Keep the separate shell exit-status repair outside FT362.

Use the approved FT358 spec next, followed by the approved landing-test-efficiency spec.
Coordinate their successor fixture and source-observation migrations with FT360 and residual FT365.
Do not create duplicate implementation graphs for these approved outcomes.

Then take C05 cleanup traversal and C06 adoption decision consumption.
Follow with C09 test-seam enforcement and C10 planning operations.
Keep C11's mechanical release cuts behind its consumer and release-evidence checks.

This sequence proposes work within the architecture candidates.
It does not replace the complete roadmap ranking or authorize a commitment change.
The other confirmed defects in that ranking keep their earlier position.

## Evidence and limits

The coordinator reread the current source behind the maps.
Source inspection established caller differences and stale roadmap premises.
It did not establish runtime compatibility, performance gains, or implementation acceptance.

During the original map pass, a requested `gpt-6.1-sol` delegate at high effort could not start.
The native agent tool returned `agent thread limit reached`.
The coordinator authored that map set; no fresh delegate review is claimed for that pass.
The earlier architecture survey retains its separate Sol 6.1/xhigh assessments.

The document checks establish the following results:

- The process map records its completed prototype and the accepted retention constraints.
- The prose check passes every map and decision ticket.
- Every map source path and ticket link resolves.
- Every cited source line in the decision tickets exists.
- The set contains one map for each architecture candidate.

The FT362 census also includes the bounds process runner.
That runner was absent from the earlier candidate summary.
The census ticket is the single source for the current migration sites.

No benchmark, production test suite, or implementation ran during this map pass.
The named assignment retains these planning artifacts; its Git history records the committed checkpoint.
No map readiness claim grants build approval or live release qualification.

## Planning status

All twenty-one outcomes have accepted specifications and ticket graphs, including reuse of landing-test-efficiency.
All are committed and integrated into this planning branch.
The preflight applicability amendment and the staleness inventory correction have independent acceptance.
Each new plan completed its spec-stage close with its review record and scorecard observation.

The table links every accepted planning outcome.
The durable-file prerequisite remains a separate deliverable within this set of twenty-one outcomes.
Acceptance here describes planning review and grants no implementation authority.
The primary handoff pins reviewed commits, close commits, and workspaces.

| Outcome | Plan | Actual state |
|---|---|---|
| Authored Markdown | [markdown-block-reader](../specs/markdown-block-reader/spec.md) | Spec and amended tickets accepted; integrated |
| Process lifetime | [process-lifetime](../specs/process-lifetime/spec.md) | Spec and thirteen tickets accepted; integrated |
| Initial source observations | [landing-test-efficiency](../specs/landing-test-efficiency/spec.md) | Approved spec and nine-ticket graph accepted for reuse |
| Residual syntax visitors | [conformance-observation-follow-on](../specs/conformance-observation-follow-on/spec.md) | Spec and three tickets accepted; integrated |
| Independent expectations | [independent-expectation-audit](../specs/independent-expectation-audit/spec.md) | Spec and three tickets accepted; integrated |
| Generic Git fixtures | [shared-test-fixtures](../specs/shared-test-fixtures/spec.md) | Spec and amended tickets accepted; integrated |
| Cleanup traversal | [cleanup-member-traversal](../specs/cleanup-member-traversal/spec.md) | Spec and four tickets accepted; integrated |
| Adoption lifecycle | [adoption-lifecycle-consumption](../specs/adoption-lifecycle-consumption/spec.md) | Spec and four tickets accepted; integrated |
| Malformed guard refusal | [degraded-guard-refusal](../specs/degraded-guard-refusal/spec.md) | Spec and one ticket accepted; integrated |
| Shared guard grammar | [shared-guard-projection](../specs/shared-guard-projection/spec.md) | Spec and three tickets accepted; integrated |
| Guard inventory scan | [guard-scan-loop](../specs/guard-scan-loop/spec.md) | Spec and one ticket accepted; integrated |
| Strict JSON reads | [strict-json-ownership](../specs/strict-json-ownership/spec.md) | Spec and five tickets accepted; integrated |
| Durable file replacement prerequisite | [durable-file-replacement](../specs/durable-file-replacement/spec.md) | Spec and one ticket accepted; integrated |
| Durable caller migrations | [durable-caller-migration](../specs/durable-caller-migration/spec.md) | Spec and thirteen tickets accepted; integrated |
| Printed shell arguments | [shell-argument-ownership](../specs/shell-argument-ownership/spec.md) | Spec and seven tickets accepted after repair; integrated |
| Test-seam admission | [production-test-seam-policy](../specs/production-test-seam-policy/spec.md) | Spec and four tickets accepted; integrated |
| Fence and premise closure | [preflight-ownership-closure](../specs/preflight-ownership-closure/spec.md) | Spec and eight tickets accepted; integrated |
| Mechanical spec staleness | [preflight-spec-staleness](../specs/preflight-spec-staleness/spec.md) | Spec and five tickets accepted after inventory repair; integrated |
| Native record operations | [native-record-operations](../specs/native-record-operations/spec.md) | Spec and three tickets accepted; integrated |
| Exact bounded readers | [exact-planning-readers](../specs/exact-planning-readers/spec.md) | Spec and five tickets accepted; integrated |
| Release forwards | [release-evidence-forwards](../specs/release-evidence-forwards/spec.md) | Spec and one ticket accepted; integrated |

The separate FT366 rg-operand residual remains outside the three guard outcomes.
Partial specs retain their source row's other obligations.
The user selected Sol 6.1/high authors and approved retained Sol 6.1/xhigh independent reviewers.
Each new graph was sliced only after its spec passed independent review.
The preflight amendment defines charge readiness while preserving manual diagnostics and all future native acceptance evidence.
Document checks and static reviews supply no runtime or performance acceptance.

## New maintenance observation

The planning integration supplied a source-verified observation of FT335's parked full-gate-selection defect.
A primary-checkout merge selected all six gate phases for a Markdown-only delta.
The next authorized target-local merge selected its prose lane.
The two deltas differ, so this is not a controlled timing comparison.

Read the [FT335 observation](ft335-maintenance-merge-observation-2026-10-07.md) for source identities and retained evidence.
Reassess its parked disposition before optional architecture refactors in the fresh ranking.
This recommendation changes neither ROADMAP.md nor the active FT290 commitment.

## Continuation

Next command: `$bench-mainenance`

The authorized planning batch is complete.
Read the primary checkout's `capture/session-handoff.md` for source pins, accepted reviews, and preserved decisions.
A later maintenance invocation can reuse these accepted plans and the full-roadmap assessment.
No implementation or landing of this planning branch was performed.
