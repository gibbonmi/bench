# Pin the end-of-input flush and the marker's whitespace tolerance

Blocked by: 02-report-unaccounted-content-below-the-entries-marker.md
Writes: internal/learnings/learnings_test.go

## What to build

Two behaviors the tree already has and no test asserts. This is a review repair
ticket: it adds coverage only, and `internal/learnings/learnings.go` must not
change.

`Parse` flushes an open unaccounted run once more after its main walk ends. No
test reaches that statement, because every existing case ends its input with a
trailing newline, and the empty final element `strings.Split` then produces is
itself in-region and classified blank, which closes the run inside the loop. The
review deleted the whole post-loop `flushRun()` call and all four packages
stayed green. A journal that ends without a trailing newline while a run is open
therefore drops its last record with nothing to catch it.

`unaccountedRegion` recognizes the marker after `strings.TrimSpace`, so a marker
line carrying trailing whitespace still opens the rule. That leniency is
deliberate — an exact match would let one invisible space disable the diagnostic
silently, which is the failure class this spec exists to close — but nothing
asserts it, so a later tightening to exact equality would pass the suite.

Add both cases to the existing unaccounted-region table in
`internal/learnings/learnings_test.go`. Do not add a new test function where the
table already covers the seam.

## Acceptance

- [ ] DL35: a journal ending `<marker>\nplain note` with no trailing newline yields one unaccounted record carrying `plain note` and its 1-based line.
- [ ] DL36: a marker line written with trailing spaces still opens the rule, and the line below it is reported.
- [ ] Deleting the post-loop `flushRun()` call in `Parse` turns DL35 red.
- [ ] Changing the marker match from `strings.TrimSpace(line) == JournalEntriesMarker` to exact equality turns DL36 red.
- [ ] `internal/learnings/learnings.go` is unchanged by this ticket.
