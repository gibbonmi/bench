# Carry a schema descriptor through the parser

Blocked by: none
Writes: internal/coverage/coverage.go, internal/coverage/coverage_test.go

## What to build

Prefactoring only — make the change easy before making it. `internal/coverage`
today derives `width`, `fields`, and `storyOffset` from one opt-in boolean, and
several checks index cells by literal offset. Replace that with one schema
descriptor value carrying the header's field-name list, and have every
offset-taking check read its cell through the descriptor. No header is added and
no message changes: every existing violation string, `Rows` output, and
`ParseSpec` verdict is byte-identical before and after.

## Acceptance

- [ ] `bench coverage --check` reports the same violations, in the same order and
      wording, for a six-column and a legacy five-column spec as it does today.
- [ ] `bench coverage <spec>` renders the identical `rows[N]{story,seam,red_signal}`
      table it renders today.
- [ ] `ParseSpec` returns the same opt-in verdict, ordered row IDs, and violations
      for both existing headers.
- [ ] Adding a header means adding one descriptor, with no other call site edited —
      demonstrated by the next ticket landing without touching a check body.
