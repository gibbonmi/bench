# Resume idempotent mutations

Blocked by: Enforce shared mutator preconditions

Ownership fence: `internal/specbuild`
Assumptions: every external lifecycle mutation has one durable request identity and prepared state

## What to build

Surround lifecycle side effects with prepared and completed state so exact
re-entry finishes or reports one durable next action without duplicating runs,
worktrees, commits, or refs. Reused identities with different inputs conflict.
Long-running Git and owner calls obey cancellation as one process group and leave
recoverable state with no surviving child.

## Acceptance

- [ ] [R25, R27] Re-entering interrupted start, assign, checkpoint, integrate, and review transitions is request-idempotent; exact repeats are no-ops and reused identities with different inputs conflict without duplicate side effects.
- [ ] [R57] SIGINT during checkpoint and integration leaves recoverable state and no surviving child, through a shared runner seam later promote and abandon can consume.
