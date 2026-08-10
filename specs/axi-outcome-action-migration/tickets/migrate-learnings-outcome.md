# Migrate learnings outcomes

Blocked by: migrate-output-adapter.md
Ownership fence: `internal/learnings`
Integration surfaces: returned-output adapter→migrate-output-adapter.md; shared outcome carrier→`internal/axi` exercised by LEA1; compatibility oracle→`cmd/bench/axi_compatibility_test.go` exercised by LEA1; final contraction→contract-outcome-action-routes.md
Contracts: the result kind, exact exit integer, and payload bytes cross `internal/learnings`'s `Command`→the shared outcome carrier, membership is the six result classes `learnings.Command` returns (usage, not-in-repo, absent-journal empty table, record-error, malformed-heading rows, ok table), ordering is owner-carry-render, and absence is the zero-row `learnings` table at exit 0, asserted by LEA1 against the real `learnings.Command`
Closure: LEA1/usage, LEA1/not-in-repo, LEA1/absent-empty, LEA1/record-error, LEA1/malformed, LEA1/ok

## What to build

`learnings.Command` constructs a shared outcome carrying its own kind and exit for each of its six result classes before the existing `toon` renderer runs. Stdout bytes and exit codes are unchanged.

This ticket is OA2's accountable first claimant: the tie-break rule is that accountability goes to the ticket covering the spec map's first-enumerated OA2 family (learnings, before maps, diff, and coverage), and `migrate-maps-outcome.md`, `migrate-diff-outcome.md`, and `migrate-coverage-outcome.md` claim OA2 as defense in depth.

Refresh precondition (OA2 is not TDD-able until OA1 and the prerequisite carrier land): `internal/axi` must exist as Go package `axi` with `go doc ./internal/axi Outcome` resolving, and `migrate-output-adapter.md` must have landed so `cmd/bench/main.go`'s `outputCommand` consumes a carrier. `cmd/bench/axi_compatibility_test.go` is an unchanged input this ticket only runs.

Write the owner test `TestLearningsOutcomeRouteByResultClass` in `internal/learnings/learnings_test.go`: one subtest per result class, each driving the real `learnings.Command` over a `t.TempDir()` repo and asserting the constructed outcome's kind and exit. A byte-equal bypass is invisible to the oracle, so this route test is the one that must go red.

## Acceptance

- [ ] [LEA1] (covers OA2) `learnings.Command` constructs a shared outcome with its own kind and exit for each of its six result classes, with byte-identical stdout and exits.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| LEA1/usage | return the `usage.Parse` line and code directly on the usage branch without constructing the outcome | `TestLearningsOutcomeRouteByResultClass/usage` | apply the mutation, run `go test ./internal/learnings -run TestLearningsOutcomeRouteByResultClass/usage -timeout 60s`; the kind assertion fails with `outcome kind = "" want "usage"`; pure in-process call with no child or wait, bounded by the `-timeout 60s` deadline |
| LEA1/not-in-repo | return `toon.NotInRepo()` and 1 directly when `git.Root()` errors, skipping the outcome | `TestLearningsOutcomeRouteByResultClass/not-in-repo` | apply the mutation, run `go test ./internal/learnings -run TestLearningsOutcomeRouteByResultClass/not-in-repo -timeout 60s`; the kind assertion fails with `outcome kind = "" want "not-in-repo"`; the `git rev-parse` child is bounded by `bounds.TestDeadline(bounds.GitRefreshTimeout)` in the fixture |
| LEA1/absent-empty | return the zero-row `toon.Table` result directly on `bounds.StateAbsent`, skipping the outcome | `TestLearningsOutcomeRouteByResultClass/absent-empty` | apply the mutation, run `go test ./internal/learnings -run TestLearningsOutcomeRouteByResultClass/absent-empty -timeout 60s`; the kind assertion fails with `outcome kind = "" want "empty"` while the rendered zero-row table is unchanged; in-process, bounded by the `-timeout 60s` deadline |
| LEA1/record-error | return `toon.RecordError(...)` and 1 directly on the failed-classification branch, skipping the outcome | `TestLearningsOutcomeRouteByResultClass/record-error` | apply the mutation, run `go test ./internal/learnings -run TestLearningsOutcomeRouteByResultClass/record-error -timeout 60s`; the kind assertion fails with `outcome kind = "" want "record-error"` on the oversized-journal fixture; the classifier reads at most `bounds.ControlRecordLimit`, so the read cannot stall |
| LEA1/malformed | return the malformed-rows table and 1 directly, skipping the outcome on the malformed-heading branch | `TestLearningsOutcomeRouteByResultClass/malformed` | apply the mutation, run `go test ./internal/learnings -run TestLearningsOutcomeRouteByResultClass/malformed -timeout 60s`; the kind assertion fails with `outcome kind = "" want "malformed"` while exit stays 1; in-process, bounded by the `-timeout 60s` deadline |
| LEA1/ok | return the rendered table and 0 directly on the clean-parse branch, skipping the outcome | `TestLearningsOutcomeRouteByResultClass/ok` | apply the mutation, run `go test ./internal/learnings -run TestLearningsOutcomeRouteByResultClass/ok -timeout 60s`; the kind assertion fails with `outcome kind = "" want "learnings"`; in-process, bounded by the `-timeout 60s` deadline |
