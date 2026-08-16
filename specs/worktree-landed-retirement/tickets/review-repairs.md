# Repair semantic-review findings

Blocked by: make-release-a-workflow-step.md
Writes: internal/worktree/clean_landed.go, internal/worktree/worktree.go, internal/worktree/clean_landed_apply_test.go, internal/worktree/clean_landed_hostile_test.go, reviews/worktree-landed-retirement.md

## What to build

Close the accepted review findings without changing the selected set or transaction:
stale set fingerprints emit the existing exact stale-fingerprint diagnostic; LR8's
post-apply list assertion crosses a fresh CLI process; and an unsafe retained target
advertises one safe pointer-style per-path remedy without putting the raw byte in output.
Delete the transient review pickup in the same green repair commit.

## Acceptance

- [ ] `(covers LR9)` Every initial-drift partition exits 1, emits the exact stale-fingerprint diagnostic, and mutates and settles nothing.
- [ ] `(covers LR8)` A fresh CLI process lists only the retained assignment after apply.
- [ ] `(covers LR18)` The ESC target remains present and its help contains one safe pointer-style per-path remedy with no raw ESC byte.
