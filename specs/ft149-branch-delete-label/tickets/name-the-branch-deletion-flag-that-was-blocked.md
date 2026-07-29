# Name the branch deletion flag that was blocked

Blocked by: none

## What to build

The destructive-git guard keeps blocking reviewer-owned branch deletions while
its refusal names the exact safe or force deletion form the caller used.

## Acceptance

- [x] `git branch -d` and `git branch --delete` refusals name
  `git branch -d`, while `git branch -D` refusals name `git branch -D`.
- [x] The delete labels remain owned by the ordered deny table and are selected
  from the flag the classifier already matched; the call site formats no label.
- [x] Delegate-worktree deletion carve-outs and all existing branch block/allow
  verdicts remain unchanged through the real hook-envelope seam.
