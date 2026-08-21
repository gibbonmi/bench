# Report a dated learnings line that misses heading shape

Blocked by: none
Writes: internal/learnings/learnings.go, internal/learnings/learnings_test.go, internal/learnings/testdata, internal/roadmap/context_test.go, internal/status/status_test.go

## What to build

A writer who appends a dated entry to `capture/learnings.md` as a bullet
instead of a `## <date> — <title> [open]` heading gets a diagnostic instead of
silence. `learnings.Parse` gains one rule: a line that leads with a
`YYYY-MM-DD` date, after a bounded walk past a run of runes that each either
satisfy `unicode.IsSpace` or are one of the markdown marker bytes `-`, `*`,
`+`, `>`, and `#` — a run of length zero included, so a date flush at column
one is reported too — and that does not already start `## `, becomes a `Malformed`
record with the reason `dated learning entry is not a heading`, its 1-based
source line, and the line text with any trailing carriage return removed.
Records stay in ascending source-line order.

`unicode.IsSpace` is the predicate exactly: it carries the whole `White_Space`
property, so every `Zs` separator is stripped and the zero-width characters
U+200B and U+FEFF are not. Do not borrow `learnings.isSpace`, which is ASCII-only
because it serves a TOON whitespace class.

The date test that `isDatedHeading` performs today becomes one named helper
both rules call — the grammar of a date keeps one definition, digit-shape only.
Do not write the new rule with a calendar parse.

This is the complete path from journal bytes to every reader. No consumer
changes: `bench learnings`, `bench roadmap --context`, and
`roadmap.learningCount` (which `bench status` reads) already consume the
malformed list, so the ticket proves each of them stops reporting a false zero
rather than editing them.

Undated content stays quiet in this ticket — the marker-anchored rule that
judges it is the sibling ticket's, and a freshly adopted repo must read green
after this one lands on its own.

## Acceptance

- [ ] DL1: `Parse` returns one malformed record for a `- <date> — <title>` bullet under the schema heading.
- [ ] DL2: that record carries the bullet's 1-based source line.
- [ ] DL3: that record's reason is `dated learning entry is not a heading`.
- [ ] DL4: `bench learnings` exits 1 on the two-bullet journal that produced the 2026-08-21 drop.
- [ ] DL5: `bench learnings` renders one `line <n>` row per lost dated line, byte-exact against a checked-in stdout fixture.
- [ ] DL6: a dated line preceded by `-`, `*`, `+`, `>`, or `#` is reported.
- [ ] DL32: a dated line flush at column one, with no leading marker or whitespace, is reported.
- [ ] DL7: a dated line whose marker is separated from the date by U+00A0 is reported.
- [ ] DL8: a dated line separated by U+3000 is reported, and one separated by U+200B is not.
- [ ] DL9: a dated bullet appended inside an existing open entry's body is reported.
- [ ] DL10: `## broken` still yields exactly one record reading `malformed learning heading`.
- [ ] DL11: `## 2026-01-01 — x` without `[open]` still yields exactly one record reading `dated learning heading must end with [open]`.
- [ ] DL12: `- 2026-88-88 — x` is reported as a lost dated line.
- [ ] DL13: `bench roadmap --context` renders one `parse_failures` row sourced at `capture/learnings.md` for the lost line.
- [ ] DL14: `bench status` renders the drain row's learnings component as `unknown (capture/learnings.md is malformed)`.
- [ ] DL15: a freshly scaffolded journal renders the empty table at exit 0.
- [ ] DL16: a journal holding only its schema heading renders the empty table at exit 0.
- [ ] DL17: a well-formed open entry keeps its `date,title` row at exit 0.
- [ ] DL20: a dated bullet on a final line with no trailing newline is reported.
- [ ] DL21: a CRLF-terminated dated bullet's record carries the line without its trailing carriage return.
