# Record review and recover runs

Blocked by: Checkpoint and integrate attributed patches

## What to build

Extend `internal/specbuild` with exact-candidate three-axis review receipts,
full retained-evidence status, shared mutator preconditions, request-idempotent
fault recovery, and fingerprinted abandon plan/apply through the existing
worktree recovery owner. Candidate changes invalidate review, and accepted
repairs remain ordinary assignments.

## Acceptance

- [ ] [R21-R23] Review requires Standards, Spec, Coverage, and finding dispositions for the exact candidate; integration invalidates it and repairs re-enter assign/checkpoint/integrate before a fresh review.
- [ ] [R25-R28] Every external mutation is recoverable and idempotent, abandon preserves unlanded work before owned cleanup, and every mutator shares fail-closed preconditions.
- [ ] [R33] Full status resolves retained assignment, checkpoint, cleanup, review, disposition, and digest relationships without leaking receipt bodies.
- [ ] [R57] SIGINT during long lifecycle operations leaves recoverable state and no surviving child process.

