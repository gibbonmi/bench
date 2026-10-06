# Sweep changed-function callers with bench consumers

Blocked by: 4-trace-workflow-step-writes.md
Writes: .agents/skills/bench-craft-spec/references/map-discipline.md, internal/anchors/registry_spec_trace.go (new), internal/anchors/registry_spec_trace_test.go (new), tests/canary/workflow-guidance-anchors/map-discipline-unexported-callers (new), CHANGELOG.md, tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns, tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/workflow-guidance-anchors/map-discipline-pin-operator-trace (new), tests/canary/workflow-guidance-anchors/map-discipline-entry-read-grader (new), tests/canary/workflow-guidance-anchors/map-discipline-consolidation-both-rules (new), tests/canary/workflow-guidance-anchors/map-discipline-consolidation-cells (new), tests/canary/workflow-guidance-anchors/map-discipline-quantified-tickets (new), tests/canary/workflow-guidance-anchors/map-discipline-workflow-step-digests (new), tests/canary/workflow-guidance-anchors/craft-spec-transaction-failure-rows, tests/canary/workflow-guidance-anchors/craft-spec-two-audience-inventory, tests/canary/workflow-guidance-anchors/map-discipline-addition-disposition, tests/canary/workflow-guidance-anchors/map-discipline-copy-survival-proof, tests/canary/workflow-guidance-anchors/map-discipline-either-side-rows, tests/canary/workflow-guidance-anchors/map-discipline-excluded-edge-caller, tests/canary/workflow-guidance-anchors/map-discipline-executed-root-trace, tests/canary/workflow-guidance-anchors/map-discipline-fixture-reachable-state, tests/canary/workflow-guidance-anchors/map-discipline-flagged-additions, tests/canary/workflow-guidance-anchors/map-discipline-moved-bytes-sweep, tests/canary/workflow-guidance-anchors/map-discipline-pre-review-proof, tests/canary/workflow-guidance-anchors/map-discipline-pre-review-source-consumers, tests/canary/workflow-guidance-anchors/map-discipline-promise-rows, tests/canary/workflow-guidance-anchors/map-discipline-quoted-operands, tests/canary/workflow-guidance-anchors/map-discipline-shipped-surface-claim, tests/canary/workflow-guidance-anchors/map-discipline-source-sentence-table, tests/canary/workflow-guidance-anchors/map-discipline-sweep-depth-bound, tests/canary/workflow-guidance-anchors/map-discipline-sweep-direct-helpers, tests/canary/workflow-guidance-anchors/map-discipline-sweep-named-consumers, tests/canary/workflow-guidance-anchors/map-discipline-sweep-reader-fence, tests/canary/workflow-guidance-anchors/reader-sweep-term
Covers: GT9, GT10, GT12, GT13, GT17

## What to build

This ticket is the last ticket of chunk GT-C1. It delivers the `FT376.callers` outcome: the caller sweep runs `bench consumers` for each changed function, unexported functions included. It also completes the nine-rule set, so it owns the conformant-tree row, the live-kit row, and the changelog entry.

Ticket 1 supplies `specTraceAnchors` and `TestSpecGraderTraceAnchors`. Tickets 1 to 4 supply the six checklist sub-bullets, and ticket 4 supplies the last one, `Workflow-step writes`.

In the `Before the map locks` section of `map-discipline.md`, add one top-level bullet directly after the last checklist sub-bullet. It stays on one physical line and holds N8, exactly as the spec's `Guidance text` section states it.

Append the N8 anchor to `specTraceAnchors`, with the spec's diagnostic, and append its rule to `TestSpecGraderTraceAnchors`. The test then holds all nine rules. Add the canary `map-discipline-unexported-callers` from the spec's canary table.

Add one `### Spec grader trace` entry under `## [Unreleased]` in `CHANGELOG.md`. The entry names the six new checklist classes and the caller-sweep rule.

The other `Writes:` paths are closure entries that build preflight requires. The ticket does not expect to edit them.

## Acceptance

- [ ] `TestSpecGraderTraceAnchors` shows that a `Before the map locks` section without N8 raises the N8 diagnostic and no other (GT9).
- [ ] The test is red before the N8 anchor exists in `specTraceAnchors`, and green after.
- [ ] The `map-discipline-unexported-callers` mutation raises the N8 diagnostic (GT10).
- [ ] The new canary reds `TestEveryRetainedFixtureBitesThroughRegisteredOwner` before its anchor exists, and passes after.
- [ ] `TestSpecGraderTraceAnchors` holds N1 to N9, and its conformant tree raises none of the nine diagnostics (GT12).
- [ ] `bench test --check docs-currency-workflow` over the kit root raises none of the nine diagnostics (GT13).
- [ ] The `## [Unreleased]` section of `CHANGELOG.md` holds a `### Spec grader trace` entry that names the six classes and the caller-sweep rule (GT17, review-owned).
- [ ] Mutation probe: shorten the N8 needle in `specTraceAnchors` so that it stops after "for each changed function". `TestEveryRetainedFixtureBitesThroughRegisteredOwner` then reds on `map-discipline-unexported-callers`. Restore the needle.
- [ ] The GT-C1 verification commands in the spec's completion plan are green.
