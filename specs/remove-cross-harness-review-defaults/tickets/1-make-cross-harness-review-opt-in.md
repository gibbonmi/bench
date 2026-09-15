# Make cross-harness review opt-in

Blocked by: none
Writes: .agents/commands/bench-implement-spec.md, .agents/commands/bench-review-implementation.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, tests/canary/workflow-guidance-anchors/review-cross-harness-opt-in (new), tests/canary/workflow-guidance-anchors/review-standing-falsification, tests/canary/workflow-guidance-anchors/implement-spec-cross-harness-pointer (new), tests/canary/workflow-guidance-anchors/implement-spec-offer-scope, tests/canary/workflow-guidance-anchors/review-kit-guidance-set, CHANGELOG.md, specs/remove-cross-harness-review-defaults/tickets/1-make-cross-harness-review-opt-in.md (new)
Covers: none

## What to build

Implementation review runs its native Standards, Spec, and Coverage axes by
default. It adds a cross-harness falsification pass only when the reviewer
requests one.

## Acceptance

- [x] Review guidance makes the cross-harness falsification pass opt-in for every diff.
- [x] Full implementation guidance points to the review phase for cross-harness opt-in.
- [x] Conformance mutations detect a missing opt-in rule and a restored standing pass.
