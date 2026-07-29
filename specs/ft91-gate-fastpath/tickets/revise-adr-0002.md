# Revise ADR 0002 posture 5 for gate-execution reuse

Blocked by: Short-circuit gate execution on a reusable green

## What to build

Story 14 of `specs/ft91-gate-fastpath/spec.md`: ADR 0002 posture 5 revised to
name gate execution (not only the gated commit) as the consumer that reuses a
fresh green for the identical closed subject, same reopen trigger, and to
record the new accepted residual — an in-place `shellcheck` upgrade inside the
freshness window under a reused green. No other posture moves. Current-state
prose per the ADR standard (`craft-adr`): no file paths, no code snippets, no
history-of-the-change.

## Acceptance

- [ ] Posture 5 names gate execution as a reusing consumer with the same
      reopen trigger.
- [ ] The shellcheck in-place-upgrade residual is recorded as accepted.
- [ ] No other posture changed; the revision reads as current decided state.
