# Research links remove the retired FT360 row target

Blocked by: none
Writes: docs/research/bench-efficiency-assessment-2026-10-06.md
Covers: none

## What to build

The `shared-test-fixtures` build (`e746bfe51`) and its retirement (`6a62ff258`) completed `FT360` with no ADR and no live absorbing roadmap row.
The drain retires `roadmap/FT360.md`.
Two links in `docs/research/bench-efficiency-assessment-2026-10-06.md` still point to `../../roadmap/FT360.md`.

Under the 2026-10-27 research-link rule (`556a6af3e`), a retired row with no ADR and no live absorbing row keeps its link text as plain text.
Replace `[FT360](../../roadmap/FT360.md:3)` on line 281 and `[FT360](../../roadmap/FT360.md:1)` on line 389 with `FT360`.

## Acceptance

- [ ] `docs/research/bench-efficiency-assessment-2026-10-06.md` contains no link to `roadmap/FT360.md`.
- [ ] Lines 281 and 389 keep the plain text `FT360` and change nothing else.
- [ ] The `bench commit` lane passes.
