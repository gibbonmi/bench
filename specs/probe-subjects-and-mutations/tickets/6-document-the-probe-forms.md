# Document the probe forms and update the known issues

Blocked by: 1-add-probe-mutation-forms.md, 2-probe-prose-and-system-checks.md, 3-name-the-first-failure-and-expected-witness.md, 4-name-a-mutant-build-failure.md, 5-lock-the-probe-subject.md
Writes: .agents/skills/bench-craft-delegate/SKILL.md, .bench/BENCH-reference.md, CHANGELOG.md, CONTEXT.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, docs/known-issues.md, internal/anchors/registry_calibration.go, internal/anchors/registry_calibration_test.go, internal/anchors/registry_charge_binding.go, internal/anchors/registry_charge_binding_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_decision_maps.go, internal/anchors/registry_decision_maps_test.go, internal/anchors/registry_ft311_preparation.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, roadmap/FT168.md, tests/canary/claude-agent-definitions/agent-unnamed-in-skill, tests/canary/claude-agent-definitions/model-declared, tests/canary/claude-agent-definitions/name-mismatch, tests/canary/claude-agent-definitions/shell-tool-absent, tests/canary/claude-agent-definitions/skill-names-missing-agent, tests/canary/claude-agent-definitions/spawning-tool, tests/canary/claude-agent-definitions/tools-absent, tests/canary/docs-currency-token-diet/benchref-imported, tests/canary/docs-currency-token-diet/benchref-pointer-dropped, tests/canary/docs-currency-token-diet/benchref-section-duplicated, tests/canary/docs-currency-token-diet/signal-vocabulary-drift, tests/canary/skills-index-command-adapters/adapter-inert-invocation-key, tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy, tests/canary/skills-index-command-adapters/dangling-index, tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted, tests/canary/skills-index-command-adapters/missing-index-field, tests/canary/skills-index-command-adapters/stale-index-wording, tests/canary/skills-index-command-adapters/unindexed-skill, tests/canary/workflow-guidance-anchors/agents-handoff-section-rule, tests/canary/workflow-guidance-anchors/calibration-delegate-pointer, tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns, tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary, tests/canary/workflow-guidance-anchors/context-acceptance-row-vocabulary, tests/canary/workflow-guidance-anchors/context-coverage-map-term, tests/canary/workflow-guidance-anchors/context-coverage-row-parts, tests/canary/workflow-guidance-anchors/context-coverage-row-vocabulary, tests/canary/workflow-guidance-anchors/context-decision-map-term, tests/canary/workflow-guidance-anchors/context-reader-sweep-term, tests/canary/workflow-guidance-anchors/context-ticket-vocabulary, tests/canary/workflow-guidance-anchors/delegate-cap-change-pinning-package, tests/canary/workflow-guidance-anchors/delegate-charge-effort-cap, tests/canary/workflow-guidance-anchors/delegate-coverage-row-charge, tests/canary/workflow-guidance-anchors/delegate-coverage-row-red-green, tests/canary/workflow-guidance-anchors/delegate-cross-harness-reviewer-pointer, tests/canary/workflow-guidance-anchors/delegate-exec-only-every-caller, tests/canary/workflow-guidance-anchors/delegate-model-id-escalation, tests/canary/workflow-guidance-anchors/delegate-own-family-native-surface, tests/canary/workflow-guidance-anchors/delegate-parallel-route-anchor, tests/canary/workflow-guidance-anchors/delegate-release-at-acceptance, tests/canary/workflow-guidance-anchors/delegate-resume-handoff-contents, tests/canary/workflow-guidance-anchors/delegate-self-probe-missing-row, tests/canary/workflow-guidance-anchors/delegate-stash-refusal-anchor, tests/canary/workflow-guidance-anchors/delegated-per-ticket-author, tests/canary/workflow-guidance-anchors/dg-6-owner-reference, tests/canary/workflow-guidance-anchors/fix-pass-sentinel-anchor, tests/canary/workflow-guidance-anchors/prepared-review-delegate-handoff-route, tests/canary/workflow-guidance-anchors/prepared-triage-bounds, tests/canary/workflow-guidance-anchors/reference-agent-push-rule, tests/canary/workflow-guidance-anchors/reference-bench-operational-layer, tests/canary/workflow-guidance-anchors/reference-category-context, tests/canary/workflow-guidance-anchors/reference-category-oracle, tests/canary/workflow-guidance-anchors/reference-category-setup, tests/canary/workflow-guidance-anchors/reference-category-work, tests/canary/workflow-guidance-anchors/reference-gate-authority, tests/canary/workflow-guidance-anchors/reference-kit-only-ship, tests/canary/workflow-guidance-anchors/reference-no-path-fallback, tests/canary/workflow-guidance-anchors/reference-progressive-loading-term, tests/canary/workflow-guidance-anchors/reference-refusal-route-shape, tests/canary/workflow-guidance-anchors/reference-retro-capture-owner, tests/canary/workflow-guidance-anchors/reference-retro-drain-owner, tests/canary/workflow-guidance-anchors/reference-skills-guidance, tests/canary/workflow-guidance-anchors/reference-upgrade-route, tests/canary/workflow-guidance-anchors/shared-worktree-path-pin
Covers: PS69, PS70, PS71, PS72, PS73, PS74

## What to build

A reader finds every new probe form in the reference guide, the glossary, the delegate skill, and the changelog.
The known-issues table lists only open issues.

Tickets 1 to 5 supply every form, subject, verdict face, and the subject lock.
This ticket states that delivered behavior, so it starts after all of them commit.

The reference guide's probe paragraphs name each new form, the subject lock, and the `mutant` table.
`CONTEXT.md` defines **combined mutant** and **expected witness**, each with an Avoid list.
The delegate skill's probe sentence at line 104 names the new forms.
The anchored needle at line 81 of the delegate skill stays byte-identical, because a self-probe still reports one omission or one swap.
`CHANGELOG.md` names each new form and the schema v2 change.

`docs/known-issues.md` loses the unused-import row.
The probe-cost row's retire condition changes to the batch obligation.
`roadmap/FT168.md` keeps that batch obligation in its body.

Each guidance change keeps the current anchor placement.
The anchor registry files, the canary units, and the command registry files join this `Writes:` line as closure.
Change a closure file only when an anchored sentence moves, and keep each canary unit's meaning intact.

## Acceptance

- [ ] The reference guide names each new form, the subject lock, and the `mutant` table.
- [ ] `CONTEXT.md` defines **combined mutant** and **expected witness** with Avoid lists.
- [ ] The delegate skill's probe sentence names the new forms, and the line 81 needle stays byte-identical.
- [ ] `CHANGELOG.md` names each new form and the schema v2 change.
- [ ] `docs/known-issues.md` holds no row that names an unused probe import.
- [ ] The probe-cost row names the batch obligation as its retire condition, and `roadmap/FT168.md` keeps that obligation.
- [ ] Each anchor needle and each registry assertion over the edited guidance stays green.
