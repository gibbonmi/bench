# FT311 recoverable reset review

Frozen base: `cf2caa6d5c195c6b6ed5852fb4478902d1d99fbd`

Reviewed tip: `2e50ae20d8162392d251494e5651918f507cb52b`

Raw findings: Standards 2, Spec 2, Coverage 1.

De-duplicated repair or decision targets: 3.

## Standards

Count: 2. Worst issue: the exported reset-ref constructor has no public documentation.

- `auto-fix` — `internal/intent/ledger/validate.go:36` exports `ResetRefPrefix` without the full-sentence symbol-leading Go doc required by the comment standard. Add the contract comment without changing behavior.
- `auto-fix` — `cmd/bench/worktree_leaves.go:72` restates the private `worktreeSuffix` name and one-line body. Delete the redundant comment.

## Spec

Count: 2 raw, 1 actionable. Worst issue: the ignored-versus-tracked collision has no authorized behavior.

- `ask-user` — `specs/ft311-recoverable-reset/spec.md:66` promises that ignored files are kept, while line 416 says ignored files are neither preserved nor removed. `internal/worktree/reset_apply.go:108` can overwrite an ignored untracked object when the checkpoint tracks the same path, and `internal/worktree/reset_apply_test.go:143` covers only non-colliding ignored paths. Choose refusal, deletion, or another recoverable collision policy; the benchmark prompt explicitly withholds that decision.
- `no-op` — the proposed exit-3 guidance omission is refuted by `.bench/BENCH-reference.md:264`, which says that a fault after preservation exits 3 and names the restore command. No repair target remains.

## Coverage

Count: 1. Worst issue: the same ignored-versus-tracked collision partition has no oracle.

- `ask-user` — add the collision partition to RR20 only after its behavior is decided. The existing nested fixture proves survival for ignored paths absent from the checkpoint, not a same-path collision. This is the same decision target as the Spec finding.
