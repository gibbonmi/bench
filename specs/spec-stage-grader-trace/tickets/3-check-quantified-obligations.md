# Check each quantified obligation at each site and in each ticket

Blocked by: 2-map-consolidated-site-rules.md
Writes: .agents/skills/bench-craft-spec/references/map-discipline.md, internal/anchors/registry_spec_trace.go (new), internal/anchors/registry_spec_trace_test.go (new), tests/canary/workflow-guidance-anchors/map-discipline-quantified-tickets (new), internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/workflow-guidance-anchors/map-discipline-pin-operator-trace (new), tests/canary/workflow-guidance-anchors/map-discipline-entry-read-grader (new), tests/canary/workflow-guidance-anchors/map-discipline-consolidation-both-rules (new), tests/canary/workflow-guidance-anchors/map-discipline-consolidation-cells (new), tests/canary/workflow-guidance-anchors/craft-spec-transaction-failure-rows, tests/canary/workflow-guidance-anchors/craft-spec-two-audience-inventory, tests/canary/workflow-guidance-anchors/map-discipline-addition-disposition, tests/canary/workflow-guidance-anchors/map-discipline-copy-survival-proof, tests/canary/workflow-guidance-anchors/map-discipline-either-side-rows, tests/canary/workflow-guidance-anchors/map-discipline-excluded-edge-caller, tests/canary/workflow-guidance-anchors/map-discipline-executed-root-trace, tests/canary/workflow-guidance-anchors/map-discipline-fixture-reachable-state, tests/canary/workflow-guidance-anchors/map-discipline-flagged-additions, tests/canary/workflow-guidance-anchors/map-discipline-moved-bytes-sweep, tests/canary/workflow-guidance-anchors/map-discipline-pre-review-proof, tests/canary/workflow-guidance-anchors/map-discipline-pre-review-source-consumers, tests/canary/workflow-guidance-anchors/map-discipline-promise-rows, tests/canary/workflow-guidance-anchors/map-discipline-quoted-operands, tests/canary/workflow-guidance-anchors/map-discipline-shipped-surface-claim, tests/canary/workflow-guidance-anchors/map-discipline-source-sentence-table, tests/canary/workflow-guidance-anchors/map-discipline-sweep-depth-bound, tests/canary/workflow-guidance-anchors/map-discipline-sweep-direct-helpers, tests/canary/workflow-guidance-anchors/map-discipline-sweep-named-consumers, tests/canary/workflow-guidance-anchors/map-discipline-sweep-reader-fence, tests/canary/workflow-guidance-anchors/reader-sweep-term
Covers: GT7, GT8

## What to build

This ticket is in chunk GT-C1. It delivers the `FT376.quantified` outcome: a spec author checks every quantified obligation at each affected site and across all tickets.

Ticket 1 supplies `specTraceAnchors` and `TestSpecGraderTraceAnchors`. Ticket 2 supplies the `Consolidated rules` sub-bullet, which is the insertion point of this ticket.

In the `Before the map locks` section of `map-discipline.md`, add one sub-bullet directly after the `Consolidated rules` sub-bullet. It stays on one physical line and holds N7, exactly as the spec's `Guidance text` section states it.

Append the N7 anchor to `specTraceAnchors`, with the spec's diagnostic, and append its rule to `TestSpecGraderTraceAnchors`. Add the canary `map-discipline-quantified-tickets` from the spec's canary table.

The other `Writes:` paths are closure entries that build preflight requires. The ticket does not expect to edit them.

## Acceptance

- [ ] `TestSpecGraderTraceAnchors` shows that a `Before the map locks` section without N7 raises the N7 diagnostic and no other (GT7).
- [ ] The test is red before the N7 anchor exists in `specTraceAnchors`, and green after.
- [ ] The `map-discipline-quantified-tickets` mutation raises the N7 diagnostic (GT8).
- [ ] The new canary reds `TestEveryRetainedFixtureBitesThroughRegisteredOwner` before its anchor exists, and passes after.
- [ ] Mutation probe: shorten the N7 needle in `specTraceAnchors` so that it stops after "at each affected site". `TestEveryRetainedFixtureBitesThroughRegisteredOwner` then reds on `map-discipline-quantified-tickets`. Restore the needle.
- [ ] The GT-C1 verification commands in the spec's completion plan are green.
