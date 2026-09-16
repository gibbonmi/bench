# Enable the unified-review trial

Blocked by: 1-retain-debug-authorship.md
Writes: .agents/commands/bench-review-implementation.md, .agents/commands/bench-implement-spec.md, internal/reviewrecord/delegated.go, internal/reviewrecord/plan.go, internal/reviewrecord/record.go, internal/reviewrecord/coverage.go, internal/reviewrecord/delegated_test.go, internal/gate/delegated_checkpoint_test.go, tests/canary/workflow-guidance-anchors, reviews/debug-loop-guidance.md, CHANGELOG.md
Covers: DG44, DG45

## What to build

Implement DG-CR as an opt-in review topology for this delegated run.
Keep the existing three-distinct-reviewer behavior as the default.
Add a source-bound completion-plan field that selects one independent reviewer for all three axes.
Reject unknown values and continue to exclude the orchestrator and every ticket author.
Require the unified reviewer to report Standards, Spec, and Coverage separately.
Require that reviewer to identify whether each issue exposes an improvement to the implementation-command prose.

Use TDD at the checkpoint seam.
First prove that the existing default still refuses one performer for multiple axes.
Then prove that the explicit unified mode accepts one independent performer for all three completed axes.
Prove that an author or orchestrator remains ineligible in either mode.

This ticket enables only the review topology.
Do not add comparative-trial controls, a metrics schema, pricing logic, or a global default change.
The coordinator retains the trial measurements in the existing review artifact.

## Acceptance

- [ ] Omitted review mode still requires three distinct independent review sessions.
- [ ] Explicit unified review mode accepts one independent session for Standards, Spec, and Coverage.
- [ ] The unified reviewer cannot be the orchestrator or any current or former ticket author.
- [ ] Unknown review modes fail closed.
- [ ] Unified review reports keep the three axes separately attributable.
- [ ] Unified review calls out any implement-spec prose improvement exposed by a finding or miss.
- [ ] The checkpoint tests demonstrate both the default refusal and the opt-in acceptance.
- [ ] This outcome verifies while the phase-guidance successor tickets remain unbuilt.

## Focused checks

- `bench test --package ./internal/reviewrecord --run 'TestDelegated.*Review'`
- `bench test --package ./internal/gate --run TestDelegatedDistinctAxes`
- `bench test --check prose-mechanics`
- `bench test --check ticket-grammar`
