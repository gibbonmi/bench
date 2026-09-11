# 4. Apply reviewed owner budgets

Blocked by: 1-select-worktrees.md, 2-select-histories.md, 3-guide-relevant-reads.md
Writes: internal/worktree/list.go, internal/worktree/list_selected_test.go (new), internal/spec/history.go, internal/spec/history_command_test.go (new), .agents/skills/bench-craft-cli/SKILL.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: QU14, QU15

## What to build

Entry requires the measurement report and an explicit per-surface reviewer budget decision.
A spec-authoring pass must record concrete values and complete owner fences before this ticket starts.
Stop before product writes when either prerequisite is absent.
Apply only the approved defaults through their existing command owners.
Preserve current unrelated policies and explicit complete-detail reads.

## Acceptance

- [ ] The approved policy identifies every changed surface and its numeric byte value.
- [ ] Boundary fixtures use below, equal, and above-budget results with multibyte text.
- [ ] Each omitted result retains the approved metadata and exact detail route.
- [ ] The build chooses no budget value and changes no unrelated owner policy.
