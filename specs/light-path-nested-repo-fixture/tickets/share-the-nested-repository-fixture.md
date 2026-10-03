# Share the nested repository fixture

Blocked by: none
Writes: internal/worktree/journey_test.go, internal/worktree/resume_test.go, internal/worktree/reset_refusal_test.go, internal/worktree/land_fixtures_test.go, internal/worktree/land_resume_refusal_test.go, internal/worktree/lifecycle_acquire_test.go
Covers: none

## What to build

The worktree tests build a nested repository in four places. A nested
repository is a separate repository with one commit inside a checkout. The
reset refusal and the landing reconcile fault share one helper, but its name
says that it is for the reset refusal. The landing resume refusal, the
automatic cleanup planner, and the explicit apply revalidation each build the
same repository inline.

One helper with a neutral name builds the nested repository. It sits with the
other disposable repository builders in `internal/worktree/journey_test.go`. One
constant names the directory and one names the file that the commit tracks. A
test that makes the repository dirty writes that file. Each of the three
inline copies calls the helper. The helper keeps the existing commit
identity.

This change adds no behavior. Each affected test keeps its assertion.

## Acceptance

- [ ] One helper defines the nested repository fixture in `internal/worktree`.
- [ ] No `internal/worktree` test builds a nested repository inline.
- [ ] The affected tests and the `internal/worktree` package suites stay green.
