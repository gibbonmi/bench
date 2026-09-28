# 6. Repair the glossary shift namespace

Blocked by: none
Writes: CONTEXT.md, internal/anchors/registry_data.go, internal/anchors/registry_decision_maps.go, internal/anchors/registry_decision_maps_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/docs-currency-token-diet/signal-vocabulary-drift, tests/canary/workflow-guidance-anchors/context-acceptance-row-vocabulary, tests/canary/workflow-guidance-anchors/context-coverage-map-term, tests/canary/workflow-guidance-anchors/context-coverage-row-parts, tests/canary/workflow-guidance-anchors/context-coverage-row-vocabulary, tests/canary/workflow-guidance-anchors/context-decision-map-term, tests/canary/workflow-guidance-anchors/context-reader-sweep-term, tests/canary/workflow-guidance-anchors/context-ticket-vocabulary, .agents/skills/bench-craft-spec/references/map-discipline.md, projects/benchkit.md
Covers: RI57, RI82

## What to build

Chunk: RI-C1a.

Repair the glossary term `unclaimed ref` in `CONTEXT.md` so that it names the shift namespace as `refs/heads/bench/shift-`.
The ledger declares that prefix.
Rewrite the glossary term `holder` so that it names an active or cleanup-pending recorded assignment branch, or a unique root, and no landed ref.
That is the rule the reviewer fixed on 2026-09-27, because a ref beneath an ancestry-landed ref is itself landed.

Run `bench anchors CONTEXT.md` before the edit and keep every anchored line intact.
The anchor registries and the canary fixtures in `Writes:` are the closure the preflight names.
This ticket changes them only when the edit moves an anchored line.

## Acceptance

- [ ] The glossary term `unclaimed ref` names `refs/heads/bench/shift-`.
- [ ] The glossary term `holder` names a recorded assignment branch or a unique root, and no landed ref.
- [ ] `bench anchors CONTEXT.md` prints the same eight anchors after the edit.
