# Reconcile late untracked descendants

Blocked by: reconcile-authorized-snapshot-after-publication.md
Ownership fence: `internal/landing/landing.go`, `internal/landing/landing_test.go`
Integration surfaces: named directory snapshot and post-authorization worktree reconciliation→`internal/landing/landing.go`; deterministic late-descendant race proof→`internal/landing/landing_test.go` plus LU1; unnamed-state preservation→`internal/landing/landing_test.go` plus LU2
Contracts: untracked descendants attributed to a named directory cross prospective authorization→`internal/landing/landing.go` reconciliation, asserted by LU1-LU2 against a held-authorization race in `internal/landing/landing_test.go`

## What to build

Close the remaining directory form of the authorization race. If an untracked
descendant appears under a named directory after composition but before a winning
destination CAS, a successful return must not leave that descendant outside the
authorized and published tree. Reconcile only the named path and preserve every
unnamed staged, unstaged, untracked, and ignored path. Retain the existing
landed-but-checkout-incomplete result if exact selective reconciliation fails.

## Acceptance

- [ ] [LU1] A held authorization followed by creation of a new untracked descendant under a named directory publishes the authorized tree and returns with the named directory's index and worktree clean against the landed commit.
- [ ] [LU2] Removing the late attributed descendant does not change unnamed staged, unstaged, untracked, or ignored state.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| LU1 | restore the authorized tracked entries without removing a later untracked descendant | deterministic named-directory authorization-race test | run `go test ./internal/landing -run '^TestLandReconcilesLateUntrackedDescendantUnderNamedDirectory$' -count=1`; expect the late descendant or dirty named status assertion to fail |
| LU2 | clean the whole checkout instead of the literal named path | deterministic named-directory authorization-race test | run `go test ./internal/landing -run '^TestLandReconcilesLateUntrackedDescendantUnderNamedDirectory$' -count=1`; expect an unnamed staged, unstaged, untracked, or ignored snapshot assertion to fail |
