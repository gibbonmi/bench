# Evaluate anchor paths

Blocked by: none
Writes: internal/anchors/registry.go, internal/anchors/anchor_harness_diagnostics_test.go, cmd/bench/anchors_command.go, cmd/bench/anchor_help_test.go, internal/gate/lane_select.go, internal/gate/lane_select_test.go, projects/benchkit.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/guidance-prose-budgets/over-budget-skill, tests/canary/line-routing/line-binding-prose-drift, tests/canary/skill-description-budgets/budget-table-missing, tests/canary/skill-description-budgets/description-folded, tests/canary/skill-description-budgets/description-missing, tests/canary/skill-description-budgets/over-budget-command, tests/canary/skill-description-budgets/over-budget-description, tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading, tests/canary/workflow-guidance-anchors/benchkit-review-round-owner, tests/canary/workflow-guidance-anchors/benchkit-review-round-routing, tests/canary/workflow-guidance-anchors/benchkit-spec-ownership, tests/canary/workflow-guidance-anchors/benchkit-system-suite-route
Covers: HP6, HP7, HP8, HP9, HP10

## What to build

Add one anchors-package path evaluation that reads a registered file once and
derives its ordered locations and registry diagnostics from the same bytes. Make
the existing group checks and `bench anchors <path>` use that evaluation logic.
The command keeps its anchor rows. A satisfied registered path returns success.
An unsatisfied registered path returns exit 1 with the exact registry diagnostics.

Derive an anchor-registry lane class from the registry's file fields. A modified
or deleted registered path selects `docs-currency-workflow` beside the existing
path classes. A path with no registry entries keeps the empty successful anchor
query and does not gain the docs check from this class.

## Acceptance

- [ ] A satisfied registered path prints its ordered anchor locations and exits 0.
- [ ] A missing required anchor or present forbidden anchor prints its registry diagnostic and exits 1.
- [ ] An unregistered regular path prints an empty anchor table and exits 0.
- [ ] A modified registered path selects `docs-currency-workflow` in the lane.
- [ ] Deletion of a registered path still selects `docs-currency-workflow`.
