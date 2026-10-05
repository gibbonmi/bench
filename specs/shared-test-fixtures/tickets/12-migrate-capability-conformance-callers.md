# Migrate capability callers in conformance and CLI tests

Blocked by: 11-migrate-capability-workflow-callers.md
Writes: cmd/bench/anchor_help_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/bounds_policy_test.go, internal/conformance/checks_test.go, internal/conformance/claude_agent_definitions_test.go, internal/conformance/fixture_bite_test.go, internal/conformance/guidance_token_sweep_test.go, internal/conformance/harness_record_test.go, internal/conformance/prose_budget_test.go, internal/conformance/skill_description_budgets_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: GF22, GF29, GF30

## What to build

Consume GF-C9 Unavailable for conformance and CLI sites. Do not rewrite source examples inside strings. The real fixture operation still determines when the skip occurs.

Exercise a conformance fixture that cannot create a symlink and one that cannot create a FIFO. Preserve the diagnostics payload and structured skip classification.

## Acceptance

- [ ] Migrated FIFO and symlink sites preserve kind and class (GF22).
- [ ] The same fixture input family has equivalent repository facts before and after migration (GF29).
- [ ] Each migration chunk passes without its successor implementation (GF30).

## Checkpoint verification

- `bench test --package ./cmd/bench`
- `bench test --package ./internal/conformance`
- `bench diff`

Use the frozen baseline and the candidate on the same fixture inputs. Exclude commit timestamps and object IDs from comparisons unless the scenario grades them.
Record the observed facts and the named omission or swap in the review pickup. A changed outcome stops the migration.
Compare every changed call site with its baseline class before commit. The final ownership chunk adds the persistent class-preservation oracle.
