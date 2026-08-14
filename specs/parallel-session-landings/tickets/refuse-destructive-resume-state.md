# Refuse destructive resume destination state

Blocked by: none
Writes: internal/worktree

## What to build

Before reconciliation, fail closed on every caller-owned destination state that
would be overwritten: staged changes, tracked-worktree changes, untracked
collisions, ignored residue, and nested repositories.

## Acceptance

- [ ] PL25 and PL30 resume journeys refuse every enumerated destructive state before reset.
- [ ] Destination refs, marker, index, worktree, and assignment remain unchanged on refusal.
