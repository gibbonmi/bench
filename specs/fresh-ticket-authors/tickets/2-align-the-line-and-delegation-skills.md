# 2. Align the line and delegation skills

Blocked by: 1-route-each-ticket-to-a-fresh-author.md
Writes: .agents/skills/bench-craft-line/SKILL.md, .agents/skills/bench-craft-delegate/SKILL.md, .agents/skills/bench-craft-delegate/references/delegation-discipline.md, .agents/skills/bench-craft-tickets/SKILL.md, internal/anchors/registry_retained_workflow.go, internal/anchors/registry_data.go, internal/conformance/retained_workflow_test.go, internal/conformance/implementation_continuation_test.go, tests/canary/workflow-guidance-anchors/, internal/anchors/registry_calibration.go, internal/anchors/registry_calibration_test.go, internal/anchors/registry_charge_binding.go, internal/anchors/registry_charge_binding_test.go, internal/anchors/registry_data_test.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_ft311_preparation.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, tests/canary/claude-agent-definitions/, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: FA15, FA16, FA17, FA18, FA19

## What to build

Chunk: FA.

Align the skills to the owner rule that ticket 1 wrote in `.bench/BENCH.md`, and restate no rule.

- `craft-line`: the continuation section governs each ticket author. Outside `--delegate`, a tier move of a fresh author asks the reviewer first. The repair line keeps the ticket's declared model. The `--delegate` section keeps its concurrent authors and its tier range.
- `craft-delegate` and its delegation discipline: point to `.bench/BENCH.md` as the owner of ticket author sessions. A spec-backed ticket goes to a fresh author session on its integration source, in `Blocked by:` order.
- `craft-tickets`: size a ticket to one fresh author context, and work the frontier with one fresh author for each ticket.

Add the rows of FA15 to FA19, change each pin of a retired sentence, and update each canary whose mutation names a changed sentence.

## Acceptance

- [ ] `bench test --check docs-currency-workflow` passes with the FA15 to FA19 needles in place.
- [ ] The retired frontier sentence of FA19 reds its Forbid needle when a probe restores it.
- [ ] The `delegated-per-ticket-author` canary still reds its mutation.
