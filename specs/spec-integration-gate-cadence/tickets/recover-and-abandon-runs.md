# Recover and abandon runs

Blocked by: Record exact-candidate reviews

Ownership fence: `internal/specbuild`
Assumptions: review and assignment provenance remain durable and queryable

## What to build

Put every mutator behind one fail-closed precondition owner, surround external
mutations with request-idempotent recovery state, and abandon only through a
fingerprinted plan/apply transition that preserves unlanded work before owned
cleanup. Interrupted long-running lifecycle work must leave recoverable state
and no surviving child process.

## Acceptance

- [ ] [R25-R28] Every external mutation is recoverable and idempotent, abandon preserves unlanded work before owned cleanup, and every mutator shares fail-closed preconditions.
- [ ] [R57] SIGINT during long lifecycle operations leaves recoverable state and no surviving child process.
