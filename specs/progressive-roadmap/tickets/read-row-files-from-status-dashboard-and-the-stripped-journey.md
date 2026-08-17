# Read row files from status, the dashboard, and the stripped journey

Blocked by: split-the-board-parser-and-migration-in-one-green.md
Writes: internal/status, internal/systemtest/owner_test.go

## What to build

`bench status`'s roadmap-reconcile signal classifies spec paths the loader found in
row files; the dashboard is pinned by a regression test only (both readers already parse
the index); the stripped-distribution system journey removes and then asserts the
absence of `roadmap/` beside `ROADMAP.md`. Coverage rows PR17, PR18, PR20.

## Acceptance

- [ ] A merged spec named only in `roadmap/FT7.md` counts in the reconcile row.
- [ ] The dashboard renders roadmap text and sequence for a split tree.
- [ ] The stripped subject carries no `roadmap/` directory and the journey's excluded-path list names it.
