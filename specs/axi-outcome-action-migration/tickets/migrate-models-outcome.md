# Migrate models outcomes

Blocked by: migrate-output-adapter.md
Ownership fence: `internal/models`
Integration surfaces: returned-output adapter→migrate-output-adapter.md; shared outcome carrier→`internal/axi` exercised by MOD1; compatibility oracle→`cmd/bench/axi_compatibility_test.go` exercised by MOD1; final contraction→contract-outcome-action-routes.md
Contracts: the result kind, exact exit integer, and payload bytes cross `internal/models`'s `Command`→the shared outcome carrier, membership is the five result classes `models.Command` returns (argument misuse at exit 2, help spelling at exit 0, render failure at exit 1, inventory with reachable providers at exit 0, inventory with unavailable provider rows at exit 0), ordering is owner-carry-render, and absence is a zero-row inventory at exit 0, asserted by MOD1 against the real `models.Command`
Closure: MOD1/misuse, MOD1/help, MOD1/render-error, MOD1/ok, MOD1/unavailable

## What to build

`models.Command` constructs a shared outcome carrying its own kind and exit for each of its five result classes before the existing `render` output is returned. The discovery tolerance stays intact: an unreachable provider still yields `unavailable` rows at exit 0 and must not collapse into the render-failure kind.

`migrate-structure-outcome.md` is OA3's accountable first claimant, under the tie-break rule that accountability goes to the ticket covering the spec map's first-enumerated OA3 family (structure, before models, testreport, guards, outline, roadmap, status, handoff, dashboard, and worktree query), so this ticket claims OA3 as defense in depth.

Refresh precondition (OA3 is not TDD-able until OA1 and the prerequisite carrier land): `internal/axi` must exist as Go package `axi` with `go doc ./internal/axi Outcome` resolving, and `migrate-output-adapter.md` must have landed. `cmd/bench/axi_compatibility_test.go` is an unchanged input this ticket only runs.

Write the owner test `TestModelsOutcomeRouteByResultClass` in `internal/models/models_test.go`: one subtest per result class driving the real `models.Command` and asserting the constructed outcome's kind and exit.

## Acceptance

- [ ] [MOD1] (covers OA3) `models.Command` constructs a shared outcome with its own kind and exit for each of its five result classes, with byte-identical stdout and exits.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| MOD1/misuse | return the `usage.Parse` line and 2 directly for a positional argument, skipping the outcome | `TestModelsOutcomeRouteByResultClass/misuse` | apply the mutation, run `go test ./internal/models -run TestModelsOutcomeRouteByResultClass/misuse -timeout 60s`; the kind assertion fails with `outcome kind = "" want "usage"` while exit stays 2; pure in-process call, bounded by the `-timeout 60s` deadline |
| MOD1/help | return the `usage.Parse` help line and 0 directly for `--help`, skipping the outcome | `TestModelsOutcomeRouteByResultClass/help` | apply the mutation, run `go test ./internal/models -run TestModelsOutcomeRouteByResultClass/help -timeout 60s`; the kind assertion fails with `outcome kind = "" want "help"` while exit stays 0; in-process, bounded by the `-timeout 60s` deadline |
| MOD1/render-error | return `toon.RenderError(err)` and 1 directly, skipping the outcome | `TestModelsOutcomeRouteByResultClass/render-error` | apply the mutation, run `go test ./internal/models -run TestModelsOutcomeRouteByResultClass/render-error -timeout 60s` with a control-bearing fixture model id that `toon` refuses; the kind assertion fails with `outcome kind = "" want "render-failed"`; in-process, bounded by the `-timeout 60s` deadline |
| MOD1/ok | return the rendered inventory and 0 directly, skipping the outcome | `TestModelsOutcomeRouteByResultClass/ok` | apply the mutation, run `go test ./internal/models -run TestModelsOutcomeRouteByResultClass/ok -timeout 60s`; the kind assertion fails with `outcome kind = "" want "models"`; provider discovery is bounded by `bounds.ProviderTimeout` per source and the fixture uses stubbed sources |
| MOD1/unavailable | give the unavailable-provider path the render-failure kind and exit 1 instead of carrying the inventory kind at exit 0 | `TestModelsOutcomeRouteByResultClass/unavailable` | apply the mutation, run `go test ./internal/models -run TestModelsOutcomeRouteByResultClass/unavailable -timeout 60s` with an unreachable stub provider; the exit assertion fails with `exit = 1, want 0` and the kind assertion with `outcome kind = "render-failed" want "models"`; the unreachable probe is bounded by `bounds.ProviderTimeout` |
