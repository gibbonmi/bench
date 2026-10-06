# Trace each spec row to its grader before the first review

Blocked by: none
Writes: .agents/skills/bench-craft-spec/references/map-discipline.md, internal/anchors/registry_spec_trace.go (new), internal/anchors/registry_spec_trace_test.go (new), tests/canary/workflow-guidance-anchors/map-discipline-pin-operator-trace (new), tests/canary/workflow-guidance-anchors/map-discipline-entry-read-grader (new), tests/canary/workflow-guidance-anchors/map-discipline-consolidation-both-rules (new), tests/canary/workflow-guidance-anchors/map-discipline-consolidation-cells (new), tests/canary/workflow-guidance-anchors/map-discipline-quantified-tickets (new), tests/canary/workflow-guidance-anchors/map-discipline-workflow-step-digests (new), tests/canary/workflow-guidance-anchors/map-discipline-unexported-callers (new), tests/canary/workflow-guidance-anchors/map-discipline-derived-grader (new), tests/canary/workflow-guidance-anchors/map-discipline-no-expectation-under-test (new), CHANGELOG.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/workflow-guidance-anchors/craft-spec-transaction-failure-rows, tests/canary/workflow-guidance-anchors/craft-spec-two-audience-inventory, tests/canary/workflow-guidance-anchors/map-discipline-addition-disposition, tests/canary/workflow-guidance-anchors/map-discipline-copy-survival-proof, tests/canary/workflow-guidance-anchors/map-discipline-either-side-rows, tests/canary/workflow-guidance-anchors/map-discipline-excluded-edge-caller, tests/canary/workflow-guidance-anchors/map-discipline-executed-root-trace, tests/canary/workflow-guidance-anchors/map-discipline-fixture-reachable-state, tests/canary/workflow-guidance-anchors/map-discipline-flagged-additions, tests/canary/workflow-guidance-anchors/map-discipline-moved-bytes-sweep, tests/canary/workflow-guidance-anchors/map-discipline-pre-review-proof, tests/canary/workflow-guidance-anchors/map-discipline-pre-review-source-consumers, tests/canary/workflow-guidance-anchors/map-discipline-promise-rows, tests/canary/workflow-guidance-anchors/map-discipline-quoted-operands, tests/canary/workflow-guidance-anchors/map-discipline-shipped-surface-claim, tests/canary/workflow-guidance-anchors/map-discipline-source-sentence-table, tests/canary/workflow-guidance-anchors/map-discipline-sweep-depth-bound, tests/canary/workflow-guidance-anchors/map-discipline-sweep-direct-helpers, tests/canary/workflow-guidance-anchors/map-discipline-sweep-named-consumers, tests/canary/workflow-guidance-anchors/map-discipline-sweep-reader-fence, tests/canary/workflow-guidance-anchors/reader-sweep-term, tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns, tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary
Covers: GT1, GT2, GT3, GT4, GT5, GT6, GT7, GT8, GT9, GT10, GT11, GT12, GT13, GT14, GT15, GT16, GT17, GT18, GT19, GT20, GT21

## What to build

This ticket is the one ticket of chunk GT-C1. It delivers all five FT376 commitment criteria. Before the first review charge, a spec author traces each pinned check, entry read, and derived expectation to its grader. The author also maps consolidated site rules, checks quantified obligations, traces workflow-step writes, and sweeps callers with `bench consumers`.

In the `Before the map locks` section of `map-discipline.md`, add the six checklist sub-bullets directly after the `Rendered-shape readers` sub-bullet. Then add the N8 top-level bullet directly after the last sub-bullet. Each line stays on one physical line, and every current sentence stays byte for byte. The spec's `Guidance text` section is the one source of the exact sentences, their order, and the unanchored gloss.

Create `internal/anchors/registry_spec_trace.go`. It declares `specTraceAnchors` with the nine anchors N1 to N9. Each anchor has group `AfterImplementSpec`, kind `RequireInSection`, file `mapDiscipline`, section `Before the map locks`, and the spec's diagnostic for its needle. Append `specTraceAnchors` at the end of the `registry` chain in `registry_data.go`. That edit changes one line and adds none, because the growth ratchet reds any line that the file gains.

Create `internal/anchors/registry_spec_trace_test.go` with `TestSpecGraderTraceAnchors`. The test writes its needles, sections, and diagnostics independently of the registry, after the precedent of `TestReviewRuleAnchors`, and runs `anchorHarness`.

Add the nine canaries from the spec's canary table. Each `BASE` holds the `map-discipline.md` path. Each `MUTATE.json` holds the one replacement from that table. Each `EXPECT` holds the diagnostic of the mutated needle.

Add one `### Spec grader trace` entry under `## [Unreleased]` in `CHANGELOG.md`. The entry names the six new checklist classes and the caller-sweep rule.

The other `Writes:` paths are closure entries that build preflight requires. The ticket does not expect to edit them.

## Acceptance

- [ ] In `TestSpecGraderTraceAnchors`, a section without one needle raises that needle's diagnostic and no other. This holds for N1 to N9 (GT1, GT2, GT3, GT4, GT5, GT6, GT7, GT9, GT11).
- [ ] `TestSpecGraderTraceAnchors` is red before the nine anchors exist in `specTraceAnchors`, and green after.
- [ ] The conformant tree of `TestSpecGraderTraceAnchors` holds N1 to N9 and raises none of the nine diagnostics (GT12).
- [ ] `bench test --check docs-currency-workflow` over the kit root raises none of the nine diagnostics (GT13).
- [ ] Each canary mutation raises the diagnostic of its needle: `map-discipline-pin-operator-trace` (GT14), `map-discipline-entry-read-grader` (GT18), `map-discipline-consolidation-both-rules` (GT19), `map-discipline-consolidation-cells` (GT15), `map-discipline-quantified-tickets` (GT8), `map-discipline-workflow-step-digests` (GT16), and `map-discipline-unexported-callers` (GT10). The same holds for `map-discipline-derived-grader` (GT20) and `map-discipline-no-expectation-under-test` (GT21).
- [ ] Each new canary reds `TestEveryRetainedFixtureBitesThroughRegisteredOwner` before its anchor exists, and passes after.
- [ ] The `## [Unreleased]` section of `CHANGELOG.md` holds a `### Spec grader trace` entry that names the six classes and the caller-sweep rule (GT17, review-owned).
- [ ] `registry_data.go` and `registry_data_test.go` gain no line.
- [ ] Mutation probes, one for each canary: shorten one needle in `specTraceAnchors` at a time. Run `TestEveryRetainedFixtureBitesThroughRegisteredOwner`, observe the red on the named canary, and restore the needle. `TestSpecGraderTraceAnchors` stays green for each probe.
  - N1 stops after "`Pin operators`", and `map-discipline-pin-operator-trace` reds.
  - N2 stops before ", and names its grader", and `map-discipline-entry-read-grader` reds.
  - N3 stops after "`Derived expectations` names", and `map-discipline-derived-grader` reds.
  - N4 stops after "No expectation comes from the code", and `map-discipline-no-expectation-under-test` reds.
  - N5 stops after "gives a consolidation table", and `map-discipline-consolidation-both-rules` reds.
  - N6 stops after "an acceptance row", and `map-discipline-consolidation-cells` reds.
  - N7 stops after "at each affected site", and `map-discipline-quantified-tickets` reds.
  - N9 stops after "for each new workflow step,", and `map-discipline-workflow-step-digests` reds.
  - N8 stops after "for each changed function", and `map-discipline-unexported-callers` reds.
- [ ] The GT-C1 verification commands in the spec's completion plan are green.
