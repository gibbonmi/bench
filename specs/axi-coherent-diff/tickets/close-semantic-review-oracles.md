# Close the semantic review oracles

Blocked by: render-one-coherent-diff-snapshot.md
Writes: `internal/diff/identity_test.go`, `internal/diff/matrix_test.go`, `internal/diff/compatibility_test.go`, `reviews/axi-coherent-diff.md`

Closure: RV1/symlink-target-drift, RV1/nonregular-stat-drift, RV2/hostile-commit-subject, RV3/base-resolution-errors, RV3/toon-refusal-errors, RV3/error-exits

## What to build

Close the exact semantic-review gaps without changing production behavior. Extend the identity matrix with symlink-target and non-regular stat movement, extend the hostile Git-sourced-text matrix with a control-bearing commit subject through `bench diff --full`, and make the compatibility oracle preserve the existing live base-resolution and TOON-refusal error bytes and exits. Delete `reviews/axi-coherent-diff.md` in the same green repair commit so resolved pickup state cannot resurface.

These findings remain one repair because all three targets are test-oracle additions under `internal/diff`; splitting them would create three test-only landings over the same already-shipped production tracer without an independently demonstrable user behavior between them.

## Acceptance

- [ ] [RV1] (covers CD5) changing an untracked symlink target or an untracked non-regular entry's stat identity between captures retries once and then refuses with the exact drift dimension and invocation.
- [ ] [RV2] (covers CD6) a reachable commit subject carrying a refused control byte makes `bench diff --full` return the structured unrepresentable-TOON error and exit 1.
- [ ] [RV3] (covers CD9) the compatibility oracle pins the existing live base-resolution and TOON-refusal error bytes, kinds, and exits in addition to its successful and argv cases.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RV1/symlink-target-drift | omit symlink targets from per-path content identity | identity drift matrix | change the target between both capture pairs and require two retries followed by the named refusal |
| RV1/nonregular-stat-drift | hardcode other-kind stat identity | identity drift matrix | move a FIFO's stat identity between both capture pairs and require two retries followed by the named refusal |
| RV2/hostile-commit-subject | bypass TOON validation for log subjects | hostile commit-subject process fixture | create a reachable commit with an ESC in its subject, run `bench diff --full`, and require exit 1 with the existing structured error |
| RV3/base-resolution-errors | change or omit one existing live base-resolution error | compatibility oracle | run the captured unresolvable-default and no-merge-base fixtures and require exact bytes and exit 1 |
| RV3/toon-refusal-errors | change or omit an existing unrepresentable path or subject error | compatibility oracle | run the captured hostile-cell fixtures and require exact bytes and exit 1 |
