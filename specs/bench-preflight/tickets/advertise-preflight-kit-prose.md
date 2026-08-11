# Advertise preflight across kit prose, anchors, and conformance

Blocked by: implement-preflight-review.md, harden-preflight-bootstrap-errors.md, implement-preflight-build.md, repair-routing-checker.md
Ownership fence: `internal/conformance/subcommand_routing_test.go`, `internal/anchors/registry_data.go`, `tests/canary/workflow-guidance-anchors/`, `.bench/BENCH.md`, `bin/bench.sh`, `.agents/commands/bench-implement-spec.md`, `.agents/commands/bench-review-implementation.md`, `projects/benchkit.md`, `CHANGELOG.md`
Integration surfaces: dispatched verb→implement-preflight-review.md's `commandRegistry` row (exercised by D1); build-mode verb→implement-preflight-build.md (named in the implement entry sentence, exercised by D2); repaired routing check→repair-routing-checker.md; cold-pickup doc check→existing `internal/conformance/docs_workflow_checks_test.go` `checkColdPickupCLILists` exercised by D1; anchor family fixtures→`tests/canary/workflow-guidance-anchors/`; anchor rows→`internal/anchors/registry_data.go`
Contracts: the two pinned phase-entry sentences cross `.agents/commands/bench-implement-spec.md` and `.agents/commands/bench-review-implementation.md`→`internal/anchors/registry_data.go`, asserted by D2 against the real command files via the direct root-conformance invocation; the `preflight` routing row crosses `internal/conformance/subcommand_routing_test.go`→the repaired check's real-tree comparison against `commandRegistry`, asserted by D1; the changelog line crosses `CHANGELOG.md`→its pinning anchor row in `internal/anchors/registry_data.go`, asserted by D3
Closure: D1/routing-row, D1/grammar-reach, D1/cold-pickup-doc-line, D2/implement-entry-anchor, D2/review-entry-anchor, D3/changelog-line-anchored

## What to build

The advertisement half of the feature, landed as one green change so the
docs-currency sweep sees code and docs together: the `"preflight":
routed("internal/preflight")` row in the `subcommandRouting` map
(`internal/conformance/subcommand_routing_test.go` — the map is the routing
registry; `internal/conformance/registry/registry.go` registers checks, not
commands, and is not touched); the `.bench/BENCH.md` CLI-inventory line and
the `bin/bench.sh` help line; the phase-entry step in each of the two phase
commands — implement's entry names `bench preflight build <slug>`, review's
entry names `bench preflight review <slug>`, each carrying the clause "a red
preflight stops the phase"; three anchor rows in
`internal/anchors/registry_data.go`: one per phase command pinning exactly its
entry sentence, and one pinning the `CHANGELOG.md` line naming the command
(the Spec A precedent at `registry_data.go:280` is the family convention —
unconditional, not optional). Add fixtures under
`tests/canary/workflow-guidance-anchors/` mirroring an existing fixture
directory if the family's bite coverage requires them; the
`projects/benchkit.md` AXI seam bullet gains the command. Each conformance red
this ticket owns is demonstrated red-then-green via `BENCH_CONFORMANCE_ROOT=<root>
go test -count=1 -run '^TestRootConformance$' ./internal/conformance` (other
pre-existing ship-tier reds are out of fence and stay).

## Acceptance

- [ ] [D1] (covers PF18) the routing registry carries the `preflight` row, the package's AST reaches `usage.Parse`, and `checkColdPickupCLILists` passes with the new `.bench/BENCH.md` inventory line — routing, grammar-reach, and cold-pickup each report no preflight violation under the direct root-conformance invocation.
- [ ] [D2] (covers PF19) each phase command's entry step carries its exact preflight sentence with the stop clause, and one anchor row per file pins exactly that sentence — anchors report no violation for the two new rows.
- [ ] [D3] (covers PF20) `CHANGELOG.md` carries one line naming `bench preflight`, pinned by its own anchor row.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| D1/routing-row | remove the `preflight` row from the `subcommandRouting` map | the repaired routing check | direct root-conformance invocation, expect the dispatched-but-unregistered violation |
| D1/grammar-reach | remove the `usage.Parse` call from the preflight package (transient probe, reverted) | the grammar-reach predicate | direct root-conformance invocation, expect the does-not-reach-grammar violation |
| D1/cold-pickup-doc-line | remove the `.bench/BENCH.md` inventory line | `checkColdPickupCLILists` | direct root-conformance invocation, expect the undocumented-route violation |
| D2/implement-entry-anchor | reword the pinned sentence in `bench-implement-spec.md` | the new implement-entry anchor row | direct root-conformance invocation, expect the anchor violation |
| D2/review-entry-anchor | reword the pinned sentence in `bench-review-implementation.md` | the new review-entry anchor row | direct root-conformance invocation, expect the anchor violation |
| D3/changelog-line-anchored | drop the changelog line | the new changelog anchor row | direct root-conformance invocation, expect the anchor violation |
