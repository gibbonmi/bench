# Lift the shared field grammar out of the maps package

Blocked by: none
Writes: internal/maps/fields.go (new), internal/maps/schema.go, internal/maps/validation.go, internal/maps/maps_parse_test.go, internal/maps/maps_graph_test.go
Covers: TG4, TG5

## What to build

One exported field scan and one exported graph walk serve both schemas. The
scan covers the fence skip, the prefix fields, and the duplicate-field
diagnostics. The walk covers the blocker-list split, the duplicate edge, the
dangling edge, the self-edge, and one cycle edge.

The decision-map package keeps every delivered diagnostic unchanged. This
ticket adds no diagnostic and changes no verdict. The `internal/tickets`
package consumes the same two symbols in a later ticket.

TG4 closes by diff inspection during review. The reviewer reads the diff and
confirms that one symbol drives both callers.

## Acceptance

- [ ] TG4 — one lifted symbol drives the field scan for both schemas.
- [ ] TG5 — `TestDecisionMapDiagnosticsGolden` asserts the exact ordered slice across the lift.
- [ ] The lifted scan reports a duplicate field for a repeated prefix line.
- [ ] The lifted scan parses a field line inside a fence as no field.

## Delegate charge

You work in the Bench repo on the `ticket-grammar` spec. Read
`specs/ticket-grammar/spec.md` first. Then read `internal/maps/schema.go`,
`internal/maps/validation.go`, and `internal/maps/maps.go` in full. Read
`internal/maps/maps_parse_test.go` and `internal/maps/maps_graph_test.go` next.

Put the lifted code in a new file `internal/maps/fields.go`. Export one field
scan that takes a field table and immutable bytes. Return the parsed field
values and the duplicate-field diagnostics. Keep the fence skip inside the
scan.

Export one graph walk that takes the nodes, the edges, and a diagnostic
formatter. Report a duplicate edge, a dangling edge, a self-edge, and one cycle
edge. Keep decision-map vocabulary out of the walk.

Rewrite `ParseDecisionMap` and `graphDiagnostics` to call the two exported
symbols. Do not copy the scan. Keep every decision-map diagnostic string and
its order unchanged.

First add `TestDecisionMapDiagnosticsGolden` in `internal/maps`. Drive the
parser over the existing test inputs. Assert the exact ordered diagnostic
slice. See it green before the lift, and keep it green after the lift.

Add `TestLiftedFieldScanReportsDuplicateField` and
`TestLiftedFieldScanSkipsFencedLines` in `internal/maps`. Assert the exact
diagnostic text in each test.

Run only `bench worktree exec ft174-ticket-grammar -- go test ./internal/maps/`.
Do not commit. Do not edit the spec.
