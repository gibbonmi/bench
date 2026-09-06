# Write the craft-research skill and register it

Blocked by: none
Writes: .agents/skills/bench-craft-research/SKILL.md (new), .claude/skills/bench-craft-research (new), .bench/BENCH-reference.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, tests/canary/workflow-guidance-anchors/, tests/canary/docs-currency-token-diet/benchref-imported, tests/canary/docs-currency-token-diet/benchref-pointer-dropped, tests/canary/docs-currency-token-diet/benchref-section-duplicated, tests/canary/skills-index-command-adapters/dangling-index, tests/canary/skills-index-command-adapters/missing-index-field, tests/canary/skills-index-command-adapters/stale-index-wording, tests/canary/skills-index-command-adapters/unindexed-skill, tests/canary/workflow-guidance-anchors/reference-bench-operational-layer, tests/canary/workflow-guidance-anchors/reference-category-context, tests/canary/workflow-guidance-anchors/reference-category-oracle, tests/canary/workflow-guidance-anchors/reference-category-setup, tests/canary/workflow-guidance-anchors/reference-category-work, tests/canary/workflow-guidance-anchors/reference-gate-authority, tests/canary/workflow-guidance-anchors/reference-kit-only-ship, tests/canary/workflow-guidance-anchors/reference-no-path-fallback, tests/canary/workflow-guidance-anchors/reference-progressive-loading-term, tests/canary/workflow-guidance-anchors/reference-refusal-route-shape, tests/canary/workflow-guidance-anchors/reference-retro-capture-owner, tests/canary/workflow-guidance-anchors/reference-retro-drain-owner, tests/canary/workflow-guidance-anchors/reference-skills-guidance, tests/canary/workflow-guidance-anchors/reference-upgrade-route, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/skills-index-command-adapters/adapter-inert-invocation-key, tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy, tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted, tests/canary/workflow-guidance-anchors/reference-agent-push-rule
Covers: CR1, CR2, CR3, CR4, CR5, CR6, CR7, CR8, CR9, CR10, CR11, CR12, CR13, CR18, CR19, CR20

## What to build

Write `.agents/skills/bench-craft-research/SKILL.md` in ASD-STE100 within 120 lines. The frontmatter carries `name: craft-research`, a description that names the trigger, and an `index:` line.

The body owns the trigger, the question graph, the primary-source standard, and adaptive round-based fan-out. It owns coordinator synthesis and verification, the destination precedence, and the metadata labels. It owns the report contract with its checklist, the compatibility-probe rule, and the read-side boundary. It points at the fan-out clause and the delegate discipline instead of restating them. It carries one contrastive pair: an independent fan-out against a dependent serial question.

Create the relative symlink `.claude/skills/bench-craft-research` as its siblings are made. Run `bench skills-index --write`. Register one needle per sentence listed under the spec's Further notes, add a `TestCraftResearchAnchorsRedOnRemoval` mutation table, and add the fixture under the workflow-guidance-anchors family.

Before the return, read the draft for its silences and list each omission you filled and each you left open. Quote the two pointers and the contrastive pair in the return.

## Acceptance

- [ ] [S1] `bench test --check load-validity-metadata` and `bench test --check skills-index-command-adapters` pass.
- [ ] [S2] `bench gate-prose` passes on the skill, and the `guidance-prose-budgets` check passes at or under 120 lines.
- [ ] [S3] Every listed needle is present, and `TestCraftResearchAnchorsRedOnRemoval` reds on each removal.
- [ ] [S4] The skill states no tier or effort and names no Codex adapter.
- [ ] [S5] The report section lists every contract element and the checklist's landing place.
