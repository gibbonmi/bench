# Migrate maps outcomes

Blocked by: migrate-output-adapter.md
Ownership fence: `internal/maps`
Integration surfaces: returned-output adapter→migrate-output-adapter.md; shared outcome carrier→`internal/axi` exercised by MAP1; compatibility oracle→`cmd/bench/axi_compatibility_test.go` exercised by MAP1; final contraction→contract-outcome-action-routes.md
Contracts: the result kind, exact exit integer, and payload bytes cross `internal/maps`'s `Command`→the shared outcome carrier, membership is the eight result classes `maps.Command` returns (usage, `--count`/`--template` mutual exclusion, template, count, not-in-repo, record-error, invalid-row table, ok table), ordering is owner-carry-render, and absence is the zero-row `maps` table at exit 0, asserted by MAP1 against the real `maps.Command`
Closure: MAP1/usage, MAP1/mutually-exclusive, MAP1/template, MAP1/count, MAP1/not-in-repo, MAP1/record-error, MAP1/invalid-row, MAP1/ok

## What to build

`maps.Command` constructs a shared outcome carrying its own kind and exit for each of its eight result classes before the existing renderer runs. Stdout bytes and exit codes are unchanged, including the two distinct not-in-repo behaviours (`--count` prints `0` at exit 0; the bare form prints `toon.NotInRepo()` at exit 1).

`migrate-learnings-outcome.md` is OA2's accountable first claimant, under the tie-break rule that accountability goes to the ticket covering the spec map's first-enumerated OA2 family (learnings, before maps, diff, and coverage), so this ticket claims OA2 as defense in depth.

Refresh precondition (OA2 is not TDD-able until OA1 and the prerequisite carrier land): `internal/axi` must exist as Go package `axi` with `go doc ./internal/axi Outcome` resolving, and `migrate-output-adapter.md` must have landed. `cmd/bench/axi_compatibility_test.go` is an unchanged input this ticket only runs.

Write the owner test `TestMapsOutcomeRouteByResultClass` in `internal/maps/maps_test.go`: one subtest per result class driving the real `maps.Command` over a `t.TempDir()` repo and asserting the constructed outcome's kind and exit.

## Acceptance

- [ ] [MAP1] (covers OA2) `maps.Command` constructs a shared outcome with its own kind and exit for each of its eight result classes, with byte-identical stdout and exits.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| MAP1/usage | return the `usage.Parse` line and code directly on the usage branch without constructing the outcome | `TestMapsOutcomeRouteByResultClass/usage` | apply the mutation, run `go test ./internal/maps -run TestMapsOutcomeRouteByResultClass/usage -timeout 60s`; the kind assertion fails with `outcome kind = "" want "usage"`; pure in-process call, bounded by the `-timeout 60s` deadline |
| MAP1/mutually-exclusive | return the `--count and --template are mutually exclusive` help line and 2 directly, skipping the outcome | `TestMapsOutcomeRouteByResultClass/mutually-exclusive` | apply the mutation, run `go test ./internal/maps -run TestMapsOutcomeRouteByResultClass/mutually-exclusive -timeout 60s`; the kind assertion fails with `outcome kind = "" want "usage"` while exit stays 2; in-process, bounded by the `-timeout 60s` deadline |
| MAP1/template | return `DecisionMapTemplate()` and 0 directly, skipping the outcome | `TestMapsOutcomeRouteByResultClass/template` | apply the mutation, run `go test ./internal/maps -run TestMapsOutcomeRouteByResultClass/template -timeout 60s`; the kind assertion fails with `outcome kind = "" want "template"`; in-process with no filesystem read, bounded by the `-timeout 60s` deadline |
| MAP1/count | return `strconv.Itoa(s.count)+"\n"` and 0 directly on the `--count` branch, skipping the outcome | `TestMapsOutcomeRouteByResultClass/count` | apply the mutation, run `go test ./internal/maps -run TestMapsOutcomeRouteByResultClass/count -timeout 60s`; the kind assertion fails with `outcome kind = "" want "count"`; the decisions-tree scan reads a fixed `t.TempDir()` fixture, bounded by the `-timeout 60s` deadline |
| MAP1/not-in-repo | return `toon.NotInRepo()` and 1 directly when `git.Root()` errors, skipping the outcome | `TestMapsOutcomeRouteByResultClass/not-in-repo` | apply the mutation, run `go test ./internal/maps -run TestMapsOutcomeRouteByResultClass/not-in-repo -timeout 60s`; the kind assertion fails with `outcome kind = "" want "not-in-repo"`; the `git rev-parse` child is bounded by `bounds.TestDeadline(bounds.GitRefreshTimeout)` in the fixture |
| MAP1/record-error | return `toon.RecordError(...)` and 1 directly on `s.state.Failed()`, skipping the outcome | `TestMapsOutcomeRouteByResultClass/record-error` | apply the mutation, run `go test ./internal/maps -run TestMapsOutcomeRouteByResultClass/record-error -timeout 60s`; the kind assertion fails with `outcome kind = "" want "record-error"`; the classifier reads at most `bounds.ControlRecordLimit`, so the read cannot stall |
| MAP1/invalid-row | return the rendered table and 1 directly on the `row[3] == "invalid"` branch, skipping the outcome | `TestMapsOutcomeRouteByResultClass/invalid-row` | apply the mutation, run `go test ./internal/maps -run TestMapsOutcomeRouteByResultClass/invalid-row -timeout 60s`; the kind assertion fails with `outcome kind = "" want "invalid"` while exit stays 1; in-process, bounded by the `-timeout 60s` deadline |
| MAP1/ok | return the rendered table and 0 directly on the clean branch, skipping the outcome | `TestMapsOutcomeRouteByResultClass/ok` | apply the mutation, run `go test ./internal/maps -run TestMapsOutcomeRouteByResultClass/ok -timeout 60s`; the kind assertion fails with `outcome kind = "" want "maps"`; in-process, bounded by the `-timeout 60s` deadline |
