# Share the stack owner rule

Blocked by: none
Writes: internal/landing/stackpolicy/stackpolicy.go (new), internal/landing/stackpolicy/stackpolicy_test.go (new), .bench/commitment.json, internal/landing/stackpolicy/purity_census_test.go (new)
Covers: LS1, LS2, LS3, LS4, LS5, LS6

## What to build

Add the pure policy package `internal/landing/stackpolicy`. It imports only the standard library and the settle policy in `internal/landing/settlepolicy`, and its purity census grades it as a policy package.

The package exports two functions:

- `Owners` takes the entries of the destination tree and the entries of each member tree, in member order. It returns, for each path whose entry differs between the destination and a member, the 1-based indexes of each member that changes that path.
- `Overlaps` takes that answer and returns each path outside `capture/` that two or more members change, with those member indexes.

An entry is a path and one metadata string, such as the mode and the object id. Two entries differ when their metadata differs, or when one side has no entry. So a mode-only change counts, and a rename gives its member two paths with no rename detection. The caller passes every tree in one metadata format.

Ticket 02 calls `Owners` from the gate proof, and ticket 10 calls `Owners` and `Overlaps` from the stack preflight. Both callers read this one rule.

`.bench/commitment.json` is on this line only because the spec-stage close commits the delivery binding. The author does not edit or read it.

## Acceptance

- [ ] `Owners` maps a path that two members change to both member indexes, and a path that one member changes to that member alone (LS1).
- [ ] `Overlaps` leaves out a `capture/` path that two members change (LS2).
- [ ] `Overlaps` keeps a path outside `capture/` that two members change (LS3).
- [ ] A member that changes only the mode of one path owns that path (LS4).
- [ ] A member that renames one path owns the old path and the new path (LS5).
- [ ] A changed path that holds a space and a `*` maps to its exact spelling (LS6).

## Checkpoint verification

- `bench test --package ./internal/landing/stackpolicy`
- `bench test --check root-conformance`

## Mutation proof

Make `Owners` keep only the last member that changes a path. The LS1 test must fail. Spell the probe against the code that the author writes, and require `bench probe` to report `bit` with a byte-exact restore.
