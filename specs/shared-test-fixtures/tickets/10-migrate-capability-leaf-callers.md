# Migrate capability callers in leaf packages

Blocked by: 09-own-capability-reasons.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/canary/mutation_test.go, internal/census/census_fifo_test.go, internal/census/census_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/env/git_policy_test.go, internal/git/admin_readers_test.go, internal/git/refs_test.go, internal/git/worktree_admin_enum_test.go, internal/git/worktree_admin_hostile_test.go, internal/guards/guards_test.go, internal/learnings/learnings_test.go, internal/outline/outline_test.go, internal/prose/walk_test.go, internal/skillsindex/skillsindex_reference_test.go, internal/spec/spec_test.go
Covers: GF22, GF29, GF30

## What to build

Consume GF-C9 Unavailable for the listed leaf sites. Remove caller-authored capability prefixes and preserve useful operation detail. Keep every original skip condition and class.

Exercise the leaf callers with their existing operation conditions. Their class and skip predicate remain equal while the reason prefix comes from capability.

## Acceptance

- [ ] Migrated FIFO and symlink sites preserve kind and class (GF22).
- [ ] The same fixture input family has equivalent repository facts before and after migration (GF29).
- [ ] Each migration chunk passes without its successor implementation (GF30).

## Checkpoint verification

- `bench test --package ./internal/canary`
- `bench test --package ./internal/census`
- `bench test --package ./internal/env`
- `bench test --package ./internal/git`
- `bench test --package ./internal/guards`
- `bench test --package ./internal/learnings`
- `bench test --package ./internal/outline`
- `bench test --package ./internal/prose`
- `bench test --package ./internal/skillsindex`
- `bench test --package ./internal/spec`
- `bench diff`

Use the frozen baseline and the candidate on the same fixture inputs. Exclude commit timestamps and object IDs from comparisons unless the scenario grades them.
Record the observed facts and the named omission or swap in the review pickup. A changed outcome stops the migration.
Compare every changed call site with its baseline class before commit. The final ownership chunk adds the persistent class-preservation oracle.
