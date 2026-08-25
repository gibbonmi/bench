# Serialize the package-global stub swaps

Line: opus / medium.

A parallel test that swaps a package-level injectable races every parallel
reader, and the census predicate cannot see the swap. The repair adds one
harness edge, so the mid tier applies.

Blocked by: 12-take-explicit-home-in-every-verb.md
Writes: internal/worktree/journey_test.go, internal/worktree/parallel_census_test.go, internal/worktree/worktree_test.go, internal/worktree/lifecycle_test.go, internal/worktree/live_binary_test.go

## What to build

`TestIgnoredInventoryStatRaceRetains` swaps the package-level `ignoredLstat`
while parallel readers run through `inventoryIgnored`, and `go test -race`
reports it. The harness gains `bindGlobal(t, name)`, which records a
`global` effect and makes the test serial by construction. The census's
serial helpers gain `bindGlobal` beside `bindEnv` and `chdir`. A test that
swaps a package-level injectable and whose readers hold no lock calls
`bindGlobal` first and stays serial. A test whose readers all hold a
test-local lock, as in `lifecycle_test.go` and `live_binary_test.go`, stays
parallel; the lock is the edge there.

## Acceptance

- [ ] A synthetic test that calls `bindGlobal` is serial under the census, and a serial test that calls it and `t.Parallel()` is reported.
- [ ] `go test -race -count=1 -run 'TestIgnoredInventoryStatRaceRetains|TestLandingDestination' ./internal/worktree` passes ten times in a row.
- [ ] WF01 stays green on the live tree.
