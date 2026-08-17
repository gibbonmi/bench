# Split the board: land the tree parser and this board's migration in one green

Blocked by: none
Writes: internal/roadmap, internal/conformance/docs_workflow_checks_test.go, cmd/bench/command_registry_test.go, ROADMAP.md, roadmap/

## What to build

The roadmap package parses an index (`ROADMAP.md`: one physical heading line per
row, no bodies) plus a `roadmap/` listing (`<ID>.md`, first line equal to the index
line, body after) into the same `Document` it produces today, and returns an ordered
list of integrity diagnostics — missing detail owner, orphan file, inline body,
heading mismatch, unrecognized file, duplicate ID, wrapped heading — each beginning
with the offending repo-relative path, with the per-class row disposition the spec
fixes. One loader classifies both sources and feeds the pure parse; `roadmap/` joins
the sequence-trust source list; the inline-body shape is no longer accepted.
`bench idea --owner` validates through the loader. `checkOccurrenceLedgerMigration`
reads the loader and `checkRemovedVerbSweep` sweeps the row files; `setupAXIRoadmap`
writes a split fixture. In the same commit, a scratchpad script (never committed)
migrates this board: headings joined onto one line, bodies (inline text, ledgers,
`Sources:`) into `roadmap/FT<n>.md` under a verbatim heading, non-row prose
untouched. Before migrating, capture `bench roadmap --context --full` with the
binary built at the base commit; after, with the current binary; `roadmap_rows` and
`sequence` blocks must diff empty, and the commit message names both captures and
the empty diff. Coverage rows PR1, PR2, PR3, PR4, PR5, PR6, PR7, PR8, PR9, PR10,
PR16, PR19, PR21, PR22, PR28.

## Acceptance

- [ ] An index line with a matching row file (heading verbatim, empty body) parses to a row with the index title, empty body, and no diagnostic.
- [ ] A row-file ledger projects the same `occurrence_count` and keys the inline parser projected; a malformed ledger reports against `roadmap/<ID>.md`.
- [ ] Each of the seven fault classes yields exactly its diagnostic and its spec-fixed row disposition, with fixtures for: missing owner (row kept, `rows_total` counts it), orphan with and without an index, inline body, heading mismatch, `notes.md` and `FT7.txt` unrecognized while `AB1.md` is a row, duplicate ID, wrapped heading (no row).
- [ ] Neither file present yields no diagnostic and no row.
- [ ] `bench idea --owner FT7 text` exits non-zero with `ROADMAP.md is structurally untrusted` when `roadmap/FT7.md` is missing, and appends when clean.
- [ ] `docs-currency-workflow` reports the same occurrence counts on the migrated board and its retired-verb sweep visits every row file.
- [ ] `roadmap_rows` and `sequence` blocks are byte-identical before and after migration; every index heading is one line; the loader reports zero diagnostics on the migrated tree; non-row prose is unchanged; the commit message carries the two capture paths and the empty-diff statement.
