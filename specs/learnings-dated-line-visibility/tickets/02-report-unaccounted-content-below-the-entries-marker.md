# Report unaccounted content below the entries marker

Blocked by: 01-report-dated-lines-that-miss-heading-shape.md
Writes: internal/learnings/learnings.go, internal/learnings/learnings_test.go, internal/learnings/testdata, internal/adopt/init.go, internal/adopt/adopt_test.go, internal/conformance/docs_workflow_checks_test.go

## What to build

A writer who appends an undated note below the scaffold's
`<!-- entries below -->` marker gets a diagnostic instead of silence.

`internal/learnings` exports the marker as `JournalEntriesMarker`.
`internal/adopt`'s learnings scaffold and `internal/conformance`'s
docs-reference check both read it from there instead of holding their own copy
of the literal, so the boundary the parser enforces is the boundary the
scaffold ships. That de-duplication is the ticket's prefactoring; do it first
and keep the scaffold's rendered bytes unchanged.

`learnings.Parse` then gains its second rule. The marker opens the rule only
when it appears above the first *real* entry heading — a line starting `## `
that `isTemplatePlaceholder` does not claim. That exclusion is load-bearing, not
tidy: the scaffold prints its worked example `## <date> - <short title>  [open]`
above the marker, so an anchor that counted that line would never open the rule
on the one journal shape it exists to serve. A marker below a real entry heading
does not open the rule, a journal with a marker and no real entry heading opens
it through end of file, a journal with no marker keeps today's behavior for
undated content entirely, and a second marker below the first is an ordinary
line that joins the run.

Below an opening marker and above the first real entry heading, each line is
classified in source order: a blank line is quiet and ends any open run, a dated
line takes the sibling ticket's reason and ends any open run, and any other line
opens or continues a run. Blank is `strings.TrimSpace(line) == ""` exactly, so a
whitespace-only line is blank — a record whose text is invisible spaces gives
the writer nothing to repair. Each run yields one `Malformed` record with the
reason `learning content below the entries marker is not an entry`, the run's
first line's 1-based number, and that line's text with any trailing carriage
return removed.

Once a `## ` line appears, the existing walk owns every line to the next `## `,
so the entry body bullets the scaffold asks for are untouched. Records stay in
ascending source-line order across all four reasons.

This ticket shares one contract with its blocker: the dated rule's reason and
its record shape are already in the tree when this lands, and the run
classification defers to it for any dated line rather than re-deciding what a
date is.

## Acceptance

- [ ] DL29: the bytes `internal/adopt` scaffolds parse with zero records, and the scaffold's marker is `learnings.JournalEntriesMarker` rather than its own literal.
- [ ] The docs-reference check in `internal/conformance` splits on `learnings.JournalEntriesMarker`, and the Go sources hold the literal exactly once.
- [ ] DL22: an undated line below the entries marker yields one record reading `learning content below the entries marker is not an entry`.
- [ ] DL23: three contiguous undated lines below the marker yield exactly one record, carrying the first line's text and its 1-based number.
- [ ] DL24: a blank line below the marker produces no record.
- [ ] DL33: a whitespace-only line below the marker produces no record.
- [ ] DL25: an undated bullet inside a well-formed open entry's body produces no record.
- [ ] DL26: a journal with no entries marker produces no record for its undated lines.
- [ ] DL27: a marker appearing below a real entry heading does not open the rule.
- [ ] DL34: a second marker below the first joins the open run's record rather than starting a new region.
- [ ] DL28: a CRLF-terminated unaccounted run's record carries the line without its trailing carriage return.
- [ ] DL18: an undated, non-heading line above the entries marker produces no record.
- [ ] DL19: records render in ascending source-line order across all four reasons.
- [ ] DL31: a dated bullet inside a fenced code block below the marker is still reported.
- [ ] DL30: `bench learnings` exits 1 and renders a `line <n>` row for a scaffolded journal with one undated note appended below its marker, byte-exact against a checked-in stdout fixture.
