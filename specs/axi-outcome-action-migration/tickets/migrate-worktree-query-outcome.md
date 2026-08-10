# Migrate the worktree list query outcome

Blocked by: migrate-output-adapter.md
Ownership fence: `internal/worktree/list.go`, `internal/worktree/list_test.go`
Integration surfaces: returned-output adapter→migrate-output-adapter.md; shared outcome carrier→`internal/axi` exercised by WQ1; compatibility oracle→`cmd/bench/axi_compatibility_test.go` exercised by WQ1; lifecycle route→migrate-worktree-lifecycle-outcomes-actions.md; final contraction→contract-outcome-action-routes.md
Contracts: the result kind, exact exit integer, and payload bytes cross `internal/worktree/list.go`'s `ListCommand`→the shared outcome carrier, membership is the six result classes `ListCommand` returns (help at exit 0, argument refusal at exit 2, not-in-repo, unreadable assignments, unreadable registrations, and the rendered `worktrees` table), ordering is read-carry-render, and absence is the zero-row `worktrees` table at exit 0, asserted by WQ1 against the real `worktree.ListCommand`
Closure: WQ1/help, WQ1/usage, WQ1/not-in-repo, WQ1/assignments-unreadable, WQ1/registrations-unreadable, WQ1/ok, WQ1/empty

## What to build

`worktree.ListCommand` — the AXI-approved nested `worktree list` query — constructs a shared outcome carrying its own kind and exit for each of its result classes before `toon.TableTyped` renders the `worktrees` table. This ticket writes only `internal/worktree/list.go` and its test: the lifecycle families are `migrate-worktree-lifecycle-outcomes-actions.md`, which is blocked on this one. Splitting them is not a keep-together case — `bench worktree list` is a read-only query with no fingerprint, no durable state, and no action tokens, so it lands green with the lifecycle writers untouched, and no project-gate red is stranded by the cut.

`migrate-structure-outcome.md` is OA3's accountable first claimant, under the tie-break rule that accountability goes to the ticket covering the spec map's first-enumerated OA3 family (structure, before models, testreport, guards, outline, roadmap, status, handoff, dashboard, and worktree query), so this ticket claims OA3 as defense in depth.

Refresh precondition (OA3 is not TDD-able until OA1 and the prerequisite carrier land): `internal/axi` must exist as Go package `axi` with `go doc ./internal/axi Outcome` resolving, and `migrate-output-adapter.md` must have landed. `cmd/bench/axi_compatibility_test.go` is an unchanged input this ticket only runs.

Write the owner test `TestWorktreeListOutcomeRouteByResultClass` in `internal/worktree/list_test.go`: one subtest per closure token driving the real `ListCommand` over a `t.TempDir()` git fixture and asserting the constructed outcome's kind and exit.

## Acceptance

- [ ] [WQ1] (covers OA3) `worktree.ListCommand` constructs a shared outcome with its own kind and exit for each of its result classes, with byte-identical stdout and exits.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| WQ1/help | return `"usage: " + usage.WorktreeList + "\n"` and 0 directly for `-h`/`--help`, skipping the outcome | `TestWorktreeListOutcomeRouteByResultClass/help` | apply the mutation, run `go test ./internal/worktree -run TestWorktreeListOutcomeRouteByResultClass/help -timeout 60s`; the kind assertion fails with `outcome kind = "" want "help"` while the usage line stays byte-identical; the branch does no IO, bounded by the `-timeout 60s` deadline |
| WQ1/usage | return `toon.Usage(usage.WorktreeList, args[0])` and 2 directly, skipping the outcome | `TestWorktreeListOutcomeRouteByResultClass/usage` | apply the mutation, run `go test ./internal/worktree -run TestWorktreeListOutcomeRouteByResultClass/usage -timeout 60s`; the kind assertion fails with `outcome kind = "" want "usage"` while exit stays 2; no IO on this branch, bounded by the `-timeout 60s` deadline |
| WQ1/not-in-repo | return `toon.NotInRepo()` and 1 directly when `git.Root()` errors, skipping the outcome | `TestWorktreeListOutcomeRouteByResultClass/not-in-repo` | apply the mutation, run `go test ./internal/worktree -run TestWorktreeListOutcomeRouteByResultClass/not-in-repo -timeout 60s`; the kind assertion fails with `outcome kind = "" want "not-in-repo"`; the `git rev-parse` child is bounded by `bounds.TestDeadline(bounds.GitRefreshTimeout)` |
| WQ1/assignments-unreadable | return `toon.Errorf("cannot read worktree assignments", ...)` and 1 directly, skipping the outcome | `TestWorktreeListOutcomeRouteByResultClass/assignments-unreadable` | apply the mutation, run `go test ./internal/worktree -run TestWorktreeListOutcomeRouteByResultClass/assignments-unreadable -timeout 60s` against a fixture whose intent ledger is malformed JSON; the kind assertion fails with `outcome kind = "" want "assignments-unreadable"`; a single bounded ledger read, under `bounds.TestDeadline(bounds.GitRefreshTimeout)` |
| WQ1/registrations-unreadable | return `toon.Errorf("cannot read registered worktrees", ...)` and 1 directly, skipping the outcome | `TestWorktreeListOutcomeRouteByResultClass/registrations-unreadable` | apply the mutation, run `go test ./internal/worktree -run TestWorktreeListOutcomeRouteByResultClass/registrations-unreadable -timeout 60s` with a stub `git` exiting nonzero for `worktree list`; the kind assertion fails with `outcome kind = "" want "registrations-unreadable"`; the stub returns immediately under `bounds.TestDeadline(bounds.GitRefreshTimeout)` |
| WQ1/ok | return the rendered `worktrees` table and 0 directly, skipping the outcome | `TestWorktreeListOutcomeRouteByResultClass/ok` | apply the mutation, run `go test ./internal/worktree -run TestWorktreeListOutcomeRouteByResultClass/ok -timeout 60s` against a fixture with one assignment row and one foreign row; the kind assertion fails with `outcome kind = "" want "worktrees"` while all eight `worktreeListFields` cells stay identical; git children bounded by `bounds.TestDeadline(bounds.GitRefreshTimeout)` |
| WQ1/empty | give the zero-registration case the not-in-repo kind and exit 1 instead of carrying the query kind at exit 0 | `TestWorktreeListOutcomeRouteByResultClass/empty` | apply the mutation, run `go test ./internal/worktree -run TestWorktreeListOutcomeRouteByResultClass/empty -timeout 60s` against a fixture with no assignments and no foreign registrations; the exit assertion fails with `exit = 1, want 0` while the zero-row `worktrees` header stays byte-identical; git children bounded by `bounds.TestDeadline(bounds.GitRefreshTimeout)` |
