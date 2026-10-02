# 10. Name each table in production

Blocked by: 9-refuse-a-read-below-a-census-entry.md
Writes: internal/worktree/table_name_census_test.go (new), internal/worktree/classifier.go, internal/worktree/worktree.go, internal/worktree/reset.go, internal/worktree/list_selected.go, internal/worktree/list.go, internal/worktree/build.go, internal/worktree/ownership.go, internal/worktree/pool_reclaim.go, internal/worktree/verb_fixture_test.go, internal/worktree/verb_runner_check_test.go, internal/worktree/list_selected_test.go, internal/worktree/build_test.go, internal/worktree/clean_landed_apply_test.go, internal/worktree/clean_landed_hostile_test.go, internal/worktree/clean_landed_test.go, internal/worktree/clean_set_test.go, internal/worktree/clean_unclaimed_test.go, internal/worktree/eligibility_test.go, internal/worktree/journey_race_test.go, internal/worktree/land_release_refusal_test.go, internal/worktree/land_surface_test.go, internal/worktree/list_actions_test.go, internal/worktree/orphan_render_test.go, internal/worktree/path_identifier_test.go, internal/worktree/pool_reclaim_test.go, internal/worktree/reset_apply_test.go, internal/worktree/reset_plan_test.go, internal/worktree/reset_refusal_test.go, internal/worktree/reset_restore_refusal_test.go, internal/worktree/snapshot_test.go, internal/worktree/parallel_census_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: WS71, WS72, WS73, WS74

## What to build

Chunk: SR-C7.

Name each table that the package renders with a package constant. Each `toon.Table` and
`toon.TableTyped` call in the non-test source takes the constant. The test-only
`cleanupTable` and `selectedTable` constants move to production.

Add a table-name census in a new test file. It derives the table names from the
constants that the non-test source passes to `toon.Table` and `toon.TableTyped`. It
reports a non-test call that passes a string literal. In the test files, it reports a
string literal that equals a table name as the block argument of `mustRows`,
`readVerbRows`, `Rows`, `toon.Table`, or `toon.TableTyped`. It also reports a test string
literal that begins with a table name and `[`. Each report names the file and the line.

Replace each reported test literal with the constant. A path segment such as
`"worktrees"` in `filepath.Join` is not a block argument, so the census does not report
it. `pool_reclaim_test.go` and `worktree_test.go` are over their line budgets, so the
change does not grow them. Raise `worktreeTestCount` by the number of new top-level tests
in this commit.

## Acceptance

- [ ] Each `toon.Table` and `toon.TableTyped` call in the package's non-test source passes an identifier.
- [ ] A synthetic block-argument literal and a synthetic rendered-header literal each draw a report with the file and the line.
- [ ] The table-name census over the live test files reports nothing.
