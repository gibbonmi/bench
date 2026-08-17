# Migrate automatic shared landedness

Blocked by: 03-expand-typed-eligibility-verdict.md
Writes: internal/worktree/classifier.go, internal/worktree/landed.go, internal/worktree/clean_landed.go, internal/worktree/resume.go, internal/worktree/list.go, internal/worktree/eligibility_test.go, internal/worktree/orphan_test.go, internal/worktree/clean_branch_test.go, internal/worktree/landed_test.go, internal/worktree/list_actions_test.go, specs/worktree-cleanup-eligibility/tickets/04-migrate-automatic-shared-landedness.md

## What to build

Move automatic cleanup and the shared `assignmentLanded` reader onto the verdict's
typed landedness evidence. Update every direct reader in classifier, landed-set
selection, resume, and list so no production code parses `plan.landed` or another
formatted decision string. Preserve current automatic tuples, the list and resume
landed classifications, and the derived-after rule: explicit `DiscardBranch` can
authorize exact branch deletion but never changes the evidence automatic cleanup
uses.

## Acceptance

- [ ] CO1: `PlanAutomatic` and the shared assignment-landed readers obtain their answer from typed verdict evidence, with no production parsing of formatted landedness.
- [ ] DB1: explicit `DiscardBranch` can authorize exact deletion after derived landedness, while automatic cleanup retains the unprovable-branch fixture as `unmerged`.
