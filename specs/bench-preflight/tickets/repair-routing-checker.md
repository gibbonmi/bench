# Repair the routing checker to read commandRegistry

Blocked by: none
Ownership fence: `internal/conformance/subcommand_routing_test.go`
Integration surfaces: dispatch source→existing `cmd/bench/main.go` `commandRegistry` (read, not written; exercised by R1's real-tree run); routing registry→existing `subcommandRouting` map in `internal/conformance/subcommand_routing_test.go` (unchanged rows, consumed as today); repaired check and its `preflight` row consumed by→advertise-preflight-kit-prose.md
Contracts: the set of dispatched subcommand names crosses `cmd/bench/main.go`→`internal/conformance/subcommand_routing_test.go`, asserted by R1 against the real `commandRegistry` composite literal via the direct root-conformance invocation
Closure: R1/live-registry-extraction, R1/missing-dispatch-red, R2/registry-bite, R2/routed-claim-bite

## What to build

`dispatchNames` in `internal/conformance/subcommand_routing_test.go` still
parses the legacy `commands` map and `run()` switch (constants at the file
top), both replaced by the `commandRegistry` composite literal in
`cmd/bench/main.go` — so the real-tree routing check reports every registered
name as "no longer dispatches". Re-derive dispatch names from
`commandRegistry`'s `Name:` fields by AST — matched by that identifier, not by
any composite literal — and update both fixture bite tests
(`TestSubcommandRoutingRegistryBites`, `TestSubcommandRoutingRoutedClaimBites`)
to fixtures of the registry shape. Weaken no predicate: the `subcommandRouting`
map's rows, the `packageReachesGrammar` / `usage.Parse` reach check, and every
other surviving assertion stay as they are. This ticket adds no `preflight`
row — that lands with the advertisement ticket once the verb dispatches.
`craft-gate` discipline applies — prove the repaired check still bites before
claiming it green.

## Acceptance

- [ ] [R1] (covers PF23) on today's tree, `BENCH_CONFORMANCE_ROOT=<root> go test -count=1 -run '^TestRootConformance$' ./internal/conformance` no longer emits "no longer dispatches" for any name in `commandRegistry`, and a registry row whose name is absent from `commandRegistry` still produces exactly that violation.
- [ ] [R2] (covers PF23) both fixture bite tests drive registry-shaped fixtures and stay green in the dev suite, each still failing when its seeded violation is present.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| R1/live-registry-extraction | match any composite literal instead of the `commandRegistry` identifier | a bite fixture carrying a decoy literal with different `Name:` fields | `go test ./internal/conformance -run 'SubcommandRouting'`, expect the wrong-source failure; author the decoy fixture as part of this ticket |
| R1/missing-dispatch-red | drop the registry-name-without-dispatch comparison from the check | the registry bite fixture seeding a registered-but-undispatched name | run the bite test, expect the missed-violation failure |
| R2/registry-bite | remove the seeded violation from `TestSubcommandRoutingRegistryBites`'s fixture | that bite test | run it, expect the expected-violation-not-found failure |
| R2/routed-claim-bite | remove the seeded violation from `TestSubcommandRoutingRoutedClaimBites`'s fixture | that bite test | run it, expect the expected-violation-not-found failure |
