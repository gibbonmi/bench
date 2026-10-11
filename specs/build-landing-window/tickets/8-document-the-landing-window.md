# Document the landing window, the fold record, and the override

Blocked by: 7-grant-a-reviewer-override.md
Writes: .agents/commands/bench-drain.md, .agents/commands/bench-final-check.md, .agents/commands/bench-implement-spec.md, .agents/commands/bench-review-implementation.md, .agents/commands/bench-write-spec.md, .agents/skills/bench-craft-delegate/SKILL.md, .bench/BENCH-reference.md, CHANGELOG.md, CONTEXT.md, README.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, docs/adr/0014-main-receives-writes-only-through-landings.md, docs/known-issues.md, internal/anchors/registry_calibration.go, internal/anchors/registry_calibration_test.go, internal/anchors/registry_charge_binding.go, internal/anchors/registry_charge_binding_test.go, internal/anchors/registry_check_floor_test.go, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_commitment.go, internal/anchors/registry_craft_research.go, internal/anchors/registry_craft_research_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_decision_maps.go, internal/anchors/registry_decision_maps_test.go, internal/anchors/registry_ft311_preparation.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, projects/benchkit.md, tests/canary/claude-agent-definitions/agent-unnamed-in-skill, tests/canary/claude-agent-definitions/model-declared, tests/canary/claude-agent-definitions/name-mismatch, tests/canary/claude-agent-definitions/shell-tool-absent, tests/canary/claude-agent-definitions/skill-names-missing-agent, tests/canary/claude-agent-definitions/spawning-tool, tests/canary/claude-agent-definitions/tools-absent, tests/canary/docs-currency-token-diet, tests/canary/guidance-prose-budgets, tests/canary/line-routing, tests/canary/load-validity-metadata/readme-shared-rule-drift, tests/canary/row-next-grammar/token-table-lacks-kit-edit, tests/canary/skill-description-budgets, tests/canary/skills-index-command-adapters, tests/canary/workflow-guidance-anchors
Covers: BW85, BW86, BW87, BW88, BW89, BW90, BW91, BW92

## What to build

The guidance states the enforced route for a moved `main`.
A session then follows the lease, the window, the gap fold, the fold record, and the override.

Rewrite the fold rule in the `Land` section of `.agents/commands/bench-implement-spec.md`.
A build folds `main` before the first chunk and in a gap before the next chunk's first dispatch.
After the last freeze, `main` moves only through the override.
Each fold after the first chunk takes `bench record fold`.

Retire the sentence that folds `main` "at the start of the last chunk".
Also retire the sentence that puts that fold before the last chunk's first dispatch.
Give each retired sentence its own forbid row in `internal/anchors/registry_chunk_chain.go`, graded by `TestChunkChainAnchors`.
Move the anchors at `registry_chunk_chain.go:30-36` to the new sentences, and update that file's header comment.

`.agents/commands/bench-implement-spec.md` is 93 lines against its 93-line budget at `projects/benchkit.md:497`.
Keep the rewrite inside that budget, or raise the budget in the same change with its reason.

Update `.agents/commands/bench-review-implementation.md:99`, `:217`, and `:244` for the fold record, the resolution review, and the landing base after a fold.
Retire the clause at `:99` that says the review chain refuses a `main` merge between two chunks.
Give that clause its own forbid row, and move the anchor at `registry_chunk_chain.go:56`.
`.agents/commands/bench-final-check.md:24` names the moved-destination route and no `bench doctor` route for a moved destination.

Add Command Notes to `.bench/BENCH-reference.md` for the lease, the window, the gap, and the override.
They also name `bench commitment block` as the reviewer's stop for a stalled build that still holds the window.
They state the back-out route: `git revert --no-commit <commit>` in a fresh worktree, `bench commit`, and a landing with the override.
`main` is never reset.
A new heading there must not repeat a `.bench/BENCH.md` heading.

Define **landing lease**, **open build**, **chunk gap**, **fold record**, and **landing override** in `CONTEXT.md`, each with an Avoid list.
Update the follow-on guard entry at `CONTEXT.md:190-196`, because ticket 7 makes that guard deny the override verb.
Remove the `base-current` row at `docs/known-issues.md:14`.
Add one `CHANGELOG.md` entry that names the lease, the window, the override, the fold record, and the moved-destination face.

Keep these other guidance readers unchanged, unless a sentence there now states a false route:

- the delegate skill at `.agents/skills/bench-craft-delegate/SKILL.md:68`
- `.agents/commands/bench-write-spec.md:17`
- `.agents/commands/bench-drain.md`
- `README.md`
- ADR 0014

The delegate skill is also at its 126-line budget.
The stale-command sweep reads the new verb forms from the help inventory, which tickets 3, 4, and 7 supply.

## Acceptance

- [ ] The implement-spec `Land` section states the gap fold before the next chunk's first dispatch and the fold record.
- [ ] The implement-spec guidance holds no sentence that folds `main` at the start of the last chunk, and a forbid row names each retired sentence.
- [ ] The review guidance states the fold record, the resolution review, and the landing base after a fold.
- [ ] The final-check guidance names the moved-destination route and no `bench doctor` route for a moved destination.
- [ ] The reference guide states the lease, the window, the gap, the override, the revert back-out route, and `bench commitment block` for a stalled build.
- [ ] `CONTEXT.md` defines the five terms, each with an Avoid list.
- [ ] `docs/known-issues.md` holds no `base-current` row.
- [ ] `CHANGELOG.md` names the lease, the window, the override, the fold record, and the moved-destination face.
