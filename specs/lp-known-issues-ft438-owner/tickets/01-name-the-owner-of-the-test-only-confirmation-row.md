# Name the owner of the test-only confirmation known issue

Blocked by: none
Writes: docs/known-issues.md
Covers: none

## What to build

The known-issue table in `docs/known-issues.md` holds a row for `bench record confirmation`,
which refuses a test-only repair in a delegated build. Its retire cell says `no owner`. The
2026-10-21 drain d-036de9c826ed opens roadmap row FT438, which owns the test-only exemption.

Set that row's retire cell to the FT438 obligation, in the `FT<n>: <obligation text>` form
that the table requires. Change no other row.

## Acceptance

- [ ] The test-only confirmation row's retire cell reads `FT438:` followed by the FT438 heading text.
- [ ] Every other row of `docs/known-issues.md` is unchanged.
