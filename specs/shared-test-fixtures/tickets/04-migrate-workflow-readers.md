# Migrate workflow reader fixtures

Blocked by: 03-migrate-leaf-fixtures.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/handoff/render_test.go, internal/handoff/sections_test.go, internal/handoff/state_file_test.go, internal/handoff/state_scan_test.go, internal/maps/freshness_test.go, internal/roadmap/learning_test.go, internal/roadmap/retro_test.go, internal/roadmap/roadmap_test.go, internal/roadmapflow/flow_test.go, internal/status/handoff_test.go, internal/status/route_test.go, internal/status/status_command_test.go, internal/status/status_counters_test.go, internal/status/status_drain_learnings_test.go, internal/status/status_fixtures_test.go, internal/status/status_gatecache_test.go, internal/status/status_producible_test.go, internal/status/status_render_test.go, internal/status/status_signals_test.go, internal/status/status_test.go
Covers: GF16, GF29, GF30

## What to build

Consume the accepted GF-C1 helper contract in workflow readers. Preserve public fixture API signatures where their callers span packages. An API facade delegates behavior and contains no Git execution policy.

Build the existing status, handoff, roadmap, and map fixtures. Their branch histories and projected answers remain equal after the helper migration.

## Acceptance

- [ ] Generic helper callers retain their fixture results (GF16).
- [ ] The same fixture input family has equivalent repository facts before and after migration (GF29).
- [ ] Each migration chunk passes without its successor implementation (GF30).

## Checkpoint verification

- `bench test --package ./internal/handoff`
- `bench test --package ./internal/maps`
- `bench test --package ./internal/roadmap`
- `bench test --package ./internal/roadmapflow`
- `bench test --package ./internal/status`
- `bench diff`

Use the frozen baseline and the candidate on the same fixture inputs. Exclude commit timestamps and object IDs from comparisons unless the scenario grades them.
Record the observed facts and the named omission or swap in the review pickup. A changed outcome stops the migration.
