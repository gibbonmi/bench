# Preserve symbolic-ref targets

Blocked by: none
Writes: internal/git/git.go, internal/git/refs_test.go, CHANGELOG.md
Covers: none

## What to build

Delete the named branch ref without following a symbolic ref to its target.
Use the existing exact-delete helper for the landing prune.
Keep the old-object check that refuses a moved branch.

## Acceptance

- [ ] The landing prune removes a landed symbolic branch and preserves its target.
- [ ] Exact deletion preserves a symbolic branch's target.
- [ ] Exact deletion refuses a branch whose object changed after the plan.
