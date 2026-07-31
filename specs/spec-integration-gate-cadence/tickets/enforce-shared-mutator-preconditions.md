# Enforce shared mutator preconditions

Blocked by: Record exact-candidate reviews

Ownership fence: `internal/specbuild`
Assumptions: the working branch, candidate ref, spec tip, and assignment ownership are durable identities

## What to build

Put every existing lifecycle mutation behind one fail-closed precondition owner.
Distinguish expected dirt in the named assignment worktree from dirt or an
unrecognized move in the active working checkout, stale candidate/spec identity,
missing prerequisite evidence, and assignment ownership mismatch. Expose the
same owner for the later promote and abandon transitions.

## Acceptance

- [ ] [R28] Start, assign, checkpoint, integrate, and review share the same precondition owner and refuse before mutation across dirty or moved working checkout, stale candidate/spec tip, missing evidence, and ownership mismatch; expected assignment dirt remains admissible.
