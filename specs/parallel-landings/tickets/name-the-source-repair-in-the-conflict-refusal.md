# Name the source repair in the conflict refusal

Blocked by: union-merge-the-phase-owned-journals.md
Writes: internal/worktree/land.go, internal/worktree/land_test.go, internal/worktree/land_surface_test.go
Line: opus / medium — one renderer at the known refusal seam.

## What to build

A conflict the rule table does not settle refuses as today. The refusal keeps
its `refused{detail=composition conflict: <kind>}` record and its paths table,
and it gains `next=`. The value names the source repair in order. Merge the
destination commit into the source worktree, commit the repair, review the new
range, and re-run the landing with the new tip. A value that is not line-safe
takes the pointer form every `next=` already uses. The merge step names raw
Git because no Bench verb moves a retained worktree onto the destination yet;
the text says so.

## Acceptance

- [ ] A conflict on `ROADMAP.md` refuses and names `ROADMAP.md` in the paths table (covers WL16).
- [ ] The conflict refusal's `next=` names the source repair steps and the re-run invocation (covers WL18).
- [ ] A conflicted path that carries a control byte renders through the sanitized paths table (edge under WL16).
- [ ] A `next=` that is not line-safe takes the pointer form (edge under WL18).
