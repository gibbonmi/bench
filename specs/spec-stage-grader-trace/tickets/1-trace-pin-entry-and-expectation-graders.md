# Trace each pin operator, entry read, and derived expectation to its grader

Blocked by: none
Writes: .agents/skills/bench-craft-spec/references/map-discipline.md, internal/anchors/registry_spec_trace.go (new), internal/anchors/registry_spec_trace_test.go (new), tests/canary/workflow-guidance-anchors/map-discipline-pin-operator-trace (new), tests/canary/workflow-guidance-anchors/map-discipline-entry-read-grader (new), internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/workflow-guidance-anchors/craft-spec-transaction-failure-rows, tests/canary/workflow-guidance-anchors/craft-spec-two-audience-inventory, tests/canary/workflow-guidance-anchors/map-discipline-addition-disposition, tests/canary/workflow-guidance-anchors/map-discipline-copy-survival-proof, tests/canary/workflow-guidance-anchors/map-discipline-either-side-rows, tests/canary/workflow-guidance-anchors/map-discipline-excluded-edge-caller, tests/canary/workflow-guidance-anchors/map-discipline-executed-root-trace, tests/canary/workflow-guidance-anchors/map-discipline-fixture-reachable-state, tests/canary/workflow-guidance-anchors/map-discipline-flagged-additions, tests/canary/workflow-guidance-anchors/map-discipline-moved-bytes-sweep, tests/canary/workflow-guidance-anchors/map-discipline-pre-review-proof, tests/canary/workflow-guidance-anchors/map-discipline-pre-review-source-consumers, tests/canary/workflow-guidance-anchors/map-discipline-promise-rows, tests/canary/workflow-guidance-anchors/map-discipline-quoted-operands, tests/canary/workflow-guidance-anchors/map-discipline-shipped-surface-claim, tests/canary/workflow-guidance-anchors/map-discipline-source-sentence-table, tests/canary/workflow-guidance-anchors/map-discipline-sweep-depth-bound, tests/canary/workflow-guidance-anchors/map-discipline-sweep-direct-helpers, tests/canary/workflow-guidance-anchors/map-discipline-sweep-named-consumers, tests/canary/workflow-guidance-anchors/map-discipline-sweep-reader-fence, tests/canary/workflow-guidance-anchors/reader-sweep-term
Covers: GT1, GT2, GT3, GT4, GT14, GT18

## What to build

This ticket is the first ticket of chunk GT-C1. It delivers the `FT376.trace` outcome: a spec author records the pin operators, the entry reads, and the derived-expectation graders before the first review charge.

In the `Before the map locks` section of `map-discipline.md`, add three sub-bullets to the fixed pre-review proof checklist, directly after the `Rendered-shape readers` sub-bullet. Each sub-bullet stays on one physical line, and every current sentence stays byte for byte. The lines hold N1, then N2 with its unanchored gloss, then the shared N3 and N4 line. The spec's `Guidance text` section is the one source of the exact sentences.

Create `internal/anchors/registry_spec_trace.go`. It declares `specTraceAnchors` with the N1 to N4 anchors. Each anchor has group `AfterImplementSpec`, kind `RequireInSection`, file `mapDiscipline`, section `Before the map locks`, and the spec's diagnostic for its needle. Append `specTraceAnchors` at the end of the `registry` chain in `registry_data.go`. That edit changes one line and adds none, because the growth ratchet reds any line that the file gains.

Create `internal/anchors/registry_spec_trace_test.go` with `TestSpecGraderTraceAnchors`. The test writes its needles, sections, and diagnostics independently of the registry, after the precedent of `TestReviewRuleAnchors`, and runs `anchorHarness`. Later tickets of GT-C1 append their rules to `specTraceAnchors` and to this test.

Add the canaries `map-discipline-pin-operator-trace` and `map-discipline-entry-read-grader`. Each `BASE` holds the `map-discipline.md` path. Each `MUTATE.json` holds the one replacement from the spec's canary table. Each `EXPECT` holds the diagnostic of the mutated needle.

The other `Writes:` paths are closure entries that build preflight requires. The ticket does not expect to edit them.

## Acceptance

- [ ] `TestSpecGraderTraceAnchors` shows that a `Before the map locks` section without N1, N2, N3, or N4 raises that needle's diagnostic and no other (GT1, GT2, GT3, GT4).
- [ ] The test is red before the four anchors exist in `specTraceAnchors`, and green after.
- [ ] The `map-discipline-pin-operator-trace` mutation raises the N1 diagnostic (GT14).
- [ ] The `map-discipline-entry-read-grader` mutation raises the N2 diagnostic (GT18).
- [ ] Each new canary reds `TestEveryRetainedFixtureBitesThroughRegisteredOwner` before its anchor exists, and passes after.
- [ ] Mutation probe: shorten the N2 needle in `specTraceAnchors` so that it stops before ", and names its grader". `TestEveryRetainedFixtureBitesThroughRegisteredOwner` then reds on `map-discipline-entry-read-grader`, and `TestSpecGraderTraceAnchors` stays green. Restore the needle.
- [ ] `registry_data.go` and `registry_data_test.go` gain no line.
- [ ] The GT-C1 verification commands in the spec's completion plan are green.
