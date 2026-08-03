# Render worktree age in the status listing

Blocked by: none
Ownership fence: `internal/status`, `internal/worktree/registry.go`
Assumptions: the registry record already carries a creation time; age is computed at render time and stored nowhere; the clock is injected so the age is deterministic under test; claims re-derived from the tree at pickup

## What to build

Users see each active worktree's age beside its path in the status listing, so a
stale worktree is visible at a glance.

## Acceptance

- [ ] [WA1] the registry exposes each active worktree's creation time to the status listing.
- [ ] [WA2] status renders each active worktree's age beside its path.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| WA1 | drop the creation time from the exposed registration record | the registry record test | zero the field, run `go test ./internal/worktree`, expect the missing-creation-time failure |
| WA2 | render the worktree row with the path alone | the worktree-row render test | omit the age column, run `go test ./internal/status`, expect the missing-age failure |
