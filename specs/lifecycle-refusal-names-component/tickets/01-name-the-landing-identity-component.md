# Name the landing identity component

Blocked by: none
Writes: internal/worktree/identity_component.go (new), internal/worktree/identity_component_test.go (new), internal/worktree/classifier.go, internal/worktree/worktree.go, internal/worktree/land_identity.go, internal/worktree/land_resume.go, internal/worktree/ownership.go, internal/worktree/land_reauthorization_test.go, internal/worktree/land_resume_refusal_test.go, internal/worktree/land_identity_test.go, internal/worktree/land_surface_test.go, internal/worktree/worktree_test.go, CONTEXT.md

## What to build

An operator who lands with a wrong token, a closed assignment, a wrong path,
or a broken owner bundle reads one refusal. It names the failed component.
For the request miss, the refusal also names the recovery command.

Add the identity component registry and its constructor. Route the landing
preflight, the resume landing, and the release verb's request lookup
through it. Retire
`assignmentMismatchDetail` and the five-predicate marker string. Add the
registry walk that requires a fixture per component. Add the glossary term.
Every existing test that pins the old strings moves to the new sentences.

## Acceptance

- [ ] LR01 pins the request-token sentence on an unknown request.
- [ ] LR02 names `bench worktree reauthorize --assignment <id>` when one active assignment owns the target.
- [ ] LR03 names `bench worktree list` when zero or two active assignments own the target.
- [ ] LR04 names the observed and wanted state for an inactive assignment.
- [ ] LR05 names the assignment's worktree and the target for a path mismatch.
- [ ] LR06 names the owner marker for a rewritten marker.
- [ ] LR07 names the worktree registration for a re-pointed registration.
- [ ] LR08 names the Bench lock for an unlocked worktree.
- [ ] LR09 proves the resume landing prints the same six sentences.
- [ ] LR10 proves two independent refusals still print in one run.
- [ ] LR18 defines `identity component` in `CONTEXT.md` with an Avoid list.
- [ ] LR19 proves `bench worktree release` prints the request-token sentence, its retained clause, and the reauthorize `next=`.
