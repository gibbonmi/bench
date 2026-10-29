# The background grade reuses the wrapper path owner

Blocked by: none
Writes: internal/gate/background_grade.go
Covers: none

## What to build

`backgroundWrapper` in `internal/gate/background_launch.go` owns the path of the worktree wrapper `bin/bench.sh` in the `gate` package.
`BackgroundCommand` in `internal/gate/background_grade.go` joins the same path itself when it replaces its image with that wrapper.
The sealed learning "The wrapper path bin/bench.sh is joined in two places" (2026-10-27) records this duplicate.

Make `BackgroundCommand` call `backgroundWrapper(root)`, and remove the import that the old join used.
The behavior does not change.
The duplicate in `internal/worktree/build.go` is out of scope, because it needs a reviewer decision on a leaf owner.

## Acceptance

- [ ] `internal/gate` joins `bin/bench.sh` in production code only in `backgroundWrapper`.
- [ ] `TestBackgroundGradeRunsAsTheTargetExecutable` (BG27) passes, and a probe that swaps the wrapper path of the entry turns it red.
- [ ] The `bench commit` lane passes.
