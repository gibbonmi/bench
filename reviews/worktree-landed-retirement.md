# worktree-landed-retirement review pickup

Frozen pair: base `68ebb9cef9a39a3d35349b9dc4534dad2c044f33`, source
`8b1281c3a2f6e350f9d7d843ac7bf8e767b51c49`.

## Standards

Finding count: 0. Worst issue: none.

## Spec

Finding count: 1. Worst issue: P1.

- **P1 — auto-fix.** LR18 requires a retained control-byte row to expose a safe
  pointer-style per-path remedy (`specs/worktree-landed-retirement/spec.md:151`),
  but `internal/worktree/clean_landed.go:224` drops the help action whenever the
  raw target is unsafe. Add a safe pointer-addressed remedy without emitting the byte.

## Coverage

Finding count: 2. Worst issue: high.

- **High — auto-fix.** LR9 promises the stale-fingerprint error
  (`specs/worktree-landed-retirement/spec.md:285`), while
  `internal/worktree/clean_landed_apply_test.go:86` accepts a silent exit 1 and
  `internal/worktree/worktree.go:279` emits only the current plan. Require and emit
  the exact diagnostic for every initial-drift partition.
- **Medium — auto-fix.** LR8's post-apply ledger read
  (`specs/worktree-landed-retirement/spec.md:231`) calls `ListCommand` in the same
  process at `internal/worktree/clean_landed_apply_test.go:43`, missing the fresh-process
  serialized-state edge in `projects/benchkit.md:165`. Re-read `worktree list` through
  the built CLI after apply.
