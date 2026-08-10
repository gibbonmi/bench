# Migrate TOON and query empty classes

Blocked by: none
Ownership fence: `internal/toon`, `cmd/bench/empty_classes.go`, `cmd/bench/empty_classes_test.go`
Integration surfaces: typed empty carrier→`internal/axi/empty.go` exercised by TE1; zero-row table owner→`internal/toon/toon.go` exercised by TE1; query dispatch used to drive each member→`cmd/bench/command_registry.go` exercised by TE1; registry empty declaration→declare-empty-dispositions.md; legacy carrier contraction→contract-aggregate-empty-routes.md
Contracts: the zero-row table class and the absent refusal class cross `internal/toon/toon.go`→`internal/axi/empty.go`; the zero-row class is exactly the `name[0]{fields}:` header with its declared field list and exit 0, and the refusal class is `toon.NotInRepo()` with exit 1; order is the command's existing block order; a zero-row success may never render as a refusal and a refusal may never render as a zero-row success, asserted by TE1 against the real query commands driven through `Command.Run`
Closure: TE1/zero-row-encoder, TE1/learnings-empty, TE1/maps-empty, TE1/diff-clean, TE1/outline-empty, TE1/worktree-list-empty, TE1/spec-history-empty, TE1/absent-refusal, TE1/route

## What to build

The shared TOON zero-row table becomes a declared empty class carried by the shared typed
empty value, and each query member that reaches it keeps reaching exactly that class with
its own field list. The six members enumerated here are the ones the FT173 surface census
names as independent reachers of the zero table under principle 5 — absent learnings, no
active maps, a clean diff, no worktrees, no outline symbols, and no spec history — plus the
encoder contract itself. Rendered bytes do not change; what changes is that the class is
declared once and observed per member.

`internal/toon.Table` stays the sole owner of the `name[0]{fields}:` rendering. The per-member
observations are authored in `cmd/bench/empty_classes_test.go`, which drives each member end
to end through `Command{}.Run` against a `t.TempDir()` git repository, so no domain package
is written by this ticket.

Tree condition that must hold when this ticket is refreshed: `internal/axi/empty.go` exists
and declares the exported empty classification type `EmptyClass` with a distinct constant for
the zero-row table class and one for the absent refusal. If that path or the symbol is
absent, stop and report rather than build — the prerequisite `axi-carriers-and-registry`
build has not landed.

## Acceptance

- [ ] [TE1] (covers AE8) the zero-row TOON table is one declared empty class that each of the six enumerated query members reaches with its own field list, distinct from the absent-source refusal.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| TE1/zero-row-encoder | in `internal/toon.table`, return the empty string for a nil row slice instead of the `name[0]{fields}:` header | `TestTable` (`internal/toon`) | run `go test ./internal/toon -run TestTable -count=1 -timeout 60s`; the zero-row case's exact-output comparison fails against `""`; the encoder is pure over in-memory rows, so there is nothing that can stall |
| TE1/learnings-empty | classify the absent-learnings result as the absent refusal instead of the zero-row table class | `TestQueryMembersReachTheZeroRowEmptyClass` (`cmd/bench`, new) | run `go test ./cmd/bench -run 'TestQueryMembersReachTheZeroRowEmptyClass/learnings' -count=1 -timeout 180s`; the assertion that `bench learnings` in a repository with no journal exits 0 with the `learnings[0]{...}` header fails with exit 1 and a structured error line; the subtest runs `Command{}.Run` in-process against a `t.TempDir()` git repository, so no subprocess is spawned |
| TE1/maps-empty | classify the no-active-maps result as the absent refusal instead of the zero-row table class | `TestQueryMembersReachTheZeroRowEmptyClass` (`cmd/bench`, new) | run `go test ./cmd/bench -run 'TestQueryMembersReachTheZeroRowEmptyClass/maps' -count=1 -timeout 180s`; the assertion that `bench maps` with no active decision map exits 0 with the `maps[0]{...}` header fails with exit 1; in-process `Command{}.Run` over a `t.TempDir()` repository |
| TE1/diff-clean | classify the clean-diff result as the absent refusal instead of the zero-row table class | `TestQueryMembersReachTheZeroRowEmptyClass` (`cmd/bench`, new) | run `go test ./cmd/bench -run 'TestQueryMembersReachTheZeroRowEmptyClass/diff' -count=1 -timeout 180s`; the assertion that `bench diff` on a committed clean tree exits 0 with the `diff[0]{...}` header fails with exit 1; the fixture makes only local `add`/`commit` git calls |
| TE1/outline-empty | classify the no-symbols result as the absent refusal instead of the zero-row table class | `TestQueryMembersReachTheZeroRowEmptyClass` (`cmd/bench`, new) | run `go test ./cmd/bench -run 'TestQueryMembersReachTheZeroRowEmptyClass/outline' -count=1 -timeout 180s`; the assertion that `bench outline` over a repository with no Go symbols exits 0 with the `outline[0]{...}` header fails with exit 1; each fixture file read is capped at `bounds.OutlineFileLimit` |
| TE1/worktree-list-empty | classify the no-worktrees result as the absent refusal instead of the zero-row table class | `TestQueryMembersReachTheZeroRowEmptyClass` (`cmd/bench`, new) | run `go test ./cmd/bench -run 'TestQueryMembersReachTheZeroRowEmptyClass/worktree-list' -count=1 -timeout 180s`; the assertion that `bench worktree list` with no registered worktree exits 0 with the `worktrees[0]{...}` header fails with exit 1; only a local `git worktree list` runs |
| TE1/spec-history-empty | classify the no-history result as the absent refusal instead of the zero-row table class | `TestQueryMembersReachTheZeroRowEmptyClass` (`cmd/bench`, new) | run `go test ./cmd/bench -run 'TestQueryMembersReachTheZeroRowEmptyClass/spec-history' -count=1 -timeout 180s`; the assertion that `bench spec history` for a spec with no recorded history exits 0 with its zero-row header fails with exit 1; only a local `git log` over a single-commit fixture runs |
| TE1/absent-refusal | classify the outside-a-repository result as the zero-row table class instead of the absent refusal | `TestQueryMembersKeepTheAbsentRefusalDistinct` (`cmd/bench`, new) | run `go test ./cmd/bench -run TestQueryMembersKeepTheAbsentRefusalDistinct -count=1 -timeout 180s`; the assertion that each of the six members exits 1 with the `toon.NotInRepo()` line outside a repository fails with exit 0 and a zero-row table — success invented from a failed read; each subtest chdirs into a non-repository `t.TempDir()` and `git.Root()` returns immediately |
| TE1/route | keep the bare `toon.Table` calls and never construct the shared typed empty value for a zero-row result | `TestQueryMembersReachTheZeroRowEmptyClass` (`cmd/bench`, new) | run `go test ./cmd/bench -run 'TestQueryMembersReachTheZeroRowEmptyClass/route' -count=1 -timeout 180s`; the assertion that each member's zero-row result carried the zero-row `axi.EmptyClass` constant fails with no classification observed, even though every rendered byte is unchanged; in-process `Command{}.Run` over `t.TempDir()` fixtures |
