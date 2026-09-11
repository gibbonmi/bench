# 2. Select bounded spec histories

Blocked by: none
Writes: internal/spec/history.go, internal/spec/history_test.go, internal/spec/history_command_test.go (new), internal/spec/spec.go, cmd/bench/main.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/package-core-guard/unrouted-subcommand
Covers: QU4, QU5, QU6, QU7, QU8, QU19, QU20, QU21, QU22, QU23, QU24, QU25

## What to build

Add repeated spec selection with an explicit per-spec event limit.
Consume the existing history producer and render summaries plus selected events.
Keep the shared producer complete for the roadmap context reader.
Do not apply selection or limits inside `History`.
Preserve positional complete history and the current root command disposition.
The selected query retains all per-target outcomes.

## Acceptance

- [ ] Each spec retains exact matching, commit order, and retire/delete classification.
- [ ] Each history respects the explicit event limit and reports complete counts.
- [ ] Each omitted history names its exact complete-detail command.
- [ ] One failed spec cannot erase successful or empty spec results.
- [ ] Existing positional invocations match the baseline input matrix.
- [ ] The shared history producer returns the baseline complete sequence after selected queries run.
- [ ] The serialized byte count comes from the complete history before projection.
- [ ] One unrepresentable commit subject cannot erase other selected spec results.
