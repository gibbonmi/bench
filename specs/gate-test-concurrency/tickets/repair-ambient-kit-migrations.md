# Repair the remaining ambient-kit migrations

Blocked by: none
Ownership fence: `internal/gate/absent_tool_slot_test.go`, `internal/gate/gate_test.go`, `internal/gate/verdict_reuse_test.go`
Integration surfaces: `internal/gate/gate.go` exported entries→`internal/gate/gate_test.go` + RAK1 and `internal/gate/verdict_reuse_test.go` + RAK1; `internal/gate/engine.go` exported entry→existing `internal/gate/evaluation.go` empty prospective scoping + RAK1; `internal/gate/absent_tool_slot_test.go` fixture table→RAK1
Contracts: a parallel fixture behavior test in `internal/gate/verdict_reuse_test.go` passes its fixture kit as an explicit path wherever the exercised execution can consume kit identity; a deliberate public-entry subject stays serial and pins its input; the absent-tool fixture derives its table from its own root
Closure: RAK1/exported-entry-audit, RAK1/absent-tool-root

## What to build

Audit every `t.Parallel` test in `verdict_reuse_test.go`,
`prospective_test.go`, and `gate_test.go` that calls `Execute`,
`ExecuteReusingFreshGreen`, `RunCommand`, or another environment-resolving
entry over a fixture or throwaway project root. Where kit identity can affect
the exercised path, route the test through the existing package-private
kit-taking seam with `root` and `kit` both set to that fixture root. A test
whose subject is entry-time environment resolution remains serial and pins
the input explicitly. Record why any reviewed call site remains exported.

Make `newAbsentToolFixture` derive its phase table from the fixture root
directly rather than re-reading ambient kit. Dead-wrapper and unused-carrier
cleanup belongs to the dependent `repair-unused-kit-seams.md` ticket.

## Acceptance

- [ ] [RAK1] (covers local) no parallel fixture behavior test relies on ambient kit where kit identity can affect its execution, and the absent-tool fixture has no ambient read, repairing RP1/P2/S2.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RAK1/exported-entry-audit | route one migrated kit-sensitive fixture test back through its exported entry | a hostile-ambient focused run | apply, run the focused test with `BENCH_KIT=/nonexistent`, expect red; restore and rerun green |
| RAK1/absent-tool-root | restore `phaseTable(fixture.root, kitRoot(fixture.root))` | the absent-tool focused family under hostile ambient | apply, run with `BENCH_KIT=/nonexistent`, expect red |
