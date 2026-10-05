# Share the remaining worktree fixture identity

Blocked by: 06-share-landing-fixture-identity.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/worktree/build_test.go, internal/worktree/clean_branch_test.go, internal/worktree/clean_discard_test.go, internal/worktree/clean_discard_transaction_test.go, internal/worktree/clean_operand_test.go, internal/worktree/clean_set_apply_test.go, internal/worktree/clean_set_test.go, internal/worktree/completion_fixture_test.go, internal/worktree/delegated_integration_test.go, internal/worktree/eligibility_test.go, internal/worktree/journey_test.go, internal/worktree/lifecycle_acquire_test.go, internal/worktree/lifecycle_facts_test.go, internal/worktree/lifecycle_policy_test.go, internal/worktree/lifecycle_test.go, internal/worktree/live_binary_test.go, internal/worktree/merge_test.go, internal/worktree/ownership_test.go, internal/worktree/reauthorize_test.go, internal/worktree/recovery_retry_test.go, internal/worktree/resume_test.go, internal/worktree/show_test.go, internal/worktree/test_run_test.go, internal/worktree/worktree_test.go
Covers: GF17, GF29, GF30

## What to build

Replace default identity literals in the remaining worktree fixtures with the GF-C1 identity owner. Keep the journey process wrappers and effect census. Resolve each candidate call by operation, not by its filename.

Create a lifecycle fixture and a cleanup fixture through the worktree journey harness. Both use the canonical default identity without changing process supervision.

## Acceptance

- [ ] Worktree Git children retain descendant cleanup and census effects (GF17).
- [ ] The same fixture input family has equivalent repository facts before and after migration (GF29).
- [ ] Each migration chunk passes without its successor implementation (GF30).

## Checkpoint verification

- `bench test --package ./internal/worktree`
- `bench diff`

Use the frozen baseline and the candidate on the same fixture inputs. Exclude commit timestamps and object IDs from comparisons unless the scenario grades them.
Record the observed facts and the named omission or swap in the review pickup. A changed outcome stops the migration.
