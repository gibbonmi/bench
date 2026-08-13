# Rebuild bare bench roadmap as the top-10 TOON board

Blocked by: none
Writes: internal/roadmap, cmd/bench/main.go, bin/bench.sh, internal/conformance/package_shipped_surface_test.go

## What to build

Bare `bench roadmap` emits one TOON document — blocks `roadmap` (index rows
for the first 10 parsed rows in document order), `board` (one aggregate row:
`rows_shown`, `rows_total`, `sequence_trusted`), `sequence`, `drain`, `help`
— replacing the verbatim markdown dump, drain-status callout, and
next-action section on this surface only. `RecommendedSequence`,
`RoadmapText`, `bench status`, and `bench idea --owner` keep their current
behavior. Absent `ROADMAP.md` renders the same document with zero
`roadmap`/`sequence` rows and a `/bench-what-next` disclosure, exit 0;
present-but-empty, failed reads, and unsupported schema keep structured
exit-1 refusals. A pending drain discloses `/bench-what-next` in `help`.

## Acceptance

- [ ] Success decodes as one TOON document with the pinned block order,
      first-10 document-order rows, and a true shown-of-total aggregate
      (covers PI10).
- [ ] Boards with 0 (present, sections, no rows), 9, 10, and 11 rows show
      min(10, N) rows and the true total, exit 0 (covers PI11).
- [ ] Absent file → zero-row document, exit 0, `/bench-what-next`
      disclosure; empty/failed/unsupported → structured exit 1
      (covers PI12).
- [ ] A control byte in a row title yields the structured render error,
      never a corrupt document (covers PI14).
- [ ] `bin/bench.sh`'s roadmap help line describes the top-10 board, not
      the removed verbatim-dump-plus-drain-status surface, asserted in the
      existing launcher line-content sweep (covers PI21).
