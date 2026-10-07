# Cleanup member traversal (FT363; C05)

Status: ready

## Destination

Give cleanup modes one ordered member traversal around their existing transactions.
Preserve each mode's authority and output contract.

## Notes

The reviewer delegated architecture recommendations on 2026-10-06.
A member outcome describes one selected target after an apply attempt.
Avoid: deletion claim for an unstarted member.
Use craft-spec and craft-seams for the preserving spec.

## Decisions so far

- [Should all modes return identical failure rows?](architecture-cleanup/tickets/1.md): preserve the existing unclaimed difference.
- [What does the shared traversal own?](architecture-cleanup/tickets/2.md): share order and stop mechanics around mode authority.
- [What proves preservation across the modes?](architecture-cleanup/tickets/3.md): retain command outcomes and transaction failure tests.

## Not yet specified

## Spec-writer discretion

- Internal row adapters and traversal names.
- Mechanical migration order with independent green checkpoints.

## Out of scope

- A uniform output change for the unclaimed mode.
- A replacement for applyCleanupTransaction or exact branch deletion.
- A weaker ownership, preservation, lock, or fingerprint requirement.

## Sources

- Path: `roadmap/FT363.md`
  Supports: the existing traversal candidate.
  Drift: a change to the candidate scope.
- Path: `internal/worktree/clean_set_apply.go`
  Supports: the existing outcome helpers.
  Drift: a change to stopped-member projection.
- Path: `internal/worktree/clean_unclaimed.go`
  Supports: the unclaimed prefix-only failure result.
  Drift: a change to applyUnclaimedAssignmentSet.
- Path: `internal/worktree/clean_landed.go`
  Supports: the landed mode's unreached outcomes.
  Drift: a change to applyLandedSet.
- Path: `docs/adr/0005-worktree-cleanup-requires-verifiable-ownership.md`
  Supports: the closed ownership and preservation decision.
  Drift: a new approved decision.

