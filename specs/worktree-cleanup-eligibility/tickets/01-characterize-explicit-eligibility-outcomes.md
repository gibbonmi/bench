# Characterize explicit eligibility outcomes

Blocked by: none
Writes: internal/worktree/eligibility_test.go, specs/worktree-cleanup-eligibility/tickets/01-characterize-explicit-eligibility-outcomes.md

## What to build

Add the independent explicit-cleanup outcome matrix at the existing planning seam.
Each real-fixture case observes the unchanged `CleanupPlan` projection and asserts
one current `(Action, ReasonCode)` tuple, its relevant evidence, and that planning
does not mutate durable state. The table is an oracle for the current precedence;
it must not derive expected outcomes from production decision code.

## Acceptance

- [x] EX1: primary and unsafe explicit targets retain as `uncertain` with their current detail and no planning mutation.
- [x] EX2: an unregistered explicit target retains as `foreign`.
- [x] EX3: malformed owner or declaration evidence retains as `malformed` with the current winner.
- [x] EX4: foreign or mismatched ownership locks retain as `unexpected-lock`.
- [x] EX5: live or ambiguous lease evidence retains as `live-lease` when it is the current final refusal.
- [x] EX6: undeclared, excessive, or otherwise unauthorized ignored residue retains as `ignored`.
- [x] EX7: a clean registered removable checkout projects `remove` with no reason code.
- [x] EX8: otherwise-removable dirty or detached state projects `recover-remove` with current recovery evidence.
- [x] EX9: authorized bounded ignored residue projects `discard-remove` with no reason code.
