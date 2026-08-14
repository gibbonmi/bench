# Movement-check explicit-base preflight

Blocked by: none
Writes: internal/diff, internal/preflight

## What to build

Resolve explicit base, source tip, index, tracked worktree, and untracked state
inside one movement-checked attempt, retry once on drift, then refuse with the
existing snapshot-drift action.

## Acceptance

- [ ] HEAD, index, tracked-worktree, and untracked drift each trigger retry or bounded refusal.
- [ ] A converged retry reports one coherent base/tip/path snapshot without configuration writes.
- [ ] `AuthorizeReviewedSource` and `Gather` movement-check the supplied source root when it differs from the process working directory.
