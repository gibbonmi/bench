# worktree-landed-retirement review pickup

Frozen pair: base `68ebb9cef9a39a3d35349b9dc4534dad2c044f33`, source
`7c24f3ebedb76085cf29d0d5535746b026a5589c`.

## Standards

Finding count: 1. Worst issue: medium.

- **Medium — auto-fix.** `internal/canary/inventory_test.go:21` hard-codes the
  repository-derived fixture total, duplicating the production derivation at
  `internal/canary/inventory.go:127` without a mapped independent-expectation need.
  Preserve a useful output-shape assertion without a global literal count.

## Spec

Finding count: 0. Worst issue: none. The prior pointer-remedy finding is closed.

## Coverage

Finding count: 1. Worst issue: medium.

- **Medium — auto-fix.** LR19 promises no Git invocation against a special assignment
  path (`specs/worktree-landed-retirement/spec.md:130` and `:302`), but
  `internal/worktree/clean_landed_hostile_test.go:116` asserts only the resulting row
  and path presence. Add a seam-backed assertion that every hostile shape makes zero
  target-path planner/Git calls.
