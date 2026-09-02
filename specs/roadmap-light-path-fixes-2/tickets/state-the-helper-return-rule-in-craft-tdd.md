# State the helper return rule in craft-tdd

Blocked by: bind-the-exec-only-form-to-every-caller-in-craft-delegate.md
Writes: .agents/skills/bench-craft-tdd/SKILL.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, tests/canary/workflow-guidance-anchors/craft-tdd-helper-returns-not-skips (new), projects/benchkit.md, tests/canary/workflow-guidance-anchors/craft-tdd-edge-class-run, tests/canary/workflow-guidance-anchors/craft-tdd-light-path-seam-gate, tests/canary/workflow-guidance-anchors/craft-tdd-red-signal-classification-owner, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go, tests/canary/guidance-prose-budgets/over-budget-skill, tests/canary/line-routing/line-binding-prose-drift, tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading, tests/canary/workflow-guidance-anchors/benchkit-review-round-owner, tests/canary/workflow-guidance-anchors/benchkit-review-round-routing, tests/canary/workflow-guidance-anchors/benchkit-spec-ownership
Covers: LP18

## What to build

Ship one sentence in the five-part precedent shape, on the integration source
after the blocker ticket lands. Place it under "The oracle is the gate, not
you", beside the sentence about a skipped matched test. A re-exec helper
returns silently outside its role environment and never skips, because the kit
gate treats an environment-class skip as red. Verify the premise against
`internal/gate/capability_skips.go` and one helper such as
`internal/gocache/lock_test.go`, and cite both in the report. Keep `SKILL.md`
inside its 122-line budget row.

## Acceptance

- [ ] The registry test reds a synthetic tree that drops the sentence and stays silent on the live root.
- [ ] The fixture bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] The prose budget check passes on the worktree.
- [ ] Self-probe: change "never skips" to "may skip" in the live file, and report the registry test red.
