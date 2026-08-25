# Detect incremented shared variables

Blocked by: none
Writes: internal/worktree/parallel_census_test.go

## What to build

The parallel census treats an increment or decrement as a write to the same shared variable as an assignment.

## Acceptance

- [ ] The census refuses a test that increments or decrements a package-level variable.
- [ ] The census marks a test serial when it increments or decrements an imported package variable.
- [ ] Existing assignment and local-shadow behavior stays unchanged.
