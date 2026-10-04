# Share the landing fixture identity

Blocked by: 05-migrate-execution-fixtures.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/worktree/land_bench_home_test.go, internal/worktree/land_broker_notice_test.go, internal/worktree/land_effects_cleanup_test.go, internal/worktree/land_effects_test.go, internal/worktree/land_facts_test.go, internal/worktree/land_fixtures_test.go, internal/worktree/land_folded_base_test.go, internal/worktree/land_freshness_test.go, internal/worktree/land_journey_test.go, internal/worktree/land_local_capture_test.go, internal/worktree/land_prunes_landed_siblings_test.go, internal/worktree/land_release_refusal_test.go, internal/worktree/land_resume_refusal_test.go, internal/worktree/land_resume_test.go, internal/worktree/land_specless_test.go, internal/worktree/land_surface_test.go, internal/worktree/land_tickets_only_test.go
Covers: GF17, GF29, GF30

## What to build

Replace default identity literals in the listed landing tests with the GF-C1 identity owner. Keep journey gitRun, gitOutput, and descendant as process-supervisor composition. Do not replace their command transport with an unsupervised generic helper.

Create a landing fixture through the worktree journey harness. Its commits use the canonical default identity while descendant effects and cleanup remain owned by that harness.

## Acceptance

- [ ] Worktree Git children retain descendant cleanup and census effects (GF17).
- [ ] The same fixture input family has equivalent repository facts before and after migration (GF29).
- [ ] Each migration chunk passes without its successor implementation (GF30).

## Checkpoint verification

- `bench test --package ./internal/worktree`
- `bench diff`

Use the frozen baseline and the candidate on the same fixture inputs. Exclude commit timestamps and object IDs from comparisons unless the scenario grades them.
Record the observed facts and the named omission or swap in the review pickup. A changed outcome stops the migration.
