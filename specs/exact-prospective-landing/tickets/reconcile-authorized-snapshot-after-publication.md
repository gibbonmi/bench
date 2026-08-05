# Reconcile the authorized snapshot after publication

Blocked by: adopt-exact-landing-in-commit.md, preserve-executable-spec-mode.md
Ownership fence: `internal/landing/landing.go`, `internal/landing/landing_test.go`
Integration surfaces: composed named-path snapshot -> prospective authorization; winning landing commit -> invoking index and worktree reconciliation; transitioned spec bytes and mode -> post-publication checkout
Contracts: `internal/landing/landing.go` reconciles every named path from the exact authorized/published tree rather than re-reading live bytes after authorization, while preserving unnamed index/worktree state and the staged spec's captured filesystem permissions; `internal/landing/landing_test.go` holds authorization open, mutates named and spec paths, and proves a successful return remains clean against the published commit

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

- [ ] [MAS1] Re-running `git add -A` from the live named path after publication makes the named-path authorization-race test red.
- [ ] [MAS2] Re-reading and transforming the live spec during reconciliation makes the spec authorization-race test red.
