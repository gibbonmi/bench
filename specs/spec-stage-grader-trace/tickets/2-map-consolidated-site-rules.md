# Map each consolidated site rule to a disposition

Blocked by: 1-trace-pin-entry-and-expectation-graders.md
Writes: .agents/skills/bench-craft-spec/references/map-discipline.md, internal/anchors/registry_spec_trace.go (new), internal/anchors/registry_spec_trace_test.go (new), tests/canary/workflow-guidance-anchors/map-discipline-consolidation-both-rules (new), tests/canary/workflow-guidance-anchors/map-discipline-consolidation-cells (new), internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/workflow-guidance-anchors/map-discipline-pin-operator-trace (new), tests/canary/workflow-guidance-anchors/map-discipline-entry-read-grader (new), tests/canary/workflow-guidance-anchors/craft-spec-transaction-failure-rows, tests/canary/workflow-guidance-anchors/craft-spec-two-audience-inventory, tests/canary/workflow-guidance-anchors/map-discipline-addition-disposition, tests/canary/workflow-guidance-anchors/map-discipline-copy-survival-proof, tests/canary/workflow-guidance-anchors/map-discipline-either-side-rows, tests/canary/workflow-guidance-anchors/map-discipline-excluded-edge-caller, tests/canary/workflow-guidance-anchors/map-discipline-executed-root-trace, tests/canary/workflow-guidance-anchors/map-discipline-fixture-reachable-state, tests/canary/workflow-guidance-anchors/map-discipline-flagged-additions, tests/canary/workflow-guidance-anchors/map-discipline-moved-bytes-sweep, tests/canary/workflow-guidance-anchors/map-discipline-pre-review-proof, tests/canary/workflow-guidance-anchors/map-discipline-pre-review-source-consumers, tests/canary/workflow-guidance-anchors/map-discipline-promise-rows, tests/canary/workflow-guidance-anchors/map-discipline-quoted-operands, tests/canary/workflow-guidance-anchors/map-discipline-shipped-surface-claim, tests/canary/workflow-guidance-anchors/map-discipline-source-sentence-table, tests/canary/workflow-guidance-anchors/map-discipline-sweep-depth-bound, tests/canary/workflow-guidance-anchors/map-discipline-sweep-direct-helpers, tests/canary/workflow-guidance-anchors/map-discipline-sweep-named-consumers, tests/canary/workflow-guidance-anchors/map-discipline-sweep-reader-fence, tests/canary/workflow-guidance-anchors/reader-sweep-term
Covers: GT5, GT6, GT15, GT19

## What to build

This ticket is in chunk GT-C1. It delivers the `FT376.consolidation` outcome. A spec that consolidates repeated rules shows each site's old rule and new rule, and each changed cell takes a disposition.

Ticket 1 supplies `specTraceAnchors`, `TestSpecGraderTraceAnchors`, and the `Derived expectations` sub-bullet, which is the insertion point of this ticket.

In the `Before the map locks` section of `map-discipline.md`, add one sub-bullet directly after the `Derived expectations` sub-bullet. It stays on one physical line and holds N5, then N6, exactly as the spec's `Guidance text` section states them.

Append the N5 and N6 anchors to `specTraceAnchors`, with the spec's diagnostics, and append their rules to `TestSpecGraderTraceAnchors`. Add the canaries `map-discipline-consolidation-both-rules` and `map-discipline-consolidation-cells` from the spec's canary table.

The other `Writes:` paths are closure entries that build preflight requires. The ticket does not expect to edit them.

## Acceptance

- [ ] `TestSpecGraderTraceAnchors` shows that a `Before the map locks` section without N5 or N6 raises that needle's diagnostic and no other (GT5, GT6).
- [ ] The test is red before the two anchors exist in `specTraceAnchors`, and green after.
- [ ] The `map-discipline-consolidation-both-rules` mutation raises the N5 diagnostic (GT19).
- [ ] The `map-discipline-consolidation-cells` mutation raises the N6 diagnostic (GT15).
- [ ] Each new canary reds `TestEveryRetainedFixtureBitesThroughRegisteredOwner` before its anchor exists, and passes after.
- [ ] Mutation probe: shorten the N5 needle in `specTraceAnchors` so that it stops after "gives a consolidation table". `TestEveryRetainedFixtureBitesThroughRegisteredOwner` then reds on `map-discipline-consolidation-both-rules`. Restore the needle.
- [ ] The GT-C1 verification commands in the spec's completion plan are green.
