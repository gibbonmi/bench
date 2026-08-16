# Accept the reduced coverage header

Blocked by: carry-a-schema-descriptor-through-the-parser.md
Writes: internal/coverage/coverage.go, internal/coverage/coverage_test.go

## What to build

A spec author can write an acceptance coverage map with no `red signal` column
and have it validate. Both reduced headers parse as mapped —
`row|story|behavior|seam|why it catches the failure` opted into row IDs, and
`story|behavior|seam|why it catches the failure` without them — while both
existing headers keep working unchanged. The parser chooses a schema by the
header's cell *names*, never its cell count, so the two distinct five-cell
headers cannot be confused. Every check that names or indexes a column follows
the reduced schema: wrong-width rows report `(want 5)`, empty-cell violations
name a reduced field, the one-predicate check reads the reduced behavior offset,
and story references, the fan-out bound, row-ID validation, and orphan stories
behave identically under all four headers. A half-renamed header is refused as
missing the canonical header rather than parsed under a guess.

## Acceptance

- [ ] `(covers SR1, SR2)` Both reduced headers report state `mapped`, with row IDs present for the
      five-cell form and absent for the four-cell form.
- [ ] `(covers SR4)` Two five-cell headers differing only in whether they name `red signal`
      select different schemas from identical row bytes.
- [ ] `(covers SR5)` A reduced-map row with the wrong cell count reports `(want 5)`.
- [ ] `(covers SR6)` An empty cell in a reduced map names a field the reduced header declares.
- [ ] `(covers SR7)` A `;` in a reduced map's behavior cell is refused as more than one predicate.
- [ ] `(covers SR3, SR8)` Story-reference, fan-out-bound, duplicate-ID, malformed-ID, and orphan-story
      violations are string-identical across all four headers.
- [ ] `(covers SR13)` A header renaming only some cells reports `coverage map missing the canonical header`.
- [ ] `(covers SR14)` A reduced map with a header and no rows reports `coverage map has no data rows`.
- [ ] `(covers SR15, SR16)` The historical marker skips validation under a reduced header.
- [ ] `(covers SR8)` An escaped pipe (`\|`) inside a reduced map's behavior cell survives parsing.
