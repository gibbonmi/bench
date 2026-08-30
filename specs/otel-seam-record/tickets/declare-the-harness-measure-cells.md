# 13. Declare the harness measure cells

Blocked by: none
Line: opus / low
Rows: OT23, OT24
Writes: internal/harnesses/harnesses.go, internal/harnesses/harnesses_test.go, internal/harnesses/command.go, internal/harnesses/command_test.go

## What to build

The harness record gains four measure cells: the tokens, the tool calls, the
Read paths, and the turns. Each cell names the supplier of that measure. Each
cell reads `Unknown` until FT204 supplies a source. FT274's rule that a measure
with no source stays absent governs the span attribute, not the cell.

`bench harnesses <harness>` prints the four cells. An operator then reads the
declaration without the source.

`internal/harnesses/harnesses_test.go` gains an independent want-list
expectation for the four cells, beside the mechanics expectation that already
reds a dropped name. A record that drops a cell therefore reds.

This ticket writes no span and no record file. It runs in parallel with ticket
1.

## Acceptance

- [ ] OT23: the harness record carries the four measure cells on every row, and the want-list expectation reds a record without them.
- [ ] OT24: `bench harnesses claude` prints the four measure cells.
- [ ] each cell names its supplier, and each cell reads `Unknown` until FT204 supplies a source.
