# Use the standard slice and string helpers in tests

Blocked by: none
Writes: internal/worktree/delegated_integration_test.go, internal/writeguard/writeguard_test.go, internal/conformance/handoff_single_source_test.go, internal/git/staged_test.go, internal/shift/shift_test.go, internal/shift/fault_test.go, internal/lines/lines_fixtures_test.go, internal/lines/lines_parse_test.go, internal/lines/lines_agentline_test.go, internal/lines/lines_verdict_test.go
Covers: none

## What to build

The quality survey of 2026-09-29, in its list of small certain cuts, found six test helpers that copy a standard-library function. Each helper has the same semantics as its standard function:

- `contains` in `internal/worktree`, `internal/writeguard`, and `internal/conformance` is `slices.Contains`.
- `equalStrings` in `internal/git` is `slices.Equal`. It compares lengths and then each element, so a nil slice and an empty slice are equal, as in `slices.Equal`.
- `contains` in `internal/shift` is `strings.Contains`. It is a manual substring scan, and an empty substring matches.
- `contains` in `internal/lines` calls `strings.Contains` and does nothing more.

Replace each helper with its standard function at every call site in its package. Then delete the helper, and remove each import that becomes unused. Change no production code. Add, remove, or rename no top-level `Test` function, because `internal/worktree` pins its exact test count.

## Acceptance

- [ ] No package in the list above declares `contains` or `equalStrings`, and each former call site calls the standard function.
- [ ] `go vet ./...` passes.
- [ ] `bench test --changed` passes.
- [ ] The `internal/worktree` test-count pin passes with no change to its count.
