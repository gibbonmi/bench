## Standards

Finding count: 0. Worst issue: none.

## Spec

Finding count: 1. Worst issue: the repair graph orders WF27 backwards.

- `auto-fix` — `repair-profile-loop-routing.md` now cites WF27, but the ticket that authors WF27 is not its blocker. Add `Blocked by: repair-census-map-coverage.md` so a fresh graph walk cannot schedule the profile repair before the coverage row it requires.

## Coverage

Finding count: 0. Worst issue: none. All 204 fixtures and 27 rows were re-derived and the prior gaps remain closed.

Raw findings: Standards 0, Spec 1, Coverage 0. De-duplicated repair targets: 1.
