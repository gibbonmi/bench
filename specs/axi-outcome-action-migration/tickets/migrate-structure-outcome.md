# Migrate structure outcomes

Blocked by: migrate-output-adapter.md
Ownership fence: `internal/structure`
Integration surfaces: returned-output adapter→migrate-output-adapter.md; shared outcome carrier→`internal/axi` exercised by STR1; compatibility oracle→`cmd/bench/axi_compatibility_test.go` exercised by STR1; final contraction→contract-outcome-action-routes.md
Contracts: the result kind, exact exit integer, payload bytes, and the stderr-only git-operation diagnostic cross `internal/structure`'s `Command`→the shared outcome carrier, membership is the six result classes `structure.Command` returns (usage, not-in-repo, `--since` git failure, `--since` report, whole-tree git failure, whole-tree report), ordering is owner-carry-render, and absence is the clean report at exit 0 from `exitOf(0)`, asserted by STR1 against the real `structure.Command`
Closure: STR1/usage, STR1/not-in-repo, STR1/since-git-error, STR1/since-report, STR1/all-git-error, STR1/all-report

## What to build

`structure.Command` constructs a shared outcome carrying its own kind and exit for each of its six result classes before returning. The two git-failure classes keep returning an empty stdout payload with the diagnostic on stderr, and `exitOf(violations)` remains the exit authority — the outcome carries that exit rather than deriving one.

This ticket is OA3's accountable first claimant: the tie-break rule is that accountability goes to the ticket covering the spec map's first-enumerated OA3 family (structure, before models, testreport, guards, outline, roadmap, status, handoff, dashboard, and worktree query), and every other OA3 claimant — `migrate-models-outcome.md`, `migrate-testreport-outcome.md`, `migrate-guards-outcome.md`, `migrate-outline-outcome.md`, `migrate-roadmap-outcome.md`, `migrate-status-outcome-action.md`, `migrate-handoff-outcome-action.md`, `migrate-dashboard-outcome-composition.md`, and `migrate-worktree-query-outcome.md` — claims OA3 as defense in depth.

Refresh precondition (OA3 is not TDD-able until OA1 and the prerequisite carrier land): `internal/axi` must exist as Go package `axi` with `go doc ./internal/axi Outcome` resolving, and `migrate-output-adapter.md` must have landed. `cmd/bench/axi_compatibility_test.go` is an unchanged input this ticket only runs.

Write the owner test `TestStructureOutcomeRouteByResultClass` in `internal/structure/structure_test.go`: one subtest per result class driving the real `structure.Command` over a `t.TempDir()` git fixture and asserting the constructed outcome's kind and exit.

## Acceptance

- [ ] [STR1] (covers OA3) `structure.Command` constructs a shared outcome with its own kind and exit for each of its six result classes, with byte-identical stdout, stderr, and exits.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| STR1/usage | return the `usage.Parse` line and code directly on the usage branch without constructing the outcome | `TestStructureOutcomeRouteByResultClass/usage` | apply the mutation, run `go test ./internal/structure -run TestStructureOutcomeRouteByResultClass/usage -timeout 60s`; the kind assertion fails with `outcome kind = "" want "usage"`; pure in-process call, bounded by the `-timeout 60s` deadline |
| STR1/not-in-repo | return `toon.NotInRepo()` and 1 directly when `git.Root()` errors, skipping the outcome | `TestStructureOutcomeRouteByResultClass/not-in-repo` | apply the mutation, run `go test ./internal/structure -run TestStructureOutcomeRouteByResultClass/not-in-repo -timeout 60s`; the kind assertion fails with `outcome kind = "" want "not-in-repo"`; the `git rev-parse` child is bounded by `bounds.TestDeadline(bounds.GitRefreshTimeout)` |
| STR1/since-git-error | return `"", 1` directly after printing `gitOpError("diff", terr)`, skipping the outcome | `TestStructureOutcomeRouteByResultClass/since-git-error` | apply the mutation, run `go test ./internal/structure -run TestStructureOutcomeRouteByResultClass/since-git-error -timeout 60s` with a stub `git` on PATH that exits nonzero for `diff`; the kind assertion fails with `outcome kind = "" want "git-failed"` while the stderr line is byte-identical; the stub returns immediately under `bounds.TestDeadline(bounds.GitRefreshTimeout)` |
| STR1/since-report | return `report, exitOf(violations)` directly on the `--since` branch, skipping the outcome | `TestStructureOutcomeRouteByResultClass/since-report` | apply the mutation, run `go test ./internal/structure -run TestStructureOutcomeRouteByResultClass/since-report -timeout 60s`; the kind assertion fails with `outcome kind = "" want "structure-touched"`; git children bounded by `bounds.TestDeadline(bounds.GitRefreshTimeout)` |
| STR1/all-git-error | return `"", 1` directly after printing `gitOpError("ls-files", cerr)`, skipping the outcome | `TestStructureOutcomeRouteByResultClass/all-git-error` | apply the mutation, run `go test ./internal/structure -run TestStructureOutcomeRouteByResultClass/all-git-error -timeout 60s` with a stub `git` that exits nonzero for `ls-files`; the kind assertion fails with `outcome kind = "" want "git-failed"`; the stub returns immediately under `bounds.TestDeadline(bounds.GitRefreshTimeout)` |
| STR1/all-report | make the whole-tree branch derive its exit from `report == ""` instead of carrying `exitOf(violations)` | `TestStructureOutcomeRouteByResultClass/all-report` | apply the mutation, run `go test ./internal/structure -run TestStructureOutcomeRouteByResultClass/all-report -timeout 60s` on a fixture with one violation; the exit assertion fails with `exit = 0, want 1`; git children bounded by `bounds.TestDeadline(bounds.GitRefreshTimeout)` |
