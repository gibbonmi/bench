# Migrate diff outcomes

Blocked by: migrate-output-adapter.md
Ownership fence: `internal/diff`
Integration surfaces: returned-output adapter→migrate-output-adapter.md; shared outcome carrier→`internal/axi` exercised by DIF1; compatibility oracle→`cmd/bench/axi_compatibility_test.go` exercised by DIF1; final contraction→contract-outcome-action-routes.md
Contracts: the result kind, exact exit integer, and payload bytes cross `internal/diff`'s `Command`→the shared outcome carrier, membership is the six result classes `diff.Command` returns (usage, not-in-repo, unresolvable range, git-operation failure, default report, `--full` report), ordering is owner-carry-render, and absence is the zero-row `files` table inside the report at exit 0, asserted by DIF1 against the real `diff.Command`
Closure: DIF1/usage, DIF1/not-in-repo, DIF1/range-error, DIF1/git-op-error, DIF1/ok-default, DIF1/ok-full

## What to build

`diff.Command` constructs a shared outcome carrying its own kind and exit for each of its six result classes before the existing renderer emits the `branch:`/`base:`/`method:` preamble, the `files` table, and (under `--full`) the `log` table and `diff_body:` section. Stdout bytes and exit codes are unchanged.

`migrate-learnings-outcome.md` is OA2's accountable first claimant, under the tie-break rule that accountability goes to the ticket covering the spec map's first-enumerated OA2 family (learnings, before maps, diff, and coverage), so this ticket claims OA2 as defense in depth.

Refresh precondition (OA2 is not TDD-able until OA1 and the prerequisite carrier land): `internal/axi` must exist as Go package `axi` with `go doc ./internal/axi Outcome` resolving, and `migrate-output-adapter.md` must have landed. `cmd/bench/axi_compatibility_test.go` is an unchanged input this ticket only runs.

Write the owner test `TestDiffOutcomeRouteByResultClass` in `internal/diff/diff_test.go`: one subtest per result class driving the real `diff.Command` against a `t.TempDir()` git fixture and asserting the constructed outcome's kind and exit.

## Acceptance

- [ ] [DIF1] (covers OA2) `diff.Command` constructs a shared outcome with its own kind and exit for each of its six result classes, with byte-identical stdout and exits.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| DIF1/usage | return the `usage.Parse` line and code directly on the usage branch without constructing the outcome | `TestDiffOutcomeRouteByResultClass/usage` | apply the mutation, run `go test ./internal/diff -run TestDiffOutcomeRouteByResultClass/usage -timeout 60s`; the kind assertion fails with `outcome kind = "" want "usage"`; pure in-process call, bounded by the `-timeout 60s` deadline |
| DIF1/not-in-repo | return `toon.NotInRepo()` and 1 directly when `git.Root()` errors, skipping the outcome | `TestDiffOutcomeRouteByResultClass/not-in-repo` | apply the mutation, run `go test ./internal/diff -run TestDiffOutcomeRouteByResultClass/not-in-repo -timeout 60s`; the kind assertion fails with `outcome kind = "" want "not-in-repo"`; the `git rev-parse` child is bounded by `bounds.TestDeadline(bounds.GitRefreshTimeout)` in the fixture |
| DIF1/range-error | return `toon.Errorf(errKind, errHint)` and 1 directly when `resolveCommitRange` reports a kind, skipping the outcome | `TestDiffOutcomeRouteByResultClass/range-error` | apply the mutation, run `go test ./internal/diff -run TestDiffOutcomeRouteByResultClass/range-error -timeout 60s` against a fixture passing `--commit deadbeef`; the kind assertion fails with `outcome kind = "" want "range-unresolved"`; every `git rev-parse` child is bounded by `bounds.TestDeadline(bounds.GitRefreshTimeout)` |
| DIF1/git-op-error | return `toon.Errorf("git diff --name-status failed", ...)` and 1 directly, skipping the outcome | `TestDiffOutcomeRouteByResultClass/git-op-error` | apply the mutation, run `go test ./internal/diff -run TestDiffOutcomeRouteByResultClass/git-op-error -timeout 60s` with a stub `git` on PATH that exits nonzero; the kind assertion fails with `outcome kind = "" want "git-failed"`; the stub returns immediately, and the fixture context is bounded by `bounds.TestDeadline(bounds.GitRefreshTimeout)` |
| DIF1/ok-default | return the assembled default report and 0 directly, skipping the outcome | `TestDiffOutcomeRouteByResultClass/ok-default` | apply the mutation, run `go test ./internal/diff -run TestDiffOutcomeRouteByResultClass/ok-default -timeout 60s`; the kind assertion fails with `outcome kind = "" want "diff"`; git children bounded by `bounds.TestDeadline(bounds.GitRefreshTimeout)` |
| DIF1/ok-full | return the assembled `--full` report and 0 directly, skipping the outcome | `TestDiffOutcomeRouteByResultClass/ok-full` | apply the mutation, run `go test ./internal/diff -run TestDiffOutcomeRouteByResultClass/ok-full -timeout 60s`; the kind assertion fails with `outcome kind = "" want "diff-full"` while the `log` table and `diff_body:` bytes stay identical; git children bounded by `bounds.TestDeadline(bounds.GitRefreshTimeout)` |
