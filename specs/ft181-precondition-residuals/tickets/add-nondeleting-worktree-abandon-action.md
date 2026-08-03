# Add the non-deleting abandon plan action for present-but-not-a-checkout paths

Blocked by: none
Ownership fence: `internal/worktree`
Assumptions: `planAbandonWithPlanner` (`internal/worktree/path.go:76-101`) today routes only absent paths to `planRemovedCheckout` and everything present to `PlanExplicit`; the `CleanupAction` enum and `removes()` predicate live in `internal/worktree/classifier.go:67-90`; plan facts join the fingerprint via `fingerprintParts` in `planRemovedCheckout`; claims re-derived from the tree at pickup

## What to build

An abandon over a path that is present but not a checkout — a husk directory
without git metadata, a dangling symlink, a non-directory entry — plans a
**new non-deleting action**: it performs the same registration cross-check the
removed-checkout path performs (live worktree registration reconciled against
the intent ledger, any existing recovery ref re-asserted), and apply releases
the registration and intent entry while leaving the filesystem entry
untouched. The new action is excluded from `removes()`. The plan carries the
leftover path so disposal routes through the existing size-bounded clean
surface, and the action and path join the plan fingerprint like every existing
plan fact. A cross-check mismatch — registration absent, branch ref, path, or
request disagreeing with the ledger — refuses at the planner exactly as the
removed-checkout path refuses today (`errAbandonMismatch`). Mirror the plan
and apply prior art in `internal/worktree/recovery_retry_test.go`'s
removed-directory and recovery-survival cases. The probe never opens the path:
classify by `Lstat`/`Stat` shape, so a FIFO cannot block planning.

## Acceptance

- [ ] [WA1] `PlanAbandon` over a husk (directory present, git metadata gone) plans the new non-deleting action naming the leftover path; `ApplyAbandon` releases registration and intent entry and the husk's content is byte-identical after apply.
- [ ] [WA2] The same holds for a dangling symlink and a non-directory entry (regular file, FIFO) at the assignment path, and planning returns without opening the path.
- [ ] [WA3] The planner refuses a husk whose registration cross-check mismatches: registration absent, branch ref, path, or request disagreeing with the ledger.
- [ ] [WA4] Existing recovery refs survive plan and apply over a not-live-but-present path and still resolve to their original heads.
- [ ] [WA5] The new action and leftover path are folded into the plan fingerprint; a plan for a different path or action yields a different fingerprint.
- [ ] [WA6] A foreign checkout (present, resolvable, different git common dir) still refuses at the planner with the checkout retained.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| WA1 | route the husk through the existing `ActionRemove` path | the bytes-untouched husk test | apply, run `go test ./internal/worktree -run Abandon`, expect the husk-content-changed failure |
| WA2 | classify by directory-existence alone | the dangling-symlink case of the shape table | apply, run the table test, expect the symlink-case failure |
| WA3 | skip the registration cross-check on the new path | the mismatch-refusal test | apply, run it, expect the missing-refusal failure |
| WA4 | drop the recovery-ref re-assertion from the new action | the recovery-survival test | apply, run it, expect the lost-ref failure |
| WA5 | omit the leftover path from `fingerprintParts` | the fingerprint-distinctness test | apply, run it, expect the equal-fingerprint failure |
| WA6 | treat any present-but-probe-failing path as not-live | the foreign-checkout refusal test | apply, run it, expect the missing-refusal failure |
