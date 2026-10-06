# 1. Move the Writes grammar to one owner in internal/tickets

Blocked by: none
Writes: internal/tickets/writes.go (new), internal/tickets/writes_test.go (new), internal/tickets/registry_data.go, internal/preflight/decision.go, internal/preflight/closure.go, internal/preflight/fence_writes.go, internal/preflight/proposal.go, internal/commitment/repository/candidate.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: LP27, LP28, LP29, LP30, LP31, LP32, LP33

## What to build

This ticket is the first ticket of review chunk LP-C1. It delivers one owner for
the `Writes:` split and cover rules. Ticket 2 calls this owner from the
light-path predicate, so preflight and the commitment check read one rule.

First, read `splitWritesEntry` in `internal/preflight/decision.go`, `pathCovered`
in `internal/preflight/closure.go`, `prefixCovers` in
`internal/tickets/registry_data.go`, and `inScope` in
`internal/commitment/repository/candidate.go`. Confirm that each one applies the
same `(new)` split or the same `/` segment rule.

Add the new file `internal/tickets/writes.go` with two exported functions.
`WritesPath(entry)` returns the tree path that one `Writes:` entry names, and it
reports whether the entry carries `(new)`. It removes a trailing `/`.
`Covers(entry, path)` reports whether an entry that is already split names path
exactly, or contains path at a `/` segment boundary.

Retire `splitWritesEntry` and `pathCovered`. Each caller in `decision.go`,
`closure.go`, `fence_writes.go`, and `proposal.go` calls `WritesPath` and
`Covers` instead. A caller that holds a raw entry calls `WritesPath` first.
`prefixCovers` delegates to `Covers`, or `BoundFiles` calls `Covers` directly.

`inScope` calls `tickets.Covers` for each scope entry. This adds the import edge
`internal/commitment/repository` to `internal/tickets`. Confirm with
`go list -deps ./internal/tickets` that the edge is acyclic.

No behavior changes. The preflight rows and the legacy scope check keep their
current results.

## Acceptance

- [ ] `TestWritesEntryCover` in `internal/tickets/writes_test.go` shows that `tickets.WritesPath` returns `internal/x.go` and true for the entry `internal/x.go (new)` (LP27).
- [ ] The same test shows that `tickets.WritesPath` returns `internal/d` for the entry `internal/d/` (LP28).
- [ ] The same test shows that `tickets.Covers` reports true for the entry `internal/d` and the path `internal/d/a.go` (LP29).
- [ ] The same test shows that `tickets.Covers` reports false for the entry `internal/d` and the path `internal/dx/a.go` (LP30).
- [ ] `TestWritesResolveAcceptsNewMarker` in `internal/preflight/decision_test.go` stays green over a `(new)` entry that the tree lacks (LP31).
- [ ] `TestCommandFencePrefixBoundary` in `internal/preflight/command_review_test.go` stays green (LP32).
- [ ] `rg` over `internal/preflight`, `internal/tickets/registry_data.go`, and `internal/commitment/repository/candidate.go` finds no second copy of the `(new)` split or the segment rule (LP33, review-owned).
- [ ] `bench test --package ./internal/tickets`, `bench test --package ./internal/preflight`, and `bench test --package ./internal/commitment/repository` pass.
