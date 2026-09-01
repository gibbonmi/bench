# Verify done-claim owners

Blocked by: formalize-repair-charge-template.md
Writes: .agents/skills/bench-craft-delegate/SKILL.md, .agents/skills/bench-craft-delegate/references/delegation-discipline.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, tests/canary/workflow-guidance-anchors/delegate-charge-effort-cap, tests/canary/workflow-guidance-anchors/delegate-coverage-row-charge, tests/canary/workflow-guidance-anchors/delegate-coverage-row-red-green, tests/canary/workflow-guidance-anchors/delegate-cross-harness-reviewer-pointer, tests/canary/workflow-guidance-anchors/delegate-model-id-escalation, tests/canary/workflow-guidance-anchors/delegate-own-family-native-surface, tests/canary/workflow-guidance-anchors/delegate-parallel-route-anchor, tests/canary/workflow-guidance-anchors/delegate-release-at-acceptance, tests/canary/workflow-guidance-anchors/delegate-resume-handoff-contents, tests/canary/workflow-guidance-anchors/delegate-self-probe-missing-row, tests/canary/workflow-guidance-anchors/delegate-stash-refusal-anchor, tests/canary/workflow-guidance-anchors/fix-pass-sentinel-anchor, tests/canary/workflow-guidance-anchors/shared-worktree-path-pin, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: LF10

## What to build

Require every named Red-mutation owner to resolve to a tree artifact. Keep an
accepted finding on its original ticket unless one shared owner justifies an
umbrella repair ticket.

## Acceptance

- [ ] An absent Red-mutation owner cannot support a verified done claim.
- [ ] Attributable findings remain on their original ticket.
- [ ] Umbrella repair tickets require one genuinely shared owner.

