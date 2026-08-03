# FT181 #2/#3 — liveness and ownership classification: tree evidence (2026-08-03, delegated read)

Supports the resolved answers to map tickets #2 (husk/symlink soften to
liveness, worktree planner included) and #3 (narrow the prepared-abandon
exemption; "registration no longer found" becomes liveness).

## Classification today

- `liveCheckout` (`internal/specbuild/precondition.go:229-243`): only
  `os.ErrNotExist` softens to `errAbsentCheckout`. A husk (directory present,
  `.git` gone), a dangling symlink, a foreign checkout, and any other `Lstat`
  failure all land on `errOwnership`.
- `ownedAssignments` (`precondition.go:199-224`): identity faults are
  `errOwnership` unconditionally; only `errAbsentCheckout` is swallowed, and
  only for abandon (`precondition.go:216-219`).
- `intent.FindAssignmentByRequest` (`internal/intent/assignment.go:323-334`)
  returns `(Assignment{}, false, nil)` for a missing ledger entry; the
  `!found` disjunct at `precondition.go:211` classifies it as `errOwnership`,
  indistinguishable from a forged record — the exact state ticket #3
  reclassifies as liveness.

## Why the #2 amendment (worktree scope) is load-bearing

- For a truly absent path, the worktree layer independently proves safety:
  `planRemovedCheckout` (`internal/worktree/resume.go:113-155`) cross-checks
  the live `git worktree list` registration against the intent ledger and
  re-asserts any existing recovery ref (resume.go:142-146); pinned by
  `recovery_retry_test.go:80-131`.
- A husk never reaches that branch — it falls to `PlanExplicit`
  (`resume.go:93`), which needs git metadata the husk no longer has, so the
  specbuild-layer softening alone would route abandon into a planner refusal.
  The map's amendment authorizing the worktree abandon-planning change is what
  makes the decision deliverable.

## The #3 exemption as shipped

- `precondition.go:94-97` re-exempts ANY ownership error — identity forgery
  included — once `op == abandon` and a prepared abandon/apply operation with
  a non-empty `Result` exists: a second, broader escape hatch layered on the
  narrow inner one. Requiring liveness classification there closes the path a
  tampered record could ride through on a resumed abandon.

## Coverage gaps for the spec

Untested at this layer: husk, dangling symlink at an assignment path,
`FindAssignmentByRequest` not-found while the run record still lists the
assignment, and whether the prepared-abandon exemption swallows an identity
fault (the only softening test, `TestAbandonAppliesForRemovedWorktree`,
reaches `errAbsentCheckout`). Pinned already: removed-worktree abandon and
recovery-ref survival, forged-identity refusal, foreign-checkout refusal
(`abandon_test.go:33-140`).

Not read: `internal/worktree/path.go` beyond ~line 155;
`internal/intent` ledger locking; checkpoint test files.
