# Pin the serial ceiling and close

Line: opus / medium.

Blocked by: 14-convert-direct-home-binds-in-land-and-clean.md, 15-convert-direct-home-binds-in-lifecycle-and-core.md
Writes: internal/worktree/parallel_census_test.go, CHANGELOG.md

## What to build

The census exposes the serial set's size. A live-tree pin requires that size
at or below the count this ticket observes. The pin's comment names the tests
that stay serial and why. A fixture that falls back to a process bind
raises the count above the pin.

Record the home seam in `CHANGELOG.md` under Unreleased, beside the entry the
census ticket added.

## Acceptance

- [ ] WF18 pins the serial ceiling, and a synthetic set one above the ceiling turns it red.
- [ ] `CHANGELOG.md` names the explicit home under Unreleased.
