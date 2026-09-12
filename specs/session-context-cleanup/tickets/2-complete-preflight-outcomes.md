# 2. Preflight and report the complete set

Blocked by: 1-plan-explicit-sets.md
Writes: internal/worktree/clean_landed.go, internal/worktree/clean_unclaimed.go, internal/worktree/clean_landed_apply_test.go, internal/worktree/clean_unclaimed_test.go, internal/worktree/clean_set.go (new), internal/worktree/clean_set_apply.go (new), internal/worktree/clean_set_apply_test.go (new), internal/worktree/clean_set_outcomes_test.go (new), internal/worktree/clean_set_wiring_test.go (new), internal/worktree/clean_set_refusal_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/worktree/land_effects.go, internal/worktree/land_effects_cleanup_test.go
Covers: CL4, CL5, CL6, CL7, CL8, CL9, CL12, CL18, CL19

## What to build

Apply complete-set preflight to explicit sets and the existing applicable selectors.
Retain the under-lock per-target recheck and all existing lifecycle effects.
Report completed, failed, retained, and not-attempted outcomes for the full original selection.
For a stale plan, render the exact re-plan command with the same mode and modifiers.
Preserve the current post-first-removal drift test and add a distinct pre-existing-drift case.

## Acceptance

- [ ] Any stale selection detected before apply entry leaves every selected worktree present.
- [ ] Each removable row passes preflight before the first target transaction begins.
- [ ] A later race preserves completed effects and retains the raced target.
- [ ] Every unstarted target appears as not attempted after partial failure.
- [ ] Stale results name an exact safely quoted re-plan command.
- [ ] A spent apply cannot repeat completed side effects.
- [ ] Landing still removes its folded sibling without a user-facing fingerprint round trip.
- [ ] Landing still retains an earlier landed assignment outside its destination-base scope.
