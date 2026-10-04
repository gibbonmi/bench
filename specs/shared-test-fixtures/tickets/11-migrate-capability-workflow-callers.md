# Migrate capability callers in workflow packages

Blocked by: 10-migrate-capability-leaf-callers.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/commit/lane_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/gate/lane_select_test.go, internal/landing/composition_test.go, internal/landing/landing_helpers_test.go, internal/landing/merge_test.go, internal/landing/state_test.go, internal/reviewrecord/recordcmd/excerpt_test.go, internal/reviewrecord/recordcmd/refusal_test.go, internal/status/status_producible_test.go, internal/status/status_spec_count_test.go, internal/worktree/classifier_shape_test.go, internal/worktree/clean_landed_hostile_test.go, internal/worktree/land_flags_test.go, internal/worktree/live_binary_test.go, internal/worktree/parallel_census_test.go, internal/worktree/pool_reclaim_facts_test.go
Covers: GF22, GF29, GF30

## What to build

Consume GF-C9 Unavailable for workflow sites. Keep current FIFO classifications for socket and device fixtures. Record stable case labels where one function has several call sites.

Exercise the workflow callers, including sockets and devices currently classified as FIFO. Preserve their class instead of inferring a new one from the operation name.

## Acceptance

- [ ] Migrated FIFO and symlink sites preserve kind and class (GF22).
- [ ] The same fixture input family has equivalent repository facts before and after migration (GF29).
- [ ] Each migration chunk passes without its successor implementation (GF30).

## Checkpoint verification

- `bench test --package ./internal/commit`
- `bench test --package ./internal/gate`
- `bench test --package ./internal/landing`
- `bench test --package ./internal/reviewrecord/recordcmd`
- `bench test --package ./internal/status`
- `bench test --package ./internal/worktree`
- `bench diff`

Use the frozen baseline and the candidate on the same fixture inputs. Exclude commit timestamps and object IDs from comparisons unless the scenario grades them.
Record the observed facts and the named omission or swap in the review pickup. A changed outcome stops the migration.
Compare every changed call site with its baseline class before commit. The final ownership chunk adds the persistent class-preservation oracle.
