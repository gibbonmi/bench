# FT311 recoverable reset review

Frozen base: `cf2caa6d5c195c6b6ed5852fb4478902d1d99fbd`

Reviewed tip: `4642603b6212fa8cd9fd8ace93c30b92b3687345`

Current findings: Standards 0, Spec 1, Coverage 1.

De-duplicated decision targets: 1.

## Standards

Count: 0. Worst issue: none. The repair-scoped re-review closed both comment findings at the reviewed tip.

## Spec

Count: 1. Worst issue: the ignored-versus-tracked collision has no authorized behavior.

- `ask-user` — `specs/ft311-recoverable-reset/spec.md:66` promises that ignored files are kept, while line 416 says ignored files are neither preserved nor removed. `internal/worktree/reset_apply.go:108` can overwrite an ignored untracked object when the checkpoint tracks the same path, and `internal/worktree/reset_apply_test.go:143` covers only non-colliding ignored paths. Choose refusal, deletion, or another recoverable collision policy; the benchmark prompt explicitly withholds that decision.

## Coverage

Count: 1. Worst issue: the same ignored-versus-tracked collision partition has no oracle.

- `ask-user` — add the collision partition to RR20 only after its behavior is decided. The existing nested fixture proves survival for ignored paths absent from the checkpoint, not a same-path collision. This is the same decision target as the Spec finding.
