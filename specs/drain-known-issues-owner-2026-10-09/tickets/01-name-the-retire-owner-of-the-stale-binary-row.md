# Name the retire owner of the stale-binary known issue

Blocked by: none
Writes: docs/known-issues.md
Covers: none

## What to build

The known-issue table in `docs/known-issues.md` holds a row for the stale worktree
`dist/bench` binary after a ticket commit. Its retire cell says `no owner`. The 2026-10-09
drain d-677ac383a1a1 opens roadmap row FT406, which decides where the rebuild lives.

Set that row's retire cell to the FT406 obligation, in the `FT<n>: <obligation text>` form
that the table requires. Change no other row.

## Acceptance

- [ ] The stale-binary row's retire cell reads `FT406:` followed by the FT406 heading text.
- [ ] Every other row of `docs/known-issues.md` is unchanged.
