# Make the landing fixture's ticket Writes union-exact

Blocked by: none
Writes: internal/worktree/land_fixtures_test.go
Covers: none

## What to build

The public landing fixture in `internal/worktree` writes a spec whose fence names
`owned.txt`, the review pickup, and the sibling review path. Its prepared ticket
keeps the `recordtest` default `Writes: source.txt`. The fence therefore differs
from the union of the ticket `Writes:` paths. A landing that grades the
`fence-writes` row would go red on the fixture itself. This ticket sets the prepared
ticket's `Writes:` line to the fence less its review pickup, through the existing
`SetTicketWrites` helper.

## Acceptance

- [ ] The fixture's ticket `Writes:` line names `owned.txt` and the sibling review path, and no other path.
- [ ] The fixture reads its landing base after the `Writes:` commit.
- [ ] The `internal/worktree` tests pass with no change to a test expectation.
