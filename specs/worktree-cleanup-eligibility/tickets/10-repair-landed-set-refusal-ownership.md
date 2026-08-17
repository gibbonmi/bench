# Repair: move the landed-set preservation refusal's selection into eligibility.go

Blocked by: 09-repair-discard-branch-detached-head.md
Writes: internal/worktree/eligibility.go, internal/worktree/classifier.go, internal/worktree/clean_landed.go, specs/worktree-cleanup-eligibility/tickets/10-repair-landed-set-refusal-ownership.md

Source: reviews/worktree-cleanup-eligibility.md, Standards finding 3 and Spec finding 1 (CO3/EV2).

## What happened

Ticket 08 moved only the boolean predicate (`automaticPreservationTrigger`) into
eligibility.go. The actual Action/ReasonCode selection for the landed-set refusal
— `plan.Action, plan.ReasonCode, plan.Reason = ActionRetain, ReasonDirty, "per-path
cleanup is required to preserve work"` — still lives in clean_landed.go's
`retainForLandedPreservation`, merely relocated within the same file from ticket
05. CO3/EV2 require the selection itself, not just the trigger, to live in the
eligibility module.

## What to build

In eligibility.go, replace `automaticPreservationTrigger(plan CleanupPlan) bool`
with a function that returns the decided verdict shape (retain-or-not, Action,
ReasonCode) rather than a bare boolean — parameterized only by the one thing that
legitimately differs per caller, the operator-facing message text (both
`decideAutomatic`'s dirty-refusal branch and `retainForLandedPreservation` share
the same ReasonDirty/ActionRetain pairing, only the Reason string differs).
Neither `decideAutomatic` nor `retainForLandedPreservation` should contain a bare
`ActionRetain, ReasonDirty` literal after this change — both call the shared
function and project only the message. Preserve every existing message string and
Recovery-reset behavior exactly.

## Acceptance

- [x] CO3/EV2: the only literal `ActionRetain, ReasonDirty` construction for this refusal is inside eligibility.go; both consumers call it.
- [x] Byte-identical messages: "automatic cleanup does not preserve uncommitted work" (PlanAutomatic) and "per-path cleanup is required to preserve work" (landed-set) unchanged.
- [x] All existing tests in eligibility_test.go, clean_landed_apply_test.go remain green with unmodified assertions.
