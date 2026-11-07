# Route the resume command's request placeholder through its owner

Blocked by: none
Writes: internal/worktree/land_refusal.go, internal/worktree/land_resume_next_test.go
Covers: none

## What to build

The resume command that an incomplete landing prints names the request placeholder. Both forms of that command read the placeholder from `requestSlot`, which owns its spelling. The two forms are the form that quotes its values and the form that a value with a control byte drops to. The rendered text does not change.

## Acceptance

- [ ] The resume command in both forms names the request placeholder that `requestSlot` holds.
- [ ] The existing tests that render the resume command pass with no change.
- [ ] A swap of the placeholder in the form for a value with a control byte makes the focused test fail.
