# Reuse the protected policy fixture in the planning commit test

Blocked by: none
Writes: specs/ft383-planning-commit-fixture/tickets/01-reuse-the-protected-policy-fixture-in-the-planning-commit-test.md (new), internal/commit/commitment_test.go
Covers: none

## What to build

The `planningCommitRepo` helper in the commit package tests writes a protected policy seed by hand. The seed has two roadmap rows, two outcomes, one policy, and one board. The `commitmenttest.SeedProtected` fixture already owns the same seed, byte for byte. Two copies of one fixture can drift, so the helper must call the fixture instead.

The helper keeps its other steps. It initializes the landing repository, commits the seed with the same message, and returns the planning checkout. This change does not change the behavior of a test. The fixture file stays unchanged.

## Acceptance

- [ ] `planningCommitRepo` calls `commitmenttest.SeedProtected` and holds no row body, outcome, policy, or board text of its own.
- [ ] Each test in `internal/commit` that calls `planningCommitRepo` passes, and the set of passing test names is the same as before the change.
- [ ] A recorded probe that changes the row body in the fixture shows which `internal/commit` tests detect the change.
