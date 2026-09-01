# Codify the load stop and quiet check

Blocked by: support-installed-lane-repair-commit.md
Writes: .agents/skills/bench-craft-delegate/SKILL.md, .agents/skills/bench-craft-line/SKILL.md, .agents/skills/bench-craft-delegate/references/delegation-discipline.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, tests/canary/workflow-guidance-anchors/delegate-charge-effort-cap, tests/canary/workflow-guidance-anchors/delegate-coverage-row-charge, tests/canary/workflow-guidance-anchors/delegate-coverage-row-red-green, tests/canary/workflow-guidance-anchors/delegate-cross-harness-reviewer-pointer, tests/canary/workflow-guidance-anchors/delegate-model-id-escalation, tests/canary/workflow-guidance-anchors/delegate-own-family-native-surface, tests/canary/workflow-guidance-anchors/delegate-parallel-route-anchor, tests/canary/workflow-guidance-anchors/delegate-release-at-acceptance, tests/canary/workflow-guidance-anchors/delegate-resume-handoff-contents, tests/canary/workflow-guidance-anchors/delegate-self-probe-missing-row, tests/canary/workflow-guidance-anchors/delegate-stash-refusal-anchor, tests/canary/workflow-guidance-anchors/fix-pass-sentinel-anchor, tests/canary/workflow-guidance-anchors/shared-worktree-path-pin, tests/canary/workflow-guidance-anchors/ticket-stage-routing-anchor, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: LF4

## What to build

Require a stop after two known-flaky refusals proven green in isolation.
Before aggregate grading, require returned delegates to have no live tests and
serialize the coordinator-owned resource.

## Acceptance

- [ ] The second proven flaky refusal stops and hands evidence to the reviewer.
- [ ] Aggregate grading waits until returned delegates own no live tests.
- [ ] Canonical anchors protect both rules without copying their full prose.

