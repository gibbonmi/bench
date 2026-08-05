# Allow refresh of a clean preserved assignment

Blocked by: repair-ft78-symlink-prospective-fixture.md
Ownership fence: `internal/specbuild/integrate.go`, `internal/specbuild/refresh_test.go`
Integration surfaces: refresh preservation patch -> candidate replay; clean owned assignment -> refreshed lifecycle checkout
Contracts: an empty preservation patch in `internal/specbuild/integrate.go` is the identity transformation over the repaired candidate, while non-empty checkpoint replay retains its existing conflict and byte-identity checks; `internal/specbuild/refresh_test.go` proves the clean assignment advances and remains clean

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

- [ ] [MCR1] Always invoking `git apply` makes the clean refresh fail with `No valid patches in input`.
- [ ] [MCR2] Skipping apply for non-empty patches loses the existing consumer payload or makes its lifecycle assertions fail.
