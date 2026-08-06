# Allow refresh of a clean preserved assignment

Blocked by: repair-ft78-symlink-prospective-fixture.md
Ownership fence: `internal/specbuild/integrate.go`, `internal/specbuild/refresh_test.go`
Integration surfaces: zero-byte preservation patch→`internal/specbuild/integrate.go`; clean refresh lifecycle proof→`internal/specbuild/refresh_test.go` plus CR1; non-empty replay conflict and byte identity→`internal/specbuild/refresh_test.go` plus CR2
Contracts: preservation patch bytes cross refresh planning→`internal/specbuild/integrate.go` checkpoint replay, asserted by CR1-CR2 against the real `assign --refresh` lifecycle in `internal/specbuild/refresh_test.go`

## What to build

Teach the shared checkpoint replay helper to skip `git apply` only when the
preservation patch has zero bytes. Add focused refresh lifecycle coverage for a
clean, uncheckpointed assignment whose prerequisite repair advanced the candidate.
The refresh must move the assignment to the repaired candidate without inventing
payload, changing tracked bytes, or weakening non-empty replay validation.

## Acceptance

- [ ] [CR1] A clean uncheckpointed assignment refreshes onto the repaired candidate and its worktree remains byte-clean at that candidate.
- [ ] [CR2] Non-empty checkpoint and refresh replay still apply the preserved patch and retain conflict and byte-identity refusal behavior.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CR1 | always invoke `git apply` for a zero-byte preservation patch | `TestRefreshAdvancesACleanAssignmentOntoTheRepairedCandidate` | run `go test ./internal/specbuild -run '^TestRefreshAdvancesACleanAssignmentOntoTheRepairedCandidate$' -count=1`; expect `No valid patches in input` |
| CR2 | skip `git apply` for a non-empty preservation patch | `TestRefreshRepairTraceThroughPublicLifecycle` and `TestRefreshRefusesConflictWithoutTouchingTheWorktree` | run `go test ./internal/specbuild -run '^(TestRefreshRepairTraceThroughPublicLifecycle|TestRefreshRefusesConflictWithoutTouchingTheWorktree)$' -count=1`; expect payload loss or the conflict refusal to fail |
