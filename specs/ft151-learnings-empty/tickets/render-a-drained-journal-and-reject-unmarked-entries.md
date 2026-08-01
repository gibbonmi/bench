# Render a drained journal and reject unmarked entries

Blocked by: none

## What to build

`bench learnings` recognizes a valid journal with no entries as an empty
rowset, while surfacing every dated heading that is not an open entry.

## Acceptance

- [x] The current post-drain `capture/learnings.md` shape and the scaffold
  template both render `learnings[0]{date,title}:` at exit 0.
- [x] Arbitrary prose, a zero-byte record, unreadable input, and wrong-type
  input remain fail-closed rather than becoming false empty states.
- [x] A dated heading missing `[open]` or carrying another state is rendered as
  a line-attributed malformed row at exit 1; valid open rows still render.
- [x] `bench status` and `bench learnings` continue to derive open-entry counts
  from the same parser.
