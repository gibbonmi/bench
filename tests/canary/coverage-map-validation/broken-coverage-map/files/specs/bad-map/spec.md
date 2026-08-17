# bad map (canary fixture)

This fixture stays on the legacy header — the no-row-ID form carrying a
`red signal` column — on purpose. It is what proves the legacy branch of
coverage-map validation still runs, so a sweep migrating fixtures to the
reduced header must leave this one alone.

## User stories
1. As a reviewer, I want malformed coverage rows to fail the gate.

### Acceptance coverage map
| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 2 | a behavior |  | a red signal | why it catches |
