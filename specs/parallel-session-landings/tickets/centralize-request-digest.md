# Centralize request digest ownership

Blocked by: none
Writes: internal/intent, internal/worktree

## What to build

Expose one request-digest owner from the intent package and route assignment,
ownership, cleanup, and resume consumers through it.

## Acceptance

- [ ] Production request digests have one derivation.
- [ ] Assignment authentication and cleanup-receipt lookup retain their current identities.
