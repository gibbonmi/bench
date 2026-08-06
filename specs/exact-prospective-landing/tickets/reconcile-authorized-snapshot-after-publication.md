# Reconcile the authorized snapshot after publication

Blocked by: adopt-exact-landing-in-commit.md, preserve-executable-spec-mode.md
Ownership fence: `internal/landing/landing.go`, `internal/landing/landing_test.go`
Integration surfaces: composed named-path and spec snapshots plus post-publication reconciliation→`internal/landing/landing.go`; named-path authorization-race proof→`internal/landing/landing_test.go` plus AS1; transitioned-spec authorization-race proof→`internal/landing/landing_test.go` plus AS2
Contracts: authorized named bytes, transitioned spec bytes, Git mode, and filesystem permissions cross composition→`internal/landing/landing.go` reconciliation, asserted by AS1-AS2 against the real authorization-race proofs in `internal/landing/landing_test.go`

## What to build

Bind successful post-publication reconciliation to the same immutable snapshot that
was composed, authorized, and committed. A named path changed while authorization
runs must not replace the authorized bytes in the invoking index or worktree after
the destination CAS wins. Apply the same rule to the staged-to-implemented spec
transition: do not re-read or re-transform a later live spec version, and retain the
captured filesystem permissions while the index mode matches the published tree.
Continue to preserve every unnamed staged, unstaged, untracked, and ignored path,
and retain the existing landed-but-checkout-incomplete error if exact selective
reconciliation fails.

## Acceptance

- [ ] [AS1] If a named path changes during prospective authorization, a winning landing publishes the authorized bytes and returns with that path's index and worktree content clean against the landed commit while unnamed state is unchanged.
- [ ] [AS2] If the staged spec changes during authorization, a winning landing restores the authorized implemented bytes, published Git mode, and captured filesystem permissions without re-transforming the later live version.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AS1 | re-run `git add -A` from the live named path after publication | `TestLandReconcilesAuthorizedNamedSnapshotAfterAuthorizationMutation` | run `go test ./internal/landing -run '^TestLandReconcilesAuthorizedNamedSnapshotAfterAuthorizationMutation$' -count=1`; expect the named-path authorization-race proof to fail |
| AS2 | re-read and transform the live spec during reconciliation | `TestLandReconcilesAuthorizedSpecSnapshotAfterAuthorizationMutation` | run `go test ./internal/landing -run '^TestLandReconcilesAuthorizedSpecSnapshotAfterAuthorizationMutation$' -count=1`; expect the spec authorization-race proof to fail |
