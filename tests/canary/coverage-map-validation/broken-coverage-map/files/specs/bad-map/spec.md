# bad map (canary fixture)

This fixture uses the reduced header — the accepted form with no `red signal`
column — and proves the empty-cell diagnostic reaches the gate:
`coverage map row 1 has an empty 'seam' cell` fails a real gate run under this
header. (That the same diagnostic reads identically across accepted headers is
proven in `coverage_test.go`, not here.)

## User stories
1. As a reviewer, I want malformed coverage rows to fail the gate.

### Acceptance coverage map
| story | behavior | seam | why it catches the failure |
|---|---|---|---|
| 2 | a behavior |  | why it catches |
