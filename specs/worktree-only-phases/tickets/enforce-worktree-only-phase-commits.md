# Enforce worktree-only phase commits

Blocked by: none
Writes: .bench/BENCH.md, CHANGELOG.md, internal/git, internal/status/status.go, internal/commit, internal/anchors

## What to build

Bench must refuse phase publication from the primary checkout. `bench commit`
accepts linked worktrees and directs primary-checkout users to `bench worktree create`.

The operating guide must describe this executable boundary. Status and commit
must use one primary-checkout classifier.

## Acceptance

- [ ] `bench commit` exits 1 before publication from the primary checkout.
- [ ] The refusal directs users to `bench worktree create`.
- [ ] `bench commit` keeps its existing behavior in a linked worktree.
- [ ] Status and commit use one primary-checkout classifier.
- [ ] The operating guide states the enforced publication boundary.
