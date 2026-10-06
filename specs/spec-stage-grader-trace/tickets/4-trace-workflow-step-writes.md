# Trace each new workflow step's writes through checkpoint digests

Blocked by: 3-check-quantified-obligations.md
Writes: .agents/skills/bench-craft-spec/references/map-discipline.md, internal/anchors/registry_spec_trace.go (new), internal/anchors/registry_spec_trace_test.go (new), tests/canary/workflow-guidance-anchors/map-discipline-workflow-step-digests (new), internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/workflow-guidance-anchors/map-discipline-pin-operator-trace (new), tests/canary/workflow-guidance-anchors/map-discipline-entry-read-grader (new), tests/canary/workflow-guidance-anchors/map-discipline-consolidation-both-rules (new), tests/canary/workflow-guidance-anchors/map-discipline-consolidation-cells (new), tests/canary/workflow-guidance-anchors/map-discipline-quantified-tickets (new), tests/canary/workflow-guidance-anchors/craft-spec-transaction-failure-rows, tests/canary/workflow-guidance-anchors/craft-spec-two-audience-inventory, tests/canary/workflow-guidance-anchors/map-discipline-addition-disposition, tests/canary/workflow-guidance-anchors/map-discipline-copy-survival-proof, tests/canary/workflow-guidance-anchors/map-discipline-either-side-rows, tests/canary/workflow-guidance-anchors/map-discipline-excluded-edge-caller, tests/canary/workflow-guidance-anchors/map-discipline-executed-root-trace, tests/canary/workflow-guidance-anchors/map-discipline-fixture-reachable-state, tests/canary/workflow-guidance-anchors/map-discipline-flagged-additions, tests/canary/workflow-guidance-anchors/map-discipline-moved-bytes-sweep, tests/canary/workflow-guidance-anchors/map-discipline-pre-review-proof, tests/canary/workflow-guidance-anchors/map-discipline-pre-review-source-consumers, tests/canary/workflow-guidance-anchors/map-discipline-promise-rows, tests/canary/workflow-guidance-anchors/map-discipline-quoted-operands, tests/canary/workflow-guidance-anchors/map-discipline-shipped-surface-claim, tests/canary/workflow-guidance-anchors/map-discipline-source-sentence-table, tests/canary/workflow-guidance-anchors/map-discipline-sweep-depth-bound, tests/canary/workflow-guidance-anchors/map-discipline-sweep-direct-helpers, tests/canary/workflow-guidance-anchors/map-discipline-sweep-named-consumers, tests/canary/workflow-guidance-anchors/map-discipline-sweep-reader-fence, tests/canary/workflow-guidance-anchors/reader-sweep-term
Covers: GT11, GT16

## What to build

This ticket is in chunk GT-C1. It delivers the `FT376.digests` outcome: for each new workflow step, a spec author traces every write through the digests that later checkpoints compare.

Ticket 1 supplies `specTraceAnchors` and `TestSpecGraderTraceAnchors`. Ticket 3 supplies the `Quantified obligations` sub-bullet, which is the insertion point of this ticket.

In the `Before the map locks` section of `map-discipline.md`, add one sub-bullet directly after the `Quantified obligations` sub-bullet. This sub-bullet is the last item of the checklist. It stays on one physical line and holds N9, exactly as the spec's `Guidance text` section states it.

Append the N9 anchor to `specTraceAnchors`, with the spec's diagnostic, and append its rule to `TestSpecGraderTraceAnchors`. Add the canary `map-discipline-workflow-step-digests` from the spec's canary table.

The other `Writes:` paths are closure entries that build preflight requires. The ticket does not expect to edit them.

## Acceptance

- [ ] `TestSpecGraderTraceAnchors` shows that a `Before the map locks` section without N9 raises the N9 diagnostic and no other (GT11).
- [ ] The test is red before the N9 anchor exists in `specTraceAnchors`, and green after.
- [ ] The `map-discipline-workflow-step-digests` mutation raises the N9 diagnostic (GT16).
- [ ] The new canary reds `TestEveryRetainedFixtureBitesThroughRegisteredOwner` before its anchor exists, and passes after.
- [ ] Mutation probe: shorten the N9 needle in `specTraceAnchors` so that it stops after "for each new workflow step,". `TestEveryRetainedFixtureBitesThroughRegisteredOwner` then reds on `map-discipline-workflow-step-digests`. Restore the needle.
- [ ] The GT-C1 verification commands in the spec's completion plan are green.
