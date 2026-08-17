# Migrate landed-set planning

Blocked by: 03-expand-typed-eligibility-verdict.md
Writes: internal/worktree/clean_landed.go, internal/worktree/clean_landed_apply_test.go, specs/worktree-cleanup-eligibility/tickets/05-migrate-landed-set-planning.md

## What to build

Move only the `--landed` set planner's preservation refusal onto the eligibility
owner. Keep its membership selection, set fingerprint, and no-op rows unchanged;
the selector may continue to consume the typed landedness contract introduced by
the expansion ticket without waiting for automatic cleanup's consumer migration.

## Acceptance

- [x] CO3: `--landed` planning obtains its preservation refusal from the eligibility owner while its set fingerprint and no-op rows stay unchanged.
