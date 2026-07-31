# Plan and apply abandonment

Blocked by: Resume idempotent mutations

Ownership fence: `internal/specbuild`
Assumptions: lifecycle recovery state and worktree ownership are durable and queryable

## What to build

Make abandon read-only by default: inventory active worktrees, provisional refs,
unintegrated checkpoints, and recovery refs into one deterministic fingerprint.
Apply only that exact plan, preserve every Git-visible unlanded change through
the existing recovery owner before owned cleanup, and retain terminal evidence.

## Acceptance

- [ ] [R26] Abandon plan performs no mutation; apply revalidates its fingerprint, refuses all drift without partial cleanup, preserves unlanded work before releasing only owned worktrees, and marks the run terminal while retaining evidence.
