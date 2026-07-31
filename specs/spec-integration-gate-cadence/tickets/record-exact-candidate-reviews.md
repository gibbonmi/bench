# Record exact-candidate reviews

Blocked by: Release integrated assignments

Ownership fence: `internal/specbuild`
Assumptions: checkpoint and integration provenance remain durable and queryable

## What to build

Accept one bounded three-axis review receipt for the exact candidate, retain its
finding dispositions without leaking the receipt body, and expose the complete
assignment, checkpoint, cleanup, review, and digest relationships through the
full durable status projection. Candidate changes invalidate review evidence;
accepted repairs remain ordinary assignments and require a fresh composed review.

## Acceptance

- [ ] [R21-R23] Review requires Standards, Spec, Coverage, and finding dispositions for the exact candidate; integration invalidates it and repairs re-enter assign/checkpoint/integrate before a fresh review.
- [ ] [R33] Full status resolves retained assignment, checkpoint, cleanup, review, disposition, and digest relationships without leaking receipt bodies.
