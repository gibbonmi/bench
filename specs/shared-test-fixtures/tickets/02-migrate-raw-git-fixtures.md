# Migrate raw Git fixture callers

Blocked by: 01-share-git-fixture-owner.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/diff/command_test.go, internal/diff/compatibility_test.go, internal/diff/explicit_base_test.go, internal/diff/identity_test.go, internal/diff/matrix_test.go, internal/diff/review_base_test.go, internal/diff/source_tip_pair_test.go, internal/git/admin_readers_test.go, internal/git/checkout_test.go, internal/git/facts_test.go, internal/git/localnote_test.go, internal/git/push_destination_test.go, internal/git/refs_test.go, internal/git/staged_test.go, internal/git/testhelpers_test.go, internal/git/worktree_admin_enum_test.go, internal/git/worktree_admin_hostile_test.go
Covers: GF16, GF29, GF30

## What to build

Consume the accepted GF-C1 helper contract. Replace generic Git wrappers and their local callers. Preserve raw stdout, combined output, and intentional identity subjects at their existing consumers. Keep specialized probes only where their transport or failure contract requires it.

Create the existing raw Git fixtures. Compare untrimmed query bytes, staged content, author overrides, and deliberate failure results before and after migration.

## Acceptance

- [ ] Generic helper callers retain their fixture results (GF16).
- [ ] The same fixture input family has equivalent repository facts before and after migration (GF29).
- [ ] Each migration chunk passes without its successor implementation (GF30).

## Checkpoint verification

- `bench test --package ./internal/diff`
- `bench test --package ./internal/git`
- `bench diff`

Use the frozen baseline and the candidate on the same fixture inputs. Exclude commit timestamps and object IDs from comparisons unless the scenario grades them.
Record the observed facts and the named omission or swap in the review pickup. A changed outcome stops the migration.
