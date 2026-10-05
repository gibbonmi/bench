# 2. Admit an unbound light-path commit inside its ticket's Writes line

Blocked by: 1-own-writes-grammar.md
Writes: internal/commitment/repository/light_path.go (new), internal/commitment/repository/light_path_test.go (new), internal/commitment/repository/candidate.go, internal/commitment/repository/readiness.go, internal/commitment/commitmenttest/, internal/commit/commitment_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: LP1, LP2, LP7, LP8, LP9, LP10, LP11, LP12, LP13, LP14, LP15, LP16, LP17, LP18, LP19, LP20, LP21, LP22, LP23, LP49, LP51, LP56

## What to build

This ticket is the second ticket of review chunk LP-C1. It delivers commit
mode: `bench commit` admits an unbound light-path change whose production paths
stay inside its one ticket's `Writes:` line. It calls `tickets.WritesPath` and
`tickets.Covers` from ticket 1, and it adds no second copy of those rules.

Add the new file `internal/commitment/repository/light_path.go`. Its tree reader
runs one `git ls-tree -r -z` of `specs` in the read tree, after the precedent in
`repository.go`. From that one listing it finds each tickets-only folder, the
`tickets.Ext` entries at every depth below `tickets/`, and their modes. It
confirms each folder through `spec.TicketsOnly(spec.CommitTree(root, tree), slug)`.
It reads the one ticket through `git.ReadTreeFile` and parses it through
`tickets.ParseTicket` with the empty tag.

A folder qualifies under the four conditions of the spec section
`The light-path predicate`. An entry with another extension than `tickets.Ext`
is an asset, and the reader ignores it.

Add one unexported predicate. `authorizeCandidate` calls it only after `readyFor`
refuses with the binding fault. Each other `readyFor` error returns unchanged.
If `readyFor` gives no distinct binding fault today, give it one in
`readiness.go` and keep its refusal text. When `readyFor` admits, the predicate
reads no tree. The transition check, `protectedCandidate`, the production loop,
and the continuation scope check stay before `readyFor`.

In commit mode the predicate reads every tickets-only folder of the candidate
tree. It admits when one qualifying ticket covers every path in `production`.
Otherwise it returns the Writes refusal or the span refusal of the spec section
`Refusal words`, word for word, with each operand through Go `%q`. The Writes
refusal wins over the span refusal. With no qualifying folder, the `readyFor`
refusal returns unchanged.

Add a helper in `internal/commitment/commitmenttest` that writes a grammatical
one-ticket folder with a `Writes:` line. `WriteTickets` keeps its current bytes,
so its fixtures never qualify. The `Store.AuthorizeCandidate` rows follow
`TestAdmitPublicationFrozenIdentity`. The commit rows follow
`TestCommitmentCommitBeforeEffects` through `planningCommitRepo`.

## Acceptance

- [ ] `TestCommitmentLightPathCommit` in `internal/commit/commitment_test.go` shows that `bench commit -- change.go` in an unbound worktree exits 0 when the one ticket lists `change.go` (LP1).
- [ ] The same test shows that `bench commit -- other.go` exits 1, and stderr holds `"other.go"`, the ticket path, and `bench commitment start` (LP2).
- [ ] `TestLightPathCandidate` in `internal/commitment/repository/light_path_test.go` shows that `Store.AuthorizeCandidate` admits a light-path tree while a second assignment holds the active binding and claim (LP7).
- [ ] The same test shows admission of a new path under a `(new)` entry (LP8), and of `pkg/a.go` under the entry `pkg` (LP9).
- [ ] The same test shows that `pkgx/a.go` under the entry `pkg` refuses with `"pkgx/a.go"` in the refusal (LP10).
- [ ] The same test shows admission of a tree that deletes `reviews/x.md` and adds `docs/adr/9999-x.md` when the ticket lists both (LP11).
- [ ] The same test shows that `tickets/one.md` beside a second top-level ticket (LP12), or beside `tickets/sub/two.md` (LP51), refuses with `assignment has no current delivery binding`.
- [ ] The same test shows that two one-ticket folders that each cover one of two paths refuse with `more than one light-path ticket` (LP13).
- [ ] The same test shows that two qualifying tickets that leave one path outside both refuse with that path named (LP56).
- [ ] The same test shows that a folder approved as a deliverable refuses with `assignment has no current delivery binding` (LP14).
- [ ] The same test shows the same refusal for a ticket with no `Writes:` field (LP15), and for a ticket entry of mode `120000` (LP16).
- [ ] The same test shows that a light-path tree that edits `.bench/commitment.json` refuses with `candidate policy has no exact approval` (LP17).
- [ ] The same test shows that a light-path tree that deletes a pinned row's detail file refuses with `candidate changes protected commitment` (LP18).
- [ ] The same test shows that a listed continuation whose ticket covers a path outside its scope refuses with `legacy continuation scope excludes` (LP19).
- [ ] The same test shows admission of a tree that removes an unpinned row's index line and detail file beside a covered path (LP20).
- [ ] The same test shows admission of one ticket beside `tickets/asset.txt` (LP21), and refusal with `assignment has no current delivery binding` of `tickets/asset.txt` alone (LP22).
- [ ] The same test shows that an uncovered path with a space and an ESC byte appears in its Go-quoted spelling in the refusal (LP23).
- [ ] The same test shows admission of a bound assignment's path beside a qualifying folder that does not cover it (LP49).
- [ ] `bench test --package ./internal/commitment/repository` and `bench test --package ./internal/commit` pass.
