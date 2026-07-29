# Guard fixture marker reads against special files

Blocked by: none

## What to build

One regular-file-guarded marker reader (stat before open, reject non-regular
files with a named refusal) serving `CHECK` today and `TEST` next ticket, with
the same trim discipline `CHECK` uses now. Replaces `fixtureCheck`'s unguarded
open in `internal/canary` — story 6 of `specs/ft91-gate-fastpath/spec.md`.

## Acceptance

- [ ] A FIFO planted at a `CHECK` path is rejected before opening with a
      refusal naming the path; the unit test bounds itself with its own
      deadline so an unguarded open reds rather than wedges.
- [ ] Absent file vs present-but-empty file remain distinct behaviors, both
      asserted at the reader.
- [ ] Trailing-newline trim behavior is asserted at the reader.
- [ ] `CHECK` reads route through the one shared reader — no second open-coded
      read remains.
