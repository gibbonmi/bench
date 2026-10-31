# Give the handoff repair suffix and the Shape heading one source each

Blocked by: none
Writes: internal/handoff/state_scan.go, internal/handoff/sections.go, internal/handoffdoc/document.go, internal/conformance/handoff_single_source_test.go
Covers: none

## What to build

Two handoff facts have two spellings each. This ticket gives each fact one owner.
The rendered text does not change.

The State line repair in `internal/handoff` spells the rerun suffix that
`stateRerun` owns. The repair must build its text from `stateRerun`.

The `## Shape` heading has a spelling in `internal/handoff` and a spelling in
`internal/handoffdoc`. The `internal/handoff` package imports
`internal/handoffdoc`, so `internal/handoffdoc` is the owner. It exports the
heading. The `internal/handoff` copy goes away, and the single-source
conformance check reads the exported heading.

## Acceptance

- [ ] `stateRepair` is "rewrite or drop this line" followed by `stateRerun`, so a State refusal hint keeps the same bytes before the refused line. The `internal/handoff` State tests that compare the hint to `stateRepair` stay green.
- [ ] The handoff document parser and renderer, and the conformance Shape check, read one exported heading constant from `internal/handoffdoc`. No `internal/handoff` heading constant remains. The `internal/handoffdoc` and `internal/conformance` tests stay green, and a change to that constant turns the `internal/handoffdoc` legacy-document parse test red.
