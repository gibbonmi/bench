# Repoint the dead FT355 link in the efficiency assessment

Blocked by: none
Writes: docs/research/bench-efficiency-assessment-2026-10-06.md
Covers: none

## What to build

The 2026-10-18 drain d-0573adab5125 folded roadmap row FT355 into FT372 and deleted
`roadmap/FT355.md`. The parked-rows table in
`docs/research/bench-efficiency-assessment-2026-10-06.md` still links `roadmap/FT355.md`.
No other live file links that deleted row file.

Repoint that link to `roadmap/FT372.md`. Keep the row's FT355 name, and state that FT355
folded into FT372. Change no other row.

## Acceptance

- [ ] The FT355 row of the parked-rows table links `../../roadmap/FT372.md` and states that FT355 folded into FT372.
- [ ] No live file outside `capture/` links `roadmap/FT355.md`.
- [ ] Every other row of the assessment is unchanged.
