## Standards

Finding count: 0. Worst issue: none.

## Spec

Finding count: 1. Worst issue: CD9 does not preserve every existing error response in its compatibility oracle.

- `auto-fix` — CD9 requires existing error kinds, exits, and argv grammar to remain byte-equal (`specs/axi-coherent-diff/spec.md`, CD9), but `internal/diff/compatibility_test.go` covers successful modes and five argv cases without the live base-resolution or existing unrepresentable-TOON error responses. Add exact preserved assertions for those error modes.

## Coverage

Finding count: 3. Worst issue: CD5 has no symlink-target drift mutation.

- `auto-fix` — CD5 requires per-path symlink-target identity, implemented in `internal/diff/diff.go`, but `internal/diff/identity_test.go` does not mutate a symlink target between captures. Add a retry-then-refuse case.
- `auto-fix` — CD5 requires stat identity for other kinds, implemented in `internal/diff/diff.go`, but the FIFO coverage in `internal/diff/matrix_test.go` checks inventory and body exclusion only. Add a mid-read non-regular stat mutation.
- `auto-fix` — the hostile-input checklist requires control-bearing Git-sourced commit subjects to refuse, but `internal/diff/matrix_test.go` covers only a control-bearing path. Add a reachable commit-subject case through `bench diff --full` and require the structured TOON refusal.
