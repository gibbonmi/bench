# worktree-landed-retirement review pickup

Frozen pair: base `68ebb9cef9a39a3d35349b9dc4534dad2c044f33`, source
`d388d4134615d8d10eef2416ccee0c4c64000c4b`.

## Standards

Finding count: 0. Worst issue: none. The derived-count finding is closed.

## Spec

Finding count: 0. Worst issue: none.

## Coverage

Finding count: 1. Worst issue: high.

- **High — auto-fix.** The hostile selector test now proves special assignment paths
  never enter the explicit planner, but the new `ListCommand` and orphan-sweep consumers
  call `PlanExplicit` directly (`internal/worktree/list.go:76`,
  `internal/worktree/resume.go:405`). A checkout directory with FIFO/socket/symlinked
  `.git` metadata can therefore reach `git -C <target>` before the selector's shape
  guard, contrary to `projects/benchkit.md:128` and LR19
  (`specs/worktree-landed-retirement/spec.md:302`). Route both consumers through the
  guarded landed-planning path and add process-safe special-metadata fixtures.
