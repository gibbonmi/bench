# State the lane commit in the field guide

Blocked by: none
Writes: docs/field-guide.html
Covers: none

## What to build

The field guide's green gate pane says that a commit happens only in the green gate state.
A worktree commit happens on a lane pass, and the landing runs the whole-project gate before `main` takes the work.
The pane states that current rule, and it keeps its shift-loop and review claims.

## Acceptance

- [ ] The green gate pane in `docs/field-guide.html` no longer says that a commit happens only in the green gate state.
- [ ] The pane says that `bench commit` publishes on a lane pass and that `bench worktree land` runs the whole-project gate before `main` takes the work.
