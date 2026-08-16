# Refuse a half-applied landed set

Blocked by: apply-the-landed-plan-and-settle-records.md
Writes: internal/worktree/clean_landed_apply_test.go, internal/worktree/clean_landed.go

## What to build

Through the existing lifecycle `Fault` seam, prove the two fixed mid-apply failure
shapes: a fault after the first row's transaction completes, and a fault callback that
dirties the second row's tree at the first row's terminal step and returns nil. Both
leave the first row removed and settled, later rows untouched with nothing preserved,
exit non-zero, and refuse the original fingerprint on re-run; after the first shape a
fresh bare plan lists only the remaining row. No set-level receipt is introduced.
Demo: fault-injected apply, then re-plan.

## Acceptance

- [ ] `(covers LR14)` Fault after the first row's completion: exit non-zero, first tree removed and its record settled, second tree present, original `--apply` refused as drift, fresh bare plan lists only the second row.
- [ ] `(covers LR20)` Second row dirtied at the first row's terminal step: first row settled, second refused as drifted with nothing removed or preserved, exit non-zero, original fingerprint refused on re-run.
