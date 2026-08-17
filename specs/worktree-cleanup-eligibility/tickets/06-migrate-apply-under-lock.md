# Migrate apply under the eligibility verdict

Blocked by: 04-migrate-automatic-shared-landedness.md
Writes: internal/worktree/resume.go, internal/worktree/resume_test.go, internal/worktree/recovery_retry_test.go, internal/worktree/worktree_test.go, specs/worktree-cleanup-eligibility/tickets/06-migrate-apply-under-lock.md

## What to build

Make explicit and automatic apply consume only a freshly replanned eligible
verdict under the existing transaction lock. Preserve the current stale-fingerprint,
interruption, recovery, checkpoint, and terminal-receipt behavior, including the
rule that a retain verdict cannot execute. This ticket establishes the common
mutation-boundary contract used by the release and landing migration.

## Acceptance

- [ ] CO2: explicit and automatic apply replan under the transaction lock and execute only an eligible verdict with unchanged stale-fingerprint and interruption outcomes.
