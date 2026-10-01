# Name the pool and residue fixtures

Blocked by: 3-name-assignment-fixtures.md
Writes: internal/worktree/verb_fixture_test.go, internal/worktree/clean_branch_test.go, internal/worktree/live_binary_test.go, internal/worktree/orphan_render_test.go, internal/worktree/pool_reclaim_facts_test.go, internal/worktree/pool_reclaim_test.go, internal/worktree/pool_root_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: VR27

## What to build

Change `unprovableLandedAssignment`, `newResidueGuardFixture`, `newReclaimPool`, and `poolRootFixture` to return one named fixture value each. Declare the value types in `verb_fixture_test.go`, keep each builder in its present file, and update every call site.

This ticket writes `clean_branch_test.go` and `verb_fixture_test.go` after ticket 3, so it waits for that ticket. It moves no verb call and changes no assertion or test name. The over-budget file `pool_reclaim_test.go` stays at or below its base line count.

## Acceptance

- [ ] The tuple scan omits the four builders.
- [ ] Each value type that these builders return is declared in `verb_fixture_test.go`.
- [ ] `TestPackageTestCountPin` passes with no change to `worktreeTestCount`.
- [ ] No over-budget test file grows past its base line count.
