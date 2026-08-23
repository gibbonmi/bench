# Roadmap flow baseline (2026-08-09 to 2026-08-22)

Source: `git log -- roadmap/ ROADMAP.md` on `main` at ee9c4c10, counting
`roadmap/FT<n>.md` adds (opened), modifies (fed), and deletes (retired), and
`**FT<n>` heading lines in `ROADMAP.md` (open mass) per commit.

- Open mass: 66 rows on 2026-08-09, 72 rows on 2026-08-22. Net +6 in two weeks.
- Before the board split (2026-08-17, 5ffb0470) per-file counts do not exist;
  the mass moved 66 → 67 across 29 commits, with about 12 spec-retires absorbed.
- After the split, 11 drain-shaped commits opened 19, fed 60, and retired 20
  rows; 10 of the 20 retirements sit in one `--restructure` pass (cd355f45).
  Spec-retire commits retired 8 more rows in the same window.
- Without the restructure pass, a drain opens about 1.5 rows and feeds about 5
  rows while it retires under 1. The board grows by size as well as count.

Supports: #1 (baseline), #6 (window and target).
