# 2. Select bounded spec histories

Blocked by: none
Writes: internal/spec/history.go, internal/spec/history_test.go, internal/spec/history_command_test.go (new), internal/spec/spec.go, cmd/bench/main.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/package-core-guard/unrouted-subcommand
Covers: QU4, QU5, QU6, QU7, QU8, QU9, QU10, QU16, QU17

## What to build

Add repeated spec selection with an explicit per-spec event limit.
Consume the existing history producer and render summaries plus selected events.
Preserve positional complete history and the current root command disposition.
The selected query retains all per-target outcomes.

## Acceptance

- [ ] Each spec retains exact matching, commit order, and retire/delete classification.
- [ ] Each history respects the explicit event limit and reports complete counts.
- [ ] Each omitted history names its exact complete-detail command.
- [ ] One failed spec cannot erase successful or empty spec results.
- [ ] Existing positional invocations match the baseline input matrix.
