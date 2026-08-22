# Close the semantic review findings

Blocked by: fail-built-in-go-table-closed.md, diagnose-partial-session-environment.md
Writes: internal/gate/branch_native_phases_test.go, internal/systemtest/session_start_test.go, projects/benchkit.md, specs/inherited-toolchain-environment/spec.md, reviews/inherited-toolchain-environment.md

## What to build

The gate and SessionStart fixtures exercise each explicit edge that initial review found open. The spec remains the one source for host-specific phase-close evidence.

## Acceptance

- [ ] The TE1 fixture proves that an empty regular `go.mod` requires Go.
- [ ] The partial-environment fixture proves discovery when `ENVMAN_LOAD` is absent or empty.
- [ ] The discovery fixture accepts an executable symlink without executing its target.
- [ ] A nonzero discovery process cannot produce a recovery assignment from valid-looking output.
- [ ] The project profile keeps only the current PATH recovery rule, and the spec records the host-specific phase-close evidence.
- [ ] The repair commit removes `reviews/inherited-toolchain-environment.md`.
