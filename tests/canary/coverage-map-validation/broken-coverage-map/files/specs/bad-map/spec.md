# bad map (canary fixture)

This fixture uses the reduced header — the accepted form with no `red signal`
column — and proves the schema-independent empty-cell diagnostic still fires:
`coverage map row 1 has an empty 'seam' cell` reads the same regardless of
which accepted header a map uses.

## User stories
1. As a reviewer, I want malformed coverage rows to fail the gate.

### Acceptance coverage map
| story | behavior | seam | why it catches the failure |
|---|---|---|---|
| 2 | a behavior |  | why it catches |
