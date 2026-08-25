# Repair the review findings

Blocked by: 02-surface-the-resolver-reason-in-exec-and-path.md
Writes: internal/worktree/identity_component.go, internal/worktree/ownership.go, internal/worktree/path.go, internal/worktree/identity_component_test.go, internal/worktree/worktree_test.go, CONTEXT.md, reviews/lifecycle-refusal-names-component.md (delete)

## What to build

The five accepted repair targets in `reviews/lifecycle-refusal-names-component.md`,
then delete that file in the same commit. No spec row changes.

## Acceptance

- [ ] A bundle with a wrong owner ID and a wrong branch names the owner marker, and each identity predicate has one source.
- [ ] `resolveWorktree` calls `landingActiveState` instead of a second derivation.
- [ ] `CONTEXT.md` keeps every descriptive sentence at 25 words or fewer.
- [ ] A landing with a re-pointed registration and an unlocked worktree names the registration.
- [ ] `bench worktree release` with a rewritten owner marker prints `owner marker does not match assignment <id>; checkout retained`.
