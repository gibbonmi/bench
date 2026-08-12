# Repair final evidence closure

Blocked by: repair-worktree-compatibility-closure.md, repair-control-bearing-action-values.md, repair-guidance-precision.md, repair-learnings-refusal-evidence.md
Writes: `internal/coverage/coverage_test.go`, `reviews/axi-query-disclosure.md`, `capture/session-handoff.md`

## What to build

Remove the coverage test helper's mirrored production derivation, close every re-review finding only after all four independently green repair tickets land, delete the pickup, and refresh the full-run phase pin.

## Acceptance

- [ ] [FE1] (covers QD6) coverage compatibility expectations are independently literal rather than rebuilt through the same action derivation used in production; mutating a help action makes the exact fixture test red.
- [ ] [FE2] (covers local) every finding recorded in `reviews/axi-query-disclosure.md` is resolved by its owning repair ticket before the pickup is deleted.
- [ ] [FE3] (covers local) `capture/session-handoff.md` pins the resulting clean review boundary and the exact harness-native continuation command.
