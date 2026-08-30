# Repair: the printer comment and the WF3 row derivation

Blocked by: repair-system-check-name-reservation.md
Writes: internal/worktree/path.go, internal/worktree/build_test.go

## What to build

Re-review findings STD-1 and COV-1 in `reviews/worktree-native-forms.md`.
The `printTargetRefusal` comment names every verb that calls it: path, exec,
show, build, and create. `TestBuildPrintsTheTable` derives its expected row
through `toon.Table` over the same cells the verb renders, so the encoder's
quoting rule has one source.

## Acceptance

- [ ] STD-1: the comment's verb list equals the set `rg -n "printTargetRefusal\(" internal/worktree` returns.
- [ ] WF3: `TestBuildPrintsTheTable` passes for an id that starts with `0` and a digit, proved by a fixed id or a repeated run.
- [ ] The gate `test` phase stays green for the whole `internal/worktree` package.
