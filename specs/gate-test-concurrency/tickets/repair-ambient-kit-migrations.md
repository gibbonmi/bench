# Repair the remaining ambient-kit migrations

Blocked by: repair-ambient-red-contract.md
Ownership fence: `internal/gate/absent_tool_slot_test.go`, `internal/gate/gate_test.go`, `internal/gate/verdict_reuse_test.go`
Integration surfaces: `internal/gate/gate.go` exported entries→`internal/gate/gate_test.go` + RAK1/no-module-exported and `internal/gate/verdict_reuse_test.go` + RAK1/no-module-exported; `internal/gate/evaluation.go` prospective empty-scoping branch→existing exported calls + RAK1/prospective-exported; `internal/gate/absent_tool_slot_test.go` fixture table→RAK1/absent-tool-root
Contracts: exported calls in `internal/gate/gate_test.go` and `internal/gate/verdict_reuse_test.go` retain public-entry coverage when `internal/gate/component_decision.go` returns empty no-module scoping; prospective exported calls retain public-entry coverage when `internal/gate/evaluation.go` returns empty scoping; the root identity in `internal/gate/absent_tool_slot_test.go` crosses to `internal/gate/manifest.go`, asserted by RAK1 against the real producer branches
Closure: RAK1/no-module-exported, RAK1/prospective-exported, RAK1/absent-tool-root

## What to build

Audit every `t.Parallel` test in `verdict_reuse_test.go`,
`prospective_test.go`, and `gate_test.go` that calls `Execute`,
`ExecuteReusingFreshGreen`, `RunCommand`, or another environment-resolving
entry over a fixture or throwaway project root. Retain the exported calls for
fixtures without a regular `go.mod`: component scoping returns empty before
phase-table resolution can consume kit identity. Retain prospective exported
calls too: prospective evaluation returns empty scoping before the kit-taking
component helper. Do not route either inert class through a private seam or
claim that hostile ambient kit makes it behaviorally red.

Only migrate an exported call when the audited fixture reaches a real kit
consumer. Any behavioral mutation claim must name that kit-sensitive fixture
and demonstrate its focused hostile-ambient red. Record the source branch that
justifies every reviewed call that remains exported.

Make `newAbsentToolFixture` call `phaseTable(fixture.root, fixture.root)` so
the fixture root remains the single source of its kit identity. Its manifest is
present, so `phaseTable` does not observe the kit argument; this is source
cleanup enforced by an exact audit, not a behavior change under hostile
ambient kit. Dead-wrapper and unused-carrier cleanup belongs to the dependent
`repair-dead-kit-wrappers.md` and `repair-unused-prospective-kit-seam.md`
tickets.

## Acceptance

- [ ] [RAK1] (covers local) audited no-module and prospective calls remain exported with their empty-scoping source branches recorded, any behavior migration reaches a real kit consumer, and the absent-tool fixture passes its own root directly with the redundant `kitRoot` read rejected by an exact source audit, repairing RP1/P2/S2.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RAK1/no-module-exported | route one no-module fixture call through a private kit-taking seam | the source-ticket audit against `component_decision.go` | apply, confirm the fixture has no regular `go.mod`, and require the audit to name the migration as inert because scoping returns before phase-table resolution |
| RAK1/prospective-exported | route one prospective fixture call through a private kit-taking seam | the source-ticket audit against `evaluation.go` | apply, confirm the call uses prospective evaluation, and require the audit to name the migration as inert because prospective scoping returns empty before kit consumption |
| RAK1/absent-tool-root | restore `phaseTable(fixture.root, kitRoot(fixture.root))` | the exact source audit for `newAbsentToolFixture` | apply, run the audit, expect it to name the redundant ambient `kitRoot` read, and restore `phaseTable(fixture.root, fixture.root)` exactly |
