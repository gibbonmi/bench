# Learnings — usage journal

<!-- entries below -->

## 2026-08-23 - bench idea on main blocks the next landing  [open]
- **What happened:** The coordinator ran `bench idea` on the primary checkout, as the working agreement directs. `bench commit` then refused the primary checkout, and `bench worktree land` refused the worktree landing because the destination held that dirty `capture/IDEAS.md`. The file was byte-identical to the worktree commit. The `block-dangerous-git` hook forbids the coordinator from the discard, so the landing waited on the reviewer.
- **Right behavior:** A parked idea must reach a commit without a hand discard on `main`. Two seams can close this gap. `bench idea` can write into the active worktree when one exists. Or the landing's clean-destination preflight can accept a `capture/` path whose dirty bytes equal the source tip.
- **Proposed rule change:** Reviewer decision on which seam owns the fix. The `union` compose rule for `capture/IDEAS.md` covers merge conflicts, not a dirty destination, so it does not close this gap.
