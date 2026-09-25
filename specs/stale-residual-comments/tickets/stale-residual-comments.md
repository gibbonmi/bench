# Correct two comments that FT336 made stale

Blocked by: none
Writes: internal/worktree/build.go, internal/responsebound/responseboundtest/spill.go
Covers: none

## What to build

The comment on `buildExecNext` names `mergeReconcileNext` as its precedent for the
fallback to the assignment id. That function now names a reset command, so the
precedent is wrong. The live precedent is `conflictRepairPrefix`, which uses the
assignment id when the path is not line-safe.

The package comment of `responseboundtest` says that the package only reads the spill
path. Since FT336, the package also holds the shared `AssignmentCheckout` fixture. The
package comment states both jobs.

These edits change only comments. No behavior changes.

## Acceptance

- [ ] The `buildExecNext` comment names `conflictRepairPrefix` as its precedent, and no comment in the file names `mergeReconcileNext`.
- [ ] The `responseboundtest` package comment names the spill reader and the assignment checkout fixture.
- [ ] `bench test --package ./internal/worktree` and `bench test --package ./internal/responsebound/...` pass, and the gate is green at the landing.
