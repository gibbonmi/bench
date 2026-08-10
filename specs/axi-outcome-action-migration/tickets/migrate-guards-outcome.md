# Migrate guards outcomes

Blocked by: migrate-output-adapter.md
Ownership fence: `internal/guards`
Integration surfaces: returned-output adapter→migrate-output-adapter.md; shared outcome carrier→`internal/axi` exercised by GUA1; compatibility oracle→`cmd/bench/axi_compatibility_test.go` exercised by GUA1; final contraction→contract-outcome-action-routes.md
Contracts: the result kind, exact exit integer, and payload bytes cross `internal/guards`'s `Command`→the shared outcome carrier, membership is the six result classes `guards.Command` returns (usage, not-in-repo, `--brief` lines, full `guards`+`guard_scan` tables, a scan truncated by `guardScanTimeout`, and a `toon` render failure), ordering is scan-carry-render, and absence is a zero-row `guards` table with its `guard_scan` row still present, asserted by GUA1 against the real `guards.Command`
Closure: GUA1/usage, GUA1/not-in-repo, GUA1/brief, GUA1/full, GUA1/scan-bound, GUA1/render-error

## What to build

`guards.Command` constructs a shared outcome carrying its own kind and exit for each of its six result classes before the existing renderer emits. `ScanResult` stays the completeness authority and its `guard_scan` status/inspected/total/omitted/reason bytes are untouched — those are `axi-aggregate-empty-migration`'s scope. Every class still exits 0 except the render failure at 1.

`migrate-structure-outcome.md` is OA3's accountable first claimant, under the tie-break rule that accountability goes to the ticket covering the spec map's first-enumerated OA3 family (structure, before models, testreport, guards, outline, roadmap, status, handoff, dashboard, and worktree query), so this ticket claims OA3 as defense in depth.

Refresh precondition (OA3 is not TDD-able until OA1 and the prerequisite carrier land): `internal/axi` must exist as Go package `axi` with `go doc ./internal/axi Outcome` resolving, and `migrate-output-adapter.md` must have landed. `cmd/bench/axi_compatibility_test.go` is an unchanged input this ticket only runs.

Write the owner test `TestGuardsOutcomeRouteByResultClass` in `internal/guards/guards_test.go`: one subtest per result class driving the real `guards.Command` over a `t.TempDir()` git fixture and asserting the constructed outcome's kind and exit.

## Acceptance

- [ ] [GUA1] (covers OA3) `guards.Command` constructs a shared outcome with its own kind and exit for each of its six result classes, with byte-identical stdout and exits.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| GUA1/usage | return the `usage.Parse` line and code directly on the usage branch without constructing the outcome | `TestGuardsOutcomeRouteByResultClass/usage` | apply the mutation, run `go test ./internal/guards -run TestGuardsOutcomeRouteByResultClass/usage -timeout 60s`; the kind assertion fails with `outcome kind = "" want "usage"`; the usage branch runs no scan, bounded by the `-timeout 60s` deadline |
| GUA1/not-in-repo | return `toon.NotInRepo()` and 1 directly when `git.Root()` errors, skipping the outcome | `TestGuardsOutcomeRouteByResultClass/not-in-repo` | apply the mutation, run `go test ./internal/guards -run TestGuardsOutcomeRouteByResultClass/not-in-repo -timeout 60s`; the kind assertion fails with `outcome kind = "" want "not-in-repo"`; the `git rev-parse` child is bounded by `bounds.TestDeadline(bounds.GitRefreshTimeout)` |
| GUA1/brief | return the assembled `--brief` string and 0 directly, skipping the outcome | `TestGuardsOutcomeRouteByResultClass/brief` | apply the mutation, run `go test ./internal/guards -run TestGuardsOutcomeRouteByResultClass/brief -timeout 60s`; the kind assertion fails with `outcome kind = "" want "guards-brief"` while the `pre-push: ... [branch: ...]` lines stay byte-identical; the scan is bounded by `bounds.Context(ctx, guardScanTimeout)` (5s) inside `Command` |
| GUA1/full | return `out + meta` and 0 directly, skipping the outcome | `TestGuardsOutcomeRouteByResultClass/full` | apply the mutation, run `go test ./internal/guards -run TestGuardsOutcomeRouteByResultClass/full -timeout 60s`; the kind assertion fails with `outcome kind = "" want "guards"`; the scan is bounded by `bounds.Context(ctx, guardScanTimeout)` (5s) |
| GUA1/scan-bound | give the timeout-truncated scan the same outcome kind as a complete scan instead of its own partial kind | `TestGuardsOutcomeRouteByResultClass/scan-bound` | apply the mutation, run `go test ./internal/guards -run TestGuardsOutcomeRouteByResultClass/scan-bound -timeout 60s` against a fixture whose scan is truncated by the 5s bound; the kind assertion fails with `outcome kind = "guards" want "guards-partial"` while the `guard_scan` row stays byte-identical; the truncation is produced by that same `guardScanTimeout` bound, so the subtest cannot outlive it |
| GUA1/render-error | return `toon.RenderError(err)` and 1 directly from either `toon.Table` call, skipping the outcome | `TestGuardsOutcomeRouteByResultClass/render-error` | apply the mutation, run `go test ./internal/guards -run TestGuardsOutcomeRouteByResultClass/render-error -timeout 60s` with a control-bearing guard path `toon` refuses; the kind assertion fails with `outcome kind = "" want "render-failed"`; the scan is bounded by `guardScanTimeout` (5s) |
