# Migrate leaf command fixtures

Blocked by: 02-migrate-raw-git-fixtures.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/consumers/blast_edges_test.go, internal/consumers/blast_test.go, internal/consumers/citation_test.go, internal/consumers/refuse_test.go, internal/gitguard/checker_junction_test.go, internal/guards/guards_test.go, internal/outline/outline_test.go, internal/poolkey/poolkey_test.go, internal/probe/probe_test.go, internal/skillsindex/command_test.go, internal/spec/history_command_test.go, internal/spec/history_selected_test.go, internal/spec/spec_test.go, internal/stophook/stophook_test.go, internal/structure/structure_test.go, internal/treetarget/identify_test.go, internal/treetarget/run_test.go, tests/canary/injected-ports/unregistered-port, tests/canary/package-core-guard/reintroduced-bare-skip
Covers: GF16, GF29, GF30

## What to build

Consume the accepted GF-C1 helper contract in the listed leaf packages. Remove private generic wrappers after their last caller migrates. Domain fixture builders may compose the owner without restating execution or identity policy.

Run the existing command fixture scenarios through the shared Git helpers. Their repository roots, branch facts, and expected refusal cases remain equal.

## Acceptance

- [ ] Generic helper callers retain their fixture results (GF16).
- [ ] The same fixture input family has equivalent repository facts before and after migration (GF29).
- [ ] Each migration chunk passes without its successor implementation (GF30).

## Checkpoint verification

- `bench test --package ./internal/consumers`
- `bench test --package ./internal/gitguard`
- `bench test --package ./internal/guards`
- `bench test --package ./internal/outline`
- `bench test --package ./internal/poolkey`
- `bench test --package ./internal/probe`
- `bench test --package ./internal/skillsindex`
- `bench test --package ./internal/spec`
- `bench test --package ./internal/stophook`
- `bench test --package ./internal/structure`
- `bench test --package ./internal/treetarget`
- `bench diff`

Use the frozen baseline and the candidate on the same fixture inputs. Exclude commit timestamps and object IDs from comparisons unless the scenario grades them.
Record the observed facts and the named omission or swap in the review pickup. A changed outcome stops the migration.
