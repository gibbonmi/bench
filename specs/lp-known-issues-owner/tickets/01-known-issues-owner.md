# Name FT482 as the owner of the sealed cache generation issue

Blocked by: none
Writes: docs/known-issues.md, specs/lp-known-issues-owner/tickets/01-known-issues-owner.md
Covers: none

## What to build

Each row of `docs/known-issues.md` names a retire condition.
That condition is an obligation in the body of the owning roadmap row, or "no owner".
The row for the sealed cache generation issue says "no owner", but roadmap row `FT482` owns that seal policy.

Rewrite only the retire condition cell of that row.
The new cell names `FT482` and states its obligation in the form of the `FT341` row.
`FT482` is still `Next: decide`, so the cell states the obligation and does not state a decided outcome.
No other text in the file changes.

## Acceptance

- [ ] The retire condition cell of the sealed cache generation row in `docs/known-issues.md` starts with `FT482:` and states an obligation that agrees with the body of `roadmap/FT482.md`.
- [ ] No other line of `docs/known-issues.md` changes.
- [ ] The `bench commit` lane passes.
