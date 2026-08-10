# Migrate coverage outcomes

Blocked by: migrate-output-adapter.md
Ownership fence: `internal/coverage`
Integration surfaces: returned-output adapter→migrate-output-adapter.md; shared outcome carrier→`internal/axi` exercised by COV1; compatibility oracle→`cmd/bench/axi_compatibility_test.go` exercised by COV1; final contraction→contract-outcome-action-routes.md
Contracts: the result kind, exact exit integer, and payload bytes cross `internal/coverage`'s `Command`→the shared outcome carrier, membership is the seven result classes `coverage.Command` returns (usage, missing `<spec.md>` argument, spec unreadable, spec not found, `--check` pass on a mapped map, `--check` pass on a historical marker, `--check` violations, and the default report), ordering is owner-carry-render, and absence is the zero-row `rows` table at exit 0, asserted by COV1 against the real `coverage.Command`
Closure: COV1/usage, COV1/missing-arg, COV1/unreadable, COV1/not-found, COV1/check-mapped, COV1/check-historical, COV1/check-violations, COV1/report

## What to build

`coverage.Command` constructs a shared outcome carrying its own kind and exit for each of its result classes before the existing renderer runs. Stdout bytes and exit codes are unchanged, including the two distinct `--check` pass lines (`coverage map valid — N row(s)` and `coverage map historical — validation skipped`).

`migrate-learnings-outcome.md` is OA2's accountable first claimant, under the tie-break rule that accountability goes to the ticket covering the spec map's first-enumerated OA2 family (learnings, before maps, diff, and coverage), so this ticket claims OA2 as defense in depth.

Refresh precondition (OA2 is not TDD-able until OA1 and the prerequisite carrier land): `internal/axi` must exist as Go package `axi` with `go doc ./internal/axi Outcome` resolving, and `migrate-output-adapter.md` must have landed. `cmd/bench/axi_compatibility_test.go` is an unchanged input this ticket only runs.

Write the owner test `TestCoverageOutcomeRouteByResultClass` in `internal/coverage/coverage_test.go`: one subtest per result class driving the real `coverage.Command` over `t.TempDir()` spec fixtures and asserting the constructed outcome's kind and exit.

## Acceptance

- [ ] [COV1] (covers OA2) `coverage.Command` constructs a shared outcome with its own kind and exit for each of its result classes, with byte-identical stdout and exits.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| COV1/usage | return the `usage.Parse` line and code directly on the usage branch without constructing the outcome | `TestCoverageOutcomeRouteByResultClass/usage` | apply the mutation, run `go test ./internal/coverage -run TestCoverageOutcomeRouteByResultClass/usage -timeout 60s`; the kind assertion fails with `outcome kind = "" want "usage"`; pure in-process call, bounded by the `-timeout 60s` deadline |
| COV1/missing-arg | return `toon.MissingArg("bench coverage", ...)` and 2 directly, skipping the outcome | `TestCoverageOutcomeRouteByResultClass/missing-arg` | apply the mutation, run `go test ./internal/coverage -run TestCoverageOutcomeRouteByResultClass/missing-arg -timeout 60s`; the kind assertion fails with `outcome kind = "" want "usage"` while exit stays 2; in-process, bounded by the `-timeout 60s` deadline |
| COV1/unreadable | return `toon.Errorf("spec not readable: ...", ...)` and 1 directly, skipping the outcome | `TestCoverageOutcomeRouteByResultClass/unreadable` | apply the mutation, run `go test ./internal/coverage -run TestCoverageOutcomeRouteByResultClass/unreadable -timeout 60s` against a mode-0000 fixture spec; the kind assertion fails with `outcome kind = "" want "spec-unreadable"`; a single `os.ReadFile` with no wait, bounded by the `-timeout 60s` deadline |
| COV1/not-found | return `toon.Errorf("spec not found: ...", ...)` and 1 directly, skipping the outcome | `TestCoverageOutcomeRouteByResultClass/not-found` | apply the mutation, run `go test ./internal/coverage -run TestCoverageOutcomeRouteByResultClass/not-found -timeout 60s`; the kind assertion fails with `outcome kind = "" want "spec-not-found"`; in-process, bounded by the `-timeout 60s` deadline |
| COV1/check-mapped | return `checkOKLine(fmt.Sprintf("coverage map valid — %d row(s)", ...))` and 0 directly, skipping the outcome | `TestCoverageOutcomeRouteByResultClass/check-mapped` | apply the mutation, run `go test ./internal/coverage -run TestCoverageOutcomeRouteByResultClass/check-mapped -timeout 60s`; the kind assertion fails with `outcome kind = "" want "check-pass"` while the `coverage map valid — 3 row(s)` line is byte-identical; in-process, bounded by the `-timeout 60s` deadline |
| COV1/check-historical | return `checkOKLine("coverage map historical — validation skipped")` and 0 directly, skipping the outcome | `TestCoverageOutcomeRouteByResultClass/check-historical` | apply the mutation, run `go test ./internal/coverage -run TestCoverageOutcomeRouteByResultClass/check-historical -timeout 60s`; the kind assertion fails with `outcome kind = "" want "check-skipped"`; in-process, bounded by the `-timeout 60s` deadline |
| COV1/check-violations | return the accumulated violation lines and 1 directly, skipping the outcome | `TestCoverageOutcomeRouteByResultClass/check-violations` | apply the mutation, run `go test ./internal/coverage -run TestCoverageOutcomeRouteByResultClass/check-violations -timeout 60s`; the kind assertion fails with `outcome kind = "" want "check-failed"` while exit stays 1; in-process, bounded by the `-timeout 60s` deadline |
| COV1/report | return the assembled `spec:`/`state:`/`rows` report and 0 directly, skipping the outcome | `TestCoverageOutcomeRouteByResultClass/report` | apply the mutation, run `go test ./internal/coverage -run TestCoverageOutcomeRouteByResultClass/report -timeout 60s`; the kind assertion fails with `outcome kind = "" want "coverage"`; in-process, bounded by the `-timeout 60s` deadline |
