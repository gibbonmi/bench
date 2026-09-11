# 1. Plan and apply explicit cleanup sets

Blocked by: none
Writes: internal/worktree/worktree.go, internal/worktree/clean_landed.go, internal/worktree/clean_set.go (new), internal/worktree/clean_set_test.go (new), internal/worktree/clean_set_command_test.go (new), internal/usage/worktree.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: CL1, CL2, CL3, CL10, CL11, CL13, CL14, CL15, CL16

## What to build

Add explicit repeated target selection through the existing cleanup owner.
Resolve the complete set, deduplicate identities, and create one combined fingerprint.
Use the existing per-target lifecycle for apply and retain all current authority checks.
This ticket delivers a working set plan and apply, not a parser-only layer.

## Acceptance

- [ ] One valid explicit set produces one fingerprint and applies through existing target transactions.
- [ ] An unresolved target prevents every removal and reports the failed selection.
- [ ] Aliases cannot duplicate a target or its side effects.
- [ ] Existing single-target and selector effects match a baseline differential matrix.
- [ ] No new budget or measurement prerequisite enters cleanup.
