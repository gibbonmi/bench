# Route a spaced field name to the one-line refusal

Blocked by: derive-the-canonical-path-in-one-leaf-package.md
Writes: internal/maps/validation.go, internal/maps/maps_graph_test.go, internal/conformance/decision_map_integrity_test.go, tests/canary/decision-map-integrity/source-wrapped-field-colon (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: LQ8, LQ9, LQ10, LQ11

## What to build

Verify the premise first: `sourceDiagnostics` in internal/maps/validation.go
cuts each field line on its first colon, and a continuation with a colon
reaches the unexpected-field branch. Then add one guard between the
no-separator branch and the unknown-field branch. When the cut name holds a
space, emit the one-physical-line message. Reuse the `expected` slice. Add a
table row for a colon-bearing continuation, and add a canary fixture under
`tests/canary/decision-map-integrity` in the shape of `source-wrapped-field`.

Repair the fixture inventory in
internal/conformance/decision_map_integrity_test.go: name every fixture
directory, including `source-wrapped-field`, and derive the count from the
inventory's length so the two hand-written `49` literals collapse into one
source.

## Acceptance

- [ ] `TestMapSourcesRequireExactRecordShape` receives the one-physical-line message for a continuation with a space and a colon.
- [ ] The unknown-field, second-path, second-url, and mixed-locator rows still receive the unexpected-field message.
- [ ] The new fixture bites through the decision-map integrity check.
- [ ] `TestDecisionMapIntegrityFixtureInventoryRejectsDeletion` passes with the repaired inventory.
- [ ] Self-probe: move the guard after the unknown-field branch, and report the new table row red.
