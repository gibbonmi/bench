# Migrate outline outcomes

Blocked by: migrate-output-adapter.md
Ownership fence: `internal/outline`
Integration surfaces: returned-output adapter→migrate-output-adapter.md; shared outcome carrier→`internal/axi` exercised by OUT1; compatibility oracle→`cmd/bench/axi_compatibility_test.go` exercised by OUT1; final contraction→contract-outcome-action-routes.md
Contracts: the result kind, exact exit integer, and payload bytes cross `internal/outline`'s `Command`→the shared outcome carrier, membership is the seven result classes `outline.Command` returns (usage, not-in-repo, `git ls-files` failure, a report carrying skip rows, a report truncated at `bounds.OutlineRowLimit`, the `--full` report, and the clean default report), ordering is scan-carry-render, and absence is a zero-row `symbols` table at exit 0, asserted by OUT1 against the real `outline.Command`
Closure: OUT1/usage, OUT1/not-in-repo, OUT1/lsfiles-error, OUT1/skips, OUT1/truncated, OUT1/full, OUT1/ok

## What to build

`outline.Command` constructs a shared outcome carrying its own kind and exit for each of its seven result classes before the existing renderer emits. The counting facts stay exactly where they are — `totalSymbols`, `emitted`, `omitted`, the skip rows, and `bounds.OutlineRowLimit` are untouched, because bounded projection is `axi-bounded-projection-migration`'s scope; this ticket only moves the kind and exit onto the shared carrier and proves the rendered bytes are unchanged.

`migrate-structure-outcome.md` is OA3's accountable first claimant, under the tie-break rule that accountability goes to the ticket covering the spec map's first-enumerated OA3 family (structure, before models, testreport, guards, outline, roadmap, status, handoff, dashboard, and worktree query), so this ticket claims OA3 as defense in depth.

Refresh precondition (OA3 is not TDD-able until OA1 and the prerequisite carrier land): `internal/axi` must exist as Go package `axi` with `go doc ./internal/axi Outcome` resolving, and `migrate-output-adapter.md` must have landed. `cmd/bench/axi_compatibility_test.go` is an unchanged input this ticket only runs.

Write the owner test `TestOutlineOutcomeRouteByResultClass` in `internal/outline/outline_test.go`: one subtest per result class driving the real `outline.Command` over a `t.TempDir()` git fixture and asserting the constructed outcome's kind, exit, and carried counts.

## Acceptance

- [ ] [OUT1] (covers OA3) `outline.Command` constructs a shared outcome with its own kind and exit for each of its seven result classes, with byte-identical stdout and exits.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| OUT1/usage | return the `usage.Parse` line and code directly on the usage branch without constructing the outcome | `TestOutlineOutcomeRouteByResultClass/usage` | apply the mutation, run `go test ./internal/outline -run TestOutlineOutcomeRouteByResultClass/usage -timeout 60s`; the kind assertion fails with `outcome kind = "" want "usage"`; the usage branch reads no files, bounded by the `-timeout 60s` deadline |
| OUT1/not-in-repo | return `toon.NotInRepo()` and 1 directly when `git.Root()` errors, skipping the outcome | `TestOutlineOutcomeRouteByResultClass/not-in-repo` | apply the mutation, run `go test ./internal/outline -run TestOutlineOutcomeRouteByResultClass/not-in-repo -timeout 60s`; the kind assertion fails with `outcome kind = "" want "not-in-repo"`; the `git rev-parse` child is bounded by `bounds.TestDeadline(bounds.GitRefreshTimeout)` |
| OUT1/lsfiles-error | return `toon.Errorf("git ls-files failed", err.Error())` and 1 directly, skipping the outcome | `TestOutlineOutcomeRouteByResultClass/lsfiles-error` | apply the mutation, run `go test ./internal/outline -run TestOutlineOutcomeRouteByResultClass/lsfiles-error -timeout 60s` with a stub `git` exiting nonzero for `ls-files`; the kind assertion fails with `outcome kind = "" want "git-failed"`; the stub returns immediately under `bounds.TestDeadline(bounds.GitRefreshTimeout)` |
| OUT1/skips | return the skip-bearing report and 0 directly, skipping the outcome, whenever `len(skips) > 0` | `TestOutlineOutcomeRouteByResultClass/skips` | apply the mutation, run `go test ./internal/outline -run TestOutlineOutcomeRouteByResultClass/skips -timeout 60s` against a fixture holding an oversized file, a binary file, and a dangling symlink; the kind assertion fails with `outcome kind = "" want "outline-partial"` while the `skips` rows stay byte-identical; every read is capped at `bounds.OutlineFileLimit`, so no read can stall |
| OUT1/truncated | return the row-capped report and 0 directly, skipping the outcome, when `emitted == bounds.OutlineRowLimit` | `TestOutlineOutcomeRouteByResultClass/truncated` | apply the mutation, run `go test ./internal/outline -run TestOutlineOutcomeRouteByResultClass/truncated -timeout 60s` against a fixture emitting `bounds.OutlineRowLimit + 25` symbols; the kind assertion fails with `outcome kind = "" want "outline-partial"` while the emitted rows stay byte-identical; reads capped at `bounds.OutlineFileLimit`, bounded by the `-timeout 60s` deadline |
| OUT1/full | give the `--full` report the same kind as the capped default report instead of its own complete kind | `TestOutlineOutcomeRouteByResultClass/full` | apply the mutation, run `go test ./internal/outline -run TestOutlineOutcomeRouteByResultClass/full -timeout 60s` on the same oversized fixture; the kind assertion fails with `outcome kind = "outline-partial" want "outline"`; bounded by the `-timeout 60s` deadline |
| OUT1/ok | return the rendered report and 0 directly on the clean branch, skipping the outcome | `TestOutlineOutcomeRouteByResultClass/ok` | apply the mutation, run `go test ./internal/outline -run TestOutlineOutcomeRouteByResultClass/ok -timeout 60s`; the kind assertion fails with `outcome kind = "" want "outline"`; bounded by the `-timeout 60s` deadline |
